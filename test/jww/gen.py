"""Writes the JWW test drawings of converter/jww into converter/jww/testdata.

Standard library only. The files follow Jw_cad's published data format
(jwdatafmt.txt, Jw_cad 7.02): an MFC CArchive with a header, the list of
figures and the list of block definitions. Text uses only characters of the
test fonts (converter/pptx/testdata/fonts).

- shapes.jww: data version 700 (Jw_cad 7) on A3: lines in every color, line
  type and width, arcs and ellipses, points, text (full- and half-width,
  turned, spaced, italic and bold, vertical), a dimension, solids (a
  quadrangle, one in any color, circle solids), nested blocks, a hidden
  layer group and layer, and auxiliary lines, which are not printed.
- old.jww: data version 300 (Jw_cad 3): no line widths per figure, no SXF
  tables, 4-byte block times, on A4.
"""

import math
import os
import struct

OUT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "converter", "jww", "testdata")


class Archive:
    """An MFC CArchive writer."""

    def __init__(self, version):
        self.b = bytearray()
        self.version = version
        self.classes = {}
        self.next = 1  # classes and objects share one index

    def u8(self, v):
        self.b += struct.pack("<B", v)

    def u16(self, v):
        self.b += struct.pack("<H", v)

    def u32(self, v):
        self.b += struct.pack("<I", v & 0xFFFFFFFF)

    def f64(self, v):
        self.b += struct.pack("<d", v)

    def str(self, s):
        raw = s.encode("cp932")
        if len(raw) < 0xFF:
            self.u8(len(raw))
        else:
            self.u8(0xFF)
            self.u16(len(raw))
        self.b += raw

    def count(self, n):
        self.u16(n)

    def obj(self, cls, body):
        if cls in self.classes:
            self.u16(0x8000 | self.classes[cls])
        else:
            self.u16(0xFFFF)
            self.u16(0)  # schema
            self.u16(len(cls))
            self.b += cls.encode()
            self.classes[cls] = self.next
            self.next += 1
        self.next += 1  # the object
        body()


def rgb(r, g, b):
    """A Windows COLORREF."""
    return r | g << 8 | b << 16


SCREEN = [rgb(255, 255, 255), rgb(0, 160, 200), rgb(0, 0, 0), rgb(0, 150, 0), rgb(200, 150, 0),
          rgb(200, 0, 200), rgb(0, 0, 255), rgb(0, 128, 128), rgb(230, 0, 0), rgb(160, 160, 160)]
PRINT = [rgb(255, 255, 255), rgb(0, 0, 0), rgb(0, 0, 0), rgb(0, 128, 0), rgb(128, 96, 0),
         rgb(160, 0, 160), rgb(0, 0, 200), rgb(0, 96, 96), rgb(200, 0, 0), rgb(128, 128, 128)]
# line types 2-9: bit patterns of Jw_cad's defaults
TYPES = [(0x99999999, 4), (0xC3C3C3C3, 8), (0xE7E7E7E7, 8), (0xF99FF99F, 16),
         (0xFFF8FFF8, 32), (0xF3F3F3F3, 8), (0xFF9FF9FF, 32), (0x22222222, 4)]


