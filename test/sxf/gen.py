"""Writes the SXF test drawings of converter/sxf into converter/sxf/testdata.

Standard library only. One drawing is written three times: as SXF Ver.3.1
feature comments (shapes.sfc), as STEP AP202 (shapes.p21) and as that P21
file in a zip archive (shapes.p2z). The P21 file is written the way the
SCADEC library writes SXF, which the SXF writers of CAD programs use, quirks
included: text points moved off their anchors, leader arrows turned towards
the leader, angular dimension arcs flagged clockwise (see converter/sxf).
Both files must draw the same (TestSFCMatchesP21). Text uses only
characters of the test fonts (converter/pptx/testdata/fonts).

The drawing, on an A3 sheet on a dark blue background (a background color
attribute): the 15 predefined line types, the predefined colors and the 9
widths, a user-defined color, line type and width; circles, arcs both ways,
ellipses and an elliptic arc, a polyline, a spline and a clothoid, the point
markers; text at the 9 anchors, turned, slanted, spaced and vertical;
compound figures placed turned and scaled, nested, and in geodetic
coordinates; linear, angular, radius and diameter dimensions, a leader and
a balloon; a color fill, hatching with a hole and a blank area; and a hidden
layer, which is not drawn.
"""

import io
import math
import os
import zipfile

OUT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "converter", "sxf", "testdata")

# predefined colors by code (1-16), as SXF names them; ISO names the first 8
COLORS = ["black", "red", "green", "blue", "yellow", "magenta", "cyan", "white",
          "deeppink", "brown", "orange", "lightgreen", "lightblue", "lavender", "lightgray", "darkgray"]
RGB = [(0, 0, 0), (255, 0, 0), (0, 255, 0), (0, 0, 255), (255, 255, 0), (255, 0, 255), (0, 255, 255),
       (255, 255, 255), (192, 0, 128), (192, 128, 64), (255, 128, 0), (128, 192, 128), (0, 128, 255),
       (128, 64, 255), (192, 192, 192), (128, 128, 128)]
LINE_TYPES = ["continuous", "dashed", "dashed spaced", "long dashed dotted", "long dashed double-dotted",
              "long dashed triplicate-dotted", "dotted", "chain", "chain double dash", "dashed dotted",
              "double-dashed dotted", "dashed double-dotted", "double-dashed double-dotted",
              "dashed triplicate-dotted", "double-dashed triplicate-dotted"]
WIDTHS = [0.13, 0.18, 0.25, 0.35, 0.5, 0.7, 1.0, 1.4, 2.0]
ARROWS = ["blanked arrow", "blanked box", "blanked dot", "dimension origin", "filled box", "filled arrow",
          "filled dot", "integral symbol", "open arrow", "slash", "unfilled arrow"]
MARKERS = ["asterisk", "circle", "dot", "plus", "square", "triangle", "x"]

USER_COLOR = (0, 160, 160)  # code 17
USER_TYPE = ("鎖線", [3.0, 1.0, 0.5, 1.0])  # code 17
USER_WIDTH = 0.3  # code 10 + 1
FONT = "ＭＳ ゴシック"
BACKGROUND = "$$ATRU$$1$$背景色$$色$$16_24_48"
LAYERS = [("図枠", True), ("図形", True), ("文字", True), ("寸法", True), ("非表示", False)]
FRAME, SHAPES, TEXT, DIMS, HIDDEN = 1, 2, 3, 4, 5


def num(v):
    return "%.6f" % v if isinstance(v, float) else str(v)


class S(str):
    """A string argument of a feature."""


class SFC:
    """SXF feature comments.

    A composite curve or a compound figure is made of the elements written
    since the previous one, so the definitions come first (in "defs"), then
    the elements of the sheet ("body"): the elements of a figure being
    defined are kept apart until it is.
    """

    def __init__(self):
        self.sections = {"tables": [], "defs": [], "body": [], "tail": []}
        self.where = "tables"
        self.frames = []

    def add(self, name, *args, section=None):
        parts = []
        for a in args:
            if isinstance(a, S):
                parts.append("\\'%s\\'" % a)
            elif isinstance(a, (list, tuple)):
                parts.append("'(%s)'" % ",".join(num(v) if not isinstance(v, tuple) else "(%s)" % ",".join(num(x) for x in v) for v in a))
            else:
                parts.append("'%s'" % num(a))
        line = "%s(%s)" % (name, ",".join(parts))
        if section:
            self.sections[section].append(line)
        elif self.frames:
            self.frames[-1].append(line)
        else:
            self.sections[self.where].append(line)

    def begin(self):
        self.frames.append([])

    def end(self):
        self.sections["defs"] += self.frames.pop()

    def text(self, name):
        head = ("ISO-10303-21;\nHEADER;\nFILE_DESCRIPTION(('SCADEC level2 feature_mode'),\n        '2;1');\n"
                "FILE_NAME('%s',\n        '2026-9-26T00:00:00',\n        (''),\n        (''),\n"
                "        'SCADEC_API_Ver3.30$$3.1',\n        'bdf test/sxf/gen.py',\n        '');\n"
                "FILE_SCHEMA(('ASSOCIATIVE_DRAUGHTING'));\nENDSEC;\nDATA;\n\n" % name)
        lines = [ln for k in ("tables", "defs", "body", "tail") for ln in self.sections[k]]
        body = "".join("/*SXF\n#%d = %s\nSXF*/\n\n" % (10 * (i + 1), ln) for i, ln in enumerate(lines))
        return head + body + "ENDSEC;\nEND-ISO-10303-21;\n"


def step_string(s):
    """A STEP string: non-ASCII characters as \\X2\\ runs."""
    out, run = [], []
    for c in s + "\0":
        if c != "\0" and ord(c) > 0x7E:
            run.append("%04X" % ord(c))
            continue
        if run:
            out.append("\\X2\\" + "".join(run) + "\\X0\\")
            run = []
        if c == "'":
            out.append("''")
        elif c != "\0":
            out.append(c)
    return "'" + "".join(out) + "'"


