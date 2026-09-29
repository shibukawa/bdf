import { createConverter } from "@bdfkit/convert";
export type { ConvertOptions, Converted, ConvertedPage, Opened, Format } from "@bdfkit/convert";
export { Converter, ConvertError } from "@bdfkit/convert";

/** PDF and EPUB only. The wasm asset is resolved by the consumer's bundler. */
export function createPdfEpubConverter() {
  return createConverter(new URL("./converter.wasm", import.meta.url));
}
