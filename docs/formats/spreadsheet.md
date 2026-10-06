# Excel, CSV, Parquet

Excel workbooks, CSV/TSV tables and Apache Parquet files all end up as the same thing in BDF: a `sheet` View, the endless-plane layout BDF uses for tabular data ([spec.md §4.1](../spec.md#41-view-の種類)). CSV/TSV and Parquet have no drawing code of their own — `converter/csv` and `converter/parquet` build a plain grid of values and hand it to `converter/xlsx`'s grid layout and drawing code, the same code that lays out and draws worksheet cells. That means a CSV or Parquet file looks and behaves exactly like an Excel sheet in the viewer: a bold, frozen header row, gridlines and frozen panes, and cells you select the way you'd select cells in a spreadsheet (drag, Shift, row/column headers, arrow keys). They copy out as tab-separated values and an HTML table rather than as a picture.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try a sample: [Excel workbook](https://shibukawa.github.io/bdf/viewer/?file=samples/features.xlsx), [a sheet with equation objects](https://shibukawa.github.io/bdf/viewer/?file=samples/math.xlsx), [Shift_JIS TSV file](https://shibukawa.github.io/bdf/viewer/?file=samples/japanese.tsv), or [Parquet table](https://shibukawa.github.io/bdf/viewer/?file=samples/basic.parquet).

## Excel (.xlsx)

`converter/xlsx` reads the OPC zip package (.xlsx, .xlsm, .xltx, .xltm) and turns each worksheet into a `sheet` View. Cell layout, formatting and number-format rendering are all done by the converter itself. BDF has no relayout, so this is the one place that interpretation happens; the sheet is then drawn into Tiles. Cell values are the ones stored in the file (a formula's cached result); formulas themselves are not recalculated.

What's read:
- Number formats (dates, Japanese eras, fractions, accounting formats and the rest of Excel's format-code language), fonts and rich text, fills and borders, alignment (wrapping with Japanese line-breaking rules, overflow into empty neighboring cells, rotation, shrink-to-fit), and merged cells.
- Conditional formatting: color scales, data bars, icon sets, and rules with formulas, evaluated by a small formula evaluator BDF includes for this purpose (arithmetic, cell references and ranges, logic, and common lookup/text/date functions); rules that depend on the current date or time are left unevaluated, with a warning.
- Tables and their built-in or custom styles, and pictures, shapes and charts — each drawn once and reused by every Tile it overlaps. Chart sheets become their own `fixed`-View page rather than a sheet.
- Column widths and row heights, computed by the same rules Excel itself uses (based on the default font's digit width); frozen panes and gridlines are written to the manifest for the viewer to draw.
- Cell-level structure for accessibility: each cell with a value becomes a `MARK CELL`, and cells in a frozen header row or a table's header row are marked as column headers.

Not read: right-to-left sheets (drawn left to right), pivot-table styles (pivot values are drawn as ordinary cells), sparklines, form controls, legacy VML-format shapes, images stored inside a cell, and the old binary .xls format.

Options (`-param`, from `bdf generate -h`):

| Option | Values | Default |
|---|---|---|
| `-param hidden=` | `true` | hidden sheets excluded |

`-hidden` (shorthand for `-param hidden=true`) includes hidden sheets; macro and dialog sheets are always skipped, with a warning. See [design.md §3.6](../design.md#36-excel--bdf-変換器converterxlsxの構造) for the internals, and [design.md §6](../design.md#6-excel-シートの-tile-化) for how a sheet's cells are split into the Tiles a `sheet` View is drawn from.

### Equation objects

An equation placed on a sheet with Insert → Equation (Office Math inside a text box or a shape) is laid out from its Office Math in `a14:m`, as in PowerPoint. A sample is [`converter/xlsx/testdata/math.xlsx`](https://github.com/shibukawa/bdf/blob/main/converter/xlsx/testdata/math.xlsx). Cell formulas (`=SUM(A1:A3)`) are something else: the converter shows the results saved in the file. See [formulas in Word, HTML and Markdown](document.md#formulas) for the layout.

## CSV and TSV

`converter/csv` reads a text table and produces one `sheet` View that looks like the file opened in Excel — because it's drawn by Excel's own converter. This package's job is entirely about the things a CSV file doesn't say explicitly: the character encoding, the dialect, and whether there's a header row. Every guess can be overridden with `-param`, and the guesses actually made are reported in the conversion summary (`Result`, the CLI's last line).

What's guessed:
- **Character encoding**: a byte-order mark if present (UTF-8, UTF-16LE/BE, UTF-32LE/BE), otherwise inferred from the first 64 KiB — byte patterns for BOM-less UTF-16 and ISO-2022-JP escapes, valid UTF-8, and otherwise whichever of Shift_JIS, EUC-JP or Windows-1252 scores highest for "looking like real text." Chinese and Korean legacy encodings (GB18030, Big5, EUC-KR) aren't guessed and need `-param charset=`.
- **Dialect**: the delimiter (comma, tab, semicolon, pipe) and the quoting style (double, single or no quotes; doubled or backslash-escaped), chosen by trying each combination against the first 64 KiB and picking the one that gives the most consistent field count.
- **Header row**: scored column by column — a string sitting above a column of numbers, dates, booleans, emails or URLs scores as a header; a value that also occurs lower in the same column scores as data.

Values are shown exactly as written (unlike Excel, which normalizes `1.50` to `1.5`); numbers with leading zeros (`007`, postal codes) stay strings rather than becoming numbers, and are aligned right the same as real numbers. Column widths follow the content, and cells with line breaks wrap and grow taller. Because it's drawn by the Excel converter, a detected header row is bold, frozen, and marked as a column header for the accessible text layer — exactly as in an Excel sheet with frozen panes.

Options (`-param`, from `bdf generate -h`):

| Option | Values | Default |
|---|---|---|
| `-param charset=` | `utf-8`, `shift_jis`, `euc-jp`, `euc-kr`, `gb18030`, `big5`, `windows-1252` … | detected |
| `-param delimiter=` | comma, tab, semicolon, pipe or another character | detected |
| `-param quote=` | `double`, `single`, `none` | detected |
| `-param header=` | `true`, `false` | guessed |
| `-param table=` | a built-in Excel table style, e.g. `TableStyleMedium2` | none |

See [design.md §3.10](../design.md#310-csvtsv--bdf-変換器convertercsvの構造) for the internals.

## Apache Parquet

`converter/parquet` reads a Parquet file — the column-oriented format pandas, Polars, DuckDB and Spark write, and the format Hugging Face datasets ship in — and produces one `sheet` View of the table, again drawn by `converter/xlsx`. Unlike CSV, there's nothing to guess: Parquet's schema carries the column types, so the only real question is how to print a value.

What's read:
- The reader is written from Parquet's own specification rather than built on Arrow: linking apache/arrow-go's or parquet-go's parquet package would roughly double the size of BDF's Office-format WebAssembly module, so only what's needed to display a table is implemented. Snappy, gzip, Zstandard (via klauspost's implementation — fast, and the default compression for Polars and Databricks) and LZ4 decode everywhere; Brotli decoding is left out of the browser build to keep it small.
- Dictionary pages and data pages (v1 and v2), and every encoding in real-world use: PLAIN, the RLE/bit-packing hybrid, PLAIN/RLE dictionary encoding, the DELTA encodings, and BYTE_STREAM_SPLIT.
- Nested values (lists, maps, structs, and lists of lists) reassembled from each leaf column's repetition and definition levels, including the older, non-standard encodings some older writers used.
- Variant columns (Spark 4, Iceberg v3, DuckDB's variant type), including values stored in shredded form.
- Only the row groups and pages needed for the rows actually shown are read; the file is never loaded into memory as a whole.

Values are printed the way data tools print them: decimals in full precision, timestamps in UTC with as many fractional-second digits as the column needs, UUIDs and intervals in their standard notations, lists/maps/structs as JSON, and geometry columns (`GEOMETRY`, `GEOGRAPHY`, GeoParquet's WKB columns) as WKT. Column names are bold and frozen in row 1, the same as a detected CSV header row, with a gray row of each column's Parquet type (`int64`, `decimal(10, 2)`, `timestamp[ms, UTC]`, `list<string>` and so on) frozen under them in row 2.

By default only the first 10,000 rows are shown (`-param rows=` picks a different number, or `all` up to the sheet limit of 1,048,574 rows). Reading itself is fast — about 0.6 seconds for a million rows and eight columns — but laying the sheet's cells out is not: about 3 seconds and 690 MB for 100,000 rows, about 45 seconds and 2.7 GB for a million, which is why the default keeps small. There's also a 4-million-cell cap regardless of `rows=`, so a table with many columns gets fewer rows; columns are capped at the sheet limit of 16,384, same as Excel and CSV.

Not read: Parquet Modular Encryption (a file with an encrypted footer is detected but not opened), Brotli-compressed columns in the browser build, LZO, and the still-experimental ALP encoding.

Options (`-param`, from `bdf generate -h`):

| Option | Values | Default |
|---|---|---|
| `-param rows=` | a number, or `all` | `10000` |
| `-param types=` | `true`, `false` | `true` |
| `-param table=` | a built-in Excel table style, e.g. `TableStyleMedium2` | none |

See [design.md §3.26](../design.md#326-parquet--bdf-変換器converterparquetの構造) for the internals.
