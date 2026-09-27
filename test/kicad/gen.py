#!/usr/bin/env python3
"""Writes the KiCad test projects of converter/kicad/testdata.

demo/      a schematic of a root sheet and a hierarchical sheet used twice,
           with every kind of label, pin shapes, units, rotated and mirrored
           symbols, a symbol that is not placed (DNP), markup, a text box, a
           table, an image and graphics; and a four-layer board with every
           pad shape, tracks and arcs, vias, zones, text on both sides and
           dimensions. demo.zip holds the project.
frame/     a one-sheet schematic whose project names its own drawing sheet
           (frame.kicad_wks, with a logo polygon and project variables).

The files are written here (not made with KiCad) so that they hold only what
the tests need; they follow the format of KiCad 9 (version 20250114 for
schematics, 20241229 for boards). Run from the repository root:

    python3 test/kicad/gen.py
"""
import base64
import os
import struct
import uuid as uuidlib
import zipfile
import zlib

ROOT = os.path.join(os.path.dirname(__file__), "..", "..")
OUT = os.path.join(ROOT, "converter", "kicad", "testdata")

_n = 0


def uid():
    """Deterministic UUIDs."""
    global _n
    _n += 1
    return str(uuidlib.UUID(int=_n + 0x1000))


def q(s):
    return '"' + s.replace("\\", "\\\\").replace('"', '\\"').replace("\n", "\\n") + '"'


def font(size=1.27, bold=False, italic=False, extra=""):
    b = " (bold yes)" if bold else ""
    i = " (italic yes)" if italic else ""
    return f"(font (size {size} {size}){b}{i}{extra})"


def effects(size=1.27, justify="", hide=False, bold=False, italic=False):
    j = f" (justify {justify})" if justify else ""
    h = " (hide yes)" if hide else ""
    return f"(effects {font(size, bold, italic)}{j}{h})"


def prop(name, value, x, y, angle=0, justify="", hide=False, show_name=False):
    sn = " (show_name yes)" if show_name else ""
    return f"(property {q(name)} {q(value)} (at {x} {y} {angle}){sn} {effects(justify=justify, hide=hide)})"


def stroke(width=0, typ="default"):
    return f"(stroke (width {width}) (type {typ}))"


def fill(typ="none", color=None):
    c = f" (color {color})" if color else ""
    return f"(fill (type {typ}){c})"


def pin(etype, shape, x, y, angle, length, name, number, hide=False, name_size=1.27):
    h = " (hide yes)" if hide else ""
    return (f"(pin {etype} {shape} (at {x} {y} {angle}) (length {length}){h}"
            f" (name {q(name)} (effects {font(name_size)})) (number {q(number)} (effects {font()})))")


# --- library symbols (y up, as symbol libraries write them) ---

LIB_R = f"""(symbol "Device:R" (pin_numbers (hide yes)) (pin_names (offset 0)) (exclude_from_sim no) (in_bom yes) (on_board yes)
  {prop("Reference", "R", 2.032, 0, 90)}
  {prop("Value", "R", 0, 0, 90)}
  {prop("Footprint", "", -1.778, 0, 90, hide=True)}
  (symbol "R_0_1" (rectangle (start -1.016 -2.54) (end 1.016 2.54) {stroke(0.254)} {fill()}))
  (symbol "R_1_1" {pin("passive", "line", 0, 3.81, 270, 1.27, "~", "1")} {pin("passive", "line", 0, -3.81, 90, 1.27, "~", "2")}))"""

LIB_C = f"""(symbol "Device:C_Polarized" (pin_numbers (hide yes)) (pin_names (offset 0.254)) (in_bom yes) (on_board yes)
  {prop("Reference", "C", 0.635, 2.54, 0, justify="left")}
  {prop("Value", "C_Polarized", 0.635, -2.54, 0, justify="left")}
  (symbol "C_Polarized_0_1"
    (rectangle (start -2.286 0.508) (end 2.286 1.016) {stroke(0)} {fill()})
    (polyline (pts (xy -1.778 2.286) (xy -0.762 2.286)) {stroke(0)} {fill()})
    (polyline (pts (xy -1.27 2.794) (xy -1.27 1.778)) {stroke(0)} {fill()})
    (rectangle (start 2.286 -0.508) (end -2.286 -1.016) {stroke(0)} {fill("outline")}))
  (symbol "C_Polarized_1_1" {pin("passive", "line", 0, 3.81, 270, 2.794, "~", "1")} {pin("passive", "line", 0, -3.81, 90, 2.794, "~", "2")}))"""

LIB_GND = f"""(symbol "power:GND" (power) (pin_numbers (hide yes)) (pin_names (offset 0) (hide yes)) (in_bom yes) (on_board yes)
  {prop("Reference", "#PWR", 0, -6.35, 0, hide=True)}
  {prop("Value", "GND", 0, -3.81, 0)}
  (symbol "GND_0_1" (polyline (pts (xy 0 0) (xy 0 -1.27) (xy 1.27 -1.27) (xy 0 -2.54) (xy -1.27 -1.27) (xy 0 -1.27)) {stroke(0)} {fill()}))
  (symbol "GND_1_1" {pin("power_in", "line", 0, 0, 270, 0, "GND", "1", hide=True)}))"""

