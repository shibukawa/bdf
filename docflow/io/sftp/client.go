// Package sftp is the SFTP input and output of docflow: the kind sftp
// writes what reaches the node to a directory of an SFTP server (an
// exit), and the kind sftp-poll takes the files that appear in one (an
// entry, which polls at every tick). A program that imports the package
// has the kinds. SSH is golang.org/x/crypto/ssh; SFTP itself (version
// 3, draft-ietf-secsh-filexfer-02, with OpenSSH's posix-rename) is
// spoken here, by Client, over any transport.
//
// A node's URL is sftp://user@host:22/base/dir: the user in the URL,
// never the password. Its credential is the user's password, or a PEM
// private key (a file of the credential directory, which begins with
// -----BEGIN). Its host key is pinned (host_key: the server's public key
// as an OpenSSH line, or its SHA256: fingerprint); the editor fetches
// it from the server (Discover) for the administrator to look at before
// the node is saved. There is no way around the host key.
//
// The exit uploads under a temporary name (.<delivery>.part) and renames
// it into place, so that a reader never sees half a file; a file at the
// final name already is taken as written by a try whose outcome was not
// known (unless the node overwrites). The entry claims a file by
// renaming it (<name>.docflow-<node>.taking) before reading it, so that
// two instances never take the same file, then deletes it or moves it
// to a done directory. Both are docflow/io/transfer's, over Conn.
package sftp

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"sync"
	"time"
)

// The protocol's packet types (draft-ietf-secsh-filexfer-02).
const (
	fxpInit          = 1
	fxpVersion       = 2
	fxpOpen          = 3
	fxpClose         = 4
	fxpRead          = 5
	fxpWrite         = 6
	fxpLstat         = 7
	fxpFstat         = 8
	fxpSetstat       = 9
	fxpFsetstat      = 10
	fxpOpendir       = 11
	fxpReaddir       = 12
	fxpRemove        = 13
	fxpMkdir         = 14
	fxpRmdir         = 15
	fxpRealpath      = 16
	fxpStat          = 17
	fxpRename        = 18
	fxpReadlink      = 19
	fxpSymlink       = 20
	fxpStatus        = 101
	fxpHandle        = 102
	fxpData          = 103
	fxpName          = 104
	fxpAttrs         = 105
	fxpExtended      = 200
	fxpExtendedReply = 201
)

// The flags of an open (pflags).
const (
	openRead   = 0x01
	openWrite  = 0x02
	openAppend = 0x04
	openCreat  = 0x08
	openTrunc  = 0x10
	openExcl   = 0x20
)

// The flags of the attributes.
const (
	attrSize        = 0x01
	attrUIDGID      = 0x02
	attrPermissions = 0x04
	attrACModTime   = 0x08
	attrExtended    = 0x80000000
)

// The status codes.
const (
	statusOK               = 0
	statusEOF              = 1
	statusNoSuchFile       = 2
	statusPermissionDenied = 3
	statusFailure          = 4
	statusBadMessage       = 5
	statusNoConnection     = 6
	statusConnectionLost   = 7
	statusOpUnsupported    = 8
)

// Version is the protocol version spoken.
const Version = 3

// maxPacket bounds a packet read from the server (OpenSSH's own bound
// is 256 KiB).
const maxPacket = 1 << 20

// chunk is how much of a file one read or write carries (OpenSSH takes
// 256 KiB, every server 32 KiB); window how many reads or writes are in
// flight at once.
const (
	chunk  = 32 << 10
	window = 16
)

// StatusError is a request the server refused (SSH_FXP_STATUS other than
// OK). It is fs.ErrNotExist for a file not there, fs.ErrPermission for
// one it may not touch.
type StatusError struct {
	Code uint32
	Msg  string
	Op   string
}

