# IDML documents exported by InDesign

Real InDesign exports, from the regression tests of SimpleIDML
(https://github.com/Starou/SimpleIDML, `tests/regressiontests/IDML/`),
copyright (c) 2012 Stanislas Guerra, BSD 3-clause license (see LICENSE).
The test documents written by test/idml/gen.py exercise the converter on
markup written by hand; these check it on what InDesign itself writes
(DOMVersion 7.5 = CS5.5, 10.0 = CC 2014, 15.1 = 2020):

- 4-pages.idml, 4-pages-layers-with-guides.idml: four pages in three spreads
  with a master (the second on two layers, with guides)
- 2articles-1photo.idml: two articles and two placed photos, linked to files
  that are not here (their frames stay empty, with a warning)
- magazineA-courrier-des-lecteurs-3pages.idml: letters to the editor over
  three pages, stories threaded between frames and across pages
- magazineA-bloc-notes.idml: a text frame whose insets leave no room for its
  text, a frame anchored in text
- interview.idml: a rotated banner, a rounded quote box, white text on a
  colored frame
- page-9modules.idml: empty graphic frames (placeholders) and one headline

The photos of the SimpleIDML media directory are not included. The Go tests
lay the text out with the test fonts of the PowerPoint converter; the fonts
the documents name (Minion Pro, Arial) are replaced.
