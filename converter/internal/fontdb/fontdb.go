// Package fontdb finds font files, resolves the family names documents ask
// for to faces that are available (directly, through metric-compatible
// substitutes, or through generic fallbacks), measures glyphs for text
// layout, and builds the subset font programs that get embedded.
package fontdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf16"

	"github.com/shibukawa/bdf/converter/internal/sfnt"
	"golang.org/x/text/unicode/norm"
)

// Face is one font in a file (a collection holds several).
type Face struct {
	Path   string
	Index  int    // index in a collection
	Family string // English family name (typographic family when present)
	Style  string // English subfamily name
	Weight int    // OS/2 usWeightClass
	Italic bool
	Mono   bool
	CFF    bool

	names []string // normalized family names in every language
	fsys  fs.FS    // the file system Path is in; nil for the local one

	once   sync.Once
	loaded *Loaded
	err    error
}

// DB is a set of faces indexed by normalized family name.
type DB struct {
	Faces  []*Face
	byName map[string][]*Face
}

var (
	scanMu    sync.Mutex
	scanCache = map[string][]*Face{}
)

// SystemDirs returns the usual font directories of the platform.
func SystemDirs() []string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		dirs := []string{filepath.Join(os.Getenv("WINDIR"), "Fonts")}
		if la := os.Getenv("LOCALAPPDATA"); la != "" {
			dirs = append(dirs, filepath.Join(la, "Microsoft", "Windows", "Fonts"))
		}
		return dirs
	case "darwin":
		return []string{"/System/Library/Fonts", "/Library/Fonts", filepath.Join(home, "Library", "Fonts")}
	default:
		dirs := []string{"/usr/share/fonts", "/usr/local/share/fonts"}
		if home != "" {
			dirs = append(dirs, filepath.Join(home, ".fonts"), filepath.Join(home, ".local", "share", "fonts"))
		}
		return dirs
	}
}

// New scans the fonts in fsys (when it is not nil), then those under dirs
// and, when system is true, the system directories. Scans of a directory
// are cached for the life of the process; fsys is scanned on every call.
func New(fsys fs.FS, dirs []string, system bool) *DB {
	all := append([]string(nil), dirs...)
	if system {
		all = append(all, SystemDirs()...)
	}
	db := &DB{byName: map[string][]*Face{}}
	if fsys != nil {
		db.Faces = scanFS(fsys)
	}
	seen := map[string]bool{}
	for _, d := range all {
		for _, f := range scanDir(d) {
			key := f.Path + "#" + strconv.Itoa(f.Index)
			if seen[key] {
				continue
			}
			seen[key] = true
			db.Faces = append(db.Faces, f)
		}
	}
	for _, f := range db.Faces {
		for _, n := range f.names {
			db.byName[n] = append(db.byName[n], f)
		}
	}
	return db
}

// isFontFile reports whether a file name extension is that of a font file.
func isFontFile(ext string) bool {
	switch strings.ToLower(ext) {
	case ".ttf", ".otf", ".ttc", ".otc":
		return true
	}
	return false
}

func scanDir(dir string) []*Face {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	scanMu.Lock()
	defer scanMu.Unlock()
	if faces, ok := scanCache[abs]; ok {
		return faces
	}
	var paths []string
	filepath.WalkDir(abs, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && isFontFile(filepath.Ext(p)) {
			paths = append(paths, p)
		}
		return nil
	})
	sort.Strings(paths)
	var faces []*Face
	for _, p := range paths {
		faces = append(faces, scanFile(p)...)
	}
	scanCache[abs] = faces
	return faces
}

// scanFS is scanDir for the fonts of a file system other than the local
// one (embedded in the program, or fetched). Its files are read with ReadAt
// when they have it, so that only the tables scanFaces needs are read.
func scanFS(fsys fs.FS) []*Face {
	var paths []string
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && isFontFile(path.Ext(p)) {
			paths = append(paths, p)
		}
		return nil
	})
	sort.Strings(paths)
	var faces []*Face
	for _, p := range paths {
		f, err := fsys.Open(p)
		if err != nil {
			continue
		}
		r, ok := f.(io.ReaderAt)
		if !ok {
			data, err := io.ReadAll(f)
			if err != nil {
				f.Close()
				continue
			}
			r = bytes.NewReader(data)
		}
		for _, face := range scanFaces(r, p) {
			face.fsys = fsys
			faces = append(faces, face)
		}
		f.Close()
	}
	return faces
}

// scanFile reads the naming and style information of every face in a file
// without loading the glyph data.
func scanFile(path string) []*Face {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	return scanFaces(f, path)
}

