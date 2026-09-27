"""Generates the Illustrator test files in converter/ai/testdata.

Usage: python3 test/ai/gen.py   (standard library only)

An .ai file of Illustrator 9 and later is a PDF with Illustrator's own data
beside it. These files have the parts of it the converter reads, in the
shape Illustrator writes them:

- artboards.ai: three artboards of different sizes, one page each. A page's
  media and bleed boxes include a 9 pt bleed; its trim box is the artboard.
  Three layers are optional content groups: a background that runs into
  the bleed, the artwork (paths, CMYK and RGB colours, a gradient,
  transparency, dashes, a clip), and a layer of red notes that is hidden.
  The XMP metadata comes first and says the file is an Illustrator
  document; each page carries Illustrator's private data under /PieceInfo
  (a stand-in: the converter does not read it).
- nopdf.ai: a file saved without "Create PDF Compatible File": the only
  page is the placeholder Illustrator writes instead of the artwork.
- legacy.ai: the head of a PostScript-based .ai of Illustrator 8.
"""
import math
import os
import zlib

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "ai", "testdata")
BLEED = 9

XMP = """<?xpacket begin="﻿" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
   <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
      <rdf:Description rdf:about=""
            xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:xmp="http://ns.adobe.com/xap/1.0/"
            xmlns:illustrator="http://ns.adobe.com/illustrator/1.0/">
         <dc:format>application/pdf</dc:format>
         <dc:title>
            <rdf:Alt>
               <rdf:li xml:lang="x-default">{title}</rdf:li>
            </rdf:Alt>
         </dc:title>
         <xmp:CreatorTool>test/ai/gen.py</xmp:CreatorTool>
         <xmp:CreateDate>2026-09-26T12:00:00+09:00</xmp:CreateDate>
         <illustrator:Type>Document</illustrator:Type>
      </rdf:Description>
   </rdf:RDF>
</x:xmpmeta>
<?xpacket end="w"?>"""


class PDF:
    def __init__(self):
        self.objs = []

    def add(self, body=None):
        self.objs.append(body)
        return len(self.objs)

    def set(self, n, body):
        self.objs[n - 1] = body

    def ref(self, n):
        return "%d 0 R" % n

    def stream(self, dict_entries, data, compress=True):
        if isinstance(data, str):
            data = data.encode("latin-1")
        if compress:
            data = zlib.compress(data, 9)
            dict_entries += "/Filter/FlateDecode"
        return (dict_entries, data)

    def write(self, path, root, info):
        out = bytearray(b"%PDF-1.6\n%\xe2\xe3\xcf\xd3\n")
        offsets = []
        for i, o in enumerate(self.objs):
            offsets.append(len(out))
            out += b"%d 0 obj\n" % (i + 1)
            if isinstance(o, tuple):
                entries, data = o
                out += b"<<%s/Length %d>>stream\n" % (entries.encode("latin-1"), len(data))
                out += data + b"\nendstream"
            else:
                out += o.encode("utf-8") if isinstance(o, str) else o
            out += b"\nendobj\n"
        xref = len(out)
        out += b"xref\n0 %d\n0000000000 65535 f \n" % (len(self.objs) + 1)
        for off in offsets:
            out += b"%010d 00000 n \n" % off
        out += b"trailer\n<</Size %d/Root %d 0 R/Info %d 0 R>>\nstartxref\n%d\n%%%%EOF\n" % (len(self.objs) + 1, root, info, xref)
        with open(path, "wb") as f:
            f.write(out)


def num(v):
    s = "%.3f" % v
    return s.rstrip("0").rstrip(".") if "." in s else s


def circle(cx, cy, r):
    k = 0.5523 * r
    return " ".join([
        "%s %s m" % (num(cx + r), num(cy)),
        "%s %s %s %s %s %s c" % (num(cx + r), num(cy + k), num(cx + k), num(cy + r), num(cx), num(cy + r)),
        "%s %s %s %s %s %s c" % (num(cx - k), num(cy + r), num(cx - r), num(cy + k), num(cx - r), num(cy)),
        "%s %s %s %s %s %s c" % (num(cx - r), num(cy - k), num(cx - k), num(cy - r), num(cx), num(cy - r)),
        "%s %s %s %s %s %s c" % (num(cx + k), num(cy - r), num(cx + r), num(cy - k), num(cx + r), num(cy)),
        "h",
    ])


def star(cx, cy, r1, r2, n=5):
    pts = []
    for i in range(2 * n):
        r = r1 if i % 2 == 0 else r2
        a = math.pi / 2 + i * math.pi / n
        pts.append((cx + r * math.cos(a), cy + r * math.sin(a)))
    return " ".join(("%s %s %s" % (num(x), num(y), "m" if i == 0 else "l")) for i, (x, y) in enumerate(pts)) + " h"


