package manager

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"io"
	"time"

	"github.com/shibukawa/bdf/sign"
	"github.com/shibukawa/yuruidb"
)

// A sign node signs the .bdf document that reaches it as a witness does
// (package sign): with the key of the manager, a signature whose claims
// say what the manager saw — the node, the job, who sent it as the
// manager verified it, and when it was received — and sends the signed
// document on. The person who sent the job holds no key: the manager
// vouches that the job came from them (a witness signature, as the
// electronic contract services that sign for their users do). The step
// is a delivery, as those of template and convert nodes are.

// Signer is the key that sign nodes sign with (Options.Signer).
type Signer struct {
	// Key signs (a key in memory, or one in a KMS behind crypto.Signer),
	// and Chain is its certificate followed by those of its chain.
	Key   crypto.Signer
	Chain []*x509.Certificate
	// TSA is the URL of a time-stamp authority (RFC 3161) that stamps the
	// signatures ("" none).
	TSA string
}

// RoleWitness is the role a sign node signs in unless it says another.
const RoleWitness = "witness"

// signPayload signs the .bdf payload of a job at a sign node and keeps the
// signed document as a new payload under the delivery's id.
func (m *Manager) signPayload(ctx context.Context, d *Delivery, n *Node, j *Job, in *Payload) (*Payload, error) {
	s := m.o.Signer
	if s == nil || s.Key == nil || len(s.Chain) == 0 {
		return nil, final("the manager has no key to sign with (manager.sign_key_file and manager.sign_cert_file)")
	}
	if in.Kind != PayloadBDF {
		return nil, final("%s signs .bdf documents; this is %s", n.Name, in.Kind)
	}
	data, err := readBlob(ctx, m.st.docs, in.ID, m.o.MaxDocument)
	if err != nil {
		return nil, err
	}
	doc, err := sign.Read(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, final("the document: %v", err)
	}
	claims := map[string]any{
		"role":     firstOr(n.SignRole, RoleWitness),
		"witness":  "BDF Docflow",
		"node":     n.Name,
		"job":      j.ID,
		"number":   j.Number,
		"received": j.Received.UTC().Format(time.RFC3339),
	}
	if p := j.SubmittedBy; p != nil {
		by := map[string]any{"kind": p.Kind, "via": p.Via}
		if p.ID != "" {
			by["id"] = p.ID
		}
		if p.Name != "" {
			by["name"] = p.Name
		}
		claims["submitted_by"] = by
	}
	sig, err := doc.Sign(s.Key, s.Chain, &sign.Options{Claims: claims, Time: m.now(), TSA: s.TSA})
	if err != nil {
		// (a TSA that did not answer may answer the next try)
		return nil, err
	}
	var out bytes.Buffer
	if err := doc.WriteSingle(&out, sig); err != nil {
		return nil, err
	}
	return m.keepBytes(ctx, d.ID, n.Name, in.Name, in.Pages, out.Bytes())
}

// keepBytes keeps the bytes of a .bdf document a node made as a payload,
// as keepDocument does a document.
func (m *Manager) keepBytes(ctx context.Context, id, node, name string, pages int, b []byte) (*Payload, error) {
	h := sha256.Sum256(b)
	sum, size := hex.EncodeToString(h[:]), int64(len(b))
	_, err := m.st.docs.Insert(ctx, id, bytes.NewReader(b))
	switch {
	case errors.Is(err, yuruidb.ErrAlreadyExists):
		// (kept by a try before, whose outcome was not known)
		r, err := m.st.docs.Open(ctx, id)
		if err != nil {
			return nil, mapErr(err)
		}
		defer r.Close()
		hh := sha256.New()
		if _, err := io.Copy(hh, r); err != nil {
			return nil, err
		}
		sum, size = hex.EncodeToString(hh.Sum(nil)), r.Size
	case err != nil:
		return nil, mapErr(err)
	}
	return &Payload{ID: id, Kind: PayloadBDF, Format: "bdf", Name: name, Pages: pages, Size: size, SHA256: sum, Made: node, At: m.now()}, nil
}
