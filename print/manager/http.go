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
//	POST /ipp/print/{node}                     the node's IPP printer
//	POST /ipp/print                            the server (CUPS-Get-Printers)
//	/api/…                                     the API (bearer AdminToken, or a session of the pages)
//	/device/…                                  the devices' endpoints (their session)
//
// It is mounted on the web application's mux; it holds nothing between
// requests.
func Handler(m *Manager) http.Handler {
	h := &handler{m: m}
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+IPPPath+"/{dest}", h.serveIPP)
	mux.HandleFunc("POST "+IPPPath, h.serveIPP)
	mux.HandleFunc("POST "+IPPPath+"/{$}", h.serveIPP)

	// the nodes, of every kind; a display's board, devices and enrollment
	// codes, and a template node's templates, are under their node
	mux.HandleFunc("GET /api/nodes", h.adminOnly(h.listNodes))
	mux.HandleFunc("GET /api/nodes/{name}", h.adminOnly(h.getNode))
	mux.HandleFunc("PUT /api/nodes/{name}", h.adminOnly(h.putNode))
	mux.HandleFunc("DELETE /api/nodes/{name}", h.adminOnly(h.deleteNode))
	mux.HandleFunc("POST /api/nodes/{name}/key", h.adminOnly(h.rotateKey))
	mux.HandleFunc("GET /api/nodes/{name}/board", h.adminOnly(h.board))
	mux.HandleFunc("GET /api/nodes/{name}/devices", h.adminOnly(h.listDevices))
	mux.HandleFunc("POST /api/nodes/{name}/enrollments", h.adminOnly(h.enrollment))
	mux.HandleFunc("GET /api/nodes/{name}/templates", h.adminOnly(h.listTemplates))
	mux.HandleFunc("GET /api/nodes/{name}/templates/{id}", h.adminOnly(h.getTemplate))
	mux.HandleFunc("GET /api/nodes/{name}/templates/{id}/file", h.adminOnly(h.templateFile))
	mux.HandleFunc("PUT /api/nodes/{name}/templates/{id}", h.adminOnly(h.putTemplate))
	mux.HandleFunc("DELETE /api/nodes/{name}/templates/{id}", h.adminOnly(h.deleteTemplate))
	mux.HandleFunc("GET /api/formats", h.adminOnly(h.formats))
	// the apps and their tokens (apps.go)
	mux.HandleFunc("GET /api/apps", h.adminOnly(h.listApps))
	mux.HandleFunc("GET /api/apps/{id}", h.adminOnly(h.getApp))
	mux.HandleFunc("PUT /api/apps/{id}", h.adminOnly(h.putApp))
	mux.HandleFunc("DELETE /api/apps/{id}", h.adminOnly(h.deleteApp))
	mux.HandleFunc("POST /api/apps/{id}/secret", h.adminOnly(h.rotateAppSecret))
	mux.HandleFunc("POST "+OAuthTokenPath, h.oauthToken)
	mux.HandleFunc("GET /api/nodes/{name}/latest", h.userOnly(h.latest))
	mux.HandleFunc("GET /api/nodes/{name}/preview", h.adminOnly(h.preview))
	mux.HandleFunc("GET /api/nodes/{name}/preview/{job}/pages/{page}", h.adminOnly(h.previewPage))
	mux.HandleFunc("POST /api/locations/{name}/slips", h.submit)
	mux.HandleFunc("DELETE /api/devices/{id}", h.adminOnly(h.removeDevice))
	mux.HandleFunc("GET /api/me", h.whoami)
	mux.HandleFunc("GET /api/my/jobs", h.userOnly(h.myJobs))
	mux.HandleFunc("GET /api/destinations", h.userOnly(h.destinations))
	mux.HandleFunc("GET /api/audit", h.auditOnly(h.auditLog))
	mux.HandleFunc("GET /api/slips/{id}", h.adminOnly(h.getSlip))
	mux.HandleFunc("GET /api/slips/{id}/document", h.adminOnly(h.slipDocument))
	mux.HandleFunc("POST /api/slips/{id}/cancel", h.adminOnly(h.cancelSlip))
	// the topology is everyone's to see (an administrator sees its
	// settings too), as is where it is drawn; its activity is the viewer's
	mux.HandleFunc("GET /api/topology", h.userOnly(h.topology))
	mux.HandleFunc("GET /api/topology/live", h.userOnly(h.live))
	mux.HandleFunc("GET /api/layout", h.userOnly(func(w http.ResponseWriter, r *http.Request, _ *Viewer) { h.getLayout(w, r) }))
	mux.HandleFunc("PUT /api/layout", h.adminOnly(h.putLayout))
	mux.HandleFunc("POST "+InPath+"/{name}", h.serveIn)
	mux.HandleFunc("GET /api/destinations/{name}/jobs", h.adminOnly(h.destJobs))
	mux.HandleFunc("POST /api/destinations/{name}/jobs", h.submit)
	mux.HandleFunc("GET /api/jobs/{id}", h.ownJob(h.getJob))
	mux.HandleFunc("GET /api/jobs/{id}/document", h.ownJob(h.jobDocument))
	mux.HandleFunc("GET /api/jobs/{id}/pages/{page}", h.ownJob(h.jobPage))
	mux.HandleFunc("POST /api/jobs/{id}/cancel", h.ownJob(h.cancelJob))
	mux.HandleFunc("POST /api/jobs/{id}/reprint", h.ownJob(h.reprint))
	mux.HandleFunc("GET /api/inboxes/{name}", h.adminOnly(h.inbox))
	mux.HandleFunc("POST /api/items/{id}/release", h.adminOnly(h.release))
	mux.HandleFunc("POST /api/items/{id}/discard", h.adminOnly(h.discardItem))
	mux.HandleFunc("GET /api/deliveries", h.adminOnly(h.pending))
	mux.HandleFunc("POST /api/deliveries/{id}/retry", h.adminOnly(h.retry))
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
	case errors.Is(err, ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, ErrConflict), errors.Is(err, ErrLeased):
		status = http.StatusConflict
	case errors.Is(err, ErrFormat):
		status = http.StatusUnsupportedMediaType
	case errors.Is(err, ErrBranch):
		status = http.StatusNotImplemented
	case errors.Is(err, ErrPrint):
		status = http.StatusBadGateway
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

