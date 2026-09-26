package jbig2

// Regenerating the fixtures:
//
//	go test ./converter/internal/jbig2/ -run TestGenerate -regen
//
// builds each stream with the test encoder, renders it with Ghostscript
// (which decodes JBIG2 with jbig2dec) wrapped in a PDF, checks that
// Ghostscript's page matches the reference rendering, and writes the stream
// (name.jb2, with name.glob for global segments) and Ghostscript's page
// (name.pbm) to testdata.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

var regen = flag.Bool("regen", false, "regenerate testdata with Ghostscript")

// jbig2PDF wraps a JBIG2 stream in a one-page PDF that draws it pixel for pixel
// at 72 dpi. decode is the image's Decode array, or "" for the default.
func jbig2PDF(page, globals []byte, w, h int, decode string) []byte {
	var b bytes.Buffer
	var offsets []int
	obj := func(body string, stream []byte) {
		offsets = append(offsets, b.Len())
		fmt.Fprintf(&b, "%d 0 obj\n%s", len(offsets), body)
		if stream != nil {
			fmt.Fprintf(&b, "\nstream\n%s\nendstream", stream)
		}
		b.WriteString("\nendobj\n")
	}
	b.WriteString("%PDF-1.5\n")
	content := fmt.Sprintf("q %d 0 0 %d 0 0 cm /Im0 Do Q", w, h)
	parms := ""
	if globals != nil {
		parms = " /DecodeParms << /JBIG2Globals 6 0 R >>"
	}
	if decode != "" {
		decode = " /Decode " + decode
	}
	obj("<< /Type /Catalog /Pages 2 0 R >>", nil)
	obj("<< /Type /Pages /Kids [3 0 R] /Count 1 >>", nil)
	obj(fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %d %d] /Contents 4 0 R /Resources << /XObject << /Im0 5 0 R >> >> >>", w, h), nil)
	obj(fmt.Sprintf("<< /Length %d >>", len(content)), []byte(content))
	obj(fmt.Sprintf("<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceGray /BitsPerComponent 1 /Filter /JBIG2Decode%s%s /Length %d >>",
		w, h, parms, decode, len(page)), page)
	if globals != nil {
		obj(fmt.Sprintf("<< /Length %d >>", len(globals)), globals)
	}
	xref := b.Len()
	fmt.Fprintf(&b, "xref\n0 %d\n0000000000 65535 f \n", len(offsets)+1)
	for _, o := range offsets {
		fmt.Fprintf(&b, "%010d 00000 n \n", o)
	}
	fmt.Fprintf(&b, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(offsets)+1, xref)
	return b.Bytes()
}

// ghostscript renders a PDF at 72 dpi and returns the page with dark pixels as 1.
func ghostscript(t *testing.T, pdf []byte) *Bitmap {
	t.Helper()
	dir := t.TempDir()
	in, out := filepath.Join(dir, "in.pdf"), filepath.Join(dir, "out.pgm")
	if err := os.WriteFile(in, pdf, 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("gs", "-q", "-dNOPAUSE", "-dBATCH", "-dSAFER", "-sDEVICE=pgmraw", "-r72", "-dNOINTERPOLATE", "-o", out, in)
	msg, err := cmd.CombinedOutput()
	if err != nil || bytes.Contains(msg, []byte("FATAL")) || bytes.Contains(msg, []byte("failed")) {
		t.Fatalf("gs: %v: %s", err, msg)
	}
	if len(msg) > 0 {
		t.Logf("gs: %s", msg)
	}
	pgm, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var w, h, maxv int
	fields, rest := pnmHeader(pgm, 4)
	if fields[0] != "P5" {
		t.Fatalf("gs wrote %q", fields[0])
	}
	w, _ = strconv.Atoi(fields[1])
	h, _ = strconv.Atoi(fields[2])
	maxv, _ = strconv.Atoi(fields[3])
	bm := newBitmap(w, h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if int(rest[y*w+x])*2 < maxv {
				bm.set(x, y)
			}
		}
	}
	return bm
}

// pnmHeader splits off n whitespace separated header fields and the single
// whitespace byte after them.
func pnmHeader(b []byte, n int) ([]string, []byte) {
	var f []string
	i := 0
	for len(f) < n {
		for b[i] == ' ' || b[i] == '\n' || b[i] == '\r' || b[i] == '\t' || b[i] == '#' {
			if b[i] == '#' {
				for b[i] != '\n' {
					i++
				}
			}
			i++
		}
		j := i
		for b[j] != ' ' && b[j] != '\n' && b[j] != '\r' && b[j] != '\t' {
			j++
		}
		f = append(f, string(b[i:j]))
		i = j
	}
	return f, b[i+1:]
}

func encodePBM(b *Bitmap) []byte {
	return append([]byte(fmt.Sprintf("P4\n%d %d\n", b.Width, b.Height)), b.Data...)
}

