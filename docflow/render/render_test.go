package render_test

import (
	"bytes"
	"context"
	"errors"
	"image/png"
	"testing"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	_ "github.com/shibukawa/bdf/docflow/io/http"
	"github.com/shibukawa/bdf/docflow/render"
)

// An endpoint's .bdf drawn as a PDF and as a PNG of its second page, and
// sent on to a web service.
func TestRender(t *testing.T) {
	e := flowtest.New(t, nil)
	hook := flowtest.NewHook(t)
	e.Put(docflow.Node{Name: "out", Kind: "http", URL: hook.URL + "/files"},
		docflow.Node{Name: "pdf", Kind: render.KindRender, To: []string{"out"}},
		docflow.Node{Name: "png", Kind: render.KindRender, Format: "png", Page: 2, To: []string{"out"}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"pdf", "png"}})
	j := e.Submit("in", "invoice.pdf", flowtest.Doc(t, 2), nil)
	e.Tick()
	got := hook.Got()
	if len(got) != 2 {
		t.Fatalf("%d requests: %+v", len(got), e.Job(j.Number).Legs)
	}
	var pdf, img []byte
	for _, r := range got {
		switch r.Header.Get("Content-Type") {
		case "application/pdf":
			pdf = r.Body
			if cd := r.Header.Get("Content-Disposition"); cd != `attachment; filename="invoice.pdf"` {
				t.Errorf("pdf: %s", cd)
			}
		case "image/png":
			img = r.Body
		}
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatalf("no PDF: %.20q", pdf)
	}
	im, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		t.Fatal(err)
	}
	// 595 × 842 pt at 150 dpi, rounded up
	if b := im.Bounds(); b.Dx() != 1240 || b.Dy() != 1755 {
		t.Errorf("png %v", b)
	}
	job := e.Job(j.Number)
	if job.State != docflow.JobCompleted {
		t.Errorf("job %s: %+v", job.State, job.Legs)
	}
	var made []string
	for _, p := range job.Payloads[1:] {
		made = append(made, p.Format+":"+p.Kind+":"+p.Made)
	}
	if len(made) != 2 || made[0] != "pdf:any:pdf" && made[1] != "pdf:any:pdf" {
		t.Errorf("payloads %v", made)
	}
}

// What a render node takes and makes: .bdf in, any out; svg not yet, a
// page out of the document fails for good.
func TestRenderSettings(t *testing.T) {
	e := flowtest.New(t, nil)
	hook := flowtest.NewHook(t)
	e.Put(docflow.Node{Name: "out", Kind: "http", URL: hook.URL})
	ctx := context.Background()
	for _, n := range []docflow.Node{
		{Name: "r", Kind: render.KindRender, Format: "svg", To: []string{"out"}},
		{Name: "r", Kind: render.KindRender, Format: "tiff", To: []string{"out"}},
	} {
		if _, err := e.M().PutNode(ctx, n); !errors.Is(err, docflow.ErrInvalid) {
			t.Errorf("%+v: %v", n, err)
		}
	}
	// data cannot reach a render node
	e.Put(docflow.Node{Name: "r", Kind: render.KindRender, Format: "png", Page: 3, To: []string{"out"}})
	if _, err := e.M().PutNode(ctx, docflow.Node{Name: "data", Kind: docflow.KindEndpoint, Accept: docflow.AcceptData, To: []string{"r"}}); !errors.Is(err, docflow.ErrInvalid) {
		t.Errorf("data to render: %v", err)
	}
	e.Put(docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"r"}})
	j := e.Submit("in", "one", flowtest.Doc(t, 1), nil)
	e.Tick()
	if d := e.Delivery(j, "r"); d.State != docflow.DeliveryFailed || d.Attempts != 1 {
		t.Errorf("page 3 of 1: %s after %d: %s", d.State, d.Attempts, d.LastError)
	}
	if len(hook.Got()) != 0 {
		t.Error("sent")
	}
}
