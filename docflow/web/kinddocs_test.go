package web

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/docflow"
)

// docky is a kind a program registers with a chapter of the manual.
func init() {
	docflow.RegisterKind(docflow.Kind{Name: "docky", Role: docflow.RoleExit, Titles: docflow.Text{"en": "Docky", "ja": "ドッキー"},
		Doc:  docflow.Text{"en": "A **docky** node does what the program says.", "ja": "**ドッキー**のノードはプログラムの言うとおりにします。"},
		Send: func(ctx context.Context, o *docflow.Out) (string, error) { return "", nil }})
}

// The chapter of the nodes ends with the chapters of the kinds a
// program registered with a Doc, in the language of the page, each
// with the kind's name as its id (the editor links /help/nodes#<kind>).
func TestKindDocs(t *testing.T) {
	s := newSite(t)
	alice := s.login("alice")
	for lang, want := range map[string][]string{
		"en": {`<h2 id="docky">Docky</h2>`, "<strong>docky</strong> node does what the program says"},
		"ja": {`<h2 id="docky">ドッキー</h2>`, "<strong>ドッキー</strong>のノードはプログラムの言うとおりにします"},
	} {
		code, body := alice.get("/help/" + lang + "/nodes")
		if code != http.StatusOK {
			t.Fatalf("/help/%s/nodes: %d", lang, code)
		}
		for _, w := range want {
			if !strings.Contains(body, w) {
				t.Errorf("/help/%s/nodes lacks %q", lang, w)
			}
		}
		// after the core's, and in the sidebar's sections
		if i, j := strings.Index(body, `id="printer"`), strings.Index(body, `id="docky"`); i < 0 || j < i {
			t.Errorf("/help/%s/nodes: the kind's chapter is not after the core's (%d, %d)", lang, i, j)
		}
		if !strings.Contains(body, `href="/help/`+lang+`/nodes#docky"`) {
			t.Errorf("/help/%s/nodes: no link to the section in the sidebar", lang)
		}
	}
}
