// Command bdf generates, inspects and converts BDF files.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/fixture"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  bdf generate [flags] <input> <out.bdf | dir/>
                                     convert a document or a drawing (bdf generate -h for flags and formats)
  bdf ls <file.bdf | dir>            list views and parts
  bdf manifest <file.bdf | dir>      print the manifest as JSON
  bdf disasm <file.bdf | dir> <hash> disassemble an object part
  bdf extract <file.bdf | dir> <hash> <out>
  bdf split <file.bdf> <dir>         write the split form
  bdf join <dir> <file.bdf>          write the single-file form
  bdf encrypt [-password-file f] <file.bdf | dir> <out.bdf | dir/>
                                     encrypt a document with a password
  bdf decrypt [-password-file f] <file.bdf | dir> <out.bdf | dir/>
                                     write an encrypted document without its encryption
  bdf demo <file.bdf | dir/>         write the fixture document (dir/ ends with a slash)

An encrypted document is read with the password in $BDF_PASSWORD; split and
join copy it without the password.`)
	os.Exit(2)
}

// passwordEnv names the environment variable that holds a password.
const passwordEnv = "BDF_PASSWORD"

// readPassword reads a password from a file (- for the standard input), or
// from $BDF_PASSWORD when file is "". A line break at the end is dropped.
func readPassword(file string) (string, error) {
	if file == "" {
		return os.Getenv(passwordEnv), nil
	}
	var b []byte
	var err error
	if file == "-" {
		b, err = io.ReadAll(os.Stdin)
	} else {
		b, err = os.ReadFile(file)
	}
	return strings.TrimRight(string(b), "\r\n"), err
}

// openStored opens a document in either form as it is stored.
func openStored(path string) (*bdf.Reader, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if st.IsDir() {
		return bdf.OpenSplit(path)
	}
	return bdf.OpenSingleFile(path)
}

// open opens a document and unlocks an encrypted one with $BDF_PASSWORD, if
// set; otherwise an encrypted document stays locked.
func open(path string) (*bdf.Reader, error) {
	r, err := openStored(path)
	if err != nil {
		return nil, err
	}
	if pw := os.Getenv(passwordEnv); r.Locked() && pw != "" {
		if err := unlock(r, pw, path); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func unlock(r *bdf.Reader, password, path string) error {
	err := r.Unlock(password)
	if errors.Is(err, bdf.ErrWrongPassword) {
		return fmt.Errorf("the password does not open %s", path)
	}
	return err
}

// part reads a part, explaining an unknown part of a locked document.
func part(r *bdf.Reader, hash string) []byte {
	h, err := bdf.ParseHash(hash)
	check(err)
	if _, ok := r.Entry(h); !ok && r.Locked() {
		check(fmt.Errorf("%w ($%s)", bdf.ErrLocked, passwordEnv))
	}
	b, err := r.Part(h)
	check(err)
	return b
}

// isDir reports whether an output path names the split form.
func isDir(out string) bool {
	return strings.HasSuffix(out, "/") || strings.HasSuffix(out, string(filepath.Separator))
}

// write writes a document in the form the output path names.
func write(d *bdf.Document, out string) error {
	if isDir(out) {
		return d.WriteSplit(out)
	}
	return writeSingle(d, out)
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "bdf:", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	switch os.Args[1] {
	case "generate":
		generate(os.Args[2:])
	case "ls":
		r, err := open(os.Args[2])
		check(err)
		m := r.Manifest
		switch {
		case r.Locked():
			fmt.Printf("bdf %d encrypted: the parts are sealed; set $%s to read the document\n", m.BDF, passwordEnv)
		case r.Encrypted():
			fmt.Printf("bdf %d opset %d unit %s title %q (encrypted)\n", m.BDF, m.Opset, m.Unit, m.Meta.DC.Title.First())
		default:
			fmt.Printf("bdf %d opset %d unit %s title %q\n", m.BDF, m.Opset, m.Unit, m.Meta.DC.Title.First())
		}
		for _, v := range m.Views {
			switch v.Kind {
			case bdf.ViewSheet:
				fmt.Printf("view %s (%s) %q tiles=%d\n", v.ID, v.Kind, v.Title, len(v.Tiles))
			default:
				fmt.Printf("view %s (%s) %q pages=%d\n", v.ID, v.Kind, v.Title, len(v.Pages))
			}
		}
		var stored, size int
		for _, p := range m.Parts {
			fmt.Printf("%s %-4s %-11s %8d -> %8d\n", p.H, p.T, p.Enc, p.Len, p.Size)
			stored += p.Len
			size += p.Size
		}
		fmt.Printf("%d parts, %d bytes stored, %d bytes decoded\n", len(m.Parts), stored, size)
	case "manifest":
		r, err := open(os.Args[2])
		check(err)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		check(enc.Encode(r.Manifest))
	case "disasm":
		if len(os.Args) < 4 {
			usage()
		}
		r, err := open(os.Args[2])
		check(err)
		s, err := bdf.Disassemble(part(r, os.Args[3]))
		check(err)
		fmt.Print(s)
	case "extract":
		if len(os.Args) < 5 {
			usage()
		}
		r, err := open(os.Args[2])
		check(err)
		check(os.WriteFile(os.Args[4], part(r, os.Args[3]), 0o644))
	case "split":
		if len(os.Args) < 4 {
			usage()
		}
		r, err := bdf.OpenSingleFile(os.Args[2])
		check(err)
		check(r.WriteSplit(os.Args[3]))
	case "join":
		if len(os.Args) < 4 {
			usage()
		}
		r, err := bdf.OpenSplit(os.Args[2])
		check(err)
		f, err := os.Create(os.Args[3])
		check(err)
		check(errors.Join(r.WriteSingle(f), f.Close()))
	case "encrypt", "decrypt":
		crypt(os.Args[1], os.Args[2:])
	case "demo":
		d, err := fixture.Demo()
		check(err)
		check(write(d, os.Args[2]))
	default:
		usage()
	}
}

// crypt encrypts a document with a password, or writes an encrypted one
// without its encryption.
func crypt(cmd string, args []string) {
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	passwordFile := fs.String("password-file", "", "read the password from this file (- for the standard input; default: $"+passwordEnv+")")
	fs.Parse(args)
	if fs.NArg() != 2 {
		usage()
	}
	password, err := readPassword(*passwordFile)
	check(err)
	if password == "" {
		check(fmt.Errorf("%s: no password (-password-file or $%s)", cmd, passwordEnv))
	}
	r, err := openStored(fs.Arg(0))
	check(err)
	if cmd == "encrypt" {
		if r.Encrypted() {
			check(fmt.Errorf("%s is encrypted already", fs.Arg(0)))
		}
	} else {
		if !r.Encrypted() {
			check(fmt.Errorf("%s is not encrypted", fs.Arg(0)))
		}
		check(unlock(r, password, fs.Arg(0)))
	}
	d, err := r.ToDocument()
	check(err)
	if cmd == "encrypt" {
		d.Lock, err = bdf.NewPasswordLock(password, 0)
		check(err)
	}
	check(write(d, fs.Arg(1)))
}

func writeSingle(d *bdf.Document, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := d.WriteSingle(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
