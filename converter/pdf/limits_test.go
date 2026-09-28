package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/filter"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
)

// Damaged and hostile files: what a file states (sizes, counts, references
// to itself) must end in an error or a warning, in the time and the memory
// of a file of its size.

// pagePDF is a one-page PDF with the entries of its resources, its content
// stream and more objects, numbered from 5.
func pagePDF(resources string, content [2]string, more ...any) []byte {
	objs := []any{
		/* 1 */ `<< /Type /Catalog /Pages 2 0 R >>`,
		/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
		/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << ` + resources + ` >> >>`,
		/* 4 */ content,
	}
	return buildPDF(append(objs, more...), "/Root 1 0 R")
}

func deflated(b []byte) string {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(b)
	w.Close()
	return buf.String()
}

// convertIn converts data, and fails when the conversion panics or is not
// done within the limit. set changes the converter before the pages are
// converted.
func convertIn(t *testing.T, limit time.Duration, data []byte, opts *Options, set func(c *converter)) (*Result, error) {
	t.Helper()
	type outcome struct {
		res   *Result
		err   error
		panic any
	}
	done := make(chan outcome, 1)
	go func() {
		var o outcome
		defer func() {
			o.panic = recover()
			done <- o
		}()
		var s *Stream
		if s, o.err = NewStream(bytes.NewReader(data), opts); o.err != nil {
			return
		}
		if set != nil {
			set(s.c)
		}
		o.res, o.err = s.Finish()
	}()
	select {
	case o := <-done:
		if o.panic != nil {
			t.Fatalf("panic: %v", o.panic)
		}
		return o.res, o.err
	case <-time.After(limit):
		t.Fatalf("not converted within %v", limit)
	}
	return nil, nil
}

// converted is convertIn for a file that converts.
func converted(t *testing.T, data []byte) *Result {
	t.Helper()
	res, err := convertIn(t, 5*time.Second, data, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func warned(res *Result, part string) bool {
	for _, w := range res.Warnings {
		if strings.Contains(w, part) {
			return true
		}
	}
	return false
}

// opened writes the document and opens it; it returns the body of the
// first page too.
func opened(t *testing.T, res *Result) (*bdf.Reader, bdf.Hash) {
	t.Helper()
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return r, r.Manifest.Views[0].Pages[0].Layers[0].Obj
}

// bodyOps counts the instructions of the body of the first page and of the
// objects it uses; the paths filled, however they are written, are counted
// as FILL_PATH.
func bodyOps(t *testing.T, res *Result) map[byte]int {
	t.Helper()
	r, body := opened(t, res)
	counts := map[byte]int{}
	var walk func(h bdf.Hash)
	walk = func(h bdf.Hash) {
		o, err := r.Object(h)
		if err != nil {
			t.Fatal(err)
		}
		err = o.Walk(func(in bdf.Instr) {
			switch in.Op {
			case bdf.OpFillPathAt, bdf.OpFillPathRun, bdf.OpFillRect:
				counts[bdf.OpFillPath]++
			case bdf.OpUse, bdf.OpUseAt:
				counts[in.Op]++
				walk(o.Objects[int(in.Args[0].(uint64))])
			default:
				counts[in.Op]++
			}
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	walk(body)
	return counts
}

func TestDamagedStreams(t *testing.T) {
	draw := "0 0 10 10 re f\n"
	z := deflated([]byte(draw))
	for _, c := range []struct {
		name, dict, data string
		fills            int    // FILL_PATH of the page or the form
		warning          string // "" for none
	}{
		{"run length, whole", "/Filter /RunLengthDecode", "\x0e" + draw + "\x80", 1, ""},
		{"run length, cut within its bytes", "/Filter /RunLengthDecode", "\x7f" + draw, 1, ""},
		{"run length, cut before the byte to repeat", "/Filter /RunLengthDecode", "\x0e" + draw + "\xf0", 1, ""},
		{"ASCII85 without data", "/Filter /ASCII85Decode", "", 0, ""},
		{"predictor, no column", "/Filter /FlateDecode /DecodeParms << /Predictor 12 /Columns 0 >>", z, 0, "predictor with 0 columns"},
		{"predictor, columns below zero", "/Filter /FlateDecode /DecodeParms << /Predictor 12 /Columns -5 >>", z, 0, "predictor with -5 columns"},
		{"TIFF predictor, no column", "/Filter /FlateDecode /DecodeParms << /Predictor 2 /Columns 0 >>", z, 0, "predictor with 0 columns"},
		{"predictor, colours below zero", "/Filter /FlateDecode /DecodeParms << /Predictor 11 /Columns 5 /Colors -1 >>", z, 0, "predictor with -1 colours"},
		{"predictor, rows longer than the data", "/Filter /FlateDecode /DecodeParms << /Predictor 12 /Columns 16777216 >>", z, 0, "predictor with rows of 16777216 bytes"},
	} {
		for _, form := range []bool{false, true} {
			name := c.name + ", page content"
			data := pagePDF("", [2]string{c.dict, c.data})
			if form {
				name = c.name + ", form"
				data = pagePDF("/XObject << /Fm 5 0 R >>", [2]string{"", "/Fm Do"},
					[2]string{"/Type /XObject /Subtype /Form /BBox [0 0 10 10] " + c.dict, c.data})
			}
			t.Run(name, func(t *testing.T) {
				res := converted(t, data)
				if c.warning == "" && len(res.Warnings) > 0 || c.warning != "" && !warned(res, c.warning) {
					t.Errorf("warnings: %q, want %q", res.Warnings, c.warning)
				}
				if fills := bodyOps(t, res)[bdf.OpFillPath]; fills != c.fills {
					t.Errorf("%d paths filled, want %d", fills, c.fills)
				}
			})
		}
	}
}

func TestDecodeRunLength(t *testing.T) {
	for _, c := range []struct {
		src  string
		max  int64
		want string
	}{
		{"\x02abc\xfed\x00e\x80f", 100, "abcddde"},
		{"\x02abc\xfed", 100, "abcddd"},
		{"\x02ab", 100, "ab"},
		{"\x81a\x81b", 100, "a" + strings.Repeat("a", 127)}, // a little more than max, not all
	} {
		if got := string(decodeRunLength([]byte(c.src), c.max)); got != c.want && !(len(c.want) == 128 && strings.HasPrefix(got, c.want) && len(got) < 100+128+1) {
			t.Errorf("%q: %q, want %q", c.src, got, c.want)
		}
	}
}

// TestDamagedFileIsAnError: what panics in pdfcpu is an error of the
// conversion.
func TestDamagedFileIsAnError(t *testing.T) {
	data := []byte("%PDF-1.7\n1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n" +
		"3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\ntrailer\n<< /Root 1 0 R /Size 4 >>\nstartxref\n0\n%%EOF\n")
	if _, err := convertIn(t, 5*time.Second, data, nil, nil); err == nil || !strings.Contains(err.Error(), "damaged") {
		t.Errorf("error: %v", err)
	}
	s, err := NewStream(bytes.NewReader(pagePDF("", [2]string{"", "0 0 10 10 re f"})), nil)
	if err != nil {
		t.Fatal(err)
	}
	s.c.pageBodies[0].attrs = nil // convertPage reads the resources of the page: a panic
	if _, err := s.Page(0); err == nil || !strings.Contains(err.Error(), "damaged") {
		t.Errorf("Page: error: %v", err)
	}
	if _, err := s.Finish(); err == nil || !strings.Contains(err.Error(), "damaged") {
		t.Errorf("Finish: error: %v", err)
	}
}

// fixturePDFs returns the PDFs of testdata and those the tests make.
func fixturePDFs(t *testing.T) map[string][]byte {
	t.Helper()
	files, err := filepath.Glob("testdata/*.pdf")
	if err != nil || len(files) == 0 {
		t.Fatalf("no fixtures: %v", err)
	}
	out := map[string][]byte{"tagged": taggedPDF(true), "softmask": softMaskPDF(), "cjk": cjkPDF("BT /F1 12 Tf <82a0> Tj ET")}
	for _, f := range files {
		if out[filepath.Base(f)], err = os.ReadFile(f); err != nil {
			t.Fatal(err)
		}
	}
	return out
}

// TestPagesAsPdfcpu: the pages, the attributes they inherit and their
// content are those pdfcpu gives.
func TestPagesAsPdfcpu(t *testing.T) {
	for name, data := range fixturePDFs(t) {
		ctx, _, err := readContext(bytes.NewReader(data), "")
		if err != nil {
			t.Fatal(name, err)
		}
		p := &pdf{ctx: ctx}
		leaves, _, err := p.pageLeaves()
		if err != nil {
			t.Fatal(name, err)
		}
		if err := ctx.EnsurePageCount(); err != nil || ctx.PageCount != len(leaves) || len(leaves) == 0 {
			t.Fatalf("%s: %d pages, pdfcpu has %d (%v)", name, len(leaves), ctx.PageCount, err)
		}
		for i, leaf := range leaves {
			dict, ref, attrs, err := ctx.PageDict(i+1, false)
			if err != nil || leaf.err != nil {
				t.Fatalf("%s page %d: %v, %v", name, i+1, err, leaf.err)
			}
			if !reflect.DeepEqual(dict, leaf.dict) || !reflect.DeepEqual(ref, leaf.ref) || !reflect.DeepEqual(attrs, leaf.attrs) {
				t.Errorf("%s page %d: %v %v %+v, pdfcpu has %v %v %+v", name, i+1, leaf.dict, leaf.ref, leaf.attrs, dict, ref, attrs)
			}
			want, err := ctx.PageContent(dict, i+1)
			if err != nil && err != model.ErrNoContent {
				t.Fatal(name, err)
			}
			got, err := p.pageContent(leaf.dict)
			if err != nil || !bytes.Equal(got, want) || len(got) == 0 {
				t.Errorf("%s page %d: content of %d bytes (%v), pdfcpu has %d", name, i+1, len(got), err, len(want))
			}
		}
	}
}

func TestPageContentStreams(t *testing.T) {
	// streams of an array follow one another without a separator
	data := buildPDF([]any{
		`<< /Type /Catalog /Pages 2 0 R >>`,
		`<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
		`<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents [4 0 R null 5 0 R 9 0 R] >>`,
		[2]string{"", "0 0 10 10 re"},
		[2]string{"/Filter /FlateDecode", deflated([]byte(" f"))},
	}, "/Root 1 0 R")
	ctx, _, err := readContext(bytes.NewReader(data), "")
	if err != nil {
		t.Fatal(err)
	}
	p := &pdf{ctx: ctx}
	leaves, _, _ := p.pageLeaves()
	if got, err := p.pageContent(leaves[0].dict); err != nil || string(got) != "0 0 10 10 re f" {
		t.Errorf("content %q, %v", got, err)
	}
	p.maxStream = 13
	if _, err := p.pageContent(leaves[0].dict); err == nil {
		t.Error("content of 14 bytes within a limit of 13")
	}
}

