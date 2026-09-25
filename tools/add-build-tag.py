#!/usr/bin/env python3
"""Adds a build constraint to every Go file in a directory (merging with an existing one).
Usage: tools/add-build-tag.py DIR 'bdf_avif'"""
import os, re, sys
d, tag = sys.argv[1], sys.argv[2]
for root, _, files in os.walk(d):
    for f in files:
        if not f.endswith(".go"):
            continue
        p = os.path.join(root, f)
        s = open(p).read()
        m = re.match(r"//go:build (.*)\n", s)
        if m:
            expr = m.group(1)
            if tag in expr:
                continue
            s = s[m.end():]
            new = f"//go:build {tag} && ({expr})\n"
        else:
            new = f"//go:build {tag}\n\n"
        open(p, "w").write(new + s.lstrip("\n") if not m else new + s)
print("tagged", d)
