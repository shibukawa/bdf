import { BdfFormatError } from "./bytes.js";
import { decode, type PartSource } from "./container.js";
import type { Encryption, Hash, Manifest, PartEntry } from "./types.js";

// Encrypted documents (spec §3.5), decrypted with WebCrypto: PBKDF2 derives
// the key-encryption key from the password, AES-KW unwraps the content key
// and AES-GCM opens each sealed part. The content key is not extractable.

/** Iteration counts above this are refused, so a file cannot keep a reader busy for long. */
export const MAX_ITERATIONS = 10_000_000;

const NONCE_SIZE = 12;
const TAG_SIZE = 16;
const MANIFEST_AAD = new TextEncoder().encode("manifest");
const utf8 = new TextDecoder();

/**
 * Thrown when an encrypted document is opened without a password
 * ("required") or with one that opens none of its key slots ("wrong").
 */
export class BdfPasswordError extends Error {
  constructor(readonly reason: "required" | "wrong") {
    super(reason === "required" ? "bdf: the document is encrypted and needs a password" : "bdf: wrong password");
    this.name = "BdfPasswordError";
  }
}

function base64(s: string): Uint8Array<ArrayBuffer> {
  const bin = atob(s);
  const out = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
  return out;
}

function hexBytes(h: Hash): Uint8Array<ArrayBuffer> {
  const out = new Uint8Array(h.length / 2);
  for (let i = 0; i < out.length; i++) out[i] = parseInt(h.slice(2 * i, 2 * i + 2), 16);
  return out;
}

/** The content key a password unwraps from one of the key slots. */
async function contentKey(enc: Encryption, password: string): Promise<CryptoKey> {
  if (enc.cipher !== "A256GCM") throw new BdfFormatError(`unknown cipher ${enc.cipher}`);
  const pw = new TextEncoder().encode(password.normalize("NFC"));
  for (const slot of enc.keys) {
    if (slot.type !== "password") continue;
    if (slot.kdf !== "PBKDF2-SHA256") throw new BdfFormatError(`unknown key derivation ${slot.kdf}`);
    if (!Number.isInteger(slot.iter) || slot.iter < 1 || slot.iter > MAX_ITERATIONS) throw new BdfFormatError(`iteration count ${slot.iter} out of range`);
    if (pw.length === 0) break; // writers refuse empty passwords, and some browsers refuse empty PBKDF2 keys
    const base = await crypto.subtle.importKey("raw", pw, "PBKDF2", false, ["deriveKey"]);
    const kek = await crypto.subtle.deriveKey(
      { name: "PBKDF2", hash: "SHA-256", salt: base64(slot.salt), iterations: slot.iter },
      base, { name: "AES-KW", length: 256 }, false, ["unwrapKey"]);
    try {
      return await crypto.subtle.unwrapKey("raw", base64(slot.key), kek, "AES-KW", { name: "AES-GCM", length: 256 }, false, ["decrypt"]);
    } catch {
      // the integrity check of the unwrap failed: not this slot's password
    }
  }
  throw new BdfPasswordError("wrong");
}

/** Decrypt a sealed part: a nonce, then the ciphertext and its tag. */
async function open(key: CryptoKey, sealed: Uint8Array, aad: Uint8Array<ArrayBuffer>): Promise<Uint8Array> {
  if (sealed.length < NONCE_SIZE + TAG_SIZE) throw new BdfFormatError("sealed part too short");
  // WebCrypto takes ArrayBuffer-backed views; copy only what is not one
  const bytes = sealed.buffer instanceof ArrayBuffer ? (sealed as Uint8Array<ArrayBuffer>) : sealed.slice();
  try {
    return new Uint8Array(await crypto.subtle.decrypt({ name: "AES-GCM", iv: bytes.subarray(0, NONCE_SIZE), additionalData: aad }, key, bytes.subarray(NONCE_SIZE)));
  } catch {
    throw new BdfFormatError("sealed part fails authentication");
  }
}

/**
 * The parts of an encrypted document, decrypted, over the source of its
 * stored (sealed) parts. Made by unlocking with a password; BdfDocument.open
 * does that when given one.
 */
export class SealedSource implements PartSource {
  private constructor(
    private readonly source: PartSource,
    private readonly outer: Map<Hash, PartEntry>,
    private readonly key: CryptoKey,
    private readonly inner: Manifest,
  ) {}

  /** Unlock an encrypted document; throws BdfPasswordError("wrong") when the password opens no key slot. */
  static async unlock(source: PartSource, password: string): Promise<SealedSource> {
    const stored = await source.manifest();
    const enc = stored.encryption;
    if (!enc) throw new Error("bdf: the document is not encrypted");
    const key = await contentKey(enc, password);
    const outer = new Map(stored.parts.map((e) => [e.h, e]));
    const me = outer.get(enc.manifest.part);
    if (!me) throw new BdfFormatError("sealed manifest part missing");
    const json = await decode(await open(key, await source.stored(me), MANIFEST_AAD), enc.manifest.enc);
    const inner = JSON.parse(utf8.decode(json)) as Manifest;
    if (inner.encryption) throw new BdfFormatError("sealed manifest is encrypted again");
    for (const e of inner.parts) {
      if (!e.sealed || !outer.has(e.sealed)) throw new BdfFormatError(`part ${e.h}: sealed part missing`);
    }
    return new SealedSource(source, outer, key, inner);
  }

  manifest(): Promise<Manifest> {
    return Promise.resolve(this.inner);
  }

  async stored(e: PartEntry): Promise<Uint8Array> {
    const o = e.sealed && this.outer.get(e.sealed);
    if (!o) throw new BdfFormatError(`part ${e.h}: sealed part missing`);
    return open(this.key, await this.source.stored(o), hexBytes(e.h));
  }
}
