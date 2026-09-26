"""Generates converter/pdf/testdata/images-jpx-jbig2.pdf: JPEG 2000 and JBIG2 images.

The JPEG 2000 images are synthetic pictures encoded with OpenJPEG's
opj_compress (lossy JP2 in RGB, a lossless grayscale codestream, and RGBA
with the alpha as SMaskInData). The JBIG2 images are streams from the
decoder's fixtures (converter/internal/jbig2/testdata): a halftone region
drawn as an image, and text whose symbol dictionary is in JBIG2Globals
drawn as a stencil mask. The PDF is written by hand
with the standard library; only opj_compress must be installed.
"""
import math
import os
import struct
import subprocess
import tempfile
import zlib

root = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
out = os.path.join(root, "converter", "pdf", "testdata", "images-jpx-jbig2.pdf")
jbig2_dir = os.path.join(root, "converter", "internal", "jbig2", "testdata")
tmp = tempfile.mkdtemp()


def picture(w, h, alpha=False):
    """Gradients, rings and a checker corner: smooth areas and sharp edges."""
    px = bytearray()
    for y in range(h):
        for x in range(w):
            r = int(255 * x / (w - 1))
            g = int(255 * y / (h - 1))
            d = math.hypot(x - w * 0.6, y - h * 0.45)
            b = int(127 + 127 * math.sin(d / 5))
            if x < w // 4 and y < h // 4:
                r = g = b = 255 if (x // 8 + y // 8) % 2 else 20
            px += bytes((r, g, b))
            if alpha:
                px.append(max(0, min(255, int(255 - 255 * d / (0.55 * w)))))
    return bytes(px)


def png(path, w, h, rgba):
    def chunk(tag, data):
        return struct.pack(">I", len(data)) + tag + data + struct.pack(">I", zlib.crc32(tag + data))
    rows = b"".join(b"\0" + rgba[y * w * 4:(y + 1) * w * 4] for y in range(h))
    with open(path, "wb") as f:
        f.write(b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 6, 0, 0, 0))
                + chunk(b"IDAT", zlib.compress(rows, 9)) + chunk(b"IEND", b""))


def opj(src, dst, *args):
    subprocess.run(["opj_compress", "-i", src, "-o", dst, *args], check=True, stdout=subprocess.DEVNULL)
    return open(dst, "rb").read()


W, H = 180, 120
rgb = picture(W, H)
with open(os.path.join(tmp, "rgb.ppm"), "wb") as f:
    f.write(b"P6\n%d %d\n255\n" % (W, H) + rgb)
gray = bytes(rgb[i] // 3 + rgb[i + 1] // 3 + rgb[i + 2] // 3 for i in range(0, len(rgb), 3))
with open(os.path.join(tmp, "gray.pgm"), "wb") as f:
    f.write(b"P5\n%d %d\n255\n" % (W, H) + gray)
png(os.path.join(tmp, "rgba.png"), W, H, picture(W, H, alpha=True))

jp2_lossy = opj(os.path.join(tmp, "rgb.ppm"), os.path.join(tmp, "lossy.jp2"), "-I", "-r", "30", "-n", "4")
j2k_gray = opj(os.path.join(tmp, "gray.pgm"), os.path.join(tmp, "gray.j2k"), "-p", "RPCL", "-c", "[64,64]", "-SOP", "-EPH")
jp2_rgba = opj(os.path.join(tmp, "rgba.png"), os.path.join(tmp, "rgba.jp2"), "-r", "10")

objs = []


def add(o):
    objs.append(o)
    return len(objs)


def stream(d, data):
    return ("<< %s /Length %d >>" % (d, len(data))).encode() + b"\nstream\n" + data + b"\nendstream"


images = {}
images["J1"] = add(stream("/Type /XObject /Subtype /Image /Width %d /Height %d /Filter /JPXDecode" % (W, H), jp2_lossy))
images["J2"] = add(stream("/Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceGray /BitsPerComponent 8 /Filter /JPXDecode" % (W, H), j2k_gray))
images["J3"] = add(stream("/Type /XObject /Subtype /Image /Width %d /Height %d /SMaskInData 1 /Filter /JPXDecode" % (W, H), jp2_rgba))

draw = []
labels = []
x = 30
for name, label in (("J1", "JP2, lossy 9/7"), ("J2", "J2K, lossless gray"), ("J3", "JP2, alpha (SMaskInData)")):
    draw.append("q %d 0 0 %d %d 250 cm /%s Do Q" % (W, H, x, name))
    labels.append("BT /F1 10 Tf %d 236 Td (%s) Tj ET" % (x, label))
    x += W + 20

# JBIG2 fixtures (the embedded streams; the size is in the .pbm beside them).
jbig2 = []
for base, label, mask in (("halftone_arith", "JBIG2 halftone region", False),
                          ("globals", "JBIG2 text, JBIG2Globals, stencil mask", True)):
    data = open(os.path.join(jbig2_dir, base + ".jb2"), "rb").read()
    w, h = map(int, open(os.path.join(jbig2_dir, base + ".pbm"), "rb").read().split(b"\n")[1].split())
    parms = ""
    glob = os.path.join(jbig2_dir, base + ".glob")
    if os.path.exists(glob):
        g = add(stream("", open(glob, "rb").read()))
        parms = " /DecodeParms << /JBIG2Globals %d 0 R >>" % g
    kind = "/ImageMask true" if mask else "/ColorSpace /DeviceGray /BitsPerComponent 1"
    n = add(stream("/Type /XObject /Subtype /Image /Width %d /Height %d %s /Filter /JBIG2Decode%s" % (w, h, kind, parms), data))
    jbig2.append((n, w, h, label, mask))
y = 20
x = 30
for i, (n, w, h, label, mask) in enumerate(jbig2):
    name = "B%d" % i
    images[name] = n
    scale = min(270 / w, 180 / h)
    color = "0.6 0.1 0.1 rg " if mask else ""
    draw.append("q %s%.3f 0 0 %.3f %d %d cm /%s Do Q" % (color, w * scale, h * scale, x, y + 20, name))
    labels.append("BT /F1 10 Tf %d %d Td (%s) Tj ET" % (x, y + 6, label))
    x += 290

content = ("\n".join(draw + ["0 g"] + labels)).encode()
contents = add(stream("", content))
font = add(b"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>")
xobjs = " ".join("/%s %d 0 R" % (k, v) for k, v in images.items())
page = add(("<< /Type /Page /Parent %d 0 R /MediaBox [0 0 620 400] /Contents %d 0 R /Resources << /Font << /F1 %d 0 R >> /XObject << %s >> >> >>"
            % (len(objs) + 2, contents, font, xobjs)).encode())
pages = add(("<< /Type /Pages /Kids [%d 0 R] /Count 1 >>" % page).encode())
catalog = add(("<< /Type /Catalog /Pages %d 0 R >>" % pages).encode())

buf = bytearray(b"%PDF-1.7\n%\xe2\xe3\xcf\xd3\n")
offsets = []
for i, o in enumerate(objs):
    offsets.append(len(buf))
    buf += b"%d 0 obj\n" % (i + 1) + o + b"\nendobj\n"
xref = len(buf)
buf += b"xref\n0 %d\n0000000000 65535 f \n" % (len(objs) + 1)
for off in offsets:
    buf += b"%010d 00000 n \n" % off
buf += b"trailer\n<< /Size %d /Root %d 0 R >>\nstartxref\n%d\n%%%%EOF\n" % (len(objs) + 1, catalog, xref)
open(out, "wb").write(buf)
print("wrote", out, len(buf), "bytes;", len(jbig2), "JBIG2 image(s)")
