package offcrypto

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
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

func TestAgileInvalidSaltSize(t *testing.T) {
	// Keep the compound stream lengths intact while replacing the password
	// salt size with a negative value. Verify and Decrypt must return an error
	// before deriving a key, rather than panicking when slicing the verifier.
	b := read(t, "testdata/agile.pptx")
	old := []byte(`saltSize="16"`)
	at := bytes.LastIndex(b, old)
	if at < 0 {
		t.Fatal("password salt size not found in fixture")
	}
	copy(b[at:], []byte(`saltSize="-1"`))
	r := bytes.NewReader(b)
	if err := Verify(r, r.Size(), password); err == nil || errors.Is(err, ErrWrongPassword) {
		t.Fatalf("Verify invalid salt size: %v", err)
	}
	if _, err := Decrypt(r, r.Size(), password); err == nil || errors.Is(err, ErrWrongPassword) {
		t.Fatalf("Decrypt invalid salt size: %v", err)
	}
}

func TestAgileSaltMetadata(t *testing.T) {
	for _, tc := range []struct {
		name string
		size int
		data int
		ok   bool
	}{
		{"minimum", 1, 1, true},
		{"usual", 16, 16, true},
		{"maximum", 65536, 65536, true},
		{"negative", -1, 16, false},
		{"zero", 0, 16, false},
		{"too large", 65537, 65537, false},
		{"truncated", 16, 15, false},
		{"extra bytes", 16, 17, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, part := range []string{"keyData", "encryptedKey"} {
				params := func(size, length int) string {
					return fmt.Sprintf(`saltSize="%d" blockSize="16" keyBits="128" hashSize="20" cipherAlgorithm="AES" cipherChaining="ChainingModeCBC" hashAlgorithm="SHA1" saltValue="%s"`, size, base64.StdEncoding.EncodeToString(make([]byte, length)))
				}
				keyData, encryptedKey := params(16, 16), params(16, 16)
				if part == "keyData" {
					keyData = params(tc.size, tc.data)
				} else {
					encryptedKey = params(tc.size, tc.data)
				}
				xml := `<encryption><keyData ` + keyData + `/><keyEncryptors><keyEncryptor uri="` + passwordEncryptor + `"><encryptedKey ` + encryptedKey + `/></keyEncryptor></keyEncryptors></encryption>`
				if _, err := parseAgile([]byte(xml)); (err == nil) != tc.ok {
					t.Errorf("%s: parseAgile error = %v, want success %v", part, err, tc.ok)
				}
			}
		})
	}
}
