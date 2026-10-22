package main

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shibukawa/bdf/form"
	"github.com/shibukawa/bdf/sign"
)

// readSigned reads a document and its signatures: a single file, or the
// directory of the split form.
func readSigned(path string) (*sign.Document, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if st.IsDir() {
		return sign.ReadSplit(path)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return sign.Read(bytes.NewReader(b), int64(len(b)))
}

// loadSigner reads a private key and the certificates of its chain from
// files of PEM (sign.ParsePEM).
func loadSigner(keyFile, certFile string) (crypto.Signer, []*x509.Certificate, error) {
	kb, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	cb, err := os.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}
	key, certs, err := sign.ParsePEM(kb, cb)
	if err != nil {
		return nil, nil, fmt.Errorf("%s, %s: %w", keyFile, certFile, err)
	}
	return key, certs, nil
}

func loadCerts(file string) ([]*x509.Certificate, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return sign.ParseCertificates(b)
}

// signCmd adds a signature to a document.
func signCmd(args []string) {
	fs := flag.NewFlagSet("sign", flag.ExitOnError)
	keyFile := fs.String("key", "", "the private key of the signer, in PEM (PKCS #8, EC or PKCS #1)")
	certFile := fs.String("cert", "", "the certificate of the key, followed by those of its chain (PEM)")
	tsa := fs.String("tsa", "", "URL of a time-stamp authority (RFC 3161) that stamps the signature")
	role := fs.String("role", "", "the role of the signer, signed with the document (author, counterparty, witness …)")
	claims := fs.String("claims", "", "a JSON file of more claims of the signer, signed with the document")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf sign -key key.pem -cert chain.pem [flags] <file.bdf | dir> [out.bdf]\n  adds a signature to a document, after those it has. Without out, a single file is signed in its place;\n  a document of the split form is signed in its directory (signatures.json)")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() < 1 || fs.NArg() > 2 || *keyFile == "" || *certFile == "" {
		fs.Usage()
		os.Exit(2)
	}
	in := fs.Arg(0)
	out := in
	if fs.NArg() == 2 {
		out = fs.Arg(1)
	}
	key, chain, err := loadSigner(*keyFile, *certFile)
	check(err)
	o := &sign.Options{TSA: *tsa, Claims: map[string]any{}}
	if *claims != "" {
		b, err := os.ReadFile(*claims)
		check(err)
		check(json.Unmarshal(b, &o.Claims))
	}
	if *role != "" {
		o.Claims["role"] = *role
	}
	if len(o.Claims) == 0 {
		o.Claims = nil
	}
	d, err := readSigned(in)
	check(err)
	sig, err := d.Sign(key, chain, o)
	check(err)
	if st, err := os.Stat(in); err == nil && st.IsDir() {
		if out != in {
			check(errors.New("a document of the split form is signed in its directory"))
		}
		check(d.WriteSplit(sig))
	} else {
		var b bytes.Buffer
		check(d.WriteSingle(&b, sig))
		check(os.WriteFile(out, b.Bytes(), 0o644))
	}
	fmt.Fprintf(os.Stderr, "%s: signature %d by %s\n", out, len(d.Signatures)+1, chain[0].Subject)
}

