package jpx

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "regenerate testdata with OpenJPEG (opj_compress, opj_decompress)")

// fixture is a codestream in testdata, made by opj_compress from a
// synthetic source; its reference decode by opj_decompress is either a
// SHA-256 of the samples (exact) or the samples themselves (lossy, in
// name.ref.gz).
type fixture struct {
	name, src, args string
	lossy           bool
}

var fixtures = []fixture{
	{"rgb8.j2k", "rgb8", "-n 5", false},
	{"rgb8_97.j2k", "rgb8", "-I -r 10 -n 5", true},
	{"tiles.j2k", "rgb8", "-t 20,17 -T 3,2 -d 5,4 -n 4 -c [16,16],[8,8] -b 8,8 -p RPCL -SOP -EPH -r 30,10,2", false},
	{"tiles_97.j2k", "rgb8", "-I -t 32,32 -T 7,11 -d 9,13 -r 15 -M 9 -p PCRL", true},
	{"styles.j2k", "gray16", "-n 5 -M 63 -r 20,5,1", false},
	{"bypass.j2k", "rgb8", "-M 5 -r 12,4,1 -p RLCP", false},
	{"sub4.j2k", "sub4", "-p CPRL -n 3 -c [16,16],[8,8] -r 10,1", false},
	{"signed12.j2k", "signed12", "-t 16,16 -d 3,1 -n 3 -p PCRL", false},
	{"roi.j2k", "rgb8", "-ROI c=0,U=4 -r 20,1", false},
	{"poc.j2k", "rgb8", "-POC T1=0,0,3,2,3,CPRL/T1=2,0,3,6,3,LRCP -r 10,5,1", false},
	{"tileparts.j2k", "rgb8", "-TP R -r 20,10,1 -t 32,32 -n 4", false},
	{"gray8.j2k", "gray8", "-n 3", false},
	{"gray8.jp2", "gray8", "-I -r 10", true},
	{"gray16_97.j2k", "gray16", "-I -n 4 -r 5", true},
	{"tiny.j2k", "tiny", "-n 2", false},
}

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// decodePlanes decodes a codestream or JP2 file to component planes,
// before any palette, channel definition or resampling.
func decodePlanes(data []byte) ([]plane, error) {
	cs := data
	if bytes.HasPrefix(data, jp2Signature) {
		var err error
		if _, cs, err = parseJP2(data); err != nil {
			return nil, err
		}
	}
	c, err := parseCodestream(cs)
	if err != nil {
		return nil, err
	}
	return c.decode()
}

// canonical serializes planes: per component its size, then its samples.
func canonical(planes []plane) []byte {
	var b []byte
	for _, p := range planes {
		b = binary.LittleEndian.AppendUint32(b, uint32(p.w))
		b = binary.LittleEndian.AppendUint32(b, uint32(p.h))
		for _, v := range p.data {
			b = binary.LittleEndian.AppendUint32(b, uint32(v))
		}
	}
	return b
}

// compact serializes planes for reference files: per component its size,
// then the differences between successive samples as varints.
func compact(planes []plane) []byte {
	var b []byte
	for _, p := range planes {
		b = binary.AppendUvarint(b, uint64(p.w))
		b = binary.AppendUvarint(b, uint64(p.h))
		prev := int32(0)
		for _, v := range p.data {
			b = binary.AppendVarint(b, int64(v-prev))
			prev = v
		}
	}
	return b
}

func fromCompact(b []byte) ([]plane, error) {
	var planes []plane
	r := bytes.NewReader(b)
	for r.Len() > 0 {
		w, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}
		h, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}
		p := plane{w: int(w), h: int(h), data: make([]int32, w*h)}
		prev := int32(0)
		for i := range p.data {
			d, err := binary.ReadVarint(r)
			if err != nil {
				return nil, err
			}
			prev += int32(d)
			p.data[i] = prev
		}
		planes = append(planes, p)
	}
	return planes, nil
}

func hashPlanes(planes []plane) string {
	s := sha256.Sum256(canonical(planes))
	return hex.EncodeToString(s[:])
}

