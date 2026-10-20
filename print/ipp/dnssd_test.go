package ipp

import (
	"context"
	"net"
	"testing"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

// response builds an mDNS response of the records of one printer.
func response(t *testing.T, instance, service, host string, port uint16, txt []string, ip net.IP) []byte {
	t.Helper()
	bld := dnsmessage.NewBuilder(nil, dnsmessage.Header{Response: true, Authoritative: true})
	bld.EnableCompression()
	bld.StartAnswers()
	inst := dnsmessage.MustNewName(instance + "." + service)
	hn := dnsmessage.MustNewName(host)
	rh := func(n dnsmessage.Name, typ dnsmessage.Type) dnsmessage.ResourceHeader {
		return dnsmessage.ResourceHeader{Name: n, Type: typ, Class: dnsmessage.ClassINET, TTL: 120}
	}
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(bld.PTRResource(rh(dnsmessage.MustNewName(service), dnsmessage.TypePTR), dnsmessage.PTRResource{PTR: inst}))
	bld.StartAdditionals()
	must(bld.SRVResource(rh(inst, dnsmessage.TypeSRV), dnsmessage.SRVResource{Target: hn, Port: port}))
	must(bld.TXTResource(rh(inst, dnsmessage.TypeTXT), dnsmessage.TXTResource{TXT: txt}))
	if ip4 := ip.To4(); ip4 != nil {
		var a [4]byte
		copy(a[:], ip4)
		must(bld.AResource(rh(hn, dnsmessage.TypeA), dnsmessage.AResource{A: a}))
	}
	b, err := bld.Finish()
	must(err)
	return b
}

func TestBrowseParse(t *testing.T) {
	b := &browser{found: map[string]*found{}}
	b.parse(response(t, `Canon\032MF\032Series`, ippService, "canon.local.", 631,
		[]string{"rp=ipp/print", "pdl=application/pdf,image/pwg-raster,image/urf", "ty=Canon MF Series", "UUID=1234", "TLS=1.2"}, net.IPv4(192, 0, 2, 10)))
	b.parse(response(t, `Canon\032MF\032Series`, ippsService, "canon.local.", 443,
		[]string{"rp=ipp/print", "pdl=application/pdf,image/pwg-raster,image/urf"}, net.IPv4(192, 0, 2, 10)))
	b.parse(response(t, "Office", ippService, "office.local.", 8631, []string{"rp=printers/office"}, net.IPv4(192, 0, 2, 20)))
	b.parse(response(t, "Not a printer", "_http._tcp.local.", "web.local.", 80, nil, net.IPv4(192, 0, 2, 30)))
	b.parse([]byte{1, 2, 3}) // garbage
	got := b.services()
	if len(got) != 2 {
		t.Fatalf("services: %+v", got)
	}
	c := got[0]
	if c.Name != "Canon MF Series" || !c.Secure || c.Port != 443 || c.URL != "ipps://192.0.2.10:443/ipp/print" || c.Host != "canon.local" {
		t.Fatalf("Canon: %+v", c)
	}
	if c.TXT["ty"] != "Canon MF Series" || c.TXT["UUID"] != "1234" {
		t.Fatalf("the keys of the other record: %v", c.TXT)
	}
	if f := c.Formats(); len(f) != 3 || f[1] != "image/pwg-raster" {
		t.Fatalf("formats: %v", f)
	}
	if o := got[1]; o.Name != "Office" || o.Secure || o.URL != "ipp://192.0.2.20:8631/printers/office" {
		t.Fatalf("Office: %+v", o)
	}
}

func TestBrowseFollowUp(t *testing.T) {
	b := &browser{found: map[string]*found{}}
	// A PTR without the rest: the follow-up asks for SRV and TXT.
	bld := dnsmessage.NewBuilder(nil, dnsmessage.Header{Response: true})
	bld.StartAnswers()
	bld.PTRResource(dnsmessage.ResourceHeader{Name: dnsmessage.MustNewName(ippService), Type: dnsmessage.TypePTR, Class: dnsmessage.ClassINET},
		dnsmessage.PTRResource{PTR: dnsmessage.MustNewName("Lonely." + ippService)})
	p, _ := bld.Finish()
	b.parse(p)
	q := b.followUp()
	if q == nil {
		t.Fatal("no follow-up")
	}
	var par dnsmessage.Parser
	if _, err := par.Start(q); err != nil {
		t.Fatal(err)
	}
	qs, err := par.AllQuestions()
	if err != nil || len(qs) != 2 || qs[0].Type != dnsmessage.TypeSRV || qs[1].Type != dnsmessage.TypeTXT || qs[0].Class&qu == 0 {
		t.Fatalf("questions: %+v %v", qs, err)
	}
	if len(b.services()) != 0 {
		t.Fatal("a service without its SRV record")
	}
}

func TestInstanceName(t *testing.T) {
	for in, want := range map[string]string{
		`Canon\032MF\032Series._ipp._tcp.local.`:    "Canon MF Series",
		`a\.b._ipp._tcp.local.`:                     "a.b",
		`plain._ipp._tcp.local.`:                    "plain",
		`\230\151\165\230\156\172._ipp._tcp.local.`: "日本",
	} {
		if got := instanceName(in); got != want {
			t.Errorf("instanceName(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestBrowseLoopback finds a printer of USB that the responder of this
// machine advertises on the loopback.
func TestBrowseLoopback(t *testing.T) {
	c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Skip(err)
	}
	defer c.Close()
	resp := response(t, `HP\032LaserJet\032(USB)`, ippService, "localhost.", 60000, []string{"rp=ipp/print", "ty=HP LaserJet", "pdl=application/pdf"}, net.IPv4(127, 0, 0, 1))
	go func() {
		buf := make([]byte, 9000)
		for {
			n, from, err := c.ReadFrom(buf)
			if err != nil {
				return
			}
			var par dnsmessage.Parser
			if _, err := par.Start(buf[:n]); err == nil {
				c.WriteTo(resp, from)
			}
		}
	}()
	old := loopbackMDNS
	loopbackMDNS = c.LocalAddr().(*net.UDPAddr)
	defer func() { loopbackMDNS = old }()
	svcs, err := Browse(context.Background(), 300*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range svcs {
		if s.Name == "HP LaserJet (USB)" {
			if !s.Local || s.URL != "ipp://127.0.0.1:60000/ipp/print" || s.Host != "localhost" {
				t.Errorf("%+v", s)
			}
			return
		}
	}
	t.Fatalf("the printer of USB is not found: %+v", svcs)
}
