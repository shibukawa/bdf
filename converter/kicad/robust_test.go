package kicad

import (
	"archive/zip"
	"bytes"
	"fmt"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/converter"
)

// within fails the test when f does not return in time.
func within(t *testing.T, d time.Duration, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { f(); close(done) }()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("%s did not finish in %v", what, d)
	}
}

// convertQuietly converts an input, and fails the test on an internal
// error (a panic recovered while drawing a page).
func convertQuietly(t *testing.T, data []byte, opts *Options) (*Result, error) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	var warnings []string
	opts.Warn = func(m string) { warnings = append(warnings, m) }
	opts.NoTextIndex = true
	res, err := Convert(bytes.NewReader(data), int64(len(data)), opts)
	for _, w := range warnings {
		if strings.Contains(w, "internal error") {
			t.Errorf("%s", w)
		}
	}
	if res != nil {
		res.Warnings = warnings
	}
	return res, err
}

// TestTruncated converts the test files cut short at many places: each
// converts or fails, in time.
func TestTruncated(t *testing.T) {
	for _, name := range []string{"demo/demo.kicad_sch", "demo/sub.kicad_sch", "demo/demo.kicad_pcb", "frame/frame.kicad_wks", "demo.zip"} {
		// a drawing sheet is read through its project
		data := readTestdata(t, name)
		for i := 1; i <= 64; i++ {
			n := len(data) * i / 64
			within(t, 10*time.Second, fmt.Sprintf("%s cut at %d", name, n), func() {
				if strings.HasSuffix(name, ".kicad_wks") {
					wks := data[:n]
					files := converter.FileMap{"frame.kicad_sch": readTestdata(t, "frame/frame.kicad_sch"), "frame.kicad_wks": wks}
					o := &converter.Options{Files: files}
					opts := testOptions()
					opts.FileName, opts.Refs = "frame.kicad_pro", o.ReadRef
					convertQuietly(t, readTestdata(t, "frame/frame.kicad_pro"), opts)
					return
				}
				convertQuietly(t, data[:n], nil)
			})
		}
	}
}

// TestMutated converts the test files with numbers replaced by values no
// editor writes, a different set of numbers each time: each converts or
// fails, in time.
func TestMutated(t *testing.T) {
	number := regexp.MustCompile(`(?:^|[ (])(-?[0-9]+(?:\.[0-9]+)?)(?:[ )]|$)`)
	for _, name := range []string{"demo/demo.kicad_sch", "demo/sub.kicad_sch", "demo/demo.kicad_pcb"} {
		data := string(readTestdata(t, name))
		for _, v := range []string{"1e300", "-1e300", "0", "-1", "1e-300", "359.99"} {
			for k := range 5 {
				i := 0
				mutated := number.ReplaceAllStringFunc(data, func(m string) string {
					i++
					if i%5 != k {
						return m
					}
					sub := number.FindStringSubmatch(m)
					return strings.Replace(m, sub[1], v, 1)
				})
				within(t, 10*time.Second, fmt.Sprintf("%s with every fifth number (from %d) %s", name, k, v), func() {
					opts := testOptions()
					opts.FileName = path.Base(name)
					convertQuietly(t, []byte(mutated), opts)
				})
			}
		}
	}
}

// TestSelfSheet converts a schematic whose sheets are drawn from its own
// file: KiCad refuses it; here the sheets are left out.
func TestSelfSheet(t *testing.T) {
	sheet := func(i int) string {
		return fmt.Sprintf(`(sheet (at %d 10) (size 20 20) (uuid "00000000-0000-0000-0000-00000000010%d")
  (property "Sheetname" "S%d" (at 0 0 0) (effects (font (size 1.27 1.27))))
  (property "Sheetfile" "self.kicad_sch" (at 0 0 0) (effects (font (size 1.27 1.27)))))`, 10+30*i, i, i)
	}
	src := []byte(schematic(sheet(0), sheet(1)))
	files := converter.FileMap{"self.kicad_sch": src}
	o := &converter.Options{Files: files}
	opts := testOptions()
	opts.FileName, opts.Refs = "self.kicad_sch", o.ReadRef
	var res *Result
	var err error
	within(t, 10*time.Second, "a sheet that contains itself", func() { res, err = convertQuietly(t, src, opts) })
	if err != nil {
		t.Fatal(err)
	}
	if res.Sheets != 1 || len(res.Warnings) != 2 || !strings.Contains(res.Warnings[0], "contains itself") {
		t.Errorf("%d sheets, warnings %q", res.Sheets, res.Warnings)
	}
}