// adminOnly serves a request of whoever changes the topology: the admin
// token, an app with the scope admin, or a person with the role topology
// (an administrator has it).
func (h *handler) adminOnly(fn http.HandlerFunc) http.HandlerFunc {
	return h.roleOnly(RoleTopology, fn)
}

// auditOnly serves a request of whoever reads the audit log: the role
// audit, or an administrator.
func (h *handler) auditOnly(fn http.HandlerFunc) http.HandlerFunc {
	return h.roleOnly(RoleAudit, fn)
}

// roleOnly serves a request of a viewer with a role.
func (h *handler) roleOnly(role string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, ok := h.viewerOf(r)
		switch {
		case !ok:
			w.Header().Set("WWW-Authenticate", `Bearer realm="bdf print manager"`)
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "the admin token or a sign-in is required"})
		case !v.Can(role):
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "the role " + role + " is required"})
		case !sameOrigin(r):
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "a request of another origin"})
		default:
			fn(w, r)
		}
	}
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
	h.m.audit(r.Context(), h.who(r), "enrollment-code", r.PathValue("name"), "")
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
	h.m.audit(r.Context(), h.who(r), "remove-device", r.PathValue("id"), "")
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
	if err := h.m.WritePayload(r.Context(), firstOr(s.Payload, jobID(s.Job)), w); err != nil {
		panic(http.ErrAbortHandler)
	}
}

