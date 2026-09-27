package gerber

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// The colors of the layer views: those of KiCad's default theme, on its
// background.
var (
	background   = bdf.RGB(0, 16, 35)
	outlineColor = bdf.RGB(208, 210, 16)
	outlineDim   = bdf.RGB(84, 88, 60)
	otherColor   = bdf.RGB(194, 194, 194)
	innerColors  = []bdf.Color{
		bdf.RGB(127, 200, 127), bdf.RGB(206, 206, 125), bdf.RGB(79, 203, 203), bdf.RGB(219, 98, 139),
		bdf.RGB(167, 165, 198), bdf.RGB(40, 204, 217), bdf.RGB(232, 178, 167), bdf.RGB(242, 237, 161),
	}
)

// layerColor returns the color a layer is drawn in in its view.
func layerColor(f function) bdf.Color {
	top := f.side != sBottom
	pick := func(t, b bdf.Color) bdf.Color {
		if top {
			return t
		}
		return b
	}
	switch f.kind {
	case kCopper:
		if f.side == sInner {
			return innerColors[max(f.index-1, 0)%len(innerColors)]
		}
		return pick(bdf.RGB(200, 52, 52), bdf.RGB(77, 127, 196))
	case kMask:
		return pick(bdf.RGB(216, 100, 255), bdf.RGB(2, 255, 238))
	case kSilk:
		return pick(bdf.RGB(242, 237, 161), bdf.RGB(232, 178, 167))
	case kPaste:
		return pick(bdf.RGB(180, 160, 154), bdf.RGB(0, 194, 194))
	case kOutline:
		return outlineColor
	case kDrill:
		if f.plated == "nonplated" {
			return bdf.RGB(26, 196, 210)
		}
		return bdf.RGB(227, 183, 46)
	}
	return otherColor
}

// palette is the colors of the board views.
type palette struct {
	substrate bdf.Color
	// copper under the solder mask, and the finish of exposed copper
	copper, finish bdf.Color
	mask           bdf.Color
	maskAlpha      float32
	silk           bdf.Color
}

// maskColors are the solder mask colors of the job file format, and their
// opacity over the board.
var maskColors = map[string]struct {
	c     bdf.Color
	alpha float32
}{
	"green":  {bdf.RGB(0, 110, 40), 0.82},
	"red":    {bdf.RGB(170, 20, 20), 0.85},
	"blue":   {bdf.RGB(10, 50, 150), 0.85},
	"black":  {bdf.RGB(16, 16, 16), 0.9},
	"white":  {bdf.RGB(235, 235, 230), 0.9},
	"yellow": {bdf.RGB(220, 180, 20), 0.85},
	"purple": {bdf.RGB(80, 30, 120), 0.85},
}

var silkColors = map[string]bdf.Color{
	"white":  bdf.RGB(245, 245, 240),
	"black":  bdf.RGB(20, 20, 20),
	"yellow": bdf.RGB(240, 220, 60),
	"red":    bdf.RGB(220, 40, 40),
	"blue":   bdf.RGB(40, 80, 220),
	"green":  bdf.RGB(40, 160, 60),
}

var finishColors = map[string]bdf.Color{
	"gold":   bdf.RGB(212, 170, 80),
	"silver": bdf.RGB(200, 200, 200),
	"copper": bdf.RGB(196, 120, 64),
}

// mix returns c drawn with opacity a over the opaque color under.
func mix(c, under bdf.Color, a float32) bdf.Color {
	ch := func(shift uint) uint8 {
		v := float64(a)*float64(uint8(c>>shift)) + (1-float64(a))*float64(uint8(under>>shift))
		return uint8(math.Round(v))
	}
	return bdf.RGB(ch(24), ch(16), ch(8))
}

// hexColor reads #rrggbb.
func hexColor(s string) (bdf.Color, bool) {
	if len(s) != 7 || s[0] != '#' {
		return 0, false
	}
	v, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return 0, false
	}
	return bdf.RGB(uint8(v>>16), uint8(v>>8), uint8(v)), true
}

// finishOf maps the surface finishes of the job file format onto colors.
func finishOf(s string) string {
	s = strings.ToLower(s)
	switch {
	case strings.Contains(s, "enig") || strings.Contains(s, "gold") || strings.Contains(s, "enepig"):
		return "gold"
	case strings.Contains(s, "hal") || strings.Contains(s, "silver") || strings.Contains(s, "tin"):
		return "silver"
	case strings.Contains(s, "osp") || s == "none":
		return "copper"
	}
	return ""
}

// palette returns the colors of the board views: the options', else the
// job file's, else a green board with white silkscreen and gold pads.
func (c *conv) palette() (palette, error) {
	p := palette{substrate: bdf.RGB(74, 70, 44), copper: bdf.RGB(214, 170, 100), silk: silkColors["white"], finish: finishColors["gold"]}
	mask, silk, finish := c.opts.Mask, c.opts.Silkscreen, c.opts.Finish
	if c.job != nil {
		for _, m := range c.job.MaterialStackup {
			switch strings.ToLower(m.Type) {
			case "soldermask":
				if mask == "" {
					mask = strings.ToLower(m.Color)
				}
			case "legend":
				if silk == "" {
					silk = strings.ToLower(m.Color)
				}
			}
		}
		if finish == "" {
			finish = finishOf(c.job.GeneralSpecs.Finish)
		}
	}
	mc := maskColors["green"]
	if mask != "" {
		if v, ok := maskColors[strings.ToLower(mask)]; ok {
			mc = v
		} else if h, ok := hexColor(mask); ok {
			mc.c, mc.alpha = h, 0.85
		} else if c.opts.Mask != "" {
			return p, fmt.Errorf("gerber: unknown solder mask color %q", mask)
		}
	}
	p.mask, p.maskAlpha = mc.c, mc.alpha
	if silk != "" {
		if v, ok := silkColors[strings.ToLower(silk)]; ok {
			p.silk = v
		} else if h, ok := hexColor(silk); ok {
			p.silk = h
		} else if c.opts.Silkscreen != "" {
			return p, fmt.Errorf("gerber: unknown silkscreen color %q", silk)
		}
	}
	if finish != "" {
		if v, ok := finishColors[strings.ToLower(finish)]; ok {
			p.finish = v
		} else if c.opts.Finish != "" {
			return p, fmt.Errorf("gerber: unknown finish %q (want gold, silver or copper)", finish)
		}
	}
	return p, nil
}
