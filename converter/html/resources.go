package html

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf/converter/internal/localfile"
	"github.com/shibukawa/bdf/converter/internal/webdoc"
	"github.com/shibukawa/bdf/converter/internal/wordproc"
	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// maxImage limits the size of an image that is read or fetched.
const maxImage = 50 << 20

// fetchers is how many images are fetched at the same time.
const fetchers = 8

// What a document may fetch from the network: so many images, of so many
// bytes together, in so much time. The images after that are left out. A
// document can name any number of images, on servers that answer slowly.
var (
	maxRemote      = 256
	maxRemoteBytes = 256 << 20
	remoteTimeout  = 120 * time.Second
)

// resources resolves the references of a document: its images and links.
type resources struct {
	dir    string
	base   *url.URL
	parts  map[string][]byte // parts of an MHTML archive by Content-Location and cid: URL
	remote bool
	fetch  func(string) ([]byte, error) // nil: httpGet
	client *http.Client                 // of httpGet: public addresses only, unless the document may reach private ones

	mu       sync.Mutex
	fetched  map[string]fetched // remote images by URL
	bytes    int                // of the images fetched
	deadline time.Time          // of the fetches, from the first of them
	warnings []string
}

type fetched struct {
	data []byte
	err  error
}

func newResources(opts *Options) *resources {
	r := &resources{dir: opts.Dir, remote: !opts.NoRemote, fetch: opts.Fetch, fetched: map[string]fetched{}, client: client}
	if opts.AllowPrivate {
		r.client = privateClient
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
	if webdoc.IsDataURL(src) {
		return webdoc.DataURL(src)
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
// directory. Parent references and symlinks may not leave that directory.
func (r *resources) local(u *url.URL) ([]byte, error) {
	if r.dir == "" {
		return nil, errors.New("relative reference with no directory to resolve it in")
	}
	if u.Host != "" || strings.HasPrefix(u.Path, "/") || filepath.IsAbs(u.Path) {
		return nil, errors.New("absolute paths are not read")
	}
	f, err := localfile.Open(r.dir, filepath.FromSlash(u.Path))
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

// leftOut is the error of an image that is not fetched, or not kept,
// because the document has fetched what it may: the layout warns of the
// limit once.
func leftOut(format string, args ...any) error {
	return &wordproc.LimitError{Warning: fmt.Sprintf(format, args...) + "; the images after that are left out"}
}

// get fetches a remote image once, within what the document may fetch.
func (r *resources) get(u string) ([]byte, error) {
	r.mu.Lock()
	f, ok := r.fetched[u]
	var err error
	switch {
	case ok:
	case len(r.fetched) >= maxRemote:
		err = leftOut("the document has more than %d images on the network", maxRemote)
	case r.deadline.IsZero():
		r.deadline = time.Now().Add(remoteTimeout)
	case !time.Now().Before(r.deadline):
		err = r.late()
	}
	if !ok && err == nil {
		r.fetched[u] = fetched{err: errors.New("not fetched yet")} // counts as one of the images
	}
	deadline := r.deadline
	r.mu.Unlock()
	if ok || err != nil {
		return f.data, cmp.Or(err, f.err)
	}
	var b []byte
	if r.fetch != nil {
		b, err = r.fetch(u)
	} else {
		b, err = httpGet(r.client, u, deadline)
		if errors.Is(err, context.DeadlineExceeded) && !time.Now().Before(deadline) {
			err = r.late()
		}
	}
	r.mu.Lock()
	if err == nil && r.bytes+len(b) > maxRemoteBytes {
		b, err = nil, leftOut("the images on the network are larger than %d MiB", maxRemoteBytes>>20)
	}
	r.bytes += len(b)
	r.fetched[u] = fetched{b, err}
	r.mu.Unlock()
	return b, err
}

// late is the error of an image that is not fetched in the time the
// document has for its images.
func (r *resources) late() error {
	return leftOut("the images on the network took more than %v", remoteTimeout)
}

// prefetch fetches the remote images of a document, several at a time:
// those of the elements that are drawn, as many as the document may fetch.
func (r *resources) prefetch(doc *xhtml.Node) {
	if !r.remote {
		return
	}
	seen := map[string]bool{}
	var urls []string
	var walk func(n *xhtml.Node)
	walk = func(n *xhtml.Node) {
		if n.Type == xhtml.ElementNode {
			if !wordproc.Drawn(n) {
				return
			}
			src := ""
			switch {
			case n.DataAtom == atom.Img && n.Namespace == "":
				if src = attr(n, "src"); src == "" {
					if f := strings.Fields(attr(n, "srcset")); len(f) > 0 {
						src = strings.TrimSuffix(f[0], ",")
					}
				}
			case n.Namespace == "svg" && (n.Data == "image" || n.Data == "feImage"):
				// written into the SVG document of the svg element
				for _, a := range n.Attr {
					if a.Key == "href" && (a.Namespace == "" || a.Namespace == "xlink") && (src == "" || a.Namespace == "") {
						src = strings.TrimSpace(a.Val)
					}
				}
			}
			if u, err := resolve(r.base, src); src != "" && err == nil && (u.Scheme == "http" || u.Scheme == "https") {
				if s := u.String(); !seen[s] && r.parts[s] == nil && len(urls) < maxRemote {
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

// client fetches the images of documents from public addresses only;
// privateClient from any address (Options.AllowPrivate). See newClient.
var client, privateClient = newClient(false), newClient(true)

// httpGet fetches an image over HTTP, until the deadline at the latest.
func httpGet(client *http.Client, u string, deadline time.Time) ([]byte, error) {
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "bdf-html/0.1 (+https://github.com/shibukawa/bdf)")
	req.Header.Set("Accept", "image/avif,image/webp,image/png,image/jpeg,image/gif,image/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		// the error names the address: without its user and password
		var ue *url.Error
		if errors.As(err, &ue) {
			if a, perr := url.Parse(ue.URL); perr == nil && a.User != nil {
				a.User = nil
				ue.URL = a.String()
			}
		}
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s", resp.Status)
	}
	return readLimited(resp.Body)
}
