import { createConverter } from "@bdfkit/convert";
export type { ConvertOptions, Converted, ConvertedPage, Opened, Format } from "@bdfkit/convert";
export { Converter, ConvertError } from "@bdfkit/convert";

/** Every browser converter in this repository. */
export function createAllConverter() {
  return createConverter(new URL("./converter.wasm", import.meta.url));
}
