// A page turning over, drawn with WebGL 2 for the book layouts of the demo
// viewer (book.ts). The turning page (the sheet) is a grid of vertices held at
// the spine. Pulling its corner to a point folds it along the perpendicular
// bisector of the corner and that point. The paper beyond the fold wraps
// around a cylinder lying on the page, and beyond the cylinder it lies flat
// again, upside down (the flap). The vertex shader bends the grid this way.
// The fragment shader shades the curl by its normal and shows the back page
// where the paper faces down.
//
// Coordinates are CSS px of the canvas, in a frame where the sheet lies
// beyond the spine (x >= spine) and turns to the other side. A page of a book
// read right to left, or one turning back, is the mirror image of that, which
// the shaders mirror on the way to the screen.

export interface Point { x: number; y: number }
export interface Rect { x: number; y: number; w: number; h: number }

/** A page lying flat, and whose bitmap it shows (a page index). */
export interface Leaf { rect: Rect; page: number }

export interface Scene {
  spine: number;
  /** The scene is shown mirrored at the spine (a book read right to left, or a page turning back). */
  mirror: boolean;
  /** The turning page lying flat, beyond the spine. */
  sheet: Rect;
  front: number;
  /** The page on the back of the sheet; null: blank paper (one page at a time). */
  back: number | null;
  /** The page before the spine, which the sheet covers as it lands. */
  before?: Leaf;
  /** The page under the sheet. */
  under?: Leaf;
  /** Opacity of the sheet where its back is in view. */
  backAlpha: number;
}

/**
 * The fold: its normal n points from the flat part of the sheet toward the
 * lifted one, l is a point on it, r is the radius of the curl, and progress
 * goes from 0 (flat on its side) to 1 (turned over).
 */
export interface Fold { n: Point; l: Point; r: number; progress: number }

/** Radius of the curl halfway through a turn, in page widths. */
const CURL = 0.09;

/**
 * Where a corner pulled toward p can go, for a sheet held at the spine: a
 * corner stays within the page width of its end of the spine and the
 * diagonal of the other end; an edge pulled at a height stays at it.
 */
export function reach(sheet: Rect, spine: number, corner: Point, edge: boolean, p: Point): Point {
  const w = sheet.x + sheet.w - spine;
  if (edge) return { x: Math.min(spine + w, Math.max(spine - w, p.x)), y: corner.y };
  const top = { x: spine, y: sheet.y }, bottom = { x: spine, y: sheet.y + sheet.h };
  const atTop = corner.y <= sheet.y + sheet.h / 2;
  const q = within(p, atTop ? top : bottom, w);
  return within(q, atTop ? bottom : top, Math.hypot(w, sheet.h));
}

function within(p: Point, c: Point, r: number): Point {
  const dx = p.x - c.x, dy = p.y - c.y, d = Math.hypot(dx, dy);
  return d <= r ? p : { x: c.x + (dx * r) / d, y: c.y + (dy * r) / d };
}

/** The fold that brings a corner of the sheet to p; undefined while it is where it was. */
export function foldAt(sheet: Rect, spine: number, corner: Point, p: Point): Fold | undefined {
  const dx = corner.x - p.x, dy = corner.y - p.y, dist = Math.hypot(dx, dy);
  if (dist < 0.5) return undefined;
  const n = { x: dx / dist, y: dy / dist };
  const w = sheet.x + sheet.w - spine;
  const progress = Math.min(1, dist / (2 * w));
  // the ends of the spine stay on the flat side: a curl of radius r moves
  // the fold PI*r/2 past the bisector
  const mx = (corner.x + p.x) / 2, my = (corner.y + p.y) / 2;
  const spare = Math.min((mx - spine) * n.x + (my - sheet.y) * n.y, (mx - spine) * n.x + (my - sheet.y - sheet.h) * n.y);
  const r = Math.max(0, Math.min(w * CURL * Math.sin(Math.PI * progress), (2 * spare) / Math.PI));
  // a corner on the flap lands (dist + PI*r)/2 beyond the fold
  const dc = (dist + Math.PI * r) / 2;
  return { n, l: { x: corner.x - n.x * dc, y: corner.y - n.y * dc }, r, progress };
}