// scanFaces is scanFile for a file opened as r.
func scanFaces(r io.ReaderAt, path string) []*Face {
	var hdr [12]byte
	if _, err := r.ReadAt(hdr[:], 0); err != nil {
		return nil
	}
	offsets := []int64{0}
	if string(hdr[:4]) == "ttcf" {
		n := int(binary.BigEndian.Uint32(hdr[8:]))
		if n <= 0 || n > 256 {
			return nil
		}
		buf := make([]byte, 4*n)
		if _, err := r.ReadAt(buf, 12); err != nil {
			return nil
		}
		offsets = offsets[:0]
		for i := 0; i < n; i++ {
			offsets = append(offsets, int64(binary.BigEndian.Uint32(buf[i*4:])))
		}
	}
	var out []*Face
	for i, off := range offsets {
		if face := scanFace(r, off); face != nil {
			face.Path, face.Index = path, i
			out = append(out, face)
		}
	}
	return out
}

func scanFace(r io.ReaderAt, off int64) *Face {
	var hdr [12]byte
	if _, err := r.ReadAt(hdr[:], off); err != nil {
		return nil
	}
	tag := string(hdr[:4])
	if tag != "\x00\x01\x00\x00" && tag != "true" && tag != "OTTO" {
		return nil
	}
	n := int(binary.BigEndian.Uint16(hdr[4:]))
	dir := make([]byte, 16*n)
	if _, err := r.ReadAt(dir, off+12); err != nil {
		return nil
	}
	table := func(name string) []byte {
		for i := 0; i < n; i++ {
			rec := dir[i*16:]
			if string(rec[:4]) != name {
				continue
			}
			toff, tlen := binary.BigEndian.Uint32(rec[8:]), binary.BigEndian.Uint32(rec[12:])
			if tlen > 1<<20 {
				return nil
			}
			b := make([]byte, tlen)
			if _, err := r.ReadAt(b, int64(toff)); err != nil {
				return nil
			}
			return b
		}
		return nil
	}
	face := &Face{Weight: 400, CFF: tag == "OTTO"}
	names := parseNames(table("name"))
	if len(names[1]) == 0 {
		return nil
	}
	face.Family = pickEnglish(names[16], pickEnglish(names[1], ""))
	face.Style = pickEnglish(names[17], pickEnglish(names[2], "Regular"))
	seen := map[string]bool{}
	for _, id := range []uint16{1, 16} {
		for _, nm := range names[id] {
			if k := Normalize(nm.val); k != "" && !seen[k] {
				seen[k] = true
				face.names = append(face.names, k)
			}
		}
	}
	if os2 := table("OS/2"); len(os2) >= 64 {
		face.Weight = int(binary.BigEndian.Uint16(os2[4:]))
		face.Italic = binary.BigEndian.Uint16(os2[62:])&1 != 0
	}
	if face.Weight < 100 || face.Weight > 1000 {
		face.Weight = 400
	}
	if post := table("post"); len(post) >= 16 {
		face.Mono = binary.BigEndian.Uint32(post[12:]) != 0
	}
	st := strings.ToLower(face.Style)
	if strings.Contains(st, "italic") || strings.Contains(st, "oblique") {
		face.Italic = true
	}
	return face
}

type nameRec struct {
	lang uint16
	val  string
}

// parseNames returns the name records by name ID, decoded from the Windows
// (UTF-16) and Macintosh Roman platforms.
func parseNames(b []byte) map[uint16][]nameRec {
	out := map[uint16][]nameRec{}
	if len(b) < 6 {
		return out
	}
	count, strOff := int(binary.BigEndian.Uint16(b[2:])), int(binary.BigEndian.Uint16(b[4:]))
	for i := 0; i < count; i++ {
		rec := 6 + i*12
		if rec+12 > len(b) {
			break
		}
		pid, eid := binary.BigEndian.Uint16(b[rec:]), binary.BigEndian.Uint16(b[rec+2:])
		lang, id := binary.BigEndian.Uint16(b[rec+4:]), binary.BigEndian.Uint16(b[rec+6:])
		l, o := int(binary.BigEndian.Uint16(b[rec+8:])), int(binary.BigEndian.Uint16(b[rec+10:]))
		if id != 1 && id != 2 && id != 16 && id != 17 {
			continue
		}
		s := strOff + o
		if s+l > len(b) {
			continue
		}
		raw := b[s : s+l]
		var val string
		switch {
		case pid == 3 && (eid == 1 || eid == 10 || eid == 0), pid == 0:
			u := make([]uint16, len(raw)/2)
			for j := range u {
				u[j] = binary.BigEndian.Uint16(raw[j*2:])
			}
			val = string(utf16.Decode(u))
		case pid == 1 && eid == 0:
			val = string(raw) // Mac Roman; family names are ASCII in practice
			lang = 0x0409
		default:
			continue
		}
		if val = strings.TrimSpace(val); val != "" {
			out[id] = append(out[id], nameRec{lang, val})
		}
	}
	return out
}

