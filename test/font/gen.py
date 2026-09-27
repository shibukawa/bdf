"""Writes the test fonts of the font converter (converter/font/testdata)
with fontTools (pip install fonttools brotli), from STIX Two 2.13 b171
(https://github.com/stipub/stixfonts, the tag archive v2.13b171.zip; SIL
Open Font License 1.1, see converter/font/testdata/OFL.txt):

- stix.ttf: STIX Two Text Regular (TrueType) for Basic Latin, Latin-1,
  Greek, Cyrillic and the combining accents, with every layout feature
  (liga, smcp, c2sc, cv01–cv03, ss01, ss02, frac, onum, kern, mark …) and
  every name, without hinting.
- stix.otf: Basic Latin, Latin-1 and the combining accents of the CFF
  build.
- stix.woff, stix.woff2: Basic Latin of stix.ttf as web fonts.
- bold.ttc: a collection of Basic Latin of Regular and Bold (TrueType).
- variable.ttf: Basic Latin of the variable STIX Two Text (wght 400–700).
- restricted.ttf: Basic Latin of stix.ttf with an OS/2 fsType that forbids
  embedding (Restricted License), so that the converter draws outlines.
- color.ttf: Basic Latin of stix.ttf with COLRv0 color glyphs for A, B
  and C (layers of other glyphs in the colors of two CPAL palettes).
- features.ttf: a font of rectangles made here, whose GSUB and GPOS (from
  the feature file FEATURES below) hold a lookup of each kind: single
  substitutions of both formats, multiple, alternate, ligature, chained
  contextual in an extension lookup, reverse chaining, a character variant
  and a stylistic set with their names, pair adjustments of both formats,
  a single adjustment, cursive, mark to base and mark to mark attachment.

Usage: python3 test/font/gen.py <unzipped stixfonts-2.13b171>
"""
import os
import sys

from fontTools import subset
from fontTools.colorLib.builder import buildCOLR, buildCPAL
from fontTools.feaLib.builder import addOpenTypeFeaturesFromString
from fontTools.fontBuilder import FontBuilder
from fontTools.pens.ttGlyphPen import TTGlyphPen
from fontTools.ttLib import TTCollection, TTFont

src = sys.argv[1]
out = os.path.join(os.path.dirname(__file__), "../../converter/font/testdata")
os.makedirs(out, exist_ok=True)

LATIN = "U+0020-007E"
MORE = "U+0020-007E,U+00A0-00FF,U+0131,U+0152,U+0153,U+0300-030C,U+0327,U+0328,U+0370-03FF,U+0400-045F,U+2013,U+2014,U+2018-201E,U+2026,U+2044,U+FB01,U+FB02"


def sub(path, unicodes, dest, flavor=None, hinting=False):
    opts = subset.Options()
    opts.layout_features = ["*"]
    opts.name_IDs = ["*"]
    opts.name_languages = ["*"]
    opts.hinting = hinting
    opts.notdef_outline = True
    opts.glyph_names = True
    opts.flavor = flavor
    font = subset.load_font(path, opts)
    s = subset.Subsetter(opts)
    s.populate(unicodes=subset.parse_unicodes(unicodes))
    s.subset(font)
    if dest:
        subset.save_font(font, os.path.join(out, dest), opts)
    return font


ttf = os.path.join(src, "fonts/static_ttf/STIXTwoText-Regular.ttf")
sub(ttf, MORE, "stix.ttf")
sub(os.path.join(src, "fonts/static_otf/STIXTwoText-Regular.otf"), LATIN + ",U+00A0-00FF,U+0300-030C,U+0327,U+0328", "stix.otf")
sub(ttf, LATIN, "stix.woff", flavor="woff")
sub(ttf, LATIN, "stix.woff2", flavor="woff2")
sub(os.path.join(src, "fonts/variable_ttf/STIXTwoText[wght].ttf"), LATIN, "variable.ttf")

f = sub(ttf, LATIN, None)
f["OS/2"].fsType = 0x0002
f.save(os.path.join(out, "restricted.ttf"))

