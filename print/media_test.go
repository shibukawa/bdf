package print

import (
	"math"
	"testing"
)

func TestParseMedia(t *testing.T) {
	for name, want := range map[string][2]float64{
		"iso_a4_210x297mm":          {595.276, 841.89},
		"na_letter_8.5x11in":        {612, 792},
		"jis_b5_182x257mm":          {515.906, 728.504},
		"jpn_hagaki_100x148mm":      {283.465, 419.528},
		"custom_my-form_80x297.5mm": {226.772, 843.307},
	} {
		m, err := ParseMedia(name)
		if err != nil || math.Abs(m.Width-want[0]) > 0.01 || math.Abs(m.Height-want[1]) > 0.01 {
			t.Errorf("ParseMedia(%q) = %+v, %v", name, m, err)
		}
	}
	for _, bad := range []string{"a4", "iso_a4", "iso_a4_210x297cm", "x_0x297mm", "x_ax297mm", "x_210x99999mm"} {
		if _, err := ParseMedia(bad); err == nil {
			t.Errorf("ParseMedia(%q) succeeded", bad)
		}
	}
	if m, err := ResolveMedia("B5"); err != nil || m.Name != "jis_b5_182x257mm" {
		t.Errorf("B5: %+v %v", m, err)
	}
	if m, err := ResolveMedia("iso_a3_297x420mm"); err != nil || m.Name != "iso_a3_297x420mm" {
		t.Errorf("a PWG name: %+v %v", m, err)
	}
}

func TestChooseMedia(t *testing.T) {
	ready := []string{"iso_a4_210x297mm"}
	supported := []string{"na_letter_8.5x11in", "iso_a4_210x297mm", "jis_b5_182x257mm", "jpn_hagaki_100x148mm", "roll_max_8.5x39in"}
	for _, c := range []struct {
		w, h float64
		want string
	}{
		{595, 842, "iso_a4_210x297mm"},         // A4 in whole points
		{842, 595, "iso_a4_210x297mm"},         // landscape
		{612, 792, "na_letter_8.5x11in"},       // letter
		{516, 729, "jis_b5_182x257mm"},         // B5
		{283.5, 419.5, "jpn_hagaki_100x148mm"}, // hagaki
		{500, 500, ""},
	} {
		if got := ChooseMedia(c.w, c.h, ready, supported); got != c.want {
			t.Errorf("ChooseMedia(%g, %g) = %q, want %q", c.w, c.h, got, c.want)
		}
	}
}
