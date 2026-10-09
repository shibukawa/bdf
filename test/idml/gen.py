"""Generates the InDesign test documents in converter/idml/testdata.

Usage: python3 test/idml/gen.py   (standard library only)

The documents are IDML packages written the way InDesign writes them
(design map, spreads, master spreads, stories and resources), so that the
converter is exercised on the markup it will meet:

- basic.idml: three pages in two spreads (a single right page, then a
  facing pair) with a master whose items (a rule and a page number) every
  page inherits; swatches of every kind (CMYK, RGB, Lab, a tint, a linear
  and a radial gradient), a rounded rectangle, an oval with a dashed
  stroke, a polygon with its stroke inside, a dotted line, a rotated
  group, a hidden layer and an invisible item; a placed image embedded in
  the file and one linked by name (link.png beside the document), clipped
  by their frames; and a story of paragraph and character styles (based
  on each other), justification, indents, tabs, a bullet and a numbered
  list, tracking, superscripts and underlines, that flows through three
  threaded frames (two columns on the first page) across the pages, with
  Japanese line breaking.
- vertical.idml: a right-to-left bound document whose story is set
  vertically and flows from the first page to the second.

Text uses M PLUS 1p (the subsets in converter/pptx/testdata/fonts), so the
Japanese text sticks to the characters those subsets hold.
"""
import base64
import os
import struct
import zipfile
import zlib
from xml.sax.saxutils import escape, quoteattr

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "idml", "testdata")
FONT = "M PLUS 1p"
IDPKG = "http://ns.adobe.com/AdobeInDesign/idml/1.0/packaging"
MIMETYPE = "application/vnd.adobe.indesign-idml-package"

JA = ("日本語の文章は、単語の区切りに空白を使わないため、文字と文字の間で改行します。"
      "句読点「、」や「。」は行頭に来ません。")


def png(w, h, fn):
    """A small RGB PNG from fn(x, y) -> (r, g, b)."""
    raw = b""
    for y in range(h):
        raw += b"\x00" + b"".join(bytes(fn(x, y)) for x in range(w))

    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xFFFFFFFF)
    return (b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
            + chunk(b"IDAT", zlib.compress(raw, 9)) + chunk(b"IEND", b""))


def num(v):
    if isinstance(v, bool):
        return "true" if v else "false"
    if isinstance(v, float):
        return repr(round(v, 6))
    return str(v)


def el(name, attrs=None, kids=None, text=None):
    """An XML element as a string; kids are strings, attrs a dict."""
    out = "<" + name
    for k, v in (attrs or {}).items():
        if v is None:
            continue
        out += " %s=%s" % (k, quoteattr(num(v)))
    inner = "".join(kids or [])
    if text is not None:
        inner += escape(text)
    if not inner:
        return out + "/>"
    return out + ">" + inner + "</" + name + ">"


def props(*kids):
    return el("Properties", kids=list(kids))


def typed(name, typ, value):
    return el(name, {"type": typ}, text=num(value))


def unit_list(name, values):
    return el(name, {"type": "list"}, [typed("ListItem", "unit", v) for v in values])


def xml(root):
    return '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>\n' + root


# --- geometry -----------------------------------------------------------------

def point(x, y, l=None, r=None):
    a = "%s %s" % (num(x), num(y))
    return el("PathPointType", {"Anchor": a, "LeftDirection": l or a, "RightDirection": r or a})


def geometry(points, open_=False):
    return props(el("PathGeometry", kids=[
        el("GeometryPathType", {"PathOpen": open_}, [el("PathPointArray", kids=points)])]))


def rect_points(x0, y0, x1, y1):
    # InDesign lists a rectangle's corners counterclockwise from the top left
    return [point(x0, y0), point(x0, y1), point(x1, y1), point(x1, y0)]


def oval_points(x0, y0, x1, y1):
    k = 0.5523
    cx, cy, rx, ry = (x0 + x1) / 2, (y0 + y1) / 2, (x1 - x0) / 2, (y1 - y0) / 2

    def p(ax, ay, lx, ly, rx_, ry_):
        return point(ax, ay, "%s %s" % (num(lx), num(ly)), "%s %s" % (num(rx_), num(ry_)))
    return [
        p(x0, cy, x0, cy + ry * k, x0, cy - ry * k),
        p(cx, y0, cx - rx * k, y0, cx + rx * k, y0),
        p(x1, cy, x1, cy - ry * k, x1, cy + ry * k),
        p(cx, y1, cx + rx * k, y1, cx - rx * k, y1),
    ]


