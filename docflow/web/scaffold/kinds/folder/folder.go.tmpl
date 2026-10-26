// Package folder is a kind of node of a program's own: an exit that
// writes what reaches it to a directory of the server (a watch folder
// of another program, a share), under a path made of the job's
// metadata. It is what a kind of one's own looks like: registered with
// docflow.RegisterKind in an init, with what the editor shows of it (its
// name and help in each language, its icon, its chapter of the manual)
// and the fields of its settings (Node.Settings), which the editor lays
// out as a form and PutNode checks before Check sees them; Send reads
// them with docflow.SettingsOf. A program imports the package to have
// the kind.
package folder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shibukawa/bdf/docflow"
)

// Kind is the kind's name: the node's kind, and the kind's API path.
const Kind = "folder"

// Settings are a folder node's settings (Node.Settings), as the fields
// below name them.
type Settings struct {
	// Dir is the directory written to, an absolute path on the server.
	Dir string `json:"dir"`
	// Path is the file's path under Dir, expanding {{meta.<key>}}
	// (DefaultPath: the job's number, and the payload's name).
	Path string `json:"path"`
	// Overwrite writes over a file that is there; otherwise one that is
	// there is left (a retry after a write finds it and is done).
	Overwrite bool `json:"overwrite"`
}

// DefaultPath is the path a folder node writes by default.
const DefaultPath = "{{meta.$job}}/{{meta.$name}}"

func init() {
	docflow.RegisterKind(docflow.Kind{
		Name:   Kind,
		Role:   docflow.RoleExit,
		Titles: docflow.Text{"en": "Folder", "ja": "フォルダ"},
		Help:   docflow.Text{"en": "written to a directory of the server", "ja": "サーバーのディレクトリに書く"},
		// a folder, as the editor draws icons: strokes in a 24 by 24 box
		Icon: []string{"M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"},
		Fields: []docflow.Field{
			{Key: "dir", Type: docflow.FieldString, Required: true,
				Title:       docflow.Text{"en": "Directory", "ja": "ディレクトリ"},
				Help:        docflow.Text{"en": "An absolute path on the server, which is there already.", "ja": "サーバー上の絶対パス。既にあるもの"},
				Placeholder: docflow.Text{"en": "/var/spool/acme"}},
			{Key: "path", Type: docflow.FieldPath, Default: DefaultPath, Expands: true,
				Title: docflow.Text{"en": "Path", "ja": "パス"},
				Help:  docflow.Text{"en": "The file's path under the directory: the job's number and the file's name by default.", "ja": "ディレクトリの下のファイルのパス。既定はジョブの番号とファイルの名前"}},
			{Key: "overwrite", Type: docflow.FieldBool,
				Title: docflow.Text{"en": "Overwrite", "ja": "上書き"},
				Help:  docflow.Text{"en": "Write over a file that is there; otherwise it is left as it is.", "ja": "既にあるファイルに上書きする。しないなら、あるものをそのまま残す"}},
		},
		Doc: docflow.Text{
			"en": "A **folder** node writes what reaches it to a directory of the server, under a path made of the job's metadata (`{{meta.<key>}}`), so that another program on the server, or on a share, takes it from there. A file that is there already is left as it is unless the node overwrites.",
			"ja": "**フォルダ**のノードは、届いたものをサーバーのディレクトリに、ジョブのメタデータ（`{{meta.<key>}}`）から作ったパスで書く。サーバーの別のプログラムや共有フォルダがそこから取る。既にあるファイルは、上書きの設定がなければそのまま残す。",
		},
		Check: check,
		Send:  send,
	})
}

// check checks what the fields do not: that the directory is there.
func check(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	s, err := docflow.SettingsOf[Settings](n)
	if err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	if !filepath.IsAbs(s.Dir) {
		return fmt.Errorf("%w: dir %q is not an absolute path", docflow.ErrInvalid, s.Dir)
	}
	if st, err := os.Stat(s.Dir); err != nil || !st.IsDir() {
		return fmt.Errorf("%w: dir %s is not a directory on this server", docflow.ErrInvalid, s.Dir)
	}
	return nil
}

// send writes the payload: to a file of a temporary name first, renamed
// into place, so that a reader never sees half of it. It returns the
// file's path, which the delivery records.
func send(ctx context.Context, o *docflow.Out) (string, error) {
	s, err := docflow.SettingsOf[Settings](o.Node)
	if err != nil {
		return "", docflow.Final(err)
	}
	rel := o.ExpandPath(s.Path)
	if rel == "" {
		return "", docflow.Final(errors.New("the path is empty"))
	}
	dest := filepath.Join(s.Dir, filepath.FromSlash(rel))
	if !s.Overwrite {
		if _, err := os.Stat(dest); err == nil {
			return dest, nil // written before (a try whose outcome was not known)
		}
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}
	r, err := o.Open(ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".docflow-*")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	if err := os.Rename(tmp.Name(), dest); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return dest, nil
}
