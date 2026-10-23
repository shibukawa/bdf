// Package notify is the notify node of docflow: it tells that a job
// passed it, to Slack (an incoming webhook), by mail, or to a webhook of
// one's own (http), and sends the job on along its To, if any, with what
// reached it. Failures reach it along the on_error edges of the nodes
// (and the manager's default, Options.OnError), with $error and
// $error_node in the metadata it sees. A program that imports the
// package has the kind. (2026-10-23: the user chose a node in the graph,
// and on_error edges with a default of the whole.)
//
// What it tells is its Text, expanding {{meta.<key>}}; by default the
// job, the node and, of a failure, the error. It sends no payload. By
// slack, the credential is the webhook's URL (a secret); by http, the
// node's URL gets a JSON object (text, job, name, node, error,
// error_node, meta) with the credential as the Authorization; by mail,
// the settings are a mail's (docflow/io/smtp), without an attachment. A
// notification is sent at least once: a webhook gets the delivery's id as
// the Idempotency-Key, a mail as its Message-ID.
package notify

import (
	"context"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"strings"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/io/http"
	"github.com/shibukawa/bdf/docflow/io/smtp"
)

// KindNotify is the kind of the notify node.
const KindNotify = "notify"

// The ways a notify node tells.
const (
	ViaSlack = "slack"
	ViaMail  = "mail"
	ViaHTTP  = "http"
)

func init() {
	docflow.RegisterKind(docflow.Kind{Name: KindNotify, Title: "Notify", Role: docflow.RoleRoute, PassesOn: true, Check: check, Required: required, Step: step})
}

// required are the settings check refuses a node without: how it
// notifies, and what that needs.
var required = append([]docflow.Required{
	{Key: "via"},
	{Key: "credential", If: map[string]string{"via": ViaSlack}},
	{Key: "url", If: map[string]string{"via": ViaHTTP}},
}, smtp.Required(map[string]string{"via": ViaMail})...)

func check(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	n.To = docflow.CleanNames(in.To)
	n.Via, n.Text = strings.ToLower(strings.TrimSpace(in.Via)), in.Text
	switch n.Via {
	case ViaSlack:
		n.Credential = strings.TrimSpace(in.Credential)
		if n.Credential == "" {
			return fmt.Errorf("%w: a notify node by slack names the credential of its webhook's URL", docflow.ErrInvalid)
		}
		return docflow.CheckCredential(n.Credential)
	case ViaMail:
		text := n.Text
		if err := smtp.Check(ctx, m, in, n); err != nil {
			return err
		}
		n.Text = text
		return nil
	case ViaHTTP:
		n.URL, n.Credential = strings.TrimSpace(in.URL), strings.TrimSpace(in.Credential)
		if err := http.CheckURL(n.URL); err != nil {
			return err
		}
		return docflow.CheckCredential(n.Credential)
	}
	return fmt.Errorf("%w: via %q: slack, mail or http", docflow.ErrInvalid, in.Via)
}

// Text is what a notification says: the node's Text, expanded, or the
// default.
func Text(o *docflow.Out) string {
	if o.Node.Text != "" {
		return o.Expand(o.Node.Text)
	}
	what := fmt.Sprintf("Job %d", o.Job.Number)
	if o.Job.Name != "" {
		what += " (" + o.Job.Name + ")"
	}
	if o.Cause != nil {
		return fmt.Sprintf("%s failed at %s: %s", what, o.Cause.Node, o.Cause.Error)
	}
	return fmt.Sprintf("%s passed %s", what, o.Node.Name)
}

// step tells, and sends on what reached the node.
func step(ctx context.Context, o *docflow.Out) (*docflow.Made, error) {
	n, m, text := o.Node, o.Manager(), Text(o)
	switch n.Via {
	case ViaSlack:
		url, err := m.Credential(n.Credential)
		if err != nil {
			return nil, err
		}
		if err := http.CheckURL(url); err != nil {
			return nil, docflow.Final(fmt.Errorf("credential %s is no webhook's URL", n.Credential))
		}
		b, _ := json.Marshal(map[string]string{"text": text})
		_, err = http.Do(ctx, m, &http.Request{Method: nethttp.MethodPost, URL: url, Header: nethttp.Header{"Content-Type": {"application/json"}}, Body: b, Key: o.ID})
		return nil, err
	case ViaHTTP:
		claimed := map[string]any{}
		for k, v := range o.Meta {
			if !strings.HasPrefix(k, "$") {
				claimed[k] = v
			}
		}
		body := map[string]any{"text": text, "job": o.Job.Number, "name": o.Job.Name, "node": n.Name, "meta": claimed}
		if o.Cause != nil {
			body["error"], body["error_node"] = o.Cause.Error, o.Cause.Node
		}
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		_, err = http.Do(ctx, m, &http.Request{Method: nethttp.MethodPost, URL: n.URL, Header: nethttp.Header{"Content-Type": {"application/json"}}, Body: b, Credential: n.Credential, Key: o.ID})
		return nil, err
	case ViaMail:
		subject := o.ExpandHeader(n.MailSubject)
		if subject == "" {
			subject = strings.SplitN(text, "\n", 2)[0]
		}
		msg := &smtp.Mail{ID: o.ID, From: o.ExpandHeader(n.From), Subject: subject, Text: text}
		for _, r := range n.Recipients {
			msg.Recipients = append(msg.Recipients, o.ExpandHeader(r))
		}
		return nil, smtp.Send(ctx, m, n, msg)
	}
	return nil, docflow.Final(fmt.Errorf("via %q", n.Via))
}
