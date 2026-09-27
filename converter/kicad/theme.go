package kicad

import (
	"strings"

	"github.com/shibukawa/bdf"
)

// The colors of KiCad's default theme ("KiCad Default") in the schematic
// editor.
var (
	schBackground     = bdf.RGB(245, 244, 239)
	schBus            = bdf.RGB(0, 0, 132)
	schBusJunction    = bdf.RGB(0, 0, 132)
	schBody           = bdf.RGB(255, 255, 194)
	schOutline        = bdf.RGB(132, 0, 0)
	schERCError       = bdf.RGBA(230, 9, 13, 204)
	schFields         = bdf.RGB(132, 0, 132)
	schJunction       = bdf.RGB(0, 150, 0)
	schLabelGlobal    = bdf.RGB(132, 0, 0)
	schLabelHier      = bdf.RGB(114, 86, 0)
	schLabelLocal     = bdf.RGB(15, 15, 15)
	schNetclassFlag   = bdf.RGB(72, 72, 72)
	schNoConnect      = bdf.RGB(0, 0, 132)
	schNote           = bdf.RGB(0, 0, 194)
	schPin            = bdf.RGB(132, 0, 0)
	schPinName        = bdf.RGB(0, 100, 100)
	schPinNumber      = bdf.RGB(169, 0, 0)
	schReference      = bdf.RGB(0, 100, 100)
	schRuleArea       = bdf.RGB(255, 0, 0)
	schSheet          = bdf.RGB(132, 0, 0)
	schSheetFields    = bdf.RGB(132, 0, 132)
	schSheetFileName  = bdf.RGB(114, 86, 0)
	schSheetLabel     = bdf.RGB(0, 100, 100)
	schSheetName      = bdf.RGB(0, 100, 100)
	schValue          = bdf.RGB(0, 100, 100)
	schWire           = bdf.RGB(0, 150, 0)
	schWorksheet      = bdf.RGB(132, 0, 0)
	schSheetBackdrop  = bdf.RGBA(255, 255, 255, 0)
	schHidden         = bdf.RGB(94, 194, 194)
	schPrivateNote    = bdf.RGB(72, 72, 255)
	schNoteBackground = bdf.RGBA(0, 0, 0, 0)
)

// The colors of KiCad's default theme in the board editor.
var (
	pcbBackground     = bdf.RGB(0, 16, 35)
	pcbWorksheet      = bdf.RGB(200, 114, 171)
	pcbEdgeCuts       = bdf.RGB(208, 210, 205)
	pcbMargin         = bdf.RGB(255, 38, 226)
	pcbViaThrough     = bdf.RGB(236, 236, 236)
	pcbViaBlind       = bdf.RGB(187, 151, 38)
	pcbViaMicro       = bdf.RGB(0, 132, 132)
	pcbViaHole        = bdf.RGB(227, 183, 46)
	pcbPadThroughHole = bdf.RGB(227, 183, 46)
	pcbPlatedHole     = bdf.RGB(194, 194, 0)
	pcbNPTHole        = bdf.RGB(26, 196, 210)
	pcbHoleBackground = bdf.RGB(0, 16, 35)
)

