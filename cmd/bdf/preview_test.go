package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

var testFontFlags = []string{"-font-dir", filepath.Join("..", "..", "converter", "pptx", "testdata", "fonts"), "-no-system-fonts"}

// runEnv runs the bdf command with extra environment variables.
func runEnv(t *testing.T, env []string, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(append(os.Environ(), "BDF_TEST_MAIN=1"), env...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()
	if exit, ok := err.(*exec.ExitError); ok {
		return exit.ExitCode(), stderr.String()
	}
	if err != nil {
		t.Fatal(err)
	}
	return 0, stderr.String()
}

func pngSize(t *testing.T, path string) (int, int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	c, err := png.DecodeConfig(f)
	if err != nil {
		t.Fatal(err)
	}
	return c.Width, c.Height
}

// TestGeneratePreviews writes the thumbnail and the search text with the
// document.
func TestGeneratePreviews(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join("..", "..", "converter", "pptx", "testdata", "basic.pptx")
	thumb, text := filepath.Join(dir, "t.png"), filepath.Join(dir, "t.json")
	args := append([]string{"generate", "-q", "-thumbnail", thumb, "-thumbnail-size", "128", "-text", text}, testFontFlags...)
	status, stderr := run(t, append(args, in, filepath.Join(dir, "out.bdf"))...)
	if status != 0 {
		t.Fatalf("exit status %d: %s", status, stderr)
	}
	if w, h := pngSize(t, thumb); w != 128 || h != 96 {
		t.Errorf("thumbnail %d×%d, want the 4:3 slide at 128×96", w, h)
	}
	b, err := os.ReadFile(text)
	if err != nil {
		t.Fatal(err)
	}
	var st bdf.SearchText
	if err := json.Unmarshal(b, &st); err != nil {
		t.Fatal(err)
	}
	if len(st.Views) != 1 || len(st.Views[0].Pages) != 5 || st.Views[0].Pages[0].Text != "BDF from PowerPoint\nSlides rendered directly from DrawingML" {
		t.Errorf("text %s", b)
	}
}

// TestPreviewsOfEncrypted checks that the thumbnail and text of a
// password-protected input, whose document is encrypted, are written only
// with -allow-plaintext.
func TestPreviewsOfEncrypted(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join("..", "..", "converter", "internal", "offcrypto", "testdata", "agile.pptx")
	if _, err := os.Stat(in); err != nil {
		t.Skip(err)
	}
	env := []string{passwordEnv + "=パスワード🔑bdf"}
	thumb, out := filepath.Join(dir, "t.png"), filepath.Join(dir, "out.bdf")
	args := append([]string{"generate", "-q", "-thumbnail", thumb}, testFontFlags...)
	status, stderr := runEnv(t, env, append(args, in, out)...)
	if status != 0 || !strings.Contains(stderr, "no thumbnail or text written") {
		t.Fatalf("exit status %d: %s", status, stderr)
	}
	if _, err := os.Stat(thumb); err == nil {
		t.Fatal("a thumbnail of the encrypted document was written")
	}
	status, stderr = runEnv(t, env, append(append(args, "-allow-plaintext"), in, out)...)
	if status != 0 {
		t.Fatalf("exit status %d: %s", status, stderr)
	}
	if _, err := os.Stat(thumb); err != nil {
		t.Fatal("-allow-plaintext wrote no thumbnail")
	}
	// the text command refuses the encrypted document the same way
	if status, stderr = runEnv(t, env, "text", out); status == 0 || !strings.Contains(stderr, "-allow-plaintext") {
		t.Errorf("text of the encrypted document: exit status %d: %s", status, stderr)
	}
	if status, stderr = runEnv(t, env, "text", "-allow-plaintext", out, filepath.Join(dir, "t.json")); status != 0 {
		t.Errorf("text -allow-plaintext: exit status %d: %s", status, stderr)
	}
	if status, _ = run(t, "thumbnail", "-allow-plaintext", out, filepath.Join(dir, "u.png")); status == 0 {
		t.Error("an encrypted document was drawn without its password")
	}
}

func TestThumbnailAndRender(t *testing.T) {
	dir := t.TempDir()
	doc := filepath.Join("..", "..", "testdata", "docx", "basic.bdf")
	out := filepath.Join(dir, "t.png")
	args := append([]string{"thumbnail", "-size", "64", "-mode", "fit"}, testFontFlags...)
	if status, stderr := run(t, append(args, doc, out)...); status != 0 {
		t.Fatalf("exit status %d: %s", status, stderr)
	}
	if w, h := pngSize(t, out); w != 45 || h != 64 {
		t.Errorf("fitted A4 page %d×%d", w, h)
	}
	args = append([]string{"render", "-page", "2", "-scale", "0.5"}, testFontFlags...)
	if status, stderr := run(t, append(args, doc, out)...); status != 0 {
		t.Fatalf("exit status %d: %s", status, stderr)
	}
	if w, h := pngSize(t, out); w != 298 || h != 421 {
		t.Errorf("page 2 at half scale %d×%d", w, h)
	}
	if status, stderr := run(t, "thumbnail", "-mode", "zoom", doc, out); status != 2 || !strings.Contains(stderr, "crop or fit") {
		t.Errorf("bad mode: exit status %d: %s", status, stderr)
	}
	if status, stderr := run(t, "thumbnail", doc, filepath.Join(dir, "t.gif")); status == 0 || !strings.Contains(stderr, ".png, .jpg or .webp") {
		t.Errorf("gif: exit status %d: %s", status, stderr)
	}
}

// TestRenderAntialias checks that the flag reaches the rasterizer and that
// its default remains the existing smoothed output at low and high scales.
func TestRenderAntialias(t *testing.T) {
	dir := t.TempDir()
	doc := bdf.NewDocument()
	o := bdf.NewObject()
	o.FillColor(bdf.RGB(0, 0, 0)).FillRect(0.3, 0.3, 2.4, 2.4)
	hash, _ := doc.AddObject(o)
	v := doc.NewView("v", bdf.ViewFixed, "")
	v.AddPage(4, 4, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	var data bytes.Buffer
	if err := doc.WriteSingle(&data); err != nil {
		t.Fatal(err)
	}
	in := filepath.Join(dir, "edges.bdf")
	if err := os.WriteFile(in, data.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, scale := range []string{"1", "8"} {
		t.Run("scale="+scale, func(t *testing.T) {
			var smoothed []byte
			for _, mode := range []string{"", "default", "on", "off"} {
				out := filepath.Join(dir, "page.png")
				args := []string{"render", "-q", "-scale", scale, "-no-system-fonts"}
				if mode != "" {
					args = append(args, "-antialias="+mode)
				}
				if status, stderr := run(t, append(args, in, out)...); status != 0 {
					t.Fatalf("-antialias=%q: exit status %d: %s", mode, status, stderr)
				}
				pixels, err := os.ReadFile(out)
				if err != nil {
					t.Fatal(err)
				}
				if mode == "" {
					smoothed = pixels
				} else if mode == "off" {
					if bytes.Equal(smoothed, pixels) {
						t.Error("off produced the smoothed image")
					}
				} else if !bytes.Equal(smoothed, pixels) {
					t.Errorf("%s differs from the existing default", mode)
				}
			}
		})
	}
}

// TestRenderBadAntialias checks usage before loading the input or writing output.
func TestRenderBadAntialias(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "page.png")
	for _, mode := range []string{"auto", "false", "ON", ""} {
		status, stderr := run(t, "render", "-antialias="+mode, filepath.Join(dir, "missing.bdf"), out)
		if status != 2 || !strings.Contains(stderr, "-antialias must be default, on or off") {
			t.Errorf("-antialias=%q: exit status %d: %s", mode, status, stderr)
		}
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Errorf("invalid flag wrote output: %v", err)
	}
}

func TestRenderBatch(t *testing.T) {
	dir := t.TempDir()
	doc := bdf.NewDocument()
	v := doc.NewView("v", bdf.ViewFixed, "")
	for _, color := range []bdf.Color{0xff0000ff, 0x0000ffff} {
		o := bdf.NewObject()
		o.FillColor(color).FillRect(.3, .3, 4.5, 4.5)
		h, _ := doc.AddObject(o)
		v.AddPage(8, 8, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	var data bytes.Buffer
	if err := doc.WriteSingle(&data); err != nil {
		t.Fatal(err)
	}
	in := filepath.Join(dir, "input.bdf")
	if err := os.WriteFile(in, data.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	for _, aa := range []string{"on", "off"} {
		pattern := filepath.Join(dir, "batch-{page}.png")
		if status, stderr := run(t, "render", "-q", "-no-system-fonts", "-antialias="+aa, "-pages=2,1", "-workers=2", "-scale=1.5", in, pattern); status != 0 {
			t.Fatalf("batch: %d %s", status, stderr)
		}
		for _, page := range []string{"1", "2"} {
			single := filepath.Join(dir, "single.png")
			if status, stderr := run(t, "render", "-q", "-no-system-fonts", "-antialias="+aa, "-page="+page, "-scale=1.5", in, single); status != 0 {
				t.Fatalf("single: %d %s", status, stderr)
			}
			want, err := os.ReadFile(single)
			if err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(strings.ReplaceAll(pattern, "{page}", page))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("batch page %s differs for AA %s", page, aa)
			}
		}
	}
	if status, stderr := run(t, "render", "-q", "-no-system-fonts", "-pages=all", in, filepath.Join(dir, "all-{page}.png")); status != 0 {
		t.Fatalf("all: %d %s", status, stderr)
	}
	for _, args := range [][]string{{"-pages=all", "-page=1"}, {"-workers=2"}, {"-pages=all", "-workers=17"}, {"-pages=bad"}} {
		if status, _ := run(t, append(append([]string{"render"}, args...), filepath.Join(dir, "missing.bdf"), filepath.Join(dir, "bad-{page}.png"))...); status != 2 {
			t.Fatalf("bad flags %v: exit %d", args, status)
		}
	}
	if status, _ := run(t, "render", "-pages=all", in, filepath.Join(dir, "no-placeholder.png")); status != 2 {
		t.Fatalf("missing placeholder: %d", status)
	}
	if status, _ := run(t, "render", "-pages=1-3", in, filepath.Join(dir, "outside-{page}.png")); status != 2 {
		t.Fatalf("outside pages: %d", status)
	}
	if files, _ := filepath.Glob(filepath.Join(dir, "outside-*.png")); len(files) != 0 {
		t.Fatal("invalid selection wrote a partial batch")
	}
}

// TestThumbnailOfInput draws the thumbnail of an input that is not a bdf
// document: converted for it (only its first page), it is the thumbnail of
// the document bdf generate makes.
func TestThumbnailOfInput(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join("..", "..", "converter", "pdf", "testdata", "chrome-slides.pdf")
	doc, direct, via := filepath.Join(dir, "out.bdf"), filepath.Join(dir, "direct.png"), filepath.Join(dir, "via.png")
	if status, stderr := run(t, append(append([]string{"generate", "-q", "-images", "keep"}, testFontFlags...), in, doc)...); status != 0 {
		t.Fatalf("generate: exit status %d: %s", status, stderr)
	}
	for _, c := range [][2]string{{in, direct}, {doc, via}} {
		if status, stderr := run(t, append(append([]string{"thumbnail", "-q", "-size", "96"}, testFontFlags...), c[0], c[1])...); status != 0 {
			t.Fatalf("thumbnail %s: exit status %d: %s", c[0], status, stderr)
		}
	}
	a, err := os.ReadFile(direct)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(via)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Error("the thumbnail of the PDF differs from that of its bdf document")
	}
	if status, stderr := run(t, "thumbnail", filepath.Join("..", "..", "go.mod"), direct); status != 2 || !strings.Contains(stderr, "neither a bdf document") {
		t.Errorf("unknown input: exit status %d: %s", status, stderr)
	}
}

// TestThumbnailOfProtectedInput draws a password-protected input only with
// -allow-plaintext, as generate writes its previews.
func TestThumbnailOfProtectedInput(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join("..", "..", "converter", "internal", "offcrypto", "testdata", "agile.pptx")
	if _, err := os.Stat(in); err != nil {
		t.Skip(err)
	}
	out := filepath.Join(dir, "t.png")
	args := append([]string{"thumbnail", "-q"}, testFontFlags...)
	if status, stderr := run(t, append(args, in, out)...); status == 0 || !strings.Contains(stderr, "is encrypted") {
		t.Errorf("without the password: exit status %d: %s", status, stderr)
	}
	env := []string{passwordEnv + "=パスワード🔑bdf"}
	if status, stderr := runEnv(t, env, append(args, in, out)...); status == 0 || !strings.Contains(stderr, "-allow-plaintext") {
		t.Errorf("without -allow-plaintext: exit status %d: %s", status, stderr)
	}
	if _, err := os.Stat(out); err == nil {
		t.Fatal("a thumbnail of the protected input was written")
	}
	if status, stderr := runEnv(t, env, append(append(args, "-allow-plaintext"), in, out)...); status != 0 {
		t.Fatalf("-allow-plaintext: exit status %d: %s", status, stderr)
	}
	if w, h := pngSize(t, out); w != 256 || h != 192 {
		t.Errorf("thumbnail %d×%d, want the 4:3 slide at 256×192", w, h)
	}
}
