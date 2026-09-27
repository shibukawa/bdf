"""Writes the test publications of the EPUB converter (converter/epub/testdata)
with the standard library only:

- basic.epub: an EPUB 3 book in English: a cover (an SVG that fits a
  picture), the navigation document in the spine, two chapters (headings,
  lists, a quotation, a figure, a table, code, a footnote, links between
  chapters and to an element of another chapter, a page marker written
  as an empty element, a paragraph hidden by a class of the style sheet, a
  class that centers, an SVG file and an inline SVG that uses a gradient and
  a symbol of a hidden sprite sheet and carries attributes of other
  namespaces) and notes outside the reading order.
- vertical.epub: a Japanese book bound on the right, written like the
  Denshoken (電書協) guide: html class="vrtl" and a style sheet that imports
  the rules (vertical-rl, tate-chu-yoko, emphasis marks, gaiji pictures);
  a horizontal cover and colophon (class="hltr"), ruby.
- fixed.epub: a fixed-layout book of pictures bound on the right (a comic):
  a cover as an SVG content document's picture, pages of one img each with
  their viewports, a page drawn by an inline SVG and one by an SVG file.

The Japanese text keeps to the kanji of the Word converter's test fonts
(converter/docx/testdata/fonts), which the testdata lays text out with.
"""
import os
import struct
import zipfile
import zlib

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter/epub/testdata")
DATE = (2026, 9, 27, 0, 0, 0)


def png(w, h, pixel):
    rows = b"".join(b"\x00" + b"".join(bytes(pixel(x, y)) for x in range(w)) for y in range(h))

    def chunk(kind, data):
        c = kind + data
        return struct.pack(">I", len(data)) + c + struct.pack(">I", zlib.crc32(c) & 0xFFFFFFFF)

    data = b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
    return data + chunk(b"IDAT", zlib.compress(rows, 9)) + chunk(b"IEND", b"")


def cover(bg, fg):
    # a band of color with a frame, like a book jacket
    def pixel(x, y):
        if 40 <= x < 560 and 40 <= y < 760 and not (48 <= x < 552 and 48 <= y < 752):
            return fg
        if 200 <= y < 330:
            return fg
        return bg
    return pixel


def figure(x, y):
    # three rising bars on white
    for i, (hgt, col) in enumerate([(60, (9, 105, 218)), (110, (26, 127, 55)), (160, (207, 34, 46))]):
        x0 = 30 + i * 90
        if x0 <= x < x0 + 60 and 190 - hgt <= y < 190:
            return col
    if y == 190 and 10 <= x < 310:
        return (80, 80, 80)
    return (255, 255, 255)


def gaiji(x, y):
    # a mark the fonts lack: a ring with a dot
    d = (x - 15.5) ** 2 + (y - 15.5) ** 2
    if d < 4 ** 2 or 10 ** 2 <= d < 13 ** 2:
        return (31, 35, 40)
    return (255, 255, 255)


def comic(n):
    # a page of panels: two rows, the top one split; a shade per page
    shade = [(250, 246, 238), (238, 246, 250), (246, 250, 238)][n % 3]
    def pixel(x, y):
        panels = [(30, 30, 285, 380), (315, 30, 570, 380), (30, 410, 570, 770)]
        for i, (x0, y0, x1, y1) in enumerate(panels):
            if x0 <= x < x1 and y0 <= y < y1:
                if x < x0 + 4 or x >= x1 - 4 or y < y0 + 4 or y >= y1 - 4:
                    return (20, 20, 20)
                # a figure in each panel: a disc
                cx, cy = (x0 + x1) / 2 + (n - 1) * 30, (y0 + y1) / 2
                if (x - cx) ** 2 + (y - cy) ** 2 < (40 + 15 * i) ** 2:
                    return (60 + 60 * i, 90, 160 - 40 * i)
                return shade
        return (255, 255, 255)
    return pixel


