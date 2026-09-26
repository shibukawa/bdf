"""Generates converter/pdf/testdata/cjk-cmaps.pdf: CJK text in predefined CMaps.

Like PDFs from older Japanese, Chinese and Korean software, the CID fonts
are not embedded and have no ToUnicode: the text is only readable through
the predefined CMaps (Shift_JIS, UCS-2, the -V vertical CMaps) and the CID
collections. One vertical column uses an embedded TrueType font (the M PLUS
1p subset of the PowerPoint tests) with Identity-V and a ToUnicode CMap. Written by hand with the
standard library only.
"""
import os
import struct
import zlib

root = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
out = os.path.join(root, "converter", "pdf", "testdata", "cjk-cmaps.pdf")
font_path = os.path.join(root, "converter", "pptx", "testdata", "fonts", "MPLUS1p-Regular-subset.ttf")


def ucs2(s):
    return s.encode("utf-16-be")


def hexstr(b):
    return "<" + b.hex() + ">"


def cmap_table(ttf):
    """Unicode → glyph index from a TrueType cmap (format 4 or 12)."""
    num = struct.unpack(">H", ttf[4:6])[0]
    tables = {}
    for i in range(num):
        tag, _, off, length = struct.unpack(">4sIII", ttf[12 + 16 * i:28 + 16 * i])
        tables[tag] = ttf[off:off + length]
    cm = tables[b"cmap"]
    n = struct.unpack(">H", cm[2:4])[0]
    best = None
    for i in range(n):
        pid, eid, off = struct.unpack(">HHI", cm[4 + 8 * i:12 + 8 * i])
        fmt = struct.unpack(">H", cm[off:off + 2])[0]
        if pid == 3 and fmt in (4, 12) and (best is None or fmt == 12):
            best = (off, fmt)
    off, fmt = best
    m = {}
    if fmt == 4:
        segs = struct.unpack(">H", cm[off + 6:off + 8])[0] // 2
        ends = struct.unpack(">%dH" % segs, cm[off + 14:off + 14 + 2 * segs])
        starts = struct.unpack(">%dH" % segs, cm[off + 16 + 2 * segs:off + 16 + 4 * segs])
        deltas = struct.unpack(">%dh" % segs, cm[off + 16 + 4 * segs:off + 16 + 6 * segs])
        ro_off = off + 16 + 6 * segs
        ros = struct.unpack(">%dH" % segs, cm[ro_off:ro_off + 2 * segs])
        for k in range(segs):
            for c in range(starts[k], ends[k] + 1):
                if ros[k] == 0:
                    g = (c + deltas[k]) & 0xFFFF
                else:
                    p = ro_off + 2 * k + ros[k] + 2 * (c - starts[k])
                    g = struct.unpack(">H", cm[p:p + 2])[0]
                    if g:
                        g = (g + deltas[k]) & 0xFFFF
                if g:
                    m[c] = g
    else:
        groups = struct.unpack(">I", cm[off + 12:off + 16])[0]
        for k in range(groups):
            s, e, g = struct.unpack(">III", cm[off + 16 + 12 * k:off + 28 + 12 * k])
            for c in range(s, e + 1):
                m[c] = g + c - s
    return m


ttf = open(font_path, "rb").read()
gids = cmap_table(ttf)

objs = []


def add(o):
    objs.append(o)
    return len(objs)


def stream(d, data, compress=True):
    if compress:
        data = zlib.compress(data, 9)
        d += " /Filter /FlateDecode"
    return ("<< %s /Length %d >>" % (d, len(data))).encode() + b"\nstream\n" + data + b"\nendstream"


def cid_font(name, registry_ordering, widths, fontfile=None, cidtogid=False):
    reg, ordering = registry_ordering.split("-")
    desc = ("<< /Type /FontDescriptor /FontName /%s /Flags 6 /FontBBox [0 -141 1000 859] /ItalicAngle 0"
            " /Ascent 859 /Descent -141 /CapHeight 709 /StemV 69" % name)
    if fontfile:
        desc += " /FontFile2 %d 0 R" % fontfile
    desc += " >>"
    sub = "CIDFontType2" if fontfile else "CIDFontType0"
    d = ("<< /Type /Font /Subtype /%s /BaseFont /%s /CIDSystemInfo << /Registry (%s) /Ordering (%s) /Supplement 2 >>"
         " /FontDescriptor %s /DW 1000 /W %s" % (sub, name, reg, ordering, desc, widths))
    if cidtogid:
        d += " /CIDToGIDMap /Identity"
    return add((d + " >>").encode())


def type0(name, encoding, descendant, to_unicode=None):
    d = ("<< /Type /Font /Subtype /Type0 /BaseFont /%s-%s /Encoding /%s /DescendantFonts [%d 0 R]"
         % (name, encoding, encoding, descendant))
    if to_unicode:
        d += " /ToUnicode %d 0 R" % to_unicode
    return add((d + " >>").encode())


