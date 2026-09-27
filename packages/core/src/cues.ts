import { ByteReader, BdfFormatError } from "./bytes.js";

/** A rectangle of a page that the playing position moves through: a system of a score (docs/spec.md §4.4). */
export interface CueSystem {
  /** Page index in the view (0-based). */
  page: number;
  x: number;
  y: number;
  w: number;
  h: number;
}

/** At tick (of the view's Standard MIDI File) the playing position is at x on a system. */
export interface Cue {
  tick: number;
  /** Index into Cues.systems. */
  system: number;
  x: number;
}

/** The content of a cue index part. */
export interface Cues {
  systems: CueSystem[];
  /** In the order of tick. */
  cues: Cue[];
}

/** uint32 of Go: ticks, pages and system numbers wrap as the encoder's do. */
const u32 = (v: number) => v % 0x100000000;

/** Parse a cue index part (docs/spec.md §4.4), as bdf.DecodeCues does. */
export function decodeCues(bytes: Uint8Array): Cues {
  if (bytes.length < 4 || bytes[0] !== 0x42 || bytes[1] !== 0x43 || bytes[2] !== 0x55 || bytes[3] !== 0x45) throw new BdfFormatError("not a cue index part");
  const r = new ByteReader(bytes);
  r.pos = 4;
  const version = r.u16();
  if (version > 1) throw new BdfFormatError(`unsupported cue index version ${version}`);
  let n = r.varuint();
  if (n > bytes.length) throw new BdfFormatError("bad system count");
  const systems: CueSystem[] = [];
  for (let i = 0; i < n; i++) systems.push({ page: u32(r.varuint()), x: r.f32(), y: r.f32(), w: r.f32(), h: r.f32() });
  n = r.varuint();
  if (n > bytes.length) throw new BdfFormatError("bad cue count");
  const cues: Cue[] = [];
  let t = 0;
  for (let i = 0; i < n; i++) {
    t = u32(t + u32(r.varuint()));
    const q = { tick: t, system: u32(r.varuint()), x: r.f32() };
    if (q.system >= systems.length) throw new BdfFormatError("cue refers to a missing system");
    cues.push(q);
  }
  return { systems, cues };
}
