package main

import "embed"

// publicFS is served at /public: the pages' scripts and style, written by
// hand (pw's asset pipeline is not used).
//
//go:embed public
var publicFS embed.FS
