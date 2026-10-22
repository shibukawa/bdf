package pdfbdf

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/cms"
	"github.com/shibukawa/bdf/vector/internal/export"
)

// Signature signs a file as PAdES (ETSI EN 319 142-1) does: a signature
// field whose signature is a CMS signature of the bytes of the file around
// it (/SubFilter /ETSI.CAdES.detached), the baseline B-B, or B-T with a
// time-stamp of a TSA. Readers of PDF files (Acrobat …) check it.
type Signature struct {
	// Key signs, and Chain is its certificate followed by those of the
	// chain to the root.
	Key   crypto.Signer
	Chain []*x509.Certificate
	// Field is the name of a picture of the document whose box the
	// signature shows on: that of a signature field of a form (form,
	// tag.InputSignature), which a reader shows as the field of the
	// signature. "" makes a signature that shows on no page.
	Field string
	// Name, Reason, Location and Contact are what the signature
	// dictionary says of the signer.
	Name, Reason, Location, Contact string
	// Time is the time the signature says (/M; the clock when zero).
	Time time.Time
	// TSA is the URL of a time-stamp authority (RFC 3161), and Client the
	// client that asks it (nil: http.DefaultClient).
	TSA    string
	Client *http.Client
	// Certify makes the signature the certification of the file (DocMDP),
	// with the changes that it permits after it: 1 none, 2 filling in
	// forms and signing, 3 those and annotations. 0 is an approval
	// signature.
	Certify int
}

// sigFieldPrefix marks the rectangle of the picture of the signature field
// among the links of a stream: it goes to the pages as they do.
const sigFieldPrefix = "\x00bdf-signature-field:"

// sigPlace is where the signature field shows.
type sigPlace struct {
	pg   *pageOut
	rect string
}

// fieldRect notes the rectangle of the picture of the signature field.
func (ex *exporter) fieldRect(x, y, w, h float64) {
	s := ex.s
	m := s.toPage.Mul(s.st.ctm)
	x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, c := range [][2]float64{{x, y}, {x + w, y}, {x, y + h}, {x + w, y + h}} {
		px, py := m.Apply(c[0], c[1])
		x0, y0, x1, y1 = math.Min(x0, px), math.Min(y0, py), math.Max(x1, px), math.Max(y1, py)
	}
	if x1 > x0 && y1 > y0 && len(s.links) < maxLinks {
		s.links = append(s.links, link{x0: x0, y0: y0, x1: x1, y1: y1, url: sigFieldPrefix + ex.opts.Sign.Field})
	}
}

// placeholders of the signature dictionary, which sign fills in
const byteRangeHolder = "[0 0000000000 0000000000 0000000000]"

// sigObjects writes the objects of the signature: its dictionary, the
// widget of its field and the appearance of the widget. It returns the
// widget, and the entries of the catalog.
func (ex *exporter) sigObjects(place *sigPlace, pages []*pageOut) (widget ref, catalog string, err error) {
	sg := ex.opts.Sign
	size, err := sigSize(sg)
	if err != nil {
		return 0, "", err
	}
	pw := ex.pw
	at := sg.Time
	if at.IsZero() {
		at = time.Now()
	}
	date, _ := pdfDate(at.Format(time.RFC3339))
	sigRef := pw.alloc()
	dict := "<< /Type /Sig /Filter /Adobe.PPKLite /SubFilter /ETSI.CAdES.detached /ByteRange " + byteRangeHolder +
		" /Contents <" + strings.Repeat("0", 2*size) + "> /M " + text(date)
	for _, kv := range [][2]string{{"Name", sg.Name}, {"Reason", sg.Reason}, {"Location", sg.Location}, {"ContactInfo", sg.Contact}} {
		if kv[1] != "" {
			dict += " /" + kv[0] + " " + text(kv[1])
		}
	}
	if sg.Certify != 0 {
		if sg.Certify < 1 || sg.Certify > 3 {
			return 0, "", fmt.Errorf("pdfbdf: Certify is 1, 2 or 3, not %d", sg.Certify)
		}
		dict += fmt.Sprintf(" /Reference [<< /Type /SigRef /TransformMethod /DocMDP /TransformParams << /Type /TransformParams /P %d /V /1.2 >> >>]", sg.Certify)
	}
	pw.put(sigRef, dict+" >>")

	widget = pw.alloc()
	pg := pages[0]
	rect := "0 0 0 0"
	ap := ""
	if place != nil {
		pg, rect = place.pg, place.rect
		// (the picture of the signature is on the page: the appearance of
		// the widget draws nothing over it)
		r := strings.Fields(rect)
		x0, _ := strconv.ParseFloat(r[0], 64)
		y0, _ := strconv.ParseFloat(r[1], 64)
		x1, _ := strconv.ParseFloat(r[2], 64)
		y1, _ := strconv.ParseFloat(r[3], 64)
		form := pw.addStream(" /Type /XObject /Subtype /Form /BBox ["+export.Nums(0, 0, x1-x0, y1-y0)+"] /Resources << >>", nil, false)
		ap = " /AP << /N " + form.String() + " >>"
	}
	t := sg.Field
	if t == "" {
		t = "Signature1"
	}
	tu := ""
	if sg.Name != "" {
		tu = " /TU " + text("Signature of "+sg.Name)
	}
	pw.put(widget, fmt.Sprintf("<< /Type /Annot /Subtype /Widget /FT /Sig /T %s%s /V %s /F 132 /Rect [%s] /P %s%s >>", text(t), tu, sigRef, rect, pg.ref, ap))
	pg.sigWidget = widget
	catalog = fmt.Sprintf(" /AcroForm << /Fields [%s] /SigFlags 3 >>", widget)
	if sg.Certify != 0 {
		catalog += " /Perms << /DocMDP " + sigRef.String() + " >>"
	}
	ex.sigRef = sigRef
	return widget, catalog, nil
}

