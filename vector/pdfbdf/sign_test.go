package pdfbdf

import (
	"bytes"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/form"
	"github.com/shibukawa/bdf/internal/cms"
	"github.com/shibukawa/bdf/internal/cms/cmstest"
)

func contract(t *testing.T) *bdf.Document {
	t.Helper()
	b, err := os.ReadFile("../../testdata/form/pdf-contract.bdf")
	if err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func signedPDF(t *testing.T, doc *bdf.Document, o *Options) ([]byte, *Result) {
	t.Helper()
	o.FontDirs, o.NoSystemFonts = []string{"../../converter/pptx/testdata/fonts"}, true
	var b bytes.Buffer
	res, err := Write(&b, doc, o)
	if err != nil {
		t.Fatal(err)
	}
	return b.Bytes(), res
}

var byteRangeRE = regexp.MustCompile(`/ByteRange \[0 (\d+) (\d+) (\d+)\] */Contents <([0-9a-f]+)>`)

// verifyPDF checks the signature of a file as a reader of PDF files does:
// the byte range covers the file but the signature, whose CMS signs it.
func verifyPDF(t *testing.T, data []byte) (*cms.SignedData, error) {
	t.Helper()
	m := byteRangeRE.FindSubmatch(data)
	if m == nil {
		i := bytes.Index(data, []byte("/ByteRange"))
		t.Fatalf("no signature dictionary: %.300q", data[max(i, 0):])
	}
	b, _ := strconv.Atoi(string(m[1]))
	c, _ := strconv.Atoi(string(m[2]))
	d, _ := strconv.Atoi(string(m[3]))
	if c+d != len(data) || data[b] != '<' || data[c-1] != '>' {
		t.Fatalf("the byte range [0 %d %d %d] of a file of %d bytes", b, c, d, len(data))
	}
	der, err := hex.DecodeString(string(m[4]))
	if err != nil {
		t.Fatal(err)
	}
	sd, err := cms.Parse(der)
	if err != nil {
		t.Fatal(err)
	}
	if sd.Content != nil || sd.Signers[0].SigningTime() != (time.Time{}) {
		t.Error("the signature holds its content or a signing time (PAdES has the time in /M)")
	}
	signed := append(append([]byte{}, data[:b]...), data[c:]...)
	return sd, sd.Verify(sd.Signers[0], signed)
}

// TestSign signs a contract on the picture of the signature of party A.
func TestSign(t *testing.T) {
	ca := cmstest.NewAuthority("BDF test root")
	id := ca.Issue("Taro Yamada", false)
	at := time.Date(2026, 10, 23, 9, 30, 0, 0, time.FixedZone("JST", 9*3600))
	data, res := signedPDF(t, contract(t), &Options{Sign: &Signature{Key: id.Key, Chain: id.Chain, Field: "party_a.sign", Name: "Taro Yamada", Reason: "I agree", Location: "Tokyo", Time: at}})
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	sd, err := verifyPDF(t, data)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sd.VerifyChain(sd.Signers[0], ca.Pool(), time.Now()); err != nil {
		t.Error(err)
	}
	s := string(data)
	for _, want := range []string{"/SubFilter /ETSI.CAdES.detached", "/M (D:20261023093000+09'00')", "/Reason (I agree)", "/FT /Sig /T (party_a.sign)", "/AcroForm << /Fields [", "/SigFlags 3", "/F 132"} {
		if !strings.Contains(s, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	// the widget is on the box of the signature of party A (110 402 380 436)
	m := regexp.MustCompile(`/FT /Sig /T \(party_a\.sign\)[^>]* /Rect \[([0-9. -]+)\]`).FindStringSubmatch(s)
	if m == nil {
		t.Fatal("no widget of the field")
	}
	var r [4]float64
	for i, f := range strings.Fields(m[1]) {
		r[i], _ = strconv.ParseFloat(f, 64)
	}
	if r[0] < 110 || r[2] > 380 || r[1] < 402 || r[3] > 436 || r[2]-r[0] < 50 {
		t.Errorf("the widget is at %v", r)
	}
	if !strings.Contains(s, "/Annots [") {
		t.Error("the page does not list the widget")
	}

	// a byte changed outside the signature
	bad := bytes.Clone(data)
	i := bytes.Index(bad, []byte("/Reason (I agree)"))
	bad[i+9] = 'X'
	if _, err := verifyPDF(t, bad); err == nil {
		t.Error("a changed file verifies")
	}

	// invisible, certifying, with a time-stamp
	tsa := cmstest.NewTSA(ca)
	defer tsa.Close()
	data, res = signedPDF(t, contract(t), &Options{Sign: &Signature{Key: id.Key, Chain: id.Chain, Field: "nothing", TSA: tsa.URL, Certify: 2}})
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `no picture of "nothing"`) {
		t.Errorf("warnings: %v", res.Warnings)
	}
	sd, err = verifyPDF(t, data)
	if err != nil {
		t.Fatal(err)
	}
	ts, err := sd.Signers[0].TimeStampOf()
	if err != nil || ts == nil || time.Since(ts.Time) > time.Minute {
		t.Errorf("the time-stamp: %v %v", ts, err)
	}
	if dir := os.Getenv("BDF_PDF_SIGN_OUT"); dir != "" {
		// (for a check of the files by another validator: pyHanko)
		root, _ := (&cmstest.Identity{}).PEMOf(ca.Cert)
		os.WriteFile(filepath.Join(dir, "root.pem"), root, 0o644)
		os.WriteFile(filepath.Join(dir, "certified-timestamped.pdf"), data, 0o644)
	}
	s = string(data)
	for _, want := range []string{"/Rect [0 0 0 0]", "/TransformMethod /DocMDP", "/P 2", "/Perms << /DocMDP "} {
		if !strings.Contains(s, want) {
			t.Errorf("the certifying signature lacks %q", want)
		}
	}
	if _, err := Write(&bytes.Buffer{}, contract(t), &Options{Sign: &Signature{Key: id.Key, Chain: id.Chain, Certify: 4}}); err == nil {
		t.Error("Certify 4")
	}
	if _, err := Write(&bytes.Buffer{}, contract(t), &Options{Sign: &Signature{Key: id.Key}}); err == nil {
		t.Error("a signature without a certificate")
	}
}

// TestSignVeraPDF checks that a signed file is still PDF/A, with veraPDF
// when it is installed.
func TestSignVeraPDF(t *testing.T) {
	bin, err := exec.LookPath("verapdf")
	if err != nil {
		t.Skip("no verapdf")
	}
	ca := cmstest.NewAuthority("BDF test root")
	id := ca.Issue("Taro Yamada", true)
	for _, level := range []string{"2b", "2u", "2a"} {
		data, _ := signedPDF(t, contract(t), &Options{PDFA: level, Sign: &Signature{Key: id.Key, Chain: id.Chain, Field: "party_a.sign"}})
		if _, err := verifyPDF(t, data); err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(t.TempDir(), "signed.pdf")
		os.WriteFile(p, data, 0o644)
		out, _ := exec.Command(bin, "--flavour", level, "--format", "text", p).CombinedOutput()
		if !strings.Contains(string(out), "PASS") {
			t.Errorf("PDF/A-%s: %s", level, out)
		}
	}
}

// TestSignOnSheet shows the signature on the picture of a template of
// Excel: the name of its tag is that of the picture, in a tile of a sheet.
func TestSignOnSheet(t *testing.T) {
	tpl, err := os.ReadFile("../../converter/xlsx/testdata/form-fixed.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile("../../test/xlsx/data/form-fixed.json")
	if err != nil {
		t.Fatal(err)
	}
	data, err := form.LoadData(form.Source{FileName: "form-fixed.json", Data: b})
	if err != nil {
		t.Fatal(err)
	}
	data["photo"] = "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAzNjAgMTIwIj48cGF0aCBkPSJNMTAgMTAwIEwxODAgMjAgTDM1MCAxMDAiIGZpbGw9Im5vbmUiIHN0cm9rZT0iIzAwMCIvPjwvc3ZnPg=="
	res, err := form.Fill(bytes.NewReader(tpl), int64(len(tpl)), &form.Options{Data: data, Convert: &converter.Options{FontDirs: []string{"../../converter/pptx/testdata/fonts"}, NoSystemFonts: true}})
	if err != nil {
		t.Fatal(err)
	}
	ca := cmstest.NewAuthority("BDF test root")
	id := ca.Issue("Taro Yamada", false)
	pdf, wr := signedPDF(t, res.Doc, &Options{Sign: &Signature{Key: id.Key, Chain: id.Chain, Field: "photo"}})
	if len(wr.Warnings) != 0 {
		t.Errorf("warnings: %v", wr.Warnings)
	}
	if _, err := verifyPDF(t, pdf); err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(pdf, []byte("/Rect [0 0 0 0]")) || !bytes.Contains(pdf, []byte("/FT /Sig /T (photo)")) {
		t.Error("the signature does not show on the picture")
	}
}
