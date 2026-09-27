import { test } from "node:test";
import assert from "node:assert/strict";
import { BdfDocument, decodeCues, parseSmf, tickToSeconds, secondsToTick } from "../dist/index.js";

// Standard MIDI Files, built by hand.
const u16 = (n) => [(n >> 8) & 255, n & 255];
const u32 = (n) => [(n >>> 24) & 255, (n >> 16) & 255, (n >> 8) & 255, n & 255];
const vlq = (n) => {
  const out = [n & 0x7f];
  while ((n >>>= 7)) out.unshift((n & 0x7f) | 0x80);
  return out;
};
const chunk = (id, data) => [...Buffer.from(id, "latin1"), ...u32(data.length), ...data];
const ev = (delta, ...bytes) => [...vlq(delta), ...bytes];
const EOT = ev(0, 0xff, 0x2f, 0);
const tempo = (delta, us) => ev(delta, 0xff, 0x51, 3, (us >> 16) & 255, (us >> 8) & 255, us & 255);
const smf = (format, division, tracks) =>
  new Uint8Array([...chunk("MThd", [...u16(format), ...u16(tracks.length), ...u16(division)]), ...tracks.flatMap((t) => chunk("MTrk", t.flat()))]);

const notes = (seq) => seq.events.filter((e) => e.type === "note");
const near = (a, b, msg) => assert.ok(Math.abs(a - b) < 1e-9, `${msg ?? ""} ${a} != ${b}`);

test("parseSmf: format 0, running status, note on with velocity 0 as note off", () => {
  const seq = parseSmf(smf(0, 480, [[
    ev(0, 0xc0, 5), // program 5 on channel 1
    ev(0, 0x90, 60, 100),
    ev(0, 64, 90), // running status
    ev(480, 60, 0), // note off by velocity 0
    ev(0, 0x80, 64, 0),
    ev(240, 0x91, 67, 80),
    ev(240, 0x81, 67, 0),
    EOT,
  ]]));
  assert.equal(seq.format, 0);
  assert.equal(seq.division, 480);
  assert.equal(seq.smpte, false);
  assert.equal(seq.tracks, 1);
  const n = notes(seq);
  assert.deepEqual(n.map((e) => [e.tick, e.channel, e.key, e.velocity, e.endTick]), [[0, 0, 60, 100, 480], [0, 0, 64, 90, 480], [720, 1, 67, 80, 960]]);
  near(n[0].duration, 0.5);
  near(n[2].time, 0.75);
  near(n[2].duration, 0.25);
  assert.deepEqual(seq.events[0], { type: "program", tick: 0, time: 0, track: 0, channel: 0, program: 5 });
  assert.equal(seq.ticks, 960);
  near(seq.duration, 1);
});

test("parseSmf: the tempo map times events, and ticks and seconds convert both ways", () => {
  const seq = parseSmf(smf(1, 96, [
    [tempo(0, 1000000), tempo(96, 250000), ev(0, 0xff, 0x58, 4, 3, 3, 24, 8), ev(0, 0xff, 0x59, 2, 0xfd, 1), EOT],
    [ev(0, 0x90, 60, 64), ev(192, 0x80, 60, 64), EOT],
  ]));
  assert.deepEqual(seq.tempos, [{ tick: 0, time: 0, tempo: 1000000 }, { tick: 96, time: 1, tempo: 250000 }]);
  const [n] = notes(seq);
  near(n.duration, 1.25, "a second, then a quarter of one");
  near(tickToSeconds(seq, 48), 0.5);
  near(tickToSeconds(seq, 144), 1.125);
  near(secondsToTick(seq, 1.125), 144);
  near(secondsToTick(seq, 0.25), 24);
  for (const t of [0, 10, 95.5, 96, 150, 1000]) near(secondsToTick(seq, tickToSeconds(seq, t)), t);
  const ts = seq.events.find((e) => e.type === "timeSignature");
  assert.deepEqual([ts.numerator, ts.denominator], [3, 8]);
  const ks = seq.events.find((e) => e.type === "keySignature");
  assert.deepEqual([ks.key, ks.minor], [-3, true]);
  // the tempo change at tick 96 comes before the note in the second track at a later tick
  assert.deepEqual(seq.events.map((e) => e.type), ["tempo", "note", "tempo", "timeSignature", "keySignature"]);
});