// tsaReserve is the room kept for the time-stamp token of a TSA, whose
// certificates the token carries.
const tsaReserve = 16 << 10

// sigSize is the bytes kept for the DER of the signature.
func sigSize(sg *Signature) (int, error) {
	if sg.Key == nil || len(sg.Chain) == 0 {
		return 0, errors.New("pdfbdf: a signature needs a key and its certificate")
	}
	// a signature of nothing is as large as that of the file, but for the
	// length of the numbers of its signature value
	der, err := cms.Sign(nil, sg.Key, sg.Chain, &cms.SignOptions{Detached: true})
	if err != nil {
		return 0, fmt.Errorf("pdfbdf: %w", err)
	}
	n := len(der) + 64
	if sg.TSA != "" {
		n += tsaReserve
	}
	return n, nil
}

// sign fills in the byte range and the signature of a file written with
// the placeholders of sigObjects.
func sign(data []byte, sg *Signature) error {
	br := bytes.Index(data, []byte("/ByteRange "+byteRangeHolder))
	if br < 0 {
		return errors.New("pdfbdf: the signature dictionary is not in the file")
	}
	c := bytes.Index(data[br:], []byte("/Contents <"))
	if c < 0 {
		return errors.New("pdfbdf: the signature dictionary is not in the file")
	}
	start := br + c + len("/Contents ")
	end := bytes.IndexByte(data[start:], '>')
	if end < 0 {
		return errors.New("pdfbdf: the signature dictionary is not in the file")
	}
	end += start + 1
	// the byte range, in the room of its placeholder
	ranges := fmt.Sprintf("[0 %d %d %d]", start, end, len(data)-end)
	if len(ranges) > len(byteRangeHolder) {
		return errors.New("pdfbdf: the file is too large to sign")
	}
	ranges += strings.Repeat(" ", len(byteRangeHolder)-len(ranges))
	copy(data[br+len("/ByteRange "):], ranges)
	signed := make([]byte, 0, len(data)-(end-start))
	signed = append(append(signed, data[:start]...), data[end:]...)
	der, err := cms.Sign(signed, sg.Key, sg.Chain, &cms.SignOptions{Detached: true})
	if err != nil {
		return fmt.Errorf("pdfbdf: %w", err)
	}
	if sg.TSA != "" {
		sd, err := cms.Parse(der)
		if err != nil {
			return err
		}
		ts, err := cms.RequestTimeStamp(context.Background(), sg.Client, sg.TSA, sd.Signers[0].SignatureDigest(crypto.SHA256), crypto.SHA256)
		if err != nil {
			return fmt.Errorf("pdfbdf: %w", err)
		}
		if der, err = cms.AddUnsigned(der, cms.Attribute{Type: cms.OIDTimeStampToken, Values: ts.DER}); err != nil {
			return fmt.Errorf("pdfbdf: %w", err)
		}
	}
	h := hex.EncodeToString(der)
	if len(h) > end-start-2 {
		return fmt.Errorf("pdfbdf: the signature (%d bytes) is larger than the room kept for it", len(der))
	}
	copy(data[start+1:], h)
	return nil
}

// markField notes that the picture after a MARK is that of the signature
// field.
func (ex *exporter) markField(kind byte, payload string) {
	if sg := ex.opts.Sign; sg != nil && sg.Field != "" && kind == bdf.MarkFigure && payload == sg.Field {
		ex.sigNext = true
	}
}
