package cms

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"
)

// The time-stamps of RFC 3161: a time-stamp authority (TSA) signs the
// digest of a signature and the time it saw it, so that a signature is
// known to have been made before then.

type messageImprint struct {
	HashAlgorithm pkix.AlgorithmIdentifier
	HashedMessage []byte
}

type timeStampReq struct {
	Version        int
	MessageImprint messageImprint
	Nonce          *big.Int `asn1:"optional"`
	CertReq        bool     `asn1:"optional"`
}

type pkiStatusInfo struct {
	Status       int
	StatusString asn1.RawValue  `asn1:"optional"`
	FailInfo     asn1.BitString `asn1:"optional"`
}

type timeStampResp struct {
	Status         pkiStatusInfo
	TimeStampToken asn1.RawValue `asn1:"optional"`
}

type accuracy struct {
	Seconds int `asn1:"optional"`
	Millis  int `asn1:"optional,tag:0"`
	Micros  int `asn1:"optional,tag:1"`
}

type tstInfo struct {
	Version        int
	Policy         asn1.ObjectIdentifier
	MessageImprint messageImprint
	SerialNumber   *big.Int
	GenTime        time.Time     `asn1:"generalized"`
	Accuracy       accuracy      `asn1:"optional"`
	Ordering       bool          `asn1:"optional"`
	Nonce          *big.Int      `asn1:"optional"`
	TSA            asn1.RawValue `asn1:"optional,tag:0"`
	Extensions     asn1.RawValue `asn1:"optional,tag:1"`
}

// TimeStamp is a time-stamp token that was read.
type TimeStamp struct {
	// Time is when the TSA stamped the digest, Digest the digest, Hash
	// its algorithm, and Serial the number of the token.
	Time   time.Time
	Digest []byte
	Hash   crypto.Hash
	Serial *big.Int
	// Token is the token: a SignedData of the TSA.
	Token *SignedData
	// DER is the DER of the token.
	DER []byte
}

// TimeStampRequest makes the DER of a request of a time-stamp of a digest,
// with a nonce that the token must hold.
func TimeStampRequest(digest []byte, h crypto.Hash) (req []byte, nonce *big.Int, err error) {
	if !h.Available() || len(digest) != h.Size() {
		return nil, nil, fmt.Errorf("cms: a digest of %d bytes is not one of %v", len(digest), h)
	}
	nonce, err = rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 64))
	if err != nil {
		return nil, nil, err
	}
	req, err = asn1.Marshal(timeStampReq{Version: 1, MessageImprint: messageImprint{HashAlgorithm: pkix.AlgorithmIdentifier{Algorithm: hashOID(h), Parameters: asn1.NullRawValue}, HashedMessage: digest}, Nonce: nonce, CertReq: true})
	return req, nonce, err
}

// maxTimeStampResponse bounds the response of a TSA.
const maxTimeStampResponse = 1 << 20