def header(a, memo, paper, widths100, hidden_group=None, hidden_layer=None):
    v = a.version
    a.b += b"JwwData."
    a.u32(v)
    a.str(memo)
    a.u32(paper)
    a.u32(0)  # the group being written
    for g in range(16):
        state = 0 if g == hidden_group else (3 if g == 0 else 2)
        a.u32(state)
        a.u32(0)
        a.f64(50.0)
        a.u32(0)
        for l in range(16):
            a.u32(0 if (g, l) == hidden_layer else (3 if (g, l) == (0, 0) else 2))
            a.u32(0)
    for _ in range(14):
        a.u32(0)
    for _ in range(5):
        a.u32(0)
    a.u32(0)
    a.u32(0xFFFFFFFF - 100 if widths100 else 100)  # -101: widths in hundredths of a mm
    a.f64(0); a.f64(0); a.f64(1)
    a.u32(0)
    a.u32(0)
    a.f64(0); a.f64(0); a.f64(0); a.f64(0); a.f64(0)
    for g in range(16):
        for l in range(16):
            a.str("%X%X" % (g, l) if l == 0 else "")
    for g in range(16):
        a.str("GROUP%d" % g if g < 2 else "")
    a.f64(0); a.f64(35); a.u32(0); a.f64(0)
    if v >= 300:
        a.f64(0); a.f64(0)
    a.u32(0)
    a.f64(1); a.f64(0); a.f64(0)
    a.f64(1); a.f64(0); a.f64(0)
    if v >= 300:
        for _ in range(8):
            a.f64(1); a.f64(0); a.f64(0); a.u32(0)
        a.f64(0); a.f64(0); a.f64(0); a.u32(0); a.f64(0); a.f64(0); a.f64(0); a.u32(0)
    else:
        for _ in range(4):
            a.f64(1); a.f64(0); a.f64(0)
    for _ in range(10):
        a.f64(1)
    a.f64(0)
    for n in range(10):
        a.u32(SCREEN[n]); a.u32(1)
    for n in range(10):
        width = [0, 13, 18, 25, 35, 50, 70, 18, 25, 13][n] if widths100 else [1, 1, 1, 2, 2, 3, 4, 1, 2, 1][n]
        a.u32(PRINT[n]); a.u32(width); a.f64(0.4)
    for n in range(2, 10):
        mask, unit = TYPES[n - 2]
        a.u32(mask); a.u32(unit); a.u32(1); a.u32(10)
    for n in range(11, 16):
        a.u32(0xFFFFFFFF); a.u32(10); a.u32(1); a.u32(10); a.u32(10)
    for n in range(16, 20):
        mask, unit = TYPES[n - 12]
        a.u32(mask); a.u32(unit); a.u32(2); a.u32(20)
    a.u32(0)
    a.u32(1)  # points printed at their radius
    a.u32(0); a.u32(0); a.u32(0); a.u32(0); a.u32(0); a.u32(0); a.u32(0); a.u32(0)
    a.u32(10 if v >= 600 else 0)  # 300 dpi
    if v >= 223:
        a.u32(0); a.u32(0); a.u32(0); a.u32(0); a.u32(0)
        a.f64(0); a.f64(0); a.f64(0); a.f64(0); a.f64(0)
    if v >= 225:
        a.f64(0); a.f64(0); a.f64(0); a.f64(0)
    if v >= 230:
        a.u32(0); a.u32(rgb(255, 0, 0))
    if v >= 420:
        for n in range(257):
            a.u32(rgb(n, 0, 255 - n)); a.u32(1)
        for n in range(257):
            a.str("")
            a.u32(rgb(n, 0, 255 - n)); a.u32(35); a.f64(0.3)
        for n in range(33):
            mask, unit = TYPES[n % 8]
            a.u32(mask); a.u32(unit); a.u32(1); a.u32(10)
        for n in range(33):
            a.str("SXF%d" % n if n == 1 else "")
            a.u32(2 if n == 1 else 0)
            for j in range(10):
                a.f64([3.0, 1.0][j] if n == 1 and j < 2 else 0)
    for i in range(1, 11):
        a.f64(i); a.f64(i); a.f64(0); a.u32(2)
    a.f64(3); a.f64(3); a.f64(0); a.u32(2); a.u32(1)
    a.f64(0); a.f64(0)
    a.u32(0)
    for _ in range(6):
        a.f64(0)


def base(a, style=1, color=2, width=0, layer=0, glayer=0, flag=0):
    a.u32(0)
    a.u8(style)
    a.u16(color)
    if a.version >= 351:
        a.u16(width)
    a.u16(layer)
    a.u16(glayer)
    a.u16(flag)


def sen(a, p, q, **kw):
    def body():
        base(a, **kw)
        a.f64(p[0]); a.f64(p[1]); a.f64(q[0]); a.f64(q[1])
    return lambda: a.obj("CDataSen", body)


def enko(a, c, r, start=0, arc=2 * math.pi, tilt=0, flat=1, full=True, **kw):
    def body():
        base(a, **kw)
        a.f64(c[0]); a.f64(c[1]); a.f64(r); a.f64(start); a.f64(arc); a.f64(tilt); a.f64(flat); a.u32(1 if full else 0)
    return lambda: a.obj("CDataEnko", body)


def ten(a, p, kari=False, color=2, **kw):
    def body():
        base(a, color=color, **kw)
        a.f64(p[0]); a.f64(p[1]); a.u32(1 if kari else 0)
    return lambda: a.obj("CDataTen", body)