// verifyCmd checks the signatures of a document.
func verifyCmd(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	roots := fs.String("roots", "", "certificates to trust, in PEM (default: those of the system)")
	asJSON := fs.Bool("json", false, "write what is found as JSON")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf verify [flags] <file.bdf | dir>\n  checks the signatures of a document: that nothing changed after each, and that its signer is trusted.\n  Exits 1 when a signature does not hold or is not trusted, or when there is none")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}
	o := &sign.VerifyOptions{}
	if *roots != "" {
		certs, err := loadCerts(*roots)
		check(err)
		o.Roots = x509.NewCertPool()
		for _, c := range certs {
			o.Roots.AddCert(c)
		}
	}
	d, err := readSigned(fs.Arg(0))
	check(err)
	sigs, err := d.Verify(o)
	check(err)
	// the signatures of the values of a form (bdf fill -sign-key)
	var values []*form.ValueSignature
	if st, err := os.Stat(fs.Arg(0)); err == nil && !st.IsDir() {
		b, err := os.ReadFile(fs.Arg(0))
		check(err)
		values, err = form.VerifyValues(bytes.NewReader(b), int64(len(b)), o)
		check(err)
	}
	ok := len(sigs)+len(values) > 0
	type report struct {
		Index     int            `json:"index"`
		Signer    string         `json:"signer,omitempty"`
		Signed    string         `json:"signed,omitempty"`
		TimeStamp string         `json:"timeStamp,omitempty"`
		TSA       string         `json:"tsa,omitempty"`
		Claims    map[string]any `json:"claims,omitempty"`
		Valid     bool           `json:"valid"`
		Error     string         `json:"error,omitempty"`
		Trust     string         `json:"trust,omitempty"`
	}
	var reports []report
	for _, s := range sigs {
		ok = ok && s.Valid()
		r := report{Index: s.Index + 1, Claims: s.Claims, Valid: s.Valid()}
		if s.Certificate != nil {
			r.Signer = s.Certificate.Subject.String()
		}
		if !s.SigningTime.IsZero() {
			r.Signed = s.SigningTime.Format(time.RFC3339)
		}
		if !s.TimeStamp.IsZero() {
			r.TimeStamp = s.TimeStamp.Format(time.RFC3339)
		}
		if s.TSA != nil {
			r.TSA = s.TSA.Subject.String()
		}
		if s.Err != nil {
			r.Error = s.Err.Error()
		}
		if s.TrustErr != nil {
			r.Trust = s.TrustErr.Error()
		}
		reports = append(reports, r)
	}
	var valueReports []report
	for _, s := range values {
		ok = ok && s.Valid()
		r := report{Index: s.Index + 1, Valid: s.Valid()}
		if s.Values != nil {
			r.Claims = s.Values.Claims
		}
		if s.Certificate != nil {
			r.Signer = s.Certificate.Subject.String()
		}
		if !s.SigningTime.IsZero() {
			r.Signed = s.SigningTime.Format(time.RFC3339)
		}
		if !s.TimeStamp.IsZero() {
			r.TimeStamp = s.TimeStamp.Format(time.RFC3339)
		}
		if s.Err != nil {
			r.Error = s.Err.Error()
		}
		if s.TrustErr != nil {
			r.Trust = s.TrustErr.Error()
		}
		valueReports = append(valueReports, r)
	}
	if *asJSON {
		b, err := json.MarshalIndent(map[string]any{"signatures": reports, "values": valueReports, "valid": ok}, "", "  ")
		check(err)
		os.Stdout.Write(append(b, '\n'))
	} else {
		if len(sigs)+len(values) == 0 {
			fmt.Println("no signatures")
		}
		for i, r := range append(reports, valueReports...) {
			kind := "signature"
			if i >= len(reports) {
				kind = "signature of the values of the form"
			}
			state := "valid"
			switch {
			case r.Error != "":
				state = "BROKEN"
			case r.Trust != "":
				state = "holds, NOT TRUSTED"
			}
			fmt.Printf("%s %d: %s\n  signer: %s\n", kind, r.Index, state, r.Signer)
			if r.Signed != "" {
				fmt.Printf("  signed: %s (the signer's clock)\n", r.Signed)
			}
			if r.TimeStamp != "" {
				fmt.Printf("  time-stamp: %s by %s\n", r.TimeStamp, r.TSA)
			}
			if len(r.Claims) > 0 {
				b, _ := json.Marshal(r.Claims)
				fmt.Printf("  claims: %s\n", b)
			}
			for _, e := range []string{r.Error, r.Trust} {
				for _, line := range strings.Split(e, "\n") {
					if line != "" {
						fmt.Printf("  %s\n", line)
					}
				}
			}
		}
	}
	if !ok {
		os.Exit(1)
	}
}

// copySignatures keeps the signatures of a document when split writes the
// split form of it, and join the single file.
func copySignatures(from, to string) {
	d, err := readSigned(from)
	check(err)
	if len(d.Signatures) == 0 {
		return
	}
	st, err := os.Stat(to)
	check(err)
	if st.IsDir() {
		out, err := sign.ReadSplit(to)
		check(err)
		check(out.WriteSplit(d.Signatures...))
		return
	}
	b, err := os.ReadFile(to)
	check(err)
	single, err := sign.Read(bytes.NewReader(b), int64(len(b)))
	check(err)
	var out bytes.Buffer
	check(single.WriteSingle(&out, d.Signatures...))
	check(os.WriteFile(to, out.Bytes(), 0o644))
}
