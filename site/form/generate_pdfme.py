"""Build standard and modern pdfme invoices, independently of converter fixtures.

The template uses only pdfme primitives understood by BDF. Its growing table
pushes the totals, bank details, and notes onto later pages when necessary.
Run with Python 3; no third-party packages are required.
"""
from pathlib import Path
from copy import deepcopy
import json

OUT = Path(__file__).resolve().parent / "samples"
FONT = "M PLUS 1p"
NAVY, INK, MUTED, RULE, PALE = "#243B53", "#263342", "#667685", "#CBD5DF", "#F3F6F9"


def text(name, content, x, y, width, height, size=10, **extra):
    return {
        "name": name, "type": "text", "readOnly": True, "content": content,
        "position": {"x": x, "y": y}, "width": width, "height": height,
        "fontName": FONT, "fontSize": size, "fontColor": INK,
        "lineHeight": 1.35, **extra,
    }


def variable(name, pattern, variables, sample, x, y, width, height, **extra):
    return {
        **text(name, "", x, y, width, height, **extra),
        "type": "multiVariableText", "readOnly": False, "text": pattern,
        "variables": variables, "content": json.dumps(sample, ensure_ascii=False),
    }


def line(name, y, width=180, x=15, color=RULE, height=0.25):
    return {"name": name, "type": "line", "readOnly": True,
            "position": {"x": x, "y": y}, "width": width, "height": height, "color": color}


data = {
    "number": "INV-2026-1014", "issued": "2026/10/14", "due": "2026/11/13",
    "customer": {"company": "株式会社みなと商事", "person": "営業部 御中"},
    "items": [
        {"項目": "ノート A5", "数量": 10, "単価": "480", "金額": "4,800"},
        {"項目": "ボールペン", "数量": 20, "単価": "150", "金額": "3,000"},
        {"項目": "デスクトレー", "数量": 3, "単価": "2,800", "金額": "8,400"},
        {"項目": "ファイルボックス", "数量": 6, "単価": "650", "金額": "3,900"},
        {"項目": "付箋セット", "数量": 4, "単価": "420", "金額": "1,680"},
    ],
    "totals": {"subtotal": "21,780", "tax": "2,178", "total": "23,958"},
    "note": "お振込手数料は貴社にてご負担くださいますようお願いいたします。",
}

cell = {
    "fontName": FONT, "fontSize": 10, "lineHeight": 1.15,
    "characterSpacing": 0, "alignment": "left", "verticalAlignment": "middle",
    "padding": {"top": 2.4, "right": 3, "bottom": 2.4, "left": 3},
    "borderWidth": {"top": 0, "right": 0, "bottom": 0.2, "left": 0},
}

