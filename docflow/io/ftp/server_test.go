package ftp

// A fake FTP server for the tests: passive mode, MLSD (or not, for the
// NLST way), AUTH TLS with a self-signed certificate (or TLS from the
// start, or none), rooted at a directory of the test.

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type ftpOpts struct {
	// TLS offers AUTH TLS; Implicit speaks TLS from the start.
	TLS, Implicit bool
	// MLSD lists with MLSD and MLST; without, NLST, SIZE and MDTM.
	MLSD bool
}

type fakeFTP struct {
	Addr        string
	Fingerprint string
	root        string
	opts        ftpOpts
	tlsCfg      *tls.Config
	ln          net.Listener
	mu          sync.Mutex
	// dataTLS and dataResumed count the data connections over TLS, and
	// those that resumed the control connection's session.
	dataTLS, dataResumed int
}

const (
	testUser     = "tester"
	testPassword = "s3cret"
)

func newFTPServer(t *testing.T, opts ftpOpts) *fakeFTP {
	t.Helper()
	s := &fakeFTP{root: t.TempDir(), opts: opts}
	if opts.TLS || opts.Implicit {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		tmpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "fake ftp"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour),
			IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)}, KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}}
		der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(der)
		s.Fingerprint = hex.EncodeToString(sum[:])
		s.tlsCfg = &tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: key}}, MinVersion: tls.VersionTLS12}
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s.ln, s.Addr = ln, ln.Addr().String()
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go s.serve(c)
		}
	}()
	return s
}

func (s *fakeFTP) local(p string) string {
	p = path.Clean("/" + p)
	return filepath.Join(s.root, filepath.FromSlash(p))
}

func mdtm(t time.Time) string { return t.UTC().Format("20060102150405") }

type session struct {
	s    *fakeFTP
	conn net.Conn
	r    *bufio.Reader
	prot bool
	from string
	pasv net.Listener
}

func (ss *session) reply(code int, msg string) {
	fmt.Fprintf(ss.conn, "%d %s\r\n", code, msg)
}

func (ss *session) replyLines(code int, lines ...string) {
	fmt.Fprintf(ss.conn, "%d-%s\r\n", code, lines[0])
	for _, l := range lines[1 : len(lines)-1] {
		fmt.Fprintf(ss.conn, " %s\r\n", l)
	}
	fmt.Fprintf(ss.conn, "%d %s\r\n", code, lines[len(lines)-1])
}

func (s *fakeFTP) serve(conn net.Conn) {
	defer conn.Close()
	if s.opts.Implicit {
		conn = tls.Server(conn, s.tlsCfg)
	}
	ss := &session{s: s, conn: conn, r: bufio.NewReader(conn)}
	ss.reply(220, "fake ftp")
	for {
		line, err := ss.r.ReadString('\n')
		if err != nil {
			return
		}
		cmd, arg, _ := strings.Cut(strings.TrimRight(line, "\r\n"), " ")
		if !ss.handle(strings.ToUpper(cmd), arg) {
			return
		}
	}
}

// data accepts the data connection of a transfer (passive).
func (ss *session) data() (net.Conn, error) {
	if ss.pasv == nil {
		return nil, fmt.Errorf("no passive port")
	}
	defer func() { ss.pasv.Close(); ss.pasv = nil }()
	ss.pasv.(*net.TCPListener).SetDeadline(time.Now().Add(5 * time.Second))
	c, err := ss.pasv.Accept()
	if err != nil {
		return nil, err
	}
	if ss.prot {
		t := tls.Server(c, ss.s.tlsCfg)
		if err := t.Handshake(); err != nil {
			c.Close()
			return nil, err
		}
		ss.s.mu.Lock()
		ss.s.dataTLS++
		if t.ConnectionState().DidResume {
			ss.s.dataResumed++
		}
		ss.s.mu.Unlock()
		return t, nil
	}
	return c, nil
}