test("parseSmf: the default tempo is 120 beats per minute; a tempo at tick 0 replaces it", () => {
  const plain = parseSmf(smf(0, 480, [[ev(0, 0x90, 60, 1), ev(960, 0x80, 60, 0), EOT]]));
  assert.deepEqual(plain.tempos, [{ tick: 0, time: 0, tempo: 500000 }]);
  near(plain.duration, 1);
  const set = parseSmf(smf(0, 480, [[tempo(0, 1000000), tempo(0, 2000000), ev(480, 0x90, 60, 1), ev(480, 0x80, 60, 0), EOT]]));
  assert.deepEqual(set.tempos, [{ tick: 0, time: 0, tempo: 2000000 }]);
  near(set.duration, 4);
});

test("parseSmf: tracks merge in time order, the same tick in the order of the tracks", () => {
  const seq = parseSmf(smf(1, 480, [
    [ev(0, 0xb0, 7, 100), ev(480, 0xb0, 7, 90), EOT],
    [ev(0, 0xc1, 40), ev(0, 0x91, 60, 100), ev(480, 0xe1, 0, 0x40), ev(0, 0x81, 60, 0), EOT],
  ]));
  assert.deepEqual(seq.events.map((e) => [e.tick, e.track, e.type]), [
    [0, 0, "control"], [0, 1, "program"], [0, 1, "note"], [480, 0, "control"], [480, 1, "bend"],
  ]);
  assert.equal(seq.events.find((e) => e.type === "bend").value, 0);
  assert.equal(seq.events.find((e) => e.type === "control").value, 100);
});

test("parseSmf: overlapping notes of a key end first on, first off", () => {
  const seq = parseSmf(smf(0, 100, [[
    ev(0, 0x90, 60, 10), ev(100, 0x90, 60, 20), ev(100, 0x80, 60, 0), ev(100, 0x80, 60, 0), ev(0, 0x80, 60, 0), EOT,
  ]]));
  assert.deepEqual(notes(seq).map((n) => [n.velocity, n.tick, n.endTick]), [[10, 0, 200], [20, 100, 300]]);
});

test("parseSmf: notes left on end with the track; all notes off ends those of its channel", () => {
  const seq = parseSmf(smf(0, 480, [[
    ev(0, 0x90, 60, 100), ev(0, 0x91, 62, 100), ev(480, 0xb1, 123, 0), ev(480, 0xff, 0x2f, 0),
  ]]));
  assert.deepEqual(notes(seq).map((n) => [n.key, n.endTick]), [[60, 960], [62, 480]]);
  near(notes(seq)[0].duration, 1);
});

test("parseSmf: system exclusive and unknown meta events are skipped", () => {
  const seq = parseSmf(smf(0, 480, [[
    ev(0, 0xf0, 5, 0x7e, 0x7f, 9, 1, 0xf7),
    ev(0, 0xff, 0x03, 4, ...Buffer.from("Horn")),
    ev(0, 0xf7, 2, 1, 2),
    ev(0, 0x90, 72, 50),
    ev(0, 0xff, 0x01, 0), // an empty text between the note on and its running status
    ev(10, 72, 0),
    EOT,
  ]]));
  assert.deepEqual(notes(seq).map((n) => [n.key, n.endTick]), [[72, 10]]);
});

test("parseSmf: SMPTE division counts ticks per second and ignores tempo", () => {
  // -25 frames per second, 40 ticks per frame
  const seq = parseSmf(smf(0, 0xe728, [[tempo(0, 1000000), ev(0, 0x90, 60, 1), ev(500, 0x80, 60, 0), EOT]]));
  assert.equal(seq.smpte, true);
  assert.equal(seq.division, 1000);
  near(notes(seq)[0].duration, 0.5);
  near(tickToSeconds(seq, 250), 0.25);
  near(secondsToTick(seq, 2), 2000);
});

