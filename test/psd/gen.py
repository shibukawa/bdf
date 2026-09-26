"""Generates the Photoshop test documents in converter/psd/testdata.

Usage: python3 test/psd/gen.py   (standard library only)

The documents are written the way Photoshop writes them, with a composite
image that this script computes from the layers (straight alpha, over white
where it is transparent, as Photoshop stores it). The Go tests check the
converter's own compositing, used for files without a composite, against
it.

- layers.psd: RGB, 8 bits, 144 ppi, with a background, an isolated group
  at 80 % with a multiply layer, a layer clipped to a circle, a layer mask,
  a solid colour fill layer with a mask, a pass-through group at 60 %, a
  screen layer, a hidden layer and fill opacity; XMP metadata and Unicode
  layer names.
- artboards.psd: RGB, 8 bits, 72 ppi, four artboards (white, coloured and
  transparent backgrounds, one hidden) on a transparent canvas; a layer
  runs over the edge of its artboard.
- Small documents for the other colour modes, depths and encodings:
  gray16.psd (grayscale, 16 bits, layers in Lr16 with ZIP and prediction,
  RLE composite), rgb32.psd (32-bit float, layers in Lr32 with ZIP and
  prediction), cmyk.psd (CMYK with transparency), lab.psd, indexed.psd (a
  colour table with a transparent index), bitmap.psd (1 bit) and
  large.psb (the large document format).
"""
import math
import os
import struct
import zlib

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "psd", "testdata")

RAW, RLE, ZIP, ZIPPRED = 0, 1, 2, 3

# --- drawing ------------------------------------------------------------------


def coverage(inside, x, y, n=4):
    """Anti-aliased coverage of pixel (x, y) by a shape, n × n samples."""
    hit = 0
    for j in range(n):
        for i in range(n):
            if inside(x + (i + 0.5) / n, y + (j + 0.5) / n):
                hit += 1
    return hit / (n * n)


def circle(cx, cy, r):
    return lambda x, y: (x - cx) ** 2 + (y - cy) ** 2 <= r * r


def box(x0, y0, x1, y1):
    return lambda x, y: x0 <= x < x1 and y0 <= y < y1


def pixels(w, h, fn):
    """A layer's pixels: fn(x, y) -> (r, g, b, a), floats 0-1, row by row."""
    return [fn(x, y) for y in range(h) for x in range(w)]


class Layer:
    def __init__(self, name, rect=(0, 0, 0, 0), px=None, blend="norm", opacity=255, fill=255, clipping=False,
                 hidden=False, mask=None, fill_color=None):
        self.name = name
        self.rect = rect  # top, left, bottom, right
        self.px = px  # pixels over rect
        self.blend = blend
        self.opacity = opacity
        self.fill = fill
        self.clipping = clipping
        self.hidden = hidden
        self.mask = mask  # (rect, default, values over rect)
        self.fill_color = fill_color  # (r, g, b) 0-255: a solid colour fill layer
        self.group = None  # a list of children, bottom to top
        self.artboard = None  # (rect, background type, colour)


def group(name, children, blend="pass", opacity=255, hidden=False, mask=None, artboard=None):
    g = Layer(name, blend=blend, opacity=opacity, hidden=hidden, mask=mask)
    g.group = children
    g.artboard = artboard
    return g


# --- the reference compositor (straight colours, floats) -------------------

BLENDS = {
    "norm": lambda cb, cs: cs,
    "mul ": lambda cb, cs: cb * cs,
    "scrn": lambda cb, cs: cb + cs - cb * cs,
}


def mask_value(m, x, y):
    if m is None:
        return 1.0
    (t, l, b, r), default, values = m
    if l <= x < r and t <= y < b:
        return values[(y - t) * (r - l) + x - l] / 255
    return default / 255


def empty(w, h):
    return [[0.0, 0.0, 0.0, 0.0] for _ in range(w * h)]  # premultiplied


