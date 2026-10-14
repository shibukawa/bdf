// Browser regression test of the built form page: npm run site && npm run test:form-page.
// Uses CHROMIUM_PATH/BDF_CHROME, Google Chrome on macOS, or Playwright's Chromium
// (install with `npx playwright-core install chromium`). Kept separate from test:site,
// whose Node/Wasm smoke tests also run on hosts without a browser.
// Optional arguments: [site directory] [screenshot directory].
import assert from "node:assert/strict";
import { existsSync } from "node:fs";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { chromium } from "playwright-core";
import { BdfDocument, BufferSource, extractText } from "../packages/core/dist/index.js";
import { serve } from "./serve.mjs";

const site = resolve(process.argv[2] ?? "site/dist");
const screenshots = process.argv[3] ? resolve(process.argv[3]) : undefined;
const executablePath = process.env.CHROMIUM_PATH ?? process.env.BDF_CHROME ?? [
  "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
  "/opt/pw-browsers/chromium_headless_shell-1194/chrome-linux/headless_shell",
].find((path) => existsSync(path));
let server, browser;

/** The filled document offered for download, including edits made while Preview is hidden. */
async function filledText(page) {
  const bytes = await page.evaluate(async () => {
    const response = await fetch(document.querySelector("#saveBdf").href);
    return Array.from(new Uint8Array(await response.arrayBuffer()));
  });
  const doc = await BdfDocument.open(new BufferSource(new Uint8Array(bytes)));
  let text = "";
  for (const view of doc.manifest.views) for (const page of view.pages ?? []) for (const layer of page.layers) {
    for (const run of extractText(await doc.ensure(layer.obj), (hash) => doc.objectSync(hash))) text += run.text;
  }
  return text;
}

/** Only the selected top-level panel is visible, and it is the keyboard stop. */
async function selected(page, name) {
  assert.deepEqual(await page.evaluate(() => ["Template", "Data", "Preview"].map((name) => {
    const tab = document.querySelector(`#tab${name}`);
    return [name, tab.getAttribute("aria-selected"), tab.tabIndex, document.querySelector(`#${name.toLowerCase()}Panel`).hidden];
  })), ["Template", "Data", "Preview"].map((value) => [value, String(value === name), value === name ? 0 : -1, value !== name]));
}

/** Every Office shape and the pdfme invoice offer the selected style; shared samples stay available. */
async function sampleStyle(page, style) {
  assert.equal(await page.locator("#sampleStyle").inputValue(), style);
  const suffix = style === "modern" ? "-modern" : "";
  const expected = [
    ...["fixed", "rows", "pages"].flatMap((shape) => ["xlsx", "docx", "pptx"].map((format) => `form-${shape}${suffix}.${format}`)),
    "application.pdf", `pdfme-invoice${suffix}.json`, "pdfme-certificate.json",
  ];
  assert.deepEqual(await page.locator("#presets button").evaluateAll((buttons) => buttons.filter((button) => button.checkVisibility()).map((button) => button.title)), expected);
}

async function openSample(page, name) {
  const previous = await page.locator("#saveBdf").getAttribute("href");
  const canvas = page.locator("#preview canvas").first();
  const drawing = await canvas.count() ? await canvas.evaluate((canvas) => canvas.toDataURL()) : undefined;
  await page.locator(`#presets button[title="${name}"]`).click();
  await page.waitForFunction(([name, previous]) => {
    return document.querySelector("#status").hasAttribute("data-error")
      || document.querySelector("#status").textContent.startsWith(`${name} filled`)
      && document.querySelector("#saveBdf").getAttribute("href") !== previous
      && !document.querySelector("#savePdf").disabled;
  }, [name, previous], { timeout: 120_000 });
  assert.equal(await page.locator("#status").getAttribute("data-error"), null, await page.locator("#status").textContent());
  await selected(page, "Data");
  return drawing;
}

async function previewPainted(page, previous) {
  await page.locator("#tabPreview").click();
  await page.waitForFunction((previous) => {
    const canvas = document.querySelector("#preview canvas");
    if (!canvas || (previous && canvas.toDataURL() === previous)) return false;
    // A new canvas already has a default width; require pixels from the renderer.
    return canvas.getContext("2d").getImageData(0, 0, canvas.width, canvas.height).data.some((value, i) => i % 4 === 3 && value > 0);
  }, previous);
  return page.locator("#preview canvas").first().evaluate((canvas) => canvas.toDataURL());
}

async function previewFits(page) {
  await page.waitForFunction(() => {
    const preview = document.querySelector("#preview"), canvas = preview.querySelector("canvas");
    const style = getComputedStyle(preview);
    const available = preview.clientWidth - parseFloat(style.paddingLeft) - parseFloat(style.paddingRight);
    const width = canvas.getBoundingClientRect().width;
    return canvas.width > 0 && width <= available + 1 && Math.abs(canvas.width - width * (devicePixelRatio || 1)) <= 1;
  });
}