// maxDiff returns the largest sample difference between two decodes.
func maxDiff(a, b []plane) (int, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("%d components, want %d", len(a), len(b))
	}
	d := 0
	for i := range a {
		if a[i].w != b[i].w || a[i].h != b[i].h {
			return 0, fmt.Errorf("component %d is %dx%d, want %dx%d", i, a[i].w, a[i].h, b[i].w, b[i].h)
		}
		for j, v := range a[i].data {
			d = max(d, abs(int(v)-int(b[i].data[j])))
		}
	}
	return d, nil
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func readRefs(t *testing.T) map[string]string {
	t.Helper()
	refs := map[string]string{}
	sc := bufio.NewScanner(bytes.NewReader(readTestdata(t, "refs.txt")))
	for sc.Scan() {
		if f := strings.Fields(sc.Text()); len(f) == 2 && !strings.HasPrefix(f[0], "#") {
			refs[f[0]] = f[1]
		}
	}
	return refs
}

func readRefPlanes(t *testing.T, name string) []plane {
	t.Helper()
	r, err := gzip.NewReader(bytes.NewReader(readTestdata(t, name)))
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	planes, err := fromCompact(b)
	if err != nil {
		t.Fatal(err)
	}
	return planes
}

func TestFixtures(t *testing.T) {
	refs := readRefs(t)
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			data := readTestdata(t, f.name)
			planes, err := decodePlanes(data)
			if err != nil {
				t.Fatal(err)
			}
			if !f.lossy {
				if got := hashPlanes(planes); got != refs[f.name] {
					t.Errorf("decode differs from OpenJPEG: sha256 %s, want %s", got, refs[f.name])
				}
				return
			}
			want := readRefPlanes(t, f.name+".ref.gz")
			d, err := maxDiff(planes, want)
			if err != nil {
				t.Fatal(err)
			}
			// OpenJPEG approximates 1/K of the 9/7 filter; with its value
			// the decodes agree up to rounding.
			invK97 = openJPEGInvK
			defer func() { invK97 = 1 / k97 }()
			planes, err = decodePlanes(data)
			if err != nil {
				t.Fatal(err)
			}
			dOpj, err := maxDiff(planes, want)
			if err != nil {
				t.Fatal(err)
			}
			if dOpj > 1 || d > 4 {
				t.Errorf("max deviation from OpenJPEG %d, %d with its 1/K", d, dOpj)
			}
			t.Logf("max deviation from OpenJPEG %d, %d with its 1/K", d, dOpj)
		})
	}
}

// TestUpdateFixtures regenerates testdata:
//
//	go test ./converter/internal/jpx -run UpdateFixtures -update
func TestUpdateFixtures(t *testing.T) {
	if !*update {
		t.Skip("run with -update to regenerate testdata")
	}
	dir := t.TempDir()
	files := map[string][]string{}
	for _, s := range opjSources {
		p, extra := s.file(dir)
		files[s.name] = append([]string{"-i", p}, extra...)
	}
	var manifest strings.Builder
	manifest.WriteString("# Reference decodes by opj_decompress; regenerate with go test -run UpdateFixtures -update.\n")
	for _, f := range fixtures {
		out := filepath.Join("testdata", f.name)
		args := append(append(slices.Clone(files[f.src]), "-o", out), strings.Fields(f.args)...)
		if b, err := exec.Command("opj_compress", args...).CombinedOutput(); err != nil {
			t.Fatalf("opj_compress %v: %v\n%s", args, err, b)
		}
		ref := filepath.Join(dir, f.name+".pgx")
		if b, err := exec.Command("opj_decompress", "-i", out, "-o", ref).CombinedOutput(); err != nil {
			t.Fatalf("opj_decompress: %v\n%s", err, b)
		}
		var planes []plane
		for i := 0; ; i++ {
			w, h, _, _, data, err := readPGX(fmt.Sprintf("%s_%d.pgx", strings.TrimSuffix(ref, ".pgx"), i))
			if errors.Is(err, os.ErrNotExist) {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			planes = append(planes, plane{w: w, h: h, data: data})
		}
		if !f.lossy {
			fmt.Fprintf(&manifest, "%s %s\n", f.name, hashPlanes(planes))
			continue
		}
		var b bytes.Buffer
		zw, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)
		zw.Write(compact(planes))
		zw.Close()
		must(os.WriteFile(filepath.Join("testdata", f.name+".ref.gz"), b.Bytes(), 0o644))
	}
	must(os.WriteFile(filepath.Join("testdata", "refs.txt"), []byte(manifest.String()), 0o644))
}

