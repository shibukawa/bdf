module github.com/shibukawa/bdf/raster/ebitenginebdf

go 1.27

require (
	github.com/hajimehoshi/ebiten/v2 v2.10.4
	github.com/shibukawa/bdf v0.1.0
)

require (
	github.com/andybalholm/brotli v1.2.5 // indirect
	github.com/ebitengine/gomobile v0.0.0-20260820040257-d11f821a26a6 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.11.0 // indirect
	github.com/klauspost/compress v1.20.1 // indirect
	github.com/shibukawa/tinygodriver v1.3.4 // indirect
	golang.org/x/image v0.46.0 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
)

// The renderer is released with the module it draws the objects of; in this
// repository it is built against the checkout.
replace github.com/shibukawa/bdf => ../..
