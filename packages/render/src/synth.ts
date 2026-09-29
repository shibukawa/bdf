// A small synthesizer for the music of a view (docs/spec.md §4.4), built of
// Web Audio nodes: no samples. Each note is an oscillator (or noise) through
// an envelope, into its channel's volume and pan, into the master volume.
// The timbre follows the General MIDI program family; channel 10 is drums.
import type { MidiEvent } from "@bdf/core";

/** What a channel plays with at a moment: its program and controllers. */
export interface ChannelState {
  program: number;
  /** CC 7 and CC 11, 0 to 127. */
  volume: number;
  expression: number;
  /** CC 10: 0 left, 64 center, 127 right. */
  pan: number;
  /** CC 70: the duty of Square Lead's pulse wave (see pulseDuty). */
  duty: number;
  /** CC 64. */
  sustain: boolean;
  /** Pitch bend, -8192 to 8191 (±2 semitones). */
  bend: number;
}

/** A channel before any event: General MIDI's defaults. */
export const initialChannel = (): ChannelState => ({ program: 0, volume: 100, expression: 127, pan: 64, duty: 0, sustain: false, bend: 0 });

/** Apply a program change, controller or pitch bend to the state of its channel; false when it changes nothing played. */
export function applyEvent(s: ChannelState, e: MidiEvent): boolean {
  switch (e.type) {
    case "program":
      s.program = e.program;
      return true;
    case "bend":
      s.bend = e.value;
      return true;
    case "control":
      switch (e.controller) {
        case 7: s.volume = e.value; return true;
        case 10: s.pan = e.value; return true;
        case 11: s.expression = e.value; return true;
        case 64: s.sustain = e.value >= 64; return true;
        case 70: s.duty = e.value; return true;
        case 121: // reset all controllers: not volume, pan nor the sound controllers
          s.expression = 127;
          s.sustain = false;
          s.bend = 0;
          return true;
      }
  }
  return false;
}

/** The duty of Square Lead's pulse wave for CC 70 (docs/spec.md §4.4). */
export function pulseDuty(cc70: number): number {
  return cc70 < 32 ? 0.5 : cc70 < 64 ? 0.125 : cc70 < 96 ? 0.25 : 0.75;
}

/** A sound: a wave, an envelope (seconds; decay and release are time constants) and a lowpass. */
export interface Patch {
  wave: OscillatorType | "pulse" | "organ" | "piano" | "noise";
  gain: number;
  attack: number;
  decay: number;
  /** The level held after the decay: 0 dies away (and faster the higher the note), like a string struck or plucked. */
  sustain: number;
  release: number;
  /** A lowpass at mul times the note's frequency plus add Hz. */
  filter?: [mul: number, add: number];
}

const patch = (wave: Patch["wave"], gain: number, attack: number, decay: number, sustain: number, release: number, filter?: [number, number]): Patch =>
  ({ wave, gain, attack, decay, sustain, release, filter });

