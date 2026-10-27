package sftp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/shibukawa/bdf/docflow/io/transfer"
	"golang.org/x/crypto/ssh"
)

// HostKeyInfo is a server's host key as Discover tells it: its type,
// the OpenSSH line to pin (host_key, which the editor fills in), and
// its fingerprint, to compare with what the server's administrator says.
type HostKeyInfo struct {
	Type        string `json:"type"`
	HostKey     string `json:"host_key"`
	Fingerprint string `json:"fingerprint"`
}

// errHostKeyRead ends a handshake that was only for the host key.
var errHostKeyRead = errors.New("sftp: host key read")

// HostKey asks a server its host key, without signing in.
func HostKey(ctx context.Context, addr string) (*HostKeyInfo, error) {
	var got ssh.PublicKey
	cfg := &ssh.ClientConfig{User: "docflow", HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		got = key
		return errHostKeyRead
	}, Timeout: 15 * time.Second}
	conn, err := dialTCP(ctx, addr, cfg.Timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, _, _, err = ssh.NewClientConn(conn, addr, cfg)
	if got == nil {
		return nil, fmt.Errorf("sftp: %s: %v", addr, err)
	}
	return &HostKeyInfo{Type: got.Type(), HostKey: strings.TrimSpace(string(ssh.MarshalAuthorizedKey(got))), Fingerprint: ssh.FingerprintSHA256(got)}, nil
}

// HostKeyMismatch is the error of a server whose host key is not the
// one pinned: nothing a retry mends.
type HostKeyMismatch struct {
	Want, Got string
}

func (e *HostKeyMismatch) Error() string {
	return fmt.Sprintf("sftp: the server's host key is %s, not the %s pinned: the server changed, or this is not the server", e.Got, e.Want)
}

// ParseHostKey reads a pinned host key: an OpenSSH line (ssh-ed25519
// AAAA…, as a known_hosts entry has after the host) or a SHA256:
// fingerprint, and makes the check of a server's key against it.
func ParseHostKey(s string) (ssh.HostKeyCallback, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "SHA256:") {
		if len(s) < 20 {
			return nil, fmt.Errorf("host key %q: a SHA256: fingerprint", s)
		}
		return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			if got := ssh.FingerprintSHA256(key); got != s {
				return &HostKeyMismatch{Want: s, Got: got}
			}
			return nil
		}, nil
	}
	fields := strings.Fields(s)
	if len(fields) >= 3 && strings.Contains(fields[0], ",") || len(fields) >= 3 && !strings.HasPrefix(fields[0], "ssh-") && !strings.HasPrefix(fields[0], "ecdsa-") && !strings.HasPrefix(fields[0], "sk-") {
		s = strings.Join(fields[1:], " ") // a known_hosts line: the host first
	}
	want, _, _, _, err := ssh.ParseAuthorizedKey([]byte(s))
	if err != nil {
		return nil, fmt.Errorf("host key: an OpenSSH line (ssh-ed25519 AAAA…) or a SHA256: fingerprint")
	}
	wantFP := ssh.FingerprintSHA256(want)
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if !bytes.Equal(key.Marshal(), want.Marshal()) {
			return &HostKeyMismatch{Want: wantFP, Got: ssh.FingerprintSHA256(key)}
		}
		return nil
	}, nil
}

// Config is how Dial reaches a server.
type Config struct {
	// Addr is host:port.
	Addr string
	User string
	// Password signs in by password (and keyboard-interactive, answered
	// with it); Key is a PEM private key, which is tried instead.
	Password string
	Key      []byte
	// HostKey checks the server's host key (ParseHostKey).
	HostKey ssh.HostKeyCallback
	// Timeout bounds the connection and the handshake (15 s by default).
	Timeout time.Duration
}

// Conn is an SFTP session over SSH: transfer.FS.
type Conn struct {
	tcp    net.Conn
	client *ssh.Client
	sess   *ssh.Session
	sftp   *Client
	stop   chan struct{}
}

