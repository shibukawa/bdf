"""Writes the test files of the Parquet converter (converter/parquet/testdata)
with pyarrow 25.0.1 and DuckDB 1.5.5 (pip install pyarrow==25.0.1
duckdb==1.5.5):

- basic.parquet: a small table as pyarrow writes it by default (Snappy,
  dictionaries, data pages v1): integers, text in English and Japanese (a
  column named in Japanese, a value with a line break, nulls), a decimal,
  floating-point numbers, booleans, dates, timestamps with milliseconds, a
  list and a struct.
- types.parquet: every scalar type of pyarrow in Zstandard: signed and
  unsigned integers at their limits, float16, float, double (NaN, the
  infinities, tiny and huge numbers), decimals stored as INT32, INT64 and
  FIXED_LEN_BYTE_ARRAY, strings, binaries (text and not), fixed-size
  binaries, UUIDs, JSON, dates, times and timestamps of every unit, with
  and without a time zone, a dictionary-encoded string and a column of
  nulls.
- nested.parquet: lists (empty, null, with null elements), lists of lists,
  of structs and of lists of lists, structs of lists, maps (with null
  values, with struct values), fixed-size lists and large lists.
- encodings.parquet: data pages v2 of 300 rows in small pages, in gzip,
  without dictionaries: DELTA_BINARY_PACKED, DELTA_LENGTH_BYTE_ARRAY,
  DELTA_BYTE_ARRAY, BYTE_STREAM_SPLIT and RLE (booleans) columns, each
  beside the same values in PLAIN.
- codecs.parquet: the same values in a column per codec: uncompressed,
  Snappy, gzip, Brotli, Zstandard and LZ4 (LZ4_RAW).
- legacy.parquet: INT96 timestamps, and lists written with pyarrow's
  non-compliant element names.
- duckdb.parquet: DuckDB's output: Variants (shredded), UUIDs, intervals,
  native GEOMETRY and the deprecated PLAIN_DICTIONARY encoding.
- variant.parquet: Variants that DuckDB shreds into objects of typed
  fields and lists of typed elements, with values that do not fit them.
- geo.parquet: GeoParquet 1.0: WKB geometries (points, lines, polygons,
  multi geometries, a collection, Z coordinates, an empty point) named by
  the file's geo metadata.

The Japanese text keeps to kana and the kanji of the test fonts
(converter/pptx/testdata/fonts), which the testdata lays text out with.
"""
import datetime
import decimal
import json
import math
import os
import struct
import uuid

import duckdb
import pyarrow as pa
import pyarrow.parquet as pq

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(ROOT, "converter", "parquet", "testdata")
UTC = datetime.timezone.utc


def path(name):
    return os.path.join(OUT, name)


def basic():
    rows = [
        (1, "Ann", "本文の説明", "1200.00", 3, 0.125, True, (2026, 1, 5), (2026, 9, 27, 9, 30, 0, 0), ["red", "blue"], (1920, 1080)),
        (2, "Bob", "図表", "35.50", 12, 3.14159, False, (2026, 2, 14), (2026, 9, 27, 10, 15, 30, 250), [], (1280, 720)),
        (3, "さとう", "英語と日本語", "-8.25", 0, -2.5, None, (2026, 3, 1), (2026, 9, 27, 11, 0, 0, 500), None, None),
        (4, "すずき", "一行目\n二行目", "0.99", 7, 1e-07, True, (2026, 4, 30), None, ["green"], (800, 600)),
        (5, None, "空白", None, None, None, None, None, None, None, None),
        (6, "カタカナ", "Mixed テキスト", "99999.99", 42, 2.5e21, False, (2025, 12, 31), (2026, 9, 28, 0, 0, 0, 1), ["a", "b", "c"], (3, 4)),
        (7, "Émile", "Latin-1 à é ü", "10.00", -5, 100.0, True, (2024, 2, 29), (2026, 1, 1, 0, 0, 0, 0), ["x"], (1, 1)),
        (8, "Zoë", "", "0.00", 1000000, 0.0, False, (2000, 1, 1), (1999, 12, 31, 23, 59, 59, 999), [None, "y"], (0, 0)),
    ]
    t = pa.table({
        "id": pa.array([r[0] for r in rows], pa.int64()),
        "name": pa.array([r[1] for r in rows]),
        "説明": pa.array([r[2] for r in rows]),
        "price": pa.array([None if r[3] is None else decimal.Decimal(r[3]) for r in rows], pa.decimal128(10, 2)),
        "qty": pa.array([r[4] for r in rows], pa.int32()),
        "ratio": pa.array([r[5] for r in rows], pa.float64()),
        "active": pa.array([r[6] for r in rows], pa.bool_()),
        "day": pa.array([None if r[7] is None else datetime.date(*r[7]) for r in rows], pa.date32()),
        "updated": pa.array([None if r[8] is None else datetime.datetime(*r[8][:6], r[8][6] * 1000, tzinfo=UTC) for r in rows],
                            pa.timestamp("ms", tz="UTC")),
        "tags": pa.array([r[9] for r in rows], pa.list_(pa.string())),
        "size": pa.array([None if r[10] is None else {"w": r[10][0], "h": r[10][1]} for r in rows],
                         pa.struct([("w", pa.int32()), ("h", pa.int32())])),
    })
    pq.write_table(t, path("basic.parquet"))


