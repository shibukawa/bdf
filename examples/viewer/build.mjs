// Builds the demo viewer with esbuild: the page, and the workers it shares
// with the other pages of the site (examples/common/build.mjs).
import { join } from "node:path";
import { bundle, buildWorkers, copyPage, root } from "../common/build.mjs";

export { root };

/**
 * Build the viewer into out. defaultSrc is the document shown when the URL
 * has no ?src= ("" shows the page for opening a file). siteRoot is where the
 * page finds lib/ (the workers and the converter modules), fonts/ and
 * samples/, relative to it: the workers are built there. site says that the
 * page is one of the demo site's, which links to the others.
 */
export async function buildViewer(out, { defaultSrc, siteRoot = "./", site = false }) {
  await bundle(out, [["examples/viewer/main.ts", "main"]], { define: { DEFAULT_SRC: JSON.stringify(defaultSrc), SITE_ROOT: JSON.stringify(siteRoot) } });
  await buildWorkers(join(out, siteRoot, "lib"));
  await copyPage("examples/viewer/index.html", join(out, "index.html"), { site });
}
