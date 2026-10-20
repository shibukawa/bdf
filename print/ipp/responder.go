package ipp

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// Advertisement is a service to advertise with DNS-SD: an instance of a
// type ("_ipp._tcp"), with subtypes ("_print", "_universal"), on a port of
// this machine, with its TXT record.
type Advertisement struct {
	// Instance is the instance's name, a label of 63 bytes at most
	// ("Office @ server"); a name another host has is renamed "Office @
	// server (2)" and so on.
	Instance string
	Type     string
	Subtypes []string
	Port     int
	// TXT are the keys and values of the TXT record, "key=value" each.
	TXT []string
}

// ResponderOptions configure a Responder.
type ResponderOptions struct {
	// Host is the host name the services are on, in .local
	// ("server-bdf.local"), whose addresses the responder answers: a name
	// of its own, not the system's (whose responder would take the
	// records of another as a conflict, and rename the machine). It is
	// probed as the instances are, and renamed "server-bdf-2.local" when
	// another host has it.
	Host string
	// Log receives what happens; nothing by default.
	Log *slog.Logger
	// Port is the port of mDNS: 5353 but in tests, whose responders are
	// asked directly on the loopback.
	Port int
	// Interfaces bound the interfaces advertised on; all those up with
	// multicast by default. Loopback adds 127.0.0.1 and ::1 (for tests).
	Interfaces []net.Interface
	Loopback   bool
}

// Responder advertises services with multicast DNS (RFC 6762) and DNS-SD
// (RFC 6763): it probes for its instance names, announces its records,
// answers the questions for them, and says goodbye when it stops. It
// shares the port of mDNS with the system's responder (Bonjour, Avahi,
// Windows), as both bind it reusably; the records of the host's address
// it answers are those the system's responder answers too, so they agree.
type Responder struct {
	o   ResponderOptions
	ads []Advertisement
	log *slog.Logger

	mu    sync.Mutex
	conns []*mdnsConn
	// names are the instances' full names as they are after probing;
	// host is the host's.
	names   []string
	host    string
	probing bool
	// conflicts are names another host answered for while probing.
	conflicts map[string]bool
	done      chan struct{}
}

// mdnsConn is a socket of mDNS, joined to the group on the interfaces.
type mdnsConn struct {
	v4   *ipv4.PacketConn
	v6   *ipv6.PacketConn
	conn net.PacketConn
	// group is where multicast goes: the group of mDNS, or the address
	// the responder of a test is asked at.
	group *net.UDPAddr
	ifs   []net.Interface
}

// TTLs of the records (RFC 6762 §10): those of the host 120 seconds, the
// others 75 minutes.
const (
	ttlHost  = 120
	ttlOther = 4500
	// cacheFlush is the bit of a record's class that says the record is
	// the whole of its set (RFC 6762 §10.2).
	cacheFlush = 0x8000
)

// Advertise advertises services until ctx ends (it then says goodbye):
// it returns once the names are probed and announced.
func Advertise(ctx context.Context, ads []Advertisement, o ResponderOptions) (*Responder, error) {
	if o.Host == "" {
		return nil, errors.New("ipp: advertise: no host name")
	}
	if !strings.HasSuffix(o.Host, ".") {
		o.Host += "."
	}
	if o.Port == 0 {
		o.Port = 5353
	}
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	r := &Responder{o: o, ads: slices.Clone(ads), log: o.Log, conflicts: map[string]bool{}, done: make(chan struct{}), host: o.Host}
	for i := range r.ads {
		r.ads[i].Instance = cleanInstance(r.ads[i].Instance)
	}
	if err := r.open(); err != nil {
		return nil, err
	}
	for _, c := range r.conns {
		go r.read(ctx, c)
	}
	r.probe(ctx)
	r.announce(false)
	go func() {
		defer close(r.done)
		t := time.NewTimer(time.Second)
		defer t.Stop()
		select {
		case <-ctx.Done():
			r.Close()
			return
		case <-t.C:
		}
		r.announce(false)
		<-ctx.Done()
		r.Close()
	}()
	return r, nil
}

// Names are the services' full instance names, as advertised.
func (r *Responder) Names() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.names)
}

// Host is the host name advertised ("server-bdf.local."), as probed.
func (r *Responder) Host() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.host
}

