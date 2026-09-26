#!/usr/bin/env python3
"""Generates converter/pptx/presets.xml.gz from presetShapeDefinitions.xml.

The input is the preset shape definition file of ECMA-376 Part 1 (Annex D,
"presetShapeDefinitions.xml"), e.g. the copy in LibreOffice:
  https://raw.githubusercontent.com/LibreOffice/core/master/oox/source/drawingml/customshapes/presetShapeDefinitions.xml
Connection sites and adjust handles are dropped (the renderer only needs
guides, the text rectangle and the paths), and so are namespace
declarations and whitespace. Usage: tools/gen-presets.py presetShapeDefinitions.xml
"""
import gzip
import os
import sys
import xml.etree.ElementTree as ET

src = sys.argv[1]
root = ET.parse(src).getroot()
out = ET.Element("presets")
for shape in root:
    s = ET.SubElement(out, shape.tag.split("}")[-1])
    for part in shape:
        tag = part.tag.split("}")[-1]
        if tag in ("cxnLst", "ahLst"):
            continue
        p = ET.SubElement(s, tag, part.attrib)
        def copy(src_el, dst_el):
            for c in src_el:
                d = ET.SubElement(dst_el, c.tag.split("}")[-1], c.attrib)
                copy(c, d)
        copy(part, p)
data = ET.tostring(out, encoding="utf-8")
dst = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "converter", "pptx", "presets.xml.gz")
with open(dst, "wb") as f:
    with gzip.GzipFile(filename="", mode="wb", fileobj=f, mtime=0, compresslevel=9) as g:
        g.write(data)
print(f"{dst}: {len(root)} shapes, {len(data)} bytes, {os.path.getsize(dst)} gzipped")
