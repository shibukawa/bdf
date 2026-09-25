import { type BdfDocument, type ObjectPart, type PathData, type PartEntry, type Hash, type Font, Verb, FontKind, FONT_STYLES } from "@bdf/core";

/** Build a Path2D from path data. */
export function buildPath2D(p: PathData): Path2D {
  const path = new Path2D();
  const a = p.args;
  let i = 0;
  for (const v of p.verbs) {
    switch (v) {
      case Verb.MOVE: path.moveTo(a[i], a[i + 1]); i += 2; break;
      case Verb.LINE: path.lineTo(a[i], a[i + 1]); i += 2; break;
      case Verb.QUAD: path.quadraticCurveTo(a[i], a[i + 1], a[i + 2], a[i + 3]); i += 4; break;
      case Verb.CUBIC: path.bezierCurveTo(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4], a[i + 5]); i += 6; break;
      case Verb.CLOSE: path.closePath(); break;
      case Verb.RECT: path.rect(a[i], a[i + 1], a[i + 2], a[i + 3]); i += 4; break;
      case Verb.ELLIPSE: path.ellipse(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4], a[i + 5], a[i + 6], a[i + 7] !== 0); i += 8; break;
      case Verb.ARC_TO: path.arcTo(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4]); i += 5; break;
      case Verb.ROUND_RECT: path.roundRect(a[i], a[i + 1], a[i + 2], a[i + 3], a[i + 4]); i += 5; break;
    }
  }
  return path;
}

/** Font family name used for an embedded font part. */
export function embeddedFamily(hash: Hash): string {
  return `bdf-${hash}`;
}

/** CSS font string for a FONT instruction. */
export function fontString(f: Font, size: number): string {
  let family: string;
  if (f.kind === FontKind.EMBEDDED && f.hash) {
    family = `"${embeddedFamily(f.hash)}"`;
    if (f.family) family += `, ${f.family}`;
  } else {
    family = f.family || "sans-serif";
  }
  return `${FONT_STYLES[f.style] ?? "normal"} ${f.weight || 400} ${size}px ${family}`;
}

/**
 * Caches of browser objects derived from parts: Path2D, ImageBitmap, FontFace.
 * One cache can serve many contexts; it is bound to a document.
 */
export class ResourceCache {
  private inlinePaths = new WeakMap<ObjectPart, (Path2D | undefined)[]>();
  private extPaths = new Map<Hash, Path2D[]>();
  private images = new Map<Hash, ImageBitmap>();
  private fonts = new Map<Hash, FontFace>();
  private fontSet: FontFaceSet | undefined;

  constructor(readonly doc: BdfDocument, fontSet?: FontFaceSet) {
    this.fontSet = fontSet ?? (globalThis as { fonts?: FontFaceSet }).fonts ?? (globalThis as { document?: { fonts?: FontFaceSet } }).document?.fonts;
  }

  /** Load an object, its children, and every font/image/path part they use. */
  async prepare(hash: Hash): Promise<ObjectPart> {
    return this.doc.ensure(hash, (e, bytes) => this.load(e, bytes));
  }

  private async load(e: PartEntry, bytes: Uint8Array): Promise<void> {
    switch (e.t) {
      case "img": {
        if (this.images.has(e.h)) return;
        const bmp = await createImageBitmap(new Blob([bytes as BlobPart]));
        this.images.set(e.h, bmp);
        return;
      }
      case "font": {
        if (this.fonts.has(e.h)) return;
        const buffer = bytes.buffer.slice(bytes.byteOffset, bytes.byteOffset + bytes.byteLength) as ArrayBuffer;
        const face = new FontFace(embeddedFamily(e.h), buffer);
        await face.load();
        this.fontSet?.add(face);
        this.fonts.set(e.h, face);
        return;
      }
      case "path": {
        if (this.extPaths.has(e.h)) return;
        const paths = await this.doc.pathCollection(e.h);
        this.extPaths.set(e.h, paths.map(buildPath2D));
        return;
      }
    }
  }

  path(o: ObjectPart, index: number): Path2D {
    const entry = o.paths[index];
    if (!entry) throw new Error(`bdf: bad path ref ${index}`);
    if ("inline" in entry) {
      let list = this.inlinePaths.get(o);
      if (!list) {
        list = new Array(o.paths.length);
        this.inlinePaths.set(o, list);
      }
      let p = list[index];
      if (!p) {
        p = buildPath2D(entry.inline);
        list[index] = p;
      }
      return p;
    }
    const set = this.extPaths.get(entry.hash);
    const p = set?.[entry.index];
    if (!p) throw new Error(`bdf: path collection ${entry.hash} not loaded`);
    return p;
  }

  image(hash: Hash): ImageBitmap {
    const img = this.images.get(hash);
    if (!img) throw new Error(`bdf: image ${hash} not loaded`);
    return img;
  }

  object(hash: Hash): ObjectPart {
    return this.doc.objectSync(hash);
  }
}
