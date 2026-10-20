package fakeprinter_test

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os/exec"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
)

// startIPPEve starts CUPS's ippeveprinter as an office printer of PDF,
// PWG Raster and JPEG, both sides, or skips where it is not (it needs
// DNS-SD: Avahi on Linux).
func startIPPEve(t *testing.T) string {
	t.Helper()
	tool, err := exec.LookPath("ippeveprinter")
	if err != nil {
		t.Skip("no ippeveprinter")
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	cmd := exec.Command(tool, "-2", "-p", strconv.Itoa(port), "-n", "localhost", "-d", t.TempDir(), "-s", "120",
		"-f", "application/pdf,image/pwg-raster,image/jpeg", "BDF Diff "+strconv.Itoa(port))
	var log bytes.Buffer
	cmd.Stdout, cmd.Stderr = &log, &log
	if err := cmd.Start(); err != nil {
		t.Skip("ippeveprinter:", err)
	}
	t.Cleanup(func() { cmd.Process.Kill(); cmd.Wait() })
	uri := fmt.Sprintf("ipp://localhost:%d/ipp/print", port)
	c := &ipp.Client{URL: uri}
	for range 100 {
		if _, err := c.GetPrinterAttributes(context.Background(), "printer-state"); err == nil {
			return uri
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Skipf("ippeveprinter did not start:\n%s", log.String())
	return ""
}

// step is a request of the comparison, made for a printer from what
// earlier steps answered it (the ids of its jobs).
type step struct {
	name string
	req  func(c *ipp.Client, ids []int32) (*ipp.Message, []byte)
	// attrs must be in the answers of both
	attrs []string
	// job keeps the job-id of the answer
	job bool
	// settle waits for the jobs to be done first: ippeveprinter takes
	// one job at a time, and is busy while it prints
	settle bool
	// either are statuses both of which are right, where the two differ
	either []ipp.Status
}

func op(c *ipp.Client, o ipp.Op, attrs ...ipp.Attribute) *ipp.Message {
	m := c.NewRequest(o)
	m.Add(ipp.TagOperation, attrs...)
	return m
}

var pdf = []byte("%PDF-1.4\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj 2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj 3 0 obj<</Type/Page/Parent 2 0 R/MediaBox[0 0 595 842]>>endobj\ntrailer<</Root 1 0 R>>\n%%EOF\n")

var steps = []step{
	{name: "Get-Printer-Attributes", attrs: []string{"printer-state", "document-format-supported", "sides-supported", "media-supported", "copies-supported", "operations-supported", "media-col-database"},
		req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
			return op(c, ipp.OpGetPrinterAttributes, ipp.Attr("requested-attributes", ipp.Keyword("printer-state"), ipp.Keyword("document-format-supported"),
				ipp.Keyword("sides-supported"), ipp.Keyword("media-supported"), ipp.Keyword("copies-supported"), ipp.Keyword("operations-supported"),
				ipp.Keyword("media-col-database"))), nil
		}},
	{name: "Validate-Job two-sided", req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		m := op(c, ipp.OpValidateJob, ipp.Attr("document-format", ipp.MimeType("application/pdf")))
		m.Add(ipp.TagJob, ipp.Attr("sides", ipp.Keyword("two-sided-long-edge")))
		return m, nil
	}},
	// ippeveprinter says attributes-or-values-not-supported; RFC 8011 has
	// the more precise document-format-not-supported for it (§4.1.7)
	{name: "Validate-Job document-format image/gif", either: []ipp.Status{ipp.StatusClientErrorAttributesOrValuesNotSupported, ipp.StatusClientErrorDocumentFormatNotSupported}, req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpValidateJob, ipp.Attr("document-format", ipp.MimeType("image/gif"))), nil
	}},
	{name: "Validate-Job sides of a fidelity it lacks", req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		m := op(c, ipp.OpValidateJob, ipp.Attr("ipp-attribute-fidelity", ipp.Boolean(true)))
		m.Add(ipp.TagJob, ipp.Attr("sides", ipp.Keyword("two-sided-sideways")))
		return m, nil
	}},
	{name: "Print-Job", job: true, attrs: []string{"job-id", "job-uri", "job-state"}, req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpPrintJob, ipp.Attr("job-name", ipp.Name("diff")), ipp.Attr("document-format", ipp.MimeType("application/pdf"))), pdf
	}},
	{name: "Get-Job-Attributes", attrs: []string{"job-id", "job-state", "job-state-reasons"}, req: func(c *ipp.Client, ids []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpGetJobAttributes, ipp.Attr("job-id", ipp.Integer(int(ids[0])))), nil
	}},
	{name: "Get-Job-Attributes of no job", req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpGetJobAttributes, ipp.Attr("job-id", ipp.Integer(9999))), nil
	}},
	{name: "Create-Job", job: true, settle: true, attrs: []string{"job-id", "job-state"}, req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpCreateJob, ipp.Attr("job-name", ipp.Name("diff 2"))), nil
	}},
	{name: "Send-Document without last-document", req: func(c *ipp.Client, ids []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpSendDocument, ipp.Attr("job-id", ipp.Integer(int(ids[1]))), ipp.Attr("document-format", ipp.MimeType("application/pdf"))), pdf
	}},
	{name: "Send-Document", attrs: []string{"job-id", "job-state"}, req: func(c *ipp.Client, ids []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpSendDocument, ipp.Attr("job-id", ipp.Integer(int(ids[1]))), ipp.Attr("document-format", ipp.MimeType("application/pdf")),
			ipp.Attr("last-document", ipp.Boolean(true))), pdf
	}},
	// (a printer without jobs of more documents may answer that first:
	// the ippeveprinter of Ubuntu does, of macOS (CUPS 2.3) not)
	{name: "Send-Document after the last", either: []ipp.Status{ipp.StatusClientErrorNotPossible, ipp.StatusServerErrorMultipleDocumentJobsNotSupported}, req: func(c *ipp.Client, ids []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpSendDocument, ipp.Attr("job-id", ipp.Integer(int(ids[1]))), ipp.Attr("document-format", ipp.MimeType("application/pdf")),
			ipp.Attr("last-document", ipp.Boolean(true))), pdf
	}},
	{name: "Cancel-Job of no job", req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpCancelJob, ipp.Attr("job-id", ipp.Integer(9999))), nil
	}},
	{name: "Get-Jobs completed", req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpGetJobs, ipp.Attr("which-jobs", ipp.Keyword("completed"))), nil
	}},
	{name: "Get-Jobs all", req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpGetJobs, ipp.Attr("which-jobs", ipp.Keyword("all")), ipp.Attr("requested-attributes", ipp.Keyword("job-id"), ipp.Keyword("job-state"))), nil
	}},
	{name: "Identify-Printer", req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.OpIdentifyPrinter), nil
	}},
	{name: "an operation not supported", req: func(c *ipp.Client, _ []int32) (*ipp.Message, []byte) {
		return op(c, ipp.Op(0x3fff)), nil
	}},
}

