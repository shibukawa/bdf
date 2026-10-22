// Package sign signs BDF documents and verifies their signatures, so that
// a change made to a document after it was signed is found (spec §3.7).
//
// A signature is a CMS SignedData (RFC 5652) with an X.509 certificate of
// its signer. What it signs is a small JSON document of its own, which it
// holds: the SHA-256 of the canonical JSON (RFC 8785) of the manifest, the
// SHA-256 of each part as it is stored, in the order of the manifest, the
// SHA-256 of each signature before it, and the claims of the signer (who
// signs, in what role). The whole digests of SHA-256 are signed, not the
// names of the parts, which are 16 bytes of them.
//
// The signatures of a single-file document follow its parts, in a block
// with a footer of 16 bytes (readers that do not know of them read the
// document as it is); those of the split form are signatures.json beside
// manifest.json. They are added one after another: a later signature
// covers the ones before it, so that one cannot be taken away unseen.
//
// A signature finds changes; it does not prevent them, and it says nothing
// of whether what was signed was true. A drawn signature in a form (form,
// tag.InputSignature) is a picture, not a signature of this package.
package sign

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/cms"
)

// Version is the version of the signed document (bdf-signature).
const Version = 1

// Limits of the signatures of a document.
const (
	MaxSignatures = 64
	maxBlock      = MaxSignatures * cms.MaxSize
)

// footerMagic starts the footer of the block of signatures.
var footerMagic = [4]byte{'b', 'd', 'f', 's'}

const footerSize = 16

// SplitFile is the file of the signatures of the split form.
const SplitFile = "signatures.json"

// Payload is what a signature signs.
type Payload struct {
	Version int `json:"bdf-signature"`
	// Manifest is the SHA-256 (hex) of the canonical JSON of the stored
	// manifest, without the offsets of its parts.
	Manifest string `json:"manifest"`
	// Parts are the SHA-256 (hex) of the stored bytes of the parts, in the
	// order of the manifest.
	Parts []string `json:"parts"`
	// Previous are the SHA-256 (hex) of the DER of the signatures before
	// this one, in their order.
	Previous []string `json:"previous"`
	// Claims are what the signer says of the signature: its role
	// ("author", "counterparty", "witness" …), and what a witness saw (who
	// did what, when). They are signed with the document.
	Claims map[string]any `json:"claims,omitempty"`
}

// Options are the options of a signature.
type Options struct {
	// Claims are signed with the document (Payload.Claims).
	Claims map[string]any
	// Time is the signing time (the clock when zero). It is the signer's
	// word; a time-stamp is the word of a TSA.
	Time time.Time
	// TSA is the URL of a time-stamp authority (RFC 3161) that stamps the
	// signature, and Client the client that asks it (nil:
	// http.DefaultClient).
	TSA    string
	Client *http.Client
	// Hash is the digest of the signature (SHA-256 when 0). The digests of
	// the payload are SHA-256 always.
	Hash crypto.Hash
}

// Document is a document with its signatures, as Read reads it.
type Document struct {
	// Reader reads the document.
	Reader *bdf.Reader
	// Size is the size of the single-file document without its signatures.
	Size int64
	// Signatures are the DER of its signatures, in their order.
	Signatures [][]byte
	dir        string
	src        io.ReaderAt
}

// Read reads a single-file document and its signatures.
func Read(r io.ReaderAt, size int64) (*Document, error) {
	docSize, sigs, err := readBlock(r, size)
	if err != nil {
		return nil, err
	}
	rd, err := bdf.OpenSingle(r, docSize)
	if err != nil {
		return nil, err
	}
	return &Document{Reader: rd, Size: docSize, Signatures: sigs, src: r}, nil
}