def transform(a=1, b=0, c=0, d=1, tx=0, ty=0):
    return " ".join(num(float(v)) for v in (a, b, c, d, tx, ty))


ITEM_DEFAULTS = {"ContentType": "Unassigned", "StrokeWeight": 0, "FillColor": "Swatch/None",
                 "StrokeColor": "Swatch/None", "GradientFillStart": "0 0", "GradientFillLength": 0, "GradientFillAngle": 0,
                 "StrokeAlignment": "CenterAlignment", "EndCap": "ButtEndCap", "EndJoin": "MiterEndJoin", "MiterLimit": 4,
                 "Locked": False, "Visible": True, "Name": "$ID/"}


def item(kind, self_, pts, layer="ub3", open_=False, extra=None, **attrs):
    a = dict(ITEM_DEFAULTS)
    a.update({"Self": self_, "ItemLayer": layer, "ItemTransform": transform()})
    a.update(attrs)
    return el(kind, a, [geometry(pts, open_)] + (extra or []))


# --- resources ----------------------------------------------------------------

def color(self_, name, space, value, model="Process"):
    return el("Color", {"Self": self_, "Model": model, "Space": space, "ColorValue": value, "Name": name,
                        "ColorEditable": True, "ColorRemovable": True, "Visible": True, "SwatchCreatorID": 7937})


def graphic():
    kids = [
        color("Color/Black", "Black", "CMYK", "0 0 0 100"),
        color("Color/Paper", "Paper", "CMYK", "0 0 0 0"),
        color("Color/Registration", "Registration", "CMYK", "100 100 100 100", "Registration"),
        color("Color/C=100 M=0 Y=0 K=0", "C=100 M=0 Y=0 K=0", "CMYK", "100 0 0 0"),
        color("Color/C=0 M=100 Y=0 K=0", "C=0 M=100 Y=0 K=0", "CMYK", "0 100 0 0"),
        color("Color/C=0 M=0 Y=100 K=0", "C=0 M=0 Y=100 K=0", "CMYK", "0 0 100 0"),
        color("Color/Navy", "Navy", "CMYK", "100 80 0 30"),
        color("Color/Red", "Red", "RGB", "220 40 40"),
        color("Color/Leaf", "Leaf", "LAB", "70 -50 50"),
        el("Tint", {"Self": "Tint/u20", "Name": "Cyan 30%", "BaseColor": "Color/C=100 M=0 Y=0 K=0", "TintValue": 30,
                    "ColorEditable": True, "ColorRemovable": True, "Visible": True}),
        el("Gradient", {"Self": "Gradient/Sunset", "Type": "Linear", "Name": "Sunset", "ColorEditable": True,
                        "ColorRemovable": True, "Visible": True}, [
            el("GradientStop", {"Self": "Gradient/Sunset GradientStop 0", "StopColor": "Color/C=0 M=0 Y=100 K=0", "Location": 0, "Midpoint": 50}),
            el("GradientStop", {"Self": "Gradient/Sunset GradientStop 1", "StopColor": "Color/C=0 M=100 Y=0 K=0", "Location": 100, "Midpoint": 50})]),
        el("Gradient", {"Self": "Gradient/Glow", "Type": "Radial", "Name": "Glow", "ColorEditable": True,
                        "ColorRemovable": True, "Visible": True}, [
            el("GradientStop", {"Self": "Gradient/Glow GradientStop 0", "StopColor": "Color/Paper", "Location": 0, "Midpoint": 50}),
            el("GradientStop", {"Self": "Gradient/Glow GradientStop 1", "StopColor": "Color/Navy", "Location": 100, "Midpoint": 50})]),
        el("Swatch", {"Self": "Swatch/None", "Name": "None", "ColorEditable": False, "ColorRemovable": False, "Visible": True}),
    ]
    for name in ["Solid", "Dashed", "Dotted", "Canned Dashed 3x2", "Canned Dashed 4x4", "Canned dotted", "ThinThin", "ThickThin"]:
        kids.append(el("StrokeStyle", {"Self": "StrokeStyle/$ID/" + name, "Name": "$ID/" + name}))
    kids.append(el("DashedStrokeStyle", {"Self": "StrokeStyle/u30", "Name": "Long dash", "LineCap": "ButtEndCap",
                                         "StrokeCornerAdjustment": "None"},
                   [props(unit_list("DashArray", [9, 3, 3, 3]))]))
    kids.append(el("DottedStrokeStyle", {"Self": "StrokeStyle/u31", "Name": "Wide dots", "DotGap": 6, "StrokeCornerAdjustment": "None"}))
    return xml(el("idPkg:Graphic", {"xmlns:idPkg": IDPKG, "DOMVersion": "8.0"}, kids))


