package bdf

import (
	"bytes"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestKeyWrapRFC3394(t *testing.T) {
	// RFC 3394 §4.6: 256 bits of key data with a 256-bit KEK.
	kek, _ := hex.DecodeString("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F")
	key, _ := hex.DecodeString("00112233445566778899AABBCCDDEEFF000102030405060708090A0B0C0D0E0F")
	want, _ := hex.DecodeString("28C9F404C4B810F4CBCCB35CFB87F8263F5786E2D80ED326CBC7F0E71A99F43BFB988B9B7A02DD21")
	got, err := wrapKey(kek, key)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("wrap = %x %v", got, err)
	}
	back, err := unwrapKey(kek, got)
	if err != nil || !bytes.Equal(back, key) {
		t.Fatalf("unwrap = %x %v", back, err)
	}
	got[3] ^= 1
	if _, err := unwrapKey(kek, got); err == nil {
		t.Fatal("tampered wrap accepted")
	}
}

// encryptedSample is a document with a title, two views' worth of parts and
// a part big enough to be compressed.
func encryptedSample(t *testing.T, password string) (*Document, Hash, Hash) {
	t.Helper()
	d := NewDocument()
	d.Meta.DC.Title = DCValues{"secret title"}
	img := d.AddImage([]byte("not really a png, but opaque bytes"))
	big := NewObject()
	for i := range 200 {
		big.FillRect(float32(i), 0, 1, 1)
	}
	h, _ := d.AddObject(big)
	v := d.NewView("v", ViewFixed, "")
	v.AddPage(100, 100, Layer{Role: RoleBody, Obj: h})
	lock, err := NewPasswordLock(password, 1000)
	if err != nil {
		t.Fatal(err)
	}
	d.Lock = lock
	return d, h, img
}

func TestEncryptedSingle(t *testing.T) {
	d, h, img := encryptedSample(t, "pass")
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	file := b.Bytes()
	if bytes.Contains(file, []byte("secret title")) || bytes.Contains(file, []byte("opaque bytes")) || bytes.Contains(file, []byte(h.String())) {
		t.Fatal("plaintext or part hashes stored in the clear")
	}
	if file[6] != FlagEncrypted {
		t.Fatalf("flags = %d", file[6])
	}
	r, err := OpenSingle(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		t.Fatal(err)
	}
	if !r.Encrypted() || !r.Locked() || len(r.Manifest.Views) != 0 {
		t.Fatalf("locked reader: encrypted %v locked %v views %d", r.Encrypted(), r.Locked(), len(r.Manifest.Views))
	}
	if _, err := r.Object(h); err == nil {
		t.Fatal("inner part readable before unlock")
	}
	if _, err := r.ToDocument(); !errors.Is(err, ErrLocked) {
		t.Fatalf("ToDocument before unlock: %v", err)
	}
	if err := r.Unlock("wrong"); !errors.Is(err, ErrWrongPassword) {
		t.Fatalf("wrong password: %v", err)
	}
	if err := r.Unlock("pass"); err != nil {
		t.Fatal(err)
	}
	if r.Locked() || r.Manifest.Meta.DC.Title.First() != "secret title" || r.Manifest.Views[0].Pages[0].Layers[0].Obj != h {
		t.Fatalf("unlocked manifest = %+v", r.Manifest)
	}
	e, _ := r.Entry(h)
	if e.Enc != EncDeflateRaw || e.Sealed == (Hash{}) {
		t.Fatalf("entry = %+v", e)
	}
	if got, err := r.Part(img); err != nil || string(got) != "not really a png, but opaque bytes" {
		t.Fatalf("image part = %q %v", got, err)
	}
	if _, err := r.Object(h); err != nil {
		t.Fatal(err)
	}

	// An unlocked document converts back to a plain one.
	d2, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	var plain bytes.Buffer
	if err := d2.WriteSingle(&plain); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(plain.Bytes(), []byte("secret title")) {
		t.Fatal("decrypted copy lost its manifest")
	}
}