const PIANO = patch("piano", 0.7, 0.004, 0.9, 0, 0.08, [6, 600]);
const MALLET = patch("sine", 0.8, 0.002, 0.35, 0, 0.05);
const ORGAN = patch("organ", 0.42, 0.01, 0.1, 1, 0.04);
const PLUCK = patch("sawtooth", 0.5, 0.003, 0.45, 0, 0.06, [3, 400]);
const BASS = patch("sawtooth", 0.6, 0.005, 0.5, 0.35, 0.06, [2, 200]);
const STRINGS = patch("sawtooth", 0.3, 0.08, 0.5, 0.85, 0.2, [3, 800]);
const CHOIR = patch("triangle", 0.55, 0.1, 0.5, 0.9, 0.25);
const BRASS = patch("sawtooth", 0.35, 0.03, 0.3, 0.8, 0.08, [5, 300]);
const REED = patch("square", 0.28, 0.02, 0.2, 0.9, 0.06, [4, 300]);
const PIPE = patch("triangle", 0.6, 0.04, 0.2, 0.9, 0.08);
const OCARINA = patch("sine", 0.6, 0.03, 0.2, 0.95, 0.06);
// the chip waves of MML: sharp edges, held as they are
const SQUARE = patch("pulse", 0.28, 0.002, 0.1, 1, 0.01);
const SAW = patch("sawtooth", 0.3, 0.002, 0.1, 1, 0.01);
const TRIANGLE = patch("triangle", 0.4, 0.002, 0.1, 1, 0.01);
const LEAD = patch("sawtooth", 0.3, 0.005, 0.3, 0.9, 0.05, [4, 500]);
const PAD = patch("sawtooth", 0.22, 0.3, 1, 0.8, 0.5, [2, 600]);
const FX = patch("triangle", 0.42, 0.08, 1, 0.7, 0.5);
const PERCUSSIVE = patch("sine", 0.8, 0.002, 0.25, 0, 0.05);
const NOISE = patch("noise", 0.42, 0.005, 0.3, 0.6, 0.06);
const CHIP_NOISE = patch("noise", 0.35, 0.002, 0.1, 1, 0.01);

/** The sound of a General MIDI program. */
export function patchOf(program: number): Patch {
  if (program < 8) return PIANO;
  if (program < 16) return MALLET; // chromatic percussion
  if (program < 24) return ORGAN;
  if (program < 32) return PLUCK; // guitars
  if (program < 40) return BASS;
  if (program < 52) return STRINGS; // and ensembles
  if (program < 55) return CHOIR;
  if (program < 56) return STRINGS; // orchestra hit
  if (program < 64) return BRASS;
  if (program < 72) return REED;
  if (program < 79) return PIPE;
  if (program === 79) return OCARINA;
  if (program === 80) return SQUARE;
  if (program === 81) return SAW;
  if (program === 82) return TRIANGLE; // Calliope Lead: MML's triangle wave
  if (program < 88) return LEAD;
  if (program < 96) return PAD;
  if (program < 104) return FX;
  if (program < 112) return PLUCK; // ethnic
  if (program < 120) return PERCUSSIVE;
  if (program === 122) return CHIP_NOISE; // Seashore: MML's noise channel
  return NOISE; // sound effects
}

/** A note sounding (or scheduled), until stop. */
interface Voice {
  channel: number;
  start: number;
  stop: number;
  sources: AudioScheduledSourceNode[];
  env: GainNode;
  detune: AudioParam[];
}

interface Channel {
  input: GainNode;
  pan: StereoPannerNode | null;
  state: ChannelState;
}

/** Voices at once: a new one beyond this takes the place of the oldest. */
export const MAX_VOICES = 64;

export class Synth {
  private readonly master: GainNode;
  private readonly channels: Channel[] = [];
  /** In the order they start. */
  private voices: Voice[] = [];
  private readonly waves = new Map<string, PeriodicWave>();
  private readonly buffers = new Map<string, AudioBuffer>();

  constructor(readonly ctx: BaseAudioContext, volume: number) {
    this.master = ctx.createGain();
    this.master.gain.value = volume;
    // many voices at once stay out of clipping
    const limit = ctx.createDynamicsCompressor();
    limit.threshold.value = -8;
    limit.knee.value = 6;
    limit.ratio.value = 8;
    this.master.connect(limit).connect(ctx.destination);
    for (let ch = 0; ch < 16; ch++) {
      const input = ctx.createGain();
      const pan = typeof ctx.createStereoPanner === "function" ? ctx.createStereoPanner() : null;
      if (pan) input.connect(pan).connect(this.master);
      else input.connect(this.master);
      this.channels.push({ input, pan, state: initialChannel() });
    }
  }

  /** Put the channels in these states at when: where playing starts. */
  reset(states: ChannelState[], when: number) {
    states.forEach((s, ch) => {
      this.channels[ch].state = { ...s };
      this.level(ch, when, true);
    });
  }