def write(name, files):
    os.makedirs(OUT, exist_ok=True)
    with zipfile.ZipFile(os.path.join(OUT, name), "w") as z:
        info = zipfile.ZipInfo("mimetype", DATE)
        info.compress_type = zipfile.ZIP_STORED
        z.writestr(info, "application/epub+zip")
        for path, data in files:
            info = zipfile.ZipInfo(path, DATE)
            info.compress_type = zipfile.ZIP_DEFLATED
            info.external_attr = 0o644 << 16
            z.writestr(info, data.encode("utf-8") if isinstance(data, str) else data)


CONTAINER = """<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="{}" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>
"""


def xhtml(title, body, lang="en", cls="", head="", ns=""):
    c = ' class="{}"'.format(cls) if cls else ""
    return """<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops"{ns} xml:lang="{lang}" lang="{lang}"{c}>
<head>
<meta charset="UTF-8"/>
<title>{title}</title>
{head}</head>
{body}
</html>
""".format(title=title, body=body, lang=lang, c=c, head=head, ns=ns)


# an SVG file of a chapter: two boxes and an arrow
DIAGRAM = """<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="240" height="90" viewBox="0 0 240 90">
<rect x="10" y="20" width="60" height="50" rx="6" fill="#1a7f37"/>
<path d="M76 45H150" stroke="#1f2328" stroke-width="4"/>
<path d="M150 35L166 45L150 55Z" fill="#1f2328"/>
<rect x="170" y="20" width="60" height="50" rx="6" fill="#cf222e"/>
</svg>
"""

# the SVG file of the last page of the comic: frames and a moon
LAST_PAGE = """<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="600" height="800" viewBox="0 0 600 800">
<rect width="600" height="800" fill="#1d2b53"/>
<circle cx="300" cy="330" r="150" fill="#fff1c1"/>
<circle cx="360" cy="290" r="130" fill="#1d2b53"/>
<rect x="30" y="30" width="540" height="740" fill="none" stroke="#fff1c1" stroke-width="6"/>
</svg>
"""


