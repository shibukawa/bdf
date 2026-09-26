Test fonts for the PowerPoint converter: subsets of M PLUS 1p Regular and
Bold (https://github.com/google/fonts/tree/main/ofl/mplus1p, SIL Open Font
License 1.1, see OFL.txt) holding ASCII, Latin-1, kana, CJK punctuation and
the characters the test decks use. The tests and testdata/pptx lay text out with
these fonts only (-font-dir testdata/fonts -no-system-fonts), so the output
does not depend on the fonts installed on the machine.

Made with fontTools:

    pyftsubset MPLUS1p-Regular.ttf --text-file=chars.txt --layout-features='*' \
      --no-hinting --name-IDs='*' --name-languages='*' --output-file=MPLUS1p-Regular-subset.ttf
