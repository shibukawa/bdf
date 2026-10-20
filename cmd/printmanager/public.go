package main

import "embed"

// publicFS is served at /static/ (handlers.StaticPath): the pages'
// scripts, and their style
// built with Tailwind CSS and daisyUI (pw's asset pipeline is not used).
// go generate builds the style with the standalone CLI of the version in
// tailwind/VERSIONS; public/app.css is committed.
//
//go:generate tailwindcss --input tailwind/app.css --output public/app.css --minify
//go:embed public
var publicFS embed.FS
