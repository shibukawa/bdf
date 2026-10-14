package pdfbdf

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

// PDF/A is the standard of files that are kept (ISO 19005): a file that
// shows the same in any reader, years on, because all it needs is in it
// and its colors are those of a known space. A file of this writer is
// close to one already: its fonts are embedded, and it has no scripts, no
// files beside it and no encryption. As a PDF/A file it also says the
// space of its colors (an output intent with the profile of sRGB, the
// space of the colors of a document), names the standard in its metadata,
// and leaves out what the standard has no place for (images that readers
// are asked to smooth, glyphs that stand for characters a font lacks).
//
// The parts written are 2 and 3 (ISO 19005-2 and -3, of PDF 1.7, with
// transparency), which differ in the files a document may carry, and this
// writer carries none. Part 1 has no transparency and is not written. The
// levels are b (it shows the same), u (and its text is Unicode: searched
// and copied) and a (and it is tagged: it has its structure).

// pdfaLevel is the PDF/A that a file is written as.
type pdfaLevel struct {
	part int
	conf string // "A", "B" or "U"
}

// parsePDFA reads Options.PDFA.
func parsePDFA(s string) (*pdfaLevel, error) {
	t := strings.ToLower(strings.TrimSpace(s))
	t = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(t, "pdf/a"), "pdfa"), "-")
	if t == "" {
		if strings.TrimSpace(s) == "" {
			return nil, nil
		}
		t = "2b"
	}
	if len(t) == 2 && (t[0] == '2' || t[0] == '3') && strings.ContainsRune("abu", rune(t[1])) {
		return &pdfaLevel{part: int(t[0] - '0'), conf: strings.ToUpper(t[1:])}, nil
	}
	return nil, fmt.Errorf("pdfbdf: %q is not a PDF/A that is written (2b, 2u, 2a, 3b, 3u or 3a)", s)
}

// outputIntent writes the output intent of a PDF/A file: the colors of
// the file, which are those of the document (sRGB), are those of the
// profile.
func (ex *exporter) outputIntent() string {
	profile := ex.pw.addStream(" /N 3", srgbProfile(), true)
	const name = "sRGB IEC61966-2.1"
	return fmt.Sprintf(" /OutputIntents [<< /Type /OutputIntent /S /GTS_PDFA1 /OutputConditionIdentifier %s /Info %s /DestOutputProfile %s >>]", text(name), text(name), profile)
}

// srgbProfile makes an ICC profile of sRGB (IEC 61966-2-1): of version 2,
// of a display, with the primaries of sRGB adapted to the white of the
// profile connection space (D50) and the curve of sRGB as a table. It is
// made here so that the file needs no profile of someone else.
func srgbProfile() []byte {
	fixed := func(v float64) []byte {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(int32(math.Round(v*65536))))
		return b
	}
	xyz := func(x, y, z float64) []byte {
		b := append([]byte("XYZ \x00\x00\x00\x00"), fixed(x)...)
		return append(append(b, fixed(y)...), fixed(z)...)
	}
	pad := func(b []byte) []byte {
		for len(b)%4 != 0 {
			b = append(b, 0)
		}
		return b
	}
	const descText = "sRGB IEC61966-2.1"
	desc := append([]byte("desc\x00\x00\x00\x00"), 0, 0, 0, byte(len(descText)+1))
	desc = append(append(desc, descText...), 0)
	desc = append(desc, make([]byte, 4+4+2+1+67)...) // no Unicode nor ScriptCode description
	cprt := append([]byte("text\x00\x00\x00\x00"), "No copyright, use freely\x00"...)
	// the curve of sRGB: from a coded value to light
	const n = 1024
	curve := append([]byte("curv\x00\x00\x00\x00"), 0, 0, byte(n>>8), byte(n&0xff))
	for i := 0; i < n; i++ {
		v := float64(i) / (n - 1)
		if v <= 0.04045 {
			v /= 12.92
		} else {
			v = math.Pow((v+0.055)/1.055, 2.4)
		}
		u := uint16(math.Round(v * 65535))
		curve = append(curve, byte(u>>8), byte(u))
	}
	type tag struct {
		sig  string
		data []byte
	}
	tags := []tag{
		{"desc", pad(desc)}, {"cprt", pad(cprt)},
		{"wtpt", xyz(0.9642, 1, 0.8249)},
		{"rXYZ", xyz(0.4360, 0.2225, 0.0139)}, {"gXYZ", xyz(0.3851, 0.7169, 0.0971)}, {"bXYZ", xyz(0.1431, 0.0606, 0.7141)},
		{"rTRC", pad(curve)},
	}
	// (the three channels have the same curve: one table)
	shared := map[string]string{"gTRC": "rTRC", "bTRC": "rTRC"}
	count := len(tags) + len(shared)
	off := 128 + 4 + 12*count
	var table, body []byte
	u32 := func(v int) []byte {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(v))
		return b
	}
	at := map[string][2]int{}
	for _, t := range tags {
		at[t.sig] = [2]int{off + len(body), len(t.data)}
		table = append(append(append(table, t.sig...), u32(off+len(body))...), u32(len(t.data))...)
		body = append(body, t.data...)
	}
	for _, sig := range []string{"gTRC", "bTRC"} {
		p := at[shared[sig]]
		table = append(append(append(table, sig...), u32(p[0])...), u32(p[1])...)
	}
	head := make([]byte, 128)
	binary.BigEndian.PutUint32(head[0:], uint32(off+len(body)))
	binary.BigEndian.PutUint32(head[8:], 0x02100000) // version 2.1
	copy(head[12:], "mntrRGB XYZ ")
	// (a date that does not change: the file is the same every time)
	for i, v := range []uint16{2026, 1, 1, 0, 0, 0} {
		binary.BigEndian.PutUint16(head[24+2*i:], v)
	}
	copy(head[36:], "acsp")
	copy(head[68:], fixed(0.9642)) // the illuminant of the connection space: D50
	copy(head[72:], fixed(1))
	copy(head[76:], fixed(0.8249))
	out := append(head, u32(count)...)
	out = append(out, table...)
	return append(out, body...)
}
