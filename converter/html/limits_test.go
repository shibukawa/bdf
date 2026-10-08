package html

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// remote sets what a document may fetch, and puts it back after the test.
func remote(t *testing.T, images, bytes int, timeout time.Duration) {
	t.Helper()
	i, b, d := maxRemote, maxRemoteBytes, remoteTimeout
	t.Cleanup(func() { maxRemote, maxRemoteBytes, remoteTimeout = i, b, d })
	maxRemote, maxRemoteBytes, remoteTimeout = images, bytes, timeout
}

// imageServer serves a PNG of one pixel, padded to the size the path asks
// for ("/1000/a.png"), and records the paths asked for. Paths with "slow"
// are answered when the server closes, paths with "missing" not at all.
func imageServer(t *testing.T) (srv *httptest.Server, asked func() []string) {
	t.Helper()
	var mu sync.Mutex
	var paths []string
	closing := make(chan struct{})
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		switch {
		case strings.Contains(r.URL.Path, "slow"):
			select {
			case <-closing:
			case <-r.Context().Done():
			}
			return
		case strings.Contains(r.URL.Path, "missing"):
			http.NotFound(w, r)
			return
		}
		size := 0
		fmt.Sscanf(r.URL.Path, "/%d/", &size)
		w.Header().Set("Content-Type", "image/png")
		w.Write(onePixel())
		w.Write(make([]byte, max(size-len(onePixel()), 0)))
	}))
	t.Cleanup(func() {
		close(closing)
		srv.Close()
	})
	return srv, func() []string {
		mu.Lock()
		defer mu.Unlock()
		slices.Sort(paths)
		return slices.Clone(paths)
	}
}

func remoteOptions() *Options {
	opts := testOptions()
	opts.NoRemote = false
	opts.AllowPrivate = true // the test servers listen on the loopback address
	opts.Extract = ExtractNone
	return opts
}

// A document fetches at most maxRemote images, of at most maxRemoteBytes
// together; the images after that show their alternative text, and the
// document warns once.
func TestRemoteLimits(t *testing.T) {
	srv, asked := imageServer(t)
	var doc strings.Builder
	doc.WriteString(`<!DOCTYPE html><p>`)
	for i := range 9 {
		fmt.Fprintf(&doc, `<img src="%s/0/%d.png" alt="[picture %d]"> `, srv.URL, i, i)
	}
	fmt.Fprintf(&doc, `<img src="%s/0/0.png" alt="[again]"><img src="%s/0/missing.png" alt="[missing]"></p>`, srv.URL, srv.URL)
	remote(t, 5, 1<<20, time.Minute)
	res, r := convertHTML(t, doc.String(), remoteOptions())
	c := viewContent(t, r, 0)
	// the first five, of which the first is drawn twice
	if got := asked(); len(got) != 5 || res.Images != 5 || len(c.images) != 6 {
		t.Errorf("asked for %v, %d images, %d drawn", got, res.Images, len(c.images))
	}
	for i := range 9 {
		if got, want := strings.Contains(c.text, fmt.Sprintf("[picture %d]", i)), i >= 5; got != want {
			t.Errorf("the alternative text of picture %d is drawn: %v", i, got)
		}
	}
	if !strings.Contains(c.text, "[missing]") || strings.Contains(c.text, "[again]") {
		t.Errorf("text %q", c.text)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "more than 5 images on the network") {
		t.Errorf("warnings %q", res.Warnings)
	}

	// the bytes: three images of 4000 bytes are less than 14000, four are more
	remote(t, 256, 14000, time.Minute)
	doc.Reset()
	doc.WriteString(`<!DOCTYPE html>`)
	for i := range 6 {
		fmt.Fprintf(&doc, `<p><img src="%s/4000/b%d.png" alt="[picture %d]"></p>`, srv.URL, i, i)
	}
	opts := remoteOptions()
	opts.Fetch = func(u string) ([]byte, error) { return httpGet(privateClient, u, time.Now().Add(time.Minute)) } // one at a time, in the order of the layout
	res, r = convertHTML(t, doc.String(), opts)
	c = viewContent(t, r, 0)
	if res.Images != 3 || len(c.images) != 3 {
		t.Errorf("%d images, %d drawn", res.Images, len(c.images))
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "larger than 0 MiB") {
		t.Errorf("warnings %q", res.Warnings)
	}
}

