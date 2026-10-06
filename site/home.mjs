// The top page of the site, in English (index.html) and Japanese
// (index.ja.html): a PDF shown from bdf as a book whose pages turn
// (home.ts), pictures of other kinds of documents in the viewer, and the ways
// on to the demos, the sample architectures and the documentation.
import { esc, head, header } from "./docs.mjs";
import { REPO } from "./nav.mjs";

/** The pictures (docs/images/, made by site/shots.mjs), the sample the viewer opens, and the page about the format. */
const GALLERY = [
  {
    shot: "pdf", sample: "demo.pdf", doc: "formats/pdf",
    title: { en: "PDF", ja: "PDF" },
    text: {
      en: "Pages with their fonts, text to select and the structure of tagged PDF. Illustrator files too.",
      ja: "フォントを埋め込んだページ、選択できるテキスト、タグ付き PDF の構造。Illustrator のファイルも。",
    },
  },
  {
    shot: "epub", sample: "vertical.epub", doc: "formats/ebook",
    title: { en: "EPUB", ja: "EPUB" },
    text: {
      en: "Books too, with pages that turn: reflowable ones in vertical Japanese, and fixed layouts.",
      ja: "本も、ページをめくって。リフロー型は和文の縦書きにも、固定レイアウトにも対応。",
    },
  },
  {
    shot: "pptx", sample: "features.pptx", doc: "formats/presentation",
    title: { en: "PowerPoint", ja: "PowerPoint" },
    text: { en: "Slides with their masters, charts, SmartArt and formulas.", ja: "スライドマスター、グラフ、SmartArt、数式を含むスライド。" },
  },
  {
    shot: "xlsx", sample: "features.xlsx", doc: "formats/spreadsheet",
    title: { en: "Excel", ja: "Excel" },
    text: {
      en: "Sheets that scroll without page breaks, with frozen panes, headers and tabs. CSV, TSV and Parquet files are shown the same way.",
      ja: "ページに分断されずスクロールできるシート。ウィンドウ枠の固定、見出し、シートのタブつき。CSV・TSV・Parquet も同じ見た目で表示します。",
    },
  },
  {
    shot: "docx", sample: "basic.docx", doc: "formats/document",
    title: { en: "Word", ja: "Word" },
    text: { en: "Pages, or one long column to scroll. Vertical Japanese text too.", ja: "紙面のページでも、一続きのスクロールでも。縦書きにも対応。" },
  },
  {
    shot: "drawio", sample: "aws.drawio", doc: "formats/diagram",
    title: { en: "draw.io", ja: "draw.io" },
    text: { en: "Pages as tabs, with today's AWS icons. Visio drawings too.", ja: "ページはタブで切り替え。AWS のアイコンも描画。Visio も。" },
  },
  {
    shot: "dxf", sample: "layout.dxf", doc: "formats/cad",
    title: { en: "DXF", ja: "DXF" },
    text: { en: "Model space and layouts. Jw_cad, SXF, CGM and HP-GL/2 too.", ja: "モデル空間とレイアウト。Jw_cad、SXF、CGM、HP-GL/2 も。" },
  },
  {
    shot: "kicad", sample: "demo.zip", doc: "formats/electronics",
    title: { en: "KiCad", ja: "KiCad" },
    text: { en: "Schematic sheets and boards. Gerber and Excellon files too.", ja: "回路図のシートと基板。Gerber・Excellon のファイルも。" },
  },
  {
    shot: "music", sample: "minuet.musicxml", doc: "formats/music",
    title: { en: "Scores", ja: "楽譜" },
    text: { en: "MML, MIDI and MusicXML, engraved as scores that the viewer plays.", ja: "MML・MIDI・MusicXML を五線譜に組み、ビューアで演奏。" },
  },
];

