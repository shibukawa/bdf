STIXTwoMath-cff-subset.otf is a subset of STIX Two Math 2.13 b171
(https://github.com/stipub/stixfonts, the CFF build, SIL Open Font License
1.1, see OFL.txt) with its MATH table, for the tests of the MATH table
reader and of the bounds of CFF glyphs: ( ) x a ∑ ∫ √ 𝑥 and the combining
circumflex, with their variants and parts. Made with fontTools:

    pyftsubset STIXTwoMath-Regular.otf --unicodes='U+0028,U+0029,U+0078,U+0061,U+2211,U+222B,U+221A,U+1D465,U+0302' \
      --layout-features='' --no-hinting --name-IDs='*' --drop-tables+=GSUB,GPOS,GDEF --output-file=STIXTwoMath-cff-subset.otf

The expected values in math_test.go come from fontTools' BoundsPen and MATH
table on this file.
