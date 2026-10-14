"""Generates the templates of pdfme that converter/pdfme is tested with
(converter/pdfme/testdata) and the data that fills them (test/pdfme/data):

  invoice.json      a blank base page with fields of every kind and a table
                    that grows with its rows over pages
  certificate.json  a base PDF (made here with reportlab) with fields over it

The templates are what the designer of pdfme (https://pdfme.com) saves. Run
with a Python that has pillow and reportlab."""
import base64, io, json, os
from PIL import Image, ImageDraw
from reportlab.lib.colors import HexColor
from reportlab.pdfgen import canvas

root = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
out = os.path.join(root, "converter", "pdfme", "testdata")
data_dir = os.path.join(root, "test", "pdfme", "data")
os.makedirs(out, exist_ok=True)
os.makedirs(data_dir, exist_ok=True)


def png(draw, size):
    img = Image.new("RGBA", size, (0, 0, 0, 0))
    draw(ImageDraw.Draw(img), *size)
    buf = io.BytesIO()
    img.save(buf, "PNG", optimize=True)
    return "data:image/png;base64," + base64.b64encode(buf.getvalue()).decode()


def logo(d, w, h):
    d.rounded_rectangle([0, 0, w - 1, h - 1], 12, fill=(51, 102, 204, 255))
    d.polygon([(20, h - 14), (w // 2, 14), (w - 20, h - 14)], fill=(255, 255, 255, 255))


def sign(d, w, h):
    pts = [(10, h - 20), (40, 20), (70, h - 25), (100, 25), (130, h - 30), (170, 30), (w - 12, h // 2)]
    d.line(pts, fill=(20, 30, 120, 255), width=5, joint="curve")


def save(name, obj, where):
    with open(os.path.join(where, name), "w", encoding="utf-8") as f:
        json.dump(obj, f, ensure_ascii=False, indent=1)
        f.write("\n")


box = {"top": 0.3, "right": 0.3, "bottom": 0.3, "left": 0.3}
pad = {"top": 1.5, "right": 1.5, "bottom": 1.5, "left": 1.5}
cell = {"fontSize": 10, "lineHeight": 1, "characterSpacing": 0, "alignment": "left", "verticalAlignment": "middle",
        "borderWidth": {"top": 0.1, "right": 0.1, "bottom": 0.1, "left": 0.1}, "padding": {"top": 2, "right": 2, "bottom": 2, "left": 2}}

invoice = {
    "pdfmeVersion": "6.2.4",
    "basePdf": {
        "width": 210, "height": 297, "padding": [20, 10, 20, 10],
        "staticSchema": [
            {"name": "brand", "type": "text", "readOnly": True, "content": "ACME Trading", "position": {"x": 10, "y": 8}, "width": 100, "height": 7,
             "fontSize": 11, "fontColor": "#3366cc"},
            {"name": "headRule", "type": "line", "readOnly": True, "position": {"x": 10, "y": 16}, "width": 190, "height": 0.4, "color": "#3366cc"},
            {"name": "footRule", "type": "line", "readOnly": True, "position": {"x": 10, "y": 281}, "width": 190, "height": 0.2, "color": "#888888"},
            {"name": "pageNo", "type": "text", "readOnly": True, "content": "{currentPage} / {totalPages}", "position": {"x": 10, "y": 283}, "width": 190,
             "height": 6, "fontSize": 9, "alignment": "center", "fontColor": "#555555"},
        ],
    },
    "schemas": [[
        {"name": "title", "type": "text", "readOnly": True, "content": "インボイス", "position": {"x": 10, "y": 20}, "width": 90, "height": 12, "fontSize": 22},
        {"name": "number", "type": "text", "content": "INV-0001", "position": {"x": 140, "y": 20}, "width": 60, "height": 7, "fontSize": 12,
         "alignment": "right", "required": True},
        {"name": "issued", "type": "date", "content": "2026/10/14", "format": "yyyy.MM.dd", "locale": "ja", "position": {"x": 140, "y": 28},
         "width": 60, "height": 7, "fontSize": 11, "alignment": "right"},
        {"name": "customer", "type": "multiVariableText", "text": "{company}\n{person} さま", "variables": ["company", "person"],
         "content": json.dumps({"company": "カメラのタナカ", "person": "タナカ ハナコ"}, ensure_ascii=False),
         "position": {"x": 10, "y": 38}, "width": 110, "height": 17, "fontSize": 12, "lineHeight": 1.3, "borderColor": "#888888",
         "borderWidth": {"top": 0, "right": 0, "bottom": 0.3, "left": 0}, "padding": {"top": 1, "right": 1, "bottom": 1, "left": 1}},
        {"name": "logo", "type": "image", "readOnly": True, "content": png(logo, (160, 80)), "position": {"x": 170, "y": 37}, "width": 30, "height": 16},
        {"name": "note", "type": "text", "content": "いつもありがとうございます。", "position": {"x": 10, "y": 58}, "width": 110, "height": 14,
         "fontSize": 11, "dynamicFontSize": {"min": 6, "max": 14, "fit": "vertical"}, "verticalAlignment": "middle", "backgroundColor": "#f2f6ff",
         "borderColor": "#3366cc", "borderWidth": box, "padding": pad},
        {"name": "paid", "type": "checkbox", "content": "false", "position": {"x": 130, "y": 58}, "width": 6, "height": 6, "color": "#000000"},
        {"name": "paidLabel", "type": "text", "readOnly": True, "content": "しはらいずみ", "position": {"x": 137, "y": 58}, "width": 30, "height": 6,
         "fontSize": 10, "verticalAlignment": "middle"},
        {"name": "method", "type": "select", "content": "カード", "options": ["カード", "ふりこみ"], "position": {"x": 170, "y": 58}, "width": 30,
         "height": 6, "fontSize": 10, "verticalAlignment": "middle", "alignment": "right"},
        {"name": "express", "type": "radioGroup", "group": "ship", "content": "false", "position": {"x": 130, "y": 66}, "width": 6, "height": 6,
         "color": "#3366cc"},
        {"name": "expressLabel", "type": "text", "readOnly": True, "content": "Express", "position": {"x": 137, "y": 66}, "width": 30, "height": 6,
         "fontSize": 10, "verticalAlignment": "middle"},
        {"name": "items", "type": "table", "position": {"x": 10, "y": 78}, "width": 190, "height": 30,
         "content": json.dumps([["ノート", "2", "300", "600"], ["ペン", "10", "120", "1,200"]], ensure_ascii=False),
         "showHead": True, "repeatHead": True, "head": ["項目", "数量", "単価", "金額"], "headWidthPercentages": [46, 14, 20, 20],
         "tableStyles": {"borderWidth": 0.3, "borderColor": "#000000"},
         "headStyles": dict(cell, fontColor="#ffffff", backgroundColor="#2980ba", borderColor="#2980ba"),
         "bodyStyles": dict(cell, fontColor="#000000", backgroundColor="", alternateBackgroundColor="#f5f5f5", borderColor="#888888"),
         "columnStyles": {"alignment": {"1": "right", "2": "right", "3": "right"}}},
        {"name": "totalBox", "type": "rectangle", "readOnly": True, "position": {"x": 128, "y": 111}, "width": 72, "height": 12, "borderWidth": 0.5,
         "borderColor": "#3366cc", "color": "#f2f6ff", "radius": 2},
        {"name": "totalText", "type": "text", "readOnly": True, "content": "総額 {total}", "position": {"x": 130, "y": 113}, "width": 68, "height": 8,
         "fontSize": 14, "alignment": "right", "verticalAlignment": "middle", "underline": True},
        {"name": "qr", "type": "qrcode", "content": "https://example.com/", "position": {"x": 10, "y": 111}, "width": 22, "height": 22,
         "barColor": "#000000", "backgroundColor": "#ffffff"},
        {"name": "code", "type": "code128", "content": "INV0001", "position": {"x": 38, "y": 113}, "width": 60, "height": 14, "barColor": "#000000",
         "backgroundColor": "", "includetext": False},
        {"name": "terms", "type": "list", "content": json.dumps(["30 日まで", "\tカード", "\tふりこみ", "返すのは 7 日まで"], ensure_ascii=False),
         "listStyle": "ordered", "position": {"x": 10, "y": 138}, "width": 110, "height": 22, "fontSize": 10, "lineHeight": 1.2},
        {"name": "seal", "type": "ellipse", "readOnly": True, "position": {"x": 150, "y": 136}, "width": 26, "height": 26, "borderWidth": 0.8,
         "borderColor": "#cc0000", "color": ""},
        {"name": "sealText", "type": "text", "readOnly": True, "content": "COPY", "position": {"x": 138, "y": 143}, "width": 50, "height": 12,
         "fontSize": 20, "fontColor": "#cc0000", "alignment": "center", "verticalAlignment": "middle", "rotate": 20, "opacity": 0.5,
         "characterSpacing": 1},
        {"name": "confirmed", "type": "circleMark", "content": "false", "position": {"x": 178, "y": 141}, "width": 22, "height": 16, "borderWidth": 0.6,
         "color": "#3366cc", "rotate": 0, "opacity": 1},
    ]],
}
save("invoice.json", invoice, out)

things = ["ノート", "ペン", "カメラ", "レンズ", "フィルム", "アルバム", "ケース", "ストラップ"]
rows, total = [], 0
for i in range(34):
    n, price = i % 5 + 1, 120 * (i % 7 + 1)
    name = things[i % len(things)] + " " + chr(ord("A") + i % 26)
    if i == 3:
        name += " — ながい項目は行をおりかえして表のたかさをかえます"
    rows.append({"項目": name, "数量": n, "単価": f"{price:,}", "金額": f"{n * price:,}"})
    total += n * price
save("invoice.json", {
    "number": "INV-2026-1014", "issued": "2026/10/14", "customer": {"company": "カメラのタナカ", "person": "タナカ ハナコ"},
    "note": "いつもありがとうございます。このたびのごちゅうもんを表にまとめました。", "paid": True, "method": "ふりこみ", "express": True,
    "items": rows, "total": f"{total:,}", "qr": "https://example.com/inv/2026-1014", "code": "INV20261014",
    "terms": ["30 日まで", "\tカード", "\tふりこみ", "返すのは 7 日まで"], "confirmed": True,
}, data_dir)

# the base PDF of the certificate: A5 landscape, with frames and a title
buf = io.BytesIO()
w, h = 595.2756, 419.5276
c = canvas.Canvas(buf, pagesize=(w, h), invariant=1)
c.setTitle("Certificate")
c.setFillColor(HexColor("#fbf8ef"))
c.rect(0, 0, w, h, stroke=0, fill=1)
c.setStrokeColor(HexColor("#b08d57"))
c.setLineWidth(3)
c.rect(18, 18, w - 36, h - 36)
c.setLineWidth(0.8)
c.rect(26, 26, w - 52, h - 52)
c.setFillColor(HexColor("#3a2f1b"))
c.setFont("Times-Bold", 30)
c.drawCentredString(w / 2, h - 86, "Certificate of Completion")
c.setFont("Times-Italic", 12)
c.drawCentredString(w / 2, h - 116, "This is to certify that")
c.drawCentredString(w / 2, h - 190, "has completed the course")
c.setFont("Helvetica", 8)
c.drawString(60, 62, "Date")
c.drawString(w - 200, 62, "Signature")
c.setLineWidth(0.5)
c.line(60, 74, 200, 74)
c.line(w - 200, 74, w - 60, 74)
c.showPage()
c.save()

certificate = {
    "pdfmeVersion": "6.2.4",
    "basePdf": "data:application/pdf;base64," + base64.b64encode(buf.getvalue()).decode(),
    "schemas": [[
        {"name": "name", "type": "text", "content": "Jane Doe", "position": {"x": 30, "y": 46}, "width": 150, "height": 16, "fontSize": 28,
         "dynamicFontSize": {"min": 10, "max": 34, "fit": "horizontal"}, "alignment": "center", "verticalAlignment": "middle", "fontColor": "#1c2a5a"},
        {"name": "course", "type": "text", "content": "Document Formats 101", "position": {"x": 30, "y": 86}, "width": 150, "height": 10, "fontSize": 16,
         "alignment": "center", "verticalAlignment": "middle"},
        {"name": "score", "type": "text", "readOnly": True, "content": "Score {points} / 100 ({points >= 80 ? 'with distinction' : 'passed'})",
         "position": {"x": 30, "y": 98}, "width": 150, "height": 7, "fontSize": 10, "alignment": "center", "fontColor": "#555555"},
        {"name": "date", "type": "date", "content": "2026/10/14", "format": "MMMM d, yyyy", "locale": "en", "position": {"x": 21, "y": 113},
         "width": 50, "height": 8, "fontSize": 11, "alignment": "center"},
        {"name": "signature", "type": "image", "content": "", "position": {"x": 140, "y": 106}, "width": 48, "height": 15},
        {"name": "verify", "type": "qrcode", "content": "https://example.com/", "position": {"x": 94, "y": 106}, "width": 20, "height": 20,
         "barColor": "#3a2f1b", "backgroundColor": ""},
        {"name": "honours", "type": "checkbox", "content": "false", "position": {"x": 30, "y": 127}, "width": 5, "height": 5, "color": "#3a2f1b"},
        {"name": "honoursLabel", "type": "text", "readOnly": True, "content": "With honours", "position": {"x": 36, "y": 127}, "width": 40, "height": 5,
         "fontSize": 9, "verticalAlignment": "middle"},
    ]],
}
save("certificate.json", certificate, out)
save("certificate.json", {"name": "Hanako Tanaka", "course": "Document Formats 101", "points": 92, "date": "2026/10/14",
                          "signature": png(sign, (240, 80)), "verify": "https://example.com/cert/42", "honours": True}, data_dir)
