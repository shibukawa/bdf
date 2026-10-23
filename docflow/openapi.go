package docflow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/form"
)

// The topology publishes endpoints, under the API by each node's kind
// (NodePath): every keyed entry takes jobs at POST /api/<kind>/<name> (an
// uploader's /api/upload/<name>), every inbox is worked at
// /api/inbox/<name>, and every downloader is read at
// /api/download/<name>. OpenAPI describes them as an OpenAPI 3.1
// document, made of the topology as it is, so that an application's
// developer generates a client of the endpoints they use, or tries them.
// What each takes is the node's: the formats an entry accepts, as media
// types; of an entry of data, the JSON Schema of the templates its data
// reaches (form.Schema), and the metadata keys the nodes after it look at
// (a filter's rules, an assert's checks, a branch's cases, a template
// node's template_from, a downloader's by and fields, the {{meta.<key>}}
// of an exit's settings), as ?meta.<key>= parameters, those the asserts
// right after it check with their constraints (and the refusal, 422); of
// an inbox, the nodes it releases to; of a downloader, the fields its
// entries carry. The jobs an app sent are followed at /api/jobs/<id>. The
// document's version is a hash of what it describes: it changes when the
// topology changes what it publishes.

// OpenAPIPath is where the manager serves the OpenAPI document.
const OpenAPIPath = "/api/openapi.json"

// OpenAPIOptions are what OpenAPI describes.
type OpenAPIOptions struct {
	// Server is the URL the manager is reached at ("https://docflow.example"),
	// the document's server; "" leaves the paths relative.
	Server string
	// Viewer is who reads it: an app is shown the entries, the inboxes
	// and the downloaders that list it, with the scopes to use them; a
	// person the downloaders and the inboxes they work; nil, an
	// administrator or a person with the role topology, every one.
	Viewer *Viewer
	// Nodes narrows it to these nodes; none, every one the viewer is
	// shown. A node named that is not published to the viewer is
	// ErrNotFound.
	Nodes []string
}

// OpenAPI returns the OpenAPI 3.1 document of the endpoints the topology
// publishes (JSON).
func (m *Manager) OpenAPI(ctx context.Context, o OpenAPIOptions) ([]byte, error) {
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	tpls, err := list[templateRecord](ctx, m.st.db)
	if err != nil {
		return nil, err
	}
	d := &apiDoc{m: m, o: o, nodes: map[string]*Node{}, tpls: map[string][]Template{}, schemas: map[string]any{}, paths: map[string]any{}, ids: map[string]bool{}}
	for i := range nodes {
		d.nodes[nodes[i].Name] = &nodes[i]
	}
	for _, t := range tpls {
		d.tpls[t.Data.Node] = append(d.tpls[t.Data.Node], Template(t.Data))
	}
	for _, ts := range d.tpls {
		slices.SortFunc(ts, func(a, b Template) int { return strings.Compare(a.ID, b.ID) })
	}
	want := map[string]bool{}
	for _, name := range o.Nodes {
		want[name] = true
	}
	shown := map[string]bool{}
	sends, inboxes := false, false
	for _, n := range nodes {
		if len(want) > 0 && !want[n.Name] || !d.shows(&n) {
			continue
		}
		switch {
		case keyed(n.Kind):
			d.entry(&n)
			sends = true
		case n.Kind == KindInbox:
			d.inbox(&n)
			inboxes = true
		case n.Kind == KindDownload:
			d.download(&n)
		default:
			continue
		}
		shown[n.Name] = true
	}
	for name := range want {
		if !shown[name] {
			return nil, fmt.Errorf("%w: %s publishes no endpoint here", ErrNotFound, name)
		}
	}
	if sends {
		d.jobs()
	}
	if inboxes {
		d.inboxes()
	}
	return d.write()
}

// apiDoc is an OpenAPI document being made.
type apiDoc struct {
	m       *Manager
	o       OpenAPIOptions
	nodes   map[string]*Node
	tpls    map[string][]Template
	schemas map[string]any
	paths   map[string]any
	tags    []any
	ids     map[string]bool
}

// shows tells whether the viewer is shown a node's endpoints: an app
// those of the uploaders it may send to, the inboxes it works and the
// downloaders it may read, a person the inboxes they work and the
// downloaders.
func (d *apiDoc) shows(n *Node) bool {
	v := d.o.Viewer
	if v == nil || v.Can(RoleTopology) {
		return true
	}
	switch {
	case n.Kind == KindUpload && v.Kind == PrincipalApp:
		need := ScopePrint
		if n.Accept == AcceptData {
			need = ScopePrintTemplate
		}
		return slices.Contains(n.Apps, v.ID) && v.Has(need)
	case n.Kind == KindDownload && v.Kind == PrincipalApp:
		return slices.Contains(n.Apps, v.ID) && v.Has(ScopeDownload)
	case n.Kind == KindDownload:
		return v.Kind == PrincipalUser
	case n.Kind == KindInbox:
		return v.mayWork(n) == nil
	}
	return false
}

