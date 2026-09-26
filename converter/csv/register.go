package csv

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	conv "github.com/shibukawa/bdf/converter"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "csv",
		Description: "CSV or TSV (comma-, tab-, semicolon- or bar-separated values)",
		Extensions:  []string{".csv", ".tsv", ".tab"},
		Params: []conv.Param{
			{Name: "charset", Usage: "character encoding: utf-8, shift_jis, euc-jp, euc-kr, gb18030, big5, windows-1252 … (default: detected)"},
			{Name: "delimiter", Usage: "comma, tab, semicolon, pipe or another character (default: detected)"},
			{Name: "quote", Usage: "double, single or none (default: detected)"},
			{Name: "header", Usage: "true or false: the first row is a header row (default: guessed)"},
			{Name: "table", Usage: "format the sheet as an Excel table with a built-in style, e.g. TableStyleMedium2"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return detect(head, r, size) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			opts := &Options{Charset: o.Param("charset"), TableStyle: o.Param("table"), Title: o.Title,
				FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn}
			if o.FileName != "" {
				opts.Name = sheetName(o.FileName)
			}
			var err error
			if opts.Delimiter, err = delimiterParam(o.Param("delimiter")); err != nil {
				return nil, err
			}
			if opts.Quote, err = quoteParam(o.Param("quote")); err != nil {
				return nil, err
			}
			switch v := strings.ToLower(o.Param("header")); v {
			case "", "auto":
			default:
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("csv: parameter header: %q is not true, false or auto", v)
				}
				opts.Header = HeaderNo
				if b {
					opts.Header = HeaderYes
				}
			}
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: summary(res)}, nil
		},
	})
}

// summary describes a conversion in a line.
func summary(res *Result) string {
	d := res.Dialect
	var what []string
	what = append(what, res.Charset)
	if res.BOM {
		what[0] += " with BOM"
	}
	what = append(what, delimiterName(d.Delimiter)+"-separated")
	switch d.Quote {
	case 0:
		what = append(what, "no quotes")
	case '\'':
		what = append(what, "single quotes")
	}
	if d.Escape != 0 {
		what = append(what, "backslash escapes")
	}
	if res.Header {
		what = append(what, "header row")
	} else {
		what = append(what, "no header row")
	}
	return fmt.Sprintf("%d row(s) × %d column(s) (%s), %d embedded font(s)", res.Rows, res.Cols, strings.Join(what, ", "), res.EmbeddedFonts)
}

var delimiterNames = map[string]rune{"comma": ',', "tab": '\t', "semicolon": ';', "pipe": '|', "bar": '|', "space": ' ', `\t`: '\t'}

func delimiterName(r rune) string {
	for _, n := range []string{"comma", "tab", "semicolon", "pipe", "space"} {
		if delimiterNames[n] == r {
			return n
		}
	}
	return strconv.QuoteRune(r)
}

func delimiterParam(v string) (rune, error) {
	if v == "" || v == "auto" {
		return 0, nil
	}
	if r, ok := delimiterNames[strings.ToLower(v)]; ok {
		return r, nil
	}
	if r, size := utf8.DecodeRuneInString(v); size == len(v) && r < utf8.RuneSelf {
		return r, nil
	}
	return 0, fmt.Errorf("csv: parameter delimiter: %q is not comma, tab, semicolon, pipe or an ASCII character", v)
}

func quoteParam(v string) (rune, error) {
	switch strings.ToLower(v) {
	case "", "auto":
		return 0, nil
	case "double", `"`:
		return '"', nil
	case "single", "'":
		return '\'', nil
	case "none":
		return NoQuote, nil
	}
	return 0, fmt.Errorf("csv: parameter quote: %q is not double, single or none", v)
}

// detect recognizes text whose records have the same number of fields,
// at least two, in the first 64 KiB.
func detect(head []byte, r io.ReaderAt, size int64) bool {
	if len(head) == 0 {
		return false
	}
	for _, sig := range []string{"%PDF", "PK\x03\x04", "{\\rtf", "<?xml", "<!DOCTYPE", "<!doctype", "<html", "ISO-10303-21"} {
		if bytes.HasPrefix(head, []byte(sig)) {
			return false
		}
	}
	data := head
	if size > int64(len(head)) {
		data = make([]byte, min(size, sniffLen))
		n, _ := r.ReadAt(data, 0)
		data = data[:n]
	}
	cs, ok := bomCharset(data)
	if !ok {
		if cs, ok = utf16Guess(data); !ok {
			if bytes.IndexByte(data, 0) >= 0 {
				return false
			}
			cs = detectCharset(data)
		}
	}
	text := decode(trimPartial(data, int64(len(data)) < size), cs)
	bad, n := 0, 0
	for _, c := range text {
		n++
		if c < 0x20 && c != '\t' && c != '\n' && c != '\r' && c != '\f' || c == 0x7f || c == utf8.RuneError {
			bad++
		}
	}
	if bad*100 > n {
		return false
	}
	t := sniff(text, int64(len(data)) < size, 0, 0)
	return t.fields >= 2 && t.records >= 2 && t.score >= 0.9
}
