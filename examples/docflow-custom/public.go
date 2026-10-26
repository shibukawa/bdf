package main

import (
	"embed"
	"io/fs"
)

// The program's own public tree, which Popcorn Web serves at /public (pw
// builds public into dist/public, with its manifest, and this embeds
// it); Docflow's scripts and style are docflow/web's, at /_docflow/.

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
