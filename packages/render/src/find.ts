// A viewer's own answer to the browser's find shortcut. A browser finds the
// text a page holds, and a viewer holds the text of the pages near the
// visible area only (their text layers, or none at all when it draws
// bitmaps alone); the worker's search() reads the text of the whole view.
// findKey tells which keys ask for a search, FindHits keeps the hits of one
// and locates them as their pages are shown, and FindBar is the small bar a
// viewer without a search box of its own shows over its pages.
import type { SearchHit, View } from "@bdfkit/core";
import type { BdfWorkerClient } from "./client.js";
import type { HitRect } from "./search.js";

/** What a key asks for: the search field, or the match after or before the one shown. */
export type FindKey = "open" | "next" | "previous";

/** Whether this platform's shortcuts take Command (macOS, iOS) rather than Control. */
function applePlatform(): boolean {
  if (typeof navigator === "undefined") return false;
  const nav = navigator as Navigator & { userAgentData?: { platform?: string } };
  return /mac|iphone|ipad|ipod/i.test(nav.userAgentData?.platform || nav.platform || nav.userAgent || "");
}

/** The Latin letter of a key: the one it types, or the one of its place on a keyboard that types another script. */
function letter(e: KeyboardEvent): string {
  if (/^[a-z]$/i.test(e.key)) return e.key.toLowerCase();
  return /^Key[A-Z]$/.test(e.code) && !/^[\x20-\x7e]$/.test(e.key) ? e.code[3].toLowerCase() : "";
}

/**
 * The find shortcut a key event is, if any: Command+F on macOS and Control+F
 * elsewhere open the search, Command/Control+G and F3 go to the next match
 * and, with Shift, to the one before. Control+F on macOS is left alone: it
 * moves the caret in a text field.
 */
export function findKey(e: KeyboardEvent, apple = applePlatform()): FindKey | undefined {
  if (e.isComposing || e.altKey) return undefined;
  if (e.key === "F3" && !e.metaKey && !e.ctrlKey) return e.shiftKey ? "previous" : "next";
  if (apple ? !e.metaKey || e.ctrlKey : !e.ctrlKey || e.metaKey) return undefined;
  const key = letter(e);
  if (key === "f") return e.shiftKey ? undefined : "open";
  if (key === "g") return e.shiftKey ? "previous" : "next";
  return undefined;
}

/** The hits a search keeps; one more is asked for, to tell when there are more. */
export const FIND_LIMIT = 1000;

/**
 * The hits of a search in a view. Finding them reads the text of the view
 * alone; where a hit lies on its page takes the page itself (and its fonts),
 * so hits are located when asked for, a page (a tile of a sheet) at a time:
 * the pages a viewer shows, and the page of the hit it goes to.
 */
export class FindHits {
  /** Whether the view has more hits than these. */
  readonly more: boolean;
  private readonly places = new Map<string, number[]>();
  private readonly located: (HitRect[] | undefined)[] = [];
  private readonly asked: (Promise<void> | undefined)[] = [];

  constructor(private readonly client: Pick<BdfWorkerClient, "locate">, readonly view: View, readonly query: string, readonly hits: SearchHit[], more = false) {
    this.more = more;
    hits.forEach((hit, i) => {
      for (const s of hit.segments) {
        const key = this.key(s.a, s.b);
        const list = this.places.get(key);
        if (!list) this.places.set(key, [i]);
        else if (list[list.length - 1] !== i) list.push(i);
      }
    });
  }

  /** Search a view. A query of nothing but spaces finds nothing. */
  static async search(client: Pick<BdfWorkerClient, "search" | "locate">, view: View, query: string, limit = FIND_LIMIT): Promise<FindHits> {
    const hits = query.trim() ? await client.search(view.id, query, { limit: limit + 1 }) : [];
    return new FindHits(client, view, query, hits.slice(0, limit), hits.length > limit);
  }

  get length(): number { return this.hits.length; }

  /** How locate() names the place of text: a page, or a tile of a sheet ("x,y"). */
  private key(a: number, b: number): string {
    return this.view.kind === "sheet" ? `${a},${b}` : String(a);
  }

  /** The rectangles of a hit, in page or sheet coordinates; undefined until it is located. */
  rects(i: number): HitRect[] | undefined {
    return this.located[i];
  }

  /** Call fn with each rectangle of the located hits; i is the hit's index. */
  each(fn: (rect: HitRect, i: number) => void): void {
    this.located.forEach((rects, i) => rects?.forEach((r) => fn(r, i)));
  }

