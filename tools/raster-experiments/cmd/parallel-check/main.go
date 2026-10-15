package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"os"
	"sync"
)

func main() {
	dpi := flag.Float64("dpi", 72, "dpi")
	flag.Parse()
	results := []map[string]any{}
	for _, path := range flag.Args() {
		rd, e := bdf.OpenSingleFile(path)
		must(e)
		doc, e := rd.ToDocument()
		must(e)
		must(rd.Close())
		opts := &imagebdf.Options{NoSystemFonts: true, FontDirs: []string{"converter/pptx/testdata/fonts", "converter/docx/testdata/fonts", "fixture/testdata/fonts"}}
		base, e := imagebdf.New(doc, opts).PageDPI(doc.Views[0], 0, *dpi)
		must(e)
		want := sha256.Sum256(base.Pix)
		hashes := make([]string, 4)
		var wg sync.WaitGroup
		for i := range hashes {
			wg.Add(1)
			go func() {
				defer wg.Done()
				img, e := imagebdf.New(doc, opts).PageDPI(doc.Views[0], 0, *dpi)
				must(e)
				h := sha256.Sum256(img.Pix)
				if h != want {
					panic("parallel page changed pixels")
				}
				hashes[i] = hex.EncodeToString(h[:])
			}()
		}
		wg.Wait()
		results = append(results, map[string]any{"input": path, "dpi": *dpi, "parallel_sha256": hashes})
	}
	must(json.NewEncoder(os.Stdout).Encode(results))
}
func must(e error) {
	if e != nil {
		panic(e)
	}
}
