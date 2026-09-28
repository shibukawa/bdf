// The top page of the site, in English (index.html) and Japanese
// (index.ja.html): a PDF shown from bdf (home.ts), pictures of other kinds of
// documents in the viewer, and the ways on to the demos, the sample
// architectures and the documentation.
import { esc, head, header } from "./docs.mjs";
import { REPO } from "./nav.mjs";

/** The pictures (docs/images/, made by site/shots.mjs), the sample the viewer opens, and the page about the format. */
const GALLERY = [
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
  {
    shot: "epub", sample: "vertical.epub", doc: "formats/ebook",
    title: { en: "EPUB", ja: "EPUB" },
    text: { en: "Books with pages that turn, vertical Japanese and fixed layouts.", ja: "ページをめくって読む本。縦書きと固定レイアウトにも対応。" },
  },
];

const T = {
  en: {
    title: "bdf – document previews for browsers",
    description: "bdf is a document format for previews in browsers: PDF, Office files, diagrams and CAD drawings, converted by one Go binary and drawn on a canvas.",
    h1: "Preview documents in the browser, without an office suite",
    hero: (kb) => `bdf is a document format made for previews. PDF, Office files, diagrams, CAD drawings and more are converted by one Go binary, or inside the browser as WebAssembly, and drawn on a canvas by a renderer of ${kb} KB.`,
    try: "Open the viewer", docs: "Read the docs",
    demoTitle: "A PDF, drawn from bdf",
    demoLead: "Converted ahead of time, as a server would. Select its text, search it, zoom.",
    demoName: "demo.pdf → demo.bdf",
    demoNote: (pdf, bdf) => `${pdf} as a PDF, ${bdf} as bdf with its fonts. The pages are fetched by range requests as they come into view.`,
    search: "search…", prev: "previous match", next: "next match", zoomIn: "zoom in", zoomOut: "zoom out",
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
      ["examples/light-server", "Convert in the browser", "A file server that only makes thumbnails. Files are converted where they are viewed."],
      ["examples/preview-server", "Convert on the server", "The server makes the preview and the thumbnail once. Browsers only draw."],
      ["examples/search", "With a search engine", "The text goes to a search engine page by page. A hit opens the document at its page."],
    ],
    read: "Read more",
    whyTitle: "Why bdf",
    why: [
      ["No office suite to run", "No headless LibreOffice, no cgo, no external program: the converters are Go packages in one binary."],
      ["In the shape of the content", "Fixed pages for slides and drawings, endless sheets for workbooks, pages or one column for text."],
      ["Made for browsers", "Instructions map onto Canvas 2D. Fonts, images and compression are left to the browser."],
      ["Many formats, one renderer", "Search, text selection and the text for screen readers work the same for every format."],
    ],
    whyMore: "More about why",
    license: "MIT License",
  },
  ja: {
    title: "bdf – ブラウザのための文書プレビュー",
    description: "bdf はブラウザでのプレビューのための文書フォーマットです。PDF、Office のファイル、図、CAD の図面を Go のバイナリ 1 つで変換し、canvas に描きます。",
    h1: "オフィススイートなしで、文書をブラウザでプレビュー",
    hero: (kb) => `bdf はプレビューのための文書フォーマットです。PDF、Office のファイル、図、CAD の図面などを Go のバイナリ 1 つで、あるいは WebAssembly にしてブラウザの中で変換し、${kb} KB のレンダラが canvas に描きます。`,
    try: "ビューアを開く", docs: "ドキュメントを読む",
    demoTitle: "bdf から描いた PDF",
    demoLead: "サーバーで変換するのと同じように、あらかじめ変換してあります。テキストの選択、検索、拡大ができます。",
    demoName: "demo.pdf → demo.bdf",
    demoNote: (pdf, bdf) => `PDF で ${pdf}、フォントを含む bdf で ${bdf}。ページは表示されるときに Range リクエストで取得します。`,
    search: "検索…", prev: "前の一致", next: "次の一致", zoomIn: "拡大", zoomOut: "縮小",
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
      ["examples/light-server", "ブラウザで変換する", "サーバーはファイルを置き、サムネイルだけを作ります。変換は見る人のブラウザで。"],
      ["examples/preview-server", "サーバーで変換する", "プレビューとサムネイルをサーバーが一度だけ作ります。ブラウザは描くだけ。"],
      ["examples/search", "検索エンジンとつなぐ", "テキストをページごとに検索エンジンへ。ヒットからそのページを開きます。"],
    ],
    read: "解説を読む",
    whyTitle: "なぜ bdf か",
    why: [
      ["オフィススイートを動かさない", "ヘッドレスの LibreOffice も cgo も外部プログラムも不要。変換器は 1 つのバイナリに入る Go のパッケージです。"],
      ["内容に合った形で見せる", "スライドや図面には固定ページ、ブックには無限のシート、文章にはページか一続きの列。"],
      ["ブラウザ表示に特化", "命令は Canvas 2D に対応。フォント、画像、展開はブラウザに任せます。"],
      ["多くの形式を 1 つのレンダラで", "検索、テキスト選択、読み上げ用のテキストは、どの形式でも同じように使えます。"],
    ],
    whyMore: "もっと詳しく",
    license: "MIT ライセンス",
  },
};

const size = (n) => (n < 1024 * 1024 ? `${Math.round(n / 1024)} KB` : `${(n / 1024 / 1024).toFixed(1)} MB`);

/**
 * The top page in a language. demo gives the sizes of the PDF shown and of
 * the bdf made of it, formats how many input formats the converters take,
 * worker the size of the rendering worker, minified and gzipped.
 */
export function home(lang, { demo, formats, worker }) {
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
<span><button type="button" id="zoomOut" title="${t.zoomOut}" aria-label="${t.zoomOut}">−</button>
<button type="button" id="zoomIn" title="${t.zoomIn}" aria-label="${t.zoomIn}">+</button></span>
</div>
<div class="view" id="demo" data-src="demo.bdf"></div>
</div>
<p class="note">${esc(t.demoNote(size(demo.pdf), size(demo.bdf)))}</p>
</section>
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
<footer class="wrap">bdf · <a href="${REPO}/blob/main/LICENSE">${t.license}</a> · <a href="${REPO}">GitHub</a></footer>
<script type="module" src="home.js"></script>
</body>
</html>
`;
}