// TestAgainstIPPEve sends the same requests to ippeveprinter and to a
// fake of the same printer and compares the answers: the status of each,
// and the attributes both must tell.
func TestAgainstIPPEve(t *testing.T) {
	if testing.Short() {
		t.Skip("ippeveprinter takes seconds a job")
	}
	eve := &ipp.Client{URL: startIPPEve(t), User: "tester"}
	ps, _ := fakeprinter.LookupPreset("office-mfp")
	p := ps.New("fake")
	p.Formats = []string{"application/pdf", "image/pwg-raster", "image/jpeg", "application/octet-stream"}
	fake := &ipp.Client{URL: fakeprinter.Start(t, p), User: "tester"}
	ctx := context.Background()
	var eveIDs, fakeIDs []int32
	settle := func(c *ipp.Client, ids []int32) {
		for _, id := range ids {
			for range 300 {
				js, err := c.GetJobAttributes(ctx, id)
				if err != nil || js.State.Done() {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	for _, s := range steps {
		if s.settle {
			settle(eve, eveIDs)
			settle(fake, fakeIDs)
		}
		if len(eveIDs) != len(fakeIDs) {
			t.Fatalf("%s: the jobs differ: ippeveprinter %v, fake %v", s.name, eveIDs, fakeIDs)
		}
		run := func(c *ipp.Client, ids *[]int32) (*ipp.Message, error) {
			m, doc := s.req(c, *ids)
			var body *bytes.Reader
			if doc != nil {
				body = bytes.NewReader(doc)
			}
			var res *ipp.Message
			var err error
			if body != nil {
				res, err = c.Do(ctx, m, body)
			} else {
				res, err = c.Do(ctx, m, nil)
			}
			if res != nil && s.job {
				if g := res.Group(ipp.TagJob); g != nil {
					id, _ := g.Get("job-id").Int()
					*ids = append(*ids, id)
				}
			}
			return res, err
		}
		er, eerr := run(eve, &eveIDs)
		fr, ferr := run(fake, &fakeIDs)
		if er == nil || fr == nil {
			t.Errorf("%s: ippeveprinter %v, fake %v", s.name, eerr, ferr)
			continue
		}
		if er.Status() != fr.Status() && !(slices.Contains(s.either, er.Status()) && slices.Contains(s.either, fr.Status())) {
			t.Errorf("%s: ippeveprinter %s, fake %s (%v)", s.name, er.Status(), fr.Status(), ferr)
		}
		for _, a := range s.attrs {
			if has(er, a) && !has(fr, a) {
				t.Errorf("%s: the fake does not tell %s", s.name, a)
			}
		}
		// the jobs complete before the next step looks at them
		time.Sleep(50 * time.Millisecond)
	}
}

func has(m *ipp.Message, name string) bool {
	for _, g := range m.Groups {
		if g.Get(name) != nil {
			return true
		}
	}
	return false
}
