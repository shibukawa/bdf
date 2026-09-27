package epub

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/webdoc"
)

// maxPart limits the size of a file of the publication that is read.
const maxPart = 64 << 20

const nsDC = "http://purl.org/dc/elements/1.1/"

// fontObfuscation are the algorithms of encryption.xml that only obfuscate
// fonts (which are not used: the text takes the reader style's fonts).
var fontObfuscation = map[string]bool{"http://www.idpf.org/2008/embedding": true, "http://ns.adobe.com/pdf/enc#RC": true}

// ErrDRM is returned for publications whose content is encrypted.
var ErrDRM = errors.New("epub: the publication is encrypted (DRM)")

// publication is an opened EPUB: its files and its package document.
type publication struct {
	files     map[string]*zip.File // by path, with forward slashes
	folded    map[string]string    // paths by their lower case, for references that differ in case
	encrypted map[string]bool

	opfPath string
	pkg     *packageDoc
	items   map[string]*item // manifest items by id
	byPath  map[string]*item
}

// packageDoc is the package document (the OPF file).
type packageDoc struct {
	Version  string   `xml:"version,attr"`
	UniqueID string   `xml:"unique-identifier,attr"`
	Metadata metaNode `xml:"metadata"`
	Manifest []item   `xml:"manifest>item"`
	Spine    struct {
		PPD      string    `xml:"page-progression-direction,attr"`
		Itemrefs []itemref `xml:"itemref"`
	} `xml:"spine"`
}

type item struct {
	ID         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	MediaType  string `xml:"media-type,attr"`
	Properties string `xml:"properties,attr"`
	Fallback   string `xml:"fallback,attr"`

	path string // in the container
}

func (it *item) has(prop string) bool { return hasToken(it.Properties, prop) }

type itemref struct {
	IDRef      string `xml:"idref,attr"`
	Linear     string `xml:"linear,attr"`
	Properties string `xml:"properties,attr"`
}

// metaNode is an element of the metadata, read generically: EPUB 2 and 3
// write it differently, and older files wrap it in dc-metadata.
type metaNode struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Text    string     `xml:",chardata"`
	Nodes   []metaNode `xml:",any"`
}

func (n *metaNode) attr(local string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == local {
			return strings.TrimSpace(a.Value)
		}
	}
	return ""
}

func hasToken(list, tok string) bool {
	for _, f := range strings.Fields(list) {
		if f == tok {
			return true
		}
	}
	return false
}