// SetTXT changes the TXT records of the services (in the order they were
// given) and announces them: what a printer can do is known after it is
// advertised.
func (r *Responder) SetTXT(txt [][]string) {
	r.mu.Lock()
	changed := false
	for i := range r.ads {
		if i < len(txt) && !slices.Equal(r.ads[i].TXT, txt[i]) {
			r.ads[i].TXT = slices.Clone(txt[i])
			changed = true
		}
	}
	r.mu.Unlock()
	if changed {
		r.announce(false)
	}
}

// Close says goodbye (the records with a TTL of 0) and closes the sockets.
func (r *Responder) Close() {
	r.mu.Lock()
	conns := r.conns
	r.conns = nil
	r.mu.Unlock()
	if conns == nil {
		return
	}
	for _, c := range conns {
		r.sendAll(c, true)
		c.conn.Close()
	}
}

// cleanInstance makes a name a label: no dots, 63 bytes at most.
func cleanInstance(s string) string {
	s = strings.ReplaceAll(s, ".", " ")
	for len(s) > 63 {
		rs := []rune(s)
		s = string(rs[:len(rs)-1])
	}
	return strings.TrimSpace(s)
}

// open opens a socket of mDNS for IPv4 and one for IPv6, joined to the
// group on the interfaces that have such addresses.
func (r *Responder) open() error {
	ifs := r.o.Interfaces
	if ifs == nil {
		all, err := net.Interfaces()
		if err != nil {
			return err
		}
		for _, ifi := range all {
			if ifi.Flags&net.FlagUp != 0 && ifi.Flags&net.FlagMulticast != 0 && ifi.Flags&net.FlagLoopback == 0 && ifi.Flags&net.FlagPointToPoint == 0 {
				ifs = append(ifs, ifi)
			}
		}
	}
	port := strconv.Itoa(r.o.Port)
	// a multicast address makes the listener bind the wildcard address,
	// reusably, as the system's responder does
	if c, err := net.ListenPacket("udp4", "224.0.0.251:"+port); err == nil {
		mc := &mdnsConn{conn: c, v4: ipv4.NewPacketConn(c), group: &net.UDPAddr{IP: net.IPv4(224, 0, 0, 251), Port: r.o.Port}}
		mc.v4.SetControlMessage(ipv4.FlagInterface, true)
		mc.v4.SetMulticastTTL(255)
		mc.v4.SetMulticastLoopback(true)
		for _, ifi := range ifs {
			if hasAddr(ifi, true) && mc.v4.JoinGroup(&ifi, mc.group) == nil {
				mc.ifs = append(mc.ifs, ifi)
			}
		}
		if r.o.Loopback {
			mc.group = &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: r.o.Port}
		}
		if len(mc.ifs) > 0 || r.o.Loopback {
			r.conns = append(r.conns, mc)
		} else {
			c.Close()
		}
	}
	if c, err := net.ListenPacket("udp6", "[ff02::fb]:"+port); err == nil {
		mc := &mdnsConn{conn: c, v6: ipv6.NewPacketConn(c), group: &net.UDPAddr{IP: net.ParseIP("ff02::fb"), Port: r.o.Port}}
		mc.v6.SetControlMessage(ipv6.FlagInterface, true)
		mc.v6.SetMulticastHopLimit(255)
		mc.v6.SetMulticastLoopback(true)
		for _, ifi := range ifs {
			if hasAddr(ifi, false) && mc.v6.JoinGroup(&ifi, mc.group) == nil {
				mc.ifs = append(mc.ifs, ifi)
			}
		}
		if len(mc.ifs) > 0 && !r.o.Loopback {
			r.conns = append(r.conns, mc)
		} else {
			c.Close()
		}
	}
	if len(r.conns) == 0 {
		return errors.New("ipp: advertise: no interface for multicast DNS")
	}
	return nil
}

// hasAddr tells an interface with an IPv4 address (v4) or a link-local
// IPv6 one.
func hasAddr(ifi net.Interface, v4 bool) bool {
	addrs, err := ifi.Addrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		if ipn, ok := a.(*net.IPNet); ok && (ipn.IP.To4() != nil) == v4 && (v4 || ipn.IP.IsLinkLocalUnicast()) {
			return true
		}
	}
	return false
}

