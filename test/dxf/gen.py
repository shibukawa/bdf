"""Writes the DXF test drawings of converter/dxf into converter/dxf/testdata.

Requires ezdxf (pip install ezdxf). Text uses only characters of the test
fonts (converter/pptx/testdata/fonts: ASCII, Latin-1, kana and a few kanji).

- shapes.dxf: AutoCAD 2018, one of each kind of entity in model space:
  lines with line types and lineweights, circles, arcs, ellipses, polylines
  with bulges and widths, splines, points, solids, hatches (solid, pattern,
  gradient, with islands), single-line text in every alignment, multiline
  text with formatting, dimensions, a leader, a multileader, blocks with
  attributes and a block array, layers that are off or frozen.
- shapes-bin.dxf: shapes.dxf in the binary form; the Go tests check that
  both convert to the same objects.
- layout.dxf: a floor plan in model space and an A3 layout with a title
  block and two viewports at different scales, one with a frozen layer.
- r12-sjis.dxf: an AutoCAD R12 drawing in Shift_JIS ($DWGCODEPAGE ANSI_932)
  with old-style polylines.
"""

import math
import os

import ezdxf
from ezdxf import colors
from ezdxf.enums import TextEntityAlignment, MTextEntityAlignment
from ezdxf.math import Vec2

OUT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "converter", "dxf", "testdata")


def fixed(doc):
    """Makes the output reproducible (with ezdxf's fixed meta data)."""
    doc.header["$TDCREATE"] = 2460000.5


def layers(doc):
    doc.layers.add("外形", color=7, lineweight=50)
    doc.layers.add("中心", color=1, linetype="CENTER", lineweight=18)
    doc.layers.add("文字", color=2)
    doc.layers.add("ハッチ", color=3)
    doc.layers.add("OFF", color=4).off()
    doc.layers.add("FROZEN", color=5).freeze()
    doc.layers.add("DIM", color=4)


# dimension sizes for a drawing at 1:1 in mm
DIM = {"dimtxt": 2.5, "dimasz": 2.5, "dimexe": 1.25, "dimexo": 0.625, "dimgap": 0.625, "dimlfac": 1, "dimdec": 0}


