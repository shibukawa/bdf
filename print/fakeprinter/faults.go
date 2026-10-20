package fakeprinter

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"sync/atomic"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// Fault is something that happens to a printer: when it happens (after a
// number of jobs, after a time, or both: the time counts from the last of
// the jobs) and what happens. The tests and the farm's YAML write the same
// faults, so that a failure happens in the same order every time.
type Fault struct {
	// AfterJobs waits for so many more jobs to be received (their last
	// document) before the fault happens; 0 does not wait.
	AfterJobs int `yaml:"after-jobs" json:"after-jobs,omitempty"`
	// After waits so long, from when the fault is injected or from its
	// last job.
	After time.Duration `yaml:"after" json:"after,omitempty"`

	// Stop stops the printer for these reasons (media-empty-error,
	// media-jam-error, door-open-error, toner-empty-error …); For resumes
	// it after so long (0: until Resume).
	Stop []string      `yaml:"stop" json:"stop,omitempty"`
	For  time.Duration `yaml:"for" json:"for,omitempty"`
	// Resume makes a stopped printer idle again.
	Resume bool `yaml:"resume" json:"resume,omitempty"`
	// Reject makes the printer refuse new jobs (not-accepting-jobs, and
	// LPD's refusal); Accept takes them again.
	Reject bool `yaml:"reject" json:"reject,omitempty"`
	Accept bool `yaml:"accept" json:"accept,omitempty"`
	// Forget drops the jobs the printer knows; Restart also counts the
	// job ids from 1 again.
	Forget  bool `yaml:"forget" json:"forget,omitempty"`
	Restart bool `yaml:"restart" json:"restart,omitempty"`
	// Abort makes the next job aborted.
	Abort bool `yaml:"abort" json:"abort,omitempty"`
	// CutAfter cuts the connection of the next document after so many
	// bytes (-1 or unset: not).
	CutAfter *int64 `yaml:"cut-after" json:"cut-after,omitempty"`
	// DropCreate cuts the connection of the next Create-Job before its
	// answer, once the job is made.
	DropCreate bool `yaml:"drop-create" json:"drop-create,omitempty"`
	// Busy refuses so many of the next jobs with server-error-busy.
	Busy int `yaml:"busy" json:"busy,omitempty"`
	// Latency and Bandwidth, when set, replace the printer's.
	Latency   *time.Duration `yaml:"latency" json:"latency,omitempty"`
	Bandwidth *int           `yaml:"bandwidth" json:"bandwidth,omitempty"`
	// NeverComplete, when set, keeps the jobs processing (or lets them
	// complete again).
	NeverComplete *bool `yaml:"never-complete" json:"never-complete,omitempty"`
	// NewCertificate makes the printer present a new self-signed
	// certificate over ipps (a client that pinned the old one refuses
	// it); ExpiredCertificate one that has expired.
	NewCertificate     bool `yaml:"new-certificate" json:"new-certificate,omitempty"`
	ExpiredCertificate bool `yaml:"expired-certificate" json:"expired-certificate,omitempty"`

	jobs  int
	armed bool
}

// Inject adds a fault to the printer: it happens at once, or when its
// jobs have come and its time has passed.
func (p *Printer) Inject(f Fault) {
	p.lock()
	defer p.mu.Unlock()
	fc := f
	p.faults = append(p.faults, &fc)
	p.fire(false)
}

// fire counts a received job (jobCame) against the faults waiting for
// jobs and makes those that are due happen (now, or after their time).
// The printer is locked.
func (p *Printer) fire(jobCame bool) {
	keep := p.faults[:0]
	for _, f := range p.faults {
		if f.armed {
			keep = append(keep, f)
			continue
		}
		if f.jobs < f.AfterJobs {
			if jobCame {
				f.jobs++
			}
			if f.jobs < f.AfterJobs {
				keep = append(keep, f)
				continue
			}
		}
		if f.After > 0 {
			f.armed = true
			fc := f
			time.AfterFunc(f.After, func() {
				p.lock()
				defer p.mu.Unlock()
				p.apply(fc)
				for i, g := range p.faults {
					if g == fc {
						p.faults = append(p.faults[:i], p.faults[i+1:]...)
						break
					}
				}
			})
			keep = append(keep, f)
			continue
		}
		p.apply(f)
	}
	p.faults = keep
}

// apply makes a fault happen. The printer is locked.
func (p *Printer) apply(f *Fault) {
	if len(f.Stop) > 0 {
		p.State, p.Reasons, p.stateSince = ipp.PrinterStopped, append([]string(nil), f.Stop...), time.Now()
		if f.For > 0 {
			time.AfterFunc(f.For, func() {
				p.lock()
				defer p.mu.Unlock()
				p.State, p.Reasons, p.stateSince = ipp.PrinterIdle, []string{"none"}, time.Now()
				p.notify(Event{Kind: "state"})
			})
		}
	}
	if f.Resume {
		p.State, p.Reasons, p.stateSince = ipp.PrinterIdle, []string{"none"}, time.Now()
	}
	if f.Reject {
		p.Accepting, p.refuse = false, true
	}
	if f.Accept {
		p.Accepting, p.refuse = true, false
	}
	if f.Forget || f.Restart {
		p.jobs = nil
	}
	if f.Restart {
		p.ids.reset()
		p.started = time.Now()
	}
	if f.Abort {
		p.abortNext = true
	}
	if f.CutAfter != nil {
		p.cutAfter = *f.CutAfter
	}
	if f.DropCreate {
		p.dropCreate = true
	}
	p.Busy += f.Busy
	if f.Latency != nil {
		p.Latency = *f.Latency
	}
	if f.Bandwidth != nil {
		p.Bandwidth = *f.Bandwidth
	}
	if f.NeverComplete != nil {
		p.NeverComplete = *f.NeverComplete
	}
	if f.NewCertificate || f.ExpiredCertificate {
		p.cert().renew(p.certHosts(), f.ExpiredCertificate)
	}
	p.notify(Event{Kind: "state"})
}

// certificate is the certificate a printer presents over ipps, which a
// fault can change while it is served.
type certificate struct {
	cur atomic.Pointer[tls.Certificate]
}

// cert returns the printer's certificate holder. The printer is locked.
func (p *Printer) cert() *certificate {
	if p.certs == nil {
		p.certs = &certificate{}
	}
	return p.certs
}

// certHosts are the names the certificate is for.
func (p *Printer) certHosts() []string { return []string{"localhost", "127.0.0.1", "::1"} }

// renew makes a new self-signed certificate for the hosts; expired makes
// one whose time has passed.
func (c *certificate) renew(hosts []string, expired bool) error {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 62))
	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: serial, Subject: pkix.Name{CommonName: "BDF Fake Printer"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(24 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	if expired {
		tmpl.NotBefore, tmpl.NotAfter = now.Add(-48*time.Hour), now.Add(-24*time.Hour)
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, h)
		}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return err
	}
	c.cur.Store(&tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key})
	return nil
}

// get is tls.Config.GetCertificate.
func (c *certificate) get(*tls.ClientHelloInfo) (*tls.Certificate, error) { return c.cur.Load(), nil }
