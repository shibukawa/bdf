package docflow

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"
)

// fieldy is a kind of a program's own: an exit with settings of its own
// (fields.go), and an API of its own under its nodes.
const kindFieldy = "fieldy"

type fieldySettings struct {
	Dir   string          `json:"dir"`
	Mode  string          `json:"mode"`
	Token string          `json:"token"`
	Limit int             `json:"limit"`
	Keep  Duration        `json:"keep"`
	Quiet bool            `json:"quiet"`
	Extra json.RawMessage `json:"extra"`
}

func init() {
	RegisterKind(Kind{Name: kindFieldy, Role: RoleExit, Titles: Text{"en": "Fieldy", "ja": "フィールディ"}, Help: Text{"en": "a kind of a program's own"}, Group: "acme",
		Icon: []string{"M3 3h18v18H3z"},
		Fields: []Field{
			{Key: "dir", Type: FieldString, Required: true, Title: Text{"en": "Directory", "ja": "ディレクトリ"}},
			{Key: "mode", Type: FieldEnum, Choices: []Choice{{Value: "copy"}, {Value: "move", Label: Text{"en": "Move"}}}, Default: "copy"},
			{Key: "token", Type: FieldCredential, Required: true, If: map[string]string{"mode": "move"}, Title: Text{"en": "Token"}},
			{Key: "limit", Type: FieldInt, Default: 10},
			{Key: "keep", Type: FieldDuration},
			{Key: "quiet", Type: FieldBool},
			{Key: "extra", Type: FieldJSON},
		},
		Send: func(ctx context.Context, o *Out) (string, error) { return "sent", nil },
		Handler: func(m *Manager) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeJSON(w, http.StatusOK, map[string]string{"node": r.PathValue("name"), "path": r.URL.Path})
			})
		},
	})
}

// A node of a kind of its own has settings of its own, checked against
// the kind's fields when it is put and read with their defaults.
func TestKindSettings(t *testing.T) {
	e := newEnv(t, nil)
	put := func(settings string, status int) string {
		t.Helper()
		var n Node
		if err := json.Unmarshal([]byte(`{"kind": "fieldy", "settings": `+settings+`}`), &n); err != nil {
			t.Fatal(err)
		}
		return string(e.node("f", n, status))
	}
	put(`{"dir": " /var/out ", "mode": "move", "token": "t1", "limit": 5, "keep": "120m", "quiet": true, "extra": {"a": [1, 2]}}`, http.StatusOK)
	var n Node
	json.Unmarshal(e.admin("GET", "/api/nodes/f", nil, http.StatusOK), &n)
	if string(n.Settings) != `{"dir":"/var/out","extra":{"a":[1,2]},"keep":"2h","limit":5,"mode":"move","quiet":true,"token":"t1"}` {
		t.Errorf("settings kept: %s", n.Settings)
	}
	s, err := SettingsOf[fieldySettings](&n)
	if err != nil || s.Dir != "/var/out" || s.Mode != "move" || s.Token != "t1" || s.Limit != 5 || s.Keep.String() != "2h" || !s.Quiet || string(s.Extra) != `{"a":[1,2]}` {
		t.Errorf("read: %+v %v", s, err)
	}
	// the defaults fill what is not there; a number and a switch may be
	// written as text (a preset's YAML, a form)
	put(`{"dir": "/var/out", "limit": "7", "quiet": "no"}`, http.StatusOK)
	json.Unmarshal(e.admin("GET", "/api/nodes/f", nil, http.StatusOK), &n)
	if s, err = SettingsOf[fieldySettings](&n); err != nil || s.Mode != "copy" || s.Limit != 7 || s.Quiet || s.Keep != 0 {
		t.Errorf("defaults: %+v %v", s, err)
	}
	if n2 := (Node{Kind: kindFieldy}); true {
		if s, err = SettingsOf[fieldySettings](&n2); err != nil || s.Mode != "copy" || s.Limit != 10 {
			t.Errorf("defaults of none: %+v %v", s, err)
		}
	}
	// what is refused, naming the setting
	for settings, want := range map[string]string{
		`{"dir": "/x", "zzz": 1}`:             `no setting "zzz"`,
		`{"mode": "copy"}`:                    "dir (Directory) is required",
		`{"dir": "/x", "mode": "fly"}`:        `"fly" is not one of copy, move`,
		`{"dir": "/x", "mode": "move"}`:       "token (Token) is required",
		`{"dir": "/x", "limit": "many"}`:      `"many" is not a whole number`,
		`{"dir": "/x", "limit": 1.5}`:         "is not a whole number",
		`{"dir": "/x", "keep": "soon"}`:       "is not a duration",
		`{"dir": "/x", "quiet": "maybe"}`:     "is not true or false",
		`{"dir": "/x", "token": "no spaces"}`: "letters, digits, _ and -",
		`[1]`:                                 "settings:",
		`{"dir": 3}`:                          "is not text",
	} {
		var refused struct{ Error string }
		json.Unmarshal([]byte(put(settings, http.StatusBadRequest)), &refused)
		if !strings.Contains(refused.Error, want) {
			t.Errorf("%s: %s; want %q", settings, refused.Error, want)
		}
	}
	// a kind of the core keeps no settings
	e.node("q", Node{Kind: KindQueue, Settings: json.RawMessage(`{"x": 1}`)}, http.StatusOK)
	var q Node
	json.Unmarshal(e.admin("GET", "/api/nodes/q", nil, http.StatusOK), &q)
	if len(q.Settings) != 0 {
		t.Errorf("a queue's settings: %s", q.Settings)
	}
	// SetSettings writes them back as a kind's Check would
	n = Node{Kind: kindFieldy}
	if err := SetSettings(&n, fieldySettings{Dir: "/y", Mode: "copy"}); err != nil || !strings.Contains(string(n.Settings), `"dir":"/y"`) {
		t.Errorf("SetSettings: %s %v", n.Settings, err)
	}
	// a job sent to it is sent out
	put(`{"dir": "/var/out"}`, http.StatusOK)
	res, out := e.do(nil, "POST", "/api/destinations/f/jobs?name=a", adminToken, testDoc(t, 1))
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("send: %d %s", res.StatusCode, out)
	}
	var j Job
	json.Unmarshal(out, &j)
	e.tick()
	got := e.job(j.Number)
	if len(got.Legs) != 1 || got.Legs[0].State != "done" {
		t.Fatalf("legs: %+v", got.Legs)
	}
	if d, err := e.m().Delivery(context.Background(), got.Legs[0].Ref); err != nil || d.State != DeliveryDone || d.Result != "sent" {
		t.Errorf("delivery: %+v %v", d, err)
	}
}