def blend_onto(dst, src, w, h, mode, opacity, atop=False):
    fn = BLENDS[mode] if mode in BLENDS else BLENDS["norm"]
    for i in range(w * h):
        s, d = src[i], dst[i]
        sa = s[3] * opacity
        if sa <= 0:
            continue
        ba = d[3]
        cs = [s[c] / s[3] for c in range(3)]
        cb = [d[c] / ba if ba > 0 else 0 for c in range(3)]
        mixed = [fn(cb[c], cs[c]) if ba > 0 else cs[c] for c in range(3)]
        for c in range(3):
            v = sa * ba * mixed[c] + (1 - sa) * ba * cb[c]
            if not atop:
                v += sa * (1 - ba) * cs[c]
            d[c] = v
        if not atop:
            d[3] = sa + ba * (1 - sa)


def render(layer, w, h):
    buf = empty(w, h)
    if layer.group is not None:
        if layer.artboard:
            (t, l, b, r), bg, color = layer.artboard
            c = {1: (1, 1, 1, 1), 2: (0, 0, 0, 1), 3: (0, 0, 0, 0), 4: (color[0] / 255, color[1] / 255, color[2] / 255, 1)}[bg]
            for y in range(max(t, 0), min(b, h)):
                for x in range(max(l, 0), min(r, w)):
                    buf[y * w + x] = list(c)
            draw_list(buf, layer.group, w, h)
            for y in range(h):
                for x in range(w):
                    if not (l <= x < r and t <= y < b):
                        buf[y * w + x] = [0.0, 0.0, 0.0, 0.0]
        else:
            draw_list(buf, layer.group, w, h)
    elif layer.fill_color is not None:
        c = [v / 255 for v in layer.fill_color]
        for i in range(w * h):
            buf[i] = [c[0], c[1], c[2], 1.0]
    else:
        t, l, b, r = layer.rect
        for y in range(max(t, 0), min(b, h)):
            for x in range(max(l, 0), min(r, w)):
                p = layer.px[(y - t) * (r - l) + x - l]
                buf[y * w + x] = [p[0] * p[3], p[1] * p[3], p[2] * p[3], p[3]]
    k = layer.fill / 255 if layer.group is None else 1
    for y in range(h):
        for x in range(w):
            m = k * mask_value(layer.mask, x, y)
            if m != 1:
                buf[y * w + x] = [v * m for v in buf[y * w + x]]
    return buf


def draw_list(dst, layers, w, h):
    i = 0
    while i < len(layers):
        j = i + 1
        while j < len(layers) and layers[j].clipping:
            j += 1
        base, clips = layers[i], layers[i + 1:j]
        i = j
        if base.hidden:
            continue
        if base.group is not None and base.artboard is None and base.blend == "pass" and not any(not c.hidden for c in clips):
            if base.opacity == 255 and base.mask is None:
                draw_list(dst, base.group, w, h)
                continue
            t = [list(p) for p in dst]
            draw_list(t, base.group, w, h)
            for y in range(h):
                for x in range(w):
                    k = base.opacity / 255 * mask_value(base.mask, x, y)
                    d, s = dst[y * w + x], t[y * w + x]
                    for c in range(4):
                        d[c] += (s[c] - d[c]) * k
            continue
        buf = render(base, w, h)
        for c in clips:
            if not c.hidden:
                blend_onto(buf, render(c, w, h), w, h, c.blend, c.opacity / 255, atop=True)
        blend_onto(dst, buf, w, h, base.blend, base.opacity / 255)


def composite(layers, w, h):
    """The composite as straight RGBA bytes."""
    dst = empty(w, h)
    draw_list(dst, layers, w, h)
    out = []
    for p in dst:
        a = p[3]
        if a <= 0:
            out.append((0, 0, 0, 0))
        else:
            out.append(tuple(max(0, min(255, round(p[c] / a * 255))) for c in range(3)) + (max(0, min(255, round(a * 255))),))
    return out


# --- encoding ---------------------------------------------------------------


def packbits(row):
    out = bytearray()
    i = 0
    while i < len(row):
        j = i
        while j < len(row) and j - i < 128 and row[j] == row[i]:
            j += 1
        if j - i >= 3:
            out += bytes([257 - (j - i), row[i]])
            i = j
            continue
        j = i
        while j < len(row) and j - i < 128 and not (j + 2 < len(row) and row[j] == row[j + 1] == row[j + 2]):
            j += 1
        out += bytes([j - i - 1]) + bytes(row[i:j])
        i = j
    return bytes(out)


def srgb_to_linear(v):
    return v / 12.92 if v <= 0.04045 else ((v + 0.055) / 1.055) ** 2.4


