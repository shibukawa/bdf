#!/usr/bin/env python3
"""Writes the test data of converter/gerber (standard library only):

- board.zip: the fabrication files of a small two-layer board, as KiCad 7
  plots them (Gerber X2 files named after KiCad's layers, PTH and NPTH drill
  files, and a job file), in a folder of the archive.
- features.gbr: one Gerber file with the apertures, macro primitives,
  interpolations, regions, polarities, step and repeat, block apertures and
  aperture transformations of the Gerber specification.

Run it from anywhere: npm run test:gerber:gen
"""
import io
import json
import math
import os
import zipfile

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "gerber", "testdata")
DATE = "2026-09-27T12:00:00+09:00"


def num(v):
    """A coordinate in the 4.6 format, leading zeros left out."""
    return str(int(round(v * 1e6)))


def dec(v):
    return "%.6f" % v


class Gerber:
    """A Gerber X2 file as KiCad writes it."""

    def __init__(self, function, polarity="Positive"):
        self.function = function
        self.polarity = polarity
        self.apertures = []  # (definition, function)
        self.codes = {}
        self.macros = []
        self.body = []
        self.cur = None
        self.pol = "D"

    def macro(self, text):
        if text not in self.macros:
            self.macros.append(text)

    def ap(self, definition, function="Conductor"):
        key = (definition, function)
        if key not in self.codes:
            self.codes[key] = 10 + len(self.apertures)
            self.apertures.append(key)
        return self.codes[key]

    def use(self, code):
        if self.cur != code:
            self.body.append("D%d*" % code)
            self.cur = code

    def polarity_(self, p):
        if self.pol != p:
            self.body.append("%%LP%s*%%" % p)
            self.pol = p

    def flash(self, code, x, y):
        self.use(code)
        self.body.append("X%sY%sD03*" % (num(x), num(y)))

    def line(self, code, pts):
        self.use(code)
        self.body.append("X%sY%sD02*" % (num(pts[0][0]), num(pts[0][1])))
        for x, y in pts[1:]:
            self.body.append("X%sY%sD01*" % (num(x), num(y)))

    def arc(self, code, start, end, center, ccw):
        """An arc from start to end about center (G75)."""
        self.use(code)
        self.body.append("X%sY%sD02*" % (num(start[0]), num(start[1])))
        self.body.append("G75*")
        self.body.append("G0%dX%sY%sI%sJ%sD01*" % (3 if ccw else 2, num(end[0]), num(end[1]),
                                                   num(center[0] - start[0]), num(center[1] - start[1])))
        self.body.append("G01*")

    def region(self, contours):
        """Contours of ("L", x, y) and ("A", x, y, cx, cy, ccw) steps from a start point."""
        self.body.append("G36*")
        for c in contours:
            x0, y0 = c[0]
            self.body.append("X%sY%sD02*" % (num(x0), num(y0)))
            px, py = x0, y0
            for s in c[1:]:
                if s[0] == "L":
                    self.body.append("G01X%sY%sD01*" % (num(s[1]), num(s[2])))
                else:
                    _, x, y, cx, cy, ccw = s
                    self.body.append("G75*")
                    self.body.append("G0%dX%sY%sI%sJ%sD01*" % (3 if ccw else 2, num(x), num(y), num(cx - px), num(cy - py)))
                px, py = s[1], s[2]
        self.body.append("G37*")

    def write(self):
        out = [
            "%TF.GenerationSoftware,KiCad,Pcbnew,7.0.10*%",
            "%TF.CreationDate," + DATE + "*%",
            "%TF.ProjectId,bdf-demo,62646664-656d-6f00-0000-000000000000,rev?*%",
            "%TF.SameCoordinates,Original*%",
            "%TF.FileFunction," + self.function + "*%",
            "%TF.FilePolarity," + self.polarity + "*%",
            "%FSLAX46Y46*%",
            "G04 Gerber Fmt 4.6, Leading zero omitted, Abs format (unit mm)*",
            "G04 Created by KiCad (PCBNEW 7.0.10) date 2026-09-27 12:00:00*",
            "%MOMM*%",
            "%LPD*%",
            "G01*",
            "G04 APERTURE LIST*",
        ]
        for m in self.macros:
            out.append(m)
        for (definition, function) in self.apertures:
            out.append("%%TA.AperFunction,%s*%%" % function)
            out.append("%%ADD%d%s*%%" % (self.codes[(definition, function)], definition))
            out.append("%TD*%")
        out.append("G04 APERTURE END LIST*")
        out.extend(self.body)
        out.append("M02*")
        return "\n".join(out) + "\n"


