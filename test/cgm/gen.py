"""Writes the CGM test metafiles of converter/cgm into converter/cgm/testdata.

Standard library only. Each metafile is described once as a list of elements
with typed parameters and written in the binary encoding (ISO/IEC 8632-3);
shapes is written in the clear text encoding (ISO/IEC 8632-4) too, and the
Go tests check that both draw the same (TestClearTextMatchesBinary). Text
uses only characters of the test fonts (converter/pptx/testdata/fonts).

shapes.cgm (and shapes-text.cgm), a CAD-like sheet: A3 at 0.01 mm per VDC
unit (32-bit integer VDC), indexed colours with a colour table. Line types
1 to 5 and one defined by the metafile, widths in mm, caps; the five
markers; solid, hollow, hatched (the six hatch indexes), patterned and
interpolated fills; a polygon set with a hole and invisible edges; circles,
ellipses, the arcs (three-point, centre, reversed, closed as pies and
chords, elliptical, hyperbolic, parabolic), a Bezier curve, a B-spline and
a NURBS; a closed figure with a hole; text at the nine alignments, turned,
expanded and spaced, vertical, restricted to a box, appended in another
colour, and in Japanese (JIS X 0208 designated as the alternate character
set, in 8-bit codes); a cell array; a clip rectangle; a segment and two
transformed copies.

illustration.cgm, a WebCGM-like illustration: real VDC (32-bit floating
point) with the y axis down, abstract scaling, direct colours, UTF-8 text
(the complete code designation of WebCGM), an application structure whose
visibility is off (not drawn), a tile array of an uncompressed bitonal tile
and a PNG tile, and a second picture.

sjis.cgm, as Japanese CAD programs write: A4 portrait at 0.1 mm per unit,
Shift_JIS names and text without a character set that says so.
"""

import math
import os
import struct
import zlib

OUT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "converter", "cgm", "testdata")

# element names of the clear text encoding by (class, id)
NAMES = {
    (0, 1): "BEGMF", (0, 2): "ENDMF", (0, 3): "BEGPIC", (0, 4): "BEGPICBODY", (0, 5): "ENDPIC",
    (0, 6): "BEGSEG", (0, 7): "ENDSEG", (0, 8): "BEGFIGURE", (0, 9): "ENDFIGURE",
    (0, 19): "BEGTILEARRAY", (0, 20): "ENDTILEARRAY", (0, 21): "BEGAPS", (0, 22): "BEGAPSBODY", (0, 23): "ENDAPS",
    (1, 1): "MFVERSION", (1, 2): "MFDESC", (1, 3): "VDCTYPE", (1, 4): "INTEGERPREC", (1, 5): "REALPREC",
    (1, 7): "COLRPREC", (1, 8): "COLRINDEXPREC", (1, 9): "MAXCOLRINDEX", (1, 10): "COLRVALUEEXT",
    (1, 11): "MFELEMLIST", (1, 13): "FONTLIST", (1, 14): "CHARSETLIST", (1, 15): "CHARCODING",
    (2, 1): "SCALEMODE", (2, 2): "COLRMODE", (2, 3): "LINEWIDTHMODE", (2, 4): "MARKERSIZEMODE",
    (2, 5): "EDGEWIDTHMODE", (2, 6): "VDCEXT", (2, 7): "BACKCOLR", (2, 16): "INTSTYLEMODE",
    (2, 17): "LINEEDGETYPEDEF",
    (3, 1): "VDCINTEGERPREC", (3, 2): "VDCREALPREC", (3, 5): "CLIPRECT", (3, 6): "CLIP", (3, 10): "NEWREGION",
    (4, 1): "LINE", (4, 2): "DISJTLINE", (4, 3): "MARKER", (4, 4): "TEXT", (4, 5): "RESTRTEXT", (4, 6): "APNDTEXT",
    (4, 7): "POLYGON", (4, 8): "POLYGONSET", (4, 9): "CELLARRAY", (4, 11): "RECT", (4, 12): "CIRCLE",
    (4, 13): "ARC3PT", (4, 14): "ARC3PTCLOSE", (4, 15): "ARCCTR", (4, 16): "ARCCTRCLOSE", (4, 17): "ELLIPSE",
    (4, 18): "ELLIPARC", (4, 19): "ELLIPARCCLOSE", (4, 20): "ARCCTRREV", (4, 22): "HYPERBARC",
    (4, 23): "PARABARC", (4, 24): "NUB", (4, 25): "NURB", (4, 26): "POLYBEZIER", (4, 28): "BITONALTILE",
    (4, 29): "TILE",
    (5, 2): "LINETYPE", (5, 3): "LINEWIDTH", (5, 4): "LINECOLR", (5, 6): "MARKERTYPE", (5, 7): "MARKERSIZE",
    (5, 8): "MARKERCOLR", (5, 10): "TEXTFONTINDEX", (5, 12): "CHAREXPAN", (5, 13): "CHARSPACE",
    (5, 14): "TEXTCOLR", (5, 15): "CHARHEIGHT", (5, 16): "CHARORI", (5, 17): "TEXTPATH", (5, 18): "TEXTALIGN",
    (5, 19): "CHARSETINDEX", (5, 20): "ALTCHARSETINDEX", (5, 22): "INTSTYLE", (5, 23): "FILLCOLR",
    (5, 24): "HATCHINDEX", (5, 25): "PATINDEX", (5, 27): "EDGETYPE", (5, 28): "EDGEWIDTH", (5, 29): "EDGECOLR",
    (5, 30): "EDGEVIS", (5, 31): "FILLREFPT", (5, 32): "PATTABLE", (5, 33): "PATSIZE", (5, 34): "COLRTABLE",
    (5, 37): "LINECAP", (5, 42): "RESTRTEXTTYPE", (5, 43): "INTERPINT",
    (8, 1): "COPYSEG", (9, 1): "APSATTR",
}


