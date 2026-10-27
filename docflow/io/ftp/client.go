// Package ftp is the FTP input and output of docflow: the kind ftp
// writes what reaches the node to a directory of an FTP server (an
// exit), and the kind ftp-poll takes the files that appear in one (an
// entry, which polls at every tick). A program that imports the package
// has the kinds. FTP is spoken here (RFC 959, with EPSV of RFC 2428,
// MLSD, SIZE and MDTM of RFC 3659, and AUTH TLS of RFC 4217), by Conn,
// over the standard library alone.
//
// A node's URL is ftp://user@host:21/base/dir (anonymous, with no
// user), or ftps://user@host:990/base for TLS from the start (implicit
// TLS): the user in the URL, never the password, which is the node's
// credential. On ftp://, the control connection is upgraded with AUTH
// TLS before the password is sent, and the data connections are TLS
// too (PROT P), resuming the control connection's session as most
// servers require; a server without AUTH TLS is refused, unless the
// node says tls: none, which sends everything in the clear (an old
// device on the LAN). The server's certificate is checked against the
// system's roots, or pinned by its SHA-256 fingerprint (a self-signed
// certificate).
//
// Transfers are binary (TYPE I) in passive mode (EPSV, else PASV: the
// server tells a port, the client connects; the address the server
// tells is ignored in favour of its own, for a server behind NAT).
// Directories are listed with MLSD where the server has it, else NLST
// with SIZE and MDTM of each name. The exit and the entry are
// docflow/io/transfer's, over Conn.
package ftp

import (
	"bufio"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf/docflow/io/transfer"
)

// TLSMode is how a connection is secured.
type TLSMode string

const (
	// TLSExplicit upgrades the control connection with AUTH TLS (and
	// the data connections, PROT P): the default of ftp://.
	TLSExplicit TLSMode = "explicit"
	// TLSImplicit speaks TLS from the start: ftps://, port 990.
	TLSImplicit TLSMode = "implicit"
	// TLSNone sends everything in the clear.
	TLSNone TLSMode = "none"
)

// Config is how Dial reaches a server.
type Config struct {
	// Addr is host:port.
	Addr string
	// User signs in; "" is anonymous (with Password as the mail).
	User     string
	Password string
	TLS      TLSMode
	// ServerName is the name the certificate is checked against (the
	// host of Addr by default).
	ServerName string
	// Fingerprint pins the server's certificate by its SHA-256 (hex,
	// with or without colons, or SHA256: and base64 as OpenSSH writes
	// it) instead of checking it against the roots.
	Fingerprint string
	// Timeout bounds the connection and each command (30 s by default).
	Timeout time.Duration
}

// ReplyError is a command the server refused (a 4xx or 5xx reply). It
// is fs.ErrNotExist for 550, the reply of a file not there (and of a
// file the user may not touch: FTP does not tell them apart).
type ReplyError struct {
	Code int
	Msg  string
	Cmd  string
}

func (e *ReplyError) Error() string {
	return fmt.Sprintf("ftp: %s: %d %s", e.Cmd, e.Code, e.Msg)
}

func (e *ReplyError) Is(target error) bool {
	return target == fs.ErrNotExist && e.Code == 550
}

// Temporary tells a 4xx reply: the server may take the command later.
func (e *ReplyError) Temporary() bool { return e.Code >= 400 && e.Code < 500 }

// Conn is a control connection to a server, signed in: transfer.FS.
// One command runs at a time.
type Conn struct {
	mu      sync.Mutex
	tcp     net.Conn // the connection, as dialed
	ctrl    net.Conn // the same, or TLS over it
	r       *bufio.Reader
	tlsCfg  *tls.Config // nil: data connections in the clear
	host    string      // the server's address, for the data connections
	feat    map[string]string
	timeout time.Duration
	stop    chan struct{}
	// data is the data connection open, for ctx to close it (under dmu:
	// mu is held while it is opened).
	dmu  sync.Mutex
	data net.Conn
}

