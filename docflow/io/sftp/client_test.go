package sftp

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

// exercise runs a client through what the kinds do: directories, a
// file written in several chunks and read back, renames (with and
// without replacing), removal. Paths are relative (the root of the
// fake server, the working directory of a real sftp-server).
func exercise(t *testing.T, c *Client, root string) {
	t.Helper()
	ctx := context.Background()
	if err := c.Mkdir(ctx, "d"); err != nil {
		t.Fatal(err)
	}
	if err := c.Mkdir(ctx, "d"); err == nil {
		t.Fatal("mkdir of a directory there already: no error")
	}
	data := make([]byte, 100<<10+7)
	rand.Read(data)
	w, err := c.Create(ctx, "d/f")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(data); i += 40000 {
		if _, err := w.Write(data[i:min(i+40000, len(data))]); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(filepath.Join(root, "d", "f")); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("the file on the server: %d bytes, %v", len(got), err)
	}
	e, err := c.Stat(ctx, "d/f")
	if err != nil || e.Size != int64(len(data)) || e.Mode.IsDir() || e.ModTime.IsZero() {
		t.Fatalf("stat: %+v, %v", e, err)
	}
	r, err := c.Open(ctx, "d/f")
	if err != nil {
		t.Fatal(err)
	}
	got, err := io.ReadAll(r)
	if err != nil || !bytes.Equal(got, data) {
		t.Fatalf("read back: %d bytes, %v", len(got), err)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Open(ctx, "d/none"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("open of a file not there: %v", err)
	}
	es, err := c.ReadDir(ctx, "d")
	if err != nil || len(es) != 1 || es[0].Name != "f" || es[0].Size != int64(len(data)) {
		t.Fatalf("readdir: %+v, %v", es, err)
	}
	if err := c.Rename(ctx, "d/f", "d/g", false); err != nil {
		t.Fatal(err)
	}
	w, err = c.Create(ctx, "d/h")
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte("other"))
	w.Close()
	if err := c.Rename(ctx, "d/g", "d/h", false); err == nil {
		t.Fatal("rename over a file without replace: no error")
	}
	if err := c.Rename(ctx, "d/g", "d/h", true); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(filepath.Join(root, "d", "h")); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("after the rename over: %d bytes, %v", len(got), err)
	}
	es, _ = c.ReadDir(ctx, "d")
	if names := namesOf(es); !slices.Equal(names, []string{"h"}) {
		t.Fatalf("after the renames: %v", names)
	}
	if err := c.Remove(ctx, "d/h"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Stat(ctx, "d/h"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("stat of a file removed: %v", err)
	}
	if err := c.Remove(ctx, "d/h"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("remove of a file not there: %v", err)
	}
}

func namesOf(es []Entry) []string {
	var out []string
	for _, e := range es {
		out = append(out, e.Name)
	}
	slices.Sort(out)
	return out
}

// pipeConn is two pipes as one transport.
type pipeConn struct {
	io.Reader
	io.WriteCloser
}

// The client against the fake server, with and without posix-rename.
func TestClientFake(t *testing.T) {
	for _, posix := range []bool{true, false} {
		t.Run(map[bool]string{true: "posix-rename", false: "rename"}[posix], func(t *testing.T) {
			root := t.TempDir()
			srv := newFakeServer(root, posix)
			// over TCP: a pipe with no buffer would hold the server's
			// answers while the client still sends (pipelined writes)
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer ln.Close()
			go func() {
				sc, err := ln.Accept()
				if err != nil {
					return
				}
				defer sc.Close()
				srv.serve(sc)
			}()
			cc, err := net.Dial("tcp", ln.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			c, err := NewClient(cc)
			if err != nil {
				t.Fatal(err)
			}
			if c.PosixRename() != posix {
				t.Fatalf("posix-rename: %v", c.PosixRename())
			}
			exercise(t, c, root)
			c.Close()
		})
	}
}

// The client against OpenSSH's sftp-server, over its standard streams,
// where it is (macOS, Linux with openssh-server): the protocol as the
// real thing speaks it.
func TestClientOpenSSH(t *testing.T) {
	var bin string
	for _, p := range []string{"/usr/libexec/sftp-server", "/usr/lib/openssh/sftp-server", "/usr/lib/ssh/sftp-server", "/usr/libexec/openssh/sftp-server"} {
		if _, err := os.Stat(p); err == nil {
			bin = p
			break
		}
	}
	if bin == "" {
		t.Skip("no sftp-server here")
	}
	root := t.TempDir()
	cmd := exec.Command(bin, "-e")
	cmd.Dir = root
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Skipf("sftp-server: %v", err)
	}
	t.Cleanup(func() { cmd.Process.Kill(); cmd.Wait() })
	c, err := NewClient(pipeConn{Reader: out, WriteCloser: in})
	if err != nil {
		t.Fatal(err)
	}
	if !c.PosixRename() {
		t.Error("OpenSSH's sftp-server has posix-rename")
	}
	exercise(t, c, root)
	c.Close()
}