def pstyle(self_, name, based="$ID/[No paragraph style]", font=FONT, leading="Auto", tabs=None, bullet=None, **attrs):
    kids = [typed("BasedOn", "string" if based.startswith("$ID/") else "object", based),
            typed("Leading", "enumeration" if leading == "Auto" else "unit", leading),
            typed("AppliedFont", "string", font)]
    if tabs:
        kids.append(el("TabList", {"type": "list"}, [
            el("ListItem", {"type": "record"}, [typed("Alignment", "enumeration", a), typed("AlignmentCharacter", "string", "."),
                                                typed("Leader", "string", ""), typed("Position", "unit", p)]) for a, p in tabs]))
    if bullet:
        kids.append(el("BulletChar", {"BulletCharacterType": "UnicodeOnly", "BulletCharacterValue": bullet}))
    a = {"Self": self_, "Name": name}
    a.update(attrs)
    return el("ParagraphStyle", a, [props(*kids)])


def styles():
    root_p = pstyle("ParagraphStyle/$ID/[No paragraph style]", "$ID/[No paragraph style]", PointSize=12, Justification="LeftAlign",
                    AutoLeading=120, AppliedLanguage="$ID/Japanese", FontStyle="Regular", FillColor="Color/Black",
                    LeftIndent=0, FirstLineIndent=0, SpaceBefore=0, SpaceAfter=0, Tracking=0)
    normal = pstyle("ParagraphStyle/$ID/NormalParagraphStyle", "$ID/NormalParagraphStyle", PointSize=10, AutoLeading=150)
    body = pstyle("ParagraphStyle/Body", "Body", "ParagraphStyle/$ID/NormalParagraphStyle", Justification="LeftJustified",
                  FirstLineIndent=10, SpaceAfter=4)
    heading = pstyle("ParagraphStyle/Heading", "Heading", "ParagraphStyle/Body", PointSize=14, FontStyle="Bold",
                     FillColor="Color/Navy", FirstLineIndent=0, SpaceBefore=8, SpaceAfter=4, leading=18)
    bullets = pstyle("ParagraphStyle/Bullets", "Bullets", "ParagraphStyle/Body", FirstLineIndent=-10, LeftIndent=14,
                     BulletsAndNumberingListType="BulletList", bullet=9679)
    numbers = pstyle("ParagraphStyle/Numbers", "Numbers", "ParagraphStyle/Body", FirstLineIndent=-12, LeftIndent=14,
                     BulletsAndNumberingListType="NumberedList", NumberingFormat="^#.")
    tabbed = pstyle("ParagraphStyle/Tabbed", "Tabbed", "ParagraphStyle/Body", FirstLineIndent=0, Justification="LeftAlign",
                    tabs=[("LeftAlign", 50), ("RightAlign", 110)])
    centered = pstyle("ParagraphStyle/Centered", "Centered", "ParagraphStyle/Body", Justification="CenterAlign", FirstLineIndent=0)
    title = pstyle("ParagraphStyle/Title", "Title", "ParagraphStyle/Centered", PointSize=24, FontStyle="Bold",
                   FillColor="Color/C=100 M=0 Y=0 K=0", Tracking=100, leading=30)
    group = el("ParagraphStyleGroup", {"Self": "u40", "Name": "Notes"}, [
        pstyle("ParagraphStyle/Notes%3aSmall", "Small", "ParagraphStyle/Body", PointSize=8, FillColor="Tint/u20",
               Justification="RightAlign", FirstLineIndent=0)])
    cs = [
        el("CharacterStyle", {"Self": "CharacterStyle/$ID/[No character style]", "Name": "$ID/[No character style]"}),
        el("CharacterStyle", {"Self": "CharacterStyle/Emphasis", "Name": "Emphasis", "FontStyle": "Bold", "FillColor": "Color/Red"},
           [props(typed("BasedOn", "string", "$ID/[No character style]"))]),
        el("CharacterStyle", {"Self": "CharacterStyle/Code", "Name": "Code", "Underline": True, "Tracking": 50},
           [props(typed("BasedOn", "object", "CharacterStyle/Emphasis"), typed("AppliedFont", "string", FONT))]),
    ]
    return xml(el("idPkg:Styles", {"xmlns:idPkg": IDPKG, "DOMVersion": "8.0"}, [
        el("RootCharacterStyleGroup", {"Self": "u7a"}, cs),
        el("RootParagraphStyleGroup", {"Self": "u7b"}, [root_p, normal, body, heading, bullets, numbers, tabbed, centered, title, group]),
    ]))


