package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/shibukawa/bdf/docflow"
)

// ref is an object an S3 event names.
type ref struct{ bucket, key string }

// event is any of the shapes a reference comes in.
type event struct {
	Bucket  string `json:"bucket"`
	Key     string `json:"key"`
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
	Detail *struct {
		Bucket struct {
			Name string `json:"name"`
		} `json:"bucket"`
		Object struct {
			Key string `json:"key"`
		} `json:"object"`
	} `json:"detail"`
}

// parseRef reads the reference of a request: {bucket, key}, an S3
// notification (whose key is URL-encoded, + a space), or EventBridge's
// event. One object a request.
func parseRef(b []byte) (ref, error) {
	var e event
	if err := json.Unmarshal(b, &e); err != nil {
		return ref{}, fmt.Errorf("%w: the body is no S3 event: %v", docflow.ErrInvalid, err)
	}
	switch {
	case e.Key != "":
		return ref{e.Bucket, e.Key}, nil
	case e.Detail != nil && e.Detail.Object.Key != "":
		return ref{e.Detail.Bucket.Name, e.Detail.Object.Key}, nil
	case len(e.Records) == 1:
		key, err := url.QueryUnescape(e.Records[0].S3.Object.Key)
		if err != nil {
			return ref{}, fmt.Errorf("%w: the key: %v", docflow.ErrInvalid, err)
		}
		return ref{e.Records[0].S3.Bucket.Name, key}, nil
	case len(e.Records) > 1:
		return ref{}, fmt.Errorf("%w: %d objects; one a request", docflow.ErrInvalid, len(e.Records))
	}
	return ref{}, fmt.Errorf("%w: the body names no object (bucket and key)", docflow.ErrInvalid)
}

// fetch reads the object a request names, from the node's bucket, under
// its prefix.
func fetch(ctx context.Context, f *docflow.Fetch) (*docflow.Fetched, error) {
	r, err := parseRef(f.Body)
	if err != nil {
		return nil, err
	}
	n := f.Node
	c, err := newClient(f.Manager(), n)
	if err != nil {
		return nil, err
	}
	if r.bucket != "" && r.bucket != c.b.name {
		return nil, fmt.Errorf("%w: bucket %q: %s takes the objects of %s", docflow.ErrForbidden, r.bucket, n.Name, c.b.name)
	}
	if docflow.CleanPath(r.key) != r.key {
		return nil, fmt.Errorf("%w: key %q", docflow.ErrInvalid, r.key)
	}
	if n.Path != "" && !strings.HasPrefix(r.key, n.Path+"/") && r.key != n.Path {
		return nil, fmt.Errorf("%w: key %q: %s takes the keys under %s/", docflow.ErrForbidden, r.key, n.Name, n.Path)
	}
	res, err := c.do(ctx, http.MethodGet, r.key, nil, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%w: %s/%s", docflow.ErrNotFound, c.b.name, r.key)
		}
		return nil, errorOf(res, "get "+r.key)
	}
	return &docflow.Fetched{Name: path.Base(r.key), Meta: map[string]any{"s3.bucket": c.b.name, "s3.key": r.key}, Body: res.Body}, nil
}