// TestImage checks the public API on a codestream with offsets and on
// subsampled components.
func TestImage(t *testing.T) {
	for _, name := range []string{"tiles.j2k", "sub4.j2k", "signed12.j2k"} {
		data := readTestdata(t, name)
		img, err := Decode(data)
		if err != nil {
			t.Fatal(name, err)
		}
		c, _ := parseCodestream(data)
		planes, _ := c.decode()
		if img.Width != c.x1-c.x0 || img.Height != c.y1-c.y0 || len(img.Components) != len(planes) || img.Alpha != -1 {
			t.Fatalf("%s: %dx%d, %d components, alpha %d", name, img.Width, img.Height, len(img.Components), img.Alpha)
		}
		for i, cp := range img.Components {
			sc := c.comps[i]
			if cp.Precision != sc.prec || cp.Signed != sc.signed || len(cp.Data) != img.Width*img.Height {
				t.Fatalf("%s component %d: precision %d signed %v, %d samples", name, i, cp.Precision, cp.Signed, len(cp.Data))
			}
			pl := planes[i]
			for y := range img.Height {
				for x := range img.Width {
					sx := min(max((c.x0+x)/sc.dx-pl.x0, 0), pl.w-1)
					sy := min(max((c.y0+y)/sc.dy-pl.y0, 0), pl.h-1)
					if got, want := cp.Data[y*img.Width+x], pl.data[sy*pl.w+sx]; got != want {
						t.Fatalf("%s component %d at %d,%d: %d, want %d", name, i, x, y, got, want)
					}
				}
			}
		}
	}
}

// TestPackedHeaders moves packet headers into PPM and PPT marker segments.
func TestPackedHeaders(t *testing.T) {
	for _, name := range []string{"tiles.j2k", "tileparts.j2k", "poc.j2k", "styles.j2k", "tiles_97.j2k"} {
		data := readTestdata(t, name)
		want, err := decodePlanes(data)
		if err != nil {
			t.Fatal(err)
		}
		for _, ppm := range []bool{false, true} {
			for _, size := range []int{65000, 37} {
				packed, err := packHeaders(data, ppm, size)
				if err != nil {
					t.Fatal(err)
				}
				marker := []byte{0xff, 0x61}
				if ppm {
					marker = []byte{0xff, 0x60}
				}
				if !bytes.Contains(packed, marker) {
					t.Fatalf("%s: no packed headers", name)
				}
				got, err := decodePlanes(packed)
				if err != nil {
					t.Fatalf("%s ppm=%v size %d: %v", name, ppm, size, err)
				}
				if d, err := maxDiff(got, want); err != nil || d != 0 {
					t.Errorf("%s ppm=%v size %d: max deviation %d, %v", name, ppm, size, d, err)
				}
			}
		}
	}
}

