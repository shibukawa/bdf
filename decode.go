package bdf

import (
	"fmt"
	"strings"
)

// PathEntry is a decoded path table entry.
type PathEntry struct {
	Inline *Path // nil when external
	Hash   Hash
	Index  uint32
}

// ObjectPart is a decoded Object part.
type ObjectPart struct {
	Opset   uint16
	BBox    Rect
	Strings []string
	Paths   []PathEntry
	Paints  []Paint
	Fonts   []Font
	Images  []Hash
	Objects []Hash
	Ops     []byte
}

// DecodeObject parses an Object part.
func DecodeObject(data []byte) (*ObjectPart, error) {
	r := &reader{b: data}
	if string(r.bytes(4)) != "BOBJ" {
		return nil, &FormatError{Msg: "not an object part"}
	}
	o := &ObjectPart{Opset: r.u16()}
	if o.Opset > OpsetVersion {
		return nil, &FormatError{Msg: fmt.Sprintf("unsupported opset %d", o.Opset)}
	}
	r.u16()
	o.BBox = Rect{r.f32(), r.f32(), r.f32(), r.f32()}
	count := func() int {
		n := r.varuint()
		if n > uint64(len(data)) {
			r.fail("bad count")
			return 0
		}
		return int(n)
	}
	for i, n := 0, count(); i < n && r.err == nil; i++ {
		o.Strings = append(o.Strings, r.str())
	}
	for i, n := 0, count(); i < n && r.err == nil; i++ {
		if r.u8() == 0 {
			o.Paths = append(o.Paths, PathEntry{Inline: decodePath(r)})
		} else {
			h := r.hash()
			o.Paths = append(o.Paths, PathEntry{Hash: h, Index: uint32(r.varuint())})
		}
	}
	for i, n := 0, count(); i < n && r.err == nil; i++ {
		o.Paints = append(o.Paints, decodePaint(r))
	}
	for i, n := 0, count(); i < n && r.err == nil; i++ {
		o.Fonts = append(o.Fonts, decodeFont(r))
	}
	for i, n := 0, count(); i < n && r.err == nil; i++ {
		o.Images = append(o.Images, r.hash())
	}
	for i, n := 0, count(); i < n && r.err == nil; i++ {
		o.Objects = append(o.Objects, r.hash())
	}
	n := count()
	o.Ops = r.bytes(n)
	if r.err != nil {
		return nil, r.err
	}
	return o, nil
}

// Instr is one decoded instruction. Operands are stored in reading order;
// f32 operands as float32, u8/u32/varuint as uint64, string refs as string.
type Instr struct {
	Op   byte
	Args []any
}

// Instructions decodes the op stream into a slice.
func (o *ObjectPart) Instructions() ([]Instr, error) {
	var out []Instr
	err := o.Walk(func(in Instr) { out = append(out, in) })
	return out, err
}

// Walk decodes the op stream, calling fn for each instruction.
func (o *ObjectPart) Walk(fn func(Instr)) error {
	r := &reader{b: o.Ops}
	for !r.eof() {
		in, err := o.readInstr(r)
		if err != nil {
			return err
		}
		fn(in)
	}
	return r.err
}

// readInstr decodes the instruction at the reader's position.
func (o *ObjectPart) readInstr(r *reader) (Instr, error) {
	code := r.u8()
	info, ok := opTable[code]
	if !ok {
		return Instr{}, &FormatError{Offset: r.pos - 1, Msg: fmt.Sprintf("unknown opcode 0x%02x", code)}
	}
	in := Instr{Op: code}
	switch info.Sig {
	case "D":
		n := int(r.varuint())
		if n > len(o.Ops) {
			r.fail("bad dash count")
			break
		}
		segs := make([]float32, n)
		for i := range segs {
			segs[i] = r.f32()
		}
		in.Args = append(in.Args, segs, r.f32())
	case "R":
		rule := uint64(r.u8())
		n := int(r.varuint())
		if n > len(o.Ops) {
			r.fail("bad run count")
			break
		}
		glyphs := make([]Glyph, n)
		for i := range glyphs {
			glyphs[i] = Glyph{PathRef(r.varuint()), r.f32(), r.f32()}
		}
		in.Args = append(in.Args, rule, glyphs)
	case "X":
		n := int(r.u32())
		in.Args = append(in.Args, r.bytes(n))
	default:
		for _, c := range info.Sig {
			switch c {
			case 'f':
				in.Args = append(in.Args, r.f32())
			case 'b':
				in.Args = append(in.Args, uint64(r.u8()))
			case 'c':
				in.Args = append(in.Args, uint64(r.u32()))
			case 'v':
				in.Args = append(in.Args, r.varuint())
			case 's':
				i := r.varuint()
				if i >= uint64(len(o.Strings)) {
					r.fail("bad string ref")
				} else {
					in.Args = append(in.Args, o.Strings[i])
				}
			}
		}
	}
	if r.err != nil {
		return Instr{}, r.err
	}
	return in, nil
}

// Disassemble renders an Object part as human-readable text.
func Disassemble(data []byte) (string, error) {
	o, err := DecodeObject(data)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "object opset=%d bbox=(%g %g %g %g)\n", o.Opset, o.BBox.X, o.BBox.Y, o.BBox.W, o.BBox.H)
	for i, s := range o.Strings {
		fmt.Fprintf(&sb, "  str[%d] = %q\n", i, s)
	}
	for i, p := range o.Paths {
		if p.Inline != nil {
			fmt.Fprintf(&sb, "  path[%d] = inline %d verbs\n", i, len(p.Inline.Verbs))
		} else {
			fmt.Fprintf(&sb, "  path[%d] = ext %s#%d\n", i, p.Hash, p.Index)
		}
	}
	for i, p := range o.Paints {
		fmt.Fprintf(&sb, "  paint[%d] = kind %d %v stops=%d\n", i, p.Kind, p.Coords, len(p.Stops))
	}
	for i, f := range o.Fonts {
		if f.Kind == FontEmbedded {
			fmt.Fprintf(&sb, "  font[%d] = embedded %s w%d s%d %q\n", i, f.Hash, f.Weight, f.Style, f.Family)
		} else {
			fmt.Fprintf(&sb, "  font[%d] = system %q w%d s%d\n", i, f.Family, f.Weight, f.Style)
		}
	}
	for i, h := range o.Images {
		fmt.Fprintf(&sb, "  image[%d] = %s\n", i, h)
	}
	for i, h := range o.Objects {
		fmt.Fprintf(&sb, "  object[%d] = %s\n", i, h)
	}
	err = o.Walk(func(in Instr) {
		sb.WriteString("  ")
		sb.WriteString(opTable[in.Op].Name)
		for _, a := range in.Args {
			switch v := a.(type) {
			case float32:
				fmt.Fprintf(&sb, " %g", v)
			case uint64:
				fmt.Fprintf(&sb, " %d", v)
			case string:
				fmt.Fprintf(&sb, " %q", v)
			case []float32:
				fmt.Fprintf(&sb, " %v", v)
			case []Glyph:
				fmt.Fprintf(&sb, " %v", v)
			case []byte:
				fmt.Fprintf(&sb, " %d bytes", len(v))
			default:
				fmt.Fprintf(&sb, " %v", v)
			}
		}
		sb.WriteByte('\n')
	})
	return sb.String(), err
}
