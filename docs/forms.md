# Making forms

Forms such as invoices and delivery slips can be made as BDF documents from a template and data. A template is an Excel workbook, a Word document or a PowerPoint presentation whose text holds tags (`{{name}}`), a PDF with form fields, a template of pdfme, or a text file of the lines of a form with tags. The data is JSON, YAML, CSV or Excel files; `bdf fill` replaces the tags with the data (or fills the fields), lays the result out, and writes a `.bdf`. No Office file is written. `bdf pdf` writes the form as a PDF file: plain, tagged, or PDF/A.

```sh
bdf fill -data order.yaml -data items=lines.csv invoice.xlsx invoice.bdf
bdf pdf invoice.bdf invoice.pdf
```

| Template | The form | Good for |
|---|---|---|
| Excel (`.xlsx`) | The pages the workbook prints. More rows make more pages | Forms that are tables, forms that calculate with formulas, slips of a fixed number of rows |
| Word (`.docx`) | The pages of the document. What gets longer flows onto the next page | Forms of text and tables, notices, mail merge letters |
| PowerPoint (`.pptx`) | The slides. Shapes stay where they are | Forms with fixed places, certificates, labels, a sheet for each record |
| [PDF form](#pdf-forms) (`.pdf`) | The pages of the PDF, with the values in the places of its fields | Forms that exist as PDF files already: applications, official forms |
| [pdfme](#templates-of-pdfme) (`.json`) | The pages its designer shows; tables grow over pages | Forms designed in the browser with the designer of pdfme |
| [Text](#templates-of-text) (`.txt`) | Characters on a grid of columns and lines, as printers print in their text mode | Continuous forms and carbon-copy slips of dot matrix printers, receipts, fixed-length text for other systems |

The tags and the data are the same for the three kinds of Office files and for templates of text: the same data fills templates of different formats. A PDF form and a template of pdfme have fields instead of tags, and take the same kind of data by the names of their fields.

To try it without installing anything, drop a template on the [form page](https://shibukawa.github.io/bdf/form/): it reads the tags, makes a form of them, shows the template filled with what you type, and saves it as a PDF file. All of it runs in the browser; nothing is uploaded.

## Writing a template

Make the form as usual and write tags where values go. Excel workbooks come first here. Word and PowerPoint differ only in what the same tags repeat ([Word templates](#word-templates), [PowerPoint templates](#powerpoint-templates)).

In Excel, the paper, margins, print area, print titles, page breaks, header and footer are those of the workbook's page setup ([the print view of Excel](formats/spreadsheet.md#printed-pages-the-print-view)).

| Tag | Meaning |
|---|---|
| `{{name}}` | A value. `{{customer.name}}` is a value of a value, `{{.}}` the item itself |
| `{{total \| fmt:"#,##0"}}` | A value in a number format: the format codes of Excel (`yyyy-mm-dd`, `0.0%` …) |
| `{{#items}} … {{/items}}` | A section: once for each item of a list, or when the value is true |
| `{{^items}} … {{/items}}` | An inverted section: when the list is empty or the value is false |
| `{{%photo}}` | A picture, fitted in the cell (or its merged range) in Excel, and in the frame of a picture of the template in Word and PowerPoint. `{{%photo \| fit:cover}}` covers the frame, `fit:stretch` stretches |
| `{{%url \| qr}}`, `{{%jan \| barcode:jan}}` | The QR code or the barcode of a value ([QR codes and barcodes](#qr-codes-and-barcodes)) |
| `{{%sign \| signature}}`, `{{%seal \| seal}}` | A picture, as `{{%photo}}` is, that a person gives in a way of its own: a signature drawn by hand, the picture of a seal ([Signatures and seals](#signatures-and-seals)) |

- The braces, marks and quotes may be the full-width or curly ones that input methods and Word put in their place (`｛｛`, `＃`, `“ ”`).
- `\{{` is `{{` as text, not a tag.
- **A cell that is one tag and nothing else** takes the value as it is: a number is a number (shown in the cell's number format, read by formulas), a date a date. The text `"2026-10-13"` is read as a date only where the number format of the cell or the tag is a date format.
- A tag written with other text is replaced in the text (`To: {{customer}}`).
- Tags look names up; they do not calculate. Sums, taxes, what is carried over and their rounding are the data's to give (the formulas of the workbook aside: they are calculated again, see below).

### Repeating rows

A section whose opening and closing tags are in different cells repeats rows. The cells of the repeated rows see the item.

**Adding rows** (both tags in the same row): a row is added for each item, the rows below move down, and the pages grow when they do not fit. The rows of the print titles repeat on every page.

| A | B | C | D | E |
|---|---|---|---|---|
| `{{#items}}{{_index}}` | `{{name}}` | `{{qty}}` | `{{price}}{{/items}}` | `=C4*D4` |
| | | | Total | `=SUM(E4:E4)` |

**Filling the rows that are there** (the tags in different rows): the rows from that of the opening tag to that of the closing tag are slots, a row for an item. No row is added; the rows without an item stay as empty rows with their borders. The tags and formulas of the first row are those of the other rows that have no tags of their own. More items than rows is an error.

**Repeating as pages** (`| pages`): the rows between the tags are copied for each item, each copy on a new page. A slip of ten rows that takes three sheets for 25 rows of data is this. The data gives the rows of each page and what is carried over:

```json
{ "pages": [
    { "rows": [ … 10 rows … ], "carry_in": 0,      "carry_out": 55000 },
    { "rows": [ … 10 rows … ], "carry_in": 55000,  "carry_out": 210000 },
    { "rows": [ … 4 rows … ],  "carry_in": 210000, "carry_out": 300000 } ],
  "total": 300000 }
```

Put the slots `{{#rows}}` … `{{/rows}}` inside `{{#pages | pages}}` … `{{/pages}}`. The page number is `{{_page}} / {{_pages}}`, and a total on the last page only `{{#_last}}Total {{total}}{{/_last}}`.

`{{#groups | rows}}` repeats rows between tags in different rows by adding rows (a heading, details and a subtotal for each group). Sections nest.

A section whose two tags are in **the same cell** repeats text in the cell (`{{#tags}}{{.}} {{/tags}}`). With `| rows` it repeats the row all the same, so a table of one column can grow too (`{{#items | rows}}{{name}}{{/items}}`). That a section which says how it repeats repeats rows (paragraphs in Word, slides in PowerPoint) rather than text is a rule of all three formats. An inverted section may say it too: `{{^items | rows}}no items{{/items}}` is a row that is there only without items.

### Variables in sections

| Variable | Meaning |
|---|---|
| `_index`, `_count` | The number of the item (from 1), and the number of items |
| `_first`, `_last` | Whether the item is the first, the last |
| `_page`, `_pages` | In a `| pages` section: the page, and the number of pages |
| `_blank`, `_blank_first` | In a slot: whether the row has no item, and whether it is the first such row |
| `_slide`, `_slides` | In PowerPoint: the number of the slide, and the number of slides of the form |

### What follows the rows

When rows are added or copied, what refers to rows follows them as it does when Excel inserts rows: merged cells, conditional formats, links, the print area and titles, page breaks, drawings, and the references of formulas. A formula in a copied row refers to its own row (its own page), and a range that holds a repeated row holds all its copies (`=SUM(E4:E4)` is the sum of all the detail rows).

The formulas are calculated again after the template is filled: those with the functions the evaluator of conditional formats knows (arithmetic, `SUM`, `AVERAGE`, `COUNT`, `COUNTIF`, `IF`, `ROUND`, text and date functions …). The cell of any other formula is left empty, with a warning.

### Excel table headings and sections across pages

For an invoice, set Excel's **Rows to repeat at top** to the table heading
alone (for example `$15:$15`), leaving the title and customer outside it.
By default Excel print titles repeat on all later pages. Two optional
sheet-scoped names in **Formulas → Name Manager** refine BDF pagination:

| Name | Refers to | Behavior |
|---|---|---|
| `_bdf.Print_TitleScope` | `'Invoice'!$A$15:$H$16` | Repeat the print-title rows only while this table continues. Include the heading and the tagged detail row; the range grows with the items. A page containing only totals or notes has no table heading. |
| `_bdf.Print_KeepTogether` | `'Invoice'!$F$18:$H$20,'Invoice'!$A$22:$H$25,'Invoice'!$A$27:$H$31` | Keep each range's rows on one page, moving the section to the next page when needed. Use separate ranges for totals, payment details and notes. |

The ranges follow inserted template rows. A keep-together range inside a
repeated item applies to each copy. Overlapping ranges form one group;
adjacent ranges remain separate. A group taller than a page can split,
and explicit manual page breaks take precedence. These names control BDF
output; Excel's own print preview continues to use Excel's print rules.
The public Excel invoice samples include both names.

## Word templates

In a Word document, tags go in the text of the body, of table cells, of text boxes and of headers and footers. The document is laid out again after it is filled: a longer value makes more lines, and more rows of a table flow onto the next page. The header rows of a table ("Repeat header rows") are repeated on the pages it flows onto. The form has the page view alone.

A section repeats the smallest of these that holds both its tags.

| Where the tags are | What repeats |
|---|---|
| In the same paragraph | Text of the paragraph |
| In different paragraphs of the body (or of a cell, or of a text box) | The paragraphs and tables from the one to the other |
| In different cells of a row of a table | The row (a row for each item) |
| In different rows of a table | The rows from the one to the other are slots, a row for an item (`\| rows` repeats them by adding rows) |

- A paragraph that holds nothing but such tags (`{{#items}}` or `{{/items}}` alone) is not in the form.
- A section with `| rows` repeats the row of the table it is in or, outside tables, its paragraph, even when both tags are in one paragraph: `{{#items | rows}}{{name}}{{/items}}` in a paragraph of a list makes a list item for each item.
- A section with `| pages` repeats the paragraphs and tables from the one tag to the other, each item after the first on a new page: a one-page form for each person, a slip of ten rows for each page of the data.
- A picture is replaced where the alternative text of a picture of the document holds `{{%photo}}`. The rest of the alternative text is that of the new picture.
- Word stores a tag in pieces where it checked the spelling or the text was edited (`{{custo` + `mer}}`). The text of a paragraph is read as one, so this works as it is. A value has the format of the text where its tag starts.
- Headers, footers and text boxes outside sections see the outermost values of the data; a text box inside a section sees the values of its item. Footnotes and endnotes are not filled.
- Page numbers are the fields of Word (`PAGE`, `NUMPAGES`).

## PowerPoint templates

In a presentation, tags go in the text of shapes, text boxes and table cells. A slide is a fixed page, and shapes stay where they are.

| Where the tags are | What repeats |
|---|---|
| In the same paragraph | Text of the paragraph |
| In different paragraphs of a shape (or of a cell) | The paragraphs from the one to the other |
| In different cells of a row of a table | The row (a row for each item) |
| In different rows of a table | The rows from the one to the other are slots (`\| rows` repeats them by adding rows) |
| In different shapes of a slide | The slide (a slide for each item) |
| On different slides | The slides from the one to the other |

- A section with `| rows` repeats the row of the table it is in, or else its paragraph, even when both tags are in one paragraph. A section with `| pages` repeats the slide (`{{^items | pages}}` is a slide that is there only without items).
- **Text is fitted again** after it is filled, by the rules of PowerPoint: a shape that shrinks text on overflow shrinks it, and a shape that resizes to fit its text gets its new height. Shapes do not push each other.
- **A table that grew** past the foot of its slide goes on on copies of the slide. The rows above the repeated ones (the headings) are on every copy, those below them (totals) on the last. The table may reach down to the shape under it, or to the margin at the foot of the slide.
- `{{_slide}} / {{_slides}}` are the number of the slide and the number of slides, with the slides that sections and tables added.
- A picture is replaced where the alternative text of a picture of the slide holds `{{%photo}}`.
- PowerPoint stores a tag in pieces where the language changes; the text of a paragraph is read as one.

## PDF forms

A PDF with form fields (an AcroForm, as Acrobat, LibreOffice and most form tools make) is a template as it is. Its fields are the names of the data; `bdf fill` draws the values in the places of the fields and writes the pages of the file as a BDF document. The filled form is flat: its fields are fields no more.

```sh
bdf template schema -json-schema application.pdf   # the fields, as the JSON Schema of the data
bdf fill -data applicant.json application.pdf application.bdf
bdf pdf application.bdf filled.pdf
```

| Field | The data |
|---|---|
| Text | A text or a number. A field of several lines breaks the text into its lines, a comb field puts a character in each cell, and a field without a size gets the size that fits it |
| Checkbox | `true` or `false` |
| Radio buttons | The value of the button that is on |
| List, dropdown | The value of a choice. The field shows what a reader sees of it |
| Signature | A picture: a signature drawn by hand, drawn whole in the middle of the box of the field ([Signatures and seals](#signatures-and-seals)) |

- The name of a field is the name the data has its value by. A field `applicant.name` is `name` of the object `applicant` (the name as it is, with its dot, works too).
- The schema has what the file says of each field: the choices (`enum`, with what a reader sees of them in `x-bdf.labels`), the most characters (`maxLength`), the value the file has (`default`), and the tooltip of the field as its title.
- The text is drawn in the size, the color and the alignment the field says, in a font of the machine that has its characters. The fonts of a form are often Latin only, so Japanese text does not depend on them.
- A field that the data has no value for stays as the file has it, with a warning (an error with `-strict`). Buttons to press are not filled.
- The borders and backgrounds of the fields stay. A checkbox or a radio button shows the appearance the file has for its state.

## Templates of pdfme

[pdfme](https://pdfme.com/) is an open-source PDF generator whose designer runs in the browser. The template its designer saves (a JSON file: a base PDF or a blank page, and fields with places in millimetres) is a template here too, and fills as the generator of pdfme fills it.

```sh
bdf fill -data invoice.json template.json invoice.bdf
bdf template compile template.json invoice.bdft   # the pages show what the designer left in the fields
```

| Field of pdfme | The data |
|---|---|
| `text`, `select` | A text or a number (`select`: one of its options) |
| `multiVariableText` | An object of its variables: `"customer": {"company": "…", "person": "…"}` |
| `date`, `time`, `dateTime` | A date (`2026-10-14`, `2026/10/14`), a time (`09:30`), or both; drawn in the format of the field (English and Japanese names) |
| `checkbox`, `radioGroup`, `circleMark` | `true` or `false` |
| `image`, `signature`, `svg` | The path of a file, or a `data:` URL (`svg`: its markup too) |
| `qrcode`, `ean13`, `ean8`, `code39`, `code128`, `nw7`, `itf14`, `upca` | The text of the code |
| `table` | A list of rows: objects by the names of the columns, or lists of cells |
| `list` | A list of texts (a tab at the start of one puts it a level deeper) |

- The data is the values of JSON, not their texts: pdfme takes the rows of a table as the text of a JSON array, here they are a list.
- **On a blank base page**, tables, lists and texts that expand (`overflow: expand`) grow with the data, move what is under them, and go on over further pages; a table says the names of its columns again when the template asks (`repeatHead`), and the static fields of the base page are on every page. **Over a base PDF** the fields stay where they are, on the page of the base of their number.
- A read-only text may hold expressions (`{price * 1.1}`, `{currentPage} / {totalPages}`): arithmetic, comparisons, `? :`, the methods of texts, numbers and lists, `Math` and `JSON`. They run no other code. The names they read of the data are fields of the schema; `date` and `dateTime` (which pdfme reads off the clock) are among them, so the same data makes the same form.
- The fonts the template names are looked up among the fonts of the machine (`-font-dir`), and Roboto, the font of pdfme, where it names none. A font the machine lacks is replaced by another.
- What its designer left in the fields is the sample of the template: `bdf template compile` without `-sample` shows it, and the form page starts from it.
- Not drawn as pdfme draws it: the inline markdown of text, and the codes `japanpost`, `upce`, `gs1datamatrix` and `pdf417` (left out, with a warning). What a barcode says is a plain line under its bars. A QR code holds the same text but may pick another mask. An SVG image fits whole in its box (the generator of pdfme crops some).
- A base PDF that the template names by a URL is not fetched: it is in the template, as the designer saves it.

## Templates of text

A text file (UTF-8) is a template of the forms that a dot matrix printer prints on continuous paper and carbon-copy slips, and a receipt printer on its roll: characters on a grid of columns and lines. What the editor shows is the form: a line of the file is a line of the page, a character takes a column, and a full-width one (kanji, kana) two. The same tags and data as the other templates fill it.

```text
---
cpi: 10                  # columns to the inch
lpi: 6                   # lines to the inch
columns: 80              # the width of the paper in columns
lines: 33                # the length of a page in lines (0: a roll, as long as the form)
margin: {left: 2, top: 1}
skip: 2                  # lines left empty over the perforation
---
                  {{#_big}}納品書{{/_big}}                {{date | fmt:"yyyy年m月d日"}}
{{#preprinted}}得意先{{/preprinted}} {{customer | width:30}}{{#preprinted}}御中{{/preprinted}}  No. {{no | right:8}}
{{#preprinted}}
┌──────────────────────┬──────┬────────────┐
│ 品名                 │ 数量 │ 金額       │
{{/preprinted}}
{{#rows}}
│ {{name | width:20 | cut}} │ {{qty | right:4}} │ {{amount | fmt:"#,##0" | right:10}} │
│                      │      │            │
│                      │      │            │
{{/rows}}
```

```sh
bdf fill -data order.json slip.txt slip.bdf    # the form as a document: to view, to print, as a PDF
bdf fill -data order.json slip.txt slip.txt    # the form as plain text: its lines, the pages separated by form feeds
bdf template compile -sample order.json slip.txt slip.bdft
```

**The front matter** (YAML between lines `---`, optional) gives the grid. Without it a template is 80 columns by 66 lines at 10 CPI and 6 LPI (fanfold paper of 8 by 11 inches).

| Key | Meaning |
|---|---|
| `cpi`, `lpi` | Columns and lines to the inch: 10, 12, 15, 17.1 or 20 CPI, 6 or 8 LPI (the pitches of ESC/P; another number is warned of) |
| `columns`, `lines` | The paper, in columns and lines, the margins included. `lines: 0` is a roll: a page is as long as its lines (a receipt) |
| `margin: {left, top}` | Columns before each line, and lines above the first |
| `skip` | Lines left empty at the bottom of each page, over the perforation |
| `title` | The first lines of the form repeat at the top of each page (the title of a list) |
| `font` | The monospace font of the pages: a family (`BIZ UDGothic`), or a font file beside the template (`fonts/slip.ttf`; `-files` names the directory) |
| `ambiguous` | The columns of the characters of ambiguous width (box drawing, `○`, `※`, Greek): `narrow` (1, the default, as most editors show them) or `wide` (2) |

**Values** take the columns of their text: a longer value moves what follows it. A value with columns takes exactly them, counted with full-width characters as two: `{{name | width:20}}` on the left of 20 columns, `{{no | right:8}}` on the right, `{{title | center:30}}` in the middle. A value wider than its columns is written whole with a warning; `| cut` cuts it, and `| wrap` puts the rest on the lines under it, in the same columns (`{{address | width:30 | wrap}}`). A line break in a value goes on to the next line in the same way.

**Lines of tags.** A line with nothing but the tags of sections of lines takes no line of the form. A form feed in the file breaks the page.

**Sections** repeat as in a workbook ([Repeating rows](#repeating-rows)), with the lines of the file as its rows: a section in one line repeats characters (`{{#tags}}{{.}} {{/tags}}`); a section over lines fills the lines between its tags, an item to a line (slots: lines without items stay as they are without values, and a line of the slots without tags of its own is the first line again); `{{#items | rows}}` repeats the lines for each item, and the pages break where they are full, with the `title` lines at the top of each; `{{#pages | pages}}` repeats them on a page of its own for each item. A form without a section of `| rows` must fit its pages (a value that wraps pushes the lines under it down; past the end of the page it is an error). The variables of sections are those of the other templates.

**Three sections are not of the data.**

| Section | Meaning |
|---|---|
| `{{#preprinted}} … {{/preprinted}}` | The rules and names that the paper has printed already. They are drawn on a layer of the role `background`: printing on that paper leaves them out (`SkipRoles` of the printer's settings). Over lines or within a line |
| `{{#_wide}} … {{/_wide}}` | Characters of double width: a half-width one takes two columns, a full-width one four |
| `{{#_tall}} … {{/_tall}}`, `{{#_big}} … {{/_big}}` | Characters of double height, or of double width and height. A line with them takes two lines of the page, its characters standing on the base line of the lower one |

**The form.** A fixed view of pages of the size of the paper (`columns / cpi` by `lines / lpi` inches). Each character is in its columns, drawn with a monospace font: BIZ UDGothic (Google Fonts, SIL Open Font License), then IPAGothic, Noto Sans Mono CJK JP and the monospace fonts of the system; whatever the font, the characters keep to the grid. Each line is a paragraph, and the lines of a section of rows or slots a table, for search, copy and tagged PDF.

**The print image.** The document also holds the grid and the characters of each line, as a part of its own (`t: "text"`, [spec §4](spec.md#4-manifest); its JSON is described by the Go package `form/textpart`): which ranges are of the preprinted paper and which are of double size. Viewers pass it over. The text modes of printers (ESC/P for dot matrix printers, ESC/POS for receipt printers) print from it, and so does `bdf fill … out.txt`:

| `-param` | Meaning |
|---|---|
| `encoding=shift_jis` | The encoding of the text: `utf-8` (the default), `shift_jis` (Windows-31J), `euc-jp`, `iso-2022-jp`, or another name of the WHATWG Encoding Standard. A character that it does not have is written as `〓`, with a warning |
| `newline=crlf` | Lines end with CR LF instead of LF |
| `paper=preprinted` | The ranges of the preprinted paper are written as spaces |

**Printing.** `bdf print` sends the form to a dot-matrix printer of ESC/P in its text mode: the printer's pitch, line spacing, page length and perforation skip are those of the template, the characters are printed with its own fonts (Kanji in its Kanji mode), and `-skip background` leaves out the ranges of the preprinted paper. A receipt printer of ESC/POS prints it so with `-print-mode text`:

```bash
bdf print -printer 'socket://192.168.0.20:9100?profile=vp-f2000' -skip background slip.bdf
```

```bash
bdf print -printer 'serial:///dev/ttyUSB0?profile=tm-t88' -print-mode text receipt.bdf
```

ESC/P prints the characters of one column of a template of `ambiguous: narrow` (its box drawing, `○`, `※`) as Kanji of half width. A receipt printer has no such Kanji: its box drawing and `×` come from its tables of a byte, and the others are printed as `?` with a warning, so write `ambiguous: wide` for a receipt with `○` or `※`. ESC/P has no Kanji of double height alone (`_tall`): they are printed at their size, with a warning. ESC/P takes the pitches of 10, 12, 15, 17.1 and 20 CPI (Kanji at 10 to 15), ESC/POS those whose column is 12 dots or more (15 CPI at 180 and 203 dpi). A form the printer cannot print so is printed as dots, with a warning (with `-print-mode text`, it is an error).

A template of text has no pictures, QR codes or barcodes (the text modes of printers print characters). Characters of the private use area and pictographs (emoji) are warned of: printers do not have them.

## QR codes and barcodes

A picture tag with `qr` or `barcode:kind` draws the code of the value in the place of a picture. The place is that of a picture: in Excel the cell of the tag (or its merged range), in Word and PowerPoint the frame of the picture whose alternative text holds the tag.

| Tag | What is drawn |
|---|---|
| `{{%url \| qr}}` | A QR code (model 2, versions 1 to 40) with the error correction level M. `qr:L`, `qr:Q` and `qr:H` set another |
| `{{%jan \| barcode:jan}}` | JAN / EAN-13. A value of 12 digits gets its check digit, that of a value of 13 is checked. `barcode:ean13` is the same |
| `barcode:ean8`, `barcode:upca` | JAN / EAN-8 (7 or 8 digits), UPC-A (11 or 12 digits) |
| `barcode:code39` | Code 39: digits, capital letters, the space and `- . $ / + %`. The `*` at the ends are added |
| `barcode:code128` | Code 128: the characters of ASCII. The code sets A, B and C are chosen |
| `barcode:gs1-128` | GS1-128. The value has its application identifiers in parentheses: `(01)04901234567894(10)ABC` |
| `barcode:itf`, `barcode:itf14` | ITF (digits; an odd number of them gets a leading 0), ITF-14 (13 digits get their check digit) |
| `barcode:nw7` | NW-7 (Codabar): digits and `- $ : / . +`. A value without start and stop characters (`A` to `D`) is put between two `A` |

- A code is drawn as black vector shapes: it stays sharp at any size, and in a PDF that is printed.
- A QR code stays a square, in the middle of its place. A barcode fits its place, including its digits when it draws them: size the place for the reader (a JAN code is 37.29 mm wide at its nominal size, and may be 0.8 to 2 times that).
- A code comes with the margin readers need (the quiet zone: 4 modules around a QR code, 10 at the sides of a barcode). `quiet:0` leaves it out when the place has a margin around it, and `quiet:2` sets a number of modules.
- The value is a text or a number. Pass a code that starts with zeros (`0490…`) as a text in the data.
- Nothing is drawn for no value (empty, or `null`). A value the code cannot hold (a JAN code of another length or with a wrong check digit, a character the code lacks, a text too long for a QR code) is left out with a warning, and is an error with `-strict`.
- A code is a figure with an alternative text, for assistive technology and for search. In Excel, and for a picture whose alternative text is the tag alone, it is what the code says.
- JAN / EAN-13, EAN-8 and UPC-A draw their digits below the bars, including the check digit. The digits are vector outlines built into `image/barcode`, so they need no font or separate value tag.

The same drawing is available without a form:

```go
symbol, err := barcode.Encode("jan", "490123456789", nil) // adds the check digit: 4901234567894
// Check err, then draw the bars and digits together.
drawing := symbol.Object(113, 36)
```

Pass `&barcode.Options{NoText: true}` to draw only the bars. `Symbol.At` and `Symbol.Rects` always describe the encoded modules alone.

## Writing a PDF

`bdf pdf` writes a BDF document as a PDF file: a form, or any other document.

```sh
bdf pdf invoice.bdf invoice.pdf
bdf pdf -skip background -skip master slip.bdf slip.pdf   # to print on paper that has them already
bdf pdf -tagged report.bdf report.pdf                      # a tagged PDF
bdf pdf -pdfa 2a invoice.bdf invoice.pdf                   # a PDF/A file, to keep
```

- Text stays text. The fonts of the document are embedded with the glyphs in use, so the text of the file is searched and copied. Fonts referred to by name are embedded from the fonts of the machine (`-font-dir`). Text in a font that may not be embedded is drawn as the outlines of its glyphs, with the text they stand for.
- Shapes, lines, gradients, transparency, blend modes and soft masks stay vectors; JPEG images are stored as they are, other bitmaps as pixels, SVG images as pictures of 300 dpi. Links and the document information (title, author, dates) are written.
- The views of pages are written (the print views of Excel, the pages of Word, slides). A document without pages (a sheet, a scroll view) becomes one page. `-view` chooses views.
- Shadows (`SHADOW`) have no counterpart in a PDF: the shadow alone becomes a translucent picture of 300 dpi under its shape. The shape stays vectors.
- `FILTER` and the composite operations of Canvas that are not blend modes (`destination-out` …) have no counterpart in a PDF: they are left out or drawn over, with a warning. Encryption is not written yet.
- The same document makes the same bytes.

### Tagged PDF and PDF/A

A **tagged PDF** (`-tagged`) holds the structure of the document beside its pages: which text is a heading, a paragraph, an item of a list or a cell of a table, in which order it is read, what a figure shows. Screen readers read a file by it, and tools that take text out of a PDF keep paragraphs and tables whole.

- The structure is the one BDF has (the `MARK` instructions the converters write): headings with their levels, paragraphs, lists, tables with their header cells and merged cells, figures with their alternative text, links, and the language of the text.
- Page headers and footers, the masters of slides, backgrounds, the gridlines of a sheet and shadows are marked as artifacts: a reader skips them.

**PDF/A** (`-pdfa`) is the PDF for keeping: a file that shows the same in ten years, with nothing that depends on the machine that opens it.

| Level | What it promises |
|---|---|
| `2b` | It looks the same anywhere: every font is embedded, the colors are sRGB by an output intent the file carries, and its metadata says what it is |
| `2u` | And its text is Unicode: search and copy give the characters |
| `2a` | And it is tagged (`-tagged` is implied) |
| `3b`, `3u`, `3a` | The same of part 3 of the standard |

- What PDF/A does not allow is written another way: images are not smoothed by the viewer, and a character that no font has is left out with a warning instead of drawn as an empty box.
- At level `a`, characters of the private use areas (icon fonts) get a warning: the standard asks for a text that says what they are.
- The files are checked with [veraPDF](https://verapdf.org/) while BDF is developed. A file is not claimed to be PDF/UA: that also needs a title, a language and headings in order, which are the document's. PDF/A-1 (which has no transparency) is not written.

## Data

The data is one object of the values of JSON, put together with `-data`:

| Flag | Read as |
|---|---|
| `-data order.json`, `-data order.yaml` | An object; a later file replaces the values of the same names |
| `-data items=lines.csv` | A CSV or TSV file as a list of objects named by its first record, under the name. The character encoding and the delimiter are detected (Shift_JIS too) |
| `-data book.xlsx`, `-data items=book.xlsx` | A list for each sheet of a workbook, by the names of the sheets; with a name, a workbook of one sheet is that list |
| `-data customer.address=addr.json` | A JSON or YAML value under the name |

A value of a CSV file is a number only when it is written as JSON writes numbers (`007` and `1,234` stay text); `true` and `false` are booleans, and an empty value is none. A picture is the path of a file (in the directory of `-files`, by default that of the template) or a `data:` URL. The path stays in that directory: `..` and links that lead out of it are not followed.

A tag whose name the data does not have is left empty, with a warning; `-strict` makes it an error.

`-each` makes a form for each record; the output path names each with tags:

```sh
bdf fill -each orders.csv -data items=lines.csv invoice.xlsx 'out/{{customer}}-{{no}}.bdf'
```

For one document of all the records, use `{{#orders | pages}}` in the template instead.

## Compiled templates (.bdft)

`bdf template compile` checks a template (a tag that is not closed, a section across rows that cannot be filled: each told with its cell, its paragraph or its shape) and writes a compiled template. It is a BDF document: its pages show the template filled with sample data (or with its tags, without a sample), and it holds the template itself and the names its tags look up. It opens in the viewer and as a thumbnail like any document, and `bdf fill` takes it in place of the template.

```sh
bdf template compile -sample sample.yaml invoice.xlsx invoice.bdft
bdf template schema invoice.bdft        # the names the tags look up, as JSON
bdf fill -data order.yaml invoice.bdft out.bdf
```

## Partly filled forms (.bdff)

A form may be filled by more than one person: the sender fills in their part of a contract and signs it, and the other party fills in theirs. `-open` names what is left for later; `bdf fill` then writes a partly filled form (.bdff) of the rest:

```sh
bdf fill -data mine.yaml -open party_b contract.bdft contract.bdff   # party A's part
bdf template schema -json-schema contract.bdff                       # what is left: party_b
bdf fill -data theirs.json contract.bdff contract.bdf                # party B completes it
```

- A partly filled form is a compiled template that holds the values given so far. Its pages show the template filled with them, with what is left open empty; it opens in the viewer as any document. Its schema is that of the names left open only, so a page that asks for the data shows the other person their part and nothing else.
- Completing it fills the template once, with the values it holds and the data given: the document is the one the template filled with all of the data at once makes, byte for byte.
- The data that completes it gives the names left open and no others: a name filled already, or one the form does not leave open, is an error. Every name left open is given (`null` gives none), or left open again with `-open`, which makes another partly filled form for a third person. `-open '*'` leaves open every name the data does not give.
- Making one, every name of the template is given or named by `-open`: a name forgotten is an error, not a blank on a form someone else signs.
- A name with dots leaves a value of an object open (`-open party_b.sign`). A list (the rows of a table) is filled, or left open, as a whole.
- The pictures given by the paths of files are kept in the partly filled form as `data:` URLs: the other person does not have the files.
- Values are kept as JSON keeps them: a date of YAML becomes its text. Give dates as text (`2026-10-22`) for the same form either way.
- A partly filled form is not protected against changes: whoever has the file can change the values it holds. When the other party fills in their part on a website, keep the .bdff on the server and take only their values.

On the [form page](https://shibukawa.github.io/bdf/form/), check the names under *Leave for another person*: the page offers the document as a .bdff, and a .bdff opened there asks for the names it leaves open only.

## Signatures and seals

`{{%sign | signature}}` and `{{%seal | seal}}` are pictures, filled and drawn as `{{%photo}}` is (`| fit:` too). What differs is what a person gives for them, and the schema says so in `x-bdf.input`:

| Input | On the form page |
|---|---|
| `signature` | A pad to draw on with a pen, a finger or a mouse. The strokes become an SVG image (a `data:` URL), sharp in a PDF at any size |
| `seal` | A picture of a seal (a hanko) to pick, whose white is taken out, so that it can be laid over text |

The signature fields of a PDF form (`/FT /Sig`) and the `signature` fields of pdfme are pictures of `signature` input in the same way.

A drawn signature is a picture of a signature, not a digital signature: it proves nothing of who drew it, nor that the document was not changed after.

## The JSON Schema of the data, and forms for it

What a template asks of data comes out as a JSON Schema: for a library that makes a form of a schema, or for checking data. The fields of the [form page](https://shibukawa.github.io/bdf/form/) are made of it.

```sh
bdf template schema -json-schema invoice.xlsx invoice.schema.json
bdf template schema -json-schema -sample sample.yaml invoice.xlsx   # a sample settles the types and the places
```

Tags do not say the types of their values, so the schema is what the tags suggest.

| In the template | In the schema |
|---|---|
| `{{name}}` | A text or a number (`"type": ["string", "number"]`) |
| `{{total \| fmt:"#,##0"}}` | A number; with a date format (`yyyy/mm/dd`), a text with `"format": "date"` |
| `{{customer.name}}` | The property `name` of the object `customer` |
| A section with names or `{{.}}` in it, or with `\| rows` or `\| pages` | A list, of objects or of values. A section of slots has its number of slots as `maxItems` |
| Another section (`{{#paid}}■{{/paid}}`) | True or false |
| `{{%photo}}` | A text: the path of a file, or a `data:` URL. `\| signature` and `\| seal` add `"input"` to its `x-bdf` |
| `{{%url \| qr}}`, `{{%jan \| barcode:jan}}` | A text: what the code says |

- The properties are in the order of the template. The `x-bdf` of each holds the kind of its tags (`kind`), their number format (`format`), how a section repeats (`mode`, `slots`), the kind of a code (`code`) and where the tags are (`at`).
- A name in a section is looked up in the item and then around it, so the tags alone do not tell whose it is: it is the item's, unless the template uses the same name around the section too.
- **A sample of the data settles it.** The names are where the sample has them, the values have the types of its values, and a section is a list, an object or a condition as its value there is.
- A compiled template (.bdft) holds the schema that its sample settled: `bdf template schema -json-schema invoice.bdft` returns it, and a form is made without reading the template again.
- Nothing is `required`: a tag whose name the data lacks is left empty.
- A PDF form and a template of pdfme say more of their fields than tags do: a checkbox is a boolean, a date field has `"format": "date"`, choices are an `enum`, the value a field holds is its `default`, and the most characters of a field its `maxLength`.

## From Go

```go
data, err := form.LoadData(
	form.Source{FileName: "order.yaml", Data: order},
	form.Source{Name: "items", FileName: "lines.csv", Data: lines})
res, err := form.Fill(template, size, &form.Options{Data: data})
err = res.Doc.WriteSingle(out)
```

`form.Fill` tells the format of a template from its content (an Office file, a PDF form, a template of pdfme, a template of text, a compiled template, a partly filled form). With `Options.Open` it makes a partly filled form, described by `res.Partial`. `form.Fields` returns the names the tags look up, `form.SchemaOf` the JSON Schema of the data, and `form.Compile` a compiled template. For a known format, `xlsx.Fill`, `docx.Fill`, `pptx.Fill` and `textform.Fill` can be called directly. `textpart.Read` reads the print image of a form of text from a document, and `textform.Plain` writes it as plain text. `pdfbdf.Write` writes a PDF file.

```go
res, err := pdfbdf.Write(out, doc, &pdfbdf.Options{})                          // a PDF file
res, err = pdfbdf.Write(out, doc, &pdfbdf.Options{Tagged: true, PDFA: "2a"}) // tagged, and PDF/A-2a
```

## Good to know

- The pages of a workbook break by the rules of Excel, but Excel rounds row heights and column widths to the pixels of the printer when it prints, so an automatic page break can be a row away from Excel's. Make the breaks that matter manual page breaks.
- Text is laid out with the fonts of the machine that converts. A font the template names and the machine lacks is replaced by another, so naming the fonts with `-font-dir` makes the same form wherever it is made.
- The same template and data make the same bytes. For that, `&D` and `&T` (the date and the time) of the headers and footers of Excel are left empty.
- Word and PowerPoint templates do not calculate (they have nothing like the formulas of Excel): totals and what is carried over come with the data.
- A template is untrusted input. Tags run no code, and what a filled form may grow to (cells, slides, elements) is bounded. The expressions of a template of pdfme are bounded too: in their steps, in the text they write (1 MiB) and in the items of the lists they make. The text of a field of a PDF form is cut at 65,536 characters.