// TestMarkerPrecedence moves coding and quantization parameters between the
// main and tile-part headers (A.6: tile-part COC > tile-part COD > main COC
// > main COD, and likewise for QCC and QCD).
func TestMarkerPrecedence(t *testing.T) {
	data := readTestdata(t, "tiles.j2k")
	want, err := decodePlanes(data)
	if err != nil {
		t.Fatal(err)
	}
	s, err := splitStream(data)
	if err != nil {
		t.Fatal(err)
	}
	cod, qcd := find(s.main, mCOD), find(s.main, mQCD)
	scod, sp := spcod(cod)
	// A wrong but valid COD: 2 levels, 16x16 code-blocks, default precincts,
	// 9/7 filter.
	wrongCOD := append(slices.Clone(cod[:5]), 2, 2, 2, 0, 0)
	wrongCOD[0] = scod &^ 1
	wrongSP := wrongCOD[5:]
	wrongQCD := []byte{0x22, 0x40, 0x00}
	coc := func(c int, sp []byte) rawSeg {
		scoc := scod & 1
		if len(sp) == 5 {
			scoc = 0 // no precinct sizes
		}
		return rawSeg{mCOC, append([]byte{byte(c), scoc}, sp...)}
	}
	qcc := func(c int, q []byte) rawSeg { return rawSeg{mQCC, append([]byte{byte(c)}, q...)} }
	ncomp := 3
	tests := []struct {
		name string
		edit func(s *rawStream)
	}{
		{"main COC over main COD", func(s *rawStream) {
			s.main = without(s.main, mCOD)
			s.main = append(s.main, rawSeg{mCOD, wrongCOD})
			for c := range ncomp {
				s.main = append(s.main, coc(c, sp))
			}
		}},
		{"tile COD over main COC", func(s *rawStream) {
			for c := range ncomp {
				s.main = append(s.main, coc(c, wrongSP))
			}
			for i := range s.parts {
				if s.parts[i].tpsot == 0 {
					s.parts[i].hdr = append(s.parts[i].hdr, rawSeg{mCOD, cod})
				}
			}
		}},
		{"tile COC over tile COD", func(s *rawStream) {
			for i := range s.parts {
				if s.parts[i].tpsot == 0 {
					s.parts[i].hdr = append(s.parts[i].hdr, rawSeg{mCOD, wrongCOD})
					for c := range ncomp {
						s.parts[i].hdr = append(s.parts[i].hdr, coc(c, sp))
					}
				}
			}
		}},
		{"main QCC over main QCD", func(s *rawStream) {
			s.main = without(s.main, mQCD)
			s.main = append(s.main, rawSeg{mQCD, wrongQCD})
			for c := range ncomp {
				s.main = append(s.main, qcc(c, qcd))
			}
		}},
		{"tile QCD over main QCC", func(s *rawStream) {
			for c := range ncomp {
				s.main = append(s.main, qcc(c, wrongQCD))
			}
			for i := range s.parts {
				if s.parts[i].tpsot == 0 {
					s.parts[i].hdr = append(s.parts[i].hdr, rawSeg{mQCD, qcd})
				}
			}
		}},
		{"tile QCC over tile QCD", func(s *rawStream) {
			for i := range s.parts {
				if s.parts[i].tpsot == 0 {
					s.parts[i].hdr = append(s.parts[i].hdr, rawSeg{mQCD, wrongQCD})
					for c := range ncomp {
						s.parts[i].hdr = append(s.parts[i].hdr, qcc(c, qcd))
					}
				}
			}
		}},
	}
	for _, tt := range tests {
		s, _ := splitStream(data)
		tt.edit(&s)
		got, err := decodePlanes(s.bytes())
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if d, err := maxDiff(got, want); err != nil || d != 0 {
			t.Errorf("%s: max deviation %d, %v", tt.name, d, err)
		}
	}

	// POC and RGN work from either header.
	for _, name := range []string{"poc.j2k", "roi.j2k"} {
		data := readTestdata(t, name)
		want, err := decodePlanes(data)
		if err != nil {
			t.Fatal(err)
		}
		s, _ := splitStream(data)
		moved := false
		for _, m := range []int{mPOC, mRGN} {
			if b := find(s.main, m); b != nil {
				s.main = without(s.main, m)
				for i := range s.parts {
					if s.parts[i].tpsot == 0 {
						s.parts[i].hdr = append(s.parts[i].hdr, rawSeg{m, b})
					}
				}
				moved = true
				continue
			}
			for i := range s.parts {
				if b := find(s.parts[i].hdr, m); b != nil {
					s.parts[i].hdr = without(s.parts[i].hdr, m)
					s.main = append(s.main, rawSeg{m, b})
					moved = true
				}
			}
		}
		if !moved {
			t.Fatalf("%s: no POC or RGN marker", name)
		}
		got, err := decodePlanes(s.bytes())
		if err != nil {
			t.Fatal(err)
		}
		if d, err := maxDiff(got, want); err != nil || d != 0 {
			t.Errorf("%s moved: max deviation %d, %v", name, d, err)
		}
	}
}

