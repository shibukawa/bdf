// The pages of the documentation and their order in the sidebar. A page is
// docs/<id>.md in English and docs/<id>.ja.md in Japanese, published as
// docs/<id>.html and docs/<id>.ja.html; the reference documents (only: "ja")
// are written in Japanese alone, as docs/<id>.md.

export const REPO = "https://github.com/shibukawa/bdf";
export const SITE = "https://shibukawa.github.io/bdf/";
export const LANGS = ["en", "ja"];

const page = (id, en, ja, more = {}) => ({ id, label: { en, ja }, ...more });

export const SECTIONS = [
  {
    title: { en: "Introduction", ja: "はじめに" },
    pages: [
      page("index", "Documentation", "ドキュメント"),
      page("why", "Why bdf", "なぜ bdf か"),
      page("getting-started", "Getting started", "はじめかた"),
      page("features", "Features", "特徴"),
    ],
  },
  {
    title: { en: "Formats", ja: "対応形式" },
    pages: [
      page("formats/index", "All formats", "形式の一覧"),
      page("formats/pdf", "PDF, Illustrator", "PDF・Illustrator"),
      page("formats/document", "Word, HTML, Markdown", "Word・HTML・Markdown"),
      page("formats/presentation", "PowerPoint", "PowerPoint"),
      page("formats/spreadsheet", "Excel, CSV, Parquet", "Excel・CSV・Parquet"),
      page("formats/diagram", "Visio, draw.io", "Visio・draw.io"),
      page("formats/cad", "CAD drawings and plots", "CAD の図面とプロット"),
      page("formats/electronics", "Circuit boards, KiCad", "プリント基板・KiCad"),
      page("formats/ebook", "EPUB", "EPUB"),
      page("formats/music", "Scores", "楽譜"),
      page("formats/font", "Font files", "フォントファイル"),
      page("formats/image", "Images, Photoshop", "画像・Photoshop"),
    ],
  },
  {
    title: { en: "Sample architectures", ja: "構成のサンプル" },
    pages: [
      page("examples/index", "Choosing one", "構成の選び方"),
      page("examples/light-server", "Convert in the browser", "ブラウザで変換する"),
      page("examples/preview-server", "Convert on the server", "サーバーで変換する"),
      page("examples/search", "With a search engine", "検索エンジンとつなぐ"),
      page("examples/secure-reader", "Books behind a login", "ログインした読者に読ませる"),
    ],
  },
  {
    title: { en: "Architecture", ja: "アーキテクチャ" },
    pages: [
      page("architecture/index", "How it works", "処理の流れ"),
      page("architecture/layout", "Repository layout", "フォルダ構成"),
      page("architecture/cli", "The bdf command", "bdf コマンド"),
      page("architecture/browser", "In the browser", "ブラウザ側"),
      page("architecture/development", "Building and testing", "ビルドとテスト"),
    ],
  },
  {
    title: { en: "Reference", ja: "リファレンス" },
    pages: [
      page("api", "API", "API 一覧", { only: "ja" }),
      page("spec", "Format specification", "フォーマット仕様", { only: "ja" }),
      page("design", "Design notes", "設計メモ", { only: "ja" }),
    ],
  },
];

export const PAGES = SECTIONS.flatMap((s) => s.pages);

/** The Markdown file of a page in a language (from the repository root); undefined when it has none. */
export function source(p, lang) {
  if (p.only) return p.only === lang ? `docs/${p.id}.md` : undefined;
  return lang === "ja" ? `docs/${p.id}.ja.md` : `docs/${p.id}.md`;
}

/** The file a page is published as in a language, from the site's root. */
export function output(p, lang) {
  return p.only || lang === "en" ? `docs/${p.id}.html` : `docs/${p.id}.ja.html`;
}

/** The words of the pages around their content. */
export const UI = {
  en: {
    skip: "Skip to content", site: "Site", docs: "Documentation", menu: "Menu", toc: "On this page", source: "Source", other: "日本語",
    viewer: "Viewer", thumbnail: "Thumbnails", text: "Search text", docsNav: "Docs", onlyJa: " (Japanese)",
  },
  ja: {
    skip: "本文へ移動", site: "サイト", docs: "ドキュメント", menu: "メニュー", toc: "このページの内容", source: "原文", other: "English",
    viewer: "ビューア", thumbnail: "サムネイル", text: "検索テキスト", docsNav: "ドキュメント", onlyJa: "",
  },
};
