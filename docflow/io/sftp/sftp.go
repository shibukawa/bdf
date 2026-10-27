package sftp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/io/transfer"
)

// The kinds of node of this package.
const (
	// KindSFTP writes what reaches it to a directory of an SFTP server.
	KindSFTP = "sftp"
	// KindPoll takes the files that appear in a directory of an SFTP
	// server, polling at every tick.
	KindPoll = "sftp-poll"
)

// ports are the schemes a node's URL takes, with their default port.
var ports = map[string]int{"sftp": 22}

// Settings are an exit node's settings: transfer's, and the host key.
type Settings struct {
	transfer.ExitSettings
	HostKey string `json:"host_key"`
}

// PollSettings are an entry node's settings: transfer's, and the host
// key.
type PollSettings struct {
	transfer.PollSettings
	HostKey string `json:"host_key"`
}

func init() {
	url := transfer.URLField("sftp://user@host:22/incoming", docflow.Text{
		"en": "The server, the user and the directory: absolute on the server, /~/dir under the user's own; the port 22 unless given. No password here: name a credential.",
		"ja": "サーバー、ユーザー、ディレクトリ。ディレクトリはサーバー上の絶対パスか、/~/dir でユーザー自身のディレクトリの下。ポートは書かなければ 22。パスワードはここに書かず、資格を名前で指す"})
	credential := transfer.CredentialField(docflow.Text{
		"en": "The user's password, or a PEM private key (a file in the credential directory, beginning with -----BEGIN; not encrypted).",
		"ja": "ユーザーのパスワードか、PEM の秘密鍵（資格のディレクトリのファイル。-----BEGIN で始まる。暗号化していないもの）"}, true)
	hostKey := docflow.Field{Key: "host_key", Type: docflow.FieldString, Required: true,
		Title:       docflow.Text{"en": "Host key", "ja": "ホスト鍵"},
		Help:        docflow.Text{"en": "The server's public key, pinned: Discover fetches it from the server; compare its fingerprint with what the server's administrator says before saving. An OpenSSH line (ssh-ed25519 AAAA…) or a SHA256: fingerprint.", "ja": "サーバーの公開鍵を留める。「サーバーに聞く」でサーバーから取り、保存の前にその指紋をサーバーの管理者の言う値と比べる。OpenSSH の 1 行（ssh-ed25519 AAAA…）か SHA256: の指紋"},
		Placeholder: docflow.Text{"en": "ssh-ed25519 AAAA…"}}
	docflow.RegisterKind(docflow.Kind{Name: KindSFTP,
		Titles: docflow.Text{"en": "SFTP", "ja": "SFTP"}, Help: docflow.Text{"en": "uploaded to an SFTP server", "ja": "SFTP サーバーに上げる"},
		Icon: []string{"M4 5h16v14H4z", "M12 9v6", "M9 12l3 3 3-3", "M8 5V3h8v2"}, Role: docflow.RoleExit,
		Fields: transfer.ExitFields(url, credential, hostKey), Check: checkExit, Send: send, Discover: discover})
	docflow.RegisterKind(docflow.Kind{Name: KindPoll,
		Titles: docflow.Text{"en": "SFTP (poll)", "ja": "SFTP（見回り）"}, Help: docflow.Text{"en": "takes the files that appear on an SFTP server", "ja": "SFTP サーバーに現れたファイルを取る"},
		Icon: []string{"M4 5h16v14H4z", "M12 15V9", "M9 12l3-3 3 3", "M8 5V3h8v2"}, Role: docflow.RoleEntry, PassesOn: true,
		Fields: transfer.PollFields(url, credential, hostKey), Check: checkPoll, Tick: tick, Discover: discover})
}

// site reads a node's URL and host key.
func site(rawURL, hostKey string) (*transfer.Site, error) {
	s, err := transfer.ParseURL(rawURL, ports)
	if err != nil {
		return nil, err
	}
	if s.User == "" {
		return nil, fmt.Errorf("url %q: the user is in the URL (sftp://user@host/dir)", rawURL)
	}
	if strings.TrimSpace(hostKey) == "" {
		return nil, errors.New("the server's host key (host_key) is pinned: Discover fetches it")
	}
	if _, err := ParseHostKey(hostKey); err != nil {
		return nil, err
	}
	return s, nil
}

func checkExit(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	s, err := docflow.SettingsOf[Settings](n)
	if err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	if _, err := site(s.URL, s.HostKey); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	return nil
}

func checkPoll(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	s, err := docflow.SettingsOf[PollSettings](n)
	if err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	if _, err := site(s.URL, s.HostKey); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	if err := transfer.CheckPoll(s.PollSettings); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	return nil
}

// discover fetches the server's host key, for the editor to fill the
// field with: the URL of the settings is enough.
func discover(ctx context.Context, m *docflow.Manager, in docflow.Node) (any, error) {
	s, err := docflow.SettingsOf[Settings](&in)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	st, err := transfer.ParseURL(s.URL, ports)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	info, err := HostKey(ctx, st.Address())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", docflow.ErrPrint, err)
	}
	return info, nil
}

// dialer connects to a node's server with its credential, read at each
// dial: a credential not there is tried again later; a host key that
// does not match is final.
func dialer(m *docflow.Manager, s *transfer.Site, credential, hostKey string) transfer.Dialer {
	return func(ctx context.Context) (transfer.FS, error) {
		check, err := ParseHostKey(hostKey)
		if err != nil {
			return nil, docflow.Final(err)
		}
		v, err := m.Credential(credential)
		if err != nil {
			return nil, err
		}
		cfg := Config{Addr: s.Address(), User: s.User, HostKey: check}
		if strings.HasPrefix(strings.TrimSpace(v), "-----BEGIN") {
			cfg.Key = []byte(v)
		} else {
			cfg.Password = v
		}
		c, err := Dial(ctx, cfg)
		if err != nil {
			var mismatch *HostKeyMismatch
			if errors.As(err, &mismatch) || strings.Contains(err.Error(), "private key") {
				return nil, docflow.Final(err)
			}
			return nil, err
		}
		return c, nil
	}
}

func send(ctx context.Context, o *docflow.Out) (string, error) {
	s, err := docflow.SettingsOf[Settings](o.Node)
	if err != nil {
		return "", docflow.Final(err)
	}
	st, err := site(s.URL, s.HostKey)
	if err != nil {
		return "", docflow.Final(err)
	}
	return transfer.Send(ctx, o, s.ExitSettings, st, dialer(o.Manager(), st, s.Credential, s.HostKey))
}

func tick(ctx context.Context, m *docflow.Manager, n *docflow.Node, rep *docflow.TickReport) error {
	s, err := docflow.SettingsOf[PollSettings](n)
	if err != nil {
		return err
	}
	st, err := site(s.URL, s.HostKey)
	if err != nil {
		return err
	}
	return transfer.Poll(ctx, m, n, rep, s.PollSettings, st, KindPoll, "sftp", dialer(m, st, s.Credential, s.HostKey))
}
