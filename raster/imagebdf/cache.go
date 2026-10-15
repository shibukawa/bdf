package imagebdf

import (
	"crypto/sha256"
	"encoding/json"
	"image"
	"math"
	"strings"
	"sync"
	"unsafe"

	"github.com/shibukawa/bdf"
)

const (
	// DefaultImageCacheBytes is the completed-image cache budget per Renderer.
	DefaultImageCacheBytes = 64 << 20
	// DefaultTextCacheBytes is the text-run coverage cache budget per Renderer.
	DefaultTextCacheBytes = 32 << 20
)

func cacheLimit(n, fallback int) int {
	if n == 0 {
		return fallback
	}
	return max(n, 0)
}

// The accounting includes backing arrays, both copies of map/entry keys,
// and conservative entry/bucket overhead. Other decoded resources, scratch
// and Go allocator overhead are outside these cache budgets.
type cacheEntry[K comparable, V any] struct {
	key        K
	value      V
	cost       int
	prev, next *cacheEntry[K, V]
}
type lru[K comparable, V any] struct {
	limit, bytes int
	items        map[K]*cacheEntry[K, V]
	first, last  *cacheEntry[K, V]
}

func entryCost[K comparable, V any](payload int) int {
	return payload + int(unsafe.Sizeof(cacheEntry[K, V]{})) + int(unsafe.Sizeof(*new(K))) + 256
}
func (c *lru[K, V]) unlink(e *cacheEntry[K, V]) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		c.first = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		c.last = e.prev
	}
	e.prev, e.next = nil, nil
}
func (c *lru[K, V]) front(e *cacheEntry[K, V]) {
	e.next = c.first
	if c.first != nil {
		c.first.prev = e
	} else {
		c.last = e
	}
	c.first = e
}
func (c *lru[K, V]) get(k K) (V, bool) {
	if e := c.items[k]; e != nil {
		c.unlink(e)
		c.front(e)
		return e.value, true
	}
	var zero V
	return zero, false
}
func (c *lru[K, V]) put(k K, v V, cost int) {
	if cost > c.limit || c.limit <= 0 {
		return
	}
	if old := c.items[k]; old != nil {
		c.unlink(old)
		c.bytes -= old.cost
		delete(c.items, k)
	}
	for c.bytes > c.limit-cost {
		e := c.last
		c.unlink(e)
		delete(c.items, e.key)
		c.bytes -= e.cost
	}
	if c.items == nil {
		c.items = make(map[K]*cacheEntry[K, V])
	}
	e := &cacheEntry[K, V]{key: k, value: v, cost: cost}
	c.items[k] = e
	c.bytes += cost
	c.front(e)
}

// Invalidate releases all cached images, text masks and decoded resources,
// clears warnings, and retains the renderer's options. Call it after changing
// document part bytes or FontFS. FontDirs and system fonts retain fontdb's
// existing process-wide scan/load cache; use FontFS for replaceable fonts.
// Document/view mutation
// must never race with rendering. Page and view layout are also included in
// completed-image keys, so changes to those layouts cannot return old pixels.
func (r *Renderer) Invalidate() { *r = *New(r.doc, &r.opts); r.fonts.r = r }

type imageKey struct {
	kind                byte
	source              [32]byte
	region              bdf.Rect
	w, h                int
	scale               float64
	noAA, backgroundSet bool
	background          solid
	mono                Mono
}

func (r *Renderer) imageKey(kind byte, source any, region bdf.Rect, w, h int, scale float64, mono Mono) (imageKey, bool) {
	if plain || r.images.limit <= 0 || !finiteCacheNumbers(float64(region.X), float64(region.Y), float64(region.W), float64(region.H), scale) {
		return imageKey{}, false
	}
	data, err := json.Marshal(source)
	if err != nil {
		return imageKey{}, false
	}
	return imageKey{kind: kind, source: sha256.Sum256(data), region: region, w: w, h: h, scale: scale, noAA: r.opts.Antialias == AntialiasOff, backgroundSet: r.opts.Background != nil, background: r.background(), mono: mono}, true
}
func cloneRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Rect)
	for y := 0; y < src.Rect.Dy(); y++ {
		copy(dst.Pix[y*dst.Stride:y*dst.Stride+4*src.Rect.Dx()], src.Pix[y*src.Stride:y*src.Stride+4*src.Rect.Dx()])
	}
	return dst
}