func (h *handler) cancelSlip(w http.ResponseWriter, r *http.Request) {
	s, err := h.m.CancelSlip(r.Context(), r.PathValue("id"), Mark{Worker: "admin"})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// The topology: nodes, jobs, inboxes and deliveries. The hash of an
// endpoint's key stays in the store: the API shows nodes without it.

func (h *handler) listNodes(w http.ResponseWriter, r *http.Request) {
	ns, err := h.m.Nodes(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	for i := range ns {
		ns[i].KeyHash = ""
	}
	writeJSON(w, http.StatusOK, ns)
}

func (h *handler) getNode(w http.ResponseWriter, r *http.Request) {
	n, err := h.m.Node(r.Context(), r.PathValue("name"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	n.KeyHash = ""
	writeJSON(w, http.StatusOK, n)
}

// putNode puts a node; a new endpoint's key is in the answer, this once.
func (h *handler) putNode(w http.ResponseWriter, r *http.Request) {
	var n Node
	if err := readJSON(r, &n); err != nil {
		h.writeErr(w, err)
		return
	}
	n.Name = r.PathValue("name")
	key, err := h.m.PutNode(r.Context(), n)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	got, err := h.m.Node(r.Context(), n.Name)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	got.KeyHash = ""
	detail := got.Kind
	if key != "" {
		detail += ", a key made"
	}
	h.m.audit(r.Context(), h.who(r), "put-node", n.Name, detail)
	writeJSON(w, http.StatusOK, nodeWithKey{got, key})
}

// nodeWithKey is a node and, once, an endpoint's key.
type nodeWithKey struct {
	*Node
	Key string `json:"key,omitempty"`
}

func (h *handler) rotateKey(w http.ResponseWriter, r *http.Request) {
	key, err := h.m.RotateKey(r.Context(), r.PathValue("name"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), h.who(r), "rotate-key", r.PathValue("name"), "")
	writeJSON(w, http.StatusOK, map[string]string{"name": r.PathValue("name"), "key": key})
}

func (h *handler) deleteNode(w http.ResponseWriter, r *http.Request) {
	if err := h.m.DeleteNode(r.Context(), r.PathValue("name")); err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), h.who(r), "delete-node", r.PathValue("name"), "")
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) topology(w http.ResponseWriter, r *http.Request, v *Viewer) {
	ns, err := h.m.Topology(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	if ns == nil {
		ns = []TopoNode{}
	}
	if !v.Can(RoleTopology) {
		ns = publicTopology(ns)
	}
	writeJSON(w, http.StatusOK, ns)
}

// live answers the topology's activity the viewer sees, or 304 while its
// revision (the ETag) is the one the client has.
func (h *handler) live(w http.ResponseWriter, r *http.Request, v *Viewer) {
	l, err := h.m.Live(r.Context(), v)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	etag := `"` + l.Rev + `"`
	w.Header().Set("ETag", etag)
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	writeJSON(w, http.StatusOK, l)
}

// destJobs lists the last jobs sent to a destination: ?limit= bounds them.
func (h *handler) destJobs(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	js, err := h.m.JobsOf(r.Context(), r.PathValue("name"), limit)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, js)
}

func (h *handler) getJob(w http.ResponseWriter, r *http.Request, v *Viewer, j *Job) {
	writeJSON(w, http.StatusOK, j)
}

// payloadOf is the payload of a job a request names (?payload=), or the
// .bdf document to show (the last one made, else the one received when
// it is one), or, with bdf unset, the one received.
func payloadOf(r *http.Request, j *Job, bdfOnly bool) (*Payload, error) {
	if id := r.URL.Query().Get("payload"); id != "" {
		p := j.payload(id)
		if p == nil {
			return nil, fmt.Errorf("%w: job %s has no payload %q", ErrNotFound, j.ID, id)
		}
		return p, nil
	}
	p := j.document()
	if p == nil && !bdfOnly {
		p = j.received()
	}
	if p == nil {
		return nil, fmt.Errorf("%w: job %s has no .bdf document", ErrNotFound, j.ID)
	}
	return p, nil
}

// jobPage draws a page of a job the viewer sees:
// /api/jobs/{id}/pages/{n}.png?w= (of the .bdf document to show, or
// ?payload= one).
func (h *handler) jobPage(w http.ResponseWriter, r *http.Request, v *Viewer, j *Job) {
	n, err := strconv.Atoi(strings.TrimSuffix(r.PathValue("page"), ".png"))
	if err != nil {
		h.writeErr(w, ErrNotFound)
		return
	}
	p, err := payloadOf(r, j, true)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	if p.Kind != PayloadBDF {
		h.writeErr(w, fmt.Errorf("%w: the payload %s is %s, not a .bdf document", ErrNotFound, p.ID, p.Kind))
		return
	}
	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	if width == 0 {
		width = 800
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "private, max-age=3600")
	if err := h.m.WritePage(r.Context(), p.ID, p.Pages, n-1, width, w); err != nil {
		h.writeErr(w, err)
	}
}

// latest lists the jobs a sink keeps that the viewer sees, newest first.
func (h *handler) latest(w http.ResponseWriter, r *http.Request, v *Viewer) {
	js, err := h.m.Latest(r.Context(), r.PathValue("name"), v)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, js)
}

// jobDocument downloads a payload of a job: the .bdf document to show,
// else what was received (data, a file), or ?payload= one.
func (h *handler) jobDocument(w http.ResponseWriter, r *http.Request, v *Viewer, j *Job) {
	p, err := payloadOf(r, j, false)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	name := p.Name
	if name == "" || strings.ContainsAny(name, `"\/`) {
		name = "job-" + j.ID + extensionOf(p)
	}
	w.Header().Set("Content-Type", contentTypeOf(p.Format, p.Name))
	w.Header().Set("Content-Length", strconv.FormatInt(p.Size, 10))
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	if err := h.m.WritePayload(r.Context(), p.ID, w); err != nil {
		panic(http.ErrAbortHandler)
	}
}

// reprint sends a job's document again, as a new job of the viewer, to
// a node they may send to: {"to": "<node>"}, or the job's destination.
func (h *handler) reprint(w http.ResponseWriter, r *http.Request, v *Viewer, j *Job) {
	var in struct {
		To string `json:"to"`
	}
	if r.ContentLength != 0 {
		if err := readJSON(r, &in); err != nil {
			h.writeErr(w, err)
			return
		}
	}
	dest := firstOr(in.To, j.Dest)
	if v.Kind == PrincipalUser {
		scope, owner, _, err := h.m.destScope(r.Context(), dest)
		if err != nil {
			h.writeErr(w, err)
			return
		}
		if !v.CanSendTo(scope, owner) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "you may not send to " + dest})
			return
		}
	} else if err := h.m.appMaySend(r.Context(), v, dest, false); err != nil {
		h.writeErr(w, err)
		return
	}
	again, err := h.m.Reprint(r.Context(), j, dest, v.Principal())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), v.Label(), "reprint", j.ID, "→ "+dest+" as "+again.ID)
	writeJSON(w, http.StatusCreated, again)
}

