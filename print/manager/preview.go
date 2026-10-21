package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// A printer's preview is a fake printer of tools/fakeprinter: the farm's
// API tells the jobs the printer received and the pages it read back
// from them, and the manager shows them beside the node, so that what a
// printer got is seen on the topology while developing and demonstrating.
// The preview is the printer's URL on the farm's API:
// http://host:port/api/printers/<printer>.

// PreviewJob is a job a fake printer received, as its farm lists it.
type PreviewJob struct {
	ID      int32     `json:"id"`
	Name    string    `json:"name,omitempty"`
	User    string    `json:"user,omitempty"`
	Format  string    `json:"format,omitempty"`
	State   string    `json:"state,omitempty"`
	Created time.Time `json:"created"`
	// Pages counts the pages read back, of the latest job.
	Pages int `json:"pages,omitempty"`
}

// previewClient fetches from a farm: a few seconds, a bounded answer.
var previewClient = &http.Client{Timeout: 10 * time.Second}

const maxPreviewBytes = 32 << 20

// checkPreviewURL takes the URL of a printer on a farm's API.
func checkPreviewURL(raw string) error {
	if raw == "" {
		return nil
	}
	if _, _, err := parsePreview(raw); err != nil {
		return err
	}
	return nil
}

// parsePreview splits a preview URL into the farm's origin and the
// printer's name.
func parsePreview(raw string) (origin, printer string, err error) {
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", "", fmt.Errorf("%w: preview %q: http://host:port/api/printers/<printer> of a fake printer farm", ErrInvalid, raw)
	}
	rest, ok := strings.CutPrefix(strings.TrimSuffix(u.Path, "/"), "/api/printers/")
	if !ok || !validName(rest) {
		return "", "", fmt.Errorf("%w: preview %q: the path is /api/printers/<printer>", ErrInvalid, raw)
	}
	return u.Scheme + "://" + u.Host, rest, nil
}

// previewOf finds a printer's preview.
func (m *Manager) previewOf(ctx context.Context, name string) (origin, printer string, err error) {
	p, err := m.printer(ctx, name)
	if err != nil {
		return "", "", err
	}
	if p.Preview == "" {
		return "", "", fmt.Errorf("%w: %s has no preview", ErrNotFound, name)
	}
	return parsePreview(p.Preview)
}

// previewGet fetches from the farm.
func previewGet(ctx context.Context, u string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	res, err := previewClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: the preview: %v", ErrPrint, err)
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("%w: the preview answered %d", ErrPrint, res.StatusCode)
	}
	return res, nil
}

// Preview lists the latest jobs a printer's fake printer received, the
// newest first (at most limit), with the pages read back of the newest.
func (m *Manager) Preview(ctx context.Context, name string, limit int) ([]PreviewJob, error) {
	origin, printer, err := m.previewOf(ctx, name)
	if err != nil {
		return nil, err
	}
	res, err := previewGet(ctx, origin+"/api/jobs?printer="+url.QueryEscape(printer))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var jobs []PreviewJob
	if err := json.NewDecoder(io.LimitReader(res.Body, maxPreviewBytes)).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("%w: the preview's jobs: %v", ErrPrint, err)
	}
	if limit <= 0 {
		limit = 5
	}
	if len(jobs) > limit {
		jobs = jobs[:limit]
	}
	if len(jobs) > 0 {
		res, err := previewGet(ctx, fmt.Sprintf("%s/api/printers/%s/jobs/%d/pages", origin, printer, jobs[0].ID))
		if err == nil {
			var pages struct {
				Pages []json.RawMessage `json:"pages"`
			}
			json.NewDecoder(io.LimitReader(res.Body, maxPreviewBytes)).Decode(&pages)
			res.Body.Close()
			jobs[0].Pages = len(pages.Pages)
		}
	}
	if jobs == nil {
		jobs = []PreviewJob{}
	}
	return jobs, nil
}

// WritePreviewPage writes a page (from 1) of a job of the printer's fake
// printer, as the farm read it back (PNG).
func (m *Manager) WritePreviewPage(ctx context.Context, name string, job int32, page int, w io.Writer) (string, error) {
	origin, printer, err := m.previewOf(ctx, name)
	if err != nil {
		return "", err
	}
	res, err := previewGet(ctx, fmt.Sprintf("%s/api/printers/%s/jobs/%d/pages/%d", origin, printer, job, page))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	ct := res.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		return "", fmt.Errorf("%w: the preview's page is %s", ErrPrint, ct)
	}
	_, err = io.Copy(w, io.LimitReader(res.Body, maxPreviewBytes))
	return ct, err
}

func (h *handler) preview(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	jobs, err := h.m.Preview(r.Context(), r.PathValue("name"), limit)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, jobs)
}

func (h *handler) previewPage(w http.ResponseWriter, r *http.Request) {
	job, err1 := strconv.Atoi(r.PathValue("job"))
	page, err2 := strconv.Atoi(strings.TrimSuffix(r.PathValue("page"), ".png"))
	if err1 != nil || err2 != nil || page < 1 {
		h.writeErr(w, ErrNotFound)
		return
	}
	// the page is fetched whole first, so that an error is an error
	var buf strings.Builder
	ct, err := h.m.WritePreviewPage(r.Context(), r.PathValue("name"), int32(job), page, &buf)
	if err != nil {
		if errors.Is(err, ErrPrint) {
			h.writeErr(w, fmt.Errorf("%w: %v", ErrNotFound, err))
			return
		}
		h.writeErr(w, err)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "private, max-age=60")
	io.WriteString(w, buf.String())
}
