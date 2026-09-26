"""Generates the Excel test workbooks in converter/xlsx/testdata with openpyxl.

Usage: python3 test/xlsx/gen.py   (pip install openpyxl pillow)
basic.xlsx exercises cell formats: number formats, fonts and rich text,
fills, borders, alignment (wrapping Latin and Japanese text, overflow,
rotation, indent, shrink to fit), merged cells, row and column sizes,
frozen panes, notes and hyperlinks, a hidden sheet, and a sheet large
enough to span several tiles. features.xlsx exercises conditional formats,
tables, a chart, a picture, a text box shape and a chart sheet.
"""
import datetime
import io
import os
import re
import struct
import zipfile
import zlib

from openpyxl import Workbook
from openpyxl.cell.rich_text import CellRichText, TextBlock
from openpyxl.cell.text import InlineFont
from openpyxl.chart import BarChart, LineChart, Reference
from openpyxl.chartsheet import Chartsheet
from openpyxl.comments import Comment
from openpyxl.drawing.image import Image
from openpyxl.formatting.rule import (CellIsRule, ColorScaleRule, DataBarRule, FormulaRule, IconSetRule, Rule)
from openpyxl.styles import Alignment, Border, Font, GradientFill, PatternFill, Side
from openpyxl.styles.differential import DifferentialStyle
from openpyxl.worksheet.table import Table, TableStyleInfo

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "xlsx", "testdata")
FIXED = datetime.datetime(2026, 9, 1, 9, 0, 0)


def png(w, h, fn):
    """A small RGB PNG from fn(x, y) -> (r, g, b)."""
    raw = b""
    for y in range(h):
        raw += b"\x00" + b"".join(bytes(fn(x, y)) for x in range(w))

    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xFFFFFFFF)
    return (b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
            + chunk(b"IDAT", zlib.compress(raw, 9)) + chunk(b"IEND", b""))


def props(wb, title):
    wb.properties.title = title
    wb.properties.creator = "bdf test"
    wb.properties.lastModifiedBy = "bdf test"
    wb.properties.created = FIXED
    wb.properties.modified = FIXED


thin = Side(style="thin", color="000000")
medium = Side(style="medium", color="1F4E79")