// pcbLayerColor returns the color of a board layer by its canonical name.
func pcbLayerColor(name string) bdf.Color {
	switch name {
	case "F.Cu":
		return bdf.RGB(200, 52, 52)
	case "B.Cu":
		return bdf.RGB(77, 127, 196)
	case "F.Adhes":
		return bdf.RGB(132, 0, 132)
	case "B.Adhes":
		return bdf.RGB(0, 0, 132)
	case "F.Paste":
		return bdf.RGBA(180, 160, 154, 230)
	case "B.Paste":
		return bdf.RGBA(0, 194, 194, 230)
	case "F.SilkS":
		return bdf.RGB(242, 237, 161)
	case "B.SilkS":
		return bdf.RGB(232, 178, 167)
	case "F.Mask":
		return bdf.RGBA(216, 100, 255, 102)
	case "B.Mask":
		return bdf.RGBA(2, 255, 238, 102)
	case "Dwgs.User", "User.Drawings":
		return bdf.RGB(194, 194, 194)
	case "Cmts.User", "User.Comments":
		return bdf.RGB(89, 148, 220)
	case "Eco1.User", "User.Eco1":
		return bdf.RGB(180, 219, 210)
	case "Eco2.User", "User.Eco2":
		return bdf.RGB(216, 200, 82)
	case "Edge.Cuts":
		return pcbEdgeCuts
	case "Margin":
		return pcbMargin
	case "F.CrtYd":
		return bdf.RGB(255, 38, 226)
	case "B.CrtYd":
		return bdf.RGB(38, 233, 255)
	case "F.Fab":
		return bdf.RGB(175, 175, 175)
	case "B.Fab":
		return bdf.RGB(88, 93, 132)
	}
	if strings.HasPrefix(name, "In") && strings.HasSuffix(name, ".Cu") {
		n := 0
		for _, c := range name[2 : len(name)-3] {
			if c < '0' || c > '9' {
				n = 0
				break
			}
			n = n*10 + int(c-'0')
		}
		if n >= 1 {
			return innerCopper[(n-1)%len(innerCopper)]
		}
	}
	if strings.HasPrefix(name, "User.") {
		n := 0
		for _, c := range name[5:] {
			if c < '0' || c > '9' {
				n = 0
				break
			}
			n = n*10 + int(c-'0')
		}
		if n >= 1 {
			return userColors[(n-1)%len(userColors)]
		}
	}
	return bdf.RGB(194, 194, 194)
}

var innerCopper = []bdf.Color{
	bdf.RGB(127, 200, 127), bdf.RGB(206, 125, 44), bdf.RGB(79, 203, 203), bdf.RGB(219, 98, 139),
	bdf.RGB(167, 165, 198), bdf.RGB(40, 204, 217), bdf.RGB(232, 178, 167), bdf.RGB(242, 237, 161),
	bdf.RGB(141, 203, 129), bdf.RGB(237, 124, 51), bdf.RGB(91, 195, 235), bdf.RGB(247, 111, 142),
	bdf.RGB(167, 165, 198), bdf.RGB(40, 204, 217), bdf.RGB(232, 178, 167), bdf.RGB(242, 237, 161),
	bdf.RGB(237, 124, 51), bdf.RGB(91, 195, 235), bdf.RGB(247, 111, 142), bdf.RGB(167, 165, 198),
	bdf.RGB(40, 204, 217), bdf.RGB(232, 178, 167), bdf.RGB(242, 237, 161), bdf.RGB(237, 124, 51),
	bdf.RGB(91, 195, 235), bdf.RGB(247, 111, 142), bdf.RGB(167, 165, 198), bdf.RGB(40, 204, 217),
	bdf.RGB(232, 178, 167), bdf.RGB(242, 237, 161),
}

var userColors = []bdf.Color{
	bdf.RGB(194, 194, 194), bdf.RGB(89, 148, 220), bdf.RGB(180, 219, 210), bdf.RGB(216, 200, 82),
	bdf.RGB(194, 194, 194), bdf.RGB(89, 148, 220), bdf.RGB(180, 219, 210), bdf.RGB(216, 200, 82),
	bdf.RGB(232, 178, 167),
}

// mix blends two colors: t of b into a.
func mix(a, b bdf.Color, t float64) bdf.Color {
	ar, ag, ab, aa := rgba(a)
	br, bg, bb, _ := rgba(b)
	f := func(x, y uint8) uint8 { return uint8(float64(x)*(1-t) + float64(y)*t + 0.5) }
	return bdf.RGBA(f(ar, br), f(ag, bg), f(ab, bb), aa)
}

// dim is how KiCad draws the items of a symbol that is not placed (DNP):
// desaturated and halfway to the background.
func dim(c, background bdf.Color) bdf.Color {
	r, g, b, a := rgba(c)
	// desaturate: the average of the brightest and darkest channel
	mx := max(r, g, b)
	mn := min(r, g, b)
	l := uint8((int(mx) + int(mn)) / 2)
	return mix(bdf.RGBA(l, l, l, a), background, 0.5)
}

// rgba returns the channels of a color.
func rgba(c bdf.Color) (r, g, b, a uint8) {
	return uint8(c >> 24), uint8(c >> 16), uint8(c >> 8), uint8(c)
}
