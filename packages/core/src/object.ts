import { ByteReader, BdfFormatError } from "./bytes.js";
import { Op, OPSET_VERSION, OP_NAMES, PaintKind, FontKind, VERB_ARGS } from "./opcodes.js";
import type { ObjectPart, PathData, PathEntry, Paint, Font, Glyph, OpSink } from "./types.js";

const PAINT_COORDS = [4, 6, 3];

export function decodePath(r: ByteReader): PathData {
  const n = r.varuint();
  const verbs = new Uint8Array(r.bytesN(n));
  let total = 0;
  for (const v of verbs) {
    if (v >= VERB_ARGS.length) throw new BdfFormatError(`unknown path verb ${v}`);
    total += VERB_ARGS[v];
  }
  return { verbs, args: r.f32array(total) };
}

export function decodePathCollection(bytes: Uint8Array): PathData[] {
  const r = new ByteReader(bytes);
  const n = r.varuint();
  const out: PathData[] = [];
  for (let i = 0; i < n; i++) out.push(decodePath(r));
  return out;
}

function decodePaint(r: ByteReader): Paint {
  const kind = r.u8();
  if (kind === PaintKind.PATTERN) {
    const image = r.varuint();
    const repeat = r.u8();
    return { kind, coords: new Float32Array(0), stops: [], image, repeat, matrix: r.f32array(6) };
  }
  if (kind >= PAINT_COORDS.length) throw new BdfFormatError(`unknown paint kind ${kind}`);
  const coords = r.f32array(PAINT_COORDS[kind]);
  const n = r.varuint();
  const stops = [];
  for (let i = 0; i < n; i++) stops.push({ offset: r.f32(), color: r.u32() });
  return { kind, coords, stops, image: 0, repeat: 0, matrix: new Float32Array(0) };
}

function decodeFont(r: ByteReader): Font {
  const kind = r.u8();
  const hash = kind === FontKind.EMBEDDED ? r.hash() : undefined;
  const family = r.str();
  const weight = r.u16();
  const style = r.u8();
  return { kind, hash, family, weight, style };
}

/** Parse an Object part. */
export function decodeObject(bytes: Uint8Array): ObjectPart {
  const r = new ByteReader(bytes);
  const magic = r.bytesN(4);
  if (magic[0] !== 0x42 || magic[1] !== 0x4f || magic[2] !== 0x42 || magic[3] !== 0x4a) {
    throw new BdfFormatError("not an object part");
  }
  const opset = r.u16();
  if (opset > OPSET_VERSION) throw new BdfFormatError(`unsupported opset ${opset}`);
  r.u16();
  const bbox = { x: r.f32(), y: r.f32(), w: r.f32(), h: r.f32() };
  const strings: string[] = [];
  for (let i = 0, n = r.varuint(); i < n; i++) strings.push(r.str());
  const paths: PathEntry[] = [];
  for (let i = 0, n = r.varuint(); i < n; i++) {
    if (r.u8() === 0) paths.push({ inline: decodePath(r) });
    else {
      const hash = r.hash();
      paths.push({ hash, index: r.varuint() });
    }
  }
  const paints: Paint[] = [];
  for (let i = 0, n = r.varuint(); i < n; i++) paints.push(decodePaint(r));
  const fonts: Font[] = [];
  for (let i = 0, n = r.varuint(); i < n; i++) fonts.push(decodeFont(r));
  const images: string[] = [];
  for (let i = 0, n = r.varuint(); i < n; i++) images.push(r.hash());
  const objects: string[] = [];
  for (let i = 0, n = r.varuint(); i < n; i++) objects.push(r.hash());
  const ops = r.bytesN(r.varuint());
  return { opset, bbox, strings, paths, paints, fonts, images, objects, ops };
}

