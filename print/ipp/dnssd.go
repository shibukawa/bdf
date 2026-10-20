package ipp

import (
	"context"
	"errors"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

// Service is an IPP printer that DNS-SD (RFC 6763) on the local links
// (mDNS, RFC 6762) found.
type Service struct {
	// Name is the instance name the printer advertises ("Canon MF Series").
	Name string
	// Host and Port are from the SRV record: the printer's host name
	// ("printer.local") and the port of its IPP server.
	Host string
	Port int
	// Addrs are the addresses of the host, from its A and AAAA records.
	Addrs []net.IP
	// TXT holds the keys of the TXT record: rp (the path), pdl (the
	// formats), ty (the model), note, adminurl, UUID, URF, TLS, Color,
	// Duplex, kind, air, priority …
	TXT map[string]string
	// Secure tells the service was found as _ipps._tcp: it speaks TLS.
	Secure bool
	// Local tells a printer of this machine, at a loopback address: one
	// of USB that ipp-usb (Linux) or the system (macOS) shows as an IPP
	// printer, or a queue the machine shares.
	Local bool
	// URL is the printer's URL for a Client: ipp or ipps, the first
	// address (an IPv4 one when there is one) and the port, and the
	// path of rp. An address stands in for the host because the
	// resolution of .local names is not everywhere.
	URL string
}

// Formats returns the document formats the TXT record names (pdl).
func (s *Service) Formats() []string {
	if s.TXT["pdl"] == "" {
		return nil
	}
	return strings.Split(s.TXT["pdl"], ",")
}

// Browse looks for IPP printers on the local links for the timeout
// (two seconds when zero) and returns them, sorted by name. Both
// _ipp._tcp and _ipps._tcp are asked for; a printer that answers both
// is listed once, as secure.
//
// The queries are sent from a port of the kernel's choosing with the
// unicast-response bit, so the answers come back to it and the port of
// mDNS (5353) is not bound, which leaves Avahi and Bonjour alone;
// answers that printers multicast anyway are read from a listener on
// the group. The responder of this machine is also asked directly at
// 127.0.0.1 (RFC 6762 §5.5): printers of USB that speak IPP (IPP over
// USB) are advertised on the loopback, which has no multicast on Linux.
func Browse(ctx context.Context, timeout time.Duration) ([]Service, error) {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	b := &browser{found: map[string]*found{}}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	packets := make(chan []byte, 64)
	var conns []net.PacketConn
	var wg sync.WaitGroup
	for _, ifi := range ifaces {
		if ifi.Flags&net.FlagUp == 0 || ifi.Flags&net.FlagMulticast == 0 {
			continue
		}
		conns = append(conns, b.open(&ifi)...)
	}
	if c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}); err == nil {
		conns = append(conns, &conn{c, loopbackMDNS})
	}
	if len(conns) == 0 {
		return nil, errors.New("ipp: no network interface for multicast DNS")
	}
	for _, c := range conns {
		wg.Add(1)
		go func(c net.PacketConn) {
			defer wg.Done()
			buf := make([]byte, 9000)
			for {
				n, _, err := c.ReadFrom(buf)
				if err != nil {
					return
				}
				select {
				case packets <- append([]byte(nil), buf[:n]...):
				case <-ctx.Done():
					return
				}
			}
		}(c)
	}
	b.query(conns, ptrQuery())
	half := time.After(timeout / 2)
	asked := false
	for {
		select {
		case p := <-packets:
			b.parse(p)
		case <-half:
			if !asked {
				asked = true
				if q := b.followUp(); q != nil {
					b.query(conns, q)
				}
			}
		case <-ctx.Done():
			for _, c := range conns {
				c.Close()
			}
			wg.Wait()
			return b.services(), nil
		}
	}
}

const (
	ippService  = "_ipp._tcp.local."
	ippsService = "_ipps._tcp.local."
)

var (
	mdnsV4 = &net.UDPAddr{IP: net.IPv4(224, 0, 0, 251), Port: 5353}
	mdnsV6 = &net.UDPAddr{IP: net.ParseIP("ff02::fb"), Port: 5353}
	// loopbackMDNS is the responder of this machine (tests change it)
	loopbackMDNS = &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 5353}
)

type browser struct {
	mu    sync.Mutex
	found map[string]*found // by the instance's full name, lower-cased
	hosts map[string][]net.IP
}

type found struct {
	name   string // full escaped name, as in the records
	secure bool
	srv    *dnsmessage.SRVResource
	txt    []string
}