// url is a path of the manager's as a URL of the document's server.
func (d *apiDoc) url(path string) string { return strings.TrimSuffix(d.o.Server, "/") + path }

// opID makes an operation's id of words: letters and digits, the rest an
// underscore, unique in the document.
func (d *apiDoc) opID(words ...string) string {
	id := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			return r
		}
		return '_'
	}, strings.Join(words, "_"))
	for i, base := 2, id; d.ids[id]; i++ {
		id = base + "_" + strconv.Itoa(i)
	}
	d.ids[id] = true
	return id
}

// tag names a node's operations, with what the node is for.
func (d *apiDoc) tag(n *Node, what string) {
	desc := what
	if n.Title != "" {
		desc = n.Title + ": " + desc
	}
	if n.Description != "" {
		desc += "\n\n" + n.Description
	}
	d.tags = append(d.tags, map[string]any{"name": n.Name, "description": desc})
}

func ref(name string) map[string]any { return map[string]any{"$ref": "#/components/schemas/" + name} }

// jsonBody is a response of JSON of a schema.
func jsonBody(desc string, schema any) map[string]any {
	return map[string]any{"description": desc, "content": map[string]any{"application/json": map[string]any{"schema": schema}}}
}

// errorResponses are the responses of errors, by status.
func (d *apiDoc) errorResponses(out map[string]any, statuses ...int) {
	why := map[int]string{
		http.StatusBadRequest:            "a request the manager cannot take",
		http.StatusUnauthorized:          "no credentials, or credentials not known",
		http.StatusForbidden:             "credentials that may not do this",
		http.StatusNotFound:              "not there, or not the reader's to read",
		http.StatusConflict:              "what is no longer in a state to do it (an item released or discarded already)",
		http.StatusRequestEntityTooLarge: "a document larger than the manager takes",
		http.StatusUnsupportedMediaType:  "a format the entry does not take",
	}
	if slices.Contains(statuses, http.StatusUnprocessableEntity) {
		out["422"] = jsonBody("metadata an assert right after the entry refused: the job is kept, aborted", ref(d.goSchema(reflect.TypeFor[apiRefusal]())))
		statuses = slices.DeleteFunc(slices.Clone(statuses), func(s int) bool { return s == http.StatusUnprocessableEntity })
	}
	d.goSchema(reflect.TypeFor[apiError]())
	for _, s := range statuses {
		out[strconv.Itoa(s)] = jsonBody(why[s], ref("Error"))
	}
}

// apiError is the body of an error.
type apiError struct {
	Error string `json:"error"`
}

// The entries.

// entry describes POST /api/<kind>/<name> of a keyed entry.
func (d *apiDoc) entry(n *Node) {
	k := kindNamed(n.Kind)
	d.tag(n, "an entry: POST a document or data to it, and it flows on through the topology")
	var desc []string
	if n.Description != "" {
		desc = append(desc, n.Description)
	}
	switch {
	case k != nil && k.Fetch != nil:
		desc = append(desc, "The body names the document, which the "+k.title()+" reads.")
	case n.Accept == AcceptData:
		desc = append(desc, "The body is data, which goes on as it is: a template node fills a template with it.")
	case n.Output == OutputFile:
		desc = append(desc, "The body is a document, which goes on as it was received.")
	default:
		desc = append(desc, "The body is a document, which goes on as the .bdf document made of it.")
	}
	if len(n.To) > 0 {
		desc = append(desc, "It goes on to "+d.names(n.To)+".")
	}
	desc = append(desc, "The answer is the job, which an app follows at /api/jobs/{id}. The sender's metadata is ?meta.<key>= parameters, or X-Meta-<Key> headers (the key in lower case), or in a multipart/form-data body a JSON object in the part meta, beside the document in the part file; the parameters below are the keys the nodes after this one look at, and any other key goes along too.")
	gates := d.gates(n)
	if len(gates) > 0 {
		desc = append(desc, "The asserts right after it ("+d.names(gates)+") check the metadata first: a job that fails a check is refused (422) with what failed, and kept, aborted.")
	}
	security := []any{map[string]any{"key": []string{}}, map[string]any{"admin": []string{}}}
	if n.Kind == KindUpload {
		scope := ScopePrint
		if n.Accept == AcceptData {
			scope = ScopePrintTemplate
		}
		security = append(security, map[string]any{"app": []string{scope}})
	}
	responses := map[string]any{"201": jsonBody("the job", ref(d.goSchema(reflect.TypeFor[Job]())))}
	statuses := []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusRequestEntityTooLarge, http.StatusUnsupportedMediaType}
	if len(gates) > 0 {
		statuses = append(statuses, http.StatusUnprocessableEntity)
	}
	d.errorResponses(responses, statuses...)
	params := []any{
		map[string]any{"name": "name", "in": "query", "description": "The job's name: a file's name, whose extension tells the format.", "schema": map[string]any{"type": "string"}},
		map[string]any{"name": "user", "in": "query", "description": "Who the sender says sent it: shown, never trusted.", "schema": map[string]any{"type": "string"}},
	}
	params = append(params, d.metaParams(n)...)
	op := map[string]any{
		"operationId": d.opID("send", n.Name),
		"tags":        []string{n.Name},
		"summary":     "Send to " + or(n.Title, n.Name),
		"description": strings.Join(desc, "\n\n"),
		"parameters":  params,
		"requestBody": map[string]any{"required": true, "content": d.entryBody(n, k)},
		"responses":   responses,
		"security":    security,
	}
	d.paths[NodePath(n.Kind, n.Name)] = map[string]any{"post": op}
}

