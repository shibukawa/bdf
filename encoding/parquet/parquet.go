// Package parquet writes and reads the tables a document holds beside its
// drawing (docs/spec.md §4.6): Apache Parquet files of flat columns.
//
// Write makes the files the converters store: the columns of a table as
// optional fields of a flat schema, PLAIN values, RLE definition levels,
// one row group, no compression (the part the file becomes is compressed
// as a whole). Read takes those files back, and the files of other
// writers whose columns are flat and whose pages are PLAIN or dictionary
// encoded and uncompressed, Snappy or gzip (what pandas, Arrow and DuckDB
// write by default). Nested columns, the delta encodings and the other
// codecs are left to converter/parquet, which reads every file.
package parquet

import (
	"errors"
	"fmt"
	"time"
)

// Type is the type of a column (docs/spec.md §4.6).
type Type uint8

// The types of columns. Other is a type of Parquet that is none of the
// rest (a nested column, a decimal of bytes, …), as the file was read.
const (
	Other Type = iota
	String
	Number
	Integer
	Boolean
	Date
	DateTime
)

var typeNames = [...]string{"other", "string", "number", "integer", "boolean", "date", "datetime"}

func (t Type) String() string {
	if int(t) < len(typeNames) {
		return typeNames[t]
	}
	return fmt.Sprintf("type(%d)", uint8(t))
}

// ParseType reads the name of a type, as the manifest writes it.
func ParseType(s string) (Type, bool) {
	for i, n := range typeNames {
		if n == s {
			return Type(i), true
		}
	}
	return Other, false
}

// Column is a column of a table.
type Column struct {
	Name string
	Type Type
}

// Table is a table in memory: its columns, and its rows of values. A value
// is a string (String), a float64 (Number), an int64 (Integer), a bool
// (Boolean), a time.Time (Date, at 0:00 UTC; DateTime, UTC) or nil (a
// null); a value of a column of type Other is what Read made of it (a
// []byte, a string, a number), or nil.
type Table struct {
	Columns []Column
	Rows    [][]any
}

// Limits of what is read and written.
const (
	// MaxColumns bounds the columns of a table (docs/spec.md §4.6).
	MaxColumns = 1 << 14
	// DefaultMaxRows is the rows Read returns when Limits.MaxRows is 0.
	DefaultMaxRows = 1 << 20
	// maxFooter bounds the footer of a file.
	maxFooter = 1 << 26
	// maxPage bounds a page of a file.
	maxPage = 1 << 28
)

// Limits bounds what Read returns.
type Limits struct {
	// MaxRows is the most rows returned (DefaultMaxRows when 0; the rows
	// after them are left out).
	MaxRows int
}

var (
	// ErrNotParquet is returned for data that is not a Parquet file.
	ErrNotParquet = errors.New("parquet: not a Parquet file")
	errCorrupt    = errors.New("parquet: corrupt file")
)

// unixDay is the start of the epoch, which dates count from.
var unixDay = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