// cancelJob cancels a job: the submitter's or an administrator's (the
// members of its group see it only).
func (h *handler) cancelJob(w http.ResponseWriter, r *http.Request, v *Viewer, j *Job) {
	if !v.CanCancel(j) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "only the sender or an administrator cancels a job"})
		return
	}
	j, err := h.m.Cancel(r.Context(), j.ID, Mark{Worker: clip(v.Label(), 80)})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), v.Label(), "cancel-job", j.ID, j.Name)
	writeJSON(w, http.StatusOK, j)
}

func (h *handler) inbox(w http.ResponseWriter, r *http.Request) {
	v, err := h.m.Inbox(r.Context(), r.PathValue("name"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

// release sends an item on: {"to": [nodes], "worker": name}.
func (h *handler) release(w http.ResponseWriter, r *http.Request) {
	var in struct {
		To     []string `json:"to"`
		Worker string   `json:"worker"`
	}
	if r.ContentLength != 0 {
		if err := readJSON(r, &in); err != nil {
			h.writeErr(w, err)
			return
		}
	}
	it, err := h.m.Release(r.Context(), r.PathValue("id"), in.To, Mark{Worker: clip(strings.TrimSpace(firstOr(in.Worker, h.who(r))), 80)})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), h.who(r), "release-item", it.ID, strings.Join(it.Released.To, ", ")+firstOr(" by "+in.Worker, ""))
	writeJSON(w, http.StatusOK, it)
}

func (h *handler) discardItem(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Worker string `json:"worker"`
	}
	if r.ContentLength != 0 {
		if err := readJSON(r, &in); err != nil {
			h.writeErr(w, err)
			return
		}
	}
	it, err := h.m.Discard(r.Context(), r.PathValue("id"), Mark{Worker: clip(strings.TrimSpace(firstOr(in.Worker, h.who(r))), 80)})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), h.who(r), "discard-item", it.ID, in.Worker)
	writeJSON(w, http.StatusOK, it)
}