// The images of a document are fetched in remoteTimeout, all of them: the
// images that take longer, and those after them, are left out.
func TestRemoteDeadline(t *testing.T) {
	srv, asked := imageServer(t)
	remote(t, 256, 1<<20, 300*time.Millisecond)
	var doc strings.Builder
	doc.WriteString(`<!DOCTYPE html><p><img src="` + srv.URL + `/0/fast.png" alt="[fast]">`)
	for i := range 12 {
		fmt.Fprintf(&doc, `<img src="%s/0/slow%d.png" alt="[slow %d]"> `, srv.URL, i, i)
	}
	doc.WriteString(`</p>`)
	start := time.Now()
	res, r := convertHTML(t, doc.String(), remoteOptions())
	if d := time.Since(start); d > 5*time.Second {
		t.Errorf("took %v", d)
	}
	c := viewContent(t, r, 0)
	// the images that were being fetched at the deadline were asked for,
	// the others not
	if got := asked(); res.Images != 1 || len(got) < 2 || len(got) > 1+fetchers || !strings.Contains(c.text, "[slow 11]") || strings.Contains(c.text, "[fast]") {
		t.Errorf("%d images, asked for %v, text %q", res.Images, got, c.text)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "took more than 300ms") {
		t.Errorf("warnings %q", res.Warnings)
	}
}

// Images in elements that are not drawn are not fetched; those of svg
// elements are fetched with the others, before the layout.
func TestRemoteDrawn(t *testing.T) {
	srv, asked := imageServer(t)
	res, _ := convertHTML(t, `<!DOCTYPE html><div hidden><img src="`+srv.URL+`/0/hidden.png"></div><p style="display: none"><img src="`+srv.URL+`/0/none.png"></p>
<img style="visibility:hidden" src="`+srv.URL+`/0/invisible.png"><template><img src="`+srv.URL+`/0/template.png"></template><video><img src="`+srv.URL+`/0/video.png"></video>
<p><img src="`+srv.URL+`/0/drawn.png"></p><svg width="9" height="9"><image href="`+srv.URL+`/0/svg.png" width="9" height="9"/><rect width="1" height="1"/></svg>`, remoteOptions())
	if got := asked(); !slices.Equal(got, []string{"/0/drawn.png", "/0/svg.png"}) || res.Images != 2 || len(res.Warnings) != 0 {
		t.Errorf("asked for %v, %d images, warnings %q", got, res.Images, res.Warnings)
	}
}

// Warnings name the images without the user and the password of their
// address.
func TestRemoteUserInfo(t *testing.T) {
	srv, _ := imageServer(t)
	withUser := strings.Replace(srv.URL, "http://", "http://reader:secret@", 1)
	closed := httptest.NewServer(http.NotFoundHandler())
	gone := strings.Replace(closed.URL, "http://", "http://reader:secret@", 1)
	closed.Close()
	res, _ := convertHTML(t, `<!DOCTYPE html><p><img src="`+withUser+`/0/missing.png" alt="a"> <img src="`+gone+`/gone.png" alt="b"> <img src="`+withUser+`/0/there.png" alt="c"></p>`,
		remoteOptions())
	if res.Images != 1 || len(res.Warnings) != 2 {
		t.Fatalf("%d images, warnings %q", res.Images, res.Warnings)
	}
	for _, w := range res.Warnings {
		if strings.Contains(w, "reader") || strings.Contains(w, "secret") || strings.Contains(w, "@") {
			t.Errorf("warning %q", w)
		}
	}
	if !strings.Contains(res.Warnings[0], srv.URL+"/0/missing.png: HTTP 404") || !strings.Contains(res.Warnings[1], closed.URL+"/gone.png") {
		t.Errorf("warnings %q", res.Warnings)
	}
}

// A document whose elements nest deeper than a parser lets them is
// refused, as XHTML as it is as HTML: the readers of the tree would need a
// stack as deep.
func TestDeepDocument(t *testing.T) {
	for name, head := range map[string]string{"html": `<!DOCTYPE html><html><body>`, "xhtml": `<?xml version="1.0"?><html xmlns="http://www.w3.org/1999/xhtml"><body>`} {
		for _, c := range []struct {
			depth int
			ok    bool
		}{{500, true}, {520, false}} {
			_, err := ConvertBytes([]byte(head+strings.Repeat("<b>", c.depth)+"x"+strings.Repeat("</b>", c.depth)+`</body></html>`), testOptions())
			if (err == nil) != c.ok {
				t.Errorf("%s, %d levels: %v", name, c.depth, err)
			}
		}
	}
}

// The article that Readability picks out holds elements it has renamed,
// headings among them, which it makes divs.
func TestRenamedHeading(t *testing.T) {
	opts := testOptions()
	opts.Extract = ExtractArticle
	res, r := convertHTML(t, `<h1 ><h><h3 >aaaaaaaaaaaaaaaaaaaamath>`, opts)
	if c := viewContent(t, r, 0); !res.Extracted || !strings.Contains(c.text, "aaaaaaaaaaaaaaaaaaaamath>") {
		t.Errorf("extracted %v, text %q", res.Extracted, c.text)
	}
}