def shapes():
    doc = ezdxf.new("R2018", setup=True, units=4)  # mm
    fixed(doc)
    layers(doc)
    doc.header["$PDMODE"] = 35
    doc.header["$PDSIZE"] = 3
    doc.header["$LTSCALE"] = 1
    msp = doc.modelspace()

    # --- a frame
    msp.add_lwpolyline([(0, 0), (300, 0), (300, 200), (0, 200)], close=True, dxfattribs={"layer": "外形"})

    # --- lines with line types, colors and lineweights
    y = 190
    for i, (lt, col, lw) in enumerate([("CONTINUOUS", 1, 13), ("DASHED", 2, 25), ("CENTER", 3, 35),
                                       ("DOT", 4, 50), ("DASHDOT", 5, 70), ("PHANTOM", 6, 100)]):
        msp.add_line((10, y - i * 6), (90, y - i * 6), dxfattribs={"linetype": lt, "color": col, "lineweight": lw, "ltscale": 8})
    msp.add_line((10, 150), (90, 150), dxfattribs={"true_color": colors.rgb2int((255, 128, 0)), "lineweight": 50})
    msp.add_line((10, 146), (90, 146), dxfattribs={"color": 1, "lineweight": 50, "transparency": colors.float2transparency(0.5)})
    msp.add_line((10, 142), (90, 142), dxfattribs={"layer": "中心", "ltscale": 8})
    msp.add_line((10, 138), (90, 138), dxfattribs={"layer": "OFF"})  # not drawn
    msp.add_line((10, 138), (90, 136), dxfattribs={"layer": "FROZEN"})  # not drawn
    msp.add_text("OFF LAYER", height=3, dxfattribs={"layer": "OFF"}).set_placement((10, 130))  # not drawn
    msp.add_text("FROZEN LAYER", height=3, dxfattribs={"layer": "FROZEN"}).set_placement((50, 130))  # not drawn
    msp.add_text("DEFPOINTS", height=3, dxfattribs={"layer": "Defpoints"}).set_placement((50, 130))  # not drawn

    # --- circles, arcs, ellipses
    msp.add_circle((120, 175), 12, dxfattribs={"color": 1})
    msp.add_arc((150, 175), 12, 30, 240, dxfattribs={"color": 2, "lineweight": 35})
    msp.add_ellipse((185, 175), major_axis=(15, 5), ratio=0.5, dxfattribs={"color": 3})
    msp.add_ellipse((215, 175), major_axis=(0, 12), ratio=0.6, start_param=0, end_param=math.pi * 1.5, dxfattribs={"color": 4})
    # a circle in a mirrored OCS
    msp.add_arc((-245, 175), 10, 0, 180, dxfattribs={"extrusion": (0, 0, -1), "color": 6})

    # --- polylines
    msp.add_lwpolyline([(110, 140, 0, 0, 0.5), (130, 150, 0, 0, 0), (150, 140, 0, 0, -1), (170, 140, 0, 0, 0)], format="xyseb",
                       dxfattribs={"color": 5})
    msp.add_lwpolyline([(180, 140), (200, 150), (220, 140)], dxfattribs={"const_width": 2, "color": 1})
    msp.add_lwpolyline([(230, 140, 0, 4, 0), (260, 140, 4, 0, 0), (280, 140, 0, 0, 0)], format="xyseb", dxfattribs={"color": 2})
    p2d = msp.add_polyline2d([(240, 185), (255, 190), (270, 180), (285, 190)], dxfattribs={"color": 3})
    p2d.close(True)

    # --- splines
    msp.add_open_spline([(10, 100), (25, 125), (40, 95), (55, 120), (70, 100)], degree=3, dxfattribs={"color": 4})
    msp.add_spline(fit_points=[(75, 100), (85, 120), (100, 105), (110, 118)], dxfattribs={"color": 6})
    msp.add_rational_spline([(115, 100), (125, 125), (140, 100)], weights=[1, 3, 1], degree=2, dxfattribs={"color": 1})

    # --- points and solids
    for i in range(3):
        msp.add_point((150 + i * 8, 110), dxfattribs={"color": 2})
    msp.add_solid([(180, 100), (200, 100), (180, 115), (200, 120)], dxfattribs={"color": 5})
    msp.add_solid([(205, 100), (225, 100), (215, 118)], dxfattribs={"color": 6})

    # --- hatches
    h = msp.add_hatch(color=3, dxfattribs={"layer": "ハッチ"})
    h.paths.add_polyline_path([(10, 50), (60, 50), (60, 85), (10, 85)], is_closed=True)
    h.set_pattern_fill("ANSI31", scale=0.5)
    h = msp.add_hatch(color=1)
    h.paths.add_polyline_path([(70, 50), (120, 50), (120, 85), (70, 85)], is_closed=True, flags=1)
    edge = h.paths.add_edge_path(flags=16)
    edge.add_arc((95, 67.5), 10, 0, 360)
    h.set_pattern_fill("ANSI37", scale=0.5, angle=15)
    h = msp.add_hatch(color=5)
    h.paths.add_polyline_path([(130, 50, 0), (180, 50, 0.4), (180, 85, 0), (130, 85, 0)], is_closed=True)
    h.set_solid_fill(color=5)
    h = msp.add_hatch()
    h.paths.add_polyline_path([(190, 50), (240, 50), (240, 85), (190, 85)], is_closed=True)
    h.set_gradient((255, 0, 0), (255, 255, 0), rotation=30, name="LINEAR")
    h = msp.add_hatch()
    edge = h.paths.add_edge_path()
    edge.add_line((250, 50), (290, 50))
    edge.add_arc((290, 67.5), 17.5, -90, 90)
    edge.add_line((290, 85), (250, 85))
    edge.add_ellipse((250, 67.5), major_axis=(0, 17.5), ratio=0.5, start_angle=0, end_angle=180, ccw=True)
    h.set_gradient((0, 64, 192), (255, 255, 255), name="CYLINDER")

    # --- text
    doc.styles.add("GOTHIC", font="msgothic.ttc")
    doc.styles.add("ROMANS", font="romans.shx")
    for i, (align, label) in enumerate([(TextEntityAlignment.LEFT, "LEFT 字"), (TextEntityAlignment.CENTER, "CENTER"),
                                        (TextEntityAlignment.RIGHT, "RIGHT"), (TextEntityAlignment.MIDDLE_CENTER, "MIDDLE"),
                                        (TextEntityAlignment.TOP_LEFT, "TOP"), (TextEntityAlignment.BOTTOM_RIGHT, "Bottom")]):
        y = 40 - i * 6
        msp.add_line((5, y), (15, y), dxfattribs={"color": 8, "lineweight": 0})
        msp.add_line((10, y - 3), (10, y + 3), dxfattribs={"color": 8, "lineweight": 0})
        msp.add_text(label, height=3, dxfattribs={"layer": "文字", "style": "ROMANS"}).set_placement((10, y), align=align)
    msp.add_text("FIT テキスト", height=3, dxfattribs={"layer": "文字"}).set_placement((40, 40), (100, 40), align=TextEntityAlignment.FIT)
    msp.add_text("ALIGNED", height=3, dxfattribs={"layer": "文字"}).set_placement((40, 32), (70, 32), align=TextEntityAlignment.ALIGNED)
    msp.add_text("%%c50 %%p0.1 45%%d %%uUNDER%%u", height=3, dxfattribs={"style": "GOTHIC"}).set_placement((40, 24))
    msp.add_text("カイテン 30°", height=4, rotation=30, dxfattribs={"width": 0.8, "oblique": 15}).set_placement((45, 6))
    msp.add_text("ミラー", height=4, dxfattribs={"text_generation_flag": 2}).set_placement((95, 10))

    mt = msp.add_mtext("図の{\\C1;文字}は\\P{\\H1.5x;おおきく}、\\Lアンダーライン\\l、{\\fArial|b1|i0|c0|p34;Bold} "
                       "\\S1/2; と \\S+0.1^ -0.2; の数。長い文字はワクの幅で折り返します。",
                       dxfattribs={"char_height": 3, "width": 60, "layer": "文字", "style": "GOTHIC"})
    mt.set_location((112, 44), attachment_point=MTextEntityAlignment.TOP_LEFT)
    mt = msp.add_mtext("CENTERED\\PMTEXT", dxfattribs={"char_height": 3, "bg_fill": 1, "bg_fill_color": 254, "box_fill_scale": 1.5})
    mt.set_location((215, 38), attachment_point=MTextEntityAlignment.MIDDLE_CENTER)

    # --- dimensions and leaders
    doc.dimstyles.duplicate_entry("EZDXF", "JIS").dxf.dimasz = 2.5
    dim = msp.add_linear_dim(base=(240, 22), p1=(230, 10), p2=(290, 10), dimstyle="EZDXF", override=DIM, dxfattribs={"layer": "DIM"})
    dim.render()
    dim = msp.add_radius_dim(center=(270, 110), radius=12, angle=45, dimstyle="EZ_RADIUS", override=DIM, dxfattribs={"layer": "DIM"})
    dim.render()
    msp.add_circle((270, 110), 12, dxfattribs={"layer": "DIM"})
    msp.add_leader([(245, 115), (235, 125), (225, 125)], override={"dimasz": 3}, dxfattribs={"layer": "DIM"})
    ml = msp.add_multileader_mtext("Standard")
    ml.set_content("MULTI\\PLEADER", char_height=2.5)
    ml.add_leader_line(ezdxf.render.mleader.ConnectionSide.left, [Vec2(235, 95)])
    ml.build(insert=Vec2(250, 100))

    # --- blocks
    blk = doc.blocks.new("MARK", base_point=(0, 0))
    blk.add_circle((0, 0), 3, dxfattribs={"color": 0})  # BYBLOCK
    blk.add_line((-4, 0), (4, 0))  # layer 0: takes the insert's layer
    blk.add_line((0, -4), (0, 4), dxfattribs={"color": 1})
    blk.add_attdef("NO", (4, 1), dxfattribs={"height": 2})
    blk.add_attdef("C", (4, -3), text="*", dxfattribs={"height": 1.5, "flags": 2})  # constant: shown
    for i in range(3):
        ref = msp.add_blockref("MARK", (160 + i * 14, 20), dxfattribs={"color": i + 1, "rotation": i * 30, "xscale": 1 + i * 0.2, "yscale": 1 + i * 0.2})
        ref.add_auto_attribs({"NO": "A%d" % (i + 1)})
    msp.add_blockref("MARK", (160, 30), dxfattribs={"xscale": -1, "layer": "外形"})  # mirrored
    arr = msp.add_blockref("MARK", (180, 30), dxfattribs={"color": 4})
    arr.dxf.column_count = 3
    arr.dxf.row_count = 2
    arr.dxf.column_spacing = 8
    arr.dxf.row_spacing = 8
    return doc


