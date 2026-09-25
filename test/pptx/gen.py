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

from pptx import Presentation
from pptx.dml.color import RGBColor
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


def main():
    os.makedirs(OUT, exist_ok=True)
    basic().save(os.path.join(OUT, "basic.pptx"))


if __name__ == "__main__":
    main()
