#!/usr/bin/env python3
"""Writes the test scores of converter/musicxml: public-domain music typed
here (Beethoven's Ode to Joy with Schiller's words, the Minuet in G of the
Notebook for Anna Magdalena Bach, "Ah! vous dirai-je, maman"), and the
compressed ode_to_joy.mxl. Run it in this directory."""

import zipfile

DTD = ('<!DOCTYPE score-partwise PUBLIC "-//Recordare//DTD MusicXML 4.0 Partwise//EN" '
       '"http://www.musicxml.org/dtds/partwise.dtd">')

TYPES = {1: "16th", 2: "eighth", 3: "eighth", 4: "quarter", 6: "quarter", 8: "half", 12: "half", 16: "whole"}
DOTTED = {3, 6, 12}


def note(pitch, dur, div=4, lyric=None, beam=None, staff=None, voice="1", chord=False, extra=""):
    """A note: pitch "F#4" or "r" for a rest; dur in sixteenths."""
    out = ["<note>"]
    if chord:
        out.append("<chord/>")
    if pitch == "r":
        out.append("<rest/>")
    else:
        step, octave, alter = pitch[0], pitch[-1], 0
        if "#" in pitch:
            alter = 1
        if "b" in pitch[1:]:
            alter = -1
        p = f"<step>{step}</step>"
        if alter:
            p += f"<alter>{alter}</alter>"
        out.append(f"<pitch>{p}<octave>{octave}</octave></pitch>")
    out.append(f"<duration>{dur * div // 4}</duration>")
    out.append(f"<voice>{voice}</voice>")
    out.append(f"<type>{TYPES[dur]}</type>")
    if dur in DOTTED:
        out.append("<dot/>")
    if staff:
        out.append(f"<staff>{staff}</staff>")
    if beam:
        out.append(f'<beam number="1">{beam}</beam>')
    out.append(extra)
    if lyric:
        syl, text = lyric
        out.append(f'<lyric number="1"><syllabic>{syl}</syllabic><text>{text}</text></lyric>')
    out.append("</note>")
    return "".join(out)


def words(text):
    """Lyric syllables: "Freu- de, schö- ner" -> (syllabic, text) pairs."""
    out, cont = [], False
    for w in text.split():
        hyph = w.endswith("-")
        w = w.rstrip("-")
        if cont:
            syl = "middle" if hyph else "end"
        else:
            syl = "begin" if hyph else "single"
        out.append((syl, w))
        cont = hyph
    return out


def ode():
    tune = [
        "F#4:4 F#4:4 G4:4 A4:4", "A4:4 G4:4 F#4:4 E4:4", "D4:4 D4:4 E4:4 F#4:4", "F#4:6 E4:2 E4:8",
        "F#4:4 F#4:4 G4:4 A4:4", "A4:4 G4:4 F#4:4 E4:4", "D4:4 D4:4 E4:4 F#4:4", "E4:6 D4:2 D4:8",
        "E4:4 E4:4 F#4:4 D4:4", "E4:4 F#4:2 G4:2 F#4:4 D4:4", "E4:4 F#4:2 G4:2 F#4:4 E4:4", "D4:4 E4:4 A3:8",
        "F#4:4 F#4:4 G4:4 A4:4", "A4:4 G4:4 F#4:4 E4:4", "D4:4 D4:4 E4:4 F#4:4", "E4:6 D4:2 D4:8",
    ]
    text = words("Freu- de, schö- ner Göt- ter- fun- ken, Toch- ter aus E- ly- si- um, "
                 "wir be- tre- ten feu- er- trun- ken, Himm- li- sche, dein Hei- lig- tum! "
                 "Dei- ne Zau- ber bin- den wie- der, was die Mo- de streng ge- teilt; "
                 "al- le Men- schen wer- den Brü- der, wo dein sanf- ter Flü- gel weilt.")
    melisma = {(9, 2), (10, 2)}  # the second of a slurred pair has no syllable
    ms, t = [], 0
    for i, bar in enumerate(tune):
        notes = [(n.split(":")[0], int(n.split(":")[1])) for n in bar.split()]
        body = ""
        if i == 0:
            body += ('<attributes><divisions>2</divisions><key><fifths>2</fifths><mode>major</mode></key>'
                     '<time><beats>4</beats><beat-type>4</beat-type></time><clef><sign>G</sign><line>2</line></clef>'
                     '</attributes><direction placement="above"><direction-type><words font-weight="bold">Allegro assai'
                     '</words></direction-type><direction-type><metronome><beat-unit>quarter</beat-unit>'
                     '<per-minute>120</per-minute></metronome></direction-type><sound tempo="120"/></direction>'
                     '<direction placement="below"><direction-type><dynamics><mf/></dynamics></direction-type>'
                     '<sound dynamics="88.89"/></direction>')
        if i == 8:
            body = '<barline location="left"><bar-style>heavy-light</bar-style><repeat direction="forward"/></barline>'
        for j, (p, d) in enumerate(notes):
            beam = extra = None
            if d == 2 and j + 1 < len(notes) and notes[j + 1][1] == 2:
                beam, extra = "begin", '<notations><slur type="start" number="1"/></notations>'
            elif d == 2 and j > 0 and notes[j - 1][1] == 2:
                beam, extra = "end", '<notations><slur type="stop" number="1"/></notations>'
            lyric = None
            if (i, j) not in melisma:
                lyric = text[t]
                t += 1
            body += note(p, d, div=2, lyric=lyric, beam=beam, extra=extra or "")
        if i == 15:
            body += '<barline location="right"><bar-style>light-heavy</bar-style><repeat direction="backward"/></barline>'
        ms.append(f'<measure number="{i + 1}">{body}</measure>')
    return score("Ode to Joy", "Ludwig van Beethoven", "Friedrich Schiller",
                 '<score-part id="P1"><part-name>Voice</part-name><score-instrument id="P1-I1"><instrument-name>Voice'
                 '</instrument-name></score-instrument><midi-instrument id="P1-I1"><midi-channel>1</midi-channel>'
                 '<midi-program>53</midi-program><volume>80</volume><pan>0</pan></midi-instrument></score-part>',
                 [("P1", ms)], movement="Symphony No. 9, finale")


