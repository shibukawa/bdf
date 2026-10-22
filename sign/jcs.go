package sign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// The JSON Canonicalization Scheme (RFC 8785): the bytes of JSON that are
// the same for the same values, whoever wrote them. The manifest of a
// document is signed as its canonical JSON, so that the single and the
// split forms (which write the manifest differently) have the same
// signature, and so that a reader in another language can check it.

// maxDepth bounds the nesting of what is canonicalized.
const maxDepth = 512

// Canonical writes the canonical JSON of JSON.
func Canonical(b []byte) ([]byte, error) {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var v any
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	if d.More() {
		return nil, errors.New("sign: more than one value of JSON")
	}
	var out bytes.Buffer
	if err := canonical(&out, v, 0); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// CanonicalValue writes the canonical JSON of a value of JSON (as
// encoding/json decodes it, with json.Number or float64 numbers).
func CanonicalValue(v any) ([]byte, error) {
	var out bytes.Buffer
	if err := canonical(&out, v, 0); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func canonical(w *bytes.Buffer, v any, depth int) error {
	if depth > maxDepth {
		return errors.New("sign: the JSON is nested too deeply")
	}
	switch x := v.(type) {
	case nil:
		w.WriteString("null")
	case bool:
		w.WriteString(strconv.FormatBool(x))
	case json.Number:
		f, err := strconv.ParseFloat(string(x), 64)
		if err != nil {
			return fmt.Errorf("sign: the number %s: %w", x, err)
		}
		return number(w, f)
	case float64:
		return number(w, x)
	case int:
		return number(w, float64(x))
	case string:
		str(w, x)
	case []any:
		w.WriteByte('[')
		for i, e := range x {
			if i > 0 {
				w.WriteByte(',')
			}
			if err := canonical(w, e, depth+1); err != nil {
				return err
			}
		}
		w.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		// (by their code units of UTF-16)
		sort.Slice(keys, func(i, j int) bool { return less16(keys[i], keys[j]) })
		w.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				w.WriteByte(',')
			}
			str(w, k)
			w.WriteByte(':')
			if err := canonical(w, x[k], depth+1); err != nil {
				return err
			}
		}
		w.WriteByte('}')
	default:
		return fmt.Errorf("sign: %T is not a value of JSON", v)
	}
	return nil
}

func less16(a, b string) bool {
	ua, ub := utf16.Encode([]rune(a)), utf16.Encode([]rune(b))
	for i := 0; i < len(ua) && i < len(ub); i++ {
		if ua[i] != ub[i] {
			return ua[i] < ub[i]
		}
	}
	return len(ua) < len(ub)
}

// str writes a string as RFC 8785 does: only what must be escaped is.
func str(w *bytes.Buffer, s string) {
	w.WriteByte('"')
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		i += n
		switch r {
		case '"':
			w.WriteString(`\"`)
		case '\\':
			w.WriteString(`\\`)
		case '\b':
			w.WriteString(`\b`)
		case '\f':
			w.WriteString(`\f`)
		case '\n':
			w.WriteString(`\n`)
		case '\r':
			w.WriteString(`\r`)
		case '\t':
			w.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(w, `\u%04x`, r)
			} else {
				w.WriteRune(r)
			}
		}
	}
	w.WriteByte('"')
}

// number writes a number as ECMAScript's Number.prototype.toString does.
func number(w *bytes.Buffer, f float64) error {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return errors.New("sign: a number that JSON does not have")
	}
	if f == 0 {
		w.WriteByte('0')
		return nil
	}
	if f < 0 {
		w.WriteByte('-')
		f = -f
	}
	// the shortest digits, and where the point is after the first
	e := strconv.FormatFloat(f, 'e', -1, 64)
	mant, exp, _ := strings.Cut(e, "e")
	digits := strings.Replace(mant, ".", "", 1)
	x, _ := strconv.Atoi(exp)
	n, k := x+1, len(digits)
	switch {
	case k <= n && n <= 21:
		w.WriteString(digits + strings.Repeat("0", n-k))
	case 0 < n && n <= 21:
		w.WriteString(digits[:n] + "." + digits[n:])
	case -6 < n && n <= 0:
		w.WriteString("0." + strings.Repeat("0", -n) + digits)
	default:
		sign := "+"
		if n-1 < 0 {
			sign = "-"
		}
		abs := n - 1
		if abs < 0 {
			abs = -abs
		}
		if k == 1 {
			w.WriteString(digits + "e" + sign + strconv.Itoa(abs))
		} else {
			w.WriteString(digits[:1] + "." + digits[1:] + "e" + sign + strconv.Itoa(abs))
		}
	}
	return nil
}