// names lists nodes: a, b and c (kind).
func (d *apiDoc) names(names []string) string {
	var out []string
	for _, name := range names {
		if x := d.nodes[name]; x != nil {
			name += " (" + x.Kind + ")"
		}
		out = append(out, name)
	}
	if len(out) > 1 {
		return strings.Join(out[:len(out)-1], ", ") + " and " + out[len(out)-1]
	}
	return strings.Join(out, "")
}

// mediaTypes are the media types of the formats contentTypeOf does not
// name; what neither names is application/octet-stream.
var mediaTypes = map[string]string{
	"audio":      "audio/*",
	"cgm":        "image/cgm",
	"docx":       "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"drawio":     "application/vnd.jgraph.mxfile",
	"dxf":        "image/vnd.dxf",
	"emf":        "image/emf",
	"epub":       "application/epub+zip",
	"font":       "font/*",
	"html":       "text/html",
	"idml":       "application/vnd.adobe.indesign-idml-package",
	"image":      "image/*",
	"markdown":   "text/markdown",
	"midi":       "audio/midi",
	"musicxml":   "application/vnd.recordare.musicxml+xml",
	"pptx":       "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"psd":        "image/vnd.adobe.photoshop",
	"pwg-raster": "image/pwg-raster",
	"tiff":       "image/tiff",
	"urf":        "image/urf",
}

func mediaTypeOf(format string) string {
	if t := contentTypeOf(format, ""); t != "application/octet-stream" {
		return t
	}
	return or(mediaTypes[format], "application/octet-stream")
}

// entryBody is the content an entry takes: of a reference, its schema;
// of data, the schema of the templates it reaches; of documents, the
// media types of its formats. Either may come in multipart/form-data
// too (withMultipart).
func (d *apiDoc) entryBody(n *Node, k *Kind) map[string]any {
	binary := map[string]any{"type": "string", "format": "binary"}
	if k != nil && k.Fetch != nil {
		s := k.FetchSchema
		if s == nil {
			s = map[string]any{"type": "object"}
		}
		return map[string]any{"application/json": map[string]any{"schema": s}}
	}
	formats := n.Formats
	if n.Accept == AcceptData && len(formats) == 0 {
		formats = dataFormats
	}
	if len(formats) == 0 {
		return withMultipart(map[string]any{"application/octet-stream": map[string]any{"schema": map[string]any{
			"type": "string", "format": "binary", "description": "A document of any format the converter reads (" + strings.Join(KnownFormats(), ", ") + "), told by its content and ?name=.",
		}}})
	}
	data := d.dataSchema(n)
	byType := map[string][]string{}
	for _, f := range formats {
		t := mediaTypeOf(f)
		byType[t] = append(byType[t], f)
	}
	out := map[string]any{}
	for t, fs := range byType {
		s := binary
		switch {
		case t == "application/json" || t == "application/yaml":
			s = data
		case strings.HasPrefix(t, "text/"):
			s = map[string]any{"type": "string"}
		}
		s = maps.Clone(s)
		var what []string
		for _, f := range fs {
			if c := converter.Lookup(f); c != nil && c.Description != "" {
				f += " (" + c.Description + ")"
			}
			what = append(what, f)
		}
		if _, ok := s["$ref"]; !ok {
			s["description"] = strings.Join(what, ", ")
		}
		out[t] = map[string]any{"schema": s}
	}
	return withMultipart(out)
}

// withMultipart adds to the media types of a body the same in
// multipart/form-data (readBody): the document, of one of those types, in
// the part file, after the metadata as a JSON object in the part meta.
func withMultipart(content map[string]any) map[string]any {
	types := slices.Sorted(maps.Keys(content))
	content["multipart/form-data"] = map[string]any{
		"schema": map[string]any{
			"type":     "object",
			"required": []string{"file"},
			"properties": map[string]any{
				"meta": map[string]any{"type": "object", "description": "The sender's metadata, over the X-Meta-* headers and ?meta.<key>= parameters."},
				"file": map[string]any{"type": "string", "format": "binary", "description": "The document, as the body of the other media types; its name is the job's unless ?name= is given."},
			},
		},
		"encoding": map[string]any{
			"meta": map[string]any{"contentType": "application/json"},
			"file": map[string]any{"contentType": strings.Join(types, ", ")},
		},
	}
	return content
}

