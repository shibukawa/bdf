package server

import (
	"context"
	"sync"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// PrinterCheck is what a printer of a configuration answered.
type PrinterCheck struct {
	Name string
	// URL is the printer's, without a user and password.
	URL     string
	Printer *ipp.Printer
	Err     error
}

// CheckPrinters asks each printer of a configuration what it supports,
// without opening the spool. The certificate of an ipps printer is
// checked against its pinned fingerprint; one not pinned is taken as the
// first one seen (the fingerprints the spool recorded are not read).
func CheckPrinters(ctx context.Context, cfg *Config) []PrinterCheck {
	ts := &ipp.TrustStore{Backend: memoryTrust{}, Pinned: cfg.pinned()}
	names := sortedKeys(cfg.Printers)
	out := make([]PrinterCheck, len(names))
	var wg sync.WaitGroup
	for i, name := range names {
		pc := cfg.Printers[name]
		out[i] = PrinterCheck{Name: name, URL: redact(pc.URL)}
		c, err := pc.client(cfg.User)
		if err != nil {
			out[i].Err = err
			continue
		}
		if !c.TLS.Insecure {
			c.TLS.Store = ts
		}
		c.Timeout = 15 * time.Second
		wg.Go(func() {
			out[i].Printer, out[i].Err = c.GetPrinterAttributes(ctx)
		})
	}
	wg.Wait()
	return out
}

// memoryTrust keeps fingerprints for a while only.
type memoryTrust struct{}

func (memoryTrust) Load() (map[string]string, error) { return map[string]string{}, nil }
func (memoryTrust) Save(map[string]string) error     { return nil }
