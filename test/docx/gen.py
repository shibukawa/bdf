"""Generates the Word test documents in converter/docx/testdata.

Usage: python3 test/docx/gen.py   (standard library only)

The documents are written as WordprocessingML directly, so that each
feature is in the markup exactly as Word writes it:

- basic.docx: styles and headings, character formatting, hyperlinks and
  bookmarks, East Asian text (line breaking rules, spacing between Japanese
  and Latin text, justification), bullets and numbering with nested levels,
  tab stops with leaders and PAGEREF fields, a table with a header row,
  merged cells and a table style, a footnote, inline and floating pictures,
  a text box, page and column breaks, a two-column section, paragraph
  borders and shading, headers and footers with page fields.
- grid.docx: a Japanese document on the document grid (lines snap to its
  pitch), with 1.5 and exact line spacing, character indents, a vertical
  merge and a long table whose rows split across pages and whose header
  row repeats.
- vertical.docx: East Asian vertical text (tbRl) on a character grid:
  upright characters, punctuation, small kana, turned Latin text and
  numbers, horizontal-in-vertical numbers, emphasis marks and underlines, numbered articles, an inline
  picture, a table, a footnote, and horizontal headers and footers.
"""
import io
import os
import struct
import zipfile
import zlib

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "docx", "testdata")

W = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
NS = (
    'xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" '
    'xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" '
    'xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" '
    'xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" '
    'xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture" '
    'xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape" '
    'xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" '
    'xmlns:v="urn:schemas-microsoft-com:vml" '
    'xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" mc:Ignorable="w14"'
)
REL = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"


def png(w, h, fn):
    """A small RGB PNG from fn(x, y) -> (r, g, b)."""
    raw = b""
    for y in range(h):
        raw += b"\x00" + b"".join(bytes(fn(x, y)) for x in range(w))

    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xFFFFFFFF)

    return (b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
            + chunk(b"IDAT", zlib.compress(raw, 9)) + chunk(b"IEND", b""))


def esc(s):
    return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;").replace('"', "&quot;")


def run(text, rpr=""):
    """A run; text may hold \t for tabs."""
    parts = text.split("\t")
    out = []
    for i, p in enumerate(parts):
        if i:
            out.append("<w:tab/>")
        if p:
            out.append(f'<w:t xml:space="preserve">{esc(p)}</w:t>')
    rp = f"<w:rPr>{rpr}</w:rPr>" if rpr else ""
    return f"<w:r>{rp}{''.join(out)}</w:r>"


def para(content, ppr="", style=None):
    s = f'<w:pStyle w:val="{style}"/>' if style else ""
    p = f"<w:pPr>{s}{ppr}</w:pPr>" if (s or ppr) else ""
    return f"<w:p>{p}{content}</w:p>"


def field(instr, result, rpr=""):
    rp = f"<w:rPr>{rpr}</w:rPr>" if rpr else ""
    return (f'<w:r>{rp}<w:fldChar w:fldCharType="begin"/></w:r>'
            f'<w:r>{rp}<w:instrText xml:space="preserve"> {esc(instr)} </w:instrText></w:r>'
            f'<w:r>{rp}<w:fldChar w:fldCharType="separate"/></w:r>'
            f'{run(result, rpr)}'
            f'<w:r>{rp}<w:fldChar w:fldCharType="end"/></w:r>')


def inline_pic(rid, cx, cy, pid, name, descr):
    return (f'<w:r><w:drawing><wp:inline distT="0" distB="0" distL="0" distR="0">'
            f'<wp:extent cx="{cx}" cy="{cy}"/><wp:effectExtent l="0" t="0" r="0" b="0"/>'
            f'<wp:docPr id="{pid}" name="{name}" descr="{esc(descr)}"/>'
            f'<a:graphic><a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">'
            f'<pic:pic><pic:nvPicPr><pic:cNvPr id="{pid}" name="{name}"/><pic:cNvPicPr/></pic:nvPicPr>'
            f'<pic:blipFill><a:blip r:embed="{rid}"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>'
            f'<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="{cx}" cy="{cy}"/></a:xfrm>'
            f'<a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr></pic:pic>'
            f'</a:graphicData></a:graphic></wp:inline></w:drawing></w:r>')


def anchor(inner, cx, cy, pid, name, descr, h, v, wrap, behind=False, z=1):
    """A floating drawing; h and v are (relativeFrom, align or offset in EMU)."""
    def pos(tag, rel, p):
        if isinstance(p, str):
            return f'<wp:{tag} relativeFrom="{rel}"><wp:align>{p}</wp:align></wp:{tag}>'
        return f'<wp:{tag} relativeFrom="{rel}"><wp:posOffset>{p}</wp:posOffset></wp:{tag}>'
    return (f'<w:r><w:drawing><wp:anchor distT="0" distB="0" distL="114300" distR="114300" simplePos="0" '
            f'relativeHeight="{z}" behindDoc="{1 if behind else 0}" locked="0" layoutInCell="1" allowOverlap="1">'
            f'<wp:simplePos x="0" y="0"/>{pos("positionH", *h)}{pos("positionV", *v)}'
            f'<wp:extent cx="{cx}" cy="{cy}"/><wp:effectExtent l="0" t="0" r="0" b="0"/>{wrap}'
            f'<wp:docPr id="{pid}" name="{name}" descr="{esc(descr)}"/><wp:cNvGraphicFramePr/>'
            f'<a:graphic>{inner}</a:graphic></wp:anchor></w:drawing></w:r>')


