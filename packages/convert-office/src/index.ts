import { createConverter } from "@bdfkit/convert";
export type { ConvertOptions, Converted, ConvertedPage, Opened, Format } from "@bdfkit/convert";
export { Converter, ConvertError } from "@bdfkit/convert";

/** Word, Excel, PowerPoint, Visio, CSV, TSV, Parquet, PDF and EPUB. */
export function createOfficeConverter() {
  return createConverter(new URL("./converter.wasm", import.meta.url));
}