// dataSchema is the schema of the data an entry takes: that of the one
// template its data reaches, any of them when it reaches several (a
// template node that picks one by its metadata), else any object or list.
func (d *apiDoc) dataSchema(entry *Node) map[string]any {
	var refs []any
	for _, name := range d.reach(entry, true) {
		n := d.nodes[name]
		if n.Kind != KindTemplate {
			continue
		}
		for _, t := range d.tpls[n.Name] {
			if n.TemplateFrom == "" && t.ID != n.Template {
				continue
			}
			if s := d.templateSchema(n, &t); s != "" {
				refs = append(refs, ref(s))
			}
		}
	}
	switch len(refs) {
	case 0:
		return map[string]any{"type": []string{"object", "array"}}
	case 1:
		return refs[0].(map[string]any)
	}
	return map[string]any{"anyOf": refs}
}

// templateSchema puts the schema of a template's data in the components,
// and names it; "" when the template's fields make none.
func (d *apiDoc) templateSchema(n *Node, t *Template) string {
	name := "TemplateData." + n.Name + "." + t.ID
	if _, ok := d.schemas[name]; ok {
		return name
	}
	b, err := form.Schema(t.Fields, t.Sample)
	if err != nil {
		d.m.o.Log.Warn("openapi: template schema", "node", n.Name, "template", t.ID, "err", err)
		return ""
	}
	var s map[string]any
	if err := json.Unmarshal(b, &s); err != nil {
		return ""
	}
	delete(s, "$schema") // the document's dialect: 2020-12 too
	s["title"] = or(t.Name, t.ID)
	desc := "The data template " + t.ID + " of " + n.Name + " fills."
	if n.TemplateFrom != "" {
		desc += " The metadata " + n.TemplateFrom + " names it: ?meta." + n.TemplateFrom + "=" + t.ID + "."
	}
	s["description"] = desc
	d.schemas[name] = s
	return name
}

// reach walks the nodes a job may reach from a node, along the edges of
// its flow (and, unless passing, OnError): passing, only through the
// nodes that send on what reaches them, as what the entry took goes.
// The nodes are in the order found.
func (d *apiDoc) reach(from *Node, passing bool) []string {
	seen := map[string]bool{from.Name: true}
	var out []string
	queue := []*Node{from}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		edges := n.edges()
		if passing {
			edges = n.flowEdges()
		}
		for _, to := range edges {
			x := d.nodes[to]
			if x == nil || seen[to] {
				continue
			}
			seen[to] = true
			out = append(out, to)
			if k := kindNamed(x.Kind); !passing || k != nil && k.PassesOn {
				queue = append(queue, x)
			}
		}
	}
	return out
}

// gates are the asserts right after an entry, which refuse a job to its
// sender (gate).
func (d *apiDoc) gates(entry *Node) []string {
	var out []string
	for _, name := range entry.To {
		if x := d.nodes[name]; x != nil && x.Kind == KindAssert {
			out = append(out, name)
		}
	}
	return out
}

// metaParams are the ?meta.<key>= parameters of an entry: the keys the
// nodes after it look at, but those the entry sets itself (they win).
// Those the asserts right after it check are described by the checks:
// required unless optional, and their schema.
func (d *apiDoc) metaParams(entry *Node) []any {
	uses := map[string][]string{}
	enums := map[string][]string{}
	checked := map[string][]any{} // the schemas of the gates' checks, by key
	required := map[string]bool{}
	for _, name := range d.gates(entry) {
		for _, c := range d.nodes[name].Checks {
			if _, set := entry.Meta[c.Key]; set || !metaKeyRE.MatchString(c.Key) || c.Key == "name" || c.Key == "user" {
				continue
			}
			checked[c.Key] = append(checked[c.Key], checkSchema(c))
			required[c.Key] = required[c.Key] || !c.Optional
		}
	}
	use := func(key, what string) {
		if !metaKeyRE.MatchString(key) || key == "name" || key == "user" {
			return // the manager's, or the job's own (?name=, ?user=)
		}
		if _, set := entry.Meta[key]; set {
			return
		}
		if !slices.Contains(uses[key], what) {
			uses[key] = append(uses[key], what)
		}
	}
	for _, name := range d.reach(entry, false) {
		n := d.nodes[name]
		for _, r := range n.Rules {
			for key := range r.Match.Meta {
				use(key, "routes the job at the filter "+n.Name)
			}
		}
		for _, c := range n.Checks {
			if slices.Contains(entry.To, n.Name) {
				use(c.Key, "is checked before the job is taken, at the assert "+n.Name+": a job that fails is refused (422)")
			} else {
				use(c.Key, "is checked at the assert "+n.Name+": a job that fails is aborted there")
			}
		}
		for _, c := range n.Cases {
			for _, w := range c.When {
				use(w.Key, "routes the job at the branch "+n.Name)
			}
		}
		if n.Kind == KindTemplate && n.TemplateFrom != "" {
			use(n.TemplateFrom, "names the template "+n.Name+" fills")
			for _, t := range d.tpls[n.Name] {
				if !slices.Contains(enums[n.TemplateFrom], t.ID) {
					enums[n.TemplateFrom] = append(enums[n.TemplateFrom], t.ID)
				}
			}
		}
		if n.Kind == KindDownload {
			if n.By != "" {
				use(n.By, "says whose the document is at the downloader "+n.Name)
			}
			for _, f := range n.Fields {
				use(f, "is listed with the document at the downloader "+n.Name)
			}
		}
		for _, s := range n.expanded() {
			for _, m := range expandRE.FindAllStringSubmatch(s, -1) {
				use(m[1], "goes into what "+n.Name+" ("+n.Kind+") sends")
			}
		}
	}
	keys := make([]string, 0, len(uses))
	for k := range uses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var out []any
	for _, k := range keys {
		s := map[string]any{"type": "string"}
		if e := enums[k]; len(e) > 0 {
			slices.Sort(e)
			s["enum"] = e
		}
		switch cs := checked[k]; len(cs) {
		case 0:
		case 1:
			s = cs[0].(map[string]any)
		default:
			s = map[string]any{"allOf": cs}
		}
		p := map[string]any{"name": "meta." + k, "in": "query",
			"description": "The metadata " + k + ": it " + strings.Join(uses[k], "; it ") + ".", "schema": s}
		if required[k] {
			p["required"] = true
		}
		out = append(out, p)
	}
	return out
}

