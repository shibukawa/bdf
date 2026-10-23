// Package messages is the catalog of the texts of Docflow's pages, in
// English and Japanese: pages.yaml, compiled by pw generate into
// messages_pw_gen.go, which the templates of package handlers call with
// the locale of the request (popcornweb.toml's [i18n]). Go calls the
// same functions with a locale: PagesHelpTitle(pw.RequestLocale(r), …).
package messages
