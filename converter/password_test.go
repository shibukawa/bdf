package converter_test

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/pdf"
	_ "github.com/shibukawa/bdf/converter/pptx"
)

// officePassword opens the encrypted decks (test/pptx/gen_encrypted.py).
const officePassword = "パスワード🔑bdf"

func pptxOptions(password string) *converter.Options {
	return &converter.Options{FontDirs: []string{"pptx/testdata/fonts"}, NoSystemFonts: true, Password: password}
}

func convertBytes(b []byte, opts *converter.Options) (*converter.Result, error) {
	return converter.Convert(bytes.NewReader(b), int64(len(b)), "", opts)
}

func TestEncryptedOffice(t *testing.T) {
	plain, err := converter.ConvertFile("pptx/testdata/basic.pptx", "", pptxOptions(""))
	if err != nil {
		t.Fatal(err)
	}
	if plain.Protected {
		t.Fatal("plain deck reported protected")
	}
	for _, name := range []string{"agile.pptx", "standard.pptx"} {
		t.Run(name, func(t *testing.T) {
			b, err := os.ReadFile("internal/offcrypto/testdata/" + name)
			if err != nil {
				t.Fatal(err)
			}
			r := bytes.NewReader(b)
			if converter.Detect(r, r.Size()) != nil {
				t.Fatal("encrypted deck detected without its password")
			}
			if _, err := convertBytes(b, pptxOptions("")); !errors.Is(err, converter.ErrPasswordRequired) {
				t.Fatalf("no password: %v", err)
			}
			if _, err := convertBytes(b, pptxOptions("wrong")); !errors.Is(err, converter.ErrWrongPassword) {
				t.Fatalf("wrong password: %v", err)
			}
			for pw, want := range map[string]error{"": converter.ErrPasswordRequired, "wrong": converter.ErrWrongPassword, officePassword: nil} {
				if p, err := converter.CheckPassword(r, r.Size(), pw); !p || !errors.Is(err, want) {
					t.Fatalf("CheckPassword(%q) = %v, %v", pw, p, err)
				}
			}
			res, err := convertBytes(b, pptxOptions(officePassword))
			if err != nil {
				t.Fatal(err)
			}
			if !res.Protected || res.Summary != plain.Summary || len(res.Warnings) != len(plain.Warnings) {
				t.Fatalf("protected %v, %q (%d warnings), want %q (%d)", res.Protected, res.Summary, len(res.Warnings), plain.Summary, len(plain.Warnings))
			}
			// The format can be named too; decryption comes first.
			if _, err := converter.Convert(r, r.Size(), "pptx", pptxOptions(officePassword)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// encryptPDF encrypts a test PDF with pdfcpu.
func encryptPDF(t *testing.T, user, owner string, perms model.PermissionFlags) []byte {
	t.Helper()
	src, err := os.ReadFile("pdf/testdata/reportlab-mixed.pdf")
	if err != nil {
		t.Fatal(err)
	}
	conf := model.NewAESConfiguration(user, owner, 256)
	conf.Permissions = perms
	var out bytes.Buffer
	if err := api.Encrypt(bytes.NewReader(src), &out, conf); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

func TestEncryptedPDF(t *testing.T) {
	plain, err := converter.ConvertFile("pdf/testdata/reportlab-mixed.pdf", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	b := encryptPDF(t, "user-pw", "owner-pw", model.PermissionsNone)
	if _, err := convertBytes(b, nil); !errors.Is(err, converter.ErrPasswordRequired) {
		t.Fatalf("no password: %v", err)
	}
	if _, err := convertBytes(b, &converter.Options{Password: "wrong"}); !errors.Is(err, converter.ErrWrongPassword) {
		t.Fatalf("wrong password: %v", err)
	}
	for _, pw := range []string{"user-pw", "owner-pw"} {
		res, err := convertBytes(b, &converter.Options{Password: pw})
		if err != nil {
			t.Fatalf("%s: %v", pw, err)
		}
		if !res.Protected || res.Summary != plain.Summary {
			t.Fatalf("%s: protected %v, %q, want %q", pw, res.Protected, res.Summary, plain.Summary)
		}
	}
	r := bytes.NewReader(b)
	for pw, want := range map[string]error{"": converter.ErrPasswordRequired, "wrong": converter.ErrWrongPassword, "user-pw": nil} {
		if p, err := converter.CheckPassword(r, r.Size(), pw); !p || !errors.Is(err, want) {
			t.Fatalf("CheckPassword(%q) = %v, %v", pw, p, err)
		}
	}

	// Only an owner password: the PDF opens without one, and a password
	// given anyway does no harm.
	b = encryptPDF(t, "", "owner-pw", model.PermissionsNone)
	for _, pw := range []string{"", "anything"} {
		res, err := convertBytes(b, &converter.Options{Password: pw})
		if err != nil || res.Protected {
			t.Fatalf("owner password only, %q: protected %v, %v", pw, res != nil && res.Protected, err)
		}
	}
	if p, err := converter.CheckPassword(bytes.NewReader(b), int64(len(b)), ""); p || err != nil {
		t.Fatalf("CheckPassword on an owner-password PDF = %v, %v", p, err)
	}
}