type conn struct {
	net.PacketConn
	dst *net.UDPAddr
}

// open opens, for an interface, a socket of a port of its own to query
// from, and a listener on the mDNS group, for IPv4 and for IPv6 as the
// interface has addresses.
func (b *browser) open(ifi *net.Interface) []net.PacketConn {
	var out []net.PacketConn
	addrs, err := ifi.Addrs()
	if err != nil {
		return nil
	}
	var v4, v6 net.IP
	for _, a := range addrs {
		ipn, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		if ip4 := ipn.IP.To4(); ip4 != nil {
			if v4 == nil {
				v4 = ip4
			}
		} else if ipn.IP.IsLinkLocalUnicast() && v6 == nil {
			v6 = ipn.IP
		}
	}
	if v4 != nil {
		if c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: v4}); err == nil {
			out = append(out, &conn{c, mdnsV4})
		}
		if c, err := net.ListenMulticastUDP("udp4", ifi, mdnsV4); err == nil {
			out = append(out, &conn{c, nil})
		}
	}
	if v6 != nil {
		dst := &net.UDPAddr{IP: mdnsV6.IP, Port: 5353, Zone: ifi.Name}
		if c, err := net.ListenUDP("udp6", &net.UDPAddr{IP: v6, Zone: ifi.Name}); err == nil {
			out = append(out, &conn{c, dst})
		}
		if c, err := net.ListenMulticastUDP("udp6", ifi, mdnsV6); err == nil {
			out = append(out, &conn{c, nil})
		}
	}
	return out
}

func (b *browser) query(conns []net.PacketConn, q []byte) {
	for _, c := range conns {
		if cc, ok := c.(*conn); ok && cc.dst != nil {
			cc.WriteTo(q, cc.dst)
		}
	}
}

// qu is the unicast-response bit of a question's class (RFC 6762 §5.4).
const qu = 0x8000

func ptrQuery() []byte {
	bld := dnsmessage.NewBuilder(nil, dnsmessage.Header{})
	bld.EnableCompression()
	bld.StartQuestions()
	for _, s := range []string{ippService, ippsService} {
		bld.Question(dnsmessage.Question{Name: dnsmessage.MustNewName(s), Type: dnsmessage.TypePTR, Class: dnsmessage.ClassINET | qu})
	}
	b, _ := bld.Finish()
	return b
}

// followUp asks for the records of the instances that lack their SRV,
// TXT or addresses.
func (b *browser) followUp() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	bld := dnsmessage.NewBuilder(nil, dnsmessage.Header{})
	bld.EnableCompression()
	bld.StartQuestions()
	n := 0
	for _, f := range b.found {
		name, err := dnsmessage.NewName(f.name)
		if err != nil {
			continue
		}
		if f.srv == nil {
			bld.Question(dnsmessage.Question{Name: name, Type: dnsmessage.TypeSRV, Class: dnsmessage.ClassINET | qu})
			n++
		}
		if f.txt == nil {
			bld.Question(dnsmessage.Question{Name: name, Type: dnsmessage.TypeTXT, Class: dnsmessage.ClassINET | qu})
			n++
		}
		if f.srv != nil && len(b.hosts[strings.ToLower(f.srv.Target.String())]) == 0 {
			bld.Question(dnsmessage.Question{Name: f.srv.Target, Type: dnsmessage.TypeA, Class: dnsmessage.ClassINET | qu})
			bld.Question(dnsmessage.Question{Name: f.srv.Target, Type: dnsmessage.TypeAAAA, Class: dnsmessage.ClassINET | qu})
			n++
		}
		if n > 50 {
			break
		}
	}
	if n == 0 {
		return nil
	}
	q, err := bld.Finish()
	if err != nil {
		return nil
	}
	return q
}