def moji_body(a, p, s, w, h, spacing=0, angle=0, kind=1, font="ＭＳ ゴシック", color=2, flag=0, layer=0, glayer=0):
    base(a, style=1, color=color, width=0, layer=layer, glayer=glayer, flag=flag)
    cells = sum(w if len(ch.encode("cp932")) == 2 else w / 2 for ch in s) + spacing * (len(s) - 1)
    r = math.radians(angle)
    q = (p[0] + cells * math.cos(r), p[1] + cells * math.sin(r))
    a.f64(p[0]); a.f64(p[1]); a.f64(q[0]); a.f64(q[1])
    a.u32(kind); a.f64(w); a.f64(h); a.f64(spacing); a.f64(angle)
    a.str(font)
    a.str(s)


def moji(a, p, s, w, h, **kw):
    return lambda: a.obj("CDataMoji", lambda: moji_body(a, p, s, w, h, **kw))


def solid(a, pts, color=3, rgbv=None, style=1):
    def body():
        base(a, style=style, color=10 if rgbv is not None else color)
        p1, p2, p3, p4 = pts
        for q in (p1, p4, p2, p3):
            a.f64(q[0]); a.f64(q[1])
        if rgbv is not None:
            a.u32(rgbv)
    return lambda: a.obj("CDataSolid", body)


def blockref(a, p, number, sx=1, sy=1, rot=0, color=2):
    def body():
        base(a, color=color)
        a.f64(p[0]); a.f64(p[1]); a.f64(sx); a.f64(sy); a.f64(rot); a.u32(number)
    return lambda: a.obj("CDataBlock", body)


def blockdef(a, number, name, items, time_size):
    def body():
        base(a)
        a.u32(number); a.u32(1)
        a.b += b"\0" * time_size
        a.str(name)
        a.count(len(items))
        for it in items:
            it()
    return lambda: a.obj("CDataList", body)


def sunpou(a, p, q, text):
    """A dimension of version 4.20 and later: line, text, extension lines
    and points."""
    def body():
        base(a, color=1)
        base(a, color=1)
        a.f64(p[0]); a.f64(p[1] + 8); a.f64(q[0]); a.f64(q[1] + 8)
        moji_body(a, ((p[0] + q[0]) / 2 - 3, p[1] + 9), text, 2.5, 2.5, color=1)
        if a.version >= 420:
            a.u16(0)
            base(a, color=1)
            a.f64(p[0]); a.f64(p[1] + 1); a.f64(p[0]); a.f64(p[1] + 10)
            base(a, color=1)
            a.f64(q[0]); a.f64(q[1] + 1); a.f64(q[0]); a.f64(q[1] + 10)
            for pt, kari in ((p, False), (q, False), (p, True), (q, True)):
                base(a, color=1)
                a.f64(pt[0]); a.f64(pt[1] + 8 if not kari else pt[1]); a.u32(1 if kari else 0)
    return lambda: a.obj("CDataSunpou", body)


def write(path, a, figures, blocks):
    a.count(len(figures))
    for f in figures:
        f()
    a.count(len(blocks))
    for b in blocks:
        b()
    if a.version >= 700:
        a.u32(0)  # no images
    with open(path, "wb") as f:
        f.write(bytes(a.b))