/** Hashes of parts an object depends on directly. */
export function objectDeps(o: ObjectPart): string[] {
  const out: string[] = [];
  for (const f of o.fonts) if (f.hash) out.push(f.hash);
  out.push(...o.images);
  for (const p of o.paths) if ("hash" in p) out.push(p.hash);
  out.push(...o.objects);
  return out;
}

/** Execute the op stream of an object against a sink. */
export function walk(o: ObjectPart, sink: OpSink): void {
  const r = new ByteReader(o.ops);
  const S = o.strings;
  const str = () => {
    const i = r.varuint();
    if (i >= S.length) throw new BdfFormatError("bad string ref");
    return S[i];
  };
  while (!r.eof) {
    const code = r.u8();
    switch (code) {
      case Op.SAVE: sink.save(); break;
      case Op.RESTORE: sink.restore(); break;
      case Op.TRANSFORM: sink.transform(r.f32(), r.f32(), r.f32(), r.f32(), r.f32(), r.f32()); break;
      case Op.TRANSLATE: sink.translate(r.f32(), r.f32()); break;
      case Op.SCALE: sink.scale(r.f32(), r.f32()); break;
      case Op.CLIP_PATH: sink.clipPath(r.varuint(), r.u8()); break;
      case Op.CLIP_RECT: sink.clipRect(r.f32(), r.f32(), r.f32(), r.f32()); break;
      case Op.FILL_COLOR: sink.fillColor(r.u32()); break;
      case Op.FILL_PAINT: sink.fillPaint(r.varuint()); break;
      case Op.STROKE_COLOR: sink.strokeColor(r.u32()); break;
      case Op.STROKE_PAINT: sink.strokePaint(r.varuint()); break;
      case Op.LINE: sink.line(r.f32(), r.u8(), r.u8(), r.f32()); break;
      case Op.DASH: {
        const n = r.varuint();
        const segs = r.f32array(n);
        sink.dash(segs, r.f32());
        break;
      }
      case Op.ALPHA: sink.alpha(r.f32()); break;
      case Op.BLEND: sink.blend(r.u8()); break;
      case Op.SHADOW: sink.shadow(r.u32(), r.f32(), r.f32(), r.f32()); break;
      case Op.FILTER: sink.filter(str()); break;
      case Op.FONT: sink.font(r.varuint(), r.f32()); break;
      case Op.TEXT_STYLE: sink.textStyle(r.u8(), r.u8(), r.u8(), r.f32()); break;
      case Op.FILL_RECT: sink.fillRect(r.f32(), r.f32(), r.f32(), r.f32()); break;
      case Op.STROKE_RECT: sink.strokeRect(r.f32(), r.f32(), r.f32(), r.f32()); break;
      case Op.FILL_PATH: sink.fillPath(r.varuint(), r.u8()); break;
      case Op.STROKE_PATH: sink.strokePath(r.varuint()); break;
      case Op.FILL_PATH_AT: sink.fillPathAt(r.varuint(), r.u8(), r.f32(), r.f32()); break;
      case Op.FILL_PATH_RUN: {
        const rule = r.u8();
        const n = r.varuint();
        const glyphs: Glyph[] = new Array(n);
        for (let i = 0; i < n; i++) glyphs[i] = { path: r.varuint(), x: r.f32(), y: r.f32() };
        sink.fillPathRun(rule, glyphs);
        break;
      }
      case Op.CLEAR_RECT: sink.clearRect(r.f32(), r.f32(), r.f32(), r.f32()); break;
      case Op.FILL_TEXT: sink.fillText(str(), r.f32(), r.f32(), r.f32()); break;
      case Op.STROKE_TEXT: sink.strokeText(str(), r.f32(), r.f32(), r.f32()); break;
      case Op.IMAGE: sink.image(r.varuint(), r.f32(), r.f32(), r.f32(), r.f32()); break;
      case Op.IMAGE_SUB:
        sink.imageSub(r.varuint(), r.f32(), r.f32(), r.f32(), r.f32(), r.f32(), r.f32(), r.f32(), r.f32());
        break;
      case Op.SMOOTHING: sink.smoothing(r.u8() !== 0, r.u8()); break;
      case Op.USE: sink.use(r.varuint()); break;
      case Op.USE_AT: sink.useAt(r.varuint(), r.f32(), r.f32()); break;
      case Op.GROUP_BEGIN: sink.groupBegin(r.f32(), r.u8(), r.f32(), r.f32(), r.f32(), r.f32()); break;
      case Op.GROUP_END: sink.groupEnd(); break;
      case Op.LINK: sink.link(r.f32(), r.f32(), r.f32(), r.f32(), str()); break;
      case Op.MARK: sink.mark(r.u8(), str()); break;
      case Op.EXT: sink.ext(r.bytesN(r.u32())); break;
      default:
        throw new BdfFormatError(`unknown opcode 0x${code.toString(16)} at ${r.pos - 1}`);
    }
  }
}