// checkSchema is the schema of a value a check takes: a number within
// its bounds, or text that matches its pattern, one of its values when
// it has some.
func checkSchema(c Check) map[string]any {
	s := map[string]any{"type": "string"}
	if !c.Optional {
		s["minLength"] = 1
	}
	if c.Min != nil || c.Max != nil {
		s = map[string]any{"type": "number"}
		if c.Min != nil {
			s["minimum"] = *c.Min
		}
		if c.Max != nil {
			s["maximum"] = *c.Max
		}
	}
	if c.Pattern != "" {
		s["pattern"] = c.Pattern
	}
	if len(c.Enum) > 0 {
		var enum []any
		for _, v := range c.Enum {
			if x, err := strconv.ParseFloat(v, 64); err == nil && s["type"] == "number" {
				enum = append(enum, x)
			} else {
				enum = append(enum, v)
			}
		}
		s["enum"] = enum
	}
	if c.Message != "" {
		s["description"] = c.Message
	}
	return s
}

// expanded are the settings of a node that expand {{meta.<key>}}.
func (n *Node) expanded() []string {
	out := append([]string{n.URL, n.From, n.MailSubject, n.Text, n.Path}, n.Recipients...)
	for _, v := range n.Headers {
		out = append(out, v)
	}
	return out
}

// The inboxes.

// inboxSecurity is who works an inbox: an app it lists, a person of its
// own, or the administrator.
var inboxSecurity = []any{map[string]any{"app": []string{ScopeInbox}}, map[string]any{"person": []string{ScopeInbox}}, map[string]any{"admin": []string{}}}