def shapes():
    a = Archive(700)
    header(a, "テスト図", 3, True, hidden_group=5, hidden_layer=(0, 3))
    figs = []
    # a frame on the sheet (A3: ±210 × ±148.5)
    for p, q in (((-200, -140), (200, -140)), ((200, -140), (200, 140)), ((200, 140), (-200, 140)), ((-200, 140), (-200, -140))):
        figs.append(sen(a, p, q, color=2, width=50))
    # line colors 1-8 and line types 1-8, 16-19
    for i in range(8):
        y = 125 - i * 6
        figs.append(sen(a, (-190, y), (-110, y), color=i + 1, style=1))
        figs.append(sen(a, (-100, y), (-20, y), color=2, style=i + 1 if i > 0 else 1))
    for i, st in enumerate((16, 17, 18, 19, 31)):
        y = 125 - i * 6
        figs.append(sen(a, (-10, y), (60, y), color=2, style=st))
    figs.append(sen(a, (-190, 72), (-110, 72), color=8, width=100))  # 1 mm wide
    figs.append(sen(a, (-190, 66), (-110, 66), color=9))  # auxiliary color: not printed
    figs.append(sen(a, (-190, 60), (-110, 60), style=9))  # auxiliary line type: not printed
    figs.append(sen(a, (-190, 54), (-110, 54), color=150))  # an SXF color
    figs.append(sen(a, (-190, 48), (-110, 48), glayer=5, color=8))  # hidden group
    figs.append(sen(a, (-190, 42), (-110, 42), layer=3, color=8))  # hidden layer
    figs.append(moji(a, (-100, 42), "HIDDEN", 4, 4, glayer=5))
    # arcs and ellipses
    figs.append(enko(a, (90, 115), 12, color=1))
    figs.append(enko(a, (125, 115), 12, start=math.radians(30), arc=math.radians(240), full=False, color=3))
    figs.append(enko(a, (165, 115), 16, tilt=math.radians(30), flat=0.5, color=5))
    figs.append(enko(a, (90, 80), 12, start=0, arc=math.pi * 1.5, tilt=math.radians(90), flat=0.6, full=False, color=6))
    # points
    for i in range(3):
        figs.append(ten(a, (120 + i * 8, 80), color=i + 1))
    figs.append(ten(a, (150, 80), kari=True))  # temporary: not printed
    # text
    figs.append(moji(a, (-190, 20), "文字のサイズ 5mm", 5, 5, kind=3))
    figs.append(moji(a, (-190, 10), "Jw_cad 2026 ハンカク", 4, 4, font="ＭＳ Ｐゴシック"))
    figs.append(moji(a, (-190, 0), "字間 あり", 4, 4, spacing=2))
    figs.append(moji(a, (-190, -10), "ヨコナガ", 8, 4))
    figs.append(moji(a, (-190, -20), "イタリック", 4, 4, kind=10001))
    figs.append(moji(a, (-120, -20), "ボールド", 4, 4, kind=20001))
    figs.append(moji(a, (-60, 10), "カイテン 30°", 4, 4, angle=30, color=6))
    figs.append(moji(a, (-20, 0), "タテガキ", 4, 4, flag=0x0020, color=8))
    # a dimension
    figs.append(sunpou(a, (0, -60), (80, -60), "4000"))
    figs.append(sen(a, (0, -60), (80, -60), color=2))
    # solids
    figs.append(solid(a, [(100, -20), (140, -20), (140, 0), (100, 10)], color=3))
    figs.append(solid(a, [(150, -20), (190, -20), (190, 10), (150, 10)], rgbv=rgb(255, 160, 0)))
    # circle solids: (centre, (radius, flatness), (tilt, start), (arc, kind))
    figs.append(solid(a, [(110, -50), (0, math.radians(0)), (math.radians(270), 100), (15, 1)], color=6, style=101))
    figs.append(solid(a, [(145, -50), (0, 0), (math.radians(120), 0), (15, 1)], color=8, style=101))
    figs.append(solid(a, [(180, -50), (0, 0), (2 * math.pi, 7), (15, 1)], color=4, style=105))
    # blocks
    figs.append(blockref(a, (-150, -80), 1))
    figs.append(blockref(a, (-110, -80), 1, sx=1.5, sy=1.5, rot=math.radians(30)))
    figs.append(blockref(a, (-60, -80), 2))
    blocks = [
        blockdef(a, 1, "MARK", [enko(a, (0, 0), 8, color=5), sen(a, (-10, 0), (10, 0), color=5), sen(a, (0, -10), (0, 10), color=5),
                                moji(a, (2, 2), "A", 3, 3, color=5)], 8),
        blockdef(a, 2, "PAIR", [blockref(a, (0, 0), 1), blockref(a, (25, 0), 1, rot=math.radians(45))], 8),
    ]
    write(os.path.join(OUT, "shapes.jww"), a, figs, blocks)


def old():
    a = Archive(300)
    header(a, "", 4, False)
    figs = [
        sen(a, (-140, -95), (140, -95)), sen(a, (140, -95), (140, 95)), sen(a, (140, 95), (-140, 95)), sen(a, (-140, 95), (-140, -95)),
        sen(a, (-120, 60), (120, 60), color=5, style=3),
        enko(a, (0, 0), 40, color=1),
        moji(a, (-60, -70), "キュウ バージョン", 6, 6, color=2),
        blockref(a, (80, -40), 1),
    ]
    blocks = [blockdef(a, 1, "BOX", [sen(a, (-10, -10), (10, -10)), sen(a, (10, -10), (10, 10)), sen(a, (10, 10), (-10, 10)),
                                     sen(a, (-10, 10), (-10, -10))], 4)]
    write(os.path.join(OUT, "old.jww"), a, figs, blocks)


def main():
    os.makedirs(OUT, exist_ok=True)
    shapes()
    old()


if __name__ == "__main__":
    main()