def basic():
    wb = Workbook()
    props(wb, "Excel test workbook")
    ws = wb.active
    ws.title = "Sales"
    ws.column_dimensions["A"].width = 14
    ws.column_dimensions["B"].width = 12
    ws.column_dimensions["C"].width = 12
    ws.column_dimensions["D"].width = 13
    ws.column_dimensions["E"].width = 18
    ws.column_dimensions["F"].width = 8
    ws.column_dimensions["H"].hidden = True

    # title: merged, bold, filled
    ws.merge_cells("A1:F1")
    ws["A1"] = "Quarterly sales 数字の表"
    ws["A1"].font = Font(size=16, bold=True, color="FFFFFF")
    ws["A1"].fill = PatternFill("solid", fgColor="1F4E79")
    ws["A1"].alignment = Alignment(horizontal="center", vertical="center")
    ws.row_dimensions[1].height = 28

    head = ["Region", "Q1", "Q2", "Growth", "Updated", "Flag"]
    for i, h in enumerate(head):
        c = ws.cell(row=3, column=i + 1, value=h)
        c.font = Font(bold=True)
        c.fill = PatternFill("solid", fgColor="DDEBF7")
        c.border = Border(top=thin, bottom=medium, left=thin, right=thin)
        c.alignment = Alignment(horizontal="center")
    rows = [
        ("North ノース", 125000, 131250.5, 0.05, datetime.date(2026, 3, 31), True),
        ("South サウス", 98000, 91140, -0.07, datetime.date(2026, 6, 30), False),
        ("East", 1234567.891, 1300000, 0.053, datetime.date(2026, 9, 30), True),
        ("West", 0, 0.25, 1.5, datetime.datetime(2026, 12, 31, 18, 30), False),
    ]
    for r, row in enumerate(rows, start=4):
        for i, v in enumerate(row):
            c = ws.cell(row=r, column=i + 1, value=v)
            c.border = Border(left=thin, right=thin, bottom=Side(style="hair", color="808080"))
        ws.cell(row=r, column=2).number_format = "#,##0"
        ws.cell(row=r, column=3).number_format = '#,##0.00;[Red]-#,##0.00'
        ws.cell(row=r, column=4).number_format = "0.0%"
        ws.cell(row=r, column=5).number_format = "yyyy-mm-dd"
    ws["E7"].number_format = "yyyy/m/d h:mm"
    ws["C5"] = -1234.5
    ws["A8"] = "Total"
    ws["A8"].font = Font(bold=True)
    ws["B8"] = "=SUM(B4:B7)"
    ws["B8"].number_format = "#,##0"
    ws["B8"].border = Border(top=thin, bottom=Side(style="double"))
    ws.freeze_panes = "A4"

    # text: overflow, clipping, wrapping (Latin and Japanese), alignment
    ws["A10"] = "This text is long and flows over the empty cells to its right."
    ws["A11"] = "Clipped by the next cell"
    ws["B11"] = "next"
    ws["F12"] = "Right-aligned text flows left"
    ws["F12"].alignment = Alignment(horizontal="right")
    ws["A13"] = "日本語の文章は、単語の区切りに空白を使わないため、文字と文字の間で改行します。句読点「、」や「。」は行頭に来ません。"
    ws["A13"].alignment = Alignment(wrap_text=True, vertical="top")
    ws.merge_cells("A13:C13")
    ws.row_dimensions[13].height = 62
    ws["D13"] = "Wrapped Latin text breaks after spaces between words."
    ws["D13"].alignment = Alignment(wrap_text=True, vertical="center")
    ws["E13"] = "Justified text spreads the words of every line but the last to both edges."
    ws["E13"].alignment = Alignment(horizontal="justify", vertical="top")
    ws["A15"] = "indent 2"
    ws["A15"].alignment = Alignment(indent=2)
    ws["B15"] = "shrink to fit this text"
    ws["B15"].alignment = Alignment(shrink_to_fit=True)
    ws["C15"] = "centered across"
    ws["C15"].alignment = Alignment(horizontal="centerContinuous")
    ws["D15"].alignment = Alignment(horizontal="centerContinuous")
    ws["F15"] = "-~"
    ws["F15"].alignment = Alignment(horizontal="fill")
    ws["A16"] = "rotated 45°"
    ws["A16"].alignment = Alignment(text_rotation=45)
    ws["B16"] = "up 90°"
    ws["B16"].alignment = Alignment(text_rotation=90, horizontal="center")
    ws["C16"] = "down"
    ws["C16"].alignment = Alignment(text_rotation=180)
    ws["D16"] = "縦書き"
    ws["D16"].alignment = Alignment(text_rotation=255, horizontal="center")
    ws.row_dimensions[16].height = 60
    # numbers that do not fit, General in a narrow column
    ws["F17"] = 123456789012
    ws["F17"].number_format = "#,##0"
    ws["F18"] = 3.14159265358979
    ws["F19"] = "=1/3"
    # rich text, fonts
    ws["A18"] = CellRichText("Rich ", TextBlock(InlineFont(b=True, color="C00000"), "bold red"), " and ",
                             TextBlock(InlineFont(i=True, u="single", rFont="Times New Roman"), "italic"))
    ws["A19"] = "Strike"
    ws["A19"].font = Font(strike=True, color="7F7F7F")
    ws["B19"] = "H"
    ws["C19"] = "x2"
    ws["C19"].font = Font(vertAlign="superscript")
    ws["A20"] = "Big text"
    ws["A20"].font = Font(size=20, name="Georgia")
    # fills and borders
    ws["A22"] = "gray125"
    ws["A22"].fill = PatternFill("gray125", fgColor="4472C4", bgColor="FFFFFF")
    ws["B22"] = "darkGrid"
    ws["B22"].fill = PatternFill("darkGrid", fgColor="A9D08E", bgColor="FFFFFF")
    ws["C22"] = "gradient"
    ws["C22"].fill = GradientFill(stop=("FFFFFF", "F4B084"), degree=90)
    ws["D22"] = "theme"
    ws["D22"].fill = PatternFill("solid", fgColor="FFE699")
    ws["E22"] = "diagonal"
    ws["E22"].border = Border(diagonal=Side(style="thin", color="FF0000"), diagonalDown=True, diagonalUp=True)
    for col, style in zip("ABCDEF", ["thin", "medium", "thick", "double", "dashed", "dotted"]):
        c = ws[f"{col}24"]
        c.value = style
        c.border = Border(top=Side(style=style), bottom=Side(style=style), left=Side(style=style), right=Side(style=style))
    ws.row_dimensions[21].hidden = True
    ws["A21"] = "hidden row"
    # a note and a hyperlink
    ws["B4"].comment = Comment("Highest quarter", "bdf")
    ws["A26"] = "example.com"
    ws["A26"].hyperlink = "https://example.com/"
    ws["A26"].font = Font(color="0563C1", underline="single")
    ws["A27"] = "Date"
    ws["B27"] = datetime.date(2026, 5, 1)
    ws["B27"].number_format = 'd-mmm-yy'
    ws["C27"] = datetime.date(2026, 5, 1)
    ws["C27"].number_format = "dddd"
    ws["A28"] = 0.375
    ws["A28"].number_format = "# ?/8"
    ws["B28"] = 12345.678
    ws["B28"].number_format = "0.00E+00"
    ws["C28"] = 1234.5
    ws["C28"].number_format = '_("$"* #,##0.00_);_("$"* \\(#,##0.00\\);_("$"* "-"??_);_(@_)'
    ws["D28"] = -1234.5
    ws["D28"].number_format = '_("$"* #,##0.00_);_("$"* \\(#,##0.00\\);_("$"* "-"??_);_(@_)'
    ws["E28"] = "#N/A"
    ws["E28"].data_type = "e"

    hidden = wb.create_sheet("Hidden")
    hidden["A1"] = "Hidden sheet"
    hidden.sheet_state = "hidden"

    big = wb.create_sheet("Big")
    big["A1"] = "Row"
    big["B1"] = "Value"
    for r in range(2, 302):
        big.cell(row=r, column=1, value=f"row {r}")
        big.cell(row=r, column=2, value=r * r)
    # text across a vertical tile boundary (2048 pt = row 137 at 15 pt)
    big["C137"] = "straddles the tile boundary"
    big["C137"].font = Font(size=28)
    big.row_dimensions[137].height = 40
    # text flowing over a horizontal tile boundary (at 2048 pt: column AI with
    # the 60 pt columns of the test fonts)
    big["AG5"] = "flows across the tile boundary to the right"
    wb.save(os.path.join(OUT, "basic.xlsx"))
    fix_core(os.path.join(OUT, "basic.xlsx"))