  /**
   * The first hit on a page or after it, for a search that starts where the
   * reader is; the first of all when there is none (the search wraps), and
   * in a sheet, whose hits follow its tiles rather than its rows.
   */
  from(page: number): number {
    if (this.view.kind === "sheet") return 0;
    const i = this.hits.findIndex((h) => h.segments[0].a >= page);
    return i < 0 ? 0 : i;
  }

  /**
   * Locate the hits on pages (on tiles of a sheet, named "x,y") that are
   * not located yet. Resolves once they all are, to whether any had to be.
   */
  async locate(places: Iterable<number | string>): Promise<boolean> {
    const todo = new Set<number>();
    const waits = new Set<Promise<void>>();
    for (const place of places) {
      for (const i of this.places.get(String(place)) ?? []) {
        if (this.located[i]) continue;
        if (this.asked[i]) waits.add(this.asked[i]!);
        else todo.add(i);
      }
    }
    if (todo.size) {
      const list = [...todo];
      const batch = this.client.locate(this.view.id, list.map((i) => this.hits[i])).then((rects) => {
        list.forEach((i, k) => { this.located[i] = rects[k] ?? []; });
      }).finally(() => {
        // located by now, or to be asked for again
        for (const i of list) this.asked[i] = undefined;
      });
      for (const i of list) this.asked[i] = batch;
      waits.add(batch);
    }
    await Promise.all(waits);
    return waits.size > 0;
  }

  /** Locate a hit, with the others of its page, and give its rectangles. */
  async locateHit(i: number): Promise<HitRect[]> {
    const s = this.hits[i]?.segments[0];
    if (!s) return [];
    if (!this.located[i]) await this.locate([this.key(s.a, s.b)]);
    return this.located[i] ?? [];
  }
}

/** The words of a find bar. */
export interface FindBarLabels {
  /** The name and the placeholder of the field. */
  find: string;
  previous: string;
  next: string;
  close: string;
  /** Shown when a query finds nothing. */
  none: string;
  /** What the count says at match index (from 1) of total; more: the view has more matches than total. */
  count: (index: number, total: number, more: boolean) => string;
}

const LABELS: FindBarLabels = {
  find: "Find in document", previous: "Previous match", next: "Next match", close: "Close", none: "No matches",
  count: (index, total, more) => `${index} / ${total}${more ? "+" : ""}`,
};

export interface FindBarOptions {
  /**
   * The reader typed (step 0: show the match nearest the place being read),
   * or asked for the next match (1) or the one before (-1).
   */
  onFind: (query: string, step: -1 | 0 | 1) => void;
  /** The bar closed: the matches are no longer shown. */
  onClose: () => void;
  labels?: Partial<FindBarLabels>;
}

/** What a find bar says of a search. */
export interface FindResult {
  /** The match shown (from 0), -1 for none. */
  index: number;
  total: number;
  /** The view has more matches than total. */
  more?: boolean;
  /** Read by screen readers after the count: where the match is and the text around it. */
  detail?: string;
}

/** Typing searches once it pauses this long (ms). */
const TYPING_PAUSE = 150;

const BAR_CSS = "position:absolute;top:8px;right:24px;z-index:10;display:none;align-items:center;gap:4px;max-width:calc(100% - 48px);box-sizing:border-box;"
  + "padding:4px 6px;border:1px solid #b8b8b8;border-radius:6px;background:#fff;color:#222;box-shadow:0 2px 10px rgba(0,0,0,.25);"
  + "font:13px/1.4 system-ui,sans-serif;color-scheme:light;text-align:left;direction:ltr";
const FIELD_CSS = "flex:1 1 12em;width:12em;min-width:4em;box-sizing:border-box;margin:0;padding:3px 6px;border:1px solid #b8b8b8;border-radius:4px;background:#fff;color:#222;font:inherit";
const COUNT_CSS = "flex:none;min-width:4.5em;text-align:center;color:#555;white-space:nowrap;font-variant-numeric:tabular-nums";
const BUTTON_CSS = "flex:none;width:26px;height:26px;margin:0;padding:0;border:0;border-radius:4px;background:transparent;color:#333;font:inherit;line-height:1;cursor:pointer";
const HIDDEN_CSS = "position:absolute;width:1px;height:1px;overflow:hidden;clip-path:inset(50%);white-space:nowrap";

