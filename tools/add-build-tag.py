#!/usr/bin/env python3
"""Adds a build constraint to every Go file in a directory (merging with an
existing one, wherever it sits in the leading comment block).
Usage: tools/add-build-tag.py DIR '!bdf_noconv'"""
import os, re, sys
d, tag = sys.argv[1], sys.argv[2]
for root, _, files in os.walk(d):
    for f in files:
        if not f.endswith(".go"):
            continue
        p = os.path.join(root, f)
        s = open(p).read()
        head, sep, rest = s.partition("\npackage ")
        m = re.search(r"^//go:build (.*)$", head, re.M)
        if m:
            expr = m.group(1)
            if tag in expr:
                continue
            head = head[:m.start()] + f"//go:build {tag} && ({expr})" + head[m.end():]
            s = head + sep + rest
        else:
            s = f"//go:build {tag}\n\n" + s.lstrip("\n")
        open(p, "w").write(s)
print("tagged", d)
