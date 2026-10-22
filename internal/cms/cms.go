// Package cms writes and reads the signatures of the Cryptographic Message
// Syntax (RFC 5652): SignedData with one signer or more, as the signatures
// of BDF documents (package sign) and of PDF files (PAdES, vector/pdfbdf)
// are, and the time-stamp tokens of RFC 3161, which are SignedData too.
//
// It is the part of CMS that those need: a signer identified by the issuer
// and serial number of its certificate (or, read, by its key identifier),
// signed attributes (the content type, the digest of the content, the
// signing time and the signing certificate of ESS), and unsigned ones (a
// time-stamp token), with ECDSA or RSA keys and SHA-2 digests. What it
// reads is untrusted: lengths, counts and nesting are bounded.
package cms

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"
)

// Object identifiers.
var (
	OIDData                 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	OIDSignedData           = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}
	OIDContentType          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 3}
	OIDMessageDigest        = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}
	OIDSigningTime          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}
	OIDSigningCertificateV2 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 47}
	OIDTimeStampToken       = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 14}
	OIDTSTInfo              = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 1, 4}

	oidSHA256 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
	oidSHA384 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 2}
	oidSHA512 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 3}

	oidSHA256WithRSA = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11}
	oidSHA384WithRSA = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 12}
	oidSHA512WithRSA = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 13}
	oidECDSASHA256   = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}
	oidECDSASHA384   = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 3}
	oidECDSASHA512   = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 4}
	oidEd25519       = asn1.ObjectIdentifier{1, 3, 101, 112}
)

// Limits of what is read.
const (
	// MaxSize bounds the DER of a SignedData.
	MaxSize = 4 << 20
	// maxItems bounds the certificates, signers and attributes.
	maxItems = 64
)

// hashOf is the digest that an algorithm identifier names.
func hashOf(oid asn1.ObjectIdentifier) (crypto.Hash, bool) {
	switch {
	case oid.Equal(oidSHA256):
		return crypto.SHA256, true
	case oid.Equal(oidSHA384):
		return crypto.SHA384, true
	case oid.Equal(oidSHA512):
		return crypto.SHA512, true
	}
	return 0, false
}

func hashOID(h crypto.Hash) asn1.ObjectIdentifier {
	switch h {
	case crypto.SHA384:
		return oidSHA384
	case crypto.SHA512:
		return oidSHA512
	}
	return oidSHA256
}