  /** A program change, controller or pitch bend at when. */
  event(e: MidiEvent, when: number) {
    if (e.type !== "program" && e.type !== "control" && e.type !== "bend") return;
    const c = this.channels[e.channel];
    if (!applyEvent(c.state, e)) return;
    if (e.type === "control" && [7, 10, 11, 121].includes(e.controller)) this.level(e.channel, when);
    if (e.type === "bend" || (e.type === "control" && e.controller === 121)) {
      const cents = (c.state.bend / 8192) * 200;
      for (const v of this.voices) if (v.channel === e.channel && v.stop > when) for (const d of v.detune) d.setValueAtTime(cents, when);
    }
  }

  /**
   * A note from when to end (context time). into: how far into the note
   * playing starts (a note sounding where playing starts), which leaves out
   * its attack.
   */
  note(channel: number, key: number, velocity: number, when: number, end: number, into = 0) {
    this.prune(when);
    if (channel === 9) return this.drum(key, velocity, when);
    const ctx = this.ctx;
    const s = this.channels[channel].state;
    const p = patchOf(s.program);
    const f = 440 * 2 ** ((key - 69) / 12);
    let src: AudioScheduledSourceNode, detune: AudioParam;
    if (p.wave === "noise") {
      // noise that changes value at a rate that follows the key, as a chip's noise channel does
      const b = ctx.createBufferSource();
      b.buffer = this.noise(4);
      b.loop = true;
      b.playbackRate.value = Math.min(8, Math.max(1 / 32, 2 ** ((key - 60) / 12)));
      src = b;
      detune = b.detune;
    } else {
      const o = ctx.createOscillator();
      o.frequency.value = f;
      if (p.wave === "pulse") o.setPeriodicWave(this.pulse(pulseDuty(s.duty)));
      else if (p.wave === "organ" || p.wave === "piano") o.setPeriodicWave(this.wave(p.wave));
      else o.type = p.wave;
      src = o;
      detune = o.detune;
    }
    detune.value = (s.bend / 8192) * 200;
    const env = ctx.createGain();
    let out: AudioNode = src;
    if (p.filter) {
      const lp = ctx.createBiquadFilter();
      lp.type = "lowpass";
      lp.frequency.value = Math.min(ctx.sampleRate * 0.45, f * p.filter[0] + p.filter[1]);
      out = src.connect(lp);
    }
    out.connect(env).connect(this.channels[channel].input);

    const peak = p.gain * (velocity / 127) ** 2;
    // a sound that dies away dies faster the higher it is
    const decay = p.sustain === 0 ? p.decay * 2 ** ((60 - key) / 24) : p.decay;
    const g = env.gain;
    g.setValueAtTime(0, when);
    let held: number;
    if (into > 0) {
      // where the envelope has got to, faded in over a moment
      const level = peak * (p.sustain + (1 - p.sustain) * Math.exp(-Math.max(0, into - p.attack) / decay));
      held = when + 0.01;
      g.linearRampToValueAtTime(level, held);
    } else {
      held = when + Math.min(p.attack, Math.max(0.001, end - when));
      g.linearRampToValueAtTime(peak, held);
    }
    if (p.sustain < 1) g.setTargetAtTime(peak * p.sustain, held, decay);
    const off = Math.max(end, held);
    g.setTargetAtTime(0, off, p.release / 3);
    const stop = off + p.release * 2 + 0.02;
    // noise from anywhere in its buffer, so that notes do not repeat each other
    if (p.wave === "noise") (src as AudioBufferSourceNode).start(when, Math.random() * 0.5);
    else src.start(when);
    src.stop(stop);
    src.onended = () => env.disconnect();
    this.voices.push({ channel, start: when, stop, sources: [src], env, detune: [detune] });
  }

