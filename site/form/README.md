# Public form examples

The examples in `samples/` are editable business documents for the Forms
page. Each Office shape and the pdfme invoice has two designs with the same
data: `standard` (the ordinary filename) uses a centered title, black rules
and a full table grid; `modern` (`-modern`) uses navy headings, more open
spacing and horizontal rules.

The item statement includes payment status, a logo, a QR code and a JAN
code. The invoice has a growing item table, tax, a due date and payment
details. The delivery slip has ten item slots per page, page subtotals and
carried amounts. Company and bank information are sample content.

These examples are independent of `converter/*/testdata/form-*`, which
exercise converter features and provide deterministic regression fixtures.

## Regenerating

Use Python with Pillow, openpyxl, python-docx and python-pptx:

```sh
python site/form/generate_data.py
python site/form/generate_xlsx.py
python site/form/generate_office.py
python site/form/generate_pdfme.py
```

The templates use **M PLUS 1p** Regular and Bold from
[Google Fonts](https://fonts.google.com/specimen/M+PLUS+1p), under the
[SIL Open Font License 1.1](https://github.com/google/fonts/blob/main/ofl/mplus1p/OFL.txt).
`fonts/index.json` contains versioned `fonts.gstatic.com` URLs, file sizes
and the byte ranges used by the font scan. The browser fetches the font
files directly from Google when filling a form; the site does not bundle
them. The same fonts are subset and embedded into the filled BDF and PDF,
so downloaded documents retain their appearance. Requests use fixed font
URLs, without sending form text or data to Google. A first use needs a
network connection; fetched bytes are reused by the converter worker.

To update the pinned font version, obtain the full TTF URLs from
`https://fonts.googleapis.com/css2?family=M+PLUS+1p:wght@400;700` with a
client that receives TTF, download them outside the repository, and update
the sizes and `scan` ranges using `scanRanges` from
`examples/common/build.mjs`. Converter test fonts remain independent.

Build the site with `npm run site`. `npm run test:site` fills both designs
with their shared data, compares their schemas without physical locations,
and exports plain, tagged and PDF/A files. `npm run test:form-page` checks
the browser workspace and style selection. Generated review PDFs belong
under `output/pdf/`; visually inspect the rendered pages after layout
changes, including a long invoice and the second delivery-slip page.