func decodePBM(p []byte) (*Bitmap, error) {
	f, rest := pnmHeader(p, 3)
	if f[0] != "P4" {
		return nil, fmt.Errorf("not a PBM file")
	}
	w, _ := strconv.Atoi(f[1])
	h, _ := strconv.Atoi(f[2])
	b := newBitmap(w, h)
	if len(rest) != len(b.Data) {
		return nil, fmt.Errorf("PBM data size %d, want %d", len(rest), len(b.Data))
	}
	copy(b.Data, rest)
	return b, nil
}

// diff describes how two bitmaps differ, or returns "".
func diff(got, want *Bitmap) string {
	if got.Width != want.Width || got.Height != want.Height {
		return fmt.Sprintf("size %dx%d, want %dx%d", got.Width, got.Height, want.Width, want.Height)
	}
	n, fx, fy := 0, -1, -1
	for y := 0; y < want.Height; y++ {
		for x := 0; x < want.Width; x++ {
			if got.pixel(x, y) != want.pixel(x, y) {
				if n == 0 {
					fx, fy = x, y
				}
				n++
			}
		}
	}
	if n == 0 {
		return ""
	}
	return fmt.Sprintf("%d pixels differ, the first at (%d, %d)", n, fx, fy)
}

func bitmapHash(b *Bitmap) string {
	h := sha256.New()
	fmt.Fprintf(h, "%d %d\n", b.Width, b.Height)
	h.Write(b.Data)
	return hex.EncodeToString(h.Sum(nil))
}

func TestGenerate(t *testing.T) {
	if !*regen {
		t.Skip("run with -regen to regenerate testdata with Ghostscript")
	}
	if _, err := exec.LookPath("gs"); err != nil {
		t.Skip("gs not found")
	}
	old, _ := filepath.Glob(filepath.Join("testdata", "*.*"))
	for _, name := range old {
		if ext := filepath.Ext(name); ext == ".jb2" || ext == ".glob" || ext == ".pbm" {
			os.Remove(name)
		}
	}
	for _, f := range fixtures() {
		t.Run(f.name, func(t *testing.T) {
			var gs *Bitmap
			if f.reference != "" {
				t.Logf("not checked with Ghostscript: %s", f.reference)
				gs = f.want
			} else {
				gs = ghostscript(t, jbig2PDF(f.page, f.globals, f.want.Width, f.want.Height, ""))
				if d := diff(gs, f.want); d != "" {
					t.Fatalf("Ghostscript disagrees with the reference rendering: %s", d)
				}
			}
			base := filepath.Join("testdata", f.name)
			write := func(name string, data []byte) {
				if err := os.WriteFile(name, data, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			write(base+".jb2", f.page)
			if f.globals != nil {
				write(base+".glob", f.globals)
			}
			write(base+".pbm", encodePBM(gs))
		})
	}
	t.Run("a4", func(t *testing.T) {
		for _, p := range largePages() {
			gs := ghostscript(t, jbig2PDF(p.page, nil, p.want.Width, p.want.Height, ""))
			if d := diff(gs, p.want); d != "" {
				t.Fatalf("%s: Ghostscript disagrees with the reference rendering: %s", p.name, d)
			}
			t.Logf("%s: %s", p.name, bitmapHash(gs))
		}
	})
}

// TestSampleConvention confirms with Ghostscript that the JBIG2Decode filter
// delivers 1 = white: with the default Decode array a JBIG2 1 (black) pixel
// renders black, and with Decode [1 0] white.
func TestSampleConvention(t *testing.T) {
	if !*regen {
		t.Skip("run with -regen to check with Ghostscript")
	}
	img := newBitmap(40, 20)
	fillRect(img, 5, 5, 10, 10)
	s, _ := newPage(40, 20, 0)
	s.add(typeImmediateGeneric, genericRegion(img, 0, 0, opOr, false, 0, nominalAT(0), false), segOpts{})
	if d := diff(ghostscript(t, jbig2PDF(s.buf, nil, 40, 20, "")), img); d != "" {
		t.Errorf("default Decode: %s", d)
	}
	inv := newBitmap(40, 20)
	inv.fill(1)
	refCompose(inv, img, 0, 0, opXor)
	if d := diff(ghostscript(t, jbig2PDF(s.buf, nil, 40, 20, "[1 0]")), inv); d != "" {
		t.Errorf("Decode [1 0]: %s", d)
	}
}

var dump = flag.String("dump", "", "write the fixture streams to this directory, for debugging")

// TestDump writes the fixture streams and their reference pages to a
// directory, for instance to feed them to jbig2dec:
//
//	go test ./converter/internal/jbig2/ -run TestDump -dump /tmp/jb2
//	jbig2dec -e -t pbm -o out.pbm /tmp/jb2/text_arith.glob /tmp/jb2/text_arith.jb2
func TestDump(t *testing.T) {
	if *dump == "" {
		t.Skip("run with -dump dir")
	}
	for _, f := range fixtures() {
		os.WriteFile(filepath.Join(*dump, f.name+".jb2"), f.page, 0o644)
		os.WriteFile(filepath.Join(*dump, f.name+".glob"), f.globals, 0o644)
		os.WriteFile(filepath.Join(*dump, f.name+".want.pbm"), encodePBM(f.want), 0o644)
	}
}