// addrs are the addresses of the host on an interface, or on all of the
// socket's (ifIndex 0); the loopback's for the responder of a test.
func (r *Responder) addrs(c *mdnsConn, ifIndex int) []net.IP {
	var out []net.IP
	if r.o.Loopback {
		return []net.IP{net.IPv4(127, 0, 0, 1)}
	}
	for _, ifi := range c.ifs {
		if ifIndex != 0 && ifi.Index != ifIndex {
			continue
		}
		as, _ := ifi.Addrs()
		for _, a := range as {
			ipn, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			if ip4 := ipn.IP.To4(); ip4 != nil {
				out = append(out, ip4)
			} else if !ipn.IP.IsLinkLocalMulticast() {
				out = append(out, ipn.IP)
			}
		}
		if ifIndex != 0 {
			break
		}
	}
	return out
}

// read answers the questions a socket receives, and notes the answers
// for names being probed.
func (r *Responder) read(ctx context.Context, c *mdnsConn) {
	buf := make([]byte, 9000)
	for {
		var n, ifIndex int
		var src net.Addr
		var err error
		if c.v4 != nil {
			var cm *ipv4.ControlMessage
			n, cm, src, err = c.v4.ReadFrom(buf)
			if cm != nil {
				ifIndex = cm.IfIndex
			}
		} else {
			var cm *ipv6.ControlMessage
			n, cm, src, err = c.v6.ReadFrom(buf)
			if cm != nil {
				ifIndex = cm.IfIndex
			}
		}
		if err != nil {
			return
		}
		if ctx.Err() != nil {
			return
		}
		r.handle(c, buf[:n], src, ifIndex)
	}
}

// handle answers a query, or notes a conflict in a response.
func (r *Responder) handle(c *mdnsConn, p []byte, src net.Addr, ifIndex int) {
	var par dnsmessage.Parser
	h, err := par.Start(p)
	if err != nil {
		return
	}
	if h.Response {
		r.noteConflicts(&par)
		return
	}
	qs, err := par.AllQuestions()
	if err != nil || len(qs) == 0 {
		return
	}
	r.mu.Lock()
	probing := r.probing
	r.mu.Unlock()
	if probing {
		return
	}
	udp, _ := src.(*net.UDPAddr)
	// a question from a port that is not mDNS's is a legacy unicast one
	// (RFC 6762 §6.7): answered to it alone, with the questions; the
	// unicast-response bit asks the same of a multicast question
	legacy := udp != nil && udp.Port != r.o.Port
	unicast := legacy
	var answers []dnsmessage.Resource
	for _, q := range qs {
		if q.Class&qu != 0 {
			unicast = true
		}
		answers = append(answers, r.answer(c, q, ifIndex)...)
	}
	if len(answers) == 0 {
		return
	}
	extra := r.additional(c, answers, ifIndex)
	var to *net.UDPAddr
	if unicast {
		to = udp
	}
	if legacy {
		for i := range answers {
			answers[i].Header.Class &^= cacheFlush
			answers[i].Header.TTL = min(answers[i].Header.TTL, 10)
		}
		for i := range extra {
			extra[i].Header.Class &^= cacheFlush
			extra[i].Header.TTL = min(extra[i].Header.TTL, 10)
		}
		r.sendID(c, to, h.ID, qs, answers, extra, ifIndex)
		return
	}
	r.sendID(c, to, 0, nil, answers, extra, ifIndex)
}

// noteConflicts notes the answers of other hosts for names being probed.
func (r *Responder) noteConflicts(par *dnsmessage.Parser) {
	if err := par.SkipAllQuestions(); err != nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.probing {
		return
	}
	for {
		rh, err := par.AnswerHeader()
		if err != nil {
			return
		}
		name := strings.ToLower(rh.Name.String())
		for _, n := range append(slices.Clone(r.names), r.host) {
			if strings.ToLower(n) == name {
				r.conflicts[name] = true
			}
		}
		if err := par.SkipAnswer(); err != nil {
			return
		}
	}
}

