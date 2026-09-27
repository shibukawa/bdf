// Package hpglsniff recognizes HP-GL and HP-GL/2 plot files, which have no
// signature: a bare plot file is a sequence of two-letter instructions
// ("IN;SP1;PA0,0;PD1000,0;"), and a plot sent to a printer or a large-format
// plotter is wrapped in a PJL job with PCL or HP RTL escape sequences
// around the HP-GL/2 and the raster data. The CSV converter asks it too, so
// that it leaves plot files (whose instructions end with semicolons) alone.
package hpglsniff

import (
	"bytes"
	"strings"
)

// Is reports whether head, the first bytes of an input, starts an HP-GL or
// HP-GL/2 plot, bare or in a PJL, PCL or HP RTL job.
func Is(head []byte) bool {
	i := skipBlank(head, 0)
	if i >= len(head) {
		return false
	}
	if head[i] == 0x1b || bytes.HasPrefix(head[i:], []byte("@PJL")) {
		return job(head[i:])
	}
	return Instructions(head[i:]) >= 3
}

// skipBlank skips white space, NULs and the device-control sequences of
// serial plotters (ESC . and a character, with parameters up to a colon).
func skipBlank(b []byte, i int) int {
	for i < len(b) {
		switch c := b[i]; {
		case c == 0 || c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == ';' || c == 0x03:
			i++
		case c == 0x1b && i+1 < len(b) && b[i+1] == '.':
			i = SkipDeviceControl(b, i)
		default:
			return i
		}
	}
	return i
}

// SkipDeviceControl returns the index after the device-control sequence
// (ESC . and a letter, and for some letters parameters up to a colon) at i.
func SkipDeviceControl(b []byte, i int) int {
	i += 2
	if i >= len(b) {
		return i
	}
	i++
	// parameters: digits and semicolons up to a colon
	j := i
	for j < len(b) && j-i < 64 && (b[j] >= '0' && b[j] <= '9' || b[j] == ';' || b[j] == ' ') {
		j++
	}
	if j < len(b) && b[j] == ':' {
		return j + 1
	}
	return i
}

// job reports whether a PJL, PCL or HP RTL job holds HP-GL/2 or HP RTL:
// PJL that enters one of their languages, the escape sequence that enters
// HP-GL/2, or the commands that set up raster data.
func job(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] == '@' && bytes.HasPrefix(b[i:], []byte("@PJL")) {
			end := bytes.IndexAny(b[i:], "\r\n")
			if end < 0 {
				end = len(b) - i
			}
			line := strings.ToUpper(string(b[i : i+end]))
			if strings.Contains(line, "ENTER") && strings.Contains(line, "LANGUAGE") {
				lang := strings.TrimSpace(line[strings.Index(line, "LANGUAGE")+len("LANGUAGE"):])
				lang = strings.TrimSpace(strings.TrimPrefix(lang, "="))
				switch {
				case strings.HasPrefix(lang, "HPGL"), strings.HasPrefix(lang, "RTL"):
					return true
				case strings.HasPrefix(lang, "PCL"):
				default:
					return false // PostScript, PDF …
				}
			}
			i += end
			continue
		}
		if b[i] != 0x1b || i+1 >= len(b) {
			continue
		}
		switch b[i+1] {
		case '%':
			// ESC % # B enters HP-GL/2 (ESC % -12345 X is the UEL)
			j := i + 2
			if j < len(b) && (b[j] == '-' || b[j] == '+') {
				j++
			}
			k := j
			for k < len(b) && b[k] >= '0' && b[k] <= '9' {
				k++
			}
			if k > j && k < len(b) && b[k] == 'B' {
				return true
			}
		case '*':
			// raster set-up: configure image data (ESC * v # W), start
			// raster graphics (ESC * r # A), transfer (ESC * b # W)
			j := i + 2
			if j < len(b) && (b[j] == 'v' || b[j] == 'r' || b[j] == 'b') {
				k := j + 1
				for k < len(b) && (b[k] >= '0' && b[k] <= '9' || b[k] == '-' || b[k] == '+' || b[k] == '.') {
					k++
				}
				if k < len(b) && (b[j] == 'v' && b[k] == 'W' || b[j] == 'r' && b[k] == 'A' || b[j] == 'b' && b[k] == 'W') {
					return true
				}
			}
		}
	}
	return false
}