// TestHostileValues converts files of values no editor writes: nothing
// hangs or fails.
func TestHostileValues(t *testing.T) {
	big := "1e300"
	sch := schematic(
		`(lib_symbols (symbol "t:X" (symbol "X_1_1"
  (arc (start 0 0) (mid 0 0) (end 0 0) (stroke (width -1) (type dash)) (fill (type hatch)))
  (arc (start 1e300 0) (mid 0 1e300) (end -1e300 0) (stroke (width 1e300) (type dash_dot)) (fill (type cross_hatch)))
  (bezier (pts (xy 0 0) (xy nan inf) (xy 1 1) (xy 2 2)) (stroke (width 0) (type default)) (fill (type none)))
  (polyline (pts) (stroke (width 0) (type default)) (fill (type background)))
  (text "${VALUE}${REFERENCE}" (at 0 0 1e300) (effects (font (size 1e300 -1e300) (thickness 1e300))))
  (pin input clock (at 0 0 45) (length -1e300) (name "~{" (effects (font (size 0 0)))) (number "" (effects (font (size 1e300 1e300)))))
  (pin passive line (at 0 0 90) (length 2.54) (name "A\nB\n\n" (effects (font (size 1.27 1.27)))) (number "[1\n2\n3]" (effects (font (size 1.27 1.27))))))))`,
		`(symbol (lib_id "t:X") (at `+big+` 0 45) (mirror x) (unit 99) (uuid "00000000-0000-0000-0000-000000000003")
  (property "Reference" "U1" (at 0 0 1e300) (effects (font (size 1.27 1.27)) (justify right top)))
  (property "Value" "${Value}${Value}" (at 0 0 0) (effects (font (size 1.27 1.27)))))`,
		`(symbol (lib_id "missing:Y") (at 10 10 0) (uuid "00000000-0000-0000-0000-000000000004"))`,
		`(rectangle (start -1e300 -1e300) (end 1e300 1e300) (stroke (width 0.001) (type dot)) (fill (type reverse_hatch)))`,
		`(circle (center 0 0) (radius -5) (stroke (width 0) (type default)) (fill (type hatch)))`,
		`(wire (pts (xy 0 0)) (stroke (width 0) (type default)) (uuid "00000000-0000-0000-0000-000000000005"))`,
		`(label "`+strings.Repeat("~{_{^{", 300)+`x" (at 0 0 1e300) (effects (font (size 1.27 1.27))) (uuid "00000000-0000-0000-0000-000000000006"))`,
		`(global_label "G" (shape unknown) (at 0 0 7) (effects (font (size 0 0))) (uuid "00000000-0000-0000-0000-000000000007"))`,
		`(text_box "`+strings.Repeat("word ", 5000)+`" (at 0 0 0) (size 1e-9 1e300) (margins -1 -1 -1 -1) (effects (font (size 1.27 1.27))) (uuid "00000000-0000-0000-0000-000000000008"))`,
		`(table (column_count 1e9) (column_widths) (row_heights) (cells (table_cell "a" (at 0 0 0) (size -1 -1) (effects (font (size 1.27 1.27))))))`,
		`(image (at 0 0) (scale 1e300) (data "iVBORw0KGgo="))`,
		`(sheet (at 0 0) (size 1e300 -1) (uuid "00000000-0000-0000-0000-000000000009") (property "Sheetname" "S" (at 0 0 0) (effects (font (size 1.27 1.27)))) (property "Sheetfile" "" (at 0 0 0) (effects (font (size 1.27 1.27)))))`,
	)
	within(t, 10*time.Second, "a schematic", func() {
		if _, err := convertQuietly(t, []byte(sch), nil); err != nil {
			t.Error(err)
		}
	})
	pcb := `(kicad_pcb (version 20240108) (generator "pcbnew")
  (general (thickness 1e300))
  (layers (0 "F.Cu" signal) (31 "B.Cu" signal) (44 "Edge.Cuts" user) (99 "Weird" user))
  (gr_arc (start 0 0) (end 10 0) (angle 1e300) (layer "Edge.Cuts") (width 0.1))
  (gr_arc (start 0 0) (mid 0 0) (end 0 0) (layer "F.Cu") (width 0.1))
  (gr_rect (start -1e300 -1e300) (end 1e300 1e300) (layer "Edge.Cuts") (width 0))
  (gr_poly (pts (arc (start 0 0) (mid 1 1) (end 2 0)) (xy nan 0)) (layer "F.Cu") (width 0) (fill solid))
  (gr_text "` + strings.Repeat("x", 10000) + `" (at 0 0 1e300) (layer "F.Cu") (effects (font (size 1e300 1e300) (thickness 1e300))))
  (segment (start 0 0) (end 0 0) (width -1) (layer "*.Cu"))
  (via (at 0 0) (size 0) (drill 1e300) (layers "F.Cu" "B.Cu"))
  (zone (net 0) (layer "F.Cu") (polygon (pts)) (filled_polygon (layer "F.Cu") (pts (xy 0 0))))
  (dimension (type aligned) (layer "F.Cu") (pts (xy 0 0) (xy 0 0)) (height 1e300) (gr_text "d" (at 0 0 0) (layer "F.Cu") (effects (font (size 1 1)))))
  (footprint "X" (layer "F.Cu") (at 0 0 1e300)
    (pad "1" smd roundrect (at 0 0 45) (size -1 1e300) (layers "*.Cu" "*.Mask") (roundrect_rratio 99) (chamfer_ratio 99) (chamfer top_left))
    (pad "2" thru_hole oval (at 0 0) (size 1 1) (drill oval 1e300 0) (layers "*.Cu"))
    (pad "3" smd custom (at 0 0) (size 1 1) (layers "F.Cu") (primitives (gr_poly (pts) (width 0)) (gr_arc (start 0 0) (end 0 0) (angle 1e300) (width 1))))
    (pad "4" smd trapezoid (at 0 0) (size 1 1) (rect_delta 1e300 -1e300) (layers "F.Cu"))
    (fp_text reference "R1" (at 0 0) (layer "F.SilkS") (effects (font (size 1 1) (thickness 0.1))))))`
	within(t, 10*time.Second, "a board", func() {
		if _, err := convertQuietly(t, []byte(pcb), nil); err != nil {
			t.Error(err)
		}
	})
	// a drawing sheet that repeats much, drawn on many pages
	wks := `(kicad_wks (setup (textsize 1 1))` + strings.Repeat(`(line (start 0 0) (end 1 1) (repeat 1000000000) (incrx 1e-9))
  (tbtext "Z${#}" (pos 1 1) (repeat 1000000000) (incrlabel 1000000000) (incry -1e300))
  (polygon (pos 1 1) (rotate 1e300) (repeat 1000) (pts (xy 0 0) (xy 1 1) (xy 1 0)))`, 20) + `)`
	var sheets []string
	for i := range 60 {
		sheets = append(sheets, fmt.Sprintf(`(sheet (at 0 0) (size 10 10) (uuid "00000000-0000-0000-0000-0000000002%02d")
  (property "Sheetname" "S%d" (at 0 0 0) (effects (font (size 1.27 1.27))))
  (property "Sheetfile" "sub.kicad_sch" (at 0 0 0) (effects (font (size 1.27 1.27)))))`, i, i))
	}
	root := schematic(sheets...)
	sub := schematic(`(sheet (at 0 0) (size 10 10) (uuid "00000000-0000-0000-0000-000000000300")
  (property "Sheetname" "Leaf" (at 0 0 0) (effects (font (size 1.27 1.27))))
  (property "Sheetfile" "leaf.kicad_sch" (at 0 0 0) (effects (font (size 1.27 1.27)))))`)
	files := converter.FileMap{
		"root.kicad_sch": []byte(root), "sub.kicad_sch": []byte(sub), "leaf.kicad_sch": []byte(schematic()),
		"root.kicad_wks": []byte(wks),
		"root.kicad_pro": []byte(`{"schematic": {"page_layout_descr_file": "root.kicad_wks", "drawing": {"junction_size_choice": 99, "text_offset_ratio": 1e300, "pin_symbol_size": -5, "dashed_lines_dash_length_ratio": -1}}}`),
	}
	var zb bytes.Buffer
	zw := zip.NewWriter(&zb)
	for name, b := range files {
		w, _ := zw.Create("p/" + name)
		w.Write(b)
	}
	zw.Close()
	within(t, 30*time.Second, "121 pages of a drawing sheet that repeats", func() {
		res, err := convertQuietly(t, zb.Bytes(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if res.Sheets != 121 {
			t.Errorf("%d sheets", res.Sheets)
		}
	})
}
