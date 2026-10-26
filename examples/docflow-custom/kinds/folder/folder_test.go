package folder_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/flowtest"

	"github.com/shibukawa/bdf/examples/docflow-custom/kinds/folder"
)

func settings(t *testing.T, s folder.Settings) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// A job sent to an uploader in front of a folder node is written to the
// directory, under the path its metadata makes; a try again does not
// write twice.
func TestFolder(t *testing.T) {
	e := flowtest.New(t, nil)
	dir := t.TempDir()
	e.Put(
		docflow.Node{Name: "out", Kind: folder.Kind, Settings: settings(t, folder.Settings{Dir: dir, Path: "{{meta.doc_type}}/{{meta.$name}}"})},
		docflow.Node{Name: "in", Kind: docflow.KindUpload, To: []string{"out"}},
	)
	doc := flowtest.Doc(t, 2)
	j := e.Submit("in", "slip-7.bdf", doc, map[string]any{"doc_type": "invoice"})
	e.Tick()
	if d := e.Delivery(j, "out"); d.State != docflow.DeliveryDone || d.Result != filepath.Join(dir, "invoice", "slip-7.bdf") {
		t.Fatalf("delivery: %+v", d)
	}
	written, err := os.ReadFile(filepath.Join(dir, "invoice", "slip-7.bdf"))
	if err != nil || !bytes.Equal(written, doc) {
		t.Fatalf("written: %d bytes, %v", len(written), err)
	}
	// the default path is the job's number and the file's name; a value
	// cannot leave the directory
	e.Put(docflow.Node{Name: "out", Kind: folder.Kind, Settings: settings(t, folder.Settings{Dir: dir})})
	j = e.Submit("in", "../../escape.bdf", doc, nil)
	e.Tick()
	if d := e.Delivery(j, "out"); d.State != docflow.DeliveryDone || filepath.Dir(filepath.Dir(d.Result)) != dir || filepath.Base(d.Result) != "escape.bdf" {
		t.Fatalf("default path: %+v", d)
	}
	// what is there is left, unless the node overwrites
	os.WriteFile(filepath.Join(dir, "invoice", "slip-7.bdf"), []byte("kept"), 0o644)
	e.Put(docflow.Node{Name: "out", Kind: folder.Kind, Settings: settings(t, folder.Settings{Dir: dir, Path: "{{meta.doc_type}}/{{meta.$name}}"})})
	j = e.Submit("in", "slip-7.bdf", doc, map[string]any{"doc_type": "invoice"})
	e.Tick()
	if b, _ := os.ReadFile(filepath.Join(dir, "invoice", "slip-7.bdf")); string(b) != "kept" {
		t.Errorf("overwritten: %d bytes", len(b))
	}
	e.Put(docflow.Node{Name: "out", Kind: folder.Kind, Settings: settings(t, folder.Settings{Dir: dir, Path: "{{meta.doc_type}}/{{meta.$name}}", Overwrite: true})})
	j = e.Submit("in", "slip-7.bdf", doc, map[string]any{"doc_type": "invoice"})
	e.Tick()
	if b, _ := os.ReadFile(filepath.Join(dir, "invoice", "slip-7.bdf")); !bytes.Equal(b, doc) {
		t.Errorf("not overwritten: %d bytes", len(b))
	}
}

// The settings are checked: the directory is required and must be
// there, and the kind is told with its fields.
func TestFolderSettings(t *testing.T) {
	e := flowtest.New(t, nil)
	dir := t.TempDir()
	e.Required(docflow.Node{Name: "ok", Kind: folder.Kind, Settings: settings(t, folder.Settings{Dir: dir})})
	for _, s := range []folder.Settings{{Dir: "relative/dir"}, {Dir: filepath.Join(dir, "nowhere")}} {
		if _, err := e.M().PutNode(t.Context(), docflow.Node{Name: "bad", Kind: folder.Kind, Settings: settings(t, s)}); err == nil {
			t.Errorf("%+v: put", s)
		}
	}
	for _, k := range docflow.KindInfos() {
		if k.Name == folder.Kind {
			if k.Group != docflow.GroupSink || k.Titles["ja"] != "フォルダ" || len(k.Fields) != 3 || k.Fields[1].Default != folder.DefaultPath {
				t.Errorf("told: %+v", k)
			}
			return
		}
	}
	t.Error("the kind is not registered")
}