const COMMON = `#version 300 es
precision highp float;
const float PI = 3.14159265;
uniform vec4 uView;    // canvas width and height (CSS px), spine, 1 (or -1: mirrored)
uniform vec4 uFold;    // fold normal, a point on the fold
uniform float uR;      // radius of the curl
uniform float uFolded; // 1 while the sheet is folded
uniform vec4 uSheet;   // the sheet lying flat
uniform vec4 uShadow;  // shadow along the curl: strength, length; of the flap: strength, softness

vec4 clip(vec2 p, float z) {
  float x = uView.z + uView.w * (p.x - uView.z);
  return vec4(x / uView.x * 2.0 - 1.0, 1.0 - p.y / uView.y * 2.0, -z / 8192.0, 1.0);
}

float box(vec2 p, vec4 r) {
  vec2 q = abs(p - r.xy - r.zw * 0.5) - r.zw * 0.5;
  return length(max(q, 0.0)) + min(max(q.x, q.y), 0.0);
}

// The shadow the folded sheet casts at q, on what lies under it.
float curlShadow(vec2 q) {
  if (uFolded < 0.5) return 0.0;
  vec2 n = uFold.xy;
  float d = dot(q - uFold.zw, n);
  if (d >= 0.0) {
    // beyond the fold, along the curl standing over it
    float along = box(q - n * d, uSheet);
    float t = max(d - uR, 0.0) / uShadow.y;
    return uShadow.x * (1.0 - smoothstep(0.0, 1.0, t)) * (1.0 - smoothstep(0.0, uShadow.w, along));
  }
  // before it, around the flap: the far part of the sheet mirrored over the
  // line PI*r/2 beyond the fold
  vec2 p = q - 2.0 * (d - PI * uR * 0.5) * n;
  return uShadow.z * (1.0 - smoothstep(0.0, uShadow.w, box(p, uSheet)));
}
`;

const QUAD_VS = `${COMMON}
in vec2 aUV;
uniform vec4 uRect;
out vec2 vUV;
void main() {
  vUV = aUV;
  gl_Position = clip(uRect.xy + aUV * uRect.zw, 0.0);
}`;

const QUAD_FS = `${COMMON}
uniform sampler2D uTex;
in vec2 vUV;
out vec4 outColor;
void main() {
  outColor = texture(uTex, vec2(uView.w < 0.0 ? 1.0 - vUV.x : vUV.x, vUV.y));
}`;

// Shadows over the whole canvas: of pages lying flat (as their elements'
// box-shadow, 0 1px 4px rgba(0,0,0,.3)) and of the folded sheet.
const SHADOW_VS = `${COMMON}
in vec2 aUV;
void main() {
  gl_Position = vec4(aUV * 4.0 - 1.0, 0.0, 1.0);
}`;

const SHADOW_FS = `${COMMON}
uniform vec4 uBox;       // a page lying flat (w 0: none)
uniform float uSheetBox; // 1: also the flat part of the sheet
uniform float uCurl;     // 1: the shadow of the fold
uniform vec2 uPixels;    // canvas height in device px, device px per CSS px
out vec4 outColor;
float boxShadow(vec2 q, vec4 r) {
  return r.z <= 0.0 ? 0.0 : 0.3 * (1.0 - smoothstep(-4.0, 4.0, box(q - vec2(0.0, 1.0), r)));
}
void main() {
  vec2 css = vec2(gl_FragCoord.x, uPixels.x - gl_FragCoord.y) / uPixels.y;
  vec2 q = vec2(uView.z + uView.w * (css.x - uView.z), css.y);
  float a = boxShadow(q, uBox);
  if (uSheetBox > 0.5) {
    float lying = uFolded > 0.5 ? 1.0 - smoothstep(-2.0, 2.0, dot(q - uFold.zw, uFold.xy)) : 1.0;
    a = max(a, boxShadow(q, uSheet) * lying);
  }
  if (uCurl > 0.5) a = max(a, curlShadow(q));
  outColor = vec4(0.0, 0.0, 0.0, a);
}`;