// Dial connects to a server and signs in. The connection ends when ctx
// does.
func Dial(ctx context.Context, cfg Config) (*Conn, error) {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	host, _, err := net.SplitHostPort(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("ftp: address %q: %v", cfg.Addr, err)
	}
	d := net.Dialer{Timeout: timeout}
	tcp, err := d.DialContext(ctx, "tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("ftp: %w", err)
	}
	c := &Conn{tcp: tcp, ctrl: tcp, host: host, timeout: timeout, stop: make(chan struct{}), feat: map[string]string{}}
	var tlsCfg *tls.Config
	if cfg.TLS != TLSNone {
		if tlsCfg, err = tlsConfig(cfg, host); err != nil {
			tcp.Close()
			return nil, err
		}
	}
	go func() {
		select {
		case <-ctx.Done():
			c.tcp.Close()
			c.dmu.Lock()
			if c.data != nil {
				c.data.Close()
			}
			c.dmu.Unlock()
		case <-c.stop:
		}
	}()
	if cfg.TLS == TLSImplicit {
		if err := c.upgrade(ctx, tlsCfg); err != nil {
			c.Close()
			return nil, err
		}
	}
	c.r = bufio.NewReader(c.ctrl)
	if _, _, err := c.reply(ctx, "greeting"); err != nil {
		c.Close()
		return nil, err
	}
	if cfg.TLS == TLSExplicit {
		code, msg, err := c.cmd(ctx, "AUTH TLS")
		if err != nil {
			c.Close()
			return nil, err
		}
		if code != 234 {
			c.Close()
			return nil, fmt.Errorf("ftp: the server offers no TLS (AUTH TLS: %d %s); tls: none sends in the clear", code, msg)
		}
		if err := c.upgrade(ctx, tlsCfg); err != nil {
			c.Close()
			return nil, err
		}
		c.r = bufio.NewReader(c.ctrl)
	}
	if tlsCfg != nil {
		if _, _, err := c.cmd(ctx, "PBSZ 0"); err != nil {
			c.Close()
			return nil, err
		}
		if code, msg, err := c.cmd(ctx, "PROT P"); err != nil {
			c.Close()
			return nil, err
		} else if code/100 != 2 {
			c.Close()
			return nil, fmt.Errorf("ftp: PROT P: %d %s", code, msg)
		}
		c.tlsCfg = tlsCfg
	}
	user, pass := cfg.User, cfg.Password
	if user == "" {
		user, pass = "anonymous", firstOr(pass, "docflow@")
	}
	code, msg, err := c.cmd(ctx, "USER "+user)
	if err != nil {
		c.Close()
		return nil, err
	}
	if code == 331 {
		if code, msg, err = c.cmd(ctx, "PASS "+pass); err != nil {
			c.Close()
			return nil, err
		}
	}
	if code/100 != 2 {
		c.Close()
		return nil, &ReplyError{Code: code, Msg: msg, Cmd: "sign in as " + user}
	}
	c.features(ctx)
	if _, ok := c.feat["UTF8"]; ok {
		c.cmd(ctx, "OPTS UTF8 ON")
	}
	if code, msg, err := c.cmd(ctx, "TYPE I"); err != nil {
		c.Close()
		return nil, err
	} else if code/100 != 2 {
		c.Close()
		return nil, &ReplyError{Code: code, Msg: msg, Cmd: "TYPE I"}
	}
	return c, nil
}

// tlsConfig is the TLS of a connection: the roots, or a pinned
// certificate.
func tlsConfig(cfg Config, host string) (*tls.Config, error) {
	t := &tls.Config{ServerName: firstOr(cfg.ServerName, host), MinVersion: tls.VersionTLS12, ClientSessionCache: tls.NewLRUClientSessionCache(8)}
	if fp := cfg.Fingerprint; fp != "" {
		want, err := parseFingerprint(fp)
		if err != nil {
			return nil, err
		}
		t.InsecureSkipVerify = true
		t.VerifyPeerCertificate = func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return errors.New("ftp: the server sent no certificate")
			}
			sum := sha256.Sum256(rawCerts[0])
			if hex.EncodeToString(sum[:]) != want {
				return fmt.Errorf("ftp: the server's certificate is SHA256 %s, not the %s pinned", hex.EncodeToString(sum[:]), want)
			}
			return nil
		}
	}
	return t, nil
}

