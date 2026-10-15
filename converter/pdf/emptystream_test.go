package pdf

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	conv "github.com/shibukawa/bdf/converter"
)

// vectorOnlyPDF is a PDF encrypted with AES whose page has three content
// streams, the one in the middle empty and written as its initialization
// vector alone: pdfcpu, which encrypts the file, writes a block of padding
// after the vector, cut off here.
func vectorOnlyPDF(t *testing.T, keyLength int, user string) []byte {
	t.Helper()
	plain := buildPDF([]any{
		`<< /Type /Catalog /Pages 2 0 R >>`,
		`<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
		`<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents [4 0 R 5 0 R 6 0 R] >>`,
		[2]string{"", "0 0 10 10 re 20 20 10 10 re"},
		[2]string{"", ""},
		[2]string{"", " f 40 40 10 10 re f"},
	}, "/Root 1 0 R")
	conf := model.NewStatelessConfiguration()
	conf.UserPW, conf.OwnerPW, conf.EncryptKeyLength = user, "owner", keyLength
	conf.EncryptUsingAES = true
	var out bytes.Buffer
	if err := api.Encrypt(t.Context(), bytes.NewReader(plain), &out, conf); err != nil {
		t.Fatal(err)
	}
	// the vector and the padding: 32 bytes, of which the vector stays, and
	// the rest of the object is filled so that no offset changes
	b := out.Bytes()
	const head, tail = "<</Length 32>>\nstream\n", "\nendstream\nendobj\n"
	at := bytes.Index(b, []byte(head))
	if at < 0 || !bytes.HasPrefix(b[at+len(head)+32:], []byte(tail)) {
		t.Fatalf("no empty stream of 32 bytes in the file pdfcpu encrypted: %q", b)
	}
	copy(b[at:], "<</Length 16>>")
	copy(b[at+len(head)+16:], tail+strings.Repeat(" ", 16))
	return b
}

func TestEmptyStreamOfItsVectorAlone(t *testing.T) {
	for _, tc := range []struct {
		name      string
		keyLength int
		user      string
	}{
		{"AES-128", 128, ""},
		{"AES-256", 256, ""},
		{"AES-256 with a password", 256, "user"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := vectorOnlyPDF(t, tc.keyLength, tc.user)
			if tc.user != "" {
				if _, _, err := readContext(bytes.NewReader(data), ""); !errors.Is(err, conv.ErrPasswordRequired) {
					t.Fatalf("no password: %v", err)
				}
				if _, _, err := readContext(bytes.NewReader(data), "wrong"); !errors.Is(err, conv.ErrWrongPassword) {
					t.Fatalf("wrong password: %v", err)
				}
			}
			ctx, protected, err := readContext(bytes.NewReader(data), tc.user)
			if err != nil {
				t.Fatal(err)
			}
			if protected != (tc.user != "") {
				t.Errorf("protected %v", protected)
			}
			p := &pdf{ctx: ctx}
			leaves, _, _ := p.pageLeaves()
			if got, err := p.pageContent(leaves[0].dict); err != nil || string(got) != "0 0 10 10 re 20 20 10 10 re f 40 40 10 10 re f" {
				t.Errorf("content %q, %v", got, err)
			}
			if _, err := Convert(bytes.NewReader(data), &Options{Password: tc.user}); err != nil {
				t.Error(err)
			}
		})
	}
}

// patchedBytes is what emptyStreams makes of a file, read n bytes at a time.
func patchedBytes(t *testing.T, file string, n int) string {
	t.Helper()
	rs, err := emptyStreams(strings.NewReader(file))
	if err != nil {
		t.Fatal(err)
	}
	if rs == nil {
		return file
	}
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if _, err := io.CopyBuffer(&out, struct{ io.Reader }{rs}, make([]byte, n)); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestEmptyStreamsShape(t *testing.T) {
	const iv = "0123456789abcdef"
	empty := strings.Repeat("\n", 16)
	for _, tc := range []struct{ name, file, want string }{
		{"line feeds", "4 0 obj\n<< /Length 16 >>\nstream\n" + iv + "\nendstream\nendobj\n", "4 0 obj\n<< /Length 0  >>\nstream\n" + empty + "\nendstream\nendobj\n"},
		{"carriage returns, as in the file it was found in", "670 0 obj\r<</Length 16>>stream\r\n" + iv + "\r\nendstream\rendobj\r", "670 0 obj\r<</Length 0 >>stream\r\n" + empty + "\r\nendstream\rendobj\r"},
		{"no line end before endstream", "4 0 obj<</Filter/FlateDecode/Length 16>>stream\n" + iv + "endstream endobj", "4 0 obj<</Filter/FlateDecode/Length 0 >>stream\n" + empty + "endstream endobj"},
		{"a length in another object", "4 0 obj\n<< /Length 12 0 R >>\nstream\n" + iv + "\nendstream\nendobj\n", "4 0 obj\n<< /Length 0      >>\nstream\n" + empty + "\nendstream\nendobj\n"},
		{"the length after those of a font file", "4 0 obj\n<< /Length1 16 /Length 16 >>\nstream\n" + iv + "\nendstream\n", "4 0 obj\n<< /Length1 16 /Length 0  >>\nstream\n" + empty + "\nendstream\n"},
		{"the second of two objects", "3 0 obj\n<< /Length 3 >>\nstream\nq Q\nendstream\nendobj\n4 0 obj\n<< /Length 16 >>\nstream\n" + iv + "\nendstream\n", "3 0 obj\n<< /Length 3 >>\nstream\nq Q\nendstream\nendobj\n4 0 obj\n<< /Length 0  >>\nstream\n" + empty + "\nendstream\n"},
		// left as they are
		{"a longer stream", "4 0 obj\n<< /Length 32 >>\nstream\n" + iv + iv + "\nendstream\n", ""},
		{"a longer stream that ends like one", "4 0 obj\n<< /Length 39 >>\nstream\n" + iv + "stream\n" + iv + "\nendstream\n", ""},
		{"another length", "4 0 obj\n<< /Length 20 >>\nstream\n" + iv + "\nendstream\n", ""},
		{"a cross-reference stream", "4 0 obj\n<< /Type /XRef /W [1 2 1] /Size 4 /Length 16 >>\nstream\n" + iv + "\nendstream\n", ""},
		{"a stream with its own crypt filter", "4 0 obj\n<< /Filter /Crypt /Length 16 >>\nstream\n" + iv + "\nendstream\n", ""},
		{"no dictionary", "4 0 obj\n(/Length 16)\nstream\n" + iv + "\nendstream\n", ""},
		{"no object", "<< /Length 16 >>\nstream\n" + iv + "\nendstream\n", ""},
	} {
		want := tc.want
		if want == "" {
			want = tc.file
		}
		for _, n := range []int{1, 7, 4096} {
			if got := patchedBytes(t, tc.file, n); got != want {
				t.Errorf("%s, read %d bytes at a time:\n%q, want\n%q", tc.name, n, got, want)
			}
		}
	}
}

// A stream anywhere in a file larger than what the scan reads at a time.
func TestEmptyStreamsAcrossChunks(t *testing.T) {
	const object = "4 0 obj\n<< /Length 16 >>\nstream\n0123456789abcdef\nendstream\nendobj\n"
	patched := strings.Replace(strings.Replace(object, "16", "0 ", 1), "0123456789abcdef", strings.Repeat("\n", 16), 1)
	for at := scanChunk - len(object) - 2; at <= scanChunk+2; at++ {
		pad, rest := strings.Repeat("%", at), strings.Repeat("%", scanChunk+scanBack)
		for _, lead := range []string{"", strings.Repeat("%", scanChunk)} {
			if got := patchedBytes(t, lead+pad+object+rest, 64<<10); got != lead+pad+patched+rest {
				t.Fatalf("an empty stream at %d: not found, or not there", len(lead)+at)
			}
		}
	}
}
