import { test } from "node:test";
import assert from "node:assert/strict";
import { parseSmf } from "@bdf/core";
import { MusicPlayer, cursorAtTick, tickAt, eventAt, due, channelStates, soundingEnds, MIN_NOTE } from "../dist/player.js";
import { pulseDuty, patchOf, MAX_VOICES } from "../dist/synth.js";

// Standard MIDI Files, built by hand.
const u16 = (n) => [(n >> 8) & 255, n & 255];
const u32 = (n) => [(n >>> 24) & 255, (n >> 16) & 255, (n >> 8) & 255, n & 255];
const vlq = (n) => {
  const out = [n & 0x7f];
  while ((n >>>= 7)) out.unshift((n & 0x7f) | 0x80);
  return out;
};
const chunk = (id, data) => [...Buffer.from(id, "latin1"), ...u32(data.length), ...data];
const smf = (division, ...tracks) =>
  new Uint8Array([...chunk("MThd", [...u16(tracks.length > 1 ? 1 : 0), ...u16(tracks.length), ...u16(division)]), ...tracks.flatMap((t) => chunk("MTrk", t.flat()))]);
/** A track from [tick, ...bytes] events at absolute ticks. */
const track = (...events) => {
  const out = [];
  let t = 0;
  for (const [tick, ...bytes] of [...events, [events.at(-1)?.[0] ?? 0, 0xff, 0x2f, 0]]) {
    out.push(...vlq(tick - t), ...bytes);
    t = tick;
  }
  return out;
};
const near = (a, b, msg) => assert.ok(Math.abs(a - b) < 1e-6, `${msg ?? ""} ${a} != ${b}`);

// Two systems on page 0 and one on page 1, as a converter would put them.
const cues = {
  systems: [{ page: 0, x: 50, y: 100, w: 500, h: 80 }, { page: 0, x: 50, y: 300, w: 500, h: 80 }, { page: 1, x: 50, y: 100, w: 500, h: 90 }],
  cues: [
    { tick: 0, system: 0, x: 100 },
    { tick: 480, system: 0, x: 300 },
    { tick: 960, system: 0, x: 550 }, // the end of the system
    { tick: 960, system: 1, x: 100 }, // the start of the next, at the same tick: this one counts
    { tick: 1920, system: 1, x: 500 },
    { tick: 2400, system: 2, x: 80 },
  ],
};

test("cursorAtTick: on the system of the last cue at or before the tick", () => {
  assert.deepEqual(cursorAtTick(cues, 0), { page: 0, system: 0, x: 100, y: 100, h: 80 });
  near(cursorAtTick(cues, 240).x, 200, "halfway to the next cue");
  near(cursorAtTick(cues, 720).x, 425);
  assert.deepEqual(cursorAtTick(cues, 960), { page: 0, system: 1, x: 100, y: 300, h: 80 });
  near(cursorAtTick(cues, 959.5).x, 550 - 250 / 480 / 2, "the end of the system comes just before");
  near(cursorAtTick(cues, 1440).x, 300);
  // the next cue is on another system: x holds
  assert.deepEqual(cursorAtTick(cues, 2000), { page: 0, system: 1, x: 500, y: 300, h: 80 });
  assert.deepEqual(cursorAtTick(cues, 2400), { page: 1, system: 2, x: 80, y: 100, h: 90 });
  // after the last cue x holds; before the first, the first
  assert.equal(cursorAtTick(cues, 1e9).x, 80);
  const late = { systems: cues.systems, cues: [{ tick: 100, system: 1, x: 70 }, { tick: 200, system: 1, x: 90 }] };
  assert.deepEqual(cursorAtTick(late, 0), { page: 0, system: 1, x: 70, y: 300, h: 80 });
  assert.equal(cursorAtTick({ systems: [], cues: [] }, 0), null);
});

