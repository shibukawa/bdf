/** Little-endian reader over a Uint8Array. */
export class ByteReader {
  readonly view: DataView;
  pos = 0;
  private readonly utf8 = new TextDecoder();

  constructor(readonly bytes: Uint8Array) {
    this.view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
  }

  get eof(): boolean { return this.pos >= this.bytes.length; }

  private need(n: number): void {
    if (this.pos + n > this.bytes.length) throw new BdfFormatError(`unexpected end of data at ${this.pos}`);
  }

  u8(): number { this.need(1); return this.bytes[this.pos++]; }
  u16(): number { this.need(2); const v = this.view.getUint16(this.pos, true); this.pos += 2; return v; }
  u32(): number { this.need(4); const v = this.view.getUint32(this.pos, true); this.pos += 4; return v; }
  u64(): number {
    this.need(8);
    const v = this.view.getBigUint64(this.pos, true);
    this.pos += 8;
    if (v > BigInt(Number.MAX_SAFE_INTEGER)) throw new BdfFormatError("u64 out of range");
    return Number(v);
  }
  f32(): number { this.need(4); const v = this.view.getFloat32(this.pos, true); this.pos += 4; return v; }

  varuint(): number {
    let result = 0;
    let shift = 0;
    for (;;) {
      const b = this.u8();
      if (shift < 28) {
        result |= (b & 0x7f) << shift;
      } else {
        result += (b & 0x7f) * 2 ** shift;
      }
      if ((b & 0x80) === 0) break;
      shift += 7;
      if (shift > 63) throw new BdfFormatError("bad varuint");
    }
    return result >>> 0 === result ? result >>> 0 : result;
  }

  bytesN(n: number): Uint8Array {
    this.need(n);
    const v = this.bytes.subarray(this.pos, this.pos + n);
    this.pos += n;
    return v;
  }

  str(): string {
    const n = this.varuint();
    return this.utf8.decode(this.bytesN(n));
  }

  hash(): string {
    const b = this.bytesN(16);
    let s = "";
    for (let i = 0; i < 16; i++) s += b[i].toString(16).padStart(2, "0");
    return s;
  }

  f32array(n: number): Float32Array {
    const out = new Float32Array(n);
    for (let i = 0; i < n; i++) out[i] = this.f32();
    return out;
  }
}

export class BdfFormatError extends Error {
  constructor(message: string) {
    super(`bdf: ${message}`);
    this.name = "BdfFormatError";
  }
}
