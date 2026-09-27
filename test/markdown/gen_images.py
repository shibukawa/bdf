"""Writes the pictures of the Markdown and HTML test documents
(converter/markdown/testdata/images, converter/html/testdata/images) with
the standard library only: PNG files of gradients and simple shapes."""
import os
import struct
import zlib

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))


def png(path, w, h, pixel):
    rows = b"".join(b"\x00" + b"".join(bytes(pixel(x, y)) for x in range(w)) for y in range(h))

    def chunk(kind, data):
        c = kind + data
        return struct.pack(">I", len(data)) + c + struct.pack(">I", zlib.crc32(c) & 0xFFFFFFFF)

    data = b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
    data += chunk(b"IDAT", zlib.compress(rows, 9)) + chunk(b"IEND", b"")
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "wb") as f:
        f.write(data)


def gradient(x, y):
    # a landscape: sky, sun, hills
    w, h = 640, 320
    if (x - 480) ** 2 + (y - 90) ** 2 < 45 ** 2:
        return (250, 200, 60)
    hill = 230 + 30 * ((x - 320) / 320) ** 2 - 40 * (1 - ((x - 200) / 200) ** 2 if x < 400 else 0)
    if y > hill:
        return (70, 140 - (y - 230) // 3, 80)
    return (110 + y // 4, 170 + y // 5, 235)


def icon(x, y):
    # a blue circle with a white square
    if 20 <= x < 44 and 20 <= y < 44:
        return (255, 255, 255)
    if (x - 31.5) ** 2 + (y - 31.5) ** 2 < 30 ** 2:
        return (9, 105, 218)
    return (255, 255, 255)


def diagram(x, y):
    # three bars
    for i, (hgt, col) in enumerate([(120, (9, 105, 218)), (180, (26, 127, 55)), (90, (207, 34, 46))]):
        x0 = 40 + i * 110
        if x0 <= x < x0 + 80 and 220 - hgt <= y < 220:
            return col
    if y == 220 and 20 <= x < 380:
        return (80, 80, 80)
    return (255, 255, 255)


png(os.path.join(ROOT, "converter/markdown/testdata/images/photo.png"), 640, 320, gradient)
png(os.path.join(ROOT, "converter/markdown/testdata/images/icon.png"), 64, 64, icon)
png(os.path.join(ROOT, "converter/html/testdata/images/chart.png"), 400, 240, diagram)
png(os.path.join(ROOT, "converter/html/testdata/images/icon.png"), 64, 64, icon)