// open reads the container of a publication and its package document.
func open(r io.ReaderAt, size int64) (*publication, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("epub: %w", err)
	}
	p := &publication{files: map[string]*zip.File{}, folded: map[string]string{}, encrypted: map[string]bool{},
		items: map[string]*item{}, byPath: map[string]*item{}}
	for _, f := range zr.File {
		name := strings.TrimPrefix(strings.ReplaceAll(f.Name, `\`, "/"), "/")
		p.files[name] = f
		if _, ok := p.folded[strings.ToLower(name)]; !ok {
			p.folded[strings.ToLower(name)] = name
		}
	}
	if err := p.readEncryption(); err != nil {
		return nil, err
	}
	p.opfPath = p.rootfile()
	if p.opfPath == "" {
		return nil, errors.New("epub: no package document (META-INF/container.xml names none, and there is no .opf file)")
	}
	data, err := p.read(p.opfPath)
	if err != nil {
		return nil, fmt.Errorf("epub: package document: %w", err)
	}
	p.pkg = &packageDoc{}
	if err := decodeXML(data, p.pkg); err != nil {
		return nil, fmt.Errorf("epub: package document: %w", err)
	}
	for i := range p.pkg.Manifest {
		it := &p.pkg.Manifest[i]
		it.path = p.resolve(p.opfPath, it.Href)
		if it.ID != "" {
			p.items[it.ID] = it
		}
		if it.path != "" {
			p.byPath[it.path] = it
		}
	}
	return p, nil
}

func decodeXML(data []byte, v any) error {
	d := xml.NewDecoder(bytes.NewReader(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))))
	d.Strict = false
	d.Entity = xml.HTMLEntity
	d.CharsetReader = charsetReader
	return d.Decode(v)
}

// rootfile returns the path of the package document: the first rootfile
// of container.xml that is one, else the first .opf file.
func (p *publication) rootfile() string {
	if data, err := p.read("META-INF/container.xml"); err == nil {
		var c struct {
			Rootfiles []struct {
				FullPath  string `xml:"full-path,attr"`
				MediaType string `xml:"media-type,attr"`
			} `xml:"rootfiles>rootfile"`
		}
		if decodeXML(data, &c) == nil {
			for _, rf := range c.Rootfiles {
				if name := p.canonical(rf.FullPath); name != "" && (rf.MediaType == "" || rf.MediaType == "application/oebps-package+xml") {
					return name
				}
			}
		}
	}
	for name := range p.files {
		if strings.HasSuffix(strings.ToLower(name), ".opf") {
			return name
		}
	}
	return ""
}

// readEncryption reads which files META-INF/encryption.xml says are
// encrypted, other than obfuscated fonts.
func (p *publication) readEncryption() error {
	data, err := p.read("META-INF/encryption.xml")
	if err != nil {
		return nil
	}
	var enc struct {
		Data []struct {
			Method struct {
				Algorithm string `xml:"Algorithm,attr"`
			} `xml:"EncryptionMethod"`
			Ref struct {
				URI string `xml:"URI,attr"`
			} `xml:"CipherData>CipherReference"`
		} `xml:"EncryptedData"`
	}
	if err := decodeXML(data, &enc); err != nil {
		return nil
	}
	for _, d := range enc.Data {
		if fontObfuscation[d.Method.Algorithm] || d.Ref.URI == "" {
			continue
		}
		if name := p.canonical(p.resolve("", d.Ref.URI)); name != "" {
			p.encrypted[name] = true
		}
	}
	return nil
}

// canonical returns the path of a file of the container as it is stored
// ("" when there is none).
func (p *publication) canonical(name string) string {
	name = strings.TrimPrefix(name, "/")
	if _, ok := p.files[name]; ok {
		return name
	}
	return p.folded[strings.ToLower(name)]
}

// read reads a file of the container.
func (p *publication) read(name string) ([]byte, error) {
	f := p.files[p.canonical(name)]
	if f == nil {
		return nil, fmt.Errorf("%s: not in the publication", name)
	}
	name = p.canonical(name)
	if p.encrypted[name] {
		return nil, fmt.Errorf("%s: %w", name, ErrDRM)
	}
	if f.UncompressedSize64 > maxPart {
		return nil, fmt.Errorf("%s: larger than %d MiB", name, maxPart>>20)
	}
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	b, err := io.ReadAll(io.LimitReader(rc, maxPart+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxPart {
		return nil, fmt.Errorf("%s: larger than %d MiB", name, maxPart>>20)
	}
	return b, nil
}

// resolve resolves a reference (a URL, relative to the file base) to the
// path of a file in the container; "" for absolute URLs and references
// that leave the container.
func (p *publication) resolve(base, ref string) string {
	u, err := url.Parse(strings.TrimSpace(ref))
	if err != nil || u.Scheme != "" || u.Host != "" {
		return ""
	}
	name := u.Path
	if name == "" {
		return base
	}
	if !strings.HasPrefix(name, "/") {
		name = path.Join(path.Dir(base), name)
	}
	name = path.Clean(strings.TrimPrefix(name, "/"))
	if name == "." || name == ".." || strings.HasPrefix(name, "../") {
		return ""
	}
	return name
}

// metadata reads the Dublin Core metadata of the package, with the
// EPUB 3 refinements (the main title, roles) and the EPUB 2 attributes
// (opf:event of dates, opf:role).
func (p *publication) metadata() (dc bdf.DublinCore, props map[string]string) {
	props = map[string]string{} // rendition and other meta properties of the whole publication
	refines := map[string]map[string]string{}
	var nodes []*metaNode
	var walk func(n *metaNode)
	walk = func(n *metaNode) {
		for i := range n.Nodes {
			k := &n.Nodes[i]
			switch strings.ToLower(k.XMLName.Local) {
			case "dc-metadata", "x-metadata":
				walk(k)
				continue
			case "meta":
				text := strings.Join(strings.Fields(k.Text), " ")
				if prop := k.attr("property"); prop != "" {
					if ref := strings.TrimPrefix(k.attr("refines"), "#"); ref != "" {
						if refines[ref] == nil {
							refines[ref] = map[string]string{}
						}
						refines[ref][prop] = text
					} else if _, ok := props[prop]; !ok {
						props[prop] = text
					}
				} else if name := k.attr("name"); name != "" {
					if _, ok := props[name]; !ok {
						props[name] = k.attr("content")
					}
				}
				continue
			}
			nodes = append(nodes, k)
		}
	}
	walk(&p.pkg.Metadata)
	var titles, subtitles []string
	for _, n := range nodes {
		if n.XMLName.Space != nsDC && n.XMLName.Space != "" {
			continue
		}
		term := strings.ToLower(n.XMLName.Local)
		v := strings.Join(strings.Fields(n.Text), " ")
		if term == "description" && strings.ContainsRune(v, '<') {
			v = plainText(v)
		}
		if v == "" {
			continue
		}
		ref := refines[n.attr("id")]
		switch term {
		case "title":
			switch ref["title-type"] {
			case "", "main":
				titles = append(titles, v)
			default:
				subtitles = append(subtitles, v)
			}
		case "identifier":
			if n.attr("id") != "" && n.attr("id") == p.pkg.UniqueID {
				dc.Identifier = append(bdf.DCValues{v}, dc.Identifier...)
			} else {
				dc.Identifier = append(dc.Identifier, v)
			}
		case "date":
			switch strings.ToLower(n.attr("event")) {
			case "modification":
				dc.Modified = append(dc.Modified, v)
			case "creation":
				dc.Created = append(dc.Created, v)
			default:
				dc.Date = append(dc.Date, v)
			}
		default:
			if f := dc.Field(term); f != nil && term != "created" && term != "modified" {
				*f = append(*f, v)
			}
		}
	}
	dc.Title = append(titles, subtitles...)
	if m := props["dcterms:modified"]; m != "" && len(dc.Modified) == 0 {
		dc.Modified = bdf.DCValues{m}
	}
	return dc, props
}

// plainText returns the text of an HTML fragment.
func plainText(s string) string {
	doc, err := webdoc.ParseHTML([]byte(s), "text/html; charset=utf-8")
	if err != nil {
		return s
	}
	return strings.Join(strings.Fields(textOf(doc)), " ")
}
