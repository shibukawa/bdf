"""Generates converter/pdfform/testdata/application.pdf with reportlab: a
PDF with an interactive form (AcroForm), which converter/pdfform fills as a
template of a form. Two pages: text fields (of a line, of lines, with a
cell for each character, with a size of their own and with a size that
fits), checkboxes, a group of radio buttons and a choice."""
import os
from reportlab.lib.pagesizes import A5
from reportlab.lib.colors import HexColor, black, white
from reportlab.pdfgen import canvas

out = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))), "converter", "pdfform", "testdata")
os.makedirs(out, exist_ok=True)
c = canvas.Canvas(os.path.join(out, "application.pdf"), pagesize=A5, invariant=1)
c.setTitle("Application form")
c.setAuthor("BDF")
w, h = A5
form = c.acroForm
edge = HexColor("#7f7f7f")
tint = HexColor("#f4f7fb")


def label(text, x, y, size=9):
    c.setFont("Helvetica-Bold", size)
    c.drawString(x, y, text)


c.setFont("Helvetica-Bold", 16)
c.drawString(40, h - 50, "Application")
c.setFont("Helvetica", 8)
c.drawRightString(w - 40, h - 50, "Form A-1")
c.setStrokeColor(edge)
c.line(40, h - 58, w - 40, h - 58)

y = h - 90
label("Name", 40, y + 6)
form.textfield(name="applicant.name", tooltip="Name of the applicant", x=110, y=y, width=220, height=20, fontName="Helvetica", fontSize=11,
               borderColor=edge, fillColor=tint, textColor=black, borderWidth=0.5, forceBorder=True, value="")
y -= 30
label("Kana", 40, y + 6)
# (a size of 0: as large as the field has room for)
form.textfield(name="applicant.kana", x=110, y=y, width=220, height=20, fontName="Helvetica", fontSize=0,
               borderColor=edge, fillColor=tint, textColor=black, borderWidth=0.5, forceBorder=True, value="")
y -= 30
label("Zip", 40, y + 6)
form.textfield(name="zip", tooltip="Seven digits", x=110, y=y, width=126, height=20, fontName="Courier", fontSize=12, maxlen=7, fieldFlags="comb",
               borderColor=edge, fillColor=white, textColor=black, borderWidth=0.5, forceBorder=True, value="")
# the cells of the digits
for i in range(1, 7):
    c.line(110 + 18 * i, y, 110 + 18 * i, y + 20)
y -= 62
label("Address", 40, y + 38)
form.textfield(name="address", x=110, y=y, width=260, height=52, fontName="Helvetica", fontSize=10, fieldFlags="multiline",
               borderColor=edge, fillColor=tint, textColor=black, borderWidth=0.5, forceBorder=True, value="")
y -= 30
label("Amount", 40, y + 6)
form.textfield(name="amount", x=110, y=y, width=100, height=20, fontName="Helvetica", fontSize=11,
               borderColor=edge, fillColor=tint, textColor=HexColor("#003399"), borderWidth=0.5, forceBorder=True, value="0")
y -= 34
label("Plan", 40, y + 4)
x = 110
for value, text in (("basic", "Basic"), ("pro", "Pro"), ("team", "Team")):
    form.radio(name="plan", value=value, selected=value == "basic", x=x, y=y, size=14, buttonStyle="circle", shape="circle",
               borderColor=edge, fillColor=white, textColor=black, borderWidth=0.8, forceBorder=True)
    c.setFont("Helvetica", 9)
    c.drawString(x + 18, y + 4, text)
    x += 70
y -= 30
label("Country", 40, y + 6)
form.choice(name="country", x=110, y=y, width=140, height=20, fontName="Helvetica", fontSize=10, value="jp",
            options=[("Japan", "jp"), ("United States", "us"), ("Germany", "de")], fieldFlags="combo",
            borderColor=edge, fillColor=white, textColor=black, borderWidth=0.5, forceBorder=True)
y -= 34
form.checkbox(name="agree", tooltip="Agrees to the terms", checked=False, x=110, y=y, size=14, buttonStyle="check",
              borderColor=edge, fillColor=white, textColor=black, borderWidth=0.8, forceBorder=True)
c.setFont("Helvetica", 9)
c.drawString(130, y + 4, "I agree to the terms")
y -= 24
form.checkbox(name="newsletter", checked=True, x=110, y=y, size=14, buttonStyle="cross",
              borderColor=edge, fillColor=white, textColor=black, borderWidth=0.8, forceBorder=True)
c.drawString(130, y + 4, "Send me the newsletter")
c.setFont("Helvetica", 8)
c.drawCentredString(w / 2, 30, "1 / 2")
c.showPage()

# the second page: a field that the first has too (one value, two places), and notes
c.setFont("Helvetica-Bold", 16)
c.drawString(40, h - 50, "Notes")
label("Name", 40, h - 84)
form.textfield(name="applicant.name", x=110, y=h - 90, width=220, height=20, fontName="Helvetica", fontSize=11,
               borderColor=edge, fillColor=tint, textColor=black, borderWidth=0.5, forceBorder=True, value="")
label("Notes", 40, h - 114)
form.textfield(name="notes", x=40, y=h - 260, width=w - 80, height=130, fontName="Helvetica", fontSize=0, fieldFlags="multiline",
               borderColor=edge, fillColor=white, textColor=black, borderWidth=0.5, forceBorder=True, value="")
c.setFont("Helvetica", 8)
c.drawCentredString(w / 2, 30, "2 / 2")
c.showPage()
c.save()
print("wrote", os.path.join(out, "application.pdf"))