def basic():
    css = """@charset "UTF-8";
/* the reader style replaces these */
body { font-family: Georgia, serif; margin: 0 5%; }
p { text-indent: 1.5em; margin: 0; text-align: justify; }
/* read from the style sheet */
.center { text-align: center; }
.hidden-note { display: none; }
p.note { font-size: 0.8em; }
"""
    link = '<link rel="stylesheet" type="text/css" href="../css/style.css"/>\n'
    cover_page = xhtml("Cover", """<body>
<div class="cover">
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1"
     width="100%" height="100%" viewBox="0 0 600 800" preserveAspectRatio="xMidYMid meet">
<title>A Small Book</title>
<image width="600" height="800" xlink:href="../images/cover.png"/>
</svg>
</div>
</body>""", ns=' xmlns:xlink="http://www.w3.org/1999/xlink"', head=link)
    nav = xhtml("Contents", """<body>
<nav epub:type="toc" id="toc">
<h1>Contents</h1>
<ol>
<li><a href="text/ch1.xhtml">Chapter One</a></li>
<li><a href="text/ch2.xhtml">Chapter Two</a>
<ol><li><a href="text/ch2.xhtml#sec2">Tables and code</a></li></ol></li>
</ol>
</nav>
<nav epub:type="landmarks" hidden="">
<ol><li><a epub:type="bodymatter" href="text/ch1.xhtml">Start</a></li></ol>
</nav>
</body>""")
    ch1 = xhtml("Chapter One", """<body>
<section epub:type="chapter">
<h1>Chapter One</h1>
<p>This book tests the <em>EPUB converter</em>. Its chapters are laid out one after the other, each from the top of
a page, in the order of the spine. The text takes the reader style: the fonts, sizes and margins of the book's
style sheet are left out.<a epub:type="noteref" href="#fn1" id="ref1">1</a></p>
<p>A page marker written as an empty element<a id="page2"/> does not swallow the rest of the paragraph, because
the chapter is read as XML. The <a href="ch2.xhtml#sec2">section on tables</a> is in the next chapter, and
<a href="notes.xhtml#n1">a note</a> lies outside the reading order.</p>
<p class="hidden-note">This paragraph is hidden by a class of the style sheet.</p>
<ul>
<li>Lists keep their bullets.</li>
<li>Items can hold <strong>strong</strong> and <code>code</code> text.</li>
</ul>
<blockquote><p>A quotation is set in a box with a bar on its left, its text muted.</p></blockquote>
<figure>
<img src="../images/figure.png" alt="Three bars rising from left to right"/>
<figcaption>Figure 1. A picture of the publication.</figcaption>
</figure>
<aside epub:type="footnote" id="fn1"><p class="note"><a href="#ref1">1</a> A footnote, laid out where it is.</p></aside>
</section>
</body>""", head=link)
    ch2 = xhtml("Chapter Two", """<body>
<section epub:type="chapter">
<h1>Chapter Two</h1>
<p class="center">A centered line, from a class of the style sheet.</p>
<h2 id="sec2">Tables and code</h2>
<table>
<thead><tr><th>Format</th><th>Pages</th><th>Text</th></tr></thead>
<tbody>
<tr><td>Reflowable</td><td>Book pages</td><td>Laid out again</td></tr>
<tr><td>Fixed layout</td><td>One per picture</td><td>None</td></tr>
</tbody>
</table>
<pre><code>func main() {
	fmt.Println("hello, epub")
}</code></pre>
<p>Back to <a href="ch1.xhtml">the first chapter</a>, or on to <a href="https://example.com/">a web page</a>.</p>
<h2 id="svg">Pictures in SVG</h2>
<p>An SVG file: <img src="../images/diagram.svg" alt="Two boxes joined by an arrow"/></p>
<figure>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
     xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape" viewBox="0 0 200 60"
     inkscape:label="Layer 1" epub:type="illustration" role="img" aria-label="Three dots on a gradient">
<rect width="200" height="60" rx="8" fill="url(#bg)"/>
<use href="#dot" x="30" y="30"/><use href="#dot" x="100" y="30"/><use xlink:href="#dot" x="170" y="30"/>
</svg>
<figcaption>Figure 2. An inline SVG with definitions shared by the chapter.</figcaption>
</figure>
<svg xmlns="http://www.w3.org/2000/svg" style="display: none"><defs>
<linearGradient id="bg"><stop offset="0" stop-color="#cfe3ff"/><stop offset="1" stop-color="#ffe2c4"/></linearGradient>
<circle id="dot" r="14" fill="#0969da"/></defs></svg>
</section>
</body>""", head=link)
    notes = xhtml("Notes", """<body>
<section epub:type="endnotes">
<h1>Notes</h1>
<p id="n1">A note outside the reading order, reached by a link.</p>
</section>
</body>""", head=link)
    opf = """<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="uid" xml:lang="en">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="uid">urn:uuid:5b0f3c1e-8f5d-4b8e-9a51-6c1d0a7e2f10</dc:identifier>
    <dc:title id="t1">A Small Book</dc:title>
    <meta refines="#t1" property="title-type">main</meta>
    <dc:title id="t2">Tests for the EPUB converter</dc:title>
    <meta refines="#t2" property="title-type">subtitle</meta>
    <dc:creator id="c1">Test Author</dc:creator>
    <meta refines="#c1" property="role" scheme="marc:relators">aut</meta>
    <dc:language>en</dc:language>
    <dc:publisher>bdf</dc:publisher>
    <dc:date>2026-09-27</dc:date>
    <dc:subject>EPUB</dc:subject>
    <dc:subject>Testing</dc:subject>
    <dc:description>&lt;p&gt;A book that tests the &lt;em&gt;EPUB&lt;/em&gt; converter.&lt;/p&gt;</dc:description>
    <meta property="dcterms:modified">2026-09-27T00:00:00Z</meta>
    <meta name="cover" content="cover-image"/>
  </metadata>
  <manifest>
    <item id="nav" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
    <item id="cover" href="text/cover.xhtml" media-type="application/xhtml+xml" properties="svg"/>
    <item id="cover-image" href="images/cover.png" media-type="image/png" properties="cover-image"/>
    <item id="figure" href="images/figure.png" media-type="image/png"/>
    <item id="diagram" href="images/diagram.svg" media-type="image/svg+xml"/>
    <item id="css" href="css/style.css" media-type="text/css"/>
    <item id="ch1" href="text/ch1.xhtml" media-type="application/xhtml+xml"/>
    <item id="ch2" href="text/ch2.xhtml" media-type="application/xhtml+xml"/>
    <item id="notes" href="text/notes.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine>
    <itemref idref="cover"/>
    <itemref idref="nav"/>
    <itemref idref="ch1"/>
    <itemref idref="notes" linear="no"/>
    <itemref idref="ch2"/>
  </spine>
</package>
"""
    write("basic.epub", [
        ("META-INF/container.xml", CONTAINER.format("OEBPS/content.opf")),
        ("OEBPS/content.opf", opf),
        ("OEBPS/nav.xhtml", nav),
        ("OEBPS/text/cover.xhtml", cover_page),
        ("OEBPS/text/ch1.xhtml", ch1),
        ("OEBPS/text/ch2.xhtml", ch2),
        ("OEBPS/text/notes.xhtml", notes),
        ("OEBPS/css/style.css", css),
        ("OEBPS/images/cover.png", png(600, 800, cover((32, 60, 110), (240, 200, 90)))),
        ("OEBPS/images/figure.png", png(320, 200, figure)),
        ("OEBPS/images/diagram.svg", DIAGRAM),
    ])


