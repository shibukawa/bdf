"""Writes the tune of the song the audio converter's tests play, as a WAV file.

An introduction of a second and a half, four phrases of three seconds (the
lines of the lyrics: test/audio/gen.sh gives them the times 1.5, 4.5, 7.5 and
10.5 s) and a last note: fifteen seconds, mono, 22050 Hz. Each
note is a sine with a little of its octave, struck and left to die away.
"""
import math
import struct
import sys
import wave

RATE = 22050
NOTE = 0.375  # seconds: eight to a phrase

# MIDI note numbers; 0 is a rest
INTRO = [60, 64, 67, 72]
PHRASES = [
    [64, 67, 69, 67, 64, 67, 72, 72],
    [72, 69, 67, 64, 62, 64, 67, 67],
    [64, 67, 69, 67, 64, 62, 60, 0],
    [62, 64, 67, 64, 62, 62, 60, 0],
]


def note(key, seconds):
    n = int(RATE * seconds)
    if key == 0:
        return [0.0] * n
    f = 440.0 * 2 ** ((key - 69) / 12)
    out = []
    for i in range(n):
        t = i / RATE
        attack = min(1.0, t / 0.01)
        release = min(1.0, (seconds - t) / 0.02)
        level = attack * release * math.exp(-3.0 * t)
        out.append(level * (0.8 * math.sin(2 * math.pi * f * t) + 0.2 * math.sin(4 * math.pi * f * t)))
    return out


samples = []
for key in INTRO:
    samples += note(key, NOTE)
for phrase in PHRASES:
    for key in phrase:
        samples += note(key, NOTE)
samples += note(60, 1.5)

with wave.open(sys.argv[1], "wb") as w:
    w.setnchannels(1)
    w.setsampwidth(2)
    w.setframerate(RATE)
    w.writeframes(b"".join(struct.pack("<h", int(max(-1.0, min(1.0, 0.6 * s)) * 32767)) for s in samples))
