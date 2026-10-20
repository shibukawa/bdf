package ipp

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

// freeUDPPort is a port of UDP no one listens on.
func freeUDPPort(t *testing.T) int {
	t.Helper()
	c, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	return c.LocalAddr().(*net.UDPAddr).Port
}

// TestAdvertise advertises a printer on the loopback and finds it with
// Browse, as bdf printers finds the printers of this machine: its SRV,
// TXT and address, under the name and type advertised.
func TestAdvertise(t *testing.T) {
	port := freeUDPPort(t)
	old := loopbackMDNS
	loopbackMDNS = &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port}
	defer func() { loopbackMDNS = old }()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ads := []Advertisement{
		{Instance: "Office @ server", Type: "_ipp._tcp", Subtypes: []string{"_print", "_universal"}, Port: 8631,
			TXT: []string{"rp=ipp/print/office", "pdl=application/pdf,image/urf", "URF=W8,SRGB24", "ty=BDF"}},
		{Instance: "Office @ server", Type: "_ipps._tcp", Port: 8631, TXT: []string{"rp=ipp/print/office", "TLS=1.2"}},
		{Instance: "Labels v2.1", Type: "_ipp._tcp", Port: 8631, TXT: []string{"rp=ipp/print/labels"}},
	}
	r, err := Advertise(ctx, ads, ResponderOptions{Host: "server.local", Port: port, Interfaces: []net.Interface{}, Loopback: true})
	if err != nil {
		t.Fatal(err)
	}
	names := r.Names()
	if len(names) != 3 || names[0] != "Office @ server._ipp._tcp.local." || names[2] != "Labels v2 1._ipp._tcp.local." {
		t.Fatalf("names %q", names)
	}
	found, err := Browse(ctx, 700*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	var office, labels *Service
	for i := range found {
		switch found[i].Name {
		case "Office @ server":
			office = &found[i]
		case "Labels v2 1":
			labels = &found[i]
		}
	}
	if office == nil || labels == nil {
		t.Fatalf("found %+v", found)
	}
	if !office.Secure || office.Port != 8631 || office.Host != "server.local" || !office.Local || office.URL != "ipps://127.0.0.1:8631/ipp/print/office" {
		t.Fatalf("office %+v", office)
	}
	if office.TXT["URF"] != "W8,SRGB24" || office.TXT["TLS"] != "1.2" {
		t.Fatalf("the keys of both records: %v", office.TXT)
	}
	if labels.Secure || labels.URL != "ipp://127.0.0.1:8631/ipp/print/labels" {
		t.Fatalf("labels %+v", labels)
	}
	// a subtype is asked as AirPrint's clients ask
	c, _ := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	defer c.Close()
	bld := dnsmessage.NewBuilder(nil, dnsmessage.Header{ID: 7})
	bld.StartQuestions()
	bld.Question(dnsmessage.Question{Name: dnsmessage.MustNewName("_universal._sub._ipp._tcp.local."), Type: dnsmessage.TypePTR, Class: dnsmessage.ClassINET})
	q, _ := bld.Finish()
	c.WriteTo(q, loopbackMDNS)
	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 9000)
	n, _, err := c.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	var p dnsmessage.Parser
	h, err := p.Start(buf[:n])
	if err != nil || h.ID != 7 || !h.Response {
		t.Fatalf("answer %+v %v", h, err)
	}
	qs, _ := p.AllQuestions()
	ans, _ := p.AllAnswers()
	if len(qs) != 1 || len(ans) != 1 || ans[0].Body.(*dnsmessage.PTRResource).PTR.String() != "Office @ server._ipp._tcp.local." || ans[0].Header.TTL > 10 {
		t.Fatalf("questions %v answers %v", qs, ans)
	}
	p.SkipAllAuthorities()
	extra, _ := p.AllAdditionals()
	types := map[dnsmessage.Type]bool{}
	for _, e := range extra {
		types[e.Header.Type] = true
	}
	if !types[dnsmessage.TypeSRV] || !types[dnsmessage.TypeTXT] || !types[dnsmessage.TypeA] {
		t.Fatalf("additional %v", extra)
	}
	// goodbye: the records with a TTL of 0
	cancel()
	<-r.done
}

// TestProbeConflict renames an instance another host answers for while
// it is probed.
func TestProbeConflict(t *testing.T) {
	r := &Responder{o: ResponderOptions{Host: "server.local."}, ads: []Advertisement{{Instance: "Office", Type: "_ipp._tcp"}},
		names: []string{"Office._ipp._tcp.local."}, probing: true, conflicts: map[string]bool{}}
	bld := dnsmessage.NewBuilder(nil, dnsmessage.Header{Response: true})
	bld.StartAnswers()
	bld.SRVResource(dnsmessage.ResourceHeader{Name: dnsmessage.MustNewName("office._ipp._tcp.local."), Type: dnsmessage.TypeSRV, Class: dnsmessage.ClassINET},
		dnsmessage.SRVResource{Port: 631, Target: dnsmessage.MustNewName("other.local.")})
	b, _ := bld.Finish()
	r.handle(nil, b, &net.UDPAddr{IP: net.IPv4(192, 0, 2, 1), Port: 5353}, 0)
	if !r.conflicts["office._ipp._tcp.local."] {
		t.Fatal("no conflict noted")
	}
	for in, want := range map[string]string{"Office": "Office (2)", "Office (2)": "Office (3)", "A (b)": "A (b) (2)",
		strings.Repeat("x", 63): strings.Repeat("x", 59) + " (2)"} {
		if got := renamedInstance(in); got != want {
			t.Errorf("%q renamed %q, not %q", in, got, want)
		}
	}
}
