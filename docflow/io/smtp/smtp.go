// Package smtp is the mail output of docflow: the kind smtp sends what
// reaches the node as the attachment of a mail, over SMTP (net/smtp). A
// program that imports the package has the kind; docflow/io/notify sends
// its mails with Send.
//
// The node's URL is the server: smtp://host:587 says STARTTLS, when the
// server offers it, smtps://host:465 TLS from the start. Its credential
// is user:password, given only over TLS or to the machine itself (as
// net/smtp gives PLAIN). From, Recipients, MailSubject and Text expand
// {{meta.<key>}}; a header's value has no line breaks. The Message-ID is
// the delivery's id, the same on every try, for the receiving side to
// tell a second send of one mail (a send is at least once).
package smtp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"github.com/shibukawa/bdf/docflow"
)

// KindSMTP is the kind of the mail output.
const KindSMTP = "smtp"

func init() {
	docflow.RegisterKind(docflow.Kind{Name: KindSMTP,
		Titles: docflow.Text{"en": "Mail (SMTP)", "ja": "メール（SMTP）"}, Help: docflow.Text{"en": "mails it as an attachment", "ja": "添付ファイルとしてメールで送る"},
		Icon: []string{"M3 6h18v12H3z", "M3 7l9 6 9-6"}, Role: docflow.RoleExit, Check: Check, Required: Required(nil), Send: send})
}

// Required are the settings Check refuses a mail without (the server, the
// sender, a recipient), when the settings of when have their values (a
// notify node's via mail); always, with none.
func Required(when map[string]string) []docflow.Required {
	out := make([]docflow.Required, 0, 3)
	for _, k := range []string{"url", "from", "recipients"} {
		out = append(out, docflow.Required{Key: k, If: when})
	}
	return out
}

// Check checks the settings of a mail: the server (URL), the sender
// (From), the recipients and the credential. A notify node that sends
// mail checks them too.
func Check(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	n.URL, n.Credential = strings.TrimSpace(in.URL), strings.TrimSpace(in.Credential)
	if _, _, err := server(n.URL); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	if err := docflow.CheckCredential(n.Credential); err != nil {
		return err
	}
	n.From, n.MailSubject, n.Text = strings.TrimSpace(in.From), strings.TrimSpace(in.MailSubject), in.Text
	if n.From == "" {
		return fmt.Errorf("%w: a mail has a sender (from)", docflow.ErrInvalid)
	}
	for _, r := range in.Recipients {
		if r = strings.TrimSpace(r); r != "" {
			n.Recipients = append(n.Recipients, r)
		}
	}
	if len(n.Recipients) == 0 {
		return fmt.Errorf("%w: a mail has at least one recipient (recipients)", docflow.ErrInvalid)
	}
	return nil
}

// server tells the address of an SMTP server and whether it is TLS from
// the start.
func server(raw string) (addr string, implicit bool, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", false, err
	}
	if u.Host == "" || u.Path != "" && u.Path != "/" {
		return "", false, fmt.Errorf("url %q: smtp://host:port or smtps://host:port", raw)
	}
	switch u.Scheme {
	case "smtp":
		return hostPort(u, "587"), false, nil
	case "smtps":
		return hostPort(u, "465"), true, nil
	}
	return "", false, fmt.Errorf("url %q: smtp://host:port or smtps://host:port", raw)
}

func hostPort(u *url.URL, port string) string {
	if u.Port() != "" {
		return u.Host
	}
	return net.JoinHostPort(u.Hostname(), port)
}

// Mail is a mail to send.
type Mail struct {
	// ID is the Message-ID's left part: a delivery's id.
	ID         string
	From       string
	Recipients []string
	Subject    string
	Text       string
	// Attachment is attached, with its name and content type, when it is
	// not nil.
	Attachment  io.Reader
	Name        string
	ContentType string
}

// TLSConfig is the TLS of the connections to the servers; nil checks
// their certificates against the system's roots (tests set their own).
var TLSConfig *tls.Config

// Send sends a mail through the server of a node (its URL and
// credential). A refusal by the server (a 5xx reply) is final; any other
// failure is tried again.
func Send(ctx context.Context, m *docflow.Manager, n *docflow.Node, msg *Mail) error {
	addr, implicit, err := server(n.URL)
	if err != nil {
		return docflow.Final(err)
	}
	var user, password string
	if n.Credential != "" {
		v, err := m.Credential(n.Credential)
		if err != nil {
			return err
		}
		var ok bool
		if user, password, ok = strings.Cut(v, ":"); !ok {
			return docflow.Final(fmt.Errorf("credential %s: user:password", n.Credential))
		}
	}
	from, err := mail.ParseAddress(msg.From)
	if err != nil {
		return docflow.Final(fmt.Errorf("from %q: %v", msg.From, err))
	}
	var to []*mail.Address
	for _, r := range msg.Recipients {
		as, err := mail.ParseAddressList(r)
		if err != nil {
			return docflow.Final(fmt.Errorf("recipient %q: %v", r, err))
		}
		to = append(to, as...)
	}
	if len(to) == 0 {
		return docflow.Final(errors.New("no recipient"))
	}
	body, err := compose(msg, from, to, time.Now())
	if err != nil {
		return err
	}
	host, _, _ := net.SplitHostPort(addr)
	tc := &tls.Config{ServerName: host}
	if TLSConfig != nil {
		tc = TLSConfig.Clone()
		if tc.ServerName == "" {
			tc.ServerName = host
		}
	}
	d := net.Dialer{Timeout: 30 * time.Second}
	var conn net.Conn
	if implicit {
		conn, err = (&tls.Dialer{NetDialer: &d, Config: tc}).DialContext(ctx, "tcp", addr)
	} else {
		conn, err = d.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return err
	}
	defer conn.Close()
	if dl, ok := ctx.Deadline(); ok {
		conn.SetDeadline(dl)
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return reply(err)
	}
	defer c.Close()
	if err := c.Hello("bdf-docflow"); err != nil {
		return reply(err)
	}
	if !implicit {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tc); err != nil {
				return reply(err)
			}
		}
	}
	if user != "" {
		if err := c.Auth(smtp.PlainAuth("", user, password, host)); err != nil {
			return reply(err)
		}
	}
	if err := c.Mail(from.Address); err != nil {
		return reply(err)
	}
	for _, a := range to {
		if err := c.Rcpt(a.Address); err != nil {
			return reply(err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return reply(err)
	}
	if _, err := w.Write(body); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return reply(err)
	}
	return reply(c.Quit())
}

