#!/usr/bin/env python3
"""Generates converter/internal/cjkcmap/cmaps.bin.gz from Adobe's CMap resources.

The input is a directory of the predefined CJK CMaps of Adobe's
cmap-resources (https://github.com/adobe-type-tools/cmap-resources, BSD
3-Clause), e.g. the copy Ghostscript installs as Resource/CMap:
  tools/gen-cmaps.py /opt/homebrew/share/ghostscript/Resource/CMap
Every encoding CMap of the Adobe-Japan1, -Japan2, -GB1, -CNS1 and -Korea1
collections is kept (code ranges → CIDs, the parent named by usecmap, the
writing mode). The reader derives CID → Unicode from the UTF-32 CMaps (the
Adobe-*-UCS2 CMaps are not in that set).

Layout (varuint = unsigned LEB128, zigzag = signed as LEB128):
  "BCM1", varuint count, then per CMap:
    str name, str usecmap, str "Registry-Ordering", u8 wmode,
    u8 form (0 codes, 1 UTF-16 → code point, 2 UTF-8 → code point),
    str base (a CMap to look the code or code point up in after the ranges),
    varuint n, n × (u8 bytes, varuint lo, varuint hi)      codespace
    varuint n, n × (varuint lo − previous hi − 1, varuint hi − lo,
                    zigzag cid − previous run's next cid)    code ranges
"""
import gzip
import os
import re
import sys

src = sys.argv[1]
out = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "converter/internal/cjkcmap/cmaps.bin.gz")

COLLECTIONS = ["Adobe-Japan1", "Adobe-Japan2", "Adobe-GB1", "Adobe-CNS1", "Adobe-Korea1"]

TOKEN = re.compile(rb"%[^\r\n]*|<[0-9A-Fa-f\s]*>|\((?:\\.|[^)])*\)|/[^\s/<>\[\]{}()%]*|[^\s/<>\[\]{}()%]+|[\[\]{}]")


def parse(path):
    data = open(path, "rb").read()
    cm = {"name": "", "use": "", "ro": "", "wmode": 0, "space": [], "codes": {}}
    reg = ordering = ""
    toks = [t for t in TOKEN.findall(data) if not t.startswith(b"%")]
    i = 0
    prev = []
    while i < len(toks):
        t = toks[i]
        if t == b"begincodespacerange":
            i += 1
            while toks[i] != b"endcodespacerange":
                lo, hi = hexbytes(toks[i]), hexbytes(toks[i + 1])
                cm["space"].append((len(lo), int.from_bytes(lo, "big"), int.from_bytes(hi, "big")))
                i += 2
        elif t == b"begincidrange":
            i += 1
            while toks[i] != b"endcidrange":
                lo, hi, cid = hexbytes(toks[i]), hexbytes(toks[i + 1]), int(toks[i + 2])
                a, b = int.from_bytes(lo, "big"), int.from_bytes(hi, "big")
                for k in range(b - a + 1):
                    cm["codes"].setdefault(a + k, cid + k)
                i += 3
        elif t == b"begincidchar":
            i += 1
            while toks[i] != b"endcidchar":
                cm["codes"].setdefault(int.from_bytes(hexbytes(toks[i]), "big"), int(toks[i + 1]))
                i += 2
        elif t == b"beginnotdefrange":
            while toks[i] != b"endnotdefrange":
                i += 1
        elif t == b"usecmap" and prev:
            cm["use"] = prev[-1][1:].decode()
        elif t == b"def" and len(prev) >= 2:
            key, val = prev[-2], prev[-1]
            if key == b"/CMapName":
                cm["name"] = val[1:].decode()
            elif key == b"/Registry":
                reg = val[1:-1].decode()
            elif key == b"/Ordering":
                ordering = val[1:-1].decode()
            elif key == b"/WMode":
                cm["wmode"] = int(val)
        prev = (prev + [t])[-2:]
        i += 1
    cm["ro"] = reg + "-" + ordering
    return cm


def hexbytes(t):
    return bytes.fromhex(t[1:-1].decode())


def varuint(v):
    out = bytearray()
    while True:
        b = v & 0x7F
        v >>= 7
        if v:
            out.append(b | 0x80)
        else:
            out.append(b)
            return bytes(out)


def zigzag(v):
    return varuint(v * 2 if v >= 0 else -v * 2 - 1)


def string(s):
    b = s.encode()
    return varuint(len(b)) + b


def runs(pairs):
    """Consecutive (key, value) pairs with both incrementing become (start, length, value)."""
    out = []
    for k, v in pairs:
        if out and out[-1][0] + out[-1][1] == k and out[-1][2] + out[-1][1] == v:
            out[-1][1] += 1
        else:
            out.append([k, 1, v])
    return out