func TestEncryptedTamper(t *testing.T) {
	d, h, _ := encryptedSample(t, "pass")
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	file := b.Bytes()
	r, err := OpenSingle(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Unlock("pass"); err != nil {
		t.Fatal(err)
	}
	e, _ := r.Entry(h)
	o := r.outer[e.Sealed]
	hdr := &reader{b: file[8:24]}
	base := int64(hdr.u64()) + int64(hdr.u64())
	file[base+o.Off+20] ^= 0x40
	if _, err := r.Part(h); err == nil {
		t.Fatal("tampered part accepted")
	}

	// The header flag must agree with the manifest.
	var c bytes.Buffer
	d.WriteSingle(&c)
	f2 := c.Bytes()
	f2[6] = 0
	if _, err := OpenSingle(bytes.NewReader(f2), int64(len(f2))); err == nil {
		t.Fatal("missing encryption flag accepted")
	}
}

func TestEncryptedSplitSwap(t *testing.T) {
	d, h, img := encryptedSample(t, "pass")
	dir := t.TempDir()
	if err := d.WriteSplit(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "parts", h.String())); err == nil {
		t.Fatal("split part named by its content hash")
	}
	r, err := OpenSplit(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Unlock("pass"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Object(h); err != nil {
		t.Fatal(err)
	}
	// Swapping two sealed parts' files is caught: parts are sealed with
	// their hash as additional data.
	a, _ := r.Entry(h)
	b, _ := r.Entry(img)
	pa := filepath.Join(dir, "parts", a.Sealed.String())
	pb := filepath.Join(dir, "parts", b.Sealed.String())
	ba, _ := os.ReadFile(pa)
	bb, _ := os.ReadFile(pb)
	os.WriteFile(pa, bb, 0o644)
	os.WriteFile(pb, ba, 0o644)
	if _, err := r.Part(img); err == nil {
		t.Fatal("swapped part accepted")
	}
}

func TestEncryptedSplitJoinWithoutPassword(t *testing.T) {
	d, h, _ := encryptedSample(t, "pass")
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	r, err := OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := r.WriteSplit(dir); err != nil {
		t.Fatal(err)
	}
	r2, err := OpenSplit(dir)
	if err != nil {
		t.Fatal(err)
	}
	var b2 bytes.Buffer
	if err := r2.WriteSingle(&b2); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b.Bytes(), b2.Bytes()) {
		t.Fatal("single -> split -> single is not byte-identical")
	}
	if err := r2.Unlock("pass"); err != nil {
		t.Fatal(err)
	}
	if _, err := r2.Object(h); err != nil {
		t.Fatal(err)
	}

	// A corrupted sealed part is caught by its name when copying.
	e := r2.stored.Parts[1]
	p := filepath.Join(dir, "parts", e.H.String())
	bs, _ := os.ReadFile(p)
	bs[0] ^= 1
	os.WriteFile(p, bs, 0o644)
	if err := r2.WriteSingle(&bytes.Buffer{}); err == nil {
		t.Fatal("corrupted part copied")
	}
}

func TestPasswordNormalization(t *testing.T) {
	// ガ precomposed (NFC) and as カ with a combining voiced mark (NFD).
	d, _, _ := encryptedSample(t, "パスワードガ")
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	r, err := OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Unlock("パスワードガ"); err != nil {
		t.Fatalf("NFD password: %v", err)
	}
	if _, err := NewPasswordLock("", 0); err == nil {
		t.Fatal("empty password accepted")
	}
}

func TestPlainSplitJoinCopy(t *testing.T) {
	d := NewDocument()
	big := NewObject()
	for i := range 200 {
		big.FillRect(float32(i), 0, 1, 1)
	}
	h, _ := d.AddObject(big)
	d.NewView("v", ViewFixed, "").AddPage(100, 100, Layer{Role: RoleBody, Obj: h})
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	r, err := OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := r.WriteSplit(dir); err != nil {
		t.Fatal(err)
	}
	r2, err := OpenSplit(dir)
	if err != nil {
		t.Fatal(err)
	}
	var b2 bytes.Buffer
	if err := r2.WriteSingle(&b2); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b.Bytes(), b2.Bytes()) {
		t.Fatal("copying single -> split -> single is not byte-identical")
	}
}