def minuet():
    # right hand, left hand: sixteenths; measures 1-8, the first half
    rh = [
        "D5:4 G4:2 A4:2 B4:2 C5:2", "D5:4 G4:4 G4:4", "E5:4 C5:2 D5:2 E5:2 F#5:2", "G5:4 G4:4 G4:4",
        "C5:4 D5:2 C5:2 B4:2 A4:2", "B4:4 C5:2 B4:2 A4:2 G4:2", "F#4:4 G4:2 A4:2 B4:2 G4:2", "A4:12",
        "B4:4 G4:2 A4:2 B4:2 G4:2",
    ]
    lh = [
        "G3:8 A3:4", "B3:12", "C4:12", "B3:12", "A3:12", "G3:12", "D4:4 B3:4 G3:4", "D4:4 D3:8", "G3:12",
    ]
    ms_rh, ms_lh = [], []
    for i in range(len(rh)):
        body = ""
        if i == 0:
            body += ('<attributes><divisions>4</divisions><key><fifths>1</fifths></key><time><beats>3</beats>'
                     '<beat-type>4</beat-type></time><staves>2</staves><clef number="1"><sign>G</sign><line>2</line>'
                     '</clef><clef number="2"><sign>F</sign><line>4</line></clef></attributes>'
                     '<direction placement="below"><direction-type><dynamics><p/></dynamics></direction-type>'
                     '<staff>1</staff></direction>')
        if i == 7:
            body += '<barline location="left"><ending number="1" type="start"/></barline>'
        if i == 8:
            body += '<barline location="left"><ending number="2" type="start">2.</ending></barline>'
        for s, part in ((1, rh), (2, lh)):
            notes = [(n.split(":")[0], int(n.split(":")[1])) for n in part[i].split()]
            run = [j for j, (_, d) in enumerate(notes) if d == 2]
            for j, (p, d) in enumerate(notes):
                beam = None
                if d == 2:
                    k = run.index(j)
                    # eighths beamed in pairs
                    beam = "begin" if k % 2 == 0 else "end"
                body += note(p, d, div=4, beam=beam, staff=s, voice="1" if s == 1 else "5")
            if s == 1:
                body += "<backup><duration>12</duration></backup>"
        if i == 7:
            body += ('<barline location="right"><bar-style>light-heavy</bar-style><ending number="1" type="stop"/>'
                     '<repeat direction="backward"/></barline>')
        if i == 8:
            body += '<barline location="right"><ending number="2" type="discontinue"/></barline>'
        ms_rh.append(f'<measure number="{i + 1}">{body}</measure>')
    return score("Minuet in G", "Christian Petzold", None,
                 '<score-part id="P1"><part-name>Piano</part-name><midi-instrument id="P1-I1"><midi-channel>1'
                 '</midi-channel><midi-program>1</midi-program></midi-instrument></score-part>',
                 [("P1", ms_rh)])


