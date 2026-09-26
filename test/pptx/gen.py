"""Generates the PowerPoint test decks in converter/pptx/testdata with python-pptx.

Usage: python3 test/pptx/gen.py   (pip install python-pptx)
The decks exercise placeholders and their inheritance, bullets and
numbering, text alignment and wrapping (Latin and Japanese), preset shapes,
fills, lines and arrowheads, groups, pictures, tables and hyperlinks.
"""
import io
import os
import struct
import zlib

import copy

from pptx import Presentation
from pptx.chart.data import CategoryChartData
from pptx.dml.color import RGBColor
from pptx.enum.chart import XL_CHART_TYPE, XL_LEGEND_POSITION
from pptx.enum.dml import MSO_LINE_DASH_STYLE
from pptx.enum.shapes import MSO_CONNECTOR, MSO_SHAPE
from pptx.enum.text import MSO_ANCHOR, MSO_AUTO_SIZE, PP_ALIGN
from pptx.oxml.ns import qn
from pptx.util import Emu, Pt, Inches

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "pptx", "testdata")


def png(w, h, fn):
    """A small RGB PNG from fn(x, y) -> (r, g, b)."""
    raw = b""
    for y in range(h):
        raw += b"\x00" + b"".join(bytes(fn(x, y)) for x in range(w))
    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xFFFFFFFF)
    return (b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
            + chunk(b"IDAT", zlib.compress(raw, 9)) + chunk(b"IEND", b""))


def add_text(tf, text, size=None, bold=False, color=None):
    p = tf.paragraphs[0] if not tf.paragraphs[0].runs and tf.paragraphs[0].text == "" else tf.add_paragraph()
    r = p.add_run()
    r.text = text
    if size:
        r.font.size = Pt(size)
    r.font.bold = bold
    if color:
        r.font.color.rgb = RGBColor.from_string(color)
    return p


