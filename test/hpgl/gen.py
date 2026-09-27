"""Writes the HP-GL/2 test plots of converter/hpgl into converter/hpgl/testdata.

Standard library only. The plots follow The HP-GL/2 and HP RTL Reference
Guide (Hewlett-Packard, 1996). Labels use only characters of the test fonts
(converter/pptx/testdata/fonts).

- shapes.plt: a bare HP-GL/2 plot on an A3 plot size: pens (NP, PC, PW),
  the eight line types fixed and adaptive, line ends and joins, arcs,
  circles (smooth and with a coarse chord angle), Bezier curves, an encoded
  polyline (PE), polygons with holes filled even-odd and non-zero, fill
  types (solid, hatching, cross-hatching, shading, a raster pattern, a PCL
  cross-hatch), rectangles and wedges (a pie chart), screened vectors, a
  clipping window, user scaling, symbol mode, and labels: stick font sizes,
  directions, slant, label origins, text paths, a proportional font, extra
  space, carriage returns and line feeds, and Japanese in Shift_JIS and in
  16-bit JIS.
- job.plt: a job for a large-format plotter: PJL, HP-GL/2 entered with
  ESC %-1B, turned with RO 90 to a portrait A4 sheet (as drivers do), with a
  24-bit raster image (HP RTL, replacement delta row compression) among
  the vectors; a second page is a black and white raster image only
  (adaptive compression).
- hpgl1.plt: HP-GL of a pen plotter: serial device control, P1 and P2,
  user scaling, labels with SI and DI, tick marks, and no plot size.
"""

import math
import os

OUT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "converter", "hpgl", "testdata")

ESC = b"\x1b"
ETX = b"\x03"


def pe_value(v, base=64):
    """Encodes a number for PE (sign in the lowest bit, base 64 digits)."""
    v = 2 * v if v >= 0 else 2 * -v + 1
    out = bytearray()
    while v >= base:
        out.append(63 + v % base)
        v //= base
    out.append((191 if base == 64 else 95) + v)
    return bytes(out)