def vertical():
    book_css = """@charset "UTF-8";
@import url("style-reset.css");
@import url("style-standard.css");
"""
    reset_css = """html, body { margin: 0; padding: 0; }
"""
    standard_css = """@charset "UTF-8";
@font-face { font-family: "serif-ja"; src: local("Hiragino Mincho ProN"); }
.hltr { writing-mode: horizontal-tb; -webkit-writing-mode: horizontal-tb; -epub-writing-mode: horizontal-tb; }
.vrtl { writing-mode: vertical-rl; -webkit-writing-mode: vertical-rl; -epub-writing-mode: vertical-rl; }
.tcy { text-combine-upright: all; -webkit-text-combine: horizontal; -epub-text-combine: horizontal; }
.em-sesame { -epub-text-emphasis-style: filled sesame; -webkit-text-emphasis-style: filled sesame; }
.align-end { text-align: right; }
.align-center { text-align: center; }
img.gaiji { width: 1em; height: 1em; }
.main p span.tcy { color: red; }
@media amzn-kf8 { .vrtl .p-text { margin: 0; } }
"""
    link = '<link rel="stylesheet" type="text/css" href="../style/book-style.css"/>\n'
    cover_page = xhtml("表紙", """<body epub:type="cover">
<div class="main">
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1"
     width="100%" height="100%" viewBox="0 0 600 800">
<image width="600" height="800" xlink:href="../image/cover.png"/>
</svg>
</div>
</body>""", lang="ja", cls="hltr", head=link)
    p1 = xhtml("１", """<body class="p-text">
<div class="main">
<h2 id="toc-001">１　縦書き</h2>
<p>縦書きの本を読むとき、行は上から下へ、右から左へ進みます。</p>
<p>この文書は、ＢＤＦの変換のためのテスト用のものです。英字のABCは横に倒して組み、漢字と仮名は正立させます。</p>
<p>数字は縦中横で、<span class="tcy">12</span>月<span class="tcy">31</span>日のように一字の中に組みます。</p>
<p><span class="em-sesame">傍点</span>を付けたことばもあります。</p>
<p>ルビは、<ruby>漢字<rt>かんじ</rt></ruby>のように本文の文字だけを出します。フォントにない字は<img class="gaiji" src="../image/gaiji.png" alt="※"/>のように一字の図にします。</p>
<p class="align-end">右寄せの行</p>
</div>
</body>""", lang="ja", cls="vrtl", head=link)
    long_text = "".join("<p>{}</p>\n".format(t) for t in [
        "ページをめくると、つぎの行がはじまります。ひとつの行に入りきらないことばは、つぎの行に回ります。",
        "「かぎかっこ」や『ふたえかぎ』のなかの文も、縦書きのまま読めます。",
        "句読点は行の頭に来ないように、行の末に残されるか、その上の字といっしょにつぎの行に回ります。",
    ] * 6)
    p2 = xhtml("２", """<body class="p-text">
<div class="main">
<h2 id="toc-002">２　ページ</h2>
{}<p><a href="p-001.xhtml#toc-001">１に返る</a></p>
</div>
</body>""".format(long_text), lang="ja", cls="vrtl", head=link)
    colophon = xhtml("奥付", """<body>
<div class="main">
<p class="align-center">縦書きのテスト</p>
<p>２０２６年９月２７日</p>
<p>書名：縦書きのテスト</p>
</div>
</body>""", lang="ja", cls="hltr", head=link)
    nav = xhtml("目次", """<body>
<nav epub:type="toc" id="toc"><h1>目次</h1>
<ol><li><a href="text/p-001.xhtml">１</a></li><li><a href="text/p-002.xhtml">２</a></li></ol>
</nav>
</body>""", lang="ja")
    opf = """<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" xml:lang="ja" unique-identifier="unique-id"
         prefix="rendition: http://www.idpf.org/vocab/rendition/#">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title id="title">縦書きのテスト</dc:title>
    <dc:creator id="creator01">テスト</dc:creator>
    <meta refines="#creator01" property="role" scheme="marc:relators">aut</meta>
    <dc:publisher>ＢＤＦ</dc:publisher>
    <dc:language>ja</dc:language>
    <dc:identifier id="unique-id">urn:uuid:0e1b7c55-3f6a-4f58-9d0f-2a7f3c9b8e41</dc:identifier>
    <meta property="dcterms:modified">2026-09-27T00:00:00Z</meta>
    <meta property="rendition:layout">reflowable</meta>
  </metadata>
  <manifest>
    <item media-type="application/xhtml+xml" id="toc" href="navigation-documents.xhtml" properties="nav"/>
    <item media-type="text/css" id="book-style" href="style/book-style.css"/>
    <item media-type="text/css" id="style-reset" href="style/style-reset.css"/>
    <item media-type="text/css" id="style-standard" href="style/style-standard.css"/>
    <item media-type="image/png" id="cover" href="image/cover.png" properties="cover-image"/>
    <item media-type="image/png" id="gaiji" href="image/gaiji.png"/>
    <item media-type="application/xhtml+xml" id="p-cover" href="text/p-cover.xhtml" properties="svg"/>
    <item media-type="application/xhtml+xml" id="p-001" href="text/p-001.xhtml"/>
    <item media-type="application/xhtml+xml" id="p-002" href="text/p-002.xhtml"/>
    <item media-type="application/xhtml+xml" id="p-colophon" href="text/p-colophon.xhtml"/>
  </manifest>
  <spine page-progression-direction="rtl">
    <itemref linear="yes" idref="p-cover" properties="rendition:page-spread-center"/>
    <itemref linear="yes" idref="p-001" properties="page-spread-left"/>
    <itemref linear="yes" idref="p-002"/>
    <itemref linear="yes" idref="p-colophon"/>
  </spine>
</package>
"""
    write("vertical.epub", [
        ("META-INF/container.xml", CONTAINER.format("item/standard.opf")),
        ("item/standard.opf", opf),
        ("item/navigation-documents.xhtml", nav),
        ("item/style/book-style.css", book_css),
        ("item/style/style-reset.css", reset_css),
        ("item/style/style-standard.css", standard_css),
        ("item/image/cover.png", png(600, 800, cover((120, 30, 40), (250, 240, 225)))),
        ("item/image/gaiji.png", png(32, 32, gaiji)),
        ("item/text/p-cover.xhtml", cover_page),
        ("item/text/p-001.xhtml", p1),
        ("item/text/p-002.xhtml", p2),
        ("item/text/p-colophon.xhtml", colophon),
    ])


