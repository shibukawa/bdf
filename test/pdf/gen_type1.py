"""Generates test PDFs with Type 1 font programs (FontFile), which pdf2bdf
converts to CFF:

- cairo-type1.pdf: cairo subsets Type 1 fonts; one subset uses WinAnsiEncoding,
  the other (for characters outside it) only its built-in encoding. The text
  has seac accents (é, ï, Å, ö) and a ligature.
- reportlab-type1.pdf: reportlab embeds the whole .pfb.

Needs cairocffi and reportlab, and the Bitstream Charter and Courier 10 Pitch
Type 1 fonts (xfonts-scalable; freely redistributable, see its copyright)."""
import os
import cairocffi as cairo
from reportlab.pdfgen import canvas
from reportlab.pdfbase import pdfmetrics

root = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
out = os.path.join(root, "pdf2bdf", "testdata")
t1 = "/usr/share/fonts/X11/Type1/"

s = cairo.PDFSurface(os.path.join(out, "cairo-type1.pdf"), 480, 200)
c = cairo.Context(s)
c.select_font_face("Bitstream Charter")
c.set_font_size(20)
c.move_to(24, 48)
c.show_text("Type1 via cairo: Hello, world")
c.move_to(24, 90)
c.show_text("café naïve Ångström — fi ﬁ")
c.select_font_face("Courier 10 Pitch")
c.set_font_size(16)
c.move_to(24, 130)
c.show_text("Courier 10 Pitch: mono")
s.finish()

face = pdfmetrics.EmbeddedType1Face(t1 + "c0648bt_.afm", t1 + "c0648bt_.pfb")
pdfmetrics.registerTypeFace(face)
pdfmetrics.registerFont(pdfmetrics.Font("Charter", face.name, "WinAnsiEncoding"))
cv = canvas.Canvas(os.path.join(out, "reportlab-type1.pdf"), pagesize=(480, 120))
cv.setTitle("Type1 fixture")
cv.setFont("Charter", 20)
cv.drawString(24, 72, "Whole Type1 font: café naïve Ærø")
cv.setFont("Charter", 12)
cv.drawString(24, 40, "0123456789 (reportlab embeds every glyph)")
cv.save()
print("wrote", os.path.join(out, "cairo-type1.pdf"), os.path.join(out, "reportlab-type1.pdf"))
