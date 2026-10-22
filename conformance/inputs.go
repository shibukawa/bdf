package conformance

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/shibukawa/bdf"
)

// The synthetic inputs: the malformed and edge cases of this package's own
// tests (malformed_test.go, bdf_test.go), written as bytes.

// w is a little-endian byte writer.
type w struct{ b []byte }

func (x *w) u8(v byte)        { x.b = append(x.b, v) }
func (x *w) u16(v uint16)     { x.b = binary.LittleEndian.AppendUint16(x.b, v) }
func (x *w) u32(v uint32)     { x.b = binary.LittleEndian.AppendUint32(x.b, v) }
func (x *w) u64(v uint64)     { x.b = binary.LittleEndian.AppendUint64(x.b, v) }
func (x *w) f32(v float32)    { x.u32(math.Float32bits(v)) }
func (x *w) varuint(v uint64) { x.b = binary.AppendUvarint(x.b, v) }
func (x *w) bytes(p []byte)   { x.b = append(x.b, p...) }
func (x *w) str(s string)     { x.varuint(uint64(len(s))); x.bytes([]byte(s)) }
func (x *w) f32s(vs ...float32) {
	for _, v := range vs {
		x.f32(v)
	}
}

// header writes a single-file header.
func header(flags uint16, moff, mlen uint64, menc byte) []byte {
	var h w
	h.bytes(bdf.Magic[:])
	h.u16(bdf.FormatVersion)
	h.u16(flags)
	h.u64(moff)
	h.u64(mlen)
	h.u8(menc)
	h.bytes(make([]byte, 7))
	return h.b
}

// single writes a document: the header, the manifest as it is, the parts.
func single(manifest string, parts []byte) []byte {
	b := append(header(0, bdf.HeaderSize, uint64(len(manifest)), 0), manifest...)
	return append(b, parts...)
}

// deflate compresses data as deflate-raw.
func deflate(data []byte, level int) []byte {
	var buf bytes.Buffer
	zw, err := flate.NewWriter(&buf, level)
	if err != nil {
		panic(err)
	}
	zw.Write(data)
	zw.Close()
	return buf.Bytes()
}

// object encodes an object part with the given tables and op stream; the
// tables are the encoded bytes after their counts.
func object(strings []string, paths, paints [][]byte, nObjects int, ops []byte) []byte {
	var x w
	x.bytes([]byte("BOBJ"))
	x.u16(bdf.OpsetVersion)
	x.u16(0)
	x.f32s(0, 0, 10, 10)
	x.varuint(uint64(len(strings)))
	for _, s := range strings {
		x.str(s)
	}
	x.varuint(uint64(len(paths)))
	for _, p := range paths {
		x.bytes(p)
	}
	x.varuint(uint64(len(paints)))
	for _, p := range paints {
		x.bytes(p)
	}
	x.varuint(0) // fonts
	x.varuint(0) // images
	x.varuint(uint64(nObjects))
	for i := range nObjects {
		h := bdf.Hash{byte(i + 1)}
		x.bytes(h[:])
	}
	x.varuint(uint64(len(ops)))
	x.bytes(ops)
	return x.b
}

// huge is a count no data can hold.
var huge = binary.AppendUvarint(nil, math.MaxUint64)

// sample is the object of TestRoundTripObject: every signature class.
func sample() []byte {
	o := bdf.NewObject()
	p := o.AddPath((&bdf.Path{}).MoveTo(0, 0).LineTo(10, 0).LineTo(10, 10).Close())
	f := o.AddFont(bdf.SystemFont("sans-serif", 400, bdf.StyleNormal))
	o.FillColor(bdf.RGB(1, 2, 3)).FillPath(p, bdf.NonZero)
	o.Dash([]float32{1, 2}, 0.5)
	o.Font(f, 12).FillText("hi", 1, 2, 3)
	o.FillPathRun(bdf.EvenOdd, []bdf.Glyph{{Path: p, X: 20, Y: 20}, {Path: p, X: 30, Y: 20}})
	o.MaskBegin(bdf.MaskAlpha, 0, nil).MaskEnd()
	o.Ext([]byte{1, 2, 3})
	return o.Encode()
}

const minimal = `{"bdf":1,"parts":[]}`

func cat(parts ...[]byte) []byte { return bytes.Join(parts, nil) }

