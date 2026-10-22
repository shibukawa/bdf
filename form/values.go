package form

import (
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/cms"
	"github.com/shibukawa/bdf/sign"
)

// The values of a partly filled form may be signed: a person who fills in
// their part signs what they give (the values the form holds then, the
// names it leaves open, and the template), not the pages, which change as
// the others fill in theirs. The signatures go with the partly filled form
// (Template.Signed), and with the form completed of it, in a record of its
// values (Record, a part of the type bdf.PartFill), where VerifyValues
// checks them: that the values signed are those of the form, and that the
// form leaves open no name that the signer did not leave open.
//
// The pages of the completed form are not signed by them: sign them with
// package sign (bdf sign), whose signature covers the record too.

// ValuesVersion is the version of what a signature of values signs.
const ValuesVersion = 1

// ValueSigner signs the values of a fill (Options.Sign).
type ValueSigner struct {
	// Key signs, and Chain is its certificate and the chain to the root.
	Key   crypto.Signer
	Chain []*x509.Certificate
	// Claims are what the signer says of the signature (its role …).
	Claims map[string]any
	// Time is the signing time (the clock when zero).
	Time time.Time
	// TSA is the URL of a time-stamp authority (RFC 3161), and Client
	// the client that asks it.
	TSA    string
	Client *http.Client
}

// Values is what a signature of values signs.
type Values struct {
	Version int `json:"bdf-form-values"`
	// Format is the format of the template, Source the SHA-256 (hex) of
	// its bytes, and From that of the description of the compiled
	// template the form was made of (Template.From).
	Format string `json:"format"`
	Source string `json:"source"`
	From   string `json:"from,omitempty"`
	// Filled are the values of the form, and Open the names it leaves
	// open, when it was signed.
	Filled map[string]any `json:"filled"`
	Open   []string       `json:"open"`
	Claims map[string]any `json:"claims,omitempty"`
}

// RecordVersion is the version of Record.
const RecordVersion = 1

// Record is the record of a form completed of a partly filled form whose
// values were signed: the part of the type bdf.PartFill.
type Record struct {
	Version int    `json:"version"`
	Format  string `json:"format"`
	Source  string `json:"source"`
	From    string `json:"from,omitempty"`
	// Filled are all of the values of the form.
	Filled map[string]any `json:"filled"`
	// Signed are the signatures of values (the DER of CMS signatures).
	Signed [][]byte `json:"signed"`
}

func sourceDigest(src []byte) string {
	s := sha256.Sum256(src)
	return hex.EncodeToString(s[:])
}

// signValues signs the values of a fill.
func signValues(s *ValueSigner, v *Values) ([]byte, error) {
	v.Version, v.Claims = ValuesVersion, s.Claims
	if v.Open == nil {
		v.Open = []string{}
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("form: the values: %w", err)
	}
	if b, err = sign.Canonical(b); err != nil {
		return nil, err
	}
	at := s.Time
	if at.IsZero() {
		at = time.Now()
	}
	der, err := cms.Sign(b, s.Key, s.Chain, &cms.SignOptions{SigningTime: at})
	if err != nil {
		return nil, fmt.Errorf("form: %w", err)
	}
	if s.TSA != "" {
		sd, err := cms.Parse(der)
		if err != nil {
			return nil, err
		}
		ts, err := cms.RequestTimeStamp(context.Background(), s.Client, s.TSA, sd.Signers[0].SignatureDigest(crypto.SHA256), crypto.SHA256)
		if err != nil {
			return nil, fmt.Errorf("form: %w", err)
		}
		if der, err = cms.AddUnsigned(der, cms.Attribute{Type: cms.OIDTimeStampToken, Values: ts.DER}); err != nil {
			return nil, fmt.Errorf("form: %w", err)
		}
	}
	return der, nil
}

// ValueSignature is what VerifyValues finds of a signature of values.
type ValueSignature struct {
	Index       int
	Certificate *x509.Certificate
	Chain       []*x509.Certificate
	SigningTime time.Time
	TimeStamp   time.Time
	// Values are what it signed.
	Values *Values
	// Err is why it does not hold for the form (nil: it does), TrustErr
	// why its certificate is not trusted.
	Err      error
	TrustErr error
}

// Valid reports whether the signature holds and its signer is trusted.
func (s *ValueSignature) Valid() bool { return s.Err == nil && s.TrustErr == nil }

