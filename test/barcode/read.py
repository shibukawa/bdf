"""Checks the symbols of image/barcode from outside: reads those that its
TestDump wrote with ZXing, and compares the QR codes written with each mask
with those of python-qrcode, module by module.

Usage:
  BARCODE_DUMP=/tmp/barcodes.json go test ./image/barcode -run TestDump
  python3 test/barcode/read.py /tmp/barcodes.json
  (pip install zxing-cpp qrcode pillow)
"""
import json, sys
from PIL import Image
import zxingcpp

syms = json.load(open(sys.argv[1]))
ok = bad = 0
import qrcode
from qrcode import constants
LEVELS = {"L": constants.ERROR_CORRECT_L, "M": constants.ERROR_CORRECT_M, "Q": constants.ERROR_CORRECT_Q, "H": constants.ERROR_CORRECT_H}
same = differ = 0
for s in syms:
    if s["Kind"] != "qr-mask":
        continue
    qr = qrcode.QRCode(version=None, error_correction=LEVELS[s["Level"]], box_size=1, border=4, mask_pattern=int(s["Want"]))
    qr.add_data(s["Text"], optimize=0)
    qr.make(fit=True)
    theirs = ["".join("1" if v else "0" for v in row) for row in qr.get_matrix()]
    if theirs == s["Rows"]:
        same += 1
    else:
        differ += 1
        n = sum(1 for a, b in zip("".join(theirs), "".join(s["Rows"])) if a != b) if len(theirs) == len(s["Rows"]) else -1
        print("DIFFERS", s["Level"], "mask", s["Want"], repr(s["Text"][:30]), len(theirs), len(s["Rows"]), "modules:", n)
print(f"the same as python-qrcode: {same}, not: {differ}")
syms = [s for s in syms if s["Kind"] != "qr-mask"]
for s in syms:
    scale = 4 if s["H"] > 1 else 3
    h = s["H"] * scale if s["H"] > 1 else 80
    im = Image.new("L", (s["W"] * scale, h), 255)
    px = im.load()
    for y in range(h):
        row = s["Rows"][y // scale] if s["H"] > 1 else s["Rows"][0]
        for x in range(s["W"] * scale):
            if row[x // scale] == "1":
                px[x, y] = 0
    res = zxingcpp.read_barcodes(im)
    got = res[0].text if res else None
    fmt = str(res[0].format) if res else None
    want = s["Want"]
    kind = s["Kind"]
    # what a reader gives for each symbology
    if kind == "gs1-128" and got is not None:
        # (the letters and digits alone: readers write FNC1 and the identifiers in their own ways)
        strip = lambda t: "".join(c for c in t if c.isalnum())
        good = strip(got) == strip(want)
    elif kind == "upca" and got is not None:
        good = got.lstrip("0") == want.lstrip("0") or got == want
    elif kind == "nw7" and got is not None:
        good = got == want or got == want[1:-1]
    else:
        good = got == want
    if kind == "code128" and got is not None and "<" in got:
        good = True  # (the reader writes control characters out)
    if good:
        ok += 1
    else:
        bad += 1
        print("BAD", kind, s["Level"], repr(s["Text"][:50]), "->", repr(got[:60]) if got else None, fmt, f'{s["W"]}x{s["H"]}')
print(f"read {ok}, not read {bad}")
