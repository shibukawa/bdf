package parquet

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	conv "github.com/shibukawa/bdf/converter"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "parquet",
		Description: "Apache Parquet",
		Extensions:  []string{".parquet", ".parq", ".pqt"},
		Params: []conv.Param{
			{Name: "rows", Usage: fmt.Sprintf("the number of rows shown, or all (default: %d)", DefaultRows)},
			{Name: "types", Usage: "true or false: a row of the column types under their names (default: true)"},
			{Name: "table", Usage: "format the sheet as an Excel table with a built-in style, e.g. TableStyleMedium2"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return detect(head, r, size) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			opts := &Options{TableStyle: o.Param("table"), Title: o.Title,
				FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn}
			if o.FileName != "" {
				opts.Name = sheetName(o.FileName)
			}
			switch v := strings.ToLower(o.Param("rows")); v {
			case "":
			case "all":
				opts.Rows = -1
			default:
				n, err := strconv.Atoi(v)
				if err != nil || n <= 0 {
					return nil, fmt.Errorf("parquet: parameter rows: %q is not a positive number or all", v)
				}
				opts.Rows = n
			}
			if v := o.Param("types"); v != "" {
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("parquet: parameter types: %q is not true or false", v)
				}
				opts.NoTypes = !b
			}
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: summary(res)}, nil
		},
	})
}

// detect recognizes the magic number at both ends of the file (PARE is
// that of files whose footer is encrypted, which are refused as such).
func detect(head []byte, r io.ReaderAt, size int64) bool {
	if size < 12 || !bytes.HasPrefix(head, []byte("PAR1")) && !bytes.HasPrefix(head, []byte("PARE")) {
		return false
	}
	var tail [4]byte
	if _, err := r.ReadAt(tail[:], size-4); err != nil && err != io.EOF {
		return false
	}
	return string(tail[:]) == "PAR1" || string(tail[:]) == "PARE"
}

// summary describes a conversion in a line.
func summary(res *Result) string {
	rows := fmt.Sprintf("%d row(s)", max(res.Rows, res.Shown))
	if res.Shown < res.Rows {
		rows = fmt.Sprintf("%d of %d row(s)", res.Shown, res.Rows)
	}
	var what []string
	what = append(what, fmt.Sprintf("%d row group(s)", res.RowGroups))
	if len(res.Codecs) > 0 {
		what = append(what, strings.Join(res.Codecs, ", "))
	}
	if res.CreatedBy != "" {
		what = append(what, "written by "+res.CreatedBy)
	}
	return fmt.Sprintf("%s × %d column(s) (%s), %d embedded font(s)", rows, res.Cols, strings.Join(what, "; "), res.EmbeddedFonts)
}