  /** A short metronome click, independent of the document's instrument channels. */
  metronome(accent: boolean, when: number) {
    this.prune(when);
    const src = this.ctx.createOscillator();
    const env = this.ctx.createGain();
    const stop = when + 0.055;
    src.type = "sine";
    src.frequency.value = accent ? 1760 : 1320;
    env.gain.setValueAtTime(0.0001, when);
    env.gain.exponentialRampToValueAtTime(accent ? 0.18 : 0.12, when + 0.001);
    env.gain.exponentialRampToValueAtTime(0.0001, when + 0.045);
    src.connect(env).connect(this.master);
    src.start(when);
    src.stop(stop);
    src.onended = () => env.disconnect();
    this.voices.push({ channel: -1, start: when, stop, sources: [src], env, detune: [] });
  }

  /** Cancel clicks already queued when the metronome is switched off. */
  clearMetronome() {
    const now = this.ctx.currentTime;
    const keep: Voice[] = [];
    for (const v of this.voices) {
      if (v.channel === -1) this.release(v, now, 0.003);
      else keep.push(v);
    }
    this.voices = keep;
  }

  /** Silence every voice at once (with a short fade), as playing stops or jumps. */
  silence() {
    const now = this.ctx.currentTime;
    for (const v of this.voices) this.release(v, now, 0.008);
    this.voices = [];
  }

  dispose() {
    this.silence();
    this.master.disconnect();
  }

  /** Drop the voices that have ended, and make room for one more. */
  private prune(when: number) {
    const now = this.ctx.currentTime;
    this.voices = this.voices.filter((v) => v.stop > now);
    while (this.voices.length >= MAX_VOICES) this.release(this.voices.shift()!, Math.max(when, now), 0.005);
  }

  private release(v: Voice, at: number, tc: number) {
    const g = v.env.gain;
    if (v.start < at) {
      if (typeof g.cancelAndHoldAtTime === "function") g.cancelAndHoldAtTime(at);
      else g.cancelScheduledValues(at);
      g.setTargetAtTime(0, at, tc);
    }
    // stopped before it starts, a source never plays
    const stop = v.start < at ? at + tc * 8 : at;
    for (const s of v.sources) {
      try { s.stop(stop); } catch { /* stopped already (older browsers allow one stop) */ }
    }
    v.stop = Math.min(v.stop, stop);
  }

  /** Channel 10: drums synthesized by General MIDI key. */
  private drum(key: number, velocity: number, when: number) {
    const ctx = this.ctx;
    const env = ctx.createGain();
    env.connect(this.channels[9].input);
    const v = (velocity / 127) ** 2;
    const sources: AudioScheduledSourceNode[] = [];
    const fade = (gain: number, decay: number) => {
      const g = ctx.createGain();
      g.gain.setValueAtTime(gain * v, when);
      g.gain.exponentialRampToValueAtTime(0.0001, when + decay);
      g.connect(env);
      return g;
    };
    const noise = (type: BiquadFilterType, freq: number, gain: number, decay: number) => {
      const src = ctx.createBufferSource();
      src.buffer = this.noise(1);
      src.loop = true;
      const f = ctx.createBiquadFilter();
      f.type = type;
      f.frequency.value = freq;
      src.connect(f).connect(fade(gain, decay));
      src.start(when, Math.random() * 0.5);
      sources.push(src);
      return decay;
    };
    const tone = (type: OscillatorType, from: number, to: number, sweep: number, gain: number, decay: number) => {
      const o = ctx.createOscillator();
      o.type = type;
      o.frequency.setValueAtTime(from, when);
      o.frequency.exponentialRampToValueAtTime(to, when + sweep);
      o.connect(fade(gain, decay));
      o.start(when);
      sources.push(o);
      return decay;
    };
    let len: number;
    switch (key) {
      case 35: case 36: len = tone("sine", 150, 45, 0.12, 2, 0.35); break; // kicks
      case 38: case 40: len = Math.max(noise("highpass", 1500, 1, 0.18), tone("triangle", 190, 160, 0.05, 0.8, 0.09)); break; // snares
      case 37: len = noise("bandpass", 2500, 1, 0.04); break; // side stick
      case 39: len = noise("bandpass", 1200, 2.5, 0.12); break; // hand clap
      case 42: len = noise("highpass", 7000, 0.6, 0.05); break; // closed hi-hat
      case 44: len = noise("highpass", 7000, 0.6, 0.07); break; // pedal hi-hat
      case 46: len = noise("highpass", 7000, 0.6, 0.35); break; // open hi-hat
      case 41: case 43: case 45: case 47: case 48: case 50: { // toms, low to high
        const f = 70 + (key - 41) * 15;
        len = tone("sine", f * 1.6, f, 0.15, 1.6, 0.35);
        break;
      }
      case 49: case 57: len = noise("highpass", 5000, 0.6, 1.2); break; // crash
      case 51: case 59: len = noise("highpass", 6000, 0.45, 0.6); break; // ride
      case 52: case 55: len = noise("highpass", 4000, 0.6, 0.8); break; // china, splash
      case 53: len = Math.max(tone("sine", 820, 800, 0.05, 0.4, 0.6), noise("highpass", 6000, 0.3, 0.5)); break; // ride bell
      default: len = noise("bandpass", 300 * 2 ** ((key - 35) / 12), 0.8, 0.1);
    }
    const stop = when + len + 0.05;
    for (const s of sources) s.stop(stop);
    sources[0].onended = () => env.disconnect();
    this.voices.push({ channel: 9, start: when, stop, sources, env, detune: [] });
  }