def preferences(w, h, pages, binding="LeftToRight"):
    return xml(el("idPkg:Preferences", {"xmlns:idPkg": IDPKG, "DOMVersion": "8.0"}, [
        el("DocumentPreference", {"PageHeight": h, "PageWidth": w, "PagesPerDocument": pages, "FacingPages": True,
                                  "PageBinding": binding, "DocumentBleedTopOffset": 0, "PageOrientation": "Portrait"})]))


def fonts():
    return xml(el("idPkg:Fonts", {"xmlns:idPkg": IDPKG, "DOMVersion": "8.0"}, [
        el("FontFamily", {"Self": "di50", "Name": FONT}, [
            el("Font", {"Self": "di50FontnM PLUS 1p Regular", "FontFamily": FONT, "Name": FONT + " Regular", "PostScriptName": "MPLUS1p-Regular",
                        "Status": "Installed", "FontStyleName": "Regular", "FontType": "OpenTypeTT"}),
            el("Font", {"Self": "di50FontnM PLUS 1p Bold", "FontFamily": FONT, "Name": FONT + " Bold", "PostScriptName": "MPLUS1p-Bold",
                        "Status": "Installed", "FontStyleName": "Bold", "FontType": "OpenTypeTT"})])]))


def metadata(title):
    return ('<?xml version="1.0" encoding="UTF-8"?>\n'
            '<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">'
            '<rdf:Description rdf:about="" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xmp="http://ns.adobe.com/xap/1.0/">'
            '<dc:title><rdf:Alt><rdf:li xml:lang="x-default">%s</rdf:li></rdf:Alt></dc:title>'
            '<dc:creator><rdf:Seq><rdf:li>BDF tests</rdf:li></rdf:Seq></dc:creator>'
            '<xmp:CreateDate>2026-01-02T03:04:05Z</xmp:CreateDate>'
            '</rdf:Description></rdf:RDF></x:xmpmeta>' % escape(title))


# --- stories ------------------------------------------------------------------

def content(text):
    return el("Content", text=text)


BR = el("Br")


def csr(kids, style="CharacterStyle/$ID/[No character style]", **attrs):
    a = {"AppliedCharacterStyle": style}
    a.update(attrs)
    return el("CharacterStyleRange", a, kids)


def psr(kids, style="ParagraphStyle/$ID/NormalParagraphStyle", **attrs):
    a = {"AppliedParagraphStyle": style}
    a.update(attrs)
    return el("ParagraphStyleRange", a, kids)


def para(style, *runs, **attrs):
    """A paragraph of runs: strings, or (text, character style, attrs) tuples."""
    kids = []
    for r in runs:
        if isinstance(r, str):
            kids.append(csr([content(r)]))
        else:
            text, cstyle = r[0], r[1]
            extra = r[2] if len(r) > 2 else {}
            kids.append(csr([content(text)], cstyle or "CharacterStyle/$ID/[No character style]", **extra))
    kids.append(csr([BR]))
    return psr(kids, style, **attrs)


def story(self_, paras, vertical=False):
    return xml(el("idPkg:Story", {"xmlns:idPkg": IDPKG, "DOMVersion": "8.0"}, [
        el("Story", {"Self": self_, "AppliedTOCStyle": "n", "TrackChanges": False, "StoryTitle": "$ID/", "AppliedNamedGrid": "n"}, [
            el("StoryPreference", {"OpticalMarginAlignment": False, "OpticalMarginSize": 12, "FrameType": "TextFrameType",
                                   "StoryOrientation": "Vertical" if vertical else "Horizontal",
                                   "StoryDirection": "LeftToRightDirection"}),
            el("InCopyExportOption", {"IncludeGraphicProxies": True, "IncludeAllResources": False})] + paras)]))