// TestPOCVolumes adds progression order changes that split the packet
// order of multi-tile codestreams into several volumes without changing
// it; later volumes revisit packets of earlier ones. The progression order
// of COD is changed, so the result is right only if the POC is followed.
func TestPOCVolumes(t *testing.T) {
	type vol struct{ rs, cs, lye, re, ce, prog int }
	tests := []struct {
		name string
		vols []vol
	}{
		{"tiles.j2k", []vol{{0, 0, 3, 2, 3, progRPCL}, {2, 0, 3, 5, 3, progRPCL}}},
		{"tileparts.j2k", []vol{{0, 0, 1, 5, 3, progLRCP}, {0, 0, 2, 3, 3, progLRCP}, {0, 0, 3, 5, 3, progLRCP}}},
		{"bypass.j2k", []vol{{0, 0, 3, 1, 3, progRLCP}, {0, 0, 3, 4, 3, progRLCP}, {0, 0, 3, 33, 255, progRLCP}}},
		{"sub4.j2k", []vol{{0, 0, 2, 4, 1, progCPRL}, {0, 1, 2, 4, 4, progCPRL}}},
	}
	for _, tt := range tests {
		data := readTestdata(t, tt.name)
		want, err := decodePlanes(data)
		if err != nil {
			t.Fatal(err)
		}
		var poc []byte
		for _, v := range tt.vols {
			poc = append(poc, byte(v.rs), byte(v.cs), byte(v.lye>>8), byte(v.lye), byte(v.re), byte(v.ce), byte(v.prog))
		}
		for _, inTile := range []bool{false, true} {
			s, err := splitStream(data)
			if err != nil {
				t.Fatal(err)
			}
			for i, m := range s.main {
				if m.m == mCOD {
					cod := slices.Clone(m.b)
					cod[1] = (cod[1] + 1) % 5
					s.main[i].b = cod
				}
			}
			if inTile {
				// Split the volumes over the tile-part headers.
				for i := range s.parts {
					if s.parts[i].tpsot == 0 {
						s.parts[i].hdr = append(s.parts[i].hdr, rawSeg{mPOC, poc[:7]}, rawSeg{mPOC, poc[7:]})
					}
				}
			} else {
				s.main = append(s.main, rawSeg{mPOC, poc})
			}
			got, err := decodePlanes(s.bytes())
			if err != nil {
				t.Fatalf("%s: %v", tt.name, err)
			}
			if d, err := maxDiff(got, want); err != nil || d != 0 {
				t.Errorf("%s (tile %v): max deviation %d, %v", tt.name, inTile, d, err)
			}
		}
	}
}

func mkBox(typ string, content ...[]byte) []byte {
	c := bytes.Join(content, nil)
	b := binary.BigEndian.AppendUint32(nil, uint32(8+len(c)))
	return append(append(b, typ...), c...)
}

// wrapJP2 makes a JP2 file of a codestream with the given header boxes.
func wrapJP2(cs []byte, header ...[]byte) []byte {
	return bytes.Join([][]byte{
		jp2Signature,
		mkBox("ftyp", []byte("jp2 \x00\x00\x00\x00jp2 ")),
		mkBox("jp2h", header...),
		mkBox("jp2c", cs),
	}, nil)
}

func ihdr(c *codestream) []byte {
	b := binary.BigEndian.AppendUint32(nil, uint32(c.y1-c.y0))
	b = binary.BigEndian.AppendUint32(b, uint32(c.x1-c.x0))
	b = binary.BigEndian.AppendUint16(b, uint16(len(c.comps)))
	return mkBox("ihdr", b, []byte{7, 7, 0, 0})
}

func colrEnum(cs uint32) []byte {
	return mkBox("colr", binary.BigEndian.AppendUint32([]byte{1, 0, 0}, cs))
}

