// Command bdf generates, inspects and converts BDF files.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/fixture"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  bdf generate [flags] <in.pdf | in.pptx> <out.bdf | dir/>
                                     convert a PDF or PowerPoint file (bdf generate -h for flags)
  bdf ls <file.bdf | dir>            list views and parts
  bdf manifest <file.bdf | dir>      print the manifest as JSON
  bdf disasm <file.bdf | dir> <hash> disassemble an object part
  bdf extract <file.bdf | dir> <hash> <out>
  bdf split <file.bdf> <dir>         write the split form
  bdf join <dir> <file.bdf>          write the single-file form
  bdf demo <file.bdf | dir/>         write the fixture document (dir/ ends with a slash)`)
	os.Exit(2)
}

func open(path string) (*bdf.Reader, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if st.IsDir() {
		return bdf.OpenSplit(path)
	}
	return bdf.OpenSingleFile(path)
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
		fmt.Printf("bdf %d opset %d unit %s title %q\n", m.BDF, m.Opset, m.Unit, m.Meta.Title)
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
		h, err := bdf.ParseHash(os.Args[3])
		check(err)
		b, err := r.Part(h)
		check(err)
		s, err := bdf.Disassemble(b)
		check(err)
		fmt.Print(s)
	case "extract":
		if len(os.Args) < 5 {
			usage()
		}
		r, err := open(os.Args[2])
		check(err)
		h, err := bdf.ParseHash(os.Args[3])
		check(err)
		b, err := r.Part(h)
		check(err)
		check(os.WriteFile(os.Args[4], b, 0o644))
	case "split":
		if len(os.Args) < 4 {
			usage()
		}
		r, err := bdf.OpenSingleFile(os.Args[2])
		check(err)
		d, err := r.ToDocument()
		check(err)
		check(d.WriteSplit(os.Args[3]))
	case "join":
		if len(os.Args) < 4 {
			usage()
		}
		r, err := bdf.OpenSplit(os.Args[2])
		check(err)
		d, err := r.ToDocument()
		check(err)
		check(writeSingle(d, os.Args[3]))
	case "demo":
		d, err := fixture.Demo()
		check(err)
		out := os.Args[2]
		if out[len(out)-1] == '/' || out[len(out)-1] == filepath.Separator {
			check(d.WriteSplit(out))
		} else {
			check(writeSingle(d, out))
		}
	default:
		usage()
	}
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