def body_story():
    return story("u100", [
        para("ParagraphStyle/Title", "InDesign サンプル"),
        para("ParagraphStyle/Centered", ("IDML からへんかんしたページです。", None, {"PointSize": 9})),
        para("ParagraphStyle/Heading", "テキストのながれ"),
        para("ParagraphStyle/Body", "このだんらくは ", ("スタイル", "CharacterStyle/Emphasis"),
             " を使っています。テキストフレームは2だんぐみで、はいりきらない文章はつぎのページのフレームにつながります。"),
        para("ParagraphStyle/Body", JA),
        para("ParagraphStyle/Body", "Latin text wraps at spaces, and a very long paragraph of English prose keeps flowing from "
             "one column to the next until the frame is full, then continues in the frame it is threaded to on the following "
             "page. Superscript", ("2", None, {"Position": "Superscript"}), " and ", ("underline", "CharacterStyle/Code"),
             " and ", ("tracking", None, {"Tracking": 200}), " are character attributes."),
        para("ParagraphStyle/Heading", "リスト"),
        para("ParagraphStyle/Bullets", "ひとつめのこうもく"),
        para("ParagraphStyle/Bullets", "ふたつめのこうもく、すこしながいテキストで、おりかえしがおきます。"),
        para("ParagraphStyle/Numbers", "はじめに"),
        para("ParagraphStyle/Numbers", "つぎに"),
        para("ParagraphStyle/Numbers", "さいごに"),
        para("ParagraphStyle/Tabbed", "name\tvalue\t12.5"),
        para("ParagraphStyle/Tabbed", "size\tlarge\t3"),
        para("ParagraphStyle/Heading", "ながい文章"),
    ] + [para("ParagraphStyle/Body", JA + " " + JA) for _ in range(4)] + [
        para("ParagraphStyle/Notes%3aSmall", "グループのスタイル（みぎぞろえ、ティント）"),
        para("ParagraphStyle/Body", JA),
        para("ParagraphStyle/Body", "Line break after a forced break, and the story ends."),
    ])


def number_story():
    return story("u101", [psr([csr([el("Content", kids=["- ", "<?ACE 18?>", " -"])], PointSize=9), csr([BR])], "ParagraphStyle/Centered")])


def label_story(self_, text, **attrs):
    return story(self_, [para("ParagraphStyle/Centered", (text, None, attrs))])


def vertical_story():
    return story("u200", [
        para("ParagraphStyle/Heading", "縦書きのテキスト"),
        para("ParagraphStyle/Body", JA),
        para("ParagraphStyle/Body", "Latin 123 はよこにたおれ、「かぎかっこ」は縦の形になります。"),
    ] + [para("ParagraphStyle/Body", JA) for _ in range(6)], vertical=True)


# --- spreads ------------------------------------------------------------------

def page(self_, x, y, w, h, master="ub5"):
    return el("Page", {"Self": self_, "AppliedMaster": master, "GeometricBounds": "0 0 %s %s" % (num(h), num(w)),
                       "ItemTransform": transform(tx=x, ty=y), "Name": self_, "OverrideList": "", "MasterPageTransform": transform()})


def text_frame(self_, story_, x0, y0, x1, y1, prev="n", nxt="n", cols=1, gutter=12, inset=0, vjust="TopAlign", **attrs):
    pref = el("TextFramePreference", {"TextColumnCount": cols, "TextColumnGutter": gutter, "TextColumnFixedWidth": 0,
                                      "UseFixedColumnWidth": False, "VerticalJustification": vjust, "FirstBaselineOffset": "AscentOffset"},
              [props(unit_list("InsetSpacing", [inset] * 4))])
    return item("TextFrame", self_, rect_points(x0, y0, x1, y1), extra=[pref], ContentType="TextType", ParentStory=story_,
                PreviousTextFrame=prev, NextTextFrame=nxt, **attrs)


def spread(self_, pages, items, kind="Spread", binding=1, **attrs):
    a = {"Self": self_, "PageCount": len(pages), "BindingLocation": binding, "ShowMasterItems": True}
    a.update(attrs)
    return xml(el("idPkg:" + kind, {"xmlns:idPkg": IDPKG, "DOMVersion": "8.0"}, [el(kind, a, pages + items)]))


