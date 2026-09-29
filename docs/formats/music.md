# Scores

MML, MIDI and MusicXML all end up as staff notation, engraved by a shared layout engine (`converter/internal/music`) onto A4 pages with the SMuFL font Bravura, and each view carries the underlying music as a Standard MIDI File so the demo viewer can play it back — see [Conversion and output → Scores and playback](../architecture/conversion.md#scores-and-playback) for how that playback works. What the three converters share ends there: two of them have to invent the notation from raw performance data, and one just reads what's already on the page.

## Guitar TAB alternatives

In the viewer's **Display** menu, switch directly between **Score**, **Piano roll**, and **Guitar TAB**. In Guitar TAB, select a part and one of four estimated fingerings: **Low positions**, **Less shifting**, **Open strings**, or **Beginner friendly**. These are alternatives, not a record of how a MIDI performance was played. All use standard six-string tuning (E A D G B E). The first view remains the original staff score, and the TAB views can play the same music with the same speed and metronome controls.

The converter tries every string/fret position for each pitch, then chooses positions across the phrase. Simultaneous and overlapping notes cannot use the same string; fretted notes in a chord must fit a four- or five-fret hand span. MusicXML `<notations><technical><string>` and `<fret>` positions take priority; an alternative is considered only when authored positions conflict with a playable chord. Notes outside the guitar's range, impossible chords, pitched parts beyond the first four, and parts with more than 20,000 notes are left out of TAB with a conversion warning. Use `-param tab=false` to leave out the TAB views.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try a sample: an [MML chiptune](https://shibukawa.github.io/bdf/viewer/?file=samples/frere.mml), a [MIDI karaoke file](https://shibukawa.github.io/bdf/viewer/?file=samples/twinkle.kar), or a [MusicXML piano score](https://shibukawa.github.io/bdf/viewer/?file=samples/minuet.musicxml). Use **Display** to choose the notation, adjust playback speed as a percentage (the opening BPM equivalent is shown and tempo ratios stay intact), and turn on the metronome while playing.

## MML (converter/mml)

MML and MIDI are both performance data — timed notes rather than notated ones — so both go through the same notation step (`Notate`) before layout. Notes are quantized onto rhythmic values (ties across beats and bar lines, dots, triplets), split into measures and up to two voices per staff, spelled in a key (the file's own, or one estimated from the notes), and placed on a clef chosen from their pitch range (keyboard instruments get a grand staff).

`converter/mml` reads three dialects, detected from the content: generic (FlMML and the BASIC `PLAY` statement), Mabinogi's `MML@…;`, and the NES MCK/PPMCK track-line format. A signature-less `.mml` file with fewer than two recognizable MML lines converts as MML only by extension, since there's nothing in the content to detect.

| Option | Values | Default |
|---|---|---|
| `dialect` | `auto`, `generic` (FlMML, MSX BASIC), `mabinogi` (`MML@…;`), `ppmck` (NES MCK/PPMCK) | `auto` |
| `octave` | `auto`, `normal` (`>` raises the octave), `reverse` (`<` raises it) | `auto` — the dialect's own convention |
| `time` | a time signature, e.g. `3/4`, `6/8` | `4/4` |
| `key` | a key signature, e.g. `G`, `Bb`, `F#m`, `Em` | estimated from the notes |
| `program` | a General MIDI program number, `0`–`127`, for tracks that set no tone | the dialect's default |
| `tab` | `true` or `false` | `true` |

## MIDI (converter/midi)

`converter/midi` reads Standard MIDI File formats 0, 1 and 2, and RIFF `RMID`. Each track (or track/channel pair) becomes a part, named from the track name or, failing that, from its General MIDI program. Karaoke lyrics (`.kar` text, or meta event 05) attach to the notes of whichever track they best line up with. Like MML, a MIDI file's notes are performance data and go through the same `Notate` step described above before they're engraved.

| Option | Values | Default |
|---|---|---|
| `time` | a time signature (`3/4`, `6/8`, `C` for common time, `cut`) | the file's own |
| `key` | a key signature (`G`, `Bb`, `F#m`, or `-2` for two flats) | the file's own, or estimated from the notes |
| `tab` | `true` or `false` | `true` |

## MusicXML (converter/musicxml)

MusicXML, unlike MML and MIDI, is read as written rather than re-notated: parts and staves, voices, chords, beams, tuplets, grace notes, accidentals, ties and slurs, articulations and ornaments, dynamics and hairpins, lyrics by verse, clef/key/time changes, repeats and voltas, tempo and rehearsal marks, and system breaks all come straight from the file's own markup rather than being inferred. Both the partwise and timewise document forms are read, along with the compressed `.mxl` container. Playback takes repeats, and jumps like D.C./D.S./Fine/To Coda, into account.

MusicXML accepts `-param tab=false` to omit the inferred TAB views. The staff score still converts as written.

See [design.md §3.27](../design.md#327-楽譜と演奏convertermmlconvertermidiconvertermusicxmlconverterinternalmusic) for how the shared engraving engine — measure spacing, stems and beams, ties and slurs, lyric wrapping — actually works.