func (e *StatusError) Error() string {
	name := map[uint32]string{statusEOF: "eof", statusNoSuchFile: "no such file", statusPermissionDenied: "permission denied", statusFailure: "failure",
		statusBadMessage: "bad message", statusNoConnection: "no connection", statusConnectionLost: "connection lost", statusOpUnsupported: "operation unsupported"}[e.Code]
	if name == "" {
		name = fmt.Sprintf("status %d", e.Code)
	}
	if e.Msg != "" {
		return fmt.Sprintf("sftp: %s: %s (%s)", e.Op, e.Msg, name)
	}
	return fmt.Sprintf("sftp: %s: %s", e.Op, name)
}

func (e *StatusError) Is(target error) bool {
	switch target {
	case fs.ErrNotExist:
		return e.Code == statusNoSuchFile
	case fs.ErrPermission:
		return e.Code == statusPermissionDenied
	}
	return false
}

// Entry is a file of a directory as the server lists it.
type Entry struct {
	Name    string
	Size    int64
	ModTime time.Time
	Mode    fs.FileMode
}

// Client speaks SFTP over a transport: the sftp subsystem of an SSH
// session (Dial), or anything else that reads and writes packets (the
// standard streams of an sftp-server process, in the tests). One request
// is answered at a time but for the reads and writes of a file, which
// are pipelined.
type Client struct {
	mu   sync.Mutex
	rw   io.ReadWriteCloser
	next uint32
	// posixRename tells a server with posix-rename@openssh.com, which
	// replaces the target.
	posixRename bool
}

// NewClient starts the protocol over rw: SSH_FXP_INIT, and the server's
// version and extensions.
func NewClient(rw io.ReadWriteCloser) (*Client, error) {
	c := &Client{rw: rw}
	if err := c.send(fxpInit, u32(Version)); err != nil {
		return nil, err
	}
	typ, p, err := c.recv()
	if err != nil {
		return nil, err
	}
	if typ != fxpVersion {
		return nil, fmt.Errorf("sftp: the server answered %d to INIT, not VERSION", typ)
	}
	v, p, err := readU32(p)
	if err != nil {
		return nil, err
	}
	if v != Version {
		return nil, fmt.Errorf("sftp: the server speaks version %d, not %d", v, Version)
	}
	for len(p) > 0 {
		var name, data string
		if name, p, err = readString(p); err != nil {
			return nil, err
		}
		if data, p, err = readString(p); err != nil {
			return nil, err
		}
		if name == "posix-rename@openssh.com" && data == "1" {
			c.posixRename = true
		}
	}
	return c, nil
}

// Close ends the transport.
func (c *Client) Close() error { return c.rw.Close() }

// PosixRename tells a server whose rename replaces the target.
func (c *Client) PosixRename() bool { return c.posixRename }

// Stat returns a file's attributes, following links.
func (c *Client) Stat(ctx context.Context, name string) (*Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := c.id()
	if err := c.send(fxpStat, u32(id), str(name)); err != nil {
		return nil, err
	}
	typ, p, err := c.reply(id, "stat "+name)
	if err != nil {
		return nil, err
	}
	if typ != fxpAttrs {
		return nil, fmt.Errorf("sftp: stat %s: unexpected packet %d", name, typ)
	}
	e, _, err := readAttrs(p)
	if err != nil {
		return nil, err
	}
	e.Name = path.Base(name)
	return e, nil
}

// ReadDir lists a directory.
func (c *Client) ReadDir(ctx context.Context, dir string) ([]Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := c.id()
	if err := c.send(fxpOpendir, u32(id), str(dir)); err != nil {
		return nil, err
	}
	h, err := c.handle(id, "opendir "+dir)
	if err != nil {
		return nil, err
	}
	var out []Entry
	for {
		id := c.id()
		if err := c.send(fxpReaddir, u32(id), str(h)); err != nil {
			return nil, err
		}
		typ, p, err := c.reply(id, "readdir "+dir)
		if err != nil {
			var se *StatusError
			if errors.As(err, &se) && se.Code == statusEOF {
				break
			}
			c.closeHandle(h)
			return nil, err
		}
		if typ != fxpName {
			c.closeHandle(h)
			return nil, fmt.Errorf("sftp: readdir %s: unexpected packet %d", dir, typ)
		}
		n, p, err := readU32(p)
		if err != nil {
			return nil, err
		}
		for range n {
			var name string
			if name, p, err = readString(p); err != nil {
				return nil, err
			}
			if _, p, err = readString(p); err != nil { // longname
				return nil, err
			}
			var e *Entry
			if e, p, err = readAttrs(p); err != nil {
				return nil, err
			}
			if name == "." || name == ".." {
				continue
			}
			e.Name = name
			out = append(out, *e)
		}
	}
	return out, c.closeHandle(h)
}