def master_spread(w, h):
    # two pages side by side, the spine at x = 0
    pages = [page("ub6", -w, -h / 2, w, h, "n"), page("ub7", 0, -h / 2, w, h, "n")]
    items = []
    for i, x in enumerate((-w, 0)):
        y = -h / 2
        items.append(item("GraphicLine", "umr%d" % i, [point(x + 30, y + 30), point(x + w - 30, y + 30)], open_=True,
                          StrokeColor="Color/Navy", StrokeWeight=0.75))
        items.append(text_frame("umt%d" % i, "u101", x + 30, y + h - 36, x + w - 30, y + h - 20))
    return spread("ub5", pages, items, "MasterSpread", Name="A-Master", NamePrefix="A", BaseName="Master", AppliedMaster="n")


def embedded_image():
    data = png(60, 40, lambda x, y: ((x * 4) % 256, (y * 6) % 256, 160))
    return base64.b64encode(data).decode("ascii")


def image(self_, tx, ty, scale, w, h, embedded=None, link=None):
    kids = [typed("Profile", "string", "$ID/Embedded" if embedded else "$ID/None"),
            el("GraphicBounds", {"Left": 0, "Top": 0, "Right": w, "Bottom": h})]
    if embedded:
        kids.append(el("Contents", text=embedded))
    link_el = el("Link", {"Self": self_ + "l", "LinkResourceURI": link or "file:/tmp/embedded.png",
                          "StoredState": "Embedded" if embedded else "Normal", "LinkClassID": 35906, "LinkClientID": 257})
    return el("Image", {"Self": self_, "ImageTypeName": "$ID/PNG", "ActualPpi": "72 72", "EffectivePpi": "%s %s" % (num(72 / scale), num(72 / scale)),
                        "ItemTransform": transform(scale, 0, 0, scale, tx, ty)}, [props(*kids), link_el])


def basic_spread1(w, h):
    # a single right page: the spread's origin is on its left edge
    y = -h / 2
    pg = page("udb", 0, y, w, h)
    items = [
        # a rounded rectangle with a gradient fill (stroke inside)
        item("Rectangle", "u10a", rect_points(30, y + 50, 200, y + 110), FillColor="Gradient/Sunset", GradientFillAngle=30,
             StrokeColor="Color/Navy", StrokeWeight=3, StrokeAlignment="InsideAlignment", CornerRadius=12,
             TopLeftCornerOption="RoundedCorner", TopRightCornerOption="RoundedCorner",
             BottomLeftCornerOption="RoundedCorner", BottomRightCornerOption="RoundedCorner"),
        # an oval with a tinted fill and a dashed stroke
        item("Oval", "u10b", oval_points(220, y + 50, 300, y + 110), FillColor="Tint/u20", StrokeColor="Color/Red",
             StrokeWeight=2, StrokeType="StrokeStyle/$ID/Canned Dashed 3x2", EndCap="RoundEndCap"),
        # a radial gradient in an oval, half transparent
        item("Oval", "u10c", oval_points(310, y + 50, 370, y + 110), FillColor="Gradient/Glow",
             extra=[el("TransparencySetting", kids=[el("BlendingSetting", {"Opacity": 60, "BlendMode": "Normal"})])]),
        # a triangle with a Lab fill and its stroke outside
        item("Polygon", "u10d", [point(30, y + 180), point(80, y + 120), point(130, y + 180)], FillColor="Color/Leaf",
             StrokeColor="Color/Black", StrokeWeight=2, StrokeAlignment="OutsideAlignment", EndJoin="RoundEndJoin"),
        # a dotted line and a custom dashed one
        item("GraphicLine", "u10e", [point(150, y + 130), point(370, y + 130)], open_=True, StrokeColor="Color/Black",
             StrokeWeight=3, StrokeType="StrokeStyle/$ID/Dotted"),
        item("GraphicLine", "u10f", [point(150, y + 150), point(370, y + 150)], open_=True, StrokeColor="Color/C=0 M=100 Y=0 K=0",
             StrokeWeight=2, StrokeType="StrokeStyle/u30"),
        item("GraphicLine", "u110", [point(150, y + 170), point(370, y + 170)], open_=True, StrokeColor="Color/Navy",
             StrokeWeight=2, StrokeType="StrokeStyle/u31"),
        # a rotated group: a square and a circle, turned 20 degrees about the group's origin
        el("Group", {"Self": "u111", "ItemLayer": "ub3", "ItemTransform": transform(0.9397, 0.342, -0.342, 0.9397, 300, y + 130),
                     "Visible": True, "Locked": False}, [
            item("Rectangle", "u112", rect_points(0, 0, 40, 40), FillColor="Color/C=0 M=0 Y=100 K=0", StrokeColor="Color/Black", StrokeWeight=1),
            item("Oval", "u113", oval_points(20, 20, 60, 60), FillColor="Color/Red", FillTint=50)]),
        # an item on the hidden layer, and an invisible item: neither shows
        item("Rectangle", "u114", rect_points(30, y + 200, 370, y + 240), layer="ub4", FillColor="Color/Red"),
        item("Rectangle", "u115", rect_points(30, y + 200, 370, y + 240), FillColor="Color/Red", Visible=False),
        # the title, and the body text in two columns
        text_frame("u120", "u100", 30, y + 200, w - 30, y + h - 50, nxt="u121", cols=2, gutter=14, inset=4),
    ]
    return spread("ucc", [pg], items, binding=0)


