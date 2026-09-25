export * from "./types.js";
export * from "./opcodes.js";
export { ByteReader, BdfFormatError } from "./bytes.js";
export { decodeObject, decodePath, decodePathCollection, objectDeps, walk, NoopSink, opHistogram } from "./object.js";
export { parseHeader, decode, BufferSource, RangeSource, SplitSource, fetchSingle, HEADER_SIZE, type PartSource, type Header } from "./container.js";
export { BdfDocument } from "./document.js";
export { extractText, guessSep, Mark, Sep, type TextRun } from "./text.js";
export { decodeTextIndex, TextSearch, normalizeQuery, normalizeChar, type IndexRun, type SearchHit, type HitSegment, type SearchOptions } from "./search.js";