const BENCHMARK_CASES = {
  "basic.docx": { en: "Basic DOCX", ja: "小さい DOCX" },
  "basic.pptx": { en: "Basic PPTX", ja: "小さい PPTX" },
  "basic.xlsx": { en: "Basic XLSX", ja: "小さい XLSX" },
  "many-sheets.xlsx": { en: "Synthetic XLSX · 4 sheets × 3,000 rows", ja: "合成 XLSX · 4 シート × 3,000 行" },
  "many-slides.pptx": { en: "Synthetic PPTX · 100 slides", ja: "合成 PPTX · 100 スライド" },
};

const T = {
  en: {
    title: "BDF – document previews for browsers",
    description: "BDF is a browser document-preview suite with Go/WASM converters, a browser renderer of about 23 KB gzipped, and a compact intermediate format.",
    h1: "Preview documents in the browser, without an office suite",
    hero: (kb) => `BDF is a browser document-preview suite: it combines converters for many formats (Go/WASM), a browser renderer of about ${kb} KB gzipped, and a compact intermediate format.`,
    heroDetail: "PDF is a portable, PostScript-based format for printing. By contrast, BDF is a lightweight intermediate format built from Canvas 2D drawing commands for fast browser rendering. Generate BDF ahead of time on a server, or convert files in real time in the browser with WebAssembly.",
    try: "Open the viewer", docs: "Read the docs",
    demoTitle: "A PDF, drawn from BDF",
    demoLead: "Converted ahead of time, as a server would. Turn its pages, select its text, search it, zoom.",
    demoName: "demo.pdf → demo.bdf",
    demoNote: (pdf, bdf) => `${pdf} as a PDF, ${bdf} as BDF with its fonts. The pages are fetched by range requests as they are about to be shown.`,
    benchmarkTitle: "Less time, memory and energy for document previews",
    benchmarkLead: "Compare the same Office files converted by LibreOffice + Poppler and by BDF's Go converter.",
    benchmarkEnergyTitle: "Estimated SoC energy per document",
    benchmarkTimeTitle: "Conversion + thumbnail time (median)",
    benchmarkMemoryTitle: "Maximum RSS (median)",
    benchmarkHighlightsLabel: "Benchmark highlights",
    benchmarkEnergyLess: "less estimated SoC energy per document",
    benchmarkEnergyMore: "more estimated SoC energy per document",
    benchmarkEnergySame: "the same estimated SoC energy per document",
    benchmarkEnergyInput: (inputs) => `Per-document average across ${inputs}`,
    benchmarkEnergyRange: (rounds, officeMin, officeMax, bdfMin, bdfMax) => `Range across ${rounds} rounds: LibreOffice + Poppler ${officeMin}–${officeMax} J/document; BDF ${bdfMin}–${bdfMax} J/document.`,
    benchmarkFaster: (name) => `faster on ${name}`,
    benchmarkSlower: (name) => `slower on ${name}`,
    benchmarkLowerMemory: (name) => `lower maximum RSS on ${name}`,
    benchmarkHigherMemory: (name) => `higher maximum RSS on ${name}`,
    benchmarkSameSpeed: "the same conversion time",
    benchmarkSameMemory: "the same maximum RSS",
    benchmarkOffice: "LibreOffice + Poppler",
    benchmarkBdf: "BDF (Go)",
    benchmarkBarNote: "Bar lengths are scaled within each input pair; values are the measured results.",
    benchmarkMethod: (runs, rounds, inputs) => `On an Apple Silicon Mac connected to AC power. Time and maximum RSS are medians of ${runs} fresh-process runs for conversion plus a 256 px first-page thumbnail; the large XLSX and PPTX inputs are synthetic. Energy is an idle-adjusted CPU, GPU and ANE SoC estimate from powermetrics across ${rounds} rounds of ${inputs}, divided by completed documents—not a wall-outlet measurement.`,
    search: "search…", prev: "previous match", next: "next match", zoomIn: "zoom in", zoomOut: "zoom out",
    layout: "layout", layouts: { spread: "Two pages", single: "One page", scroll: "Scroll" },
    prevPage: "previous page", nextPage: "next page", pages: "pages",
    galleryTitle: "Each kind of document, in its own shape",
    galleryLead: "A sheet stays a sheet that scrolls; a book turns its pages. All of them are drawn by the same renderer.",
    open: "Open in the viewer", about: "About the format",
    more: (n) => `${n} formats in all, among them Visio, Illustrator, Photoshop, TIFF, HTML, Markdown and font files:`, all: "all formats",
    demosTitle: "Try it with your files",
    demosLead: "Files are converted inside your browser. They are not uploaded.",
    demos: [
      ["viewer/", "Viewer", "Drop a file to view it: search, text selection, pages that turn, scores that play."],
      ["thumbnail/", "Thumbnails", "Thumbnails at four sizes at once, drawn without a browser engine, as on a server."],
      ["text/", "Search text", "The text and metadata of a document, page by page, for a search index."],
    ],
    archTitle: "Fit it into your system",
    archLead: "Three sample projects, each a small Go server with its pages.",
    arch: [
      ["examples/light-server", "Convert in the browser", "For documents made in large numbers but seldom opened, such as records kept for audits. The server only makes thumbnails; a file is converted in the browser of whoever opens it."],
      ["examples/preview-server", "Convert on the server", "The server makes the preview and the thumbnail once. Browsers only draw."],
      ["examples/search", "With a search engine", "The text goes to a search engine page by page. A hit opens the document at its page."],
    ],
    read: "Read more",
    whyTitle: "Why BDF",
    why: [
      ["No office suite to run", "No headless LibreOffice, no cgo, no external program: the converters are Go packages in one binary."],
      ["In the shape of the content", "Fixed pages for slides and drawings, endless sheets for workbooks, pages or one column for text."],
      ["Made for browsers", "Instructions map onto Canvas 2D. Fonts, images and compression are left to the browser."],
      ["Many formats, one renderer", "Search, text selection and the text for screen readers work the same for every format."],
      ["Only the converters you need", "Each format is a Go package of its own: link in the ones your service needs, or the browser modules for those kinds of files."],
      ["Find in page and screen readers", "The text lies over the picture as real text, so the browser's own find (Ctrl+F) finds it, and screen readers read its headings, lists and tables."],
    ],
    whyMore: "More about why",
    license: "MIT License",
    fonts: "Fonts and licenses",
  },
  ja: {
    title: "BDF – ブラウザのための文書プレビュー",
    description: "BDF はブラウザ向けの文書プレビュースイートです。Go/WASM の多形式コンバーター、gzip 圧縮後約 23 KB のブラウザ用レンダラー、中間形式で構成されます。",
    h1: "オフィススイートなしで、文書をブラウザでプレビュー",
    hero: (kb) => `BDF は、ブラウザ向けの文書プレビュースイートです。多様な形式に対応するコンバーター（Go/WASM）、gzip 圧縮後約 ${kb} KB のブラウザ用レンダラー、中間形式で構成されます。`,
    heroDetail: "PDF が PostScript をベースにした印刷向けのポータブル文書形式であるのに対し、BDF は Canvas 2D の描画命令を基盤とした、ブラウザ表示向けの軽量な中間形式です。サーバーでの事前生成と、WASM を使ったブラウザ内でのリアルタイム変換に対応します。",
    try: "ビューアを開く", docs: "ドキュメントを読む",
    demoTitle: "BDF から描いた PDF",
    demoLead: "サーバーで変換するのと同じように、あらかじめ変換してあります。ページめくり、テキストの選択、検索、拡大ができます。",
    demoName: "demo.pdf → demo.bdf",
    demoNote: (pdf, bdf) => `PDF で ${pdf}、フォントを含む BDF で ${bdf}。ページは表示の直前に Range リクエストで取得します。`,
    benchmarkTitle: "文書変換が速く、メモリと電力も少なく",
    benchmarkLead: "同じ Office ファイルを LibreOffice + Poppler と BDF の Go 変換器で処理して比較しました。",
    benchmarkEnergyTitle: "1 文書あたりの推定 SoC エネルギー",
    benchmarkTimeTitle: "変換＋サムネイル作成の経過時間（中央値）",
    benchmarkMemoryTitle: "最大 RSS（中央値）",
    benchmarkHighlightsLabel: "主な結果",
    benchmarkEnergyLess: "文書あたりの推定 SoC エネルギーを削減",
    benchmarkEnergyMore: "文書あたりの推定 SoC エネルギーが増加",
    benchmarkEnergySame: "文書あたりの推定 SoC エネルギーは同程度",
    benchmarkEnergyInput: (inputs) => `${inputs} の文書あたり平均`,
    benchmarkEnergyRange: (rounds, officeMin, officeMax, bdfMin, bdfMax) => `${rounds} ラウンドの範囲（J/文書）: LibreOffice + Poppler ${officeMin}–${officeMax}、BDF ${bdfMin}–${bdfMax}。`,
    benchmarkFaster: (name) => `${name} の変換が高速`,
    benchmarkSlower: (name) => `${name} の変換が低速`,
    benchmarkLowerMemory: (name) => `${name} の最大 RSS が少ない`,
    benchmarkHigherMemory: (name) => `${name} の最大 RSS が多い`,
    benchmarkSameSpeed: "変換時間は同程度",
    benchmarkSameMemory: "最大 RSS は同程度",
    benchmarkOffice: "LibreOffice + Poppler",
    benchmarkBdf: "BDF（Go）",
    benchmarkBarNote: "横棒はファイルごとに比較し、長い方を 100% に正規化。数値は測定値です。",
    benchmarkMethod: (runs, rounds, inputs) => `Apple Silicon の Mac を AC 電源で使用。時間と最大 RSS は、変換と先頭ページの 256 px サムネイル作成までを新規プロセスで計測した ${runs} 回の中央値。大きな XLSX・PPTX は合成データです。エネルギーは ${inputs} の変換を ${rounds} ラウンド繰り返し、powermetrics の CPU・GPU・ANE の SoC 推定値からアイドル分を差し引いて文書あたりで算出。コンセントでの実測値ではありません。`,
    search: "検索…", prev: "前の一致", next: "次の一致", zoomIn: "拡大", zoomOut: "縮小",
    layout: "表示", layouts: { spread: "見開き", single: "1 ページずつ", scroll: "スクロール" },
    prevPage: "前のページ", nextPage: "次のページ", pages: "ページ",
    galleryTitle: "文書の種類ごとに、合った形で",
    galleryLead: "シートはスクロールできるシートのまま、本はページをめくって。どれも同じレンダラが描きます。",
    open: "ビューアで開く", about: "形式の説明",
    more: (n) => `対応する形式は ${n} 種類。Visio、Illustrator、Photoshop、TIFF、HTML、Markdown、フォントファイルなども:`, all: "形式の一覧",
    demosTitle: "手元のファイルで試す",
    demosLead: "ファイルはブラウザの中で変換します。アップロードはされません。",
    demos: [
      ["viewer/", "ビューア", "ファイルをドロップして表示。検索、テキスト選択、ページめくり、楽譜の演奏。"],
      ["thumbnail/", "サムネイル", "4 つの大きさのサムネイルを一度に。サーバーと同じく、ブラウザのエンジンなしで描きます。"],
      ["text/", "検索用テキスト", "検索エンジンに渡す、ページごとのテキストとメタデータ。"],
    ],
    archTitle: "システムに組み込む",
    archLead: "構成ごとのサンプルプロジェクトが 3 つ。どれも小さな Go のサーバーとページです。",
    arch: [
      ["examples/light-server", "ブラウザで変換する", "監査のために残す記録のように、大量に作られるが開かれるのはまれな文書に。サーバーはサムネイルだけを作り、変換は開いた人のブラウザで。"],
      ["examples/preview-server", "サーバーで変換する", "プレビューとサムネイルをサーバーが一度だけ作ります。ブラウザは描くだけ。"],
      ["examples/search", "検索エンジンとつなぐ", "テキストをページごとに検索エンジンへ。ヒットからそのページを開きます。"],
    ],
    read: "解説を読む",
    whyTitle: "なぜ BDF か",
    why: [
      ["オフィススイートを動かさない", "ヘッドレスの LibreOffice も cgo も外部プログラムも不要。変換器は 1 つのバイナリに入る Go のパッケージです。"],
      ["内容に合った形で見せる", "スライドや図面には固定ページ、ブックには無限のシート、文章にはページか一続きの列。"],
      ["ブラウザ表示に特化", "命令は Canvas 2D に対応。フォント、画像、展開はブラウザに任せます。"],
      ["多くの形式を 1 つのレンダラで", "検索、テキスト選択、読み上げ用のテキストは、どの形式でも同じように使えます。"],
      ["必要な変換器だけを選べる", "変換器は形式ごとに独立した Go のパッケージ。サービスに要るものだけを組み込めます。ブラウザ用のモジュールも種類ごとに分かれています。"],
      ["ブラウザの検索も、読み上げも", "文字は画像の上に本物のテキストとして重ねてあるので、ブラウザ自身のページ内検索（Ctrl+F）で見つかり、スクリーンリーダーは見出し・リスト・表を読み上げます。"],
    ],
    whyMore: "もっと詳しく",
    license: "MIT ライセンス",
    fonts: "フォントとライセンス",
  },
};