def floor_plan(msp):
    """A small plan in metres... in mm at 1:1, 12 m by 8 m."""
    walls = [(0, 0), (12000, 0), (12000, 8000), (0, 8000)]
    msp.add_lwpolyline(walls, close=True, dxfattribs={"layer": "外形", "const_width": 150})
    msp.add_line((5000, 0), (5000, 8000), dxfattribs={"layer": "外形"})
    msp.add_line((5000, 4000), (12000, 4000), dxfattribs={"layer": "外形"})
    msp.add_line((-1000, 4000), (13000, 4000), dxfattribs={"layer": "中心", "ltscale": 50})
    msp.add_line((6000, -1000), (6000, 9000), dxfattribs={"layer": "中心", "ltscale": 50})
    msp.add_circle((2500, 4000), 800, dxfattribs={"layer": "外形"})
    msp.add_text("区間 A", height=300, dxfattribs={"layer": "文字"}).set_placement((2500, 6000), align=TextEntityAlignment.MIDDLE_CENTER)
    msp.add_text("区間 B", height=300, dxfattribs={"layer": "文字"}).set_placement((8500, 6000), align=TextEntityAlignment.MIDDLE_CENTER)
    msp.add_text("FROZEN IN VIEWPORT 2", height=300, dxfattribs={"layer": "FROZEN2"}).set_placement((8500, 2000), align=TextEntityAlignment.MIDDLE_CENTER)
    h = msp.add_hatch(color=8)
    h.paths.add_polyline_path([(5000, 0), (12000, 0), (12000, 4000), (5000, 4000)], is_closed=True)
    h.set_pattern_fill("ANSI32", scale=40)