// RequestTimeStamp asks a TSA at a URL (HTTP, RFC 3161 §3.4) for a
// time-stamp token of a digest, and checks it (ReadTimeStamp): the digest
// it stamps and the nonce of the request. client nil is http.DefaultClient.
func RequestTimeStamp(ctx context.Context, client *http.Client, url string, digest []byte, h crypto.Hash) (*TimeStamp, error) {
	req, nonce, err := TimeStampRequest(digest, h)
	if err != nil {
		return nil, err
	}
	if client == nil {
		client = http.DefaultClient
	}
	hr, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(req))
	if err != nil {
		return nil, err
	}
	hr.Header.Set("Content-Type", "application/timestamp-query")
	res, err := client.Do(hr)
	if err != nil {
		return nil, fmt.Errorf("cms: the time-stamp authority: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cms: the time-stamp authority answered %s", res.Status)
	}
	body, err := io.ReadAll(io.LimitReader(res.Body, maxTimeStampResponse+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxTimeStampResponse {
		return nil, errors.New("cms: the answer of the time-stamp authority is too large")
	}
	ts, err := readResponse(body)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(ts.Digest, digest) || ts.Hash != h {
		return nil, errors.New("cms: the time-stamp is of another digest")
	}
	if ts.nonce == nil || ts.nonce.Cmp(nonce) != 0 {
		return nil, errors.New("cms: the time-stamp is not the answer to the request (its nonce differs)")
	}
	return &ts.TimeStamp, nil
}

type readStamp struct {
	TimeStamp
	nonce *big.Int
}

// ReadTimeStampResponse reads the DER of a TimeStampResp: the token, when
// the TSA granted it.
func ReadTimeStampResponse(der []byte) (*TimeStamp, error) {
	ts, err := readResponse(der)
	if err != nil {
		return nil, err
	}
	return &ts.TimeStamp, nil
}

func readResponse(der []byte) (*readStamp, error) {
	var resp timeStampResp
	if _, err := asn1.Unmarshal(der, &resp); err != nil {
		return nil, fmt.Errorf("cms: the answer of the time-stamp authority: %w", err)
	}
	if resp.Status.Status != 0 && resp.Status.Status != 1 {
		return nil, fmt.Errorf("cms: the time-stamp authority refused (status %d)", resp.Status.Status)
	}
	if len(resp.TimeStampToken.FullBytes) == 0 {
		return nil, errors.New("cms: the time-stamp authority gave no token")
	}
	return readToken(resp.TimeStampToken.FullBytes)
}

// ReadTimeStamp reads and checks a time-stamp token: that its TSA signed
// it (the certificate of the TSA is not checked here: VerifyChain of the
// token with x509.ExtKeyUsageTimeStamping).
func ReadTimeStamp(der []byte) (*TimeStamp, error) {
	ts, err := readToken(der)
	if err != nil {
		return nil, err
	}
	return &ts.TimeStamp, nil
}

func readToken(der []byte) (*readStamp, error) {
	sd, err := Parse(der)
	if err != nil {
		return nil, fmt.Errorf("cms: the time-stamp token: %w", err)
	}
	if !sd.ContentType.Equal(OIDTSTInfo) || sd.Content == nil {
		return nil, errors.New("cms: the time-stamp token holds no TSTInfo")
	}
	if err := sd.Verify(sd.Signers[0], nil); err != nil {
		return nil, fmt.Errorf("cms: the time-stamp token: %w", err)
	}
	var info tstInfo
	if _, err := asn1.Unmarshal(sd.Content, &info); err != nil {
		return nil, fmt.Errorf("cms: the TSTInfo: %w", err)
	}
	h, ok := hashOf(info.MessageImprint.HashAlgorithm.Algorithm)
	if !ok {
		return nil, fmt.Errorf("cms: the time-stamp is of the digest %v, which is not known", info.MessageImprint.HashAlgorithm.Algorithm)
	}
	return &readStamp{TimeStamp: TimeStamp{Time: info.GenTime, Digest: info.MessageImprint.HashedMessage, Hash: h, Serial: info.SerialNumber, Token: sd, DER: der}, nonce: info.Nonce}, nil
}

// TimeStampToken makes a token of a TSA (for tests, and for a TSA of one's
// own): the TSTInfo of a digest at a time, signed with the key of the TSA.
func TimeStampToken(digest []byte, h crypto.Hash, at time.Time, serial, nonce *big.Int, policy asn1.ObjectIdentifier, key crypto.Signer, chain []*x509.Certificate) ([]byte, error) {
	info, err := asn1.Marshal(tstInfo{Version: 1, Policy: policy, MessageImprint: messageImprint{HashAlgorithm: pkix.AlgorithmIdentifier{Algorithm: hashOID(h), Parameters: asn1.NullRawValue}, HashedMessage: digest},
		SerialNumber: serial, GenTime: at.UTC().Truncate(time.Second), Nonce: nonce})
	if err != nil {
		return nil, err
	}
	return Sign(info, key, chain, &SignOptions{ContentType: OIDTSTInfo})
}

// TimeStampResponse wraps a token in a granted TimeStampResp.
func TimeStampResponse(token []byte) ([]byte, error) {
	return asn1.Marshal(timeStampResp{Status: pkiStatusInfo{Status: 0}, TimeStampToken: asn1.RawValue{FullBytes: token}})
}

// ParseTimeStampRequest reads a request of a time-stamp: its digest and
// nonce (for a TSA of one's own).
func ParseTimeStampRequest(der []byte) (digest []byte, h crypto.Hash, nonce *big.Int, err error) {
	var req timeStampReq
	if _, err := asn1.Unmarshal(der, &req); err != nil {
		return nil, 0, nil, fmt.Errorf("cms: the request of a time-stamp: %w", err)
	}
	h, ok := hashOf(req.MessageImprint.HashAlgorithm.Algorithm)
	if !ok || len(req.MessageImprint.HashedMessage) != h.Size() {
		return nil, 0, nil, errors.New("cms: the request is of a digest that is not known")
	}
	return req.MessageImprint.HashedMessage, h, req.Nonce, nil
}

// TimeStampOf returns the time-stamp token of a signer (its unsigned
// signature-time-stamp attribute), checked against the signature: nil
// without one.
func (s *Signer) TimeStampOf() (*TimeStamp, error) {
	b := s.Unsigned(OIDTimeStampToken)
	if b == nil {
		return nil, nil
	}
	ts, err := ReadTimeStamp(b)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(ts.Digest, s.SignatureDigest(ts.Hash)) {
		return nil, errors.New("cms: the time-stamp is of another signature")
	}
	return ts, nil
}