const size = (n) => (n < 1024 * 1024 ? `${Math.round(n / 1024)} KB` : `${(n / 1024 / 1024).toFixed(1)} MB`);

function benchmarkSection(lang, data) {
  const t = T[lang];
  const process = data?.process;
  const power = data?.power;
  if (!process?.results || !power?.median_estimated_soc_j_per_document || !Array.isArray(power.inputs) || !power.inputs.length) {
    throw new Error("The homepage requires process and power results in docs/benchmarks/latest.json");
  }
  if (!Number.isInteger(process.rounds) || process.rounds < 1 || !Number.isInteger(power.rounds) || power.rounds < 1) {
    throw new Error("Invalid benchmark round count in docs/benchmarks/latest.json");
  }

  const positive = (value, label) => {
    if (!Number.isFinite(value) || value <= 0) throw new Error(`Invalid benchmark value for ${label}`);
    return value;
  };
  const number = (value, digits) => new Intl.NumberFormat(lang, {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  }).format(value);
  const labelFor = (name) => BENCHMARK_CASES[name]?.[lang] ?? name;
  const cases = Object.entries(process.results).map(([name, result]) => ({
    label: labelFor(name),
    officeTime: positive(result?.libreoffice?.elapsed_seconds, `${name} LibreOffice elapsed time`),
    bdfTime: positive(result?.bdf?.elapsed_seconds, `${name} bdf elapsed time`),
    officeMemory: positive(result?.libreoffice?.max_rss_bytes, `${name} LibreOffice maximum RSS`),
    bdfMemory: positive(result?.bdf?.max_rss_bytes, `${name} bdf maximum RSS`),
  }));
  if (!cases.length) throw new Error("No process benchmark results in docs/benchmarks/latest.json");

  const officeEnergy = positive(power.median_estimated_soc_j_per_document.libreoffice, "LibreOffice estimated SoC energy");
  const bdfEnergy = positive(power.median_estimated_soc_j_per_document.bdf, "bdf estimated SoC energy");
  const energyRanges = power.range_estimated_soc_j_per_document;
  if (!energyRanges) throw new Error("Missing energy ranges in docs/benchmarks/latest.json");
  const range = (values, label) => {
    if (!Array.isArray(values) || values.length !== 2) throw new Error(`Invalid benchmark range for ${label}`);
    const min = positive(values[0], `${label} minimum`);
    const max = positive(values[1], `${label} maximum`);
    if (min > max) throw new Error(`Invalid benchmark range for ${label}`);
    return [min, max];
  };
  const [officeEnergyMin, officeEnergyMax] = range(energyRanges.libreoffice, "LibreOffice estimated SoC energy");
  const [bdfEnergyMin, bdfEnergyMax] = range(energyRanges.bdf, "bdf estimated SoC energy");
  const energyChange = Math.round(Math.abs(1 - bdfEnergy / officeEnergy) * 100);
  const energyHighlight = {
    value: `${energyChange}%`,
    label: energyChange === 0
      ? t.benchmarkEnergySame
      : bdfEnergy < officeEnergy ? t.benchmarkEnergyLess : t.benchmarkEnergyMore,
  };

  const bestBy = (a, b) => a.ratio - b.ratio;
  const speedUps = cases.map((entry) => ({ label: entry.label, ratio: entry.officeTime / entry.bdfTime }));
  const fastest = speedUps.reduce((best, current) => bestBy(current, best) > 0 ? current : best);
  const slowest = speedUps.reduce((best, current) => current.ratio < best.ratio ? current : best);
  const speedRatio = fastest.ratio >= 1 ? fastest.ratio : 1 / slowest.ratio;
  const speedHighlight = {
    value: `${number(speedRatio, 0)}×`,
    label: fastest.ratio > 1 ? t.benchmarkFaster(fastest.label)
      : slowest.ratio < 1 ? t.benchmarkSlower(slowest.label) : t.benchmarkSameSpeed,
  };

  const memoryChanges = cases.map((entry) => ({ label: entry.label, ratio: entry.officeMemory / entry.bdfMemory }));
  const leastMemory = memoryChanges.reduce((best, current) => bestBy(current, best) > 0 ? current : best);
  const mostMemory = memoryChanges.reduce((best, current) => current.ratio < best.ratio ? current : best);
  const memoryRatio = leastMemory.ratio >= 1 ? leastMemory.ratio : 1 / mostMemory.ratio;
  const memoryHighlight = {
    value: `${number(memoryRatio, 0)}×`,
    label: leastMemory.ratio > 1 ? t.benchmarkLowerMemory(leastMemory.label)
      : mostMemory.ratio < 1 ? t.benchmarkHigherMemory(mostMemory.label) : t.benchmarkSameMemory,
  };

  const row = (label, office, bdf, format) => {
    const max = Math.max(office, bdf);
    const officeWidth = Math.max(2, office / max * 100);
    const bdfWidth = Math.max(2, bdf / max * 100);
    return `<li class="benchmark-case">
<h4>${esc(label)}</h4>
<div class="benchmark-measure office">
<span class="benchmark-series">${esc(t.benchmarkOffice)}</span>
<span class="benchmark-track" aria-hidden="true"><span style="width:${officeWidth.toFixed(1)}%"></span></span>
<strong class="benchmark-value">${esc(format(office))}</strong>
</div>
<div class="benchmark-measure bdf">
<span class="benchmark-series">${esc(t.benchmarkBdf)}</span>
<span class="benchmark-track" aria-hidden="true"><span style="width:${bdfWidth.toFixed(1)}%"></span></span>
<strong class="benchmark-value">${esc(format(bdf))}</strong>
</div>
</li>`;
  };
  const figure = (id, title, rows, additionalNote = "") => `<figure class="benchmark-card">
<h3 id="${id}">${esc(title)}</h3>
<ul class="benchmark-cases">${rows.join("\n")}</ul>
<figcaption>${esc(t.benchmarkBarNote)}${additionalNote ? `<br>${esc(additionalNote)}` : ""}</figcaption>
</figure>`;
  const inputNames = power.inputs.map(labelFor);
  const energyInput = t.benchmarkEnergyInput(inputNames.join(lang === "ja" ? "・" : ", "));
  const energyRange = t.benchmarkEnergyRange(
    power.rounds,
    number(officeEnergyMin, 2),
    number(officeEnergyMax, 2),
    number(bdfEnergyMin, 2),
    number(bdfEnergyMax, 2),
  );
  const highlights = [energyHighlight, speedHighlight, memoryHighlight]
    .map(({ value, label }) => `<li class="benchmark-highlight"><strong>${esc(value)}</strong><span>${esc(label)}</span></li>`)
    .join("\n");

  return `<section class="wrap benchmarks" aria-labelledby="benchmark-title">
<h2 id="benchmark-title">${esc(t.benchmarkTitle)}</h2>
<p class="lead">${esc(t.benchmarkLead)}</p>
<ul class="benchmark-highlights" aria-label="${esc(t.benchmarkHighlightsLabel)}">${highlights}</ul>
<div class="benchmark-grid">
${figure("benchmark-time-title", t.benchmarkTimeTitle, cases.map((entry) =>
    row(entry.label, entry.officeTime, entry.bdfTime, (value) => `${number(value, 2)} s`)))}
${figure("benchmark-memory-title", t.benchmarkMemoryTitle, cases.map((entry) =>
    row(entry.label, entry.officeMemory / (1024 * 1024), entry.bdfMemory / (1024 * 1024), (value) => `${number(value, 1)} MiB`)))}
${figure("benchmark-energy-title", t.benchmarkEnergyTitle, [
    row(energyInput, officeEnergy, bdfEnergy, (value) => `${number(value, 2)} J`),
  ], energyRange)}
</div>
<p class="benchmark-method">${esc(t.benchmarkMethod(
    process.rounds,
    power.rounds,
    inputNames.join(lang === "ja" ? "・" : ", "),
  ))}</p>
</section>`;
}

