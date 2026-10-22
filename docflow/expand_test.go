package docflow

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// {{meta.<key>}} and nothing else; a header without line breaks; a path
// that does not leave its prefix.
func TestExpand(t *testing.T) {
	o := &Out{ID: "000001.02", Node: &Node{Name: "mail"}, Name: "job",
		Job:     &Job{Number: 1, Name: "a.pdf", Meta: map[string]any{"customer": "acme", "n": 3.0, "evil": "../../x\r\nBcc: y", "$job": 1}},
		Payload: &Payload{ID: "000001.01", Kind: PayloadAny, Format: "pdf", Name: "a.pdf"},
		Cause:   &Failure{Node: "erp", Error: "400 Bad Request"}}
	for in, want := range map[string]string{
		"{{meta.customer}}-{{ meta.n }}":        "acme-3",
		"{{meta.missing}}|{{meta.$job}}":        "|1",
		"{{meta.$node}} {{meta.$delivery}}":     "mail 000001.02",
		"{{meta.$error_node}}: {{meta.$error}}": "erp: 400 Bad Request",
		"{{meta.$name}} {{meta.name}}":          "a.pdf a.pdf",
		"{{customer}} {{meta.customer":          "{{customer}} {{meta.customer",
	} {
		if got := o.Expand(in); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
	if got := o.ExpandHeader("Re: {{meta.evil}}"); strings.ContainsAny(got, "\r\n") {
		t.Errorf("header %q", got)
	}
	for in, want := range map[string]string{
		"out/{{meta.evil}}/{{meta.$name}}": "out/xBcc: y/a.pdf",
		"/{{meta.customer}}//x/../y":       "acme/x/y",
		"{{meta.missing}}":                 "",
	} {
		if got := o.ExpandPath(in); got != want {
			t.Errorf("path %q: %q, want %q", in, got, want)
		}
	}
}

// A credential is its environment variable, else its file; the error
// names it, never a value.
func TestCredential(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "from-file"), []byte("user:pw\r\n"), 0o600)
	m := &Manager{o: Options{CredentialDir: dir}}
	t.Setenv(CredentialEnv("smtp-main"), "env-secret")
	if v, err := m.Credential("smtp-main"); err != nil || v != "env-secret" || CredentialEnv("smtp-main") != "DOCFLOW_CRED_smtp_main" {
		t.Errorf("env: %q %v", v, err)
	}
	if v, err := m.Credential("from-file"); err != nil || v != "user:pw" {
		t.Errorf("file: %q %v", v, err)
	}
	for _, name := range []string{"missing", "../from-file", ""} {
		if _, err := m.Credential(name); !errors.Is(err, ErrCredential) {
			t.Errorf("%q: %v", name, err)
		}
	}
	if err := CheckCredential("a/b"); !errors.Is(err, ErrInvalid) {
		t.Errorf("a/b: %v", err)
	}
}
