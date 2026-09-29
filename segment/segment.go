// Package segment serves documents to readers a few pages at a time, each
// segment sealed for a key pair the reader makes for that one request
// (docs/spec.md §3.6, docs/design.md §3.30). A reader that forgets its
// private key after opening the segment leaves nothing that opens it: what
// was sent stays sealed even if the server's keys, the TLS keys or the
// reader's session are taken later.
//
// The request is a POST of a JSON Request; the answer is the single-file
// form of the segment of Handler.Pages pages that holds the page asked for,
// sealed for the request's key (bdf.NewECDHLock). @bdfkit/core's
// SegmentLoader is the reader's side.
package segment

import (
	"bytes"
	"crypto/ecdh"
	"encoding/json"
	"errors"
	"io/fs"
	"mime"
	"net/http"
	"strconv"

	"github.com/shibukawa/bdf"
)

// Pages is the number of pages of a segment unless Handler.Pages says.
const Pages = 10

// maxRequest bounds the size of a request body.
const maxRequest = 64 << 10

// Request is the JSON body of a request for a segment.
type Request struct {
	// Key is the reader's public key for this request: an uncompressed
	// P-256 point (65 bytes, base64 in JSON).
	Key []byte `json:"key"`
	// View and Page name a page (0-based) the segment is to hold; an empty
	// View is the document's first view.
	View string `json:"view,omitempty"`
	Page int    `json:"page"`
	// Have lists the pages the reader holds from earlier segments, whose
	// parts are left out.
	Have []bdf.Segment `json:"have,omitempty"`
}

// Handler answers requests for segments.
type Handler struct {
	// Open returns the document a request is for, once it has checked that
	// the reader may read it (the application's own login and rights). An
	// error that is a *StatusError answers with its status, fs.ErrNotExist
	// 404, fs.ErrPermission 403, others 500.
	Open func(r *http.Request) (*bdf.Reader, error)
	// Allow, if set, is asked before a segment is sent, for rules on pages
	// (a sample of the first pages, how fast pages are read). Its error
	// answers as one of Open does.
	Allow func(r *http.Request, s bdf.Segment) error
	// Pages is the number of pages of a segment (0: Pages).
	Pages int
	// Log, if set, is given the errors answered with 500.
	Log func(r *http.Request, err error)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "segments are asked for with POST", http.StatusMethodNotAllowed)
		return
	}
	// JSON cannot be sent across sites without asking (CORS), unlike a form
	if mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type")); mt != "application/json" {
		http.Error(w, "the request is JSON", http.StatusUnsupportedMediaType)
		return
	}
	var req Request
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequest)).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	key, err := ecdh.P256().NewPublicKey(req.Key)
	if err != nil {
		http.Error(w, "bad request: the key is not a P-256 public key", http.StatusBadRequest)
		return
	}
	doc, err := h.Open(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	s, err := h.segment(doc, &req)
	if err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if h.Allow != nil {
		if err := h.Allow(r, s); err != nil {
			h.fail(w, r, err)
			return
		}
	}
	lock, err := bdf.NewECDHLock(key)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var b bytes.Buffer
	if err := doc.WriteSegment(&b, bdf.SegmentOptions{Segment: s, Have: req.Have, Lock: lock}); err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	w.Header().Set("Cache-Control", "no-store")
	w.Write(b.Bytes())
}

// segment returns the segment that holds the page asked for, and checks the
// pages the reader says it holds.
func (h *Handler) segment(doc *bdf.Reader, req *Request) (bdf.Segment, error) {
	views := doc.Manifest.Views
	view := func(id string) *bdf.View {
		for _, v := range views {
			if v.ID == id {
				return v
			}
		}
		return nil
	}
	if len(views) == 0 {
		return bdf.Segment{}, errors.New("the document has no views")
	}
	v := views[0]
	if req.View != "" {
		if v = view(req.View); v == nil {
			return bdf.Segment{}, errors.New("no view " + strconv.Quote(req.View))
		}
	}
	if req.Page < 0 || req.Page >= len(v.Pages) {
		return bdf.Segment{}, errors.New("no page " + strconv.Itoa(req.Page) + " in view " + strconv.Quote(v.ID))
	}
	for _, s := range req.Have {
		hv := view(s.View)
		if hv == nil || s.From < 0 || s.From >= s.To || s.To > len(hv.Pages) {
			return bdf.Segment{}, errors.New("have: no pages " + strconv.Itoa(s.From) + " to " + strconv.Itoa(s.To) + " in view " + strconv.Quote(s.View))
		}
	}
	n := h.Pages
	if n <= 0 {
		n = Pages
	}
	return bdf.SegmentAt(v, req.Page, n), nil
}

// StatusError is an error of Handler.Open or Handler.Allow that answers
// with an HTTP status of its own: 401 for a reader who is not logged in,
// 429 for one who reads too fast.
type StatusError struct {
	Status int
	Err    error
}

func (e *StatusError) Error() string { return e.Err.Error() }
func (e *StatusError) Unwrap() error { return e.Err }

func (h *Handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	var se *StatusError
	switch {
	case errors.As(err, &se):
		http.Error(w, http.StatusText(se.Status), se.Status)
	case errors.Is(err, fs.ErrNotExist):
		http.Error(w, "no such document", http.StatusNotFound)
	case errors.Is(err, fs.ErrPermission):
		http.Error(w, "not allowed", http.StatusForbidden)
	default:
		if h.Log != nil {
			h.Log(r, err)
		}
		http.Error(w, "the segment could not be made", http.StatusInternalServerError)
	}
}
