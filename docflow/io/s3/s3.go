// Package s3 is the S3 input and output of docflow: the kind s3 writes
// what reaches the node to a bucket (an exit), and the kind s3-trigger
// takes the objects an S3 event names (an entry). A program that imports
// the package has the kinds. It speaks the S3 API itself, signed with
// yuruidb's SigV4: no SDK.
//
// A node's URL is the bucket: s3://bucket?region=…&endpoint=…, with
// endpoint for a service other than AWS's (MinIO, rustfs, R2: path-style
// requests there). Its credential is access-key:secret-key, or
// access-key:secret-key:session-token; none takes AWS's own (the
// environment, the container's or the instance's role).
//
// The exit writes once and does not look back (2026-10-21: what must be
// kept goes to an archive): a PUT with If-None-Match: *, so that a retry
// after a write whose outcome was not known finds the object there (412)
// and is done. Its Path is the key, expanding {{meta.<key>}}; a value
// cannot leave the prefix (no .., no leading /).
//
// The trigger does not listen to a bucket: an S3 event goes to an
// EventBridge rule (an API destination) or a Lambda, which posts the
// object's reference to POST /in/<node> with the node's key (it is keyed
// as an endpoint is): {"bucket": …, "key": …}, the event of an S3
// notification (Records), or the event EventBridge sends (detail). The
// node reads the object with its credential, and the job goes on as from
// an endpoint, with s3.bucket and s3.key in its metadata. Its Path is
// the prefix of the keys it takes.
package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shibukawa/bdf/docflow"
	awsauth "github.com/shibukawa/yuruidb/auth/aws"
)

// The kinds of node of this package.
const (
	// KindS3 writes what reaches it to a bucket.
	KindS3 = "s3"
	// KindTrigger takes the objects S3 events name.
	KindTrigger = "s3-trigger"
)

// DefaultPath is the key an s3 exit writes by default: the job, and the
// payload's name.
const DefaultPath = "{{meta.$job}}/{{meta.$name}}"

func init() {
	docflow.RegisterKind(docflow.Kind{Name: KindS3, Title: "S3", Role: docflow.RoleExit, Check: checkS3, Send: send})
	docflow.RegisterKind(docflow.Kind{Name: KindTrigger, Title: "S3 Trigger", Role: docflow.RoleEntry, PassesOn: true, Keyed: true, Check: checkTrigger, Fetch: fetch})
}

func checkS3(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	if err := checkBucket(in, n); err != nil {
		return err
	}
	n.Path = strings.TrimSpace(in.Path)
	return nil
}

func checkTrigger(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	if err := checkBucket(in, n); err != nil {
		return err
	}
	n.Path = docflow.CleanPath(strings.TrimSpace(in.Path))
	return nil
}

func checkBucket(in docflow.Node, n *docflow.Node) error {
	n.URL, n.Credential = strings.TrimSpace(in.URL), strings.TrimSpace(in.Credential)
	if _, err := parse(n.URL); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	return docflow.CheckCredential(n.Credential)
}

// bucket is where a node's objects are.
type bucket struct {
	name, region string
	// base is the URL a key is put after: the endpoint and the bucket
	// (path-style), or the bucket's own host (AWS).
	base string
}

func parse(raw string) (*bucket, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "s3" || u.Host == "" || strings.ContainsAny(u.Host, "/\\?#") || u.Path != "" && u.Path != "/" {
		return nil, fmt.Errorf("url %q: s3://bucket?region=…&endpoint=…", raw)
	}
	q := u.Query()
	b := &bucket{name: u.Host, region: awsauth.Region(q.Get("region"))}
	if ep := q.Get("endpoint"); ep != "" {
		e, err := url.Parse(ep)
		if err != nil || (e.Scheme != "http" && e.Scheme != "https") || e.Host == "" {
			return nil, fmt.Errorf("endpoint %q: an http(s) URL", ep)
		}
		b.base = strings.TrimRight(ep, "/") + "/" + url.PathEscape(b.name)
	} else {
		if b.region == "" {
			return nil, fmt.Errorf("url %q: no region (?region=, or AWS_REGION)", raw)
		}
		b.base = "https://" + b.name + ".s3." + b.region + ".amazonaws.com"
	}
	if b.region == "" {
		b.region = "us-east-1"
	}
	return b, nil
}

