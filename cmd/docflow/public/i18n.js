// The language of the scripts' texts: the page's (the server renders a
// page in its reader's language and says which on <html lang>). The
// texts are written in English in the scripts; ja.js has their Japanese,
// loaded on a Japanese page only.

export const lang = document.documentElement.lang === "ja" ? "ja" : "en";

const dict = lang === "ja" ? (await import("./ja.js")).default : null;

// _ is a text in the page's language: the English as written, or its
// translation, each {name} filled from args. A text is a whole sentence,
// so that a translation puts the words in its own order, and it is
// written in the call (public_test.go checks every one has its
// translation): a count says each form, n === 1 ? _("{n} page") :
// _("{n} pages").
export function _(text, args) {
  const s = (dict && dict[text]) || text;
  return args ? s.replace(/\{(\w+)\}/g, (m, k) => (k in args ? String(args[k]) : m)) : s;
}

// sentences joins sentences as the page's language writes them: with a
// space in English, none after a Japanese full stop. Empty ones go.
export function sentences(...parts) {
  return parts.filter(Boolean).join(lang === "ja" ? "" : " ");
}

// locale is how dates are written: the browser's own way when it reads
// the page's language (en-GB stays en-GB), else the page's.
export const locale = (navigator.language || "").toLowerCase().startsWith(lang) ? undefined : lang;