# a dual op amp: unit 1 and 2 are amplifiers, unit 3 the supply pins
OPAMP_BODY = f"""(polyline (pts (xy -5.08 5.08) (xy 5.08 0) (xy -5.08 -5.08) (xy -5.08 5.08)) {stroke(0.254)} {fill("background")})"""
LIB_OPAMP = f"""(symbol "Amplifier:Dual" (pin_names (offset 0.127)) (in_bom yes) (on_board yes)
  {prop("Reference", "U", 0, 5.08, 0, justify="left")}
  {prop("Value", "DualOpAmp", 0, -5.08, 0, justify="left")}
  (symbol "Dual_1_1" {OPAMP_BODY}
    {pin("output", "line", 7.62, 0, 180, 2.54, "~", "1")}
    {pin("input", "line", -7.62, -2.54, 0, 2.54, "-", "2")}
    {pin("input", "line", -7.62, 2.54, 0, 2.54, "+", "3")})
  (symbol "Dual_2_1" {OPAMP_BODY}
    {pin("output", "line", 7.62, 0, 180, 2.54, "~", "7")}
    {pin("input", "line", -7.62, -2.54, 0, 2.54, "-", "6")}
    {pin("input", "line", -7.62, 2.54, 0, 2.54, "+", "5")})
  (symbol "Dual_3_1"
    {pin("power_in", "line", -2.54, 7.62, 270, 3.81, "V+", "8")}
    {pin("power_in", "line", -2.54, -7.62, 90, 3.81, "V-", "4")}))"""

# a logic chip with every pin shape; names inside the body
SHAPES = ["line", "inverted", "clock", "inverted_clock", "input_low", "clock_low", "output_low", "edge_clock_high", "non_logic"]
logic_pins = []
for i, s in enumerate(SHAPES):
    y = 10.16 - i * 2.54
    logic_pins.append(pin("input", s, -10.16, y, 0, 2.54, s.upper().replace("_", " ")[:6] if i % 2 else f"~{{IN{i}}}", str(i + 1)))
logic_pins.append(pin("output", "line", 10.16, 7.62, 180, 2.54, "Q", "20"))
logic_pins.append(pin("output", "inverted", 10.16, 5.08, 180, 2.54, "~{Q}", "21"))
logic_pins.append(pin("bidirectional", "line", 10.16, 0, 180, 2.54, "D_{0}", "22"))
logic_pins.append(pin("tri_state", "line", 10.16, -2.54, 180, 2.54, "A^{2}", "23"))
logic_pins.append(pin("no_connect", "line", 10.16, -7.62, 180, 2.54, "NC", "24"))
logic_pins.append(pin("power_in", "line", 0, 15.24, 270, 2.54, "VCC", "30"))
logic_pins.append(pin("power_in", "line", 0, -15.24, 90, 2.54, "GND", "31"))
LIB_LOGIC = f"""(symbol "Logic:Shapes" (pin_names (offset 1.016)) (in_bom yes) (on_board yes)
  {prop("Reference", "U", -7.62, 13.97, 0, justify="left")}
  {prop("Value", "Shapes", -7.62, -13.97, 0, justify="left")}
  (symbol "Shapes_0_1"
    (rectangle (start -7.62 12.7) (end 7.62 -12.7) {stroke(0.254)} {fill("background")})
    (circle (center 3.81 -8.89) (radius 1.27) {stroke(0.254)} {fill("outline")})
    (arc (start -3.81 -10.16) (mid -2.54 -8.89) (end -1.27 -10.16) {stroke(0.254)} {fill()})
    (text "LOGIC" (at 0 10.16 0) (effects {font(1.27, bold=True)}))
    (text "IC" (at 5.08 0 900) (effects {font(1.0)}))
    (bezier (pts (xy -5.08 -6.35) (xy -3.81 -3.81) (xy -1.27 -8.89) (xy 0 -6.35)) {stroke(0.254)} {fill()}))
  (symbol "Shapes_1_1" {" ".join(logic_pins)}))"""

# pin names outside (offset 0) on a connector
LIB_CONN = f"""(symbol "Conn:Header_4" (pin_names (offset 0)) (in_bom yes) (on_board yes)
  {prop("Reference", "J", 0, 7.62, 0)}
  {prop("Value", "Header_4", 0, -7.62, 0)}
  (symbol "Header_4_1_1"
    (rectangle (start -1.27 5.08) (end 1.27 -5.08) {stroke(0.254)} {fill("color", "255 230 200 1")})
    {pin("passive", "line", -5.08, 3.81, 0, 3.81, "SDA", "1")}
    {pin("passive", "line", -5.08, 1.27, 0, 3.81, "SCL", "2")}
    {pin("passive", "line", -5.08, -1.27, 0, 3.81, "VCC", "3")}
    {pin("passive", "line", -5.08, -3.81, 0, 3.81, "GND", "4")}))"""

LIBS = [LIB_R, LIB_C, LIB_GND, LIB_OPAMP, LIB_LOGIC, LIB_CONN]


def lib_symbols(names):
    return "(lib_symbols\n" + "\n".join(l for l in LIBS if any(f'"{n}"' in l.split("\n")[0] for n in names)) + ")"


def symbol(lib, x, y, angle, ref, value, root_uuid, paths, unit=1, mirror=None, dnp=False, fields=None, lib_name=None):
    """A placed symbol. paths are the sheet instance paths it appears on,
    with the reference of each."""
    m = f" (mirror {mirror})" if mirror else ""
    ln = f" (lib_name {q(lib_name)})" if lib_name else ""
    inst = " ".join(f'(path {q(p)} (reference {q(r)}) (unit {unit}))' for p, r in paths)
    f = fields or []
    return (f"(symbol (lib_id {q(lib)}){ln} (at {x} {y} {angle}){m} (unit {unit}) (exclude_from_sim no) (in_bom yes) (on_board yes)"
            f" (dnp {'yes' if dnp else 'no'}) (uuid {q(uid())})\n  " + "\n  ".join(f) +
            f"\n  (instances (project \"demo\" {inst})))")


