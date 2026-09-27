import { BdfFormatError } from "./bytes.js";

// Standard MIDI Files: the music a view plays (docs/spec.md §4.4). The
// tracks are merged into one list of events in time order, with note ons and
// offs paired into notes and every event timed in seconds by the tempo map.

interface EventBase {
  tick: number;
  /** Seconds from the start. */
  time: number;
  /** The track the event is in (0-based). */
  track: number;
}

export interface NoteEvent extends EventBase {
  type: "note";
  channel: number;
  key: number;
  velocity: number;
  /** Seconds from the note on to its note off. */
  duration: number;
  /** The tick of the note off. */
  endTick: number;
}
export interface ProgramEvent extends EventBase { type: "program"; channel: number; program: number }
export interface ControlEvent extends EventBase { type: "control"; channel: number; controller: number; value: number }
/** Pitch bend, -8192 to 8191. */
export interface BendEvent extends EventBase { type: "bend"; channel: number; value: number }
/** A tempo change (meta FF 51), microseconds per quarter note. */
export interface TempoEvent extends EventBase { type: "tempo"; tempo: number }
export interface TimeSignatureEvent extends EventBase { type: "timeSignature"; numerator: number; denominator: number }
/** key: sharps (> 0) or flats (< 0). */
export interface KeySignatureEvent extends EventBase { type: "keySignature"; key: number; minor: boolean }

export type MidiEvent = NoteEvent | ProgramEvent | ControlEvent | BendEvent | TempoEvent | TimeSignatureEvent | KeySignatureEvent;

/** From tick on, a quarter note lasts tempo microseconds; time is the seconds at tick. */
export interface TempoPoint { tick: number; time: number; tempo: number }

export interface MidiSequence {
  /** 0, 1 or 2 (whose tracks play one after another). */
  format: number;
  /** Ticks per quarter note; with smpte, ticks per second. */
  division: number;
  smpte: boolean;
  /** The tracks read. */
  tracks: number;
  /** In time order; at the same tick in the order of the tracks, and in a track of the file. */
  events: MidiEvent[];
  /** The tempo map, from tick 0 (120 beats per minute until the file sets another tempo). */
  tempos: TempoPoint[];
  /** The end of the music: the last end of track or note off. */
  ticks: number;
  /** The end of the music in seconds. */
  duration: number;
}

/** At most this many events are read: the rest of a larger (or hostile) file is dropped. */
export const MAX_MIDI_EVENTS = 1 << 20;
const MAX_TRACKS = 4096;
/** Notes of one key sounding at once: another ends the first (a hostile file must not make pairing slow). */
const MAX_OPEN = 64;
const DEFAULT_TEMPO = 500000;

/** A note off, paired with its note on while the tracks are merged. */
interface NoteOff { type: "off"; tick: number; time: number; track: number; channel: number; key: number }
type Raw = MidiEvent | NoteOff;

/**
 * Parse a Standard MIDI File. A file cut short or damaged keeps what comes
 * before the damage in each track; notes left on end with the music.
 */
export function parseSmf(bytes: Uint8Array): MidiSequence {
  const b = bytes;
  const u32 = (p: number) => ((b[p] << 24) | (b[p + 1] << 16) | (b[p + 2] << 8) | b[p + 3]) >>> 0;
  const u16 = (p: number) => (b[p] << 8) | b[p + 1];
  const id = (p: number) => String.fromCharCode(b[p], b[p + 1], b[p + 2], b[p + 3]);
  if (b.length < 14 || id(0) !== "MThd") throw new BdfFormatError("not a Standard MIDI File");
  const hlen = u32(4);
  if (hlen < 6) throw new BdfFormatError("bad MIDI header");
  const format = u16(8);
  const raw = u16(12);
  let division: number, smpte = false;
  if (raw & 0x8000) {
    // frames per second (a negative byte; 29 is 29.97 drop frame) times ticks per frame
    const fps = 256 - (raw >> 8);
    division = (fps === 29 ? 29.97 : fps) * (raw & 0xff);
    smpte = true;
  } else {
    division = raw;
  }
  if (division <= 0) throw new BdfFormatError("bad MIDI division");

  const out: Raw[] = [];
  const order: number[] = []; // position of each event in the file, which breaks ties of tick and track
  let tracks = 0;
  let offset = 0; // format 2: each track starts where the one before ends
  let end = 0;
  for (let pos = 8 + hlen; pos + 8 <= b.length && tracks < MAX_TRACKS && out.length < MAX_MIDI_EVENTS; ) {
    const len = u32(pos + 4);
    const start = pos + 8;
    const stop = Math.min(b.length, start + len);
    if (id(pos) === "MTrk") {
      const last = readTrack(b, start, stop, tracks++, format === 2 ? offset : 0, out, order);
      end = Math.max(end, last);
      if (format === 2) offset = last;
    }
    pos = start + len;
  }

  // merge the tracks: a stable sort by tick, then track
  const idx = out.map((_, i) => i);
  idx.sort((i, j) => out[i].tick - out[j].tick || out[i].track - out[j].track || order[i] - order[j]);

  const tempos: TempoPoint[] = [{ tick: 0, time: 0, tempo: DEFAULT_TEMPO }];
  let tp = tempos[0];
  const timeAt = (tick: number) => (smpte ? tick / division : tp.time + ((tick - tp.tick) * tp.tempo) / 1e6 / division);
  const events: MidiEvent[] = [];
  /** Notes sounding, first on first, by channel * 128 + key, and how many on each channel. */
  const open = new Map<number, NoteEvent[]>();
  const sounding = new Array<number>(16).fill(0);
  const close = (n: NoteEvent, tick: number, time: number) => {
    n.endTick = tick;
    n.duration = time - n.time;
    sounding[n.channel]--;
  };
  for (const i of idx) {
    const e = out[i];
    e.time = timeAt(e.tick);
    end = Math.max(end, e.tick);
    switch (e.type) {
      case "off": {
        const q = open.get(e.channel * 128 + e.key);
        const n = q?.shift();
        if (n) close(n, e.tick, e.time);
        continue;
      }
      case "note": {
        const k = e.channel * 128 + e.key;
        const q = open.get(k) ?? [];
        if (q.length >= MAX_OPEN) close(q.shift()!, e.tick, e.time);
        q.push(e);
        open.set(k, q);
        sounding[e.channel]++;
        break;
      }
      case "tempo":
        if (smpte) break;
        if (tp.tick === e.tick) tp.tempo = e.tempo;
        else tempos.push((tp = { tick: e.tick, time: e.time, tempo: e.tempo }));
        break;
      case "control":
        // all sound off, all notes off
        if ((e.controller === 120 || e.controller === 123) && sounding[e.channel] > 0) {
          for (let k = e.channel * 128; k < e.channel * 128 + 128; k++) {
            for (const n of open.get(k) ?? []) close(n, e.tick, e.time);
            open.delete(k);
          }
        }
        break;
    }
    events.push(e);
  }
  const duration = timeAt(end);
  for (const q of open.values()) for (const n of q) close(n, end, duration);
  return { format, division, smpte, tracks, events, tempos, ticks: end, duration };
}

/**
 * Read the events of a track into out; returns the tick where it ends (of the
 * last whole event, when the track is cut short or damaged).
 */
function readTrack(b: Uint8Array, start: number, stop: number, track: number, tick: number, out: Raw[], order: number[]): number {
  let pos = start;
  let running = 0;
  // a variable-length quantity of at most 4 bytes; -1 when damaged or cut short
  const vlq = (): number => {
    let v = 0;
    for (let i = 0; i < 4 && pos < stop; i++) {
      const c = b[pos++];
      v = v * 128 + (c & 0x7f);
      if (!(c & 0x80)) return v;
    }
    return -1;
  };
  const push = (e: Raw) => {
    order.push(pos);
    out.push(e);
  };
  while (pos < stop && out.length < MAX_MIDI_EVENTS) {
    const delta = vlq();
    if (delta < 0 || pos >= stop) break;
    const now = tick + delta;
    let status = b[pos];
    if (status & 0x80) pos++;
    else if (running) status = running; // running status: the data bytes follow at once
    else break;
    const at = { tick: now, time: 0, track };
    if (status === 0xff) {
      if (pos >= stop) break;
      const type = b[pos++];
      const len = vlq();
      if (len < 0 || pos + len > stop) break;
      const d = pos;
      pos += len;
      tick = now;
      if (type === 0x2f) return tick; // end of track
      if (type === 0x51 && len >= 3) {
        const tempo = (b[d] << 16) | (b[d + 1] << 8) | b[d + 2];
        if (tempo > 0) push({ ...at, type: "tempo", tempo });
      } else if (type === 0x58 && len >= 2) {
        push({ ...at, type: "timeSignature", numerator: b[d], denominator: 2 ** Math.min(b[d + 1], 16) });
      } else if (type === 0x59 && len >= 2) {
        push({ ...at, type: "keySignature", key: (b[d] << 24) >> 24, minor: b[d + 1] === 1 });
      }
      continue;
    }
    if (status === 0xf0 || status === 0xf7) {
      // system exclusive: skipped
      const len = vlq();
      if (len < 0 || pos + len > stop) break;
      pos += len;
      tick = now;
      continue;
    }
    if (status > 0xf0) break; // system common and real time messages have no place in a file
    running = status;
    const kind = status & 0xf0, channel = status & 0x0f;
    const n = kind === 0xc0 || kind === 0xd0 ? 1 : 2;
    if (pos + n > stop) break;
    const d1 = b[pos] & 0x7f, d2 = n === 2 ? b[pos + 1] & 0x7f : 0;
    pos += n;
    tick = now;
    switch (kind) {
      case 0x80:
        push({ ...at, type: "off", channel, key: d1 });
        break;
      case 0x90:
        if (d2 === 0) push({ ...at, type: "off", channel, key: d1 });
        else push({ ...at, type: "note", channel, key: d1, velocity: d2, duration: 0, endTick: tick });
        break;
      case 0xb0:
        push({ ...at, type: "control", channel, controller: d1, value: d2 });
        break;
      case 0xc0:
        push({ ...at, type: "program", channel, program: d1 });
        break;
      case 0xe0:
        push({ ...at, type: "bend", channel, value: ((d2 << 7) | d1) - 8192 });
        break;
      // key pressure and channel pressure are not played
    }
  }
  return tick;
}

/** The last tempo point at or before a tick (or a time, with byTime). */
function pointAt(points: TempoPoint[], v: number, byTime: boolean): TempoPoint {
  let lo = 0, hi = points.length - 1;
  while (lo < hi) {
    const mid = (lo + hi + 1) >> 1;
    if ((byTime ? points[mid].time : points[mid].tick) <= v) lo = mid;
    else hi = mid - 1;
  }
  return points[lo];
}

/** The seconds at a tick of a sequence (fractional ticks too). */
export function tickToSeconds(seq: MidiSequence, tick: number): number {
  if (seq.smpte) return tick / seq.division;
  const p = pointAt(seq.tempos, tick, false);
  return p.time + ((tick - p.tick) * p.tempo) / 1e6 / seq.division;
}

/** The tick (fractional) at a time of a sequence, in seconds. */
export function secondsToTick(seq: MidiSequence, seconds: number): number {
  if (seq.smpte) return seconds * seq.division;
  const p = pointAt(seq.tempos, seconds, true);
  return p.tick + ((seconds - p.time) * 1e6 * seq.division) / p.tempo;
}