func TestPageTree(t *testing.T) {
	page := `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 100 200] >>`
	chain := func(nodes int) []any {
		objs := []any{`<< /Type /Catalog /Pages 2 0 R >>`}
		for i := 0; i < nodes; i++ {
			objs = append(objs, fmt.Sprintf(`<< /Type /Pages /Kids [%d 0 R] /Count 1 >>`, i+3))
		}
		return append(objs, page)
	}
	for _, c := range []struct {
		name  string
		objs  []any
		pages int
		w     float32
	}{
		{"a node among its own kids", []any{`<< /Type /Catalog /Pages 2 0 R >>`, `<< /Type /Pages /Kids [2 0 R 3 0 R] /Count 1 >>`, page}, 1, 100},
		{"a node among the kids of its kid", []any{`<< /Type /Catalog /Pages 2 0 R >>`, `<< /Type /Pages /Kids [4 0 R 3 0 R] /Count 1 >>`, page,
			`<< /Type /Pages /Kids [2 0 R 3 0 R] /Count 1 >>`}, 2, 100},
		{"/Count beyond the pages", []any{`<< /Type /Catalog /Pages 2 0 R >>`, `<< /Type /Pages /Kids [3 0 R] /Count 3000000 >>`, page}, 1, 100},
		{"no /Count", []any{`<< /Type /Catalog /Pages 2 0 R >>`, `<< /Type /Pages /Kids [3 0 R] >>`, page}, 1, 100},
		{"a box of two numbers", []any{`<< /Type /Catalog /Pages 2 0 R >>`, `<< /Type /Pages /Kids [3 0 R 4 0 R] /Count 2 /MediaBox [0 0] >>`, page, page}, -1, 0},
		{"inherited", []any{`<< /Type /Catalog /Pages 2 0 R >>`, `<< /Type /Pages /Kids [3 0 R 4 0 R] /Count 2 /MediaBox [0 0 300 400] /Rotate 90 >>`,
			`<< /Type /Page /Parent 2 0 R >>`, page}, 2, 400},
		{"a page listed over and over", []any{`<< /Type /Catalog /Pages 2 0 R >>`, `<< /Type /Pages /Kids [` + strings.Repeat("3 0 R ", 1000) + `] /Count 1000 >>`, page},
			1 + maxRepeatedPages, 100},
		{"nodes nested as deep as the limit", chain(maxPageTreeDepth + 1), 1, 100},
		{"nodes nested deeper than the limit", chain(maxPageTreeDepth + 2), 0, 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			res, err := convertIn(t, 5*time.Second, buildPDF(c.objs, "/Root 1 0 R"), nil, nil)
			if c.pages < 0 {
				if err == nil || !strings.Contains(err.Error(), "/MediaBox is not a rectangle") {
					t.Errorf("error: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if res.Pages != c.pages || c.pages > 0 && res.Doc.Views[0].Pages[0].W != c.w {
				t.Errorf("%d pages, want %d of width %v: %+v", res.Pages, c.pages, c.w, res.Doc.Views[0].Pages)
			}
		})
	}
}

// TestManyPages: the page tree is walked once, not from the root for
// every page (20 s for these pages before).
func TestManyPages(t *testing.T) {
	const n = 12000
	objs := []any{`<< /Type /Catalog /Pages 2 0 R >>`, ""}
	var kids strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&kids, "%d 0 R ", 3+i)
		objs = append(objs, `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>`)
	}
	objs[1] = fmt.Sprintf(`<< /Type /Pages /Kids [%s] /Count %d >>`, kids.String(), n)
	data := buildPDF(objs, "/Root 1 0 R")
	done := make(chan error, 1)
	go func() {
		s, err := NewStream(bytes.NewReader(data), &Options{NoTextIndex: true})
		if err == nil && s.Pages() != n {
			err = fmt.Errorf("%d pages", s.Pages())
		}
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("the pages of %d were not read within 5 s", n)
	}
}

func TestFunctionDepth(t *testing.T) {
	// a function among its own parts, and the same one many times
	for _, fn := range []string{
		`<< /FunctionType 3 /Domain [0 1] /Functions [5 0 R] /Bounds [] /Encode [0 1] >>`,
		`[5 0 R 5 0 R 5 0 R 5 0 R]`,
	} {
		res := converted(t, pagePDF("/ColorSpace << /CS0 [/Separation /X /DeviceRGB 5 0 R] >>", [2]string{"", "/CS0 cs 1 scn 0 0 10 10 re f"}, fn))
		if n := bodyOps(t, res)[bdf.OpFillPath]; n != 1 {
			t.Errorf("%s: %d paths filled", fn, n)
		}
	}

	// parts within parts, deeper than the limit by one, and more than the
	// limit of them
	objs := []any{`<< /Type /Catalog /Pages 2 0 R >>`, `<< /Type /Pages /Kids [] /Count 0 >>`}
	const first = 3
	for i := 0; i < maxFunctionDepth+2; i++ {
		objs = append(objs, fmt.Sprintf(`<< /FunctionType 3 /Domain [0 1] /Functions [%d 0 R] /Bounds [] /Encode [0 1] >>`, first+i+1))
	}
	objs = append(objs, `<< /FunctionType 2 /Domain [0 1] /C0 [0] /C1 [1] /N 1 >>`)
	exp := len(objs)
	objs = append(objs, "["+strings.Repeat(fmt.Sprintf("%d 0 R ", exp), maxFunctionParts+10)+"]")
	ctx, _, err := readContext(bytes.NewReader(buildPDF(objs, "/Root 1 0 R")), "")
	if err != nil {
		t.Fatal(err)
	}
	p := &pdf{ctx: ctx}
	ref := func(n int) types.Object { return *types.NewIndirectRef(n, 0) }
	depth := 0
	for f := p.loadFunction(ref(first)); f != nil && len(f.funcs) == 1; f = f.funcs[0] {
		depth++
	}
	if depth != maxFunctionDepth+1 { // the functions at depth 0 to maxFunctionDepth
		t.Errorf("functions %d deep, want %d", depth, maxFunctionDepth+1)
	}
	if f := p.loadFunction(ref(first + 2)); f == nil || len(f.eval(0.5)) != 1 {
		t.Errorf("the functions within the limit are not evaluated: %+v", f)
	}
	if f := p.loadFunction(ref(exp + 1)); f == nil || len(f.parts) != maxFunctionParts-1 {
		t.Errorf("function of %d parts, want %d", len(f.parts), maxFunctionParts-1)
	}
}

func TestColorSpaceDepth(t *testing.T) {
	// the colour space of the resources is one deep, the spaces within it deeper
	nested := func(depth int) string {
		cs := "/DeviceRGB"
		for i := 1; i < depth; i++ {
			cs = "[/Pattern " + cs + "]"
		}
		return cs
	}
	for _, c := range []struct {
		name, cs string
		warning  bool
	}{
		{"its own base", "5 0 R", true},
		{"deeper than the limit", nested(maxColorSpaceDepth + 1), true},
		{"as deep as the limit", nested(maxColorSpaceDepth), false},
	} {
		res := converted(t, pagePDF("/ColorSpace << /CS0 "+c.cs+" >>", [2]string{"", "/CS0 cs 0 0 10 10 re f"}, `[/Indexed 5 0 R 1 <000000ffffff>]`))
		if warned(res, "colour spaces nested too deeply") != c.warning {
			t.Errorf("%s: warnings %q", c.name, res.Warnings)
		}
	}
}

// nesting returns the depth of the arrays and dictionaries of an operand.
func nesting(o types.Object) int {
	depth := 0
	switch v := o.(type) {
	case types.Array:
		for _, e := range v {
			depth = max(depth, nesting(e))
		}
		return depth + 1
	case types.Dict:
		for _, e := range v {
			depth = max(depth, nesting(e))
		}
		return depth + 1
	}
	return 0
}

func TestLexerDepth(t *testing.T) {
	for _, c := range []struct {
		open, close string
		n, want     int
	}{
		{"[", "]", maxLexDepth, maxLexDepth},
		{"[", "]", maxLexDepth + 8, maxLexDepth},
		{"<</A ", ">>", maxLexDepth + 8, maxLexDepth},
		{"[<</A ", ">>]", maxLexDepth, maxLexDepth},
	} {
		l := &lexer{b: []byte(strings.Repeat(c.open, c.n) + "7 " + strings.Repeat(c.close, c.n) + " 8 Tj")}
		kind, obj, _ := l.next()
		if kind != tokOperand || nesting(obj) != c.want {
			t.Errorf("%d of %q: operand %d deep, want %d", c.n, c.open, nesting(obj), c.want)
		}
		// the closing brackets left over are skipped
		if kind, obj, _ := l.next(); kind != tokOperand || obj != types.Integer(8) {
			t.Errorf("%d of %q: then %v %v", c.n, c.open, kind, obj)
		}
		if kind, _, op := l.next(); kind != tokOperator || op != "Tj" {
			t.Errorf("%d of %q: then %v %q", c.n, c.open, kind, op)
		}
	}
	// stray delimiters are skipped one after the other
	l := &lexer{b: []byte(strings.Repeat(")]>}{", 1<<16) + "7")}
	if kind, obj, _ := l.next(); kind != tokOperand || obj != types.Integer(7) {
		t.Errorf("after stray delimiters: %v %v", kind, obj)
	}
	// the same in a content stream and in a CMap
	deep := strings.Repeat("[", 1000) + "(a)" + strings.Repeat("]", 1000)
	res := converted(t, pagePDF("/Font << /F1 5 0 R >>", [2]string{"", "BT /F1 12 Tf " + deep + " TJ ET"},
		`<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /ToUnicode 6 0 R >>`, [2]string{"", deep + " beginbfchar <61> <0062> endbfchar"}))
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
}

