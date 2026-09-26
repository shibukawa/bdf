/** Opcodes of opset 1 (docs/spec.md §7). */
export const Op = {
  SAVE: 0x01, RESTORE: 0x02, TRANSFORM: 0x03, TRANSLATE: 0x04, SCALE: 0x05,
  CLIP_PATH: 0x06, CLIP_RECT: 0x07,
  FILL_COLOR: 0x10, FILL_PAINT: 0x11, STROKE_COLOR: 0x12, STROKE_PAINT: 0x13,
  LINE: 0x14, DASH: 0x15, ALPHA: 0x16, BLEND: 0x17, SHADOW: 0x18, FILTER: 0x19,
  FONT: 0x1a, TEXT_STYLE: 0x1b,
  FILL_RECT: 0x20, STROKE_RECT: 0x21, FILL_PATH: 0x22, STROKE_PATH: 0x23,
  FILL_PATH_AT: 0x24, FILL_PATH_RUN: 0x25, CLEAR_RECT: 0x26,
  FILL_TEXT: 0x30, STROKE_TEXT: 0x31,
  IMAGE: 0x40, IMAGE_SUB: 0x41, SMOOTHING: 0x42,
  USE: 0x50, USE_AT: 0x51, GROUP_BEGIN: 0x52, GROUP_END: 0x53, MASK_BEGIN: 0x54, MASK_END: 0x55,
  LINK: 0x70, MARK: 0x71, EXT: 0xff,
} as const;

export const OPSET_VERSION = 1;
export const FORMAT_VERSION = 1;

export const OP_NAMES: Record<number, string> = Object.fromEntries(
  Object.entries(Op).map(([k, v]) => [v, k]),
);

/** Soft mask kinds of MASK_BEGIN (docs/spec.md §7.6). */
export const MaskKind = { ALPHA: 0, LUMINOSITY: 1 } as const;

export const BLEND_NAMES: GlobalCompositeOperation[] = [
  "source-over", "multiply", "screen", "overlay", "darken", "lighten",
  "color-dodge", "color-burn", "hard-light", "soft-light", "difference",
  "exclusion", "hue", "saturation", "color", "luminosity",
  "destination-over", "destination-in", "destination-out", "source-in",
  "source-out", "source-atop", "destination-atop", "xor", "copy", "lighter",
];

export const LINE_CAPS: CanvasLineCap[] = ["butt", "round", "square"];
export const LINE_JOINS: CanvasLineJoin[] = ["miter", "round", "bevel"];
export const TEXT_ALIGNS: CanvasTextAlign[] = ["left", "right", "center", "start", "end"];
export const TEXT_BASELINES: CanvasTextBaseline[] = ["alphabetic", "top", "middle", "bottom", "hanging", "ideographic"];
export const TEXT_DIRECTIONS: CanvasDirection[] = ["inherit", "ltr", "rtl"];
export const FILL_RULES: CanvasFillRule[] = ["nonzero", "evenodd"];
export const REPEATS = ["repeat", "repeat-x", "repeat-y", "no-repeat"];
export const SMOOTHING_QUALITIES: ImageSmoothingQuality[] = ["low", "medium", "high"];
export const FONT_STYLES = ["normal", "italic", "oblique"];

/** Path verbs (docs/spec.md §6.3). */
export const Verb = { MOVE: 0, LINE: 1, QUAD: 2, CUBIC: 3, CLOSE: 4, RECT: 5, ELLIPSE: 6, ARC_TO: 7, ROUND_RECT: 8 } as const;
export const VERB_ARGS = [2, 2, 4, 6, 0, 4, 8, 5, 5];

export const PaintKind = { LINEAR: 0, RADIAL: 1, CONIC: 2, PATTERN: 3 } as const;
export const FontKind = { EMBEDDED: 0, SYSTEM: 1 } as const;