// paletteJP2 wraps gray8.j2k with a palette of n entries whose columns are
// 8-bit, 12-bit and signed 5-bit, and maps them, out of order, together
// with the index component itself.
func paletteJP2(t *testing.T, n int) (file []byte, pal [][3]int32) {
	return paletteMapJP2(t, n, []byte{0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 2, 0, 0, 0, 0})
}

func paletteMapJP2(t *testing.T, n int, cmap []byte) (file []byte, pal [][3]int32) {
	cs := readTestdata(t, "gray8.j2k")
	c, err := parseCodestream(cs)
	if err != nil {
		t.Fatal(err)
	}
	pclr := binary.BigEndian.AppendUint16(nil, uint16(n))
	pclr = append(pclr, 3, 7, 11, 0x84)
	for i := range n {
		e := [3]int32{int32(255 - i*3%256), int32(i * 37 % 4096), int32(i%32 - 16)}
		pal = append(pal, e)
		pclr = append(pclr, byte(e[0]), byte(e[1]>>8), byte(e[1]), byte(e[2]))
	}
	return wrapJP2(cs, ihdr(c), colrEnum(16), mkBox("pclr", pclr), mkBox("cmap", cmap)), pal
}

func TestPalette(t *testing.T) {
	want, err := decodePlanes(readTestdata(t, "gray8.j2k"))
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range []int{256, 100} {
		file, pal := paletteJP2(t, n)
		img, err := Decode(file)
		if err != nil {
			t.Fatal(err)
		}
		if len(img.Components) != 4 || img.ColorSpace != 16 || img.Alpha != -1 {
			t.Fatalf("%d components, colour space %d, alpha %d", len(img.Components), img.ColorSpace, img.Alpha)
		}
		precs := [][2]int{{12, 0}, {8, 0}, {5, 1}, {8, 0}}
		for i, cp := range img.Components {
			if cp.Precision != precs[i][0] || cp.Signed != (precs[i][1] == 1) {
				t.Errorf("channel %d: precision %d signed %v", i, cp.Precision, cp.Signed)
			}
		}
		cols := []int{1, 0, 2}
		for j, v := range want[0].data {
			e := pal[min(int(v), n-1)]
			for i, col := range cols {
				if got := img.Components[i].Data[j]; got != e[col] {
					t.Fatalf("n=%d channel %d sample %d: %d, want %d", n, i, j, got, e[col])
				}
			}
			if got := img.Components[3].Data[j]; got != v {
				t.Fatalf("direct channel sample %d: %d, want %d", j, got, v)
			}
		}
	}
}

// cdefJP2 wraps sub4.j2k with a channel definition that stores the opacity
// first and the colours in reverse order.
func cdefJP2(t *testing.T, premultiplied bool) []byte {
	cs := readTestdata(t, "sub4.j2k")
	c, err := parseCodestream(cs)
	if err != nil {
		t.Fatal(err)
	}
	typ := byte(1)
	if premultiplied {
		typ = 2
	}
	cdef := []byte{0, 4, 0, 0, 0, typ, 0, 0, 0, 1, 0, 0, 0, 3, 0, 2, 0, 0, 0, 2, 0, 3, 0, 0, 0, 1}
	return wrapJP2(cs, ihdr(c), colrEnum(16), mkBox("cdef", cdef))
}

func TestChannelDefinition(t *testing.T) {
	plain, err := Decode(readTestdata(t, "sub4.j2k"))
	if err != nil {
		t.Fatal(err)
	}
	for _, pre := range []bool{false, true} {
		img, err := Decode(cdefJP2(t, pre))
		if err != nil {
			t.Fatal(err)
		}
		if img.Alpha != 3 || img.Premultiplied != pre {
			t.Fatalf("alpha %d premultiplied %v", img.Alpha, img.Premultiplied)
		}
		for i, src := range []int{3, 2, 1, 0} {
			if !slices.Equal(img.Components[i].Data, plain.Components[src].Data) {
				t.Errorf("channel %d is not component %d", i, src)
			}
		}
	}
}

