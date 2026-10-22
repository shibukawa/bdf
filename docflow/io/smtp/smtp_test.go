package smtp_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	"github.com/shibukawa/bdf/docflow/io/smtp"
)

// A document mailed to the recipients its metadata names: the subject
// expanded without the line break a sender put in it, the document
// attached, the delivery's id as the Message-ID, the credential's user
// and password given to the server.
func TestMail(t *testing.T) {
	srv := flowtest.NewSMTP(t)
	t.Setenv(docflow.CredentialEnv("mail-main"), "docflow:pa55")
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "mail", Kind: smtp.KindSMTP, URL: "smtp://" + srv.Addr, Credential: "mail-main", From: "Docflow <docflow@example.com>",
		Recipients: []string{"{{meta.customer}}@example.com", "office@example.com"}, MailSubject: "Invoice {{meta.order}}", Text: "Dear {{meta.customer}},\nthe invoice is attached."},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"mail"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", "invoice.bdf", doc, map[string]any{"customer": "acme", "order": "A-1\r\nBcc: evil@example.com"})
	e.Tick()
	d := e.Delivery(j, "mail")
	if d.State != docflow.DeliveryDone {
		t.Fatalf("%s: %s", d.State, d.LastError)
	}
	mails := srv.Mails()
	if len(mails) != 1 {
		t.Fatalf("%d mails", len(mails))
	}
	m := mails[0]
	if m.From != "docflow@example.com" || strings.Join(m.To, ",") != "acme@example.com,office@example.com" || m.Auth != "\x00docflow\x00pa55" {
		t.Errorf("from %s to %v auth %q", m.From, m.To, m.Auth)
	}
	msg, err := mail.ReadMessage(strings.NewReader(m.Data))
	if err != nil {
		t.Fatal(err)
	}
	if id := msg.Header.Get("Message-Id"); id != "<"+d.ID+"@bdf-docflow>" || d.Result != id {
		t.Errorf("message id %s, result %s", id, d.Result)
	}
	subject, _ := new(mime.WordDecoder).DecodeHeader(msg.Header.Get("Subject"))
	if subject != "Invoice A-1  Bcc: evil@example.com" || msg.Header.Get("Bcc") != "" {
		t.Errorf("subject %q", subject)
	}
	_, params, _ := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	mr := multipart.NewReader(msg.Body, params["boundary"])
	var parts []string
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(p)
		parts = append(parts, p.Header.Get("Content-Type"))
		if p.FileName() == "invoice.bdf" && !bytes.Equal(decode(t, b), doc) {
			t.Error("the attachment is not the document")
		}
		if strings.HasPrefix(p.Header.Get("Content-Type"), "text/plain") && string(decode(t, b)) != "Dear acme,\r\nthe invoice is attached." {
			t.Errorf("text %q", decode(t, b))
		}
	}
	if len(parts) != 2 {
		t.Errorf("parts %v", parts)
	}
}

// decode reads a part in base64.
func decode(t *testing.T, b []byte) []byte {
	t.Helper()
	b = bytes.ReplaceAll(bytes.ReplaceAll(b, []byte("\r"), nil), []byte("\n"), nil)
	out, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(b)))
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// A recipient the server refuses (550) fails the send for good; a
// credential not there is tried again; the settings are checked.
func TestMailFails(t *testing.T) {
	srv := flowtest.NewSMTP(t)
	srv.Refuse = "nobody@example.com"
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "refused", Kind: smtp.KindSMTP, URL: "smtp://" + srv.Addr, From: "d@example.com", Recipients: []string{"nobody@example.com"}},
		docflow.Node{Name: "nocred", Kind: smtp.KindSMTP, URL: "smtp://" + srv.Addr, Credential: "missing", From: "d@example.com", Recipients: []string{"a@example.com"}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"refused", "nocred"}})
	j := e.Submit("in", "x", flowtest.Doc(t, 1), nil)
	e.Tick()
	if d := e.Delivery(j, "refused"); d.State != docflow.DeliveryFailed || !strings.Contains(d.LastError, "550") {
		t.Errorf("refused: %s %s", d.State, d.LastError)
	}
	if d := e.Delivery(j, "nocred"); d.State != docflow.DeliveryPending || !strings.Contains(d.LastError, "DOCFLOW_CRED_missing") {
		t.Errorf("no credential: %s %s", d.State, d.LastError)
	}
	ctx := context.Background()
	for _, n := range []docflow.Node{
		{Name: "x", Kind: smtp.KindSMTP, URL: "http://" + srv.Addr, From: "a@b", Recipients: []string{"c@d"}},
		{Name: "x", Kind: smtp.KindSMTP, URL: "smtp://" + srv.Addr, Recipients: []string{"c@d"}},
		{Name: "x", Kind: smtp.KindSMTP, URL: "smtp://" + srv.Addr, From: "a@b"},
		{Name: "x", Kind: smtp.KindSMTP, URL: "smtp://" + srv.Addr, From: "a@b", Recipients: []string{"c@d"}, Credential: "../etc/passwd"},
	} {
		if _, err := e.M().PutNode(ctx, n); !errors.Is(err, docflow.ErrInvalid) {
			t.Errorf("%+v: %v", n, err)
		}
	}
}