test("parseSmf: format 2 plays its tracks one after another", () => {
  const seq = parseSmf(smf(2, 480, [
    [ev(0, 0x90, 60, 1), ev(480, 0x80, 60, 0), EOT],
    [ev(0, 0x90, 62, 1), ev(480, 0x80, 62, 0), EOT],
  ]));
  assert.deepEqual(notes(seq).map((n) => [n.key, n.tick, n.endTick]), [[60, 0, 480], [62, 480, 960]]);
  near(seq.duration, 1);
});

test("parseSmf: files cut short or damaged keep what comes before", () => {
  const whole = smf(0, 480, [[ev(0, 0x90, 60, 100), ev(480, 0x80, 60, 0), ev(0, 0x90, 62, 100), ev(480, 0x80, 62, 0), EOT]]);
  for (let n = 14; n < whole.length; n++) {
    const seq = parseSmf(whole.subarray(0, n));
    assert.ok(seq.events.length <= 2);
    for (const e of notes(seq)) assert.ok(e.duration >= 0 && e.endTick >= e.tick);
  }
  const cut = parseSmf(whole.subarray(0, whole.length - 6));
  assert.deepEqual(notes(cut).map((n) => [n.key, n.endTick]), [[60, 480], [62, 480]]);
  // a data byte without a status, a system real time byte, a quantity over 4 bytes
  assert.equal(parseSmf(smf(0, 480, [[ev(0, 60, 100), ev(0, 0x90, 50, 1)]])).events.length, 0);
  for (const bad of [ev(5, 0xf8), [0x80, 0x80, 0x80, 0x80, 0x00, 0x90, 60, 1]]) {
    const seq = parseSmf(smf(0, 480, [[ev(0, 0x90, 50, 1), bad, ev(10, 0x80, 50, 0)]]));
    assert.deepEqual(notes(seq).map((n) => [n.key, n.endTick]), [[50, 0]]);
    assert.equal(seq.events.length, 1);
  }
  // a chunk longer than the file, an alien chunk, a huge meta length
  const alien = new Uint8Array([...chunk("MThd", [0, 1, 0, 2, 0, 96]), ...chunk("XFIH", [1, 2, 3]), ...chunk("MTrk", [...ev(0, 0x90, 60, 1), ...ev(0, 0xff, 0x01, 0x8f, 0xff, 0xff, 0x7f)])]);
  assert.equal(notes(parseSmf(alien)).length, 1);
  const long = new Uint8Array([...chunk("MThd", [0, 0, 0, 1, 0, 96]), ...Buffer.from("MTrk"), 0xff, 0xff, 0xff, 0xff, ...ev(0, 0x90, 60, 1)]);
  assert.equal(notes(parseSmf(long)).length, 1);
  assert.throws(() => parseSmf(new Uint8Array(20)), /not a Standard MIDI File/);
  assert.throws(() => parseSmf(new Uint8Array([...chunk("MThd", [0, 0, 0, 1])])), /not a Standard MIDI File|bad MIDI header/);
  assert.throws(() => parseSmf(smf(0, 0, [[EOT]])), /bad MIDI division/);
});

test("parseSmf: a hostile file is read in bounded time and memory", () => {
  // a track of n events of 3 bytes: a note on, then running status
  const flood = (n, make) => {
    const data = new Uint8Array(4 + n * 3);
    data.set([0, 0x90, 60, 1]);
    for (let i = 0; i < n; i++) data.set(make(i), 4 + i * 3);
    return new Uint8Array([...chunk("MThd", [0, 0, 0, 1, 1, 0xe0]), ...Buffer.from("MTrk"), ...u32(data.length), ...data]);
  };
  const t0 = performance.now();
  // many notes of one key on at once, then all off
  const ons = parseSmf(flood(300000, (i) => (i < 150000 ? [0, 60, 100] : [1, 60, 0])));
  assert.equal(notes(ons).length, 150001);
  assert.ok(notes(ons).every((n) => n.endTick >= n.tick));
  // more events than are read
  const many = parseSmf(flood(1200000, (i) => [1, 40 + (i % 40), i % 2 ? 0 : 90]));
  assert.ok(many.events.length < 1 << 20);
  assert.ok(performance.now() - t0 < 10000, `${performance.now() - t0} ms`);
});

