// Playing the recording of a view (docs/spec.md §4.4) on the main thread.
// The browser's own player, an audio element, plays the file as the
// document stores it; the cues, whose ticks are milliseconds, tell which
// line of the page the recording is at: a line of the lyrics, a chapter.
// The cue logic is in plain functions, which run without a page.
import type { Cues } from "@bdfkit/core";
import type { Cursor, PlayState } from "./player.js";

/**
 * The line a recording is at, at a time in milliseconds (spec §4.4): the
 * system of the last cue at or before it (the later of cues at one time),
 * all of it. null before the first cue (the introduction) and without cues.
 */
export function lineAtTime(cues: Cues, ms: number): Cursor | null {
  const list = cues.cues;
  if (!list.length || !(ms >= list[0].tick)) return null;
  let lo = 0, hi = list.length - 1;
  while (lo < hi) {
    const mid = (lo + hi + 1) >> 1;
    if (list[mid].tick <= ms) lo = mid;
    else hi = mid - 1;
  }
  const system = list[lo].system, s = cues.systems[system];
  return { page: s.page, system, x: s.x, y: s.y, w: s.w, h: s.h };
}

/**
 * The time in milliseconds a line is at (where a click on it plays from):
 * that of its cue, of the one nearest to near when it has several. null
 * when no cue is on the system.
 */
export function timeOfLine(cues: Cues, system: number, near = 0): number | null {
  let best: number | null = null, far = Infinity;
  for (const c of cues.cues) {
    if (c.system !== system) continue;
    const d = Math.abs(c.tick - near);
    if (d < far) {
      best = c.tick;
      far = d;
    }
  }
  return best;
}

export interface AudioPlayerOptions {
  /**
   * The audio element that plays (default: a new one, outside the page). A
   * viewer that shows the browser's own controls passes an element of the
   * page that has them.
   */
  element?: HTMLAudioElement;
}

/**
 * Plays a view's recording with an audio element, and tells which line of
 * the pages it is at. It has the methods of MusicPlayer that a recording
 * has a use for, so a viewer drives both alike. play() should be called
 * from a user gesture (a click, a key) so that the browser lets it sound.
 */
export class AudioPlayer {
  readonly element: HTMLAudioElement;
  /** Called when playing starts, pauses, stops (also at the end) or jumps, and when the length or the speed is known or changes. */
  onUpdate: (() => void) | undefined;
  /** Called when the recording has played to its end; the player has stopped. */
  onEnd: (() => void) | undefined;
  /** Called when the browser cannot play the file: a format it has no decoder for, or a damaged file. */
  onError: ((e: Error) => void) | undefined;
  private readonly url: string;
  /** At the start, with no place to show: not yet played, stopped, or played to the end. */
  private idle = true;
  private failure: Error | undefined;
  private readonly listeners: [string, () => void][];

  /** data: the audio file (AudioData.blob, or its bytes); type: the media type of bytes. */
  constructor(data: Blob | ArrayBuffer | Uint8Array, type = "", readonly cues: Cues | null = null, options: AudioPlayerOptions = {}) {
    const blob = data instanceof Blob ? data : new Blob([data as BlobPart], { type });
    this.url = URL.createObjectURL(blob);
    const el = this.element = options.element ?? new Audio();
    const update = () => this.onUpdate?.();
    this.listeners = [
      // the element may be played and paused by its own controls, or by the keys of the keyboard and the headset
      ["play", () => { this.idle = false; update(); }],
      ["pause", update],
      ["seeked", update],
      ["durationchange", update],
      ["ratechange", update],
      ["ended", () => this.finish()],
      ["error", () => this.fail()],
    ];
    for (const [name, fn] of this.listeners) el.addEventListener(name, fn);
    el.preload = "metadata";
    el.src = this.url;
  }

  /** Seconds; 0 until the browser has read the length of the file. */
  get duration(): number {
    const d = this.element.duration;
    return Number.isFinite(d) && d > 0 ? d : 0;
  }
  get state(): PlayState { return this.playing ? "playing" : this.idle ? "stopped" : "paused"; }
  get playing(): boolean { return !this.element.paused && !this.element.ended; }
  /** Seconds from the start. */
  get position(): number { return this.idle && !this.playing ? 0 : this.element.currentTime; }
  /** The speed as a percentage of the recorded one. */
  get speedPercent(): number { return Math.round(this.element.playbackRate * 100); }
  /** Why the browser cannot play the file, once it has said so. */
  get error(): Error | undefined { return this.failure; }

  /** Play faster or slower (25 to 200 percent); the pitch stays. */
  setSpeedPercent(percent: number) {
    if (!Number.isFinite(percent)) return;
    const rate = Math.max(25, Math.min(200, Math.round(percent))) / 100;
    // the element keeps the rate over another load of its source
    this.element.defaultPlaybackRate = this.element.playbackRate = rate;
  }

  /** Play from the position (from the start once the recording has played to its end). */
  play(): Promise<void> {
    if (this.failure) return Promise.reject(this.failure);
    this.idle = false;
    return this.element.play();
  }

  pause() {
    this.element.pause();
  }

  /** Stop, back at the start. */
  stop() {
    if (this.state === "stopped") return;
    this.idle = true;
    this.element.pause();
    this.rewind();
    this.onUpdate?.();
  }

  /** Go to a time (seconds); a stopped player is paused there. */
  seek(seconds: number) {
    const end = this.duration;
    this.idle = false;
    this.element.currentTime = Math.max(0, end > 0 ? Math.min(end, seconds) : seconds);
    this.onUpdate?.();
  }

  /** The line of the pages the recording is at, at a time (seconds); null without cues, and before the first. */
  cursorAt(seconds: number): Cursor | null {
    // half a millisecond: the rounding of seconds must not keep a cue from taking effect at its time
    return this.cues ? lineAtTime(this.cues, seconds * 1000 + 0.5) : null;
  }

  /** The time (seconds) of a line, for playing from a line clicked on (see timeOfLine); null when it has none. */
  timeAt(system: number): number | null {
    if (!this.cues) return null;
    const ms = timeOfLine(this.cues, system, this.position * 1000);
    return ms === null ? null : ms / 1000;
  }

  /** Stop and let the file go. */
  dispose() {
    const el = this.element;
    for (const [name, fn] of this.listeners) el.removeEventListener(name, fn);
    el.pause();
    el.removeAttribute("src");
    el.load();
    URL.revokeObjectURL(this.url);
    this.idle = true;
  }

  private rewind() {
    try {
      this.element.currentTime = 0;
    } catch {
      // an element that has nothing loaded yet is at the start
    }
  }

  /** Played to the end. */
  private finish() {
    this.idle = true;
    this.rewind();
    this.onUpdate?.();
    this.onEnd?.();
  }

  private fail() {
    const e = this.element.error;
    // MEDIA_ERR_SRC_NOT_SUPPORTED (4) and MEDIA_ERR_DECODE (3): the file, not the network
    const why = e?.code === 4 ? "this browser does not play audio of this format" : e?.code === 3 ? "the audio could not be decoded" : "the audio could not be loaded";
    this.failure = new Error(e?.message ? `${why} (${e.message})` : why);
    this.idle = true;
    this.onError?.(this.failure);
    this.onUpdate?.();
  }
}
