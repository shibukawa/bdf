package manager

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/shibukawa/yuruidb"
)

// An endpoint is an entry of the network for applications, over REST:
// POST /in/<name> with the endpoint's key takes a document (PDF, Office,
// images, .bdf …) or data (JSON, YAML, CSV, Excel, Parquet: for a template
// node after it), and the job flows from the endpoint to its nodes. What
// it takes is its Accept and Formats, and what it sends on of a document
// its Output (PutNode). Each endpoint has a key of its own, kept as a
// hash (argon2id); the key is shown once, when the endpoint is made or
// its key rotated.

// InPath is where the endpoints are.
const InPath = "/in"

// OAuthTokenPath is where the apps take their tokens (apps.go).
const OAuthTokenPath = "/oauth/token"

// newKey makes an endpoint's key: 32 random bytes, in base64url.
func newKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "bdfk_" + base64.RawURLEncoding.EncodeToString(b)
}

// Parameters of argon2id (OWASP's first recommendation).
const (
	argonTime    = 2
	argonMemory  = 19 * 1024
	argonThreads = 1
	argonLen     = 32
)

// hashKey hashes a key as a PHC string of argon2id.
func hashKey(key string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	sum := argon2.IDKey([]byte(key), salt, argonTime, argonMemory, argonThreads, argonLen)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(sum)), nil
}

// checkKey tells whether a key is the one a PHC string hashes.
func checkKey(phc, key string) bool {
	parts := strings.Split(phc, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false
	}
	var memory, t uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &t, &threads); err != nil {
		return false
	}
	salt, err1 := base64.RawStdEncoding.DecodeString(parts[4])
	want, err2 := base64.RawStdEncoding.DecodeString(parts[5])
	if err1 != nil || err2 != nil || len(want) == 0 || memory == 0 || t == 0 || threads == 0 {
		return false
	}
	got := argon2.IDKey([]byte(key), salt, t, memory, threads, uint32(len(want)))
	return subtle.ConstantTimeCompare(got, want) == 1
}

// endpoint returns an endpoint node.
func (m *Manager) endpoint(ctx context.Context, name string) (*Node, error) {
	return m.nodeOfKind(ctx, name, KindEndpoint)
}

// Endpoint finds the endpoint a key opens.
func (m *Manager) Endpoint(ctx context.Context, name, key string) (*Node, error) {
	n, err := m.endpoint(ctx, name)
	if err != nil {
		return nil, err
	}
	if key == "" || n.KeyHash == "" || !checkKey(n.KeyHash, key) {
		return nil, fmt.Errorf("%w: the endpoint's key", ErrUnauthorized)
	}
	return n, nil
}

// serveIn takes a request to an endpoint: a document, or data (the body
// as it is), with the sender's metadata in X-Meta-* headers and ?meta.*
// parameters. What the endpoint refuses (ErrFormat) is recorded in the
// audit log.
func (h *handler) serveIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.PathValue("name")
	var ep *Node
	var by *Principal
	var app *Viewer // an app sends with its token: the endpoint must list it
	if v, ok := h.viewerOf(r); ok && (v.Admin || v.Kind == PrincipalApp) && r.Header.Get("Authorization") != "" {
		by = v.Principal()
		n, err := h.m.endpoint(ctx, name)
		if err != nil {
			h.writeErr(w, ErrNotFound)
			return
		}
		ep = n
		if !v.Admin {
			app = v
		}
	} else {
		n, err := h.m.Endpoint(ctx, name, bearer(r))
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				w.Header().Set("WWW-Authenticate", `Bearer realm="bdf print manager endpoint"`)
			}
			h.writeErr(w, err)
			return
		}
		ep = n
		by = &Principal{Kind: PrincipalEndpoint, ID: name, Via: ViaEndpointKey}
	}
	if err := h.m.appMaySend(ctx, app, ep.Name, ep.Accept == AcceptData); err != nil {
		h.writeErr(w, err)
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
	j, err := h.m.Submit(ctx, ep.Name, Input{Name: q.Get("name"), By: by, ClaimedUser: q.Get("user"), Source: "endpoint", Doc: f, Size: n,
		ContentType: r.Header.Get("Content-Type"), Meta: metaOfRequest(r)})
	if err != nil {
		if errors.Is(err, ErrFormat) {
			h.m.audit(ctx, by.Label(), "refuse-document", ep.Name, err.Error())
		}
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, j)
}

// formats answers the formats an endpoint may take.
func (h *handler) formats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, KnownFormats())
}

// The layout of the topology's canvas.

const layoutID = "admin"

// GetLayout returns where the nodes are drawn.
func (m *Manager) GetLayout(ctx context.Context) (*Layout, error) {
	rec, err := m.st.db.Get[layoutRecord](ctx, layoutID)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return &Layout{Positions: map[string]Point{}}, nil
	}
	if err != nil {
		return nil, mapErr(err)
	}
	l := Layout(rec.Data)
	if l.Positions == nil {
		l.Positions = map[string]Point{}
	}
	return &l, nil
}

// PutLayout moves nodes: the positions given replace those of the same
// names; a position of a node that is gone is dropped.
func (m *Manager) PutLayout(ctx context.Context, l Layout) (*Layout, error) {
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	for _, n := range nodes {
		names[n.Name] = true
	}
	rec, err := change(ctx, m.st.db, layoutID, true, func(r *layoutRecord) error {
		if r.Positions == nil {
			r.Positions = map[string]Point{}
		}
		for name, p := range l.Positions {
			r.Positions[name] = p
		}
		for name := range r.Positions {
			if !names[name] {
				delete(r.Positions, name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	out := Layout(rec)
	return &out, nil
}

func (h *handler) getLayout(w http.ResponseWriter, r *http.Request) {
	l, err := h.m.GetLayout(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func (h *handler) putLayout(w http.ResponseWriter, r *http.Request) {
	var l Layout
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&l); err != nil {
		h.writeErr(w, errors.Join(ErrInvalid, err))
		return
	}
	got, err := h.m.PutLayout(r.Context(), l)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, got)
}