func (h *handler) pending(w http.ResponseWriter, r *http.Request) {
	ds, err := h.m.Pending(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ds)
}

func (h *handler) retry(w http.ResponseWriter, r *http.Request) {
	d, err := h.m.Retry(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

// submit takes a document (the body) for a destination: ?name= names it.
// At /api/locations/{name}/slips the answer is the slip of the display;
// at /api/destinations/{name}/jobs, the job.
func (h *handler) submit(w http.ResponseWriter, r *http.Request) {
	by, ok := h.submitter(r)
	if v, signedIn := h.viewerOf(r); signedIn && v.Kind == PrincipalApp {
		// an app sends to the endpoints that list it
		if err := h.m.appMaySend(r.Context(), v, r.PathValue("name"), false); err != nil {
			h.writeErr(w, err)
			return
		}
	} else if signedIn && v.Kind == PrincipalUser {
		// a person signed in sends as themselves, to a node they may
		scope, owner, _, err := h.m.destScope(r.Context(), r.PathValue("name"))
		if err != nil {
			h.writeErr(w, err)
			return
		}
		if !v.CanSendTo(scope, owner) || !sameOrigin(r) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "you may not send to " + r.PathValue("name")})
			return
		}
		by, ok = v.Principal(), true
	}
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
	dest := r.PathValue("name")
	j, err := h.m.Submit(r.Context(), dest, Input{Name: q.Get("name"), By: by, ClaimedUser: q.Get("user"), Source: "api", Doc: f, Size: n,
		ContentType: r.Header.Get("Content-Type"), Meta: metaOfRequest(r)})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/locations/") {
		for _, l := range j.Legs {
			if l.Kind == KindDisplay && l.Node == dest && l.Ref != "" {
				if s, err := h.m.Slip(r.Context(), l.Ref); err == nil {
					writeJSON(w, http.StatusCreated, s)
					return
				}
			}
		}
	}
	writeJSON(w, http.StatusCreated, j)
}

// tick runs Tick: for the admin, the scheduler's TickToken, or an app
// with the scope tick.
func (h *handler) tick(w http.ResponseWriter, r *http.Request) {
	if sameSecret(h.m.o.TickToken, bearer(r)) {
		h.runTick(w, r)
		return
	}
	if v, ok := h.viewerOf(r); ok && v.Kind == PrincipalApp && v.Has(ScopeTick) {
		h.runTick(w, r)
		return
	}
	h.adminOnly(h.runTick)(w, r)
}

func (h *handler) runTick(w http.ResponseWriter, r *http.Request) {
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
	l, err := h.m.display(r.Context(), d.Location)
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

// deviceSlip is a slip of the device's display.
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
		l, lerr := h.m.display(r.Context(), d.Location)
		if lerr != nil {
			h.writeErr(w, lerr)
			return
		}
		if l.Fallback == nil {
			h.writeErr(w, errors.Join(ErrInvalid, errors.New("the display has no printer")))
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
	if err := h.m.WritePage(r.Context(), firstOr(s.Payload, jobID(s.Job)), s.Pages, n-1, width, w); err != nil {
		h.writeErr(w, err)
	}
}
