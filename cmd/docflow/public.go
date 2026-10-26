package main

import (
	"embed"
	"io/fs"
)

// The program's own public tree, which Popcorn Web serves at /public:
// this one has nothing of its own (Docflow's scripts and style are
// docflow/web's, served at /_docflow/), but pw builds public into
// dist/public with its manifest for every program, and a program
// composing Docflow puts its own files here (web.Options.Public).

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