def shapes():
    p = bytearray()

    def w(s):
        p.extend(s.encode("ascii") if isinstance(s, str) else s)

    def label(s, enc="ascii"):
        w(b"LB" + s.encode(enc) + ETX)

    w("IN;")
    w('BP1,"HP-GL/2 test plot";')
    w("PS16800,11880;IN;")  # A3: 420 x 297 mm
    w("NP16;PC1,0,0,0;PC2,200,0,0;PC3,0,140,0;PC4,0,0,200;PC5,230,120,0;PC6,120,0,160;PC8,90,90,90;")
    w("PW0.5,1;PW0.25,2;PW0.25,3;PW0.25,4;PW0.35,5;PW0.35,6;PW0.13,8;")
    # the sheet frame
    w("SP1;PU400,400;EA16400,11480;")

    # line types, fixed and adaptive
    w("SP2;SI0.2,0.28;")
    for i, t in enumerate(range(1, 9)):
        y = 11000 - i * 350
        w("PU800,%d;LT%d,5,1;PD3800,%d;PU;LT-%d,5,1;PD4400,%d,4400,%d;PU;LT;" % (y, t, y, t, y, y - 250))
        w("PU600,%d;" % (y - 60))
        label("%d" % t)
    # line ends and joins with a thick pen
    w("SP4;PW2,4;")
    for i, (end, join) in enumerate([(1, 1), (2, 5), (4, 4)]):
        x = 800 + i * 1300
        w("LA1,%d,2,%d;PU%d,7600;PD%d,8000,%d,7600;PU;" % (end, join, x, x + 400, x + 800))
    w("LA;PW0.25,4;")

    # arcs, circles and curves
    w("SP3;PU5200,10600;PD;AA5700,10600,-120;PU;PU5200,9900;PD;AR500,0,180,30;PU;")
    w("PU6600,10200;PD;AT7100,10800,7600,10200;PU;PU6600,9700;PD;RT500,-400,1000,0;PU;")
    w("PA8400,10200;CI500;CI350,45;")
    w("PU9200,9800;PD;BZ9400,10900,10000,9300,10300,10400;BR300,500,600,-600,900,0;PU;")
    # an encoded polyline: a star, pen 6
    pts = []
    for k in range(6):
        a = math.pi / 2 + k * 4 * math.pi / 5
        pts.append((round(12000 + 500 * math.cos(a)), round(10200 + 500 * math.sin(a))))
    pe = b"PE:" + pe_value(6) + b"<=" + pe_value(pts[0][0]) + pe_value(pts[0][1])
    for (x0, y0), (x1, y1) in zip(pts, pts[1:]):
        pe += pe_value(x1 - x0) + pe_value(y1 - y0)
    w(pe + b";")

    # polygons and fills
    w("SP1;PU5200,7200;PM0;PD6400,7200,6400,8400,5200,8400;PM1;PU5500,7500;PD6100,7500,6100,8100,5500,8100;PM2;")
    w("FT3,80,45;FP;EP;")
    w("PU6800,7200;PM0;PD8000,7200,7000,8400,7400,6900,7900,8400;PM2;FT1;SP5;FP1;SP1;EP;")
    w("PU8400,7200;PM0;CI600;CI300;PM2;FT4,100,30;SP3;FP;SP1;EP;")
    w("SP4;FT10,30;PU9400,6600;RR1000,1200;FT;SP1;ER1000,1200;")
    w("RF1,8,8," + ",".join("1" if (x + y) % 4 < 2 else "0" for y in range(8) for x in range(8)) + ";")
    w("SP6;FT11,1,1;PU10600,6600;RR1000,1200;FT;SP1;ER1000,1200;")
    w("SP2;FT21,5;PU11800,6600;RR1000,1200;FT;SP1;ER1000,1200;")
    # a pie chart
    w("PA14200,9900;")
    for start, sweep, pen, ft in [(0, 120, 2, "FT;"), (120, 90, 3, "FT3,60,0;"), (210, 150, 5, "FT10,60;")]:
        w("SP%d;%sWG900,%d,%d;SP1;FT;EW900,%d,%d;" % (pen, ft, start, sweep, start, sweep))
    # screened vectors and an opaque fill
    w("SP4;PW1.5,4;SV1,40;PU13200,8000;PD15600,8000;PU;SV;PW0.25,4;")
    w("SP5;TR0;FT10,50;PU13600,7700;RR600,600;TR;FT;")

    # user scaling: a chart in a window
    w("SP1;IP1000,1000,5000,4000;SC0,100,0,50;PU0,0;PD100,0,100,50,0,50,0,0;PU;")
    w("SP2;IW10,5,90,45;PU0,0;PD;")
    for x in range(0, 101, 5):
        w("PA%d,%d;" % (x, round(25 + 24 * math.sin(x / 100 * 2 * math.pi))))
    w("PU;IW;SP3;SM*;PU10,40;PD30,20,50,35,70,10,90,30;PU;SM;SC;IP;")

    # labels
    w("SP1;SI0.25,0.35;PU5600,5400;")
    label("Stick font 0.25 x 0.35 cm")
    w("SI0.5,0.7;PU5600,4700;")
    label("Large")
    w("SI;PU5600,4300;")
    label("Default size: 9 characters per inch")
    w("SI0.25,0.35;DI1,1;PU9400,2900;")
    label("Turned 45")
    w("DI;SL0.4;PU5600,3800;")
    label("Slanted")
    w("SL;ES0.5;PU5600,3400;")
    label("Extra space")
    w("ES;PU5600,3000;")
    label("Line one\r\nLine two")
    # label origins around a mark
    w("SP2;SI0.18,0.25;")
    for i, lo in enumerate([1, 4, 7, 3, 6, 9, 13, 19]):
        x, y = 11000 + (i % 4) * 1400, 5200 - (i // 4) * 700
        w("SP8;PU%d,%d;PD%d,%d;PU%d,%d;PD%d,%d;PU%d,%d;SP2;LO%d;" % (x - 80, y, x + 80, y, x, y - 80, x, y + 80, x, y, lo))
        label("LO%d" % lo)
    w("LO;DV1;PU15800,4600;")
    label("DOWN")
    w("DV;")
    # a proportional font, bold and italic, and the alternate font
    w("SP4;SD1,21,2,1,4,14,5,0,6,3,7,4148;SS;SI;PU11000,3200;")
    label("Univers Bold 14 pt")
    w("AD1,21,2,1,4,12,5,1,6,0,7,4101;SA;PU11000,2700;")
    label("CG Times Italic 12 pt")
    w("SS;SD;SA;")
    # Japanese: Shift_JIS in 8-bit mode and JIS in 16-bit mode
    w("SS;SP1;SI0.4,0.4;PU11000,2000;")
    label("文字 テスト", "cp932")
    w("SD1,1611,2,0,7,48;SS;LM1;PU11000,1400;LB")
    for ch in "テスト図":
        sj = ch.encode("iso2022_jp")[3:5]  # the JIS X 0208 code between the escapes
        w(sj)
    w(b"\x00" + ETX + b"LM;SD;SS;")

    w("SP0;PG;")
    return bytes(p)


def delta9(seed, row):
    """Replacement delta row compression (method 9), literal runs only."""
    out = bytearray()
    i = last = 0
    while i < len(row):
        if i < len(seed) and seed[i] == row[i]:
            i += 1
            continue
        j = i
        while j < len(row) and j - i < 8 and not (j < len(seed) and seed[j] == row[j]):
            j += 1
        off, cnt = i - last, j - i
        cmd = (min(off, 15) << 3) | min(cnt - 1, 7)
        out.append(cmd)
        if off >= 15:
            off -= 15
            while off >= 255:
                out.append(255)
                off -= 255
            out.append(off)
        if cnt - 1 >= 7:
            out.append(cnt - 1 - 7)
        out += row[i:j]
        last = j
        i = j
    return bytes(out)


def job():
    p = bytearray()
    p += ESC + b"%-12345X@PJL JOB NAME=\"Plotter job\"\r\n"
    p += b"@PJL SET RESOLUTION=300\r\n@PJL ENTER LANGUAGE=HPGL2\r\n"
    p += ESC + b"E" + ESC + b"%-1B"
    # a portrait A4 sheet on a roll: the length (297 mm) along the feed
    p += b'BP1,"Plotter job page";PS11880,8400;IN;RO90;IP;'
    p += b"NP4;PC1,0,0,0;PC2,0,0,200;PW0.4;SP1;"
    p += b"PU200,200;EA8200,11680;"  # in the turned system: 210 x 297 mm
    p += b"SI0.4,0.55;PU800,10800;LBPortrait page\x03"
    p += b"SP2;PU800,10300;PD7400,10300;PU;"
    # a raster image at the pen: 24 bits per pixel, 150 dpi
    p += b"PU1000,6000;" + ESC + b"%1A"
    p += ESC + b"*v6W" + bytes([0, 3, 8, 8, 8, 8]) + ESC + b"*t150R" + ESC + b"*r240S" + ESC + b"*r180T"
    p += ESC + b"*r1A" + ESC + b"*b9M"
    seed = b""
    for y in range(180):
        row = bytearray()
        for x in range(240):
            if (x - 120) ** 2 + (y - 90) ** 2 < 70 ** 2:
                row += bytes([230, 60 + y, 40])
            else:
                row += bytes([255 - y, 200, 120 + x // 2])
        d = delta9(seed, bytes(row))
        seed = bytes(row)
        p += ESC + b"*b%dW" % len(d) + d
    p += ESC + b"*rC" + ESC + b"%1B"
    p += b"SP1;PU1000,4400;LBRaster image above\x03;PG;"
    # page 2: a black and white raster image only, adaptive compression
    p += ESC + b"%0A" + ESC + b"E" + ESC + b"*t300R" + ESC + b"*r400S" + ESC + b"*r1A" + ESC + b"*b5M"
    block = bytearray()
    for y in range(300):
        row = bytearray(50)
        for x in range(400):
            if (x // 25 + y // 25) % 2 == 0 or abs(x - y) < 3:
                row[x // 8] |= 0x80 >> (x % 8)
        block += bytes([0, 0, len(row)]) + row
        if len(block) > 30000:
            p += ESC + b"*b%dW" % len(block) + block
            block = bytearray()
    block += bytes([4, 0, 20])  # 20 empty rows
    p += ESC + b"*b%dW" % len(block) + block + ESC + b"*rC"
    p += ESC + b"%0B" + b"PG;"
    p += ESC + b"%-12345X@PJL EOJ\r\n" + ESC + b"%-12345X"
    return bytes(p)


def hpgl1():
    p = bytearray()
    p += ESC + b".(" + ESC + b".I81;;17:" + ESC + b".N;19:"
    p += b"IN;IP250,596,10250,7796;SC0,1000,0,720;SP1;VS10;"
    p += b"PU0,0;PD1000,0,1000,720,0,720,0,0;PU;"
    p += b"PA100,100;XT;PA200,100;XT;PA300,100;XT;PA100,200;YT;PA100,300;YT;"
    p += b"SP2;PA100,100;PD200,300,300,250,400,500,500,400,600,600;PU;"
    p += b"SP1;SI0.3,0.4;PA100,650;LBHP-GL PEN PLOTTER\x03"
    p += b"DI0,1;PA950,100;LBVERTICAL\x03DI;"
    p += b"SP3;PA700,300;CI80;PA800,200;LT2;PD900,200,900,300;PU;LT;"
    p += b"SP0;"
    return bytes(p)


def main():
    os.makedirs(OUT, exist_ok=True)
    for name, data in (("shapes.plt", shapes()), ("job.plt", job()), ("hpgl1.plt", hpgl1())):
        with open(os.path.join(OUT, name), "wb") as f:
            f.write(data)


if __name__ == "__main__":
    main()
