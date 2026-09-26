#!/usr/bin/env python3
"""Generates converter/xlsx/tablestyles.xml.gz from presetTableStyles.xml.

The input is the built-in table style file of ECMA-376 Part 1 (the
SpreadsheetML styles annex, "presetTableStyles.xml"), e.g. the copy in
Apache POI:
  https://raw.githubusercontent.com/apache/poi/trunk/poi-ooxml/src/main/resources/org/apache/poi/xssf/usermodel/presetTableStyles.xml
Only the 60 table styles are kept (TableStyleLight1-21, TableStyleMedium1-28,
TableStyleDark1-11; not the pivot table styles), without namespaces,
counts and whitespace. Each style keeps its own dxfs, which its
tableStyleElements refer to by 1-based dxfId.
Usage: tools/gen-tablestyles.py presetTableStyles.xml
"""
import gzip
import os
import sys
import xml.etree.ElementTree as ET

src = sys.argv[1]
root = ET.parse(src).getroot()


def local(tag):
    return tag.split("}")[-1]


def copy(src_el, dst_el):
    for c in src_el:
        attrib = {k: v for k, v in c.attrib.items() if k != "count"}
        d = ET.SubElement(dst_el, local(c.tag), attrib)
        copy(c, d)


out = ET.Element("tableStyles")
n = 0
for style in sorted(root, key=lambda e: local(e.tag)):
    name = local(style.tag)
    if not name.startswith("TableStyle"):
        continue
    s = ET.SubElement(out, "style", {"name": name})
    copy(style, s)
    n += 1
data = ET.tostring(out, encoding="utf-8")
dst = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "converter", "xlsx", "tablestyles.xml.gz")
with open(dst, "wb") as f:
    with gzip.GzipFile(filename="", mode="wb", fileobj=f, mtime=0, compresslevel=9) as g:
        g.write(data)
print(f"{dst}: {n} styles, {len(data)} bytes, {os.path.getsize(dst)} gzipped")