cmaps = {}
for name in sorted(os.listdir(src)):
    path = os.path.join(src, name)
    head = open(path, "rb").read(64)
    if not head.startswith(b"%!PS-Adobe-3.0 Resource-CMap"):
        continue
    cm = parse(path)
    if cm["ro"] not in COLLECTIONS or not cm["name"] or cm["name"] != name:
        continue
    if cm["name"].startswith("Adobe-") or cm["name"].startswith("Identity"):
        continue  # the collections' identity CMaps; Identity-H/V are built in
    cmaps[name] = cm


def resolved(name):
    """code → CID of a CMap with its usecmap parents (the CMap's own entries win)."""
    cm = cmaps[name]
    m = dict(resolved(cm["use"])) if cm["use"] in cmaps else {}
    m.update(cm["codes"])
    return m


def space(name):
    cm = cmaps[name]
    return cm["space"] or (space(cm["use"]) if cm["use"] in cmaps else [])


def codepoint(code, form):
    n = max(1, (code.bit_length() + 7) // 8)
    if form == UTF16:
        n = 4 if code > 0xFFFF else 2
    try:
        s = code.to_bytes(n, "big").decode("utf-16-be" if form == UTF16 else "utf-8")
    except (UnicodeDecodeError, OverflowError):
        return None
    return ord(s) if len(s) == 1 else None


# The UCS-2, UTF-16 and UTF-8 CMaps map the same code points to the same
# CIDs as the UTF-32 ones (UCS-2 differs in a few dozen): they are stored as
# a decoding form, the UTF-32 CMap to look the code point up in, and the
# entries that differ from it. So are the UTF-32 CMaps of a collection that
# differ from its main one in a few entries (JIS2004, JIS X 0213).
RAW, UTF16, UTF8 = 0, 1, 2
MAIN = {"Adobe-Japan1": "UniJIS", "Adobe-Japan2": "UniHojo", "Adobe-GB1": "UniGB", "Adobe-CNS1": "UniCNS", "Adobe-Korea1": "UniKS"}
records = {}
for name, cm in cmaps.items():
    rec = {"form": RAW, "base": "", "use": cm["use"], "codes": cm["codes"], "space": cm["space"]}
    m = re.match(r"^(Uni\w+?)-(UCS2|UTF16|UTF8|UTF32)(-HW)?-([HV])$", name)
    if m:
        form = {"UCS2": UTF16, "UTF16": UTF16, "UTF8": UTF8, "UTF32": RAW}[m.group(2)]
        wm = m.group(4)
        mine = {}
        for code, cid in resolved(name).items():
            cp = code if form == RAW else codepoint(code, form)
            if cp is not None:
                mine[cp] = cid
        best = None
        for base in dict.fromkeys([f"{m.group(1)}-UTF32-{wm}", f"{MAIN[cm['ro']]}-UTF32-{wm}"]):
            if base == name or base not in cmaps:
                continue
            ref = resolved(base)
            diff = {cp: cid for cp, cid in mine.items() if ref.get(cp) != cid}
            if best is None or len(diff) < len(best[1]):
                best = (base, diff)
        if best and (form != RAW or len(best[1]) < len(mine) // 2):
            rec = {"form": form, "base": best[0], "use": "", "codes": best[1], "space": space(name)}
    records[name] = rec

# Other CMaps of a collection that differ from another one in a few entries
# (GBKp-EUC-H and GBK-EUC-H) are stored the same way, over a CMap that is
# stored whole.
whole = [n for n, r in records.items() if r["form"] == RAW and not r["base"] and not r["use"]]
for name in sorted(whole):
    cm = cmaps[name]
    mine = resolved(name)
    best = None
    for base in whole:
        if base == name or cmaps[base]["ro"] != cm["ro"] or records[base]["base"]:
            continue
        ref = resolved(base)
        diff = {code: cid for code, cid in mine.items() if ref.get(code) != cid}
        if best is None or len(diff) < len(best[1]):
            best = (base, diff)
    if best and len(best[1]) < len(mine) // 4:
        records[name] = {"form": RAW, "base": best[0], "use": "", "codes": best[1], "space": space(name)}

blob = bytearray(b"BCM1")
blob += varuint(len(records))
for name, rec in sorted(records.items()):
    cm = cmaps[name]
    blob += string(name) + string(rec["use"]) + string(cm["ro"]) + bytes([cm["wmode"], rec["form"]]) + string(rec["base"])
    blob += varuint(len(rec["space"]))
    for n, lo, hi in rec["space"]:
        blob += bytes([n]) + varuint(lo) + varuint(hi)
    rs = runs(sorted(rec["codes"].items()))
    blob += varuint(len(rs))
    prev_hi, next_cid = -1, 0
    for lo, n, cid in rs:
        blob += varuint(lo - prev_hi - 1) + varuint(n - 1) + zigzag(cid - next_cid)
        prev_hi, next_cid = lo + n - 1, cid + n


with gzip.GzipFile(out, "wb", compresslevel=9, mtime=0) as f:
    f.write(bytes(blob))
print(f"{len(cmaps)} CMaps: {len(blob)} bytes, {os.path.getsize(out)} gzipped → {out}")
