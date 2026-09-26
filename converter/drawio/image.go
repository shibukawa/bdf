package drawio

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// imageRef is an image stored in the document.
type imageRef struct {
	hash bdf.Hash
	w, h float64 // intrinsic size in pixels (0 when unknown)
}

// imageFor stores the image an image style names: a data URI (base64 or
// URL-encoded, as draw.io writes SVG images). Images on the web are not
// fetched: the conversion does not depend on the network, and the shape
// is drawn without its picture (with a warning).
func (c *converter) imageFor(src string) *imageRef {
	src = strings.TrimSpace(src)
	if src == "" {
		return nil
	}
	if r, ok := c.images[src]; ok {
		return r
	}
	var ref *imageRef
	defer func() { c.images[src] = ref }()
	data, err := decodeDataURI(src)
	if err != nil {
		if strings.HasPrefix(src, "data:") {
			c.warnOnce("imgdata:"+src[:min(len(src), 40)], "image data URI cannot be read: %v", err)
		} else {
			c.warnOnce("imgurl", "images referenced by URL are not embedded (e.g. %s); export the diagram with embedded images", truncate(src, 80))
		}
		return nil
	}
	format := imgconv.Sniff(data)
	switch format {
	case "":
		c.warnOnce("imgfmt", "unsupported image format in a data URI")
		return nil
	case "svg":
		ref = &imageRef{}
		ref.w, ref.h = svgSize(data)
	default:
		ref = &imageRef{}
		if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
			ref.w, ref.h = float64(cfg.Width), float64(cfg.Height)
		}
		if res, err := imgconv.Optimize(data, c.opts.Images); err == nil || res.Data != nil {
			data = res.Data
		}
	}
	ref.hash = c.doc.AddImage(data)
	return ref
}

// decodeDataURI returns the bytes of a data: URI.
func decodeDataURI(s string) ([]byte, error) {
	if !strings.HasPrefix(s, "data:") {
		return nil, errNotData
	}
	comma := strings.IndexByte(s, ',')
	if comma < 0 {
		return nil, errNotData
	}
	meta, payload := s[5:comma], s[comma+1:]
	if strings.HasSuffix(meta, ";base64") {
		payload = strings.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, payload)
		if u, err := url.PathUnescape(payload); err == nil {
			payload = u
		}
		b, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			b, err = base64.RawStdEncoding.DecodeString(strings.TrimRight(payload, "="))
		}
		return b, err
	}
	if u, err := url.PathUnescape(payload); err == nil {
		payload = u
	}
	// draw.io also writes "data:image/png,<base64>" without the ;base64 marker
	// (Graph.postProcessCellStyle adds it)
	if !strings.HasPrefix(meta, "image/svg") {
		if b, err := base64.StdEncoding.DecodeString(payload); err == nil {
			return b, nil
		}
	}
	return []byte(payload), nil
}

type dataErr string

func (e dataErr) Error() string { return string(e) }

const errNotData = dataErr("not a data URI")

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

var (
	svgTag     = regexp.MustCompile(`(?s)<svg\b[^>]*>`)
	svgAttr    = regexp.MustCompile(`\b(width|height|viewBox)\s*=\s*["']([^"']*)["']`)
	svgNumUnit = regexp.MustCompile(`^\s*([0-9.]+)\s*(px|pt)?\s*$`)
)

// svgSize reads the intrinsic size of an SVG document.
func svgSize(data []byte) (w, h float64) {
	tag := svgTag.Find(data)
	var vb []float64
	for _, m := range svgAttr.FindAllSubmatch(tag, -1) {
		v := string(m[2])
		switch string(m[1]) {
		case "width", "height":
			if mm := svgNumUnit.FindStringSubmatch(v); mm != nil {
				f, _ := strconv.ParseFloat(mm[1], 64)
				if mm[2] == "pt" {
					f *= 96.0 / 72
				}
				if string(m[1]) == "width" {
					w = f
				} else {
					h = f
				}
			}
		case "viewBox":
			for _, s := range strings.FieldsFunc(v, func(r rune) bool { return r == ' ' || r == ',' }) {
				f, _ := strconv.ParseFloat(s, 64)
				vb = append(vb, f)
			}
		}
	}
	if (w == 0 || h == 0) && len(vb) == 4 {
		w, h = vb[2], vb[3]
	}
	return w, h
}
