package main

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
)

// segmentCmd writes a segment of a document (docs/spec.md §3.6), as package
// segment answers a reader's request: for looking at what a server sends,
// and for the test files of @bdf/core.
func segmentCmd(args []string) {
	fs := flag.NewFlagSet("segment", flag.ExitOnError)
	view := fs.String("view", "", "the view (default: the first one)")
	page := fs.Int("page", 1, "the page the segment holds (from 1)")
	size := fs.Int("size", 10, "pages of a segment")
	have := fs.String("have", "", "pages of the view the reader holds already, e.g. 1-10,21-30: the parts they need are left out")
	keyFile := fs.String("key", "", "seal the segment for this public key (a PEM file of a P-256 key); without it the segment is in the clear")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf segment [flags] <file.bdf | dir> <out.bdf>")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 2 || *size < 1 {
		fs.Usage()
		os.Exit(2)
	}
	r, err := open(fs.Arg(0))
	check(err)
	if r.Locked() {
		check(fmt.Errorf("%w ($%s)", bdf.ErrLocked, passwordEnv))
	}
	var v *bdf.View
	for _, w := range r.Manifest.Views {
		if *view == "" || w.ID == *view {
			v = w
			break
		}
	}
	if v == nil {
		check(fmt.Errorf("no view %q", *view))
	}
	if *page < 1 || *page > len(v.Pages) {
		check(fmt.Errorf("view %q has no page %d", v.ID, *page))
	}
	o := bdf.SegmentOptions{Segment: bdf.SegmentAt(v, *page-1, *size)}
	spans, err := converter.ParsePages(*have)
	check(err)
	for _, s := range spans {
		to := s.To
		if to == 0 || to > len(v.Pages) {
			to = len(v.Pages)
		}
		o.Have = append(o.Have, bdf.Segment{View: v.ID, From: s.From - 1, To: to})
	}
	if *keyFile != "" {
		key, err := readPublicKey(*keyFile)
		check(err)
		o.Lock, err = bdf.NewECDHLock(key)
		check(err)
	}
	f, err := os.Create(fs.Arg(1))
	check(err)
	check(errors.Join(r.WriteSegment(f, o), f.Close()))
}

// readPublicKey reads a P-256 public key from a PEM file ("PUBLIC KEY",
// as openssl ec -pubout writes it).
func readPublicKey(path string) (*ecdh.PublicKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	blk, _ := pem.Decode(b)
	if blk == nil || blk.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("%s: not a PEM public key", path)
	}
	k, err := x509.ParsePKIXPublicKey(blk.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	ek, ok := k.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%s: not an elliptic curve key", path)
	}
	return ek.ECDH()
}
