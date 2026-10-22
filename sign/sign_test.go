package sign

import (
	"bytes"
	"crypto/x509"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/cms/cmstest"
)

func demo(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("../testdata/demo.bdf")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func read(t *testing.T, b []byte) *Document {
	t.Helper()
	d, err := Read(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// signed returns a document with a signature of id added.
func signed(t *testing.T, b []byte, id *cmstest.Identity, o *Options) []byte {
	t.Helper()
	d := read(t, b)
	sig, err := d.Sign(id.Key, id.Chain, o)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := d.WriteSingle(&out, sig); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

func verify(t *testing.T, b []byte, roots *x509.CertPool) []*Signature {
	t.Helper()
	sigs, err := read(t, b).Verify(&VerifyOptions{Roots: roots})
	if err != nil {
		t.Fatal(err)
	}
	return sigs
}

// TestSign signs a document twice and checks what changes after it break.
func TestSign(t *testing.T) {
	ca := cmstest.NewAuthority("BDF test root")
	alice, bob := ca.Issue("Alice", false), ca.Issue("Bob", true)
	doc := demo(t)
	one := signed(t, doc, alice, &Options{Claims: map[string]any{"role": "author"}})
	two := signed(t, one, bob, &Options{Claims: map[string]any{"role": "counterparty", "note": "同意します"}})

	// the document opens as it did: the signatures are after it
	for _, b := range [][]byte{one, two} {
		rd, err := bdf.OpenSingle(bytes.NewReader(b), int64(len(b)))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := rd.ToDocument(); err != nil {
			t.Fatal(err)
		}
	}
	if d := read(t, two); d.Size != int64(len(doc)) || len(d.Signatures) != 2 {
		t.Fatalf("the signed document: %d bytes of %d, %d signatures", d.Size, len(doc), len(d.Signatures))
	}
	sigs := verify(t, two, ca.Pool())
	if len(sigs) != 2 || !sigs[0].Valid() || !sigs[1].Valid() {
		t.Fatalf("signatures: %+v %+v", sigs[0], sigs[1])
	}
	if sigs[0].Certificate.Subject.CommonName != "Alice" || sigs[1].Claims["note"] != "同意します" || sigs[1].Claims["role"] != "counterparty" || time.Since(sigs[1].SigningTime) > time.Minute {
		t.Errorf("what the signatures say: %+v / %+v", sigs[0], sigs[1])
	}
	if len(sigs[1].Payload.Previous) != 1 || len(sigs[0].Payload.Previous) != 0 {
		t.Errorf("previous: %v %v", sigs[0].Payload.Previous, sigs[1].Payload.Previous)
	}

	// not trusted: the signatures hold, the certificates are not trusted
	for _, s := range verify(t, two, cmstest.NewAuthority("another").Pool()) {
		if s.Err != nil || s.TrustErr == nil || s.Valid() {
			t.Errorf("another root: %v %v", s.Err, s.TrustErr)
		}
	}

	// a byte of a part changed
	bad := bytes.Clone(two)
	bad[len(doc)-10] ^= 0xff
	for _, s := range verify(t, bad, ca.Pool()) {
		if s.Err == nil || !strings.Contains(s.Err.Error(), "part(s) changed") {
			t.Errorf("a changed part: %v", s.Err)
		}
	}

	// the manifest changed: another title, with the signatures of before
	rd, _ := bdf.OpenSingle(bytes.NewReader(doc), int64(len(doc)))
	d2, _ := rd.ToDocument()
	d2.Meta.Generator = "another generator"
	var retitled bytes.Buffer
	if err := d2.WriteSingle(&retitled); err != nil {
		t.Fatal(err)
	}
	orig := read(t, two)
	writeBlock(&retitled, orig.Signatures)
	for _, s := range verify(t, retitled.Bytes(), ca.Pool()) {
		if s.Err == nil || !strings.Contains(s.Err.Error(), "the manifest changed") {
			t.Errorf("a changed manifest: %v", s.Err)
		}
	}

	// the first signature taken away: the second says there was one
	var taken bytes.Buffer
	taken.Write(doc)
	writeBlock(&taken, orig.Signatures[1:])
	if s := verify(t, taken.Bytes(), ca.Pool()); s[0].Err == nil || !strings.Contains(s[0].Err.Error(), "1 signature(s) came before") {
		t.Errorf("a signature taken away: %v", s[0].Err)
	}
	// the signatures in another order
	var swapped bytes.Buffer
	swapped.Write(doc)
	writeBlock(&swapped, [][]byte{orig.Signatures[1], orig.Signatures[0]})
	if s := verify(t, swapped.Bytes(), ca.Pool()); s[0].Err == nil || s[1].Err == nil {
		t.Errorf("signatures in another order: %v / %v", s[0].Err, s[1].Err)
	}
	// a damaged block
	cut := bytes.Clone(two)
	cut[len(cut)-9] ^= 0x40
	if _, err := Read(bytes.NewReader(cut), int64(len(cut))); err == nil {
		t.Error("a damaged block of signatures was read")
	}
	// a document without signatures
	if d := read(t, doc); d.Size != int64(len(doc)) || len(d.Signatures) != 0 {
		t.Errorf("an unsigned document: %d %d", d.Size, len(d.Signatures))
	}
}

// TestSplit signs the split form, and checks that a signature of the
// single form holds for the split form of the same document.
func TestSplit(t *testing.T) {
	ca := cmstest.NewAuthority("BDF test root")
	alice := ca.Issue("Alice", false)
	one := signed(t, demo(t), alice, nil)
	d := read(t, one)
	dir := filepath.Join(t.TempDir(), "doc")
	if err := d.Reader.WriteSplit(dir); err != nil {
		t.Fatal(err)
	}
	split, err := ReadSplit(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(split.Signatures) != 0 {
		t.Fatal("signatures before they are written")
	}
	if err := split.WriteSplit(d.Signatures...); err != nil {
		t.Fatal(err)
	}
	split, _ = ReadSplit(dir)
	bob := ca.Issue("Bob", false)
	sig, err := split.Sign(bob.Key, bob.Chain, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := split.WriteSplit(sig); err != nil {
		t.Fatal(err)
	}
	split, _ = ReadSplit(dir)
	sigs, err := split.Verify(&VerifyOptions{Roots: ca.Pool()})
	if err != nil {
		t.Fatal(err)
	}
	if len(sigs) != 2 || !sigs[0].Valid() || !sigs[1].Valid() {
		t.Fatalf("the split form: %v %v", sigs[0].Err, sigs[1].Err)
	}
	// a part file changed
	parts, _ := filepath.Glob(filepath.Join(dir, "parts", "*"))
	b, _ := os.ReadFile(parts[0])
	b[len(b)/2] ^= 1
	os.WriteFile(parts[0], b, 0o644)
	split, _ = ReadSplit(dir)
	if sigs, err := split.Verify(&VerifyOptions{Roots: ca.Pool()}); err != nil || sigs[0].Err == nil {
		t.Errorf("a changed part file: %v", err)
	}
}

// TestEncrypted signs an encrypted document without its password, and
// verifies it without the password.
func TestEncrypted(t *testing.T) {
	ca := cmstest.NewAuthority("BDF test root")
	rd, _ := bdf.OpenSingle(bytes.NewReader(demo(t)), int64(len(demo(t))))
	doc, _ := rd.ToDocument()
	lock, err := bdf.NewPasswordLock("secret", 1000)
	if err != nil {
		t.Fatal(err)
	}
	doc.Lock = lock
	var enc bytes.Buffer
	if err := doc.WriteSingle(&enc); err != nil {
		t.Fatal(err)
	}
	one := signed(t, enc.Bytes(), ca.Issue("Alice", false), nil)
	if s := verify(t, one, ca.Pool()); !s[0].Valid() {
		t.Errorf("an encrypted document: %v %v", s[0].Err, s[0].TrustErr)
	}
	d := read(t, one)
	if err := d.Reader.Unlock("secret"); err != nil {
		t.Errorf("unlocking the signed document: %v", err)
	}
}

// TestTimeStamp signs with a time-stamp of a TSA, which certificates are
// checked at.
func TestTimeStamp(t *testing.T) {
	ca := cmstest.NewAuthority("BDF test root")
	tsa := cmstest.NewTSA(ca)
	defer tsa.Close()
	stamped := time.Now().Add(-30 * time.Minute).Truncate(time.Second)
	tsa.Now = func() time.Time { return stamped }
	one := signed(t, demo(t), ca.Issue("Alice", false), &Options{TSA: tsa.URL})
	s := verify(t, one, ca.Pool())[0]
	if !s.Valid() || !s.TimeStamp.Equal(stamped) || s.TSA.Subject.CommonName != "BDF test TSA" {
		t.Errorf("the time-stamp: %v %v %v %v", s.Err, s.TrustErr, s.TimeStamp, s.TSA)
	}
	if tsa.Requests.Load() != 1 {
		t.Errorf("%d requests", tsa.Requests.Load())
	}
	if _, err := read(t, demo(t)).Sign(ca.Issue("Bob", false).Key, ca.Issue("Bob", false).Chain, &Options{TSA: tsa.URL + "/"}); err == nil || !strings.Contains(err.Error(), "not that of the key") {
		t.Errorf("a certificate of another key: %v", err)
	}
}

// TestCanonical checks the canonical JSON with the examples of RFC 8785.
func TestCanonical(t *testing.T) {
	for in, want := range map[string]string{
		`{"numbers":[333333333.33333329,1E30,4.50,2e-3,0.000000000000000000000000001]}`:                                                                                                                              `{"numbers":[333333333.3333333,1e+30,4.5,0.002,1e-27]}`,
		`[0, -0, 1e21, 1e20, 0.000001, 1e-7, 9007199254740992, 295147905179352830000, -1.5e-10, 100, 123.456]`:                                                                                                       `[0,0,1e+21,100000000000000000000,0.000001,1e-7,9007199254740992,295147905179352830000,-1.5e-10,100,123.456]`,
		`{"\u20ac":"Euro Sign","\r":"Carriage Return","\ufb33":"Hebrew Letter Dalet With Dagesh","1":"One","\ud83d\ude00":"Emoji: Grinning Face","\u0080":"Control","\u00f6":"Latin Small Letter O With Diaeresis"}`: "{\"\\r\":\"Carriage Return\",\"1\":\"One\",\"\u0080\":\"Control\",\"ö\":\"Latin Small Letter O With Diaeresis\",\"€\":\"Euro Sign\",\"😀\":\"Emoji: Grinning Face\",\"\ufb33\":\"Hebrew Letter Dalet With Dagesh\"}",
		`{"string":"\u20ac$\u000F\u000aA'\u0042\u0022\u005c\\\"\/"}`:                                                                                                                                                 `{"string":"€$\u000f\nA'B\"\\\\\"/"}`,
		` { "b" : [ true , false , null ] , "a" : { } } `:                                                                                                                                                            `{"a":{},"b":[true,false,null]}`,
	} {
		got, err := Canonical([]byte(in))
		if err != nil || string(got) != want {
			t.Errorf("%s\n got %s (%v)\nwant %s", in, got, err, want)
		}
	}
	for _, in := range []string{`{"a":1} {}`, `[1e400]`, `{`} {
		if _, err := Canonical([]byte(in)); err == nil {
			t.Errorf("%s: no error", in)
		}
	}
}

// FuzzRead reads untrusted blocks of signatures.
func FuzzRead(f *testing.F) {
	ca := cmstest.NewAuthority("root")
	b, _ := os.ReadFile("../testdata/demo.bdf")
	d, _ := Read(bytes.NewReader(b), int64(len(b)))
	id := ca.Issue("x", false)
	sig, _ := d.Sign(id.Key, id.Chain, nil)
	var out bytes.Buffer
	d.WriteSingle(&out, sig)
	signed := out.Bytes()
	f.Add(signed[len(b):])
	f.Fuzz(func(t *testing.T, tail []byte) {
		all := append(bytes.Clone(b), tail...)
		d, err := Read(bytes.NewReader(all), int64(len(all)))
		if err != nil {
			return
		}
		d.Verify(&VerifyOptions{Roots: ca.Pool()})
	})
}