ROUNDRECT = """%AMRoundRect*
0 Rectangle with rounded corners*
0 $1 Rounding radius*
0 $2 $3 $4 $5 $6 $7 $8 $9 X,Y pos of 4 corners*
0 Add a 4 corners polygon primitive as pad body*
4,1,4,$2,$3,$4,$5,$6,$7,$8,$9,$2,$3,0*
0 Add four circle primitives for the rounded corners*
1,1,$1+$1,$2,$3*
1,1,$1+$1,$4,$5*
1,1,$1+$1,$6,$7*
1,1,$1+$1,$8,$9*
0 Add four rect primitives between the rounded corners*
20,1,$1+$1,$2,$3,$4,$5,0*
20,1,$1+$1,$4,$5,$6,$7,0*
20,1,$1+$1,$6,$7,$8,$9,0*
20,1,$1+$1,$8,$9,$2,$3,0*%"""

THERMAL = """%AMThermal*
0 Thermal relief: outer diameter, inner diameter, spoke width*
7,0,0,$1,$2,$3,45*%"""


def roundrect(w, h, r):
    x, y = w / 2 - r, h / 2 - r
    return "RoundRect,%sX%sX%sX%sX%sX%sX%sX%sX%sX0" % (
        dec(r), dec(-x), dec(-y), dec(x), dec(-y), dec(x), dec(y), dec(-x), dec(y))


# A stroke font on a grid 4 wide and 6 high.
GLYPHS = {
    "B": [[(0, 0), (0, 6), (3, 6), (4, 5), (4, 4), (3, 3), (0, 3)], [(3, 3), (4, 2), (4, 1), (3, 0), (0, 0)]],
    "D": [[(0, 0), (0, 6), (2.5, 6), (4, 4.5), (4, 1.5), (2.5, 0), (0, 0)]],
    "F": [[(4, 6), (0, 6), (0, 0)], [(0, 3), (3, 3)]],
    "R": [[(0, 0), (0, 6), (3, 6), (4, 5), (4, 4), (3, 3), (0, 3)], [(2, 3), (4, 0)]],
    "C": [[(4, 5), (3, 6), (1, 6), (0, 5), (0, 1), (1, 0), (3, 0), (4, 1)]],
    "U": [[(0, 6), (0, 1), (1, 0), (3, 0), (4, 1), (4, 6)]],
    "J": [[(4, 6), (4, 1), (3, 0), (1, 0), (0, 1)]],
    "O": [[(1, 0), (0, 1), (0, 5), (1, 6), (3, 6), (4, 5), (4, 1), (3, 0), (1, 0)]],
    "T": [[(0, 6), (4, 6)], [(2, 6), (2, 0)]],
    "M": [[(0, 0), (0, 6), (2, 3), (4, 6), (4, 0)]],
    "G": [[(4, 5), (3, 6), (1, 6), (0, 5), (0, 1), (1, 0), (3, 0), (4, 1), (4, 3), (2, 3)]],
    "E": [[(4, 6), (0, 6), (0, 0), (4, 0)], [(0, 3), (3, 3)]],
    "H": [[(0, 0), (0, 6)], [(4, 0), (4, 6)], [(0, 3), (4, 3)]],
    "N": [[(0, 0), (0, 6), (4, 0), (4, 6)]],
    "1": [[(1, 5), (2, 6), (2, 0)], [(1, 0), (3, 0)]],
    "2": [[(0, 5), (1, 6), (3, 6), (4, 5), (4, 4), (0, 0), (4, 0)]],
    "3": [[(0, 5), (1, 6), (3, 6), (4, 5), (4, 4), (3, 3), (4, 2), (4, 1), (3, 0), (1, 0), (0, 1)], [(1, 3), (3, 3)]],
    "4": [[(3, 0), (3, 6), (0, 2), (4, 2)]],
}