// parseFingerprint reads a certificate's SHA-256: hex, with or without
// colons or spaces, in any case.
func parseFingerprint(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "sha256:")
	s = strings.NewReplacer(":", "", " ", "", "-", "").Replace(s)
	if len(s) != 64 {
		return "", fmt.Errorf("ftp: the fingerprint is the certificate's SHA-256 in hex (64 digits)")
	}
	if _, err := hex.DecodeString(s); err != nil {
		return "", fmt.Errorf("ftp: the fingerprint is the certificate's SHA-256 in hex (64 digits)")
	}
	return s, nil
}

// upgrade makes the control connection TLS.
func (c *Conn) upgrade(ctx context.Context, cfg *tls.Config) error {
	t := tls.Client(c.ctrl, cfg)
	c.ctrl.SetDeadline(c.deadline(ctx))
	if err := t.HandshakeContext(ctx); err != nil {
		return fmt.Errorf("ftp: tls: %w", err)
	}
	c.ctrl.SetDeadline(time.Time{})
	c.ctrl = t
	return nil
}

// features asks what the server has (FEAT): MLSD, MLST, SIZE, MDTM,
// EPSV, UTF8 …, each with what follows its name.
func (c *Conn) features(ctx context.Context) {
	code, msg, err := c.cmd(ctx, "FEAT")
	if err != nil || code != 211 {
		return
	}
	for _, line := range strings.Split(msg, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(strings.ToUpper(line), "FEATURES") || strings.HasPrefix(line, "211") || strings.EqualFold(line, "END") {
			continue
		}
		name, rest, _ := strings.Cut(line, " ")
		c.feat[strings.ToUpper(name)] = rest
	}
}

// Close says goodbye and ends the connection.
func (c *Conn) Close() error {
	select {
	case <-c.stop:
		return nil
	default:
		close(c.stop)
	}
	if c.r != nil {
		c.ctrl.SetDeadline(time.Now().Add(2 * time.Second))
		fmt.Fprintf(c.ctrl, "QUIT\r\n")
	}
	c.ctrl.Close()
	return c.tcp.Close()
}

// deadline is when a command gives up: the timeout, or ctx's, whichever
// comes first.
func (c *Conn) deadline(ctx context.Context) time.Time {
	d := time.Now().Add(c.timeout)
	if t, ok := ctx.Deadline(); ok && t.Before(d) {
		d = t
	}
	return d
}

// cmd sends a command and reads its reply (the lock held by the
// caller, or no other command running).
func (c *Conn) cmd(ctx context.Context, line string) (int, string, error) {
	if strings.ContainsAny(line, "\r\n\x00") {
		return 0, "", fmt.Errorf("ftp: a line break in a name")
	}
	c.ctrl.SetDeadline(c.deadline(ctx))
	if _, err := io.WriteString(c.ctrl, line+"\r\n"); err != nil {
		return 0, "", fmt.Errorf("ftp: %w", err)
	}
	return c.reply(ctx, strings.SplitN(line, " ", 2)[0])
}

// reply reads a reply: its code and text (the lines of a multi-line
// reply joined with line breaks).
func (c *Conn) reply(ctx context.Context, what string) (int, string, error) {
	c.ctrl.SetDeadline(c.deadline(ctx))
	line, err := c.r.ReadString('\n')
	if err != nil {
		return 0, "", fmt.Errorf("ftp: %s: %w", what, err)
	}
	line = strings.TrimRight(line, "\r\n")
	if len(line) < 3 {
		return 0, "", fmt.Errorf("ftp: %s: reply %q", what, line)
	}
	code, err := strconv.Atoi(line[:3])
	if err != nil {
		return 0, "", fmt.Errorf("ftp: %s: reply %q", what, line)
	}
	if len(line) > 3 && line[3] == '-' {
		// multi-line: until a line of the code and a space
		var b strings.Builder
		b.WriteString(line[4:])
		for {
			next, err := c.r.ReadString('\n')
			if err != nil {
				return 0, "", fmt.Errorf("ftp: %s: %w", what, err)
			}
			next = strings.TrimRight(next, "\r\n")
			if strings.HasPrefix(next, line[:3]+" ") {
				b.WriteString("\n" + next[4:])
				break
			}
			b.WriteString("\n" + next)
		}
		return code, b.String(), nil
	}
	return code, strings.TrimSpace(line[3:]), nil
}