type contentInfo struct {
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"explicit,optional,tag:0"`
}

type encapContentInfo struct {
	EContentType asn1.ObjectIdentifier
	EContent     asn1.RawValue `asn1:"explicit,optional,tag:0"`
}

type signedData struct {
	Version          int
	DigestAlgorithms asn1.RawValue
	EncapContentInfo encapContentInfo
	Certificates     asn1.RawValue `asn1:"optional,tag:0"`
	CRLs             asn1.RawValue `asn1:"optional,tag:1"`
	SignerInfos      asn1.RawValue
}

type issuerAndSerial struct {
	Issuer asn1.RawValue
	Serial *big.Int
}

type signerInfo struct {
	Version            int
	SID                asn1.RawValue
	DigestAlgorithm    pkix.AlgorithmIdentifier
	SignedAttrs        asn1.RawValue `asn1:"optional,tag:0"`
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Signature          []byte
	UnsignedAttrs      asn1.RawValue `asn1:"optional,tag:1"`
}

type attribute struct {
	Type   asn1.ObjectIdentifier
	Values asn1.RawValue
}

type essCertIDv2 struct {
	// (the hash algorithm is SHA-256, the default, and so not written)
	CertHash []byte
}

type signingCertificateV2 struct {
	Certs []essCertIDv2
}

// Attribute is an attribute of a signer: its type, and the DER of its
// values (the content of their SET).
type Attribute struct {
	Type   asn1.ObjectIdentifier
	Values []byte
}

// SignOptions are the options of Sign.
type SignOptions struct {
	// Hash is the digest of the content and of the attributes (SHA-256
	// when 0; SHA-384 and SHA-512 also).
	Hash crypto.Hash
	// Detached leaves the content out: who verifies has it.
	Detached bool
	// ContentType is the type of the content (id-data when nil).
	ContentType asn1.ObjectIdentifier
	// SigningTime is written as the signing-time attribute, unless it is
	// zero (PAdES puts the time in the signature dictionary instead).
	SigningTime time.Time
	// NoSigningCertificate leaves out the signing-certificate-v2 attribute
	// of ESS (which CAdES and PAdES require).
	NoSigningCertificate bool
	// Attributes are more signed attributes.
	Attributes []Attribute
}

// Sign signs content with a key whose certificate is the first of chain;
// the other certificates of chain go with it (the intermediate ones). It
// returns the DER of a ContentInfo of SignedData.
func Sign(content []byte, key crypto.Signer, chain []*x509.Certificate, o *SignOptions) ([]byte, error) {
	if o == nil {
		o = &SignOptions{}
	}
	if len(chain) == 0 {
		return nil, errors.New("cms: no certificate of the key")
	}
	h := o.Hash
	if h == 0 {
		h = crypto.SHA256
	}
	if !h.Available() || (h != crypto.SHA256 && h != crypto.SHA384 && h != crypto.SHA512) {
		return nil, fmt.Errorf("cms: the digest %v is not one of SHA-256, SHA-384 and SHA-512", h)
	}
	ctype := o.ContentType
	if ctype == nil {
		ctype = OIDData
	}
	cert := chain[0]
	if !publicKeyEqual(cert.PublicKey, key.Public()) {
		return nil, errors.New("cms: the certificate is not that of the key")
	}
	sigAlg, err := signatureAlgorithm(key.Public(), h)
	if err != nil {
		return nil, err
	}
	digest := h.New()
	digest.Write(content)

	var attrs []attribute
	add := func(t asn1.ObjectIdentifier, v any) error {
		b, err := asn1.Marshal(v)
		if err != nil {
			return err
		}
		attrs = append(attrs, attribute{Type: t, Values: asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagSet, IsCompound: true, Bytes: b}})
		return nil
	}
	if err := add(OIDContentType, ctype); err != nil {
		return nil, err
	}
	if !o.SigningTime.IsZero() {
		if err := add(OIDSigningTime, o.SigningTime.UTC()); err != nil {
			return nil, err
		}
	}
	if err := add(OIDMessageDigest, digest.Sum(nil)); err != nil {
		return nil, err
	}
	if !o.NoSigningCertificate {
		sum := sha256.Sum256(cert.Raw)
		if err := add(OIDSigningCertificateV2, signingCertificateV2{Certs: []essCertIDv2{{CertHash: sum[:]}}}); err != nil {
			return nil, err
		}
	}
	for _, a := range o.Attributes {
		attrs = append(attrs, attribute{Type: a.Type, Values: asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagSet, IsCompound: true, Bytes: a.Values}})
	}
	signed, err := setOf(attrs)
	if err != nil {
		return nil, err
	}
	// the attributes are signed as the SET that they are, and stored as
	// the IMPLICIT [0] of the signer
	toSign := append([]byte{0x31}, signed[1:]...)
	ad := h.New()
	ad.Write(toSign)
	var sig []byte
	if _, ok := key.Public().(ed25519.PublicKey); ok {
		sig, err = key.Sign(rand.Reader, toSign, crypto.Hash(0))
	} else {
		sig, err = key.Sign(rand.Reader, ad.Sum(nil), h)
	}
	if err != nil {
		return nil, fmt.Errorf("cms: signing: %w", err)
	}
	sid, err := asn1.Marshal(issuerAndSerial{Issuer: asn1.RawValue{FullBytes: cert.RawIssuer}, Serial: cert.SerialNumber})
	if err != nil {
		return nil, err
	}
	si := signerInfo{
		Version:            1,
		SID:                asn1.RawValue{FullBytes: sid},
		DigestAlgorithm:    pkix.AlgorithmIdentifier{Algorithm: hashOID(h)},
		SignedAttrs:        asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: 0, IsCompound: true, Bytes: rawContent(signed)},
		SignatureAlgorithm: sigAlg,
		Signature:          sig,
	}
	sib, err := asn1.Marshal(si)
	if err != nil {
		return nil, err
	}
	digestAlgs, err := asn1.Marshal(pkix.AlgorithmIdentifier{Algorithm: hashOID(h)})
	if err != nil {
		return nil, err
	}
	var certs []byte
	for _, c := range chain {
		certs = append(certs, c.Raw...)
	}
	sd := signedData{
		Version:          1,
		DigestAlgorithms: asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagSet, IsCompound: true, Bytes: digestAlgs},
		EncapContentInfo: encapContentInfo{EContentType: ctype},
		Certificates:     asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: 0, IsCompound: true, Bytes: certs},
		SignerInfos:      asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagSet, IsCompound: true, Bytes: sib},
	}
	if !ctype.Equal(OIDData) {
		sd.Version = 3
	}
	if !o.Detached {
		oct, err := asn1.Marshal(content)
		if err != nil {
			return nil, err
		}
		sd.EncapContentInfo.EContent = asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: 0, IsCompound: true, Bytes: oct}
	}
	sdb, err := asn1.Marshal(sd)
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(contentInfo{ContentType: OIDSignedData, Content: asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: 0, IsCompound: true, Bytes: sdb}})
}

// setOf writes a SET OF attributes in the order of DER: of their
// encodings.
func setOf(attrs []attribute) ([]byte, error) {
	enc := make([][]byte, len(attrs))
	for i, a := range attrs {
		b, err := asn1.Marshal(a)
		if err != nil {
			return nil, err
		}
		enc[i] = b
	}
	slices.SortFunc(enc, bytes.Compare)
	return asn1.Marshal(asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagSet, IsCompound: true, Bytes: bytes.Join(enc, nil)})
}

// rawContent is the content of a DER element: what follows its tag and
// length.
func rawContent(der []byte) []byte {
	var v asn1.RawValue
	if _, err := asn1.Unmarshal(der, &v); err != nil {
		return nil
	}
	return v.Bytes
}

func signatureAlgorithm(pub crypto.PublicKey, h crypto.Hash) (pkix.AlgorithmIdentifier, error) {
	switch pub.(type) {
	case *ecdsa.PublicKey:
		oid := map[crypto.Hash]asn1.ObjectIdentifier{crypto.SHA256: oidECDSASHA256, crypto.SHA384: oidECDSASHA384, crypto.SHA512: oidECDSASHA512}[h]
		return pkix.AlgorithmIdentifier{Algorithm: oid}, nil
	case *rsa.PublicKey:
		oid := map[crypto.Hash]asn1.ObjectIdentifier{crypto.SHA256: oidSHA256WithRSA, crypto.SHA384: oidSHA384WithRSA, crypto.SHA512: oidSHA512WithRSA}[h]
		return pkix.AlgorithmIdentifier{Algorithm: oid, Parameters: asn1.NullRawValue}, nil
	case ed25519.PublicKey:
		return pkix.AlgorithmIdentifier{Algorithm: oidEd25519}, nil
	}
	return pkix.AlgorithmIdentifier{}, fmt.Errorf("cms: keys of the type %T are not used", pub)
}

func publicKeyEqual(a, b crypto.PublicKey) bool {
	type equaler interface{ Equal(crypto.PublicKey) bool }
	if e, ok := a.(equaler); ok {
		return e.Equal(b)
	}
	return false
}

// SignedData is a SignedData that was read.
type SignedData struct {
	// ContentType is the type of the content, and Content the content
	// that it holds (nil when it is detached).
	ContentType asn1.ObjectIdentifier
	Content     []byte
	// Certificates are the certificates that it carries.
	Certificates []*x509.Certificate
	Signers      []*Signer
}

// Signer is a signer of a SignedData.
type Signer struct {
	// Certificate is the certificate of the signer among those of the
	// SignedData (nil when it is not among them).
	Certificate *x509.Certificate
	// Hash is the digest of the signer.
	Hash crypto.Hash
	// SignedAttributes and UnsignedAttributes are its attributes.
	SignedAttributes   []Attribute
	UnsignedAttributes []Attribute
	// Signature is the value of its signature.
	Signature []byte

	sigAlg    asn1.ObjectIdentifier
	signedDER []byte // the SET of the signed attributes, as signed
	sidIssuer []byte
	sidSerial *big.Int
	sidKeyID  []byte
	raw       []byte // the DER of the SignerInfo
}

// Parse reads the DER of a ContentInfo of SignedData.
func Parse(der []byte) (*SignedData, error) {
	if len(der) > MaxSize {
		return nil, fmt.Errorf("cms: more than %d MiB", MaxSize>>20)
	}
	var ci contentInfo
	rest, err := asn1.Unmarshal(der, &ci)
	if err != nil {
		return nil, fmt.Errorf("cms: %w", err)
	}
	if len(rest) > 0 && len(bytes.TrimRight(rest, "\x00")) > 0 {
		return nil, errors.New("cms: data after the signature")
	}
	if !ci.ContentType.Equal(OIDSignedData) {
		return nil, errors.New("cms: not SignedData")
	}
	var sd signedData
	if _, err := asn1.Unmarshal(ci.Content.Bytes, &sd); err != nil {
		return nil, fmt.Errorf("cms: SignedData: %w", err)
	}
	out := &SignedData{ContentType: sd.EncapContentInfo.EContentType}
	if len(sd.EncapContentInfo.EContent.Bytes) > 0 {
		// (an OCTET STRING, or the constructed one of BER, of pieces)
		c, err := octets(sd.EncapContentInfo.EContent.Bytes)
		if err != nil {
			return nil, fmt.Errorf("cms: the content: %w", err)
		}
		out.Content = c
	}
	for rest := sd.Certificates.Bytes; len(rest) > 0; {
		var v asn1.RawValue
		var err error
		if rest, err = asn1.Unmarshal(rest, &v); err != nil {
			return nil, fmt.Errorf("cms: certificates: %w", err)
		}
		if len(out.Certificates) >= maxItems {
			return nil, errors.New("cms: too many certificates")
		}
		if v.Class != asn1.ClassUniversal || v.Tag != asn1.TagSequence {
			continue // (an attribute certificate, or another kind)
		}
		c, err := x509.ParseCertificate(v.FullBytes)
		if err != nil {
			return nil, fmt.Errorf("cms: a certificate: %w", err)
		}
		out.Certificates = append(out.Certificates, c)
	}
	for rest := sd.SignerInfos.Bytes; len(rest) > 0; {
		var si signerInfo
		start := rest
		var err error
		if rest, err = asn1.Unmarshal(rest, &si); err != nil {
			return nil, fmt.Errorf("cms: a signer: %w", err)
		}
		if len(out.Signers) >= maxItems {
			return nil, errors.New("cms: too many signers")
		}
		s := &Signer{Signature: si.Signature, sigAlg: si.SignatureAlgorithm.Algorithm, raw: start[:len(start)-len(rest)]}
		h, ok := hashOf(si.DigestAlgorithm.Algorithm)
		if !ok {
			return nil, fmt.Errorf("cms: the digest %v is not known", si.DigestAlgorithm.Algorithm)
		}
		s.Hash = h
		switch {
		case si.SID.Class == asn1.ClassUniversal && si.SID.Tag == asn1.TagSequence:
			var ias issuerAndSerial
			if _, err := asn1.Unmarshal(si.SID.FullBytes, &ias); err != nil {
				return nil, fmt.Errorf("cms: the signer: %w", err)
			}
			s.sidIssuer, s.sidSerial = ias.Issuer.FullBytes, ias.Serial
		case si.SID.Class == asn1.ClassContextSpecific && si.SID.Tag == 0:
			s.sidKeyID = si.SID.Bytes
		default:
			return nil, errors.New("cms: the signer is not identified")
		}
		if len(si.SignedAttrs.FullBytes) > 0 {
			s.signedDER = append([]byte{0x31}, si.SignedAttrs.FullBytes[1:]...)
			if s.SignedAttributes, err = attributes(si.SignedAttrs.Bytes); err != nil {
				return nil, err
			}
		}
		if len(si.UnsignedAttrs.FullBytes) > 0 {
			if s.UnsignedAttributes, err = attributes(si.UnsignedAttrs.Bytes); err != nil {
				return nil, err
			}
		}
		for _, c := range out.Certificates {
			if s.identifies(c) {
				s.Certificate = c
				break
			}
		}
		out.Signers = append(out.Signers, s)
	}
	if len(out.Signers) == 0 {
		return nil, errors.New("cms: no signer")
	}
	return out, nil
}

// octets reads the content of an OCTET STRING, primitive or constructed.
func octets(der []byte) ([]byte, error) {
	var v asn1.RawValue
	if _, err := asn1.Unmarshal(der, &v); err != nil {
		return nil, err
	}
	if v.Tag != asn1.TagOctetString {
		return nil, errors.New("not an OCTET STRING")
	}
	if !v.IsCompound {
		return v.Bytes, nil
	}
	var out []byte
	for rest, n := v.Bytes, 0; len(rest) > 0; n++ {
		var p asn1.RawValue
		var err error
		if rest, err = asn1.Unmarshal(rest, &p); err != nil {
			return nil, err
		}
		if p.Tag != asn1.TagOctetString || p.IsCompound || n > 1<<16 {
			return nil, errors.New("a constructed OCTET STRING of other pieces")
		}
		out = append(out, p.Bytes...)
	}
	return out, nil
}

func attributes(der []byte) ([]Attribute, error) {
	var out []Attribute
	for rest := der; len(rest) > 0; {
		var a attribute
		var err error
		if rest, err = asn1.Unmarshal(rest, &a); err != nil {
			return nil, fmt.Errorf("cms: an attribute: %w", err)
		}
		if len(out) >= maxItems {
			return nil, errors.New("cms: too many attributes")
		}
		out = append(out, Attribute{Type: a.Type, Values: a.Values.Bytes})
	}
	return out, nil
}

// identifies reports whether a certificate is that of the signer.
func (s *Signer) identifies(c *x509.Certificate) bool {
	if s.sidKeyID != nil {
		return bytes.Equal(c.SubjectKeyId, s.sidKeyID)
	}
	return bytes.Equal(c.RawIssuer, s.sidIssuer) && c.SerialNumber.Cmp(s.sidSerial) == 0
}

// Attribute returns the first value of a signed attribute (nil when the
// signer has none of the type).
func (s *Signer) Attribute(t asn1.ObjectIdentifier) []byte {
	for _, a := range s.SignedAttributes {
		if a.Type.Equal(t) {
			return a.Values
		}
	}
	return nil
}

// Unsigned returns the first value of an unsigned attribute.
func (s *Signer) Unsigned(t asn1.ObjectIdentifier) []byte {
	for _, a := range s.UnsignedAttributes {
		if a.Type.Equal(t) {
			return a.Values
		}
	}
	return nil
}

// SigningTime is the time of the signing-time attribute (zero without
// one).
func (s *Signer) SigningTime() time.Time {
	b := s.Attribute(OIDSigningTime)
	if b == nil {
		return time.Time{}
	}
	var t time.Time
	if _, err := asn1.Unmarshal(b, &t); err != nil {
		return time.Time{}
	}
	return t
}

// Verify checks the signature of the signer over content (that of the
// SignedData, or the detached content), with the certificate of the
// signer: the digest of the content in its attributes, the content type,
// and the signing certificate of ESS when there is one. It does not check
// the certificate itself (VerifyChain).
func (sd *SignedData) Verify(s *Signer, content []byte) error {
	if s.Certificate == nil {
		return errors.New("cms: the certificate of the signer is not in the signature")
	}
	if content == nil {
		content = sd.Content
	}
	if s.signedDER == nil {
		return errors.New("cms: the signer has no signed attributes")
	}
	md := s.Attribute(OIDMessageDigest)
	var want []byte
	if md == nil {
		return errors.New("cms: the signer has no message digest")
	}
	if _, err := asn1.Unmarshal(md, &want); err != nil {
		return fmt.Errorf("cms: the message digest: %w", err)
	}
	h := s.Hash.New()
	h.Write(content)
	if !bytes.Equal(h.Sum(nil), want) {
		return errors.New("cms: the content is not what was signed (the digest differs)")
	}
	if ct := s.Attribute(OIDContentType); ct != nil {
		var oid asn1.ObjectIdentifier
		if _, err := asn1.Unmarshal(ct, &oid); err != nil || !oid.Equal(sd.ContentType) {
			return errors.New("cms: the content type of the signer is not that of the content")
		}
	} else {
		return errors.New("cms: the signer has no content type")
	}
	if sc := s.Attribute(OIDSigningCertificateV2); sc != nil {
		var v signingCertificateV2
		if _, err := asn1.Unmarshal(sc, &v); err == nil && len(v.Certs) > 0 {
			sum := sha256.Sum256(s.Certificate.Raw)
			if !bytes.Equal(v.Certs[0].CertHash, sum[:]) {
				return errors.New("cms: the signing certificate is not the certificate of the signer")
			}
		}
	}
	return verifySignature(s.Certificate.PublicKey, s.sigAlg, s.Hash, s.signedDER, s.Signature)
}

func verifySignature(pub crypto.PublicKey, alg asn1.ObjectIdentifier, h crypto.Hash, signed, sig []byte) error {
	if k, ok := pub.(ed25519.PublicKey); ok {
		if !ed25519.Verify(k, signed, sig) {
			return errors.New("cms: the signature is not that of the signer")
		}
		return nil
	}
	d := h.New()
	d.Write(signed)
	sum := d.Sum(nil)
	switch k := pub.(type) {
	case *ecdsa.PublicKey:
		if !ecdsa.VerifyASN1(k, sum, sig) {
			return errors.New("cms: the signature is not that of the signer")
		}
	case *rsa.PublicKey:
		if alg.Equal(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 10}) {
			if err := rsa.VerifyPSS(k, h, sum, sig, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}); err != nil {
				return errors.New("cms: the signature is not that of the signer")
			}
		} else if err := rsa.VerifyPKCS1v15(k, h, sum, sig); err != nil {
			return errors.New("cms: the signature is not that of the signer")
		}
	default:
		return fmt.Errorf("cms: keys of the type %T are not verified", pub)
	}
	return nil
}

// VerifyChain checks the certificate of a signer against roots (nil: the
// roots of the system), at a time (that of a time-stamp, or now), with the
// intermediate certificates of the SignedData. It returns the chain.
func (sd *SignedData) VerifyChain(s *Signer, roots *x509.CertPool, at time.Time, usage ...x509.ExtKeyUsage) ([][]*x509.Certificate, error) {
	if s.Certificate == nil {
		return nil, errors.New("cms: the certificate of the signer is not in the signature")
	}
	inter := x509.NewCertPool()
	for _, c := range sd.Certificates {
		if c != s.Certificate {
			inter.AddCert(c)
		}
	}
	if len(usage) == 0 {
		usage = []x509.ExtKeyUsage{x509.ExtKeyUsageAny}
	}
	return s.Certificate.Verify(x509.VerifyOptions{Roots: roots, Intermediates: inter, CurrentTime: at, KeyUsages: usage})
}

// SignatureDigest is the digest of the value of the signature of a signer:
// what a time-stamp token for it stamps.
func (s *Signer) SignatureDigest(h crypto.Hash) []byte {
	d := h.New()
	d.Write(s.Signature)
	return d.Sum(nil)
}

// AddUnsigned adds an unsigned attribute to the first signer of the DER of
// a SignedData, and returns the new DER. The signature stays valid: the
// unsigned attributes are not signed.
func AddUnsigned(der []byte, a Attribute) ([]byte, error) {
	var ci contentInfo
	if _, err := asn1.Unmarshal(der, &ci); err != nil {
		return nil, fmt.Errorf("cms: %w", err)
	}
	var sd signedData
	if _, err := asn1.Unmarshal(ci.Content.Bytes, &sd); err != nil {
		return nil, fmt.Errorf("cms: SignedData: %w", err)
	}
	var si signerInfo
	rest, err := asn1.Unmarshal(sd.SignerInfos.Bytes, &si)
	if err != nil {
		return nil, fmt.Errorf("cms: a signer: %w", err)
	}
	var attrs []attribute
	old, err := attributes(si.UnsignedAttrs.Bytes)
	if err != nil {
		return nil, err
	}
	for _, o := range append(old, a) {
		attrs = append(attrs, attribute{Type: o.Type, Values: asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagSet, IsCompound: true, Bytes: o.Values}})
	}
	set, err := setOf(attrs)
	if err != nil {
		return nil, err
	}
	si.UnsignedAttrs = asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: 1, IsCompound: true, Bytes: rawContent(set)}
	// (the signed attributes as they were, byte for byte)
	si.SignedAttrs = asn1.RawValue{FullBytes: si.SignedAttrs.FullBytes}
	sib, err := asn1.Marshal(si)
	if err != nil {
		return nil, err
	}
	sd.SignerInfos = asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagSet, IsCompound: true, Bytes: append(sib, rest...)}
	sd.DigestAlgorithms = asn1.RawValue{FullBytes: sd.DigestAlgorithms.FullBytes}
	sd.EncapContentInfo.EContent = asn1.RawValue{FullBytes: sd.EncapContentInfo.EContent.FullBytes}
	if len(sd.Certificates.FullBytes) > 0 {
		sd.Certificates = asn1.RawValue{FullBytes: sd.Certificates.FullBytes}
	}
	if len(sd.CRLs.FullBytes) > 0 {
		sd.CRLs = asn1.RawValue{FullBytes: sd.CRLs.FullBytes}
	}
	sdb, err := asn1.Marshal(sd)
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(contentInfo{ContentType: OIDSignedData, Content: asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: 0, IsCompound: true, Bytes: sdb}})
}