// cloneImage copies an image of the kinds the renderer draws (see Mono).
func cloneImage(src image.Image) image.Image {
	switch src := src.(type) {
	case *image.RGBA:
		return cloneRGBA(src)
	case *image.Gray:
		dst := image.NewGray(src.Rect)
		copy(dst.Pix, src.Pix)
		return dst
	case *image.Paletted:
		dst := image.NewPaletted(src.Rect, src.Palette)
		copy(dst.Pix, src.Pix)
		return dst
	}
	return src
}

// imagePix is the memory of the pixels of an image the renderer draws.
func imagePix(img image.Image) int {
	switch img := img.(type) {
	case *image.RGBA:
		return len(img.Pix) + int(unsafe.Sizeof(*img))
	case *image.Gray:
		return len(img.Pix) + int(unsafe.Sizeof(*img))
	case *image.Paletted:
		return len(img.Pix) + int(unsafe.Sizeof(*img))
	}
	return 0
}
func (r *Renderer) storedImage(k imageKey, ok bool) image.Image {
	if ok {
		if img, found := r.images.get(k); found {
			return cloneImage(img)
		}
	}
	return nil
}
func (r *Renderer) keepImage(k imageKey, ok bool, img image.Image) {
	if !ok || img == nil {
		return
	}
	cost := entryCost[imageKey, image.Image](imagePix(img))
	if cost <= r.images.limit {
		r.images.put(k, cloneImage(img), cost)
	}
}
func (r *Renderer) cachedPage(p *bdf.Page, region bdf.Rect, w, h int, mono Mono) image.Image {
	k, ok := r.imageKey(1, p, region, w, h, 0, mono)
	if img := r.storedImage(k, ok); img != nil {
		return img
	}
	img := r.drawPage(p, region, w, h, mono)
	r.keepImage(k, ok, img)
	return img
}
func (r *Renderer) cachedView(v *bdf.View, region bdf.Rect, w, h int, sheet bool, mono Mono) image.Image {
	kind := byte(2)
	if sheet {
		kind = 3
	}
	k, ok := r.imageKey(kind, v, region, w, h, 0, mono)
	if img := r.storedImage(k, ok); img != nil {
		return img
	}
	var img image.Image
	if sheet {
		img = r.drawSheet(v, region, w, h, mono)
	} else {
		img = r.drawContinuous(v, region, w, h, mono)
	}
	r.keepImage(k, ok, img)
	return img
}

type textMaskKey struct {
	font          bdf.Font
	text          string
	m             matrix
	size, spacing float64
	dir           byte
	bounds        image.Rectangle
	noAA          bool
}

func finiteCacheNumbers(values ...float64) bool {
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}

func (r *Renderer) textMask(k textMaskKey) *mask {
	if plain || r.textMasks.limit <= 0 || !k.m.finite() || !finiteCacheNumbers(k.size, k.spacing) {
		return nil
	}
	m, _ := r.textMasks.get(k)
	return m
}
func (r *Renderer) keepTextMask(k textMaskKey, m *mask) {
	if plain || r.textMasks.limit <= 0 || m == nil || !k.m.finite() || !finiteCacheNumbers(k.size, k.spacing) {
		return
	}
	cost := entryCost[textMaskKey, *mask](len(m.a)*4 + len(k.text) + len(k.font.Family) + int(unsafe.Sizeof(*m)))
	if cost > r.textMasks.limit {
		return
	}
	k.text = strings.Clone(k.text)
	k.font.Family = strings.Clone(k.font.Family)
	a := make([]float32, len(m.a))
	copy(a, m.a)
	if m.a == nil {
		a = nil
	}
	r.textMasks.put(k, &mask{r: m.r, a: a}, cost)
}

// Shared by page workers; stored masks are immutable after insertion.
type textMaskCache struct {
	sync.Mutex
	lru[textMaskKey, *mask]
}

func (c *textMaskCache) get(k textMaskKey) (*mask, bool) {
	c.Lock()
	defer c.Unlock()
	return c.lru.get(k)
}
func (c *textMaskCache) put(k textMaskKey, m *mask, n int) {
	c.Lock()
	defer c.Unlock()
	c.lru.put(k, m, n)
}

// freeBudget is the memory of the canvases the renderer keeps for its
// next drawings: its share of a batch (Pages), else Options.CanvasCacheBytes.
func (r *Renderer) freeBudget() int {
	if r.freeLimit > 0 {
		return r.freeLimit
	}
	return cacheLimit(r.opts.CanvasCacheBytes, DefaultCanvasCacheBytes)
}
