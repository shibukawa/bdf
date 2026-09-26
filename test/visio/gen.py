"""Generates the Visio test drawings in converter/visio/testdata.

Usage: python3 test/visio/gen.py   (standard library only)

The drawings are written from one model into both formats the converter
reads, so that the tests can check that a .vsdx package and a .vdx XML
drawing of the same content convert to the same objects:

- shapes.vsdx: a Visio 2013 package with a dynamic theme. Themed shapes and
  variant quick styles, the geometry rows (arcs, elliptical arcs, ellipses,
  polylines, NURBS, B-splines, relative Béziers), fills (patterns,
  gradients, transparency), line patterns and weights, all 45 arrowheads,
  rounding, rotation, flips, groups, shadows, a picture, layers, links and a
  background page; a second page of text: alignment, character formats,
  bullets, tabs, Japanese line breaking, vertical and rotated text, text
  backgrounds, a connector label and line jumps where connectors cross.
- flow.vsdx and flow.vdx: the same flowchart without a theme, with masters,
  styles, a picture and a background page, in both formats.

Text uses M PLUS 1p (the subsets in converter/pptx/testdata/fonts), so the
Japanese text sticks to the characters those subsets hold.
"""
import base64
import math
import os
import struct
import zipfile
import zlib
from xml.sax.saxutils import escape, quoteattr

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "visio", "testdata")
FONT = "M PLUS 1p"
NS = "http://schemas.microsoft.com/office/visio/2012/main"
NS2003 = "http://schemas.microsoft.com/visio/2003/core"
RELNS = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
PT = 1 / 72  # a point in inches


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
        return "1" if v else "0"
    if isinstance(v, float):
        return repr(round(v, 10))
    return str(v)


# --- the sheet model ---------------------------------------------------------

class Row:
    def __init__(self, key, typ=None, **cells):
        self.key, self.typ, self.cells = str(key), typ, [(k, v) for k, v in cells.items()]


class Section:
    def __init__(self, name, ix=None, cells=None, rows=None):
        self.name, self.ix = name, ix
        self.cells = list((cells or {}).items())
        self.rows = rows or []


class Sheet:
    def __init__(self, **cells):
        self.cells = []  # (name, value) with value a number or (value, formula)
        self.sections = []
        self.text = None  # segments: str, ("cp", ix), ("pp", ix), ("fld", ix, text)
        self.foreign = None  # (type, compression, data)
        self.set(**cells)

    def set(self, **cells):
        for k, v in cells.items():
            self.cells = [(n, x) for n, x in self.cells if n != k] + [(k, v)]
        return self

    def section(self, s):
        self.sections.append(s)
        return self


class Shape(Sheet):
    def __init__(self, id, typ="Shape", name=None, master=None, master_shape=None, styles=None, **cells):
        super().__init__(**cells)
        self.id, self.typ, self.name = id, typ, name
        self.master, self.master_shape, self.styles = master, master_shape, styles
        self.kids = []


class Style(Sheet):
    def __init__(self, id, name, parents=None, **cells):
        super().__init__(**cells)
        self.id, self.name, self.parents = id, name, parents


class Master:
    def __init__(self, id, name, shapes):
        self.id, self.name, self.shapes = id, name, shapes


class Page:
    def __init__(self, id, name, sheet, shapes, background=False, back_page=None):
        self.id, self.name, self.sheet, self.shapes = id, name, sheet, shapes
        self.background, self.back_page = background, back_page


class Drawing:
    def __init__(self, title, styles, masters, pages, colors=None, faces=None, theme=None, language="en-US"):
        self.title, self.styles, self.masters, self.pages = title, styles, masters, pages
        self.colors, self.faces, self.theme, self.language = colors or {}, faces or {}, theme, language


def value(v):
    """A cell value and formula: (value, formula) or a plain value."""
    if isinstance(v, tuple):
        return num(v[0]), v[1]
    return num(v), None


# --- .vsdx -------------------------------------------------------------------

def vsdx_cells(sheet):
    out = []
    for n, v in sheet.cells:
        val, f = value(v)
        out.append("<Cell N=%s V=%s%s/>" % (quoteattr(n), quoteattr(val), " F=%s" % quoteattr(f) if f else ""))
    for s in sheet.sections:
        attrs = " N=%s" % quoteattr(s.name) + (" IX='%d'" % s.ix if s.ix is not None else "")
        body = "".join("<Cell N=%s V=%s/>" % (quoteattr(n), quoteattr(value(v)[0])) for n, v in s.cells)
        for r in s.rows:
            key = " IX='%s'" % r.key if r.key.isdigit() else " N=%s" % quoteattr(r.key)
            typ = " T=%s" % quoteattr(r.typ) if r.typ else ""
            cells = ""
            for n, v in r.cells:
                val, f = value(v)
                cells += "<Cell N=%s V=%s%s/>" % (quoteattr(n), quoteattr(val), " F=%s" % quoteattr(f) if f else "")
            body += "<Row%s%s>%s</Row>" % (typ, key, cells)
        out.append("<Section%s>%s</Section>" % (attrs, body))
    return "".join(out)


def text_xml(segments):
    out = ""
    for s in segments:
        if isinstance(s, str):
            out += escape(s)
        elif s[0] == "fld":
            out += "<fld IX='%d'>%s</fld>" % (s[1], escape(s[2]))
        else:
            out += "<%s IX='%d'/>" % (s[0], s[1])
    return "<Text>%s</Text>" % out


def shape_attrs(s):
    a = " ID='%d' Type=%s" % (s.id, quoteattr(s.typ))
    if s.name:
        a += " NameU=%s Name=%s" % (quoteattr(s.name), quoteattr(s.name))
    if s.master is not None:
        a += " Master='%d'" % s.master
    if s.master_shape is not None:
        a += " MasterShape='%d'" % s.master_shape
    if s.styles:
        a += " LineStyle='%d' FillStyle='%d' TextStyle='%d'" % s.styles
    return a


class VSDXPart:
    """A part being written with the relationships its images need."""

    def __init__(self):
        self.rels = []  # (id, type, target)
        self.media = []  # (name, data)

    def rel(self, typ, target):
        rid = "rId%d" % (len(self.rels) + 1)
        self.rels.append((rid, typ, target))
        return rid


def vsdx_shape(s, part, media_counter):
    body = vsdx_cells(s)
    if s.text is not None:
        body += text_xml(s.text)
    if s.foreign is not None:
        typ, comp, data = s.foreign
        media_counter[0] += 1
        name = "image%d.%s" % (media_counter[0], comp.lower())
        part.media.append((name, data))
        rid = part.rel("http://schemas.microsoft.com/visio/2010/relationships/image", "../media/" + name)
        body += "<ForeignData ForeignType='%s' CompressionType='%s'><Rel r:id='%s'/></ForeignData>" % (typ, comp, rid)
    if s.kids:
        body += "<Shapes>%s</Shapes>" % "".join(vsdx_shape(k, part, media_counter) for k in s.kids)
    return "<Shape%s>%s</Shape>" % (shape_attrs(s), body)


def rels_xml(rels):
    return ("<?xml version='1.0' encoding='utf-8' standalone='yes'?><Relationships xmlns='http://schemas.openxmlformats.org/package/2006/relationships'>"
            + "".join("<Relationship Id='%s' Type='%s' Target='%s'/>" % r for r in rels) + "</Relationships>")


def core_xml(d):
    return ("<?xml version='1.0' encoding='utf-8' standalone='yes'?>"
            "<cp:coreProperties xmlns:cp='http://schemas.openxmlformats.org/package/2006/metadata/core-properties' "
            "xmlns:dc='http://purl.org/dc/elements/1.1/' xmlns:dcterms='http://purl.org/dc/terms/' "
            "xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance'>"
            "<dc:title>%s</dc:title><dc:creator>bdf test</dc:creator><dc:language>%s</dc:language>"
            "</cp:coreProperties>" % (escape(d.title), d.language))