test("tickAt: the place clicked on a system, between cues or at the nearest", () => {
  near(tickAt(cues, 0, 200), 240);
  near(tickAt(cues, 1, 300), 1440);
  assert.equal(tickAt(cues, 0, 20), 0, "left of the first cue");
  assert.equal(tickAt(cues, 2, 400), 2400, "a system with one cue");
  assert.equal(tickAt(cues, 5, 100), null);
  // a system played twice (a repeat): the pass nearest to where the music is
  const twice = { systems: cues.systems, cues: [{ tick: 0, system: 0, x: 100 }, { tick: 100, system: 0, x: 200 }, { tick: 100, system: 1, x: 100 }, { tick: 200, system: 0, x: 100 }, { tick: 300, system: 0, x: 200 }] };
  near(tickAt(twice, 0, 150, 0), 50);
  near(tickAt(twice, 0, 150, 280), 250);
});

test("eventAt and due find what to schedule", () => {
  const seq = parseSmf(smf(480, track([0, 0x90, 60, 100], [480, 0x80, 60, 0], [480, 0x90, 62, 100], [960, 0x80, 62, 0], [960, 0xb0, 7, 90])));
  const ev = seq.events; // note at 0 s, note at 0.5 s, controller at 1 s
  assert.equal(ev.length, 3);
  assert.equal(eventAt(ev, 0), 0);
  assert.equal(eventAt(ev, 0.1), 1);
  assert.equal(eventAt(ev, 0.5), 1);
  assert.equal(eventAt(ev, 5), 3);
  assert.equal(due(ev, 0, 0.5), 1, "until is not included");
  assert.equal(due(ev, 1, 0.65), 2);
  assert.equal(due(ev, 2, 10), 3);
});

test("channelStates replays programs, controllers and bends", () => {
  const seq = parseSmf(smf(480, track(
    [0, 0xc3, 80], [0, 0xb3, 70, 40], [0, 0xb3, 7, 64], [0, 0xb3, 10, 0], [0, 0xb3, 11, 90], [0, 0xe3, 0, 0x60],
    [480, 0xb3, 121, 0], [480, 0xc3, 81],
  )));
  const at = (i) => channelStates(seq.events, i)[3];
  assert.deepEqual(at(0), { program: 0, volume: 100, expression: 127, pan: 64, duty: 0, sustain: false, bend: 0 });
  assert.deepEqual(at(6), { program: 80, volume: 64, expression: 90, pan: 0, duty: 40, sustain: false, bend: 4096 });
  // reset all controllers leaves volume, pan and the sound controllers
  assert.deepEqual(at(8), { program: 81, volume: 64, expression: 127, pan: 0, duty: 40, sustain: false, bend: 0 });
  assert.equal(channelStates(seq.events, 8)[0].program, 0);
});

test("soundingEnds holds notes released under the sustain pedal", () => {
  const seq = parseSmf(smf(480, track(
    [0, 0x90, 60, 100], [240, 0x80, 60, 0], // 0 to 0.25 s, no pedal
    [480, 0xb0, 64, 127], [480, 0x90, 62, 100], [600, 0x80, 62, 0], [960, 0xb0, 64, 0], // released under the pedal, up at 1 s
    [960, 0x90, 64, 100], [960, 0x80, 64, 0], // zero length
    [1200, 0xb1, 64, 127], [1200, 0x91, 65, 100], [1300, 0x81, 65, 0], // the pedal stays down to the end
    [1920, 0x90, 67, 1], [2000, 0x80, 67, 0],
  )));
  const notes = seq.events.map((e, i) => [e, i]).filter(([e]) => e.type === "note");
  const ends = soundingEnds(seq.events, seq.duration);
  const endOf = (key) => ends[notes.find(([e]) => e.key === key)[1]];
  near(endOf(60), 0.25);
  near(endOf(62), 1);
  near(endOf(64), 1 + MIN_NOTE);
  near(endOf(65), seq.duration);
  near(endOf(67), 2 + 80 / 960);
  assert.equal(ends[seq.events.findIndex((e) => e.type === "control")], 0);
});

test("the duty of Square Lead follows CC 70; programs have sounds by family", () => {
  assert.deepEqual([0, 31, 32, 63, 64, 95, 96, 127].map(pulseDuty), [0.5, 0.5, 0.125, 0.125, 0.25, 0.25, 0.75, 0.75]);
  assert.equal(patchOf(80).wave, "pulse");
  assert.equal(patchOf(81).wave, "sawtooth");
  assert.equal(patchOf(82).wave, "triangle");
  assert.equal(patchOf(79).wave, "sine");
  assert.equal(patchOf(122).wave, "noise");
  assert.equal(patchOf(0).sustain, 0, "a piano dies away");
  assert.equal(patchOf(19).sustain, 1, "an organ holds");
  for (let p = 0; p < 128; p++) assert.ok(patchOf(p).gain > 0 && patchOf(p).gain <= 1, `program ${p}`);
});