def to_unicode(chars):
    """A ToUnicode CMap for 2-byte glyph codes (what dvipdfmx and others write)."""
    body = "\n".join("<%04x> <%s>" % (gids[ord(ch)], ch.encode("utf-16-be").hex()) for ch in sorted(set(chars)))
    cm = ("/CIDInit /ProcSet findresource begin 12 dict begin begincmap\n"
          "/CIDSystemInfo << /Registry (Adobe) /Ordering (UCS) /Supplement 0 >> def\n"
          "/CMapName /Adobe-Identity-UCS def /CMapType 2 def\n"
          "1 begincodespacerange <0000> <ffff> endcodespacerange\n"
          "%d beginbfchar\n%s\nendbfchar\nendcmap CMapName currentdict /CMap defineresource pop end end" % (len(set(chars)), body))
    return add(stream("", cm.encode()))


vertical_mplus = "たてがきのテキスト"
latin = "[1 95 500 231 326 500]"
ryumin = cid_font("Ryumin-Light", "Adobe-Japan1", latin)
gothic = cid_font("GothicBBB-Medium", "Adobe-Japan1", latin)
song = cid_font("STSong-Light", "Adobe-GB1", "[1 95 500]")
sung = cid_font("MSung-Light", "Adobe-CNS1", "[1 95 500]")
myeongjo = cid_font("HYSMyeongJo-Medium", "Adobe-Korea1", "[1 100 500]")
ff = add(stream("/Length1 %d" % len(ttf), ttf))
mplus = cid_font("MPLUS1p-Regular", "Adobe-Identity", "[]", fontfile=ff, cidtogid=True)

fonts = {
    "F1": type0("Ryumin-Light", "90ms-RKSJ-H", ryumin),
    "F2": type0("GothicBBB-Medium", "UniJIS-UCS2-H", gothic),
    "F3": type0("STSong-Light", "UniGB-UCS2-H", song),
    "F4": type0("MSung-Light", "UniCNS-UCS2-H", sung),
    "F5": type0("HYSMyeongJo-Medium", "UniKS-UCS2-H", myeongjo),
    "F6": type0("Ryumin-Light", "UniJIS-UCS2-V", ryumin),
    "F7": type0("GothicBBB-Medium", "90ms-RKSJ-V", gothic),
    "F8": type0("MPLUS1p-Regular", "Identity-V", mplus, to_unicode(vertical_mplus)),
}

content = "\n".join([
    "0.1 0.2 0.4 rg",
    "BT /F1 20 Tf 30 350 Td %s Tj ET" % hexstr("明朝体の日本語 Shift_JIS".encode("cp932")),
    "BT /F2 20 Tf 30 315 Td %s Tj ET" % hexstr(ucs2("ゴシック体：UCS-2 の文字列")),
    "BT /F3 20 Tf 30 280 Td %s Tj ET" % hexstr(ucs2("简体中文的文本")),
    "BT /F4 20 Tf 30 245 Td %s Tj ET" % hexstr(ucs2("繁體中文的文字")),
    "BT /F5 20 Tf 30 210 Td %s Tj ET" % hexstr(ucs2("한국어 텍스트")),
    "0.5 0.1 0.1 rg",
    "BT /F6 20 Tf 480 370 Td %s Tj ET" % hexstr(ucs2("縦書き「かぎ括弧」、ー。")),
    "BT /F7 20 Tf 440 370 Td %s Tj ET" % hexstr("縦組み（丸括弧）ー".encode("cp932")),
    "0.1 0.4 0.2 rg",
    "BT /F8 20 Tf 400 370 Td <%s> Tj ET" % "".join("%04x" % gids[ord(ch)] for ch in vertical_mplus),
    "0 g",
    "BT /F2 10 Tf 30 30 Td %s Tj ET" % hexstr(ucs2("横書きと縦書き、埋め込みなしのCIDフォント")),
]).encode()

contents = add(stream("", content))
font_res = " ".join("/%s %d 0 R" % (k, v) for k, v in fonts.items())
page = add(("<< /Type /Page /Parent %d 0 R /MediaBox [0 0 520 400] /Contents %d 0 R /Resources << /Font << %s >> >> >>"
            % (len(objs) + 2, contents, font_res)).encode())
pages = add(("<< /Type /Pages /Kids [%d 0 R] /Count 1 >>" % page).encode())
assert pages == page + 1
catalog = add(("<< /Type /Catalog /Pages %d 0 R /Lang (ja) >>" % pages).encode())
info = add(b"<< /Title (CJK predefined CMaps) /Producer (test/pdf/gen_cjk.py) >>")

buf = bytearray(b"%PDF-1.7\n%\xe2\xe3\xcf\xd3\n")
offsets = []
for i, o in enumerate(objs):
    offsets.append(len(buf))
    buf += b"%d 0 obj\n" % (i + 1) + o + b"\nendobj\n"
xref = len(buf)
buf += b"xref\n0 %d\n0000000000 65535 f \n" % (len(objs) + 1)
for off in offsets:
    buf += b"%010d 00000 n \n" % off
buf += b"trailer\n<< /Size %d /Root %d 0 R /Info %d 0 R >>\nstartxref\n%d\n%%%%EOF\n" % (len(objs) + 1, catalog, info, xref)
open(out, "wb").write(buf)
print("wrote", out, len(buf), "bytes")