# typed parameters
def I(v): return ("I", v)
def IX(v): return ("IX", v)
def E(v, word): return ("E", v, word)
def R(v): return ("R", v)
def V(v): return ("VDC", v)
def P(x, y): return ("P", (x, y))
def CI(v): return ("CI", v)
def CD(r, g, b): return ("CD", (r, g, b))
def S(b): return ("S", b if isinstance(b, bytes) else b.encode("latin-1"))
def N(v): return ("N", v)
def F(v): return ("FLOAT", v)  # SCALING MODE's factor: floating point in the binary encoding
def BYTES(b, text=""): return ("RAW", b, text)  # parameters given encoded (colour lists, bitstreams)
def BINARY(p): return ("BONLY", p)  # a parameter of the binary encoding only


class Meta:
    """A metafile: its elements and the precisions they are written with."""

    def __init__(self):
        self.elems = []

    def e(self, cls, eid, *params):
        self.elems.append((cls, eid, list(params)))



def real_format(form, a, b):
    if form == 0:
        return ("float", 32 if (a, b) == (9, 23) else 64)
    return ("fixed", a + b)


class Binary:
    def __init__(self):
        self.int_bits, self.vdc_real, self.vdc_int_bits = 16, False, 16
        self.real = ("fixed", 32)
        self.vdc_realf = ("fixed", 32)
        self.colour_bits, self.index_bits = 8, 8
        self.direct = False

    def real_bytes(self, v, f):
        kind, bits = f
        if kind == "float":
            return struct.pack(">f" if bits == 32 else ">d", v)
        whole = math.floor(v)
        frac = int(round((v - whole) * (1 << (bits // 2))))
        if frac == 1 << (bits // 2):
            whole, frac = whole + 1, 0
        return struct.pack(">hH" if bits == 32 else ">iI", whole, frac)

    def int_bytes(self, v, bits):
        return v.to_bytes(bits // 8, "big", signed=True)

    def param(self, p):
        k = p[0]
        if k == "I":
            return self.int_bytes(p[1], self.int_bits)
        if k in ("IX", "E", "N"):
            return self.int_bytes(p[1], 16)
        if k == "R":
            return self.real_bytes(p[1], self.real)
        if k == "FLOAT":
            return struct.pack(">f", p[1])
        if k == "VDC":
            if self.vdc_real:
                return self.real_bytes(p[1], self.vdc_realf)
            return self.int_bytes(p[1], self.vdc_int_bits)
        if k == "P":
            return self.param(V(p[1][0])) + self.param(V(p[1][1]))
        if k == "CI":
            return p[1].to_bytes(self.index_bits // 8, "big")
        if k == "CD":
            return b"".join(c.to_bytes(self.colour_bits // 8, "big") for c in p[1])
        if k == "S":
            b = p[1]
            if len(b) < 255:
                return bytes([len(b)]) + b
            return b"\xff" + struct.pack(">H", len(b)) + b
        if k == "RAW":
            return p[1]
        if k == "BONLY":
            return self.param(p[1])
        raise ValueError(k)

    def element(self, cls, eid, params):
        data = b"".join(self.param(p) for p in params)
        # the precisions the element sets apply to the elements after it
        if (cls, eid) == (1, 3):
            self.vdc_real = params[0][1] == 1
        elif (cls, eid) == (1, 4):
            self.int_bits = params[0][1]
        elif (cls, eid) == (1, 5):
            self.real = real_format(params[0][1], params[1][1], params[2][1])
        elif (cls, eid) == (3, 1):
            self.vdc_int_bits = params[0][1]
        elif (cls, eid) == (3, 2):
            self.vdc_realf = real_format(params[0][1], params[1][1], params[2][1])
        head = cls << 12 | eid << 5
        if len(data) < 31:
            out = struct.pack(">H", head | len(data)) + data
        else:
            out = struct.pack(">H", head | 31)
            rest = data
            while True:
                part, rest = rest[:32766], rest[32766:]
                out += struct.pack(">H", len(part) | (0x8000 if rest else 0)) + part
                if not rest:
                    break
        if len(out) % 2:
            out += b"\0"
        return out

    def write(self, meta):
        return b"".join(self.element(c, i, p) for c, i, p in meta.elems)


def text_param(p):
    k = p[0]
    if k in ("I", "IX", "N", "CI"):
        return str(p[1])
    if k == "E":
        return p[2]
    if k in ("R", "FLOAT"):
        return repr(float(p[1]))
    if k == "VDC":
        return str(p[1])
    if k == "P":
        return "(%s,%s)" % (p[1][0], p[1][1])
    if k == "CD":
        return "%d %d %d" % p[1]
    if k == "S":
        return "'" + p[1].decode("latin-1").replace("'", "''") + "'"
    if k == "RAW":
        return p[2]
    if k == "BONLY":
        return ""
    raise ValueError(k)


def clear_text(meta):
    lines = []
    for cls, eid, params in meta.elems:
        name = NAMES[(cls, eid)]
        # the precisions of the binary encoding are sizes in bits; clear
        # text only needs the largest colour component value
        if (cls, eid) in ((1, 4), (1, 5), (3, 1), (3, 2), (1, 8)):
            continue
        if (cls, eid) == (1, 7):
            lines.append("COLRPREC %d;" % ((1 << params[0][1]) - 1))
            continue
        args = " ".join(a for a in (text_param(p) for p in params) if a)
        lines.append((name + " " + args).strip() + ";")
    return ("\n".join(lines) + "\n").encode("latin-1")


# --- shapes.cgm -----------------------------------------------------------

U = 100  # VDC units per mm


def mmp(x, y):
    return P(int(round(x * U)), int(round(y * U)))


def mmv(v):
    return V(int(round(v * U)))


PALETTE = [(255, 255, 255), (0, 0, 0), (220, 30, 30), (20, 140, 40), (30, 70, 200), (230, 160, 0),
           (0, 150, 160), (160, 40, 160), (120, 120, 120), (250, 220, 200), (200, 225, 250)]
WHITE, BLACK, RED, GREEN, BLUE, ORANGE, TEAL, PURPLE, GREY, PEACH, SKY = range(11)


def shapes():
    m = Meta()
    m.e(0, 1, S("shapes"))
    m.e(1, 1, I(3))
    m.e(1, 2, S("CGM test drawing of converter/cgm"))
    m.e(1, 11, I(1), IX(-1), IX(1))  # the drawing set
    m.e(1, 9, CI(len(PALETTE) - 1))
    m.e(1, 13, S("Helvetica"), S("Times-Bold"), S("MS Gothic"))
    m.e(1, 14, E(0, "STD94"), S("B"), E(2, "STD94MULTIBYTE"), S("B"))
    m.e(1, 15, E(1, "BASIC8BIT"))
    m.e(0, 3, S("sheet"))
    m.e(2, 1, E(1, "METRIC"), F(0.01))
    m.e(2, 3, E(3, "MM"))
    m.e(2, 4, E(3, "MM"))
    m.e(2, 5, E(3, "MM"))
    m.e(3, 1, I(32))
    m.e(2, 6, P(0, 0), P(420 * U, 297 * U))
    m.e(2, 7, CD(255, 255, 255))
    m.e(2, 17, IX(-1), R(8.0), I(6), I(1), I(1), I(1))  # a chain line of 8 mm
    m.e(0, 4)
    m.e(5, 34, CI(0), *[CD(*c) for c in PALETTE])

    # the border and a title block
    m.e(5, 22, E(4, "EMPTY"))
    m.e(5, 30, E(1, "ON"))
    m.e(5, 28, R(0.7))
    m.e(5, 29, CI(BLACK))
    m.e(4, 11, mmp(10, 10), mmp(410, 287))
    m.e(5, 28, R(0.35))
    m.e(4, 11, mmp(290, 10), mmp(410, 40))
    m.e(5, 30, E(0, "OFF"))
    m.e(5, 3, R(0.25))
    m.e(5, 4, CI(BLACK))
    m.e(4, 2, mmp(290, 25), mmp(410, 25), mmp(330, 10), mmp(330, 25))

    # line types 1 to 5, the metafile's own, and widths
    for i, t in enumerate([1, 2, 3, 4, 5, -1]):
        m.e(5, 2, IX(t))
        m.e(5, 3, R([0.25, 0.35, 0.5, 0.35, 0.25, 0.5][i]))
        m.e(5, 4, CI([BLACK, RED, GREEN, BLUE, PURPLE, TEAL][i]))
        y = 275 - 6 * i
        m.e(4, 1, mmp(20, y), mmp(110, y))
    m.e(5, 2, IX(1))
    m.e(5, 4, CI(GREY))
    m.e(5, 3, R(3.0))
    for i, cap in enumerate([2, 3, 4]):
        m.e(5, 37, IX(cap), IX(1))
        m.e(4, 1, mmp(30, 236 - 7 * i), mmp(100, 236 - 7 * i))
    m.e(5, 37, IX(1), IX(1))
    m.e(5, 3, R(0.25))
    m.e(5, 4, CI(BLACK))
    m.e(4, 1, mmp(20, 243), mmp(20, 214), mmp(110, 214))

    # markers
    m.e(5, 7, R(4.0))
    for i in range(5):
        m.e(5, 6, IX(i + 1))
        m.e(5, 8, CI([BLACK, RED, GREEN, BLUE, PURPLE][i]))
        m.e(4, 3, mmp(130 + 12 * i, 275), mmp(130 + 12 * i, 265))

    # fills
    m.e(5, 30, E(1, "ON"))
    m.e(5, 28, R(0.35))
    m.e(5, 29, CI(BLACK))
    m.e(5, 22, E(1, "SOLID"))
    m.e(5, 23, CI(RED))
    m.e(4, 7, mmp(200, 250), mmp(230, 250), mmp(215, 278))
    m.e(5, 22, E(0, "HOLLOW"))
    m.e(5, 23, CI(BLUE))
    m.e(5, 30, E(0, "OFF"))
    m.e(4, 11, mmp(240, 250), mmp(265, 278))
    m.e(5, 30, E(1, "ON"))
    m.e(5, 22, E(3, "HATCH"))
    m.e(5, 23, CI(GREEN))
    for i in range(6):
        m.e(5, 24, IX(i + 1))
        x = 275 + 20 * (i % 3)
        y = 257 if i < 3 else 237
        m.e(4, 11, mmp(x, y), mmp(x + 16, y + 16))
    # a 4 x 4 checker pattern of 2 mm cells
    cells = [ORANGE if (x + y) % 2 == 0 else SKY for y in range(4) for x in range(4)]
    m.e(5, 32, IX(1), I(4), I(4), I(8), BYTES(bytes(cells), " ".join(str(c) for c in cells)))
    m.e(5, 33, mmv(0), mmv(8), mmv(8), mmv(0))
    m.e(5, 31, mmp(340, 230))
    m.e(5, 22, E(2, "PAT"))
    m.e(5, 25, IX(1))
    m.e(4, 11, mmp(340, 230), mmp(372, 262))
    m.e(5, 22, E(6, "INTERP"))
    m.e(5, 31, mmp(378, 230))
    m.e(5, 43, IX(1), mmv(0), mmv(32), I(2), R(0.0), R(1.0), CI(ORANGE), CI(BLUE))
    m.e(4, 11, mmp(378, 230), mmp(402, 262))
    # a polygon set: a square with a square hole, one edge invisible
    m.e(5, 22, E(1, "SOLID"))
    m.e(5, 23, CI(PEACH))
    m.e(4, 8, mmp(200, 205), E(1, "VIS"), mmp(240, 205), E(0, "INVIS"), mmp(240, 240), E(1, "VIS"),
        mmp(200, 240), E(3, "CLOSEVIS"),
        mmp(210, 215), E(1, "VIS"), mmp(230, 215), E(1, "VIS"), mmp(230, 230), E(1, "VIS"), mmp(210, 230), E(3, "CLOSEVIS"))

    # curves
    m.e(5, 30, E(0, "OFF"))
    m.e(5, 22, E(0, "HOLLOW"))
    m.e(5, 23, CI(BLACK))
    m.e(5, 3, R(0.35))
    m.e(5, 4, CI(BLUE))
    m.e(4, 12, mmp(30, 180), mmv(12))
    m.e(4, 17, mmp(65, 180), mmp(80, 188), mmp(61, 188))
    m.e(4, 13, mmp(90, 170), mmp(100, 192), mmp(115, 172))
    m.e(4, 15, mmp(135, 178), mmv(10), mmv(0), mmv(0), mmv(10), mmv(12))
    m.e(5, 4, CI(RED))
    m.e(4, 20, mmp(135, 178), mmv(10), mmv(0), mmv(0), mmv(10), mmv(8))
    m.e(5, 4, CI(PURPLE))
    m.e(4, 18, mmp(165, 180), mmp(180, 180), mmp(165, 188), mmv(0), mmv(1), mmv(-1), mmv(-1))
    m.e(4, 22, mmp(195, 180), mmp(200, 180), mmp(195, 185), mmv(10), mmv(-5), mmv(10), mmv(5))
    m.e(4, 23, mmp(235, 195), mmp(215, 170), mmp(255, 170))
    m.e(5, 4, CI(GREEN))
    m.e(4, 26, IX(2), mmp(265, 170), mmp(270, 195), mmp(285, 165), mmp(290, 185), mmp(295, 200), mmp(305, 170), mmp(310, 185))
    m.e(5, 4, CI(ORANGE))
    ctrl = [(320, 170), (325, 195), (340, 165), (350, 195), (360, 170)]
    m.e(4, 24, I(4), I(5), *[mmp(*c) for c in ctrl], *[R(k) for k in (0, 0, 0, 0, 0.5, 1, 1, 1, 1)], R(0.0), R(1.0))
    m.e(5, 4, CI(TEAL))
    s = math.sqrt(0.5)
    m.e(4, 25, I(3), I(3), mmp(380, 170), mmp(380, 190), mmp(400, 190),
        *[R(k) for k in (0, 0, 0, 1, 1, 1)], R(1.0), R(s), R(1.0), R(0.0), R(1.0))
    # closed arcs
    m.e(5, 22, E(1, "SOLID"))
    m.e(5, 23, CI(ORANGE))
    m.e(5, 30, E(1, "ON"))
    m.e(4, 14, mmp(20, 140), mmp(32, 152), mmp(44, 140), E(0, "PIE"))
    m.e(5, 23, CI(SKY))
    m.e(4, 16, mmp(65, 140), mmv(10), mmv(10), mmv(-10), mmv(10), mmv(12), E(1, "CHORD"))
    m.e(5, 23, CI(PEACH))
    m.e(4, 19, mmp(95, 140), mmp(110, 140), mmp(95, 148), mmv(1), mmv(0), mmv(0), mmv(-1), E(0, "PIE"))
    # a closed figure: a rounded shape with a round hole
    m.e(5, 23, CI(GREEN))
    m.e(0, 8)
    m.e(4, 1, mmp(125, 130), mmp(155, 130))
    m.e(4, 15, mmp(155, 140), mmv(0), mmv(-10), mmv(0), mmv(10), mmv(10))
    m.e(4, 1, mmp(155, 150), mmp(125, 150))
    m.e(3, 10)
    m.e(4, 12, mmp(140, 140), mmv(5))
    m.e(0, 9)

    # a cell array: 8 x 8 cells of the colour table
    cells = [(x + y) % 8 + 2 if x != y else BLACK for y in range(8) for x in range(8)]
    m.e(4, 9, mmp(180, 150), mmp(212, 118), mmp(212, 150), I(8), I(8), I(8), BINARY(E(1, "PACKED")),
        BYTES(bytes(cells), " ".join(str(c) for c in cells)))

    # a clip rectangle
    m.e(3, 5, mmp(225, 120), mmp(250, 150))
    m.e(5, 23, CI(PURPLE))
    m.e(5, 30, E(0, "OFF"))
    m.e(4, 12, mmp(235, 125), mmv(15))
    m.e(3, 5, mmp(0, 0), mmp(420, 297))

    # a segment and two copies of it
    m.e(0, 6, N(1))
    m.e(5, 23, CI(BLUE))
    m.e(4, 7, mmp(270, 130), mmp(285, 130), mmp(285, 125), mmp(295, 135), mmp(285, 145), mmp(285, 140), mmp(270, 140))
    m.e(0, 7)
    c30, s30 = math.cos(math.radians(30)), math.sin(math.radians(30))
    m.e(8, 1, N(1), R(1.0), R(0.0), R(0.0), R(1.0), mmv(40), mmv(0), E(0, "NO"))
    m.e(8, 1, N(1), R(round(c30 * 0.5, 4)), R(round(s30 * 0.5, 4)), R(round(-s30 * 0.5, 4)), R(round(c30 * 0.5, 4)),
        mmv(230), mmv(-45), E(0, "NO"))

    # text
    m.e(5, 10, IX(1))
    m.e(5, 15, mmv(4))
    m.e(5, 14, CI(BLACK))
    m.e(5, 7, R(2.0))
    m.e(5, 6, IX(2))
    m.e(5, 8, CI(RED))
    for j, (v, vw) in enumerate([(1, "TOP"), (3, "HALF"), (4, "BASE")]):
        for i, (h, hw) in enumerate([(1, "LEFT"), (2, "CTR"), (3, "RIGHT")]):
            x, y = 30 + 30 * i, 95 - 12 * j
            m.e(5, 18, E(h, hw), E(v, vw), R(0.0), R(0.0))
            m.e(4, 4, mmp(x, y), E(1, "FINAL"), S("Ag"))
            m.e(4, 3, mmp(x, y))
    m.e(5, 18, E(0, "NORMHORIZ"), E(0, "NORMVERT"), R(0.0), R(0.0))
    m.e(5, 16, mmv(-1), mmv(1), mmv(1), mmv(1))
    m.e(5, 10, IX(2))
    m.e(4, 4, mmp(120, 70), E(1, "FINAL"), S("Turned 45"))
    m.e(5, 16, mmv(0), mmv(1), mmv(1), mmv(0))
    m.e(5, 10, IX(1))
    m.e(5, 12, R(1.5))
    m.e(5, 13, R(0.25))
    m.e(4, 4, mmp(160, 90), E(1, "FINAL"), S("Wide"))
    m.e(5, 12, R(1.0))
    m.e(5, 13, R(0.0))
    m.e(5, 14, CI(RED))
    m.e(4, 4, mmp(160, 78), E(0, "NOTFINAL"), S("Red "))
    m.e(5, 14, CI(BLUE))
    m.e(4, 6, E(1, "FINAL"), S("Blue"))
    m.e(5, 14, CI(BLACK))
    # restricted to a box of 60 x 5 mm, drawn around it
    m.e(5, 42, IX(2))
    m.e(4, 5, mmv(60), mmv(5), mmp(160, 60), E(1, "FINAL"), S("Boxed text"))
    m.e(5, 22, E(4, "EMPTY"))
    m.e(5, 30, E(1, "ON"))
    m.e(5, 28, R(0.13))
    m.e(5, 29, CI(GREY))
    m.e(4, 11, mmp(160, 60), mmp(220, 65))
    m.e(5, 30, E(0, "OFF"))
    # Japanese: JIS X 0208 as the alternate set, its bytes with the high bit
    jis = "日本語の文字".encode("euc_jp")
    m.e(5, 10, IX(3))
    m.e(5, 20, IX(2))
    m.e(4, 4, mmp(160, 40), E(1, "FINAL"), S(b"JIS: " + jis))
    # vertical text
    m.e(5, 17, E(3, "DOWN"))
    m.e(4, 4, mmp(250, 100), E(1, "FINAL"), S("縦書き".encode("euc_jp")))
    m.e(5, 17, E(0, "RIGHT"))
    m.e(5, 10, IX(1))
    m.e(5, 15, mmv(5))
    m.e(4, 4, mmp(295, 30), E(1, "FINAL"), S("CGM test drawing"))
    m.e(5, 15, mmv(3))
    m.e(5, 20, IX(1))
    m.e(4, 4, mmp(295, 15), E(1, "FINAL"), S(b"\xa9 converter/cgm"))
    m.e(5, 20, IX(2))
    m.e(4, 4, mmp(335, 15), E(1, "FINAL"), S("図枠 A3".encode("euc_jp")))
    m.e(0, 5)
    m.e(0, 2)
    return m


# --- illustration.cgm -------------------------------------------------------

def png(w, h, rgb):
    """An RGB PNG of w x h pixels from a function of (x, y)."""
    raw = b"".join(b"\0" + b"".join(bytes(rgb(x, y)) for x in range(w)) for y in range(h))

    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xffffffff)
    return (b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0)) +
            chunk(b"IDAT", zlib.compress(raw, 9)) + chunk(b"IEND", b""))


def illustration():
    m = Meta()
    m.e(0, 1, S("illustration"))
    m.e(1, 1, I(4))
    m.e(1, 2, S("ProfileId:WebCGM ProfileEd:2.1"))
    m.e(1, 3, E(1, "REAL"))
    m.e(1, 5, E(0, "FLOATING"), I(9), I(23))
    m.e(1, 13, S("Helvetica"), S("Courier"))
    m.e(1, 14, E(4, "COMPLETECODE"), S(b"\x2f\x49"))
    m.e(1, 15, E(1, "BASIC8BIT"))
    m.e(0, 3, S("figure 1"))
    m.e(2, 2, E(1, "DIRECT"))
    m.e(2, 3, E(0, "ABS"))
    m.e(2, 5, E(0, "ABS"))
    m.e(3, 2, E(0, "FLOATING"), I(9), I(23))
    # the y axis down: the first corner is the bottom left
    m.e(2, 6, P(0.0, 600.0), P(800.0, 0.0))
    m.e(2, 7, CD(250, 250, 245))
    m.e(0, 4)
    m.e(5, 16, V(0.0), V(-1.0), V(1.0), V(0.0))
    m.e(5, 15, V(24.0))
    m.e(5, 14, CD(20, 40, 90))
    m.e(4, 4, P(40.0, 60.0), E(1, "FINAL"), S("Übersicht — ×2 °C".encode("utf-8")))
    m.e(5, 3, V(3.0))
    m.e(5, 4, CD(20, 40, 90))
    m.e(4, 1, P(40.0, 75.0), P(760.0, 75.0))
    m.e(5, 22, E(1, "SOLID"))
    m.e(5, 30, E(1, "ON"))
    m.e(5, 28, V(2.0))
    m.e(5, 29, CD(60, 60, 60))
    m.e(5, 23, CD(255, 200, 80))
    m.e(4, 11, P(60.0, 120.0), P(360.0, 300.0))
    # a layer whose visibility is off, and one that is on
    m.e(0, 21, S("hidden-layer"), S("layer"), E(0, "STLIST"))
    m.e(9, 1, S("visibility"), BYTES(b"\x06" + b"\x00\x05" + b"\x00\x01" + b"\x00\x00", "'5 1 0'"))
    m.e(0, 22)
    m.e(5, 23, CD(255, 0, 0))
    m.e(4, 12, P(210.0, 210.0), V(60.0))
    m.e(4, 4, P(80.0, 330.0), E(1, "FINAL"), S("hidden text"))
    m.e(0, 23)
    m.e(0, 21, S("callout"), S("grobject"), E(0, "STLIST"))
    m.e(9, 1, S("visibility"), BYTES(b"\x06" + b"\x00\x05" + b"\x00\x01" + b"\x00\x01", "'5 1 1'"))
    m.e(0, 22)
    m.e(5, 23, CD(90, 160, 220))
    m.e(4, 12, P(210.0, 210.0), V(40.0))
    m.e(5, 18, E(2, "CTR"), E(3, "HALF"), R(0.0), R(0.0))
    m.e(5, 14, CD(255, 255, 255))
    m.e(5, 10, IX(2))
    m.e(4, 4, P(210.0, 210.0), E(1, "FINAL"), S("A1"))
    m.e(0, 23)
    # a tile array: an uncompressed bitonal tile and a PNG tile of 16 x 16
    # cells of 10 units, along the x axis, lines progressing at 90 degrees
    # from it: down the page, whose y axis is down
    m.e(0, 19, P(420.0, 120.0), E(0, "0"), E(0, "90"), I(2), I(1), I(16), I(16), R(10.0), R(10.0),
        I(0), I(0), I(32), I(16))
    bits = b"".join(struct.pack(">H", sum(1 << (15 - x) for x in range(16) if (x // 4 + y // 4) % 2 == 0))
                    for y in range(16))
    m.e(4, 28, E(5, "5"), I(16), CD(255, 255, 255), CD(0, 0, 0), S(b""), BYTES(bits))
    tile = png(16, 16, lambda x, y: (x * 16, y * 16, 128))
    m.e(4, 29, E(9, "9"), I(0), I(8), S(b""), BYTES(tile))
    m.e(0, 20)
    m.e(0, 5)
    m.e(0, 3, S("figure 2"))
    m.e(2, 2, E(1, "DIRECT"))
    m.e(2, 6, P(0.0, 0.0), P(400.0, 300.0))
    m.e(0, 4)
    m.e(5, 22, E(1, "SOLID"))
    m.e(5, 23, CD(40, 150, 90))
    m.e(4, 17, P(200.0, 150.0), P(350.0, 150.0), P(200.0, 250.0))
    m.e(5, 14, CD(0, 0, 0))
    m.e(5, 15, V(20.0))
    m.e(5, 18, E(2, "CTR"), E(3, "HALF"), R(0.0), R(0.0))
    m.e(4, 4, P(200.0, 150.0), E(1, "FINAL"), S("Page 2"))
    m.e(0, 5)
    m.e(0, 2)
    return m


# --- sjis.cgm ---------------------------------------------------------------

def sjis():
    m = Meta()
    m.e(0, 1, S("テスト図".encode("shift_jis")))
    m.e(1, 1, I(1))
    m.e(1, 13, S("ＭＳ ゴシック".encode("shift_jis")), S("ＭＳ 明朝".encode("shift_jis")))
    m.e(0, 3, S("1"))
    m.e(2, 1, E(1, "METRIC"), F(0.1))
    m.e(2, 3, E(3, "MM"))
    m.e(2, 6, P(0, 0), P(2100, 2970))
    m.e(0, 4)
    m.e(5, 3, R(0.5))
    m.e(4, 1, P(100, 100), P(2000, 100), P(2000, 2870), P(100, 2870), P(100, 100))
    m.e(5, 3, R(0.25))
    m.e(4, 1, P(1200, 100), P(1200, 400), P(2000, 400))
    m.e(4, 2, P(1200, 250), P(2000, 250))
    m.e(5, 15, V(50))
    m.e(4, 4, P(1250, 300), E(1, "FINAL"), S("図枠".encode("shift_jis")))
    m.e(5, 10, IX(2))
    m.e(4, 4, P(1250, 150), E(1, "FINAL"), S("説明文字 テスト".encode("shift_jis")))
    m.e(5, 3, R(0.35))
    m.e(4, 12, P(1050, 1500), V(500))
    m.e(4, 4, P(900, 1500), E(1, "FINAL"), S("日本語の文字".encode("shift_jis")))
    m.e(0, 5)
    m.e(0, 2)
    return m


def main():
    os.makedirs(OUT, exist_ok=True)
    for name, meta in [("shapes", shapes()), ("illustration", illustration()), ("sjis", sjis())]:
        with open(os.path.join(OUT, name + ".cgm"), "wb") as f:
            f.write(Binary().write(meta))
    with open(os.path.join(OUT, "shapes-text.cgm"), "wb") as f:
        f.write(clear_text(shapes()))


if __name__ == "__main__":
    main()
