package html

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

// mhtml is a web page saved as one MIME archive (RFC 2557): its HTML
// document and the resources it refers to.
type mhtml struct {
	html        []byte
	contentType string
	location    string            // the address the page was saved from
	parts       map[string][]byte // resources by Content-Location and cid: URL
}

// IsMHTML reports whether data starts like a web archive (MHTML): MIME
// headers with a multipart/related content type.
func IsMHTML(data []byte) bool {
	head := data[:min(len(data), 2048)]
	end := bytes.Index(head, []byte("\n\n"))
	if e := bytes.Index(head, []byte("\r\n\r\n")); e >= 0 && (end < 0 || e < end) {
		end = e
	}
	if end < 0 {
		return false
	}
	h := bytes.ToLower(head[:end])
	return bytes.Contains(h, []byte("mime-version:")) && bytes.Contains(h, []byte("multipart/related"))
}

// readMHTML reads an archive. The HTML document is the part the archive's
// start parameter names, else the first text/html part.
func readMHTML(data []byte) (*mhtml, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	mt, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil || !strings.HasPrefix(mt, "multipart/") || params["boundary"] == "" {
		return nil, errors.New("MHTML: not a multipart archive")
	}
	m := &mhtml{parts: map[string][]byte{}, location: strings.TrimSpace(msg.Header.Get("Snapshot-Content-Location"))}
	start := strings.Trim(params["start"], "<>")
	mr := multipart.NewReader(msg.Body, params["boundary"])
	found := false
	for {
		p, err := mr.NextRawPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			if found {
				break // a truncated archive: keep what was read
			}
			return nil, err
		}
		body, err := readPart(p)
		if err != nil {
			continue
		}
		ct := p.Header.Get("Content-Type")
		loc := strings.TrimSpace(p.Header.Get("Content-Location"))
		id := strings.Trim(strings.TrimSpace(p.Header.Get("Content-ID")), "<>")
		pt, _, _ := mime.ParseMediaType(ct)
		main := !found && (start != "" && id == start || start == "" && pt == "text/html")
		if main {
			m.html, m.contentType, found = body, ct, true
			if loc != "" && m.location == "" {
				m.location = loc
			}
			continue
		}
		if loc != "" {
			m.parts[loc] = body
		}
		if id != "" {
			m.parts["cid:"+id] = body
		}
	}
	if !found {
		return nil, errors.New("MHTML: no HTML document in the archive")
	}
	return m, nil
}

// readPart reads a part, undoing its transfer encoding.
func readPart(p *multipart.Part) ([]byte, error) {
	var r io.Reader = p
	switch strings.ToLower(strings.TrimSpace(p.Header.Get("Content-Transfer-Encoding"))) {
	case "base64":
		r = base64.NewDecoder(base64.StdEncoding, p) // line breaks are skipped
	case "quoted-printable":
		r = quotedprintable.NewReader(p)
	}
	return readLimited(r)
}
