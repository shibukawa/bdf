package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/shibukawa/bdf/vector/pdfbdf"
)

// pdfCmd writes a document as a PDF file.
func pdfCmd(args []string) {
	fs := flag.NewFlagSet("pdf", flag.ExitOnError)
	var views, skip stringList
	fs.Var(&views, "view", "id of a view to write (repeatable; default: the views of pages, or all views of a document without pages)")
	fs.Var(&skip, "skip", "role of the layers to leave out, e.g. master or background, for paper that has them printed already (repeatable)")
	ignore := fs.Bool("ignore-fstype", false, "embed fonts of the system whose OS/2 fsType forbids embedding or subsetting (only with the rights to do so; default: their text is drawn as outlines)")
	plain := fs.Bool("no-compression", false, "write the streams of the file as they are, to read them")
	allow := fs.Bool("allow-plaintext", false, "write an encrypted document (read with $"+passwordEnv+"); the PDF is not encrypted")
	quiet := fs.Bool("q", false, "do not print warnings")
	tagged := fs.Bool("tagged", false, "write a tagged PDF: the structure of the document (headings, paragraphs, lists, tables, figures with their alternative texts) for screen readers and the tools that take text out")
	pdfa := fs.String("pdfa", "", "write a PDF/A file, for keeping: 2b (it shows the same in any reader), 2u (and its text is Unicode) or 2a (and it is tagged); 3b, 3u and 3a are the same of part 3")
	signKey := fs.String("sign-key", "", "sign the file (PAdES) with this private key, in PEM (as bdf sign takes it)")
	signCert := fs.String("sign-cert", "", "with -sign-key: the certificate of the key, followed by those of its chain (PEM)")
	signField := fs.String("sign-field", "", "with -sign-key: the name of the picture of a signature field of the form (party_a.sign), which the signature shows on (default: a signature that shows on no page)")
	signName := fs.String("sign-name", "", "with -sign-key: the name of the signer, as the signature says it")
	signReason := fs.String("sign-reason", "", "with -sign-key: why the file is signed")
	signLocation := fs.String("sign-location", "", "with -sign-key: where it is signed")
	tsa := fs.String("tsa", "", "with -sign-key: URL of a time-stamp authority (RFC 3161) that stamps the signature (PAdES B-T)")
	certify := fs.Int("certify", 0, "with -sign-key: certify the file, permitting after it 1 no changes, 2 filling in forms and signing, 3 those and annotations (default 0: an approval signature)")
	ropts := fontFlags(fs)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf pdf [flags] <file.bdf | dir> <out.pdf>\n  writes the pages of a document as a PDF file: its text as text in the embedded fonts, its drawings as vectors")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 2 {
		fs.Usage()
		os.Exit(2)
	}
	doc := loadDocument(fs.Arg(0), *allow)
	o := ropts()
	var sig *pdfbdf.Signature
	if *signKey != "" {
		if *signCert == "" {
			check(fmt.Errorf("-sign-key needs -sign-cert"))
		}
		key, chain, err := loadSigner(*signKey, *signCert)
		check(err)
		sig = &pdfbdf.Signature{Key: key, Chain: chain, Field: *signField, Name: *signName, Reason: *signReason, Location: *signLocation, TSA: *tsa, Certify: *certify}
	}
	f, err := os.Create(fs.Arg(1))
	check(err)
	w := bufio.NewWriter(f)
	res, err := pdfbdf.Write(w, doc, &pdfbdf.Options{Views: views, SkipRoles: skip, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts,
		IgnoreFSType: *ignore, NoCompression: *plain, Tagged: *tagged, PDFA: *pdfa, Sign: sig})
	if err == nil {
		err = w.Flush()
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(fs.Arg(1))
		check(err)
	}
	if !*quiet {
		for _, w := range res.Warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	fmt.Fprintf(os.Stderr, "%s: %d page(s), %d embedded font(s), %d warning(s)\n", fs.Arg(1), res.Pages, res.Fonts, len(res.Warnings))
}
