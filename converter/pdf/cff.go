package pdf

import (
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/internal/cff"
)

// The CFF font programs of PDF files are read and subset by
// internal/cff, which the font embedding of the other converters
// shares.
type (
	cffFont      = cff.Font
	cffDictEntry = cff.DictEntry
)

var (
	parseCFF           = cff.Parse
	subsetCFF          = cff.Subset
	cffEncodeDict      = cff.EncodeDict
	cffIndex           = cff.Index
	cffStandardStrings = cff.StandardStrings
	standardEncoding   = cff.StandardEncoding
)

// glyphNameToRune maps a glyph name to a Unicode code point using the AGL and
// the uniXXXX / uXXXX[XX] conventions.
func glyphNameToRune(name string) (rune, bool) {
	if name == "" {
		return 0, false
	}
	if i := strings.IndexByte(name, '.'); i > 0 {
		name = name[:i]
	}
	if r, ok := glyphList[name]; ok {
		return r, true
	}
	if strings.HasPrefix(name, "uni") && len(name) >= 7 {
		if v, err := strconv.ParseUint(name[3:7], 16, 32); err == nil {
			return rune(v), true
		}
	}
	if strings.HasPrefix(name, "u") && len(name) >= 5 && len(name) <= 7 {
		if v, err := strconv.ParseUint(name[1:], 16, 32); err == nil {
			return rune(v), true
		}
	}
	// Names like "Cxx"/"Gxx" carry no Unicode.
	return 0, false
}

// glyphNameIndex parses names of the form gNN, glyphNN, cidNN, GNN, index NN.
func glyphNameIndex(name string) (int, bool) {
	for _, prefix := range []string{"glyph", "cid", "index", "g", "G", "c", "C"} {
		if strings.HasPrefix(name, prefix) {
			if v, err := strconv.Atoi(name[len(prefix):]); err == nil {
				return v, true
			}
		}
	}
	return 0, false
}