// inbox describes /api/inbox/<name>: what it holds, and their release or
// discard.
func (d *apiDoc) inbox(n *Node) {
	var release []string
	for _, to := range n.To {
		if x := d.nodes[to]; x == nil || !n.stockedTo(x) {
			release = append(release, to)
		}
	}
	what := "an inbox: it holds the jobs that reach it until someone releases or discards them"
	if len(release) > 0 {
		what += "; it releases to " + d.names(release)
	}
	d.tag(n, what)
	job := map[string]any{"name": "job", "in": "path", "required": true, "description": "The job's number.", "schema": map[string]any{"type": "integer", "format": "int32"}}
	item := ref(d.goSchema(reflect.TypeFor[Item]()))
	png := map[string]any{"description": "the image", "content": map[string]any{"image/png": map[string]any{"schema": map[string]any{"type": "string", "format": "binary"}}}}
	op := func(id, summary, desc string, params []any, ok map[string]any, statuses ...int) map[string]any {
		responses := map[string]any{"200": ok}
		d.errorResponses(responses, append([]int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound}, statuses...)...)
		o := map[string]any{"operationId": d.opID(id, n.Name), "tags": []string{n.Name}, "summary": summary, "description": desc,
			"responses": responses, "security": inboxSecurity}
		if len(params) > 0 {
			o["parameters"] = params
		}
		return o
	}
	base := NodePath(KindInbox, n.Name)
	d.paths[base] = map[string]any{"get": op("get_inbox", "What "+or(n.Title, n.Name)+" holds",
		"The items held, oldest first, the ones released or discarded still shown, latest first, and the nodes it releases to (to). The ETag is of the answer: with If-None-Match, 304 while nothing changed.",
		nil, jsonBody("the inbox", ref(d.goSchema(reflect.TypeFor[InboxView]()))))}
	d.paths[base+"/{job}"] = map[string]any{"get": op("get_item", "An item of "+or(n.Title, n.Name), "An item, by its job's number: who sent it, its pages, its state, and where it went.",
		[]any{job}, jsonBody("the item", item))}
	d.paths[base+"/{job}/thumbnail"] = map[string]any{"get": op("get_item_thumbnail", "The thumbnail of an item",
		"Its first page, of a .bdf document: an item of a file as received or of data has none (404), only its document.",
		[]any{job,
			map[string]any{"name": "size", "in": "query", "description": "The longer side, in pixels (256 by default).", "schema": map[string]any{"type": "integer", "minimum": 32, "maximum": 1024}},
			map[string]any{"name": "mode", "in": "query", "description": "fit, the whole page (the default); crop, the square at its top; auto, as the kind of document asks.", "schema": map[string]any{"type": "string", "enum": []string{"fit", "crop", "auto"}}},
		}, png)}
	d.paths[base+"/{job}/pages/{page}"] = map[string]any{"get": op("get_item_page", "A page of an item", "A page of its .bdf document, from 1 (…/pages/1.png too).",
		[]any{job,
			map[string]any{"name": "page", "in": "path", "required": true, "description": "The page's number, from 1.", "schema": map[string]any{"type": "string"}},
			map[string]any{"name": "w", "in": "query", "description": "The width, in pixels (800 by default).", "schema": map[string]any{"type": "integer", "minimum": 64, "maximum": MaxPageWidth}},
		}, png)}
	d.paths[base+"/{job}/document"] = map[string]any{"get": op("get_item_document", "The document of an item", "The document the inbox holds, of its format; the download is in the audit log.",
		[]any{job}, map[string]any{"description": "the document", "content": map[string]any{"*/*": map[string]any{"schema": map[string]any{"type": "string", "format": "binary"}}}})}
	to := map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "The nodes to release to; none, every node the inbox releases to."}
	if len(release) > 0 {
		to["items"] = map[string]any{"type": "string", "enum": release}
	}
	workerProp := map[string]any{"type": "string", "description": "Who decided, written on the item and in the audit log; the reader's name by default."}
	body := func(props map[string]any) map[string]any {
		return map[string]any{"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": props, "additionalProperties": false}}}}
	}
	rel := op("release_item", "Release an item", "Sends the item on, as if the inbox had forwarded it. Only the role topology releases what an inbox that names no nodes holds; the others release to its nodes alone.",
		[]any{job}, jsonBody("the item, released", item), http.StatusBadRequest, http.StatusConflict)
	rel["requestBody"] = body(map[string]any{"to": to, "worker": workerProp})
	d.paths[base+"/{job}/release"] = map[string]any{"post": rel}
	dis := op("discard_item", "Discard an item", "Ends the item: the job goes no further from the inbox.",
		[]any{job}, jsonBody("the item, discarded", item), http.StatusBadRequest, http.StatusConflict)
	dis["requestBody"] = body(map[string]any{"worker": workerProp})
	d.paths[base+"/{job}/discard"] = map[string]any{"post": dis}
}

// inboxes describes GET /api/inbox: the inboxes the reader works.
func (d *apiDoc) inboxes() {
	d.tags = append(d.tags, map[string]any{"name": "inboxes", "description": "The inboxes the reader works, and what each holds."})
	responses := map[string]any{"200": jsonBody("the inboxes", map[string]any{"type": "array", "items": ref(d.goSchema(reflect.TypeFor[InboxSummary]()))})}
	d.errorResponses(responses, http.StatusUnauthorized)
	d.paths[APIPath+"/"+KindInbox] = map[string]any{"get": map[string]any{"operationId": d.opID("list_inboxes"), "tags": []string{"inboxes"}, "summary": "The inboxes",
		"description": "The inboxes the reader works: an app's that list it, a person's of their groups or their own, every one for the role topology; with the nodes each releases to and how many items it holds.",
		"responses":   responses, "security": inboxSecurity}}
}

// The downloaders.

