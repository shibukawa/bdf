package docflow

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"
)

// The OpenAPI document describes what the topology publishes, as the
// reader may use it: the entries' bodies (the templates' data), the
// metadata the nodes after them look at, the exports' entries.
func TestOpenAPI(t *testing.T) {
	e := newEnv(t, withViewers)
	store := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "vault.db"))
	secrets := map[string]string{}
	for id, scopes := range map[string][]string{"erp": {ScopePrintTemplate}, "backend": {ScopeExport}, "stranger": {ScopePrint, ScopeExport}} {
		var made struct {
			Secret string `json:"secret"`
		}
		json.Unmarshal(e.admin("PUT", "/api/apps/"+id, App{Scopes: scopes}, http.StatusOK), &made)
		secrets[id] = made.Secret
	}
	// data → a filter → a template node that picks its template by the
	// metadata → an archive → an export; documents → the archive
	e.node("out", Node{Kind: KindExport, By: "customer", Fields: []string{"invoice"}, Apps: []string{"backend"}}, http.StatusOK)
	e.node("vault", Node{Kind: KindArchive, Store: store, To: []string{"out"}}, http.StatusOK)
	e.node("forms", Node{Kind: KindTemplate, TemplateFrom: "doc_type", To: []string{"vault"}}, http.StatusOK)
	tpl, err := os.ReadFile("../testdata/form/pages.bdft")
	if err != nil {
		t.Fatal(err)
	}
	e.admin("PUT", "/api/nodes/forms/templates/slip?name=pages.bdft", tpl, http.StatusOK)
	e.admin("PUT", "/api/nodes/forms/templates/label?name=pages.bdft", tpl, http.StatusOK)
	e.node("pick", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Meta: map[string]string{"region": "^jp", "$sender": "erp"}}, To: []string{"forms"}}}, Otherwise: []string{"forms"}}, http.StatusOK)
	e.node("data-in", Node{Kind: KindEndpoint, Title: "Invoices", Description: "The ERP's invoices.", Accept: AcceptData, To: []string{"pick"},
		Apps: []string{"erp", "stranger"}, Meta: map[string]any{"region": "jp-east"}, Port: 18080}, http.StatusOK)
	e.node("docs-in", Node{Kind: KindEndpoint, Formats: []string{"pdf", "docx", "image", "json"}, To: []string{"vault"}}, http.StatusOK)
	e.node("seen", Node{Kind: KindSink}, http.StatusOK)

	get := func(do func(*http.Request), query string, status int) map[string]any {
		t.Helper()
		res, out := e.doWith(do, "GET", OpenAPIPath+query, nil)
		if res.StatusCode != status {
			t.Fatalf("GET %s: %d %s; want %d", OpenAPIPath+query, res.StatusCode, out, status)
		}
		var doc map[string]any
		json.Unmarshal(out, &doc)
		return doc
	}
	bearer := func(tok string) func(*http.Request) {
		return func(r *http.Request) { r.Header.Set("Authorization", "Bearer "+tok) }
	}
	admin := bearer(adminToken)
	doc := get(admin, "", http.StatusOK)
	if doc["openapi"] != "3.1.0" || dig(doc, "servers", 0, "url") != e.srv.URL {
		t.Fatalf("document: %v %v", doc["openapi"], doc["servers"])
	}
	want := []string{"/api/jobs/{id}", "/api/jobs/{id}/cancel", "/api/jobs/{id}/document", "/in/data-in", "/in/docs-in",
		"/out/out", "/out/out/{job}", "/out/out/{job}/document", "/out/out/{job}/record"}
	if got := keysOf(doc["paths"]); !slices.Equal(got, want) {
		t.Fatalf("paths %v; want %v", got, want)
	}
	if got := dig(doc, "components", "securitySchemes", "app", "flows", "clientCredentials", "tokenUrl"); got != e.srv.URL+OAuthTokenPath {
		t.Fatalf("token URL %v", got)
	}

	// the data an entry takes is any of the templates its data reaches
	in := dig(doc, "paths", "/in/data-in", "post").(map[string]any)
	if in["summary"] != "Send to Invoices" || in["operationId"] != "send_data_in" || !strings.Contains(in["description"].(string), "The ERP's invoices.") {
		t.Fatalf("operation %v %v %v", in["summary"], in["operationId"], in["description"])
	}
	anyOf := dig(in, "requestBody", "content", "application/json", "schema", "anyOf")
	if !reflect.DeepEqual(anyOf, []any{
		map[string]any{"$ref": "#/components/schemas/TemplateData.forms.label"},
		map[string]any{"$ref": "#/components/schemas/TemplateData.forms.slip"},
	}) {
		t.Fatalf("data schema %v", anyOf)
	}
	if content := keysOf(dig(in, "requestBody", "content")); !slices.Equal(content, []string{"application/json", "application/vnd.apache.parquet",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/yaml", "text/csv"}) {
		t.Fatalf("data content %v", content)
	}
	slip := dig(doc, "components", "schemas", "TemplateData.forms.slip").(map[string]any)
	if slip["type"] != "object" || dig(slip, "properties", "pages") == nil || slip["$schema"] != nil || !strings.Contains(slip["description"].(string), "?meta.doc_type=slip") {
		t.Fatalf("template schema %v", slip)
	}
	// the metadata the nodes after it look at, but what the entry sets
	// itself and the manager's keys
	params := map[string]any{}
	for _, p := range in["parameters"].([]any) {
		params[p.(map[string]any)["name"].(string)] = p
	}
	if got := keysOf(params); !slices.Equal(got, []string{"meta.customer", "meta.doc_type", "meta.invoice", "name", "user"}) {
		t.Fatalf("parameters %v", got)
	}
	if enum := dig(params, "meta.doc_type", "schema", "enum"); !reflect.DeepEqual(enum, []any{"label", "slip"}) {
		t.Fatalf("doc_type %v", enum)
	}
	if desc := dig(params, "meta.customer", "description").(string); !strings.Contains(desc, "whose the document is at the export out") {
		t.Fatalf("customer: %s", desc)
	}
	// its own port takes it too; a person's token does not send
	if got := dig(doc, "paths", "/in/data-in", "servers", 1, "url"); got != "http://127.0.0.1:18080" {
		t.Fatalf("port %v", got)
	}
	if got := dig(in, "security"); !reflect.DeepEqual(got, []any{map[string]any{"key": []any{}}, map[string]any{"admin": []any{}}, map[string]any{"app": []any{ScopePrintTemplate}}}) {
		t.Fatalf("security %v", got)
	}
	// documents: the media types of the formats it takes
	docs := dig(doc, "paths", "/in/docs-in", "post", "requestBody", "content")
	if got := keysOf(docs); !slices.Equal(got, []string{"application/json", "application/pdf", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "image/*"}) {
		t.Fatalf("document content %v", got)
	}
	// an export's entries carry its fields
	entry := dig(doc, "paths", "/out/out", "get", "responses", "200", "content", "application/json", "schema", "items", "$ref")
	if entry != "#/components/schemas/ExportEntry.out" || dig(doc, "components", "schemas", "ExportEntry.out", "allOf", 1, "properties", "fields", "properties", "invoice") == nil {
		t.Fatalf("entry %v", entry)
	}
	// the schemas of the answers are their Go types'
	job := dig(doc, "components", "schemas", "Job").(map[string]any)
	if dig(job, "properties", "state", "enum") == nil || dig(job, "properties", "received", "format") != "date-time" ||
		!reflect.DeepEqual(dig(job, "properties", "legs", "type"), []any{"array", "null"}) || slices.Contains(job["required"].([]any), "meta") {
		t.Fatalf("job %v", job)
	}
	for _, name := range []string{"Payload", "Leg", "Principal", "Mark", "ArchiveRecord", "ExportEntry", "Trace", "Fingerprint", "Error"} {
		if dig(doc, "components", "schemas", name) == nil {
			t.Fatalf("no schema %s", name)
		}
	}
	// every reference is to a schema there
	var refs func(v any)
	refs = func(v any) {
		switch x := v.(type) {
		case map[string]any:
			if r, ok := x["$ref"].(string); ok && dig(doc, "components", "schemas", strings.TrimPrefix(r, "#/components/schemas/")) == nil {
				t.Errorf("dangling %s", r)
			}
			for _, y := range x {
				refs(y)
			}
		case []any:
			for _, y := range x {
				refs(y)
			}
		}
	}
	refs(doc)

	// the version follows what is published
	if again := get(admin, "", http.StatusOK); again["info"].(map[string]any)["version"] != doc["info"].(map[string]any)["version"] {
		t.Fatal("the version changed with nothing published changed")
	}
	e.node("seen", Node{Kind: KindSink, Latest: 3}, http.StatusOK)
	if again := get(admin, "", http.StatusOK); again["info"].(map[string]any)["version"] != doc["info"].(map[string]any)["version"] {
		t.Fatal("the version changed with a sink")
	}
	e.node("docs-in", Node{Kind: KindEndpoint, Formats: []string{"pdf"}, To: []string{"vault"}}, http.StatusOK)
	if again := get(admin, "", http.StatusOK); again["info"].(map[string]any)["version"] == doc["info"].(map[string]any)["version"] {
		t.Fatal("the version stayed with an endpoint changed")
	}

	// ?node= narrows it; a node that publishes nothing is not found
	if got := keysOf(get(admin, "?node=docs-in", http.StatusOK)["paths"]); !slices.Equal(got, []string{"/api/jobs/{id}", "/api/jobs/{id}/cancel", "/api/jobs/{id}/document", "/in/docs-in"}) {
		t.Fatalf("docs-in %v", got)
	}
	if got := keysOf(get(admin, "?node=out", http.StatusOK)["paths"]); !slices.Equal(got, []string{"/out/out", "/out/out/{job}", "/out/out/{job}/document", "/out/out/{job}/record"}) {
		t.Fatalf("out %v", got)
	}
	get(admin, "?node=seen", http.StatusNotFound)
	get(admin, "?node=nobody", http.StatusNotFound)

	// an app reads what it may use: the endpoints and the exports that
	// list it, with the scopes; a person the exports; nobody nothing
	token := func(id string) string {
		res, out := e.doWith(func(r *http.Request) {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.SetBasicAuth(id, secrets[id])
		}, "POST", OAuthTokenPath, []byte("grant_type=client_credentials"))
		var body struct {
			Token string `json:"access_token"`
		}
		if json.Unmarshal(out, &body); res.StatusCode != http.StatusOK || body.Token == "" {
			t.Fatalf("token of %s: %d %s", id, res.StatusCode, out)
		}
		return body.Token
	}
	if got := keysOf(get(bearer(token("erp")), "", http.StatusOK)["paths"]); !slices.Equal(got, []string{"/api/jobs/{id}", "/api/jobs/{id}/cancel", "/api/jobs/{id}/document", "/in/data-in"}) {
		t.Fatalf("erp %v", got)
	}
	get(bearer(token("erp")), "?node=docs-in", http.StatusNotFound)
	if got := keysOf(get(bearer(token("backend")), "", http.StatusOK)["paths"]); !slices.Equal(got, []string{"/out/out", "/out/out/{job}", "/out/out/{job}/document", "/out/out/{job}/record"}) {
		t.Fatalf("backend %v", got)
	}
	// listed with the wrong scope (print, not print.template), or not
	// listed (the export)
	if got := keysOf(get(bearer(token("stranger")), "", http.StatusOK)["paths"]); len(got) != 0 {
		t.Fatalf("stranger %v", got)
	}
	alice := func(r *http.Request) { r.Header.Set("X-Test-Viewer", "alice") }
	if got := keysOf(get(alice, "", http.StatusOK)["paths"]); !slices.Equal(got, []string{"/out/out", "/out/out/{job}", "/out/out/{job}/document", "/out/out/{job}/record"}) {
		t.Fatalf("alice %v", got)
	}
	get(func(*http.Request) {}, "", http.StatusUnauthorized)
}