def text(g, code, s, cx, cy, height, mirror=False):
    """Text centred on (cx, cy); mirrored as bottom silkscreen is plotted."""
    k = height / 6
    width = (len(s) * 5 - 1) * k
    x0 = cx - width / 2
    for i, ch in enumerate(s):
        for stroke in GLYPHS.get(ch, []):
            pts = []
            for gx, gy in stroke:
                x = x0 + (i * 5 + gx) * k
                if mirror:
                    x = 2 * cx - x
                pts.append((x, cy - height / 2 + gy * k))
            g.line(code, pts)


def rect_outline(g, code, x0, y0, x1, y1):
    g.line(code, [(x0, y0), (x1, y0), (x1, y1), (x0, y1), (x0, y0)])


def board():
    # parts: SOIC-8 U1, 0805 R1 R2, vertical 0805 C1, 2x3 header J1, 2-pin
    # terminal J2 with slotted oval pads, four mounting holes, fiducials, a
    # test point and vias
    soic = [(127.525, -116.095), (127.525, -117.365), (127.525, -118.635), (127.525, -119.905),
            (132.475, -119.905), (132.475, -118.635), (132.475, -117.365), (132.475, -116.095)]
    r1 = [(115.05, -110), (116.95, -110)]
    r2 = [(115.05, -114), (116.95, -114)]
    c1 = [(144, -109.05), (144, -110.95)]
    j1 = [(144.73, -123.46), (147.27, -123.46), (144.73, -126), (147.27, -126), (144.73, -128.54), (147.27, -128.54)]
    j2 = [(108, -130), (112, -130)]
    holes = [(104, -104), (156, -104), (104, -136), (156, -136)]
    fids = [(104, -120), (156, -120)]
    tp = (130, -133)
    vias = [(122, -110), (113, -110), (138, -124), (126, -130), (150, -116)]

    files = {}

    # top copper
    g = Gerber("Copper,L1,Top")
    g.macro(ROUNDRECT)
    soic_pad = g.ap(roundrect(1.95, 0.6, 0.15), "SMDPad,CuDef")
    r_pad = g.ap(roundrect(1.0, 1.45, 0.25), "SMDPad,CuDef")
    c_pad = g.ap(roundrect(1.45, 1.0, 0.25), "SMDPad,CuDef")
    sq = g.ap("R,1.700000X1.700000", "ComponentPad")
    rd = g.ap("C,1.700000", "ComponentPad")
    ov = g.ap("O,2.000000X3.000000", "ComponentPad")
    fid = g.ap("C,1.000000", "FiducialPad,Global")
    tpp = g.ap("C,1.500000", "TestPad")
    via = g.ap("C,0.600000", "ViaPad")
    sig = g.ap("C,0.250000")
    pwr = g.ap("C,0.500000")
    g.line(sig, [r1[1], (122, -110)])
    g.line(sig, [soic[0], (125, -116.095), (122, -113.095), (122, -110)])
    g.line(sig, [r2[1], (119, -114), (122.365, -117.365), soic[1]])
    g.line(sig, [soic[7], (137, -116.095), (142.145, -110.95), c1[1]])
    g.line(sig, [soic[6], (138, -117.365), (142, -113.365), (148, -113.365), (150, -115.365), (150, -116)])
    g.line(sig, [r1[0], (113, -110)])
    g.line(sig, [r2[0], (110, -114), (108, -116), j2[0]])
    g.line(sig, [soic[2], (124.5, -118.635), (124.5, -128.5), (126, -130)])
    g.line(pwr, [j1[0], (138.5, -123.46), (134.945, -119.905), soic[4]])
    g.line(pwr, [c1[0], (144, -107.5), (158, -107.5), (158, -126), (147.27, -126)])
    g.line(sig, [j1[2], (141, -126)])
    g.arc(sig, (141, -126), (134, -133), (141, -133), True)
    g.line(sig, [(134, -133), tp])
    # a copper area with a hole made by a cut-in
    g.region([[(146, -136), ("L", 154, -136), ("L", 154, -133), ("L", 151.5, -133),
               ("A", 151.5, -133, 150, -133, False), ("L", 154, -133), ("L", 154, -130), ("L", 146, -130), ("L", 146, -136)]])
    for p in soic:
        g.flash(soic_pad, *p)
    for p in r1 + r2:
        g.flash(r_pad, *p)
    for p in c1:
        g.flash(c_pad, *p)
    g.flash(sq, *j1[0])
    for p in j1[1:]:
        g.flash(rd, *p)
    for p in j2:
        g.flash(ov, *p)
    for p in fids:
        g.flash(fid, *p)
    g.flash(tpp, *tp)
    for p in vias:
        g.flash(via, *p)
    files["bdf-demo-F_Cu.gbr"] = g.write()

    # bottom copper: a ground pour with clearances and a thermal relief
    g = Gerber("Copper,L2,Bot")
    g.macro(THERMAL)
    sq = g.ap("R,1.700000X1.700000", "ComponentPad")
    rd = g.ap("C,1.700000", "ComponentPad")
    ov = g.ap("O,2.000000X3.000000", "ComponentPad")
    via = g.ap("C,0.600000", "ViaPad")
    sig = g.ap("C,0.250000")
    clr_pad = g.ap("C,2.300000", "Other,Clearance")
    clr_ov = g.ap("O,2.600000X3.600000", "Other,Clearance")
    clr_via = g.ap("C,1.200000", "Other,Clearance")
    clr_hole = g.ap("C,4.000000", "Other,Clearance")
    clr_sig = g.ap("C,0.850000", "Other,Clearance")
    thermal = g.ap("Thermal,2.600000X1.900000X0.500000", "Other,ThermalRelief")
    r, x0, y0, x1, y1 = 2.5, 100.5, -139.5, 159.5, -100.5
    g.region([[(x0 + r, y0), ("L", x1 - r, y0), ("A", x1, y0 + r, x1 - r, y0 + r, True), ("L", x1, y1 - r),
               ("A", x1 - r, y1, x1 - r, y1 - r, True), ("L", x0 + r, y1), ("A", x0, y1 - r, x0 + r, y1 - r, True),
               ("L", x0, y0 + r), ("A", x0 + r, y0, x0 + r, y0 + r, True)]])
    routes = [[(113, -110), (113, -126), (112, -127), j2[1]],
              [(122, -110), (122, -107), (147, -107), (150, -110), (150, -116)],
              [(138, -124), (132, -130), (126, -130)]]
    g.polarity_("C")
    for p in [j1[0]] + j1[2:]:
        g.flash(clr_pad, *p)
    g.flash(thermal, *j1[1])
    g.flash(clr_ov, *j2[0])
    g.flash(clr_ov, *j2[1])
    for p in vias:
        g.flash(clr_via, *p)
    for p in holes:
        g.flash(clr_hole, *p)
    for rt in routes:
        g.line(clr_sig, rt)
    g.polarity_("D")
    for rt in routes:
        g.line(sig, rt)
    g.flash(sq, *j1[0])
    for p in j1[1:]:
        g.flash(rd, *p)
    for p in j2:
        g.flash(ov, *p)
    for p in vias:
        g.flash(via, *p)
    files["bdf-demo-B_Cu.gbr"] = g.write()

    # solder masks: openings 0.05 mm larger than the pads
    g = Gerber("Soldermask,Top", "Negative")
    g.macro(ROUNDRECT)
    for p in soic:
        g.flash(g.ap(roundrect(2.05, 0.7, 0.15), "SMDPad,CuDef"), *p)
    for p in r1 + r2:
        g.flash(g.ap(roundrect(1.1, 1.55, 0.25), "SMDPad,CuDef"), *p)
    for p in c1:
        g.flash(g.ap(roundrect(1.55, 1.1, 0.25), "SMDPad,CuDef"), *p)
    g.flash(g.ap("R,1.800000X1.800000", "ComponentPad"), *j1[0])
    for p in j1[1:]:
        g.flash(g.ap("C,1.800000", "ComponentPad"), *p)
    for p in j2:
        g.flash(g.ap("O,2.100000X3.100000", "ComponentPad"), *p)
    for p in fids:
        g.flash(g.ap("C,2.000000", "FiducialPad,Global"), *p)
    g.flash(g.ap("C,1.600000", "TestPad"), *tp)
    for p in holes:
        g.flash(g.ap("C,3.400000", "ComponentPad"), *p)
    files["bdf-demo-F_Mask.gbr"] = g.write()

    g = Gerber("Soldermask,Bot", "Negative")
    g.flash(g.ap("R,1.800000X1.800000", "ComponentPad"), *j1[0])
    for p in j1[1:]:
        g.flash(g.ap("C,1.800000", "ComponentPad"), *p)
    for p in j2:
        g.flash(g.ap("O,2.100000X3.100000", "ComponentPad"), *p)
    for p in holes:
        g.flash(g.ap("C,3.400000", "ComponentPad"), *p)
    files["bdf-demo-B_Mask.gbr"] = g.write()

    # paste: the SMD pads
    g = Gerber("Paste,Top")
    g.macro(ROUNDRECT)
    for p in soic:
        g.flash(g.ap(roundrect(1.95, 0.6, 0.15), "SMDPad,CuDef"), *p)
    for p in r1 + r2:
        g.flash(g.ap(roundrect(1.0, 1.45, 0.25), "SMDPad,CuDef"), *p)
    for p in c1:
        g.flash(g.ap(roundrect(1.45, 1.0, 0.25), "SMDPad,CuDef"), *p)
    files["bdf-demo-F_Paste.gbr"] = g.write()

    # silkscreen
    g = Gerber("Legend,Top")
    thin = g.ap("C,0.120000", "Other,Legend")
    txt = g.ap("C,0.150000", "Other,Legend")
    big = g.ap("C,0.300000", "Other,Legend")
    dot = g.ap("C,0.300000", "Other,Legend")
    rect_outline(g, thin, 128.1, -120.55, 131.9, -115.45)
    g.flash(dot, 126.6, -115.0)
    rect_outline(g, thin, 143.4, -129.81, 148.6, -122.19)
    rect_outline(g, thin, 106.5, -132.2, 113.5, -127.8)
    for cx, cy in [(116, -110), (116, -114)]:
        g.line(thin, [(115.4, cy + 0.85), (116.6, cy + 0.85)])
        g.line(thin, [(115.4, cy - 0.85), (116.6, cy - 0.85)])
    text(g, txt, "R1", 116, -108.2, 1.0)
    text(g, txt, "R2", 116, -115.8, 1.0)
    text(g, txt, "C1", 146.2, -110, 1.0)
    text(g, txt, "U1", 130, -122, 1.0)
    text(g, txt, "J1", 146, -121, 1.0)
    text(g, txt, "J2", 110, -126.8, 1.0)
    text(g, big, "BDF", 130, -103.5, 2.5)
    files["bdf-demo-F_Silkscreen.gbr"] = g.write()

    g = Gerber("Legend,Bot")
    txt = g.ap("C,0.200000", "Other,Legend")
    text(g, txt, "BOTTOM", 130, -110, 2.0, mirror=True)
    text(g, txt, "GND", 136, -134, 1.5, mirror=True)
    files["bdf-demo-B_Silkscreen.gbr"] = g.write()

    # the outline: a rounded rectangle and a slot cut out of the board,
    # drawn in pieces out of order and in both directions
    g = Gerber("Profile,NP")
    e = g.ap("C,0.100000", "Profile")
    r = 3
    g.line(e, [(103, -140), (157, -140)])
    g.line(e, [(160, -103), (160, -137)])
    g.line(e, [(100, -103), (100, -137)])
    g.line(e, [(157, -100), (103, -100)])
    g.arc(e, (157, -140), (160, -137), (157, -137), True)
    g.arc(e, (157, -100), (160, -103), (157, -103), False)
    g.arc(e, (103, -100), (100, -103), (103, -103), True)
    g.arc(e, (100, -137), (103, -140), (103, -137), True)
    # the slot: 6 × 1.5 mm with round ends
    g.line(e, [(119, -136.75), (125, -136.75)])
    g.line(e, [(119, -135.25), (125, -135.25)])
    g.arc(e, (125, -136.75), (125, -135.25), (125, -136), True)
    g.arc(e, (119, -135.25), (119, -136.75), (119, -136), True)
    files["bdf-demo-Edge_Cuts.gbr"] = g.write()

    # drill files
    def drill(function, tools, body):
        out = ["M48", "; DRILL file {KiCad 7.0.10} date 2026-09-27T12:00:00+0900",
               "; FORMAT={-:-/ absolute / metric / decimal}",
               "; #@! TF.CreationDate," + DATE,
               "; #@! TF.GenerationSoftware,Kicad,Pcbnew,7.0.10",
               "; #@! TF.FileFunction," + function,
               "FMAT,2", "METRIC"]
        for n, (d, f) in enumerate(tools, 1):
            out.append("; #@! TA.AperFunction," + f)
            out.append("T%dC%.3f" % (n, d))
        out += ["%", "G90", "G05"] + body + ["M30"]
        return "\n".join(out) + "\n"

    def xy(p):
        # KiCad's decimal format: at least one decimal
        f = lambda v: ("%.3f" % v).rstrip("0").rstrip(".") + (".0" if v == int(v) else "")
        return "X%sY%s" % (f(p[0]), f(p[1]))

    body = ["T1"] + [xy(p) for p in vias] + ["T2"] + [xy(p) for p in j1] + ["T3"]
    for x, y in j2:
        body += ["G00" + xy((x, y + 0.5)), "M15", "G01" + xy((x, y - 0.5)), "M16", "G05"]
    body.append("T0")
    files["bdf-demo-PTH.drl"] = drill("Plated,1,2,PTH", [
        (0.3, "Plated,PTH,ViaDrill"), (1.0, "Plated,PTH,ComponentDrill"), (1.0, "Plated,PTH,ComponentDrill")], body)
    files["bdf-demo-NPTH.drl"] = drill("NonPlated,1,2,NPTH", [(3.2, "NonPlated,NPTH,ComponentDrill")],
                                       ["T1"] + [xy(p) for p in holes] + ["T0"])

    job = {
        "Header": {"GenerationSoftware": {"Vendor": "KiCad", "Application": "Pcbnew", "Version": "7.0.10"},
                   "CreationDate": DATE},
        "GeneralSpecs": {"ProjectId": {"Name": "bdf-demo", "GUID": "62646664-656d-6f00-0000-000000000000", "Revision": "rev?"},
                         "Size": {"X": 60.1, "Y": 40.1}, "LayerNumber": 2, "BoardThickness": 1.6, "Finish": "ENIG"},
        "FilesAttributes": [{"Path": n, "FileFunction": f, "FilePolarity": p} for n, f, p in [
            ("bdf-demo-F_Cu.gbr", "Copper,L1,Top", "Positive"), ("bdf-demo-B_Cu.gbr", "Copper,L2,Bot", "Positive"),
            ("bdf-demo-F_Paste.gbr", "SolderPaste,Top", "Positive"),
            ("bdf-demo-F_Silkscreen.gbr", "Legend,Top", "Positive"), ("bdf-demo-B_Silkscreen.gbr", "Legend,Bot", "Positive"),
            ("bdf-demo-F_Mask.gbr", "SolderMask,Top", "Negative"), ("bdf-demo-B_Mask.gbr", "SolderMask,Bot", "Negative"),
            ("bdf-demo-Edge_Cuts.gbr", "Profile", "Positive")]],
        "MaterialStackup": [
            {"Type": "Legend", "Color": "White", "Name": "Top Silk Screen"},
            {"Type": "SolderPaste", "Name": "Top Solder Paste"},
            {"Type": "SolderMask", "Color": "Green", "Thickness": 0.01, "Name": "Top Solder Mask"},
            {"Type": "Copper", "Thickness": 0.035, "Name": "F.Cu"},
            {"Type": "Dielectric", "Thickness": 1.51, "Material": "FR4", "Name": "F.Cu/B.Cu"},
            {"Type": "Copper", "Thickness": 0.035, "Name": "B.Cu"},
            {"Type": "SolderMask", "Color": "Green", "Thickness": 0.01, "Name": "Bottom Solder Mask"},
            {"Type": "Legend", "Color": "White", "Name": "Bottom Silk Screen"}],
    }
    files["bdf-demo-job.gbrjob"] = json.dumps(job, indent=2) + "\n"

    buf = io.BytesIO()
    with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as z:
        for name in sorted(files):
            info = zipfile.ZipInfo("bdf-demo/" + name, date_time=(2026, 9, 27, 12, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            info.external_attr = 0o644 << 16
            z.writestr(info, files[name])
    with open(os.path.join(OUT, "board.zip"), "wb") as f:
        f.write(buf.getvalue())


def features():
    """Every construct of the specification, in rows."""
    L = []
    w = L.append
    w("G04 Gerber features for the bdf converter tests*")
    w("%TF.FileFunction,Other,Features*%")
    w("%FSLAX46Y46*%")
    w("%MOMM*%")
    # macros
    w("%AMCIRCLES*1,1,1.5,0,0*1,1,0.8,1.2,0,45*%")
    w("%AMLINES*20,1,0.4,-1,-1,1,1,0*21,1,2.0,0.5,0,0,30*22,1,0.6,1.6,-1.2,-0.8,0*%")
    w("%AMSTAR*4,1,10,")
    pts = []
    for i in range(11):
        a = math.pi / 2 + i * math.pi / 5
        r = 1.2 if i % 2 == 0 else 0.5
        pts.append("%.4f,%.4f" % (r * math.cos(a), r * math.sin(a)))
    L[-1] += ",".join(pts) + ",0*%"
    w("%AMHEX*5,1,6,0,0,2.2,15*%")
    w("%AMMOIRE*6,0,0,2.4,0.2,0.2,3,0.1,2.8,0*%")
    w("%AMTHERM*7,0,0,2.4,1.6,0.4,0*%")
    w("%AMDONUT*0 a ring made by taking a circle away from a circle*$3=$1-$2*$4=$3/2*1,1,$1,0,0*1,0,$2,0,0*21,1,$1,$4,0,0,0*%")
    w("%AMEXPR*0 arithmetic: (1+2)x0.5-0.25/0.5 = 1.0*1,1,(1+2)x0.5-0.25/0.5,0,0*1,1,$1x-1+2.2,1.3,0*%")
    # standard apertures, some with holes
    w("%ADD10C,1.5*%")
    w("%ADD11C,1.5X0.6*%")
    w("%ADD12R,2.0X1.2*%")
    w("%ADD13R,2.0X1.2X0.5*%")
    w("%ADD14O,2.0X1.0*%")
    w("%ADD15O,1.0X2.0X0.4*%")
    w("%ADD16P,2.0X3*%")
    w("%ADD17P,2.0X5X18*%")
    w("%ADD18P,2.0X8X22.5X0.8*%")
    w("%ADD20CIRCLES*%")
    w("%ADD21LINES*%")
    w("%ADD22STAR*%")
    w("%ADD23HEX*%")
    w("%ADD24MOIRE*%")
    w("%ADD25THERM*%")
    w("%ADD26DONUT,2.4X1.2*%")
    w("%ADD27EXPR,0.6*%")
    w("%ADD30C,0.3*%")
    w("%ADD31R,0.6X0.3*%")
    w("%ADD32C,0*%")
    # row 1: standard apertures
    for i, d in enumerate([10, 11, 12, 13, 14, 15, 16, 17, 18]):
        w("D%d*" % d)
        w("X%sY%sD03*" % (num(2 + 3 * i), num(22)))
    # row 2: macros (the first selected the old way, with G54)
    for i, d in enumerate([20, 21, 22, 23, 24, 25, 26, 27]):
        w(("G54D%d*" if i == 0 else "D%d*") % d)
        w("X%sY%sD03*" % (num(2 + 3 * i), num(18)))
    # row 3: lines and arcs
    w("D30*")
    w("X%sY%sD02*" % (num(0.5), num(14)))
    w("G01X%sY%sD01*" % (num(3), num(14)))
    w("X%sY%sD01*" % (num(4), num(15)))
    w("Y%sD01*" % num(13))  # X kept from the point before
    w("D31*")  # a rectangle drawn: swept
    w("X%sY%sD02*" % (num(5.5), num(13)))
    w("X%sY%sD01*" % (num(7.5), num(15)))
    w("D30*")
    w("G75*")
    w("X%sY%sD02*" % (num(9), num(14)))
    w("G03X%sY%sI%sJ%sD01*" % (num(11), num(14), num(1), num(0)))  # a half circle, counterclockwise
    w("G02X%sY%sI%sJ%sD01*" % (num(13), num(14), num(1), num(0)))  # and one clockwise
    w("X%sY%sD02*" % (num(15), num(14)))
    w("G03X%sY%sI%sJ%sD01*" % (num(15), num(14), num(1), num(0)))  # a full circle
    w("G74*")  # single quadrant: the signs of I and J are left out
    w("X%sY%sD02*" % (num(18), num(13)))
    w("G02X%sY%sI%sJ%sD01*" % (num(19), num(14), num(1), num(0)))
    w("G01*")
    w("G75*")
    w("X%sY%sD02*" % (num(21), num(14)))
    w("X%sY%sD01*" % (num(21), num(14)))  # a dot: a draw of no length
    w("D32*")
    w("X%sY%sD02*" % (num(22.5), num(13)))
    w("X%sY%sD01*" % (num(24.5), num(15)))  # a hairline: an aperture of no size
    # row 4: regions and polarity
    w("G36*")
    w("X%sY%sD02*" % (num(0.5), num(8.5)))
    w("G01X%sY%sD01*" % (num(3.5), num(8.5)))
    w("G03X%sY%sI%sJ%sD01*" % (num(3.5), num(11.5), num(0), num(1.5)))
    w("G01X%sY%sD01*" % (num(0.5), num(11.5)))
    w("X%sY%sD01*" % (num(0.5), num(8.5)))
    w("G37*")
    # a region with a hole made by a cut-in, and two contours in one region
    w("G36*")
    for c in ["X6000000Y8500000D02*", "X9000000Y8500000D01*", "X9000000Y10000000D01*", "X8000000Y10000000D01*",
              "G02X8000000Y10000000I-500000J0D01*", "G01X9000000Y10000000D01*", "X9000000Y11500000D01*",
              "X6000000Y11500000D01*", "X6000000Y8500000D01*",
              "X10000000Y8500000D02*", "X11500000Y8500000D01*", "X10750000Y11500000D01*", "X10000000Y8500000D01*"]:
        w(c)
    w("G37*")
    # clear polarity: a hole through a square, and an island in the hole
    w("D12*")
    w("%LPD*%")
    w("G36*")
    w("X13000000Y8500000D02*")
    w("X16000000Y8500000D01*")
    w("X16000000Y11500000D01*")
    w("X13000000Y11500000D01*")
    w("X13000000Y8500000D01*")
    w("G37*")
    w("%LPC*%")
    w("D10*")
    w("X14500000Y10000000D03*")
    w("%LPD*%")
    w("D30*")
    w("X14500000Y10000000D03*")
    # row 5: step and repeat, a block aperture, aperture transformations
    w("%SRX3Y2I1.2J1.2*%")
    w("D31*")
    w("X500000Y2500000D03*")
    w("D30*")
    w("X500000Y2500000D02*")
    w("X1200000Y2500000D01*")
    w("%SR*%")
    w("%ABD40*%")
    w("D12*")
    w("X0Y0D03*")
    w("%LPC*%")
    w("D30*")
    w("X-500000Y0D03*")
    w("%LPD*%")
    w("D14*")
    w("X600000Y600000D03*")
    w("%AB*%")
    w("D40*")
    w("X6000000Y3000000D03*")
    w("%LR45*%")
    w("X9000000Y3000000D03*")
    w("%LMX*%")
    w("%LR0*%")
    w("X12000000Y3000000D03*")
    w("%LMN*%")
    w("%LS0.5*%")
    w("D22*")
    w("X15000000Y3000000D03*")
    w("%LS1*%")
    w("%LR30*%")
    w("D12*")
    w("X18000000Y3000000D03*")
    w("%LR0*%")
    w("M02*")
    with open(os.path.join(OUT, "features.gbr"), "w") as f:
        f.write("\n".join(L) + "\n")


if __name__ == "__main__":
    os.makedirs(OUT, exist_ok=True)
    board()
    features()
