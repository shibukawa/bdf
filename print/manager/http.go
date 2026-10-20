package manager

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/image/barcode"
)

// SessionCookie is the cookie of a device's session.
const SessionCookie = "bdf_pm_device"

type handler struct {
	m *Manager
}

// Handler serves a manager over HTTP:
//
//	POST /ipp/print/{location}                 the location's IPP printer
//	/api/…                                     the API (bearer AdminToken)
//	/device/…                                  the devices' endpoints (their session)
//
// It is mounted on the web application's mux; it holds nothing between
// requests.
func Handler(m *Manager) http.Handler {
	h := &handler{m: m}
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+IPPPath+"/{location}", h.serveIPP)

	mux.HandleFunc("GET /api/locations", h.adminOnly(h.listLocations))
	mux.HandleFunc("GET /api/locations/{name}", h.adminOnly(h.getLocation))
	mux.HandleFunc("PUT /api/locations/{name}", h.adminOnly(h.putLocation))
	mux.HandleFunc("DELETE /api/locations/{name}", h.adminOnly(h.deleteLocation))
	mux.HandleFunc("GET /api/locations/{name}/board", h.adminOnly(h.board))
	mux.HandleFunc("GET /api/locations/{name}/devices", h.adminOnly(h.listDevices))
	mux.HandleFunc("POST /api/locations/{name}/enrollments", h.adminOnly(h.enrollment))
	mux.HandleFunc("POST /api/locations/{name}/slips", h.submit)
	mux.HandleFunc("DELETE /api/devices/{id}", h.adminOnly(h.removeDevice))
	mux.HandleFunc("GET /api/printers", h.adminOnly(h.listPrinters))
	mux.HandleFunc("GET /api/printers/{name}", h.adminOnly(h.getPrinter))
	mux.HandleFunc("PUT /api/printers/{name}", h.adminOnly(h.putPrinter))
	mux.HandleFunc("DELETE /api/printers/{name}", h.adminOnly(h.deletePrinter))
	mux.HandleFunc("GET /api/slips/{id}", h.adminOnly(h.getSlip))
	mux.HandleFunc("GET /api/slips/{id}/document", h.adminOnly(h.slipDocument))
	mux.HandleFunc("POST /api/slips/{id}/cancel", h.adminOnly(h.cancelSlip))
	mux.HandleFunc("POST /api/tick", h.tick)

	mux.HandleFunc("POST /device/enroll", h.enroll)
	mux.HandleFunc("GET /device/challenge", h.challenge)
	mux.HandleFunc("POST /device/session", h.session)
	mux.HandleFunc("GET /device/me", h.deviceOnly(h.me))
	mux.HandleFunc("GET /device/board", h.deviceOnly(h.deviceBoard))
	mux.HandleFunc("POST /device/slips/{id}/{action}", h.deviceOnly(h.deviceAction))
	mux.HandleFunc("GET /device/slips/{id}/pages/{page}", h.deviceOnly(h.devicePage))
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeErr answers an error with its status and {"error": …}.
func (h *handler) writeErr(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrInvalid):
		status = http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, ErrConflict), errors.Is(err, ErrLeased):
		status = http.StatusConflict
	case errors.Is(err, ErrFormat):
		status = http.StatusUnsupportedMediaType
	case errors.Is(err, ErrBranch):
		status = http.StatusNotImplemented
	}
	msg := err.Error()
	if status == http.StatusInternalServerError {
		h.m.o.Log.Error("request", "err", err)
		msg = "internal error"
	}
	writeJSON(w, status, map[string]string{"error": msg})
}

func readJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return errors.Join(ErrInvalid, err)
	}
	return nil
}

func bearer(r *http.Request) string {
	tok, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	return tok
}