func pickEnglish(recs []nameRec, dflt string) string {
	for _, r := range recs {
		if r.lang == 0x0409 {
			return r.val
		}
	}
	if len(recs) > 0 {
		return recs[0].val
	}
	return dflt
}

// Normalize folds a family name for comparison: NFKC (full-width letters
// become ASCII), lower case, without spaces, hyphens and underscores.
func Normalize(name string) string {
	s := strings.ToLower(norm.NFKC.String(name))
	return strings.Map(func(r rune) rune {
		switch r {
		case ' ', '-', '_', '　':
			return -1
		}
		return r
	}, s)
}

// Family returns the faces of a family, or nil.
func (db *DB) Family(name string) []*Face { return db.byName[Normalize(name)] }

// Match picks the face of a family closest to the requested weight and
// slant. synthBold / synthItalic report what the renderer has to fake.
func Match(faces []*Face, bold, italic bool) (best *Face, synthBold, synthItalic bool) {
	want := 400
	if bold {
		want = 700
	}
	bestScore := 1 << 30
	for _, f := range faces {
		score := abs(f.Weight-want) * 2
		if bold && f.Weight < 600 || !bold && f.Weight >= 600 {
			score += 1000
		}
		if f.Italic != italic {
			score += 5000
		}
		if score < bestScore {
			best, bestScore = f, score
		}
	}
	if best == nil {
		return nil, false, false
	}
	return best, bold && best.Weight < 600, italic && !best.Italic
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// Loaded is a face with its glyph data.
type Loaded struct {
	Face *Face
	Font *sfnt.Font
	Data []byte // the whole file (the font's tables point into it)
	upem float64
	// vertical metrics in em
	Ascent, Descent       float64
	UnderlinePos, UnderTh float64
	// CapHeight is the height of capital letters in em (estimated when
	// the font does not say).
	CapHeight     float64
	embed, subset bool
}

// ErrNoGlyphs is returned for faces whose glyph data cannot be read.
var ErrNoGlyphs = errors.New("fontdb: font has no usable glyph data")

// Load reads and parses the face's font program (once).
func (f *Face) Load() (*Loaded, error) {
	f.once.Do(func() {
		var data []byte
		var err error
		if f.fsys != nil {
			data, err = fs.ReadFile(f.fsys, f.Path)
		} else {
			data, err = os.ReadFile(f.Path)
		}
		if err != nil {
			f.err = err
			return
		}
		sf, err := sfnt.ParseIndex(data, f.Index)
		if err != nil {
			f.err = err
			return
		}
		if sf.Cmap == nil || len(sf.Advances) == 0 {
			f.err = ErrNoGlyphs
			return
		}
		l := &Loaded{Face: f, Font: sf, Data: data, upem: float64(sf.UnitsPerEm), embed: true, subset: true}
		os2, ok := sf.OS2()
		switch {
		case ok && os2.WinAscent+os2.WinDescent > 0:
			// Office lays lines out with the Windows (GDI) metrics.
			l.Ascent, l.Descent = float64(os2.WinAscent)/l.upem, float64(os2.WinDescent)/l.upem
		case sf.Ascent != 0 || sf.Descent != 0:
			l.Ascent, l.Descent = float64(sf.Ascent)/l.upem, float64(-sf.Descent+sf.LineGap)/l.upem
		default:
			l.Ascent, l.Descent = 0.9, 0.3
		}
		if ok {
			l.embed, l.subset = os2.Embeddable()
		}
		l.CapHeight = 0.7
		if ok && os2.CapHeight > 0 {
			l.CapHeight = float64(os2.CapHeight) / l.upem
		}
		pos, th := sf.Underline()
		l.UnderlinePos, l.UnderTh = float64(pos)/l.upem, float64(th)/l.upem
		f.loaded = l
	})
	return f.loaded, f.err
}

// Glyph returns the glyph index for r.
func (l *Loaded) Glyph(r rune) (uint16, bool) {
	g, ok := l.Font.Cmap[uint32(r)]
	return g, ok && g != 0
}

// Has reports whether the font maps r to a glyph.
func (l *Loaded) Has(r rune) bool {
	_, ok := l.Glyph(r)
	return ok
}

// Advance returns the advance of glyph g in em.
func (l *Loaded) Advance(g uint16) float64 { return float64(l.Font.Advance(g)) / l.upem }

// Embeddable reports whether the license bits allow embedding, and subsetting.
func (l *Loaded) Embeddable() (embed, subset bool) { return l.embed, l.subset }