/** A sink that ignores everything; extend it and override what you need. */
export class NoopSink implements OpSink {
  save(): void {}
  restore(): void {}
  transform(_a: number, _b: number, _c: number, _d: number, _e: number, _f: number): void {}
  translate(_x: number, _y: number): void {}
  scale(_x: number, _y: number): void {}
  clipPath(_path: number, _rule: number): void {}
  clipRect(_x: number, _y: number, _w: number, _h: number): void {}
  fillColor(_rgba: number): void {}
  fillPaint(_paint: number): void {}
  strokeColor(_rgba: number): void {}
  strokePaint(_paint: number): void {}
  line(_width: number, _cap: number, _join: number, _miter: number): void {}
  dash(_segments: Float32Array, _offset: number): void {}
  alpha(_a: number): void {}
  blend(_mode: number): void {}
  shadow(_rgba: number, _blur: number, _dx: number, _dy: number): void {}
  filter(_css: string): void {}
  font(_font: number, _size: number): void {}
  textStyle(_align: number, _baseline: number, _dir: number, _letterSpacing: number): void {}
  fillRect(_x: number, _y: number, _w: number, _h: number): void {}
  strokeRect(_x: number, _y: number, _w: number, _h: number): void {}
  fillPath(_path: number, _rule: number): void {}
  strokePath(_path: number): void {}
  fillPathAt(_path: number, _rule: number, _x: number, _y: number): void {}
  fillPathRun(_rule: number, _glyphs: Glyph[]): void {}
  clearRect(_x: number, _y: number, _w: number, _h: number): void {}
  fillText(_text: string, _x: number, _y: number, _advance: number): void {}
  strokeText(_text: string, _x: number, _y: number, _advance: number): void {}
  image(_img: number, _x: number, _y: number, _w: number, _h: number): void {}
  imageSub(_img: number, _sx: number, _sy: number, _sw: number, _sh: number, _dx: number, _dy: number, _dw: number, _dh: number): void {}
  smoothing(_enabled: boolean, _quality: number): void {}
  use(_obj: number): void {}
  useAt(_obj: number, _x: number, _y: number): void {}
  groupBegin(_alpha: number, _blend: number, _x: number, _y: number, _w: number, _h: number): void {}
  groupEnd(): void {}
  link(_x: number, _y: number, _w: number, _h: number, _url: string): void {}
  mark(_kind: number, _payload: string): void {}
  ext(_payload: Uint8Array): void {}
}

/** Counts instructions by opcode name; useful for tests and tooling. */
export function opHistogram(o: ObjectPart): Record<string, number> {
  const counts: Record<string, number> = {};
  const sink = new Proxy(new NoopSink(), {
    get(target, prop) {
      return (..._args: unknown[]) => {
        const name = String(prop);
        counts[name] = (counts[name] ?? 0) + 1;
      };
    },
  }) as OpSink;
  walk(o, sink);
  return counts;
}

export { OP_NAMES };