// reply makes a refusal of the server (a permanent, 5xx reply) final.
func reply(err error) error {
	var te *textproto.Error
	if errors.As(err, &te) && te.Code >= 500 {
		return docflow.Final(err)
	}
	return err
}

// compose writes a mail: a text, and the attachment.
func compose(msg *Mail, from *mail.Address, to []*mail.Address, now time.Time) ([]byte, error) {
	var b bytes.Buffer
	var rcpts []string
	for _, a := range to {
		rcpts = append(rcpts, a.String())
	}
	h := func(k, v string) { fmt.Fprintf(&b, "%s: %s\r\n", k, v) }
	h("From", from.String())
	h("To", strings.Join(rcpts, ", "))
	h("Subject", mime.QEncoding.Encode("utf-8", oneLine(msg.Subject)))
	h("Date", now.Format(time.RFC1123Z))
	h("Message-ID", "<"+msg.ID+"@bdf-docflow>")
	h("MIME-Version", "1.0")
	text := strings.ReplaceAll(strings.ReplaceAll(msg.Text, "\r\n", "\n"), "\n", "\r\n")
	if msg.Attachment == nil {
		h("Content-Type", `text/plain; charset="utf-8"`)
		h("Content-Transfer-Encoding", "base64")
		b.WriteString("\r\n")
		writeBase64(&b, []byte(text))
		return b.Bytes(), nil
	}
	mw := multipart.NewWriter(&b)
	h("Content-Type", `multipart/mixed; boundary="`+mw.Boundary()+`"`)
	b.WriteString("\r\n")
	pw, err := mw.CreatePart(textproto.MIMEHeader{"Content-Type": {`text/plain; charset="utf-8"`}, "Content-Transfer-Encoding": {"base64"}})
	if err != nil {
		return nil, err
	}
	writeBase64(pw, []byte(text))
	ct := msg.ContentType
	if ct == "" {
		ct = "application/octet-stream"
	}
	name := oneLine(msg.Name)
	disp := mime.FormatMediaType("attachment", map[string]string{"filename": name})
	pw, err = mw.CreatePart(textproto.MIMEHeader{"Content-Type": {mime.FormatMediaType(ct, map[string]string{"name": name})}, "Content-Transfer-Encoding": {"base64"},
		"Content-Disposition": {disp}})
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(msg.Attachment)
	if err != nil {
		return nil, err
	}
	writeBase64(pw, data)
	if err := mw.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// writeBase64 writes data in base64, in lines of 76.
func writeBase64(w io.Writer, data []byte) {
	s := base64.StdEncoding.EncodeToString(data)
	for len(s) > 76 {
		io.WriteString(w, s[:76]+"\r\n")
		s = s[76:]
	}
	io.WriteString(w, s+"\r\n")
}

// oneLine is a header's value without line breaks and control characters.
func oneLine(s string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, s))
}

// ContentType is the content type of a payload's format.
func ContentType(format string) string {
	switch format {
	case "bdf":
		return "application/vnd.bdf"
	case "json":
		return "application/json"
	case "yaml":
		return "application/yaml"
	case "csv":
		return "text/csv"
	}
	if t := mime.TypeByExtension("." + format); t != "" {
		return t
	}
	return "application/octet-stream"
}

// send mails what reaches the node, attached.
func send(ctx context.Context, o *docflow.Out) (string, error) {
	r, err := o.Open(ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()
	n := o.Node
	msg := &Mail{ID: o.ID, From: o.ExpandHeader(n.From), Subject: o.ExpandHeader(n.MailSubject), Text: o.Expand(n.Text),
		Attachment: r, Name: AttachmentName(o), ContentType: ContentType(o.Payload.Format)}
	for _, rc := range n.Recipients {
		msg.Recipients = append(msg.Recipients, o.ExpandHeader(rc))
	}
	if err := Send(ctx, o.Manager(), n, msg); err != nil {
		return "", err
	}
	return "<" + o.ID + "@bdf-docflow>", nil
}

// AttachmentName is the name a payload is sent by: its own, else the job's
// with the extension of its format.
func AttachmentName(o *docflow.Out) string {
	if name := docflow.CleanPath(o.Payload.Name); name != "" {
		return name[strings.LastIndex(name, "/")+1:]
	}
	name := docflow.CleanPath(o.Name)
	name = name[strings.LastIndex(name, "/")+1:]
	if name == "" {
		name = "document"
	}
	if o.Payload.Format != "" && !strings.HasSuffix(strings.ToLower(name), "."+o.Payload.Format) {
		name += "." + o.Payload.Format
	}
	return name
}
