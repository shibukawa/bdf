Test font for the templates of text: a subset of BIZ UDGothic Regular
(https://github.com/google/fonts/tree/main/ofl/bizudgothic, SIL Open Font
License 1.1, see OFL.txt), the monospace font that textform looks for first.
It holds ASCII, the box drawing characters (U+2500–U+257F) and the characters
of chars.txt: those of the test templates (../*.txt) and of their data
(test/text/data/*.json). The tests and testdata/form/text-*.bdf lay text out
with this font only (-font-dir converter/textform/testdata/fonts
-no-system-fonts), so the output does not depend on the fonts installed on
the machine.

Made with fontTools from BIZUDGothic-Regular.ttf of google/fonts (SHA-256
709fcd41e3209fb765da750472f55ccdf925653e9fa7e1eb007cb65c8f749c75):

    pyftsubset BIZUDGothic-Regular.ttf --text-file=chars.txt --name-IDs='*' \
      --name-languages='*' --drop-tables+=DSIG,meta --output-file=BIZUDGothic-Regular-subset.ttf

A test template with characters that chars.txt lacks needs the subset made
again with them added.