// probe asks three times, a quarter of a second apart, whether another
// host has the instance names (RFC 6762 §8.1), and renames those that one
// has, until none has.
func (r *Responder) probe(ctx context.Context) {
	r.mu.Lock()
	r.names = make([]string, len(r.ads))
	for i, a := range r.ads {
		r.names[i] = a.Instance + "." + a.Type + ".local."
	}
	r.mu.Unlock()
	for round := 0; round < 10; round++ {
		r.mu.Lock()
		r.probing, r.conflicts = true, map[string]bool{}
		r.mu.Unlock()
		// a random wait first, so that hosts that start together do not
		// probe together
		time.Sleep(time.Duration(rand.IntN(250)) * time.Millisecond)
		for range 3 {
			for _, c := range r.conns {
				r.sendProbe(c)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(250 * time.Millisecond):
			}
		}
		r.mu.Lock()
		renamed := false
		if r.conflicts[strings.ToLower(r.host)] {
			r.host = renamedHost(r.host)
			r.log.Info("mDNS: the host name is another host's: renamed", "host", r.host)
			renamed = true
		}
		for i, a := range r.ads {
			if r.conflicts[strings.ToLower(r.names[i])] {
				r.ads[i].Instance = renamedInstance(a.Instance)
				r.names[i] = r.ads[i].Instance + "." + a.Type + ".local."
				r.log.Info("mDNS: the name is another host's: renamed", "instance", r.ads[i].Instance)
				renamed = true
			}
		}
		r.probing = false
		r.mu.Unlock()
		if !renamed {
			return
		}
	}
}

// renamedInstance is "name (2)", "name (3)" … of a name another host has.
func renamedInstance(s string) string {
	n := 2
	if i := strings.LastIndex(s, " ("); i >= 0 && strings.HasSuffix(s, ")") {
		if k, err := strconv.Atoi(s[i+2 : len(s)-1]); err == nil {
			n, s = k+1, s[:i]
		}
	}
	suffix := " (" + strconv.Itoa(n) + ")"
	return cleanInstance(s[:min(len(s), 63-len(suffix))] + suffix)
}

// renamedHost is "name-2.local.", "name-3.local." … of a host name another
// host has.
func renamedHost(h string) string {
	label, rest, _ := strings.Cut(h, ".")
	n := 2
	if i := strings.LastIndexByte(label, '-'); i >= 0 {
		if k, err := strconv.Atoi(label[i+1:]); err == nil && k >= 2 {
			n, label = k+1, label[:i]
		}
	}
	return label + "-" + strconv.Itoa(n) + "." + rest
}

// sendProbe asks for the instance names, with the records the responder
// would answer in the authority section (RFC 6762 §8.2).
func (r *Responder) sendProbe(c *mdnsConn) {
	if len(c.ifs) == 0 {
		r.sendProbeOn(c, 0)
		return
	}
	for _, ifi := range c.ifs {
		r.sendProbeOn(c, ifi.Index)
	}
}

func (r *Responder) sendProbeOn(c *mdnsConn, ifIndex int) {
	bld := dnsmessage.NewBuilder(nil, dnsmessage.Header{})
	bld.EnableCompression()
	bld.StartQuestions()
	r.mu.Lock()
	names := append(slices.Clone(r.names), r.host)
	r.mu.Unlock()
	for _, n := range names {
		name, err := dnsmessage.NewName(n)
		if err != nil {
			continue
		}
		bld.Question(dnsmessage.Question{Name: name, Type: dnsmessage.TypeALL, Class: dnsmessage.ClassINET | qu})
	}
	bld.StartAuthorities()
	for i := range r.ads {
		for _, rr := range r.instanceRecords(i) {
			addResource(&bld, rr)
		}
	}
	for _, rr := range r.hostRecords(c, ifIndex) {
		addResource(&bld, rr)
	}
	b, err := bld.Finish()
	if err != nil {
		return
	}
	r.write(c, nil, b, ifIndex)
}

// announce sends all the records unsolicited (RFC 6762 §8.3), or with a
// TTL of 0 to say goodbye.
func (r *Responder) announce(goodbye bool) {
	r.mu.Lock()
	conns := slices.Clone(r.conns)
	r.mu.Unlock()
	for _, c := range conns {
		r.sendAll(c, goodbye)
	}
}

// sendAll sends every record on each interface of a socket, with the
// host's addresses on that interface.
func (r *Responder) sendAll(c *mdnsConn, goodbye bool) {
	if len(c.ifs) == 0 {
		r.sendID(c, nil, 0, nil, r.answerAll(c, 0, goodbye), nil, 0)
		return
	}
	for _, ifi := range c.ifs {
		r.sendID(c, nil, 0, nil, r.answerAll(c, ifi.Index, goodbye), nil, ifi.Index)
	}
}

