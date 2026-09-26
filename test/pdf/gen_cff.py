"""Generates test PDFs whose fonts are CFF/OpenType rather than TrueType, the
way HTML-to-PDF tools and cairo embed web fonts and OpenType fonts:

- weasyprint-webfonts.pdf: @font-face web fonts. An OpenType CFF font loaded
  as .otf and as .woff2 is embedded as FontFile3/OpenType (CIDFontType0).
  TrueType copies whose OS/2 fsType says Restricted License, No subsetting
  and Preview & Print exercise the license checks.
- cairo-cff.pdf: the same CFF font embedded by cairo as a bare Type1C simple
  font (Latin text) and a bare CIDFontType0C composite font (Thai text).

Needs weasyprint, cairocffi and fontTools (with brotli), the Loma OpenType
font (fonts-tlwg-loma-otf; GPL-2+ with the font embedding exception) and the
fixture DejaVu subset."""
import os, tempfile
import cairocffi as cairo
import weasyprint
from fontTools.ttLib import TTFont

root = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
out = os.path.join(root, "converter", "pdf", "testdata")
loma = "/usr/share/fonts/opentype/tlwg/Loma.otf"
dejavu = os.path.join(root, "fixture", "fonts", "DejaVuSans-sub.ttf")

tmp = tempfile.mkdtemp()
f = TTFont(loma)
f.flavor = "woff2"
f.save(os.path.join(tmp, "Loma.woff2"))
licenses = {"LicRestricted": 0x0002, "LicNoSubset": 0x0100, "LicPreview": 0x0004}
for name, fs in licenses.items():
    f = TTFont(dejavu)
    f["OS/2"].fsType = fs
    f.save(os.path.join(tmp, name + ".ttf"))

faces = [f'@font-face {{ font-family: LomaOTF; src: url("file://{loma}") format("opentype"); }}',
         '@font-face { font-family: LomaWOFF2; src: url(Loma.woff2) format("woff2"); }']
faces += [f'@font-face {{ font-family: {n}; src: url({n}.ttf) format("truetype"); }}' for n in licenses]
html = """<!doctype html><html><head><meta charset="utf-8"><style>
@page { size: 480pt 300pt; margin: 24pt }
%s
body { font-size: 18pt; margin: 0 }
p { margin: 0 0 10pt }
</style></head><body>
<p style="font-family: LomaOTF">OpenType CFF (.otf): The quick brown fox</p>
<p style="font-family: LomaWOFF2">Web font (.woff2): jumps over the lazy dog</p>
<p style="font-family: LomaOTF; font-size: 14pt">Accents: café naïve résumé — 1234567890</p>
<p style="font-family: LicRestricted">Restricted License: system font</p>
<p style="font-family: LicNoSubset">No subsetting: embedded whole</p>
<p style="font-family: LicPreview">Preview &amp; Print: embedded</p>
</body></html>""" % "\n".join(faces)
with open(os.path.join(tmp, "webfonts.html"), "w") as fp:
    fp.write(html)
weasyprint.HTML(os.path.join(tmp, "webfonts.html")).write_pdf(os.path.join(out, "weasyprint-webfonts.pdf"))

s = cairo.PDFSurface(os.path.join(out, "cairo-cff.pdf"), 480, 200)
c = cairo.Context(s)
c.select_font_face("Loma")
c.set_font_size(20)
c.move_to(24, 48)
c.show_text("Type1C via cairo: Hello, world")
c.move_to(24, 96)
c.show_text("CIDFontType0C: ภาษาไทย")
c.set_font_size(14)
c.move_to(24, 140)
c.show_text("café — “quotes” ½")
s.finish()
print("wrote", os.path.join(out, "weasyprint-webfonts.pdf"), os.path.join(out, "cairo-cff.pdf"))