// Cue index parts (docs/spec.md §4.4), built as bdf.EncodeCues does.
function encodeCues({ systems, cues }, version = 1) {
  const out = [...Buffer.from("BCUE"), version & 255, version >> 8];
  const varuint = (n) => {
    do {
      let b = n % 128;
      n = Math.floor(n / 128);
      if (n) b |= 0x80;
      out.push(b);
    } while (n);
  };
  const f32 = (v) => out.push(...new Uint8Array(new Float32Array([v]).buffer));
  varuint(systems.length);
  for (const s of systems) {
    varuint(s.page);
    for (const v of [s.x, s.y, s.w, s.h]) f32(v);
  }
  varuint(cues.length);
  let t = 0;
  for (const q of cues) {
    varuint(q.tick - t);
    t = q.tick;
    varuint(q.system);
    f32(q.x);
  }
  return new Uint8Array(out);
}

const sample = {
  systems: [{ page: 0, x: 40, y: 100, w: 500, h: 80.5 }, { page: 1, x: 40, y: 60, w: 500, h: 80 }],
  cues: [{ tick: 0, system: 0, x: 60 }, { tick: 1920, system: 0, x: 540 }, { tick: 1920, system: 1, x: 60 }, { tick: 200000, system: 1, x: 300.25 }],
};

test("decodeCues reads what bdf.EncodeCues writes", () => {
  assert.deepEqual(decodeCues(encodeCues(sample)), sample);
  assert.deepEqual(decodeCues(encodeCues({ systems: [], cues: [] })), { systems: [], cues: [] });
});

test("decodeCues rejects what bdf.DecodeCues rejects", () => {
  const good = encodeCues(sample);
  assert.throws(() => decodeCues(new Uint8Array([0x42, 0x43])), /not a cue index part/);
  assert.throws(() => decodeCues(Uint8Array.from(good, (b, i) => (i === 0 ? 0x41 : b))), /not a cue index part/);
  assert.throws(() => decodeCues(encodeCues(sample, 2)), /unsupported cue index version 2/);
  assert.throws(() => decodeCues(encodeCues({ systems: sample.systems, cues: [{ tick: 0, system: 2, x: 0 }] })), /missing system/);
  for (let n = 4; n < good.length; n++) assert.throws(() => decodeCues(good.subarray(0, n)), /unexpected end of data/);
  // counts larger than the part could hold
  assert.throws(() => decodeCues(new Uint8Array([...Buffer.from("BCUE"), 1, 0, 0xff, 0x7f])), /bad system count/);
  assert.throws(() => decodeCues(new Uint8Array([...Buffer.from("BCUE"), 1, 0, 0, 0xff, 0x7f])), /bad cue count/);
});

test("BdfDocument.play reads a view's music and its cues", async () => {
  const seq = smf(0, 480, [[ev(0, 0x90, 60, 100), ev(480, 0x80, 60, 0), EOT]]);
  const cues = encodeCues(sample);
  const parts = new Map([["a".repeat(32), seq], ["b".repeat(32), cues]]);
  const entry = (h, t) => ({ h, t, enc: "identity", len: parts.get(h).length, size: parts.get(h).length });
  const manifest = {
    bdf: 1, opset: 1, unit: "pt",
    views: [
      { id: "score", kind: "fixed", pages: [], play: { seq: "a".repeat(32), cues: "b".repeat(32) } },
      { id: "nocues", kind: "fixed", pages: [], play: { seq: "a".repeat(32) } },
      { id: "plain", kind: "fixed", pages: [] },
    ],
    parts: [entry("a".repeat(32), "seq"), entry("b".repeat(32), "idx")],
  };
  const doc = await BdfDocument.open({ manifest: async () => manifest, stored: async (e) => parts.get(e.h) });
  const p = await doc.play(doc.view("score"));
  assert.deepEqual([...p.seq], [...seq]);
  assert.deepEqual(p.cues, sample);
  assert.equal((await doc.play(doc.view("nocues"))).cues, null);
  assert.equal(await doc.play(doc.view("plain")), null);
});