SHAPE = """<xdr:twoCellAnchor editAs="oneCell"><xdr:from><xdr:col>8</xdr:col><xdr:colOff>0</xdr:colOff><xdr:row>31</xdr:row><xdr:rowOff>0</xdr:rowOff></xdr:from><xdr:to><xdr:col>11</xdr:col><xdr:colOff>304800</xdr:colOff><xdr:row>36</xdr:row><xdr:rowOff>0</xdr:rowOff></xdr:to><xdr:sp macro="" textlink=""><xdr:nvSpPr><xdr:cNvPr id="100" name="Note box" descr="A rounded note box"/><xdr:cNvSpPr/></xdr:nvSpPr><xdr:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="2438400" cy="952500"/></a:xfrm><a:prstGeom prst="roundRect"><a:avLst/></a:prstGeom><a:solidFill><a:schemeClr val="accent1"><a:lumMod val="20000"/><a:lumOff val="80000"/></a:schemeClr></a:solidFill><a:ln w="12700"><a:solidFill><a:schemeClr val="accent1"/></a:solidFill></a:ln></xdr:spPr><xdr:txBody><a:bodyPr vertOverflow="clip" horzOverflow="clip" rtlCol="0" anchor="ctr"/><a:lstStyle/><a:p><a:pPr algn="ctr"/><a:r><a:rPr lang="ja-JP" altLang="en-US" sz="1100"/><a:t>図形のテキスト shape text</a:t></a:r></a:p></xdr:txBody></xdr:sp><xdr:clientData/></xdr:twoCellAnchor>"""