def document_xml(d):
    colors = "".join("<ColorEntry IX='%d' RGB='%s'/>" % (i, c) for i, c in sorted(d.colors.items()))
    faces = "".join("<FaceName NameU=%s/>" % quoteattr(n) for i, n in sorted(d.faces.items()))
    styles = ""
    for st in d.styles:
        par = " LineStyle='%d' FillStyle='%d' TextStyle='%d'" % st.parents if st.parents else ""
        styles += "<StyleSheet ID='%d' NameU=%s Name=%s%s>%s</StyleSheet>" % (st.id, quoteattr(st.name), quoteattr(st.name), par, vsdx_cells(st))
    return ("<?xml version='1.0' encoding='utf-8' ?><VisioDocument xmlns='%s' xmlns:r='%s' xml:space='preserve'>"
            "<DocumentSettings TopPage='%d' DefaultTextStyle='0' DefaultLineStyle='0' DefaultFillStyle='0'/>"
            "<Colors>%s</Colors><FaceNames>%s</FaceNames><StyleSheets>%s</StyleSheets></VisioDocument>"
            % (NS, RELNS, d.pages[0].id, colors, faces, styles))


def write_vsdx(d, path):
    files = {}
    overrides = [("/visio/document.xml", "application/vnd.ms-visio.drawing.main+xml"),
                 ("/docProps/core.xml", "application/vnd.openxmlformats-package.core-properties+xml")]
    files["_rels/.rels"] = rels_xml([
        ("rId1", "http://schemas.microsoft.com/visio/2010/relationships/document", "visio/document.xml"),
        ("rId2", "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties", "docProps/core.xml")])
    files["docProps/core.xml"] = core_xml(d)
    files["visio/document.xml"] = document_xml(d)
    doc_rels = [("rId1", "http://schemas.microsoft.com/visio/2010/relationships/masters", "masters/masters.xml"),
                ("rId2", "http://schemas.microsoft.com/visio/2010/relationships/pages", "pages/pages.xml")]
    if d.theme:
        doc_rels.append(("rId3", "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme", "theme/theme1.xml"))
        files["visio/theme/theme1.xml"] = d.theme
        overrides.append(("/visio/theme/theme1.xml", "application/vnd.openxmlformats-officedocument.theme+xml"))
    files["visio/_rels/document.xml.rels"] = rels_xml(doc_rels)
    media_counter = [0]
    media = []
    # masters
    mrels, mlist = [], ""
    for i, m in enumerate(d.masters):
        part = VSDXPart()
        name = "master%d.xml" % (i + 1)
        files["visio/masters/" + name] = ("<?xml version='1.0' encoding='utf-8' ?><MasterContents xmlns='%s' xmlns:r='%s' xml:space='preserve'><Shapes>%s</Shapes></MasterContents>"
                                          % (NS, RELNS, "".join(vsdx_shape(s, part, media_counter) for s in m.shapes)))
        if part.rels:
            files["visio/masters/_rels/%s.rels" % name] = rels_xml(part.rels)
        media += part.media
        rid = "rId%d" % (i + 1)
        mrels.append((rid, "http://schemas.microsoft.com/visio/2010/relationships/master", name))
        mlist += "<Master ID='%d' NameU=%s Name=%s><Rel r:id='%s'/></Master>" % (m.id, quoteattr(m.name), quoteattr(m.name), rid)
        overrides.append(("/visio/masters/" + name, "application/vnd.ms-visio.master+xml"))
    files["visio/masters/masters.xml"] = "<?xml version='1.0' encoding='utf-8' ?><Masters xmlns='%s' xmlns:r='%s' xml:space='preserve'>%s</Masters>" % (NS, RELNS, mlist)
    files["visio/masters/_rels/masters.xml.rels"] = rels_xml(mrels)
    overrides.append(("/visio/masters/masters.xml", "application/vnd.ms-visio.masters+xml"))
    # pages
    prels, plist = [], ""
    for i, pg in enumerate(d.pages):
        part = VSDXPart()
        for m in d.masters:
            part.rel("http://schemas.microsoft.com/visio/2010/relationships/master", "../masters/master%d.xml" % (d.masters.index(m) + 1))
        name = "page%d.xml" % (i + 1)
        files["visio/pages/" + name] = ("<?xml version='1.0' encoding='utf-8' ?><PageContents xmlns='%s' xmlns:r='%s' xml:space='preserve'><Shapes>%s</Shapes></PageContents>"
                                        % (NS, RELNS, "".join(vsdx_shape(s, part, media_counter) for s in pg.shapes)))
        files["visio/pages/_rels/%s.rels" % name] = rels_xml(part.rels)
        media += part.media
        rid = "rId%d" % (i + 1)
        prels.append((rid, "http://schemas.microsoft.com/visio/2010/relationships/page", name))
        attrs = " ID='%d' NameU=%s Name=%s" % (pg.id, quoteattr(pg.name), quoteattr(pg.name))
        if pg.background:
            attrs += " Background='1'"
        if pg.back_page is not None:
            attrs += " BackPage='%d'" % pg.back_page
        plist += "<Page%s><PageSheet LineStyle='0' FillStyle='0' TextStyle='0'>%s</PageSheet><Rel r:id='%s'/></Page>" % (attrs, vsdx_cells(pg.sheet), rid)
        overrides.append(("/visio/pages/" + name, "application/vnd.ms-visio.page+xml"))
    files["visio/pages/pages.xml"] = "<?xml version='1.0' encoding='utf-8' ?><Pages xmlns='%s' xmlns:r='%s' xml:space='preserve'>%s</Pages>" % (NS, RELNS, plist)
    files["visio/pages/_rels/pages.xml.rels"] = rels_xml(prels)
    overrides.append(("/visio/pages/pages.xml", "application/vnd.ms-visio.pages+xml"))
    for name, data in media:
        files["visio/media/" + name] = data
    files["[Content_Types].xml"] = ("<?xml version='1.0' encoding='utf-8' standalone='yes'?><Types xmlns='http://schemas.openxmlformats.org/package/2006/content-types'>"
                                    "<Default Extension='rels' ContentType='application/vnd.openxmlformats-package.relationships+xml'/>"
                                    "<Default Extension='xml' ContentType='application/xml'/><Default Extension='png' ContentType='image/png'/>"
                                    + "".join("<Override PartName='%s' ContentType='%s'/>" % o for o in overrides) + "</Types>")
    with zipfile.ZipFile(path, "w", zipfile.ZIP_DEFLATED) as z:
        for name in sorted(files, key=lambda n: (n != "[Content_Types].xml", n)):
            info = zipfile.ZipInfo(name, date_time=(2026, 1, 1, 0, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            data = files[name]
            z.writestr(info, data if isinstance(data, bytes) else data.encode("utf-8"))


# --- .vdx --------------------------------------------------------------------

# the element that holds each top-level cell in an XML drawing
VDX_GROUPS = {}
for group, names in {
    "XForm": "PinX PinY Width Height LocPinX LocPinY Angle FlipX FlipY ResizeMode",
    "XForm1D": "BeginX BeginY EndX EndY",
    "Line": "LineWeight LineColor LinePattern Rounding EndArrowSize BeginArrow EndArrow LineCap BeginArrowSize LineColorTrans",
    "Fill": "FillForegnd FillBkgnd FillPattern ShdwForegnd ShdwPattern FillForegndTrans FillBkgndTrans ShdwForegndTrans "
            "ShapeShdwType ShapeShdwOffsetX ShapeShdwOffsetY ShapeShdwBlur",
    "TextBlock": "LeftMargin RightMargin TopMargin BottomMargin VerticalAlign TextBkgnd DefaultTabStop TextDirection TextBkgndTrans",
    "TextXForm": "TxtPinX TxtPinY TxtWidth TxtHeight TxtLocPinX TxtLocPinY TxtAngle",
    "Misc": "HideText LangID ObjType",
    "Layout": "ConLineJumpCode ConLineJumpStyle ConLineJumpDirX ConLineJumpDirY",
    "PageLayout": "LineJumpCode LineJumpStyle LineToLineX LineToLineY LineJumpFactorX LineJumpFactorY PageLineJumpDirX PageLineJumpDirY",
    "Group": "DisplayMode",
    "LayerMem": "LayerMember",
    "Foreign": "ImgOffsetX ImgOffsetY ImgWidth ImgHeight",
    "PageProps": "PageWidth PageHeight ShdwOffsetX ShdwOffsetY PageScale DrawingScale",
}.items():
    for n in names.split():
        VDX_GROUPS[n] = group

VDX_SECTIONS = {"Character": "Char", "Paragraph": "Para", "Geometry": "Geom", "Hyperlink": "Hyperlink", "Layer": "Layer",
                "FillGradient": "FillGradient", "Tabs": "Tabs"}


def vdx_value(n, v):
    val, f = value(v)
    return "<%s%s>%s</%s>" % (n, " F=%s" % quoteattr(f) if f else "", escape(val), n)


def vdx_cells(sheet, faces):
    groups = {}
    order = []
    for n, v in sheet.cells:
        g = VDX_GROUPS[n]
        if g not in groups:
            groups[g] = ""
            order.append(g)
        groups[g] += vdx_value(n, v)
    out = "".join("<%s>%s</%s>" % (g, groups[g], g) for g in order)
    for s in sheet.sections:
        el = VDX_SECTIONS[s.name]
        if s.name == "Geometry":
            body = "".join(vdx_value(n, v) for n, v in s.cells)
            for r in s.rows:
                body += "<%s IX='%s'>%s</%s>" % (r.typ, r.key, "".join(vdx_value(n, v) for n, v in r.cells), r.typ)
            out += "<Geom IX='%d'>%s</Geom>" % (s.ix, body)
            continue
        for r in s.rows:
            key = " IX='%s'" % r.key if r.key.isdigit() else " NameU=%s" % quoteattr(r.key)
            cells = ""
            for n, v in r.cells:
                if n in ("Font", "AsianFont") and not isinstance(v, (int, float)):
                    v = faces[v]  # XML drawings refer to fonts by ID
                cells += vdx_value(n, v)
            out += "<%s%s>%s</%s>" % (el, key, cells, el)
    return out


def vdx_shape(s, faces):
    body = vdx_cells(s, faces)
    if s.foreign is not None:
        typ, comp, data = s.foreign
        body += "<ForeignData ForeignType='%s' CompressionType='%s'>%s</ForeignData>" % (typ, comp, base64.b64encode(data).decode())
    if s.text is not None:
        body += text_xml(s.text)
    if s.kids:
        body += "<Shapes>%s</Shapes>" % "".join(vdx_shape(k, faces) for k in s.kids)
    return "<Shape%s>%s</Shape>" % (shape_attrs(s), body)


def write_vdx(d, path):
    faces = {n: i for i, n in d.faces.items()}
    out = ["<?xml version='1.0' encoding='utf-8' ?><VisioDocument xmlns='%s' xml:space='preserve'>" % NS2003,
           "<DocumentProperties><Title>%s</Title><Creator>bdf test</Creator></DocumentProperties>" % escape(d.title),
           "<Colors>%s</Colors>" % "".join("<ColorEntry IX='%d' RGB='%s'/>" % (i, c) for i, c in sorted(d.colors.items())),
           "<FaceNames>%s</FaceNames>" % "".join("<FaceName ID='%d' Name=%s/>" % (i, quoteattr(n)) for i, n in sorted(d.faces.items())),
           "<StyleSheets>"]
    for st in d.styles:
        par = " LineStyle='%d' FillStyle='%d' TextStyle='%d'" % st.parents if st.parents else ""
        out.append("<StyleSheet ID='%d' NameU=%s Name=%s%s>%s</StyleSheet>" % (st.id, quoteattr(st.name), quoteattr(st.name), par, vdx_cells(st, faces)))
    out.append("</StyleSheets><Masters>")
    for m in d.masters:
        out.append("<Master ID='%d' NameU=%s Name=%s><Shapes>%s</Shapes></Master>" % (m.id, quoteattr(m.name), quoteattr(m.name), "".join(vdx_shape(s, faces) for s in m.shapes)))
    out.append("</Masters><Pages>")
    for pg in d.pages:
        attrs = " ID='%d' NameU=%s Name=%s" % (pg.id, quoteattr(pg.name), quoteattr(pg.name))
        if pg.background:
            attrs += " Background='1'"
        if pg.back_page is not None:
            attrs += " BackPage='%d'" % pg.back_page
        out.append("<Page%s><PageSheet LineStyle='0' FillStyle='0' TextStyle='0'>%s</PageSheet><Shapes>%s</Shapes></Page>"
                   % (attrs, vdx_cells(pg.sheet, faces), "".join(vdx_shape(s, faces) for s in pg.shapes)))
    out.append("</Pages></VisioDocument>")
    with open(path, "w", encoding="utf-8") as f:
        f.write("".join(out))


# --- building blocks ---------------------------------------------------------

def rect_geometry(ix=0, **flags):
    return Section("Geometry", ix, flags, [
        Row(1, "RelMoveTo", X=0, Y=0), Row(2, "RelLineTo", X=1, Y=0), Row(3, "RelLineTo", X=1, Y=1),
        Row(4, "RelLineTo", X=0, Y=1), Row(5, "RelLineTo", X=0, Y=0)])


def ellipse_geometry(w, h, ix=0):
    return Section("Geometry", ix, None, [Row(1, "Ellipse", X=w / 2, Y=h / 2, A=w, B=h / 2, C=w / 2, D=h)])


def box(id, x, y, w, h, **cells):
    """A 2-D shape placed by its center."""
    s = Shape(id, PinX=x, PinY=y, Width=w, Height=h, LocPinX=w / 2, LocPinY=h / 2, **cells)
    return s


def line(id, x1, y1, x2, y2, **cells):
    """A 1-D shape from (x1, y1) to (x2, y2)."""
    dx, dy = x2 - x1, y2 - y1
    length = math.hypot(dx, dy)
    s = Shape(id, PinX=(x1 + x2) / 2, PinY=(y1 + y2) / 2, Width=length, Height=0, LocPinX=length / 2, LocPinY=0,
              Angle=math.atan2(dy, dx), BeginX=x1, BeginY=y1, EndX=x2, EndY=y2, **cells)
    s.section(Section("Geometry", 0, {"NoFill": 1}, [Row(1, "MoveTo", X=0, Y=0), Row(2, "LineTo", X=length, Y=0)]))
    return s


def char_row(ix=0, **cells):
    return Row(ix, **cells)


def text(s, *segments, **char):
    s.text = list(segments) if segments else None
    if char:
        s.section(Section("Character", rows=[char_row(0, **char)]))
    return s


def page_sheet(w, h, layers=None):
    sh = Sheet(PageWidth=w, PageHeight=h, ShdwOffsetX=0.06, ShdwOffsetY=-0.06, PageScale=1, DrawingScale=1,
               LineJumpCode=1, LineJumpStyle=0, LineToLineX=0.125, LineToLineY=0.125, LineJumpFactorX=2 / 3, LineJumpFactorY=2 / 3)
    if layers:
        sh.section(Section("Layer", rows=[Row(i, Name=n, Visible=int(v), Print=int(v)) for i, (n, v) in enumerate(layers)]))
    return sh


def root_style(font):
    """The root style sheet ("No Style"): the defaults of every cell."""
    st = Style(0, "No Style", None, LineWeight=0.75 * PT, LineColor=0, LinePattern=1, Rounding=0, EndArrowSize=2,
               BeginArrow=0, EndArrow=0, LineCap=0, BeginArrowSize=2, LineColorTrans=0, FillForegnd=1, FillBkgnd=0,
               FillPattern=1, ShdwForegnd=0, ShdwPattern=0, FillForegndTrans=0, FillBkgndTrans=0, ShdwForegndTrans=0,
               ShapeShdwType=0, ShapeShdwOffsetX=0, ShapeShdwOffsetY=0, ShapeShdwBlur=0, LeftMargin=4 * PT,
               RightMargin=4 * PT, TopMargin=4 * PT, BottomMargin=4 * PT, VerticalAlign=1, TextBkgnd=0,
               DefaultTabStop=0.5, TextDirection=0, TextBkgndTrans=0, HideText=0, LangID="en-US")
    st.section(Section("Character", rows=[Row(0, Font=font, AsianFont=font, Color=0, Style=0, Case=0, Pos=0,
                                              Size=12 * PT, Letterspace=0, LangID="en-US")]))
    st.section(Section("Paragraph", rows=[Row(0, IndFirst=0, IndLeft=0, IndRight=0, SpLine=-1.2, SpBefore=0,
                                              SpAfter=0, HorzAlign=1, Bullet=0)]))
    return st


# --- the dynamic theme of shapes.vsdx -----------------------------------------

def theme_xml():
    def clr(v):
        return "<a:srgbClr val='%s'/>" % v
    ph = "<a:schemeClr val='phClr'/>"
    fills = ("<a:solidFill>%s</a:solidFill>" % clr("FFFFFF") * 2
             + "<a:gradFill rotWithShape='1'><a:gsLst><a:gs pos='0'><a:schemeClr val='phClr'><a:tint val='40000'/></a:schemeClr></a:gs>"
               "<a:gs pos='100000'>%s</a:gs></a:gsLst><a:lin ang='5400000' scaled='0'/></a:gradFill>" % ph
             + "<a:solidFill>%s</a:solidFill>" % ph
             + "<a:solidFill><a:schemeClr val='phClr'><a:shade val='75000'/></a:schemeClr></a:solidFill>"
             + "<a:gradFill rotWithShape='1'><a:gsLst><a:gs pos='0'><a:schemeClr val='phClr'><a:shade val='60000'/></a:schemeClr></a:gs>"
               "<a:gs pos='100000'>%s</a:gs></a:gsLst><a:path path='circle'><a:fillToRect l='50000' t='50000' r='50000' b='50000'/></a:path></a:gradFill>" % ph)
    lines = "".join("<a:ln w='%d' cap='%s'><a:solidFill><a:schemeClr val='phClr'>%s</a:schemeClr></a:solidFill><a:prstDash val='%s'/></a:ln>"
                    % (w, cap, mod, dash) for w, cap, mod, dash in [
                        (9525, "flat", "<a:tint val='40000'/>", "solid"), (9525, "flat", "", "solid"),
                        (12700, "flat", "<a:shade val='75000'/>", "solid"), (19050, "flat", "<a:shade val='50000'/>", "solid"),
                        (12700, "rnd", "<a:shade val='50000'/>", "dash"), (28575, "flat", "<a:shade val='50000'/>", "solid")])
    effects = ("<a:effectStyle><a:effectLst/></a:effectStyle>" * 3
               + "<a:effectStyle><a:effectLst><a:outerShdw blurRad='38100' dist='38100' dir='2700000'><a:schemeClr val='phClr'><a:alpha val='40000'/></a:schemeClr></a:outerShdw></a:effectLst></a:effectStyle>" * 3)
    conn_lines = "".join("<a:ln w='%d' cap='rnd'><a:solidFill><a:schemeClr val='phClr'/></a:solidFill><a:prstDash val='%s'/></a:ln>" % (w, d)
                         for w, d in [(12700, "solid"), (12700, "sysDash"), (25400, "solid")])
    var_colors = "".join("<vt:variationClrScheme%s>%s</vt:variationClrScheme>" % (
        " monotone='1'" if i == 0 else "", "".join("<vt:varColor%d>%s</vt:varColor%d>" % (k + 1, clr(c), k + 1) for k, c in enumerate(cs)))
        for i, cs in enumerate([["2E75B6", "6FA0D0", "1F4E79", "BDD7EE", "ED7D31", "70AD47", "FFC000"],
                                ["C55A11", "70AD47", "7030A0", "FFC000", "2E75B6", "A5A5A5", "5B9BD5"]]))
    var_styles = ("<vt:variationStyleScheme embellishment='1'>"
                  + "".join("<vt:varStyle fillIdx='%d' lineIdx='%d' effectIdx='%d' fontIdx='%d'/>" % (f, l, e, f) for f, l, e in [(4, 3, 1), (3, 2, 1), (5, 4, 4), (6, 6, 4)])
                  + "</vt:variationStyleScheme>")
    line_ex = "".join("<vt:lineStyle><vt:lineEx rndg='%d' start='0' startSize='2' end='%d' endSize='2' pattern='%d'/></vt:lineStyle>" % x
                      for x in [(0, 0, 1), (0, 0, 1), (0, 0, 1), (0, 0, 1), (0, 0, 1), (0, 0, 2), (0, 0, 1)])
    conn_ex = "".join("<vt:lineStyle><vt:lineEx rndg='0' start='0' startSize='2' end='%d' endSize='2' pattern='%d'/></vt:lineStyle>" % x
                      for x in [(0, 1), (4, 1), (4, 9), (13, 1)])
    font_props = "".join("<vt:fontProps style='%d'><vt:color>%s</vt:color></vt:fontProps>" % x for x in [
        (0, "<a:schemeClr val='phClr'/>"), (0, "<a:schemeClr val='phClr'/>"), (0, "<a:schemeClr val='phClr'><a:shade val='50000'/></a:schemeClr>"),
        (1, "<a:schemeClr val='lt1'/>"), (1, "<a:schemeClr val='lt1'/>"), (1, "<a:schemeClr val='lt1'/>")])
    conn_font = "<vt:fontProps style='0'><vt:color><a:schemeClr val='phClr'><a:shade val='50000'/></a:schemeClr></vt:color></vt:fontProps>" * 3
    return ("<?xml version='1.0' encoding='utf-8' standalone='yes'?>"
            "<a:theme xmlns:a='http://schemas.openxmlformats.org/drawingml/2006/main' xmlns:vt='http://schemas.microsoft.com/office/visio/2012/theme' name='Test'>"
            "<a:themeElements><a:clrScheme name='Test'>"
            + "".join("<a:%s>%s</a:%s>" % (n, clr(c), n) for n, c in [
                ("dk1", "000000"), ("lt1", "FFFFFF"), ("dk2", "44546A"), ("lt2", "E7E6E6"), ("accent1", "2E75B6"),
                ("accent2", "ED7D31"), ("accent3", "A5A5A5"), ("accent4", "FFC000"), ("accent5", "4472C4"),
                ("accent6", "70AD47"), ("hlink", "0563C1"), ("folHlink", "954F72")])
            + "<a:extLst><a:ext uri='{093E89EA-6996-430E-BFF9-83A9FAAAAB73}'><vt:bkgnd>%s</vt:bkgnd></a:ext>" % clr("FFFFFF")
            + "<a:ext uri='{DDD2D869-C2EF-471E-B8FA-914AFA308C9F}'><vt:variationClrSchemeLst>%s</vt:variationClrSchemeLst></a:ext></a:extLst>" % var_colors
            + "</a:clrScheme><a:fontScheme name='Test'><a:majorFont><a:latin typeface='%s'/><a:ea typeface='%s'/><a:cs typeface=''/></a:majorFont>"
              "<a:minorFont><a:latin typeface='%s'/><a:ea typeface='%s'/><a:cs typeface=''/></a:minorFont></a:fontScheme>" % (FONT, FONT, FONT, FONT)
            + "<a:fmtScheme name='Test'><a:fillStyleLst>%s</a:fillStyleLst><a:lnStyleLst>%s</a:lnStyleLst><a:effectStyleLst>%s</a:effectStyleLst>"
              "<a:bgFillStyleLst><a:solidFill>%s</a:solidFill></a:bgFillStyleLst></a:fmtScheme>" % (fills, lines, effects, ph)
            + "<a:extLst><a:ext uri='{1342405F-259F-4C95-8CDF-A9DCE2D418A2}'><vt:fmtConnectorScheme name='Test'>"
              "<a:fillStyleLst>%s</a:fillStyleLst><a:lnStyleLst>%s</a:lnStyleLst><a:effectStyleLst>%s</a:effectStyleLst>"
              "<a:bgFillStyleLst><a:solidFill>%s</a:solidFill></a:bgFillStyleLst></vt:fmtConnectorScheme></a:ext>"
              % ("<a:solidFill>%s</a:solidFill>" % ph * 3, conn_lines, "<a:effectStyle><a:effectLst/></a:effectStyle>" * 3, ph)
            + "<a:ext uri='{56243398-1771-4C39-BF73-A5702A9C147F}'><vt:fillStyles>%s</vt:fillStyles></a:ext>" % ("<vt:fillProps pattern='1'/>" * 6)
            + "<a:ext uri='{6CAB99AB-0A78-4BAB-B597-526D62367CB4}'><vt:lineStyles><vt:fmtConnectorSchemeLineStyles>%s</vt:fmtConnectorSchemeLineStyles>"
              "<vt:fmtSchemeLineStyles>%s</vt:fmtSchemeLineStyles></vt:lineStyles></a:ext>" % (conn_ex, line_ex)
            + "<a:ext uri='{EBE24D50-EC5C-4D6F-A1A3-C5F0A18B936A}'><vt:fontStylesGroup><vt:connectorFontStyles>%s</vt:connectorFontStyles>"
              "<vt:fontStyles>%s</vt:fontStyles></vt:fontStylesGroup></a:ext>" % (conn_font, font_props)
            + "<a:ext uri='{494CE47F-D151-47DC-95E8-85652EA8A67E}'><vt:variationStyleSchemeLst>%s%s</vt:variationStyleSchemeLst></a:ext>" % (var_styles, var_styles)
            + "</a:extLst></a:themeElements></a:theme>")


THEMED = ["LineWeight", "LineColor", "LinePattern", "LineCap", "LineColorTrans", "Rounding", "FillForegnd", "FillBkgnd",
          "FillPattern", "FillForegndTrans", "FillBkgndTrans", "ShdwForegnd", "ShdwPattern", "ShdwForegndTrans",
          "ShapeShdwType", "ShapeShdwOffsetX", "ShapeShdwOffsetY", "ShapeShdwBlur"]


def themed_styles():
    """No Style, Theme (THEMEVAL), Normal and Connector, as Visio 2013 writes them."""
    root = root_style(FONT)
    theme = Style(1, "Theme", (0, 0, 0), **{n: ("Themed", "THEMEVAL()") for n in THEMED})
    theme.set(QuickStyleLineColor=100, QuickStyleFillColor=100, QuickStyleShadowColor=100, QuickStyleFontColor=100,
              QuickStyleLineMatrix=100, QuickStyleFillMatrix=100, QuickStyleEffectsMatrix=100, QuickStyleFontMatrix=100,
              FillGradientEnabled=("Themed", "THEMEVAL()"), FillGradientDir=("Themed", "THEMEVAL()"),
              FillGradientAngle=("Themed", "THEMEVAL()"))
    theme.section(Section("Character", rows=[Row(0, Font=("Themed", "THEMEVAL()"), Color=("Themed", "THEMEVAL()"))]))
    theme.section(Section("FillGradient", rows=[Row(i, GradientStopColor=("Themed", "THEMEVAL()"),
                                                    GradientStopColorTrans=("Themed", "THEMEVAL()"),
                                                    GradientStopPosition=("Themed", "THEMEVAL()")) for i in range(2)]))
    normal = Style(2, "Normal", (1, 1, 1))
    connector = Style(3, "Connector", (2, 2, 2), EndArrow=("Themed", "THEMEVAL()"), EndArrowSize=("Themed", "THEMEVAL()"),
                      TextBkgnd="#FFFFFF", VerticalAlign=1)
    return [root, theme, normal, connector]


# --- shapes.vsdx ----------------------------------------------------------------

def shapes_drawing():
    W, H = 11.0, 8.5
    rect_master = Master(2, "Rectangle", [rect_master_shape(1)])
    circle = Shape(1, styles=(2, 2, 2), PinX=0.5, PinY=0.5, Width=1, Height=1, LocPinX=0.5, LocPinY=0.5)
    circle.section(Section("Geometry", 0, None, [Row(1, "Ellipse", X=("0.5", "Width*0.5"), Y=("0.5", "Height*0.5"),
                                                    A=("1", "Width"), B=("0.5", "Height*0.5"), C=("0.5", "Width*0.5"), D=("1", "Height"))]))
    circle_master = Master(3, "Circle", [circle])
    conn = Shape(1, styles=(3, 3, 3), PinX=0.5, PinY=0, Width=1, Height=0, LocPinX=0.5, LocPinY=0, BeginX=0, BeginY=0, EndX=1, EndY=0,
                 ObjType=2)
    conn.section(Section("Geometry", 0, {"NoFill": 1}, [Row(1, "MoveTo", X=0, Y=0), Row(2, "LineTo", X=("1", "Width"), Y=0)]))
    conn_master = Master(4, "Dynamic connector", [conn])

    ids = iter(range(10, 10000))
    shapes = []

    def label(x, y, s, size=9, w=1.6, **char):
        t = box(next(ids), x, y, w, 0.3, LinePattern=0, FillPattern=0)
        t.section(rect_geometry(NoFill=1, NoLine=1))
        return text(t, s, Size=size * PT, **char)

    # 1: variant quick styles and theme colors
    shapes.append(label(1.6, 8.1, "Theme quick styles", 12, 3, Style=1))
    for i, qs in enumerate([100, 101, 102, 103]):
        s = Shape(next(ids), master=2, PinX=0.8 + i * 0.9, PinY=7.55, Width=0.7, Height=0.5, LocPinX=0.35, LocPinY=0.25)
        s.set(**{n: qs for n in ["QuickStyleLineColor", "QuickStyleFillColor", "QuickStyleShadowColor", "QuickStyleFontColor",
                                 "QuickStyleLineMatrix", "QuickStyleFillMatrix", "QuickStyleEffectsMatrix", "QuickStyleFontMatrix"]})
        shapes.append(text(s, "V%d" % (qs - 99)))
    for i, qs in enumerate([2, 3, 4, 5, 6, 7, 204]):
        s = Shape(next(ids), master=3, PinX=0.6 + i * 0.5, PinY=6.85, Width=0.4, Height=0.4, LocPinX=0.2, LocPinY=0.2,
                  QuickStyleFillColor=qs, QuickStyleLineColor=qs, QuickStyleFillMatrix=4, QuickStyleLineMatrix=3)
        # like Visio, an instance of another size holds the geometry of its size
        s.section(Section("Geometry", 0, None, [Row(1, "Ellipse", X=(0.2, "Inh"), Y=(0.2, "Inh"), A=(0.4, "Inh"),
                                                    B=(0.2, "Inh"), C=(0.2, "Inh"), D=(0.4, "Inh"))]))
        shapes.append(s)

    # 2: geometry rows
    shapes.append(label(1.6, 6.35, "Geometry rows", 12, 3, Style=1))
    x0, y0 = 0.3, 5.3
    geoms = []
    arc = Shape(next(ids), styles=(2, 2, 2), PinX=x0 + 0.4, PinY=y0 + 0.4, Width=0.8, Height=0.8, LocPinX=0.4, LocPinY=0.4)
    arc.section(Section("Geometry", 0, None, [Row(1, "MoveTo", X=0.2, Y=0), Row(2, "LineTo", X=0.6, Y=0), Row(3, "ArcTo", X=0.8, Y=0.2, A=0.0586),
                                              Row(4, "LineTo", X=0.8, Y=0.6), Row(5, "ArcTo", X=0.6, Y=0.8, A=0.0586), Row(6, "LineTo", X=0.2, Y=0.8),
                                              Row(7, "ArcTo", X=0, Y=0.6, A=0.0586), Row(8, "LineTo", X=0, Y=0.2), Row(9, "ArcTo", X=0.2, Y=0, A=0.0586)]))
    geoms.append(arc)
    pie = Shape(next(ids), styles=(2, 2, 2), PinX=x0 + 1.4, PinY=y0 + 0.4, Width=0.8, Height=0.8, LocPinX=0.4, LocPinY=0.4)
    pie.section(Section("Geometry", 0, None, [Row(1, "MoveTo", X=0.4, Y=0.4), Row(2, "LineTo", X=0.8, Y=0.4),
                                              Row(3, "EllipticalArcTo", X=0.4, Y=0.8, A=0.6828, B=0.6828, C=0, D=1), Row(4, "LineTo", X=0.4, Y=0.4)]))
    geoms.append(pie)
    half = Shape(next(ids), styles=(2, 2, 2), PinX=x0 + 2.4, PinY=y0 + 0.4, Width=0.9, Height=0.6, LocPinX=0.45, LocPinY=0.3)
    half.section(Section("Geometry", 0, None, [Row(1, "MoveTo", X=0, Y=0.3), Row(2, "EllipticalArcTo", X=0.9, Y=0.3, A=0.45, B=0.6, C=0, D=1.5),
                                               Row(3, "LineTo", X=0, Y=0.3)]))
    geoms.append(half)
    star = Shape(next(ids), styles=(2, 2, 2), PinX=x0 + 3.4, PinY=y0 + 0.4, Width=0.8, Height=0.8, LocPinX=0.4, LocPinY=0.4)
    pts = [(0.5 + 0.5 * math.sin(i * math.pi / 5) * (1 if i % 2 == 0 else 0.4), 0.5 + 0.5 * math.cos(i * math.pi / 5) * (1 if i % 2 == 0 else 0.4)) for i in range(10)]
    star.section(Section("Geometry", 0, None, [Row(1, "RelMoveTo", X=pts[0][0], Y=pts[0][1]),
                                               Row(2, "PolylineTo", X=("0.4", "Width*0.5"), Y=("0.8", "Height*1"),
                                                   A="POLYLINE(0, 0, " + ", ".join("%.4f,%.4f" % p for p in pts[1:]) + ")")]))
    # the POLYLINE's points are relative to the box; the end is the start
    geoms.append(star)
    nurbs = Shape(next(ids), styles=(2, 2, 2), PinX=x0 + 4.4, PinY=y0 + 0.4, Width=0.8, Height=0.8, LocPinX=0.4, LocPinY=0.4, FillPattern=0)
    nurbs.section(Section("Geometry", 0, {"NoFill": 1}, [Row(1, "MoveTo", X=0, Y=0),
                                                        Row(2, "NURBSTo", X=0.8, Y=0, A=1, B=1, C=0, D=1,
                                                            E="NURBS(1, 3, 0, 0, 0,1,0,1, 0.5,0,0.5,1, 1,1,1,1)")]))
    geoms.append(nurbs)
    spline = Shape(next(ids), styles=(2, 2, 2), PinX=x0 + 5.4, PinY=y0 + 0.4, Width=0.8, Height=0.8, LocPinX=0.4, LocPinY=0.4, FillPattern=0)
    spline.section(Section("Geometry", 0, {"NoFill": 1}, [Row(1, "MoveTo", X=0, Y=0.4), Row(2, "SplineStart", X=0.2, Y=0.8, A=0, B=0, C=1, D=3),
                                                         Row(3, "SplineKnot", X=0.6, Y=0, A=0.5), Row(4, "SplineKnot", X=0.8, Y=0.4, A=1)]))
    geoms.append(spline)
    bez = Shape(next(ids), styles=(2, 2, 2), PinX=x0 + 6.4, PinY=y0 + 0.4, Width=0.8, Height=0.8, LocPinX=0.4, LocPinY=0.4)
    bez.section(Section("Geometry", 0, None, [Row(1, "RelMoveTo", X=0, Y=0), Row(2, "RelCubBezTo", X=1, Y=0, A=0.2, B=1, C=0.8, D=1),
                                              Row(3, "RelQuadBezTo", X=0, Y=0, A=0.5, B=-0.5)]))
    geoms.append(bez)
    donut = Shape(next(ids), styles=(2, 2, 2), PinX=x0 + 7.4, PinY=y0 + 0.4, Width=0.8, Height=0.8, LocPinX=0.4, LocPinY=0.4)
    donut.section(Section("Geometry", 0, None, [Row(1, "Ellipse", X=0.4, Y=0.4, A=0.8, B=0.4, C=0.4, D=0.8),
                                                Row(2, "Ellipse", X=0.4, Y=0.4, A=0.6, B=0.4, C=0.4, D=0.6)]))
    geoms.append(donut)
    shapes += geoms

    # 3: fills
    shapes.append(label(1.6, 4.8, "Fills", 12, 3, Style=1))
    fills = [dict(FillForegnd="#2E75B6"), dict(FillForegnd="#2E75B6", FillForegndTrans=0.6)]
    fills += [dict(FillPattern=p, FillForegnd="#1F4E79", FillBkgnd="#DDEBF7") for p in (2, 4, 9, 13, 17, 23)]
    fills += [dict(FillPattern=25, FillForegnd="#1F4E79", FillBkgnd="#FFFFFF"), dict(FillPattern=40, FillForegnd="#C55A11", FillBkgnd="#FFFFFF")]
    for i, f in enumerate(fills):
        s = box(next(ids), 0.5 + i * 0.65, 4.3, 0.5, 0.5, styles=(0, 0, 0), **f)
        shapes.append(s.section(rect_geometry()))
    grad = box(next(ids), 7.3, 4.3, 0.9, 0.5, styles=(0, 0, 0), FillGradientEnabled=1, FillGradientDir=0, FillGradientAngle=math.pi / 4)
    grad.section(Section("FillGradient", rows=[Row(0, GradientStopColor="#FFC000", GradientStopColorTrans=0, GradientStopPosition=0),
                                               Row(1, GradientStopColor="#C00000", GradientStopColorTrans=0, GradientStopPosition=1)]))
    shapes.append(grad.section(rect_geometry()))
    radial = box(next(ids), 8.4, 4.3, 0.9, 0.5, styles=(0, 0, 0), FillGradientEnabled=1, FillGradientDir=3)
    radial.section(Section("FillGradient", rows=[Row(0, GradientStopColor="#FFFFFF", GradientStopColorTrans=0, GradientStopPosition=0),
                                                 Row(1, GradientStopColor="#70AD47", GradientStopColorTrans=0, GradientStopPosition=1)]))
    shapes.append(radial.section(ellipse_geometry(0.9, 0.5)))

    # 4: lines, rounding, rotation, flips, groups, shadows
    shapes.append(label(1.6, 3.8, "Lines", 12, 3, Style=1))
    for i, p in enumerate(range(2, 24)):
        shapes.append(line(next(ids), 0.3 + (i % 11) * 0.45, 3.4 - (i // 11) * 0.25, 0.65 + (i % 11) * 0.45, 3.4 - (i // 11) * 0.25,
                           styles=(0, 0, 0), LinePattern=p, LineWeight=1.5 * PT, LineColor="#1F4E79"))
    for i, (w, cap) in enumerate([(0.25, 0), (1, 1), (3, 2), (6, 0)]):
        shapes.append(line(next(ids), 0.3, 2.7 - i * 0.2, 1.8, 2.7 - i * 0.2, styles=(0, 0, 0), LineWeight=w * PT, LineCap=cap))
    rounded = box(next(ids), 2.6, 2.4, 1.0, 0.6, styles=(2, 2, 2), Rounding=0.15)
    shapes.append(rounded.section(rect_geometry()))
    rot = box(next(ids), 3.9, 2.4, 1.0, 0.4, styles=(2, 2, 2), Angle=math.pi / 6)
    shapes.append(text(rot.section(rect_geometry()), "Turned"))
    flip = Shape(next(ids), styles=(2, 2, 2), PinX=5.1, PinY=2.4, Width=0.8, Height=0.6, LocPinX=0.4, LocPinY=0.3, FlipX=1)
    flip.section(Section("Geometry", 0, None, [Row(1, "MoveTo", X=0, Y=0), Row(2, "LineTo", X=0.8, Y=0), Row(3, "LineTo", X=0, Y=0.6), Row(4, "LineTo", X=0, Y=0)]))
    shapes.append(text(flip, "FlipX"))
    flipy = Shape(next(ids), styles=(2, 2, 2), PinX=6.1, PinY=2.4, Width=0.8, Height=0.6, LocPinX=0.4, LocPinY=0.3, FlipY=1)
    flipy.section(Section("Geometry", 0, None, [Row(1, "MoveTo", X=0, Y=0), Row(2, "LineTo", X=0.8, Y=0), Row(3, "LineTo", X=0, Y=0.6), Row(4, "LineTo", X=0, Y=0)]))
    shapes.append(text(flipy, "FlipY"))
    group = Shape(next(ids), "Group", PinX=7.4, PinY=2.4, Width=1.2, Height=0.8, LocPinX=0.6, LocPinY=0.4, Angle=-math.pi / 12)
    g1 = box(next(ids), 0.3, 0.4, 0.5, 0.7, styles=(0, 0, 0), FillForegnd="#FFC000")
    g2 = Shape(next(ids), styles=(0, 0, 0), PinX=0.9, PinY=0.4, Width=0.5, Height=0.5, LocPinX=0.25, LocPinY=0.25, FillForegnd="#70AD47")
    g2.section(ellipse_geometry(0.5, 0.5))
    group.kids = [g1.section(rect_geometry()), g2]
    shapes.append(group)
    shadow = box(next(ids), 8.8, 2.4, 0.8, 0.6, styles=(0, 0, 0), FillForegnd="#FFFFFF", ShdwPattern=1, ShdwForegnd="#7F7F7F",
                 ShapeShdwType=1, ShapeShdwOffsetX=0.06, ShapeShdwOffsetY=-0.06)
    shapes.append(text(shadow.section(rect_geometry()), "Shadow"))
    hidden = box(next(ids), 9.8, 2.4, 0.6, 0.6, styles=(0, 0, 0), FillForegnd="#FF0000", LayerMember="1")
    shapes.append(text(hidden.section(rect_geometry()), "Hidden layer"))
    picture = Shape(next(ids), "Foreign", PinX=9.9, PinY=4.3, Width=0.8, Height=0.6, LocPinX=0.4, LocPinY=0.3,
                    ImgOffsetX=0, ImgOffsetY=0, ImgWidth=0.8, ImgHeight=0.6)
    picture.foreign = ("Bitmap", "PNG", png(8, 6, lambda x, y: ((x * 31) % 256, (y * 42) % 256, 160) if (x + y) % 2 else (255, 255, 255)))
    shapes.append(picture)
    link = box(next(ids), 9.9, 3.3, 1.0, 0.4, styles=(2, 2, 2))
    link.section(Section("Hyperlink", rows=[Row("Row_1", Address="", SubAddress="Text")]))
    shapes.append(text(link.section(rect_geometry()), "To page 2"))

    # 5: arrowheads
    shapes.append(label(1.6, 1.75, "Arrowheads", 12, 3, Style=1))
    for k in range(1, 46):
        col, r = (k - 1) % 15, (k - 1) // 15
        x = 0.3 + col * 0.7
        y = 1.3 - r * 0.4
        shapes.append(line(next(ids), x, y, x + 0.55, y, styles=(0, 0, 0), EndArrow=k, BeginArrow=(k % 3) * 4, LineWeight=1 * PT,
                           LineColor="#000000", EndArrowSize=2, BeginArrowSize=1))

    background = Page(3, "Background", page_sheet(W, H), [
        text(box(90, W / 2, H - 0.2, W, 0.4, styles=(0, 0, 0), FillForegnd="#DEEBF7", LinePattern=0).section(rect_geometry()),
             ("cp", 0), "BDF from Visio: ", ("fld", 0, "shapes.vsdx"), Size=10 * PT, Color="#1F4E79")], background=True)
    pages = [Page(0, "Shapes", page_sheet(W, H, [("Visible", True), ("Hidden", False)]), shapes, back_page=3),
             Page(1, "Text", page_sheet(W, H), text_page_shapes(), back_page=3), background]
    return Drawing("Visio shapes", themed_styles(), [rect_master, circle_master, conn_master], pages, theme=theme_xml(),
                   faces={1: FONT})


def rect_master_shape(id):
    s = Shape(id, styles=(2, 2, 2), PinX=0.5, PinY=0.375, Width=1, Height=0.75, LocPinX=("0.5", "Width*0.5"), LocPinY=("0.375", "Height*0.5"))
    return s.section(rect_geometry())


def text_page_shapes():
    ids = iter(range(500, 10000))
    shapes = []
    kinds = [(0, "Left aligned text wraps at the edge of its block."), (1, "Centered text wraps at the edge of its block."),
             (2, "Right aligned text wraps at the edge of its block."), (3, "Justified text spreads the words of lines that wrap.")]
    for i, (a, t) in enumerate(kinds):
        s = box(next(ids), 1.2 + i * 2.2, 7.4, 2.0, 0.9, styles=(0, 0, 0), FillForegnd="#F2F2F2", LineColor="#A6A6A6", VerticalAlign=0)
        s.section(rect_geometry())
        s.section(Section("Paragraph", rows=[Row(0, HorzAlign=a)]))
        shapes.append(text(s, t))
    # character formats
    s = box(next(ids), 3.0, 6.3, 5.6, 0.9, styles=(0, 0, 0), FillForegnd="#FFFFFF", LineColor="#A6A6A6")
    s.section(rect_geometry())
    s.section(Section("Character", rows=[Row(0), Row(1, Style=1), Row(2, Style=2), Row(3, Style=4), Row(4, Strikethru=1),
                                         Row(5, Pos=1, Size=12 * PT), Row(6, Pos=2), Row(7, Color="#C00000", Size=16 * PT),
                                         Row(8, Case=1), Row(9, Style=8), Row(10, Letterspace=2 * PT)]))
    shapes.append(text(s, ("cp", 0), "Plain ", ("cp", 1), "bold ", ("cp", 2), "italic ", ("cp", 3), "underline ",
                       ("cp", 4), "strike ", ("cp", 0), "x", ("cp", 5), "2", ("cp", 0), " H", ("cp", 6), "2", ("cp", 0), "O ",
                       ("cp", 7), "red 16pt ", ("cp", 8), "all caps ", ("cp", 9), "Small Caps ", ("cp", 10), "spaced"))
    # bullets and paragraph spacing
    s = box(next(ids), 8.6, 6.3, 3.4, 0.9, styles=(0, 0, 0), FillForegnd="#FFFFFF", LineColor="#A6A6A6", VerticalAlign=0)
    s.section(rect_geometry())
    s.section(Section("Paragraph", rows=[Row(0, HorzAlign=0, Bullet=1, IndLeft=0.1, SpAfter=3 * PT),
                                         Row(1, HorzAlign=0, Bullet=4, IndLeft=0.35, SpLine=-1.0)]))
    shapes.append(text(s, ("pp", 0), "First bullet\nSecond bullet\n", ("pp", 1), "Nested item, tighter lines\n"))
    # Japanese line breaking and vertical text
    s = box(next(ids), 1.6, 4.9, 2.8, 1.2, styles=(0, 0, 0), FillForegnd="#FFFFFF", LineColor="#A6A6A6")
    s.section(rect_geometry())
    s.section(Section("Paragraph", rows=[Row(0, HorzAlign=0)]))
    shapes.append(text(s, "日本語の文章は、単語の区切りに空白を使わないため、文字と文字の間で改行します。句読点「、」や「。」は行頭に来ません。",
                       LangID="ja-JP"))
    s = box(next(ids), 3.7, 4.9, 0.6, 1.6, styles=(0, 0, 0), FillForegnd="#FFFFFF", LineColor="#A6A6A6", TextDirection=1)
    s.section(rect_geometry())
    shapes.append(text(s, "縦書きのテキスト", LangID="ja-JP"))
    # rotated text block, text background, tabs, line break
    s = box(next(ids), 5.4, 4.9, 1.8, 1.0, styles=(0, 0, 0), FillForegnd="#FFFFFF", LineColor="#A6A6A6", TxtAngle=math.pi / 2,
            TxtWidth=1.0, TxtHeight=1.8, TxtLocPinX=0.5, TxtLocPinY=0.9)
    shapes.append(text(s.section(rect_geometry()), "Text turned by TxtAngle"))
    s = box(next(ids), 7.4, 4.9, 1.8, 1.0, styles=(0, 0, 0), FillForegnd="#FCE4D6", LinePattern=0, TextBkgnd="#FFFFFF")
    shapes.append(text(s.section(rect_geometry()), "Text on a white background"))
    s = box(next(ids), 9.6, 4.9, 2.5, 1.0, styles=(0, 0, 0), FillForegnd="#FFFFFF", LineColor="#A6A6A6", VerticalAlign=0)
    s.section(rect_geometry())
    s.section(Section("Paragraph", rows=[Row(0, HorzAlign=0)]))
    s.section(Section("Tabs", rows=[Row(0, Position1=0.8, Alignment1=0, Position2=2.0, Alignment2=2)]))
    shapes.append(text(s, "Name\tValue\t9\nTabs and a soft break"))
    # a themed shape and a connector with a label between two masters
    a = Shape(next(ids), master=2, PinX=1.5, PinY=2.2, Width=1.4, Height=0.8, LocPinX=0.7, LocPinY=0.4)
    shapes.append(text(a, "Start"))
    b = Shape(next(ids), master=3, PinX=5.5, PinY=2.2, Width=1.0, Height=1.0, LocPinX=0.5, LocPinY=0.5)
    shapes.append(text(b, "End"))
    c = Shape(next(ids), master=4, PinX=3.6, PinY=2.2, Width=2.3, Height=0, LocPinX=1.15, LocPinY=0, BeginX=2.2, BeginY=2.2, EndX=5.0, EndY=2.2)
    c.section(Section("Geometry", 0, None, [Row(1, "MoveTo", X=0, Y=0), Row(2, "LineTo", X=2.8, Y=0)]))
    c.set(Width=2.8, LocPinX=1.4)
    shapes.append(text(c, "connector label"))
    link = box(next(ids), 8.5, 2.2, 2.2, 0.6, styles=(2, 2, 2))
    link.section(Section("Hyperlink", rows=[Row("Row_1", Address="https://example.com/", SubAddress="")]))
    shapes.append(text(link.section(rect_geometry()), "https://example.com/"))
    # line jumps where connectors cross: the horizontal ones jump, in the
    # page's style (an arc) or their own
    t = box(next(ids), 1.5, 1.0, 2.0, 0.3, LinePattern=0, FillPattern=0)
    t.section(rect_geometry(NoFill=1, NoLine=1))
    shapes.append(text(t, "Line jumps", Size=12 * PT, Style=1))

    def connector(x1, y1, x2, y2, **cells):
        c = Shape(next(ids), master=4, PinX=x1, PinY=y1, Width=max(abs(x2 - x1), 0.01), Height=max(abs(y2 - y1), 0.01),
                  LocPinX=0, LocPinY=0, BeginX=x1, BeginY=y1, EndX=x2, EndY=y2, **cells)
        c.section(Section("Geometry", 0, {"NoFill": 1}, [Row(1, "MoveTo", X=0, Y=0), Row(2, "LineTo", X=x2 - x1, Y=y2 - y1)]))
        return c
    for i, (style, down) in enumerate([(0, 0), (2, 0), (3, 0), (5, 2)]):
        y = 1.45 - i * 0.3
        cells = {"ConLineJumpStyle": style}
        if down:
            cells["ConLineJumpDirX"] = down
        shapes.append(connector(3.0, y, 6.2, y, **cells))
    for x in (4.0, 5.2):
        shapes.append(connector(x, 1.7, x, 0.2))
    # a connector that never jumps
    shapes.append(connector(6.8, 1.0, 8.6, 1.0, ConLineJumpCode=1))
    shapes.append(connector(7.7, 1.7, 7.7, 0.2))
    return shapes


# --- flow.vsdx / flow.vdx ----------------------------------------------------------

def flow_drawing():
    """A flowchart without a theme, written in both formats."""
    W, H = 8.5, 11.0
    root = root_style(FONT)
    flow = Style(1, "Flow Normal", (0, 0, 0), LineColor=10, LineWeight=1 * PT, FillForegnd="#DEEBF7", FillPattern=1)
    flow.section(Section("Character", rows=[Row(0, Color="#1F4E79", Size=11 * PT)]))
    conn = Style(2, "Connector", (0, 0, 0), EndArrow=13, EndArrowSize=2, LineColor="#595959", TextBkgnd="#FFFFFF")
    process = Shape(1, styles=(1, 1, 1), PinX=0.75, PinY=0.375, Width=1.5, Height=0.75, LocPinX=0.75, LocPinY=0.375)
    process.section(rect_geometry())
    decision = Shape(1, styles=(1, 1, 1), PinX=0.75, PinY=0.5, Width=1.5, Height=1.0, LocPinX=0.75, LocPinY=0.5, FillForegnd="#FFF2CC")
    decision.section(Section("Geometry", 0, None, [Row(1, "RelMoveTo", X=0.5, Y=0), Row(2, "RelLineTo", X=1, Y=0.5),
                                                   Row(3, "RelLineTo", X=0.5, Y=1), Row(4, "RelLineTo", X=0, Y=0.5), Row(5, "RelLineTo", X=0.5, Y=0)]))
    terminator = Shape(1, styles=(1, 1, 1), PinX=0.75, PinY=0.25, Width=1.5, Height=0.5, LocPinX=0.75, LocPinY=0.25, FillForegnd="#E2EFDA")
    terminator.section(Section("Geometry", 0, None, [Row(1, "MoveTo", X=0.25, Y=0), Row(2, "LineTo", X=1.25, Y=0),
                                                     Row(3, "EllipticalArcTo", X=1.25, Y=0.5, A=1.5, B=0.25, C=0, D=1),
                                                     Row(4, "LineTo", X=0.25, Y=0.5), Row(5, "EllipticalArcTo", X=0.25, Y=0, A=0, B=0.25, C=0, D=1)]))
    connector = Shape(1, styles=(2, 2, 2), PinX=0, PinY=0, Width=1, Height=0, LocPinX=0.5, LocPinY=0, ObjType=2)
    connector.section(Section("Geometry", 0, {"NoFill": 1}, [Row(1, "MoveTo", X=0, Y=0), Row(2, "LineTo", X=1, Y=0)]))
    masters = [Master(1, "Process", [process]), Master(2, "Decision", [decision]), Master(3, "Terminator", [terminator]),
               Master(4, "Dynamic connector", [connector])]

    ids = iter(range(10, 1000))
    shapes = []

    def node(master, x, y, label, **cells):
        s = Shape(next(ids), master=master, PinX=x, PinY=y, **cells)
        return text(s, label)

    def arrow(points, label=None):
        (x1, y1), (x2, y2) = points[0], points[-1]
        s = Shape(next(ids), master=4, PinX=(x1 + x2) / 2, PinY=(y1 + y2) / 2, Width=1, Height=0, LocPinX=0.5, LocPinY=0,
                  BeginX=x1, BeginY=y1, EndX=x2, EndY=y2)
        # local coordinates start at the begin point
        s.set(PinX=x1, PinY=y1, LocPinX=0, LocPinY=0, Width=max(abs(x2 - x1), 0.01), Height=max(abs(y2 - y1), 0.01))
        rows = [Row(1, "MoveTo", X=0, Y=0)] + [Row(i + 2, "LineTo", X=px - x1, Y=py - y1) for i, (px, py) in enumerate(points[1:])]
        s.section(Section("Geometry", 0, {"NoFill": 1}, rows))
        if label:
            mx = (points[len(points) // 2 - 1][0] + points[len(points) // 2][0]) / 2 - x1
            my = (points[len(points) // 2 - 1][1] + points[len(points) // 2][1]) / 2 - y1
            s.set(TxtPinX=mx, TxtPinY=my, TxtWidth=0.5, TxtHeight=0.25, TxtLocPinX=0.25, TxtLocPinY=0.125)
            text(s, label)
        return s

    shapes.append(node(3, 4.25, 10.0, "Start"))
    shapes.append(node(1, 4.25, 8.9, "Collect input"))
    shapes.append(node(2, 4.25, 7.55, "Valid?"))
    shapes.append(node(1, 4.25, 5.95, "Process data"))
    shapes.append(node(1, 6.95, 7.55, "Report error", FillForegnd="#FCE4D6"))
    shapes.append(node(3, 4.25, 4.75, "End"))
    shapes.append(arrow([(4.25, 9.75), (4.25, 9.275)]))
    shapes.append(arrow([(4.25, 8.525), (4.25, 8.05)]))
    shapes.append(arrow([(4.25, 7.05), (4.25, 6.325)], "Yes"))
    shapes.append(arrow([(5.0, 7.55), (6.2, 7.55)], "No"))
    shapes.append(arrow([(6.95, 7.925), (6.95, 9.45), (3.0, 9.45), (3.0, 8.9), (3.5, 8.9)]))
    shapes.append(arrow([(4.25, 5.575), (4.25, 5.0)]))
    logo = Shape(next(ids), "Foreign", PinX=7.4, PinY=10.2, Width=0.6, Height=0.6, LocPinX=0.3, LocPinY=0.3,
                 ImgOffsetX=0, ImgOffsetY=0, ImgWidth=0.6, ImgHeight=0.6)
    logo.foreign = ("Bitmap", "PNG", png(4, 4, lambda x, y: (46, 117, 182) if (x + y) % 2 else (255, 192, 0)))
    shapes.append(logo)
    note = box(next(ids), 1.4, 7.6, 1.8, 1.2, styles=(0, 0, 0), FillForegnd="#FFFFFF", LineColor="#A6A6A6", LinePattern=2, LayerMember="0")
    note.section(rect_geometry())
    note.section(Section("Paragraph", rows=[Row(0, HorzAlign=0)]))
    shapes.append(text(note, "Notes live on a layer; the hidden one is not drawn."))
    hidden = box(next(ids), 1.4, 6.0, 1.2, 0.6, styles=(0, 0, 0), FillForegnd="#FF0000", LayerMember="1")
    shapes.append(text(hidden.section(rect_geometry()), "Hidden layer"))
    background = Page(1, "Background", page_sheet(W, H), [
        text(box(90, W / 2, 0.3, W, 0.6, styles=(0, 0, 0), FillForegnd="#F2F2F2", LinePattern=0).section(rect_geometry()),
             "Flowchart test drawing", Size=9 * PT, Color="#7F7F7F")], background=True)
    pages = [Page(0, "Flow", page_sheet(W, H, [("Notes", True), ("Hidden", False)]), shapes, back_page=1), background]
    return Drawing("Visio flowchart", [root, flow, conn], masters, pages, colors={24: "#1F4E79"}, faces={1: FONT})


def main():
    os.makedirs(OUT, exist_ok=True)
    write_vsdx(shapes_drawing(), os.path.join(OUT, "shapes.vsdx"))
    flow = flow_drawing()
    write_vsdx(flow, os.path.join(OUT, "flow.vsdx"))
    write_vdx(flow, os.path.join(OUT, "flow.vdx"))


if __name__ == "__main__":
    main()