// expect sends a command and refuses a reply other than a success.
func (c *Conn) expect(ctx context.Context, line string, ok ...int) (string, error) {
	code, msg, err := c.cmd(ctx, line)
	if err != nil {
		return "", err
	}
	if len(ok) == 0 && code/100 == 2 {
		return msg, nil
	}
	for _, k := range ok {
		if code == k {
			return msg, nil
		}
	}
	return "", &ReplyError{Code: code, Msg: msg, Cmd: strings.SplitN(line, " ", 2)[0]}
}

// passive opens a data connection: EPSV, else PASV, then the command
// of the transfer, whose preliminary reply (150, 125) says the data
// flows.
func (c *Conn) passive(ctx context.Context, command string) (net.Conn, error) {
	var port int
	if code, msg, err := c.cmd(ctx, "EPSV"); err != nil {
		return nil, err
	} else if code == 229 {
		// 229 Entering Extended Passive Mode (|||6446|)
		i, j := strings.LastIndex(msg, "(|||"), strings.LastIndex(msg, "|)")
		if i < 0 || j < i+4 {
			return nil, fmt.Errorf("ftp: EPSV: %q", msg)
		}
		if port, err = strconv.Atoi(msg[i+4 : j]); err != nil {
			return nil, fmt.Errorf("ftp: EPSV: %q", msg)
		}
	} else {
		code, msg, err := c.cmd(ctx, "PASV")
		if err != nil {
			return nil, err
		}
		if code != 227 {
			return nil, &ReplyError{Code: code, Msg: msg, Cmd: "PASV"}
		}
		// 227 Entering Passive Mode (192,0,2,1,25,46)
		i, j := strings.LastIndex(msg, "("), strings.LastIndex(msg, ")")
		if i < 0 || j < i {
			i, j = strings.IndexAny(msg, "0123456789")-1, len(msg)
		}
		parts := strings.Split(msg[i+1:j], ",")
		if len(parts) != 6 {
			return nil, fmt.Errorf("ftp: PASV: %q", msg)
		}
		hi, err1 := strconv.Atoi(strings.TrimSpace(parts[4]))
		lo, err2 := strconv.Atoi(strings.TrimSpace(parts[5]))
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("ftp: PASV: %q", msg)
		}
		port = hi*256 + lo
	}
	d := net.Dialer{Timeout: c.timeout}
	conn, err := d.DialContext(ctx, "tcp", net.JoinHostPort(c.host, strconv.Itoa(port)))
	if err != nil {
		return nil, fmt.Errorf("ftp: data connection: %w", err)
	}
	if c.tlsCfg != nil {
		t := tls.Client(conn, c.tlsCfg)
		conn = t
	}
	c.dmu.Lock()
	c.data = conn
	c.dmu.Unlock()
	code, msg, err := c.cmd(ctx, command)
	if err != nil {
		c.endData(conn)
		return nil, err
	}
	if code/100 != 1 {
		c.endData(conn)
		return nil, &ReplyError{Code: code, Msg: msg, Cmd: strings.SplitN(command, " ", 2)[0]}
	}
	if t, ok := conn.(*tls.Conn); ok {
		conn.SetDeadline(c.deadline(ctx))
		if err := t.HandshakeContext(ctx); err != nil {
			c.endData(conn)
			c.reply(ctx, command) // the transfer's end
			return nil, fmt.Errorf("ftp: data connection: tls: %w", err)
		}
		conn.SetDeadline(time.Time{})
	}
	return conn, nil
}

// endData closes a data connection.
func (c *Conn) endData(conn net.Conn) {
	conn.Close()
	c.dmu.Lock()
	if c.data == conn {
		c.data = nil
	}
	c.dmu.Unlock()
}

// finish ends a transfer: the data connection closed, the final reply
// (226, 250) read.
func (c *Conn) finish(ctx context.Context, conn net.Conn, what string) error {
	c.endData(conn)
	code, msg, err := c.reply(ctx, what)
	if err != nil {
		return err
	}
	if code/100 != 2 {
		return &ReplyError{Code: code, Msg: msg, Cmd: what}
	}
	return nil
}

