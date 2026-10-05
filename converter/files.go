package converter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/shibukawa/bdf/converter/internal/localfile"
)

// FileMap is a file system of files in memory, by name (Options.Files): a
// server lists the files uploaded with an input.
//
//	opts.Files = converter.FileMap{"power.kicad_sch": power, "io.kicad_sch": io}
//
// Names are slash-separated paths; the directories they imply are listed
// too.
type FileMap map[string][]byte

// Open implements fs.FS.
func (m FileMap) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	if b, ok := m[name]; ok {
		return &memFile{Reader: bytes.NewReader(b), info: memInfo{name: path.Base(name), size: int64(len(b))}}, nil
	}
	// a directory: the names under it
	prefix := name + "/"
	if name == "." {
		prefix = ""
	}
	seen := map[string]bool{}
	var entries []fs.DirEntry
	for n, b := range m {
		n = path.Clean(strings.TrimPrefix(n, "./"))
		if !strings.HasPrefix(n, prefix) {
			continue
		}
		rest := n[len(prefix):]
		first, _, isDir := strings.Cut(rest, "/")
		if first == "" || seen[first] {
			continue
		}
		seen[first] = true
		info := memInfo{name: first, size: int64(len(b)), dir: isDir}
		entries = append(entries, fs.FileInfoToDirEntry(info))
	}
	if len(entries) == 0 && name != "." {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	slices.SortFunc(entries, func(a, b fs.DirEntry) int { return strings.Compare(a.Name(), b.Name()) })
	return &memDir{info: memInfo{name: path.Base(name), dir: true}, entries: entries}, nil
}

type memInfo struct {
	name string
	size int64
	dir  bool
}

func (i memInfo) Name() string { return i.name }
func (i memInfo) Size() int64 {
	if i.dir {
		return 0
	}
	return i.size
}
func (i memInfo) Mode() fs.FileMode {
	if i.dir {
		return fs.ModeDir | 0o555
	}
	return 0o444
}
func (i memInfo) ModTime() time.Time { return time.Time{} }
func (i memInfo) IsDir() bool        { return i.dir }
func (i memInfo) Sys() any           { return nil }

type memFile struct {
	*bytes.Reader
	info memInfo
}

func (f *memFile) Stat() (fs.FileInfo, error) { return f.info, nil }
func (f *memFile) Close() error               { return nil }

type memDir struct {
	info    memInfo
	entries []fs.DirEntry
	pos     int
}

func (d *memDir) Stat() (fs.FileInfo, error) { return d.info, nil }
func (d *memDir) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.info.name, Err: errors.New("is a directory")}
}
func (d *memDir) Close() error { return nil }
func (d *memDir) ReadDir(n int) ([]fs.DirEntry, error) {
	rest := d.entries[d.pos:]
	if n <= 0 {
		d.pos = len(d.entries)
		return rest, nil
	}
	if len(rest) == 0 {
		return nil, io.EOF
	}
	n = min(n, len(rest))
	d.pos += n
	return rest[:n], nil
}

// MaxRefSize bounds the files ReadRef reads.
const MaxRefSize = 256 << 20

// ReadRef reads a file an input refers to. rel is the reference as the
// input writes it, relative to the directory of the file that makes it:
// from, a slash path relative to the input's directory ("" for the input
// itself). It looks in o.Files and then in o.Dir, and in Files also by the
// base name of rel when no file has its path. It returns the contents and
// the file's path relative to the input's directory (slash-separated), for
// the references the file makes in turn. Absolute references are not read;
// in Dir, references and symlinks must stay within the directory, while
// in Files a reference outside them matches by base name only.
func (o *Options) ReadRef(from, rel string) ([]byte, string, error) {
	rel = strings.ReplaceAll(rel, `\`, "/")
	if rel == "" || path.IsAbs(rel) || filepath.IsAbs(rel) || (len(rel) > 1 && rel[1] == ':') {
		return nil, "", fmt.Errorf("%q: only files relative to the input are read", rel)
	}
	want := path.Clean(path.Join(path.Dir(from), rel))
	if o.Files != nil {
		if strings.HasPrefix(want, "../") || want == ".." {
			// outside the listed files: only the base name can match
		} else if b, err := readAll(o.Files, want); err == nil {
			return b, want, nil
		}
		if name := findBase(o.Files, path.Base(rel)); name != "" {
			if b, err := readAll(o.Files, name); err == nil {
				return b, name, nil
			}
		}
	}
	if o.Dir != "" {
		p := filepath.Join(o.Dir, filepath.FromSlash(want))
		f, err := localfile.Open(o.Dir, filepath.FromSlash(want))
		if err == nil {
			defer f.Close()
			b, err := io.ReadAll(io.LimitReader(f, MaxRefSize+1))
			if err != nil {
				return nil, "", err
			}
			if len(b) > MaxRefSize {
				return nil, "", fmt.Errorf("%s: larger than %d bytes", p, MaxRefSize)
			}
			return b, want, nil
		}
	}
	return nil, "", fmt.Errorf("%s: %w", rel, fs.ErrNotExist)
}

func readAll(fsys fs.FS, name string) ([]byte, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if st, err := f.Stat(); err == nil && st.IsDir() {
		return nil, fs.ErrNotExist
	}
	b, err := io.ReadAll(io.LimitReader(f, MaxRefSize+1))
	if err != nil {
		return nil, err
	}
	if len(b) > MaxRefSize {
		return nil, fmt.Errorf("%s: larger than %d bytes", name, MaxRefSize)
	}
	return b, nil
}

// findBase returns the first file (in lexical order) with the base name,
// or "".
func findBase(fsys fs.FS, base string) string {
	found := ""
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && d.Name() == base {
			found = p
			return fs.SkipAll
		}
		return nil
	})
	return found
}