// download describes GET /api/download/<name> and its entries.
func (d *apiDoc) download(n *Node) {
	d.tag(n, "a downloader: the documents an inbox holds or an archive wrote, for their owners and the apps it lists")
	entry := d.goSchema(reflect.TypeFor[DownloadEntry]())
	if len(n.Fields) > 0 {
		props := map[string]any{}
		for _, f := range n.Fields {
			props[f] = map[string]any{"type": "string"}
		}
		name := "DownloadEntry." + n.Name
		d.schemas[name] = map[string]any{"allOf": []any{ref(entry), map[string]any{"properties": map[string]any{
			"fields": map[string]any{"type": "object", "properties": props, "description": "The metadata " + n.Name + " lists with its entries."},
		}}}}
		entry = name
	}
	security := []any{map[string]any{"app": []string{ScopeDownload}}, map[string]any{"person": []string{ScopeDownload}}, map[string]any{"admin": []string{}}}
	owner := "the person a reader reads as"
	if n.By != "" {
		owner = "the value of the metadata " + n.By
	}
	job := map[string]any{"name": "job", "in": "path", "required": true, "description": "The job's number.", "schema": map[string]any{"type": "integer", "format": "int32"}}
	base := NodePath(KindDownload, n.Name)
	get := func(id, summary, desc string, params []any, ok map[string]any) map[string]any {
		responses := map[string]any{"200": ok}
		d.errorResponses(responses, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound)
		op := map[string]any{"operationId": d.opID(id, n.Name), "tags": []string{n.Name}, "summary": summary, "description": desc,
			"responses": responses, "security": security}
		if len(params) > 0 {
			op["parameters"] = params
		}
		return map[string]any{"get": op}
	}
	d.paths[base] = get("list", "List "+or(n.Title, n.Name),
		"The entries the reader may read, newest first: every one's for an app the downloader lists, an administrator and an auditor, a person's own (whose owner is "+owner+").",
		[]any{
			map[string]any{"name": "owner", "in": "query", "description": "One owner's entries alone.", "schema": map[string]any{"type": "string"}},
			map[string]any{"name": "limit", "in": "query", "description": "How many entries at most.", "schema": map[string]any{"type": "integer", "minimum": 1}},
		},
		jsonBody("the entries", map[string]any{"type": "array", "items": ref(entry)}))
	d.paths[base+"/{job}"] = get("get_entry", "An entry of "+or(n.Title, n.Name), "One entry, by its job's number.",
		[]any{job}, jsonBody("the entry", ref(entry)))
	record := map[string]any{"oneOf": []any{ref(d.goSchema(reflect.TypeFor[ArchiveRecord]())), ref(d.goSchema(reflect.TypeFor[Job]()))}}
	d.paths[base+"/{job}/record"] = get("get_record", "The record of an entry", "The record the archive wrote of the job: its metadata, payloads and legs then; of an inbox's entry, the job as it is now.",
		[]any{job}, jsonBody("the record", record))
	d.paths[base+"/{job}/document"] = get("get_document", "The document of an entry", "The document where it is kept (as the archive wrote it, or as the inbox holds it), of the entry's format; the download is in the audit log.",
		[]any{job}, map[string]any{"description": "the document", "content": map[string]any{"*/*": map[string]any{"schema": map[string]any{"type": "string", "format": "binary"}}}})
}

// The jobs an app sent.

func (d *apiDoc) jobs() {
	d.tags = append(d.tags, map[string]any{"name": "jobs", "description": "The jobs an app sent: their state and documents. A sender with an entry's key has none of these."})
	security := []any{map[string]any{"app": []string{}}, map[string]any{"admin": []string{}}}
	job := map[string]any{"name": "id", "in": "path", "required": true, "description": "The job's id, or its number.", "schema": map[string]any{"type": "string"}}
	jobRef := ref(d.goSchema(reflect.TypeFor[Job]()))
	op := func(id, summary, desc string, sec []any, ok map[string]any, params ...any) map[string]any {
		responses := map[string]any{"200": ok}
		d.errorResponses(responses, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound)
		return map[string]any{"operationId": d.opID(id), "tags": []string{"jobs"}, "summary": summary, "description": desc,
			"parameters": append([]any{job}, params...), "responses": responses, "security": sec}
	}
	d.paths["/api/jobs/{id}"] = map[string]any{"get": op("get_job", "A job", "A job the app sent: its state, payloads, metadata and legs.", security, jsonBody("the job", jobRef))}
	d.paths["/api/jobs/{id}/document"] = map[string]any{"get": op("get_job_document", "A job's document",
		"Its .bdf document (the last one a node made, or the one received), else what was received; ?payload= one by its id.", security,
		map[string]any{"description": "the payload", "content": map[string]any{"*/*": map[string]any{"schema": map[string]any{"type": "string", "format": "binary"}}}},
		map[string]any{"name": "payload", "in": "query", "description": "A payload's id (Job.payloads).", "schema": map[string]any{"type": "string"}})}
	d.paths["/api/jobs/{id}/cancel"] = map[string]any{"post": op("cancel_job", "Cancel a job", "Cancels a job the app sent.",
		[]any{map[string]any{"app": []string{ScopeJobsCancel}}, map[string]any{"admin": []string{}}}, jsonBody("the job", jobRef))}
}