// ReadSplit reads a document of the split form and its signatures.
func ReadSplit(dir string) (*Document, error) {
	rd, err := bdf.OpenSplit(dir)
	if err != nil {
		return nil, err
	}
	d := &Document{Reader: rd, dir: dir}
	b, err := os.ReadFile(filepath.Join(dir, SplitFile))
	if errors.Is(err, os.ErrNotExist) {
		return d, nil
	}
	if err != nil {
		return nil, err
	}
	var f struct {
		Version    int      `json:"version"`
		Signatures [][]byte `json:"signatures"`
	}
	if len(b) > 2*maxBlock {
		return nil, errors.New("sign: signatures.json is too large")
	}
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, fmt.Errorf("sign: %s: %w", SplitFile, err)
	}
	if f.Version != 1 || len(f.Signatures) > MaxSignatures {
		return nil, fmt.Errorf("sign: %s is of version %d, with %d signatures", SplitFile, f.Version, len(f.Signatures))
	}
	d.Signatures = f.Signatures
	return d, nil
}

// readBlock finds the block of signatures at the end of a single file: the
// size of the document before it, and the signatures.
func readBlock(r io.ReaderAt, size int64) (int64, [][]byte, error) {
	if size < bdf.HeaderSize+footerSize {
		return size, nil, nil
	}
	var foot [footerSize]byte
	if _, err := r.ReadAt(foot[:], size-footerSize); err != nil {
		return 0, nil, err
	}
	if !bytes.Equal(foot[:4], footerMagic[:]) {
		return size, nil, nil
	}
	version := binary.LittleEndian.Uint16(foot[4:])
	count := int(binary.LittleEndian.Uint16(foot[6:]))
	blockLen := binary.LittleEndian.Uint64(foot[8:])
	if version != 1 {
		return 0, nil, fmt.Errorf("sign: the signatures are of version %d", version)
	}
	if count == 0 || count > MaxSignatures || blockLen > maxBlock || blockLen > uint64(size-footerSize-bdf.HeaderSize) {
		return 0, nil, errors.New("sign: the block of signatures is damaged")
	}
	block := make([]byte, blockLen)
	start := size - footerSize - int64(blockLen)
	if _, err := r.ReadAt(block, start); err != nil {
		return 0, nil, err
	}
	var sigs [][]byte
	for len(block) > 0 {
		if len(block) < 4 {
			return 0, nil, errors.New("sign: the block of signatures is damaged")
		}
		n := binary.LittleEndian.Uint32(block)
		if uint64(n) > uint64(len(block)-4) || n > cms.MaxSize {
			return 0, nil, errors.New("sign: the block of signatures is damaged")
		}
		sigs = append(sigs, block[4:4+n])
		block = block[4+n:]
	}
	if len(sigs) != count {
		return 0, nil, errors.New("sign: the block of signatures is damaged")
	}
	return start, sigs, nil
}

// writeBlock writes a block of signatures and its footer.
func writeBlock(w io.Writer, sigs [][]byte) error {
	if len(sigs) == 0 {
		return nil
	}
	var block []byte
	for _, s := range sigs {
		block = binary.LittleEndian.AppendUint32(block, uint32(len(s)))
		block = append(block, s...)
	}
	foot := append([]byte{}, footerMagic[:]...)
	foot = binary.LittleEndian.AppendUint16(foot, 1)
	foot = binary.LittleEndian.AppendUint16(foot, uint16(len(sigs)))
	foot = binary.LittleEndian.AppendUint64(foot, uint64(len(block)))
	if _, err := w.Write(block); err != nil {
		return err
	}
	_, err := w.Write(foot)
	return err
}

// digests computes the digests of a document: of its manifest, and of its
// parts.
func digests(rd *bdf.Reader) (manifest string, parts []string, err error) {
	m, mj := rd.StoredManifest()
	var v map[string]any
	d := json.NewDecoder(bytes.NewReader(mj))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		return "", nil, fmt.Errorf("sign: the manifest: %w", err)
	}
	// (the offsets are those of the single form, which the split form
	// does not have)
	if list, ok := v["parts"].([]any); ok {
		for _, e := range list {
			if o, ok := e.(map[string]any); ok {
				delete(o, "off")
			}
		}
	}
	c, err := CanonicalValue(v)
	if err != nil {
		return "", nil, err
	}
	sum := sha256.Sum256(c)
	manifest = hex.EncodeToString(sum[:])
	for _, e := range m.Parts {
		b, err := rd.StoredPart(e)
		if err != nil {
			return "", nil, fmt.Errorf("sign: part %s: %w", e.H, err)
		}
		s := sha256.Sum256(b)
		parts = append(parts, hex.EncodeToString(s[:]))
	}
	return manifest, parts, nil
}

