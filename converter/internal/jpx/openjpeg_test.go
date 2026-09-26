package jpx

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// This file compares the decoder with OpenJPEG over a matrix of encoder
// options. It needs opj_compress and opj_decompress and runs only when
// JPX_OPENJPEG=1:
//
//	JPX_OPENJPEG=1 go test ./converter/internal/jpx -run OpenJPEG -v

// source is a synthetic test image.
type source struct {
	name   string
	w, h   int
	prec   int
	signed bool
	sub    [][2]int // subsampling per component
}

func (s source) comps() int { return len(s.sub) }

// file writes the image in a format opj_compress reads and returns the
// extra options it needs.
func (s source) file(dir string) (string, []string) {
	planes := synth(s)
	plain := !s.signed
	for _, d := range s.sub {
		plain = plain && d == [2]int{1, 1}
	}
	if plain && (s.comps() == 1 || s.comps() == 3) {
		var b bytes.Buffer
		magic := "P5"
		if s.comps() == 3 {
			magic = "P6"
		}
		fmt.Fprintf(&b, "%s\n%d %d\n%d\n", magic, s.w, s.h, 1<<s.prec-1)
		for i := range s.w * s.h {
			for c := range s.comps() {
				v := planes[c][i]
				if s.prec > 8 {
					b.WriteByte(byte(v >> 8))
				}
				b.WriteByte(byte(v))
			}
		}
		ext := ".pgm"
		if s.comps() == 3 {
			ext = ".ppm"
		}
		path := filepath.Join(dir, s.name+ext)
		must(os.WriteFile(path, b.Bytes(), 0o644))
		return path, nil
	}
	var b bytes.Buffer
	sign := "u"
	if s.signed {
		sign = "s"
	}
	var subs []string
	for c, d := range s.sub {
		subs = append(subs, fmt.Sprintf("%dx%d", d[0], d[1]))
		for _, v := range planes[c] {
			if s.prec > 8 {
				b.WriteByte(byte(v >> 8))
			}
			b.WriteByte(byte(v))
		}
	}
	path := filepath.Join(dir, s.name+".raw")
	must(os.WriteFile(path, b.Bytes(), 0o644))
	return path, []string{"-F", fmt.Sprintf("%d,%d,%d,%d,%s@%s", s.w, s.h, s.comps(), s.prec, sign, strings.Join(subs, ":"))}
}

