package bdf

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// IndexRun is one entry of a text index part.
type IndexRun struct {
	A, B    uint32 // page, layer (fixed/flow) or tile x, y (sheet)
	Ordinal uint32
	Sep     byte
	Text    string
}

// EncodeTextIndex serializes index runs as a text index part (type "idx").
func EncodeTextIndex(runs []IndexRun) []byte {
	var w buf
	w.bytes([]byte("BTXT"))
	w.u16(1)
	w.varuint(uint64(len(runs)))
	for _, r := range runs {
		w.varuint(uint64(r.A))
		w.varuint(uint64(r.B))
		w.varuint(uint64(r.Ordinal))
		w.u8(r.Sep)
		w.str(r.Text)
	}
	return w.b
}

// DecodeTextIndex parses a text index part.
func DecodeTextIndex(data []byte) ([]IndexRun, error) {
	r := &reader{b: data}
	if string(r.bytes(4)) != "BTXT" {
		return nil, &FormatError{Msg: "not a text index part"}
	}
	if v := r.u16(); v > 1 {
		return nil, &FormatError{Msg: fmt.Sprintf("unsupported text index version %d", v)}
	}
	n := r.varuint()
	if n > uint64(len(data)) {
		return nil, &FormatError{Msg: "bad run count"}
	}
	out := make([]IndexRun, 0, n)
	for i := uint64(0); i < n && r.err == nil; i++ {
		out = append(out, IndexRun{A: uint32(r.varuint()), B: uint32(r.varuint()), Ordinal: uint32(r.varuint()), Sep: r.u8(), Text: r.str()})
	}
	return out, r.err
}

// BuildTextIndex extracts the text of a view from the document's objects,
// stores it as a text index part and records it in v.TextIndex.
func (d *Document) BuildTextIndex(v *View) (Hash, error) {
	decoded := map[Hash]*ObjectPart{}
	resolve := func(h Hash) *ObjectPart {
		if o, ok := decoded[h]; ok {
			return o
		}
		p := d.parts[h]
		if p == nil {
			return nil
		}
		o, err := DecodeObject(p.Data)
		if err != nil {
			return nil
		}
		decoded[h] = o
		return o
	}
	var runs []IndexRun
	add := func(a, b uint32, h Hash, first bool) error {
		o := resolve(h)
		if o == nil {
			return fmt.Errorf("bdf: object %s missing or invalid", h)
		}
		trs, err := ExtractText(o, resolve)
		if err != nil {
			return err
		}
		for i, tr := range trs {
			sep := tr.Sep
			if i == 0 {
				sep = SepBreak
			}
			runs = append(runs, IndexRun{A: a, B: b, Ordinal: uint32(tr.Ordinal), Sep: sep, Text: tr.Text})
		}
		return nil
	}
	switch v.Kind {
	case ViewSheet:
		type key struct{ x, y int }
		keys := make([]key, 0, len(v.Tiles))
		for k := range v.Tiles {
			var kk key
			if _, err := fmt.Sscanf(k, "%d,%d", &kk.x, &kk.y); err != nil {
				return Hash{}, fmt.Errorf("bdf: bad tile key %q", k)
			}
			keys = append(keys, kk)
		}
		sort.Slice(keys, func(i, j int) bool {
			if keys[i].y != keys[j].y {
				return keys[i].y < keys[j].y
			}
			return keys[i].x < keys[j].x
		})
		for _, k := range keys {
			h, err := ParseHash(v.Tiles[strconv.Itoa(k.x)+","+strconv.Itoa(k.y)])
			if err != nil {
				return Hash{}, err
			}
			if err := add(uint32(k.x), uint32(k.y), h, true); err != nil {
				return Hash{}, err
			}
		}
	default:
		for pi, p := range v.Pages {
			for li, l := range p.Layers {
				if err := add(uint32(pi), uint32(li), l.Obj, true); err != nil {
					return Hash{}, err
				}
			}
		}
	}
	h := d.AddPart(PartIndex, EncodeTextIndex(runs))
	v.TextIndex = h.String()
	return h, nil
}

// PlainText joins index runs into a readable string (for tools and tests).
func PlainText(runs []IndexRun) string {
	var sb strings.Builder
	for i, r := range runs {
		if i > 0 {
			switch r.Sep {
			case SepSpace:
				sb.WriteByte(' ')
			case SepBreak:
				sb.WriteByte('\n')
			}
		}
		sb.WriteString(r.Text)
	}
	return sb.String()
}