def pic_graphic(rid, cx, cy, pid, name):
    return (f'<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">'
            f'<pic:pic><pic:nvPicPr><pic:cNvPr id="{pid}" name="{name}"/><pic:cNvPicPr/></pic:nvPicPr>'
            f'<pic:blipFill><a:blip r:embed="{rid}"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>'
            f'<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="{cx}" cy="{cy}"/></a:xfrm>'
            f'<a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr></pic:pic></a:graphicData>')


def textbox_graphic(cx, cy, content):
    """A rounded rectangle shape holding a text box, as Word writes it
    (inside mc:AlternateContent, with a VML fallback)."""
    return (f'<a:graphicData uri="http://schemas.microsoft.com/office/word/2010/wordprocessingShape">'
            f'<wps:wsp><wps:cNvSpPr txBox="1"/><wps:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="{cx}" cy="{cy}"/></a:xfrm>'
            f'<a:prstGeom prst="roundRect"><a:avLst/></a:prstGeom><a:solidFill><a:srgbClr val="FFF2CC"/></a:solidFill>'
            f'<a:ln w="12700"><a:solidFill><a:srgbClr val="BF9000"/></a:solidFill></a:ln></wps:spPr>'
            f'<wps:txbx><w:txbxContent>{content}</w:txbxContent></wps:txbx>'
            f'<wps:bodyPr rot="0" vert="horz" wrap="square" lIns="91440" tIns="45720" rIns="91440" bIns="45720" anchor="ctr">'
            f'<a:noAutofit/></wps:bodyPr></wps:wsp></a:graphicData>')


def alternate(choice_run, fallback_run):
    return (f'<mc:AlternateContent><mc:Choice Requires="wps">{choice_run}</mc:Choice>'
            f'<mc:Fallback>{fallback_run}</mc:Fallback></mc:AlternateContent>')


STYLES = f"""<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles {NS}>
<w:docDefaults>
 <w:rPrDefault><w:rPr><w:rFonts w:asciiTheme="minorHAnsi" w:eastAsiaTheme="minorEastAsia" w:hAnsiTheme="minorHAnsi" w:cstheme="minorBidi"/>
  <w:kern w:val="2"/><w:sz w:val="21"/><w:szCs w:val="22"/><w:lang w:val="en-US" w:eastAsia="ja-JP" w:bidi="ar-SA"/></w:rPr></w:rPrDefault>
 <w:pPrDefault><w:pPr><w:spacing w:after="120" w:line="259" w:lineRule="auto"/></w:pPr></w:pPrDefault>
</w:docDefaults>
<w:style w:type="paragraph" w:default="1" w:styleId="a"><w:name w:val="Normal"/><w:qFormat/><w:pPr><w:widowControl/><w:jc w:val="both"/></w:pPr></w:style>
<w:style w:type="character" w:default="1" w:styleId="a0"><w:name w:val="Default Paragraph Font"/></w:style>
<w:style w:type="table" w:default="1" w:styleId="a1"><w:name w:val="Normal Table"/>
 <w:tblPr><w:tblInd w:w="0" w:type="dxa"/><w:tblCellMar><w:top w:w="0" w:type="dxa"/><w:left w:w="108" w:type="dxa"/><w:bottom w:w="0" w:type="dxa"/><w:right w:w="108" w:type="dxa"/></w:tblCellMar></w:tblPr></w:style>
<w:style w:type="paragraph" w:styleId="Title"><w:name w:val="Title"/><w:basedOn w:val="a"/><w:next w:val="a"/>
 <w:pPr><w:spacing w:after="240"/><w:jc w:val="center"/></w:pPr><w:rPr><w:rFonts w:asciiTheme="majorHAnsi" w:eastAsiaTheme="majorEastAsia"/><w:b/><w:color w:val="1F3864"/><w:sz w:val="40"/></w:rPr></w:style>
<w:style w:type="paragraph" w:styleId="1"><w:name w:val="heading 1"/><w:basedOn w:val="a"/><w:next w:val="a"/>
 <w:pPr><w:keepNext/><w:spacing w:before="240" w:after="120"/><w:outlineLvl w:val="0"/></w:pPr>
 <w:rPr><w:rFonts w:asciiTheme="majorHAnsi" w:eastAsiaTheme="majorEastAsia"/><w:b/><w:color w:themeColor="accent1" w:themeShade="BF"/><w:sz w:val="28"/></w:rPr></w:style>
<w:style w:type="paragraph" w:styleId="2"><w:name w:val="heading 2"/><w:basedOn w:val="a"/><w:next w:val="a"/>
 <w:pPr><w:keepNext/><w:spacing w:before="160" w:after="80"/><w:outlineLvl w:val="1"/></w:pPr>
 <w:rPr><w:b/><w:sz w:val="24"/></w:rPr></w:style>
<w:style w:type="paragraph" w:styleId="List"><w:name w:val="List Paragraph"/><w:basedOn w:val="a"/><w:pPr><w:ind w:left="720"/><w:contextualSpacing/></w:pPr></w:style>
<w:style w:type="paragraph" w:styleId="FootnoteText"><w:name w:val="footnote text"/><w:basedOn w:val="a"/><w:pPr><w:spacing w:after="0"/><w:jc w:val="left"/></w:pPr><w:rPr><w:sz w:val="18"/></w:rPr></w:style>
<w:style w:type="character" w:styleId="FootnoteReference"><w:name w:val="footnote reference"/><w:rPr><w:vertAlign w:val="superscript"/></w:rPr></w:style>
<w:style w:type="character" w:styleId="Hyperlink"><w:name w:val="Hyperlink"/><w:rPr><w:color w:themeColor="hyperlink"/><w:u w:val="single"/></w:rPr></w:style>
<w:style w:type="character" w:styleId="Strong"><w:name w:val="Strong"/><w:rPr><w:b/></w:rPr></w:style>
<w:style w:type="paragraph" w:styleId="Header"><w:name w:val="header"/><w:basedOn w:val="a"/>
 <w:pPr><w:tabs><w:tab w:val="center" w:pos="4513"/><w:tab w:val="right" w:pos="9026"/></w:tabs><w:spacing w:after="0"/></w:pPr><w:rPr><w:sz w:val="18"/></w:rPr></w:style>
<w:style w:type="paragraph" w:styleId="TOC1"><w:name w:val="toc 1"/><w:basedOn w:val="a"/>
 <w:pPr><w:tabs><w:tab w:val="right" w:leader="dot" w:pos="9016"/></w:tabs><w:spacing w:after="60"/></w:pPr></w:style>
<w:style w:type="table" w:styleId="Grid"><w:name w:val="Grid Table 4 Accent 1"/><w:basedOn w:val="a1"/>
 <w:tblPr><w:tblStyleRowBandSize w:val="1"/><w:tblStyleColBandSize w:val="1"/>
  <w:tblBorders><w:top w:val="single" w:sz="4" w:color="8EAADB"/><w:left w:val="single" w:sz="4" w:color="8EAADB"/><w:bottom w:val="single" w:sz="4" w:color="8EAADB"/>
  <w:right w:val="single" w:sz="4" w:color="8EAADB"/><w:insideH w:val="single" w:sz="4" w:color="8EAADB"/><w:insideV w:val="single" w:sz="4" w:color="8EAADB"/></w:tblBorders></w:tblPr>
 <w:tblStylePr w:type="firstRow"><w:rPr><w:b/><w:color w:val="FFFFFF"/></w:rPr><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="4472C4"/></w:tcPr></w:tblStylePr>
 <w:tblStylePr w:type="band1Horz"><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="D9E2F3"/></w:tcPr></w:tblStylePr>
</w:style>
</w:styles>"""