// transfer runs a command whose data the server sends, and gives the
// lines (a listing).
func (c *Conn) lines(ctx context.Context, command string) ([]string, error) {
	conn, err := c.passive(ctx, command)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(c.deadline(ctx))
	b, rerr := io.ReadAll(io.LimitReader(conn, 64<<20))
	if err := c.finish(ctx, conn, strings.SplitN(command, " ", 2)[0]); err != nil {
		return nil, err
	}
	if rerr != nil {
		return nil, fmt.Errorf("ftp: %s: %w", command, rerr)
	}
	var out []string
	for _, l := range strings.Split(string(b), "\n") {
		if l = strings.TrimRight(l, "\r"); l != "" {
			out = append(out, l)
		}
	}
	return out, nil
}

// List lists a directory: MLSD, else NLST with SIZE and MDTM of each
// name.
func (c *Conn) List(ctx context.Context, dir string) ([]transfer.Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.feat["MLSD"]; ok {
		lines, err := c.lines(ctx, "MLSD "+dir)
		if err != nil {
			return nil, err
		}
		var out []transfer.Entry
		for _, l := range lines {
			e, ok := parseMLSx(l)
			if ok {
				out = append(out, e)
			}
		}
		return out, nil
	}
	lines, err := c.lines(ctx, "NLST "+dir)
	if err != nil {
		return nil, err
	}
	var out []transfer.Entry
	for _, l := range lines {
		name := path.Base(l)
		if name == "." || name == ".." || name == "/" {
			continue
		}
		e := transfer.Entry{Name: name}
		full := path.Join(dir, name)
		if dir == "." {
			full = name
		}
		if msg, err := c.expect(ctx, "SIZE "+full, 213); err == nil {
			e.Size, _ = strconv.ParseInt(strings.TrimSpace(msg), 10, 64)
		} else {
			e.Dir = true // SIZE is of files: a directory has none
		}
		if msg, err := c.expect(ctx, "MDTM "+full, 213); err == nil {
			e.ModTime = parseMDTM(msg)
		}
		out = append(out, e)
	}
	return out, nil
}

// parseMLSx reads a line of MLSD or MLST: facts, a space, the name.
func parseMLSx(line string) (transfer.Entry, bool) {
	facts, name, ok := strings.Cut(strings.TrimLeft(line, " "), " ")
	if !ok || name == "" {
		return transfer.Entry{}, false
	}
	e := transfer.Entry{Name: path.Base(name)}
	for _, f := range strings.Split(facts, ";") {
		k, v, _ := strings.Cut(f, "=")
		switch strings.ToLower(k) {
		case "type":
			switch strings.ToLower(v) {
			case "cdir", "pdir":
				return transfer.Entry{}, false
			case "dir":
				e.Dir = true
			case "file":
			default:
				return transfer.Entry{}, false // a link, a device
			}
		case "size":
			e.Size, _ = strconv.ParseInt(v, 10, 64)
		case "modify":
			e.ModTime = parseMDTM(v)
		}
	}
	if e.Name == "." || e.Name == ".." {
		return transfer.Entry{}, false
	}
	return e, true
}

// parseMDTM reads a time as MDTM and the modify fact write it:
// YYYYMMDDHHMMSS[.sss], UTC.
func parseMDTM(s string) time.Time {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '.'); i >= 0 {
		s = s[:i]
	}
	t, err := time.ParseInLocation("20060102150405", s, time.UTC)
	if err != nil {
		return time.Time{}
	}
	return t
}

