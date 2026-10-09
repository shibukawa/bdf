package epub

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"
	"github.com/shibukawa/tinygodriver/encoding/xmlro/htmlentity"

	"github.com/shibukawa/bdf/converter/internal/ziputil"
	"github.com/shibukawa/bdf/internal/xmltree"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/webdoc"
)

// maxPart limits the size of a file of the publication that is read, and
// maxText that of a content document, a style sheet or the package
// document, which take many times their size when they are read.
var (
	maxPart int64 = 64 << 20
	maxText int64 = 16 << 20
)

// maxInflated is how much larger than the publication the files read from
// it may be together, decompressed; what is read after that is left out. A
// publication is a ZIP file: a small one can hold any number of files that
// decompress to the size they may have.
var maxInflated int64 = 256 << 20

// errInflated is the error of a file that is not read because the files
// read before it have reached maxInflated.
var errInflated = errors.New("the files of the publication are too large, decompressed")

// maxTags is how many tags the content documents of a publication may
// have together (some two million elements, several times those of a
// large dictionary); the content documents after them are left out. The
// elements are what reading a document takes memory for, a hundred times
// the size of the shortest of them, and the documents of a publication
// are all read before they are laid out.
var maxTags = 4 << 20

// errTags is the error of a content document that is not read because
// those read before it have reached maxTags.
var errTags = errors.New("the content documents of the publication have too many elements")

const nsDC = "http://purl.org/dc/elements/1.1/"

// fontObfuscation are the algorithms of encryption.xml that only obfuscate
// fonts (which are not used: the text takes the reader style's fonts).
var fontObfuscation = map[string]bool{"http://www.idpf.org/2008/embedding": true, "http://ns.adobe.com/pdf/enc#RC": true}

// ErrDRM is returned for publications whose content is encrypted.
var ErrDRM = errors.New("epub: the publication is encrypted (DRM)")

// publication is an opened EPUB: its files and its package document.
type publication struct {
	size      int64                // of the publication
	read      int64                // bytes read from its files, decompressed
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
	Version  string
	UniqueID string
	Metadata metaNode
	Manifest []item
	Spine    struct {
		PPD      string // page-progression-direction
		Itemrefs []itemref
	}
}

type item struct {
	ID         string
	Href       string
	MediaType  string
	Properties string
	Fallback   string

	path string // in the container
}

func (it *item) has(prop string) bool { return hasToken(it.Properties, prop) }

type itemref struct {
	IDRef      string
	Linear     string
	Properties string
}

// metaNode is an element of the XML files of a publication, read
// generically: the metadata above all, which EPUB 2 and 3 write
// differently, and older files wrap in dc-metadata.
type metaNode struct {
	XMLName xmltree.Name
	Attrs   []xmltree.Attr
	Text    string // the character data directly in the element
	Nodes   []metaNode
}

func (n *metaNode) attr(local string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == local {
			return strings.TrimSpace(a.Value)
		}
	}
	return ""
}

// value returns an attribute as it is written, whatever its namespace: the
// last of the name.
func (n *metaNode) value(local string) string {
	v := ""
	for _, a := range n.Attrs {
		if a.Name.Local == local {
			v = a.Value
		}
	}
	return v
}

// each calls f for the child elements with the name, whatever their
// namespace.
func (n *metaNode) each(local string, f func(k *metaNode)) {
	for i := range n.Nodes {
		if n.Nodes[i].XMLName.Local == local {
			f(&n.Nodes[i])
		}
	}
}

// maxNodes is how many elements an XML file of a publication other than a
// content document may have (the package document of a large publication
// has some thousands, an item for each of its files).
const maxNodes = 1 << 20

// readXML reads the root element of the package document, the container
// or the encryption file. It is lenient, and the entities of HTML are
// known, as package documents use them in their metadata.
func readXML(data []byte) (*metaNode, error) {
	r := xmltree.Open(data, xmlro.Options{Lenient: true, Entities: htmlentity.Lookup, CharsetReader: charsetReader})
	for {
		k, err := r.Next()
		if err != nil {
			return nil, err
		}
		switch k {
		case xmlro.EOF:
			return nil, io.ErrUnexpectedEOF
		case xmlro.StartElement:
			budget := maxNodes
			n, err := readNode(r, &budget)
			if err != nil {
				return nil, err
			}
			return &n, nil
		}
	}
}

// readNode reads the element whose start a reader is on. The elements
// within one another are bounded by the reader.
func readNode(r *xmlro.Reader, budget *int) (metaNode, error) {
	if *budget--; *budget < 0 {
		return metaNode{}, fmt.Errorf("more than %d elements", maxNodes)
	}
	n := metaNode{XMLName: xmltree.ElementName(r), Attrs: xmltree.Attrs(r)}
	var text []byte
	for {
		k, err := r.Next()
		if err != nil {
			return n, err
		}
		switch k {
		case xmlro.StartElement:
			c, err := readNode(r, budget)
			if err != nil {
				return n, err
			}
			n.Nodes = append(n.Nodes, c)
		case xmlro.Text, xmlro.CData:
			text = xmltree.AppendText(text, r)
		case xmlro.EndElement:
			n.Text = string(text)
			return n, nil
		case xmlro.EOF:
			return n, io.ErrUnexpectedEOF
		}
	}
}

