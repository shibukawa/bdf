package image

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/shibukawa/bdf"
)

// convert converts an image of the test data.
func convert(t *testing.T, name string) *Result {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), nil)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// dcJSON writes Dublin Core as the manifest does.
func dcJSON(t *testing.T, dc bdf.DublinCore) string {
	t.Helper()
	b, err := json.Marshal(dc)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestConvert(t *testing.T) {
	for _, c := range []struct {
		file, format string
		w, h         float64
		animated     bool
		warnings     int
		dc           string
	}{
		// EXIF orientation 6 swaps the sides. XMP comes before EXIF, EXIF
		// before IPTC: the title and subjects are XMP's (its default
		// language), the creators and rights EXIF's; the camera's
		// placeholder description gives way to XMP's.
		{"photo.jpg", "jpeg", 48, 64, false, 0, `{"title":"Sunny field","creator":["Taro Yamada","Hanako Sato"],"subject":["sun","field","晴天"],` +
			`"description":"A red sun over a green field","rights":"© 2026 Taro Yamada","created":"2026-09-01T10:20:30+09:00","modified":"2026-09-02T08:00:00"}`},
		{"iptc.jpg", "jpeg", 8, 8, false, 0, `{"title":"夕焼け","creator":["Ann","Bob"],"subject":["sky","evening"],"description":"The sky at dusk",` +
			`"rights":"Ann and Bob","created":"2026-08-31T18:45:00+09:00"}`},
		{"tags.png", "png", 32, 24, false, 0, `{"title":"PNG の題名","creator":"Ann","description":"A blue rounded box","rights":"CC0","created":"2026-09-01T10:20:30+09:00"}`},
		{"anim.gif", "gif", 16, 16, true, 1, `{"description":"A blinking square"}`},
		// browsers do not turn WebP images by their EXIF orientation
		{"scene.webp", "webp", 64, 48, false, 1, `{"title":"Scene (WebP)","creator":"Ann"}`},
		// irot turns the image a quarter
		{"rotated.avif", "avif", 48, 64, false, 0, `{"title":"Scene (AVIF)","creator":"Bob","created":"2026-09-03T12:00:00Z"}`},
		{"flag.bmp", "bmp", 12, 8, false, 0, `{}`},
		{"icon.ico", "ico", 32, 32, false, 0, `{}`},
		// the RDF metadata comes before the title element
		{"drawing.svg", "svg", 320, 240, false, 0, `{"title":"Shapes and hairlines","creator":"Hanako Sato","subject":["test","vector"],` +
			`"description":"A gradient card with a circle, a star and hairlines","date":"2026-09-26","language":"en","rights":"CC0"}`},
	} {
		t.Run(c.file, func(t *testing.T) {
			res := convert(t, c.file)
			if res.Format != c.format || res.Width != c.w || res.Height != c.h || res.Animated != c.animated {
				t.Errorf("got %s %v × %v animated %v, want %s %v × %v animated %v", res.Format, res.Width, res.Height, res.Animated, c.format, c.w, c.h, c.animated)
			}
			if len(res.Warnings) != c.warnings {
				t.Errorf("warnings %q, want %d", res.Warnings, c.warnings)
			}
			if got := dcJSON(t, res.Doc.Meta.DC); got != c.dc {
				t.Errorf("dc\n got %s\nwant %s", got, c.dc)
			}
			doc := res.Doc
			if doc.Meta.Source != c.format {
				t.Errorf("source %q", doc.Meta.Source)
			}
			// one page the size of the image (96 px per inch), which draws
			// the stored image, as it is, over the whole page
			page := doc.Views[0].Pages[0]
			if page.W != float32(c.w*0.75) || page.H != float32(c.h*0.75) {
				t.Errorf("page %v × %v", page.W, page.H)
			}
			data, _ := os.ReadFile(filepath.Join("testdata", c.file))
			obj, err := bdf.DecodeObject(doc.Part(page.Layers[0].Obj).Data)
			if err != nil {
				t.Fatal(err)
			}
			if len(obj.Images) != 1 || !bytes.Equal(doc.Part(obj.Images[0]).Data, data) {
				t.Fatal("the image is not stored as it is")
			}
			ins, _ := obj.Instructions()
			var img *bdf.Instr
			for i := range ins {
				if ins[i].Op == bdf.OpImage {
					img = &ins[i]
				}
			}
			if img == nil || img.Args[3] != page.W || img.Args[4] != page.H {
				t.Errorf("image instruction %v", img)
			}
		})
	}
}

// TestAltText: the image is a figure whose alternative text is the
// description, else the title, else the file name.
func TestAltText(t *testing.T) {
	alt := func(res *Result) string {
		doc := res.Doc
		obj, _ := bdf.DecodeObject(doc.Part(doc.Views[0].Pages[0].Layers[0].Obj).Data)
		ins, _ := obj.Instructions()
		if len(ins) == 0 || ins[0].Op != bdf.OpMark || ins[0].Args[0] != uint64(bdf.MarkFigure) {
			t.Fatalf("no figure: %v", ins)
		}
		return ins[0].Args[1].(string)
	}
	if got := alt(convert(t, "photo.jpg")); got != "A red sun over a green field" {
		t.Errorf("photo: %q", got)
	}
	if got := alt(convert(t, "scene.webp")); got != "Scene (WebP)" {
		t.Errorf("webp: %q", got)
	}
	if got := alt(convert(t, "flag.bmp")); got != "flag.bmp" {
		t.Errorf("bmp: %q", got)
	}
	res, err := ConvertFile(filepath.Join("testdata", "flag.bmp"), &Options{Title: "Flag"})
	if err != nil {
		t.Fatal(err)
	}
	if got := alt(res); got != "Flag" || res.Doc.Meta.DC.Title.First() != "Flag" || res.Doc.Views[0].Title != "Flag" {
		t.Errorf("title option: alt %q, title %q", got, res.Doc.Meta.DC.Title)
	}
}

// TestTruncated: images cut anywhere fail or convert, without panicking.
func TestTruncated(t *testing.T) {
	files, _ := filepath.Glob("testdata/*")
	for _, f := range files {
		data, _ := os.ReadFile(f)
		for n := range len(data) {
			Convert(bytes.NewReader(data[:n]), int64(n), nil)
		}
	}
}

func TestSniff(t *testing.T) {
	for _, c := range []struct {
		name, data, want string
	}{
		{"svg", `<svg xmlns="http://www.w3.org/2000/svg"/>`, "svg"},
		{"prolog", "\xef\xbb\xbf<?xml version=\"1.0\"?>\n<!-- a comment -->\n<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"x.dtd\" [\n<!ENTITY a \"<b>\">\n]>\n<svg>", "svg"},
		{"prefixed", `<svg:svg xmlns:svg="http://www.w3.org/2000/svg"/>`, "svg"},
		{"html", `<!DOCTYPE html><html><svg/></html>`, ""},
		{"xml", `<?xml version="1.0"?><mxfile/>`, ""},
		{"text", `svg`, ""},
		{"bm text", "BM is not a bitmap", ""},
		{"heic", "\x00\x00\x00\x18ftypheic\x00\x00\x00\x00mif1heic", ""},
		{"avif brand", "\x00\x00\x00\x1cftypmif1\x00\x00\x00\x00mif1avifmiaf", "avif"},
	} {
		if got := sniff([]byte(c.data), strings.NewReader(c.data), int64(len(c.data))); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
	// the root element after a long comment
	long := "<!--" + strings.Repeat("x", 3000) + "--><svg/>"
	if got := sniff([]byte(long[:1024]), strings.NewReader(long), int64(len(long))); got != "svg" {
		t.Errorf("long comment: %q", got)
	}
}

func TestSVGSize(t *testing.T) {
	for _, c := range []struct {
		width, height, viewBox string
		w, h                   float64
	}{
		{"200", "100", "", 200, 100},
		{"200px", "1in", "0 0 10 10", 200, 96},
		{"72pt", "2.54cm", "", 96, 96},
		{"1e2", "5em", "", 100, 80},
		{"", "", "0 0 40 30", 40, 30},
		{"80", "", "0,0,40,30", 80, 60},
		{"100%", "50", "0 0 40 20", 100, 50},
		{"", "", "", 300, 150},
		{"50", "", "", 50, 150},
		{"-5", "abc", "0 0 0 0", 300, 150},
	} {
		if w, h := svgSize(c.width, c.height, c.viewBox); w != c.w || h != c.h {
			t.Errorf("svgSize(%q, %q, %q) = %v × %v, want %v × %v", c.width, c.height, c.viewBox, w, h, c.w, c.h)
		}
	}
}

// TestSVGVariants: an Illustrator SVG that names its namespace with an
// entity, a UTF-16 document, a Latin-1 one, and an SVG without the SVG
// namespace, which browsers do not draw.
func TestSVGVariants(t *testing.T) {
	illustrator := `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd" [
	<!ENTITY ns_svg "http://www.w3.org/2000/svg">
]>
<svg version="1.1" xmlns="&ns_svg;" width="64" height="32"><title>Illustrated</title></svg>`
	p, err := readSVG([]byte(illustrator))
	if err != nil || len(p.warnings) != 0 || p.w != 64 || p.h != 32 || p.native.Title.First() != "Illustrated" {
		t.Errorf("illustrator: %v %+v", err, p)
	}
	u := utf16.Encode([]rune(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 8 4"><desc>UTF-16 の説明</desc></svg>`))
	b := []byte{0xff, 0xfe}
	for _, c := range u {
		b = binary.LittleEndian.AppendUint16(b, c)
	}
	if p, err := readSVG(b); err != nil || p.w != 8 || p.native.Description.First() != "UTF-16 の説明" {
		t.Errorf("utf-16: %v %+v", err, p)
	}
	latin := []byte("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><svg xmlns=\"http://www.w3.org/2000/svg\"><title>Caf\xe9</title></svg>")
	if p, err := readSVG(latin); err != nil || p.native.Title.First() != "Café" {
		t.Errorf("latin-1: %v %+v", err, p)
	}
	if p, err := readSVG([]byte(`<svg width="1" height="1"/>`)); err != nil || len(p.warnings) != 1 {
		t.Errorf("no namespace: %v %+v", err, p)
	}
}

func TestDates(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"2026", "2026"},
		{"2026-09", "2026-09"},
		{"2026-09-01T10:20", "2026-09-01T10:20"},
		{"2026-09-01T10:20:30.25+0900", "2026-09-01T10:20:30.25+09:00"},
		{"2026:09:01 10:20:30Z", "2026-09-01T10:20:30Z"},
		{"2026-13-01", ""},
		{"yesterday", ""},
	} {
		if got := isoDate(c.in); got != c.want {
			t.Errorf("isoDate(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	for _, c := range []struct{ in, offset, want string }{
		{"2026:09:01 10:20:30", "+09:00", "2026-09-01T10:20:30+09:00"},
		{"2026:09:01 10:20:30", "", "2026-09-01T10:20:30"},
		{"    :  :     :  :  ", "", ""},
		{"0000:00:00 00:00:00", "", ""},
		{"2026:09:01 10:20:30", "+9", "2026-09-01T10:20:30"},
	} {
		if got := exifDate(c.in, c.offset); got != c.want {
			t.Errorf("exifDate(%q, %q) = %q, want %q", c.in, c.offset, got, c.want)
		}
	}
	if got := textDate("Tue, 01 Sep 2026 10:20:30 GMT"); got != "2026-09-01T10:20:30Z" {
		t.Errorf("textDate: %q", got)
	}
}

// TestCompressedParts: SVG and BMP parts are compressed (they are not
// compressed data), other images are stored as they are.
func TestCompressedParts(t *testing.T) {
	for file, want := range map[string]string{"drawing.svg": bdf.EncDeflateRaw, "icon.ico": bdf.EncDeflateRaw, "photo.jpg": bdf.EncIdentity} {
		res := convert(t, file)
		var buf bytes.Buffer
		if err := res.Doc.WriteSingle(&buf); err != nil {
			t.Fatal(err)
		}
		r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range r.Manifest.Parts {
			if p.T == bdf.PartImage && p.Enc != want {
				t.Errorf("%s: image part %s, want %s", file, p.Enc, want)
			}
		}
	}
}

// FuzzConvert: any input fails or converts, without panicking. The seeds
// are the test images.
func FuzzConvert(f *testing.F) {
	files, _ := filepath.Glob("testdata/*")
	for _, file := range files {
		data, _ := os.ReadFile(file)
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		Convert(bytes.NewReader(data), int64(len(data)), nil)
	})
}
