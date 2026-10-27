package ftp

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
	// KindFTP writes what reaches it to a directory of an FTP server.
	KindFTP = "ftp"
	// KindPoll takes the files that appear in a directory of an FTP
	// server, polling at every tick.
	KindPoll = "ftp-poll"
)

// ports are the schemes a node's URL takes, with their default port.
var ports = map[string]int{"ftp": 21, "ftps": 990}

// Settings are an exit node's settings: transfer's, and the TLS.
type Settings struct {
	transfer.ExitSettings
	TLS            string `json:"tls"`
	TLSFingerprint string `json:"tls_fingerprint"`
}

// PollSettings are an entry node's settings: transfer's, and the TLS.
type PollSettings struct {
	transfer.PollSettings
	TLS            string `json:"tls"`
	TLSFingerprint string `json:"tls_fingerprint"`
}

func init() {
	url := transfer.URLField("ftp://user@host:21/incoming", docflow.Text{
		"en": "The server, the user and the directory: absolute on the server, /~/dir under the login's own; ftps:// is TLS from the start (port 990). No user: anonymous. No password here: name a credential.",
		"ja": "サーバー、ユーザー、ディレクトリ。ディレクトリはサーバー上の絶対パスか、/~/dir でログインのディレクトリの下。ftps:// は最初から TLS（ポート 990）。ユーザーが無ければ anonymous。パスワードはここに書かず、資格を名前で指す"})
	credential := transfer.CredentialField(docflow.Text{"en": "The user's password.", "ja": "ユーザーのパスワード"}, false)
	tls := docflow.Field{Key: "tls", Type: docflow.FieldEnum, Default: string(TLSExplicit),
		Choices: []docflow.Choice{{Value: string(TLSExplicit), Label: docflow.Text{"en": "required (AUTH TLS)", "ja": "必須（AUTH TLS）"}}, {Value: string(TLSNone), Label: docflow.Text{"en": "none: in the clear", "ja": "なし：平文"}}},
		Title:   docflow.Text{"en": "TLS", "ja": "TLS"},
		Help:    docflow.Text{"en": "On ftp://, the connection is upgraded with AUTH TLS before the password is sent, and a server without it is refused. None sends the password and the files in the clear: only for an old device on the LAN.", "ja": "ftp:// では、パスワードを送る前に AUTH TLS で TLS にし、TLS の無いサーバーは断る。「なし」はパスワードもファイルも平文で流す。LAN の中の古い機器にだけ"}}
	fingerprint := docflow.Field{Key: "tls_fingerprint", Type: docflow.FieldString,
		Title:       docflow.Text{"en": "Certificate fingerprint", "ja": "証明書の指紋"},
		Help:        docflow.Text{"en": "The SHA-256 of the server's certificate (hex), pinned instead of the system's roots: for a self-signed certificate. Empty: the roots.", "ja": "サーバーの証明書の SHA-256（16 進）。システムのルート証明書の代わりに留める。自己署名の証明書のために。空ならルート証明書で確かめる"},
		Placeholder: docflow.Text{"en": "ab:cd:…"}}
	docflow.RegisterKind(docflow.Kind{Name: KindFTP,
		Titles: docflow.Text{"en": "FTP", "ja": "FTP"}, Help: docflow.Text{"en": "uploaded to an FTP server (FTPS)", "ja": "FTP サーバーに上げる（FTPS）"},
		Icon: []string{"M4 5h16v14H4z", "M12 9v6", "M9 12l3 3 3-3", "M4 9h16"}, Role: docflow.RoleExit,
		Fields: transfer.ExitFields(url, credential, tls, fingerprint), Check: checkExit, Send: send})
	docflow.RegisterKind(docflow.Kind{Name: KindPoll,
		Titles: docflow.Text{"en": "FTP (poll)", "ja": "FTP（見回り）"}, Help: docflow.Text{"en": "takes the files that appear on an FTP server", "ja": "FTP サーバーに現れたファイルを取る"},
		Icon: []string{"M4 5h16v14H4z", "M12 15V9", "M9 12l3-3 3 3", "M4 9h16"}, Role: docflow.RoleEntry, PassesOn: true,
		Fields: transfer.PollFields(url, credential, tls, fingerprint), Check: checkPoll, Tick: tick})
}

// site reads a node's URL and TLS settings.
func site(rawURL, tlsMode, fingerprint string) (*transfer.Site, TLSMode, error) {
	s, err := transfer.ParseURL(rawURL, ports)
	if err != nil {
		return nil, "", err
	}
	mode := TLSMode(firstOr(tlsMode, string(TLSExplicit)))
	if s.Scheme == "ftps" {
		mode = TLSImplicit
	} else if mode != TLSExplicit && mode != TLSNone {
		return nil, "", fmt.Errorf("tls %q: explicit or none", tlsMode)
	}
	if fingerprint != "" {
		if _, err := parseFingerprint(fingerprint); err != nil {
			return nil, "", err
		}
	}
	return s, mode, nil
}

func checkExit(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	s, err := docflow.SettingsOf[Settings](n)
	if err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	if _, _, err := site(s.URL, s.TLS, s.TLSFingerprint); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	return nil
}

func checkPoll(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	s, err := docflow.SettingsOf[PollSettings](n)
	if err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	if _, _, err := site(s.URL, s.TLS, s.TLSFingerprint); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	if err := transfer.CheckPoll(s.PollSettings); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	return nil
}

// dialer connects to a node's server with its credential, read at each
// dial: a credential not there is tried again later; a server without
// the TLS asked, or whose certificate is not the one pinned, is final.
func dialer(m *docflow.Manager, s *transfer.Site, mode TLSMode, fingerprint, credential string) transfer.Dialer {
	return func(ctx context.Context) (transfer.FS, error) {
		cfg := Config{Addr: s.Address(), User: s.User, TLS: mode, Fingerprint: fingerprint}
		if credential != "" {
			v, err := m.Credential(credential)
			if err != nil {
				return nil, err
			}
			cfg.Password = v
		}
		c, err := Dial(ctx, cfg)
		if err != nil {
			msg := err.Error()
			if strings.Contains(msg, "offers no TLS") || strings.Contains(msg, "certificate") || strings.Contains(msg, "tls:") {
				return nil, docflow.Final(err)
			}
			var reply *ReplyError
			if errors.As(err, &reply) && strings.HasPrefix(reply.Cmd, "sign in") && reply.Code/100 == 5 {
				return nil, docflow.Final(err) // refused: the credential is not the user's
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
	st, mode, err := site(s.URL, s.TLS, s.TLSFingerprint)
	if err != nil {
		return "", docflow.Final(err)
	}
	return transfer.Send(ctx, o, s.ExitSettings, st, dialer(o.Manager(), st, mode, s.TLSFingerprint, s.Credential))
}

func tick(ctx context.Context, m *docflow.Manager, n *docflow.Node, rep *docflow.TickReport) error {
	s, err := docflow.SettingsOf[PollSettings](n)
	if err != nil {
		return err
	}
	st, mode, err := site(s.URL, s.TLS, s.TLSFingerprint)
	if err != nil {
		return err
	}
	return transfer.Poll(ctx, m, n, rep, s.PollSettings, st, KindPoll, "ftp", dialer(m, st, mode, s.TLSFingerprint, s.Credential))
}
