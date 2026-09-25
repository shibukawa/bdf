package bdf

// Opcodes of opset 1. See docs/spec.md §7.
const (
	OpSave        byte = 0x01
	OpRestore     byte = 0x02
	OpTransform   byte = 0x03
	OpTranslate   byte = 0x04
	OpScale       byte = 0x05
	OpClipPath    byte = 0x06
	OpClipRect    byte = 0x07
	OpFillColor   byte = 0x10
	OpFillPaint   byte = 0x11
	OpStrokeColor byte = 0x12
	OpStrokePaint byte = 0x13
	OpLine        byte = 0x14
	OpDash        byte = 0x15
	OpAlpha       byte = 0x16
	OpBlend       byte = 0x17
	OpShadow      byte = 0x18
	OpFilter      byte = 0x19
	OpFont        byte = 0x1A
	OpTextStyle   byte = 0x1B
	OpFillRect    byte = 0x20
	OpStrokeRect  byte = 0x21
	OpFillPath    byte = 0x22
	OpStrokePath  byte = 0x23
	OpFillPathAt  byte = 0x24
	OpFillPathRun byte = 0x25
	OpClearRect   byte = 0x26
	OpFillText    byte = 0x30
	OpStrokeText  byte = 0x31
	OpImage       byte = 0x40
	OpImageSub    byte = 0x41
	OpSmoothing   byte = 0x42
	OpUse         byte = 0x50
	OpUseAt       byte = 0x51
	OpGroupBegin  byte = 0x52
	OpGroupEnd    byte = 0x53
	OpLink        byte = 0x70
	OpMark        byte = 0x71
	OpExt         byte = 0xFF
)

// OpsetVersion is the instruction set version written by this package.
const OpsetVersion = 1

// FormatVersion is the container format version written by this package.
const FormatVersion = 1

// Fill rules.
const (
	NonZero byte = 0
	EvenOdd byte = 1
)

// Line caps.
const (
	CapButt   byte = 0
	CapRound  byte = 1
	CapSquare byte = 2
)

// Line joins.
const (
	JoinMiter byte = 0
	JoinRound byte = 1
	JoinBevel byte = 2
)

// Text alignment.
const (
	AlignLeft   byte = 0
	AlignRight  byte = 1
	AlignCenter byte = 2
	AlignStart  byte = 3
	AlignEnd    byte = 4
)

// Text baseline.
const (
	BaselineAlphabetic  byte = 0
	BaselineTop         byte = 1
	BaselineMiddle      byte = 2
	BaselineBottom      byte = 3
	BaselineHanging     byte = 4
	BaselineIdeographic byte = 5
)

// Text direction.
const (
	DirInherit byte = 0
	DirLTR     byte = 1
	DirRTL     byte = 2
)

// Blend modes (globalCompositeOperation). See docs/spec.md appendix A.
const (
	BlendSourceOver      byte = 0
	BlendMultiply        byte = 1
	BlendScreen          byte = 2
	BlendOverlay         byte = 3
	BlendDarken          byte = 4
	BlendLighten         byte = 5
	BlendColorDodge      byte = 6
	BlendColorBurn       byte = 7
	BlendHardLight       byte = 8
	BlendSoftLight       byte = 9
	BlendDifference      byte = 10
	BlendExclusion       byte = 11
	BlendHue             byte = 12
	BlendSaturation      byte = 13
	BlendColor           byte = 14
	BlendLuminosity      byte = 15
	BlendDestinationOver byte = 16
	BlendDestinationIn   byte = 17
	BlendDestinationOut  byte = 18
	BlendSourceIn        byte = 19
	BlendSourceOut       byte = 20
	BlendSourceAtop      byte = 21
	BlendDestinationAtop byte = 22
	BlendXor             byte = 23
	BlendCopy            byte = 24
	BlendLighter         byte = 25
)

// BlendNames maps blend mode values to globalCompositeOperation strings.
var BlendNames = []string{
	"source-over", "multiply", "screen", "overlay", "darken", "lighten",
	"color-dodge", "color-burn", "hard-light", "soft-light", "difference",
	"exclusion", "hue", "saturation", "color", "luminosity",
	"destination-over", "destination-in", "destination-out", "source-in",
	"source-out", "source-atop", "destination-atop", "xor", "copy", "lighter",
}

// opInfo describes the operand signature of an opcode.
//
// Signature characters: f=f32, b=u8, c=u32 color, v=varuint, s=string ref.
// Special signatures: "D" (DASH), "R" (FILL_PATH_RUN), "X" (EXT).
type opInfo struct {
	Name string
	Sig  string
}

var opTable = map[byte]opInfo{
	OpSave:        {"SAVE", ""},
	OpRestore:     {"RESTORE", ""},
	OpTransform:   {"TRANSFORM", "ffffff"},
	OpTranslate:   {"TRANSLATE", "ff"},
	OpScale:       {"SCALE", "ff"},
	OpClipPath:    {"CLIP_PATH", "vb"},
	OpClipRect:    {"CLIP_RECT", "ffff"},
	OpFillColor:   {"FILL_COLOR", "c"},
	OpFillPaint:   {"FILL_PAINT", "v"},
	OpStrokeColor: {"STROKE_COLOR", "c"},
	OpStrokePaint: {"STROKE_PAINT", "v"},
	OpLine:        {"LINE", "fbbf"},
	OpDash:        {"DASH", "D"},
	OpAlpha:       {"ALPHA", "f"},
	OpBlend:       {"BLEND", "b"},
	OpShadow:      {"SHADOW", "cfff"},
	OpFilter:      {"FILTER", "s"},
	OpFont:        {"FONT", "vf"},
	OpTextStyle:   {"TEXT_STYLE", "bbbf"},
	OpFillRect:    {"FILL_RECT", "ffff"},
	OpStrokeRect:  {"STROKE_RECT", "ffff"},
	OpFillPath:    {"FILL_PATH", "vb"},
	OpStrokePath:  {"STROKE_PATH", "v"},
	OpFillPathAt:  {"FILL_PATH_AT", "vbff"},
	OpFillPathRun: {"FILL_PATH_RUN", "R"},
	OpClearRect:   {"CLEAR_RECT", "ffff"},
	OpFillText:    {"FILL_TEXT", "sfff"},
	OpStrokeText:  {"STROKE_TEXT", "sfff"},
	OpImage:       {"IMAGE", "vffff"},
	OpImageSub:    {"IMAGE_SUB", "vffffffff"},
	OpSmoothing:   {"SMOOTHING", "bb"},
	OpUse:         {"USE", "v"},
	OpUseAt:       {"USE_AT", "vff"},
	OpGroupBegin:  {"GROUP_BEGIN", "fbffff"},
	OpGroupEnd:    {"GROUP_END", ""},
	OpLink:        {"LINK", "ffffs"},
	OpMark:        {"MARK", "bs"},
	OpExt:         {"EXT", "X"},
}
