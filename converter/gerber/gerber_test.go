package gerber

import (
	"archive/zip"
	"bytes"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
)

// parse reads a Gerber file and fails the test on warnings unless allowed.
func parse(t *testing.T, src string, allowWarnings bool) *layer {
	t.Helper()
	var warnings []string
	l, err := parseGerber([]byte(src), nil, func(f string, a ...any) { warnings = append(warnings, fmt.Sprintf(f, a...)) })
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) > 0 && !allowWarnings {
		t.Errorf("warnings: %q", warnings)
	}
	return l
}

func near(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func boxIs(t *testing.T, what string, b box, x0, y0, x1, y1 float64) {
	t.Helper()
	if !b.ok || !near(b.Min.X, x0) || !near(b.Min.Y, y0) || !near(b.Max.X, x1) || !near(b.Max.Y, y1) {
		t.Errorf("%s: bounds %v, want (%g %g) (%g %g)", what, b, x0, y0, x1, y1)
	}
}

const header = "%FSLAX46Y46*%\n%MOMM*%\n"

func TestCoordinates(t *testing.T) {
	for _, c := range []struct {
		name, src      string
		x0, y0, x1, y1 float64
	}{
		{"mm 4.6", header + "%ADD10C,1*%\nD10*\nX1000000Y-2000000D03*\nM02*", 0.5, -2.5, 1.5, -1.5},
		{"inch 2.4", "%FSLAX24Y24*%\n%MOIN*%\n%ADD10C,0.1*%\nD10*\nX10000Y5000D03*\nM02*", 25.4 - 1.27, 12.7 - 1.27, 25.4 + 1.27, 12.7 + 1.27},
		// trailing zeros left out: X15 in 2.4 is 15.0000... padded to 150000 → 15
		{"trailing zeros", "%FSTAX24Y24*%\n%MOMM*%\n%ADD10C,1*%\nD10*\nX15Y01D03*\nM02*", 14.5, 0.5, 15.5, 1.5},
		{"decimal point", header + "%ADD10C,1*%\nD10*\nX1.5Y2.5D03*\nM02*", 1, 2, 2, 3},
		// G70/G71 and modal coordinates: Y kept from the point before
		{"G71 and modal Y", "%FSLAX33Y33*%\nG71*\n%ADD10C,1*%\nD10*\nX1000Y1000D03*\nX3000D03*\nM02*", 0.5, 0.5, 3.5, 1.5},
		{"incremental", "%FSLIX33Y33*%\n%MOMM*%\n%ADD10C,1*%\nD10*\nX1000Y1000D03*\nX1000D03*\nX1000D03*\nM02*", 0.5, 0.5, 3.5, 1.5},
		// an operation code carried over from the block before (deprecated)
		{"modal D01", header + "%ADD10C,1*%\nD10*\nX0Y0D02*\nX1000000Y0D01*\nX2000000Y0*\nM02*", -0.5, -0.5, 2.5, 0.5},
		{"no MO", "%FSLAX24Y24*%\n%ADD10C,0.1*%\nD10*\nX10000Y0D03*\nM02*", 25.4 - 1.27, -1.27, 25.4 + 1.27, 1.27},
	} {
		t.Run(c.name, func(t *testing.T) {
			l := parse(t, c.src, c.name == "no MO")
			boxIs(t, c.name, l.img.bounds, c.x0, c.y0, c.x1, c.y1)
		})
	}
}

func TestArcs(t *testing.T) {
	sweepOf := func(src string) (float64, vec) {
		l := parse(t, header+"%ADD10C,0.1*%\nD10*\n"+src+"M02*", false)
		st := l.img.runs[0].strokes[0]
		s := st.chains[0].segs[0]
		return s.sweep, s.c
	}
	for _, c := range []struct {
		name, src string
		sweep     float64
		c         vec
	}{
		{"G75 counterclockwise half", "G75*\nX0Y0D02*\nG03X2000000Y0I1000000J0D01*\n", math.Pi, vec{1, 0}},
		{"G75 clockwise half", "G75*\nX0Y0D02*\nG02X2000000Y0I1000000J0D01*\n", -math.Pi, vec{1, 0}},
		{"G75 full circle", "G75*\nX0Y0D02*\nG02X0Y0I1000000J0D01*\n", -2 * math.Pi, vec{1, 0}},
		{"G75 three quarters", "G75*\nX1000000Y0D02*\nG03X0Y-1000000I-1000000J0D01*\n", 3 * math.Pi / 2, vec{}},
		// single quadrant: the signs of I and J come from the arc of 90° or less
		{"G74 quadrant", "G74*\nX0Y0D02*\nG02X1000000Y1000000I1000000J0D01*\n", -math.Pi / 2, vec{1, 0}},
		{"G74 quadrant down", "G74*\nX0Y0D02*\nG03X1000000Y-1000000I1000000J0D01*\n", math.Pi / 2, vec{1, 0}},
		{"G74 other center", "G74*\nX1000000Y0D02*\nG03X0Y1000000I1000000J0D01*\n", math.Pi / 2, vec{}},
	} {
		sw, cc := sweepOf(c.src)
		if !near(sw, c.sweep) || !near(cc.X, c.c.X) || !near(cc.Y, c.c.Y) {
			t.Errorf("%s: sweep %g about %v, want %g about %v", c.name, sw, cc, c.sweep, c.c)
		}
	}
	// the bounds of a half circle
	l := parse(t, header+"%ADD10C,0.2*%\nD10*\nG75*\nX0Y0D02*\nG02X2000000Y0I1000000J0D01*\nM02*", false)
	boxIs(t, "half circle", l.img.bounds, -0.1, -0.1, 2.1, 1.1)
}

func TestApertures(t *testing.T) {
	for _, c := range []struct {
		ad             string
		x0, y0, x1, y1 float64
		contours       int
	}{
		{"C,2", -1, -1, 1, 1, 1},
		{"C,2X1", -1, -1, 1, 1, 2},
		{"R,2X1", -1, -0.5, 1, 0.5, 1},
		{"R,2X1X0.5", -1, -0.5, 1, 0.5, 2},
		{"R,2X1X0.5X0.2", -1, -0.5, 1, 0.5, 2}, // a rectangular hole, of older files
		{"O,3X1", -1.5, -0.5, 1.5, 0.5, 1},
		{"O,1X3X0.4", -0.5, -1.5, 0.5, 1.5, 2},
		{"P,2X4", -1, -1, 1, 1, 1},
		{"P,2X4X45", -math.Sqrt2 / 2, -math.Sqrt2 / 2, math.Sqrt2 / 2, math.Sqrt2 / 2, 1},
		{"P,2X3", -0.5, -math.Sqrt(3) / 2, 1, math.Sqrt(3) / 2, 1},
	} {
		l := parse(t, header+"%ADD10"+c.ad+"*%\nD10*\nX0Y0D03*\nM02*", false)
		f := l.img.runs[0].flashes[0]
		boxIs(t, c.ad, f.s.bounds, c.x0, c.y0, c.x1, c.y1)
		n := 0
		for _, p := range f.s.prims {
			n += len(p.contours)
			for i, co := range p.contours {
				if (co.area() > 0) != (i == 0) {
					t.Errorf("%s: contour %d has area %g", c.ad, i, co.area())
				}
			}
		}
		if n != c.contours {
			t.Errorf("%s: %d contours, want %d", c.ad, n, c.contours)
		}
	}
	// an inch aperture is in mm
	l := parse(t, "%FSLAX24Y24*%\n%MOIN*%\n%ADD10R,0.1X0.05*%\nD10*\nX0Y0D03*\nM02*", false)
	boxIs(t, "inch rectangle", l.img.runs[0].flashes[0].s.bounds, -1.27, -0.635, 1.27, 0.635)
}

func TestMacros(t *testing.T) {
	vars := map[int]float64{1: 2, 2: 0.5}
	for _, c := range []struct {
		expr string
		want float64
	}{
		{"1+2x3", 7}, {"(1+2)x3", 9}, {"$1X$2", 1}, {"-$1+4", 2}, {"$1/4", 0.5}, {"2x-1", -2},
		{"10-2-3", 5}, {"8/2/2", 2}, {".5+1.", 1.5}, {"$9", 0}, {"+3", 3}, {"1/0", 0},
	} {
		got, err := eval(c.expr, vars)
		if err != nil || !near(got, c.want) {
			t.Errorf("eval(%q) = %g, %v, want %g", c.expr, got, err, c.want)
		}
	}
	for _, bad := range []string{"", "1+", "(1", "1)", "a", "$", "1..2"} {
		if _, err := eval(bad, vars); err == nil {
			t.Errorf("eval(%q) accepted", bad)
		}
	}
	for _, c := range []struct {
		name, am, ad   string
		x0, y0, x1, y1 float64
		simple         bool
	}{
		{"circle", "1,1,2,1,0", "", 0, -1, 2, 1, true},
		{"circle turned", "1,1,2,1,0,90", "", -1, 0, 1, 2, true},
		{"vector line", "20,1,1,0,0,4,0,0", "", 0, -0.5, 4, 0.5, true},
		{"center line turned", "21,1,4,2,0,0,90", "", -1, -2, 1, 2, true},
		{"lower left line", "22,1,2,1,1,1,0", "", 1, 1, 3, 2, true},
		{"outline", "4,1,3,0,0,2,0,2,2,0,0,0", "", 0, 0, 2, 2, true},
		{"outline clockwise", "4,1,3,0,0,2,2,2,0,0,0,0", "", 0, 0, 2, 2, true},
		{"polygon", "5,1,4,0,0,2,0", "", -1, -1, 1, 1, true},
		{"moire", "6,0,0,4,0.5,0.5,2,0.2,5,0", "", -2.5, -2.5, 2.5, 2.5, true},
		// the gaps of the cross cut the circle's extremes
		{"thermal", "7,0,0,4,2,0.5,0", "", -math.Sqrt(4 - 0.0625), -math.Sqrt(4 - 0.0625), math.Sqrt(4 - 0.0625), math.Sqrt(4 - 0.0625), true},
		{"variables", "$3=$1x2*1,1,$3,0,0", ",1", -1, -1, 1, 1, true},
		{"exposure off", "1,1,4,0,0*1,0,2,0,0", "", -2, -2, 2, 2, false},
		{"comments", "0 a comment*0 another, with a comma*1,1,2,0,0", "", -1, -1, 1, 1, true},
	} {
		l := parse(t, header+"%AMM*"+c.am+"*%\n%ADD10M"+c.ad+"*%\nD10*\nX0Y0D03*\nM02*", false)
		s := l.img.runs[0].flashes[0].s
		boxIs(t, c.name, s.bounds, c.x0, c.y0, c.x1, c.y1)
		if s.simple != c.simple {
			t.Errorf("%s: simple %v", c.name, s.simple)
		}
		for _, p := range s.prims {
			if p.on && len(p.contours) > 0 && p.contours[0].area() <= 0 {
				t.Errorf("%s: outer contour clockwise", c.name)
			}
		}
	}
	// the thermal's pieces: four, each within a quadrant
	l := parse(t, header+"%AMT*7,0,0,4,2,0.5,0*%\n%ADD10T*%\nD10*\nX0Y0D03*\nM02*", false)
	cs := l.img.runs[0].flashes[0].s.prims[0].contours
	if len(cs) != 4 {
		t.Fatalf("thermal: %d pieces", len(cs))
	}
	area := 0.0
	for _, c := range cs {
		area += c.area()
	}
	// the ring less the four arms of the cross (0.5 wide, 1 long)
	if want := math.Pi*(4-1) - 4*0.5*1; math.Abs(area-want) > 0.02 {
		t.Errorf("thermal: area %g, want about %g", area, want)
	}
	// KiCad's rounded rectangle
	src := header + "%AMRoundRect*\n0 Rectangle with rounded corners*\n0 $1 Rounding radius*\n4,1,4,$2,$3,$4,$5,$6,$7,$8,$9,$2,$3,0*\n1,1,$1+$1,$2,$3*\n1,1,$1+$1,$4,$5*\n1,1,$1+$1,$6,$7*\n1,1,$1+$1,$8,$9*\n20,1,$1+$1,$2,$3,$4,$5,0*\n20,1,$1+$1,$4,$5,$6,$7,0*\n20,1,$1+$1,$6,$7,$8,$9,0*\n20,1,$1+$1,$8,$9,$2,$3,0*%\n" +
		"%ADD10RoundRect,0.250000X-0.450000X-0.600000X0.450000X-0.600000X0.450000X0.600000X-0.450000X0.600000X0*%\nD10*\nX0Y0D03*\nM02*"
	l = parse(t, src, false)
	boxIs(t, "RoundRect", l.img.runs[0].flashes[0].s.bounds, -0.7, -0.85, 0.7, 0.85)
}

func TestPolarityAndBlocks(t *testing.T) {
	l := parse(t, header+"%ADD10C,1*%\nD10*\nX0Y0D03*\n%LPC*%\nX0Y0D03*\nX1000000Y0D03*\n%LPD*%\nX2000000Y0D03*\nM02*", false)
	var pol []bool
	for _, r := range l.img.runs {
		pol = append(pol, r.clear)
	}
	if !slices.Equal(pol, []bool{false, true, false}) || !l.img.hasClear() {
		t.Errorf("runs %v", pol)
	}
	// the bounds hold only what is dark
	boxIs(t, "polarity", l.img.bounds, -0.5, -0.5, 2.5, 0.5)

	// step and repeat: 3 × 2 copies
	l = parse(t, header+"%ADD10C,1*%\n%SRX3Y2I2.0J3.0*%\nD10*\nX0Y0D03*\n%SR*%\nM02*", false)
	if n := len(l.img.runs[0].flashes); n != 6 {
		t.Errorf("step and repeat: %d flashes", n)
	}
	boxIs(t, "step and repeat", l.img.bounds, -0.5, -0.5, 4.5, 3.5)

	// a block aperture flashed turned, and flashed in clear polarity (its
	// objects' polarities inverted: the dark rectangle joins the clear run
	// before it)
	src := header + "%ADD10R,2X1*%\n%ADD11C,0.5*%\n%ABD20*%\nD10*\nX1000000Y0D03*\n%LPC*%\nD11*\nX1000000Y0D03*\n%LPD*%\n%AB*%\n" +
		"D20*\nX0Y0D03*\n%LR90*%\nX10000000Y0D03*\n%LPC*%\nX20000000Y0D03*\nM02*"
	l = parse(t, src, false)
	var runs []string
	for _, r := range l.img.runs {
		runs = append(runs, fmt.Sprintf("%v:%d", r.clear, len(r.flashes)))
	}
	if want := []string{"false:1", "true:1", "false:1", "true:2", "false:1"}; !slices.Equal(runs, want) {
		t.Errorf("block runs %v, want %v", runs, want)
	}
	// turned 90°: the rectangle at (10, 1), 1 wide and 2 high
	boxIs(t, "turned block", l.img.runs[2].flashes[0].s.bounds.offset(l.img.runs[2].flashes[0].at), 9.5, 0, 10.5, 2)
}

func TestApertureTransformations(t *testing.T) {
	// a rectangle 2 × 1 whose centre is off the flash point: a macro
	flash := func(tr string) box {
		l := parse(t, header+"%AMOFF*21,1,2,1,1,0,0*%\n%ADD10OFF*%\n"+tr+"D10*\nX0Y0D03*\nM02*", false)
		return l.img.bounds
	}
	boxIs(t, "none", flash(""), 0, -0.5, 2, 0.5)
	boxIs(t, "LMX", flash("%LMX*%\n"), -2, -0.5, 0, 0.5)
	boxIs(t, "LR90", flash("%LR90*%\n"), -0.5, 0, 0.5, 2)
	boxIs(t, "LS2", flash("%LS2*%\n"), 0, -1, 4, 1)
	// mirrored, then turned
	boxIs(t, "LMX LR90", flash("%LMX*%\n%LR90*%\n"), -0.5, -2, 0.5, 0)
	// a round aperture scaled draws wider lines
	l := parse(t, header+"%ADD10C,1*%\n%LS0.5*%\nD10*\nX0Y0D02*\nX1000000Y0D01*\nM02*", false)
	if w := l.img.runs[0].strokes[0].width; !near(w, 0.5) {
		t.Errorf("scaled line width %g", w)
	}
}

func TestDraws(t *testing.T) {
	// consecutive draws of one aperture make one chain; a move starts another
	l := parse(t, header+"%ADD10C,0.1*%\n%ADD11C,0.2*%\nD10*\nX0Y0D02*\nX1000000Y0D01*\nX1000000Y1000000D01*\nX3000000Y0D02*\nX4000000Y0D01*\nD11*\nX5000000Y0D01*\nM02*", false)
	r := l.img.runs[0]
	if len(r.strokes) != 2 || len(r.strokes[0].chains) != 2 || len(r.strokes[0].chains[0].segs) != 2 {
		t.Errorf("strokes %+v", r.strokes)
	}
	// a rectangle drawn: its sweep, filled
	l = parse(t, header+"%ADD10R,1X0.5*%\nD10*\nX0Y0D02*\nX2000000Y0D01*\nM02*", false)
	if r := l.img.runs[0]; len(r.fills) != 1 || len(r.strokes) != 0 {
		t.Fatalf("rectangle draw: %d fills, %d strokes", len(r.fills), len(r.strokes))
	}
	boxIs(t, "rectangle draw", l.img.bounds, -0.5, -0.25, 2.5, 0.25)
	// a draw of no length is a dot
	l = parse(t, header+"%ADD10C,1*%\nD10*\nX0Y0D02*\nX0Y0D01*\nM02*", false)
	boxIs(t, "dot", l.img.bounds, -0.5, -0.5, 0.5, 0.5)
}

func TestRegions(t *testing.T) {
	// a clockwise contour is filled counterclockwise; a cut-in hole stays
	// clockwise inside it; a second contour of the region is filled too
	src := header + "G36*\nX0Y0D02*\nG01X0Y4000000D01*\nX4000000Y4000000D01*\nX4000000Y2000000D01*\nX3000000Y2000000D01*\n" +
		"G75*\nG03X3000000Y2000000I-1000000J0D01*\nG01X4000000Y2000000D01*\nX4000000Y0D01*\nX0Y0D01*\n" +
		"X5000000Y0D02*\nX6000000Y0D01*\nX6000000Y1000000D01*\nG37*\nM02*"
	l := parse(t, src, false)
	fills := l.img.runs[0].fills
	if len(fills) != 2 {
		t.Fatalf("%d contours", len(fills))
	}
	if a := fills[0].area(); math.Abs(a-(16-math.Pi)) > 1e-6 {
		t.Errorf("area with the hole %g, want %g", a, 16-math.Pi)
	}
	if a := fills[1].area(); !near(a, 0.5) {
		t.Errorf("the triangle, closed: area %g", a)
	}
}

func TestLegacy(t *testing.T) {
	// %IPNEG*%, several commands in one extended block, G54, an image
	// offset and a comment written as G4
	src := "%FSLAX23Y23*MOIN*IPNEG*OFA0.5B0*%\n%ADD10C,0.100*%\nG4 old comment*\nG54D10*\nX0Y0D03*\nM02*"
	l := parse(t, src, false)
	if !l.negative {
		t.Error("IPNEG not read")
	}
	boxIs(t, "offset", l.img.bounds, 12.7-1.27, -1.27, 12.7+1.27, 1.27)
	// attributes, in KiCad's comments too
	l = parse(t, "G04 #@! TF.FileFunction,Copper,L2,Bot*\n%TF.ProjectId,demo,uuid,rev*%\n"+header+"%ADD10C,1*%\nD10*\nX0Y0D03*\nM02*", false)
	if l.attrs["FileFunction"] != "Copper,L2,Bot" || l.attrs["ProjectId"] != "demo,uuid,rev" {
		t.Errorf("attributes %v", l.attrs)
	}
	// a negative copper layer
	l = parse(t, "%TF.FileFunction,Copper,L2,Inr,Plane*%\n%TF.FilePolarity,Negative*%\n"+header+"%ADD10C,1*%\nD10*\nX0Y0D03*\nM02*", false)
	if !l.negative {
		t.Error("negative copper not read")
	}
}

func TestWarnings(t *testing.T) {
	var warnings []string
	parseGerber([]byte(header+"%ADD10Q,1*%\nD99*\nX0Y0D03*\nG99*\n%XX1*%\nG36*\nX0Y0D02*\nX1000000Y0D01*\nX0Y1000000D03*\nG37*\nM02*"), nil,
		func(f string, a ...any) { warnings = append(warnings, fmt.Sprintf(f, a...)) })
	for _, want := range []string{"aperture D10", "aperture D99 is not defined", "a flash without an aperture", "unknown code G99", "unknown extended command \"XX\"", "flashes in regions"} {
		if !slices.ContainsFunc(warnings, func(w string) bool { return strings.Contains(w, want) }) {
			t.Errorf("no warning %q in %q", want, warnings)
		}
	}
}

func TestExcellon(t *testing.T) {
	for _, c := range []struct {
		name, src      string
		holes          int
		x0, y0, x1, y1 float64
	}{
		{"KiCad decimal", "M48\n; FORMAT={-:-/ absolute / metric / decimal}\nFMAT,2\nMETRIC\nT1C1.000\n%\nG90\nG05\nT1\nX10.0Y-5.0\nX20Y-5\nT0\nM30\n", 2, 9.5, -5.5, 20.5, -4.5},
		// Eagle: inches, 2.4, leading zeros left out
		{"Eagle", "%\nM48\nM72\nT01C0.0400\n%\nT01\nX5000Y1500\nX10000Y1500\nM30\n", 2, 12.7 - 0.508, 3.81 - 0.508, 25.4 + 0.508, 3.81 + 0.508},
		// Altium: FILE_FORMAT 2:5 with leading zeros kept
		{"Altium", "M48\n;Layer_Color=9474304\n;FILE_FORMAT=2:5\nINCH,LZ\n;TYPE=PLATED\nT1F00S00C0.0400\n%\nT01\nX0112Y0175\nM30\n", 1, 28.448 - 0.508, 44.45 - 0.508, 28.448 + 0.508, 44.45 + 0.508},
		// metric with trailing zeros kept, 3.3 and a format in the header
		{"metric TZ", "M48\nMETRIC,TZ,000.000\nT1C0.8\n%\nT1\nX12500Y-3000\nM30\n", 1, 12.1, -3.4, 12.9, -2.6},
		{"metric LZ", "M48\nMETRIC,LZ\nT1C0.8\n%\nT1\nX0125Y-003\nM30\n", 1, 12.1, -3.4, 12.9, -2.6},
		{"incremental and repeat", "M48\nMETRIC\nICI,ON\nT1C1.0\n%\nT1\nX1.0Y1.0\nX1.0\nR3X1.0Y0\nM30\n", 5, 0.5, 0.5, 5.5, 1.5},
		// a tool of the inch header, and coordinates in mm after M71
		{"units in the body", "M48\nINCH\nT1C0.05\n%\nM71\nT1\nX1.0Y1.0\nM30\n", 1, 0.365, 0.365, 1.635, 1.635},
		{"tool defined in the body", "M48\nMETRIC\n%\nT1C2.0\nX0Y0\nM30\n", 1, -1, -1, 1, 1},
	} {
		t.Run(c.name, func(t *testing.T) {
			var warnings []string
			l := parseExcellon([]byte(c.src), nil, func(f string, a ...any) { warnings = append(warnings, fmt.Sprintf(f, a...)) })
			if len(warnings) > 0 {
				t.Errorf("warnings %q", warnings)
			}
			n := 0
			for _, r := range l.img.runs {
				n += len(r.flashes)
			}
			if n != c.holes {
				t.Errorf("%d holes, want %d", n, c.holes)
			}
			boxIs(t, c.name, l.img.bounds, c.x0, c.y0, c.x1, c.y1)
		})
	}
	// slots: G85, and routed with the tool down (KiCad), and an arc
	src := "M48\nMETRIC\nT1C1.0\n%\nT1\nX0Y0G85X0Y2.0\nG00X5.0Y0\nM15\nG01X5.0Y2.0\nG03X7.0Y2.0A1.0\nM16\nG05\nX10.0Y10.0\nM30\n"
	l := parseExcellon([]byte(src), nil, func(string, ...any) {})
	r := l.img.runs[0]
	if len(r.flashes) != 1 || len(r.strokes) != 1 || len(r.strokes[0].chains) != 2 {
		t.Fatalf("slots: %d holes, %+v", len(r.flashes), r.strokes)
	}
	if s := r.strokes[0].chains[1].segs; len(s) != 2 || !s[1].arc || !near(s[1].sweep, math.Pi) || !near(s[1].c.X, 6) {
		t.Errorf("routed arc %+v", s)
	}
	boxIs(t, "slots", l.img.bounds, -0.5, -0.5, 10.5, 10.5)
	if l.plated != "" {
		t.Errorf("plated %q", l.plated)
	}
	l = parseExcellon([]byte("M48\n; #@! TF.FileFunction,NonPlated,1,2,NPTH\n;TYPE=NON_PLATED\nMETRIC\nT1C3.2\n%\nT1\nX0Y0\nM30\n"), nil, func(string, ...any) {})
	if l.plated != "nonplated" || l.attrs["FileFunction"] != "NonPlated,1,2,NPTH" {
		t.Errorf("plating %q %v", l.plated, l.attrs)
	}
}

func TestIdentify(t *testing.T) {
	type want struct {
		k kind
		s side
		i int
	}
	for name, w := range map[string]want{
		// KiCad
		"board-F_Cu.gbr": {kCopper, sTop, 0}, "board-B_Cu.gbr": {kCopper, sBottom, 0}, "board-In2_Cu.gbr": {kCopper, sInner, 2},
		"board-F_Mask.gbr": {kMask, sTop, 0}, "board-B_Mask.gbs": {kMask, sBottom, 0}, "board-F_SilkS.gbr": {kSilk, sTop, 0},
		"board-B_Silkscreen.gbr": {kSilk, sBottom, 0}, "board-F_Paste.gbr": {kPaste, sTop, 0}, "board-Edge_Cuts.gbr": {kOutline, sNone, 0},
		"board-F.Cu.gbr": {kCopper, sTop, 0}, "silk-test-B_Cu.gbr": {kCopper, sBottom, 0},
		// Protel / Altium
		"board.GTL": {kCopper, sTop, 0}, "board.GBL": {kCopper, sBottom, 0}, "board.G2": {kCopper, sInner, 2}, "board.GP1": {kCopper, sInner, 1},
		"board.GTO": {kSilk, sTop, 0}, "board.GBS": {kMask, sBottom, 0}, "board.GTP": {kPaste, sTop, 0}, "board.GKO": {kOutline, sNone, 0},
		"board.GM1": {kOutline, sNone, 0},
		// EasyEDA
		"Gerber_TopLayer.GTL": {kCopper, sTop, 0}, "Gerber_BoardOutlineLayer.GKO": {kOutline, sNone, 0},
		"Gerber_TopSolderMaskLayer.gbr": {kMask, sTop, 0}, "Gerber_BottomPasteMaskLayer.gbr": {kPaste, sBottom, 0},
		"Gerber_TopSilkscreenLayer.gbr": {kSilk, sTop, 0}, "Gerber_InnerLayer1.gbr": {kCopper, sInner, 1},
		// Eagle
		"board.cmp": {kCopper, sTop, 0}, "board.sol": {kCopper, sBottom, 0}, "board.plc": {kSilk, sTop, 0}, "board.stc": {kMask, sTop, 0},
		"board.crs": {kPaste, sBottom, 0}, "board.dim": {kOutline, sNone, 0},
		"copper_top.gbr": {kCopper, sTop, 0}, "copper_bottom.gbr": {kCopper, sBottom, 0}, "soldermask_top.gbr": {kMask, sTop, 0},
		"silkscreen_bottom.gbr": {kSilk, sBottom, 0}, "solderpaste_top.gbr": {kPaste, sTop, 0}, "profile.gbr": {kOutline, sNone, 0},
		"copper_l2.gbr": {kCopper, sInner, 1},
		// OrCAD and words
		"board.smb": {kMask, sBottom, 0}, "outline.gbr": {kOutline, sNone, 0}, "BottomLayer.ger": {kCopper, sBottom, 0},
		// not layers of the board
		"board-drl_map.gbr": {kOther, sNone, 0}, "board-F_Fab.gbr": {kOther, sNone, 0}, "board-User_Drawings.gbr": {kOther, sNone, 0},
		"notes.gbr": {kOther, sNone, 0},
	} {
		f, _ := byName(name, false)
		if f.kind != w.k || f.side != w.s || f.index != w.i {
			t.Errorf("%s: %+v, want %+v", name, f, w)
		}
	}
	for name, plated := range map[string]string{"board-PTH.drl": "plated", "board-NPTH.drl": "nonplated", "board.drl": "", "Drill_NPTH_Through.DRL": "nonplated"} {
		if f, _ := byName(name, true); f.kind != kDrill || f.plated != plated {
			t.Errorf("%s: %+v", name, f)
		}
	}
	for v, w := range map[string]function{
		"Copper,L1,Top":            {kind: kCopper, side: sTop},
		"Copper,L3,Inr,Plane":      {kind: kCopper, side: sInner, index: 2},
		"Copper,L4,Bot,Signal":     {kind: kCopper, side: sBottom, index: 3},
		"Soldermask,Bot":           {kind: kMask, side: sBottom},
		"SolderMask,Top":           {kind: kMask, side: sTop},
		"Legend,Top":               {kind: kSilk, side: sTop},
		"SolderPaste,Bot":          {kind: kPaste, side: sBottom},
		"Profile,NP":               {kind: kOutline},
		"Plated,1,4,PTH":           {kind: kDrill, plated: "plated"},
		"NonPlated,1,2,NPTH,Drill": {kind: kDrill, plated: "nonplated"},
		"Drillmap":                 {kind: kOther, other: "Drillmap"},
	} {
		if f, ok := fileFunction(v); !ok || f != w {
			t.Errorf("%s: %+v, want %+v", v, f, w)
		}
	}
	if w := words("Gerber_TopSilkLayer-In1.Cu"); !slices.Equal(w, []string{"gerber", "top", "silk", "layer", "in", "1", "cu"}) {
		t.Errorf("words %q", w)
	}
}

func TestOtherTitle(t *testing.T) {
	for v, want := range map[string]string{
		"Other,Features": "Features", "OtherDrawing,Title block": "Title block", "Drillmap": "Drill map",
		"AssemblyDrawing,Top": "Top assembly drawing", "Keep-out,Bot": "Bottom keep-out", "Component,L1,Top": "Top components",
		"FabricationDrawing": "Fabrication drawing", "Other": "Other",
	} {
		if got := otherTitle(v); got != want {
			t.Errorf("%s: %q, want %q", v, got, want)
		}
	}
}

func TestOutline(t *testing.T) {
	// a square drawn in pieces, one reversed and one a little short; a
	// round cut-out; a line that closes nothing; a frame around nothing
	// outside the board
	src := header + "%ADD10C,0.1*%\nD10*\nG75*\n" +
		"X0Y0D02*\nX10000000Y0D01*\n" +
		"X10000000Y10000000D02*\nX10000000Y0D01*\n" +
		"X0Y10000000D02*\nX9990000Y10000000D01*\n" +
		"X0Y0D02*\nX0Y10000000D01*\n" +
		"X5000000Y4000000D02*\nG02X5000000Y4000000I0J1000000D01*\n" +
		"X2000000Y2000000D02*\nG01X3000000Y3000000D01*\n" +
		"X20000000Y0D02*\nX21000000Y0D01*\nX21000000Y1000000D01*\nX20000000Y0D01*\nM02*"
	l := parse(t, src, false)
	loops, ok := boardShape([]*layer{l})
	if !ok || len(loops) != 2 {
		t.Fatalf("%d loops, closed %v", len(loops), ok)
	}
	if a := math.Abs(loops[0].area()); math.Abs(a-100) > 0.2 {
		t.Errorf("board area %g", a)
	}
	if a := math.Abs(loops[1].area()); !near(a, math.Pi) {
		t.Errorf("cut-out area %g", a)
	}
	// lines that close nothing
	l = parse(t, header+"%ADD10C,0.1*%\nD10*\nX0Y0D02*\nX10000000Y0D01*\nX10000000Y10000000D01*\nM02*", false)
	if _, ok := boardShape([]*layer{l}); ok {
		t.Error("an open outline closed")
	}
}

func readZip(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestConvertBoard(t *testing.T) {
	data := readZip(t, "board.zip")
	res, err := Convert(bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) > 0 {
		t.Errorf("warnings %q", res.Warnings)
	}
	var views []string
	for _, v := range res.Doc.Views {
		views = append(views, v.ID+"="+v.Title)
	}
	want := []string{"top=Top", "bottom=Bottom", "top-silkscreen=Top silkscreen", "top-paste=Top paste", "top-solder-mask=Top solder mask",
		"top-copper=Top copper", "bottom-copper=Bottom copper", "bottom-solder-mask=Bottom solder mask", "bottom-silkscreen=Bottom silkscreen",
		"outline=Outline", "drill-non-plated=Drill (non-plated)", "drill-plated=Drill (plated)"}
	if !slices.Equal(views, want) {
		t.Errorf("views\n%q\nwant\n%q", views, want)
	}
	if !res.Board || res.Layers != 10 || res.Scale != 6 || math.Abs(res.W-60.1) > 1e-6 {
		t.Errorf("result %+v", res)
	}
	m := res.Doc.Meta
	if m.Source != "gerber" || m.DC.Title.First() != "bdf-demo" || m.DC.Created.First() != "2026-09-27T12:00:00+09:00" {
		t.Errorf("meta %+v", m)
	}
	p := res.Doc.Views[0].Pages[0]
	// the board, 60.1 mm wide with 1.8 mm margins, at 6:1
	if w := float64(p.W) / (72 / 25.4); math.Abs(w-(60.1+2*1.803)*6) > 0.01 {
		t.Errorf("page width %g mm", w)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}

	// views=board, and the colors of the options
	res, err = Convert(bytes.NewReader(data), int64(len(data)), &Options{Views: ViewsBoard, Mask: "red", Finish: "silver", Silkscreen: "#ffff00"})
	if err != nil || len(res.Doc.Views) != 2 {
		t.Fatalf("views=board: %v", err)
	}
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), &Options{Mask: "pink"}); err == nil {
		t.Error("unknown mask color accepted")
	}
	res, err = Convert(bytes.NewReader(data), int64(len(data)), &Options{Views: ViewsLayers, Title: "T"})
	if err != nil || len(res.Doc.Views) != 10 || res.Doc.Meta.DC.Title.First() != "T" {
		t.Fatalf("views=layers: %v", err)
	}
}