/**
 * The top page in a language. demo gives the sizes of the PDF shown and of
 * the bdf made of it, formats how many input formats the converters take,
 * worker the size of the rendering worker, minified and gzipped, and benchmark
 * the latest comparison results.
 */
export function home(lang, { demo, formats, worker, benchmark }) {
  const t = T[lang];
  const out = lang === "ja" ? "index.ja.html" : "index.html";
  const other = lang === "ja" ? "index.html" : "index.ja.html";
  const doc = (id) => `docs/${id}${lang === "ja" ? ".ja" : ""}.html`;
  const gallery = GALLERY.map((g) => `<li class="card">
<a class="shot" href="viewer/?file=samples/${g.sample}" tabindex="-1" aria-hidden="true"><img src="docs/images/${g.shot}.webp" alt="" width="1600" height="1000" loading="lazy"></a>
<div class="body">
<h3>${esc(g.title[lang])}</h3>
<p>${esc(g.text[lang])}</p>
<p class="links"><a href="viewer/?file=samples/${g.sample}" aria-label="${esc(`${t.open}: ${g.title[lang]}`)}">${t.open}</a> <a href="${doc(g.doc)}" aria-label="${esc(`${t.about}: ${g.title[lang]}`)}">${t.about}</a></p>
</div>
</li>`).join("\n");
  const plain = (items, href, label) => items.map(([to, title, text]) => `<li class="card">
<div class="body">
<h3><a href="${href(to)}">${esc(title)}</a></h3>
<p>${esc(text)}</p>${label ? `\n<p class="links"><a href="${href(to)}" aria-label="${esc(`${label}: ${title}`)}">${label}</a></p>` : ""}
</div>
</li>`).join("\n");
  return `<!doctype html>
<html lang="${lang}">
<head>
${head({ out, other, lang })}
<meta name="description" content="${esc(t.description)}">
<title>${esc(t.title)}</title>
<script>
// the viewer was the first page of the site: its addresses (?src=…, ?layout=…) go on to it
if (/[?&](src|layout|range)\\b/.test(location.search)) location.replace("viewer/" + location.search + location.hash);
</script>
</head>
<body class="home">
${header({ lang, out, other })}
<main id="content" tabindex="-1">
<section class="hero wrap">
<h1>${esc(t.h1)}</h1>
<p>${esc(t.hero(Math.round(worker / 1024)))}</p>
<p class="format-note">${esc(t.heroDetail)}</p>
<p class="actions"><a class="button primary" href="viewer/">${t.try}</a> <a class="button" href="${doc("index")}">${t.docs}</a> <a class="button" href="${REPO}">GitHub</a></p>
</section>
<section class="wrap" aria-labelledby="demo-title">
<h2 id="demo-title">${esc(t.demoTitle)}</h2>
<p class="lead">${esc(t.demoLead)}</p>
<div class="demo">
<div class="tools">
<span class="name">${esc(t.demoName)}</span>
<span role="search"><input id="q" type="search" placeholder="${t.search}" aria-label="${t.search}">
<button type="button" id="prev" title="${t.prev}" aria-label="${t.prev}">▲</button>
<button type="button" id="next" title="${t.next}" aria-label="${t.next}">▼</button></span>
<span id="hits" role="status"></span>
<select id="layout" aria-label="${t.layout}" title="${t.layout}">${Object.entries(t.layouts).map(([v, label]) => `<option value="${v}">${label}</option>`).join("")}</select>
<span class="pager" role="group" aria-label="${t.pages}"><button type="button" id="prevPage" title="${t.prevPage}" aria-label="${t.prevPage}">‹</button>
<span id="pageNum" aria-live="polite"></span>
<button type="button" id="nextPage" title="${t.nextPage}" aria-label="${t.nextPage}">›</button></span>
<span><button type="button" id="zoomOut" title="${t.zoomOut}" aria-label="${t.zoomOut}">−</button>
<button type="button" id="zoomIn" title="${t.zoomIn}" aria-label="${t.zoomIn}">+</button></span>
</div>
<div class="view" id="demo" data-src="demo.bdf" role="region" aria-label="demo.pdf" tabindex="0"></div>
</div>
<p class="note">${esc(t.demoNote(size(demo.pdf), size(demo.bdf)))}</p>
</section>
${benchmarkSection(lang, benchmark)}
<section class="wrap" aria-labelledby="gallery-title">
<h2 id="gallery-title">${esc(t.galleryTitle)}</h2>
<p class="lead">${esc(t.galleryLead)}</p>
<ul class="cards">
${gallery}
</ul>
<p class="more">${esc(t.more(formats))} <a href="${doc("formats/index")}">${t.all}</a></p>
</section>
<section class="wrap" aria-labelledby="demos-title">
<h2 id="demos-title">${esc(t.demosTitle)}</h2>
<p class="lead">${esc(t.demosLead)}</p>
<ul class="cards plain">
${plain(t.demos, (to) => to)}
</ul>
</section>
<section class="wrap" aria-labelledby="arch-title">
<h2 id="arch-title">${esc(t.archTitle)}</h2>
<p class="lead">${esc(t.archLead)}</p>
<ul class="cards plain">
${plain(t.arch, doc, t.read)}
</ul>
</section>
<section class="wrap" aria-labelledby="why-title">
<h2 id="why-title">${esc(t.whyTitle)}</h2>
<ul class="points">
${t.why.map(([title, text]) => `<li><h3>${esc(title)}</h3><p>${esc(text)}</p></li>`).join("\n")}
</ul>
<p class="more"><a href="${doc("why")}">${t.whyMore}</a></p>
</section>
</main>
<footer class="wrap">BDF · <a href="${REPO}/blob/main/LICENSE">${t.license}</a> · <a href="${doc("licenses")}">${t.fonts}</a> · <a href="${REPO}">GitHub</a></footer>
<script type="module" src="home.js"></script>
</body>
</html>
`;
}
