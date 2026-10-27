package sftp

// A fake SFTP server for the tests: version 3 over any transport,
// rooted at a directory of the test (every path, absolute or not, is
// under it), with OpenSSH's posix-rename when asked; and an SSH server
// that serves it as the sftp subsystem, for the kinds' tests.

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

type fakeServer struct {
	root    string
	posix   bool
	files   map[string]*os.File
	dirs    map[string][]fs.DirEntry
	dirPath map[string]string
	next    int
	// renames counts the renames served (the claims of a poll).
	renames int
}

func newFakeServer(root string, posix bool) *fakeServer {
	return &fakeServer{root: root, posix: posix, files: map[string]*os.File{}, dirs: map[string][]fs.DirEntry{}, dirPath: map[string]string{}}
}

// local is a path of the protocol under the root.
func (s *fakeServer) local(p string) string {
	p = strings.TrimPrefix(p, "/")
	if p == "" || p == "." {
		return s.root
	}
	return filepath.Join(s.root, filepath.FromSlash(p))
}

// serve answers requests until the transport ends.
func (s *fakeServer) serve(rw io.ReadWriter) error {
	for {
		var head [4]byte
		if _, err := io.ReadFull(rw, head[:]); err != nil {
			return err
		}
		n := binary.BigEndian.Uint32(head[:])
		p := make([]byte, n)
		if _, err := io.ReadFull(rw, p); err != nil {
			return err
		}
		typ, p := p[0], p[1:]
		if typ == fxpInit {
			ext := []byte{}
			if s.posix {
				ext = append(str("posix-rename@openssh.com"), str("1")...)
			}
			if err := write(rw, fxpVersion, append(u32(3), ext...)); err != nil {
				return err
			}
			continue
		}
		id, p, _ := readU32(p)
		if err := writeParts(rw, s.handle(typ, id, p)); err != nil {
			return err
		}
	}
}

func write(w io.Writer, typ byte, fields ...[]byte) error {
	n := 1
	for _, f := range fields {
		n += len(f)
	}
	buf := make([]byte, 4, 4+n)
	binary.BigEndian.PutUint32(buf, uint32(n))
	buf = append(buf, typ)
	for _, f := range fields {
		buf = append(buf, f...)
	}
	_, err := w.Write(buf)
	return err
}

// writeParts writes a packet handle made: its type, then its fields.
func writeParts(w io.Writer, parts []any) error {
	fields := make([][]byte, 0, len(parts)-1)
	for _, p := range parts[1:] {
		fields = append(fields, p.([]byte))
	}
	return write(w, parts[0].(byte), fields...)
}

func status(id uint32, code uint32, msg string) []any {
	return []any{byte(fxpStatus), u32(id), u32(code), str(msg), str("")}
}

func statusOf(id uint32, err error) []any {
	switch {
	case err == nil:
		return status(id, statusOK, "")
	case errors.Is(err, fs.ErrNotExist):
		return status(id, statusNoSuchFile, err.Error())
	case errors.Is(err, fs.ErrPermission):
		return status(id, statusPermissionDenied, err.Error())
	}
	return status(id, statusFailure, err.Error())
}

func attrs(fi fs.FileInfo) []byte {
	perm := uint32(fi.Mode().Perm())
	if fi.IsDir() {
		perm |= 0o040000
	} else {
		perm |= 0o100000
	}
	out := u32(attrSize | attrPermissions | attrACModTime)
	out = append(out, u64(uint64(fi.Size()))...)
	out = append(out, u32(perm)...)
	out = append(out, u32(uint32(fi.ModTime().Unix()))...)
	return append(out, u32(uint32(fi.ModTime().Unix()))...)
}

