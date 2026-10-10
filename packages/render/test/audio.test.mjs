import { test } from "node:test";
import assert from "node:assert/strict";
import { AudioPlayer, lineAtTime, timeOfLine } from "../dist/audio.js";

// The lines of a song's lyrics on a card, as the audio converter puts them:
// a cue a line, in milliseconds; the third line is sung again at the end.
const cues = {
  systems: [
    { page: 0, x: 32, y: 400, w: 296, h: 16 },
    { page: 0, x: 32, y: 416, w: 296, h: 32 },
    { page: 0, x: 32, y: 456, w: 296, h: 16 },
  ],
  cues: [
    { tick: 1500, system: 0, x: 32 },
    { tick: 4500, system: 1, x: 32 },
    { tick: 8000, system: 2, x: 32 },
    { tick: 12000, system: 0, x: 32 },
    { tick: 20000, system: 2, x: 32 },
  ],
};

test("lineAtTime: the line of the last cue at or before the time, all of it", () => {
  assert.equal(lineAtTime(cues, 0), null, "nothing during the introduction");
  assert.equal(lineAtTime(cues, 1499.9), null);
  assert.deepEqual(lineAtTime(cues, 1500), { page: 0, system: 0, x: 32, y: 400, w: 296, h: 16 });
  assert.equal(lineAtTime(cues, 4499).system, 0);
  assert.deepEqual(lineAtTime(cues, 4500), { page: 0, system: 1, x: 32, y: 416, w: 296, h: 32 });
  assert.equal(lineAtTime(cues, 11999).system, 2);
  assert.equal(lineAtTime(cues, 12000).system, 0, "a line sung again");
  assert.equal(lineAtTime(cues, 1e9).system, 2, "the last line stays");
  assert.equal(lineAtTime({ systems: [], cues: [] }, 0), null);
  assert.equal(lineAtTime(cues, NaN), null);
  // of cues at one time the later counts
  const same = { systems: cues.systems, cues: [{ tick: 1000, system: 0, x: 32 }, { tick: 1000, system: 1, x: 32 }] };
  assert.equal(lineAtTime(same, 1000).system, 1);
});

test("timeOfLine: the time of a line, the nearest when it is sung more than once", () => {
  assert.equal(timeOfLine(cues, 1), 4500);
  assert.equal(timeOfLine(cues, 0, 0), 1500);
  assert.equal(timeOfLine(cues, 0, 9000), 12000);
  assert.equal(timeOfLine(cues, 2, 15000), 20000);
  assert.equal(timeOfLine(cues, 2, 9000), 8000);
  assert.equal(timeOfLine(cues, 7), null);
});

/** What AudioPlayer uses of an audio element. */
class FakeAudio extends EventTarget {
  paused = true;
  ended = false;
  duration = NaN;
  playbackRate = 1;
  defaultPlaybackRate = 1;
  preload = "";
  src = "";
  error = null;
  loads = 0;
  time = 0;
  /** What play() answers with: a rejection for a browser that refuses to sound. */
  refuse = null;
  get currentTime() { return this.time; }
  set currentTime(t) {
    this.time = t;
    this.ended = false;
    this.dispatchEvent(new Event("seeked"));
  }
  play() {
    if (this.refuse) return Promise.reject(this.refuse);
    if (this.ended) this.time = 0;
    this.ended = false;
    this.paused = false;
    this.dispatchEvent(new Event("play"));
    return Promise.resolve();
  }
  pause() {
    if (this.paused) return;
    this.paused = true;
    this.dispatchEvent(new Event("pause"));
  }
  load() { this.loads++; }
  removeAttribute(name) { if (name === "src") this.src = ""; }
  /** The browser read the length of the file. */
  loaded(duration) {
    this.duration = duration;
    this.dispatchEvent(new Event("durationchange"));
  }
  /** The file played to its end. */
  end() {
    this.time = this.duration;
    this.paused = true;
    this.ended = true;
    this.dispatchEvent(new Event("pause"));
    this.dispatchEvent(new Event("ended"));
  }
}