func TestStreamLimit(t *testing.T) {
	text := bytes.Repeat([]byte("0 0 1 1 re f\n"), 100) // 1300 bytes
	lzw, err := func() ([]byte, error) {
		fi, _ := filter.NewFilter(filter.LZW, nil)
		r, err := fi.Encode(bytes.NewReader(text))
		if err != nil {
			return nil, err
		}
		return readAll(r)
	}()
	if err != nil {
		t.Fatal(err)
	}
	var runs []byte
	for i := 0; i < len(text); i += 100 {
		runs = append(append(runs, 99), text[i:i+100]...)
	}
	z := []byte(deflated(text))
	ctx, _, err := readContext(bytes.NewReader(pagePDF("", [2]string{"", ""})), "")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name    string
		data    []byte
		filters []string
		parms   types.Dict
	}{
		{"flate", z, []string{filter.Flate}, nil},
		{"flate, predictor", []byte(deflated(append([]byte{0}, text...))), []string{filter.Flate}, types.Dict{"Predictor": types.Integer(12), "Columns": types.Integer(len(text))}},
		{"flate twice", []byte(deflated(z)), []string{filter.Flate, filter.Flate}, nil},
		{"LZW", lzw, []string{filter.LZW}, nil},
		{"run length", runs, []string{filter.RunLength}, nil},
		{"ASCII hex of flate", []byte(fmt.Sprintf("%x>", z)), []string{filter.ASCIIHex, filter.Flate}, nil},
		{"none", text, nil, nil},
	} {
		sd := &types.StreamDict{Raw: c.data}
		for _, f := range c.filters {
			sd.FilterPipeline = append(sd.FilterPipeline, types.PDFFilter{Name: f, DecodeParms: c.parms})
		}
		p := &pdf{ctx: ctx}
		if got, _, err := p.decodeStream(sd); err != nil || !bytes.Equal(got, text) {
			t.Errorf("%s: %d bytes, %v", c.name, len(got), err)
		}
		if len(c.filters) == 0 {
			continue
		}
		p.maxStream = int64(len(text))
		if got, _, err := p.decodeStream(sd); err != nil || !bytes.Equal(got, text) {
			t.Errorf("%s, limit of its length: %d bytes, %v", c.name, len(got), err)
		}
		p.maxStream = int64(len(text)) - 1
		if got, _, err := p.decodeStream(sd); err == nil || !strings.Contains(err.Error(), "decodes to more than 1299 bytes") {
			t.Errorf("%s, limit below its length: %d bytes, %v", c.name, len(got), err)
		}
		if got, _, err := p.decodeStreamUpTo(sd, 1000, true); err != nil || !bytes.Equal(got, text[:1000]) {
			t.Errorf("%s, cut: %d bytes, %v", c.name, len(got), err)
		}
	}

	// the content of a page and of a form; the rest of the page is converted
	const limit = 1 << 16
	for _, n := range []int{limit, limit + 1} {
		big := deflated(append(bytes.Repeat([]byte(" "), n-len("0 0 9 9 re f")), "0 0 9 9 re f"...))
		res, err := convertIn(t, 5*time.Second, pagePDF("/XObject << /Fm 5 0 R >>", [2]string{"", "/Fm Do 0 0 10 10 re f"},
			[2]string{"/Type /XObject /Subtype /Form /BBox [0 0 10 10] /Filter /FlateDecode", big}), nil, func(c *converter) { c.pdf.maxStream = limit })
		if err != nil {
			t.Fatal(err)
		}
		fills, uses := 2, 1 // the form and its path, and the path of the page
		if n > limit {
			fills, uses = 1, 0
		}
		if ops := bodyOps(t, res); warned(res, "form XObject: FlateDecode: the stream decodes to more than 65536 bytes") != (n > limit) ||
			ops[bdf.OpFillPath] != fills || ops[bdf.OpUse] != uses {
			t.Errorf("form of %d bytes: warnings %q, instructions %v", n, res.Warnings, ops)
		}
	}
}

