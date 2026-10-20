package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf/print/spool"
)

func TestPrintConverter(t *testing.T) {
	c := printConverter{fontDirs: []string{filepath.Join("..", "..", "converter", "pptx", "testdata", "fonts")}, noSystem: true}
	ctx := context.Background()
	read := func(path string) *bytes.Reader {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return bytes.NewReader(b)
	}

	// an encrypted deck: its name tells it, or the encryption does
	deck := read(filepath.Join("..", "..", "converter", "internal", "offcrypto", "testdata", "agile.pptx"))
	if f, ok := c.Detect(deck, deck.Size(), "deck.pptx"); !ok || f != "pptx" {
		t.Fatalf("by name: %q %v", f, ok)
	}
	if f, ok := c.Detect(deck, deck.Size(), "upload"); !ok || f != encryptedOffice {
		t.Fatalf("by its encryption: %q %v", f, ok)
	}
	const password = "パスワード🔑bdf"
	for _, format := range []string{"pptx", encryptedOffice} {
		if _, _, err := c.Convert(ctx, deck, deck.Size(), format, "deck.pptx", ""); !errors.Is(err, spool.ErrPasswordRequired) {
			t.Fatalf("%s without the password: %v", format, err)
		}
		if _, _, err := c.Convert(ctx, deck, deck.Size(), format, "deck.pptx", "wrong"); !errors.Is(err, spool.ErrWrongPassword) {
			t.Fatalf("%s with a wrong password: %v", format, err)
		}
		doc, _, err := c.Convert(ctx, deck, deck.Size(), format, "deck.pptx", password)
		if err != nil {
			t.Fatal(err)
		}
		// what needed the password stays locked in the spool
		if doc.Lock == nil {
			t.Fatalf("%s: the converted deck is not locked", format)
		}
	}

	// Markdown: by its name, as pages
	md := bytes.NewReader([]byte("# Notes\n\nA line.\n"))
	f, ok := c.Detect(md, md.Size(), "notes.md")
	if !ok || f != "markdown" {
		t.Fatalf("markdown: %q %v", f, ok)
	}
	doc, _, err := c.Convert(ctx, md, md.Size(), f, "notes.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Views) == 0 || len(doc.Views[0].Pages) == 0 {
		t.Fatalf("markdown views %+v", doc.Views)
	}
	if _, ok := c.Detect(bytes.NewReader([]byte{0, 1, 2, 3}), 4, "blob"); ok {
		t.Fatal("bytes detected")
	}
	// a conversion stops waiting when asked
	cctx, cancel := context.WithCancel(ctx)
	cancel()
	if _, _, err := c.Convert(cctx, md, md.Size(), f, "notes.md", ""); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled: %v", err)
	}
}
