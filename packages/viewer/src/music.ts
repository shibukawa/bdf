/** Optional music support. Import this entry only when piano roll or playback is needed. */
import { parseSmf } from "@bdfkit/core";
import type { MusicRenderContext } from "./index.js";

export { MusicPlayer, cursorAtTick, tickAt } from "@bdfkit/render";

export async function pianoRoll({ stage, renderer, view, zoom, isCurrent }: MusicRenderContext): Promise<void> {
  const play = await renderer.play(view.id);
  if (!isCurrent() || !play) return;
  const notes = parseSmf(new Uint8Array(play.seq)).events.filter((event) => event.type === "note");
  const width = Math.min(16000, Math.max(stage.clientWidth, Math.ceil(notes.reduce((end, note) => Math.max(end, note.time + note.duration), 0) * 64 * zoom)));
  const height = 128 * 12 * zoom;
  const pixelRatio = Math.min(devicePixelRatio, Math.sqrt(32_000_000 / Math.max(1, width * height)));
  const canvas = document.createElement("canvas");
  canvas.width = Math.ceil(width * pixelRatio);
  canvas.height = Math.ceil(height * pixelRatio);
  canvas.style.cssText = `display:block;width:${width}px;height:${height}px;background:#fff`;
  canvas.setAttribute("aria-label", `Piano roll with ${notes.length} notes`);
  const ctx = canvas.getContext("2d")!;
  ctx.scale(pixelRatio, pixelRatio);
  for (let key = 0; key < 128; key++) {
    const y = (127 - key) * 12 * zoom;
    ctx.fillStyle = [1, 3, 6, 8, 10].includes(key % 12) ? "#f1f1f1" : "#fff";
    ctx.fillRect(0, y, width, 12 * zoom);
  }
  ctx.fillStyle = "#3569a3";
  for (const note of notes) ctx.fillRect(note.time * 64 * zoom, (127 - note.key) * 12 * zoom + 1, Math.max(2, note.duration * 64 * zoom), Math.max(2, 12 * zoom - 2));
  if (isCurrent()) stage.replaceChildren(canvas);
}