def sample_bytes(values, depth, linear):
    """Samples 0-255 (floats allowed) to bytes of the depth, row-major."""
    if depth == 8:
        return bytes(int(round(v)) for v in values)
    if depth == 16:
        return b"".join(struct.pack(">H", int(round(v / 255 * 65535))) for v in values)
    if depth == 32:
        return b"".join(struct.pack(">f", srgb_to_linear(v / 255) if linear else v / 255) for v in values)
    raise ValueError(depth)


def bits(values, w, h):
    """1-bit rows: 1 is black (values below 128)."""
    out = bytearray()
    for y in range(h):
        row = bytearray((w + 7) // 8)
        for x in range(w):
            if values[y * w + x] < 128:
                row[x >> 3] |= 0x80 >> (x & 7)
        out += row
    return bytes(out)


def predict(raw, w, h, depth):
    rb = len(raw) // h if h else 0
    out = bytearray()
    for y in range(h):
        row = bytearray(raw[y * rb:(y + 1) * rb])
        if depth == 8:
            enc = bytearray(row)
            for x in range(len(row) - 1, 0, -1):
                enc[x] = (row[x] - row[x - 1]) & 0xFF
        elif depth == 16:
            vals = [struct.unpack(">H", row[2 * x:2 * x + 2])[0] for x in range(w)]
            enc = bytearray()
            for x in range(w):
                enc += struct.pack(">H", (vals[x] - (vals[x - 1] if x else 0)) & 0xFFFF)
        else:  # 32: bytes plane by plane, then deltas over the row
            planes = bytearray()
            for b in range(4):
                planes += bytes(row[4 * x + b] for x in range(w))
            enc = bytearray(planes)
            for x in range(len(planes) - 1, 0, -1):
                enc[x] = (planes[x] - planes[x - 1]) & 0xFF
        out += enc
    return bytes(out)


def encode_channel(raw, w, h, depth, comp, psb):
    """A layer channel: compression code and data."""
    if comp == RAW or w == 0 or h == 0:
        return struct.pack(">H", RAW) + raw
    if comp == RLE:
        rb = len(raw) // h
        rows = [packbits(raw[y * rb:(y + 1) * rb]) for y in range(h)]
        fmt = ">I" if psb else ">H"
        return struct.pack(">H", RLE) + b"".join(struct.pack(fmt, len(r)) for r in rows) + b"".join(rows)
    if comp == ZIP:
        return struct.pack(">H", ZIP) + zlib.compress(raw, 9)
    return struct.pack(">H", ZIPPRED) + zlib.compress(predict(raw, w, h, depth), 9)


def pascal(s, pad):
    b = s.encode("latin-1", "replace")[:255]
    out = bytes([len(b)]) + b
    while len(out) % pad:
        out += b"\0"
    return out


def unicode_string(s):
    u = s.encode("utf-16-be")
    return struct.pack(">I", len(u) // 2) + u


def key_id(k):
    if len(k) == 4:
        return struct.pack(">I", 0) + k.encode()
    return struct.pack(">I", len(k)) + k.encode()


def descriptor(cls, items):
    out = unicode_string("") + key_id(cls) + struct.pack(">I", len(items))
    for k, t, v in items:
        out += key_id(k) + t.encode()
        if t == "Objc":
            out += descriptor(*v)
        elif t == "doub":
            out += struct.pack(">d", v)
        elif t == "long":
            out += struct.pack(">i", v)
        elif t == "TEXT":
            out += unicode_string(v)
        elif t == "bool":
            out += bytes([1 if v else 0])
        elif t == "VlLs":
            out += struct.pack(">I", 0)
    return out


def rgbc(c):
    return ("RGBC", [("Rd  ", "doub", float(c[0])), ("Grn ", "doub", float(c[1])), ("Bl  ", "doub", float(c[2]))])


def block(key, data, psb=False, wide=False):
    n = struct.pack(">Q", len(data)) if psb and wide else struct.pack(">I", len(data))
    out = b"8BIM" + key.encode() + n + data
    if len(data) % 2:
        out += b"\0"
    return out


# --- the file ---------------------------------------------------------------

MODE_CHANNELS = {0: 1, 1: 1, 2: 1, 3: 3, 4: 4, 9: 3}


def flatten(layers):
    """Layer records bottom to top: groups as a divider, children, the group."""
    out = []
    for l in layers:
        if l.group is not None:
            out.append(("divider", l))
            out += flatten(l.group)
            out.append(("group", l))
        else:
            out.append(("layer", l))
    return out


def to_mode(rgb, mode):
    """An RGB colour (0-255) as samples of the mode."""
    r, g, b = rgb
    if mode == 3:
        return [r, g, b]
    if mode == 1 or mode == 0:
        return [round(0.299 * r + 0.587 * g + 0.114 * b)]
    if mode == 4:  # inverted CMYK: 255 is no ink
        c, m, y = 1 - r / 255, 1 - g / 255, 1 - b / 255
        k = min(c, m, y)
        if k >= 1:
            return [255, 255, 255, 0]
        return [round((1 - (c - k) / (1 - k)) * 255), round((1 - (m - k) / (1 - k)) * 255),
                round((1 - (y - k) / (1 - k)) * 255), round((1 - k) * 255)]
    raise ValueError(mode)


def layer_record(kind, l, mode, depth, comp, psb):
    """The record of a layer and its channel data."""
    nch = MODE_CHANNELS[mode]
    t, left, b, r = l.rect if kind == "layer" and l.fill_color is None else (0, 0, 0, 0)
    w, h = r - left, b - t
    channels = []
    if kind == "layer" and l.fill_color is None:
        alpha = sample_bytes([p[3] * 255 for p in l.px], depth, False)
        cols = [to_mode([round(p[0] * 255), round(p[1] * 255), round(p[2] * 255)], mode) for p in l.px]
        channels.append((-1, encode_channel(alpha, w, h, depth, comp, psb)))
        for c in range(nch):
            channels.append((c, encode_channel(sample_bytes([v[c] for v in cols], depth, True), w, h, depth, comp, psb)))
    elif kind == "layer":
        for c in [-1] + list(range(nch)):
            channels.append((c, struct.pack(">H", RAW)))
    else:
        for c in [-1] + list(range(nch)):
            channels.append((c, struct.pack(">H", RAW)))
    mask = b""
    if l.mask is not None and kind != "divider":
        (mt, ml, mb, mr), default, values = l.mask
        mask = struct.pack(">iiiiBB", mt, ml, mb, mr, default, 0) + b"\0\0"
        channels.append((-2, encode_channel(sample_bytes(values, depth, False), mr - ml, mb - mt, depth, comp, psb)))
    lenfmt = ">Q" if psb else ">I"
    rec = struct.pack(">iiiiH", t, left, b, r, len(channels))
    for cid, data in channels:
        rec += struct.pack(">h", cid) + struct.pack(lenfmt, len(data))
    blend = l.blend
    name = l.name
    extra_blocks = b""
    if kind == "divider":
        blend, name = "norm", "</Layer group>"
        extra_blocks += block("lsct", struct.pack(">I", 3))
    elif kind == "group":
        extra_blocks += block("lsct", struct.pack(">I", 1) + b"8BIM" + l.blend.encode())
        if l.artboard:
            (at, al, ab, ar), bg, color = l.artboard
            d = descriptor("artboard", [
                ("artboardRect", "Objc", ("classFloatRect", [("Top ", "doub", float(at)), ("Left", "doub", float(al)),
                                                              ("Btom", "doub", float(ab)), ("Rght", "doub", float(ar))])),
                ("guideIndeces", "VlLs", None),
                ("artboardPresetName", "TEXT", ""),
                ("Clr ", "Objc", rgbc(color)),
                ("artboardBackgroundType", "long", bg),
            ])
            extra_blocks += block("artb", struct.pack(">I", 16) + d)
    if kind != "divider":
        extra_blocks = block("luni", unicode_string(name)) + extra_blocks
        if l.fill != 255:
            extra_blocks += block("iOpa", bytes([l.fill, 0, 0, 0]))
        if l.fill_color is not None:
            extra_blocks += block("SoCo", struct.pack(">I", 16) + descriptor("null", [("Clr ", "Objc", rgbc(l.fill_color))]))
    flags = 8 | (2 if l.hidden and kind != "divider" else 0)
    opacity = 255 if kind == "divider" else l.opacity
    clipping = 1 if l.clipping and kind != "divider" else 0
    extra = struct.pack(">I", len(mask)) + mask + struct.pack(">I", 0) + pascal(name, 4) + extra_blocks
    rec += b"8BIM" + blend.encode() + struct.pack(">BBBB", opacity, clipping, flags, 0) + struct.pack(">I", len(extra)) + extra
    return rec, b"".join(d for _, d in channels)


def layer_info(layers, mode, depth, comp, psb, global_alpha):
    recs = flatten(layers)
    body = struct.pack(">h", -len(recs) if global_alpha else len(recs))
    data = b""
    for kind, l in recs:
        rec, ch = layer_record(kind, l, mode, depth, comp, psb)
        body += rec
        data += ch
    body += data
    if len(body) % 2:
        body += b"\0"
    return body


def resources(res):
    out = b""
    for rid in sorted(res):
        data = res[rid]
        out += b"8BIM" + struct.pack(">H", rid) + b"\0\0" + struct.pack(">I", len(data)) + data
        if len(data) % 2:
            out += b"\0"
    return out


def version_info(merged):
    return struct.pack(">IB", 1, 1 if merged else 0) + unicode_string("Adobe Photoshop") + unicode_string("test/psd/gen.py") + struct.pack(">I", 1)


def resolution(ppi):
    v = int(ppi * 65536)
    return struct.pack(">IHHIHH", v, 1, 1, v, 1, 1)


XMP = """<?xpacket begin="﻿" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
 <rdf:Description rdf:about="" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xmp="http://ns.adobe.com/xap/1.0/"
   xmp:CreateDate="2026-09-26T12:00:00+09:00">
  <dc:format>image/vnd.adobe.photoshop</dc:format>
  <dc:title><rdf:Alt><rdf:li xml:lang="x-default">{title}</rdf:li></rdf:Alt></dc:title>
  <dc:creator><rdf:Seq><rdf:li>test/psd/gen.py</rdf:li></rdf:Seq></dc:creator>
 </rdf:Description>
</rdf:RDF></x:xmpmeta>
<?xpacket end="w"?>"""


def write(name, w, h, mode, depth, planes, layers=None, res=None, psb=False, comp=RLE, layer_comp=RLE,
          global_alpha=False, deep_block=None, palette=None):
    """Writes a document; planes are the composite's channels (sample lists)."""
    out = b"8BPS" + struct.pack(">H", 2 if psb else 1) + b"\0" * 6
    out += struct.pack(">HIIHH", len(planes), h, w, depth, mode)
    cm = palette or b""
    out += struct.pack(">I", len(cm)) + cm
    r = resources(res or {})
    out += struct.pack(">I", len(r)) + r
    lenfmt = ">Q" if psb else ">I"
    if layers:
        info = layer_info(layers, mode, depth, layer_comp, psb, global_alpha)
        if deep_block:
            # 16- and 32-bit documents keep their layers in Lr16 or Lr32.
            lm = struct.pack(lenfmt, 0) + struct.pack(">I", 0)
            b = block(deep_block, info, psb, wide=True)
            while len(b) % 4:
                b += b"\0"
            lm += b
        else:
            lm = struct.pack(lenfmt, len(info)) + info + struct.pack(">I", 0)
        out += struct.pack(lenfmt, len(lm)) + lm
    else:
        out += struct.pack(lenfmt, 0)
    # Image data: the composite, channel after channel.
    raws = []
    for i, p in enumerate(planes):
        linear = i < MODE_CHANNELS.get(mode, 1) and mode in (1, 3)
        raws.append(bits(p, w, h) if depth == 1 else sample_bytes(p, depth, linear))
    out += struct.pack(">H", comp)
    if comp == RAW:
        out += b"".join(raws)
    elif comp == RLE:
        rb = (w * depth + 7) // 8
        rows = [packbits(raw[y * rb:(y + 1) * rb]) for raw in raws for y in range(h)]
        fmt = ">I" if psb else ">H"
        out += b"".join(struct.pack(fmt, len(x)) for x in rows) + b"".join(rows)
    else:
        out += zlib.compress(b"".join(raws), 9)
    with open(os.path.join(OUT, name), "wb") as f:
        f.write(out)


def composite_planes(layers, w, h, mode, alpha):
    """The composite's channels: colour (over white where transparent) and alpha."""
    px = composite(layers, w, h)
    planes = [[] for _ in range(MODE_CHANNELS[mode] + (1 if alpha else 0))]
    for r, g, b, a in px:
        k = a / 255
        matted = [r * k + 255 * (1 - k), g * k + 255 * (1 - k), b * k + 255 * (1 - k)] if alpha else [r, g, b]
        vals = to_mode([round(v) for v in matted], mode)
        for c, v in enumerate(vals):
            planes[c].append(v)
        if alpha:
            planes[-1].append(a)
    return planes


# --- the documents ------------------------------------------------------------


def layers_psd():
    w, h = 160, 120

    def bg(x, y):
        return (0.93 - 0.3 * y / h, 0.95 - 0.1 * x / w, 1.0, 1.0)

    def disc(cx, cy, r, color, alpha=1.0):
        inside = circle(cx, cy, r)
        x0, y0 = int(cx - r) - 1, int(cy - r) - 1
        size = int(2 * r) + 3
        return (y0, x0, y0 + size, x0 + size), pixels(size, size, lambda x, y: color + (alpha * coverage(inside, x0 + x, y0 + y),))

    def stripes(x, y):
        return (1.0, 0.85, 0.2, 1.0) if (x + y) // 6 % 2 == 0 else (0.95, 0.4, 0.1, 1.0)

    rect, px = disc(50, 55, 34, (0.2, 0.4, 0.9))
    circle_layer = Layer("円", rect, px)
    clipped = Layer("Stripes", (0, 0, h, 100), pixels(100, h, stripes), clipping=True, opacity=200)
    mrect, mpx = disc(95, 45, 22, (0.9, 0.2, 0.5))
    multiply = Layer("Multiply", mrect, mpx, blend="mul ")
    shapes = group("Shapes", [circle_layer, clipped, multiply], blend="norm", opacity=204)
    masked = Layer("Masked", (70, 10, 110, 150), pixels(140, 40, lambda x, y: (0.1, 0.6, 0.3, 1.0)),
                   mask=((70, 10, 110, 150), 0, [round(255 * x / 139) for y in range(40) for x in range(140)]))
    fill = Layer("Fill", fill_color=(250, 200, 40), opacity=128,
                 mask=((10, 120, 60, 155), 0, [255 if (x - 17) ** 2 + (y - 25) ** 2 < 300 else 0 for y in range(50) for x in range(35)]))
    srect, spx = disc(130, 95, 20, (0.3, 0.3, 0.8))
    screen = Layer("Screen", srect, spx, blend="scrn", fill=153)
    hidden = Layer("Hidden", (0, 0, h, w), pixels(w, h, lambda x, y: (1.0, 0.0, 0.0, 1.0)), hidden=True)
    prect, ppx = disc(25, 100, 16, (0.5, 0.1, 0.6))
    passthrough = group("Pass through", [Layer("Purple", prect, ppx)], blend="pass", opacity=153)
    layers = [Layer("Background", (0, 0, h, w), pixels(w, h, bg)), shapes, masked, fill, screen, passthrough, hidden]
    res = {1005: resolution(144), 1057: version_info(True), 1060: XMP.format(title="Layers").encode("utf-8")}
    write("layers.psd", w, h, 3, 8, composite_planes(layers, w, h, 3, False), layers, res)


def artboards_psd():
    w, h = 400, 260

    def disc_layer(name, cx, cy, r, color):
        inside = circle(cx, cy, r)
        x0, y0 = int(cx - r) - 1, int(cy - r) - 1
        size = int(2 * r) + 3
        return Layer(name, (y0, x0, y0 + size, x0 + size), pixels(size, size, lambda x, y: color + (coverage(inside, x0 + x, y0 + y),)))

    def bar_layer(name, x0, y0, x1, y1, color):
        return Layer(name, (y0, x0, y1, x1), pixels(x1 - x0, y1 - y0, lambda x, y: color + (1.0,)))

    a1 = group("Artboard 1", [disc_layer("Sun", 60, 60, 35, (0.95, 0.6, 0.1)), bar_layer("Ground", 0, 110, 200, 150, (0.2, 0.6, 0.25)),
                              disc_layer("Over the edge", 190, 90, 30, (0.3, 0.3, 0.9))],
               artboard=((0, 0, 150, 200), 1, (255, 255, 255)))
    a2 = group("Artboard 2", [bar_layer("Bar 1", 240, 20, 380, 40, (1, 1, 1)), bar_layer("Bar 2", 240, 55, 340, 75, (1, 1, 1)),
                              disc_layer("Dot", 360, 95, 14, (0.9, 0.2, 0.3))],
               artboard=((0, 220, 120, 400), 4, (40, 70, 140)))
    a3 = group("Artboard 3", [disc_layer("Ring", 300, 200, 45, (0.1, 0.7, 0.7)), disc_layer("Hole", 300, 200, 25, (1, 1, 1))],
               artboard=((140, 220, 260, 380), 3, (255, 255, 255)))
    a4 = group("Hidden artboard", [bar_layer("Hidden bar", 10, 180, 190, 250, (0.8, 0.1, 0.1))],
               artboard=((170, 0, 260, 200), 2, (0, 0, 0)), hidden=True)
    layers = [a1, a2, a3, a4]
    res = {1005: resolution(72), 1057: version_info(True), 1060: XMP.format(title="Artboards").encode("utf-8")}
    write("artboards.psd", w, h, 3, 8, composite_planes(layers, w, h, 3, True), layers, res, global_alpha=True)


def small_layers(w, h):
    """Two layers for the small documents: an opaque gradient and a disc."""
    inside = circle(w * 0.6, h * 0.5, h * 0.4)
    return [
        Layer("Gradient", (0, 0, h, w), pixels(w, h, lambda x, y: (x / (w - 1), 0.5, 1 - x / (w - 1), 1.0))),
        Layer("Disc", (0, 0, h, w), pixels(w, h, lambda x, y: (0.9, 0.9, 0.2, coverage(inside, x, y))), opacity=191),
    ]


def small_psds():
    w, h = 24, 12
    res = {1005: resolution(72), 1057: version_info(True)}
    layers = small_layers(w, h)
    # grayscale, 16 bits: layers in Lr16, ZIP with prediction; RLE composite
    gl = composite_planes(layers, w, h, 1, False)
    write("gray16.psd", w, h, 1, 16, gl, layers, res, comp=RLE, layer_comp=ZIPPRED, deep_block="Lr16")
    # RGB, 32 bits: layers in Lr32, ZIP with prediction; raw composite
    write("rgb32.psd", w, h, 3, 32, composite_planes(layers, w, h, 3, False), layers, res, comp=RAW, layer_comp=ZIPPRED,
          deep_block="Lr32")
    # CMYK with transparency: only the disc, on nothing
    disc = [layers[1]]
    write("cmyk.psd", w, h, 4, 8, composite_planes(disc, w, h, 4, True), disc, res, comp=RLE, layer_comp=ZIP, global_alpha=True)
    # the large document format, with an RLE composite and RLE layers
    write("large.psb", w, h, 3, 8, composite_planes(layers, w, h, 3, False), layers, res, psb=True)
    # Lab, no layers: white, black, and a strong red (L 54, a 81, b 70 in D50)
    lab = [[255, 0, round(54 * 255 / 100)] * (w * h // 3), [128, 128, 128 + 81] * (w * h // 3), [128, 128, 128 + 70] * (w * h // 3)]
    write("lab.psd", w, h, 9, 8, lab, None, res, comp=RAW)
    # indexed: a four-colour table, index 3 transparent
    table = bytearray(768)
    for i, (r, g, b) in enumerate([(255, 0, 0), (0, 255, 0), (0, 0, 255), (255, 255, 255)]):
        table[i], table[256 + i], table[512 + i] = r, g, b
    idx = [[(x // 6) % 4 for y in range(h) for x in range(w)]]
    write("indexed.psd", w, h, 2, 8, idx, None, {**res, 1047: struct.pack(">H", 3)}, comp=RLE, palette=bytes(table))
    # bitmap: a checkerboard of 4-pixel squares, black top left
    bm = [[0 if (x // 4 + y // 4) % 2 == 0 else 255 for y in range(h) for x in range(w)]]
    write("bitmap.psd", w, h, 0, 1, bm, None, res, comp=RAW)


if __name__ == "__main__":
    os.makedirs(OUT, exist_ok=True)
    layers_psd()
    artboards_psd()
    small_psds()
