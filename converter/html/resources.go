package html

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// maxImage limits the size of an image that is read or fetched.
const maxImage = 50 << 20

// fetchers is how many images are fetched at the same time.
const fetchers = 8

// resources resolves the references of a document: its images and links.
type resources struct {
	dir    string
	base   *url.URL
	parts  map[string][]byte // parts of an MHTML archive by Content-Location and cid: URL
	remote bool
	fetch  func(string) ([]byte, error)

	mu       sync.Mutex
	fetched  map[string]fetched // remote images by URL
	warnings []string
}

type fetched struct {
	data []byte
	err  error
}

func newResources(opts *Options) *resources {
	r := &resources{dir: opts.Dir, remote: !opts.NoRemote, fetch: opts.Fetch, fetched: map[string]fetched{}}
	if r.fetch == nil {
		r.fetch = httpGet
	}
	if opts.BaseURL != "" {
		r.setBase(opts.BaseURL)
	}
	return r
}

// setBase makes an absolute URL (or one relative to the base in effect)
// the base of the references.
func (r *resources) setBase(s string) {
	u, err := resolve(r.base, s)
	if err == nil && u.IsAbs() {
		r.base = u
	}
}

// link resolves the target of a link to an absolute http:, https: or
// mailto: URL; "" when it has none (a relative link without a base URL).
func (r *resources) link(href string) string {
	u, err := resolve(r.base, href)
	if err != nil {
		return ""
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https", "mailto":
		return u.String()
	}
	return ""
}

// image returns the bytes of an image by its src.
func (r *resources) image(src string) ([]byte, error) {
	src = strings.TrimSpace(src)
	if strings.HasPrefix(strings.ToLower(src), "data:") {
		return dataURL(src)
	}
	if b, ok := r.parts[src]; ok {
		return b, nil
	}
	u, err := resolve(r.base, src)
	if err != nil {
		return nil, err
	}
	if b, ok := r.parts[u.String()]; ok {
		return b, nil
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		if !r.remote {
			return nil, errors.New("images on the network are not fetched")
		}
		return r.get(u.String())
	case "":
		return r.local(u)
	}
	return nil, fmt.Errorf("unsupported URL scheme %q", u.Scheme)
}

// local reads a file that a relative reference names, in the document's
// directory (".." may leave it; absolute paths are not read).
func (r *resources) local(u *url.URL) ([]byte, error) {
	if r.dir == "" {
		return nil, errors.New("relative reference with no directory to resolve it in")
	}
	if u.Host != "" || strings.HasPrefix(u.Path, "/") || filepath.IsAbs(u.Path) {
		return nil, errors.New("absolute paths are not read")
	}
	f, err := os.Open(filepath.Join(r.dir, filepath.FromSlash(u.Path)))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLimited(f)
}

func readLimited(rd io.Reader) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(rd, maxImage+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxImage {
		return nil, fmt.Errorf("larger than %d MiB", maxImage>>20)
	}
	return b, nil
}

// get fetches a remote image once.
func (r *resources) get(u string) ([]byte, error) {
	r.mu.Lock()
	f, ok := r.fetched[u]
	r.mu.Unlock()
	if ok {
		return f.data, f.err
	}
	b, err := r.fetch(u)
	r.mu.Lock()
	r.fetched[u] = fetched{b, err}
	r.mu.Unlock()
	return b, err
}

// prefetch fetches the remote images of a document, several at a time.
func (r *resources) prefetch(doc *xhtml.Node) {
	if !r.remote {
		return
	}
	seen := map[string]bool{}
	var urls []string
	var walk func(n *xhtml.Node)
	walk = func(n *xhtml.Node) {
		if n.Type == xhtml.ElementNode && n.DataAtom == atom.Img {
			src := attr(n, "src")
			if src == "" {
				if f := strings.Fields(attr(n, "srcset")); len(f) > 0 {
					src = strings.TrimSuffix(f[0], ",")
				}
			}
			if u, err := resolve(r.base, src); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
				if s := u.String(); !seen[s] && r.parts[s] == nil {
					seen[s] = true
					urls = append(urls, s)
				}
			}
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			walk(k)
		}
	}
	walk(doc)
	jobs := make(chan string)
	var wg sync.WaitGroup
	for range min(fetchers, len(urls)) {
		wg.Go(func() {
			for u := range jobs {
				r.get(u)
			}
		})
	}
	for _, u := range urls {
		jobs <- u
	}
	close(jobs)
	wg.Wait()
}

var client = &http.Client{Timeout: 30 * time.Second}

// httpGet fetches an image over HTTP.
func httpGet(u string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "bdf-html/0.1 (+https://github.com/shibukawa/bdf)")
	req.Header.Set("Accept", "image/avif,image/webp,image/png,image/jpeg,image/gif,image/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s", resp.Status)
	}
	return readLimited(resp.Body)
}

// dataURL decodes a data: URL.
func dataURL(s string) ([]byte, error) {
	_, rest, _ := strings.Cut(s, ":")
	meta, data, ok := strings.Cut(rest, ",")
	if !ok {
		return nil, errors.New("malformed data: URL")
	}
	if strings.HasSuffix(strings.ToLower(meta), ";base64") {
		data = strings.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, data)
		if b, err := base64.StdEncoding.DecodeString(data); err == nil {
			return b, nil
		}
		return base64.RawStdEncoding.DecodeString(strings.TrimRight(data, "="))
	}
	d, err := url.PathUnescape(data)
	return []byte(d), err
}