// answerAll is every record of the services and the host (its addresses
// on an interface).
func (r *Responder) answerAll(c *mdnsConn, ifIndex int, goodbye bool) []dnsmessage.Resource {
	var out []dnsmessage.Resource
	r.mu.Lock()
	n := len(r.ads)
	r.mu.Unlock()
	for i := 0; i < n; i++ {
		out = append(out, r.ptrRecords(i)...)
		out = append(out, r.instanceRecords(i)...)
	}
	out = append(out, r.hostRecords(c, ifIndex)...)
	if goodbye {
		for i := range out {
			out[i].Header.TTL = 0
		}
	}
	return out
}

// ptrRecords are the shared records that point at a service: its type
// (and the type in the list of types), and its subtypes.
func (r *Responder) ptrRecords(i int) []dnsmessage.Resource {
	r.mu.Lock()
	a, inst := r.ads[i], r.names[i]
	r.mu.Unlock()
	typ := a.Type + ".local."
	var out []dnsmessage.Resource
	add := func(name, target string) {
		n, err1 := dnsmessage.NewName(name)
		t, err2 := dnsmessage.NewName(target)
		if err1 != nil || err2 != nil {
			return
		}
		out = append(out, dnsmessage.Resource{Header: dnsmessage.ResourceHeader{Name: n, Type: dnsmessage.TypePTR, Class: dnsmessage.ClassINET, TTL: ttlOther},
			Body: &dnsmessage.PTRResource{PTR: t}})
	}
	add(typ, inst)
	add("_services._dns-sd._udp.local.", typ)
	for _, s := range a.Subtypes {
		add(s+"._sub."+typ, inst)
	}
	return out
}

// instanceRecords are the unique records of a service: its SRV and TXT.
func (r *Responder) instanceRecords(i int) []dnsmessage.Resource {
	r.mu.Lock()
	a, inst := r.ads[i], r.names[i]
	r.mu.Unlock()
	r.mu.Lock()
	hn := r.host
	r.mu.Unlock()
	name, err := dnsmessage.NewName(inst)
	host, err2 := dnsmessage.NewName(hn)
	if err != nil || err2 != nil {
		return nil
	}
	txt := a.TXT
	if len(txt) == 0 {
		txt = []string{""}
	}
	return []dnsmessage.Resource{
		{Header: dnsmessage.ResourceHeader{Name: name, Type: dnsmessage.TypeSRV, Class: dnsmessage.ClassINET | cacheFlush, TTL: ttlHost},
			Body: &dnsmessage.SRVResource{Port: uint16(a.Port), Target: host}},
		{Header: dnsmessage.ResourceHeader{Name: name, Type: dnsmessage.TypeTXT, Class: dnsmessage.ClassINET | cacheFlush, TTL: ttlOther},
			Body: &dnsmessage.TXTResource{TXT: txt}},
	}
}

// hostRecords are the addresses of the host on an interface (all of the
// socket's when 0).
func (r *Responder) hostRecords(c *mdnsConn, ifIndex int) []dnsmessage.Resource {
	r.mu.Lock()
	hn := r.host
	r.mu.Unlock()
	host, err := dnsmessage.NewName(hn)
	if err != nil {
		return nil
	}
	var out []dnsmessage.Resource
	for _, ip := range r.addrs(c, ifIndex) {
		hd := dnsmessage.ResourceHeader{Name: host, Class: dnsmessage.ClassINET | cacheFlush, TTL: ttlHost}
		if ip4 := ip.To4(); ip4 != nil {
			hd.Type = dnsmessage.TypeA
			out = append(out, dnsmessage.Resource{Header: hd, Body: &dnsmessage.AResource{A: [4]byte(ip4)}})
		} else {
			hd.Type = dnsmessage.TypeAAAA
			out = append(out, dnsmessage.Resource{Header: hd, Body: &dnsmessage.AAAAResource{AAAA: [16]byte(ip.To16())}})
		}
	}
	return out
}

// answer returns the records that answer a question.
func (r *Responder) answer(c *mdnsConn, q dnsmessage.Question, ifIndex int) []dnsmessage.Resource {
	name := strings.ToLower(q.Name.String())
	var all []dnsmessage.Resource
	r.mu.Lock()
	n := len(r.ads)
	r.mu.Unlock()
	for i := 0; i < n; i++ {
		all = append(all, r.ptrRecords(i)...)
		all = append(all, r.instanceRecords(i)...)
	}
	all = append(all, r.hostRecords(c, ifIndex)...)
	var out []dnsmessage.Resource
	for _, rr := range all {
		if strings.ToLower(rr.Header.Name.String()) != name || (q.Type != dnsmessage.TypeALL && q.Type != rr.Header.Type) {
			continue
		}
		dup := false
		for _, o := range out {
			if sameRecord(o, rr) {
				dup = true
			}
		}
		if !dup {
			out = append(out, rr)
		}
	}
	return out
}

