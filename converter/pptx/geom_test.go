package pptx

import (
	"math"
	"testing"
)

func TestPresetsEvaluate(t *testing.T) {
	defs := loadPresets()
	if len(defs) < 180 {
		t.Fatalf("only %d presets", len(defs))
	}
	for name, def := range defs {
		g := evalGeometry(def, nil, 200, 100)
		if len(g.paths) == 0 {
			t.Errorf("%s: no paths", name)
			continue
		}
		drawn := false
		for _, p := range g.paths {
			for _, a := range p.path.Args {
				if math.IsNaN(float64(a)) || math.IsInf(float64(a), 0) {
					t.Errorf("%s: bad coordinate %v", name, a)
				}
			}
			if !p.empty {
				drawn = true
			}
		}
		if !drawn {
			t.Errorf("%s: nothing drawn", name)
		}
	}
}

func TestArcEnd(t *testing.T) {
	// A quarter of a circle from 0° to 90° ends at the bottom.
	x, y, a0, a1, ok := arcEnd(100, 50, 50, 50, 0, math.Pi/2)
	if !ok || math.Abs(x-50) > 1e-9 || math.Abs(y-100) > 1e-9 || a0 != 0 || math.Abs(a1-math.Pi/2) > 1e-9 {
		t.Fatalf("arc = %v %v %v %v", x, y, a0, a1)
	}
	// On an ellipse the visual angle 45° hits the point on the line to the
	// corner of the bounding box scaled by the radii.
	_, _, _, a1, _ = arcEnd(200, 50, 100, 50, 0, math.Pi/4)
	want := math.Atan2(100*math.Sin(math.Pi/4), 50*math.Cos(math.Pi/4))
	if math.Abs(a1-want) > 1e-9 {
		t.Fatalf("a1 = %v, want %v", a1, want)
	}
	// Full sweep stays a full turn; negative sweeps go the other way.
	_, _, a0, a1, _ = arcEnd(200, 50, 100, 50, 0.3, 2*math.Pi)
	if math.Abs(a1-a0-2*math.Pi) > 1e-9 {
		t.Fatalf("full sweep = %v", a1-a0)
	}
	_, _, a0, a1, _ = arcEnd(200, 50, 100, 50, math.Pi, -math.Pi/2)
	if a1 >= a0 {
		t.Fatalf("negative sweep %v -> %v", a0, a1)
	}
}

func TestGuideFormulas(t *testing.T) {
	e := newGuideEnv(200, 100)
	for fmla, want := range map[string]float64{
		"*/ w 1 4": 50, "+- w h 50": 250, "+/ w h 3": 100, "?: -1 w h": 100, "abs -3": 3,
		"at2 1 1": 2700000, "cos 10 0": 10, "sin 10 5400000": 10, "max w h": 200, "min w h": 100,
		"mod 3 4 0": 5, "pin 0 150 100": 100, "sqrt 16": 4, "val 7": 7, "cat2 10 1 0": 10, "sat2 10 0 1": 10,
	} {
		if got := e.eval(fmla); math.Abs(got-want) > 1e-6 {
			t.Errorf("%q = %v, want %v", fmla, got, want)
		}
	}
}
