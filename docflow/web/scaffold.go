package web

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"text/template"

	"github.com/shibukawa/popcornweb/pw"
)

// The subcommand gen-project writes a project of one's own that composes
// Docflow: a Popcorn Web project (popcornweb.toml beside go.mod) whose
// main is RegisterConfig and Main, with a kind of node of its own (the
// folder exit) to start from. What it writes is scaffold/, the files of
// examples/docflow-custom with the project's name and module path as
// parameters (scaffold_test.go keeps the two the same).

//go:embed all:scaffold
var scaffoldFS embed.FS

// GenProject is the subcommand gen-project: docflow gen-project <dir>
// [--module …] [--name …] [--bdf …] [--force].
type GenProject struct {
	Dir    string `arg:"required" help:"directory of the new project, made if it is not there"`
	Module string `help:"module path of the project (example.com/acme/docflow); the directory's name by default"`
	Name   string `help:"name of the program, its binary and its pw project; the last element of the module path by default"`
	Bdf    string `help:"a checkout of bdf the module is replaced with, as go.mod writes it (../bdf); none takes the module as released"`
	Force  bool   `default:"false" help:"write over files that are there"`
}

// RegisterCommands registers the subcommands of the program (gen-project),
// before the configuration is parsed, as RegisterConfig is called; Run
// answers them.
func RegisterCommands() {
	pw.RegisterSubCommand[GenProject]("gen-project", "write a project of one's own that composes Docflow")
}

// scaffoldData is what the templates of scaffold/ take.
type scaffoldData struct {
	Module, Name, GoVersion, BdfVersion, PopcornwebVersion, Replace string
}

var (
	projectNameRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)
	modulePathRE  = regexp.MustCompile(`^[A-Za-z0-9._~-]+(/[A-Za-z0-9._~-]+)*$`)
)

// GenerateProject writes the project GenProject names, telling out each
// file written and what to do next.
func GenerateProject(c GenProject, out io.Writer) error {
	dir := strings.TrimSpace(c.Dir)
	if dir == "" {
		return errors.New("gen-project: a directory")
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	d := scaffoldData{Module: strings.TrimSpace(c.Module), Name: strings.TrimSpace(c.Name), Replace: strings.TrimSpace(c.Bdf)}
	if d.Module == "" {
		d.Module = filepath.Base(abs)
	}
	if !modulePathRE.MatchString(d.Module) || strings.Contains(d.Module, "..") {
		return fmt.Errorf("gen-project: module path %q: letters, digits, . _ ~ - and / (example.com/acme/docflow)", d.Module)
	}
	if d.Name == "" {
		d.Name = path.Base(d.Module)
	}
	if !projectNameRE.MatchString(d.Name) {
		return fmt.Errorf("gen-project: name %q: letters, digits, _ and -", d.Name)
	}
	d.GoVersion, d.BdfVersion, d.PopcornwebVersion = moduleVersions()
	if d.Replace != "" {
		d.Replace = filepath.ToSlash(d.Replace)
	}
	if st, err := os.Stat(abs); err == nil && !st.IsDir() {
		return fmt.Errorf("gen-project: %s is not a directory", dir)
	}
	// nothing is written before every file is known to be free (or
	// --force), so that a refusal leaves no half of a project
	var files []string
	err = fs.WalkDir(scaffoldFS, "scaffold", func(p string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() {
			return err
		}
		files = append(files, p)
		return nil
	})
	if err != nil {
		return err
	}
	if !c.Force {
		for _, p := range files {
			target := filepath.Join(abs, filepath.FromSlash(strings.TrimSuffix(strings.TrimPrefix(p, "scaffold/"), ".tmpl")))
			if _, err := os.Stat(target); err == nil {
				return fmt.Errorf("gen-project: %s is there already (--force writes over it)", target)
			}
		}
	}
	for _, p := range files {
		rel := strings.TrimSuffix(strings.TrimPrefix(p, "scaffold/"), ".tmpl")
		src, err := scaffoldFS.ReadFile(p)
		if err != nil {
			return err
		}
		// the templates' own delimiters: the files hold {{meta.<key>}}
		t, err := template.New(rel).Delims("[[", "]]").Parse(string(src))
		if err != nil {
			return fmt.Errorf("gen-project: %s: %w", rel, err)
		}
		var b strings.Builder
		if err := t.Execute(&b, d); err != nil {
			return fmt.Errorf("gen-project: %s: %w", rel, err)
		}
		target := filepath.Join(abs, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, []byte(b.String()), 0o644); err != nil {
			return err
		}
		fmt.Fprintf(out, "wrote %s\n", filepath.Join(dir, filepath.FromSlash(rel)))
	}
	fmt.Fprintf(out, "\n%s is a project of its own: module %s, built with pw build to bin/%s.\n", dir, d.Module, d.Name)
	if d.BdfVersion == "v0.0.0" && d.Replace == "" {
		fmt.Fprintf(out, "This docflow was built from a checkout rather than a released version of bdf, so go.mod says bdf v0.0.0: set the version in go.mod, or write the project again with --bdf <checkout>.\n")
	}
	fmt.Fprintf(out, "Next:\n  cd %s\n  go mod tidy\n  pw dev\n", dir)
	return nil
}

// moduleVersions are the versions the project is written for: the Go of
// this binary, and bdf's and popcornweb's as this binary was built with
// them (bdf v0.0.0 when it was built from a checkout: the project needs
// a replace, or a version set by hand).
func moduleVersions() (goVersion, bdf, popcornweb string) {
	goVersion = strings.TrimPrefix(strings.Fields(runtime.Version())[0], "go")
	if !regexp.MustCompile(`^\d+\.\d+`).MatchString(goVersion) {
		goVersion = "1.27"
	}
	bdf, popcornweb = "v0.0.0", "v0.6.0"
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	version := func(m *debug.Module) string {
		for m.Replace != nil {
			m = m.Replace
		}
		if v := m.Version; v != "" && v != "(devel)" && !strings.Contains(v, "-") {
			return v
		}
		return ""
	}
	if bi.Main.Path == "github.com/shibukawa/bdf" {
		if v := version(&bi.Main); v != "" {
			bdf = v
		}
	}
	for _, m := range bi.Deps {
		switch m.Path {
		case "github.com/shibukawa/bdf":
			if v := version(m); v != "" {
				bdf = v
			}
		case "github.com/shibukawa/popcornweb":
			if v := version(m); v != "" {
				popcornweb = v
			}
		}
	}
	return
}