func (b *browser) parse(p []byte) {
	var par dnsmessage.Parser
	h, err := par.Start(p)
	if err != nil || !h.Response {
		return
	}
	if err := par.SkipAllQuestions(); err != nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	for section := 0; section < 3; section++ {
		for {
			var rh dnsmessage.ResourceHeader
			var err error
			switch section {
			case 0:
				rh, err = par.AnswerHeader()
			case 1:
				rh, err = par.AuthorityHeader()
			default:
				rh, err = par.AdditionalHeader()
			}
			if err != nil {
				break
			}
			if len(b.found) > 1000 {
				return
			}
			name := strings.ToLower(rh.Name.String())
			switch rh.Type {
			case dnsmessage.TypePTR:
				r, err := par.PTRResource()
				if err != nil {
					return
				}
				secure := name == ippsService
				if name != ippService && !secure {
					continue
				}
				inst := strings.ToLower(r.PTR.String())
				f := b.found[inst]
				if f == nil {
					f = &found{name: r.PTR.String()}
					b.found[inst] = f
				}
				f.secure = f.secure || secure
			case dnsmessage.TypeSRV:
				r, err := par.SRVResource()
				if err != nil {
					return
				}
				if f := b.found[name]; f != nil {
					r := r
					f.srv = &r
				}
			case dnsmessage.TypeTXT:
				r, err := par.TXTResource()
				if err != nil {
					return
				}
				if f := b.found[name]; f != nil {
					f.txt = r.TXT
				}
			case dnsmessage.TypeA:
				r, err := par.AResource()
				if err != nil {
					return
				}
				b.addHost(name, net.IP(r.A[:]))
			case dnsmessage.TypeAAAA:
				r, err := par.AAAAResource()
				if err != nil {
					return
				}
				b.addHost(name, net.IP(r.AAAA[:]))
			default:
				if err := par.SkipAnswer(); err != nil {
					return
				}
			}
		}
	}
}

func (b *browser) addHost(name string, ip net.IP) {
	if b.hosts == nil {
		b.hosts = map[string][]net.IP{}
	}
	for _, have := range b.hosts[name] {
		if have.Equal(ip) {
			return
		}
	}
	if len(b.hosts[name]) < 16 {
		b.hosts[name] = append(b.hosts[name], ip)
	}
}

func (b *browser) services() []Service {
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []Service
	for _, f := range b.found {
		if f.srv == nil {
			continue
		}
		s := Service{Name: instanceName(f.name), Host: strings.TrimSuffix(f.srv.Target.String(), "."), Port: int(f.srv.Port),
			Secure: f.secure, TXT: map[string]string{}}
		for _, kv := range f.txt {
			k, v, _ := strings.Cut(kv, "=")
			if k != "" {
				s.TXT[k] = v
			}
		}
		s.Addrs = append(s.Addrs, b.hosts[strings.ToLower(f.srv.Target.String())]...)
		sort.SliceStable(s.Addrs, func(i, j int) bool { return s.Addrs[i].To4() != nil && s.Addrs[j].To4() == nil })
		s.Local = strings.EqualFold(s.Host, "localhost")
		if len(s.Addrs) > 0 {
			s.Local = true
			for _, a := range s.Addrs {
				s.Local = s.Local && a.IsLoopback()
			}
		}
		host := s.Host
		if len(s.Addrs) > 0 {
			host = s.Addrs[0].String()
		}
		scheme := "ipp"
		if s.Secure {
			scheme = "ipps"
		}
		rp := strings.TrimPrefix(s.TXT["rp"], "/")
		if rp == "" {
			rp = "ipp/print"
		}
		s.URL = scheme + "://" + net.JoinHostPort(host, strconv.Itoa(s.Port)) + "/" + rp
		out = append(out, s)
	}
	return merge(out)
}

// merge lists a printer that advertises both _ipp._tcp and _ipps._tcp
// (two instances of one name on one host) once, as the secure one, with
// the keys of both TXT records (the secure one's where both have a key).
func merge(in []Service) []Service {
	byKey := map[string]int{}
	var out []Service
	for _, s := range in {
		k := strings.ToLower(s.Name) + "\x00" + strings.ToLower(s.Host)
		i, ok := byKey[k]
		if !ok {
			byKey[k] = len(out)
			out = append(out, s)
			continue
		}
		keep, other := out[i], s
		if s.Secure && !keep.Secure {
			keep, other = s, out[i]
		}
		for key, v := range other.TXT {
			if _, ok := keep.TXT[key]; !ok {
				keep.TXT[key] = v
			}
		}
		out[i] = keep
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// instanceName returns the first label of a service instance name,
// unescaped: "Canon MF Series" of "Canon\032MF\032Series._ipp._tcp.local.".
func instanceName(full string) string {
	var b strings.Builder
	for i := 0; i < len(full); i++ {
		c := full[i]
		switch {
		case c == '\\' && i+3 < len(full) && isDigit(full[i+1]) && isDigit(full[i+2]) && isDigit(full[i+3]):
			n, _ := strconv.Atoi(full[i+1 : i+4])
			b.WriteByte(byte(n))
			i += 3
		case c == '\\' && i+1 < len(full):
			b.WriteByte(full[i+1])
			i++
		case c == '.':
			return b.String()
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }
