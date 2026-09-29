// Playing the music of a view (docs/spec.md §4.4) on the main thread. The
// Standard MIDI File is parsed once; while it plays, a timer schedules the
// events of the next moment on the audio clock (Synth), and the cues tell
// where on the pages the music is. The timing and cue logic is in plain
// functions, which run without Web Audio.
import { parseSmf, secondsToTick, tickToSeconds, type Cues, type MidiEvent, type MidiSequence } from "@bdf/core";
import { Synth, applyEvent, initialChannel, type ChannelState } from "./synth.js";

/** Where the music is on the pages: at x on a system, which spans y to y + h of a page (page units). */
export interface Cursor {
  page: number;
  /** Index into Cues.systems. */
  system: number;
  x: number;
  y: number;
  h: number;
}

/**
 * The playing position at a tick (spec §4.4): on the system of the last cue
 * at or before it (the later of cues at one tick), x going linearly to the
 * next cue when that is on the same system, else held. Before the first cue,
 * at the first. null without cues.
 */
export function cursorAtTick(cues: Cues, tick: number): Cursor | null {
  const list = cues.cues;
  if (!list.length) return null;
  let lo = 0, hi = list.length - 1;
  while (lo < hi) {
    const mid = (lo + hi + 1) >> 1;
    if (list[mid].tick <= tick) lo = mid;
    else hi = mid - 1;
  }
  const c = list[lo], n = list[lo + 1];
  let x = c.x;
  if (n && n.system === c.system && n.tick > c.tick && tick > c.tick) x += (n.x - c.x) * Math.min(1, (tick - c.tick) / (n.tick - c.tick));
  const s = cues.systems[c.system];
  return { page: s.page, system: c.system, x, y: s.y, h: s.h };
}

/**
 * The tick at x on a system (where a click on the score plays from): between
 * two cues of the system, as cursorAtTick goes; otherwise at the nearest cue.
 * A place played more than once (repeats) is taken nearest to tick near.
 * null when no cue is on the system.
 */
export function tickAt(cues: Cues, system: number, x: number, near = 0): number | null {
  const list = cues.cues;
  let best: number | null = null, off = Infinity, far = Infinity;
  const take = (tick: number, d: number) => {
    const dn = Math.abs(tick - near);
    if (d < off - 1e-6 || (d <= off + 1e-6 && dn < far)) {
      best = tick;
      off = d;
      far = dn;
    }
  };
  for (let i = 0; i < list.length; i++) {
    const c = list[i];
    if (c.system !== system) continue;
    take(c.tick, Math.abs(c.x - x));
    const n = list[i + 1];
    if (n && n.system === system && n.tick > c.tick && n.x !== c.x) {
      const f = (x - c.x) / (n.x - c.x);
      if (f > 0 && f < 1) take(c.tick + f * (n.tick - c.tick), 0);
    }
  }
  return best;
}

/** The first event at or after a time. */
export function eventAt(events: MidiEvent[], time: number): number {
  let lo = 0, hi = events.length;
  while (lo < hi) {
    const mid = (lo + hi) >> 1;
    if (events[mid].time < time) lo = mid + 1;
    else hi = mid;
  }
  return lo;
}

/** The end of the events from index from that start before until: those to schedule now. */
export function due(events: MidiEvent[], from: number, until: number): number {
  let i = from;
  while (i < events.length && events[i].time < until) i++;
  return i;
}

/** The state of the 16 channels after the events before index end. */
export function channelStates(events: MidiEvent[], end: number): ChannelState[] {
  const states = Array.from({ length: 16 }, initialChannel);
  for (let i = 0; i < end; i++) {
    const e = events[i];
    if (e.type === "program" || e.type === "control" || e.type === "bend") applyEvent(states[e.channel], e);
  }
  return states;
}

/** The shortest a note sounds, seconds. */
export const MIN_NOTE = 0.02;

/**
 * When each note stops sounding (seconds; 0 for other events): at its note
 * off, or, released while the sustain pedal (CC 64) is down, when the pedal
 * comes up (at the latest at end).
 */
export function soundingEnds(events: MidiEvent[], end = Infinity): Float64Array {
  // the times each channel's pedal goes down and comes up
  const pedal: number[][] = Array.from({ length: 16 }, () => []);
  const down: (number | undefined)[] = new Array(16);
  for (const e of events) {
    if (e.type !== "control" || (e.controller !== 64 && e.controller !== 121)) continue;
    const on = e.controller === 64 && e.value >= 64;
    if (on && down[e.channel] === undefined) down[e.channel] = e.time;
    if (!on && down[e.channel] !== undefined) {
      pedal[e.channel].push(down[e.channel]!, e.time);
      down[e.channel] = undefined;
    }
  }
  down.forEach((t, ch) => { if (t !== undefined) pedal[ch].push(t, end); });
  const ends = new Float64Array(events.length);
  events.forEach((e, i) => {
    if (e.type !== "note") return;
    let t = e.time + Math.max(MIN_NOTE, e.duration);
    const p = pedal[e.channel];
    let lo = 0, hi = p.length / 2 - 1, k = -1;
    while (lo <= hi) {
      const mid = (lo + hi) >> 1;
      if (p[2 * mid] <= t) { k = mid; lo = mid + 1; } else hi = mid - 1;
    }
    if (k >= 0 && t < p[2 * k + 1]) t = p[2 * k + 1];
    ends[i] = t;
  });
  return ends;
}

export type PlayState = "stopped" | "playing" | "paused";

export interface MusicPlayerOptions {
  /** Seconds of music scheduled ahead of the audio clock (default 0.15). */
  lookahead?: number;
  /** Milliseconds between scheduling (default 30). */
  interval?: number;
  /** The master volume (default 0.5): a note at velocity 100 and channel volume 100 peaks near 0.15. */
  volume?: number;
  /** Makes the audio context, on the first play (default: a new AudioContext). */
  context?: () => AudioContext;
}

/** Seconds from asking to play to the first sound: time to schedule it. */
const START = 0.05;
interface MeterPoint { tick: number; numerator: number; denominator: number }

/**
 * Plays a view's music with Web Audio, and tells where it is on the pages.
 * The audio context is made on the first play(), which should be called
 * from a user gesture (a click, a key) so that the browser lets it sound.
 */
export class MusicPlayer {
  readonly sequence: MidiSequence;
  /** The first tempo in the file, used to show the opening BPM equivalent. */
  readonly baseBpm: number;
  /** Called when playing starts, pauses, stops (also at the end) or jumps. */
  onUpdate: (() => void) | undefined;
  /** Called when the music has played to its end; the player has stopped. */
  onEnd: (() => void) | undefined;
  private readonly ends: Float64Array;
  private readonly meters: MeterPoint[];
  private ctx: AudioContext | undefined;
  private synth: Synth | undefined;
  private current: PlayState = "stopped";
  /** Where playing started, and the context time that plays at. */
  private from = 0;
  private at = 0;
  /** The position while not playing. */
  private held = 0;
  /** The next event to schedule. */
  private next = 0;
  /** The next metronome beat in musical ticks. */
  private clickTick = 0;
  private clickMeter = 0;
  private clickBeat = 0;
  private clickEnabled = false;
  private speedPercentValue = 100;
  /** Playback seconds per second in the file's original tempo map. */
  private playbackScale = 1;
  private timer: ReturnType<typeof setInterval> | undefined;

  /** seq: a Standard MIDI File, or one parsed already. */
  constructor(seq: Uint8Array | MidiSequence, readonly cues: Cues | null = null, private readonly options: MusicPlayerOptions = {}) {
    this.sequence = seq instanceof Uint8Array ? parseSmf(seq) : seq;
    this.ends = soundingEnds(this.sequence.events, this.sequence.duration);
    this.baseBpm = 60_000_000 / (this.sequence.tempos[0]?.tempo ?? 500_000);
    this.meters = [{ tick: 0, numerator: 4, denominator: 4 }];
    for (const e of this.sequence.events) {
      if (e.type !== "timeSignature") continue;
      const point = { tick: e.tick, numerator: Math.max(1, e.numerator), denominator: Math.max(1, e.denominator) };
      const last = this.meters[this.meters.length - 1];
      if (last.tick === point.tick) this.meters[this.meters.length - 1] = point;
      else this.meters.push(point);
    }
  }

  /** Seconds. */
  get duration(): number { return this.sequence.duration * this.playbackScale; }
  get state(): PlayState { return this.current; }
  get playing(): boolean { return this.current === "playing"; }
  /** The opening BPM after applying the speed percentage. */
  get bpm(): number { return this.baseBpm * this.speedPercentValue / 100; }
  get speedPercent(): number { return this.speedPercentValue; }
  get tempoScale(): number { return this.playbackScale; }
  get bpmAvailable(): boolean { return !this.sequence.smpte; }
  get metronomeEnabled(): boolean { return this.clickEnabled; }
  get metronomeAvailable(): boolean { return !this.sequence.smpte; }

  /** Scale playback speed; all tempo changes in the file keep their relative ratios. */
  setSpeedPercent(percent: number) {
    if (!Number.isFinite(percent)) return;
    const next = Math.max(25, Math.min(200, Math.round(percent)));
    if (next === this.speedPercentValue) return;
    const wasPlaying = this.playing;
    const wasPaused = this.current === "paused";
    const tick = secondsToTick(this.sequence, this.position / this.playbackScale);
    this.speedPercentValue = next;
    this.playbackScale = 100 / next;
    const at = tickToSeconds(this.sequence, tick) * this.playbackScale;
    if (wasPlaying) {
      this.synth?.silence();
      this.start(at);
    } else if (wasPaused) this.held = at;
    this.onUpdate?.();
  }

  /** Turn the beat click on or off; SMPTE-timed MIDI files have no musical beat to follow. */
  setMetronome(enabled: boolean) {
    this.clickEnabled = enabled && this.metronomeAvailable;
    if (this.clickEnabled && this.playing) this.resetClick(secondsToTick(this.sequence, this.position / this.playbackScale));
    if (!this.clickEnabled) this.synth?.clearMetronome();
    else this.schedule();
    this.onUpdate?.();
  }

  /** Seconds from the start: of what is heard now while playing. */
  get position(): number {
    if (this.current !== "playing" || !this.ctx) return this.held;
    return this.clock(this.ctx.outputLatency || 0);
  }

  /** Play from the position (from the start once the music has played to its end). */
  play(): Promise<void> {
    if (this.current === "playing") return Promise.resolve();
    const ctx = this.context();
    // within the user's gesture: the browser lets the context start
    const running = ctx.state === "suspended" ? ctx.resume() : Promise.resolve();
    this.current = "playing";
    this.start(this.held >= this.duration ? 0 : this.held);
    // (music without length has ended already, and said so)
    if (this.playing) this.onUpdate?.();
    return running;
  }

  pause() {
    if (this.current !== "playing") return;
    this.held = this.position;
    this.halt();
    this.current = "paused";
    this.onUpdate?.();
  }

  /** Stop, back at the start. */
  stop() {
    if (this.current === "stopped") return;
    this.halt();
    this.current = "stopped";
    this.held = 0;
    this.onUpdate?.();
  }

  /** Go to a time (seconds); a stopped player is paused there. */
  seek(seconds: number) {
    const t = Math.max(0, Math.min(this.duration, seconds));
    if (this.current === "playing") {
      this.synth!.silence();
      this.start(t);
    } else {
      this.held = t;
      this.current = "paused";
    }
    this.onUpdate?.();
  }

  /** Where the music is on the pages at a time (seconds); null without cues. */
  cursorAt(seconds: number): Cursor | null {
    if (!this.cues) return null;
    // the rounding of seconds must not keep a cue from taking effect at its tick
    return cursorAtTick(this.cues, secondsToTick(this.sequence, seconds / this.playbackScale) + 1e-6);
  }

  /** The time (seconds) at x on a system, for playing from a place clicked on (see tickAt); null when none. */
  timeAt(system: number, x: number): number | null {
    if (!this.cues) return null;
    const tick = tickAt(this.cues, system, x, secondsToTick(this.sequence, this.position / this.playbackScale));
    return tick === null ? null : Math.min(this.duration, tickToSeconds(this.sequence, tick) * this.playbackScale);
  }

  /** Stop and let the audio context go. */
  dispose() {
    this.halt();
    this.current = "stopped";
    this.held = 0;
    this.synth?.dispose();
    this.ctx?.close().catch(() => {});
    this.ctx = this.synth = undefined;
  }

  private context(): AudioContext {
    if (!this.ctx) {
      this.ctx = this.options.context?.() ?? new AudioContext({ latencyHint: "interactive" });
      this.synth = new Synth(this.ctx, this.options.volume ?? 0.5);
    }
    return this.ctx;
  }

  /** The position on the audio clock (lag: seconds the output is behind it). */
  private clock(lag = 0): number {
    return Math.min(this.duration, this.from + Math.max(0, this.ctx!.currentTime - lag - this.at));
  }

  /** Play from a time: the channels as they are there, and the notes sounding there go on. */
  private start(from: number) {
    const ctx = this.ctx!, synth = this.synth!, events = this.sequence.events;
    clearInterval(this.timer);
    this.from = from;
    this.at = ctx.currentTime + START;
    const fileTime = from / this.playbackScale;
    this.next = eventAt(events, fileTime);
    this.resetClick(secondsToTick(this.sequence, fileTime));
    synth.reset(channelStates(events, this.next), this.at);
    // drums are struck: those before are over
    for (let i = 0; i < this.next; i++) {
      const e = events[i];
      if (e.type === "note" && e.channel !== 9 && this.ends[i] > fileTime + MIN_NOTE) {
        synth.note(e.channel, e.key, e.velocity, this.at, this.at + (this.ends[i] - fileTime) * this.playbackScale, (fileTime - e.time) * this.playbackScale);
      }
    }
    this.timer = setInterval(() => this.schedule(), this.options.interval ?? 30);
    this.schedule();
  }

  private schedule() {
    if (!this.ctx || this.current !== "playing") return;
    if (this.position >= this.duration) return this.finish();
    // a hidden page's timers may run late: look further ahead there
    const hidden = typeof document !== "undefined" && document.hidden;
    const ahead = this.options.lookahead ?? 0.15;
    const until = this.clock() + (hidden ? Math.max(1, ahead) : ahead);
    const events = this.sequence.events;
    for (const end = due(events, this.next, until / this.playbackScale); this.next < end; this.next++) {
      const e = events[this.next];
      const when = this.at + e.time * this.playbackScale - this.from;
      if (e.type === "note") this.synth!.note(e.channel, e.key, e.velocity, when, this.at + this.ends[this.next] * this.playbackScale - this.from);
      else this.synth!.event(e, when);
    }
    this.scheduleMetronome(until);
  }

  /** Set the next beat at or after tick, respecting the active time signature. */
  private resetClick(tick: number) {
    let lo = 0, hi = this.meters.length;
    while (lo < hi) {
      const mid = (lo + hi) >> 1;
      if (this.meters[mid].tick <= tick) lo = mid + 1;
      else hi = mid;
    }
    this.clickMeter = Math.max(0, lo - 1);
    const meter = this.meters[this.clickMeter];
    const step = Math.max(1, this.sequence.division * 4 / meter.denominator);
    this.clickBeat = Math.max(0, Math.ceil((tick - meter.tick) / step - 1e-9));
    this.clickTick = meter.tick + this.clickBeat * step;
    const next = this.meters[this.clickMeter + 1];
    if (next && this.clickTick >= next.tick - 1e-9) {
      this.clickMeter++;
      this.clickBeat = 0;
      this.clickTick = next.tick;
    }
  }

  /** Put beat clicks on the same audio clock as notes, including tempo and meter changes. */
  private scheduleMetronome(until: number) {
    if (!this.clickEnabled || this.sequence.smpte || !this.synth) return;
    // A very short lookahead and fast tempo can contain many beats; cap one
    // timer turn so a hostile file cannot monopolize the UI thread.
    for (let count = 0; count < 4096; count++) {
      const time = tickToSeconds(this.sequence, this.clickTick) * this.playbackScale;
      if (time >= this.duration || time >= until) return;
      this.synth.metronome(this.clickBeat % this.meters[this.clickMeter].numerator === 0, this.at + time - this.from);
      const meter = this.meters[this.clickMeter];
      this.clickBeat++;
      this.clickTick += Math.max(1, this.sequence.division * 4 / meter.denominator);
      const next = this.meters[this.clickMeter + 1];
      if (next && this.clickTick >= next.tick - 1e-9) {
        this.clickMeter++;
        this.clickBeat = 0;
        this.clickTick = next.tick;
      }
    }
  }

  /** Played to the end: the last notes die away by themselves. */
  private finish() {
    clearInterval(this.timer);
    this.timer = undefined;
    this.current = "stopped";
    this.held = 0;
    this.onUpdate?.();
    this.onEnd?.();
  }

  private halt() {
    clearInterval(this.timer);
    this.timer = undefined;
    this.synth?.silence();
  }
}