async function captureStyle(page, style) {
  if (screenshots) {
    await page.screenshot({ path: join(screenshots, `form-${style}-desktop.png`), fullPage: true });
    const png = await page.locator("#preview canvas").first().evaluate((canvas) => canvas.toDataURL());
    await writeFile(join(screenshots, `form-${style}-page.png`), Buffer.from(png.split(",")[1], "base64"));
  }
  await page.emulateMedia({ colorScheme: "dark" });
  assert.equal(await page.locator("#preview .page").first().evaluate((paper) => getComputedStyle(paper).backgroundColor), "rgb(255, 255, 255)");
  if (screenshots) await page.screenshot({ path: join(screenshots, `form-${style}-dark.png`), fullPage: true });
  await page.emulateMedia({ colorScheme: "light" });
  await page.setViewportSize({ width: 390, height: 844 });
  await previewFits(page);
  if (screenshots) await page.screenshot({ path: join(screenshots, `form-${style}-mobile.png`), fullPage: true });
  await page.setViewportSize({ width: 1280, height: 900 });
  await previewFits(page);
}

try {
  let port;
  ({ server, port } = await serve(site));
  browser = await chromium.launch(executablePath ? { executablePath } : {});
  const page = await browser.newPage({ viewport: { width: 1280, height: 900 } });
  page.setDefaultTimeout(30_000);
  const fontManifest = JSON.parse(await readFile(join(site, "form/fonts/index.json"), "utf8"));
  const fontURLs = new Set(fontManifest.map((font) => font.url));
  assert.ok(fontURLs.size > 0, "the font manifest supplies remote URLs");
  for (const url of fontURLs) assert.equal(new URL(url).origin, "https://fonts.gstatic.com");
  const fontRequests = [];
  let failFontRequest;
  const fontFailure = new Promise((_, reject) => { failFontRequest = reject; });
  // Dedicated-worker requests are reported on the context, rather than on Page.
  page.context().on("request", (request) => {
    if (fontURLs.has(request.url())) fontRequests.push(request.url());
  });
  page.context().on("requestfailed", (request) => {
    if (fontURLs.has(request.url())) failFontRequest(new Error(`Google Fonts request failed: ${request.url()}: ${request.failure()?.errorText}`));
  });
  page.context().on("response", (response) => {
    if (fontURLs.has(response.url()) && response.status() >= 400) failFontRequest(new Error(`Google Fonts returned HTTP ${response.status()}: ${response.url()}`));
  });
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  if (screenshots) await mkdir(screenshots, { recursive: true });
  await page.goto(`http://127.0.0.1:${port}/form/`);
  await page.locator('#presets button[title="form-fixed.xlsx"]').waitFor();
  await selected(page, "Template");
  assert.equal(await page.locator("#tabData").isDisabled(), true);
  assert.equal(await page.locator("#tabPreview").isDisabled(), true);
  await sampleStyle(page, "standard");
  await page.locator("#sampleStyle").selectOption("modern");
  await sampleStyle(page, "modern");
  await page.locator("#sampleStyle").selectOption("standard");
  await sampleStyle(page, "standard");

  await Promise.race([openSample(page, "form-fixed.xlsx"), fontFailure]);
  for (const url of fontURLs) assert.ok(fontRequests.includes(url), `the converter fetches the Google Font: ${url}`);
  const initialFontRequests = fontRequests.length;
  assert.equal(await page.locator("#tabPreview").isDisabled(), false);

  // The JSON tab still edits the data; hiding Preview must not prevent the new document from being filled.
  await page.locator("#tabJson").click();
  const defaultData = JSON.parse(await page.locator("#json").inputValue());
  assert.equal(defaultData.no, 1024);
  const data = structuredClone(defaultData);
  data.product = "デスクトレー";
  const previous = await page.locator("#saveBdf").getAttribute("href");
  await page.locator("#json").fill(JSON.stringify(data, null, 2));
  await page.waitForFunction((href) => document.querySelector("#saveBdf").getAttribute("href") !== href, previous, { timeout: 120_000 });
  assert.ok((await filledText(page)).includes("デスクトレー"), "editing hidden Preview updates the filled document");

  const firstDrawing = await previewPainted(page);
  await selected(page, "Preview");
  const surface = await page.locator("#preview").evaluate((preview) => {
    const page = preview.querySelector(".page"), canvas = page.querySelector("canvas");
    const style = getComputedStyle(preview);
    return {
      background: style.backgroundColor,
      paper: getComputedStyle(page).backgroundColor,
      shadow: getComputedStyle(page).boxShadow,
      canvasWidth: canvas.getBoundingClientRect().width,
      available: preview.clientWidth - parseFloat(style.paddingLeft) - parseFloat(style.paddingRight),
    };
  });
  assert.ok(!["rgba(0, 0, 0, 0)", "rgb(255, 255, 255)"].includes(surface.background), "Preview contrasts with the paper");
  assert.equal(surface.paper, "rgb(255, 255, 255)");
  assert.notEqual(surface.shadow, "none", "paper has a shadow");
  assert.ok(surface.canvasWidth > 100 && surface.canvasWidth <= surface.available + 1, "the newly visible page fits Preview");

  // Outer tabs and the nested data tabs retain their own state.
  await page.locator("#tabPreview").press("ArrowLeft");
  await selected(page, "Data");
  assert.equal(await page.locator("#tabJson").getAttribute("aria-selected"), "true");
  await page.locator("#tabForm").click();
  assert.equal(await page.getByLabel("product", { exact: true }).inputValue(), "デスクトレー");
  assert.equal(await page.getByLabel("no", { exact: true }).inputValue(), "1024");
  // Reuse the existing page frames: a second edit must replace pixels when Preview is shown again.
  const beforeSecondEdit = await page.locator("#saveBdf").getAttribute("href");
  await page.getByLabel("product", { exact: true }).fill("ファイルボックス");
  await page.waitForFunction((href) => document.querySelector("#saveBdf").getAttribute("href") !== href, beforeSecondEdit, { timeout: 120_000 });
  assert.ok((await filledText(page)).includes("ファイルボックス"), "a second edit updates the filled document");
  assert.equal(fontRequests.length, initialFontRequests, "editing data reuses fonts already fetched by the same worker");
  await previewPainted(page, firstDrawing);
  const [pdfDownload] = await Promise.all([
    page.waitForEvent("download", { timeout: 120_000 }),
    page.locator("#savePdf").click(),
  ]);
  const pdfPath = await pdfDownload.path();
  assert.ok(pdfPath, "the filled PDF downloads");
  const pdf = await readFile(pdfPath);
  assert.match(pdf.toString("latin1"), /\/FontFile[23]\b/, "the downloaded PDF embeds the remotely fetched fonts");
  await page.emulateMedia({ colorScheme: "dark" });
  assert.equal(await page.locator("#preview .page").first().evaluate((paper) => getComputedStyle(paper).backgroundColor), "rgb(255, 255, 255)");
  assert.notEqual(await page.locator("#preview").evaluate((preview) => getComputedStyle(preview).backgroundColor), surface.background);
  await page.emulateMedia({ colorScheme: "light" });
  await page.locator("#tabTemplate").click();
  await selected(page, "Template");
  await page.locator("#tabData").click();
  await selected(page, "Data");
  assert.equal(await page.getByLabel("product", { exact: true }).inputValue(), "ファイルボックス");
  assert.equal(await page.getByLabel("no", { exact: true }).inputValue(), "1024");

  await page.setViewportSize({ width: 390, height: 844 });
  await page.locator("#tabPreview").click();
  await previewFits(page);

  // Capture and compare both styles with the same untouched sample data.
  await page.setViewportSize({ width: 1280, height: 900 });
  await previewFits(page);
  await page.locator("#tabTemplate").click();
  const beforeStandard = await openSample(page, "form-fixed.xlsx");
  assert.deepEqual(JSON.parse(await page.locator("#json").inputValue()), defaultData);
  const standardDrawing = await previewPainted(page, beforeStandard);
  await captureStyle(page, "standard");
  await page.locator("#tabTemplate").click();
  await page.locator("#sampleStyle").selectOption("modern");
  await sampleStyle(page, "modern");
  const beforeModern = await openSample(page, "form-fixed-modern.xlsx");
  assert.deepEqual(JSON.parse(await page.locator("#json").inputValue()), defaultData, "both styles start from identical data");
  const modernText = await filledText(page);
  for (const value of ["1024", defaultData.customer.name, defaultData.product]) assert.ok(modernText.includes(value), `modern form contains ${value}`);
  const modernDrawing = await previewPainted(page, beforeModern);
  assert.notEqual(modernDrawing, standardDrawing, "the two styles render different page designs");
  await captureStyle(page, "modern");

  const linked = await browser.newPage({ viewport: { width: 1280, height: 900 } });
  linked.on("pageerror", (error) => errors.push(error.message));
  await linked.goto(`http://127.0.0.1:${port}/form/?template=form-fixed-modern.xlsx`);
  await linked.waitForFunction(() => document.querySelector("#status").textContent.startsWith("form-fixed-modern.xlsx filled"), null, { timeout: 120_000 });
  await selected(linked, "Data");
  assert.equal(await linked.locator("#sampleStyle").inputValue(), "modern", "a modern deep link selects Modern");
  assert.deepEqual(JSON.parse(await linked.locator("#json").inputValue()), defaultData);
  await linked.locator("#tabTemplate").click();
  await sampleStyle(linked, "modern");
  if (screenshots) await linked.screenshot({ path: join(screenshots, "form-modern-template-selector.png"), fullPage: true });
  await linked.close();
  assert.deepEqual(errors, []);
  console.log("form page browser regression test passed: remote font fetch/cache/PDF embedding, tabs, edits, preview paper, both sample styles, modern deep link");
} finally {
  await browser?.close();
  await new Promise((resolve) => server?.close(resolve) ?? resolve());
}