// handle answers one request: the packet's type and fields.
func (s *fakeServer) handle(typ byte, id uint32, p []byte) []any {
	switch typ {
	case fxpOpen:
		name, p, _ := readString(p)
		pflags, _, _ := readU32(p)
		flag := 0
		switch {
		case pflags&openRead != 0 && pflags&openWrite != 0:
			flag = os.O_RDWR
		case pflags&openWrite != 0:
			flag = os.O_WRONLY
		}
		if pflags&openCreat != 0 {
			flag |= os.O_CREATE
		}
		if pflags&openTrunc != 0 {
			flag |= os.O_TRUNC
		}
		if pflags&openExcl != 0 {
			flag |= os.O_EXCL
		}
		f, err := os.OpenFile(s.local(name), flag, 0o644)
		if err != nil {
			return statusOf(id, err)
		}
		s.next++
		h := "f" + strconv.Itoa(s.next)
		s.files[h] = f
		return []any{byte(fxpHandle), u32(id), str(h)}
	case fxpClose:
		h, _, _ := readString(p)
		if f, ok := s.files[h]; ok {
			delete(s.files, h)
			return statusOf(id, f.Close())
		}
		if _, ok := s.dirs[h]; ok {
			delete(s.dirs, h)
			return status(id, statusOK, "")
		}
		return status(id, statusFailure, "no such handle")
	case fxpRead:
		h, p, _ := readString(p)
		off, p, _ := readU64(p)
		n, _, _ := readU32(p)
		f, ok := s.files[h]
		if !ok {
			return status(id, statusFailure, "no such handle")
		}
		buf := make([]byte, n)
		got, err := f.ReadAt(buf, int64(off))
		if got == 0 && (err == io.EOF || err == nil) {
			return status(id, statusEOF, "eof")
		}
		return []any{byte(fxpData), u32(id), bstr(buf[:got])}
	case fxpWrite:
		h, p, _ := readString(p)
		off, p, _ := readU64(p)
		data, _, _ := readBytes(p)
		f, ok := s.files[h]
		if !ok {
			return status(id, statusFailure, "no such handle")
		}
		_, err := f.WriteAt(data, int64(off))
		return statusOf(id, err)
	case fxpLstat, fxpStat:
		name, _, _ := readString(p)
		fi, err := os.Stat(s.local(name))
		if err != nil {
			return statusOf(id, err)
		}
		return []any{byte(fxpAttrs), u32(id), attrs(fi)}
	case fxpFstat:
		h, _, _ := readString(p)
		f, ok := s.files[h]
		if !ok {
			return status(id, statusFailure, "no such handle")
		}
		fi, err := f.Stat()
		if err != nil {
			return statusOf(id, err)
		}
		return []any{byte(fxpAttrs), u32(id), attrs(fi)}
	case fxpOpendir:
		name, _, _ := readString(p)
		es, err := os.ReadDir(s.local(name))
		if err != nil {
			return statusOf(id, err)
		}
		s.next++
		h := "d" + strconv.Itoa(s.next)
		s.dirs[h], s.dirPath[h] = es, s.local(name)
		return []any{byte(fxpHandle), u32(id), str(h)}
	case fxpReaddir:
		h, _, _ := readString(p)
		es, ok := s.dirs[h]
		if !ok {
			return status(id, statusFailure, "no such handle")
		}
		if len(es) == 0 {
			return status(id, statusEOF, "eof")
		}
		s.dirs[h] = nil
		out := u32(uint32(len(es)))
		for _, e := range es {
			fi, err := e.Info()
			if err != nil {
				return statusOf(id, err)
			}
			out = append(out, str(e.Name())...)
			out = append(out, str(e.Name())...)
			out = append(out, attrs(fi)...)
		}
		return []any{byte(fxpName), u32(id), out}
	case fxpRemove:
		name, _, _ := readString(p)
		return statusOf(id, os.Remove(s.local(name)))
	case fxpMkdir:
		name, _, _ := readString(p)
		return statusOf(id, os.Mkdir(s.local(name), 0o755))
	case fxpRmdir:
		name, _, _ := readString(p)
		return statusOf(id, os.Remove(s.local(name)))
	case fxpRealpath:
		name, _, _ := readString(p)
		if !strings.HasPrefix(name, "/") {
			name = "/" + strings.TrimPrefix(name, ".")
		}
		fi, err := os.Stat(s.local(name))
		if err != nil {
			return statusOf(id, err)
		}
		return []any{byte(fxpName), u32(id), u32(1), str(name), str(name), attrs(fi)}
	case fxpRename:
		from, p, _ := readString(p)
		to, _, _ := readString(p)
		s.renames++
		if _, err := os.Lstat(s.local(to)); err == nil {
			return status(id, statusFailure, "a file is there already")
		}
		return statusOf(id, os.Rename(s.local(from), s.local(to)))
	case fxpExtended:
		name, p, _ := readString(p)
		if name == "posix-rename@openssh.com" && s.posix {
			from, p, _ := readString(p)
			to, _, _ := readString(p)
			s.renames++
			return statusOf(id, os.Rename(s.local(from), s.local(to)))
		}
		return status(id, statusOpUnsupported, name)
	}
	return status(id, statusOpUnsupported, "type "+strconv.Itoa(int(typ)))
}