// synth makes planes with gradients, texture, noise and sharp edges.
func synth(s source) [][]int32 {
	seed := uint32(2463534242)
	rnd := func() int {
		seed ^= seed << 13
		seed ^= seed >> 17
		seed ^= seed << 5
		return int(seed % 1000)
	}
	top := 1<<s.prec - 1
	planes := make([][]int32, s.comps())
	for c, d := range s.sub {
		w, h := ceilDiv(s.w, d[0]), ceilDiv(s.h, d[1])
		p := make([]int32, w*h)
		for y := range h {
			for x := range w {
				gx, gy := x*d[0], y*d[1]
				v := (gx*top/s.w + gy*top/s.h*(c+1)) / (c + 2)
				v += ((gx*7 + gy*3 + c*50) % 97) * top / 400
				if gx > s.w/3 && gx < s.w*2/3 && gy > s.h/4 && gy < s.h*3/4 {
					v = top - v/3
				}
				if (gx+gy)%23 == 0 {
					v = top * c / max(1, s.comps()-1)
				}
				v += (rnd() - 500) * top / 5000
				v = min(max(v, 0), top)
				if s.signed {
					v -= 1 << (s.prec - 1)
				}
				p[y*w+x] = int32(v)
			}
		}
		planes[c] = p
	}
	return planes
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// readPGX reads a component written by opj_decompress.
func readPGX(path string) (w, h, prec int, signed bool, data []int32, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	r := bufio.NewReader(bytes.NewReader(b))
	line, err := r.ReadString('\n')
	if err != nil {
		return
	}
	var endian, sign string
	if _, err = fmt.Sscanf(line, "PG %s %s %d %d %d", &endian, &sign, &prec, &w, &h); err != nil {
		return
	}
	signed = sign == "-"
	body := b[len(line):]
	n := 1
	if prec > 16 {
		n = 4
	} else if prec > 8 {
		n = 2
	}
	if len(body) < w*h*n {
		err = fmt.Errorf("short pgx")
		return
	}
	data = make([]int32, w*h)
	for i := range data {
		var v uint32
		for j := range n {
			if endian == "ML" {
				v = v<<8 | uint32(body[i*n+j])
			} else {
				v |= uint32(body[i*n+j]) << (8 * j)
			}
		}
		x := int32(v)
		if signed && n < 4 && v&(1<<(8*n-1)) != 0 {
			x -= 1 << (8 * n)
		}
		data[i] = x
	}
	return
}

// openJPEGInvK is the 1/K of OpenJPEG's 9/7 synthesis (its two_invK / 2).
const openJPEGInvK = 1.625732422 / 2

// opjCase is a codestream made by opj_compress with args from a source; the
// decodes of the reversible filter must be exact, even when layers are cut
// by a rate.
type opjCase struct {
	src   string
	ext   string
	args  string
	exact bool
}

// opjSources are the synthetic images. sub3 keeps its second component at
// full size: opj_decompress converts three components as sYCC when it is
// subsampled.
var opjSources = []source{
	{name: "rgb8", w: 67, h: 45, prec: 8, sub: [][2]int{{1, 1}, {1, 1}, {1, 1}}},
	{name: "gray8", w: 67, h: 45, prec: 8, sub: [][2]int{{1, 1}}},
	{name: "gray16", w: 37, h: 29, prec: 16, sub: [][2]int{{1, 1}}},
	{name: "gray12", w: 41, h: 33, prec: 12, sub: [][2]int{{1, 1}}},
	{name: "signed12", w: 39, h: 27, prec: 12, signed: true, sub: [][2]int{{1, 1}}},
	{name: "sub3", w: 67, h: 45, prec: 8, sub: [][2]int{{1, 1}, {1, 1}, {2, 2}}},
	{name: "sub4", w: 67, h: 45, prec: 8, sub: [][2]int{{1, 1}, {2, 2}, {2, 1}, {1, 1}}},
	{name: "tiny", w: 5, h: 3, prec: 8, sub: [][2]int{{1, 1}}},
	{name: "rgb8big", w: 301, h: 203, prec: 8, sub: [][2]int{{1, 1}, {1, 1}, {1, 1}}},
}

func opjCases() []opjCase {
	var cs []opjCase
	add := func(src, ext, args string, exact bool) {
		cs = append(cs, opjCase{src, ext, args, exact})
	}
	// OpenJPEG needs tiles of at least 2^(n-1) samples and precincts of at
	// least 2x2 above resolution 0.
	for _, src := range []string{"rgb8", "gray8", "gray16", "gray12", "signed12", "sub3", "sub4"} {
		add(src, ".j2k", "-n 5", true)
		add(src, ".j2k", "-n 5 -I -r 10", false)
	}
	add("tiny", ".j2k", "-n 2", true)
	add("tiny", ".j2k", "-n 2 -I -r 3", false)
	add("tiny", ".j2k", "-n 1", true)
	for _, p := range []string{"LRCP", "RLCP", "RPCL", "PCRL", "CPRL"} {
		add("rgb8", ".j2k", "-p "+p+" -c [16,16],[16,16],[8,8] -r 40,20,1 -n 3", true)
		add("rgb8", ".j2k", "-p "+p+" -t 20,17 -c [32,32],[16,16] -b 8,8 -n 4", true)
		add("sub4", ".j2k", "-p "+p+" -c [16,16],[8,8] -n 3", true)
		add("rgb8", ".j2k", "-p "+p+" -I -r 30,10 -c [16,16] -b 16,16 -n 4", false)
		add("rgb8big", ".j2k", "-p "+p+" -t 64,48 -T 3,5 -d 11,7 -c [32,32],[16,16] -r 30,10,3,1.5 -SOP -EPH", true)
		add("sub3", ".j2k", "-p "+p+" -t 32,24 -d 5,3 -n 4 -c [16,16],[8,8] -r 20,5,1", true)
	}
	for _, m := range []string{"1", "2", "4", "8", "16", "32", "3", "5", "12", "63", "36", "9", "17"} {
		add("rgb8", ".j2k", "-M "+m, true)
		add("gray16", ".j2k", "-n 5 -M "+m+" -r 20,5,1", true)
		add("rgb8", ".j2k", "-M "+m+" -I -r 8", false)
		add("rgb8big", ".j2k", "-M "+m+" -I -r 20,8,3", false)
	}
	for _, o := range []string{
		"-n 1", "-n 2", "-n 6",
		"-b 4,4", "-b 64,16", "-b 16,64", "-b 1024,4", "-b 4,1024", "-b 4,8 -n 3",
		"-c [4,4] -n 3 -r 2", "-c [8,8],[4,4],[2,2] -n 3 -r 3", "-c [64,32],[32,16]", "-c [2,2] -n 1 -r 1.5",
		"-t 16,16 -n 5", "-t 19,13 -n 4", "-t 16,16 -T 5,3 -d 7,5 -n 4", "-t 32,32 -T 7,11 -d 9,13",
		"-d 1,1", "-d 3,5", "-d 17,9 -n 4",
		"-SOP", "-EPH", "-SOP -EPH", "-SOP -EPH -t 30,20 -r 10,5,1 -n 5",
		"-r 20,10,1", "-r 50,30,10,5,2,1 -n 3",
		"-mct 0", "-mct 1",
		"-ROI c=0,U=4", "-ROI c=1,U=10 -r 20,1",
		// OpenJPEG's encoder repeats the packets of overlapping progression
		// volumes, so these do not overlap; tile numbers start at 1.
		"-POC T1=0,0,3,2,3,CPRL/T1=2,0,3,6,3,LRCP -r 10,5,1",
		"-POC T1=0,0,3,6,1,RPCL/T1=0,1,3,6,3,PCRL -r 10,5,1",
		"-POC T1=0,0,3,3,3,RLCP/T1=3,0,3,6,2,CPRL/T1=3,2,3,6,3,RPCL -r 10,5,1",
		"-PLT", "-TLM -t 32,32", "-TP R", "-TP L -r 20,10,1", "-TP C -t 32,32",
	} {
		add("rgb8", ".j2k", o, true)
	}
	for _, o := range []string{
		"-t 16,16 -T 5,3 -d 7,5 -n 4", "-t 32,32 -T 7,11 -d 9,13", "-d 3,5", "-n 6", "-b 64,16", "-c [16,16] -n 4",
		"-r 50,30,10", "-q 30,40,50", "-M 9", "-mct 0", "-SOP -EPH -p RPCL -c [16,16] -n 4",
		"-ROI c=0,U=4", "-n 1", "-POC T1=0,0,3,2,3,CPRL/T1=2,0,3,6,3,LRCP -r 10,5,1",
	} {
		add("rgb8", ".j2k", "-I "+o, false)
	}
	for _, src := range []string{"gray16", "gray12", "signed12", "sub3", "sub4"} {
		add(src, ".j2k", "-I -r 5 -t 16,16 -d 3,1 -n 4", false)
		add(src, ".j2k", "-t 16,16 -d 3,1 -n 3 -p PCRL", true)
	}
	add("rgb8", ".jp2", "", true)
	add("gray8", ".jp2", "-I -r 10", false)
	add("sub4", ".jp2", "-p CPRL", true)
	add("rgb8big", ".j2k", "", true)
	add("rgb8big", ".j2k", "-n 7", true)
	add("rgb8big", ".j2k", "-I -n 7", false)
	add("rgb8big", ".j2k", "-I -r 20 -n 6 -t 128,96 -T 3,2 -d 5,7 -p PCRL -c [32,32]", false)
	add("rgb8big", ".jp2", "-I -q 35", false)
	return cs
}

func TestOpenJPEG(t *testing.T) {
	if os.Getenv("JPX_OPENJPEG") == "" {
		t.Skip("set JPX_OPENJPEG=1 to compare with OpenJPEG")
	}
	dir := t.TempDir()
	srcs := map[string]source{}
	files := map[string][]string{}
	for _, s := range opjSources {
		srcs[s.name] = s
		p, extra := s.file(dir)
		files[s.name] = append([]string{"-i", p}, extra...)
	}
	worst := map[bool]int{}
	worstOpj := 0
	for i, c := range opjCases() {
		name := fmt.Sprintf("%03d_%s_%s", i, c.src, strings.NewReplacer(" ", "", "/", "_", ",", "_", "[", "", "]", "", "=", "").Replace(c.args))
		t.Run(name, func(t *testing.T) {
			out := filepath.Join(dir, name+c.ext)
			args := append(append([]string{}, files[c.src]...), "-o", out)
			args = append(args, strings.Fields(c.args)...)
			if b, err := exec.Command("opj_compress", args...).CombinedOutput(); err != nil {
				t.Fatalf("opj_compress %v: %v\n%s", args, err, b)
			}
			ref := filepath.Join(dir, name+"_ref.pgx")
			if b, err := exec.Command("opj_decompress", "-i", out, "-o", ref).CombinedOutput(); err != nil {
				t.Fatalf("opj_decompress: %v\n%s", err, b)
			}
			data, err := os.ReadFile(out)
			must(err)
			base := strings.TrimSuffix(ref, ".pgx")
			dev, err := compareWithPGX(data, base, srcs[c.src].comps())
			if err != nil {
				t.Fatal(err)
			}
			if c.exact {
				if dev != 0 {
					t.Errorf("max deviation %d", dev)
				}
				return
			}
			// OpenJPEG scales the 9/7 high-pass samples by an approximation
			// of 1/K; with it, the results must agree up to rounding.
			invK97 = openJPEGInvK
			devOpj, err := compareWithPGX(data, base, srcs[c.src].comps())
			invK97 = 1 / k97
			if err != nil {
				t.Fatal(err)
			}
			if devOpj > 1 || srcs[c.src].prec <= 8 && dev > 1 {
				t.Errorf("max deviation %d, %d with OpenJPEG's 1/K", dev, devOpj)
			}
			worst[srcs[c.src].prec > 8] = max(worst[srcs[c.src].prec > 8], dev)
			worstOpj = max(worstOpj, devOpj)
		})
	}
	t.Logf("max deviation of lossy cases: %d (8 bits), %d (more bits); %d with OpenJPEG's 1/K", worst[false], worst[true], worstOpj)
}

// compareWithPGX decodes data and returns the largest difference to the
// component files base_<i>.pgx.
func compareWithPGX(data []byte, base string, n int) (int, error) {
	cs := data
	if bytes.HasPrefix(data, jp2Signature) {
		var err error
		if _, cs, err = parseJP2(data); err != nil {
			return 0, err
		}
	}
	c, err := parseCodestream(cs)
	if err != nil {
		return 0, err
	}
	planes, err := c.decode()
	if err != nil {
		return 0, err
	}
	if len(planes) != n {
		return 0, fmt.Errorf("%d components, want %d", len(planes), n)
	}
	dev := 0
	for i, pl := range planes {
		w, h, prec, signed, ref, err := readPGX(fmt.Sprintf("%s_%d.pgx", base, i))
		if err != nil {
			return 0, err
		}
		if w != pl.w || h != pl.h || prec != c.comps[i].prec || signed != c.comps[i].signed {
			return 0, fmt.Errorf("component %d: %dx%d prec %d signed %v, want %dx%d prec %d signed %v",
				i, pl.w, pl.h, c.comps[i].prec, c.comps[i].signed, w, h, prec, signed)
		}
		for j, v := range ref {
			d := int(v - pl.data[j])
			if d < 0 {
				d = -d
			}
			if d > dev {
				dev = d
			}
		}
	}
	return dev, nil
}

// opjDecode decodes a file with opj_decompress and returns its components.
func opjDecode(t *testing.T, dir, name string, data []byte) []plane {
	t.Helper()
	in := filepath.Join(dir, name)
	must(os.WriteFile(in, data, 0o644))
	ref := filepath.Join(dir, strings.TrimSuffix(name, filepath.Ext(name))+"_ref.pgx")
	if b, err := exec.Command("opj_decompress", "-i", in, "-o", ref).CombinedOutput(); err != nil {
		t.Fatalf("opj_decompress %s: %v\n%s", name, err, b)
	}
	var planes []plane
	for i := 0; ; i++ {
		w, h, _, _, d, err := readPGX(fmt.Sprintf("%s_%d.pgx", strings.TrimSuffix(ref, ".pgx"), i))
		if err != nil {
			break
		}
		planes = append(planes, plane{w: w, h: h, data: d})
	}
	return planes
}

// TestOpenJPEGSyntax checks syntax that opj_compress does not write, made
// by rewriting its output, against opj_decompress.
func TestOpenJPEGSyntax(t *testing.T) {
	if os.Getenv("JPX_OPENJPEG") == "" {
		t.Skip("set JPX_OPENJPEG=1 to compare with OpenJPEG")
	}
	dir := t.TempDir()
	for _, f := range fixtures {
		if !strings.HasSuffix(f.name, ".j2k") {
			continue
		}
		data := readTestdata(t, f.name)
		for _, ppm := range []bool{false, true} {
			for _, size := range []int{65000, 29} {
				packed, err := packHeaders(data, ppm, size)
				if err != nil {
					t.Fatal(err)
				}
				name := fmt.Sprintf("%s_%v_%d.j2k", strings.TrimSuffix(f.name, ".j2k"), ppm, size)
				want := opjDecode(t, dir, name, packed)
				invK97 = openJPEGInvK
				got, err := decodePlanes(packed)
				invK97 = 1 / k97
				if err != nil {
					t.Fatal(err)
				}
				if d, err := maxDiff(got, want); err != nil || d > 1 || !f.lossy && d != 0 {
					t.Errorf("%s: max deviation %d, %v", name, d, err)
				}
			}
		}
	}
	// Palettes and channel definitions.
	src := source{name: "rgba8", w: 67, h: 45, prec: 8, sub: [][2]int{{1, 1}, {1, 1}, {1, 1}, {1, 1}}}
	p, extra := src.file(dir)
	rgba := filepath.Join(dir, "rgba8.j2k")
	if b, err := exec.Command("opj_compress", append(append([]string{"-i", p}, extra...), "-o", rgba)...).CombinedOutput(); err != nil {
		t.Fatalf("opj_compress: %v\n%s", err, b)
	}
	rgbaCS, err := os.ReadFile(rgba)
	must(err)
	c, err := parseCodestream(rgbaCS)
	must(err)
	cdef := []byte{0, 4, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 3, 0, 2, 0, 0, 0, 2, 0, 3, 0, 0, 0, 1}
	// OpenJPEG maps palette columns only in order.
	identity := []byte{0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 2}
	pal256, _ := paletteMapJP2(t, 256, identity)
	pal100, _ := paletteMapJP2(t, 100, identity)
	for name, file := range map[string][]byte{
		"pal256.jp2": pal256,
		"pal100.jp2": pal100,
		"cdef.jp2":   wrapJP2(rgbaCS, ihdr(c), colrEnum(16), mkBox("cdef", cdef)),
	} {
		want := opjDecode(t, dir, name, file)
		img, err := Decode(file)
		if err != nil {
			t.Fatal(err)
		}
		if len(img.Components) != len(want) {
			t.Fatalf("%s: %d channels, want %d", name, len(img.Components), len(want))
		}
		for i, cp := range img.Components {
			if cp.Signed {
				continue // OpenJPEG does not sign-extend palette entries
			}
			if !slices.Equal(cp.Data, want[i].data) {
				t.Errorf("%s: channel %d differs", name, i)
			}
		}
	}
}

// TestOpenJPEGMountain decodes the JPX test image of pdfcpu, from the module
// cache, and compares its codestream with opj_decompress.
func TestOpenJPEGMountain(t *testing.T) {
	if os.Getenv("JPX_OPENJPEG") == "" {
		t.Skip("set JPX_OPENJPEG=1 to compare with OpenJPEG")
	}
	out, err := exec.Command("go", "env", "GOMODCACHE").Output()
	must(err)
	path := filepath.Join(strings.TrimSpace(string(out)), "github.com/pdfcpu/pdfcpu@v0.11.0/pkg/testdata/resources/mountain.jpx")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skip(err)
	}
	start := time.Now()
	img, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%dx%d, %d channels, colour space %d, ICC %d bytes, alpha %d, decoded in %v",
		img.Width, img.Height, len(img.Components), img.ColorSpace, len(img.ICC), img.Alpha, time.Since(start))
	_, cs, err := parseJP2(data)
	must(err)
	dir := t.TempDir()
	want := opjDecode(t, dir, "mountain.j2k", cs)
	got, err := decodePlanes(cs)
	must(err)
	d, err := maxDiff(got, want)
	must(err)
	invK97 = openJPEGInvK
	got, err = decodePlanes(cs)
	invK97 = 1 / k97
	must(err)
	dOpj, err := maxDiff(got, want)
	must(err)
	t.Logf("max deviation from OpenJPEG %d, %d with its 1/K", d, dOpj)
	if d > 1 || dOpj > 1 {
		t.Errorf("max deviation from OpenJPEG %d, %d with its 1/K", d, dOpj)
	}
}
