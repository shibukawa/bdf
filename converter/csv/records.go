package csv

import "fmt"

// Records reads the records of a CSV or TSV file as text, without
// converting it: the character encoding and the dialect are found as
// Convert finds them (of the options, Charset, Delimiter, Quote and Header
// apply). The result tells what was found and whether the first record is
// a header row; its Doc is nil.
func Records(data []byte, opts *Options) ([][]string, *Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if err := checkRune("delimiter", opts.Delimiter); err != nil {
		return nil, nil, err
	}
	if opts.Quote != NoQuote {
		if err := checkRune("quote", opts.Quote); err != nil {
			return nil, nil, err
		}
	}
	res := &Result{}
	cs, bom := bomCharset(data)
	switch {
	case bom:
	case opts.Charset != "" && opts.Charset != "auto":
		var err error
		if cs, err = namedCharset(opts.Charset); err != nil {
			return nil, nil, fmt.Errorf("csv: %w", err)
		}
	default:
		cs = detectCharset(data)
	}
	res.Charset, res.BOM = cs.name, bom
	text := decode(data, cs)
	sample, cut := sampleOf(text)
	res.Dialect = sniff(sample, cut, opts.Delimiter, opts.Quote).d
	rd := newReader(text, res.Dialect)
	var rows [][]string
	for {
		f, ok := rd.record(nil)
		if !ok {
			break
		}
		if rd.unfinished {
			res.Warnings = append(res.Warnings, fmt.Sprintf("record %d: a quote is never closed; the rest of the file is one value", len(rows)+1))
		}
		rows = append(rows, f)
		res.Cols = max(res.Cols, len(f))
	}
	res.Rows = len(rows)
	switch opts.Header {
	case HeaderYes:
		res.Header = len(rows) > 0
	case HeaderAuto:
		res.Header = hasHeader(rows[:min(len(rows), maxHeaderSample)])
	}
	return rows, res, nil
}