/** An AudioContext that plays nothing: it records the sources started, whose clock the test moves. */
function fakeAudio() {
  const started = [];
  const param = (value = 0) => {
    const p = { value, calls: [] };
    for (const m of ["setValueAtTime", "linearRampToValueAtTime", "exponentialRampToValueAtTime", "setTargetAtTime", "cancelScheduledValues", "cancelAndHoldAtTime"]) {
      p[m] = (...args) => { p.calls.push([m, ...args]); return p; };
    }
    return p;
  };
  const node = (props = {}) => ({ connect: (n) => n, disconnect() {}, ...props });
  const source = (kind, props) => {
    const s = node({ kind, detune: param(), onended: null, ...props });
    s.start = (when, offset) => { s.when = when; s.offset = offset; started.push(s); };
    s.stop = (when) => { s.end = when; };
    return s;
  };
  const ctx = {
    currentTime: 0,
    sampleRate: 8000,
    state: "suspended",
    outputLatency: 0,
    destination: node(),
    resumed: 0,
    closed: false,
    createGain: () => node({ gain: param(1) }),
    createDynamicsCompressor: () => node({ threshold: param(), knee: param(), ratio: param() }),
    createStereoPanner: () => node({ pan: param() }),
    createBiquadFilter: () => node({ type: "lowpass", frequency: param(350), Q: param(1) }),
    createOscillator: () => source("osc", { type: "sine", frequency: param(440), setPeriodicWave(w) { this.wave = w; } }),
    createBufferSource: () => source("buffer", { buffer: null, loop: false, playbackRate: param(1) }),
    createBuffer: (channels, length) => ({ getChannelData: () => new Float32Array(length) }),
    createPeriodicWave: (real, imag) => ({ real, imag }),
    resume() { this.resumed++; this.state = "running"; return Promise.resolve(); },
    close() { this.closed = true; return Promise.resolve(); },
  };
  return { ctx, started };
}

// Quarter notes of 0.4 s every half second (120 beats per minute), up an octave from middle C.
const scale = smf(480, track(
  [0, 0xc0, 80], [0, 0xb0, 70, 40], // Square Lead at 12.5%
  ...[0, 1, 2, 3, 4, 5, 6, 7].flatMap((i) => [[i * 480, 0x90, 60 + i, 100], [i * 480 + 384, 0x80, 60 + i, 0]]),
));

