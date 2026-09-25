"""Generates reportlab-master.pdf: pages drawn from one "master" (header band,
logo, footer rule) followed by page-specific content, the shape a PDF exported
from a slide layout or a letterhead has. pdf2bdf's prefix sharing must turn
the common start of every page into one shared object."""
import os
from reportlab.lib.pagesizes import A4
from reportlab.pdfgen import canvas
from reportlab.lib.colors import HexColor
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont

out = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))), "pdf2bdf", "testdata")
pdfmetrics.registerFont(TTFont("DejaVu", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"))

c = canvas.Canvas(os.path.join(out, "reportlab-master.pdf"), pagesize=A4)
c.setTitle("master fixture")
w, h = A4


def master():
    """The part every page starts with."""
    c.setFillColor(HexColor("#1f3a5f"))
    c.rect(0, h - 60, w, 60, stroke=0, fill=1)
    c.setFillColor(HexColor("#4c9be8"))
    c.circle(40, h - 30, 16, stroke=0, fill=1)
    c.setFillColor(HexColor("#ffffff"))
    c.setFont("DejaVu", 16)
    c.drawString(70, h - 36, "Quarterly report — master header")
    # a decorative strip of dots and a small "logo" of paths, as layouts have
    for i in range(24):
        c.setFillColor(HexColor("#4c9be8" if i % 3 else "#2ea06c"))
        c.circle(72 + i * 19, h - 76, 4, stroke=0, fill=1)
    c.setFillColor(HexColor("#2ea06c"))
    p = c.beginPath()
    p.moveTo(w - 110, h - 46); p.lineTo(w - 90, h - 14); p.lineTo(w - 70, h - 46); p.close()
    c.drawPath(p, stroke=0, fill=1)
    c.setStrokeColor(HexColor("#ffffff"))
    c.setLineWidth(2)
    c.line(w - 100, h - 36, w - 80, h - 36)
    c.setStrokeColor(HexColor("#1f3a5f"))
    c.setLineWidth(1.5)
    c.line(72, 60, w - 72, 60)
    c.setFont("Helvetica", 8)
    c.setFillColor(HexColor("#999999"))
    for i, label in enumerate(["Finance", "Operations", "People", "Product", "Markets"]):
        c.drawString(72 + i * 90, 30, label)
    c.setFillColor(HexColor("#666666"))
    c.setFont("Helvetica", 9)
    c.drawString(72, 46, "Confidential — do not distribute")


bodies = [
    ("Overview", ["Revenue grew in every region.", "Costs stayed flat.", "Headcount: 120."]),
    ("Regions", ["EMEA: +12%", "APAC: +8%", "Americas: +5%"]),
    ("Outlook", ["Guidance unchanged.", "Two launches planned.", "Hiring resumes in Q3."]),
]
for i, (title, lines) in enumerate(bodies):
    master()
    c.setFillColor(HexColor("#111111"))
    c.setFont("Helvetica-Bold", 20)
    c.drawString(72, h - 120, title)
    c.setFont("Helvetica", 12)
    y = h - 160
    for line in lines:
        c.drawString(72, y, line)
        y -= 18
    c.setFont("Helvetica", 9)
    c.drawRightString(w - 72, 46, f"page {i + 1} of {len(bodies)}")
    c.showPage()
c.save()
print("wrote", os.path.join(out, "reportlab-master.pdf"))