// Mkdir makes a directory; one there already is an error (a *StatusError
// of statusFailure on most servers).
func (c *Client) Mkdir(ctx context.Context, dir string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := c.id()
	if err := c.send(fxpMkdir, u32(id), str(dir), u32(0)); err != nil {
		return err
	}
	return c.status(id, "mkdir "+dir)
}

// Remove removes a file.
func (c *Client) Remove(ctx context.Context, name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := c.id()
	if err := c.send(fxpRemove, u32(id), str(name)); err != nil {
		return err
	}
	return c.status(id, "remove "+name)
}

// Rename renames a file. With replace, a file at the new name is
// replaced: posix-rename@openssh.com where the server has it, else the
// file there is removed first (not atomically).
func (c *Client) Rename(ctx context.Context, from, to string, replace bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if replace && c.posixRename {
		id := c.id()
		if err := c.send(fxpExtended, u32(id), str("posix-rename@openssh.com"), str(from), str(to)); err != nil {
			return err
		}
		return c.status(id, "rename "+from)
	}
	if replace {
		id := c.id()
		if err := c.send(fxpRemove, u32(id), str(to)); err != nil {
			return err
		}
		if err := c.status(id, "remove "+to); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	id := c.id()
	if err := c.send(fxpRename, u32(id), str(from), str(to)); err != nil {
		return err
	}
	return c.status(id, "rename "+from)
}

// Open opens a file to read, as a stream: reads are pipelined, window at
// a time.
func (c *Client) Open(ctx context.Context, name string) (io.ReadCloser, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := c.id()
	if err := c.send(fxpOpen, u32(id), str(name), u32(openRead), u32(0)); err != nil {
		return nil, err
	}
	h, err := c.handle(id, "open "+name)
	if err != nil {
		return nil, err
	}
	return &reader{c: c, h: h, name: name}, nil
}

// Create opens a file to write, made or truncated: writes are pipelined.
func (c *Client) Create(ctx context.Context, name string) (io.WriteCloser, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := c.id()
	if err := c.send(fxpOpen, u32(id), str(name), u32(openWrite|openCreat|openTrunc), u32(0)); err != nil {
		return nil, err
	}
	h, err := c.handle(id, "create "+name)
	if err != nil {
		return nil, err
	}
	return &writer{c: c, h: h, name: name}, nil
}

// reader reads a file by its handle, window reads ahead.
type reader struct {
	c      *Client
	h      string
	name   string
	off    uint64
	buf    []byte
	eof    bool
	closed bool
}

func (r *reader) Read(p []byte) (int, error) {
	if len(r.buf) == 0 && !r.eof {
		if err := r.fill(); err != nil {
			return 0, err
		}
	}
	if len(r.buf) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

// fill asks for window chunks from the offset and keeps what came, in
// order; a short or empty chunk ends the file.
func (r *reader) fill() error {
	c := r.c
	c.mu.Lock()
	defer c.mu.Unlock()
	ids := make([]uint32, 0, window)
	for i := range window {
		id := c.id()
		if err := c.send(fxpRead, u32(id), str(r.h), u64(r.off+uint64(i*chunk)), u32(chunk)); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	for _, id := range ids {
		typ, p, err := c.reply(id, "read "+r.name)
		if err != nil {
			var se *StatusError
			if errors.As(err, &se) && se.Code == statusEOF {
				r.eof = true
				continue
			}
			return err
		}
		if r.eof {
			continue // after the end: nothing
		}
		if typ != fxpData {
			return fmt.Errorf("sftp: read %s: unexpected packet %d", r.name, typ)
		}
		data, _, err := readBytes(p)
		if err != nil {
			return err
		}
		r.buf = append(r.buf, data...)
		r.off += uint64(len(data))
		if len(data) < chunk {
			r.eof = true
		}
	}
	return nil
}

func (r *reader) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true
	r.c.mu.Lock()
	defer r.c.mu.Unlock()
	return r.c.closeHandle(r.h)
}

// writer writes a file by its handle, window writes in flight.
type writer struct {
	c      *Client
	h      string
	name   string
	off    uint64
	ids    []uint32
	closed bool
	err    error
}

func (w *writer) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	c := w.c
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for len(p) > 0 {
		m := min(len(p), chunk)
		id := c.id()
		if err := c.send(fxpWrite, u32(id), str(w.h), u64(w.off), bstr(p[:m])); err != nil {
			w.err = err
			return n, err
		}
		w.ids = append(w.ids, id)
		w.off += uint64(m)
		p = p[m:]
		n += m
		if len(w.ids) >= window {
			if err := w.drain(); err != nil {
				w.err = err
				return n, err
			}
		}
	}
	return n, nil
}

// drain takes the statuses of the writes in flight (the lock held).
func (w *writer) drain() error {
	var first error
	for _, id := range w.ids {
		if err := w.c.status(id, "write "+w.name); err != nil && first == nil {
			first = err
		}
	}
	w.ids = w.ids[:0]
	return first
}

func (w *writer) Close() error {
	if w.closed {
		return w.err
	}
	w.closed = true
	w.c.mu.Lock()
	defer w.c.mu.Unlock()
	err := w.drain()
	if cerr := w.c.closeHandle(w.h); err == nil {
		err = cerr
	}
	if w.err == nil {
		w.err = err
	}
	return err
}

// The packets: a length, a type, and the fields.

func (c *Client) id() uint32 {
	c.next++
	return c.next
}

// send writes a packet of a type and fields.
func (c *Client) send(typ byte, fields ...[]byte) error {
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
	_, err := c.rw.Write(buf)
	return err
}

// recv reads a packet: its type and what follows.
func (c *Client) recv() (byte, []byte, error) {
	var head [4]byte
	if _, err := io.ReadFull(c.rw, head[:]); err != nil {
		return 0, nil, fmt.Errorf("sftp: %w", err)
	}
	n := binary.BigEndian.Uint32(head[:])
	if n == 0 || n > maxPacket {
		return 0, nil, fmt.Errorf("sftp: a packet of %d bytes", n)
	}
	p := make([]byte, n)
	if _, err := io.ReadFull(c.rw, p); err != nil {
		return 0, nil, fmt.Errorf("sftp: %w", err)
	}
	return p[0], p[1:], nil
}

// reply reads the answer to a request: its type and what follows the
// id; a status other than OK is a *StatusError, OK is returned as
// fxpStatus with nothing.
func (c *Client) reply(id uint32, op string) (byte, []byte, error) {
	typ, p, err := c.recv()
	if err != nil {
		return 0, nil, err
	}
	got, p, err := readU32(p)
	if err != nil {
		return 0, nil, err
	}
	if got != id {
		return 0, nil, fmt.Errorf("sftp: %s: the answer to request %d came for %d", op, got, id)
	}
	if typ == fxpStatus {
		code, p, err := readU32(p)
		if err != nil {
			return 0, nil, err
		}
		if code != statusOK {
			msg, _, _ := readString(p)
			return 0, nil, &StatusError{Code: code, Msg: msg, Op: op}
		}
		return typ, nil, nil
	}
	return typ, p, nil
}

// status reads the status of a request that answers nothing else.
func (c *Client) status(id uint32, op string) error {
	typ, _, err := c.reply(id, op)
	if err != nil {
		return err
	}
	if typ != fxpStatus {
		return fmt.Errorf("sftp: %s: unexpected packet %d", op, typ)
	}
	return nil
}

// handle reads the handle a request opened.
func (c *Client) handle(id uint32, op string) (string, error) {
	typ, p, err := c.reply(id, op)
	if err != nil {
		return "", err
	}
	if typ != fxpHandle {
		return "", fmt.Errorf("sftp: %s: unexpected packet %d", op, typ)
	}
	h, _, err := readString(p)
	return h, err
}

func (c *Client) closeHandle(h string) error {
	id := c.id()
	if err := c.send(fxpClose, u32(id), str(h)); err != nil {
		return err
	}
	return c.status(id, "close")
}

// readAttrs reads the attributes of a file.
func readAttrs(p []byte) (*Entry, []byte, error) {
	flags, p, err := readU32(p)
	if err != nil {
		return nil, nil, err
	}
	e := &Entry{}
	if flags&attrSize != 0 {
		var size uint64
		if size, p, err = readU64(p); err != nil {
			return nil, nil, err
		}
		e.Size = int64(size)
	}
	if flags&attrUIDGID != 0 {
		if _, p, err = readU32(p); err != nil {
			return nil, nil, err
		}
		if _, p, err = readU32(p); err != nil {
			return nil, nil, err
		}
	}
	if flags&attrPermissions != 0 {
		var perm uint32
		if perm, p, err = readU32(p); err != nil {
			return nil, nil, err
		}
		e.Mode = modeOf(perm)
	}
	if flags&attrACModTime != 0 {
		var mtime uint32
		if _, p, err = readU32(p); err != nil {
			return nil, nil, err
		}
		if mtime, p, err = readU32(p); err != nil {
			return nil, nil, err
		}
		e.ModTime = time.Unix(int64(mtime), 0).UTC()
	}
	if flags&attrExtended != 0 {
		var n uint32
		if n, p, err = readU32(p); err != nil {
			return nil, nil, err
		}
		for range n {
			if _, p, err = readString(p); err != nil {
				return nil, nil, err
			}
			if _, p, err = readString(p); err != nil {
				return nil, nil, err
			}
		}
	}
	return e, p, nil
}

// modeOf is the file mode of POSIX permission bits.
func modeOf(perm uint32) fs.FileMode {
	m := fs.FileMode(perm & 0o777)
	switch perm & 0o170000 {
	case 0o040000:
		m |= fs.ModeDir
	case 0o120000:
		m |= fs.ModeSymlink
	case 0o100000, 0:
	default:
		m |= fs.ModeIrregular
	}
	return m
}

// The fields of a packet.

func u32(v uint32) []byte { return binary.BigEndian.AppendUint32(nil, v) }
func u64(v uint64) []byte { return binary.BigEndian.AppendUint64(nil, v) }
func str(s string) []byte { return append(u32(uint32(len(s))), s...) }
func bstr(b []byte) []byte {
	return append(u32(uint32(len(b))), b...)
}

func readU32(p []byte) (uint32, []byte, error) {
	if len(p) < 4 {
		return 0, nil, errors.New("sftp: a packet ends in a number")
	}
	return binary.BigEndian.Uint32(p), p[4:], nil
}

func readU64(p []byte) (uint64, []byte, error) {
	if len(p) < 8 {
		return 0, nil, errors.New("sftp: a packet ends in a number")
	}
	return binary.BigEndian.Uint64(p), p[8:], nil
}

func readBytes(p []byte) ([]byte, []byte, error) {
	n, p, err := readU32(p)
	if err != nil {
		return nil, nil, err
	}
	if uint64(n) > uint64(len(p)) {
		return nil, nil, errors.New("sftp: a packet ends in a string")
	}
	return p[:n], p[n:], nil
}

func readString(p []byte) (string, []byte, error) {
	b, p, err := readBytes(p)
	return string(b), p, err
}
