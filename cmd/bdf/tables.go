package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shibukawa/bdf"
	fullparquet "github.com/shibukawa/bdf/converter/parquet"
	"github.com/shibukawa/bdf/encoding/parquet"
)

// tablesCmd lists the tables a document holds beside its drawing.
func tablesCmd(args []string) {
	fs := flag.NewFlagSet("tables", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "write the list as JSON (the tables member of the manifest)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf tables [flags] <file.bdf | dir>\n  lists the tables the document holds beside its drawing (the data of charts, sheets, CSV and Parquet files): their ids,\n  rows and columns, titles and places")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}
	r, err := open(fs.Arg(0))
	check(err)
	defer r.Close()
	if r.Locked() {
		check(fmt.Errorf("%s is encrypted: set $%s", fs.Arg(0), passwordEnv))
	}
	if *asJSON {
		writeJSON(os.Stdout, r.Manifest.Tables)
		return
	}
	if len(r.Manifest.Tables) == 0 {
		fmt.Println("no tables")
		return
	}
	for _, t := range r.Manifest.Tables {
		where := ""
		switch {
		case t.View != "" && t.Page > 0:
			where = fmt.Sprintf(" (%s, page %d)", t.View, t.Page)
		case t.View != "":
			where = fmt.Sprintf(" (%s)", t.View)
		}
		title := ""
		if t.Title != "" {
			title = " " + quoted(t.Title)
		}
		fmt.Printf("%s: %d rows × %d columns%s%s\n", t.ID, t.Rows, len(t.Columns), title, where)
		cols := make([]string, len(t.Columns))
		for i, c := range t.Columns {
			cols[i] = c.Name + ":" + c.Type
		}
		fmt.Printf("  %s\n", strings.Join(cols, ", "))
	}
}

func quoted(s string) string { return fmt.Sprintf("%q", s) }

// tableCmd writes a table of a document.
func tableCmd(args []string) {
	fs := flag.NewFlagSet("table", flag.ExitOnError)
	format := fs.String("format", "csv", "csv, json, or parquet (the file as the document holds it)")
	rows := fs.Int("rows", 0, "the most rows written (default: all, up to 1048576)")
	allow := fs.Bool("allow-plaintext", false, "read an encrypted document (with $"+passwordEnv+"); the table is written in plain text")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf table [flags] <file.bdf | dir> <id> [out]\n  writes a table the document holds (bdf tables lists them) as CSV, JSON or Parquet (default: CSV to the standard output)")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() < 2 || fs.NArg() > 3 {
		fs.Usage()
		os.Exit(2)
	}
	switch *format {
	case "csv", "json", "parquet":
	default:
		badUsage("table", "-format must be csv, json or parquet")
	}
	doc := loadDocument(fs.Arg(0), *allow)
	t := doc.Table(fs.Arg(1))
	if t == nil {
		check(fmt.Errorf("%s has no table %q (bdf tables lists them)", fs.Arg(0), fs.Arg(1)))
	}
	h, err := bdf.ParseHash(t.Part)
	check(err)
	part := doc.Part(h)
	if part == nil {
		check(fmt.Errorf("table %s: its part %s is missing", t.ID, t.Part))
	}
	var w io.Writer = os.Stdout
	var f *os.File
	if fs.NArg() == 3 {
		f, err = os.Create(fs.Arg(2))
		check(err)
		w = f
	}
	bw := bufio.NewWriter(w)
	switch *format {
	case "parquet":
		_, err = bw.Write(part.Data)
	default:
		tb, rerr := readTable(part.Data, *rows)
		check(rerr)
		if *format == "csv" {
			err = parquet.WriteCSV(bw, tb)
		} else {
			err = parquet.WriteJSON(bw, tb)
		}
	}
	if err == nil {
		err = bw.Flush()
	}
	if f != nil {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}
	check(err)
}

// readTable reads the rows of a table's Parquet file: with encoding/parquet,
// else, for a file of encodings or types it does not read (one stored as
// it came), with the converter's reader, as text.
func readTable(data []byte, rows int) (*parquet.Table, error) {
	t, err := parquet.Read(data, parquet.Limits{MaxRows: rows})
	if err == nil {
		return t, nil
	}
	if !strings.Contains(err.Error(), "converter/parquet") && !strings.Contains(err.Error(), "repeated") {
		return nil, err
	}
	if rows <= 0 {
		rows = -1
	}
	return fullparquet.ReadTable(bytes.NewReader(data), int64(len(data)), rows)
}
