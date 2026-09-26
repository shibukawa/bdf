package offcrypto

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

// password of the test decks (test/pptx/gen_encrypted.py).
const password = "パスワード🔑bdf"

func read(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestDecrypt(t *testing.T) {
	plain := read(t, "../../pptx/testdata/basic.pptx")
	for _, name := range []string{"agile.pptx", "standard.pptx"} {
		t.Run(name, func(t *testing.T) {
			b := read(t, "testdata/"+name)
			r := bytes.NewReader(b)
			if !IsEncrypted(r, r.Size()) {
				t.Fatal("not recognized as encrypted")
			}
			if err := Verify(r, r.Size(), password); err != nil {
				t.Fatal(err)
			}
			for _, pw := range []string{"", "wrong", "パスワード🔑bd"} {
				if err := Verify(r, r.Size(), pw); !errors.Is(err, ErrWrongPassword) {
					t.Fatalf("Verify(%q) = %v", pw, err)
				}
				if _, err := Decrypt(r, r.Size(), pw); !errors.Is(err, ErrWrongPassword) {
					t.Fatalf("Decrypt(%q) = %v", pw, err)
				}
			}
			got, err := Decrypt(r, r.Size(), password)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, plain) {
				t.Fatal("decrypted package differs from the original")
			}
		})
	}
	r := bytes.NewReader(plain)
	if IsEncrypted(r, r.Size()) {
		t.Fatal("plain package recognized as encrypted")
	}
}

func TestIntegrity(t *testing.T) {
	b := read(t, "testdata/agile.pptx")
	// The package is the last stream: flip a bit of its last segment.
	b[len(b)-600] ^= 0x10
	r := bytes.NewReader(b)
	got, err := Decrypt(r, r.Size(), password)
	if !errors.Is(err, ErrIntegrity) || len(got) == 0 {
		t.Fatalf("tampered package: %d bytes, %v", len(got), err)
	}
}