def score(title, composer, lyricist, partlist, parts, movement=None):
    out = ['<?xml version="1.0" encoding="UTF-8" standalone="no"?>', DTD, '<score-partwise version="4.0">']
    out.append(f"<work><work-title>{title}</work-title></work>")
    if movement:
        out.append(f"<movement-title>{movement}</movement-title>")
    out.append("<identification>")
    out.append(f'<creator type="composer">{composer}</creator>')
    if lyricist:
        out.append(f'<creator type="lyricist">{lyricist}</creator>')
    out.append("<rights>Public domain</rights></identification>")
    out.append(f"<part-list>{partlist}</part-list>")
    for pid, ms in parts:
        out.append(f'<part id="{pid}">')
        out.extend(ms)
        out.append("</part>")
    out.append("</score-partwise>")
    return "\n".join(out) + "\n"


def twinkle_timewise():
    # "Ah! vous dirai-je, maman": melody and a bass line, timewise
    mel = ["C4:4 C4:4 G4:4 G4:4", "A4:4 A4:4 G4:8", "F4:4 F4:4 E4:4 E4:4", "D4:4 D4:4 C4:8"]
    bass = ["C3:8 E3:8", "F3:8 E3:8", "D3:8 C3:8", "G3:8 C3:8"]
    out = ['<?xml version="1.0" encoding="UTF-8"?>',
           '<!DOCTYPE score-timewise PUBLIC "-//Recordare//DTD MusicXML 4.0 Timewise//EN" '
           '"http://www.musicxml.org/dtds/timewise.dtd">',
           '<score-timewise version="4.0">', '<movement-title>Ah! vous dirai-je, maman</movement-title>',
           '<part-list><score-part id="P1"><part-name>Flute</part-name></score-part>'
           '<score-part id="P2"><part-name>Cello</part-name></score-part></part-list>']
    for i in range(4):
        out.append(f'<measure number="{i + 1}">')
        for pid, line, clef in (("P1", mel, ("G", 2)), ("P2", bass, ("F", 4))):
            body = ""
            if i == 0:
                body = (f'<attributes><divisions>1</divisions><key><fifths>0</fifths></key><time><beats>4</beats>'
                        f'<beat-type>4</beat-type></time><clef><sign>{clef[0]}</sign><line>{clef[1]}</line></clef>'
                        f'</attributes>')
            for n in line[i].split():
                p, d = n.split(":")
                body += note(p, int(d), div=1)
            if i == 3:
                body += '<barline location="right"><bar-style>light-heavy</bar-style></barline>'
            out.append(f'<part id="{pid}">{body}</part>')
        out.append("</measure>")
    out.append("</score-timewise>")
    return "\n".join(out) + "\n"


def main():
    files = {"ode_to_joy.musicxml": ode(), "minuet.musicxml": minuet(), "twinkle_timewise.musicxml": twinkle_timewise()}
    for name, text in files.items():
        with open(name, "w", encoding="utf-8") as f:
            f.write(text)
    container = ('<?xml version="1.0" encoding="UTF-8"?>\n<container><rootfiles>'
                 '<rootfile full-path="ode_to_joy.musicxml" media-type="application/vnd.recordare.musicxml+xml"/>'
                 '</rootfiles></container>\n')
    with zipfile.ZipFile("ode_to_joy.mxl", "w") as z:
        info = zipfile.ZipInfo("mimetype", date_time=(2026, 1, 1, 0, 0, 0))
        z.writestr(info, "application/vnd.recordare.musicxml", compress_type=zipfile.ZIP_STORED)
        info = zipfile.ZipInfo("META-INF/container.xml", date_time=(2026, 1, 1, 0, 0, 0))
        z.writestr(info, container, compress_type=zipfile.ZIP_DEFLATED)
        info = zipfile.ZipInfo("ode_to_joy.musicxml", date_time=(2026, 1, 1, 0, 0, 0))
        z.writestr(info, files["ode_to_joy.musicxml"], compress_type=zipfile.ZIP_DEFLATED)


if __name__ == "__main__":
    main()
