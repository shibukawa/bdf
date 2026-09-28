// Documents written by hand for the tests: what no writer writes (docs/spec.md §3.3, §5).
import { deflateRawSync } from "node:zlib";

export const cat = (...parts) => {
  const out = new Uint8Array(parts.reduce((n, p) => n + p.length, 0));
  let at = 0;
  for (const p of parts) {
    out.set(p, at);
    at += p.length;
  }
  return out;
};
export const bytes = (...b) => new Uint8Array(b);
export const ascii = (s) => new TextEncoder().encode(s);
export const varuint = (n) => {
  const out = [];
  for (let v = BigInt(n); ; ) {
    const b = Number(v & 0x7fn);
    v >>= 7n;
    if (v === 0n) return new Uint8Array([...out, b]);
    out.push(b | 0x80);
  }
};
export const f32 = (...vs) => new Uint8Array(new Float32Array(vs).buffer);
export const u16 = (v) => bytes(v & 0xff, v >> 8);
export const str = (s) => cat(varuint(ascii(s).length), ascii(s));
const hex = (h) => Uint8Array.from(h.match(/../g), (x) => parseInt(x, 16));

/** A part name: 32 hex characters that say i. */
export const name = (i) => i.toString(16).padStart(32, "0");

export const Op = {
  SAVE: 0x01, RESTORE: 0x02, FILTER: 0x19, FONT: 0x1a, FILL_RECT: 0x20, FILL_TEXT: 0x30, IMAGE: 0x40,
  USE: 0x50, GROUP_BEGIN: 0x52, GROUP_END: 0x53, MASK_BEGIN: 0x54, MASK_END: 0x55, MARK: 0x71,
};

/** A font of the system, as an entry of an object's font table. */
export const systemFont = (family = "sans-serif") => cat(bytes(1), str(family), u16(400), bytes(0));

/** An Object part; images and objects are part names. */
export function object({ strings = [], fonts = [], images = [], objects = [], ops = [] } = {}) {
  const code = cat(...ops);
  return cat(
    ascii("BOBJ"), u16(1), u16(0), f32(0, 0, 100, 100),
    varuint(strings.length), ...strings.map(str),
    varuint(0), // paths
    varuint(0), // paints
    varuint(fonts.length), ...fonts,
    varuint(images.length), ...images.map(hex),
    varuint(objects.length), ...objects.map(hex),
    varuint(code.length), code,
  );
}

/** A page of 100 by 100 units that draws an object. */
export const page = (obj) => ({ w: 100, h: 100, layers: [{ role: "body", obj }] });

/**
 * A single-file document. parts: [{ h, t, bytes, deflate, size }]; size is
 * the size the manifest states (default: the real one). The manifest is
 * deflated with deflateManifest.
 */
export function single(manifest, parts = [], { deflateManifest = false } = {}) {
  const stored = parts.map((p) => (p.deflate ? new Uint8Array(deflateRawSync(p.bytes, { level: 9 })) : p.bytes));
  let off = 0;
  const entries = parts.map((p, i) => {
    const e = { h: p.h, t: p.t, enc: p.deflate ? "deflate-raw" : "identity", len: stored[i].length, size: p.size ?? p.bytes.length, off };
    off += stored[i].length;
    return { ...e, ...p.entry };
  });
  let json = ascii(JSON.stringify({ bdf: 1, opset: 1, unit: "pt", views: [], ...manifest, parts: entries }));
  if (deflateManifest) json = new Uint8Array(deflateRawSync(json, { level: 1 }));
  const header = new Uint8Array(32);
  header.set([0x62, 0x64, 0x66, 0]);
  const dv = new DataView(header.buffer);
  dv.setUint16(4, 1, true);
  dv.setBigUint64(8, 32n, true);
  dv.setBigUint64(16, BigInt(json.length), true);
  header[24] = deflateManifest ? 1 : 0;
  return cat(header, json, ...stored);
}

/** Objects that each draw the next one fan times, depth deep, and a last one of leaf ops. */
export function fanOut(depth, fan, leaf) {
  const parts = [];
  for (let i = 0; i < depth; i++) {
    const ops = Array.from({ length: fan }, () => cat(bytes(Op.USE), varuint(0)));
    parts.push({ h: name(i + 1), t: "obj", bytes: object({ objects: [name(i + 2)], ops }), deflate: true });
  }
  parts.push({ h: name(depth + 1), t: "obj", bytes: object(leaf) });
  return single({ views: [{ id: "v", kind: "fixed", pages: [page(name(1))] }] }, parts);
}
