package cms_test

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/internal/cms"
	"github.com/shibukawa/bdf/internal/cms/cmstest"
)

// TestSignVerify signs with ECDSA and RSA keys, attached and detached, and
// checks what a change to the content or the signature does.
func TestSignVerify(t *testing.T) {
	ca := cmstest.NewAuthority("BDF test root")
	content := []byte(`{"bdf-signature":1}`)
	at := time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)
	for _, rsaKey := range []bool{false, true} {
		id := ca.Issue("Taro Yamada", rsaKey)
		for _, detached := range []bool{false, true} {
			der, err := cms.Sign(content, id.Key, id.Chain, &cms.SignOptions{Detached: detached, SigningTime: at})
			if err != nil {
				t.Fatal(err)
			}
			sd, err := cms.Parse(der)
			if err != nil {
				t.Fatal(err)
			}
			if detached != (sd.Content == nil) || len(sd.Certificates) != 2 || len(sd.Signers) != 1 {
				t.Fatalf("rsa %v detached %v: %+v", rsaKey, detached, sd)
			}
			s := sd.Signers[0]
			var with []byte
			if detached {
				with = content
			}
			if err := sd.Verify(s, with); err != nil {
				t.Errorf("rsa %v detached %v: %v", rsaKey, detached, err)
			}
			if !s.SigningTime().Equal(at) || s.Certificate.Subject.CommonName != "Taro Yamada" {
				t.Errorf("the signer: %v %v", s.SigningTime(), s.Certificate.Subject)
			}
			if _, err := sd.VerifyChain(s, ca.Pool(), at.Add(-time.Minute)); err == nil {
				t.Error("a certificate is valid before it was issued")
			}
			if _, err := sd.VerifyChain(s, ca.Pool(), time.Now()); err != nil {
				t.Errorf("the chain: %v", err)
			}
			if _, err := sd.VerifyChain(s, cmstest.NewAuthority("another").Pool(), time.Now()); err == nil {
				t.Error("a chain to another root")
			}
			if err := sd.Verify(s, []byte(`{"bdf-signature":2}`)); err == nil || !strings.Contains(err.Error(), "digest differs") {
				t.Errorf("another content: %v", err)
			}
			// a byte of the signature changed
			bad := bytes.Clone(der)
			i := bytes.Index(bad, s.Signature)
			bad[i+len(s.Signature)/2] ^= 1
			if sd2, err := cms.Parse(bad); err == nil {
				if err := sd2.Verify(sd2.Signers[0], with); err == nil {
					t.Error("a changed signature verifies")
				}
			}
		}
	}
	other := ca.Issue("other", false)
	if _, err := cms.Sign(content, other.Key, ca.Issue("x", false).Chain, nil); err == nil {
		t.Error("signing with the certificate of another key")
	}
}

