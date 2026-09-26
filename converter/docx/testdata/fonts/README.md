Test fonts for the Word converter: subsets of M PLUS 1p Regular and Bold
(https://github.com/google/fonts/tree/main/ofl/mplus1p, SIL Open Font
License 1.1, see OFL.txt) holding ASCII, Latin-1, general punctuation,
arrows, geometric shapes, circled numbers, kana, CJK punctuation, the
vertical forms (U+FE10–FE19, of which the font has only a few), the
full-width forms and the characters the test documents use. The tests and
testdata/docx lay text out with these fonts only (-font-dir testdata/fonts
-no-system-fonts), so the output does not depend on the fonts installed on
the machine. Every family the documents ask for (Century, Arial, 游明朝 …)
falls back to them.

Made with fontTools, chars.txt holding the ranges above and the characters
of converter/docx/testdata/*.docx:

    pyftsubset MPLUS1p-Regular.ttf --text-file=chars.txt --layout-features='*' \
      --no-hinting --name-IDs='*' --name-languages='*' --output-file=MPLUS1p-Regular-subset.ttf
