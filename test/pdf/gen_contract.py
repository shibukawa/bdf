"""Generates converter/pdfform/testdata/contract.pdf: a PDF of one page with
an interactive form (AcroForm) of two parties, each with a text field of
their name and a signature field (/FT /Sig) that converter/pdfform draws a
picture of a signature in. reportlab makes no signature fields, so the file
is written here by hand, without dependencies."""
import os

out = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))), "converter", "pdfform", "testdata")
os.makedirs(out, exist_ok=True)

W, H = 420, 595  # A5
objs = {}


def obj(n, body):
    objs[n] = body


content = b"""BT /F1 18 Tf 40 540 Td (Contract) Tj ET
BT /F1 9 Tf 40 500 Td (The parties agree to the terms of this contract.) Tj ET
BT /F1 9 Tf 40 446 Td (Party A) Tj ET
BT /F1 9 Tf 40 366 Td (Party B) Tj ET
0.5 g 110 440 m 380 440 l S 110 400 m 380 400 l S 110 360 m 380 360 l S 110 320 m 380 320 l S
"""

obj(1, b"<< /Type /Catalog /Pages 2 0 R /AcroForm << /Fields [5 0 R 6 0 R 7 0 R 8 0 R] /DA (/Helv 0 Tf 0 g) /DR << /Font << /Helv 4 0 R >> >> /SigFlags 0 >> >>")
obj(2, b"<< /Type /Pages /Kids [3 0 R] /Count 1 >>")
obj(3, b"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %d %d] /Resources << /Font << /F1 4 0 R >> >> /Contents 9 0 R /Annots [5 0 R 6 0 R 7 0 R 8 0 R] >>" % (W, H))
obj(4, b"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>")
obj(5, b"<< /Type /Annot /Subtype /Widget /FT /Tx /T (party_a.name) /TU (Name of party A) /Rect [110 442 380 462] /F 4 /P 3 0 R /DA (/Helv 11 Tf 0 g) /V () >>")
obj(6, b"<< /Type /Annot /Subtype /Widget /FT /Sig /T (party_a.sign) /TU (Signature of party A) /Rect [110 402 380 436] /F 4 /P 3 0 R >>")
obj(7, b"<< /Type /Annot /Subtype /Widget /FT /Tx /T (party_b.name) /TU (Name of party B) /Rect [110 362 380 382] /F 4 /P 3 0 R /DA (/Helv 11 Tf 0 g) /V () >>")
obj(8, b"<< /Type /Annot /Subtype /Widget /FT /Sig /T (party_b.sign) /TU (Signature of party B) /Rect [110 322 380 356] /F 4 /P 3 0 R >>")
obj(9, b"<< /Length %d >>\nstream\n" % len(content) + content + b"\nendstream")

data = bytearray(b"%PDF-1.7\n%\xe2\xe3\xcf\xd3\n")
offsets = {}
for n in sorted(objs):
    offsets[n] = len(data)
    data += b"%d 0 obj\n" % n + objs[n] + b"\nendobj\n"
xref = len(data)
data += b"xref\n0 %d\n0000000000 65535 f \n" % (len(objs) + 1)
for n in sorted(objs):
    data += b"%010d 00000 n \n" % offsets[n]
data += b"trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n" % (len(objs) + 1, xref)
with open(os.path.join(out, "contract.pdf"), "wb") as f:
    f.write(data)