// TestTimeStamp asks a time-stamp authority for a token of a signature,
// adds it, and reads it back.
func TestTimeStamp(t *testing.T) {
	ca := cmstest.NewAuthority("BDF test root")
	tsa := cmstest.NewTSA(ca)
	defer tsa.Close()
	stamped := time.Now().Add(-time.Minute).Truncate(time.Second).UTC()
	tsa.Now = func() time.Time { return stamped }
	id := ca.Issue("Taro Yamada", false)
	der, err := cms.Sign([]byte("content"), id.Key, id.Chain, nil)
	if err != nil {
		t.Fatal(err)
	}
	sd, _ := cms.Parse(der)
	ts, err := cms.RequestTimeStamp(context.Background(), nil, tsa.URL, sd.Signers[0].SignatureDigest(crypto.SHA256), crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	if !ts.Time.Equal(stamped) {
		t.Errorf("stamped at %v", ts.Time)
	}
	if _, err := ts.Token.VerifyChain(ts.Token.Signers[0], ca.Pool(), ts.Time, x509.ExtKeyUsageTimeStamping); err != nil {
		t.Errorf("the chain of the TSA: %v", err)
	}
	der2, err := cms.AddUnsigned(der, cms.Attribute{Type: cms.OIDTimeStampToken, Values: ts.DER})
	if err != nil {
		t.Fatal(err)
	}
	sd2, err := cms.Parse(der2)
	if err != nil {
		t.Fatal(err)
	}
	if err := sd2.Verify(sd2.Signers[0], nil); err != nil {
		t.Errorf("the signature after the time-stamp is added: %v", err)
	}
	got, err := sd2.Signers[0].TimeStampOf()
	if err != nil || got == nil || !got.Time.Equal(stamped) {
		t.Errorf("the time-stamp of the signer: %v %v", got, err)
	}
	// a token of another signature
	other, _ := cms.Sign([]byte("other"), id.Key, id.Chain, nil)
	der3, _ := cms.AddUnsigned(other, cms.Attribute{Type: cms.OIDTimeStampToken, Values: ts.DER})
	sd3, _ := cms.Parse(der3)
	if _, err := sd3.Signers[0].TimeStampOf(); err == nil {
		t.Error("the time-stamp of another signature")
	}
	sum := sha256.Sum256([]byte("x"))
	if _, err := cms.RequestTimeStamp(context.Background(), nil, tsa.URL, sum[:], crypto.SHA384); err == nil {
		t.Error("a digest of another length was stamped as SHA-384")
	}
}

// TestOpenSSL checks the signatures with OpenSSL, when it is installed.
func TestOpenSSL(t *testing.T) {
	bin, err := exec.LookPath("openssl")
	if err != nil {
		t.Skip("no openssl")
	}
	dir := t.TempDir()
	ca := cmstest.NewAuthority("BDF test root")
	write := func(name string, b []byte) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, b, 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	root, _ := (&cmstest.Identity{Key: ca.Key, Chain: []*x509.Certificate{ca.Cert}}).PEM()
	write("root.pem", root)
	content := write("content.json", []byte(`{"bdf-signature":1,"x":"テスト"}`))
	for _, rsaKey := range []bool{false, true} {
		id := ca.Issue("Signer", rsaKey)
		for _, detached := range []bool{false, true} {
			der, err := cms.Sign([]byte(`{"bdf-signature":1,"x":"テスト"}`), id.Key, id.Chain, &cms.SignOptions{Detached: detached, SigningTime: time.Now()})
			if err != nil {
				t.Fatal(err)
			}
			sig := write("sig.der", der)
			args := []string{"cms", "-verify", "-binary", "-inform", "DER", "-in", sig, "-CAfile", filepath.Join(dir, "root.pem"), "-purpose", "any", "-out", filepath.Join(dir, "out")}
			if detached {
				args = append(args, "-content", content)
			}
			if out, err := exec.Command(bin, args...).CombinedOutput(); err != nil {
				t.Errorf("rsa %v detached %v: openssl: %v\n%s", rsaKey, detached, err, out)
			}
		}
	}
	// and a time-stamp token of OpenSSL's own TSA, read here
	tsaID := ca.Issue("OpenSSL TSA", false, x509.ExtKeyUsageTimeStamping)
	certs, key := tsaID.PEM()
	write("tsa.pem", certs)
	write("tsa.key", key)
	write("tsa.cnf", []byte(`[ tsa ]
default_tsa = tsa_config
[ tsa_config ]
serial = `+filepath.Join(dir, "serial")+`
signer_cert = `+filepath.Join(dir, "tsa.pem")+`
certs = `+filepath.Join(dir, "tsa.pem")+`
signer_key = `+filepath.Join(dir, "tsa.key")+`
signer_digest = sha256
default_policy = 1.2.3.4.1
digests = sha256
ess_cert_id_alg = sha256
`))
	write("serial", []byte("01\n"))
	sum := sha256.Sum256([]byte("signature"))
	req, _, err := cms.TimeStampRequest(sum[:], crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	write("req.tsq", req)
	if out, err := exec.Command(bin, "ts", "-reply", "-config", filepath.Join(dir, "tsa.cnf"), "-queryfile", filepath.Join(dir, "req.tsq"), "-out", filepath.Join(dir, "resp.tsr")).CombinedOutput(); err != nil {
		t.Fatalf("openssl ts -reply: %v\n%s", err, out)
	}
	resp, _ := os.ReadFile(filepath.Join(dir, "resp.tsr"))
	ts, err := cms.ReadTimeStampResponse(resp)
	if err != nil {
		t.Fatalf("the answer of OpenSSL's TSA: %v", err)
	}
	if !bytes.Equal(ts.Digest, sum[:]) || time.Since(ts.Time) > time.Minute {
		t.Errorf("the time-stamp of OpenSSL: %v %x", ts.Time, ts.Digest)
	}
	if _, err := ts.Token.VerifyChain(ts.Token.Signers[0], ca.Pool(), ts.Time, x509.ExtKeyUsageTimeStamping); err != nil {
		t.Errorf("the chain of OpenSSL's TSA: %v", err)
	}
	// and the token of our TSA, verified by OpenSSL
	ours := cmstest.NewTSA(ca)
	defer ours.Close()
	got, err := cms.RequestTimeStamp(context.Background(), nil, ours.URL, sum[:], crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	reply, _ := cms.TimeStampResponse(got.DER)
	write("ours.tsr", reply)
	tsaCerts, _ := ours.Identity.PEM()
	write("ours-tsa.pem", tsaCerts)
	write("data", []byte("signature"))
	if out, err := exec.Command(bin, "ts", "-verify", "-in", filepath.Join(dir, "ours.tsr"), "-data", filepath.Join(dir, "data"), "-CAfile", filepath.Join(dir, "root.pem"), "-untrusted", filepath.Join(dir, "ours-tsa.pem")).CombinedOutput(); err != nil {
		t.Errorf("openssl ts -verify of our token: %v\n%s", err, out)
	}
}

// FuzzParse reads untrusted DER.
func FuzzParse(f *testing.F) {
	ca := cmstest.NewAuthority("root")
	id := ca.Issue("x", false)
	der, _ := cms.Sign([]byte("content"), id.Key, id.Chain, &cms.SignOptions{SigningTime: time.Now()})
	f.Add(der)
	f.Fuzz(func(t *testing.T, b []byte) {
		sd, err := cms.Parse(b)
		if err != nil {
			return
		}
		for _, s := range sd.Signers {
			sd.Verify(s, nil)
			s.TimeStampOf()
			s.SigningTime()
		}
	})
}