func sameRecord(a, b dnsmessage.Resource) bool {
	return a.Header.Type == b.Header.Type && strings.EqualFold(a.Header.Name.String(), b.Header.Name.String()) && a.Body.GoString() == b.Body.GoString()
}

// additional are the records that go with answers (RFC 6763 §12): the SRV,
// TXT and addresses of the instances a PTR points at, the addresses of
// an SRV's host.
func (r *Responder) additional(c *mdnsConn, answers []dnsmessage.Resource, ifIndex int) []dnsmessage.Resource {
	var out []dnsmessage.Resource
	add := func(rrs []dnsmessage.Resource) {
		for _, rr := range rrs {
			dup := false
			for _, o := range append(answers, out...) {
				if sameRecord(o, rr) {
					dup = true
					break
				}
			}
			if !dup {
				out = append(out, rr)
			}
		}
	}
	hostAdded := false
	for _, a := range answers {
		switch body := a.Body.(type) {
		case *dnsmessage.PTRResource:
			target := strings.ToLower(body.PTR.String())
			r.mu.Lock()
			idx := slices.IndexFunc(r.names, func(n string) bool { return strings.ToLower(n) == target })
			r.mu.Unlock()
			if idx >= 0 {
				add(r.instanceRecords(idx))
				if !hostAdded {
					add(r.hostRecords(c, ifIndex))
					hostAdded = true
				}
			}
		case *dnsmessage.SRVResource:
			if !hostAdded {
				add(r.hostRecords(c, ifIndex))
				hostAdded = true
			}
		}
	}
	return out
}

// send sends a response of records: multicast, or to an address.
func (r *Responder) send(c *mdnsConn, to *net.UDPAddr, answers, extra []dnsmessage.Resource) {
	r.sendID(c, to, 0, nil, answers, extra, 0)
}

func (r *Responder) sendID(c *mdnsConn, to *net.UDPAddr, id uint16, qs []dnsmessage.Question, answers, extra []dnsmessage.Resource, ifIndex int) {
	if len(answers) == 0 {
		return
	}
	bld := dnsmessage.NewBuilder(nil, dnsmessage.Header{ID: id, Response: true, Authoritative: true})
	bld.EnableCompression()
	bld.StartQuestions()
	for _, q := range qs {
		bld.Question(q)
	}
	bld.StartAnswers()
	for _, rr := range answers {
		addResource(&bld, rr)
	}
	bld.StartAdditionals()
	for _, rr := range extra {
		addResource(&bld, rr)
	}
	b, err := bld.Finish()
	if err != nil {
		r.log.Debug("mDNS: a response cannot be built", "err", err)
		return
	}
	r.write(c, to, b, ifIndex)
}

// write sends a message to an address, or to the group on each interface
// (or the one of ifIndex).
func (r *Responder) write(c *mdnsConn, to *net.UDPAddr, b []byte, ifIndex int) {
	if to != nil {
		c.conn.WriteTo(b, to)
		return
	}
	if r.o.Loopback && c.v4 != nil {
		c.conn.WriteTo(b, c.group)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ifi := range c.ifs {
		if ifIndex != 0 && ifi.Index != ifIndex {
			continue
		}
		if c.v4 != nil {
			c.v4.SetMulticastInterface(&ifi)
		} else {
			c.v6.SetMulticastInterface(&ifi)
		}
		dst := *c.group
		if c.v6 != nil {
			dst.Zone = ifi.Name
		}
		c.conn.WriteTo(b, &dst)
	}
}

func addResource(b *dnsmessage.Builder, rr dnsmessage.Resource) {
	switch body := rr.Body.(type) {
	case *dnsmessage.PTRResource:
		b.PTRResource(rr.Header, *body)
	case *dnsmessage.SRVResource:
		b.SRVResource(rr.Header, *body)
	case *dnsmessage.TXTResource:
		b.TXTResource(rr.Header, *body)
	case *dnsmessage.AResource:
		b.AResource(rr.Header, *body)
	case *dnsmessage.AAAAResource:
		b.AAAAResource(rr.Header, *body)
	}
}