test("MusicPlayer schedules the notes of the next moment on the audio clock", (t) => {
  t.mock.timers.enable({ apis: ["setInterval"] });
  const { ctx, started } = fakeAudio();
  const player = new MusicPlayer(scale, cues, { context: () => ctx });
  const updates = [];
  player.onUpdate = () => updates.push(player.state);
  near(player.duration, 3.9);
  assert.equal(player.state, "stopped");
  assert.equal(started.length, 0, "no audio before play");
  player.play();
  assert.equal(ctx.resumed, 1, "the context starts in the gesture");
  assert.deepEqual(updates, ["playing"]);
  // the first note starts a moment after play, on a pulse wave of the duty CC 70 chose
  assert.equal(started.length, 1);
  near(started[0].when, 0.05);
  near(started[0].frequency.value, 261.6255653);
  near(started[0].end, 0.05 + 0.4 + 0.01 * 2 + 0.02);
  assert.ok(started[0].wave.imag instanceof Float32Array);
  near(started[0].wave.imag[1], (2 / Math.PI) * (1 - Math.cos(Math.PI / 4)));

  ctx.currentTime = 0.45; // the position is 0.4: the note at 0.5 is due
  t.mock.timers.tick(30);
  assert.equal(started.length, 2);
  near(started[1].when, 0.55);
  near(player.position, 0.4);
  near(player.cursorAt(player.position).x, 100 + 200 * (0.4 / 0.5));

  ctx.currentTime = 0.75;
  t.mock.timers.tick(30);
  assert.equal(started.length, 2, "nothing more before 0.85");

  // pausing silences what sounds and holds the position
  player.pause();
  assert.equal(player.state, "paused");
  near(player.position, 0.7);
  assert.ok(started[1].end <= 0.75 + 0.1, "the voice fades out at once");
  ctx.currentTime = 2;
  near(player.position, 0.7);
  t.mock.timers.tick(300);
  assert.equal(started.length, 2);

  // playing again goes on from there: the note sounding there starts where it is
  player.play();
  assert.equal(started.length, 3);
  near(started[2].when, 2.05);
  near(started[2].frequency.value, 277.182631, "C sharp, into its second 0.2 s");
  near(started[2].end, 2.05 + 0.2 + 0.04);

  // seeking while playing jumps
  player.seek(3);
  near(player.position, 3);
  const after = started.slice(3);
  assert.equal(after.length, 1);
  near(after[0].when, 2.05);
  near(after[0].frequency.value, 440 * 2 ** ((66 - 69) / 12));
  assert.deepEqual(updates, ["playing", "paused", "playing", "playing"]);

  // to the end: the player stops and says so
  let ended = 0;
  player.onEnd = () => ended++;
  ctx.currentTime = 2.05 + 0.9;
  t.mock.timers.tick(30);
  assert.equal(player.state, "stopped");
  assert.equal(ended, 1);
  assert.equal(player.position, 0);
  player.dispose();
  assert.ok(ctx.closed);
});

test("MusicPlayer: seek and stop without playing, and the time of a place clicked", () => {
  const player = new MusicPlayer(scale, cues, { context: () => { throw new Error("no audio"); } });
  player.seek(1.25);
  assert.equal(player.state, "paused", "a cursor at the place, ready to play from it");
  near(player.position, 1.25);
  const c = player.cursorAt(1.25);
  assert.deepEqual({ ...c, x: 0 }, { page: 0, system: 1, x: 0, y: 300, h: 80 });
  near(c.x, 100 + (400 * (1200 - 960)) / 960);
  // exactly at a cue's time the cue counts, rounding or not
  assert.equal(player.cursorAt(0.5 * 960 / 480).system, 1);
  near(player.timeAt(0, 300), 0.5);
  near(player.timeAt(1, 300), 1.5);
  assert.equal(player.timeAt(4, 300), null);
  player.seek(99);
  near(player.position, player.duration);
  player.stop();
  assert.equal(player.state, "stopped");
  assert.equal(player.position, 0);
  assert.equal(new MusicPlayer(scale).cursorAt(1), null, "no cues, no cursor");
});

test("MusicPlayer: drums on channel 10, and a limit on the voices at once", (t) => {
  t.mock.timers.enable({ apis: ["setInterval"] });
  const on = [], off = [];
  for (let k = 0; k < 100; k++) {
    on.push([480, 0x91, 20 + k, 100]);
    off.push([960, 0x81, 20 + k, 0]);
  }
  const { ctx, started } = fakeAudio();
  const player = new MusicPlayer(smf(480, track([0, 0x99, 36, 100], [0, 0x99, 42, 100], [10, 0x89, 36, 0], [10, 0x89, 42, 0], ...on, ...off)), null, { context: () => ctx });
  player.play();
  // a kick is a sine sweeping down, a hi-hat noise
  assert.equal(started.length, 2);
  const [kick, hat] = started;
  assert.equal(kick.kind, "osc");
  assert.equal(kick.frequency.calls[0][1], 150);
  assert.equal(hat.kind, "buffer");
  assert.ok(kick.end > 0.35 && hat.end < 0.2, "drums ring by themselves, not for the length of the note");
  // a chord of 100 notes: the oldest voices make room for the newest
  ctx.currentTime = 0.45;
  t.mock.timers.tick(30);
  const chord = started.slice(2);
  assert.equal(chord.length, 100);
  assert.equal(chord.filter((s) => s.end < 0.8).length, 100 - MAX_VOICES);
  assert.ok(chord.slice(100 - MAX_VOICES).every((s) => s.end > 1));
  player.dispose();
});