func TestConvertSingle(t *testing.T) {
	data := readZip(t, "features.gbr")
	res, err := Convert(bytes.NewReader(data), int64(len(data)), &Options{FileName: "features.gbr"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) > 0 {
		t.Errorf("warnings %q", res.Warnings)
	}
	if v := res.Doc.Views; len(v) != 1 || v[0].ID != "layer" || v[0].Title != "Features" || res.Board {
		t.Errorf("views %+v", v[0])
	}
	// a drill file alone, told by its content
	drill := []byte("M48\nMETRIC\nT1C1.0\n%\nT1\nX0Y0\nX10.0Y0\nM30\n")
	res, err = Convert(bytes.NewReader(drill), int64(len(drill)), &Options{FileName: "x.txt"})
	if err != nil || res.Doc.Views[0].ID != "drill" {
		t.Fatalf("drill: %v", err)
	}
	// a copper layer by its name: its view, in its color
	top := []byte(header + "%ADD10C,1*%\nD10*\nX0Y0D03*\nM02*")
	res, err = Convert(bytes.NewReader(top), int64(len(top)), &Options{FileName: "board.GTL"})
	if err != nil || res.Doc.Views[0].Title != "Top copper" {
		t.Fatalf("copper: %v", err)
	}
	for _, bad := range [][]byte{[]byte("hello"), []byte(header + "M02*")} {
		if _, err := Convert(bytes.NewReader(bad), int64(len(bad)), nil); err == nil {
			t.Errorf("%q converted", bad)
		}
	}
	// an archive without Gerber files
	var zb bytes.Buffer
	zw := zip.NewWriter(&zb)
	w, _ := zw.Create("readme.txt")
	w.Write([]byte("hello"))
	zw.Close()
	if _, err := Convert(bytes.NewReader(zb.Bytes()), int64(zb.Len()), nil); err == nil {
		t.Error("an archive without Gerber files converted")
	}
}

func TestDetect(t *testing.T) {
	for _, c := range []struct {
		name string
		data []byte
		want bool
	}{
		{"gerber", []byte(header + "%ADD10C,1*%\n"), true},
		{"gerber after comments", []byte(strings.Repeat("G04 a long header of comments*\n", 100) + header), true},
		{"excellon", []byte("M48\n;comment\nMETRIC\n"), true},
		{"excellon after a rewind stop", []byte("%\nM48\nM72\n"), true},
		{"board", readZip(t, "board.zip"), true},
		{"text", []byte("G04 is not enough\n"), false},
		{"binary", []byte("%FS\x00\x01"), false},
		{"pdf", []byte("%PDF-1.7\n%FS"), false},
	} {
		if got := Detect(c.data[:min(len(c.data), 1024)], bytes.NewReader(c.data), int64(len(c.data))); got != c.want {
			t.Errorf("%s: %v", c.name, got)
		}
	}
}

// TestHostile checks that malformed files end quickly without panicking.
func TestHostile(t *testing.T) {
	for _, src := range []string{
		header + "%ADD10C,1*%\n%SRX100000Y100000I0.001J0.001*%\nD10*\nX0Y0D03*\n%SR*%\nM02*",
		header + "%ADD10C,1*%\n" + strings.Repeat("%ABD11*%\n", 100) + "D10*\nX0Y0D03*\n" + strings.Repeat("%AB*%\n", 100) + "M02*",
		header + "%AMM*4,1,1000000000,0,0*%\n%ADD10M*%\nD10*\nX0Y0D03*\nM02*",
		header + "%AMM*1,1,((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((1*%\n%ADD10M*%\nM02*",
		header + "%AMM*6,0,0,1e300,1,1,1000000,1,1,0*%\n%ADD10M*%\nD10*\nX0Y0D03*\nM02*",
		header + "%ADD10P,1X1000*%\n%ADD11C,NaN*%\n%ADD12C,1e308*%\nD12*\nX0Y0D03*\nM02*",
		header + "%ADD10C,1*%\nD10*\nX99999999999999999999999Y1D03*\nG75*\nG03X1Y1I0J0D01*\nM02*",
		header + "G36*\nG36*\nX0Y0D01*\nG37*\n%LPC*%\n%ABD10*%\n%SRX2Y2I1J1*%\n",
		"%FSLAX99Y99*%\n%MOMM*%\n%ADD10C,1*%\nD10*\nX1Y1D03*\n",
		"%\n%%%%%%%*****%",
		header + "%ADD10C,1*%\n%ABD20*%\nD10*\nX0Y0D03*\n%AB*%\n%SRX1000Y1000I1J1*%\nD20*\nX0Y0D03*\n%SR*%\n%SRX1000Y1000I1J1*%\nD20*\nX0Y0D03*\n%SR*%\nM02*",
	} {
		start := time.Now()
		for _, name := range []string{"x.gbr", "x.drl"} {
			Convert(bytes.NewReader([]byte(src)), int64(len(src)), &Options{FileName: name})
		}
		if d := time.Since(start); d > 20*time.Second {
			t.Errorf("%.40q took %v", src, d)
		}
	}
}

func FuzzGerber(f *testing.F) {
	f.Add([]byte(header + "%ADD10C,1*%\nD10*\nX0Y0D03*\nM02*"))
	f.Add(readFile(f, "features.gbr"))
	f.Add([]byte("M48\nMETRIC\nT1C1.0\n%\nT1\nX0Y0G85X1Y1\nM30\n"))
	f.Fuzz(func(t *testing.T, data []byte) {
		res, err := Convert(bytes.NewReader(data), int64(len(data)), &Options{FileName: "x.gbr"})
		if err == nil {
			var buf bytes.Buffer
			if err := res.Doc.WriteSingle(&buf); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func readFile(f *testing.F, name string) []byte {
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		f.Fatal(err)
	}
	return b
}

// TestObjects checks what the layer objects hold: clear objects drawn with
// destination-out, flashes as runs of paths, macros that take parts away
// in groups of their own.
func TestObjects(t *testing.T) {
	l := parse(t, header+"%AMD*1,1,2,0,0*1,0,1,0,0*%\n%ADD10C,1*%\n%ADD11D*%\nD10*\nX0Y0D03*\nX2000000Y0D03*\n%LPC*%\nX0Y0D03*\n%LPD*%\nD11*\nX5000000Y0D03*\nM02*", false)
	doc := bdf.NewDocument()
	e := &encoder{doc: doc, m: mapping{k: 1, origin: vec{-10, 10}}}
	lo := e.object(l.img)
	if !lo.clear {
		t.Error("clear not reported")
	}
	part := doc.Part(lo.hash)
	if part == nil {
		t.Fatal("no object")
	}
	obj, err := bdf.DecodeObject(part.Data)
	if err != nil {
		t.Fatal(err)
	}
	ins, err := obj.Instructions()
	if err != nil {
		t.Fatal(err)
	}
	var ops []byte
	for _, in := range ins {
		ops = append(ops, in.Op)
	}
	want := []byte{bdf.OpFillPathRun, bdf.OpSave, bdf.OpBlend, bdf.OpFillPathRun, bdf.OpRestore, bdf.OpGroupBegin, bdf.OpUseAt, bdf.OpGroupEnd}
	if !slices.Equal(ops, want) {
		t.Errorf("ops %v, want %v", ops, want)
	}
}