// write writes the document: its version is a hash of the rest.
func (d *apiDoc) write() ([]byte, error) {
	slices.SortFunc(d.tags, func(a, b any) int {
		return strings.Compare(a.(map[string]any)["name"].(string), b.(map[string]any)["name"].(string))
	})
	scopes := map[string]any{
		ScopePrint:         "send documents to the uploaders that list the app",
		ScopePrintTemplate: "send data to the uploaders that list the app",
		ScopeDownload:      "read the downloaders that list the app",
		ScopeInbox:         "work the inboxes that list the app: release and discard what they hold",
		ScopeJobsCancel:    "cancel the app's jobs",
	}
	body := map[string]any{
		"paths": d.paths,
		"tags":  d.tags,
		"components": map[string]any{
			"schemas": d.schemas,
			"securitySchemes": map[string]any{
				"key": map[string]any{"type": "http", "scheme": "bearer",
					"description": "The entry's key (bdfk_…), shown once when the node was made or its key rotated."},
				"admin": map[string]any{"type": "http", "scheme": "bearer",
					"description": "The operator's admin token."},
				"app": map[string]any{"type": "oauth2",
					"description": "An app's token: the client credentials grant, the client's id and secret in Basic or in the form. A token lasts " + AppTokenFor.String() + ".",
					"flows":       map[string]any{"clientCredentials": map[string]any{"tokenUrl": d.url(OAuthTokenPath), "scopes": scopes}}},
				"person": map[string]any{"type": "oauth2",
					"description": "A person's token, which a frontend (a public client) gets with the authorization code grant and PKCE (S256).",
					"flows": map[string]any{"authorizationCode": map[string]any{"authorizationUrl": d.url(OAuthAuthorizePath), "tokenUrl": d.url(OAuthTokenPath),
						"scopes": map[string]any{ScopeDownload: "read the downloaders, as the person", ScopeInbox: "work the inboxes of the person"}}}},
			},
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(b)
	body["openapi"] = "3.1.0"
	body["info"] = map[string]any{
		"title":       "BDF Docflow",
		"version":     hex.EncodeToString(sum[:6]),
		"description": "The endpoints this Docflow's topology publishes. Errors are {\"error\": \"…\"}.",
	}
	if d.o.Server != "" {
		body["servers"] = []any{map[string]any{"url": d.o.Server}}
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Schemas of the types the API answers, made of their Go types: a struct
// is an object of its JSON fields, required unless omitempty or omitzero
// (a slice or a map that may be nil is null too); the types are named in
// the components by their names.

// apiEnums are the values of the types that have a fixed set.
var apiEnums = map[reflect.Type][]string{
	reflect.TypeFor[JobState](): {string(JobProcessing), string(JobHeld), string(JobCompleted), string(JobCanceled), string(JobAborted)},
}

// goSchema puts the schema of a struct type in the components, and
// names it.
func (d *apiDoc) goSchema(t reflect.Type) string {
	name := strings.TrimPrefix(t.Name(), "api")
	if _, ok := d.schemas[name]; ok {
		return name
	}
	d.schemas[name] = nil // a type that refers to itself
	props := map[string]any{}
	var required []string
	var fields func(t reflect.Type)
	fields = func(t reflect.Type) {
		for i := range t.NumField() {
			f := t.Field(i)
			tag := f.Tag.Get("json")
			if tag == "-" || !f.IsExported() && !f.Anonymous {
				continue
			}
			key, opts, _ := strings.Cut(tag, ",")
			if f.Anonymous && key == "" && f.Type.Kind() == reflect.Struct {
				fields(f.Type)
				continue
			}
			key = or(key, f.Name)
			omit := strings.Contains(opts, "omitempty") || strings.Contains(opts, "omitzero")
			s := d.typeSchema(f.Type)
			if !omit {
				required = append(required, key)
				if k := f.Type.Kind(); k == reflect.Slice || k == reflect.Map {
					s = nullable(s)
				}
			}
			props[key] = s
		}
	}
	fields(t)
	s := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		s["required"] = required
	}
	d.schemas[name] = s
	return name
}

// nullable is a schema that is also null.
func nullable(s map[string]any) map[string]any {
	if t, ok := s["type"].(string); ok {
		s = maps.Clone(s)
		s["type"] = []string{t, "null"}
		return s
	}
	return map[string]any{"anyOf": []any{s, map[string]any{"type": "null"}}}
}

// typeSchema is the schema of a value of a Go type.
func (d *apiDoc) typeSchema(t reflect.Type) map[string]any {
	if e, ok := apiEnums[t]; ok {
		return map[string]any{"type": "string", "enum": e}
	}
	switch t {
	case reflect.TypeFor[time.Time]():
		return map[string]any{"type": "string", "format": "date-time"}
	case reflect.TypeFor[Duration]():
		return map[string]any{"type": "string", "description": `A duration: "90s", "10m", "2h", "14d".`}
	}
	switch t.Kind() {
	case reflect.Pointer:
		return d.typeSchema(t.Elem())
	case reflect.String:
		return map[string]any{"type": "string"}
	case reflect.Bool:
		return map[string]any{"type": "boolean"}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16:
		return map[string]any{"type": "integer", "format": "int32"}
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return map[string]any{"type": "integer", "format": "int64"}
	case reflect.Float32, reflect.Float64:
		return map[string]any{"type": "number"}
	case reflect.Slice, reflect.Array:
		return map[string]any{"type": "array", "items": d.typeSchema(t.Elem())}
	case reflect.Map:
		s := map[string]any{"type": "object"}
		if t.Elem().Kind() != reflect.Interface {
			s["additionalProperties"] = d.typeSchema(t.Elem())
		}
		return s
	case reflect.Struct:
		return ref(d.goSchema(t))
	}
	return map[string]any{} // any value
}

// openAPI serves the OpenAPI document: ?node= narrows it to nodes (the
// parameter again, or names with commas).
func (h *handler) openAPI(w http.ResponseWriter, r *http.Request, v *Viewer) {
	var nodes []string
	for _, s := range r.URL.Query()["node"] {
		nodes = append(nodes, cleanNames(strings.Split(s, ","))...)
	}
	b, err := h.m.OpenAPI(r.Context(), OpenAPIOptions{Server: origin(r), Viewer: v, Nodes: nodes})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(b)
}

// origin is the scheme and host a request came to.
func origin(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