func synthetic() []input {
	var in []input
	add := func(id, kind, note string, data []byte, want string) *input {
		in = append(in, input{id: id, kind: kind, note: note, data: data, want: want})
		return &in[len(in)-1]
	}

	// ---- the single-file header -------------------------------------------
	add("header/minimal", "single", "the smallest document: a header and an empty manifest", single(minimal, nil), OK)
	bad := single(minimal, nil)
	bad[0] = 'B'
	add("header/bad-magic", "single", "", bad, Err)
	bad = single(minimal, nil)
	bad[4] = 2
	add("header/version-2", "single", "a format version this reader does not know", bad, Err)
	bad = single(minimal, nil)
	bad[6] = 2
	add("header/undefined-flag", "single", "a reserved flag bit set", bad, Err)
	bad = single(minimal, nil)
	bad[24] = 2
	add("header/manifest-encoding-2", "single", "", bad, Err)
	bad = single(minimal, nil)
	bad[31] = 1
	add("header/reserved-byte", "single", "a reserved header byte that is not zero", bad, Err)
	add("header/short", "single", "less than a header", single(minimal, nil)[:20], Err)
	for _, c := range []struct {
		name       string
		moff, mlen uint64
	}{
		{"negative-length", bdf.HeaderSize, 1 << 63},
		{"negative-offset", 1 << 63, 8},
		{"sum-overflows", math.MaxInt64, math.MaxInt64},
		{"offset-and-length-cancel", 1<<64 - 8, 16},
		{"past-the-end", bdf.HeaderSize, 1 << 20},
	} {
		add("header/manifest-range-"+c.name, "single", "the manifest's range lies outside the file",
			append(header(0, c.moff, c.mlen, 0), minimal...), Err)
	}
	packed := deflate([]byte(minimal), flate.BestCompression)
	add("header/manifest-deflated", "single", "a manifest stored as deflate-raw",
		append(header(0, bdf.HeaderSize, uint64(len(packed)), 1), packed...), OK)
	corrupt := packed[:len(packed)-2]
	add("header/manifest-deflated-corrupt", "single", "a deflated manifest that ends before its final block",
		append(header(0, bdf.HeaderSize, uint64(len(corrupt)), 1), corrupt...), Err)
	add("header/encrypted-flag-without-block", "single", "the encrypted flag set, no encryption in the manifest",
		append(header(bdf.FlagEncrypted, bdf.HeaderSize, uint64(len(minimal)), 0), minimal...), Err)
	outer := `{"bdf":1,"encryption":{"cipher":"A256GCM","keys":[{"type":"password","kdf":"PBKDF2-SHA256","iter":600000,"salt":"AAECAwQFBgcICQoLDA0ODw==","key":"QUJD"}],"manifest":{"part":"abababababababababababababababab","enc":"deflate-raw"}},"parts":[{"h":"abababababababababababababababab","t":"sealed","enc":"identity","len":3,"size":3,"off":0}]}`
	add("header/encryption-block-without-flag", "single", "an encryption block, the flag clear", single(outer, []byte("xyz")), Err)
	locked := add("header/encrypted-outer", "single", "an encrypted document's outer manifest: it opens locked, with its sealed parts",
		append(append(header(bdf.FlagEncrypted, bdf.HeaderSize, uint64(len(outer)), 0), outer...), "xyz"...), OK)
	locked.wantReads = map[string]string{"abababababababababababababababab": OK}

	// ---- the manifest's numbers (TestManifestNumbers) ----------------------
	layers := `"layers":[]`
	for name, views := range map[string]string{
		"null-view":         `[null]`,
		"null-page":         `[{"id":"v","kind":"fixed","pages":[null]}]`,
		"negative-page":     `[{"id":"v","kind":"fixed","pages":[{"w":-5,"h":80000000,` + layers + `}]}]`,
		"huge-page":         `[{"id":"v","kind":"fixed","pages":[{"w":100,"h":1e9,` + layers + `}]}]`,
		"page-body":         `[{"id":"v","kind":"flow","pages":[{"w":100,"h":100,"body":{"x":0,"y":-1e30,"w":10,"h":10},` + layers + `}]}]`,
		"gap":               `[{"id":"v","kind":"flow","continuous":{"gap":1e30}}]`,
		"tiny-tile":         `[{"id":"v","kind":"sheet","tile":0.000001}]`,
		"many-rows":         `[{"id":"v","kind":"sheet","rows":[[16777216,10],[1,10]]}]`,
		"rows-beyond-f32":   `[{"id":"v","kind":"sheet","rows":[[1e11,0]]}]`,
		"half-a-column":     `[{"id":"v","kind":"sheet","cols":[[1.5,10]]}]`,
		"negative-columns":  `[{"id":"v","kind":"sheet","cols":[[-1,10]]}]`,
		"thin-columns":      `[{"id":"v","kind":"sheet","cols":[[400000,0.0000001]]}]`,
		"negative-row-size": `[{"id":"v","kind":"sheet","rows":[[4,-10]]}]`,
		"negative-freeze":   `[{"id":"v","kind":"sheet","freeze":{"rows":-1}}]`,
	} {
		add("manifest/"+name, "single", "a number out of the range the manifest allows",
			single(`{"bdf":1,"views":`+views+`,"parts":[]}`, nil), Err)
	}
	add("manifest/writers", "single", "what writers write: pages, hidden rows, the tile of the default size",
		single(`{"bdf":1,"views":[{"id":"p","kind":"flow","continuous":{"gap":12},"pages":[{"w":595.28,"h":841.89,"body":{"x":72,"y":72,"w":451,"h":697},`+layers+`}]},`+
			`{"id":"s","kind":"sheet","tile":2048,"cols":[[3,64],[16381,48]],"rows":[[1,0],[1048575,15]],"freeze":{"rows":1}},{"id":"t","kind":"sheet"}],"parts":[]}`, nil), OK)

	// ---- the manifest's JSON ------------------------------------------------
	for name, mj := range map[string]string{
		"trailing-comma":       `{"bdf":1,"parts":[],}`,
		"data-after-value":     `{"bdf":1,"parts":[]} x`,
		"empty":                ``,
		"array":                `[]`,
		"version-2":            `{"bdf":2,"parts":[]}`,
		"version-not-integer":  `{"bdf":1.0,"parts":[]}`,
		"version-string":       `{"bdf":"1","parts":[]}`,
		"opset-2":              `{"bdf":1,"opset":2}`,
		"views-object":         `{"bdf":1,"views":{}}`,
		"view-number":          `{"bdf":1,"views":[5]}`,
		"page-size-string":     `{"bdf":1,"views":[{"pages":[{"w":"5"}]}]}`,
		"page-size-beyond-f32": `{"bdf":1,"views":[{"pages":[{"w":1e39}]}]}`,
		"part-name-short":      `{"bdf":1,"parts":[{"h":"xyz"}]}`,
		"part-name-number":     `{"bdf":1,"parts":[{"h":5}]}`,
		"part-length-float":    `{"bdf":1,"parts":[{"h":"00000000000000000000000000000000","len":1e3}]}`,
		"dc-title-number":      `{"bdf":1,"meta":{"dc":{"title":5}}}`,
		"key-salt-not-base64":  `{"bdf":1,"encryption":{"keys":[{"salt":"not base64!"}]}}`,
	} {
		add("json/"+name, "single", "JSON a manifest reader refuses", single(mj, nil), Err)
	}
	for name, mj := range map[string]string{
		"duplicate-member-last-wins": `{"bdf":2,"bdf":1,"parts":[]}`,
		"unknown-members":            `{"bdf":1,"future":{"deep":[[[[1]]]]},"views":null,"parts":[]}`,
		"opset-1":                    `{"bdf":1,"opset":1}`,
		"part-name-uppercase":        `{"bdf":1,"parts":[{"h":"0123456789ABCDEF0123456789abcdef"}]}`,
		"short-run":                  `{"bdf":1,"views":[{"cols":[[26]],"rows":[[1,2,3]],"pages":[{"layers":[null]}]}]}`,
		"null-part":                  `{"bdf":1,"parts":[null]}`,
		"whitespace":                 " \t\r\n{\"bdf\":1}\n",
	} {
		add("json/"+name, "single", "JSON a manifest reader accepts", single(mj, nil), OK)
	}

	// ---- parts --------------------------------------------------------------
	data := []byte("0123456789")
	h := bdf.HashOf(data)
	for name, e := range map[string]string{
		"negative-length":  `"len":-1,"size":10`,
		"negative-offset":  `"len":10,"size":10,"off":-4`,
		"offset-overflows": `"len":10,"size":10,"off":9223372036854775807`,
		"past-the-end":     `"len":11,"size":11`,
	} {
		c := add("part/range-"+name, "single", "a part entry that reaches outside the parts region: the document opens, the part does not read",
			single(fmt.Sprintf(`{"bdf":1,"parts":[{"h":"%s","t":"img","enc":"identity",%s}]}`, h, e), data), OK)
		c.wantReads = map[string]string{h.String(): Err}
	}
	c := add("part/identity", "single", "a part stored as it is", single(fmt.Sprintf(`{"bdf":1,"parts":[{"h":"%s","t":"img","enc":"identity","len":10,"size":10,"off":0}]}`, h), data), OK)
	c.wantReads = map[string]string{h.String(): OK}
	// the block of signatures that follows the parts of a signed document
	// (spec §3.7): a reader that does not check signatures reads the
	// document as it is
	sig := []byte("0\x80 a signature, not read here")
	block := &w{}
	block.u32(uint32(len(sig)))
	block.bytes(sig)
	n := len(block.b)
	block.bytes([]byte("bdfs"))
	block.u16(1)
	block.u16(1)
	block.u64(uint64(n))
	c = add("part/signatures-after-parts", "single", "a signed document: its block of signatures follows the parts",
		cat(single(fmt.Sprintf(`{"bdf":1,"parts":[{"h":"%s","t":"img","enc":"identity","len":10,"size":10,"off":0}]}`, h), data), block.b), OK)
	c.wantReads = map[string]string{h.String(): OK}
	plain := make([]byte, 1<<20)
	zeros := deflate(plain, flate.BestCompression)
	hz := bdf.HashOf(plain)
	c = add("part/inflates-to-more-than-stated", "single", "a megabyte of zeros, stated as a kilobyte: the part does not read",
		single(fmt.Sprintf(`{"bdf":1,"parts":[{"h":"%s","t":"img","enc":"deflate-raw","len":%d,"size":1024,"off":0}]}`, hz, len(zeros)), zeros), OK)
	c.wantReads = map[string]string{hz.String(): Err}
	c = add("part/inflates-as-stated", "single", "a megabyte of zeros, stated as such",
		single(fmt.Sprintf(`{"bdf":1,"parts":[{"h":"%s","t":"img","enc":"deflate-raw","len":%d,"size":%d,"off":0}]}`, hz, len(zeros), len(plain)), zeros), OK)
	c.wantReads = map[string]string{hz.String(): OK}
	c = add("part/unknown-encoding", "single", "", single(fmt.Sprintf(`{"bdf":1,"parts":[{"h":"%s","t":"img","enc":"brotli","len":10,"size":10,"off":0}]}`, h), data), OK)
	c.wantReads = map[string]string{h.String(): Err}
	c = add("part/negative-size", "single", "a deflated part with a negative decoded size", single(fmt.Sprintf(`{"bdf":1,"parts":[{"h":"%s","t":"img","enc":"deflate-raw","len":10,"size":-1,"off":0}]}`, h), data), OK)
	c.wantReads = map[string]string{h.String(): Err}
	c = add("part/unknown", "single", "the manifest names no parts; nothing reads", single(minimal, data), OK)
	// parts that no view refers to: a reader opens the document whatever
	// their types (the print image of a form of text, a type to come)
	image := []byte(`{"version":1,"cpi":10,"lpi":6,"columns":80,"lines":66,"margin":{"left":0,"top":0},"pages":[{"lines":[{"text":"No. 1"}]}]}`)
	hi := bdf.HashOf(image)
	future := []byte("a part of a type to come")
	hf := bdf.HashOf(future)
	c = add("part/unreferenced-types", "single", "a page, the print image of a form of text (t \"text\") and a part of a type this reader does not know; neither stops the document from opening",
		single(fmt.Sprintf(`{"bdf":1,"views":[{"id":"pages","kind":"fixed","pages":[{"w":576,"h":792,"layers":[]}]}],"parts":[{"h":"%s","t":"text","enc":"identity","len":%d,"size":%d,"off":0},{"h":"%s","t":"x-future","enc":"identity","len":%d,"size":%d,"off":%d}]}`,
			hi, len(image), len(image), hf, len(future), len(future), len(image)), cat(image, future)), OK)
	c.wantReads = map[string]string{hi.String(): OK, hf.String(): OK}

	big := bytes.Repeat([]byte{' '}, bdf.MaxManifestSize+1)
	copy(big, minimal)
	packedBig := deflate(big, flate.BestSpeed)
	add("manifest/too-large", "single", "a manifest that inflates to more than 256 MiB",
		append(header(0, bdf.HeaderSize, uint64(len(packedBig)), 1), packedBig...), Err)

	// ---- a document a writer makes ------------------------------------------
	d := bdf.NewDocument()
	d.Meta.Source = "fixture"
	v := d.NewView("pages", bdf.ViewFixed, "Pages")
	for i := range 24 {
		o := bdf.NewObject()
		o.FillColor(bdf.RGB(byte(i*10), 0, 0)).FillRect(10, 10, 80, 80)
		o.Font(o.AddFont(bdf.SystemFont("sans-serif", 400, bdf.StyleNormal)), 12).FillText(fmt.Sprintf("page %d", i+1), 10, 50, 40)
		oh, _ := d.AddObject(o)
		v.AddPage(100, 100, bdf.Layer{Role: bdf.RoleBody, Obj: oh})
	}
	var buf bytes.Buffer
	if err := d.WriteSingle(&buf); err != nil {
		panic(err)
	}
	add("document/pages", "single", "24 pages written by this implementation: the manifest is deflated, the objects shared", buf.Bytes(), OK)

	// ---- objects ------------------------------------------------------------
	s := sample()
	c = add("object/sample", "object", "every operand signature", s, OK)
	c.wantWalk = OK
	for n := 0; n < len(s); n++ {
		add(fmt.Sprintf("object/truncated-%03d", n), "object", "a prefix of object/sample", s[:n], Err)
	}
	unknown := append([]byte(nil), s...)
	unknown[len(unknown)-8] = 0x7e // the EXT opcode
	c = add("object/unknown-opcode", "object", "an opcode that is not in opset 1 ends the walk", unknown, OK)
	c.wantWalk = Err
	pad := make([]byte, 64)
	for name, data := range map[string][]byte{
		"path-verbs":  object(nil, [][]byte{cat([]byte{0}, huge, pad)}, nil, 0, nil),
		"paint-stops": object(nil, nil, [][]byte{cat([]byte{bdf.PaintLinear}, make([]byte, 16), huge, pad)}, 0, nil),
		"strings":     cat([]byte("BOBJ\x01\x00\x00\x00"), make([]byte, 16), huge, pad),
	} {
		add("object/huge-count-"+name, "object", "a count no data can hold", data, Err)
	}
	for name, data := range map[string][]byte{
		"dash":       object(nil, nil, nil, 0, cat([]byte{bdf.OpDash}, huge, pad)),
		"glyph-run":  object(nil, nil, nil, 0, cat([]byte{bdf.OpFillPathRun, bdf.NonZero}, huge, pad)),
		"mask-bytes": object(nil, nil, nil, 0, cat([]byte{bdf.OpMaskBegin, 0, 0, 0, 0, 0}, huge, pad)),
	} {
		c = add("object/huge-count-"+name, "object", "a count no data can hold, in an instruction", data, OK)
		c.wantWalk = Err
	}
	c = add("object/bad-string-ref", "object", "a FILL_TEXT past the string table",
		object([]string{"text"}, nil, nil, 0, cat([]byte{bdf.OpFillText}, binary.AppendUvarint(nil, 5), make([]byte, 12))), OK)
	c.wantWalk = Err
	add("object/unknown-verb", "object", "", object(nil, [][]byte{{0, 1, 9}}, nil, 0, nil), Err)
	add("object/unknown-paint-kind", "object", "", object(nil, nil, [][]byte{{4}}, 0, nil), Err)
	add("object/opset-2", "object", "", cat([]byte("BOBJ\x02\x00\x00\x00"), make([]byte, 16), make([]byte, 7)), Err)
	add("object/flags", "object", "reserved object flags set", cat([]byte("BOBJ\x01\x00\x01\x00"), make([]byte, 16), make([]byte, 7)), Err)
	c = add("object/empty", "object", "empty tables and no instructions", object(nil, nil, nil, 0, nil), OK)
	c.wantWalk = OK
	c = add("object/trailing-bytes", "object", "bytes after the op stream are not read", append(object(nil, nil, nil, 0, nil), 0xab), OK)
	c.wantWalk = OK
	c = add("object/short-operands", "object", "", object(nil, nil, nil, 0, []byte{bdf.OpFillRect, 0, 0, 0, 0}), OK)
	c.wantWalk = Err

	// ---- path collections ---------------------------------------------------
	coll := bdf.EncodePathCollection([]*bdf.Path{(&bdf.Path{}).Rect(1, 2, 3, 4), (&bdf.Path{}).MoveTo(0, 0).QuadTo(1, 1, 2, 0)})
	add("paths/two", "paths", "", coll, OK)
	add("paths/truncated", "paths", "", coll[:len(coll)-1], Err)
	add("paths/huge-count", "paths", "", cat(huge, pad), Err)
	add("paths/unknown-verb", "paths", "", []byte{1, 1, 9}, Err)
	return in
}