modern = {
    "pdfmeVersion": "6.2.4",
    "basePdf": {
        "width": 210, "height": 297, "padding": [37, 15, 24, 15],
        "staticSchema": [
            text("brand", "BDF STUDIO", 15, 12, 100, 7, 14, fontColor=NAVY, characterSpacing=0.6),
            text("brandDetail", "Document Design & Supply", 15, 20, 100, 5, 8, fontColor=MUTED),
            text("issuer", "〒150-0001 東京都渋谷区神宮前1-2-3\nTEL 03-5550-0140 / hello@bdf.example", 115, 14, 80, 12,
                 8.5, fontColor=MUTED, alignment="right", lineHeight=1.5),
            line("headerRule", 30, color=NAVY, height=0.4),
            line("footerRule", 279),
            text("footerBrand", "BDF STUDIO  /  INVOICE", 15, 282, 130, 5, 8, fontColor=MUTED, characterSpacing=0.3),
            text("pageNumber", "{currentPage} / {totalPages}", 165, 282, 30, 5, 8, fontColor=MUTED, alignment="right"),
        ],
    },
    "schemas": [[
        text("title", "請求書", 15, 38, 100, 14, 24, fontColor=NAVY, characterSpacing=1.2),
        text("numberLabel", "請求番号", 128, 39, 24, 5, 8.5, fontColor=MUTED),
        text("number", data["number"], 150, 38.5, 45, 6, 10, readOnly=False, alignment="right", required=True),
        text("issuedLabel", "Date", 128, 47, 24, 5, 8.5, fontColor=MUTED),
        {**text("issued", data["issued"], 150, 46.5, 45, 6, 10, readOnly=False, alignment="right"),
         "type": "date", "format": "yyyy/MM/dd", "locale": "ja"},
        text("dueLabel", "支払期限", 128, 55, 24, 5, 8.5, fontColor=MUTED),
        {**text("due", data["due"], 150, 54.5, 45, 6, 10, readOnly=False, alignment="right"),
         "type": "date", "format": "yyyy/MM/dd", "locale": "ja"},
        variable("customer", "{company}\n{person}", ["company", "person"], data["customer"],
                 15, 65, 110, 18, size=12, lineHeight=1.65, fontColor=NAVY, overflow="expand"),
        text("greeting", "下記のとおりご請求申し上げます。", 15, 87, 180, 6, 9.5),
        text("amountLabel", "ご請求金額（税込）", 15, 99, 105, 5, 8.5, fontColor=MUTED),
        text("amount", "¥ {totals.total}", 15, 106, 110, 13, 22, fontColor=NAVY),
        text("currency", "金額はすべて日本円です。", 130, 114, 65, 5, 8, fontColor=MUTED, alignment="right"),
        {
            "name": "items", "type": "table", "position": {"x": 15, "y": 127}, "width": 180, "height": 55,
            "content": json.dumps([[str(row[key]) for key in ["項目", "数量", "単価", "金額"]] for row in data["items"]], ensure_ascii=False),
            "showHead": True, "repeatHead": True, "head": ["項目", "数量", "単価", "金額"],
            "headWidthPercentages": [48, 12, 20, 20],
            "tableStyles": {"borderWidth": 0, "borderColor": RULE},
            "headStyles": {**cell, "fontColor": NAVY, "backgroundColor": PALE, "borderColor": NAVY,
                           "borderWidth": {"top": 0, "right": 0, "bottom": 0.35, "left": 0}},
            "bodyStyles": {**cell, "fontColor": INK, "backgroundColor": PALE, "alternateBackgroundColor": "#FFFFFF", "borderColor": RULE},
            "columnStyles": {"alignment": {"1": "right", "2": "right", "3": "right"}},
        },
        variable("totals", "小計　{subtotal} 円\n消費税（10%）　{tax} 円\n総額　{total} 円", ["subtotal", "tax", "total"], data["totals"],
                 120, 190, 75, 30, size=11, lineHeight=1.8, alignment="right", fontColor=NAVY,
                 backgroundColor=PALE, padding={"top": 4, "right": 4, "bottom": 4, "left": 4}),
        text("payment", "お振込先\nみらい銀行　渋谷支店　普通 1234567\n口座名義：ビーディーエフスタジオ", 15, 232, 180, 21,
             9.5, lineHeight=1.65, fontColor=INK),
        text("remarks", "備考\n{note}", 15, 259, 180, 12, 9, lineHeight=1.6, fontColor=MUTED, overflow="expand"),
    ]],
}

def standard_invoice(source):
    """A monochrome Japanese business form with a centered title and full grid.

    Field names, defaults, and data bindings stay identical to the modern
    invoice, so readers can compare the document styles with the same data.
    """
    colors = {NAVY: "#111111", INK: "#111111", MUTED: "#222222", RULE: "#333333", PALE: "#FFFFFF"}

    def monochrome(value):
        if isinstance(value, dict):
            return {key: monochrome(item) for key, item in value.items()}
        if isinstance(value, list):
            return [monochrome(item) for item in value]
        return colors.get(value, value) if isinstance(value, str) else value

    invoice = monochrome(deepcopy(source))
    fields = {field["name"]: field for field in invoice["schemas"][0]}
    fields["title"].update(width=180, alignment="center", fontName=FONT + "-Bold", characterSpacing=6)
    # Keep the document identity below the centered title and beside the
    # addressee, rather than competing with the title on the same line.
    for name in ["numberLabel", "number", "issuedLabel", "issued", "dueLabel", "due"]:
        fields[name]["position"]["y"] += 20
    fields["amountLabel"].update(fontName=FONT + "-Bold", fontSize=9)
    fields["amount"].update(
        fontName=FONT + "-Bold", borderColor="#111111",
        borderWidth={"top": 0, "right": 0, "bottom": 0.5, "left": 0},
    )
    table = fields["items"]
    table["tableStyles"].update(borderWidth=0.35, borderColor="#111111")
    for part in ["headStyles", "bodyStyles"]:
        table[part].update(
            backgroundColor="#FFFFFF", alternateBackgroundColor="#FFFFFF", borderColor="#333333",
            borderWidth={"top": 0.15, "right": 0.15, "bottom": 0.15, "left": 0.15},
        )
    table["headStyles"]["fontName"] = FONT + "-Bold"
    fields["totals"].update(fontName=FONT + "-Bold", borderColor="#111111", borderWidth=0.4)
    for field in invoice["basePdf"]["staticSchema"]:
        if field["name"] == "brand":
            field["fontName"] = FONT + "-Bold"
    return invoice


OUT.mkdir(parents=True, exist_ok=True)
themes = {"standard": standard_invoice(modern), "modern": modern}
for theme, invoice in themes.items():
    suffix = "" if theme == "standard" else "-" + theme
    (OUT / ("pdfme-invoice" + suffix + ".json")).write_text(
        json.dumps(invoice, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
(OUT / "pdfme-invoice.data.json").write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