// Instructions returns how many HP-GL instructions b starts with, 0 when
// something in it is not HP-GL: an unknown mnemonic, or characters that
// are not parameters. A sequence cut short at the end of b counts.
func Instructions(b []byte) int {
	n := 0
	numeric := false // an instruction had parameters
	term := byte(0x03)
	i := 0
	for {
		i = skipBlank(b, i)
		if i >= len(b) {
			break
		}
		c := b[i]
		if c == 0x1b {
			break // the HP-GL/2 part ends (PCL, HP RTL)
		}
		if c == '!' {
			// device-specific instructions of cutting plotters (Roland)
			i++
			for i < len(b) && isLetter(b[i]) {
				i++
			}
			i = skipParams(b, i)
			if i < 0 {
				return 0
			}
			continue
		}
		if !isLetter(c) {
			return 0
		}
		if i+1 >= len(b) {
			break
		}
		if !isLetter(b[i+1]) {
			return 0
		}
		m := strings.ToUpper(string(b[i : i+2]))
		if !Known(m) {
			return 0
		}
		n++
		i += 2
		switch m {
		case "LB", "WD", "BL":
			k := bytes.IndexByte(b[i:], term)
			if k < 0 {
				return n
			}
			i += k + 1
		case "DT":
			if i < len(b) && b[i] != ';' {
				term = b[i]
				i++
			} else {
				term = 0x03
			}
			if i = skipParams(b, i); i < 0 {
				return 0
			}
		case "SM":
			if i < len(b) && b[i] != ';' {
				i++
			}
		case "CO", "MG":
			for i < len(b) && (b[i] == ' ' || b[i] == '\t') {
				i++
			}
			if i < len(b) && b[i] == '"' {
				k := bytes.IndexByte(b[i+1:], '"')
				if k < 0 {
					return n
				}
				i += k + 2
			} else if k := bytes.IndexByte(b[i:], ';'); k >= 0 {
				i += k
			} else {
				return n
			}
		case "PE":
			k := bytes.IndexByte(b[i:], ';')
			if k < 0 {
				return n
			}
			i += k + 1
			numeric = true
		case "BP":
			// kind, value pairs; the picture name is a quoted string
			for i < len(b) && b[i] != ';' && !isLetter(b[i]) && b[i] != 0x1b {
				if b[i] == '"' {
					k := bytes.IndexByte(b[i+1:], '"')
					if k < 0 {
						return n
					}
					i += k + 2
					continue
				}
				if !isParam(b[i]) {
					return 0
				}
				i++
			}
		default:
			j := skipParams(b, i)
			if j < 0 {
				return 0
			}
			if j > i {
				numeric = true
			}
			i = j
		}
	}
	if !numeric && n < 6 {
		return 0
	}
	return n
}

// skipParams skips numeric parameters and separators up to the next
// instruction; -1 when a character cannot be part of them.
func skipParams(b []byte, i int) int {
	for i < len(b) {
		c := b[i]
		switch {
		case isParam(c):
			i++
		case c == ';' || isLetter(c) || c == 0x1b || c == 0x03:
			return i
		default:
			return -1
		}
	}
	return i
}

func isParam(c byte) bool {
	return c >= '0' && c <= '9' || c == '.' || c == ',' || c == '+' || c == '-' || c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func isLetter(c byte) bool { return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' }

// Known reports whether m (upper case) is an HP-GL or HP-GL/2 mnemonic.
func Known(m string) bool {
	return strings.Contains(mnemonics, " "+m+" ")
}

// mnemonics are the instructions of HP-GL/2 with its extensions, and those
// of HP-GL that HP-GL/2 dropped.
const mnemonics = " AA AC AD AR AT BP BR BZ CF CI CO CP CR CT DC DF DI DL DP DR DT DV EA EC EP ER ES EW FI FN FP FR FT IN IP IR IW LA LB LM LO LT MC MG MT NP NR OD OE OH OI OP OS PA PC PD PE PG PM PP PR PS PU PW QL RA RF RO RP RR RT SA SB SC SD SI SL SM SP SR SS ST SV TD TR UL VS WG WU" +
	" AF AH AP AS BF BL CA CC CM CS CV DS EL GC GM GP IM IV KY LI OA OC OF OG OK OL OO OT OW PB PT SG TL UC VA VN WD XT YT "