def rounded(x, y, w, h, r):
    k = 0.5523 * r
    return " ".join([
        "%s %s m" % (num(x + r), num(y)),
        "%s %s l" % (num(x + w - r), num(y)),
        "%s %s %s %s %s %s c" % (num(x + w - r + k), num(y), num(x + w), num(y + r - k), num(x + w), num(y + r)),
        "%s %s l" % (num(x + w), num(y + h - r)),
        "%s %s %s %s %s %s c" % (num(x + w), num(y + h - r + k), num(x + w - r + k), num(y + h), num(x + w - r), num(y + h)),
        "%s %s l" % (num(x + r), num(y + h)),
        "%s %s %s %s %s %s c" % (num(x + r - k), num(y + h), num(x), num(y + h - r + k), num(x), num(y + h - r)),
        "%s %s l" % (num(x), num(y + r)),
        "%s %s %s %s %s %s c" % (num(x), num(y + r - k), num(x + r - k), num(y), num(x + r), num(y)),
        "h",
    ])


def layer(mc, body):
    return "/OC /%s BDC\n%s\nEMC\n" % (mc, body.strip())


def artboards():
    pdf = PDF()
    catalog = pdf.add()
    meta = pdf.add(pdf.stream("/Type/Metadata/Subtype/XML", XMP.format(title="Artboards").encode("utf-8"), compress=False))
    pages = pdf.add()
    info = pdf.add("<</Title(artboards)/Creator(test/ai/gen.py)/CreationDate(D:20260926120000+09'00')>>")
    names = ["Background", "Artwork", "Notes"]
    ocgs = []
    for name in names:
        intent = pdf.add("[/View/Design]")
        usage = pdf.add("<</CreatorInfo<</Creator(Adobe Illustrator 24.0)/Subtype/Artwork>>>>")
        ocgs.append(pdf.add("<</Intent %s/Name(%s)/Type/OCG/Usage %s>>" % (pdf.ref(intent), name, pdf.ref(usage))))
    order = pdf.add("[%s]" % " ".join(pdf.ref(o) for o in ocgs))
    ai_meta = pdf.add(pdf.stream("", "%!PS-Adobe-3.0 \r%%Creator: test/ai/gen.py (a stand-in for Illustrator's private data)\r%%EndComments\r"))
    private = pdf.add("<</AIMetaData %s/ContainerVersion 12/CreatorVersion 24/NumBlock 0/RoundtripVersion 24>>" % pdf.ref(ai_meta))
    piece = pdf.add("<</Illustrator<</LastModified(D:20260926120000+09'00')/Private %s>>>>" % pdf.ref(private))
    gs_multiply = pdf.add("<</Type/ExtGState/ca 0.6/CA 0.6/BM/Multiply>>")
    gs_normal = pdf.add("<</Type/ExtGState/ca 1/CA 1/BM/Normal>>")
    gradient = pdf.add("<</ShadingType 2/ColorSpace/DeviceRGB/Coords[0 0 1 0]/Extend[true true]"
                       "/Function<</FunctionType 2/Domain[0 1]/C0[0.1 0.45 0.9]/C1[0.95 0.35 0.55]/N 1>>>>")

    boards = [(200, 200), (300, 150), (120, 240)]
    kids = []
    for i, (w, h) in enumerate(boards):
        b = BLEED
        # Background: runs into the bleed on every side.
        bg = "0.04 0.02 0 0 k %s %s %s %s re f" % (num(0), num(0), num(w + 2 * b), num(h + 2 * b))
        notes = ("1 0 0 RG 6 w %s %s m %s %s l %s %s m %s %s l S\n"
                 "1 0 0 rg %s %s 60 20 re f") % (num(b), num(b), num(b + w), num(b + h), num(b), num(b + h), num(b + w), num(b),
                                                num(b + 10), num(b + 10))
        if i == 0:
            art = "\n".join([
                "0.8 0 0.05 0 k %s f" % circle(b + 80, b + 110, 55),
                "0.13 0.55 0.13 RG 0.99 0.8 0.1 rg 3 w 1 j %s B" % star(b + 140, b + 70, 50, 22),
                # Runs over the right edge of the artboard into the bleed: cut at the trim box.
                "0 0.9 0.8 0 k %s %s 40 30 re f" % (num(b + w - 20), num(b + 20)),
            ])
        elif i == 1:
            art = "\n".join([
                "q %s W n %s 0 0 %s %s %s cm /Sh0 sh Q" % (rounded(b + 15, b + 15, 170, 120, 18), num(170), num(120), num(b + 15), num(b + 15)),
                "/GS0 gs 0 0.5 1 0 k %s f 1 0.3 0 0 k %s f /GS1 gs" % (circle(b + 225, b + 85, 42), circle(b + 250, b + 60, 42)),
            ])
        else:
            art = "\n".join([
                "0 0 0 0.85 K 4 w [8 5] 0 d 1 J %s S [] 0 d" % rounded(b + 12, b + 12, w - 24, h - 24, 14),
                "0 G 10 w 1 J %s %s m %s %s l S 2 J %s %s m %s %s l S" % (num(b + 30), num(b + 180), num(b + 90), num(b + 180),
                                                                        num(b + 30), num(b + 150), num(b + 90), num(b + 150)),
                "0.2 0.2 0.2 rg %s f" % star(b + 60, b + 80, 40, 16, 6),
            ])
        content = layer("MC0", bg) + layer("MC1", art) + layer("MC2", notes)
        stream = pdf.add(pdf.stream("", content))
        res = ("<</ExtGState<</GS0 %s/GS1 %s>>/Shading<</Sh0 %s>>/Properties<</MC0 %s/MC1 %s/MC2 %s>>>>"
               % (pdf.ref(gs_multiply), pdf.ref(gs_normal), pdf.ref(gradient), pdf.ref(ocgs[0]), pdf.ref(ocgs[1]), pdf.ref(ocgs[2])))
        box = lambda x0, y0, x1, y1: "[%s %s %s %s]" % (num(x0), num(y0), num(x1), num(y1))
        page = pdf.add("<</Type/Page/Parent %s/MediaBox %s/BleedBox %s/TrimBox %s/ArtBox %s/Contents %s/Resources %s/PieceInfo %s"
                       "/LastModified(D:20260926120000+09'00')>>"
                       % (pdf.ref(pages), box(0, 0, w + 2 * b, h + 2 * b), box(0, 0, w + 2 * b, h + 2 * b), box(b, b, b + w, b + h),
                          box(b, b, b + w, b + h), pdf.ref(stream), res, pdf.ref(piece)))
        kids.append(page)
    pdf.set(pages, "<</Count %d/Kids[%s]/Type/Pages>>" % (len(kids), " ".join(pdf.ref(k) for k in kids)))
    pdf.set(catalog, "<</Metadata %s/OCProperties<</D<</OFF[%s]/ON[%s %s]/Order %s/RBGroups[]>>/OCGs[%s]>>/Pages %s/Type/Catalog>>"
            % (pdf.ref(meta), pdf.ref(ocgs[2]), pdf.ref(ocgs[0]), pdf.ref(ocgs[1]), pdf.ref(order),
               " ".join(pdf.ref(o) for o in ocgs), pdf.ref(pages)))
    pdf.write(os.path.join(OUT, "artboards.ai"), catalog, info)


