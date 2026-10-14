package form

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/csv"
	"github.com/shibukawa/bdf/converter/xlsx"
	"github.com/shibukawa/bdf/form/tag"
)

// Source is a file of data for a template.
type Source struct {
	// Name is where the data of the file goes: "" makes it the data
	// itself (a JSON or YAML object, whose values the values of the
	// sources before it with the same names give way to), and a name
	// ("items", "customer.address") makes it the value of that name.
	//
	// A CSV or TSV file is a list of objects: the first record names the
	// values of the others. Without a name it is the list named after the
	// file ("lines" for lines.csv). A workbook (.xlsx) is such a list for
	// each of its sheets, by the names of the sheets; with a name, a
	// workbook of one sheet is that list.
	Name string
	// FileName tells the format of the data by its extension: .json,
	// .yaml or .yml, .csv or .tsv, .xlsx.
	FileName string
	Data     []byte
}

// MaxDataSize bounds the bytes of a file of data.
const MaxDataSize = 256 << 20

// LoadData reads files of data into the data of a template: the values of
// JSON, with numbers as float64 and dates (of YAML, and of workbooks) as
// time.Time.
func LoadData(sources ...Source) (map[string]any, error) {
	root := map[string]any{}
	for _, s := range sources {
		if len(s.Data) > MaxDataSize {
			return nil, fmt.Errorf("form: %s is larger than %d MiB", s.FileName, MaxDataSize>>20)
		}
		v, err := readSource(s)
		if err != nil {
			return nil, fmt.Errorf("form: %s: %w", s.FileName, err)
		}
		if s.Name != "" {
			if err := put(root, s.Name, v); err != nil {
				return nil, fmt.Errorf("form: %s: %w", s.FileName, err)
			}
			continue
		}
		obj, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("form: %s: the data of a template is an object; give a list or a value a name (name=%s)", s.FileName, s.FileName)
		}
		for k, e := range obj {
			root[k] = e
		}
	}
	return root, nil
}

// put sets the value of a name ("a.b") in an object.
func put(root map[string]any, name string, v any) error {
	parts := strings.Split(name, ".")
	for i, p := range parts {
		if p = strings.TrimSpace(p); p == "" {
			return fmt.Errorf("the name %q has an empty part", name)
		}
		if i == len(parts)-1 {
			root[p] = v
			break
		}
		next, ok := root[p].(map[string]any)
		if !ok {
			next = map[string]any{}
			root[p] = next
		}
		root = next
	}
	return nil
}

// readSource reads a file of data: an object for a file that is the data
// itself, else its value.
func readSource(s Source) (any, error) {
	ext := strings.ToLower(path.Ext(strings.ReplaceAll(s.FileName, "\\", "/")))
	base := strings.TrimSuffix(path.Base(strings.ReplaceAll(s.FileName, "\\", "/")), path.Ext(s.FileName))
	switch ext {
	case ".json":
		dec := json.NewDecoder(bytes.NewReader(s.Data))
		dec.UseNumber()
		var v any
		if err := dec.Decode(&v); err != nil {
			return nil, err
		}
		if dec.More() {
			return nil, fmt.Errorf("more than one JSON value")
		}
		return tag.Normalize(v), nil
	case ".yaml", ".yml":
		v, err := decodeYAML(s.Data)
		if err != nil {
			return nil, err
		}
		return tag.Normalize(v), nil
	case ".csv", ".tsv", ".txt":
		recs, _, err := csv.Records(s.Data, &csv.Options{Header: csv.HeaderYes})
		if err != nil {
			return nil, err
		}
		rows := make([][]any, len(recs))
		for i, rec := range recs {
			rows[i] = make([]any, len(rec))
			for j, f := range rec {
				rows[i][j] = csvValue(f)
			}
		}
		list, err := table(rows)
		if err != nil {
			return nil, err
		}
		if s.Name == "" {
			return map[string]any{base: list}, nil
		}
		return list, nil
	case ".xlsx", ".xlsm":
		sheets, err := xlsx.Values(bytes.NewReader(s.Data), int64(len(s.Data)), nil)
		if err != nil {
			return nil, err
		}
		obj := map[string]any{}
		var only any
		for _, sh := range sheets {
			list, err := table(sh.Rows)
			if err != nil {
				return nil, fmt.Errorf("sheet %q: %w", sh.Name, err)
			}
			obj[sh.Name], only = list, list
		}
		if s.Name != "" && len(sheets) == 1 {
			return only, nil
		}
		return obj, nil
	}
	return nil, fmt.Errorf("data is JSON (.json), YAML (.yaml), CSV (.csv, .tsv) or a workbook (.xlsx)")
}

// csvValue is a value of a CSV file as data: a number that is written as
// numbers are (not one with zeros before it, or separators: a code), true
// and false, nothing for an empty value, else the text.
func csvValue(s string) any {
	t := strings.TrimSpace(s)
	switch strings.ToLower(t) {
	case "":
		return nil
	case "true":
		return true
	case "false":
		return false
	}
	if plainNumber(t) {
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f
		}
	}
	return s
}

// plainNumber reports whether a text is a number as JSON writes one,
// without an exponent.
func plainNumber(s string) bool {
	if strings.HasPrefix(s, "-") {
		s = s[1:]
	}
	whole, frac, dot := strings.Cut(s, ".")
	if whole == "" || len(whole) > 15 || len(whole) > 1 && whole[0] == '0' || dot && frac == "" {
		return false
	}
	for _, r := range whole + frac {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// table makes the rows of a table a list of objects: the first row that
// has values names the values of the rows after it. Rows without values
// are left out.
func table(rows [][]any) ([]any, error) {
	var keys []string
	out := []any{}
	for _, row := range rows {
		empty := true
		for _, v := range row {
			empty = empty && v == nil
		}
		if empty {
			continue
		}
		if keys == nil {
			seen := map[string]bool{}
			for _, v := range row {
				k := ""
				if v != nil {
					k = strings.TrimSpace(fmt.Sprint(v))
				}
				if k != "" && seen[k] {
					return nil, fmt.Errorf("two columns are named %q", k)
				}
				seen[k] = true
				keys = append(keys, k)
			}
			continue
		}
		obj := make(map[string]any, len(keys))
		for i, k := range keys {
			if k == "" {
				continue
			}
			if i < len(row) {
				obj[k] = row[i]
			} else {
				obj[k] = nil
			}
		}
		out = append(out, obj)
	}
	return out, nil
}
