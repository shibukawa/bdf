package main

import (
	"embed"
	"io/fs"
)

// The pages' scripts and style are written in public (the style built
// there by Tailwind CSS with daisyUI, from tailwind/app.css, by go
// generate with the standalone CLI of tailwind/VERSIONS). pw generate
// builds them into dist/public, with the manifest beside this file, and
// that tree is what is embedded and served at /public. public/app.css,
// dist/public and the manifest are committed, so that a plain go build
// carries the files: CI builds them again and fails on a difference.

//go:generate tailwindcss --input tailwind/app.css --output public/app.css --minify

//go:embed all:dist/public
var embeddedPublic embed.FS

// publicTree is the tree served at /public.
func publicTree() fs.FS {
	tree, err := fs.Sub(embeddedPublic, "dist/public")
	if err != nil {
		panic(err)
	}
	return tree
}