def layout():
    doc = ezdxf.new("R2018", setup=True, units=4)
    fixed(doc)
    layers(doc)
    frozen2 = doc.layers.add("FROZEN2", color=1)
    floor_plan(doc.modelspace())
    ps = doc.layouts.new("A3 図面")
    doc.layouts.delete("Layout1")  # empty layouts are left out anyway
    ps.page_setup(size=(420, 297), margins=(10, 10, 10, 10), units="mm")
    # a title block
    ps.add_lwpolyline([(0, 0), (400, 0), (400, 277), (0, 277)], close=True, dxfattribs={"lineweight": 70})
    ps.add_lwpolyline([(300, 0), (400, 0), (400, 30), (300, 30)], close=True, dxfattribs={"lineweight": 35})
    ps.add_text("テスト図 1:100 / 1:200", height=4).set_placement((350, 15), align=TextEntityAlignment.MIDDLE_CENTER)
    vp = ps.add_viewport(center=(140, 150), size=(260, 220), view_center_point=(6000, 4000), view_height=11000)
    vp.dxf.status = 2
    vp = ps.add_viewport(center=(340, 150), size=(110, 100), view_center_point=(8500, 4000), view_height=20000)
    vp.dxf.status = 3
    vp.frozen_layers = ["FROZEN2"]
    return doc


def r12():
    doc = ezdxf.new("R12", setup=False)
    doc.encoding = "cp932"
    fixed(doc)
    msp = doc.modelspace()
    doc.layers.add("文字", color=3)
    msp.add_polyline2d([(0, 0), (100, 0), (100, 60), (0, 60)], close=True, dxfattribs={"color": 1})
    pl = msp.add_polyline2d([(10, 10), (40, 30), (70, 10)], format="xyb", dxfattribs={"color": 2})
    pl.vertices[0].dxf.bulge = 0.5
    msp.add_polyline3d([(10, 50, 0), (50, 55, 10), (90, 50, 0)], dxfattribs={"color": 4})
    msp.add_text("シフトJIS の文字", height=5, dxfattribs={"layer": "文字"}).set_placement((50, 35), align=TextEntityAlignment.MIDDLE_CENTER)
    msp.add_circle((80, 20), 8, dxfattribs={"color": 5})
    return doc


def main():
    # fixed dates and GUIDs, so that the files are the same every time
    ezdxf.options.write_fixed_meta_data_for_testing = True
    os.makedirs(OUT, exist_ok=True)
    d = shapes()
    d.saveas(os.path.join(OUT, "shapes.dxf"))
    d.saveas(os.path.join(OUT, "shapes-bin.dxf"), fmt="bin")
    layout().saveas(os.path.join(OUT, "layout.dxf"))
    r12().saveas(os.path.join(OUT, "r12-sjis.dxf"))


if __name__ == "__main__":
    main()