// allocated returns the bytes that f allocates.
func allocated(f func()) uint64 {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	f()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

// TestImageData: the samples of an image are decoded as far as its rows.
func TestImageData(t *testing.T) {
	data := pagePDF("/XObject << /Im 5 0 R /Mk 6 0 R >>", [2]string{"", "/Im Do /Mk Do"},
		[2]string{"/Type /XObject /Subtype /Image /Width 2 /Height 2 /BitsPerComponent 8 /ColorSpace /DeviceGray /Filter /FlateDecode", deflated(make([]byte, 32<<20))},
		[2]string{"/Type /XObject /Subtype /Image /Width 2 /Height 2 /ImageMask true /Filter /FlateDecode", deflated(make([]byte, 32<<20))})
	var res *Result
	if n := allocated(func() { res = converted(t, data) }); n > 8<<20 {
		t.Errorf("%d bytes allocated for two images of four pixels", n)
	}
	if len(res.Warnings) != 0 || bodyOps(t, res)[bdf.OpImage] != 2 {
		t.Errorf("warnings %q, instructions %v", res.Warnings, bodyOps(t, res))
	}
}

// TestImageSize: the images that are left out, whatever their data.
func TestImageSize(t *testing.T) {
	for _, c := range []struct {
		dict    string
		samples int64 // the limit of the bytes of samples (0: the default)
		warning string
	}{
		{"/Width 2 /Height 2 /BitsPerComponent 8 /ColorSpace /DeviceGray", 0, ""},
		{"/Width 65536 /Height 1 /BitsPerComponent 8 /ColorSpace /DeviceGray", 0, "image of 65536 × 1 pixels is too large"},
		{"/Width 1 /Height 65536 /BitsPerComponent 8 /ColorSpace /DeviceGray", 0, "image of 1 × 65536 pixels is too large"},
		{"/Width 1099511627776 /Height 1 /BitsPerComponent 8 /ColorSpace /DeviceGray", 0, "pixels is too large"},
		{"/Width 65535 /Height 65535 /BitsPerComponent 8 /ColorSpace /DeviceGray", 0, "image of 65535 × 65535 pixels is too large (4294836225 bytes of samples)"},
		{"/Width 65535 /Height 65535 /BitsPerComponent 1 /ColorSpace /DeviceRGB", 0, "image of 65535 × 65535 pixels is too large (1610588160 bytes of samples)"},
		{"/Width 30 /Height 10 /BitsPerComponent 8 /ColorSpace /DeviceRGB", 900, ""},
		{"/Width 30 /Height 10 /BitsPerComponent 8 /ColorSpace /DeviceRGB", 899, "image of 30 × 10 pixels is too large (900 bytes of samples)"},
		{"/Width 30 /Height 10 /BitsPerComponent 1 /ColorSpace /DeviceGray", 40, ""},
		{"/Width 30 /Height 10 /BitsPerComponent 1 /ColorSpace /DeviceGray", 39, "image of 30 × 10 pixels is too large (40 bytes of samples)"},
		{"/Width 30 /Height 10 /ImageMask true", 39, "image of 30 × 10 pixels is too large (40 bytes of samples)"},
		{"/Width 30 /Height 10 /BitsPerComponent 1 /Filter /CCITTFaxDecode /DecodeParms << /K -1 /Columns 30 >>", 39, "image of 30 × 10 pixels is too large (40 bytes of samples)"},
		{"/Width 30 /Height 10 /BitsPerComponent 1 /Filter /JBIG2Decode", 39, "image of 30 × 10 pixels is too large (40 bytes of samples)"},
		{"/Width 2 /Height 2 /BitsPerComponent 0 /ColorSpace /DeviceGray", 0, "image with 0 bits per component"},
		{"/Width 2 /Height 2 /BitsPerComponent 4611686018427387904 /ColorSpace /DeviceGray", 0, "bits per component"},
		{"/Width 2 /Height 2 /BitsPerComponent 8 /ColorSpace [/DeviceN [" + strings.Repeat("/A ", 33) + "] /DeviceGray 6 0 R]", 0, "image with 33 colour components"},
		{"/Width 2 /Height 2 /BitsPerComponent 1 /Filter /CCITTFaxDecode /DecodeParms << /K -1 /Columns 65536 /Rows 2 >>", 0, "CCITTFaxDecode with 65536 columns"},
	} {
		var res *Result
		n := allocated(func() {
			var err error
			res, err = convertIn(t, 5*time.Second, pagePDF("/XObject << /Im 5 0 R >>", [2]string{"", "/Im Do"},
				[2]string{"/Type /XObject /Subtype /Image " + c.dict, "\xff\xff\xff\xff"},
				`<< /FunctionType 2 /Domain [0 1] /N 1 >>`), nil, func(conv *converter) { conv.pdf.maxSamples = c.samples })
			if err != nil {
				t.Fatal(err)
			}
		})
		if n > 8<<20 {
			t.Errorf("%s: %d bytes allocated", c.dict, n)
		}
		images := bodyOps(t, res)[bdf.OpImage]
		if c.warning == "" && (len(res.Warnings) != 0 || images != 1) || c.warning != "" && (!warned(res, c.warning) || images != 0) {
			t.Errorf("%s: warnings %q, %d images", c.dict, res.Warnings, images)
		}
	}
	// A JPEG is stored as it is, whatever size the dictionary states: the
	// size is that of the JPEG, and nothing is made of it here.
	var jpg bytes.Buffer
	if err := jpeg.Encode(&jpg, image.NewGray(image.Rect(0, 0, 16, 16)), nil); err != nil {
		t.Fatal(err)
	}
	stored := converted(t, pagePDF("/XObject << /Im 5 0 R >>", [2]string{"", "/Im Do"},
		[2]string{"/Type /XObject /Subtype /Image /Width 20000 /Height 20000 /BitsPerComponent 8 /ColorSpace /DeviceGray /Filter /DCTDecode", jpg.String()}))
	if len(stored.Warnings) != 0 || bodyOps(t, stored)[bdf.OpImage] != 1 {
		t.Errorf("JPEG: warnings %q", stored.Warnings)
	}
	// The rows of a CCITT image are those of the image at most: a bit of
	// Group 4 data is a white row after a white row, of 512 bytes here.
	var res *Result
	n := allocated(func() {
		res = converted(t, pagePDF("/XObject << /Im 5 0 R >>", [2]string{"", "/Im Do"},
			[2]string{"/Type /XObject /Subtype /Image /Width 4096 /Height 2 /BitsPerComponent 1 /Filter /CCITTFaxDecode /DecodeParms << /K -1 /Columns 4096 /Rows 2000000000 >>",
				strings.Repeat("\xff", 8<<10)}))
	})
	if n > 8<<20 || len(res.Warnings) != 0 || bodyOps(t, res)[bdf.OpImage] != 1 {
		t.Errorf("CCITT: %d bytes allocated, warnings %q", n, res.Warnings)
	}
}

func TestReduction(t *testing.T) {
	for _, c := range []struct {
		w, h  int
		limit int64
		k     int
	}{
		{10, 10, 100, 1},
		{10, 10, 99, 2},
		{10, 10, 25, 2},
		{10, 10, 24, 3}, // 4 × 4
		{13, 11, 42, 2}, // 7 × 6
		{13, 11, 41, 3}, // 5 × 4
		{1, 1000, 10, 100},
		{7, 5, 1, 7},
		{7, 5, 0, 7},
		{19866, 28087, imgconv.MaxDecodePixels, 3}, // A0 at 600 dpi: 6622 × 9363
		{65535, 65535, imgconv.MaxDecodePixels, 7},
	} {
		if k := reduction(c.w, c.h, c.limit); k != c.k {
			t.Errorf("%d × %d within %d pixels: reduced by %d, want %d", c.w, c.h, c.limit, k, c.k)
		}
	}
}

// bitsOf packs a text of 0 and 1 into bytes.
func bitsOf(text string) string {
	var out []byte
	for i := 0; i < len(text); i += 8 {
		var b byte
		for k := 0; k < 8; k++ {
			b <<= 1
			if i+k < len(text) && text[i+k] == '1' {
				b |= 1
			}
		}
		out = append(out, b)
	}
	return string(out)
}

// noise returns n bytes that look random, the same every time.
func noise(n int, seed uint32) []byte {
	out := make([]byte, n)
	for i := range out {
		seed = seed*1664525 + 1013904223
		out[i] = byte(seed >> 24)
	}
	return out
}

// averaged reduces an image by k, a pixel of the result at a time: the
// average of the pixels it stands for, the colours weighted by the opacity.
func averaged(img *image.NRGBA, k int) *image.NRGBA {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	out := image.NewNRGBA(image.Rect(0, 0, (w+k-1)/k, (h+k-1)/k))
	for oy := 0; oy*k < h; oy++ {
		for ox := 0; ox*k < w; ox++ {
			var r, g, b, a, n int
			for y := oy * k; y < (oy+1)*k && y < h; y++ {
				for x := ox * k; x < (ox+1)*k && x < w; x++ {
					c := img.NRGBAAt(x, y)
					r, g, b = r+int(c.R)*int(c.A), g+int(c.G)*int(c.A), b+int(c.B)*int(c.A)
					a += int(c.A)
					n++
				}
			}
			if alpha := (2*a + n) / (2 * n); alpha > 0 { // rounded
				out.SetNRGBA(ox, oy, color.NRGBA{uint8((2*r + a) / (2 * a)), uint8((2*g + a) / (2 * a)), uint8((2*b + a) / (2 * a)), uint8(alpha)})
			}
		}
	}
	return out
}

// imageOf reads the pixels of the image that is the object nr of a PDF,
// as an image of the first page; an image of more than limit pixels is
// reduced (0: the limit there is).
func imageOf(t *testing.T, data []byte, nr int, limit int64) (*imagePixels, *converter) {
	t.Helper()
	s, err := NewStream(bytes.NewReader(data), nil)
	if err != nil {
		t.Fatal(err)
	}
	c := s.c
	c.pdf.maxPixels = limit
	sd := c.pdf.stream(*types.NewIndirectRef(nr, 0))
	if sd == nil {
		t.Fatalf("object %d is no stream", nr)
	}
	px, err := c.loadPixels(sd.Dict, sd.Raw, sd.FilterPipeline, c.pageBodies[0].attrs.Resources, bdf.RGB(10, 200, 30))
	if err != nil {
		t.Fatal(err)
	}
	return px, c
}

func samePixels(t *testing.T, name string, got, want *image.NRGBA) {
	t.Helper()
	if got.Bounds() != want.Bounds() {
		t.Errorf("%s: %v, want %v", name, got.Bounds(), want.Bounds())
		return
	}
	for y := 0; y < want.Bounds().Dy(); y++ {
		for x := 0; x < want.Bounds().Dx(); x++ {
			if got.NRGBAAt(x, y) != want.NRGBAAt(x, y) {
				t.Errorf("%s: pixel %d,%d is %v, want %v", name, x, y, got.NRGBAAt(x, y), want.NRGBAAt(x, y))
				return
			}
		}
	}
}

// TestImageReduced: an image of more pixels than the limit is stored with
// fewer, each the average of the pixels of the image that it stands for,
// whatever the samples are encoded with and whatever colours they are.
func TestImageReduced(t *testing.T) {
	const w, h = 13, 11
	size := fmt.Sprintf("/Type /XObject /Subtype /Image /Width %d /Height %d ", w, h)
	rgb := noise(w*h*3, 1)
	// the rows of a PNG predictor (Up) and of the TIFF predictor
	var up, tiff []byte
	for y := 0; y < h; y++ {
		up = append(up, 2)
		for x := 0; x < 3*w; x++ {
			v := rgb[y*3*w+x]
			if y > 0 {
				up = append(up, v-rgb[(y-1)*3*w+x])
			} else {
				up = append(up, v)
			}
			if x >= 3 {
				tiff = append(tiff, v-rgb[y*3*w+x-3])
			} else {
				tiff = append(tiff, v)
			}
		}
	}
	lzw := func(b []byte) string {
		fi, _ := filter.NewFilter(filter.LZW, nil)
		r, err := fi.Encode(bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		out, _ := readAll(r)
		return string(out)
	}
	runs := func(b []byte) string {
		var out []byte
		for len(b) > 0 {
			n := min(len(b), 100)
			out = append(append(out, byte(n-1)), b[:n]...)
			b = b[n:]
		}
		return string(append(out, 0x80))
	}
	fixture, err := os.ReadFile("testdata/images-jpx-jbig2.pdf")
	if err != nil {
		t.Fatal(err)
	}
	jbig2 := map[bool]int{} // the JBIG2 images of the fixture: a mask and an image
	{
		s, err := NewStream(bytes.NewReader(fixture), nil)
		if err != nil {
			t.Fatal(err)
		}
		for nr := 1; nr <= *s.c.pdf.ctx.XRefTable.Size; nr++ {
			if sd := s.c.pdf.stream(*types.NewIndirectRef(nr, 0)); sd != nil && imageCodec(sd.FilterPipeline) == filter.JBIG2 {
				jbig2[s.c.pdf.boolOr(sd.Dict["ImageMask"], false)] = nr
			}
		}
		if len(jbig2) != 2 {
			t.Fatalf("JBIG2 images of the fixture: %v", jbig2)
		}
	}
	for _, c := range []struct {
		name, dict, data string
		more             []any
		pdf              []byte // instead of a PDF made of the above, with the image nr
		nr               int
		opaque           bool // every pixel of the image is
		fast             bool // the colours come from a table
	}{
		{name: "grey, not encoded", dict: "/BitsPerComponent 8 /ColorSpace /DeviceGray", data: string(noise(w*h, 2)), opaque: true, fast: true},
		{name: "RGB, Flate", dict: "/BitsPerComponent 8 /ColorSpace /DeviceRGB /Filter /FlateDecode", data: deflated(rgb), opaque: true},
		{name: "RGB, Flate with a PNG predictor", dict: fmt.Sprintf("/BitsPerComponent 8 /ColorSpace /DeviceRGB /Filter /FlateDecode /DecodeParms << /Predictor 12 /Colors 3 /Columns %d >>", w),
			data: deflated(up), opaque: true},
		{name: "RGB, Flate with the TIFF predictor", dict: fmt.Sprintf("/BitsPerComponent 8 /ColorSpace /DeviceRGB /Filter /FlateDecode /DecodeParms << /Predictor 2 /Colors 3 /Columns %d >>", w),
			data: deflated(tiff), opaque: true},
		{name: "grey of 4 bits, LZW", dict: "/BitsPerComponent 4 /ColorSpace /DeviceGray /Filter /LZWDecode", data: lzw(noise((w*4+7)/8*h, 3)), opaque: true, fast: true},
		{name: "grey of 2 bits, run length in ASCII85", dict: "/BitsPerComponent 2 /ColorSpace /DeviceGray /Filter [/ASCIIHexDecode /RunLengthDecode]",
			data: fmt.Sprintf("%x>", runs(noise((w*2+7)/8*h, 4))), opaque: true, fast: true},
		{name: "grey of 1 bit, inverted", dict: "/BitsPerComponent 1 /ColorSpace /DeviceGray /Decode [1 0]", data: string(noise((w+7)/8*h, 5)), opaque: true, fast: true},
		{name: "grey of 16 bits", dict: "/BitsPerComponent 16 /ColorSpace /DeviceGray", data: string(noise(w*h*2, 6)), opaque: true},
		{name: "CMYK", dict: "/BitsPerComponent 8 /ColorSpace /DeviceCMYK /Filter /FlateDecode", data: deflated(noise(w*h*4, 7)), opaque: true},
		{name: "indexed colours of 4 bits", dict: "/BitsPerComponent 4 /ColorSpace [/Indexed /DeviceRGB 15 <" + fmt.Sprintf("%x", noise(48, 8)) + ">]",
			data: string(noise((w*4+7)/8*h, 9)), opaque: true, fast: true},
		{name: "indexed colours of 8 bits, Flate", dict: "/BitsPerComponent 8 /ColorSpace [/Indexed /DeviceRGB 255 6 0 R] /Filter /FlateDecode",
			data: deflated(noise(w*h, 10)), more: []any{[2]string{"", string(noise(768, 11))}}, opaque: true, fast: true},
		{name: "ICC colours", dict: "/BitsPerComponent 8 /ColorSpace [/ICCBased 6 0 R] /Filter /FlateDecode", data: deflated(rgb),
			more: []any{[2]string{"/N 3", "profile"}}, opaque: true},
		{name: "a separation", dict: "/BitsPerComponent 8 /ColorSpace [/Separation /Spot /DeviceRGB 6 0 R]", data: string(noise(w*h, 12)),
			more: []any{`<< /FunctionType 2 /Domain [0 1] /C0 [1 1 1] /C1 [0.9 0.1 0.4] /N 1.5 >>`}, opaque: true, fast: true},
		{name: "Lab", dict: "/BitsPerComponent 8 /ColorSpace [/Lab << /WhitePoint [0.9505 1 1.089] /Range [-100 100 -100 100] >>]", data: string(rgb), opaque: true},
		{name: "data of half the rows", dict: "/BitsPerComponent 8 /ColorSpace /DeviceGray", data: string(noise(w*5+w/2, 13)), fast: true},
		{name: "stencil mask", dict: "/ImageMask true", data: string(noise((w+7)/8*h, 14))},
		{name: "stencil mask, inverted, Flate, of half the rows", dict: "/ImageMask true /Decode [1 0] /Filter /FlateDecode", data: deflated(noise((w+7)/8*5, 15))},
		{name: "CCITT Group 4", dict: "/Width 8 /Height 7 /BitsPerComponent 1 /ColorSpace /DeviceGray /Filter /CCITTFaxDecode /DecodeParms << /K -1 /Columns 8 /Rows 7 >>",
			// 3 white and 5 black pixels, a white pixel more in each of the next four rows
			data: bitsOf("001" + "1000" + "0011" + "0111" + "0111" + "0111" + "0111" + "11" + "11"), opaque: true, fast: true},
		{name: "CCITT Group 4, stencil mask", dict: "/Width 8 /Height 7 /ImageMask true /Filter /CCITTFaxDecode /DecodeParms << /K -1 /Columns 8 /Rows 7 >>",
			data: bitsOf("001" + "1000" + "0011" + "0111" + "0111" + "0111" + "0111" + "11" + "11")},
		{name: "JBIG2", pdf: fixture, nr: jbig2[false], opaque: true, fast: true},
		{name: "JBIG2, stencil mask", pdf: fixture, nr: jbig2[true]},
	} {
		data, nr := c.pdf, c.nr
		if data == nil {
			// (the size given first holds unless the dictionary has another)
			dict := size + c.dict
			if strings.Contains(c.dict, "/Width") {
				dict = "/Type /XObject /Subtype /Image " + c.dict
			}
			data, nr = pagePDF("/XObject << /Im 5 0 R >>", [2]string{"", "/Im Do"}, append([]any{[2]string{dict, c.data}}, c.more...)...), 5
		}
		full, conv := imageOf(t, data, nr, 0)
		if full.reduced || len(conv.warnings) != 0 {
			t.Fatalf("%s: reduced within the limit; warnings %q", c.name, conv.warnings)
		}
		iw, ih := full.img.Bounds().Dx(), full.img.Bounds().Dy()
		opaque, colours := true, map[color.NRGBA]bool{}
		for i := 0; i < len(full.img.Pix); i += 4 {
			opaque = opaque && full.img.Pix[i+3] == 255
			colours[color.NRGBA{full.img.Pix[i], full.img.Pix[i+1], full.img.Pix[i+2], full.img.Pix[i+3]}] = true
		}
		if opaque != c.opaque || len(colours) < 2 {
			t.Errorf("%s: opaque %v, %d colours", c.name, opaque, len(colours))
		}
		// the limit of the pixels of the image itself, and those at which it is reduced by 2, 3, 4 and to a pixel
		for _, k := range []int{1, 2, 3, 4, max(iw, ih)} {
			limit := int64((iw+k-1)/k) * int64((ih+k-1)/k)
			got, conv := imageOf(t, data, nr, limit)
			name := fmt.Sprintf("%s, within %d pixels", c.name, limit)
			if k == 1 {
				if got.reduced || len(conv.warnings) != 0 || !bytes.Equal(got.img.Pix, full.img.Pix) {
					t.Errorf("%s: reduced %v, warnings %q", name, got.reduced, conv.warnings)
				}
				continue
			}
			want := averaged(full.img, reduction(iw, ih, limit))
			samePixels(t, name, got.img, want)
			warning := fmt.Sprintf("an image of %d × %d pixels is stored at %d × %d", iw, ih, want.Bounds().Dx(), want.Bounds().Dy())
			if !got.reduced || !got.lossless || got.isMask != full.isMask || len(conv.warnings) != 1 || conv.warnings[0] != warning {
				t.Errorf("%s: %+v, warnings %q, want %q", name, got, conv.warnings, warning)
			}
		}
		// the table of the colours of the sample values gives the pixels that the samples give one by one
		if sd := conv.pdf.stream(*types.NewIndirectRef(nr, 0)); !full.isMask && imageCodec(sd.FilterPipeline) == "" {
			cs, err := conv.sampleColorSpace(sd.Dict, conv.pageBodies[0].attrs.Resources)
			if err != nil {
				t.Fatal(c.name, err)
			}
			samples, _, err := conv.pdf.decodeStream(sd)
			if err != nil {
				t.Fatal(c.name, err)
			}
			rows := sampleRows(samples, iw, ih, conv.pdf.intOr(sd.Dict["BitsPerComponent"], 8), cs, conv.pdf.nums(sd.Dict["Decode"]))
			if (rows.fast != nil) != c.fast {
				t.Errorf("%s: table of colours: %v", c.name, rows.fast != nil)
			}
			for y := 0; y < ih && rows.fast != nil; y++ {
				a, b := bytes.Repeat([]byte{7}, 4*iw), bytes.Repeat([]byte{9}, 4*iw)
				if ra, rb := rows.row(y, a), rows.fast(y, b); ra != rb || !bytes.Equal(a, b) || ra && !bytes.Equal(a, full.img.Pix[y*full.img.Stride:y*full.img.Stride+4*iw]) {
					t.Errorf("%s: row %d is %v (%v) from the table, %v (%v) from the samples", c.name, y, b, rb, a, ra)
					break
				}
			}
		}
	}
}

// TestReducedAfterData: the rows after those of the data are all the same,
// and are read once: an image of many rows and little data is little work.
func TestReducedAfterData(t *testing.T) {
	const w, h, k = 10, 100001, 4
	data := noise(w*6+3, 41) // six rows and a part of the seventh
	read := 0
	rows := sampleRows(data, w, h, 8, csGray, nil)
	row := rows.row
	rows.row, rows.fast = func(y int, dst []uint8) bool { read++; return row(y, dst) }, nil
	got := reduced(rows, k)
	if rows.inData != 7 || read > 4*k {
		t.Errorf("%d rows of data, %d rows read", rows.inData, read)
	}
	full := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h && row(y, full.Pix[y*full.Stride:y*full.Stride+4*w]); y++ {
	}
	samePixels(t, "grey", got, averaged(full, k))

	// a stencil mask paints the rows it has no data for, or none of them
	for _, decode := range [][]float64{nil, {1, 0}} {
		read = 0
		rows = stencilRows(data[:5], w, h, decode, bdf.RGB(1, 2, 3))
		row := rows.row
		rows.row = func(y int, dst []uint8) bool { read++; return row(y, dst) }
		got = reduced(rows, k)
		if rows.inData != 3 || read > 4*k {
			t.Errorf("stencil: %d rows of data, %d rows read", rows.inData, read)
		}
		full = image.NewNRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			row(y, full.Pix[y*full.Stride:y*full.Stride+4*w])
		}
		samePixels(t, fmt.Sprintf("stencil %v", decode), got, averaged(full, k))
		if last := got.NRGBAAt(0, got.Bounds().Dy()-1); (last == color.NRGBA{1, 2, 3, 255}) != (decode == nil) {
			t.Errorf("stencil %v: the last pixel is %v", decode, last)
		}
	}
}

// TestReducedWeights: the colours of the pixels count as much as the
// pixels are opaque.
func TestReducedWeights(t *testing.T) {
	const w, h = 9, 7
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	copy(img.Pix, noise(len(img.Pix), 51))
	copy(img.Pix, []uint8{255, 0, 0, 255, 0, 0, 255, 0}) // red, and blue that is not seen
	copy(img.Pix[img.Stride:], []uint8{0, 255, 0, 0, 0, 255, 0, 0})
	rows := imageRows{w: w, h: h, inData: h, row: func(y int, dst []uint8) bool {
		copy(dst, img.Pix[y*img.Stride:y*img.Stride+4*w])
		return true
	}}
	got := reduced(rows, 2)
	samePixels(t, "image", got, averaged(img, 2))
	if px := got.NRGBAAt(0, 0); px != (color.NRGBA{255, 0, 0, 64}) {
		t.Errorf("a red pixel of four: %v", px)
	}
	samePixels(t, "image, by 3", reduced(rows, 3), averaged(img, 3))
}

// TestSampleTable: the samples of 1 to 8 bits of an image of one component,
// read from a table of their colours.
func TestSampleTable(t *testing.T) {
	for bpc := 1; bpc <= 8; bpc++ {
		const w, h = 11, 4
		rowBytes := (w*bpc + 7) / 8
		data := noise(rowBytes*h-rowBytes/2, uint32(bpc)) // the last row is cut
		cs := &colorSpace{family: "Indexed", n: 1, base: csRGB, hival: 1<<bpc - 1, lookup: noise(3<<bpc, 77)}
		for _, rows := range []imageRows{sampleRows(data, w, h, bpc, csGray, nil), sampleRows(data, w, h, bpc, csGray, []float64{1, 0}), sampleRows(data, w, h, bpc, cs, nil)} {
			if rows.fast == nil || rows.inData != h {
				t.Fatalf("%d bits: no table, or %d rows", bpc, rows.inData)
			}
			for y := 0; y <= h; y++ {
				a, b := make([]byte, 4*w), bytes.Repeat([]byte{9}, 4*w)
				if ra, rb := rows.row(y, a), rows.fast(y, b); ra != rb || ra != (y < h) || !bytes.Equal(a, b) {
					t.Errorf("%d bits, row %d: %v (%v) from the table, %v (%v) from the samples", bpc, y, b, rb, a, ra)
				}
			}
		}
	}
}

// TestImageReducedMasks: the masks of an image that is reduced, and the
// masks that are reduced themselves.
func TestImageReducedMasks(t *testing.T) {
	const w, h = 13, 11
	rgb := noise(w*h*3, 21)
	soft, softSmall, stencil := noise(w*h, 22), noise(5*4, 23), noise((w+7)/8*h, 24)
	masked := func(masks string) [2]string {
		return [2]string{fmt.Sprintf("/Type /XObject /Subtype /Image /Width %d /Height %d /BitsPerComponent 8 /ColorSpace /DeviceRGB %s", w, h, masks), string(rgb)}
	}
	data := pagePDF("/XObject << /Im 5 0 R >>", [2]string{"", "/Im Do"},
		/* 5 */ masked("/SMask 8 0 R"),
		/* 6 */ masked("/SMask 9 0 R"),
		/* 7 */ masked("/Mask 10 0 R"),
		/* 8 */ [2]string{fmt.Sprintf("/Type /XObject /Subtype /Image /Width %d /Height %d /BitsPerComponent 8 /ColorSpace /DeviceGray", w, h), string(soft)},
		/* 9 */ [2]string{"/Type /XObject /Subtype /Image /Width 5 /Height 4 /BitsPerComponent 8 /ColorSpace /DeviceGray", string(softSmall)},
		/* 10 */ [2]string{fmt.Sprintf("/Type /XObject /Subtype /Image /Width %d /Height %d /ImageMask true", w, h), string(stencil)})
	// what the masks leave of a pixel, from the pixels of the mask (reduced or not)
	softOf := func(mask *image.NRGBA) func(x, y, w, h int) uint8 {
		mw, mh := mask.Bounds().Dx(), mask.Bounds().Dy()
		return func(x, y, w, h int) uint8 { return mask.NRGBAAt(x*mw/w, y*mh/h).R }
	}
	stencilOf := func(mask *image.NRGBA) func(x, y, w, h int) uint8 {
		mw, mh := mask.Bounds().Dx(), mask.Bounds().Dy()
		return func(x, y, w, h int) uint8 { return 255 - mask.NRGBAAt(x*mw/w, y*mh/h).A }
	}
	for _, c := range []struct {
		name     string
		nr, mask int
		left     func(mask *image.NRGBA) func(x, y, w, h int) uint8
	}{
		{"soft mask of the size of the image", 5, 8, softOf},
		{"soft mask of a size of its own", 6, 9, softOf},
		{"stencil mask", 7, 10, stencilOf},
	} {
		plain, _ := imageOf(t, data, 5, 0)
		for i := 3; i < len(plain.img.Pix); i += 4 {
			plain.img.Pix[i] = 255 // the image without its mask
		}
		fullMask, _ := imageOf(t, data, c.mask, 0)
		for _, k := range []int{1, 2, 3, 5} {
			limit := int64((w+k-1)/k) * int64((h+k-1)/k)
			got, conv := imageOf(t, data, c.nr, limit)
			// the image and the mask are reduced each to the limit
			want := plain.img
			if k > 1 {
				want = averaged(plain.img, k)
			}
			mask := fullMask.img
			if mk := reduction(mask.Bounds().Dx(), mask.Bounds().Dy(), limit); mk > 1 {
				mask = averaged(mask, mk)
			}
			left := c.left(mask)
			ow, oh := want.Bounds().Dx(), want.Bounds().Dy()
			for y := 0; y < oh; y++ {
				for x := 0; x < ow; x++ {
					px := want.NRGBAAt(x, y)
					px.A = left(x, y, ow, oh)
					want.SetNRGBA(x, y, px)
				}
			}
			samePixels(t, fmt.Sprintf("%s, within %d pixels", c.name, limit), got.img, want)
			if (len(conv.warnings) == 1) != (k > 1) {
				t.Errorf("%s, within %d pixels: warnings %q", c.name, limit, conv.warnings)
			}
			if k == 1 {
				want = nil // (made of plain: not to be written to again)
				plain, _ = imageOf(t, data, 5, 0)
				for i := 3; i < len(plain.img.Pix); i += 4 {
					plain.img.Pix[i] = 255
				}
			}
		}
	}

	// a mask of more pixels than the limit, of an image within it
	big := pagePDF("/XObject << /Im 5 0 R >>", [2]string{"", "/Im Do"},
		/* 5 */ [2]string{"/Type /XObject /Subtype /Image /Width 3 /Height 2 /BitsPerComponent 8 /ColorSpace /DeviceGray /Mask 6 0 R", string(noise(6, 25))},
		/* 6 */ [2]string{"/Type /XObject /Subtype /Image /Width 6 /Height 4 /ImageMask true", "\x0c\x0c\xf0\x00"})
	got, _ := imageOf(t, big, 5, 6)
	// The mask paints the zeros of 000011, 000011, 111100 and 000000: all of
	// the first two blocks of four pixels and nothing of the third, then
	// half of two blocks and all of the last.
	for i, a := range []uint8{0, 0, 255, 127, 127, 0} {
		if got.img.Pix[4*i+3] != a {
			t.Errorf("mask reduced, pixel %d: opacity %d, want %d", i, got.img.Pix[4*i+3], a)
		}
	}
}

// TestMaskOfMask: the masks that the image of a mask has are not read.
func TestMaskOfMask(t *testing.T) {
	grey := "/Type /XObject /Subtype /Image /Width 2 /Height 2 /BitsPerComponent 8 /ColorSpace /DeviceGray "
	data := pagePDF("/XObject << /Im 5 0 R >>", [2]string{"", "/Im Do"},
		/* 5 */ [2]string{grey + "/SMask 6 0 R", "\x00\x40\x80\xff"},
		/* 6 */ [2]string{grey + "/SMask 7 0 R", "\x10\x20\x30\x40"},
		/* 7 */ [2]string{"/Type /XObject /Subtype /Image /Height 2 /BitsPerComponent 8 /ColorSpace /DeviceGray", "\x00\x00"}) // without a width
	px, c := imageOf(t, data, 5, 0)
	if len(c.warnings) != 0 || c.inMask != 0 {
		t.Errorf("warnings %q", c.warnings)
	}
	for i, a := range []uint8{0x10, 0x20, 0x30, 0x40} {
		if px.img.Pix[4*i+3] != a {
			t.Errorf("pixel %d: opacity %#x, want %#x", i, px.img.Pix[4*i+3], a)
		}
	}
}

// TestImageReducedStored: images that are reduced are drawn where they
// were, and the document says once that there are such images.
func TestImageReducedStored(t *testing.T) {
	data := pagePDF("/XObject << /Im 5 0 R /Mk 6 0 R >>", [2]string{"", "q 300 0 0 200 50 60 cm /Im Do Q 0 0 1 rg q 100 0 0 100 0 0 cm /Mk Do Q BI /W 12 /H 8 /BPC 8 /CS /G ID " + string(noise(96, 31)) + " EI"},
		[2]string{"/Type /XObject /Subtype /Image /Width 12 /Height 8 /BitsPerComponent 8 /ColorSpace /DeviceRGB /Interpolate true", string(noise(288, 32))},
		[2]string{"/Type /XObject /Subtype /Image /Width 12 /Height 8 /ImageMask true /Interpolate true", string(noise(16, 33))})
	sizes := func(res *Result) []string {
		var out []string
		for _, p := range res.Doc.Parts() {
			if p.Type == bdf.PartImage {
				cfg, err := png.DecodeConfig(bytes.NewReader(p.Data))
				if err != nil {
					t.Fatal(err)
				}
				out = append(out, fmt.Sprintf("%d × %d", cfg.Width, cfg.Height))
			}
		}
		slices.Sort(out)
		return out
	}
	full := converted(t, data)
	r, body := opened(t, full)
	if got := sizes(full); len(full.Warnings) != 0 || !slices.Equal(got, []string{"12 × 8", "12 × 8", "12 × 8"}) {
		t.Fatalf("warnings %q, images of %q", full.Warnings, got)
	}
	drawn := markSeq(t, r, body)
	res, err := convertIn(t, 5*time.Second, data, nil, func(c *converter) { c.pdf.maxPixels = 24 })
	if err != nil {
		t.Fatal(err)
	}
	if got := sizes(res); len(res.Warnings) != 1 || res.Warnings[0] != "an image of 12 × 8 pixels is stored at 6 × 4" || !slices.Equal(got, []string{"6 × 4", "6 × 4", "6 × 4"}) {
		t.Errorf("warnings %q, images of %q", res.Warnings, got)
	}
	r, body = opened(t, res)
	if got := markSeq(t, r, body); got != drawn || strings.Count(got, "image") != 3 {
		t.Errorf("drawn %s, the images as they are %s", got, drawn)
	}
	for _, op := range []byte{bdf.OpTransform, bdf.OpImage, bdf.OpSave} {
		if got, want := opArgs(t, res, op), opArgs(t, full, op); got != want || want == "" {
			t.Errorf("instruction %#x: %s, the images as they are %s", op, got, want)
		}
	}
}

// opArgs lists the operands of the instructions op of the body of the first page.
func opArgs(t *testing.T, res *Result, op byte) string {
	t.Helper()
	r, body := opened(t, res)
	o, err := r.Object(body)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	o.Walk(func(in bdf.Instr) {
		if in.Op == op {
			out = append(out, fmt.Sprint(in.Args...))
		}
	})
	return strings.Join(out, "|")
}

// TestSharedReferences: what names the same object many times, at every
// depth, is read as often as there are objects, not references.
func TestSharedReferences(t *testing.T) {
	refs := func(obj, n int) string { return strings.TrimSpace(strings.Repeat(fmt.Sprintf("%d 0 R ", obj), n)) }

	t.Run("visibility expression", func(t *testing.T) {
		objs := []any{
			/* 1 */ `<< /Type /Catalog /Pages 2 0 R /OCProperties << /OCGs [6 0 R] /D << /BaseState /OFF >> >> >>`,
			/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
			/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Properties << /P0 5 0 R /P1 7 0 R >> >> >>`,
			/* 4 */ [2]string{"", "/OC /P0 BDC 0 0 10 10 re f EMC /OC /P1 BDC 0 0 10 10 re f EMC"},
			/* 5 */ `<< /Type /OCMD /VE 8 0 R >>`, // more expressions than the limit: not used, and no group: shown
			/* 6 */ `<< /Type /OCG /Name (g) >>`,
			/* 7 */ `<< /Type /OCMD /VE [/Not [/Or 6 0 R 6 0 R]] >>`, // off, off: shown
		}
		for i := 0; i < 16; i++ {
			objs = append(objs, "[/And "+refs(len(objs)+2, 8)+"]")
		}
		objs = append(objs, `[/And 6 0 R]`)
		res, err := convertIn(t, 5*time.Second, buildPDF(objs, "/Root 1 0 R"), nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		if n := bodyOps(t, res)[bdf.OpFillPath]; n != 2 {
			t.Errorf("%d paths filled, want 2", n)
		}
	})

	t.Run("table", func(t *testing.T) {
		objs := []any{
			/* 1 */ `<< /Type /Catalog /Pages 2 0 R /StructTreeRoot 5 0 R >>`,
			/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
			/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /StructParents 0 /Resources << /Font << /F1 6 0 R >> >> >>`,
			/* 4 */ [2]string{"", "/TD <</MCID 0>> BDC BT /F1 12 Tf (x) Tj ET EMC"},
			/* 5 */ `<< /Type /StructTreeRoot /K 8 0 R /ParentTree 7 0 R >>`,
			/* 6 */ `<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`,
			/* 7 */ `<< /Nums [0 [14 0 R]] >>`,
			/* 8 */ `<< /Type /StructElem /S /Table /P 5 0 R /K [` + refs(9, 40) + `] >>`,
			/* 9 */ `<< /Type /StructElem /S /THead /P 8 0 R /K [` + refs(10, 40) + `] >>`,
			/* 10 */ `<< /Type /StructElem /S /TBody /P 9 0 R /K [` + refs(11, 40) + `] >>`,
			/* 11 */ `<< /Type /StructElem /S /TFoot /P 10 0 R /K [` + refs(12, 40) + `] >>`,
			/* 12 */ `<< /Type /StructElem /S /THead /P 11 0 R /K [` + refs(13, 40) + ` 15 0 R] >>`,
			/* 13 */ `<< /Type /StructElem /S /TR /P 12 0 R /K [16 0 R] >>`,
			/* 14 */ `<< /Type /StructElem /S /TD /P 15 0 R /K 0 /Pg 3 0 R >>`,
			/* 15 */ `<< /Type /StructElem /S /TR /P 12 0 R /K [16 0 R 14 0 R] >>`,
			/* 16 */ `<< /Type /StructElem /S /TD /P 13 0 R >>`,
		}
		data := buildPDF(objs, "/Root 1 0 R")
		res, err := convertIn(t, 5*time.Second, data, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		// the second row of the two there are, its second cell
		r, body := opened(t, res)
		if got := markSeq(t, r, body); !strings.Contains(got, "B2") {
			t.Errorf("marks %s", got)
		}
	})

	t.Run("form of a tagged document", func(t *testing.T) {
		objs := []any{
			/* 1 */ `<< /Type /Catalog /Pages 2 0 R /StructTreeRoot 5 0 R >>`,
			/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
			/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 6 0 R >> /XObject << /Fm 7 0 R >> >> >>`,
			/* 4 */ [2]string{"", "/Fm Do\n/Span <</Lang (fr)>> BDC /Fm Do EMC\n"},
			/* 5 */ `<< /Type /StructTreeRoot >>`,
			/* 6 */ `<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`,
			/* 7 */ [2]string{"/Type /XObject /Subtype /Form /BBox [0 0 100 100] /Resources << /Font << /F1 6 0 R >> /XObject << /Fm 7 0 R >> >>",
				"BT /F1 12 Tf (x) Tj ET\n" + strings.Repeat("/Fm Do\n", 4)},
		}
		var c *converter
		res, err := convertIn(t, 10*time.Second, buildPDF(objs, "/Root 1 0 R"), &Options{NoTextIndex: true}, func(conv *converter) { c = conv })
		if err != nil {
			t.Fatal(err)
		}
		if !warned(res, "form XObjects converted for too many structures") || c.reconverted <= maxReconversions {
			t.Errorf("warnings %q, %d forms converted again", res.Warnings, c.reconverted)
		}
	})
}

func TestNesting(t *testing.T) {
	// q: the states beyond the limit are not kept, and their Q restore nothing
	content := strings.Repeat("q 1 0 0 1 1 0 cm ", maxSaveDepth+50) + "0 0 10 10 re f " + strings.Repeat("Q ", maxSaveDepth+50) + "0 0 10 10 re f"
	res := converted(t, pagePDF("", [2]string{"", content}))
	ops := bodyOps(t, res)
	if !warned(res, "q operators nested more than 1024 deep") || ops[bdf.OpSave] != maxSaveDepth || ops[bdf.OpRestore] != maxSaveDepth || ops[bdf.OpFillPath] != 2 {
		t.Errorf("warnings %q, instructions %v", res.Warnings, ops)
	}
	res = converted(t, pagePDF("", [2]string{"", strings.Repeat("q ", maxSaveDepth) + "0 0 10 10 re f " + strings.Repeat("Q ", maxSaveDepth)}))
	if len(res.Warnings) != 0 {
		t.Errorf("q within the limit: warnings %q", res.Warnings)
	}

	// marked content: the target and the text of the enclosing sequences
	// are kept, not searched for (30 s for these sequences before)
	const n = 200000
	content = strings.Repeat("/Span <</ActualText (A)>> BDC /P BMC ", n/2) + "BT /F1 12 Tf (x) Tj ET " + strings.Repeat("EMC ", n) + "BT (y) Tj ET"
	res, err := convertIn(t, 5*time.Second, pagePDF("/Font << /F1 5 0 R >>", [2]string{"/Filter /FlateDecode", deflated([]byte(content))},
		`<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	r, _ := opened(t, res)
	if got := pageText(t, r); !strings.Contains(got, "A") || !strings.Contains(got, "y") || strings.Contains(got, "x") {
		t.Errorf("text %q", got)
	}
}

func t1Encrypted(plain []byte, r uint16) []byte {
	out := make([]byte, len(plain))
	for i, p := range plain {
		c := p ^ byte(r>>8)
		r = (uint16(c)+r)*52845 + 22719
		out[i] = c
	}
	return out
}

// cffWithEncoding is a CFF program of two glyphs whose encoding has the
// codes 0 to 99 for the glyphs 1 to 100.
func cffWithEncoding() []byte {
	var top []cffDictEntry
	for _, op := range []int{15, 16, 17} {
		top = append(top, cffDictEntry{Op: op, Args: []float64{0}, Raw: appendDictNumber(nil, 0)})
	}
	top = append(top, cffDictEntry{Op: 18, Args: []float64{0, 0}, Raw: appendDictNumber(appendDictNumber(nil, 0), 0)})
	offsets := map[int][]int{}
	encodeTop := func() []byte {
		return cffEncodeDict(top, func(e cffDictEntry) ([]int, bool) {
			if v, ok := offsets[e.Op]; ok {
				return v, true
			}
			return make([]int, len(e.Args)), true
		})
	}
	head, names, strs, gsubrs := []byte{1, 0, 4, 4}, cffIndex([][]byte{[]byte("Few")}), cffIndex(nil), cffIndex(nil)
	charset := []byte{0, 0, 66} // glyph 1 is "a"
	enc := []byte{0, 100}
	for code := 0; code < 100; code++ {
		enc = append(enc, byte(code))
	}
	glyphs, private := cffIndex([][]byte{{14}, {14}}), []byte{139, 20}
	at := len(head) + len(names) + len(cffIndex([][]byte{encodeTop()})) + len(strs) + len(gsubrs)
	offsets[15] = []int{at}
	offsets[16] = []int{at + len(charset)}
	offsets[17] = []int{at + len(charset) + len(enc)}
	offsets[18] = []int{len(private), at + len(charset) + len(enc) + len(glyphs)}
	var out []byte
	for _, part := range [][]byte{head, names, cffIndex([][]byte{encodeTop()}), strs, gsubrs, charset, enc, glyphs, private} {
		out = append(out, part...)
	}
	return out
}

// TestDamagedValues: values that no file should have do not end the
// conversion.
func TestDamagedValues(t *testing.T) {
	type1 := func(private string) []byte {
		clear := "%!PS-AdobeFont-1.0: Test 1.0\n/FontName /Test def\n/FontMatrix [0.001 0 0 0.001 0 0] def\n/Encoding StandardEncoding def\ncurrentfile eexec\n"
		prog := append([]byte(clear), t1Encrypted([]byte("XXXX"+private), 55665)...)
		return pagePDF("/Font << /F1 5 0 R >>", [2]string{"", "BT /F1 12 Tf (a) Tj ET"},
			`<< /Type /Font /Subtype /Type1 /BaseFont /Test /FirstChar 97 /LastChar 97 /Widths [500] /FontDescriptor 6 0 R >>`,
			`<< /Type /FontDescriptor /FontName /Test /Flags 32 /FontFile 7 0 R >>`,
			[2]string{fmt.Sprintf("/Length1 %d", len(prog)), string(prog)})
	}
	cid := func(widths string) []byte {
		return pagePDF("/Font << /F1 5 0 R >>", [2]string{"", "BT /F1 12 Tf <0001> Tj ET"},
			`<< /Type /Font /Subtype /Type0 /BaseFont /X /Encoding /Identity-V /DescendantFonts [6 0 R] >>`,
			`<< /Type /Font /Subtype /CIDFontType2 /BaseFont /X /CIDSystemInfo << /Registry (Adobe) /Ordering (Identity) /Supplement 0 >> `+widths+` >>`)
	}
	sampled := func(space, dict string) []byte {
		return pagePDF("/ColorSpace << /CS0 "+space+" >>", [2]string{"", "/CS0 cs 1 1 1 scn 0 0 10 10 re f"},
			[2]string{"/FunctionType 0 /Range [0 1 0 1 0 1] /BitsPerSample 8 " + dict, strings.Repeat("\x00", 12)})
	}
	for _, c := range []struct {
		name string
		data []byte
	}{
		{"DeviceN without colorants", pagePDF("/ColorSpace << /CS0 [/DeviceN [] /DeviceRGB << /FunctionType 2 /Domain [0 1] /C0 [0 0 0] /C1 [1 1 1] /N 1 >>] >>",
			[2]string{"", "/CS0 cs 0 0 10 10 re f"})},
		{"sampled function without samples", sampled("[/Separation /X /DeviceRGB 5 0 R]", "/Domain [0 1] /Size [0]")},
		{"sampled function with fewer domains than inputs", sampled("[/DeviceN [/A /B] /DeviceRGB 5 0 R]", "/Domain [0 1] /Size [2 2]")},
		{"sampled function with more samples than an int counts", sampled("[/DeviceN [/A /B /C] /DeviceRGB 5 0 R]", "/Domain [0 1 0 1 0 1] /Size [3037000500 3037000500 4]")},
		{"sampled function of no bits", sampled("[/Separation /X /DeviceRGB 5 0 R]", "/Domain [0 1] /Size [2] /BitsPerSample 0")},
		{"inline image with a length beyond int", pagePDF("", [2]string{"", strings.Repeat("0 0 1 1 re f\n", 100) + "BI /W 1 /H 1 /BPC 8 /CS /G /F /AHx /L 9223372036854774784 ID 00> EI"})},
		{"inline image without data", pagePDF("", [2]string{"", "BI /W 1 /H 1 /BPC 8 /CS /G ID EI 0 0 1 1 re f"})},
		{"Type 1 charstring longer than the program", type1("dup /Private 8 dict dup begin /lenIV 4 def\n/CharStrings 1 dict dup begin\n/a 1e30 RD abcdefgh ND\nend\n")},
		{"Type 1 charstring of no length", type1("dup /Private 8 dict dup begin /lenIV 4 def\n/CharStrings 1 dict dup begin\n/a nan RD abcdefgh ND\nend\n")},
		{"Type 1 subroutine longer than the program", type1("dup /Private 8 dict dup begin /lenIV 4 def\n/Subrs 1 array\ndup 0 9223372036854775808 RD abcdefgh NP\n/CharStrings 1 dict dup begin\n/a 4 RD abcd ND\nend\n")},
		{"CFF encoding of more glyphs than the program has", pagePDF("/Font << /F1 5 0 R >>", [2]string{"", "BT /F1 12 Tf (2) Tj ET"},
			`<< /Type /Font /Subtype /Type1 /BaseFont /Few /FontDescriptor 6 0 R >>`,
			`<< /Type /FontDescriptor /FontName /Few /Flags 4 /FontFile3 7 0 R >>`,
			[2]string{"/Subtype /Type1C", string(cffWithEncoding())})},
		{"widths of CIDs beyond 2^53", cid("/W [9007199254740992 9007199254740994 500]")},
		{"vertical metrics of CIDs beyond 2^53", cid("/W2 [9007199254740992 9007199254740994 -1000 500 880]")},
		{"widths of CIDs below zero", cid("/W [-70000 -69990 500 1 3 600]")},
	} {
		t.Run(c.name, func(t *testing.T) {
			res := converted(t, c.data)
			if res.Pages != 1 || len(bodyOps(t, res)) == 0 {
				t.Errorf("%d pages, warnings %q", res.Pages, res.Warnings)
			}
		})
	}
	if first, last, ok := cidRangeOf(1, 3); !ok || first != 1 || last != 3 {
		t.Errorf("CIDs 1 to 3: %d to %d, %v", first, last, ok)
	}
	if _, _, ok := cidRangeOf(0, 65536); ok {
		t.Error("65537 CIDs are a range")
	}
}

// TestLookups: what is looked up for every link, form or character is
// not searched for from its start every time.
func TestLookups(t *testing.T) {
	t.Run("named destinations", func(t *testing.T) {
		const n = 300
		var annots, names strings.Builder
		objs := []any{
			/* 1 */ `<< /Type /Catalog /Pages 2 0 R /Names << /Dests 5 0 R >> >>`,
			/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
			/* 3 */ "",
			/* 4 */ [2]string{"", "0 0 10 10 re f"},
			/* 5 */ `<< /Kids [6 0 R 7 0 R] >>`,
			/* 6 */ "",
			/* 7 */ `<< /Names [(d000001) [3 0 R /Fit] (last) [3 0 R /Fit] (none) null] >>`,
		}
		for i := 0; i < n; i++ {
			fmt.Fprintf(&annots, "%d 0 R ", len(objs)+1)
			fmt.Fprintf(&names, "(d%06d) [3 0 R /Fit] ", i)
			objs = append(objs, fmt.Sprintf(`<< /Type /Annot /Subtype /Link /Rect [0 %d 5 %d] /Dest (d%06d) >>`, i, i+1, n-1-i))
		}
		objs[2] = `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Annots [` + annots.String() + `] >>`
		objs[5] = `<< /Names [` + names.String() + `] >>`
		var c *converter
		res, err := convertIn(t, 5*time.Second, buildPDF(objs, "/Root 1 0 R"), nil, func(conv *converter) { c = conv })
		if err != nil {
			t.Fatal(err)
		}
		if links := bodyOps(t, res)[bdf.OpLink]; links != n || c.pdf.treeReads != 3 {
			t.Errorf("%d links, %d nodes of the name tree read for them (it has 3)", links, c.pdf.treeReads)
		}
		if c.destTarget(types.StringLiteral("last")) != "#page=1" || c.destTarget(types.StringLiteral("none")) != "" || c.destTarget(types.StringLiteral("other")) != "" {
			t.Errorf("destinations: %v", c.dests)
		}
	})

	t.Run("objects of an object", func(t *testing.T) {
		const n = 30000 // 450 million steps before
		done := make(chan string, 1)
		go func() {
			parent := newPending(bdf.Rect{})
			kids := make([]*pending, n)
			for i := range kids {
				kids[i] = newPending(bdf.Rect{})
				if r := parent.childRef(kids[i]); int(r) != i {
					done <- fmt.Sprintf("object %d has the reference %d", i, r)
					return
				}
			}
			for i, k := range kids {
				if r := parent.childRef(k); int(r) != i || parent.children[r] != k {
					done <- fmt.Sprintf("object %d has the reference %d again", i, r)
					return
				}
			}
			// an object made of another has the references it was given
			again := (&converter{}).rebindPending(bdf.NewObject(), parent, nil, []int{2, -1, 2}, kids[0])
			if again.childRef(kids[2]) != 0 || again.childRef(kids[0]) != 1 {
				done <- fmt.Sprintf("references of the object made of another: %v", again.childRefs)
				return
			}
			done <- ""
		}()
		select {
		case msg := <-done:
			if msg != "" {
				t.Error(msg)
			}
		case <-time.After(3 * time.Second):
			t.Errorf("no references of %d objects within 3 s", n)
		}
	})

	t.Run("overlapping ranges of a CMap", func(t *testing.T) {
		var b strings.Builder
		b.WriteString("1 begincodespacerange <0000> <FFFF> endcodespacerange\n")
		const ranges = 60000
		for i := 0; i < ranges; i += 100 {
			b.WriteString("100 begincidrange\n")
			for k := i; k < i+100; k++ {
				fmt.Fprintf(&b, "<0000> <%04X> %d\n", k%0x8000, k%0x8000)
			}
			b.WriteString("endcidrange\n")
		}
		b.WriteString("1 begincidrange <FFF0> <FFFF> 7 endcidrange\n")
		cm := parseCMap([]byte(b.String()))
		if !cm.overlap || len(cm.cidRange) != ranges+1 {
			t.Fatalf("%d ranges, overlap %v", len(cm.cidRange), cm.overlap)
		}
		text := bytes.Repeat([]byte{0xff, 0xf1, 0x00, 0x05, 0x90, 0x00}, 100000) // 12 × 10^9 steps before
		done := make(chan []glyphCode, 1)
		go func() { done <- cm.decode(text) }()
		select {
		case codes := <-done:
			if len(codes) != 300000 || codes[0].cid != 8 || codes[1].cid < 5 || codes[1].cid != codes[4].cid || codes[2].cid != 0 || len(cm.found) != 3 {
				t.Errorf("%d codes, %v …, %d codes searched for", len(codes), codes[:3], len(cm.found))
			}
		case <-time.After(3 * time.Second):
			t.Error("not decoded within 3 s")
		}
	})
}

// TestMaskPixels: the pixels of a mask are those that were read back from
// the image stored for it before.
func TestMaskPixels(t *testing.T) {
	// masks of greys, of two colours, of many colours (which imgconv may
	// store lossy) and a JPEG
	gradient := make([]byte, 80*80)
	for i := range gradient {
		gradient[i] = byte(i)
	}
	photo := make([]byte, 80*80*3)
	for i, v := 0, uint32(1); i < len(photo); i++ {
		v = v*1664525 + 1013904223
		photo[i] = byte(v >> 24)
	}
	grey := image.NewGray(image.Rect(0, 0, 16, 16))
	copy(grey.Pix, gradient)
	var jpg bytes.Buffer
	if err := jpeg.Encode(&jpg, grey, nil); err != nil {
		t.Fatal(err)
	}
	fixture, err := os.ReadFile("testdata/chrome-masks.pdf")
	if err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{"chrome-masks.pdf": fixture}
	image4 := "/Type /XObject /Subtype /Image /Width 2 /Height 2 /BitsPerComponent 8 /ColorSpace /DeviceGray "
	files["masks"] = pagePDF("/XObject << /Im 5 0 R /Ip 8 0 R /Ij 10 0 R >>", [2]string{"", "/Im Do /Ip Do /Ij Do"},
		/* 5 */ [2]string{image4 + "/SMask 6 0 R /Mask 7 0 R", "\x00\x40\x80\xff"},
		/* 6 */ [2]string{"/Type /XObject /Subtype /Image /Width 80 /Height 80 /BitsPerComponent 8 /ColorSpace /DeviceGray", string(gradient)},
		/* 7 */ [2]string{"/Type /XObject /Subtype /Image /Width 8 /Height 2 /ImageMask true", "\xa5\x0f"},
		/* 8 */ [2]string{image4 + "/SMask 9 0 R", "\x00\x40\x80\xff"},
		/* 9 */ [2]string{"/Type /XObject /Subtype /Image /Width 80 /Height 80 /BitsPerComponent 8 /ColorSpace /DeviceRGB", string(photo)},
		/* 10 */ [2]string{image4 + "/SMask 11 0 R", "\x00\x40\x80\xff"},
		/* 11 */ [2]string{"/Type /XObject /Subtype /Image /Width 16 /Height 16 /BitsPerComponent 8 /ColorSpace /DeviceGray /Filter /DCTDecode", jpg.String()})
	masks, direct := 0, 0
	for name, data := range files {
		for _, mode := range []imgconv.Mode{imgconv.Keep, imgconv.Convert} {
			if name != "masks" && mode == imgconv.Keep {
				continue // the fixture as testdata has it, which is converted
			}
			s, err := NewStream(bytes.NewReader(data), &Options{Images: imgconv.Options{Mode: mode}})
			if err != nil {
				t.Fatal(name, err)
			}
			c := s.c
			for nr := 1; nr <= *c.pdf.ctx.XRefTable.Size; nr++ {
				sd := c.pdf.stream(*types.NewIndirectRef(nr, 0))
				if sd == nil || c.pdf.name(sd.Dict["Subtype"]) != "Image" {
					continue
				}
				for _, key := range []string{"SMask", "Mask"} {
					m := c.pdf.stream(sd.Dict[key])
					if m == nil {
						continue
					}
					masks++
					stored := c.stored
					mask, err := c.maskPixels(m, nil, bdf.RGB(0, 0, 0))
					if err != nil || mask == nil {
						t.Fatal(name, nr, err)
					}
					got, w, h := mask.img, mask.w, mask.h
					if c.stored == stored {
						direct++
					}
					// as before: the image is stored, and what is stored decoded
					px, err := c.loadPixels(m.Dict, m.Raw, m.FilterPipeline, nil, bdf.RGB(0, 0, 0))
					if err != nil {
						t.Fatal(name, nr, err)
					}
					di := c.store(px)
					want, err := decodeStored(di)
					if err != nil {
						t.Fatal(name, nr, err)
					}
					if w != di.w || h != di.h {
						t.Fatalf("%s %d %s: %d × %d, stored %d × %d", name, nr, key, w, h, di.w, di.h)
					}
					for y := 0; y < h; y++ {
						for x := 0; x < w; x++ {
							r0, g0, b0, a0 := got.At(x, y).RGBA()
							r1, g1, b1, a1 := want.At(x, y).RGBA()
							if r, a := mask.at(x, y); r != r0 || a != a0 {
								t.Fatalf("%s %d %s: pixel %d,%d of the mask is %d, %d, of its image %d, %d", name, nr, key, x, y, r, a, r0, a0)
							}
							if r0 != r1 || g0 != g1 || b0 != b1 || a0 != a1 {
								t.Fatalf("%s %d %s, mode %v: pixel %d,%d is %v, stored %v", name, nr, key, mode, x, y, got.At(x, y), want.At(x, y))
							}
						}
					}
				}
			}
		}
	}
	if masks < 8 || direct == 0 || direct == masks {
		t.Errorf("%d masks, %d taken as they were decoded", masks, direct)
	}
}