def features():
    wb = Workbook()
    props(wb, "Excel features")
    ws = wb.active
    ws.title = "Features"
    ws["A1"] = "Conditional formats"
    ws["A1"].font = Font(bold=True, size=13)
    vals = [3, 18, 7, 42, 25, 18, 9, 31, 12, 36]
    ws["A2"], ws["B2"], ws["C2"], ws["D2"], ws["E2"] = "scale", "bar", "icons", "cellIs", "top/dup"
    for i, v in enumerate(vals, start=3):
        for col in "ABCDE":
            ws[f"{col}{i}"] = v
    ws.conditional_formatting.add("A3:A12", ColorScaleRule(start_type="min", start_color="F8696B", mid_type="percentile", mid_value=50,
                                                          mid_color="FFEB84", end_type="max", end_color="63BE7B"))
    ws.conditional_formatting.add("B3:B12", DataBarRule(start_type="min", end_type="max", color="638EC6"))
    ws.conditional_formatting.add("C3:C12", IconSetRule("3Arrows", "percent", [0, 33, 67]))
    red = DifferentialStyle(font=Font(color="9C0006"), fill=PatternFill(bgColor="FFC7CE"))
    ws.conditional_formatting.add("D3:D12", Rule(type="cellIs", operator="greaterThan", formula=["20"], dxf=red))
    green = DifferentialStyle(font=Font(bold=True, color="006100"), fill=PatternFill(bgColor="C6EFCE"))
    ws.conditional_formatting.add("E3:E12", Rule(type="top10", rank=3, dxf=green))
    yellow = DifferentialStyle(fill=PatternFill(bgColor="FFEB9C"))
    ws.conditional_formatting.add("E3:E12", Rule(type="duplicateValues", dxf=yellow))
    ws.column_dimensions["G"].width = 12
    ws["G2"] = "text rule"
    for i, s in enumerate(["apple", "banana", "cherry", "grape"], start=3):
        ws[f"G{i}"] = s
    ws.conditional_formatting.add("G3:G6", Rule(type="containsText", operator="containsText", text="an",
                                                formula=['NOT(ISERROR(SEARCH("an",G3)))'], dxf=red))

    # a table with totals
    ws["A15"] = "Table"
    ws["A15"].font = Font(bold=True, size=13)
    data = [("Item", "Units", "Price"), ("Pens", 120, 1.5), ("Paper", 40, 4.25), ("Ink", 15, 22.0), ("Staples", 300, 0.02)]
    for r, row in enumerate(data, start=16):
        for c, v in enumerate(row, start=1):
            ws.cell(row=r, column=c, value=v)
    ws["A21"] = "Total"
    ws["B21"] = "=SUBTOTAL(109,B17:B20)"
    tab = Table(displayName="Supplies", ref="A16:C21", totalsRowCount=1)
    tab.tableStyleInfo = TableStyleInfo(name="TableStyleMedium2", showRowStripes=True)
    ws.add_table(tab)

    # a chart and a picture
    chart = BarChart()
    chart.title = "Units"
    chart.y_axis.title = "units"
    chart.add_data(Reference(ws, min_col=2, min_row=16, max_row=20), titles_from_data=True)
    chart.set_categories(Reference(ws, min_col=1, min_row=17, max_row=20))
    chart.width, chart.height = 12, 7
    ws.add_chart(chart, "E15")
    img = Image(io.BytesIO(png(64, 48, lambda x, y: ((x * 4) % 256, (y * 5) % 256, 160 if (x // 8 + y // 8) % 2 else 60))))
    img.anchor = "I3"
    ws.add_image(img)

    cs = wb.create_chartsheet("Chart")
    line = LineChart()
    line.title = "Values"
    line.add_data(Reference(ws, min_col=1, min_row=2, max_row=12), titles_from_data=True)
    cs.add_chart(line)
    wb.save(os.path.join(OUT, "features.xlsx"))
    add_shape(os.path.join(OUT, "features.xlsx"))
    fix_core(os.path.join(OUT, "features.xlsx"))


def fix_core(path):
    """Sets the modified date openpyxl stamps with the time of saving."""
    rewrite(path, {"docProps/core.xml": lambda s: re.sub(r"(<dcterms:modified[^>]*>)[^<]*", r"\g<1>" + FIXED.isoformat() + "Z", s)})


def rewrite(path, edits):
    """Rewrites parts of a package with functions of their text."""
    with zipfile.ZipFile(path) as z:
        items = [(i, z.read(i.filename)) for i in z.infolist()]
    out = io.BytesIO()
    with zipfile.ZipFile(out, "w", zipfile.ZIP_DEFLATED) as z:
        for info, data in items:
            if info.filename in edits:
                data = edits[info.filename](data.decode("utf-8")).encode("utf-8")
            info.date_time = FIXED.timetuple()[:6]
            z.writestr(info, data)
    with open(path, "wb") as f:
        f.write(out.getvalue())


def add_shape(path):
    """Adds a text box shape to the first drawing (openpyxl writes none)."""
    def edit(s):
        s = s.rstrip()
        # openpyxl writes the spreadsheet drawing namespace as the default one
        shape = SHAPE.replace("<xdr:", "<").replace("</xdr:", "</")
        s = re.sub(r"</wsDr>$", lambda m: shape + m.group(0), s)
        return s.replace("<wsDr ", '<wsDr xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ', 1)
    rewrite(path, {"xl/drawings/drawing1.xml": edit})


if __name__ == "__main__":
    os.makedirs(OUT, exist_ok=True)
    basic()
    features()
    print("wrote", OUT)