const SHEET_VS = `${COMMON}
in vec2 aUV;
out vec2 vUV;
out float vTheta; // angle around the curl: 0 flat, PI turned over
out vec2 vPos;
void main() {
  vec2 p = uSheet.xy + aUV * uSheet.zw;
  vec2 pos = p;
  float z = 0.0, th = 0.0;
  if (uFolded > 0.5) {
    vec2 n = uFold.xy;
    float d = dot(p - uFold.zw, n);
    if (d > 0.0) {
      if (uR < 0.01 || d > PI * uR) {
        pos = p - n * (2.0 * d - PI * uR);
        z = 2.0 * uR + 1.0;
        th = PI;
      } else {
        th = d / uR;
        pos = p + n * (uR * sin(th) - d);
        z = uR * (1.0 - cos(th));
      }
    }
  }
  vUV = aUV;
  vTheta = th;
  vPos = pos;
  gl_Position = clip(pos, z);
}`;

const SHEET_FS = `${COMMON}
uniform sampler2D uFront;
uniform sampler2D uBack;
uniform float uHasBack;
uniform float uBackAlpha;
in vec2 vUV;
in float vTheta;
in vec2 vPos;
out vec4 outColor;
const vec3 PAPER = vec3(0.99, 0.988, 0.982);
void main() {
  bool back = vTheta > PI * 0.5;
  float s = uView.w < 0.0 ? 1.0 - vUV.x : vUV.x;
  vec4 front = texture(uFront, vec2(s, vUV.y));
  vec3 c = front.rgb;
  float a = 1.0;
  if (back) {
    // a page alone has a blank back, with the front showing through a little
    c = uHasBack > 0.5 ? texture(uBack, vec2(1.0 - s, vUV.y)).rgb : mix(PAPER, front.rgb, 0.06);
    a = uBackAlpha;
  }
  // light from the top left of the screen on the side in view: 1 where flat
  float th = min(vTheta, PI);
  vec3 N = vec3(-uFold.xy * sin(th), cos(th));
  if (back) N = -N;
  N.x *= uView.w;
  vec3 L = normalize(vec3(-0.3, -0.5, 1.0));
  vec3 H = normalize(L + vec3(0.0, 0.0, 1.0));
  float spec = max(pow(max(dot(N, H), 0.0), 24.0) - pow(H.z, 24.0), 0.0);
  float shade = clamp(0.3 + 0.7 * dot(N, L) / L.z, 0.35, 1.1) + 0.25 * spec;
  // the flap's shadow on the flat part
  if (vTheta <= 0.0) shade *= 1.0 - curlShadow(vPos);
  outColor = vec4(c * shade * a, a);
}`;

interface Program { prog: WebGLProgram; loc: Map<string, WebGLUniformLocation | null> }

/** Columns and rows of the sheet's grid. */
const GRID_X = 48, GRID_Y = 64;

export class FlipRenderer {
  private readonly quad: Program;
  private readonly shadow: Program;
  private readonly sheet: Program;
  private readonly quadVao: WebGLVertexArrayObject;
  private readonly gridVao: WebGLVertexArrayObject;
  private readonly gridCount: number;
  private readonly textures = new Map<number, { tex: WebGLTexture; src: ImageBitmap }>();
  private readonly blank: WebGLTexture;
  private width = 1;
  private height = 1;
  /** Whether the context is usable (it can be lost, as when the GPU resets). */
  ok = true;

  /** A renderer on canvas; undefined without WebGL 2. */
  static create(canvas: HTMLCanvasElement): FlipRenderer | undefined {
    const gl = canvas.getContext("webgl2", { alpha: true, premultipliedAlpha: true, antialias: true, depth: true });
    if (!gl) return undefined;
    try {
      return new FlipRenderer(canvas, gl);
    } catch (e) {
      console.warn("bdf viewer: no page turning:", e);
      return undefined;
    }
  }