// readPackage reads the package document: the metadata, the items of the
// manifest and the item references of the spine, whatever the namespace of
// their elements and wherever among the package's children they are.
func readPackage(data []byte) (*packageDoc, error) {
	root, err := readXML(data)
	if err != nil {
		return nil, err
	}
	pkg := &packageDoc{Version: root.value("version"), UniqueID: root.value("unique-identifier")}
	first := true
	root.each("metadata", func(m *metaNode) {
		if first {
			pkg.Metadata, first = *m, false
		} else {
			pkg.Metadata.Nodes = append(pkg.Metadata.Nodes, m.Nodes...)
		}
	})
	root.each("manifest", func(m *metaNode) {
		m.each("item", func(k *metaNode) {
			pkg.Manifest = append(pkg.Manifest, item{ID: k.value("id"), Href: k.value("href"), MediaType: k.value("media-type"),
				Properties: k.value("properties"), Fallback: k.value("fallback")})
		})
	})
	root.each("spine", func(sp *metaNode) {
		if ppd := sp.value("page-progression-direction"); ppd != "" {
			pkg.Spine.PPD = ppd
		}
		sp.each("itemref", func(k *metaNode) {
			pkg.Spine.Itemrefs = append(pkg.Spine.Itemrefs, itemref{IDRef: k.value("idref"), Linear: k.value("linear"), Properties: k.value("properties")})
		})
	})
	return pkg, nil
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
	zr, err := ziputil.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("epub: %w", err)
	}
	p := &publication{size: size, files: map[string]*zip.File{}, folded: map[string]string{}, encrypted: map[string]bool{},
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
	data, err := p.text(p.opfPath)
	if err != nil {
		return nil, fmt.Errorf("epub: package document: %w", err)
	}
	if p.pkg, err = readPackage(data); err != nil {
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

// rootfile returns the path of the package document: the first rootfile
// of container.xml that is one, else the first .opf file.
func (p *publication) rootfile() string {
	if data, err := p.text("META-INF/container.xml"); err == nil {
		if c, err := readXML(data); err == nil {
			found := ""
			c.each("rootfiles", func(rfs *metaNode) {
				rfs.each("rootfile", func(rf *metaNode) {
					if found != "" {
						return
					}
					if name, typ := p.canonical(rf.value("full-path")), rf.value("media-type"); name != "" && (typ == "" || typ == "application/oebps-package+xml") {
						found = name
					}
				})
			})
			if found != "" {
				return found
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
	data, err := p.text("META-INF/encryption.xml")
	if err != nil {
		return nil
	}
	enc, err := readXML(data)
	if err != nil {
		return nil
	}
	enc.each("EncryptedData", func(d *metaNode) {
		algorithm, uri := "", ""
		d.each("EncryptionMethod", func(m *metaNode) {
			if v := m.value("Algorithm"); v != "" {
				algorithm = v
			}
		})
		d.each("CipherData", func(c *metaNode) {
			c.each("CipherReference", func(ref *metaNode) {
				if v := ref.value("URI"); v != "" {
					uri = v
				}
			})
		})
		if fontObfuscation[algorithm] || uri == "" {
			return
		}
		if name := p.canonical(p.resolve("", uri)); name != "" {
			p.encrypted[name] = true
		}
	})
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

// picture reads a file of the container of at most maxPart.
func (p *publication) picture(name string) ([]byte, error) { return p.file(name, maxPart) }

// text reads a file of the container of at most maxText: a content
// document, a style sheet, the package document.
func (p *publication) text(name string) ([]byte, error) { return p.file(name, maxText) }

// file reads a file of the container that is not larger than limit, while
// the files read are not maxInflated larger than the publication
// (errInflated).
func (p *publication) file(name string, limit int64) ([]byte, error) {
	f := p.files[p.canonical(name)]
	if f == nil {
		return nil, fmt.Errorf("%s: not in the publication", name)
	}
	name = p.canonical(name)
	if p.encrypted[name] {
		return nil, fmt.Errorf("%s: %w", name, ErrDRM)
	}
	if f.UncompressedSize64 > uint64(limit) {
		return nil, fmt.Errorf("%s: larger than %d MiB", name, limit>>20)
	}
	// the size the file says is the most the ZIP reader gives
	if p.read+int64(f.UncompressedSize64) > p.size+maxInflated {
		return nil, fmt.Errorf("%s: %w", name, errInflated)
	}
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	b, err := io.ReadAll(io.LimitReader(rc, limit+1))
	p.read += int64(len(b))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("%s: larger than %d MiB", name, limit>>20)
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