def f(v):
    return "%.6f" % v


def a(v):
    return "%.14f" % v


class P21:
    """STEP AP202 entity instances."""

    def __init__(self):
        self.out = []
        self.n = 0

    def e(self, text):
        self.n += 10
        self.out.append("#%d=%s;" % (self.n, text))
        return "#%d" % self.n

    def complex(self, *records):
        self.n += 10
        self.out.append("#%d=(\n%s\n);" % (self.n, "\n".join(records)))
        return "#%d" % self.n

    def point(self, x, y):
        return self.e("CARTESIAN_POINT(' ',(%s,%s))" % (f(x), f(y)))

    def direction(self, x, y):
        return self.e("DIRECTION(' ',(%s,%s))" % (a(x), a(y)))

    def placement(self, x, y, angle=None):
        """An AXIS2_PLACEMENT_2D at (x, y) turned by angle (degrees)."""
        p = self.point(x, y)
        if angle is None:
            return self.e("AXIS2_PLACEMENT_2D(' ',%s,$)" % p)
        d = self.direction(math.cos(math.radians(angle)), math.sin(math.radians(angle)))
        return self.e("AXIS2_PLACEMENT_2D(' ',%s,%s)" % (p, d))

    def text(self, name):
        head = ("ISO-10303-21;\nHEADER;\nFILE_DESCRIPTION(('SCADEC level2 AP202_mode'),\n\t\t'2;1');\n"
                "FILE_NAME('%s',\n\t\t   '2026-9-26T00:00:00',\n\t\t   (''),\n\t\t   (''),\n"
                "\t\t   'SCADEC_API_Ver3.30$$3.1',\n\t\t   'bdf test/sxf/gen.py',\n\t\t'');\n"
                "FILE_SCHEMA(('ASSOCIATIVE_DRAUGHTING'));\nENDSEC;\nDATA;\n" % name)
        return head + "\n".join(self.out) + "\nENDSEC;\nEND-ISO-10303-21;\n"


def units(s):
    """The width of a text in half-width units."""
    return sum(2 if len(c.encode("cp932")) == 2 else 1 for c in s)