  private constructor(private readonly canvas: HTMLCanvasElement, private readonly gl: WebGL2RenderingContext) {
    canvas.addEventListener("webglcontextlost", () => { this.ok = false; });
    this.quad = this.program(QUAD_VS, QUAD_FS);
    this.shadow = this.program(SHADOW_VS, SHADOW_FS);
    this.sheet = this.program(SHEET_VS, SHEET_FS);
    this.quadVao = this.vao(new Float32Array([0, 0, 1, 0, 0, 1, 1, 1]));
    const uv: number[] = [];
    for (let y = 0; y <= GRID_Y; y++) for (let x = 0; x <= GRID_X; x++) uv.push(x / GRID_X, y / GRID_Y);
    const idx: number[] = [];
    for (let y = 0; y < GRID_Y; y++) {
      for (let x = 0; x < GRID_X; x++) {
        const i = y * (GRID_X + 1) + x;
        idx.push(i, i + 1, i + GRID_X + 1, i + 1, i + GRID_X + 2, i + GRID_X + 1);
      }
    }
    this.gridVao = this.vao(new Float32Array(uv), new Uint16Array(idx));
    this.gridCount = idx.length;
    this.blank = this.texture();
    gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, new Uint8Array([255, 255, 255, 255]));
  }

  private program(vs: string, fs: string): Program {
    const gl = this.gl;
    const shader = (type: number, src: string) => {
      const s = gl.createShader(type)!;
      gl.shaderSource(s, src);
      gl.compileShader(s);
      if (!gl.getShaderParameter(s, gl.COMPILE_STATUS)) throw new Error(gl.getShaderInfoLog(s) ?? "shader");
      return s;
    };
    const prog = gl.createProgram()!;
    gl.attachShader(prog, shader(gl.VERTEX_SHADER, vs));
    gl.attachShader(prog, shader(gl.FRAGMENT_SHADER, fs));
    gl.bindAttribLocation(prog, 0, "aUV");
    gl.linkProgram(prog);
    if (!gl.getProgramParameter(prog, gl.LINK_STATUS)) throw new Error(gl.getProgramInfoLog(prog) ?? "program");
    return { prog, loc: new Map() };
  }

  private vao(uv: Float32Array, idx?: Uint16Array): WebGLVertexArrayObject {
    const gl = this.gl;
    const vao = gl.createVertexArray()!;
    gl.bindVertexArray(vao);
    gl.bindBuffer(gl.ARRAY_BUFFER, gl.createBuffer());
    gl.bufferData(gl.ARRAY_BUFFER, uv, gl.STATIC_DRAW);
    gl.enableVertexAttribArray(0);
    gl.vertexAttribPointer(0, 2, gl.FLOAT, false, 0, 0);
    if (idx) {
      gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.createBuffer());
      gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, idx, gl.STATIC_DRAW);
    }
    gl.bindVertexArray(null);
    return vao;
  }

  private texture(): WebGLTexture {
    const gl = this.gl;
    const tex = gl.createTexture()!;
    gl.bindTexture(gl.TEXTURE_2D, tex);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    return tex;
  }

  /** Size the canvas: CSS px, and device px per CSS px. */
  resize(width: number, height: number, dpr: number) {
    this.width = Math.max(1, width);
    this.height = Math.max(1, height);
    this.canvas.width = Math.max(1, Math.round(width * dpr));
    this.canvas.height = Math.max(1, Math.round(height * dpr));
    this.canvas.style.width = `${width}px`;
    this.canvas.style.height = `${height}px`;
  }

  /** Give a page its bitmap (again only when it is another one). */
  upload(page: number, bmp: ImageBitmap) {
    if (this.textures.get(page)?.src === bmp) return;
    const gl = this.gl;
    const tex = this.textures.get(page)?.tex ?? this.texture();
    gl.bindTexture(gl.TEXTURE_2D, tex);
    gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, bmp);
    this.textures.set(page, { tex, src: bmp });
  }

  forget(page: number) {
    const t = this.textures.get(page);
    if (!t) return;
    this.gl.deleteTexture(t.tex);
    this.textures.delete(page);
  }

  forgetAll() {
    for (const page of [...this.textures.keys()]) this.forget(page);
  }

  private use(p: Program, scene: Scene, fold: Fold | undefined) {
    const gl = this.gl;
    gl.useProgram(p.prog);
    const s = scene.sheet;
    const k = fold ? Math.sin(Math.PI * fold.progress) : 0;
    const r = fold?.r ?? 0;
    this.set(p, "uView", [this.width, this.height, scene.spine, scene.mirror ? -1 : 1]);
    this.set(p, "uFold", fold ? [fold.n.x, fold.n.y, fold.l.x, fold.l.y] : [1, 0, 0, 0]);
    this.set(p, "uR", r);
    this.set(p, "uFolded", fold ? 1 : 0);
    this.set(p, "uSheet", [s.x, s.y, s.w, s.h]);
    this.set(p, "uShadow", [0.35 * k, s.w * 0.15 * k + r + 1, 0.3 * Math.sqrt(k), 2 + 1.2 * r + 8 * k]);
  }

  private set(p: Program, name: string, v: number | number[]) {
    const gl = this.gl;
    const loc = this.uniform(p, name);
    if (typeof v === "number") gl.uniform1f(loc, v);
    else if (v.length === 2) gl.uniform2fv(loc, v);
    else gl.uniform4fv(loc, v);
  }

  private bind(unit: number, page: number | null) {
    const gl = this.gl;
    gl.activeTexture(gl.TEXTURE0 + unit);
    gl.bindTexture(gl.TEXTURE_2D, (page !== null && this.textures.get(page)?.tex) || this.blank);
  }

  draw(scene: Scene, fold: Fold | undefined) {
    const gl = this.gl;
    if (!this.ok) return;
    gl.viewport(0, 0, this.canvas.width, this.canvas.height);
    gl.clearColor(0, 0, 0, 0);
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
    gl.enable(gl.BLEND);
    gl.blendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA);
    gl.disable(gl.DEPTH_TEST);
    // in the order of the page elements: the page before the spine, then those beyond it
    const shadow = (box: Rect | undefined, sheetBox: boolean, curl: boolean) => {
      this.use(this.shadow, scene, fold);
      this.set(this.shadow, "uBox", box ? [box.x, box.y, box.w, box.h] : [0, 0, 0, 0]);
      this.set(this.shadow, "uSheetBox", sheetBox ? 1 : 0);
      this.set(this.shadow, "uCurl", curl ? 1 : 0);
      this.set(this.shadow, "uPixels", [this.canvas.height, this.canvas.width / this.width]);
      gl.bindVertexArray(this.quadVao);
      gl.drawArrays(gl.TRIANGLES, 0, 3);
    };
    const leaf = (l: Leaf | undefined) => {
      if (!l) return;
      this.use(this.quad, scene, fold);
      this.set(this.quad, "uRect", [l.rect.x, l.rect.y, l.rect.w, l.rect.h]);
      this.bind(0, l.page);
      gl.uniform1i(this.uniform(this.quad, "uTex"), 0);
      gl.bindVertexArray(this.quadVao);
      gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);
    };
    if (scene.before) shadow(scene.before.rect, false, false);
    leaf(scene.before);
    shadow(scene.under?.rect, true, false);
    leaf(scene.under);
    if (fold) shadow(undefined, false, true);

    this.use(this.sheet, scene, fold);
    this.bind(0, scene.front);
    this.bind(1, scene.back);
    gl.uniform1i(this.uniform(this.sheet, "uFront"), 0);
    gl.uniform1i(this.uniform(this.sheet, "uBack"), 1);
    this.set(this.sheet, "uHasBack", scene.back !== null ? 1 : 0);
    this.set(this.sheet, "uBackAlpha", scene.backAlpha);
    gl.enable(gl.DEPTH_TEST);
    gl.depthFunc(gl.LESS);
    gl.bindVertexArray(this.gridVao);
    gl.drawElements(gl.TRIANGLES, this.gridCount, gl.UNSIGNED_SHORT, 0);
    gl.bindVertexArray(null);
  }

  private uniform(p: Program, name: string): WebGLUniformLocation | null {
    if (!p.loc.has(name)) p.loc.set(name, this.gl.getUniformLocation(p.prog, name));
    return p.loc.get(name)!;
  }

  destroy() {
    this.forgetAll();
    this.gl.getExtension("WEBGL_lose_context")?.loseContext();
  }
}
