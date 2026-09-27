package main

import (
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