NUMBERING = f"""<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:numbering {NS}>
<w:abstractNum w:abstractNumId="0"><w:multiLevelType w:val="hybridMultilevel"/>
 <w:lvl w:ilvl="0"><w:start w:val="1"/><w:numFmt w:val="bullet"/><w:lvlText w:val="●"/><w:lvlJc w:val="left"/><w:pPr><w:ind w:left="420" w:hanging="420"/></w:pPr><w:rPr><w:rFonts w:ascii="Wingdings" w:hAnsi="Wingdings" w:hint="default"/></w:rPr></w:lvl>
 <w:lvl w:ilvl="1"><w:start w:val="1"/><w:numFmt w:val="bullet"/><w:lvlText w:val="◆"/><w:lvlJc w:val="left"/><w:pPr><w:ind w:left="840" w:hanging="420"/></w:pPr></w:lvl>
</w:abstractNum>
<w:abstractNum w:abstractNumId="1"><w:multiLevelType w:val="hybridMultilevel"/>
 <w:lvl w:ilvl="0"><w:start w:val="1"/><w:numFmt w:val="decimal"/><w:lvlText w:val="%1."/><w:lvlJc w:val="left"/><w:pPr><w:ind w:left="420" w:hanging="420"/></w:pPr></w:lvl>
 <w:lvl w:ilvl="1"><w:start w:val="1"/><w:numFmt w:val="lowerLetter"/><w:lvlText w:val="%2)"/><w:lvlJc w:val="left"/><w:pPr><w:ind w:left="840" w:hanging="420"/></w:pPr></w:lvl>
</w:abstractNum>
<w:abstractNum w:abstractNumId="2"><w:multiLevelType w:val="hybridMultilevel"/>
 <w:lvl w:ilvl="0"><w:start w:val="1"/><w:numFmt w:val="decimalFullWidth"/><w:lvlText w:val="第%1条"/><w:lvlJc w:val="left"/><w:suff w:val="space"/><w:pPr><w:ind w:left="0" w:firstLine="0"/></w:pPr></w:lvl>
</w:abstractNum>
<w:num w:numId="1"><w:abstractNumId w:val="0"/></w:num>
<w:num w:numId="2"><w:abstractNumId w:val="1"/></w:num>
<w:num w:numId="3"><w:abstractNumId w:val="2"/></w:num>
</w:numbering>"""

THEME = """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme"><a:themeElements>
<a:clrScheme name="Office"><a:dk1><a:sysClr val="windowText" lastClr="000000"/></a:dk1><a:lt1><a:sysClr val="window" lastClr="FFFFFF"/></a:lt1>
<a:dk2><a:srgbClr val="44546A"/></a:dk2><a:lt2><a:srgbClr val="E7E6E6"/></a:lt2><a:accent1><a:srgbClr val="4472C4"/></a:accent1>
<a:accent2><a:srgbClr val="ED7D31"/></a:accent2><a:accent3><a:srgbClr val="A5A5A5"/></a:accent3><a:accent4><a:srgbClr val="FFC000"/></a:accent4>
<a:accent5><a:srgbClr val="5B9BD5"/></a:accent5><a:accent6><a:srgbClr val="70AD47"/></a:accent6><a:hlink><a:srgbClr val="0563C1"/></a:hlink>
<a:folHlink><a:srgbClr val="954F72"/></a:folHlink></a:clrScheme>
<a:fontScheme name="Office"><a:majorFont><a:latin typeface="Arial"/><a:ea typeface=""/><a:cs typeface=""/><a:font script="Jpan" typeface="游ゴシック Light"/></a:majorFont>
<a:minorFont><a:latin typeface="Century"/><a:ea typeface=""/><a:cs typeface=""/><a:font script="Jpan" typeface="游明朝"/></a:minorFont></a:fontScheme>
<a:fmtScheme name="Office"><a:fillStyleLst><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:fillStyleLst>
<a:lnStyleLst><a:ln w="6350"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln><a:ln w="12700"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln><a:ln w="19050"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln></a:lnStyleLst>
<a:effectStyleLst><a:effectStyle><a:effectLst/></a:effectStyle><a:effectStyle><a:effectLst/></a:effectStyle><a:effectStyle><a:effectLst/></a:effectStyle></a:effectStyleLst>
<a:bgFillStyleLst><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:bgFillStyleLst></a:fmtScheme>
</a:themeElements></a:theme>"""