// Stat tells a file: MLST, else SIZE and MDTM, else CWD for a
// directory.
func (c *Conn) Stat(ctx context.Context, name string) (*transfer.Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.feat["MLST"]; ok {
		msg, err := c.expect(ctx, "MLST "+name, 250)
		if err != nil {
			return nil, err
		}
		for _, l := range strings.Split(msg, "\n") {
			if e, ok := parseMLSx(l); ok {
				e.Name = path.Base(name)
				return &e, nil
			}
		}
		return nil, fmt.Errorf("ftp: MLST %s: %q", name, msg)
	}
	e := &transfer.Entry{Name: path.Base(name)}
	if msg, err := c.expect(ctx, "SIZE "+name, 213); err == nil {
		e.Size, _ = strconv.ParseInt(strings.TrimSpace(msg), 10, 64)
		if msg, err := c.expect(ctx, "MDTM "+name, 213); err == nil {
			e.ModTime = parseMDTM(msg)
		}
		return e, nil
	}
	// not a file: a directory, if it can be entered (and left again)
	pwd, _ := c.expect(ctx, "PWD", 257)
	if _, err := c.expect(ctx, "CWD "+name, 250); err != nil {
		return nil, &ReplyError{Code: 550, Msg: "no such file or directory", Cmd: "stat " + name}
	}
	e.Dir = true
	if dir := pwdOf(pwd); dir != "" {
		c.expect(ctx, "CWD "+dir, 250)
	}
	return e, nil
}

// pwdOf reads the directory of a PWD reply: "/dir" created.
func pwdOf(msg string) string {
	if i := strings.IndexByte(msg, '"'); i >= 0 {
		if j := strings.LastIndexByte(msg, '"'); j > i {
			return strings.ReplaceAll(msg[i+1:j], `""`, `"`)
		}
	}
	return ""
}

// Open reads a file (RETR).
func (c *Conn) Open(ctx context.Context, name string) (io.ReadCloser, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn, err := c.passive(ctx, "RETR "+name)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(c.deadline(ctx))
	return &dataReader{c: c, conn: conn, ctx: ctx}, nil
}

// dataReader reads a transfer; Close reads its end.
type dataReader struct {
	c    *Conn
	conn net.Conn
	ctx  context.Context
	done bool
}

func (r *dataReader) Read(p []byte) (int, error) {
	r.conn.SetDeadline(r.c.deadline(r.ctx))
	return r.conn.Read(p)
}

func (r *dataReader) Close() error {
	if r.done {
		return nil
	}
	r.done = true
	r.c.mu.Lock()
	defer r.c.mu.Unlock()
	return r.c.finish(r.ctx, r.conn, "RETR")
}

// Create writes a file (STOR); Close ends the transfer and tells
// whether the server took it all.
func (c *Conn) Create(ctx context.Context, name string) (io.WriteCloser, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn, err := c.passive(ctx, "STOR "+name)
	if err != nil {
		return nil, err
	}
	return &dataWriter{c: c, conn: conn, ctx: ctx}, nil
}

// dataWriter writes a transfer; Close reads its end.
type dataWriter struct {
	c    *Conn
	conn net.Conn
	ctx  context.Context
	done bool
}

func (w *dataWriter) Write(p []byte) (int, error) {
	w.conn.SetDeadline(w.c.deadline(w.ctx))
	return w.conn.Write(p)
}

func (w *dataWriter) Close() error {
	if w.done {
		return nil
	}
	w.done = true
	w.c.mu.Lock()
	defer w.c.mu.Unlock()
	if t, ok := w.conn.(*tls.Conn); ok {
		t.CloseWrite()
	}
	return w.c.finish(w.ctx, w.conn, "STOR")
}

// Rename renames a file (RNFR, RNTO). Without replace, a file at the
// new name is an error (FTP would replace it: it is looked for first;
// a race of two renames to the same name is the server's).
func (c *Conn) Rename(ctx context.Context, from, to string, replace bool) error {
	if !replace {
		if _, err := c.Stat(ctx, to); err == nil {
			return &ReplyError{Code: 553, Msg: "a file is there already", Cmd: "RNTO"}
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, err := c.expect(ctx, "RNFR "+from, 350); err != nil {
		return err
	}
	_, err := c.expect(ctx, "RNTO "+to, 250)
	return err
}

// Remove removes a file (DELE).
func (c *Conn) Remove(ctx context.Context, name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.expect(ctx, "DELE "+name, 250)
	return err
}

// Mkdir makes a directory (MKD).
func (c *Conn) Mkdir(ctx context.Context, dir string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.expect(ctx, "MKD "+dir, 257)
	return err
}

// Features are what the server said it has (FEAT), by name.
func (c *Conn) Features() map[string]string { return c.feat }

func firstOr(s, def string) string {
	if s != "" {
		return s
	}
	return def
}