// Dial connects to a server, checks its host key, signs in and opens
// the sftp subsystem. The connection ends when ctx does.
func Dial(ctx context.Context, cfg Config) (*Conn, error) {
	if cfg.HostKey == nil {
		return nil, errors.New("sftp: no host key to check the server against")
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	var auth []ssh.AuthMethod
	switch {
	case len(cfg.Key) > 0:
		signer, err := ssh.ParsePrivateKey(cfg.Key)
		if err != nil {
			var missing *ssh.PassphraseMissingError
			if errors.As(err, &missing) {
				return nil, errors.New("sftp: the private key is encrypted: an unencrypted key")
			}
			return nil, fmt.Errorf("sftp: the private key: %v", err)
		}
		auth = append(auth, ssh.PublicKeys(signer))
	case cfg.Password != "":
		pw := cfg.Password
		auth = append(auth, ssh.Password(pw), ssh.KeyboardInteractive(func(name, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			for i := range answers {
				answers[i] = pw
			}
			return answers, nil
		}))
	default:
		return nil, errors.New("sftp: no password and no private key")
	}
	sshCfg := &ssh.ClientConfig{User: cfg.User, Auth: auth, HostKeyCallback: cfg.HostKey, Timeout: timeout}
	tcp, err := dialTCP(ctx, cfg.Addr, timeout)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(timeout)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	tcp.SetDeadline(deadline)
	cc, chans, reqs, err := ssh.NewClientConn(tcp, cfg.Addr, sshCfg)
	if err != nil {
		tcp.Close()
		var mismatch *HostKeyMismatch
		if errors.As(err, &mismatch) {
			return nil, mismatch
		}
		return nil, fmt.Errorf("sftp: %s: %w", cfg.Addr, err)
	}
	tcp.SetDeadline(time.Time{})
	c := &Conn{tcp: tcp, client: ssh.NewClient(cc, chans, reqs), stop: make(chan struct{})}
	if c.sess, err = c.client.NewSession(); err != nil {
		c.Close()
		return nil, fmt.Errorf("sftp: %s: %w", cfg.Addr, err)
	}
	in, err := c.sess.StdinPipe()
	if err != nil {
		c.Close()
		return nil, err
	}
	out, err := c.sess.StdoutPipe()
	if err != nil {
		c.Close()
		return nil, err
	}
	if err := c.sess.RequestSubsystem("sftp"); err != nil {
		c.Close()
		return nil, fmt.Errorf("sftp: %s: the sftp subsystem: %w", cfg.Addr, err)
	}
	if c.sftp, err = NewClient(&pipe{Reader: out, WriteCloser: in}); err != nil {
		c.Close()
		return nil, err
	}
	go func() {
		select {
		case <-ctx.Done():
			c.tcp.Close()
		case <-c.stop:
		}
	}()
	return c, nil
}

// pipe is the subsystem's streams as one transport.
type pipe struct {
	io.Reader
	io.WriteCloser
}

func dialTCP(ctx context.Context, addr string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("sftp: %w", err)
	}
	return conn, nil
}

// Client is the protocol's client, for what transfer.FS does not do.
func (c *Conn) Client() *Client { return c.sftp }

// Close ends the session and the connection.
func (c *Conn) Close() error {
	select {
	case <-c.stop:
		return nil
	default:
		close(c.stop)
	}
	if c.sess != nil {
		c.sess.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
	return c.tcp.Close()
}

func (c *Conn) List(ctx context.Context, dir string) ([]transfer.Entry, error) {
	es, err := c.sftp.ReadDir(ctx, dir)
	if err != nil {
		return nil, err
	}
	out := make([]transfer.Entry, 0, len(es))
	for _, e := range es {
		out = append(out, transfer.Entry{Name: e.Name, Size: e.Size, ModTime: e.ModTime, Dir: e.Mode.IsDir()})
	}
	return out, nil
}

func (c *Conn) Stat(ctx context.Context, name string) (*transfer.Entry, error) {
	e, err := c.sftp.Stat(ctx, name)
	if err != nil {
		return nil, err
	}
	return &transfer.Entry{Name: e.Name, Size: e.Size, ModTime: e.ModTime, Dir: e.Mode.IsDir()}, nil
}

func (c *Conn) Open(ctx context.Context, name string) (io.ReadCloser, error) {
	return c.sftp.Open(ctx, name)
}

func (c *Conn) Create(ctx context.Context, name string) (io.WriteCloser, error) {
	return c.sftp.Create(ctx, name)
}

func (c *Conn) Rename(ctx context.Context, from, to string, replace bool) error {
	return c.sftp.Rename(ctx, from, to, replace)
}

func (c *Conn) Remove(ctx context.Context, name string) error { return c.sftp.Remove(ctx, name) }

func (c *Conn) Mkdir(ctx context.Context, dir string) error { return c.sftp.Mkdir(ctx, dir) }