def settings(extra=""):
    return f"""<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:settings {NS}><w:defaultTabStop w:val="840"/>{extra}
<w:compat><w:compatSetting w:name="compatibilityMode" w:uri="http://schemas.microsoft.com/office/word" w:val="15"/></w:compat>
<w:themeFontLang w:val="en-US" w:eastAsia="ja-JP"/>
<w:clrSchemeMapping w:bg1="light1" w:t1="dark1" w:bg2="light2" w:t2="dark2" w:accent1="accent1" w:accent2="accent2" w:accent3="accent3" w:accent4="accent4" w:accent5="accent5" w:accent6="accent6" w:hyperlink="hyperlink" w:followedHyperlink="followedHyperlink"/>
</w:settings>"""


def core(title, creator):
    return f"""<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/"
 xmlns:dcterms="http://purl.org/dc/terms/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><dc:title>{esc(title)}</dc:title>
<dc:creator>{esc(creator)}</dc:creator><dcterms:created xsi:type="dcterms:W3CDTF">2026-09-26T09:00:00Z</dcterms:created></cp:coreProperties>"""


def package(path, document, parts, rels, extra_types=""):
    """Writes a .docx: document.xml with the given part files and
    relationships (id -> (type suffix, target[, external]))."""
    ct = f"""<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/><Default Extension="png" ContentType="image/png"/>
<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
<Override PartName="/word/numbering.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"/>
<Override PartName="/word/settings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"/>
<Override PartName="/word/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>
<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
{extra_types}</Types>"""
    root_rels = f"""<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="{REL}/officeDocument" Target="word/document.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
</Relationships>"""
    doc_rels = ['<?xml version="1.0" encoding="UTF-8" standalone="yes"?>',
                '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">']
    for rid, (typ, target, *ext) in rels.items():
        mode = ' TargetMode="External"' if ext and ext[0] else ""
        doc_rels.append(f'<Relationship Id="{rid}" Type="{REL}/{typ}" Target="{esc(target)}"{mode}/>')
    doc_rels.append("</Relationships>")
    with zipfile.ZipFile(path, "w", zipfile.ZIP_DEFLATED) as z:
        def add(name, data):
            info = zipfile.ZipInfo(name, date_time=(2026, 9, 26, 0, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            z.writestr(info, data)
        add("[Content_Types].xml", ct)
        add("_rels/.rels", root_rels)
        add("word/_rels/document.xml.rels", "\n".join(doc_rels))
        add("word/document.xml", document)
        for name, data in parts.items():
            add(name, data)


def hdr_ftr(tag, body):
    return f'<?xml version="1.0" encoding="UTF-8" standalone="yes"?>\n<w:{tag} {NS}>{body}</w:{tag}>'


LOREM = ("Portable previews of office documents need a layout engine: Word files are laid out when they are "
         "opened, so the converter breaks the lines itself, with the metrics of the fonts it then embeds. "
         "The lines of this justified paragraph end at the same right edge.")
JA = ("文書の行は、変換器が折り返します。"
      "句読点や閉じ括弧は行頭に来ず、「開き括弧」"
      "は行末に残りません。Englishの単語と日本語の間には"
      "少し間隔が入り、2026年のような数字も同じです。"
      "両端揃えの行は文字の間を広げて右端を揃えます。")


def basic():
    photo = png(64, 48, lambda x, y: ((x * 4) % 256, (y * 5) % 256, 160))
    logo = png(24, 24, lambda x, y: (30, 90, 200) if (x // 6 + y // 6) % 2 else (240, 200, 40))
    hl = 'w:rStyle w:val="Hyperlink"'
    body = []
    body.append(para(run("BDF Word Converter"), style="Title"))
    # contents with dot leaders and page references
    for i, (text, bm) in enumerate([("1. Text and paragraphs", "_Toc1"), ("2. Lists and tables", "_Toc2"), ("3. Pages and columns", "_Toc3")]):
        body.append(para(f'<w:hyperlink w:anchor="{bm}">{run(text, f"<{hl}/>")}{run(chr(9))}{field("PAGEREF " + bm + " \\h", "9")}</w:hyperlink>', style="TOC1"))
    body.append(para(f'<w:bookmarkStart w:id="1" w:name="_Toc1"/>{run("Text and paragraphs")}<w:bookmarkEnd w:id="1"/>', style="1"))
    body.append(para(run(LOREM)))
    body.append(para(
        run("Runs can be ") + run("bold", "<w:b/>") + run(", ") + run("italic", "<w:i/>") + run(", ")
        + run("underlined", '<w:u w:val="single"/>') + run(", ") + run("double underlined", '<w:u w:val="double"/>') + run(", ")
        + run("struck", "<w:strike/>") + run(", ") + run("colored", '<w:color w:val="C00000"/>') + run(", ")
        + run("highlighted", '<w:highlight w:val="yellow"/>') + run(" or ") + run("small caps", "<w:smallCaps/>")
        + run("; x") + run("2", '<w:vertAlign w:val="superscript"/>') + run(" and H") + run("2", '<w:vertAlign w:val="subscript"/>')
        + run("O. A link to ") + f'<w:hyperlink r:id="rIdLink">{run("example.com", f"<{hl}/>")}</w:hyperlink>'
        + run(" and a note.") + f'<w:r><w:rPr><w:rStyle w:val="FootnoteReference"/></w:rPr><w:footnoteReference w:id="1"/></w:r>'))
    body.append(para(run(JA), '<w:ind w:firstLineChars="100" w:firstLine="210"/>'))
    body.append(para(run("傑出した傳点を突すと強調", '<w:em w:val="dot"/>')
                     + run("、均等割付けの行。"), '<w:jc w:val="distribute"/>'))
    body.append(para(f'<w:bookmarkStart w:id="2" w:name="_Toc2"/>{run("Lists and tables")}<w:bookmarkEnd w:id="2"/>', style="1"))
    for lvl, text in [(0, "A bulleted item"), (1, "A nested item"), (1, "Another nested item"), (0, "Back at the first level")]:
        body.append(para(run(text), f'<w:numPr><w:ilvl w:val="{lvl}"/><w:numId w:val="1"/></w:numPr>', style="List"))
    for lvl, text in [(0, "First step"), (1, "a detail"), (1, "another detail"), (0, "Second step, long enough to wrap onto a second line so that the hanging indent shows under the number")]:
        body.append(para(run(text), f'<w:numPr><w:ilvl w:val="{lvl}"/><w:numId w:val="2"/></w:numPr>', style="List"))
    for text in ["目的を定める。", "適用範囲を定める。"]:
        body.append(para(run(text), '<w:numPr><w:ilvl w:val="0"/><w:numId w:val="3"/></w:numPr>'))
    # table: header row, style banding, a horizontal and a vertical merge
    def tc(text, w, extra="", rpr=""):
        return f'<w:tc><w:tcPr><w:tcW w:w="{w}" w:type="dxa"/>{extra}</w:tcPr>{para(run(text, rpr), "<w:spacing w:after=\"0\"/><w:jc w:val=\"left\"/>")}</w:tc>'
    rows = [
        '<w:tr><w:trPr><w:tblHeader/></w:trPr>' + tc("Item", 2400) + tc("Quantity", 2000) + tc("Note", 4600) + '</w:tr>',
        '<w:tr>' + tc("Apples", 2400) + tc("12", 2000, '', '') + tc("Fresh from the orchard", 4600, '<w:vMerge w:val="restart"/><w:vAlign w:val="center"/>') + '</w:tr>',
        '<w:tr>' + tc("Pears", 2400) + tc("7", 2000) + tc("", 4600, '<w:vMerge/>') + '</w:tr>',
        '<w:tr>' + tc("Total of the fruit above, in one merged cell", 4400, '<w:gridSpan w:val="2"/>') + tc("合計は19個", 4600) + '</w:tr>',
    ]
    body.append('<w:tbl><w:tblPr><w:tblStyle w:val="Grid"/><w:tblW w:w="9000" w:type="dxa"/><w:tblLook w:val="04A0" w:firstRow="1" w:lastRow="0" w:firstColumn="1" w:lastColumn="0" w:noHBand="0" w:noVBand="1"/></w:tblPr>'
                '<w:tblGrid><w:gridCol w:w="2400"/><w:gridCol w:w="2000"/><w:gridCol w:w="4600"/></w:tblGrid>' + "".join(rows) + '</w:tbl>')
    body.append(para(run("An inline picture ") + inline_pic("rIdPhoto", 64 * 12700, 48 * 12700, 10, "Picture 1", "A gradient photo") + run(" sits on the baseline.")))
    body.append(para('<w:r><w:br w:type="page"/></w:r>'))
    body.append(para(f'<w:bookmarkStart w:id="3" w:name="_Toc3"/>{run("Pages and columns")}<w:bookmarkEnd w:id="3"/>', style="1"))
    body.append(para(anchor(pic_graphic("rIdLogo", 72 * 12700, 72 * 12700, 11, "Logo"), 72 * 12700, 72 * 12700, 11, "Logo", "A checkered logo",
                            ("column", "right"), ("paragraph", 0), '<wp:wrapSquare wrapText="bothSides"/>')
                     + run(LOREM + " " + LOREM)))
    box = para(run("A text box with centered text", "<w:b/>"), '<w:jc w:val="center"/><w:spacing w:after="0"/>')
    tb_run = anchor(textbox_graphic(200 * 12700, 50 * 12700, box), 200 * 12700, 50 * 12700, 12, "Text Box 1", "",
                    ("margin", "center"), ("paragraph", 6 * 12700), '<wp:wrapTopAndBottom/>', z=2)
    body.append(para(alternate(tb_run, '<w:r><w:pict><v:rect style="width:200pt;height:50pt"/></w:pict></w:r>') + run("The shape above holds a text box.")))
    body.append(para(run("A paragraph with borders and shading."),
                     '<w:pBdr><w:top w:val="single" w:sz="8" w:space="4" w:color="2F5496"/><w:left w:val="single" w:sz="8" w:space="4" w:color="2F5496"/>'
                     '<w:bottom w:val="single" w:sz="8" w:space="4" w:color="2F5496"/><w:right w:val="single" w:sz="8" w:space="4" w:color="2F5496"/></w:pBdr>'
                     '<w:shd w:val="clear" w:color="auto" w:fill="DEEAF6"/><w:ind w:left="240" w:right="240"/>'))
    # end of the first section: continuous break into two columns
    sect1 = ('<w:sectPr><w:headerReference w:type="default" r:id="rIdH1"/><w:footerReference w:type="default" r:id="rIdF1"/>'
             '<w:pgSz w:w="11906" w:h="16838"/><w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/>'
             '<w:cols w:space="720"/></w:sectPr>')
    body.append(para(run("Two columns follow."), sect1))
    for i in range(3):
        body.append(para(run(f"Column paragraph {i + 1}. " + LOREM)))
    body.append(para('<w:r><w:br w:type="column"/></w:r>' + run("After a column break.")))
    body.append(para(run(JA)))
    sect2 = ('<w:sectPr><w:type w:val="continuous"/><w:pgSz w:w="11906" w:h="16838"/>'
             '<w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/>'
             '<w:cols w:num="2" w:space="425"/></w:sectPr>')
    body.append(para(run("The end of the two columns."), sect2))
    sect3 = ('<w:sectPr><w:type w:val="continuous"/><w:pgSz w:w="11906" w:h="16838"/>'
             '<w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/>'
             '<w:cols w:space="720"/></w:sectPr>')
    body.append(para(run("Back to one column, below the columns."), sect3))
    # a landscape section on a new page, with a header of its own on its
    # first page and a picture behind the text
    body.append(para(run("A landscape section"), style="1"))
    body.append(para(anchor(pic_graphic("rIdLogo", 144 * 12700, 144 * 12700, 13, "Watermark"), 144 * 12700, 144 * 12700, 13, "Watermark", "",
                            ("margin", "center"), ("paragraph", 0), '<wp:wrapNone/>', behind=True)
                     + run(LOREM + " " + JA)))
    last = ('<w:sectPr><w:headerReference w:type="default" r:id="rIdH1"/><w:headerReference w:type="first" r:id="rIdH2"/>'
            '<w:footerReference w:type="default" r:id="rIdF1"/>'
            '<w:pgSz w:w="16838" w:h="11906" w:orient="landscape"/>'
            '<w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/>'
            '<w:cols w:space="720"/><w:titlePg/></w:sectPr>')
    document = f'<?xml version="1.0" encoding="UTF-8" standalone="yes"?>\n<w:document {NS}><w:body>{"".join(body)}{last}</w:body></w:document>'
    header = hdr_ftr("hdr", para(run("BDF test document\t\tPage ") + field("PAGE", "1") + run(" of ") + field("NUMPAGES", "3"), style="Header"))
    first_header = hdr_ftr("hdr", para(run("The first page of the landscape section"), '<w:jc w:val="center"/>', style="Header"))
    footer = hdr_ftr("ftr", para(field("PAGE", "1"), '<w:jc w:val="center"/>', style="Header"))
    footnotes = hdr_ftr("footnotes",
                        '<w:footnote w:type="separator" w:id="-1"><w:p><w:r><w:separator/></w:r></w:p></w:footnote>'
                        '<w:footnote w:type="continuationSeparator" w:id="0"><w:p><w:r><w:continuationSeparator/></w:r></w:p></w:footnote>'
                        '<w:footnote w:id="1">' + para('<w:r><w:rPr><w:rStyle w:val="FootnoteReference"/></w:rPr><w:footnoteRef/></w:r>'
                                                      + run(" Footnotes go at the bottom of the page."), style="FootnoteText") + '</w:footnote>')
    ct = ('<Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>'
          '<Override PartName="/word/header2.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>'
          '<Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>'
          '<Override PartName="/word/footnotes.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footnotes+xml"/>')
    package(os.path.join(OUT, "basic.docx"), document, {
        "word/styles.xml": STYLES, "word/numbering.xml": NUMBERING, "word/settings.xml": settings(), "word/theme/theme1.xml": THEME,
        "word/header1.xml": header, "word/header2.xml": first_header, "word/footer1.xml": footer, "word/footnotes.xml": footnotes,
        "word/media/photo.png": photo, "word/media/logo.png": logo, "docProps/core.xml": core("BDF Word test", "BDF"),
    }, {
        "rId1": ("styles", "styles.xml"), "rId2": ("numbering", "numbering.xml"), "rId3": ("settings", "settings.xml"),
        "rId4": ("theme", "theme/theme1.xml"), "rIdH1": ("header", "header1.xml"), "rIdH2": ("header", "header2.xml"), "rIdF1": ("footer", "footer1.xml"),
        "rIdFn": ("footnotes", "footnotes.xml"), "rIdPhoto": ("image", "media/photo.png"), "rIdLogo": ("image", "media/logo.png"),
        "rIdLink": ("hyperlink", "https://example.com/", True),
    }, ct)


def grid():
    """A Japanese document on an 18 pt line grid with a long table."""
    body = []
    body.append(para(run("文書グリッドの文書", "<w:b/><w:sz w:val=\"28\"/>"), '<w:jc w:val="center"/>'))
    body.append(para(run(JA)))
    body.append(para(run(JA), '<w:spacing w:line="360" w:lineRule="auto"/>'))
    body.append(para(run(JA), '<w:spacing w:line="300" w:lineRule="exact"/><w:ind w:leftChars="200" w:left="420"/>'))
    rows = []

    def tc(text, w, extra=""):
        return f'<w:tc><w:tcPr><w:tcW w:w="{w}" w:type="dxa"/>{extra}</w:tcPr>{para(run(text), "<w:jc w:val=\"left\"/>")}</w:tc>'
    rows.append('<w:tr><w:trPr><w:tblHeader/></w:trPr>' + tc("項目", 2000, '<w:shd w:val="clear" w:color="auto" w:fill="D9D9D9"/>')
                + tc("説明", 6500, '<w:shd w:val="clear" w:color="auto" w:fill="D9D9D9"/>') + '</w:tr>')
    for i in range(14):
        text = JA * (1 + i % 3)
        rows.append('<w:tr>' + tc(f"項目{i + 1}", 2000) + tc(text, 6500) + '</w:tr>')
    borders = ('<w:tblBorders><w:top w:val="single" w:sz="4" w:color="000000"/><w:left w:val="single" w:sz="4" w:color="000000"/>'
               '<w:bottom w:val="single" w:sz="4" w:color="000000"/><w:right w:val="single" w:sz="4" w:color="000000"/>'
               '<w:insideH w:val="single" w:sz="4" w:color="000000"/><w:insideV w:val="single" w:sz="4" w:color="000000"/></w:tblBorders>')
    body.append(f'<w:tbl><w:tblPr><w:tblW w:w="8500" w:type="dxa"/>{borders}</w:tblPr><w:tblGrid><w:gridCol w:w="2000"/><w:gridCol w:w="6500"/></w:tblGrid>{"".join(rows)}</w:tbl>')
    body.append(para(run("以上")))
    last = ('<w:sectPr><w:footerReference w:type="default" r:id="rIdF1"/><w:pgSz w:w="11906" w:h="16838"/>'
            '<w:pgMar w:top="1985" w:right="1701" w:bottom="1701" w:left="1701" w:header="851" w:footer="992" w:gutter="0"/>'
            '<w:cols w:space="425"/><w:docGrid w:type="lines" w:linePitch="360"/></w:sectPr>')
    document = f'<?xml version="1.0" encoding="UTF-8" standalone="yes"?>\n<w:document {NS}><w:body>{"".join(body)}{last}</w:body></w:document>'
    styles = STYLES.replace('<w:spacing w:after="120" w:line="259" w:lineRule="auto"/>', '')
    footer = hdr_ftr("ftr", para(run("- ") + field("PAGE", "1") + run(" -"), '<w:jc w:val="center"/>'))
    ct = '<Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>'
    package(os.path.join(OUT, "grid.docx"), document, {
        "word/styles.xml": styles, "word/numbering.xml": NUMBERING, "word/settings.xml": settings(), "word/theme/theme1.xml": THEME,
        "word/footer1.xml": footer, "docProps/core.xml": core("文書グリッド", "BDF"),
    }, {
        "rId1": ("styles", "styles.xml"), "rId2": ("numbering", "numbering.xml"), "rId3": ("settings", "settings.xml"),
        "rId4": ("theme", "theme/theme1.xml"), "rIdF1": ("footer", "footer1.xml"),
    }, ct)


VJA = ("\u7e26\u66f8\u304d\u3067\u306f\u3001\u884c\u306f\u4e0a\u304b\u3089\u4e0b\u3078\u9032\u307f\u3001\u53f3\u304b\u3089\u5de6\u3078"
       "\u4e26\u3073\u307e\u3059\u3002\u6f22\u5b57\u3084\u4eee\u540d\u306f\u6b63\u7acb\u3057\u3001\u53e5\u8aad\u70b9\u306f\u53f3\u4e0a\u306b"
       "\u5bc4\u308a\u307e\u3059\u3002\u9577\u97f3\u8a18\u53f7\u300c\u30fc\u300d\u3084\u304b\u3063\u3053\u306f\u884c\u3068\u4e00\u7dd2\u306b"
       "\u56de\u308a\u307e\u3059\u3002")
VJA2 = ("\u82f1\u5b57\u306eWord\u3084\u6570\u5b57\u306e2026\u306f\u3001\u6a2a\u306b\u5012\u3057\u3066\u7d44\u307f\u307e\u3059\u3002"
        "\u5c0f\u3055\u306a\u300c\u3063\u300d\u3084\u300c\u3083\u300d\u3082\u53f3\u4e0a\u306b\u5bc4\u305b\u307e\u3059\u3002")


def vertical():
    """East Asian vertical text on a character grid."""
    photo = png(48, 32, lambda x, y: (200, (x * 5) % 256, (y * 7) % 256))
    body = []
    body.append(para(run("\u7e26\u66f8\u304d\u306e\u6587\u66f8"), style="1"))
    body.append(para(run(VJA), '<w:ind w:firstLineChars="100" w:firstLine="210"/>'))
    body.append(para(run(VJA2) + f'<w:r><w:rPr><w:rStyle w:val="FootnoteReference"/></w:rPr><w:footnoteReference w:id="1"/></w:r>',
                     '<w:ind w:firstLineChars="100" w:firstLine="210"/>'))
    body.append(para(run("\u508d\u70b9", '<w:em w:val="comma"/>') + run("\u3092\u4ed8\u3051\u305f\u8a00\u8449\u3068\u3001")
                     + run("\u508d\u7dda", '<w:u w:val="single"/>') + run("\u3092\u5f15\u3044\u305f\u8a00\u8449\u3002"),
                     '<w:ind w:firstLineChars="100" w:firstLine="210"/>'))
    for text in ["\u76ee\u7684\u3092\u5b9a\u3081\u308b\u3002", "\u9069\u7528\u7bc4\u56f2\u3092\u5b9a\u3081\u308b\u3002"]:
        body.append(para(run(text), '<w:numPr><w:ilvl w:val="0"/><w:numId w:val="3"/></w:numPr>'))
    # horizontal in vertical: two digits in one square, and three compressed into it
    tcy = '<w:eastAsianLayout w:id="1" w:vert="1"/>'
    fit = '<w:eastAsianLayout w:id="2" w:vert="1" w:vertCompress="1"/>'
    body.append(para(run("2026\u5e74") + run("10", tcy) + run("\u6708") + run("26", tcy) + run("\u65e5\u3001")
                     + run("100", fit) + run("\u500b\u306e\u4f8b\u3002"), '<w:ind w:firstLineChars="100" w:firstLine="210"/>'))
    body.append(para(run("\u56f3\u306f\u884c\u306e\u4e2d\u3067\u6b63\u7acb\u3057\u307e\u3059") + inline_pic("rIdPhoto", 48 * 12700, 32 * 12700, 10, "Picture 1", "A small picture")
                     + run("\u3002")))

    def tc(text, w):
        return f'<w:tc><w:tcPr><w:tcW w:w="{w}" w:type="dxa"/></w:tcPr>{para(run(text), "<w:jc w:val=\"left\"/>")}</w:tc>'
    borders = ('<w:tblBorders><w:top w:val="single" w:sz="4" w:color="000000"/><w:left w:val="single" w:sz="4" w:color="000000"/>'
               '<w:bottom w:val="single" w:sz="4" w:color="000000"/><w:right w:val="single" w:sz="4" w:color="000000"/>'
               '<w:insideH w:val="single" w:sz="4" w:color="000000"/><w:insideV w:val="single" w:sz="4" w:color="000000"/></w:tblBorders>')
    rows = ['<w:tr>' + tc("\u9805\u76ee", 1500) + tc("\u8aac\u660e", 4000) + '</w:tr>',
            '<w:tr>' + tc("\u7e26\u66f8\u304d", 1500) + tc("\u8868\u3082\u7e26\u306b\u7d44\u307e\u308c\u307e\u3059\u3002", 4000) + '</w:tr>']
    body.append(f'<w:tbl><w:tblPr><w:tblW w:w="5500" w:type="dxa"/>{borders}</w:tblPr><w:tblGrid><w:gridCol w:w="1500"/><w:gridCol w:w="4000"/></w:tblGrid>{"".join(rows)}</w:tbl>')
    for _ in range(6):
        body.append(para(run(VJA + VJA2), '<w:ind w:firstLineChars="100" w:firstLine="210"/>'))
    last = ('<w:sectPr><w:headerReference w:type="default" r:id="rIdH1"/><w:footerReference w:type="default" r:id="rIdF1"/>'
            '<w:pgSz w:w="11906" w:h="16838"/><w:pgMar w:top="1985" w:right="1701" w:bottom="1701" w:left="1701" w:header="851" w:footer="992" w:gutter="0"/>'
            '<w:cols w:space="425"/><w:textDirection w:val="tbRl"/><w:docGrid w:type="linesAndChars" w:linePitch="420" w:charSpace="819"/></w:sectPr>')
    document = f'<?xml version="1.0" encoding="UTF-8" standalone="yes"?>\n<w:document {NS}><w:body>{"".join(body)}{last}</w:body></w:document>'
    header = hdr_ftr("hdr", para(run("\u7e26\u66f8\u304d\u306e\u6587\u66f8"), '<w:jc w:val="right"/>', style="Header"))
    footer = hdr_ftr("ftr", para(run("- ") + field("PAGE", "1") + run(" -"), '<w:jc w:val="center"/>'))
    footnotes = hdr_ftr("footnotes",
                        '<w:footnote w:type="separator" w:id="-1"><w:p><w:r><w:separator/></w:r></w:p></w:footnote>'
                        '<w:footnote w:id="1">' + para('<w:r><w:rPr><w:rStyle w:val="FootnoteReference"/></w:rPr><w:footnoteRef/></w:r>'
                                                      + run("\u811a\u6ce8\u3082\u7e26\u66f8\u304d\u3067\u3059\u3002"), style="FootnoteText") + '</w:footnote>')
    styles = STYLES.replace('<w:spacing w:after="120" w:line="259" w:lineRule="auto"/>', '')
    ct = ('<Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>'
          '<Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>'
          '<Override PartName="/word/footnotes.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footnotes+xml"/>')
    package(os.path.join(OUT, "vertical.docx"), document, {
        "word/styles.xml": styles, "word/numbering.xml": NUMBERING, "word/settings.xml": settings(), "word/theme/theme1.xml": THEME,
        "word/header1.xml": header, "word/footer1.xml": footer, "word/footnotes.xml": footnotes,
        "word/media/photo.png": photo, "docProps/core.xml": core("\u7e26\u66f8\u304d\u306e\u6587\u66f8", "BDF"),
    }, {
        "rId1": ("styles", "styles.xml"), "rId2": ("numbering", "numbering.xml"), "rId3": ("settings", "settings.xml"),
        "rId4": ("theme", "theme/theme1.xml"), "rIdH1": ("header", "header1.xml"), "rIdF1": ("footer", "footer1.xml"),
        "rIdFn": ("footnotes", "footnotes.xml"), "rIdPhoto": ("image", "media/photo.png"),
    }, ct)


if __name__ == "__main__":
    os.makedirs(OUT, exist_ok=True)
    basic()
    grid()
    vertical()
