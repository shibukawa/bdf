package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/print/spool"
)

// printViews are the views of the formats whose default views have no
// pages to print: an Excel workbook as Excel prints it, HTML and Markdown
// as A4 pages instead of one long column.
var printViews = map[string]string{"xlsx": "print", "html": "pages", "markdown": "pages"}

// printConverter converts the documents bdf printd receives with the
// converters of bdf generate, into views of pages (printViews). HTML,
// Markdown and EPUB leave out the images of the network (remote=false):
// the server is no way into the networks it reaches. Nothing beside a
// document is read.
type printConverter struct {
	fontDirs []string
	noSystem bool
}

func (c printConverter) Detect(r io.ReaderAt, size int64, name string) (string, bool) {
	// the rasters of driverless clients, which the IPP server takes
	head := make([]byte, 8)
	if n, _ := r.ReadAt(head, 0); n > 0 {
		switch raster.Sniff(head[:n]) {
		case raster.FormatPWG:
			return formatPWG, true
		case raster.FormatURF:
			return formatURF, true
		}
	}
	f := converter.Detect(r, size)
	if f == nil || f.Guess {
		// the content does not tell (a CSV file of one line, Markdown):
		// the extension does
		ext := strings.ToLower(filepath.Ext(name))
		for _, g := range converter.Formats() {
			if ext != "" && slices.Contains(g.Extensions, ext) {
				f = g
				break
			}
		}
	}
	if f == nil {
		// an encrypted Office document tells its format once decrypted
		if protected, _ := converter.CheckPassword(r, size, ""); protected {
			return encryptedOffice, true
		}
		return "", false
	}
	return f.Name, true
}

// encryptedOffice is the format of an encrypted Office document whose
// name does not tell which it is.
const encryptedOffice = "office-encrypted"

// The raster formats, which print/raster reads (converter does not: the
// converters are built without print).
const (
	formatPWG = "pwg-raster"
	formatURF = "urf"
)

func (c printConverter) Convert(ctx context.Context, r io.ReaderAt, size int64, format, name, password string) (*bdf.Document, []string, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if format == formatPWG || format == formatURF {
		doc, _, err := raster.Document(ctx, io.NewSectionReader(r, 0, size), &raster.DocumentOptions{Title: name})
		return doc, nil, err
	}
	opts := &converter.Options{FontDirs: c.fontDirs, NoSystemFonts: c.noSystem, Password: password, FileName: name, Params: map[string]string{}}
	if format == encryptedOffice {
		format = "" // detected once decrypted
	} else {
		f := converter.Lookup(format)
		if f == nil {
			return nil, nil, fmt.Errorf("no converter of %s", format)
		}
		for _, p := range f.Params {
			switch {
			case p.Name == "remote":
				opts.Params["remote"] = "false"
			case p.Name == "views" && printViews[format] != "":
				opts.Params["views"] = printViews[format]
			}
		}
	}
	type outcome struct {
		res *converter.Result
		err error
	}
	done := make(chan outcome, 1)
	// a conversion does not stop when asked: the spool that stops does
	// not wait for it
	go func() {
		res, err := converter.Convert(r, size, format, opts)
		done <- outcome{res, err}
	}()
	var o outcome
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case o = <-done:
	}
	switch {
	case errors.Is(o.err, converter.ErrPasswordRequired):
		return nil, nil, spool.ErrPasswordRequired
	case errors.Is(o.err, converter.ErrWrongPassword):
		return nil, nil, spool.ErrWrongPassword
	case o.err != nil:
		return nil, nil, o.err
	}
	if o.res.Protected && password != "" {
		// what needed the password to read needs it in the spool too
		lock, err := bdf.NewPasswordLock(password, 0)
		if err != nil {
			return nil, nil, err
		}
		o.res.Doc.Lock = lock
	}
	return o.res.Doc, o.res.Warnings, nil
}