def basic():
    prs = Presentation()
    layouts = prs.slide_layouts
    # 1: title slide
    s = prs.slides.add_slide(layouts[0])
    s.shapes.title.text = "BDF from PowerPoint"
    s.placeholders[1].text = "Slides rendered directly from DrawingML"

    # 2: title and content with bullets, numbering and a hyperlink
    s = prs.slides.add_slide(layouts[1])
    s.shapes.title.text = "Bullets and levels"
    tf = s.placeholders[1].text_frame
    tf.text = "First level bullet with enough words to wrap onto a second line inside the placeholder"
    p = tf.add_paragraph()
    p.text = "Second level"
    p.level = 1
    p = tf.add_paragraph()
    p.text = "Third level"
    p.level = 2
    p = tf.add_paragraph()
    r = p.add_run()
    r.text = "A link to example.com"
    r.hyperlink.address = "https://example.com/"
    for i, t in enumerate(["numbered one", "numbered two", "numbered three"]):
        p = tf.add_paragraph()
        p.text = t
        p.level = 1
        pPr = p._p.get_or_add_pPr()
        bu = pPr.makeelement(qn("a:buAutoNum"), {"type": "arabicPeriod"})
        pPr.append(bu)

    # 3: shapes, fills and lines
    s = prs.slides.add_slide(layouts[5])
    s.shapes.title.text = "Shapes"
    x = Inches(0.4)
    kinds = [MSO_SHAPE.RECTANGLE, MSO_SHAPE.ROUNDED_RECTANGLE, MSO_SHAPE.OVAL, MSO_SHAPE.ISOSCELES_TRIANGLE,
             MSO_SHAPE.RIGHT_ARROW, MSO_SHAPE.STAR_5_POINT, MSO_SHAPE.HEART, MSO_SHAPE.CAN]
    for i, k in enumerate(kinds):
        sh = s.shapes.add_shape(k, Inches(0.4 + (i % 4) * 2.3), Inches(1.6 + (i // 4) * 1.6), Inches(1.9), Inches(1.2))
        if i == 1:
            sh.fill.gradient()
            sh.fill.gradient_angle = 45
            sh.fill.gradient_stops[0].color.rgb = RGBColor(0xFF, 0xC0, 0x00)
            sh.fill.gradient_stops[1].color.rgb = RGBColor(0xC0, 0x00, 0x40)
        elif i == 2:
            sh.fill.solid()
            sh.fill.fore_color.rgb = RGBColor(0x2E, 0x8B, 0x57)
            sh.line.color.rgb = RGBColor(0, 0, 0)
            sh.line.width = Pt(3)
            sh.line.dash_style = MSO_LINE_DASH_STYLE.DASH
        elif i == 5:
            sh.rotation = 20
        if i in (0, 4):
            sh.text_frame.text = "Text %d" % i
    ln = s.shapes.add_connector(MSO_CONNECTOR.STRAIGHT, Inches(0.5), Inches(5.0), Inches(4.5), Inches(6.4))
    ln.line.width = Pt(2)
    ln.line.color.rgb = RGBColor(0xC0, 0x00, 0x00)
    lnx = ln.line._get_or_add_ln()
    lnx.append(lnx.makeelement(qn("a:tailEnd"), {"type": "triangle"}))
    grp = s.shapes.add_group_shape()
    a = grp.shapes.add_shape(MSO_SHAPE.RECTANGLE, Inches(5.2), Inches(5.0), Inches(1.2), Inches(0.8))
    b = grp.shapes.add_shape(MSO_SHAPE.OVAL, Inches(6.6), Inches(5.3), Inches(1.6), Inches(1.0))
    a.text_frame.text = "grp"
    grp.rotation = 10

    # 4: text alignment, wrapping and Japanese
    s = prs.slides.add_slide(layouts[6])
    y = Inches(0.3)
    for algn, label in [(PP_ALIGN.LEFT, "left"), (PP_ALIGN.CENTER, "center"), (PP_ALIGN.RIGHT, "right"), (PP_ALIGN.JUSTIFY, "justify")]:
        tb = s.shapes.add_textbox(Inches(0.4), y, Inches(4.2), Inches(1.0))
        tf = tb.text_frame
        tf.word_wrap = True
        tf.text = "Aligned %s: the quick brown fox jumps over the lazy dog and keeps running." % label
        tf.paragraphs[0].alignment = algn
        tb.line.color.rgb = RGBColor(0xBB, 0xBB, 0xBB)
        y += Inches(1.15)
    tb = s.shapes.add_textbox(Inches(4.9), Inches(0.3), Inches(4.7), Inches(2.4))
    tf = tb.text_frame
    tf.word_wrap = True
    tf.text = "日本語の文章は、単語の区切りに空白を使わないため、文字と文字の間で改行します。句読点「、」や「。」は行頭に来ません。"
    tf.paragraphs[0].font.size = Pt(18)
    p = tf.add_paragraph()
    r = p.add_run()
    r.text = "混在 mixed テキスト"
    r.font.bold = True
    r = p.add_run()
    r.text = " x"
    r = p.add_run()
    r.text = "2"
    r.font._rPr.set("baseline", "30000")
    r = p.add_run()
    r.text = " underline"
    r.font.underline = True
    tb = s.shapes.add_textbox(Inches(8.4), Inches(3.0), Inches(1.2), Inches(3.8))
    tb.text_frame.text = "縦書きのテキスト"
    tb.text_frame._txBody.bodyPr.set("vert", "eaVert")
    tb.text_frame.paragraphs[0].font.size = Pt(24)
    tb.fill.solid()
    tb.fill.fore_color.rgb = RGBColor(0xF2, 0xF2, 0xF2)
    tb = s.shapes.add_textbox(Inches(4.9), Inches(3.0), Inches(3.2), Inches(1.4))
    tf = tb.text_frame
    tf.word_wrap = True
    tf.vertical_anchor = MSO_ANCHOR.MIDDLE
    tf.text = "Middle anchored"
    tf.paragraphs[0].alignment = PP_ALIGN.CENTER
    tb.fill.solid()
    tb.fill.fore_color.rgb = RGBColor(0xDD, 0xEB, 0xF7)

    # 5: picture and table
    s = prs.slides.add_slide(layouts[5])
    s.shapes.title.text = "Picture and table"
    img = png(64, 48, lambda x, y: (x * 4, y * 5, 160))
    pic = s.shapes.add_picture(io.BytesIO(img), Inches(0.5), Inches(1.6), Inches(3.2), Inches(2.4))
    pic.crop_left = 0.1
    pic.crop_right = 0.1
    rows, cols = 4, 3
    t = s.shapes.add_table(rows, cols, Inches(4.2), Inches(1.6), Inches(5.2), Inches(2.0)).table
    for r in range(rows):
        for c in range(cols):
            t.cell(r, c).text = "Header %d" % c if r == 0 else "r%dc%d" % (r, c)
    t.cell(3, 0).merge(t.cell(3, 1))
    t.cell(3, 0).text = "merged cell with longer text that wraps"

    # 6: hidden slide (skipped by default)
    s = prs.slides.add_slide(layouts[5])
    s.shapes.title.text = "Hidden slide"
    s._element.set("show", "0")
    prs.core_properties.title = "PowerPoint test deck"
    return prs


def move_to(shape, tree):
    """Moves a shape created on a slide into another shape tree (a master or layout)."""
    el = shape._element
    el.getparent().remove(el)
    tree.append(el)


def features():
    prs = Presentation()
    prs.slide_width, prs.slide_height = Emu(12192000), Emu(6858000)
    master = prs.slide_master
    # master background: a vertical gradient
    cSld = master._element.find(qn("p:cSld"))
    bg = cSld.makeelement(qn("p:bg"), {})
    bg.append(bg.makeelement(qn("p:bgPr"), {}))
    grad = etree_fromstring('<a:gradFill xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" rotWithShape="1">'
                            '<a:gsLst><a:gs pos="0"><a:schemeClr val="bg1"/></a:gs><a:gs pos="100000"><a:schemeClr val="accent1"><a:tint val="30000"/></a:schemeClr></a:gs></a:gsLst>'
                            '<a:lin ang="5400000" scaled="0"/></a:gradFill>')
    bg[0].append(grad)
    bg[0].append(bg.makeelement(qn("a:effectLst"), {}))
    cSld.insert(0, bg)
    # the template is 4:3; widen the placeholders of the master and layouts
    sx = 12192000 / 9144000
    for owner in [master] + list(prs.slide_layouts):
        for sh in owner.placeholders:
            left, top, width, height = sh.left, sh.top, sh.width, sh.height
            sh.left, sh.top, sh.width, sh.height = int(left * sx), top, int(width * sx), height
    # master graphics: a band at the bottom and a label, drawn on every slide
    tmp = prs.slides.add_slide(prs.slide_layouts[6])
    band = tmp.shapes.add_shape(MSO_SHAPE.RECTANGLE, 0, Inches(7.0), prs.slide_width, Inches(0.5))
    band.fill.solid()
    band.fill.fore_color.theme_color = 5  # accent1
    band.line.fill.background()
    band.text_frame.text = "BDF features deck"
    band.text_frame.paragraphs[0].font.size = Pt(14)
    move_to(band, master.shapes._spTree)
    # the Title and Content layout adds an accent bar
    bar = tmp.shapes.add_shape(MSO_SHAPE.RECTANGLE, Inches(0.5), Inches(1.45), Inches(3), Inches(0.08))
    bar.fill.solid()
    bar.fill.fore_color.rgb = RGBColor(0xED, 0x7D, 0x31)
    bar.line.fill.background()
    move_to(bar, prs.slide_layouts[1].shapes._spTree)
    rId = prs.slides._sldIdLst[-1].rId
    prs.part.drop_rel(rId)
    del prs.slides._sldIdLst[-1]

    # 1 and 2: same layout; master and layout layers are shared
    for title, items in [("Master and layout", ["The band at the bottom comes from the slide master",
                                                   "The orange bar comes from the layout"]),
                         ("Second slide", ["Same master and layout layers", "Only the body differs"])]:
        s = prs.slides.add_slide(prs.slide_layouts[1])
        s.shapes.title.text = title
        tf = s.placeholders[1].text_frame
        tf.text = items[0]
        for it in items[1:]:
            tf.add_paragraph().text = it

    # 3: charts
    s = prs.slides.add_slide(prs.slide_layouts[5])
    s.shapes.title.text = "Charts"
    cd = CategoryChartData()
    cd.categories = ["Q1", "Q2", "Q3", "Q4"]
    cd.add_series("East", (20.4, 27.4, 90, 20.4))
    cd.add_series("West", (30.6, 38.6, 34.6, 31.6))
    cd.add_series("North", (45.9, 46.9, 45, 43.9))
    ch = s.shapes.add_chart(XL_CHART_TYPE.COLUMN_CLUSTERED, Inches(0.4), Inches(1.6), Inches(4.2), Inches(3.4), cd).chart
    ch.has_legend = True
    ch.legend.position = XL_LEGEND_POSITION.BOTTOM
    ch.legend.include_in_layout = False
    ch = s.shapes.add_chart(XL_CHART_TYPE.LINE_MARKERS, Inches(4.7), Inches(1.6), Inches(4.2), Inches(3.4), cd).chart
    ch.has_legend = True
    pd = CategoryChartData()
    pd.categories = ["Alpha", "Beta", "Gamma", "Delta"]
    pd.add_series("Share", (0.45, 0.25, 0.2, 0.1))
    ch = s.shapes.add_chart(XL_CHART_TYPE.PIE, Inches(9.0), Inches(1.6), Inches(4.0), Inches(3.4), pd).chart
    ch.has_legend = True
    ch.plots[0].has_data_labels = True
    ch.plots[0].data_labels.show_percentage = True
    ch.plots[0].data_labels.show_value = False
    ch.plots[0].data_labels.number_format = "0%"
    sd = CategoryChartData()
    sd.categories = ["A", "B", "C"]
    sd.add_series("Stacked 1", (3, 5, 2))
    sd.add_series("Stacked 2", (2, 1, 4))
    ch = s.shapes.add_chart(XL_CHART_TYPE.BAR_STACKED, Inches(0.4), Inches(5.1), Inches(6), Inches(1.8), sd).chart
    ch.has_legend = False

    # 4: Japanese
    s = prs.slides.add_slide(prs.slide_layouts[1])
    s.shapes.title.text = "日本語のスライド"
    tf = s.placeholders[1].text_frame
    tf.text = "箇条書きの一行目です。長い文章は枠の幅で折り返され、句読点は行頭に来ません。"
    p = tf.add_paragraph()
    p.text = "二段目の項目"
    p.level = 1
    p = tf.add_paragraph()
    p.text = "英数字 ABC と 123 の混在"
    body = s.placeholders[1]
    left, top, height = body.left, body.top, body.height
    body.left, body.top, body.width, body.height = left, top, Inches(8.5), height
    tb = s.shapes.add_textbox(Inches(11.2), Inches(1.6), Inches(1.4), Inches(5.0))
    tb.text_frame.text = "縦書きで、読みます。"
    tb.text_frame._txBody.bodyPr.set("vert", "eaVert")
    tb.text_frame.paragraphs[0].font.size = Pt(28)
    tb.line.color.rgb = RGBColor(0x99, 0x99, 0x99)
    t = s.shapes.add_table(3, 2, Inches(0.8), Inches(5.4), Inches(6), Inches(1.2)).table
    for r, (a, b) in enumerate([("項目", "説明"), ("図形", "プリセット形状"), ("表", "表のスタイル")]):
        t.cell(r, 0).text = a
        t.cell(r, 1).text = b

    # 5: fills, lines and effects
    s = prs.slides.add_slide(prs.slide_layouts[5])
    s.shapes.title.text = "Fills and lines"
    sh = s.shapes.add_shape(MSO_SHAPE.ROUNDED_RECTANGLE, Inches(0.5), Inches(1.7), Inches(2.6), Inches(1.6))
    sh.fill.patterned()
    sh.fill.fore_color.rgb = RGBColor(0x44, 0x72, 0xC4)
    sh.fill.back_color.rgb = RGBColor(0xFF, 0xFF, 0xFF)
    sh.text_frame.text = "Pattern"
    sh.text_frame.paragraphs[0].font.color.rgb = RGBColor(0, 0, 0)
    sh = s.shapes.add_shape(MSO_SHAPE.OVAL, Inches(3.5), Inches(1.7), Inches(2.6), Inches(1.6))
    img = png(32, 32, lambda x, y: (255, 200, 0) if (x // 8 + y // 8) % 2 else (40, 80, 160))
    part, rid = s.part.get_or_add_image_part(io.BytesIO(img))
    spPr = sh._element.spPr
    for f in spPr.findall(qn("a:solidFill")):
        spPr.remove(f)
    bf = etree_fromstring('<a:blipFill xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" '
                          'xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" rotWithShape="1">'
                          '<a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></a:blipFill>' % rid)
    spPr.insert(2, bf)
    for i, dash in enumerate([MSO_LINE_DASH_STYLE.DASH, MSO_LINE_DASH_STYLE.ROUND_DOT, MSO_LINE_DASH_STYLE.LONG_DASH_DOT]):
        ln = s.shapes.add_connector(MSO_CONNECTOR.STRAIGHT, Inches(6.6), Inches(1.9 + i * 0.6), Inches(9.6), Inches(1.9 + i * 0.6))
        ln.line.width = Pt(3)
        ln.line.dash_style = dash
        ln.line.color.rgb = RGBColor(0x20, 0x20, 0x20)
        x = ln.line._get_or_add_ln()
        x.append(x.makeelement(qn("a:headEnd"), {"type": ["oval", "diamond", "stealth"][i]}))
        x.append(x.makeelement(qn("a:tailEnd"), {"type": ["triangle", "arrow", "triangle"][i], "w": "lg", "len": "lg"}))
    ff = s.shapes.build_freeform(10.2, 3.2, scale=Inches(1))
    ff.add_line_segments([(11.0, 1.7), (11.8, 3.2), (11.4, 2.6), (10.6, 2.6)], close=True)
    shp = ff.convert_to_shape()
    shp.fill.solid()
    shp.fill.fore_color.rgb = RGBColor(0x70, 0xAD, 0x47)
    for i, k in enumerate([MSO_SHAPE.CHEVRON, MSO_SHAPE.DONUT, MSO_SHAPE.LIGHTNING_BOLT, MSO_SHAPE.CLOUD, MSO_SHAPE.BLOCK_ARC, MSO_SHAPE.RECTANGULAR_CALLOUT]):
        sh = s.shapes.add_shape(k, Inches(0.5 + i * 2.1), Inches(4.2), Inches(1.8), Inches(1.4))
        if k == MSO_SHAPE.RECTANGULAR_CALLOUT:
            sh.text_frame.text = "Callout"
    prs.core_properties.title = "Features"
    return prs


def etree_fromstring(xml):
    from lxml import etree
    return etree.fromstring(xml)


def main():
    os.makedirs(OUT, exist_ok=True)
    basic().save(os.path.join(OUT, "basic.pptx"))
    features().save(os.path.join(OUT, "features.pptx"))


if __name__ == "__main__":
    main()