def basic_spread2(w, h, embedded, link):
    y = -h / 2
    pages = [page("udc", -w, y, w, h), page("udd", 0, y, w, h)]
    items = [
        # left page: the story goes on, beside a placed image clipped by an oval
        text_frame("u121", "u100", -w + 30, y + 50, -w / 2 - 6, y + h - 50, prev="u120", nxt="u122", inset=0),
        item("Oval", "u130", oval_points(-w / 2 + 10, y + 50, -30, y + 150), StrokeColor="Color/Navy", StrokeWeight=1,
             extra=[image("u131", -w / 2 + 10, y + 50, 2.5, 60, 40, embedded=embedded)]),
        # the same image linked by name, scaled down and offset in its frame, with a Paper fill
        item("Rectangle", "u132", rect_points(-w / 2 + 10, y + 170, -30, y + 260), FillColor="Color/Paper",
             StrokeColor="Color/Black", StrokeWeight=1,
             extra=[image("u133", -w / 2 + 20, y + 180, 1.2, 60, 40, link=link)]),
        # a frame whose text is overset: the frame is too small for it
        text_frame("u134", "u102", -w / 2 + 10, y + 280, -30, y + 310, inset=2, StrokeColor="Color/Red", StrokeWeight=0.5),
        # right page: the end of the story, centered vertically in its frame
        text_frame("u122", "u100", 30, y + 50, w - 30, y + h - 50, prev="u121", vjust="CenterAlign",
                   FillColor="Tint/u20", FillTint=40),
    ]
    return spread("udx", pages, items, binding=1)


def designmap(spreads, masters, stories, sections, layers):
    kids = [
        el("Language", {"Self": "Language/$ID/Japanese", "Name": "$ID/Japanese", "SingleQuotes": "‘’", "DoubleQuotes": "“”"}),
        el("idPkg:Graphic", {"src": "Resources/Graphic.xml"}),
        el("idPkg:Fonts", {"src": "Resources/Fonts.xml"}),
        el("idPkg:Styles", {"src": "Resources/Styles.xml"}),
        el("idPkg:Preferences", {"src": "Resources/Preferences.xml"}),
    ]
    for self_, name, visible in layers:
        kids.append(el("Layer", {"Self": self_, "Name": name, "Visible": visible, "Locked": False, "IgnoreWrap": False,
                                 "ShowGuides": True, "LockGuides": False, "UI": True, "Expendable": True, "Printable": True}))
    for m in masters:
        kids.append(el("idPkg:MasterSpread", {"src": "MasterSpreads/MasterSpread_%s.xml" % m}))
    for s in spreads:
        kids.append(el("idPkg:Spread", {"src": "Spreads/Spread_%s.xml" % s}))
    for self_, start, number, length in sections:
        kids.append(el("Section", {"Self": self_, "Length": length, "Name": "", "ContinueNumbering": False,
                                   "PageNumberStart": number, "PageStart": start, "PageNumberStyle": "Arabic", "Marker": ""}))
    for s in stories:
        kids.append(el("idPkg:Story", {"src": "Stories/Story_%s.xml" % s}))
    return ('<?xml version="1.0" encoding="UTF-8" standalone="yes"?>\n'
            '<?aid style="50" type="document" readerVersion="6.0" featureSet="257" product="8.0(370)" ?>\n'
            + el("Document", {"xmlns:idPkg": IDPKG, "DOMVersion": "8.0", "Self": "d", "StoryList": " ".join(stories),
                              "Name": "test.indd", "ZeroPoint": "0 0", "ActiveLayer": layers[0][0]}, kids))