// sshServer is an SSH server of the tests: one user by password, one
// by a key, the sftp subsystem served by a fakeServer, on 127.0.0.1.
type sshServer struct {
	Addr    string
	HostKey ssh.PublicKey
	// Signer is the key of the user by key (its PEM in KeyPEM).
	KeyPEM []byte
	root   string
	ln     net.Listener
	srv    *fakeServer
}

const (
	testUser     = "tester"
	testPassword = "s3cret"
)

// newSSHServer starts an SSH server rooted at a directory of the test.
func newSSHServer(t *testing.T, posix bool) *sshServer {
	t.Helper()
	root := t.TempDir()
	_, hostPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	hostSigner, err := ssh.NewSignerFromKey(hostPriv)
	if err != nil {
		t.Fatal(err)
	}
	_, userPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	userSigner, err := ssh.NewSignerFromKey(userPriv)
	if err != nil {
		t.Fatal(err)
	}
	block, err := ssh.MarshalPrivateKey(userPriv, "")
	if err != nil {
		t.Fatal(err)
	}
	cfg := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pw []byte) (*ssh.Permissions, error) {
			if c.User() == testUser && string(pw) == testPassword {
				return nil, nil
			}
			return nil, errors.New("wrong password")
		},
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if c.User() == testUser && string(key.Marshal()) == string(userSigner.PublicKey().Marshal()) {
				return nil, nil
			}
			return nil, errors.New("unknown key")
		},
	}
	cfg.AddHostKey(hostSigner)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := &sshServer{Addr: ln.Addr().String(), HostKey: hostSigner.PublicKey(), KeyPEM: pem.EncodeToMemory(block), root: root, ln: ln, srv: newFakeServer(root, posix)}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go s.serveConn(c, cfg)
		}
	}()
	return s
}

func (s *sshServer) serveConn(c net.Conn, cfg *ssh.ServerConfig) {
	defer c.Close()
	conn, chans, reqs, err := ssh.NewServerConn(c, cfg)
	if err != nil {
		return
	}
	defer conn.Close()
	go ssh.DiscardRequests(reqs)
	for nc := range chans {
		if nc.ChannelType() != "session" {
			nc.Reject(ssh.UnknownChannelType, "session only")
			continue
		}
		ch, requests, err := nc.Accept()
		if err != nil {
			return
		}
		go func() {
			defer ch.Close()
			for req := range requests {
				if req.Type == "subsystem" && len(req.Payload) >= 4 && string(req.Payload[4:]) == "sftp" {
					req.Reply(true, nil)
					s.srv.serve(ch)
					return
				}
				req.Reply(false, nil)
			}
		}()
	}
}

// writeFile puts a file under the root, with a time.
func (s *sshServer) writeFile(t *testing.T, name string, data []byte, at time.Time) {
	t.Helper()
	p := filepath.Join(s.root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if !at.IsZero() {
		if err := os.Chtimes(p, at, at); err != nil {
			t.Fatal(err)
		}
	}
}

// names lists the root's files, sorted, with a directory's as dir/name.
func (s *sshServer) names(t *testing.T) []string {
	t.Helper()
	var out []string
	filepath.WalkDir(s.root, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(s.root, p)
			out = append(out, filepath.ToSlash(rel))
		}
		return nil
	})
	return out
}