func digestOf(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

// Sign makes a signature of a document after the signatures it has: the
// DER of a SignedData. Add writes the document with it.
func (d *Document) Sign(key crypto.Signer, chain []*x509.Certificate, o *Options) ([]byte, error) {
	if o == nil {
		o = &Options{}
	}
	if len(d.Signatures) >= MaxSignatures {
		return nil, fmt.Errorf("sign: a document has %d signatures at most", MaxSignatures)
	}
	p := Payload{Version: Version, Claims: o.Claims, Previous: []string{}}
	var err error
	if p.Manifest, p.Parts, err = digests(d.Reader); err != nil {
		return nil, err
	}
	for _, s := range d.Signatures {
		p.Previous = append(p.Previous, digestOf(s))
	}
	pj, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("sign: the claims: %w", err)
	}
	if pj, err = Canonical(pj); err != nil {
		return nil, err
	}
	at := o.Time
	if at.IsZero() {
		at = time.Now()
	}
	der, err := cms.Sign(pj, key, chain, &cms.SignOptions{Hash: o.Hash, SigningTime: at})
	if err != nil {
		return nil, err
	}
	if o.TSA != "" {
		sd, err := cms.Parse(der)
		if err != nil {
			return nil, err
		}
		ts, err := cms.RequestTimeStamp(context.Background(), o.Client, o.TSA, sd.Signers[0].SignatureDigest(crypto.SHA256), crypto.SHA256)
		if err != nil {
			return nil, err
		}
		if der, err = cms.AddUnsigned(der, cms.Attribute{Type: cms.OIDTimeStampToken, Values: ts.DER}); err != nil {
			return nil, err
		}
	}
	return der, nil
}

// WriteSingle writes a single-file document with its signatures and more
// (Sign). The bytes of the document are copied as they are.
func (d *Document) WriteSingle(w io.Writer, more ...[]byte) error {
	if d.src == nil {
		return errors.New("sign: the document is of the split form")
	}
	if _, err := io.Copy(w, io.NewSectionReader(d.src, 0, d.Size)); err != nil {
		return err
	}
	return writeBlock(w, append(append([][]byte{}, d.Signatures...), more...))
}

// WriteSplit writes the signatures of a document of the split form, with
// more, to its signatures.json.
func (d *Document) WriteSplit(more ...[]byte) error {
	if d.dir == "" {
		return errors.New("sign: the document is of the single-file form")
	}
	all := append(append([][]byte{}, d.Signatures...), more...)
	b, err := json.MarshalIndent(map[string]any{"version": 1, "signatures": all}, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d.dir, SplitFile), append(b, '\n'), 0o644)
}

// VerifyOptions are the options of Verify.
type VerifyOptions struct {
	// Roots are the certificates trusted (nil: those of the system).
	Roots *x509.CertPool
	// Now is the time to check certificates at when a signature has no
	// time-stamp (the clock when zero).
	Now time.Time
}

// Signature is what Verify finds of a signature.
type Signature struct {
	// Index is the place of the signature among those of the document.
	Index int
	// Certificate is the certificate of the signer, and Chain its chain
	// to a trusted root (none when it is not trusted).
	Certificate *x509.Certificate
	Chain       []*x509.Certificate
	// SigningTime is the time the signer says; TimeStamp the time a TSA
	// stamped the signature at (zero without one), and TSA its
	// certificate.
	SigningTime time.Time
	TimeStamp   time.Time
	TSA         *x509.Certificate
	// Claims are the claims of the signer.
	Claims map[string]any
	// Covers tells what the signature covers: the whole of the document
	// with the signatures before it.
	Payload *Payload
	// Err is why the signature does not hold for the document (nil: it
	// does): the document or a signature before it changed after it was
	// signed, or the signature is damaged. TrustErr is why its
	// certificate is not trusted, while the signature itself holds.
	Err      error
	TrustErr error
}

