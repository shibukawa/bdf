export type Hash = string; // 32 lowercase hex characters

export interface Manifest {
  bdf: number;
  /**
   * Present only in the outer manifest of an encrypted document (spec §3.5),
   * which has no views: BdfDocument.open reads the sealed manifest with a password.
   */
  encryption?: Encryption;
  opset: number;
  unit: string;
  meta?: Meta;
  views: View[];
  parts: PartEntry[];
}

/** How an encrypted document is sealed (spec §3.5). */
export interface Encryption {
  cipher: "A256GCM";
  keys: KeySlot[];
  /** The sealed part that holds the manifest, and its encoding inside the seal. */
  manifest: { part: Hash; enc: Encoding };
}

/** The content key, wrapped with AES-KW under a key derived from a password. */
export interface KeySlot {
  type: "password";
  kdf: "PBKDF2-SHA256";
  iter: number;
  /** base64 */
  salt: string;
  /** base64: the wrapped content key */
  key: string;
}

export interface Meta {
  /** Dublin Core description of the document (spec §4.3). */
  dc?: DublinCore;
  /** Input format the document was converted from ("pdf", "pptx", ...). */
  source?: string;
  generator?: string;
}

/** A Dublin Core element: one value, or an array when the element repeats. Read it with dcValues(). */
export type DCValue = string | string[];

/** The fifteen Dublin Core Metadata Element Set 1.1 elements plus the DCMI Terms created and modified. */
export interface DublinCore {
  title?: DCValue;
  creator?: DCValue;
  subject?: DCValue;
  description?: DCValue;
  publisher?: DCValue;
  contributor?: DCValue;
  date?: DCValue;
  type?: DCValue;
  format?: DCValue;
  identifier?: DCValue;
  source?: DCValue;
  language?: DCValue;
  relation?: DCValue;
  coverage?: DCValue;
  rights?: DCValue;
  created?: DCValue;
  modified?: DCValue;
}

/**
 * fixed: pages; flow: pages with body rectangles, readable as one continuous
 * scroll; sheet: an unbounded plane of tiles; scroll: one long column without
 * pages, stored as strips that are shown stacked (docs/spec.md §4.1).
 */
export type ViewKind = "fixed" | "flow" | "sheet" | "scroll";

export interface View {
  id: string;
  kind: ViewKind;
  title?: string;
  textIndex?: Hash;
  pages?: Page[];
  continuous?: { gap: number };
  tile?: number;
  cols?: [number, number][];
  rows?: [number, number][];
  freeze?: { cols?: number; rows?: number };
  gridlines?: boolean;
  tiles?: Record<string, Hash>;
  tilesRef?: Hash;
}

export interface Page {
  w: number;
  h: number;
  body?: RectDef;
  layers: Layer[];
}

export interface RectDef { x: number; y: number; w: number; h: number }

export interface Layer { role: string; obj: Hash }

export type PartType = "obj" | "font" | "img" | "path" | "idx" | "sealed";
export type Encoding = "identity" | "deflate-raw";

export interface PartEntry {
  h: Hash;
  t: PartType;
  enc: Encoding;
  len: number;
  size: number;
  off?: number;
  /** Encrypted documents: the outer part that holds this part sealed. */
  sealed?: Hash;
}

export interface PathData { verbs: Uint8Array; args: Float32Array }

export type PathEntry = { inline: PathData } | { hash: Hash; index: number };

export interface Stop { offset: number; color: number }

export interface Paint {
  kind: number;
  coords: Float32Array; // gradients
  stops: Stop[];
  image: number; // pattern: image ref
  repeat: number;
  matrix: Float32Array;
}

export interface Font {
  kind: number;
  hash?: Hash;
  family: string;
  weight: number;
  style: number;
}

export interface Rect { x: number; y: number; w: number; h: number }

/** A decoded Object part. */
export interface ObjectPart {
  opset: number;
  bbox: Rect;
  strings: string[];
  paths: PathEntry[];
  paints: Paint[];
  fonts: Font[];
  images: Hash[];
  objects: Hash[];
  ops: Uint8Array;
}

export interface Glyph { path: number; x: number; y: number }

/** Receiver of decoded instructions. Every method has a no-op default in NoopSink. */
export interface OpSink {
  save(): void;
  restore(): void;
  transform(a: number, b: number, c: number, d: number, e: number, f: number): void;
  translate(x: number, y: number): void;
  scale(x: number, y: number): void;
  clipPath(path: number, rule: number): void;
  clipRect(x: number, y: number, w: number, h: number): void;
  fillColor(rgba: number): void;
  fillPaint(paint: number): void;
  strokeColor(rgba: number): void;
  strokePaint(paint: number): void;
  line(width: number, cap: number, join: number, miter: number): void;
  dash(segments: Float32Array, offset: number): void;
  alpha(a: number): void;
  blend(mode: number): void;
  shadow(rgba: number, blur: number, dx: number, dy: number): void;
  filter(css: string): void;
  font(font: number, size: number): void;
  textStyle(align: number, baseline: number, dir: number, letterSpacing: number): void;
  fillRect(x: number, y: number, w: number, h: number): void;
  strokeRect(x: number, y: number, w: number, h: number): void;
  fillPath(path: number, rule: number): void;
  strokePath(path: number): void;
  fillPathAt(path: number, rule: number, x: number, y: number): void;
  fillPathRun(rule: number, glyphs: Glyph[]): void;
  clearRect(x: number, y: number, w: number, h: number): void;
  fillText(text: string, x: number, y: number, advance: number): void;
  strokeText(text: string, x: number, y: number, advance: number): void;
  image(img: number, x: number, y: number, w: number, h: number): void;
  imageSub(img: number, sx: number, sy: number, sw: number, sh: number, dx: number, dy: number, dw: number, dh: number): void;
  smoothing(enabled: boolean, quality: number): void;
  use(obj: number): void;
  useAt(obj: number, x: number, y: number): void;
  groupBegin(alpha: number, blend: number, x: number, y: number, w: number, h: number): void;
  groupEnd(): void;
  link(x: number, y: number, w: number, h: number, url: string): void;
  mark(kind: number, payload: string): void;
  ext(payload: Uint8Array): void;
}
