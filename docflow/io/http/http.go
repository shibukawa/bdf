// Package http is the HTTP output of docflow: the kind http, shown as
// HTTP (push), sends what reaches the node to a web service in a request
// (the export node, HTTP (pull), is where applications fetch documents
// instead). A program that imports the package has the kind;
// docflow/io/notify posts its webhooks with Do.
//
// The node's URL expands {{meta.<key>}} (each value URL-escaped), its
// Method is POST (the default) or PUT, and its Headers expand
// {{meta.<key>}} too, without line breaks. Its credential is the
// request's Authorization: a token sent as Bearer, or a whole value with
// its scheme ("Basic …"). Its Body is how the request carries the
// document and the sender's metadata:
//
//   - direct (the default): the document is the body, as it is, and the
//     metadata goes along as X-Meta-<key> headers, so that another
//     docflow's endpoint (POST /in/<node>, with an app's token as the
//     credential) takes it as its sender's: .bdf goes to another site as
//     it is (2026-10-21).
//   - multipart: the body is multipart/form-data, the metadata a JSON
//     object in the part meta, then the document in the part file, with
//     its name and type (2026-10-23).
//
// Every request has an Idempotency-Key, the delivery's id, the same on
// every try (a send is at least once). 2xx is done; 408, 429 and 5xx are
// tried again; any other answer, a redirect too, fails for good.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	nethttp "net/http"
	"net/textproto"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/io/smtp"
)

// KindHTTP is the kind of the HTTP output.
const KindHTTP = "http"

// The ways a request carries the document (Node.Body).
const (
	// BodyDirect sends the document as the body, the metadata in
	// X-Meta-* headers.
	BodyDirect = "direct"
	// BodyMultipart sends multipart/form-data: the metadata as JSON in
	// the part meta, the document in the part file.
	BodyMultipart = "multipart"
)

func init() {
	docflow.RegisterKind(docflow.Kind{Name: KindHTTP, Title: "HTTP (push)", Role: docflow.RoleExit, Check: check, Send: send})
}

// headerRE is the shape of a header's name.
var headerRE = regexp.MustCompile("^[A-Za-z0-9!#$%&'*+.^_`|~-]{1,64}$")

// reserved are the headers a node does not set: the request's own.
var reserved = map[string]bool{"Authorization": true, "Host": true, "Content-Length": true, "Idempotency-Key": true, "Transfer-Encoding": true, "Connection": true}

func check(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	n.URL, n.Credential = strings.TrimSpace(in.URL), strings.TrimSpace(in.Credential)
	if err := CheckURL(n.URL); err != nil {
		return err
	}
	if err := docflow.CheckCredential(n.Credential); err != nil {
		return err
	}
	n.Method = strings.ToUpper(strings.TrimSpace(in.Method))
	if n.Method == "" {
		n.Method = nethttp.MethodPost
	}
	if n.Method != nethttp.MethodPost && n.Method != nethttp.MethodPut {
		return fmt.Errorf("%w: method %q: POST or PUT", docflow.ErrInvalid, in.Method)
	}
	n.Body = strings.ToLower(strings.TrimSpace(in.Body))
	if n.Body == "" {
		n.Body = BodyDirect
	}
	if n.Body != BodyDirect && n.Body != BodyMultipart {
		return fmt.Errorf("%w: body %q: direct (the document as the body) or multipart (form-data)", docflow.ErrInvalid, in.Body)
	}
	for k, v := range in.Headers {
		k = nethttp.CanonicalHeaderKey(strings.TrimSpace(k))
		if !headerRE.MatchString(k) || reserved[k] {
			return fmt.Errorf("%w: header %q", docflow.ErrInvalid, k)
		}
		if n.Body == BodyMultipart && (k == "Content-Type" || k == "Content-Disposition") {
			return fmt.Errorf("%w: header %q: a multipart body sets its own", docflow.ErrInvalid, k)
		}
		if n.Headers == nil {
			n.Headers = map[string]string{}
		}
		n.Headers[k] = v
	}
	return nil
}

// CheckURL takes an http or https URL (that may expand {{meta.<key>}}).
func CheckURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("%w: url %q: an http(s) URL", docflow.ErrInvalid, raw)
	}
	return nil
}