// The metadata an entry's job needs at the exits: the {{meta.<key>}} of
// their settings, but the manager's keys.
func TestOpenAPIMetaOfExits(t *testing.T) {
	d := &apiDoc{nodes: map[string]*Node{
		"in":   {Name: "in", Kind: KindEndpoint, To: []string{"hook"}, OnError: []string{"tell"}},
		"hook": {Name: "hook", Kind: "http", URL: "https://example.com/{{meta.customer}}?job={{ meta.$job }}", Headers: map[string]string{"X-Order": "{{meta.order}}"}},
		"tell": {Name: "tell", Kind: "notify", Text: "{{meta.$error}} for {{meta.owner}}", Recipients: []string{"{{meta.contact}}"}},
	}}
	var got []string
	for _, p := range d.metaParams(d.nodes["in"]) {
		got = append(got, p.(map[string]any)["name"].(string))
	}
	if !slices.Equal(got, []string{"meta.contact", "meta.customer", "meta.order", "meta.owner"}) {
		t.Fatalf("parameters %v", got)
	}
}

// dig walks a JSON value by keys and indexes; nil when it is not there.
func dig(v any, path ...any) any {
	for _, p := range path {
		switch k := p.(type) {
		case string:
			m, ok := v.(map[string]any)
			if !ok {
				return nil
			}
			v = m[k]
		case int:
			a, ok := v.([]any)
			if !ok || k >= len(a) {
				return nil
			}
			v = a[k]
		}
	}
	return v
}

// keysOf are the keys of a JSON object, sorted.
func keysOf(v any) []string {
	m, _ := v.(map[string]any)
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