func sameSecret(a, b string) bool {
	return a != "" && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (h *handler) admin(r *http.Request) bool { return sameSecret(h.m.o.AdminToken, bearer(r)) }

func (h *handler) adminOnly(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.admin(r) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="bdf print manager"`)
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "the admin token is required"})
			return
		}
		fn(w, r)
	}
}

func (h *handler) listLocations(w http.ResponseWriter, r *http.Request) {
	ls, err := h.m.Locations(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ls)
}

func (h *handler) getLocation(w http.ResponseWriter, r *http.Request) {
	l, err := h.m.Location(r.Context(), r.PathValue("name"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func (h *handler) putLocation(w http.ResponseWriter, r *http.Request) {
	var l Location
	if err := readJSON(r, &l); err != nil {
		h.writeErr(w, err)
		return
	}
	l.Name = r.PathValue("name")
	if err := h.m.PutLocation(r.Context(), l); err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func (h *handler) deleteLocation(w http.ResponseWriter, r *http.Request) {
	if err := h.m.DeleteLocation(r.Context(), r.PathValue("name")); err != nil {
		h.writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) board(w http.ResponseWriter, r *http.Request) {
	b, err := h.m.Board(r.Context(), r.PathValue("name"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *handler) listDevices(w http.ResponseWriter, r *http.Request) {
	ds, err := h.m.Devices(r.Context(), r.PathValue("name"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ds)
}

// enrollment makes a one-time code; ttl is a duration ("15m").
func (h *handler) enrollment(w http.ResponseWriter, r *http.Request) {
	var ttl time.Duration
	if s := r.URL.Query().Get("ttl"); s != "" {
		d, err := ParseDuration(s)
		if err != nil {
			h.writeErr(w, errors.Join(ErrInvalid, err))
			return
		}
		ttl = d.d()
	}
	e, err := h.m.NewEnrollment(r.Context(), r.PathValue("name"), ttl)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	u := scheme + "://" + r.Host + h.m.o.TabletPath + "?code=" + e.Code
	writeJSON(w, http.StatusCreated, struct {
		*Enrollment
		URL string `json:"url"`
		QR  string `json:"qr,omitempty"`
	}{e, u, qrSVG(u)})
}

// qrSVG draws a QR code of a text as SVG, for a device's camera.
func qrSVG(text string) string {
	sym, err := barcode.Encode(barcode.QR, text, nil)
	if err != nil {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" shape-rendering="crispEdges"><rect width="100%%" height="100%%" fill="#fff"/><path fill="#000" d="`, sym.W, sym.H)
	for _, r := range sym.Rects() {
		fmt.Fprintf(&b, "M%d %dh%dv%dh-%dz", r[0], r[1], r[2], r[3], r[2])
	}
	b.WriteString(`"/></svg>`)
	return b.String()
}

func (h *handler) removeDevice(w http.ResponseWriter, r *http.Request) {
	if err := h.m.RemoveDevice(r.Context(), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) listPrinters(w http.ResponseWriter, r *http.Request) {
	ps, err := h.m.Printers(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ps)
}

func (h *handler) getPrinter(w http.ResponseWriter, r *http.Request) {
	p, err := h.m.Printer(r.Context(), r.PathValue("name"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *handler) putPrinter(w http.ResponseWriter, r *http.Request) {
	var p Printer
	if err := readJSON(r, &p); err != nil {
		h.writeErr(w, err)
		return
	}
	p.Name = r.PathValue("name")
	if err := h.m.PutPrinter(r.Context(), p); err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *handler) deletePrinter(w http.ResponseWriter, r *http.Request) {
	if err := h.m.DeletePrinter(r.Context(), r.PathValue("name")); err != nil {
		h.writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) getSlip(w http.ResponseWriter, r *http.Request) {
	s, err := h.m.Slip(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *handler) slipDocument(w http.ResponseWriter, r *http.Request) {
	s, err := h.m.Slip(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/x-bdf")
	w.Header().Set("Content-Length", strconv.FormatInt(s.Size, 10))
	w.Header().Set("Content-Disposition", `attachment; filename="slip-`+s.ID+`.bdf"`)
	if err := h.m.WriteDocument(r.Context(), s, w); err != nil {
		panic(http.ErrAbortHandler)
	}
}

func (h *handler) cancelSlip(w http.ResponseWriter, r *http.Request) {
	s, err := h.m.Cancel(r.Context(), r.PathValue("id"), Mark{Worker: "admin"})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// submit takes a document (the body) for a location: ?name= names it.
func (h *handler) submit(w http.ResponseWriter, r *http.Request) {
	user, ok := h.submitter(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "a submit token is required"})
		return
	}
	f, err := os.CreateTemp(h.m.o.TempDir, "bdf-pm-in-*")
	if err != nil {
		h.writeErr(w, err)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	n, err := io.Copy(f, io.LimitReader(r.Body, h.m.o.MaxDocument+1))
	if err != nil {
		h.writeErr(w, errors.Join(ErrInvalid, err))
		return
	}
	if n > h.m.o.MaxDocument {
		writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "the document is too large"})
		return
	}
	q := r.URL.Query()
	s, err := h.m.Submit(r.Context(), r.PathValue("name"), Input{Name: q.Get("name"), User: firstOr(user, q.Get("user")), Source: "api", Doc: f, Size: n})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

// tick runs Tick: for the admin, or the scheduler's TickToken.
func (h *handler) tick(w http.ResponseWriter, r *http.Request) {
	if !h.admin(r) && !sameSecret(h.m.o.TickToken, bearer(r)) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "the tick token is required"})
		return
	}
	rep, err := h.m.Tick(r.Context())
	if err != nil {
		h.m.o.Log.Warn("tick", "err", err)
	}
	writeJSON(w, http.StatusOK, rep)
}

// The devices' endpoints.

func (h *handler) enroll(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Key  string `json:"key"`
	}
	if err := readJSON(r, &in); err != nil {
		h.writeErr(w, err)
		return
	}
	d, err := h.m.Enroll(r.Context(), in.Code, in.Name, in.Key)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"device": d.ID, "location": d.Location, "name": d.Name})
}