// Valid reports whether the signature holds and its certificate is
// trusted.
func (s *Signature) Valid() bool { return s.Err == nil && s.TrustErr == nil }

// Verify checks every signature of a document. It fails only when the
// document cannot be read; what it finds of each signature is in its
// Signature.
func (d *Document) Verify(o *VerifyOptions) ([]*Signature, error) {
	if o == nil {
		o = &VerifyOptions{}
	}
	manifest, parts, err := digests(d.Reader)
	if err != nil {
		return nil, err
	}
	now := o.Now
	if now.IsZero() {
		now = time.Now()
	}
	var out []*Signature
	for i, der := range d.Signatures {
		s := &Signature{Index: i}
		out = append(out, s)
		sd, err := cms.Parse(der)
		if err != nil {
			s.Err = err
			continue
		}
		signer := sd.Signers[0]
		s.Certificate, s.SigningTime = signer.Certificate, signer.SigningTime()
		if !sd.ContentType.Equal(cms.OIDData) || sd.Content == nil {
			s.Err = errors.New("sign: the signature holds no payload")
			continue
		}
		if err := sd.Verify(signer, nil); err != nil {
			s.Err = err
			continue
		}
		var p Payload
		if err := json.Unmarshal(sd.Content, &p); err != nil {
			s.Err = fmt.Errorf("sign: the payload: %w", err)
			continue
		}
		s.Payload, s.Claims = &p, p.Claims
		s.Err = p.check(manifest, parts, d.Signatures[:i])
		at := now
		ts, err := signer.TimeStampOf()
		switch {
		case err != nil:
			s.Err = errors.Join(s.Err, err)
		case ts != nil:
			s.TimeStamp, at = ts.Time, ts.Time
			if len(ts.Token.Signers) > 0 {
				s.TSA = ts.Token.Signers[0].Certificate
				if _, err := ts.Token.VerifyChain(ts.Token.Signers[0], o.Roots, ts.Time, x509.ExtKeyUsageTimeStamping); err != nil {
					s.TrustErr = fmt.Errorf("sign: the certificate of the time-stamp authority: %w", err)
				}
			}
		}
		chains, err := sd.VerifyChain(signer, o.Roots, at)
		if err != nil {
			s.TrustErr = errors.Join(s.TrustErr, fmt.Errorf("sign: the certificate of the signer: %w", err))
		} else {
			s.Chain = chains[0]
		}
	}
	return out, nil
}

// check compares a payload with the document it is verified against.
func (p *Payload) check(manifest string, parts []string, before [][]byte) error {
	if p.Version != Version {
		return fmt.Errorf("sign: the payload is of version %d", p.Version)
	}
	var errs []error
	if p.Manifest != manifest {
		errs = append(errs, errors.New("sign: the manifest changed after the signature (its views, pages, metadata or list of parts)"))
	}
	if len(p.Parts) != len(parts) {
		errs = append(errs, fmt.Errorf("sign: the document has %d parts; %d were signed", len(parts), len(p.Parts)))
	} else {
		changed := 0
		for i := range parts {
			if p.Parts[i] != parts[i] {
				changed++
			}
		}
		if changed > 0 {
			errs = append(errs, fmt.Errorf("sign: %d part(s) changed after the signature", changed))
		}
	}
	if len(p.Previous) != len(before) {
		errs = append(errs, fmt.Errorf("sign: %d signature(s) came before this one when it was made; the document has %d", len(p.Previous), len(before)))
	} else {
		for i, b := range before {
			if p.Previous[i] != digestOf(b) {
				errs = append(errs, fmt.Errorf("sign: signature %d before this one changed", i+1))
			}
		}
	}
	return errors.Join(errs...)
}