/**
 * A search field with the count of matches and buttons for the match
 * before and after, shown over the top right corner of a viewer (put
 * element in something positioned). Enter goes to the next match,
 * Shift+Enter to the one before, Escape closes. It is styled through the
 * style property of its elements, which a content security policy without
 * inline styles allows.
 */
export class FindBar {
  readonly element: HTMLDivElement;
  readonly input: HTMLInputElement;
  private readonly status: HTMLSpanElement;
  private readonly labels: FindBarLabels;
  private timer: ReturnType<typeof setTimeout> | undefined;

  constructor(private readonly options: FindBarOptions) {
    const labels = this.labels = { ...LABELS, ...options.labels };
    const bar = this.element = document.createElement("div");
    bar.style.cssText = BAR_CSS;
    bar.setAttribute("role", "search");
    const input = this.input = document.createElement("input");
    input.type = "search";
    input.placeholder = labels.find;
    input.setAttribute("aria-label", labels.find);
    input.autocomplete = "off";
    input.spellcheck = false;
    input.style.cssText = FIELD_CSS;
    this.status = document.createElement("span");
    this.status.style.cssText = COUNT_CSS;
    this.status.setAttribute("role", "status");
    const button = (text: string, label: string, click: () => void) => {
      const b = document.createElement("button");
      b.type = "button";
      b.textContent = text;
      b.title = label;
      b.setAttribute("aria-label", label);
      b.style.cssText = BUTTON_CSS;
      b.onpointerenter = () => { b.style.background = "#e8e8e8"; };
      b.onpointerleave = () => { b.style.background = "transparent"; };
      b.onclick = click;
      return b;
    };
    bar.append(input, this.status, button("▲", labels.previous, () => this.step(-1)), button("▼", labels.next, () => this.step(1)), button("✕", labels.close, () => this.close()));

    // the text of a composition (an input method) is searched once it is entered
    input.addEventListener("input", (e) => { if (!(e as InputEvent).isComposing) this.typed(); });
    input.addEventListener("compositionend", () => this.typed());
    bar.addEventListener("keydown", (e) => {
      // 229: the key that enters a composition, in the browsers that end it before the key is told
      if (e.isComposing || e.keyCode === 229) return;
      if (e.key === "Escape") {
        // the bar closes, not what holds the viewer (a lightbox)
        e.preventDefault();
        e.stopPropagation();
        this.close();
      } else if (e.key === "Enter" && e.target === input) {
        e.preventDefault();
        this.step(e.shiftKey ? -1 : 1);
      }
    });
  }

  /** Whether the bar is shown. */
  get shown(): boolean { return this.element.style.display !== "none"; }
  /** What the field holds. */
  get query(): string { return this.input.value; }
  /** Whether the reader is in the field. */
  get focused(): boolean {
    const root = this.input.getRootNode() as Document | ShadowRoot;
    return root.activeElement === this.input;
  }

  /** Show the bar with the reader in its field, what it holds selected: typing replaces it. */
  open(): void {
    this.element.style.display = "flex";
    this.input.focus({ preventScroll: true });
    this.input.select();
  }

  /** Hide the bar and tell the viewer, which shows the matches no longer. */
  close(): void {
    if (!this.shown) return;
    this.hide();
    this.options.onClose();
  }

  /** Hide the bar and forget its query, without a word to the viewer: another document or view is shown. */
  reset(): void {
    this.hide();
    this.input.value = "";
  }

  private hide(): void {
    clearTimeout(this.timer);
    this.element.style.display = "none";
    this.status.textContent = "";
  }

  /** Say what a search found; nothing for no search. */
  result(result?: FindResult): void {
    if (!result) { this.status.textContent = ""; return; }
    if (!result.total) { this.status.textContent = this.labels.none; return; }
    this.status.textContent = this.labels.count(result.index + 1, result.total, !!result.more);
    if (result.detail) {
      const detail = document.createElement("span");
      detail.style.cssText = HIDDEN_CSS;
      detail.textContent = ` ${result.detail}`;
      this.status.append(detail);
    }
  }

  private typed(): void {
    clearTimeout(this.timer);
    // an emptied field clears the matches at once
    if (!this.input.value) { this.options.onFind("", 0); return; }
    this.timer = setTimeout(() => this.options.onFind(this.input.value, 0), TYPING_PAUSE);
  }

  private step(delta: -1 | 1): void {
    clearTimeout(this.timer);
    this.options.onFind(this.input.value, delta);
  }
}