def container():
    return ('<?xml version="1.0" encoding="UTF-8"?>\n'
            '<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container" version="1.0">'
            '<rootfiles><rootfile full-path="designmap.xml" media-type="text/xml"/></rootfiles></container>')


def write(name, parts):
    os.makedirs(OUT, exist_ok=True)
    path = os.path.join(OUT, name)
    with zipfile.ZipFile(path, "w") as z:
        # the mimetype comes first, stored, as InDesign writes it
        info = zipfile.ZipInfo("mimetype", date_time=(2026, 1, 2, 3, 4, 5))
        z.writestr(info, MIMETYPE, compress_type=zipfile.ZIP_STORED)
        for n, data in parts:
            info = zipfile.ZipInfo(n, date_time=(2026, 1, 2, 3, 4, 5))
            info.compress_type = zipfile.ZIP_DEFLATED
            z.writestr(info, data.encode("utf-8") if isinstance(data, str) else data)
    print("wrote", os.path.relpath(path, ROOT))


def basic():
    w, h = 400, 560
    os.makedirs(OUT, exist_ok=True)
    link_png = png(60, 40, lambda x, y: (40, (x * 4) % 256, (y * 6) % 256))
    with open(os.path.join(OUT, "link.png"), "wb") as f:
        f.write(link_png)
    print("wrote", os.path.relpath(os.path.join(OUT, "link.png"), ROOT))
    parts = [
        ("META-INF/container.xml", container()),
        ("META-INF/metadata.xml", metadata("IDML のテスト")),
        ("designmap.xml", designmap(["ucc", "udx"], ["ub5"], ["u100", "u101", "u102"],
                                    [("ud8", "udb", 1, 3)], [("ub3", "Layer 1", True), ("ub4", "Hidden", False)])),
        ("Resources/Graphic.xml", graphic()),
        ("Resources/Fonts.xml", fonts()),
        ("Resources/Styles.xml", styles()),
        ("Resources/Preferences.xml", preferences(w, h, 3)),
        ("MasterSpreads/MasterSpread_ub5.xml", master_spread(w, h)),
        ("Spreads/Spread_ucc.xml", basic_spread1(w, h)),
        ("Spreads/Spread_udx.xml", basic_spread2(w, h, embedded_image(), "file:/Users/someone/Links/link.png")),
        ("Stories/Story_u100.xml", body_story()),
        ("Stories/Story_u101.xml", number_story()),
        ("Stories/Story_u102.xml", label_story("u102", "このフレームにははいりきらないながいテキストがあります。", PointSize=9)),
    ]
    write("basic.idml", parts)


def vertical():
    w, h = 300, 400
    y = -h / 2
    # bound on the right: the first page is a left page, the second a right one
    s1 = spread("ucc", [page("udb", -w, y, w, h)],
                [text_frame("u220", "u200", -w + 30, y + 30, -30, y + h - 40, nxt="u221", cols=2, gutter=16)], binding=1)
    s2 = spread("udx", [page("udc", 0, y, w, h)],
                [text_frame("u221", "u200", 30, y + 30, w - 30, y + h - 40, prev="u220")], binding=0)
    parts = [
        ("META-INF/container.xml", container()),
        ("META-INF/metadata.xml", metadata("縦書き")),
        ("designmap.xml", designmap(["ucc", "udx"], ["ub5"], ["u200", "u101"],
                                    [("ud8", "udb", 1, 2)], [("ub3", "Layer 1", True)])),
        ("Resources/Graphic.xml", graphic()),
        ("Resources/Fonts.xml", fonts()),
        ("Resources/Styles.xml", styles()),
        ("Resources/Preferences.xml", preferences(w, h, 2, "RightToLeft")),
        ("MasterSpreads/MasterSpread_ub5.xml", master_spread(w, h)),
        ("Spreads/Spread_ucc.xml", s1),
        ("Spreads/Spread_udx.xml", s2),
        ("Stories/Story_u200.xml", vertical_story()),
        ("Stories/Story_u101.xml", number_story()),
    ]
    write("vertical.idml", parts)


if __name__ == "__main__":
    basic()
    vertical()