def nopdf():
    pdf = PDF()
    catalog = pdf.add()
    meta = pdf.add(pdf.stream("/Type/Metadata/Subtype/XML", XMP.format(title="No PDF content").encode("utf-8"), compress=False))
    pages = pdf.add()
    info = pdf.add("<</Title(nopdf)/Creator(test/ai/gen.py)>>")
    font = pdf.add("<</Type/Font/Subtype/Type1/BaseFont/Helvetica/Encoding/WinAnsiEncoding>>")
    lines = [
        "This is an Adobe\xae Illustrator\xae File that was",
        "saved without PDF Content.",
        "To place or open this file in other",
        "applications, it should be re-saved from",
        "Adobe Illustrator with the \"Create PDF",
        "Compatible File\" option turned on.",
    ]
    text = "BT /F1 14 Tf 16 TL 40 740 Td " + " ".join("(%s) '" % l.replace("(", "\\(").replace(")", "\\)") for l in lines) + " ET"
    stream = pdf.add(pdf.stream("", text))
    page = pdf.add("<</Type/Page/Parent %s/MediaBox[0 0 612 792]/Contents %s/Resources<</Font<</F1 %s>>>>>>"
                   % (pdf.ref(pages), pdf.ref(stream), pdf.ref(font)))
    pdf.set(pages, "<</Count 1/Kids[%s]/Type/Pages>>" % pdf.ref(page))
    pdf.set(catalog, "<</Metadata %s/Pages %s/Type/Catalog>>" % (pdf.ref(meta), pdf.ref(pages)))
    pdf.write(os.path.join(OUT, "nopdf.ai"), catalog, info)


def legacy():
    with open(os.path.join(OUT, "legacy.ai"), "wb") as f:
        f.write(b"%!PS-Adobe-3.0 \r%%Creator: Adobe Illustrator(R) 8.0\r%%AI8_CreatorVersion: 8.0.1\r"
                b"%%For: (test/ai/gen.py) ()\r%%Title: (legacy.ai)\r%%BoundingBox: 0 0 100 100\r%%EndComments\r"
                b"%%BeginProlog\r%%EndProlog\r%%BeginSetup\r%%EndSetup\r"
                b"0 0 100 100 re f\r%%PageTrailer\r%%Trailer\r%%EOF\r")


if __name__ == "__main__":
    os.makedirs(OUT, exist_ok=True)
    artboards()
    nopdf()
    legacy()