// VerifyValues checks the signatures of the values of a partly filled form,
// or of a form completed of one (its Record). A document with neither has
// none.
func VerifyValues(r io.ReaderAt, size int64, o *sign.VerifyOptions) ([]*ValueSignature, error) {
	if o == nil {
		o = &sign.VerifyOptions{}
	}
	// (the signatures of a document, after its parts, are no matter here)
	d, err := sign.Read(r, size)
	if err != nil {
		return nil, err
	}
	rd := d.Reader
	if rd.Locked() {
		return nil, errors.New("form: the document is encrypted")
	}
	var cur Values
	var signed [][]byte
	found := false
	for _, e := range rd.Manifest.Parts {
		switch e.T {
		case bdf.PartTemplate:
			def, src, err := ReadTemplate(io.NewSectionReader(r, 0, d.Size), d.Size)
			if err != nil {
				return nil, err
			}
			cur = Values{Format: def.Format, Source: sourceDigest(src), From: def.From, Filled: def.Filled, Open: def.Open}
			signed, found = def.Signed, true
		case bdf.PartFill:
			b, err := templatePart(rd, e.H, bdf.PartFill, "record of the form", maxTemplatePart)
			if err != nil {
				return nil, err
			}
			var rec Record
			if err := json.Unmarshal(b, &rec); err != nil {
				return nil, fmt.Errorf("form: the record of the form: %w", err)
			}
			if rec.Version != RecordVersion {
				return nil, fmt.Errorf("form: the record of the form is of version %d", rec.Version)
			}
			cur = Values{Format: rec.Format, Source: rec.Source, From: rec.From, Filled: rec.Filled}
			signed, found = rec.Signed, true
		}
		if found {
			break
		}
	}
	now := o.Now
	if now.IsZero() {
		now = time.Now()
	}
	var out []*ValueSignature
	for i, der := range signed {
		s := &ValueSignature{Index: i}
		out = append(out, s)
		sd, err := cms.Parse(der)
		if err != nil {
			s.Err = err
			continue
		}
		signer := sd.Signers[0]
		s.Certificate, s.SigningTime = signer.Certificate, signer.SigningTime()
		if err := sd.Verify(signer, nil); err != nil {
			s.Err = err
			continue
		}
		var v Values
		if err := json.Unmarshal(sd.Content, &v); err != nil {
			s.Err = fmt.Errorf("form: what the signature signs: %w", err)
			continue
		}
		s.Values = &v
		s.Err = v.holds(&cur)
		at := now
		if ts, err := signer.TimeStampOf(); err != nil {
			s.Err = errors.Join(s.Err, err)
		} else if ts != nil {
			s.TimeStamp, at = ts.Time, ts.Time
			if _, err := ts.Token.VerifyChain(ts.Token.Signers[0], o.Roots, ts.Time, x509.ExtKeyUsageTimeStamping); err != nil {
				s.TrustErr = fmt.Errorf("form: the certificate of the time-stamp authority: %w", err)
			}
		}
		if chains, err := sd.VerifyChain(signer, o.Roots, at); err != nil {
			s.TrustErr = errors.Join(s.TrustErr, fmt.Errorf("form: the certificate of the signer: %w", err))
		} else {
			s.Chain = chains[0]
		}
	}
	return out, nil
}

// holds checks signed values against those of a form: the same template,
// every value signed among those of the form as it was signed, and no
// name open in the form that the signer did not leave open.
func (v *Values) holds(cur *Values) error {
	var errs []error
	if v.Version != ValuesVersion {
		return fmt.Errorf("form: the signed values are of version %d", v.Version)
	}
	if v.Format != cur.Format || v.Source != cur.Source || v.From != cur.From {
		errs = append(errs, errors.New("form: the values were signed for another template"))
	}
	var changed []string
	diff(v.Filled, cur.Filled, "", &changed)
	for _, name := range changed {
		errs = append(errs, fmt.Errorf("form: %q is not what was signed", name))
	}
	for _, name := range cur.Open {
		if !under(v.Open, name) {
			errs = append(errs, fmt.Errorf("form: %q was filled when the values were signed, and is open now", name))
		}
	}
	return errors.Join(errs...)
}

// diff lists the names whose values signed differ in a form.
func diff(signed, now map[string]any, prefix string, out *[]string) {
	for k, sv := range signed {
		name := join(prefix, k)
		nv, ok := now[k]
		if sm, isObj := sv.(map[string]any); isObj {
			if nm, ok := nv.(map[string]any); ok {
				diff(sm, nm, name, out)
				continue
			}
		}
		if !ok || !reflect.DeepEqual(sv, nv) {
			*out = append(*out, name)
		}
	}
}

// record adds the record of a completed form to its document.
func record(doc *bdf.Document, format, source, from string, filled map[string]any, signed [][]byte) error {
	b, err := json.Marshal(Record{Version: RecordVersion, Format: format, Source: source, From: from, Filled: filled, Signed: signed})
	if err != nil {
		return fmt.Errorf("form: %w", err)
	}
	doc.AddPart(bdf.PartFill, b)
	return nil
}
