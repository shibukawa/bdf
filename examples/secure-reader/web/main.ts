// The front end of secure-reader: opens the book ?book= names with the mini
// viewer (examples/miniviewer), as segments the server seals for each
// request (@bdf/core's SegmentLoader, in the rendering worker). The page
// holds no key and no page of the book: the worker asks for ten pages at a
// time as they come into view, and opens each answer with a key pair it
// made for that request and does not keep.
import { dcValues } from "@bdf/core";
import { BdfWorkerError } from "@bdf/render";
import { MiniViewer } from "../../miniviewer/miniviewer.js";

const book = new URLSearchParams(location.search).get("book");
const status = document.getElementById("status")!;
const notice = document.getElementById("notice")!;
const stage = document.getElementById("stage")!;

function tell(text: string, link?: { href: string; text: string }) {
  notice.textContent = text;
  if (link) {
    const a = document.createElement("a");
    a.href = link.href;
    a.textContent = link.text;
    notice.append(" ", a);
  }
  notice.hidden = false;
}

/** What the server said when it did not send pages. */
function refused(e: unknown): boolean {
  if (!(e instanceof BdfWorkerError)) return false;
  if (e.code === "not-allowed") tell("The pages after the sample are for the book's owners.", { href: "/", text: "Back to your books" });
  else if (e.code === "rate-limited") tell("Pages are asked for faster than anyone reads them: wait a minute and scroll again.");
  else return false;
  return true;
}

async function main() {
  if (!book) {
    status.textContent = "no ?book= given: pick a book from your books";
    return;
  }
  let title = book, pages = 0;
  const viewer = new MiniViewer(stage, {
    worker: new URL("lib/worker.js", location.href).href,
    onError: (e) => refused(e) || console.error(e),
    onChange: ({ page }) => (status.textContent = `${title} · page ${page + 1} of ${pages}`),
  });
  try {
    const manifest = await viewer.open({ kind: "segments", url: new URL(`/segments/${encodeURIComponent(book)}`, location.href).href });
    if (!manifest) return;
    title = dcValues(manifest.meta?.dc?.title)[0] ?? book;
    pages = manifest.views[0]?.pages?.length ?? 0;
    document.title = `${title} – secure-reader`;
    status.textContent = `${title} · ${pages} pages, ten at a time`;
  } catch (e) {
    // the first segment: 401 when the session is gone
    if (e instanceof BdfWorkerError && e.code === "not-allowed") {
      location.href = "/login";
      return;
    }
    if (!refused(e)) throw e;
  }
}

main().catch((e) => (status.textContent = `error: ${(e as Error).message ?? e}`));