// Client makes the requests (tests set their own). It follows no
// redirect.
var Client = &nethttp.Client{Timeout: 2 * time.Minute, CheckRedirect: func(*nethttp.Request, []*nethttp.Request) error { return nethttp.ErrUseLastResponse }}

// Request is a request Do sends.
type Request struct {
	Method, URL string
	Header      nethttp.Header
	Body        []byte
	// Credential names the Authorization's value (Node.Credential).
	Credential string
	// Key is the Idempotency-Key.
	Key string
}

// Do sends a request, and tells its status: 2xx is nil, 408, 429 and 5xx
// an error tried again, any other a final one.
func Do(ctx context.Context, m *docflow.Manager, r *Request) (string, error) {
	req, err := nethttp.NewRequestWithContext(ctx, r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil {
		return "", docflow.Final(err)
	}
	for k, v := range r.Header {
		req.Header[k] = v
	}
	req.Header.Set("User-Agent", "bdf-docflow")
	if r.Key != "" {
		req.Header.Set("Idempotency-Key", r.Key)
	}
	if r.Credential != "" {
		v, err := m.Credential(r.Credential)
		if err != nil {
			return "", err
		}
		if !strings.Contains(v, " ") {
			v = "Bearer " + v
		}
		req.Header.Set("Authorization", v)
	}
	res, err := Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
	switch {
	case res.StatusCode/100 == 2:
		return res.Status, nil
	case res.StatusCode == nethttp.StatusRequestTimeout || res.StatusCode == nethttp.StatusTooManyRequests || res.StatusCode >= 500:
		return "", fmt.Errorf("%s %s: %s: %s", r.Method, redact(r.URL), res.Status, strings.TrimSpace(string(b)))
	}
	return "", docflow.Final(fmt.Errorf("%s %s: %s: %s", r.Method, redact(r.URL), res.Status, strings.TrimSpace(string(b))))
}

// redact is a URL as an error tells it: without its query and user.
func redact(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "the URL"
	}
	u.User, u.RawQuery = nil, ""
	return u.String()
}

// send sends what reaches the node.
func send(ctx context.Context, o *docflow.Out) (string, error) {
	n := o.Node
	r, err := o.Open(ctx)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(r)
	r.Close()
	if err != nil {
		return "", err
	}
	h := nethttp.Header{}
	if n.Body == BodyMultipart {
		if body, err = multipartBody(o, body, h); err != nil {
			return "", docflow.Final(err)
		}
	} else {
		for k, v := range o.Meta {
			if !strings.HasPrefix(k, "$") && headerRE.MatchString("X-Meta-"+k) {
				h.Set("X-Meta-"+k, oneLine(docflow.MetaText(v)))
			}
		}
	}
	for k, v := range n.Headers {
		h.Set(k, o.ExpandHeader(v))
	}
	if n.Body != BodyMultipart {
		if h.Get("Content-Type") == "" {
			h.Set("Content-Type", smtp.ContentType(o.Payload.Format))
		}
		h.Set("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(smtp.AttachmentName(o), `"`, "")+`"`)
	}
	return Do(ctx, o.Manager(), &Request{Method: n.Method, URL: o.ExpandFunc(n.URL, url.QueryEscape), Header: h, Body: body, Credential: n.Credential, Key: o.ID})
}

// multipartBody wraps the document in multipart/form-data, after the
// sender's metadata: the part meta, a JSON object of it (the values as
// they were given), and the part file, the document with its name and
// type. It sets the request's Content-Type.
func multipartBody(o *docflow.Out, doc []byte, h nethttp.Header) ([]byte, error) {
	meta := map[string]any{}
	for k, v := range o.Meta {
		if !strings.HasPrefix(k, "$") {
			meta[k] = v
		}
	}
	js, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreatePart(textproto.MIMEHeader{"Content-Disposition": {`form-data; name="meta"`}, "Content-Type": {"application/json"}})
	if err != nil {
		return nil, err
	}
	part.Write(js)
	part, err = w.CreatePart(textproto.MIMEHeader{"Content-Disposition": {multipart.FileContentDisposition("file", oneLine(smtp.AttachmentName(o)))},
		"Content-Type": {smtp.ContentType(o.Payload.Format)}})
	if err != nil {
		return nil, err
	}
	part.Write(doc)
	if err := w.Close(); err != nil {
		return nil, err
	}
	h.Set("Content-Type", w.FormDataContentType())
	return b.Bytes(), nil
}

func oneLine(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, s)
}