func (h *handler) challenge(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"challenge": h.m.Challenge()})
}

func (h *handler) session(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Device    string `json:"device"`
		Challenge string `json:"challenge"`
		Signature string `json:"signature"`
	}
	if err := readJSON(r, &in); err != nil {
		h.writeErr(w, err)
		return
	}
	tok, until, err := h.m.Authenticate(r.Context(), in.Device, in.Challenge, in.Signature)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: SessionCookie, Value: tok, Path: "/", MaxAge: int(h.m.o.SessionFor / time.Second), HttpOnly: true,
		SameSite: http.SameSiteStrictMode, Secure: h.m.o.SecureCookies})
	writeJSON(w, http.StatusOK, map[string]any{"until": until})
}

type deviceHandler func(http.ResponseWriter, *http.Request, *Device)

func (h *handler) deviceOnly(fn deviceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SessionCookie)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "no session"})
			return
		}
		d, err := h.m.Session(r.Context(), c.Value)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "the session ended"})
				return
			}
			h.writeErr(w, err)
			return
		}
		fn(w, r, d)
	}
}

func (h *handler) me(w http.ResponseWriter, r *http.Request, d *Device) {
	l, err := h.m.Location(r.Context(), d.Location)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"device": d.ID, "name": d.Name, "location": d.Location,
		"title": or(l.Title, l.Name), "can_print": l.Fallback != nil})
}

// deviceBoard answers the board, or 304 while its revision (the ETag) is
// the one the device has.
func (h *handler) deviceBoard(w http.ResponseWriter, r *http.Request, d *Device) {
	rev, err := h.m.BoardRev(r.Context(), d.Location)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	etag := `"` + rev + `"`
	if r.Header.Get("If-None-Match") == etag {
		w.Header().Set("ETag", etag)
		w.WriteHeader(http.StatusNotModified)
		return
	}
	b, err := h.m.Board(r.Context(), d.Location)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	w.Header().Set("ETag", `"`+b.Rev+`"`)
	writeJSON(w, http.StatusOK, b)
}

// deviceSlip is a slip of the device's location.
func (h *handler) deviceSlip(w http.ResponseWriter, r *http.Request, d *Device) *Slip {
	s, err := h.m.Slip(r.Context(), r.PathValue("id"))
	if err == nil && s.Location != d.Location {
		err = ErrNotFound
	}
	if err != nil {
		h.writeErr(w, err)
		return nil
	}
	return s
}

// deviceAction is open, done or print (on the location's printer).
func (h *handler) deviceAction(w http.ResponseWriter, r *http.Request, d *Device) {
	s := h.deviceSlip(w, r, d)
	if s == nil {
		return
	}
	var in struct {
		Worker string `json:"worker"`
	}
	if r.ContentLength != 0 {
		if err := readJSON(r, &in); err != nil {
			h.writeErr(w, err)
			return
		}
	}
	by := Mark{Device: d.ID, Worker: clip(strings.TrimSpace(in.Worker), 80)}
	var err error
	switch r.PathValue("action") {
	case "open":
		s, err = h.m.Open(r.Context(), s.ID, by)
	case "done":
		s, err = h.m.Done(r.Context(), s.ID, by)
	case "print":
		l, lerr := h.m.Location(r.Context(), d.Location)
		if lerr != nil {
			h.writeErr(w, lerr)
			return
		}
		if l.Fallback == nil {
			h.writeErr(w, errors.Join(ErrInvalid, errors.New("the location has no printer")))
			return
		}
		s, err = h.m.PrintSlip(r.Context(), s.ID, l.Fallback.Printer, by)
	default:
		h.writeErr(w, ErrNotFound)
		return
	}
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// devicePage draws a page of a slip: /device/slips/{id}/pages/{n}.png?w=.
func (h *handler) devicePage(w http.ResponseWriter, r *http.Request, d *Device) {
	s := h.deviceSlip(w, r, d)
	if s == nil {
		return
	}
	n, err := strconv.Atoi(strings.TrimSuffix(r.PathValue("page"), ".png"))
	if err != nil {
		h.writeErr(w, ErrNotFound)
		return
	}
	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	if width == 0 {
		width = 1200
	}
	w.Header().Set("Content-Type", "image/png")
	// the page does not change, but it is only the device's to see
	w.Header().Set("Cache-Control", "private, max-age=3600")
	if err := h.m.WritePage(r.Context(), s, n-1, width, w); err != nil {
		h.writeErr(w, err)
	}
}