class Writer:
    """Writes the drawing as SFC and P21 at once."""

    def __init__(self):
        self.sfc = SFC()
        self.p = P21()
        p = self.p
        self.mm = p.complex("LENGTH_UNIT()", "NAMED_UNIT(*)", "SI_UNIT(.MILLI.,.METRE.)")
        rad = p.complex("NAMED_UNIT(*)", "PLANE_ANGLE_UNIT()", "SI_UNIT($,.RADIAN.)")
        self.context = p.complex("GEOMETRIC_REPRESENTATION_CONTEXT(2)", "GLOBAL_UNIT_ASSIGNED_CONTEXT((%s,%s))" % (rad, self.mm),
                                 "REPRESENTATION_CONTEXT('ID1','2D')")
        self.source = p.e("EXTERNAL_SOURCE(IDENTIFIER('scadec'))")
        self.null_style = p.e("PRESENTATION_STYLE_ASSIGNMENT((NULL_STYLE(.NULL.)))")
        self.colours, self.fonts, self.widths = {}, {}, {}
        # the tables
        for code in range(2, 17):
            self.sfc.add("pre_defined_colour_feature", S(COLORS[code - 1]))
            name = COLORS[code - 1]
            if code <= 8:
                self.colours[code] = p.e("DRAUGHTING_PRE_DEFINED_COLOUR('%s')" % name)
            else:
                r, g, b = RGB[code - 1]
                self.colours[code] = p.e("COLOUR_RGB('$$SXF_%s',%s,%s,%s)" % (name, f(r / 255), f(g / 255), f(b / 255)))
        self.sfc.add("user_defined_colour_feature", *USER_COLOR)
        self.colours[17] = p.e("COLOUR_RGB(' ',%s,%s,%s)" % tuple(f(v / 255) for v in USER_COLOR))
        for code, name in enumerate(LINE_TYPES, 1):
            self.sfc.add("pre_defined_font_feature", S(name))
            self.fonts[code] = p.e("DRAUGHTING_PRE_DEFINED_CURVE_FONT('%s')" % name)
        name, pitch = USER_TYPE
        self.sfc.add("user_defined_font_feature", S(name), len(pitch), [float(v) for v in pitch])
        pats = [p.e("CURVE_STYLE_FONT_PATTERN(%s,%s)" % (f(pitch[i]), f(pitch[i + 1]))) for i in range(0, len(pitch), 2)]
        self.fonts[17] = p.e("CURVE_STYLE_FONT(%s,(%s))" % (step_string(name), ",".join(pats)))
        for code, w in enumerate(WIDTHS + [USER_WIDTH], 1):
            self.sfc.add("width_feature", float(w))
            self.widths[code if code <= 9 else 11] = p.e("LENGTH_MEASURE_WITH_UNIT(POSITIVE_LENGTH_MEASURE(%s),%s)" % (f(w), self.mm))
        self.sfc.add("text_font_feature", S(FONT))
        self.text_font = p.e("EXTERNALLY_DEFINED_TEXT_FONT(IDENTIFIER(%s),%s)" % (step_string(FONT), self.source))
        self.sfc.where = "body"
        self.layer_items = {i: [] for i in range(len(LAYERS) + 1)}
        # the elements of the sheet, or of the compound figure being defined
        self.items = []
        self.frames = []
        self.curves = []  # composite curves by code: P21 curve, style, visible

    # frames: the members of a compound figure are gathered until it is
    # defined
    def begin(self):
        self.frames.append(self.items)
        self.items = []
        self.sfc.begin()

    def end(self):
        items = self.items
        self.items = self.frames.pop()
        self.sfc.end()
        return items

    def occurrence(self, layer, kind, style, item, name="' '"):
        rec = {"curve": ["ANNOTATION_CURVE_OCCURRENCE()"], "text": ["ANNOTATION_TEXT_OCCURRENCE()"],
               "symbol": ["ANNOTATION_SYMBOL_OCCURRENCE()"], "fill": ["ANNOTATION_FILL_AREA_OCCURRENCE(%s)" % self.p.point(0.0, 0.0)],
               "dim": ["ANNOTATION_CURVE_OCCURRENCE()", "DIMENSION_CURVE()"],
               "projection": ["ANNOTATION_CURVE_OCCURRENCE()", "PROJECTION_CURVE()"],
               "leader": ["ANNOTATION_CURVE_OCCURRENCE()", "LEADER_CURVE()"]}[kind]
        records = rec + ["ANNOTATION_OCCURRENCE()", "DRAUGHTING_ANNOTATION_OCCURRENCE()", "GEOMETRIC_REPRESENTATION_ITEM()",
                         "REPRESENTATION_ITEM(%s)" % name, "STYLED_ITEM((%s),%s)" % (style, item)]
        occ = self.p.complex(*sorted(records))
        if layer is not None:
            self.layer_items[layer].append(occ)
        return occ

    def curve_style(self, color, typ, width):
        cs = self.p.e("CURVE_STYLE(' ',%s,%s,%s)" % (self.fonts[typ], self.widths[width], self.colours[color]))
        return self.p.e("PRESENTATION_STYLE_ASSIGNMENT((%s))" % cs)

    def curve(self, layer, color, typ, width, curve, kind="curve"):
        occ = self.occurrence(layer, kind, self.curve_style(color, typ, width), curve)
        return occ

    # geometry, in P21
    def line_curve(self, x0, y0, x1, y1):
        p, q = self.p.point(x0, y0), self.p.point(x1, y1)
        d = self.p.e("DIRECTION(' ',(%s,%s))" % (f(x1 - x0), f(y1 - y0)))
        v = self.p.e("VECTOR(' ',%s,1.000000)" % d)
        ln = self.p.e("LINE(' ',%s,%s)" % (p, v))
        return self.p.e("TRIMMED_CURVE(' ',%s,(%s),(%s),.T.,.CARTESIAN.)" % (ln, p, q))

    def polyline_curve(self, pts):
        ids = [self.p.point(x, y) for x, y in pts]
        return self.p.e("POLYLINE(' ',(%s))" % ",".join(ids))

    def arc_curve(self, cx, cy, rx, ry, rot, a0, a1, sense):
        pl = self.p.placement(cx, cy, rot if rx != ry else None)
        if rx == ry:
            basis = self.p.e("CIRCLE(' ',%s,%s)" % (pl, f(rx)))
        else:
            basis = self.p.e("ELLIPSE(' ',%s,%s,%s)" % (pl, f(rx), f(ry)))
        return self.p.e("TRIMMED_CURVE(' ',%s,(PARAMETER_VALUE(%s)),(PARAMETER_VALUE(%s)),%s,.PARAMETER.)"
                        % (basis, a(math.radians(a0)), a(math.radians(a1)), ".T." if sense else ".F."))

    # features
    def line(self, layer, color, typ, width, x0, y0, x1, y1):
        self.sfc.add("line_feature", layer, color, typ, width, x0, y0, x1, y1)
        self.items.append(self.curve(layer, color, typ, width, self.line_curve(x0, y0, x1, y1)))

    def polyline(self, layer, color, typ, width, pts):
        self.sfc.add("polyline_feature", layer, color, typ, width, len(pts), [x for x, _ in pts], [y for _, y in pts])
        self.items.append(self.curve(layer, color, typ, width, self.polyline_curve(pts)))

    def circle(self, layer, color, typ, width, cx, cy, r):
        self.sfc.add("circle_feature", layer, color, typ, width, cx, cy, r)
        c = self.p.e("CIRCLE(' ',%s,%s)" % (self.p.placement(cx, cy), f(r)))
        self.items.append(self.curve(layer, color, typ, width, c))

    def arc(self, layer, color, typ, width, cx, cy, r, cw, a0, a1):
        self.sfc.add("arc_feature", layer, color, typ, width, cx, cy, r, 1 if cw else 0, float(a0), float(a1))
        self.items.append(self.curve(layer, color, typ, width, self.arc_curve(cx, cy, r, r, 0, a0, a1, not cw)))

    def ellipse(self, layer, color, typ, width, cx, cy, rx, ry, rot):
        self.sfc.add("ellipse_feature", layer, color, typ, width, cx, cy, rx, ry, float(rot))
        e = self.p.e("ELLIPSE(' ',%s,%s,%s)" % (self.p.placement(cx, cy, rot), f(rx), f(ry)))
        self.items.append(self.curve(layer, color, typ, width, e))

    def ellipse_arc(self, layer, color, typ, width, cx, cy, rx, ry, cw, rot, a0, a1):
        self.sfc.add("ellipse_arc_feature", layer, color, typ, width, cx, cy, rx, ry, 1 if cw else 0, float(rot), float(a0), float(a1))
        self.items.append(self.curve(layer, color, typ, width, self.arc_curve(cx, cy, rx, ry, rot, a0, a1, not cw)))

    def spline(self, layer, color, typ, width, pts):
        self.sfc.add("spline_feature", layer, color, typ, width, 0, len(pts), [x for x, _ in pts], [y for _, y in pts])
        ids = [self.p.point(x, y) for x, y in pts]
        b = self.p.e("BEZIER_CURVE(' ',3,(%s),.UNSPECIFIED.,.F.,.U.)" % ",".join(ids))
        self.items.append(self.curve(layer, color, typ, width, b))

    def clothoid(self, layer, color, typ, width, x, y, param, cw, angle, s0, s1):
        self.sfc.add("clothoid_feature", layer, color, typ, width, x, y, param, 1 if cw else 0, float(angle), s0, s1)
        c = self.p.e("CLOTHOID('%d',%s,%s)" % (1 if cw else 0, self.p.placement(x, y, angle), f(param)))
        t = self.p.e("TRIMMED_CURVE('%s,%s',%s,(PARAMETER_VALUE((%s))),(PARAMETER_VALUE((%s))),.T.,.PARAMETER.)"
                     % (f(s0), f(s1), c, f(s0), f(s1)))
        self.items.append(self.curve(layer, color, typ, width, t))

    def symbol(self, layer, color, definition, x, y, dx, dy, scale, kind=None):
        d = self.p.direction(dx, dy)
        pl = self.p.e("AXIS2_PLACEMENT_2D(' ',%s,%s)" % (self.p.point(x, y), d))
        target = self.p.e("SYMBOL_TARGET(' ',%s,%s,%s)" % (pl, a(scale), a(scale)))
        ds = self.p.e("DEFINED_SYMBOL(' ',%s,%s)" % (definition, target))
        sc = self.p.e("SYMBOL_COLOUR(%s)" % self.colours[color])
        ss = self.p.e("SYMBOL_STYLE(' ',%s)" % sc)
        style = self.p.e("PRESENTATION_STYLE_ASSIGNMENT((%s))" % ss)
        records = ["ANNOTATION_OCCURRENCE()", "ANNOTATION_SYMBOL_OCCURRENCE()", "DRAUGHTING_ANNOTATION_OCCURRENCE()",
                   "GEOMETRIC_REPRESENTATION_ITEM()", "REPRESENTATION_ITEM(' ')", "STYLED_ITEM((%s),%s)" % (style, ds)]
        if kind:
            records += kind
        occ = self.p.complex(*sorted(records))
        if layer is not None:
            self.layer_items[layer].append(occ)
        return occ

    def marker(self, layer, color, x, y, code, angle, scale):
        self.sfc.add("point_marker_feature", layer, color, x, y, code, float(angle), float(scale))
        d = self.p.e("PRE_DEFINED_POINT_MARKER_SYMBOL('%s')" % MARKERS[code - 1])
        t = math.radians(angle)
        self.items.append(self.symbol(layer, color, d, x, y, math.cos(t), math.sin(t), scale))

    def text_args(self, font, s, x, y, h, w, spc, angle, slant, b, direct):
        return (font, S(s), x, y, h, w, spc, float(angle), float(slant), b, direct)

    def p21_text(self, layer, color, s, x, y, h, w, spc, angle, slant, b, direct, name="' '"):
        row, col = (b - 1) // 3, (b - 1) % 3
        # SCADEC moves the point: up half the height on the bottom row, down
        # half of it on the middle row and all of it on the top row
        k = [0.5, -0.5, -1.0][row] * h
        t = math.radians(angle)
        px, py = x - math.sin(t) * k, y + math.cos(t) * k
        ext = self.p.e("PLANAR_EXTENT(' ',%s,%s)" % (f(w), f(h)))
        pl = self.p.e("AXIS2_PLACEMENT_2D(' ',%s,%s)" % (self.p.point(px, py), self.p.direction(math.cos(t), math.sin(t))))
        cols = ["left", "centre", "right"]
        lit = self.p.e("TEXT_LITERAL_WITH_EXTENT('$$SXF_%s %s',%s,%s,'baseline %s',%s,%s,%s)"
                       % (["baseline", "middleline", "topline"][row], cols[col], step_string(s), pl, cols[col],
                          ".DOWN." if direct == 2 else ".RIGHT.", self.text_font, ext))
        fd = self.p.e("TEXT_STYLE_FOR_DEFINED_FONT(%s)" % self.colours[color])
        n = max(len(s), 1)
        ts = self.p.complex("TEXT_STYLE(' ',%s)" % fd,
                            "TEXT_STYLE_WITH_BOX_CHARACTERISTICS((BOX_HEIGHT(%s),BOX_WIDTH(%s),BOX_SLANT_ANGLE(%s),BOX_ROTATE_ANGLE(%s)))"
                            % (f(h), f(w / n), a(math.radians(slant)), a(0.0)),
                            "TEXT_STYLE_WITH_SPACING(LENGTH_MEASURE(%s))" % f(spc))
        style = self.p.e("PRESENTATION_STYLE_ASSIGNMENT((%s))" % ts)
        return self.occurrence(layer, "text", style, lit, name)

    def text(self, layer, color, s, x, y, h, spc=0.0, angle=0, slant=0, b=1, direct=1, w=None):
        if w is None:
            w = units(s) * h / 2 + spc * (len(s) - 1)
        self.sfc.add("text_string_feature", layer, color, *self.text_args(1, s, x, y, h, w, spc, angle, slant, b, direct))
        self.items.append(self.p21_text(layer, color, s, x, y, h, w, spc, angle, slant, b, direct))

    def arrow(self, layer, color, code, x, y, dx, dy, scale, kind):
        d = self.p.e("PRE_DEFINED_TERMINATOR_SYMBOL('%s')" % ARROWS[code - 1])
        return self.symbol(None, color, d, x, y, dx, dy, scale, kind)

    def callout(self, layer, kinds, members):
        c = self.p.complex(*sorted(["DRAUGHTING_CALLOUT((%s))" % ",".join(members), "DRAUGHTING_ELEMENTS()",
                                    "GEOMETRIC_REPRESENTATION_ITEM()", "REPRESENTATION_ITEM(' ')"] + kinds))
        si = self.p.e("STYLED_ITEM(' ',(%s),%s)" % (self.null_style, c))
        self.layer_items[layer].append(si)
        self.items.append(si)

    def linear_dim(self, layer, color, typ, width, a, b, ext1, ext2, arr1, arr2, text):
        """ext: (base, start, end) or None; arr: (code, dir, x, y, scale); text: (s, x, y, h, b)."""
        args = [layer, color, typ, width, *a, *b]
        for e in (ext1, ext2):
            args += [1, *e[0], *e[1], *e[2]] if e else [0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
        for r in (arr1, arr2):
            args += [r[0], r[1], r[2], r[3], float(r[4])]
        s, tx, ty, h, tb = text
        w = units(s) * h / 2
        args += [1, *self.text_args(1, s, tx, ty, h, w, 0.0, 0, 0, tb, 1)]
        self.sfc.add("linear_dim_feature", *args)
        style = self.curve_style(color, typ, width)
        members = [self.occurrence(None, "dim", style, self.line_curve(*a, *b))]
        for i, e in enumerate((ext1, ext2), 1):
            if e:
                members.append(self.occurrence(None, "projection", style, self.line_curve(*e[1], *e[2]), "'$$SXF_prj_%d'" % i))
        members += self.dim_arrows(color, [arr1, arr2], a, b, members[0])
        members.append(self.p21_text(None, color, s, tx, ty, h, w, 0.0, 0, 0, tb, 1, "'dimension value'"))
        self.callout(layer, ["DIMENSION_CURVE_DIRECTED_CALLOUT()", "LINEAR_DIMENSION()"], members)

    def dim_arrows(self, color, arrows, a, b, curve):
        out = []
        for i, (code, d, x, y, scale) in enumerate(arrows):
            # pointing away from the farther end of the line (d 1) or at it (2)
            other = a if math.hypot(x - a[0], y - a[1]) >= math.hypot(x - b[0], y - b[1]) else b
            dx, dy = x - other[0], y - other[1]
            if d == 2:
                dx, dy = -dx, -dy
            l = math.hypot(dx, dy)
            out.append(self.arrow(None, color, code, x, y, dx / l, dy / l, scale,
                                  ["DIMENSION_CURVE_TERMINATOR(%s)" % (".ORIGIN." if i == 0 else ".TARGET."), "TERMINATOR_SYMBOL(%s)" % curve]))
        return out

    def angular_dim(self, layer, color, typ, width, ctr, r, a0, a1, ext1, ext2, arr1, arr2, text):
        args = [layer, color, typ, width, *ctr, r, float(a0), float(a1)]
        for e in (ext1, ext2):
            args += [1, *e[0], *e[1], *e[2]] if e else [0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
        for q in (arr1, arr2):
            args += [q[0], q[1], q[2], q[3], float(q[4])]
        s, tx, ty, h, angle, tb = text
        w = units(s) * h / 2
        args += [1, *self.text_args(1, s, tx, ty, h, w, 0.0, angle, 0, tb, 1)]
        self.sfc.add("angular_dim_feature", *args)
        style = self.curve_style(color, typ, width)
        # SCADEC flags the arc clockwise, though it runs counter-clockwise
        members = [self.occurrence(None, "dim", style, self.arc_curve(*ctr, r, r, 0, a0, a1, False))]
        for i, e in enumerate((ext1, ext2), 1):
            if e:
                members.append(self.occurrence(None, "projection", style, self.line_curve(*e[1], *e[2]), "'$$SXF_prj_%d'" % i))
        for k, (code, d, x, y, scale) in enumerate((arr1, arr2)):
            t = math.atan2(y - ctr[1], x - ctr[0])
            dx, dy = math.sin(t), -math.cos(t)
            if k == 1:
                dx, dy = -dx, -dy
            if d == 2:
                dx, dy = -dx, -dy
            members.append(self.arrow(None, color, code, x, y, dx, dy, scale,
                                      ["DIMENSION_CURVE_TERMINATOR(%s)" % (".ORIGIN." if k == 0 else ".TARGET."), "TERMINATOR_SYMBOL(%s)" % members[0]]))
        members.append(self.p21_text(None, color, s, tx, ty, h, w, 0.0, angle, 0, tb, 1, "'dimension value'"))
        self.callout(layer, ["ANGULAR_DIMENSION()", "DIMENSION_CURVE_DIRECTED_CALLOUT()"], members)

    def radius_dim(self, layer, color, typ, width, a, b, arr, text, diameter=None):
        s, tx, ty, h, angle, tb = text
        w = units(s) * h / 2
        args = [layer, color, typ, width, *a, *b, arr[0], arr[1], arr[2], arr[3], float(arr[4])]
        if diameter:
            args += [diameter[0], diameter[1], diameter[2], diameter[3], float(diameter[4])]
        args += [1, *self.text_args(1, s, tx, ty, h, w, 0.0, angle, 0, tb, 1)]
        self.sfc.add("diameter_dim_feature" if diameter else "radius_dim_feature", *args)
        style = self.curve_style(color, typ, width)
        members = [self.occurrence(None, "dim", style, self.line_curve(*a, *b))]
        members += self.dim_arrows(color, [arr] + ([diameter] if diameter else []), a, b, members[0])
        members.append(self.p21_text(None, color, s, tx, ty, h, w, 0.0, angle, 0, tb, 1, "'dimension value'"))
        self.callout(layer, ["DIAMETER_DIMENSION()" if diameter else "RADIUS_DIMENSION()", "DIMENSION_CURVE_DIRECTED_CALLOUT()"], members)

    def label(self, layer, color, typ, width, pts, code, scale, text, balloon=None):
        s, tx, ty, h, tb = text
        w = units(s) * h / 2
        args = [layer, color, typ, width, len(pts), [x for x, _ in pts], [y for _, y in pts]]
        if balloon:
            args += list(balloon)
        args += [code, float(scale), 1, *self.text_args(1, s, tx, ty, h, w, 0.0, 0, 0, tb, 1)]
        self.sfc.add("balloon_feature" if balloon else "label_feature", *args)
        style = self.curve_style(color, typ, width)
        leader = self.occurrence(None, "leader", style, self.polyline_curve(pts))
        members = [leader]
        if balloon:
            c = self.p.e("CIRCLE(' ',%s,%s)" % (self.p.placement(balloon[0], balloon[1]), f(balloon[2])))
            members.append(self.occurrence(None, "curve", style, c))
        # SCADEC turns the arrow towards the leader
        (x0, y0), (x1, y1) = pts[0], pts[1]
        l = math.hypot(x1 - x0, y1 - y0)
        members.append(self.arrow(None, color, code, x0, y0, (x1 - x0) / l, (y1 - y0) / l, scale,
                                  ["LEADER_TERMINATOR()", "TERMINATOR_SYMBOL(%s)" % leader]))
        members.append(self.p21_text(None, color, s, tx, ty, h, w, 0.0, 0, 0, tb, 1, "'dimension value'"))
        self.callout(layer, ["LEADER_DIRECTED_CALLOUT()"], members)

    def composite(self, color, typ, width, visible, segments):
        """A composite curve of (kind, args) segments: its code."""
        segs = []
        for kind, args in segments:
            if kind == "polyline":
                self.sfc.add("polyline_feature", SHAPES, color, typ, width, len(args), [x for x, _ in args], [y for _, y in args], section="defs")
                c = self.polyline_curve(args)
            elif kind == "arc":
                cx, cy, r, cw, a0, a1 = args
                self.sfc.add("arc_feature", SHAPES, color, typ, width, cx, cy, r, 1 if cw else 0, float(a0), float(a1), section="defs")
                c = self.arc_curve(cx, cy, r, r, 0, a0, a1, not cw)
            segs.append(self.p.e("COMPOSITE_CURVE_SEGMENT(.CONTINUOUS.,.T.,%s)" % c))
        self.sfc.add("composite_curve_org_feature", color, typ, width, 1 if visible else 0, section="defs")
        cc = self.p.e("COMPOSITE_CURVE(' ',(%s),.T.)" % ",".join(segs))
        self.curves.append((cc, self.curve_style(color, typ, width), visible))
        return len(self.curves)

    def fill(self, layer, name, out, holes, sfc_args, styles):
        self.sfc.add(name, layer, *sfc_args, out, len(holes), list(holes))
        area = self.p.e("ANNOTATION_FILL_AREA(' ',(%s))" % ",".join(self.curves[c - 1][0] for c in [out] + list(holes)))
        fs = self.p.e("FILL_AREA_STYLE(' ',(%s))" % ",".join(styles))
        style = self.p.e("PRESENTATION_STYLE_ASSIGNMENT((%s))" % fs)
        self.items.append(self.occurrence(layer, "fill", style, area))
        # the outlines that show
        for c in [out] + list(holes):
            cc, cs, visible = self.curves[c - 1]
            if visible:
                self.items.append(self.occurrence(layer, "curve", cs, cc))

    def colour_fill(self, layer, color, out, holes=()):
        st = self.p.e("FILL_AREA_STYLE_COLOUR(' ',%s)" % self.colours[color])
        self.fill(layer, "fill_area_style_colour_feature", out, holes, [color], [st])

    def hatching(self, layer, lines, out, holes=()):
        styles = []
        for color, typ, width, x, y, spacing, angle in lines:
            cs = self.p.e("CURVE_STYLE(' ',%s,%s,%s)" % (self.fonts[typ], self.widths[width], self.colours[color]))
            t = math.radians(angle)
            d = self.p.direction(-math.sin(t), math.cos(t))
            v = self.p.e("VECTOR(' ',%s,%s)" % (d, f(spacing)))
            rf = self.p.e("ONE_DIRECTION_REPEAT_FACTOR(' ',%s)" % v)
            pt = self.p.point(x, y)
            styles.append(self.p.e("FILL_AREA_STYLE_HATCHING(' ',%s,%s,%s,%s,%s)" % (cs, rf, pt, pt, a(t))))
        self.fill(layer, "fill_area_style_hatching_feature", out, holes,
                  [len(lines)] + [tuple(float(v) if i >= 3 else v for i, v in enumerate(ln)) for ln in lines], styles)

    def blank(self, layer, out):
        st = self.p.e("EXTERNALLY_DEFINED_HATCH_STYLE(IDENTIFIER('Area_control'),%s,' ')" % self.source)
        self.fill(layer, "externally_defined_hatch_feature", out, (), [S("Area_control")], [st])

    def subfigure(self, name, flag, items):
        """Defines a compound figure of the items; its P21 map."""
        self.sfc.add("sfig_org_feature", S(name), flag, section="defs")
        prefix = {1: "$$SXF_FM_", 2: "$$SXF_FG_", 3: "$$SXF_G_", 4: "$$SXF_P_"}[flag]
        origin = self.p.placement(0.0, 0.0)
        rep = self.p.e("DRAUGHTING_SUBFIGURE_REPRESENTATION(%s,(%s),%s)" % (step_string(prefix + name), ",".join(items + [origin]), self.context))
        return self.p.e("SYMBOL_REPRESENTATION_MAP(%s,%s)" % (origin, rep)), prefix + name

    def locate(self, figure, name, x, y, angle, sx, sy):
        self.sfc.add("sfig_locate_feature", 0, S(name), x, y, float(angle), float(sx), float(sy))
        mapping, full = figure
        target = self.p.e("SYMBOL_TARGET(' ',%s,%s,%s)" % (self.p.placement(x, y, angle), a(sx), a(sy)))
        mi = self.p.complex("ANNOTATION_SYMBOL()", "GEOMETRIC_REPRESENTATION_ITEM()", "MAPPED_ITEM(%s,%s)" % (mapping, target), "REPRESENTATION_ITEM(' ')")
        occ = self.p.complex(*sorted(["ANNOTATION_OCCURRENCE()", "ANNOTATION_SUBFIGURE_OCCURRENCE()", "ANNOTATION_SYMBOL_OCCURRENCE()",
                                      "DRAUGHTING_ANNOTATION_OCCURRENCE()", "GEOMETRIC_REPRESENTATION_ITEM()",
                                      "REPRESENTATION_ITEM(%s)" % step_string(full), "STYLED_ITEM((%s),%s)" % (self.null_style, mi)]))
        self.layer_items[0].append(occ)
        self.items.append(occ)

    def finish(self, sheet, title):
        self.sfc.where = "tail"
        for name, visible in LAYERS:
            self.sfc.add("layer_feature", S(name), 1 if visible else 0)
        self.sfc.add("drawing_attribute_feature", S(" "), S(" "), S(" "), S(title), S(" "), S(" "), S(" "), 0, 1, 1, S(" "), S(" "))
        # A3 (code 3), landscape
        self.sfc.add("drawing_sheet_feature", S(sheet), 3, 1, 420, 297)
        p = self.p
        box = p.e("PLANAR_BOX(' ',420.000000,297.000000,%s)" % p.placement(0.0, 0.0))
        definition = p.e("DRAWING_DEFINITION(' ',' ')")
        rev = p.e("DRAUGHTING_DRAWING_REVISION('01',%s,' ')" % definition)
        p.e("DRAUGHTING_TITLE((%s),'JAPANESE',%s)" % (rev, step_string(title)))
        sh = p.e("DRAWING_SHEET_REVISION(%s,(%s),%s,'01')" % (step_string(sheet), ",".join(self.items + [box]), self.context))
        p.e("DRAWING_SHEET_REVISION_USAGE(%s,%s,'01')" % (sh, rev))
        p.e("PRESENTATION_SIZE(%s,%s)" % (sh, box))
        hidden = []
        for i, (name, visible) in enumerate([("$$SXF_dummy_layer_for_subfigure", True)] + LAYERS):
            items = self.layer_items[i]
            if not items:
                continue
            la = p.e("PRESENTATION_LAYER_ASSIGNMENT(%s,' ',(%s))" % (step_string(name), ",".join(items)))
            p.e("PRESENTATION_LAYER_USAGE(%s,%s)" % (la, sh))
            if not visible:
                hidden.append(la)
        if hidden:
            p.e("INVISIBILITY((%s))" % ",".join(hidden))


def drawing():
    w = Writer()
    # the background color, carried by an empty group as SXF writers do
    w.begin()
    bg = w.composite(8, 1, 1, False, [("polyline", [(-30.0, -30.0), (-10.0, -30.0), (-10.0, -10.0), (-30.0, -10.0), (-30.0, -30.0)])])
    w.blank(FRAME, bg)
    w.locate(w.subfigure(BACKGROUND, 3, w.end()), BACKGROUND, 0.0, 0.0, 0, 1, 1)

    # the frame and the title
    w.polyline(FRAME, 8, 1, 5, [(10.0, 10.0), (410.0, 10.0), (410.0, 287.0), (10.0, 287.0), (10.0, 10.0)])
    w.text(TEXT, 8, "SXF テスト図", 210.0, 275.0, 7.0, b=5)

    # the line types in the predefined colors, and the user-defined ones
    colors = [8, 2, 3, 4, 5, 6, 7, 9, 10, 11, 12, 13, 14, 15, 16]
    for i in range(15):
        w.line(SHAPES, colors[i], i + 1, 3, 20.0, 250.0 - 6 * i, 110.0, 250.0 - 6 * i)
    w.line(SHAPES, 17, 17, 11, 20.0, 160.0, 110.0, 160.0)
    # the widths
    for i in range(9):
        w.line(SHAPES, 8, 1, i + 1, 125.0, 250.0 - 6 * i, 200.0, 250.0 - 6 * i)

    # curves
    w.circle(SHAPES, 7, 1, 3, 240.0, 240.0, 15.0)
    w.arc(SHAPES, 2, 1, 4, 280.0, 240.0, 15.0, False, 30, 150)
    w.arc(SHAPES, 3, 2, 3, 280.0, 235.0, 8.0, True, 30, 150)
    w.ellipse(SHAPES, 5, 1, 3, 330.0, 240.0, 20.0, 10.0, 30)
    w.ellipse_arc(SHAPES, 6, 4, 3, 380.0, 240.0, 20.0, 10.0, False, 0, 0, 200)
    w.polyline(SHAPES, 9, 1, 3, [(220.0, 200.0), (230.0, 215.0), (240.0, 200.0), (250.0, 215.0), (260.0, 200.0)])
    w.spline(SHAPES, 11, 1, 3, [(270.0, 200.0), (280.0, 220.0), (290.0, 190.0), (300.0, 205.0), (310.0, 220.0), (320.0, 190.0), (330.0, 210.0)])
    w.clothoid(SHAPES, 12, 1, 3, 345.0, 195.0, 30.0, False, 10, 0.0, 45.0)
    for i in range(7):
        w.marker(SHAPES, 13, 225.0 + 20 * i, 178.0, i + 1, 15 if i == 6 else 0, 1.2)

    # text at the nine anchors, each on a cross
    for b in range(1, 10):
        x, y = 30.0 + 45 * ((b - 1) % 3), 110.0 + 16 * ((b - 1) // 3)
        w.line(TEXT, 16, 1, 1, x - 3, y, x + 3, y)
        w.line(TEXT, 16, 1, 1, x, y - 3, x, y + 3)
        w.text(TEXT, 8, "文字%d" % b, x, y, 4.0, b=b)
    w.text(TEXT, 5, "カイテン 30°", 20.0, 90.0, 5.0, angle=30)
    w.text(TEXT, 7, "スラント", 80.0, 95.0, 5.0, slant=15)
    w.text(TEXT, 9, "字間 アリ", 125.0, 95.0, 4.0, spc=1.5)
    w.text(TEXT, 14, "縦書き", 185.0, 150.0, 20.0, b=7, direct=2, w=6.0)

    # compound figures: a part placed twice, a group that places the part,
    # and a part in geodetic coordinates (x north)
    w.begin()
    w.polyline(SHAPES, 7, 1, 3, [(-5.0, -5.0), (5.0, -5.0), (5.0, 5.0), (-5.0, 5.0), (-5.0, -5.0)])
    w.circle(SHAPES, 2, 1, 3, 0.0, 0.0, 3.0)
    w.text(TEXT, 8, "A", 0.0, -6.0, 3.0, b=8)
    part = w.subfigure("部品A", 1, w.end())
    w.locate(part, "部品A", 240.0, 125.0, 30, 1.5, 1.5)
    w.locate(part, "部品A", 270.0, 125.0, 0, 1, 1)
    w.begin()
    w.locate(part, "部品A", 0.0, 0.0, 0, 0.5, 0.5)
    w.line(SHAPES, 12, 1, 3, -8.0, -8.0, 8.0, 8.0)
    group = w.subfigure("組", 3, w.end())
    w.locate(group, "組", 300.0, 125.0, 0, 2, 2)
    w.begin()
    w.polyline(SHAPES, 11, 1, 4, [(0.0, 0.0), (12.0, 0.0), (9.0, 2.0)])
    north = w.subfigure("方位", 2, w.end())
    w.locate(north, "方位", 340.0, 115.0, 0, 1, 1)

    # dimensions
    w.linear_dim(DIMS, 8, 1, 1, (250.0, 60.0), (310.0, 60.0),
                 ((250.0, 45.0), (250.0, 46.0), (250.0, 61.0)), ((310.0, 45.0), (310.0, 46.0), (310.0, 61.0)),
                 (6, 1, 250.0, 60.0, 1), (6, 1, 310.0, 60.0, 1), ("60", 280.0, 61.0, 3.5, 2))
    w.angular_dim(DIMS, 8, 1, 1, (325.0, 20.0), 22.0, 0, 60,
                  ((325.0, 20.0), (327.0, 20.0), (349.0, 20.0)), ((325.0, 20.0), (326.0, 21.7320508), (337.0, 40.7846097)),
                  (9, 1, 347.0, 20.0, 1), (9, 1, 336.0, 39.0525589, 1), ("60°", 340.0, 30.0, 3.5, 0, 1))
    w.circle(SHAPES, 8, 1, 3, 378.0, 62.0, 12.0)
    w.radius_dim(DIMS, 8, 1, 1, (378.0, 62.0), (386.4852814, 70.4852814), (6, 1, 386.4852814, 70.4852814, 1),
                 ("R12", 380.0, 68.0, 3.0, 45, 1))
    w.circle(SHAPES, 8, 1, 3, 378.0, 30.0, 10.0)
    w.radius_dim(DIMS, 8, 1, 1, (370.9289322, 22.9289322), (385.0710678, 37.0710678), (6, 1, 370.9289322, 22.9289322, 1),
                 ("20", 374.0, 32.0, 3.0, 45, 1), diameter=(6, 1, 385.0710678, 37.0710678, 1))
    w.label(DIMS, 8, 1, 1, [(200.0, 50.0), (212.0, 62.0), (237.0, 62.0)], 7, 1, ("ラベル", 213.0, 63.0, 3.5, 1))
    w.label(DIMS, 8, 1, 1, [(200.0, 25.0), (212.0, 32.0)], 6, 1, ("1", 217.0, 34.0, 4.0, 5), balloon=(217.0, 34.0, 5.0))

    # a color fill, hatching around a hole, and a blank area in a fill
    c1 = w.composite(8, 1, 3, True, [("polyline", [(20.0, 20.0), (60.0, 20.0), (60.0, 60.0), (20.0, 60.0), (20.0, 20.0)])])
    w.colour_fill(SHAPES, 17, c1)
    c2 = w.composite(8, 1, 3, True, [("polyline", [(70.0, 20.0), (120.0, 20.0), (120.0, 60.0)]),
                                     ("arc", (95.0, 60.0, 25.0, False, 0, 180)),
                                     ("polyline", [(70.0, 60.0), (70.0, 20.0)])])
    c3 = w.composite(8, 1, 3, True, [("arc", (95.0, 40.0, 8.0, False, 0, 360))])
    w.hatching(SHAPES, [(3, 1, 1, 70.0, 20.0, 4.0, 45), (6, 2, 1, 70.0, 20.0, 6.0, 135)], c2, [c3])
    c4 = w.composite(8, 1, 3, False, [("polyline", [(130.0, 20.0), (190.0, 20.0), (160.0, 75.0), (130.0, 20.0)])])
    w.colour_fill(SHAPES, 10, c4)
    c5 = w.composite(8, 1, 3, False, [("polyline", [(150.0, 30.0), (170.0, 30.0), (170.0, 45.0), (150.0, 45.0), (150.0, 30.0)])])
    w.blank(SHAPES, c5)

    # on the hidden layer
    w.line(HIDDEN, 2, 1, 7, 10.0, 10.0, 410.0, 287.0)
    w.text(HIDDEN, 2, "ヒヒョウジ", 200.0, 150.0, 10.0)

    w.finish("図-1", "SXF テスト図")
    return w


def main():
    os.makedirs(OUT, exist_ok=True)
    w = drawing()
    sfc = w.sfc.text("shapes.sfc").encode("cp932")
    p21 = w.p.text("shapes.p21").encode("cp932")
    with open(os.path.join(OUT, "shapes.sfc"), "wb") as fp:
        fp.write(sfc)
    with open(os.path.join(OUT, "shapes.p21"), "wb") as fp:
        fp.write(p21)
    buf = io.BytesIO()
    with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as z:
        info = zipfile.ZipInfo("shapes.p21", date_time=(2026, 9, 26, 0, 0, 0))
        info.compress_type = zipfile.ZIP_DEFLATED
        z.writestr(info, p21)
    with open(os.path.join(OUT, "shapes.p2z"), "wb") as fp:
        fp.write(buf.getvalue())


if __name__ == "__main__":
    main()