def fixed():
    svg_cover = """<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1"
     width="600" height="800" viewBox="0 0 600 800">
<image width="600" height="800" xlink:href="image/cover.png"/>
</svg>
"""
    pages = []
    viewport = '<meta name="viewport" content="width=600, height=800"/>\n'
    for i in range(1, 4):
        pages.append(("item/xhtml/p-{:03d}.xhtml".format(i), xhtml("{}".format(i), """<body>
<div class="main"><img src="../image/i-{:03d}.png" alt=""/></div>
</body>""".format(i), lang="ja", head=viewport)))
    # a page drawn in SVG, and one whose picture is an SVG file
    pages.append(("item/xhtml/p-004.xhtml", xhtml("4", """<body>
<svg xmlns="http://www.w3.org/2000/svg" width="100%" height="100%" viewBox="0 0 600 800">
<defs><linearGradient id="sky" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#9ecbff"/><stop offset="1" stop-color="#ffffff"/></linearGradient></defs>
<rect width="600" height="800" fill="url(#sky)"/>
<circle cx="440" cy="170" r="70" fill="#f2c94c"/>
<path d="M0 620 Q150 520 300 600 T600 580 V800 H0 Z" fill="#2e7d32"/>
<rect x="30" y="30" width="540" height="740" fill="none" stroke="#1f2328" stroke-width="6"/>
</svg>
</body>""", lang="ja", head=viewport)))
    pages.append(("item/xhtml/p-005.xhtml", xhtml("5", """<body>
<div class="main"><img src="../image/i-005.svg" alt="The last page"/></div>
</body>""", lang="ja", head=viewport)))
    opf = """<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" xml:lang="ja" unique-identifier="unique-id"
         prefix="rendition: http://www.idpf.org/vocab/rendition/#">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>固定レイアウトのテスト</dc:title>
    <dc:creator>テスト</dc:creator>
    <dc:language>ja</dc:language>
    <dc:identifier id="unique-id">urn:uuid:9a3e0c2b-7d41-4e1f-8b6a-5c2d1e0f3a97</dc:identifier>
    <meta property="dcterms:modified">2026-09-27T00:00:00Z</meta>
    <meta property="rendition:layout">pre-paginated</meta>
    <meta property="rendition:spread">landscape</meta>
  </metadata>
  <manifest>
    <item media-type="application/xhtml+xml" id="toc" href="navigation-documents.xhtml" properties="nav"/>
    <item media-type="image/svg+xml" id="p-cover" href="cover.svg"/>
    <item media-type="image/png" id="cover" href="image/cover.png" properties="cover-image"/>
""" + "".join("""    <item media-type="image/png" id="i-{0:03d}" href="image/i-{0:03d}.png"/>
    <item media-type="application/xhtml+xml" id="p-{0:03d}" href="xhtml/p-{0:03d}.xhtml"/>
""".format(i) for i in range(1, 4)) + """    <item media-type="application/xhtml+xml" id="p-004" href="xhtml/p-004.xhtml" properties="svg"/>
    <item media-type="image/svg+xml" id="i-005" href="image/i-005.svg"/>
    <item media-type="application/xhtml+xml" id="p-005" href="xhtml/p-005.xhtml"/>
  </manifest>
  <spine page-progression-direction="rtl">
    <itemref linear="yes" idref="p-cover" properties="rendition:page-spread-center"/>
    <itemref linear="yes" idref="p-001" properties="page-spread-left"/>
    <itemref linear="yes" idref="p-002" properties="page-spread-right"/>
    <itemref linear="yes" idref="p-003" properties="page-spread-left"/>
    <itemref linear="yes" idref="p-004" properties="page-spread-right"/>
    <itemref linear="yes" idref="p-005" properties="page-spread-left"/>
    <itemref linear="no" idref="toc"/>
  </spine>
</package>
"""
    nav = xhtml("目次", """<body>
<nav epub:type="toc" id="toc"><h1>目次</h1><ol><li><a href="xhtml/p-001.xhtml">本文</a></li></ol></nav>
</body>""", lang="ja")
    files = [
        ("META-INF/container.xml", CONTAINER.format("item/standard.opf")),
        ("item/standard.opf", opf),
        ("item/navigation-documents.xhtml", nav),
        ("item/cover.svg", svg_cover),
        ("item/image/cover.png", png(600, 800, cover((30, 90, 60), (250, 250, 250)))),
    ]
    for i in range(1, 4):
        files.append(("item/image/i-{:03d}.png".format(i), png(600, 800, comic(i))))
    files.append(("item/image/i-005.svg", LAST_PAGE))
    write("fixed.epub", files + pages)


basic()
vertical()
fixed()