def types():
    n = 4
    t = pa.table({
        "i8": pa.array([-128, 0, 127, None], pa.int8()),
        "i16": pa.array([-32768, 1, 32767, None], pa.int16()),
        "i32": pa.array([-2**31, 2, 2**31 - 1, None], pa.int32()),
        "i64": pa.array([-2**63, 3, 2**63 - 1, None], pa.int64()),
        "u8": pa.array([0, 1, 255, None], pa.uint8()),
        "u16": pa.array([0, 1, 65535, None], pa.uint16()),
        "u32": pa.array([0, 1, 2**32 - 1, None], pa.uint32()),
        "u64": pa.array([0, 1, 2**64 - 1, None], pa.uint64()),
        "f16": pa.array([0.1, -65504.0, 6e-08, None], pa.float16()),
        "f32": pa.array([0.1, -3.4028235e38, 1.5, None], pa.float32()),
        "f64": pa.array([math.nan, math.inf, -math.inf, None], pa.float64()),
        "f64b": pa.array([1e-7, 123456789.125, 1e21, None], pa.float64()),
        "dec5": pa.array([decimal.Decimal("-1.23"), decimal.Decimal("0.05"), decimal.Decimal("999.99"), None], pa.decimal128(5, 2)),
        "dec15": pa.array([decimal.Decimal("-123456789012.345"), decimal.Decimal("0.001"), decimal.Decimal("1"), None], pa.decimal128(15, 3)),
        "dec30": pa.array([decimal.Decimal("-12345678901234567890123456.7890"), decimal.Decimal("0.0001"), decimal.Decimal("42"), None],
                          pa.decimal128(30, 4)),
        "str": pa.array(["text", "かな", "", None]),
        "lstr": pa.array(["large", "string", "テキスト", None], pa.large_string()),
        "bin": pa.array([b"plain text", b"\x00\x01\x02\xff", bytes(range(40)), None], pa.binary()),
        "fixed": pa.array([b"ABCD", b"\x00\x00\x00\x01", b"\xde\xad\xbe\xef", None], pa.binary(4)),
        "uuid": pa.array([uuid.UUID("12345678-1234-5678-1234-567812345678").bytes, uuid.UUID(int=0).bytes,
                          uuid.UUID("f24f9b64-81fa-49d1-b74e-8c09a6e31c56").bytes, None], pa.uuid()),
        "json": pa.array(['{"a": 1}', "[1, 2]", '"s"', None], pa.json_()),
        "date": pa.array([datetime.date(1970, 1, 1), datetime.date(1969, 12, 31), datetime.date(9999, 12, 31), None], pa.date32()),
        "time_ms": pa.array([0, 45296789, 86399999, None], pa.time32("ms")),
        "time_us": pa.array([0, 45296789012, 86399999999, None], pa.time64("us")),
        "time_ns": pa.array([0, 45296789012345, 86399999999999, None], pa.time64("ns")),
        "ts_s": pa.array([0, -1, 1790505000, None], pa.timestamp("s")),
        "ts_ms": pa.array([0, -1, 1790505000123, None], pa.timestamp("ms")),
        "ts_us": pa.array([0, -1, 1790505000123456, None], pa.timestamp("us", tz="UTC")),
        "ts_ns": pa.array([0, -1, 1790505000123456789, None], pa.timestamp("ns", tz="Asia/Tokyo")),
        "dict": pa.array(["red", "green", "red", None]).dictionary_encode(),
        "null": pa.nulls(n),
    })
    pq.write_table(t, path("types.parquet"), compression="zstd", store_decimal_as_integer=True)