def wire(*pts, bus=False):
    """Wires (or buses) through points: KiCad keeps one segment per wire."""
    kind = "bus" if bus else "wire"
    return "\n".join(f"({kind} (pts (xy {a[0]} {a[1]}) (xy {b[0]} {b[1]})) {stroke()} (uuid {q(uid())}))"
                     for a, b in zip(pts, pts[1:]))


def png(w, h, pixel):
    """A small RGBA PNG."""
    raw = b"".join(b"\x00" + b"".join(bytes(pixel(x, y)) for x in range(w)) for y in range(h))

    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xffffffff)
    return (b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 6, 0, 0, 0)) +
            chunk(b"IDAT", zlib.compress(raw, 9)) + chunk(b"IEND", b""))


def image(x, y, scale, data):
    b = base64.b64encode(data).decode()
    lines = " ".join(q(b[i:i + 76]) for i in range(0, len(b), 76))
    return f"(image (at {x} {y}) (scale {scale}) (uuid {q(uid())}) (data {lines}))"


def schematic_demo():
    root = uid()
    sub_a, sub_b = uid(), uid()
    rp = "/" + root
    items = []
    # a resistor divider, rotated and mirrored symbols
    items.append(symbol("Device:R", 50.8, 50.8, 0, "R1", "10k", root, [(rp, "R1")],
                        fields=[prop("Reference", "R1", 53.34, 49.53, 0, justify="left"), prop("Value", "10k", 53.34, 52.07, 0, justify="left"),
                                prop("Footprint", "Resistor_SMD:R_0603", 50.8, 50.8, 90, hide=True)]))
    items.append(symbol("Device:R", 66.04, 45.72, 90, "R2", "4k7", root, [(rp, "R2")],
                        fields=[prop("Reference", "R2", 66.04, 40.64, 0), prop("Value", "4k7", 66.04, 43.18, 0)]))
    items.append(symbol("Device:R", 88.9, 50.8, 0, "R3", "1k", root, [(rp, "R3")], dnp=True,
                        fields=[prop("Reference", "R3", 91.44, 49.53, 0, justify="left"), prop("Value", "1k", 91.44, 52.07, 0, justify="left"),
                                prop("Tolerance", "1%", 91.44, 54.61, 0, justify="left", show_name=True)]))
    items.append(symbol("Device:C_Polarized", 50.8, 71.12, 180, "C1", "10u", root, [(rp, "C1")], mirror="y",
                        fields=[prop("Reference", "C1", 53.34, 72.39, 0, justify="left"), prop("Value", "10u", 53.34, 69.85, 0, justify="left")]))
    items.append(symbol("power:GND", 50.8, 80.01, 0, "#PWR01", "GND", root, [(rp, "#PWR01")],
                        fields=[prop("Reference", "#PWR01", 50.8, 86.36, 0, hide=True), prop("Value", "GND", 50.8, 83.82, 0)]))
    # the op amp's three units, one mirrored
    items.append(symbol("Amplifier:Dual", 127, 50.8, 0, "U1", "LM358", root, [(rp, "U1")], unit=1,
                        fields=[prop("Reference", "U1", 127, 43.18, 0), prop("Value", "LM358", 127, 45.72, 0)]))
    items.append(symbol("Amplifier:Dual", 127, 76.2, 0, "U1", "LM358", root, [(rp, "U1")], unit=2, mirror="x",
                        fields=[prop("Reference", "U1", 127, 68.58, 0), prop("Value", "LM358", 127, 71.12, 0)]))
    items.append(symbol("Amplifier:Dual", 149.86, 63.5, 0, "U1", "LM358", root, [(rp, "U1")], unit=3,
                        fields=[prop("Reference", "U1", 152.4, 62.23, 0, justify="left"), prop("Value", "LM358", 152.4, 64.77, 0, justify="left")]))
    # the logic chip: every pin shape; rotated
    items.append(symbol("Logic:Shapes", 203.2, 63.5, 0, "U2", "Shapes", root, [(rp, "U2")],
                        fields=[prop("Reference", "U2", 195.58, 49.53, 0, justify="left"), prop("Value", "Shapes", 195.58, 77.47, 0, justify="left")]))
    items.append(symbol("Conn:Header_4", 250.19, 60.96, 270, "J1", "I2C", root, [(rp, "J1")],
                        fields=[prop("Reference", "J1", 257.81, 60.96, 90), prop("Value", "I2C", 255.27, 60.96, 90)]))
    items += [
        wire((50.8, 46.99), (50.8, 45.72), (62.23, 45.72)),
        wire((50.8, 54.61), (50.8, 67.31)),
        wire((50.8, 74.93), (50.8, 80.01)),
        wire((69.85, 45.72), (88.9, 45.72), (88.9, 46.99)),
        wire((88.9, 54.61), (88.9, 60.96), (119.38, 60.96), (119.38, 53.34)),
        wire((101.6, 25.4), (101.6, 35.56), bus=True),
        wire((101.6, 35.56), (160.02, 35.56), bus=True),
        f"(bus_entry (at 106.68 35.56) (size 2.54 2.54) {stroke()} (uuid {q(uid())}))",
        f"(junction (at 88.9 45.72) (diameter 0) (color 0 0 0 0) (uuid {q(uid())}))",
        f"(junction (at 50.8 60.96) (diameter 0) (color 0 0 0 0) (uuid {q(uid())}))",
        f"(no_connect (at 134.62 76.2) (uuid {q(uid())}))",
    ]
    # labels of every kind and spin
    for i, a in enumerate([0, 90, 180, 270]):
        items.append(f"(label {q('NET_' + str(a))} (at {30 + i * 10} 100 {a}) (effects {font()} (justify left bottom)) (uuid {q(uid())}))")
    for i, s in enumerate(["input", "output", "bidirectional", "tri_state", "passive"]):
        a = [0, 180, 90, 270, 0][i]
        items.append(f"(global_label {q('G_' + s.upper())} (shape {s}) (at {80 + i * 25} 100 {a}) (effects {font()} (justify left))"
                     f" (uuid {q(uid())}) (property \"Intersheetrefs\" \"${{INTERSHEET_REFS}}\" (at {80 + i * 25} 100 0) {effects(hide=True)}))")
    for i, s in enumerate(["input", "output", "bidirectional", "passive"]):
        a = [0, 180, 90, 270][i]
        items.append(f"(hierarchical_label {q('H_' + s.upper())} (shape {s}) (at {30 + i * 20} 115 {a}) (effects {font()} (justify left)) (uuid {q(uid())}))")
    for i, s in enumerate(["dot", "round", "diamond", "rectangle"]):
        items.append(f"(netclass_flag \"\" (length 2.54) (shape {s}) (at {120 + i * 12} 120 0) (effects {font()} (justify left bottom))"
                     f" (uuid {q(uid())}) (property \"Netclass\" \"Power\" (at {121 + i * 12} 116.84 0) {effects(justify='left', italic=True)}))")
    # notes: markup, styles, lines, rotation, Japanese
    items += [
        f"(text \"~{{RESET}} D_{{0}} x^{{2}} and plain\" (exclude_from_sim no) (at 30 140 0) (effects {font(1.524)} (justify left bottom)) (uuid {q(uid())}))",
        f"(text \"Bold italic note\\nsecond line\" (exclude_from_sim no) (at 30 150 0) (effects {font(1.524, True, True)} (justify left bottom)) (uuid {q(uid())}))",
        f"(text \"Vertical\" (exclude_from_sim no) (at 25 165 90) (effects {font(1.27)} (justify left bottom)) (uuid {q(uid())}))",
        f"(text \"カナとかな、テスト。 kana\" (exclude_from_sim no) (at 30 160 0) (effects {font(1.524)} (justify left bottom)) (uuid {q(uid())}))",
        f"(text \"${{TITLE}} rev ${{REVISION}} (${{PROJECTNAME}})\" (exclude_from_sim no) (at 30 170 0) (effects {font(1.27)} (justify left bottom)) (uuid {q(uid())}))",
        f"(text_box \"A text box wraps its text to the width of the box, inside its margins.\" (exclude_from_sim no) (at 80 135 0) (size 40 15)"
        f" (margins 0.9525 0.9525 0.9525 0.9525) {stroke(0.254, 'dash')} {fill('color', '200 220 255 1')} (effects {font()} (justify left top)) (uuid {q(uid())}))",
        f"""(table (column_count 2) (border (external yes) (header yes) {stroke(0.254)}) (separators (rows yes) (cols yes) {stroke(0.1524)})
  (column_widths 20 15) (row_heights 5 5 5)
  (cells
    (table_cell "Rail" (exclude_from_sim no) (at 130 135 0) (size 20 5) (margins 0.9525 0.9525 0.9525 0.9525) (span 1 1) {fill()} (effects {font(bold=True)} (justify left top)) (uuid {q(uid())}))
    (table_cell "Volts" (exclude_from_sim no) (at 150 135 0) (size 15 5) (margins 0.9525 0.9525 0.9525 0.9525) (span 1 1) {fill()} (effects {font(bold=True)} (justify left top)) (uuid {q(uid())}))
    (table_cell "VCC" (exclude_from_sim no) (at 130 140 0) (size 20 5) (margins 0.9525 0.9525 0.9525 0.9525) (span 1 1) {fill()} (effects {font()} (justify left top)) (uuid {q(uid())}))
    (table_cell "5.0" (exclude_from_sim no) (at 150 140 0) (size 15 5) (margins 0.9525 0.9525 0.9525 0.9525) (span 1 1) {fill()} (effects {font()} (justify right top)) (uuid {q(uid())}))
    (table_cell "VDD" (exclude_from_sim no) (at 130 145 0) (size 20 5) (margins 0.9525 0.9525 0.9525 0.9525) (span 1 1) {fill()} (effects {font()} (justify left top)) (uuid {q(uid())}))
    (table_cell "3.3" (exclude_from_sim no) (at 150 145 0) (size 15 5) (margins 0.9525 0.9525 0.9525 0.9525) (span 1 1) {fill()} (effects {font()} (justify right top)) (uuid {q(uid())}))))""",
        f"(rectangle (start 180 130) (end 200 145) {stroke(0.254, 'dash_dot')} {fill('hatch', '0 132 0 1')} (uuid {q(uid())}))",
        f"(circle (center 212 137.5) (radius 6) {stroke(0.254)} {fill('color', '255 200 200 1')} (uuid {q(uid())}))",
        f"(arc (start 222 145) (mid 228 131) (end 234 145) {stroke(0.254, 'dot')} {fill()} (uuid {q(uid())}))",
        f"(polyline (pts (xy 180 150) (xy 190 155) (xy 200 150) (xy 210 155)) {stroke(0.3, 'dash')} {fill()} (uuid {q(uid())}))",
        f"(bezier (pts (xy 215 150) (xy 220 160) (xy 230 145) (xy 235 155)) {stroke(0.254)} {fill()} (uuid {q(uid())}))",
        f"(rule_area (polyline (pts (xy 240 130) (xy 270 130) (xy 270 150) (xy 240 150)) {stroke(0.254, 'dash')} {fill()} (uuid {q(uid())})))",
        image(255, 165, 1, png(48, 32, lambda x, y: (255 * x // 47, 255 * y // 31, 160, 255))),
    ]

    def sheet(uuid, x, y, name, page):
        pins = (f"(pin \"IN\" input (at {x} {y + 5.08} 180) (uuid {q(uid())}) (effects {font()} (justify left)))"
                f" (pin \"OUT\" output (at {x + 25.4} {y + 5.08} 0) (uuid {q(uid())}) (effects {font()} (justify right)))"
                f" (pin \"BUS\" bidirectional (at {x + 12.7} {y + 15.24} 270) (uuid {q(uid())}) (effects {font()} (justify left)))")
        return (f"(sheet (at {x} {y}) (size 25.4 15.24) (exclude_from_sim no) (in_bom yes) (on_board yes) (dnp no) {stroke(0.1524, 'solid')}"
                f" (fill (color 255 255 225 1)) (uuid {q(uuid)})"
                f" (property \"Sheetname\" {q(name)} (at {x} {y - 0.7} 0) {effects(justify='left bottom')})"
                f" (property \"Sheetfile\" \"sub.kicad_sch\" (at {x} {y + 15.84} 0) {effects(justify='left top')})"
                f" {pins} (instances (project \"demo\" (path {q(rp)} (page {q(page)})))))")
    items.append(sheet(sub_a, 180, 90, "Channel A", "2"))
    items.append(sheet(sub_b, 230, 90, "Channel B", "3"))
    tb = ('(title_block (title "Converter test") (date "2026-09-28") (rev "B") (company "bdf")'
          ' (comment 1 "First comment") (comment 2 "Second comment"))')
    root_sch = (f"(kicad_sch (version 20250114) (generator \"eeschema\") (generator_version \"9.0\") (uuid {q(root)}) (paper \"A4\")\n{tb}\n"
                f"{lib_symbols(['Device:R', 'Device:C_Polarized', 'power:GND', 'Amplifier:Dual', 'Logic:Shapes', 'Conn:Header_4'])}\n"
                + "\n".join(items) + '\n(sheet_instances (path "/" (page "1"))) (embedded_fonts no))\n')
    # the sheet used twice: its symbols have a reference in each instance
    pa, pb = f"{rp}/{sub_a}", f"{rp}/{sub_b}"
    sub = [
        symbol("Device:R", 50.8, 50.8, 0, "R101", "22k", root, [(pa, "R101"), (pb, "R201")],
               fields=[prop("Reference", "R101", 53.34, 49.53, 0, justify="left"), prop("Value", "22k", 53.34, 52.07, 0, justify="left")]),
        symbol("Device:R", 76.2, 50.8, 90, "R102", "1k", root, [(pa, "R102"), (pb, "R202")],
               fields=[prop("Reference", "R102", 76.2, 45.72, 0), prop("Value", "1k", 76.2, 48.26, 0)]),
        wire((50.8, 46.99), (50.8, 40.64), (38.1, 40.64)),
        wire((50.8, 54.61), (50.8, 60.96), (72.39, 60.96), (72.39, 50.8)),
        wire((80.01, 50.8), (101.6, 50.8)),
        f"(hierarchical_label \"IN\" (shape input) (at 38.1 40.64 180) (effects {font()} (justify right)) (uuid {q(uid())}))",
        f"(hierarchical_label \"OUT\" (shape output) (at 101.6 50.8 0) (effects {font()} (justify left)) (uuid {q(uid())}))",
        f"(hierarchical_label \"BUS\" (shape bidirectional) (at 60.96 76.2 270) (effects {font()} (justify right)) (uuid {q(uid())}))",
        f"(text \"Sheet ${{#}} of ${{##}}: ${{SHEETNAME}}\" (exclude_from_sim no) (at 38.1 90 0) (effects {font(1.524)} (justify left bottom)) (uuid {q(uid())}))",
    ]
    sub_sch = (f"(kicad_sch (version 20250114) (generator \"eeschema\") (generator_version \"9.0\") (uuid {q(uid())}) (paper \"A5\")\n"
               '(title_block (title "Channel") (date "2026-09-28") (rev "B") (company "bdf"))\n'
               f"{lib_symbols(['Device:R'])}\n" + "\n".join(sub) + "\n(embedded_fonts no))\n")
    return root_sch, sub_sch


def project(name, extra_schematic=""):
    return ('{\n  "meta": {"filename": "' + name + '.kicad_pro", "version": 3},\n'
            '  "schematic": {' + extra_schematic + '},\n'
            '  "text_variables": {"PROJECT": "Demo project", "DESIGNER": "bdf tests"}\n}\n')


# --- the board ---

def fp_pad(num, typ, shape, x, y, angle, w, h, layers, extra=""):
    ls = " ".join(q(l) for l in layers)
    return f"(pad {q(num)} {typ} {shape} (at {x} {y} {angle}) (size {w} {h}) (layers {ls}){extra} (uuid {q(uid())}))"


def fp_text(kind, text, x, y, angle, layer, hide=False, mirror=False, unlocked=False):
    j = " (justify mirror)" if mirror else ""
    h = " (hide yes)" if hide else ""
    u = " (unlocked yes)" if unlocked else ""
    return (f"(property {q(kind)} {q(text)} (at {x} {y} {angle}){u} (layer {q(layer)}){h} (uuid {q(uid())})"
            f" (effects (font (size 1 1) (thickness 0.15)){j}))")


def fp_line(x1, y1, x2, y2, layer, w=0.12):
    return f"(fp_line (start {x1} {y1}) (end {x2} {y2}) (stroke (width {w}) (type solid)) (layer {q(layer)}) (uuid {q(uid())}))"


def footprint(name, layer, x, y, angle, items):
    return f"(footprint {q(name)} (layer {q(layer)}) (uuid {q(uid())}) (at {x} {y} {angle})\n  " + "\n  ".join(items) + ")"


def board():
    items = []
    # the outline: a rectangle with rounded corners
    x0, y0, x1, y1, r = 100, 80, 160, 120, 3
    items += [
        f"(gr_line (start {x0 + r} {y0}) (end {x1 - r} {y0}) (stroke (width 0.1) (type default)) (layer \"Edge.Cuts\") (uuid {q(uid())}))",
        f"(gr_line (start {x1} {y0 + r}) (end {x1} {y1 - r}) (stroke (width 0.1) (type default)) (layer \"Edge.Cuts\") (uuid {q(uid())}))",
        f"(gr_line (start {x1 - r} {y1}) (end {x0 + r} {y1}) (stroke (width 0.1) (type default)) (layer \"Edge.Cuts\") (uuid {q(uid())}))",
        f"(gr_line (start {x0} {y1 - r}) (end {x0} {y0 + r}) (stroke (width 0.1) (type default)) (layer \"Edge.Cuts\") (uuid {q(uid())}))",
    ]
    d = r * (1 - 0.70710678)
    for cx, cy, sx, sy in [(x0 + r, y0 + r, -1, -1), (x1 - r, y0 + r, 1, -1), (x1 - r, y1 - r, 1, 1), (x0 + r, y1 - r, -1, 1)]:
        start = (cx + sx * r, cy) if sx * sy > 0 else (cx, cy + sy * r)
        end = (cx, cy + sy * r) if sx * sy > 0 else (cx + sx * r, cy)
        mid = (cx + sx * r * 0.70710678, cy + sy * r * 0.70710678)
        items.append(f"(gr_arc (start {start[0]} {start[1]}) (mid {mid[0]:.4f} {mid[1]:.4f}) (end {end[0]} {end[1]}) (stroke (width 0.1) (type default)) (layer \"Edge.Cuts\") (uuid {q(uid())}))")
    # an SMD resistor, two of them turned
    for i, (x, y, a) in enumerate([(110, 90, 0), (110, 96, 90), (110, 104, 180)]):
        items.append(footprint("Resistor_SMD:R_0603", "F.Cu", x, y, a, [
            fp_text("Reference", f"R{i + 1}", 0, -1.43, a, "F.SilkS"),
            fp_text("Value", "10k", 0, 1.43, a, "F.Fab"),
            fp_line(-0.8, -0.4, 0.8, -0.4, "F.Fab", 0.1), fp_line(0.8, -0.4, 0.8, 0.4, "F.Fab", 0.1),
            fp_line(0.8, 0.4, -0.8, 0.4, "F.Fab", 0.1), fp_line(-0.8, 0.4, -0.8, -0.4, "F.Fab", 0.1),
            f"(fp_rect (start -1.48 -0.73) (end 1.48 0.73) (stroke (width 0.05) (type solid)) (fill no) (layer \"F.CrtYd\") (uuid {q(uid())}))",
            fp_pad("1", "smd", "roundrect", -0.825, 0, a, 0.8, 0.95, ["F.Cu", "F.Paste", "F.Mask"], " (roundrect_rratio 0.25)"),
            fp_pad("2", "smd", "roundrect", 0.825, 0, a, 0.8, 0.95, ["F.Cu", "F.Paste", "F.Mask"], " (roundrect_rratio 0.25)"),
        ]))
    # a SOIC-8 with rect and oval pads and a pin 1 mark
    soic = [fp_text("Reference", "U1", 0, -3.4, 0, "F.SilkS"), fp_text("Value", "SOIC-8", 0, 3.4, 0, "F.Fab"),
            f"(fp_poly (pts (xy -2.7 -2.55) (xy -2.94 -2.88) (xy -2.46 -2.88)) (stroke (width 0.12) (type solid)) (fill yes) (layer \"F.SilkS\") (uuid {q(uid())}))",
            f"(fp_circle (center -1.2 -1.2) (end -0.9 -1.2) (stroke (width 0.1) (type solid)) (fill no) (layer \"F.Fab\") (uuid {q(uid())}))",
            fp_line(-1.95, -2.45, 1.95, -2.45, "F.SilkS"), fp_line(-1.95, 2.45, 1.95, 2.45, "F.SilkS")]
    for i in range(4):
        shape = "rect" if i == 0 else "oval"
        soic.append(fp_pad(str(i + 1), "smd", shape, -2.475, -1.905 + i * 1.27, 0, 1.95, 0.6, ["F.Cu", "F.Paste", "F.Mask"]))
        soic.append(fp_pad(str(8 - i), "smd", "oval", 2.475, -1.905 + i * 1.27, 0, 1.95, 0.6, ["F.Cu", "F.Paste", "F.Mask"]))
    items.append(footprint("Package_SO:SOIC-8", "F.Cu", 125, 95, 0, soic))
    # a through-hole header: a rect pad, round pads, an oval drill
    hdr = [fp_text("Reference", "J1", 0, -2.33, 0, "F.SilkS"), fp_text("Value", "Header", 0, 10, 0, "F.Fab")]
    for i in range(4):
        shape = "rect" if i == 0 else "circle"
        hdr.append(fp_pad(str(i + 1), "thru_hole", shape, 0, i * 2.54, 0, 1.7, 1.7, ["*.Cu", "*.Mask"], " (drill 1)"))
    hdr.append(fp_pad("5", "thru_hole", "oval", 0, 10.5, 0, 1.7, 2.6, ["*.Cu", "*.Mask"], " (drill oval 1 1.8)"))
    hdr.append(fp_pad("", "np_thru_hole", "circle", 3, 5, 0, 2.2, 2.2, ["*.Cu", "*.Mask"], " (drill 2.2)"))
    items.append(footprint("Connector:Header_1x04", "F.Cu", 150, 88, 0, hdr))
    # odd pad shapes: trapezoid, chamfered, custom
    odd = [fp_text("Reference", "X1", 0, -3, 0, "F.SilkS"),
           fp_pad("1", "smd", "trapezoid", -3, 0, 0, 1.5, 2, ["F.Cu", "F.Mask"], " (rect_delta 0 0.6)"),
           fp_pad("2", "smd", "roundrect", 0, 0, 0, 1.8, 1.8, ["F.Cu", "F.Mask"], " (roundrect_rratio 0) (chamfer_ratio 0.3) (chamfer top_left bottom_right)"),
           fp_pad("3", "smd", "custom", 3, 0, 0, 0.8, 0.8, ["F.Cu", "F.Mask"],
                  " (options (clearance outline) (anchor circle)) (primitives (gr_poly (pts (xy 0 -1) (xy 1.2 0) (xy 0 1)) (width 0.1) (fill yes)))")]
    items.append(footprint("Test:OddPads", "F.Cu", 128, 110, 0, odd))
    # a footprint on the back: its text mirrored, the part turned 180°
    items.append(footprint("Resistor_SMD:R_0603", "B.Cu", 140, 110, 180, [
        fp_text("Reference", "R10", 0, -1.43, 180, "B.SilkS", mirror=True),
        fp_text("Value", "0R", 0, 1.43, 180, "B.Fab", mirror=True),
        fp_pad("1", "smd", "roundrect", -0.825, 0, 180, 0.8, 0.95, ["B.Cu", "B.Paste", "B.Mask"], " (roundrect_rratio 0.25)"),
        fp_pad("2", "smd", "roundrect", 0.825, 0, 180, 0.8, 0.95, ["B.Cu", "B.Paste", "B.Mask"], " (roundrect_rratio 0.25)"),
    ]))
    # tracks, an arc, vias
    items += [
        f"(segment (start 110.825 90) (end 118 90) (width 0.25) (layer \"F.Cu\") (net 1) (uuid {q(uid())}))",
        f"(segment (start 118 90) (end 122.525 93.095) (width 0.25) (layer \"F.Cu\") (net 1) (uuid {q(uid())}))",
        f"(arc (start 127.475 93.095) (mid 135 88) (end 142 92) (width 0.4) (layer \"F.Cu\") (net 2) (uuid {q(uid())}))",
        f"(segment (start 142 92) (end 142 108) (width 0.4) (layer \"B.Cu\") (net 2) (uuid {q(uid())}))",
        f"(segment (start 118 100) (end 135 100) (width 0.3) (layer \"In1.Cu\") (net 3) (uuid {q(uid())}))",
        f"(via (at 142 92) (size 0.8) (drill 0.4) (layers \"F.Cu\" \"B.Cu\") (net 2) (uuid {q(uid())}))",
        f"(via blind (at 118 100) (size 0.6) (drill 0.3) (layers \"F.Cu\" \"In1.Cu\") (net 3) (uuid {q(uid())}))",
        f"(via micro (at 135 100) (size 0.45) (drill 0.2) (layers \"In1.Cu\" \"In2.Cu\") (net 3) (uuid {q(uid())}))",
    ]
    # zones, as KiCad saves them filled
    def zone(layer, pts, holes=()):
        p = " ".join(f"(xy {x} {y})" for x, y in pts)
        return (f"(zone (net 4) (net_name \"GND\") (layer {q(layer)}) (uuid {q(uid())}) (hatch edge 0.5) (connect_pads (clearance 0.5))"
                f" (min_thickness 0.25) (fill yes (thermal_gap 0.5) (thermal_bridge_width 0.5)) (polygon (pts {p}))"
                f" (filled_polygon (layer {q(layer)}) (pts {p})))")
    items.append(zone("B.Cu", [(101, 81), (159, 81), (159, 119), (101, 119)]))
    items.append(zone("In2.Cu", [(102, 82), (130, 82), (130, 118), (102, 118)]))
    items.append(zone("F.Cu", [(145, 112), (158, 112), (158, 118), (145, 118)]))
    # text: front, back (mirrored), multiline, and graphics on silkscreen
    items += [
        f"(gr_text \"bdf KiCad test\" (at 130 84 0) (layer \"F.SilkS\") (uuid {q(uid())}) (effects (font (size 1.5 1.5) (thickness 0.3) bold) (justify left)))",
        f"(gr_text \"BACK SIDE\" (at 130 116 0) (layer \"B.SilkS\") (uuid {q(uid())}) (effects (font (size 1.2 1.2) (thickness 0.2)) (justify mirror)))",
        f"(gr_text \"Copper\\ntext\" (at 104 112 90) (layer \"F.Cu\") (uuid {q(uid())}) (effects (font (size 1 1) (thickness 0.2))))",
        f"(gr_text \"${{TITLE}} rev ${{REVISION}}\" (at 104 84 0) (layer \"F.Fab\") (uuid {q(uid())}) (effects (font (size 1 1) (thickness 0.15)) (justify left)))",
        f"(gr_circle (center 155 116) (end 156.5 116) (stroke (width 0.15) (type default)) (fill no) (layer \"F.SilkS\") (uuid {q(uid())}))",
        f"(gr_rect (start 103 100) (end 106 106) (stroke (width 0.15) (type default)) (fill yes) (layer \"F.SilkS\") (uuid {q(uid())}))",
        f"(gr_poly (pts (xy 150 100) (xy 156 100) (xy 153 105)) (stroke (width 0.1) (type default)) (fill yes) (layer \"B.SilkS\") (uuid {q(uid())}))",
    ]
    # dimensions of the board
    items += [
        f"""(dimension (type aligned) (layer "Dwgs.User") (uuid {q(uid())}) (pts (xy 100 120) (xy 160 120)) (height 6)
  (gr_text "60.0000 mm" (at 130 124.85 0) (layer "Dwgs.User") (uuid {q(uid())}) (effects (font (size 1 1) (thickness 0.15))))
  (format (prefix "") (suffix "") (units 3) (units_format 1) (precision 4))
  (style (thickness 0.1) (arrow_length 1.27) (text_position_mode 0) (arrow_direction outward) (extension_height 0.58642) (extension_offset 0.5) (keep_text_aligned yes)))""",
        f"""(dimension (type orthogonal) (layer "Dwgs.User") (uuid {q(uid())}) (pts (xy 160 80) (xy 160 120)) (height 6) (orientation 1)
  (gr_text "40.0000 mm" (at 164.85 100 90) (layer "Dwgs.User") (uuid {q(uid())}) (effects (font (size 1 1) (thickness 0.15))))
  (format (prefix "") (suffix "") (units 3) (units_format 1) (precision 4))
  (style (thickness 0.1) (arrow_length 1.27) (text_position_mode 0) (arrow_direction outward) (extension_height 0.58642) (extension_offset 0.5) (keep_text_aligned yes)))""",
    ]
    layers = """(layers
  (0 "F.Cu" signal) (4 "In1.Cu" signal) (6 "In2.Cu" signal) (2 "B.Cu" signal)
  (9 "F.Adhes" user "F.Adhesive") (11 "B.Adhes" user "B.Adhesive") (13 "F.Paste" user) (15 "B.Paste" user)
  (5 "F.SilkS" user "F.Silkscreen") (7 "B.SilkS" user "B.Silkscreen") (1 "F.Mask" user) (3 "B.Mask" user)
  (17 "Dwgs.User" user "User.Drawings") (19 "Cmts.User" user "User.Comments") (25 "Edge.Cuts" user)
  (27 "Margin" user) (31 "F.CrtYd" user "F.Courtyard") (29 "B.CrtYd" user "B.Courtyard") (35 "F.Fab" user) (33 "B.Fab" user))"""
    nets = '(net 0 "") (net 1 "A") (net 2 "B") (net 3 "C") (net 4 "GND")'
    return (f"(kicad_pcb (version 20241229) (generator \"pcbnew\") (generator_version \"9.0\")\n"
            '(general (thickness 1.6) (legacy_teardrops no)) (paper "A4")\n'
            '(title_block (title "Converter test") (rev "B") (company "bdf"))\n'
            f"{layers}\n(setup (pad_to_mask_clearance 0.05))\n{nets}\n" + "\n".join(items) + "\n(embedded_fonts no))\n")


# --- a project with its own drawing sheet ---

FRAME_WKS = """(kicad_wks (version 20231118) (generator "pl_editor")
  (setup (textsize 1.5 1.5) (linewidth 0.15) (textlinewidth 0.15)
    (left_margin 5) (right_margin 5) (top_margin 5) (bottom_margin 5))
  (rect (start 0 0 ltcorner) (end 0 0))
  (line (start 0 20 ltcorner) (end 0 20 rtcorner))
  (tbtext "${PROJECT}" (pos 5 10 ltcorner) (font (size 4 4) bold))
  (tbtext "Designer: ${DESIGNER}" (pos 5 10 rtcorner) (justify right))
  (tbtext "Sheet ${#}/${##}  ${TITLE}  ${ISSUE_DATE}" (pos 5 5) (justify right))
  (tbtext "1" (pos 20 2.5 lbcorner) (font (size 1.3 1.3)) (repeat 10) (incrx 40))
  (polygon (pos 30 30) (rotate 30) (linewidth 0.01)
    (pts (xy 0 0) (xy 8 0) (xy 8 8) (xy 0 8))
    (pts (xy 2 2) (xy 6 2) (xy 4 6))))
"""


def frame():
    root = uid()
    return (f"(kicad_sch (version 20250114) (generator \"eeschema\") (generator_version \"9.0\") (uuid {q(root)}) (paper \"A5\")\n"
            '(title_block (title "Custom frame") (date "2026-09-28"))\n(lib_symbols)\n'
            f"(text \"A sheet with the project's own drawing sheet\" (exclude_from_sim no) (at 30 70 0) (effects {font(2)} (justify left bottom)) (uuid {q(uid())}))\n"
            '(sheet_instances (path "/" (page "1"))) (embedded_fonts no))\n')


def main():
    root_sch, sub_sch = schematic_demo()
    files = {
        "demo/demo.kicad_sch": root_sch,
        "demo/sub.kicad_sch": sub_sch,
        "demo/demo.kicad_pcb": board(),
        "demo/demo.kicad_pro": project("demo"),
        "frame/frame.kicad_sch": frame(),
        "frame/frame.kicad_wks": FRAME_WKS,
        "frame/frame.kicad_pro": project("frame", '"page_layout_descr_file": "frame.kicad_wks"'),
    }
    for name, text in files.items():
        path = os.path.join(OUT, name)
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w", encoding="utf-8", newline="\n") as f:
            f.write(text)
    # the demo project zipped (fixed times, so the archive does not change)
    with zipfile.ZipFile(os.path.join(OUT, "demo.zip"), "w", zipfile.ZIP_DEFLATED) as z:
        for name in ["demo/demo.kicad_pro", "demo/demo.kicad_sch", "demo/sub.kicad_sch", "demo/demo.kicad_pcb"]:
            info = zipfile.ZipInfo(name, date_time=(2026, 9, 28, 0, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            z.writestr(info, files[name])


if __name__ == "__main__":
    main()