test("AudioPlayer plays, pauses, seeks and stops an audio element", async () => {
  const el = new FakeAudio();
  const player = new AudioPlayer(new Uint8Array([1, 2, 3]), "audio/mpeg", cues, { element: el });
  const states = [];
  let ended = 0;
  player.onUpdate = () => states.push(player.state);
  player.onEnd = () => ended++;
  assert.match(el.src, /^blob:/);
  assert.equal(el.preload, "metadata");
  assert.equal(player.state, "stopped");
  assert.equal(player.duration, 0, "not known yet");
  el.loaded(24);
  assert.equal(player.duration, 24);
  assert.deepEqual(states, ["stopped"]);

  await player.play();
  assert.equal(player.state, "playing");
  assert.ok(player.playing);
  el.time = 5;
  assert.equal(player.position, 5);
  assert.equal(player.cursorAt(player.position).system, 1);

  player.pause();
  assert.equal(player.state, "paused");
  assert.equal(player.position, 5, "held where it paused");

  // a line clicked on: its time, the nearest when it is sung twice
  assert.equal(player.timeAt(0), 1.5);
  player.seek(11);
  assert.equal(player.timeAt(0), 12);
  assert.equal(player.timeAt(9), null);
  player.seek(99);
  assert.equal(el.time, 24, "not past the end");
  player.seek(-3);
  assert.equal(el.time, 0);

  player.seek(8);
  player.stop();
  assert.equal(player.state, "stopped");
  assert.equal(player.position, 0);
  assert.equal(el.time, 0);
  assert.equal(player.cursorAt(player.position), null);
  const n = states.length;
  player.stop();
  assert.equal(states.length, n, "stopping a stopped player says nothing");

  // a stopped player goes to a time paused
  player.seek(4.5);
  assert.equal(player.state, "paused");
  assert.equal(player.cursorAt(player.position).system, 1, "a cue takes effect at its time");

  // played to its end: stopped, at the start
  await player.play();
  el.end();
  assert.equal(player.state, "stopped");
  assert.equal(player.position, 0);
  assert.equal(ended, 1);
  assert.equal(states.at(-1), "stopped");
  await player.play();
  assert.equal(player.state, "playing");
});

test("AudioPlayer follows an element played by its own controls", () => {
  const el = new FakeAudio();
  const player = new AudioPlayer(new Blob([new Uint8Array(4)], { type: "audio/flac" }), "", null, { element: el });
  let updates = 0;
  player.onUpdate = () => updates++;
  el.play();
  assert.equal(player.state, "playing");
  el.time = 3;
  el.pause();
  assert.equal(player.state, "paused");
  assert.equal(player.position, 3);
  assert.equal(updates, 2);
  // no cues: no line, no time of a line
  assert.equal(player.cursorAt(3), null);
  assert.equal(player.timeAt(0), null);
});

test("AudioPlayer changes the speed", () => {
  const el = new FakeAudio();
  const player = new AudioPlayer(new Uint8Array(1), "audio/mpeg", null, { element: el });
  assert.equal(player.speedPercent, 100);
  player.setSpeedPercent(150);
  assert.equal(el.playbackRate, 1.5);
  assert.equal(el.defaultPlaybackRate, 1.5);
  assert.equal(player.speedPercent, 150);
  player.setSpeedPercent(5);
  assert.equal(player.speedPercent, 25);
  player.setSpeedPercent(1000);
  assert.equal(player.speedPercent, 200);
  player.setSpeedPercent(NaN);
  assert.equal(player.speedPercent, 200);
});

test("AudioPlayer reports a file the browser cannot play", async () => {
  const el = new FakeAudio();
  const player = new AudioPlayer(new Uint8Array(1), "audio/aiff", null, { element: el });
  const errors = [];
  player.onError = (e) => errors.push(e.message);
  el.error = { code: 4, message: "" };
  el.dispatchEvent(new Event("error"));
  assert.deepEqual(errors, ["this browser does not play audio of this format"]);
  assert.equal(player.error.message, errors[0]);
  assert.equal(player.state, "stopped");
  await assert.rejects(player.play(), /does not play audio of this format/);

  // a browser that refuses to start: the rejection is the caller's, and the player is not playing
  const el2 = new FakeAudio();
  const other = new AudioPlayer(new Uint8Array(1), "audio/mpeg", null, { element: el2 });
  el2.refuse = new Error("NotAllowedError");
  await assert.rejects(other.play(), /NotAllowedError/);
  assert.ok(!other.playing);
});

test("AudioPlayer.dispose lets the file and the element go", async () => {
  const el = new FakeAudio();
  const player = new AudioPlayer(new Uint8Array(1), "audio/mpeg", cues, { element: el });
  let updates = 0;
  player.onUpdate = () => updates++;
  await player.play();
  const url = el.src;
  player.dispose();
  assert.ok(el.paused);
  assert.equal(el.src, "");
  assert.equal(el.loads, 1);
  assert.equal(player.state, "stopped");
  const n = updates;
  el.play();
  assert.equal(updates, n, "no longer listening");
  // the blob is revoked: fetching its address fails
  await assert.rejects(fetch(url));
});