def nested():
    t = pa.table({
        "list": pa.array([[1, 2], [], None, [None, 3]], pa.list_(pa.int32())),
        "list_list": pa.array([[["a"], ["b", "c"]], [[], None], None, [[None]]], pa.list_(pa.list_(pa.string()))),
        "list_struct": pa.array([[{"k": "x", "v": 1.5}], [{"k": None, "v": None}, None], [], None],
                                pa.list_(pa.struct([("k", pa.string()), ("v", pa.float64())]))),
        "struct_list": pa.array([{"name": "p", "scores": [1, 2]}, {"name": None, "scores": []}, None, {"name": "q", "scores": None}],
                                pa.struct([("name", pa.string()), ("scores", pa.list_(pa.int16()))])),
        "map": pa.array([[("a", 1), ("b", None)], [], None, [("c", 3)]], pa.map_(pa.string(), pa.int64())),
        "map_struct": pa.array([[(1, {"ok": True})], [(2, None)], None, []], pa.map_(pa.int32(), pa.struct([("ok", pa.bool_())]))),
        "fixed_list": pa.array([[1.0, 2.0, 3.0], None, [0.5, -0.5, 0.0], [1e-7, 1e21, 2.0]], pa.list_(pa.float32(), 3)),
        "large_list": pa.array([[10], [20, 30], [], None], pa.large_list(pa.int64())),
        "deep": pa.array([[[[1, 2], [3]], [[4]]], [[[]]], [[None]], None], pa.list_(pa.list_(pa.list_(pa.int8())))),
    })
    pq.write_table(t, path("nested.parquet"))