func TestICC(t *testing.T) {
	cs := readTestdata(t, "gray8.j2k")
	c, _ := parseCodestream(cs)
	icc := []byte("not really an ICC profile")
	file := wrapJP2(cs, ihdr(c), mkBox("colr", []byte{2, 0, 0}, icc), colrEnum(17))
	img, err := Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(img.ICC, icc) || img.ColorSpace != 0 {
		t.Errorf("ICC %q, colour space %d", img.ICC, img.ColorSpace)
	}
}

// TestTruncated checks that a codestream cut anywhere in its tile data
// still decodes, to what its data allows.
func TestTruncated(t *testing.T) {
	for _, name := range []string{"rgb8.j2k", "tiles.j2k", "tileparts.j2k", "bypass.j2k"} {
		data := readTestdata(t, name)
		full, err := decodePlanes(data)
		if err != nil {
			t.Fatal(err)
		}
		c, _ := parseCodestream(data)
		first := bytes.Index(data, []byte{0xff, 0x90})
		prev := -1
		for _, frac := range []int{10, 30, 50, 70, 90, 99} {
			n := first + 14 + (len(data)-first-14)*frac/100
			planes, err := decodePlanes(data[:n])
			if err != nil {
				t.Fatalf("%s cut at %d: %v", name, n, err)
			}
			d, err := maxDiff(planes, full)
			if err != nil {
				t.Fatal(err)
			}
			top := 1<<c.comps[0].prec - 1
			if frac == 99 && d > top/2 || prev >= 0 && d > max(prev, top/8)*2 {
				t.Errorf("%s cut at %d%%: max deviation %d", name, frac, d)
			}
			prev = d
		}
	}
}

// TestCorrupt feeds truncated and damaged files to Decode: it must return,
// without panicking, in bounded time.
func TestCorrupt(t *testing.T) {
	names := []string{"rgb8.j2k", "tiles.j2k", "styles.j2k", "bypass.j2k", "sub4.j2k", "poc.j2k", "gray8.jp2", "tiles_97.j2k"}
	var files [][]byte
	for _, n := range names {
		files = append(files, readTestdata(t, n))
	}
	packed, err := packHeaders(files[1], true, 100)
	if err != nil {
		t.Fatal(err)
	}
	pal, _ := paletteJP2(t, 100)
	files = append(files, packed, pal, cdefJP2(t, false))
	seed := uint32(1)
	rnd := func(n int) int {
		seed ^= seed << 13
		seed ^= seed >> 17
		seed ^= seed << 5
		return int(seed % uint32(n))
	}
	check := func(b []byte, what string) {
		start := time.Now()
		_, err := Decode(b)
		if err != nil && strings.Contains(err.Error(), "internal error") {
			t.Fatalf("%s: %v", what, err)
		}
		if d := time.Since(start); d > 2*time.Second {
			t.Errorf("%s: took %v", what, d)
		}
	}
	for fi, data := range files {
		for n := 0; n < len(data); n += 1 + len(data)/150 {
			check(data[:n], fmt.Sprintf("file %d cut at %d", fi, n))
		}
		for range 300 {
			b := slices.Clone(data)
			for range 1 + rnd(4) {
				i := rnd(len(b))
				switch rnd(3) {
				case 0:
					b[i] ^= 1 << rnd(8)
				case 1:
					b[i] = byte(rnd(256))
				default:
					b[i] = 0xff
				}
			}
			check(b, fmt.Sprintf("file %d damaged", fi))
		}
	}
}

// FuzzDecode looks for inputs that make Decode panic or hang:
//
//	go test ./converter/internal/jpx -run '^$' -fuzz FuzzDecode
func FuzzDecode(f *testing.F) {
	for _, fx := range fixtures {
		b, err := os.ReadFile(filepath.Join("testdata", fx.name))
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		_, err := Decode(data)
		if err != nil && strings.Contains(err.Error(), "internal error") {
			t.Fatal(err)
		}
	})
}

func BenchmarkDecode(b *testing.B) {
	for _, name := range []string{"rgb8.j2k", "rgb8_97.j2k", "tiles.j2k"} {
		data, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			b.Fatal(err)
		}
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			for b.Loop() {
				if _, err := Decode(data); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