// The API tells of every kind what the editor needs: its names and help
// by language, its heading, its icon, its fields, and the roles of the
// core's kinds.
func TestKindInfos(t *testing.T) {
	e := newEnv(t, nil)
	var ks []KindInfo
	if err := json.Unmarshal(e.admin("GET", "/api/kinds", nil, http.StatusOK), &ks); err != nil {
		t.Fatal(err)
	}
	by := map[string]KindInfo{}
	for _, k := range ks {
		by[k.Name] = k
		if k.Titles["en"] == "" || k.Group == "" {
			t.Errorf("%s: no English title or group: %+v", k.Name, k)
		}
	}
	if f := by[kindFieldy]; f.Titles["ja"] != "フィールディ" || f.Group != "acme" || len(f.Icon) != 1 || len(f.Fields) != 7 || f.Fields[1].Choices[1].Label["en"] != "Move" || !f.Fields[0].Required {
		t.Errorf("fieldy: %+v", f)
	}
	if r := by[kindFieldy].Required; len(r) != 2 || r[0].Key != "dir" || r[1].Key != "token" || r[1].If["mode"] != "move" {
		t.Errorf("fieldy's required: %+v", r)
	}
	if q := by[KindQueue]; q.Group != GroupControl || q.Titles["ja"] != "キュー" || len(q.Icon) == 0 || q.Picks {
		t.Errorf("queue: %+v", q)
	}
	if f := by[KindFilter]; !f.Picks || f.Group != GroupControl {
		t.Errorf("filter: %+v", f)
	}
	if p := by[KindPrinter]; !p.Prints || p.Group != GroupSink {
		t.Errorf("printer: %+v", p)
	}
	if u := by[KindUpload]; u.Group != GroupSource || u.Prints {
		t.Errorf("upload: %+v", u)
	}
	if c := by[KindConvert]; c.Group != GroupConvert {
		t.Errorf("convert: %+v", c)
	}
	if !slices.ContainsFunc(ks, func(k KindInfo) bool { return k.Name == KindCupsOut && k.Group == GroupSink && k.Titles["ja"] != "" }) {
		t.Error("the adapters' kinds are told too (adapters_test.go)")
	}
}

// A kind's own API is under its nodes, for the role topology by default.
func TestKindHandler(t *testing.T) {
	e := newEnv(t, nil)
	var n Node
	json.Unmarshal([]byte(`{"kind": "fieldy", "settings": {"dir": "/x"}}`), &n)
	e.node("f", n, http.StatusOK)
	var out map[string]string
	json.Unmarshal(e.admin("GET", "/api/fieldy/f/hello/there", nil, http.StatusOK), &out)
	if out["node"] != "f" || out["path"] != "/api/fieldy/f/hello/there" {
		t.Errorf("own API: %v", out)
	}
	if res, _ := e.do(nil, "GET", "/api/fieldy/f/hello", "", nil); res.StatusCode != http.StatusUnauthorized {
		t.Errorf("without a sign-in: %d", res.StatusCode)
	}
}

// What a kind registers is checked: its fields and its access.
func TestRegisterKindChecks(t *testing.T) {
	for name, k := range map[string]Kind{
		"key":     {Name: "bad-key", Role: RoleExit, Fields: []Field{{Key: "Dir", Type: FieldString}}},
		"twice":   {Name: "bad-twice", Role: RoleExit, Fields: []Field{{Key: "a", Type: FieldString}, {Key: "a", Type: FieldInt}}},
		"type":    {Name: "bad-type", Role: RoleExit, Fields: []Field{{Key: "a", Type: "float"}}},
		"enum":    {Name: "bad-enum", Role: RoleExit, Fields: []Field{{Key: "a", Type: FieldEnum}}},
		"default": {Name: "bad-default", Role: RoleExit, Fields: []Field{{Key: "a", Type: FieldInt, Default: "ten"}}},
		"access":  {Name: "bad-access", Role: RoleExit, Access: "anyone"},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: registered", name)
				}
			}()
			RegisterKind(k)
		}()
	}
	if (Text{"en": "Queue", "ja": "キュー"}).In("ja-JP") != "キュー" || (Text{"en": "Queue"}).In("ja") != "Queue" || (Text{"fr": "File"}).In("ja") != "File" {
		t.Error("Text.In")
	}
}