func (ss *session) handle(cmd, arg string) bool {
	s := ss.s
	switch cmd {
	case "USER":
		if arg == "anonymous" {
			ss.reply(230, "anonymous ok")
		} else {
			ss.reply(331, "password")
		}
	case "PASS":
		if arg == testPassword {
			ss.reply(230, "signed in")
		} else {
			ss.reply(530, "wrong password")
		}
	case "FEAT":
		feats := []string{"Features:", "UTF8", "SIZE", "MDTM", "EPSV", "PBSZ", "PROT"}
		if s.opts.TLS {
			feats = append(feats, "AUTH TLS")
		}
		if s.opts.MLSD {
			feats = append(feats, "MLST type*;size*;modify*;", "MLSD")
		}
		ss.replyLines(211, append(feats, "End")...)
	case "AUTH":
		if !s.opts.TLS || strings.ToUpper(arg) != "TLS" {
			ss.reply(502, "no tls here")
			return true
		}
		ss.reply(234, "go ahead")
		t := tls.Server(ss.conn, s.tlsCfg)
		if err := t.Handshake(); err != nil {
			return false
		}
		ss.conn, ss.r = t, bufio.NewReader(t)
	case "PBSZ":
		ss.reply(200, "pbsz 0")
	case "PROT":
		ss.prot = strings.ToUpper(arg) == "P"
		ss.reply(200, "prot "+arg)
	case "TYPE", "OPTS", "NOOP":
		ss.reply(200, "ok")
	case "PWD":
		ss.reply(257, `"/" is here`)
	case "CWD":
		if fi, err := os.Stat(s.local(arg)); err == nil && fi.IsDir() {
			ss.reply(250, "ok")
		} else {
			ss.reply(550, "no such directory")
		}
	case "EPSV", "PASV":
		if ss.pasv != nil {
			ss.pasv.Close()
		}
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			ss.reply(425, err.Error())
			return true
		}
		ss.pasv = ln
		port := ln.Addr().(*net.TCPAddr).Port
		if cmd == "EPSV" {
			ss.reply(229, fmt.Sprintf("Entering Extended Passive Mode (|||%d|)", port))
		} else {
			ss.reply(227, fmt.Sprintf("Entering Passive Mode (127,0,0,1,%d,%d)", port/256, port%256))
		}
	case "MLSD", "NLST":
		if cmd == "MLSD" && !s.opts.MLSD {
			ss.reply(502, "no mlsd")
			return true
		}
		es, err := os.ReadDir(s.local(firstOr(arg, ".")))
		if err != nil {
			ss.reply(550, "no such directory")
			return true
		}
		sort.Slice(es, func(i, j int) bool { return es[i].Name() < es[j].Name() })
		ss.reply(150, "here it comes")
		c, err := ss.data()
		if err != nil {
			ss.reply(425, err.Error())
			return true
		}
		for _, e := range es {
			if cmd == "NLST" {
				fmt.Fprintf(c, "%s\r\n", e.Name())
				continue
			}
			fi, _ := e.Info()
			typ := "file"
			if e.IsDir() {
				typ = "dir"
			}
			fmt.Fprintf(c, "type=%s;size=%d;modify=%s; %s\r\n", typ, fi.Size(), mdtm(fi.ModTime()), e.Name())
		}
		c.Close()
		ss.reply(226, "done")
	case "MLST":
		if !s.opts.MLSD {
			ss.reply(502, "no mlst")
			return true
		}
		fi, err := os.Stat(s.local(arg))
		if err != nil {
			ss.reply(550, "no such file")
			return true
		}
		typ := "file"
		if fi.IsDir() {
			typ = "dir"
		}
		ss.replyLines(250, "Listing", fmt.Sprintf("type=%s;size=%d;modify=%s; %s", typ, fi.Size(), mdtm(fi.ModTime()), arg), "End")
	case "SIZE":
		fi, err := os.Stat(s.local(arg))
		if err != nil || fi.IsDir() {
			ss.reply(550, "no such file")
		} else {
			ss.reply(213, strconv.FormatInt(fi.Size(), 10))
		}
	case "MDTM":
		fi, err := os.Stat(s.local(arg))
		if err != nil {
			ss.reply(550, "no such file")
		} else {
			ss.reply(213, mdtm(fi.ModTime()))
		}
	case "RETR":
		f, err := os.Open(s.local(arg))
		if err != nil {
			ss.reply(550, "no such file")
			return true
		}
		defer f.Close()
		ss.reply(150, "here it comes")
		c, err := ss.data()
		if err != nil {
			ss.reply(425, err.Error())
			return true
		}
		io.Copy(c, f)
		c.Close()
		ss.reply(226, "done")
	case "STOR":
		f, err := os.Create(s.local(arg))
		if err != nil {
			ss.reply(550, err.Error())
			return true
		}
		ss.reply(150, "send it")
		c, err := ss.data()
		if err != nil {
			f.Close()
			ss.reply(425, err.Error())
			return true
		}
		_, err = io.Copy(f, c)
		c.Close()
		f.Close()
		if err != nil {
			ss.reply(451, err.Error())
		} else {
			ss.reply(226, "stored")
		}
	case "RNFR":
		if _, err := os.Lstat(s.local(arg)); err != nil {
			ss.reply(550, "no such file")
		} else {
			ss.from = arg
			ss.reply(350, "and to?")
		}
	case "RNTO":
		if ss.from == "" {
			ss.reply(503, "RNFR first")
		} else if err := os.Rename(s.local(ss.from), s.local(arg)); err != nil {
			ss.reply(550, err.Error())
		} else {
			ss.reply(250, "renamed")
		}
		ss.from = ""
	case "DELE":
		if err := os.Remove(s.local(arg)); err != nil {
			ss.reply(550, "no such file")
		} else {
			ss.reply(250, "deleted")
		}
	case "MKD":
		if err := os.Mkdir(s.local(arg), 0o755); err != nil {
			ss.reply(550, err.Error())
		} else {
			ss.reply(257, `"`+path.Clean("/"+arg)+`" made`)
		}
	case "QUIT":
		ss.reply(221, "bye")
		return false
	default:
		ss.reply(502, "not here: "+cmd)
	}
	return true
}

// writeFile puts a file under the root, with a time.
func (s *fakeFTP) writeFile(t *testing.T, name string, data []byte, at time.Time) {
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
func (s *fakeFTP) names(t *testing.T) []string {
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
