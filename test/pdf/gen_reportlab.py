"""Generates test PDFs with reportlab: standard 14 fonts, TrueType subsets, shapes, images, links."""
import io, os, sys
from reportlab.lib.pagesizes import A4
from reportlab.pdfgen import canvas
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.lib.colors import Color, HexColor
from reportlab.lib.utils import ImageReader

out = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))), "pdf2bdf", "testdata")
pdfmetrics.registerFont(TTFont("DejaVu", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"))
pdfmetrics.registerFont(TTFont("DejaVu-Bold", "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"))

c = canvas.Canvas(os.path.join(out, "reportlab-mixed.pdf"), pagesize=A4)
c.setTitle("reportlab fixture")
w, h = A4
# standard fonts (not embedded)
y = h - 72
for font in ["Helvetica", "Helvetica-Bold", "Times-Roman", "Times-Italic", "Courier"]:
    c.setFont(font, 14)
    c.drawString(72, y, f"{font}: The quick brown fox jumps over the lazy dog")
    y -= 22
# TrueType subset (simple font, symbolic cmap) with non-ASCII
c.setFont("DejaVu", 14)
c.drawString(72, y, "DejaVu TrueType: café naïve — “quotes” • bullet ½ ©")
y -= 22
c.setFont("DejaVu-Bold", 14)
c.drawString(72, y, "DejaVu Bold, right-aligned below:")
c.drawRightString(w - 72, y - 22, "right aligned text")
y -= 60
# shapes
c.setFillColor(HexColor("#4c9be8")); c.setStrokeColor(HexColor("#1f3a5f")); c.setLineWidth(3); c.setDash(6, 3)
c.roundRect(72, y - 120, 200, 120, 12, stroke=1, fill=1)
c.setDash()
c.setFillColor(Color(0.9, 0.3, 0.2, alpha=0.5))
c.circle(360, y - 60, 60, stroke=0, fill=1)
c.setFillColor(HexColor("#2ea06c"))
p = c.beginPath(); p.moveTo(430, y - 120); p.curveTo(470, y, 520, y - 140, 560, y - 20); p.lineTo(560, y - 120); p.close()
c.drawPath(p, stroke=1, fill=1)
y -= 150
# image (PNG with alpha and JPEG)
from PIL import Image
img = Image.new("RGBA", (64, 64))
for yy in range(64):
    for xx in range(64):
        d = 1 if ((xx >> 3) + (yy >> 3)) % 2 else 0.5
        img.putpixel((xx, yy), (int(255 * xx / 64 * d), int(255 * yy / 64 * d), 160 if d == 1 else 90, 255 if xx < 48 else 96))
buf = io.BytesIO(); img.save(buf, "PNG"); buf.seek(0)
c.drawImage(ImageReader(buf), 72, y - 150, 150, 150, mask="auto")
jbuf = io.BytesIO(); img.convert("RGB").save(jbuf, "JPEG", quality=90); jbuf.seek(0)
c.drawImage(ImageReader(jbuf), 250, y - 150, 150, 150)
# rotated text and link
c.saveState(); c.translate(470, y - 150); c.rotate(30); c.setFont("Helvetica", 12); c.setFillColor(HexColor("#1f3a5f"))
c.drawString(0, 0, "rotated 30 degrees"); c.restoreState()
c.linkURL("https://example.com/rl", (72, y - 190, 300, y - 170), relative=0)
c.setFont("Helvetica", 12); c.setFillColor(HexColor("#0645ad")); c.drawString(72, y - 185, "link annotation over this text")
c.showPage()
# page 2: text state variants
c.setFont("DejaVu", 12)
t = c.beginText(72, h - 72)
t.setCharSpace(2); t.textLine("char spacing 2")
t.setCharSpace(0); t.setWordSpace(6); t.textLine("word spacing six between words")
t.setWordSpace(0); t.setHorizScale(70); t.textLine("horizontal scale 70 percent")
t.setHorizScale(100); t.setRise(6); t.textOut("rise "); t.setRise(0); t.textLine("normal")
t.setTextRenderMode(1); t.textLine("stroke only outline text")
t.setTextRenderMode(3); t.textLine("invisible text (searchable)")
t.setTextRenderMode(0)
c.drawText(t)
c.save()
print("wrote", os.path.join(out, "reportlab-mixed.pdf"))