// object is the URL of an object.
func (b *bucket) object(key string) string {
	parts := strings.Split(key, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return b.base + "/" + strings.Join(parts, "/")
}

// client is a node's way to its bucket.
type client struct {
	b     *bucket
	creds awsauth.Provider
	http  *http.Client
}

// HTTPClient makes the requests to S3 (tests set their own).
var HTTPClient = &http.Client{Timeout: 5 * time.Minute}

func newClient(m *docflow.Manager, n *docflow.Node) (*client, error) {
	b, err := parse(n.URL)
	if err != nil {
		return nil, docflow.Final(err)
	}
	c := &client{b: b, http: HTTPClient, creds: awsauth.Default(nil)}
	if n.Credential != "" {
		v, err := m.Credential(n.Credential)
		if err != nil {
			return nil, err
		}
		parts := strings.SplitN(v, ":", 3)
		if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
			return nil, docflow.Final(fmt.Errorf("credential %s: access-key:secret-key[:session-token]", n.Credential))
		}
		cr := awsauth.Credentials{AccessKey: parts[0], SecretKey: parts[1]}
		if len(parts) == 3 {
			cr.SessionToken = parts[2]
		}
		c.creds = awsauth.Static(cr)
	}
	return c, nil
}

// do sends a signed request.
func (c *client) do(ctx context.Context, method, key string, body []byte, header http.Header) (*http.Response, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.b.object(key), r)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	cr, err := c.creds(ctx)
	if err != nil {
		return nil, err
	}
	if err := awsauth.Sign(req, cr, c.b.region, "s3", time.Now()); err != nil {
		return nil, err
	}
	return c.http.Do(req)
}

// errorOf is the error of a response that is not a success: final for a
// refusal (4xx, but a timeout or too many requests), else tried again.
func errorOf(res *http.Response, what string) error {
	b, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
	err := fmt.Errorf("s3: %s: %s: %s", what, res.Status, strings.TrimSpace(string(b)))
	if res.StatusCode >= 400 && res.StatusCode < 500 && res.StatusCode != http.StatusRequestTimeout && res.StatusCode != http.StatusTooManyRequests {
		return docflow.Final(err)
	}
	return err
}

// send writes what reaches the node, once.
func send(ctx context.Context, o *docflow.Out) (string, error) {
	c, err := newClient(o.Manager(), o.Node)
	if err != nil {
		return "", err
	}
	key := Key(o)
	if key == "" {
		return "", docflow.Final(fmt.Errorf("the path %q is empty for this job", o.Node.Path))
	}
	r, err := o.Open(ctx)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(r)
	r.Close()
	if err != nil {
		return "", err
	}
	h := http.Header{"If-None-Match": {"*"}, "Content-Type": {contentType(o.Payload.Format)}}
	res, err := c.do(ctx, http.MethodPut, key, body, h)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	switch {
	case res.StatusCode/100 == 2:
	case res.StatusCode == http.StatusPreconditionFailed:
		// written before: by this delivery, whose outcome was not known
	default:
		return "", errorOf(res, "put "+key)
	}
	return "s3://" + c.b.name + "/" + key, nil
}

// Key is the key an s3 exit writes a payload under.
func Key(o *docflow.Out) string {
	path := o.Node.Path
	if path == "" {
		path = DefaultPath
	}
	return o.ExpandPath(path)
}

func contentType(format string) string {
	switch format {
	case "pdf":
		return "application/pdf"
	case "png":
		return "image/png"
	case "json":
		return "application/json"
	case "csv":
		return "text/csv"
	}
	return "application/octet-stream"
}