f = sub(ttf, LATIN, None)
cmap = f.getBestCmap()
f["CPAL"] = buildCPAL([[(0.85, 0.2, 0.2, 1.0), (0.1, 0.4, 0.8, 1.0)], [(0.1, 0.6, 0.3, 1.0), (0.9, 0.6, 0.1, 1.0)]])
f["COLR"] = buildCOLR({
    cmap[ord("A")]: [(cmap[ord("A")], 0), (cmap[ord("-")], 1)],
    cmap[ord("B")]: [(cmap[ord("B")], 1), (cmap[ord(".")], 0)],
    cmap[ord("C")]: [(cmap[ord("C")], 0), (cmap[ord("o")], 1)],
})
f.save(os.path.join(out, "color.ttf"))

c = TTCollection()
c.fonts = [sub(ttf, LATIN, None), sub(os.path.join(src, "fonts/static_ttf/STIXTwoText-Bold.ttf"), LATIN, None)]
c.save(os.path.join(out, "bold.ttc"))


FEATURES = """
languagesystem DFLT dflt;
languagesystem latn dflt;
languagesystem latn TRK;

table GDEF {
  GlyphClassDef [a b c d e f i a.alt b.alt c.alt], [f_i f_f_i], [acutecomb], ;
} GDEF;

lookup DELTA { sub [a b c] by [a.alt b.alt c.alt]; } DELTA;
lookup LIST { sub a by e; sub b by c; } LIST;
lookup MULTI { sub d by a b; } MULTI;
lookup ALT { sub a from [a.alt b.alt c.alt]; } ALT;
lookup LIGA { sub f i by f_i; sub f f i by f_f_i; } LIGA;
lookup CTX useExtension { sub a' b by a.alt; } CTX;
lookup RSUB { rsub [a b] c' d by e; } RSUB;

feature smcp { lookup DELTA; } smcp;
feature locl { script latn; language TRK; lookup LIST; } locl;
feature ccmp { lookup MULTI; } ccmp;
feature salt { lookup ALT; } salt;
feature liga { lookup LIGA; } liga;
feature calt { lookup CTX; } calt;
feature rclt { lookup RSUB; } rclt;
feature cv01 {
  cvParameters { FeatUILabelNameID { name "Open a"; }; Character 0x61; Character 0x62; };
  sub a by a.alt;
} cv01;
feature ss01 {
  featureNames { name "Round dots"; };
  sub b by b.alt;
} ss01;

markClass acutecomb <anchor 100 500> @TOP;
feature kern { pos a b -50; pos [c d] [e f] -30; } kern;
feature palt { pos a <-10 0 -20 0>; } palt;
feature curs { pos cursive f <anchor 0 0> <anchor 500 0>; } curs;
feature mark { pos base [a b] <anchor 250 450> mark @TOP; } mark;
feature mkmk { pos mark acutecomb <anchor 100 700> mark @TOP; } mkmk;
"""


def rect(w):
    pen = TTGlyphPen(None)
    pen.moveTo((50, 0))
    pen.lineTo((50, 600))
    pen.lineTo((w - 50, 600))
    pen.lineTo((w - 50, 0))
    pen.closePath()
    return pen.glyph()


order = [".notdef", "space", "a", "b", "c", "d", "e", "f", "i", "a.alt", "b.alt", "c.alt", "f_i", "f_f_i", "acutecomb"]
fb = FontBuilder(1000, isTTF=True)
fb.setupGlyphOrder(order)
fb.setupCharacterMap({0x20: "space", 0x61: "a", 0x62: "b", 0x63: "c", 0x64: "d", 0x65: "e", 0x66: "f", 0x69: "i", 0x301: "acutecomb"})
glyphs = {g: rect(500) for g in order}
glyphs["space"] = TTGlyphPen(None).glyph()
glyphs["f_i"], glyphs["f_f_i"] = rect(900), rect(1300)
fb.setupGlyf(glyphs)
metrics = {g: (500, 50) for g in order}
metrics["space"] = (250, 0)
metrics["f_i"], metrics["f_f_i"], metrics["acutecomb"] = (900, 50), (1300, 50), (0, 50)
fb.setupHorizontalMetrics(metrics)
fb.setupHorizontalHeader(ascent=800, descent=-200)
fb.setupNameTable({"familyName": "Features Test", "styleName": "Regular"})
fb.setupOS2(sTypoAscender=800, sTypoDescender=-200, usWinAscent=800, usWinDescent=200, fsType=0)
fb.setupPost()
addOpenTypeFeaturesFromString(fb.font, FEATURES)
fb.save(os.path.join(out, "features.ttf"))