  /** Set a channel's volume (CC 7 × CC 11) and pan at when; jump: at once, otherwise over a few ms. */
  private level(ch: number, when: number, jump = false) {
    const c = this.channels[ch], s = c.state;
    const gain = (s.volume / 127) ** 2 * (s.expression / 127) ** 2;
    const pan = Math.max(-1, Math.min(1, (s.pan - 64) / 63));
    const set = (p: AudioParam, v: number) => {
      if (!jump) return p.setTargetAtTime(v, when, 0.005);
      p.cancelScheduledValues(0);
      p.setValueAtTime(v, when);
    };
    set(c.input.gain, gain);
    if (c.pan) set(c.pan.pan, pan);
  }

  /** One second of noise that holds each value for hold samples. */
  private noise(hold: number): AudioBuffer {
    const name = `noise${hold}`;
    let b = this.buffers.get(name);
    if (!b) {
      const n = this.ctx.sampleRate;
      b = this.ctx.createBuffer(1, n, n);
      const d = b.getChannelData(0);
      for (let i = 0; i < n; i += hold) d.fill(Math.random() * 2 - 1, i, i + hold);
      this.buffers.set(name, b);
    }
    return b;
  }

  /** A pulse wave of a duty (0 to 1). */
  private pulse(duty: number): PeriodicWave {
    const name = `pulse${duty}`;
    let w = this.waves.get(name);
    if (!w) {
      const n = 64;
      const re = new Float32Array(n), im = new Float32Array(n);
      for (let k = 1; k < n; k++) {
        re[k] = (2 / (k * Math.PI)) * Math.sin(2 * Math.PI * k * duty);
        im[k] = (2 / (k * Math.PI)) * (1 - Math.cos(2 * Math.PI * k * duty));
      }
      w = this.ctx.createPeriodicWave(re, im);
      this.waves.set(name, w);
    }
    return w;
  }

  /** The drawbars of an organ, or the partials of a piano's string. */
  private wave(name: "organ" | "piano"): PeriodicWave {
    let w = this.waves.get(name);
    if (!w) {
      const im = name === "organ"
        ? new Float32Array([0, 1, 0.8, 0.5, 0.35, 0, 0.25, 0, 0.18])
        : Float32Array.from({ length: 16 }, (_, k) => (k === 0 ? 0 : 1 / k ** 1.6));
      w = this.ctx.createPeriodicWave(new Float32Array(im.length), im);
      this.waves.set(name, w);
    }
    return w;
  }
}