def encodings():
    n = 300
    ints = [None if i % 17 == 5 else (i * i * 7919) % 100003 - 50000 for i in range(n)]
    bigs = [None if i % 23 == 3 else (i - 150) * 12345678901 for i in range(n)]
    words = [None if i % 19 == 7 else "word-%d-%s" % (i % 37, "x" * (i % 5)) for i in range(n)]
    sorted_words = sorted("prefix/%04d/%s" % (i // 3, "abc"[i % 3]) for i in range(n))
    floats = [None if i % 13 == 4 else i * 0.25 - 10 for i in range(n)]
    bools = [None if i % 11 == 2 else i % 3 == 0 for i in range(n)]
    cols = {
        "delta_i32": (pa.array(ints, pa.int32()), "DELTA_BINARY_PACKED"),
        "delta_i64": (pa.array(bigs, pa.int64()), "DELTA_BINARY_PACKED"),
        "delta_length": (pa.array(words), "DELTA_LENGTH_BYTE_ARRAY"),
        "delta_bytes": (pa.array(sorted_words), "DELTA_BYTE_ARRAY"),
        "split_f32": (pa.array(floats, pa.float32()), "BYTE_STREAM_SPLIT"),
        "split_f64": (pa.array(floats, pa.float64()), "BYTE_STREAM_SPLIT"),
        "split_i32": (pa.array(ints, pa.int32()), "BYTE_STREAM_SPLIT"),
        "rle_bool": (pa.array(bools, pa.bool_()), "RLE"),
    }
    data, enc = {}, {}
    for name, (arr, e) in cols.items():
        data[name], enc[name] = arr, e
        data["plain_" + name], enc["plain_" + name] = arr, "PLAIN"
    pq.write_table(pa.table(data), path("encodings.parquet"), compression="gzip", use_dictionary=False,
                   column_encoding=enc, data_page_version="2.0", data_page_size=256)


def codecs():
    values = ["value %d %s" % (i % 10, "repeated text " * (i % 4)) if i % 9 else None for i in range(200)]
    kinds = ["NONE", "SNAPPY", "GZIP", "BROTLI", "ZSTD", "LZ4"]
    t = pa.table({k.lower(): pa.array(values) for k in kinds})
    pq.write_table(t, path("codecs.parquet"), compression={k.lower(): k for k in kinds})


def legacy():
    t = pa.table({
        "ts": pa.array([datetime.datetime(2024, 1, 1, 20, 34, 56, 123456), datetime.datetime(1900, 1, 1), None], pa.timestamp("us")),
        "list": pa.array([[1, 2], [], None], pa.list_(pa.int32())),
    })
    pq.write_table(t, path("legacy.parquet"), use_deprecated_int96_timestamps=True, use_compliant_nested_type=False)


def duck():
    con = duckdb.connect()
    con.execute("""COPY (
        SELECT * FROM (VALUES
            (1, 42::VARIANT, 'b5f3c2de-0a41-4d4b-9f3e-1c2d3e4f5a6b'::UUID, INTERVAL 3 DAY, 'POINT (1 2)'::GEOMETRY, 'red'),
            (2, 'text'::VARIANT, NULL, INTERVAL '1 year 2 months 3 days 04:05:06.5', 'LINESTRING (0 0, 1 1)'::GEOMETRY, 'blue'),
            (3, {'a': 1, 'b': [true, NULL]}::VARIANT, '00000000-0000-0000-0000-000000000000'::UUID, INTERVAL 0 SECOND, NULL, 'red'),
            (4, NULL, NULL, NULL, 'POLYGON ((0 0, 4 0, 4 4, 0 0))'::GEOMETRY, NULL),
            (5, [1::VARIANT, 'two'::VARIANT, 3.5::VARIANT]::VARIANT, NULL, INTERVAL 90 MINUTE, 'POINT EMPTY'::GEOMETRY, 'blue')
        ) AS t(id, v, u, iv, g, color) ORDER BY id
    ) TO '%s' (FORMAT parquet, GEOPARQUET_VERSION 'NONE')""" % path("duckdb.parquet"))


def variant():
    con = duckdb.connect()
    con.execute("""COPY (
        SELECT * FROM (VALUES
            (1, {'a': 1, 'b': 'x'}::VARIANT, [1, 2]::VARIANT),
            (2, {'a': 2, 'b': 'y', 'c': true}::VARIANT, [3]::VARIANT),
            (3, {'a': NULL, 'b': 'z'}::VARIANT, []::INTEGER[]::VARIANT),
            (4, 5::VARIANT, 'not a list'::VARIANT),
            (5, {'a': 'text', 'b': {'deep': [1.5]}}::VARIANT, [4, NULL]::VARIANT),
            (6, NULL, NULL)
        ) AS t(id, obj, arr) ORDER BY id
    ) TO '%s' (FORMAT parquet)""" % path("variant.parquet"))


def wkb_point(x, y, z=None):
    if z is None:
        return struct.pack("<BIdd", 1, 1, x, y)
    return struct.pack("<BIddd", 1, 1001, x, y, z)


def wkb_line(pts):
    return struct.pack("<BII", 1, 2, len(pts)) + b"".join(struct.pack("<dd", *p) for p in pts)


def wkb_polygon(rings):
    out = struct.pack("<BII", 1, 3, len(rings))
    for r in rings:
        out += struct.pack("<I", len(r)) + b"".join(struct.pack("<dd", *p) for p in r)
    return out


def geo():
    geoms = [
        wkb_point(139.767, 35.681),
        wkb_line([(0, 0), (1.5, 2), (3, 0)]),
        wkb_polygon([[(0, 0), (10, 0), (10, 10), (0, 10), (0, 0)], [(2, 2), (3, 2), (3, 3), (2, 2)]]),
        struct.pack("<BII", 1, 4, 2) + wkb_point(1, 2) + wkb_point(3, 4),
        struct.pack("<BII", 1, 6, 1) + wkb_polygon([[(0, 0), (1, 0), (0, 1), (0, 0)]]),
        struct.pack("<BII", 1, 7, 2) + wkb_point(5, 6) + wkb_line([(7, 8), (9, 10)]),
        wkb_point(1, 2, 3),
        struct.pack(">BIdd", 0, 1, math.nan, math.nan),
        None,
    ]
    meta = {"version": "1.0.0", "primary_column": "geometry",
            "columns": {"geometry": {"encoding": "WKB", "geometry_types": []}}}
    t = pa.table({"name": pa.array(["point", "line", "polygon", "multipoint", "multipolygon", "collection", "point z", "empty", "null"]),
                  "geometry": pa.array(geoms, pa.binary())})
    t = t.replace_schema_metadata({"geo": json.dumps(meta)})
    pq.write_table(t, path("geo.parquet"))


os.makedirs(OUT, exist_ok=True)
basic()
types()
nested()
encodings()
codecs()
legacy()
duck()
variant()
geo()
