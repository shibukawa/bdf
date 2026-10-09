package idml

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

// pack makes an IDML package of parts.
func pack(t *testing.T, parts map[string]string) *bytes.Reader {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, data := range parts {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(data)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(buf.Bytes())
}

const pkgNS = `xmlns:idPkg="http://ns.adobe.com/AdobeInDesign/idml/1.0/packaging"`

func designmap(body string) string {
	return `<?xml version="1.0" encoding="UTF-8"?><Document ` + pkgNS + ` DOMVersion="8.0" Self="d">` +
		`<Layer Self="ub3" Name="Layer 1" Visible="true"/>` + body + `</Document>`
}

func spreadPart(items string) string {
	return `<idPkg:Spread ` + pkgNS + `><Spread Self="ucc" BindingLocation="0" PageCount="1">` +
		`<Page Self="udb" AppliedMaster="n" GeometricBounds="0 0 200 100" ItemTransform="1 0 0 1 0 -100"/>` + items + `</Spread></idPkg:Spread>`
}

func rectItem(self, extra string) string {
	return `<Rectangle Self="` + self + `" ItemLayer="ub3" ItemTransform="1 0 0 1 0 0" FillColor="Color/Black" ` + extra + `>` +
		`<Properties><PathGeometry><GeometryPathType PathOpen="false"><PathPointArray>` +
		`<PathPointType Anchor="10 -90"/><PathPointType Anchor="10 -10"/><PathPointType Anchor="90 -10"/><PathPointType Anchor="90 -90"/>` +
		`</PathPointArray></GeometryPathType></PathGeometry></Properties></Rectangle>`
}

func convertParts(t *testing.T, parts map[string]string) (*Result, error) {
	t.Helper()
	r := pack(t, parts)
	return Convert(r, r.Size(), testOptions())
}

func TestMalformed(t *testing.T) {
	// not a document
	if _, err := convertParts(t, map[string]string{"designmap.xml": `<Spread/>`}); err == nil || !strings.Contains(err.Error(), "not an InDesign document") {
		t.Errorf("a spread as design map: %v", err)
	}
	// no pages
	if _, err := convertParts(t, map[string]string{"designmap.xml": designmap("")}); err == nil || !strings.Contains(err.Error(), "no pages") {
		t.Errorf("no pages: %v", err)
	}
	// a page without a usable size is skipped; a missing spread part warns
	res, err := convertParts(t, map[string]string{
		"designmap.xml": designmap(`<idPkg:Spread src="Spreads/missing.xml"/><idPkg:Spread src="Spreads/Spread_ucc.xml"/>`),
		"Spreads/Spread_ucc.xml": `<idPkg:Spread ` + pkgNS + `><Spread Self="ucc"><Page Self="u1" GeometricBounds="0 0 NaN 100"/>` +
			`<Page Self="u2" GeometricBounds="0 0 -5 100"/><Page Self="u3" GeometricBounds="0 0 100 100" ItemTransform="1 0 0 1 0 0"/></Spread></idPkg:Spread>`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Pages != 1 || len(res.Warnings) != 3 {
		t.Errorf("%d pages, warnings %q", res.Pages, res.Warnings)
	}
	// malformed XML in a story: the frame is empty, the page converts
	res, err = convertParts(t, map[string]string{
		"designmap.xml":          designmap(`<idPkg:Spread src="Spreads/Spread_ucc.xml"/><idPkg:Story src="Stories/Story_u100.xml"/>`),
		"Spreads/Spread_ucc.xml": spreadPart(strings.Replace(rectItem("u10", `ParentStory="u100"`), "Rectangle", "TextFrame", 2)),
		"Stories/Story_u100.xml": `<idPkg:Story><Story Self="u100"><ParagraphStyleRange><CharacterStyleRange><Content>x`,
	})
	if err != nil || res.Pages != 1 {
		t.Errorf("malformed story: %v, %d pages", err, res.Pages)
	}
}

func TestDeepAndCyclic(t *testing.T) {
	// groups nested beyond the limit are cut off, not followed forever
	deep := rectItem("u10", "")
	for i := 0; i < 200; i++ {
		deep = `<Group Self="g` + string(rune('a'+i%26)) + `" ItemLayer="ub3" ItemTransform="1 0 0 1 0 0">` + deep + `</Group>`
	}
	res, err := convertParts(t, map[string]string{
		"designmap.xml":          designmap(`<idPkg:Spread src="Spreads/Spread_ucc.xml"/>`),
		"Spreads/Spread_ucc.xml": spreadPart(deep),
	})
	if err != nil || res.Pages != 1 {
		t.Errorf("deep groups: %v", err)
	}
	// styles based on each other in a ring, a master applying itself, and
	// frames threaded in a ring
	res, err = convertParts(t, map[string]string{
		"designmap.xml": designmap(`<idPkg:Styles src="Resources/Styles.xml"/><idPkg:MasterSpread src="MasterSpreads/M.xml"/>` +
			`<idPkg:Spread src="Spreads/Spread_ucc.xml"/><idPkg:Story src="Stories/Story_u100.xml"/>`),
		"Resources/Styles.xml": `<idPkg:Styles><RootParagraphStyleGroup><ParagraphStyle Self="ParagraphStyle/A" PointSize="20">` +
			`<Properties><BasedOn type="object">ParagraphStyle/B</BasedOn></Properties></ParagraphStyle>` +
			`<ParagraphStyle Self="ParagraphStyle/B"><Properties><BasedOn type="object">ParagraphStyle/A</BasedOn></Properties></ParagraphStyle>` +
			`</RootParagraphStyleGroup></idPkg:Styles>`,
		"MasterSpreads/M.xml": `<idPkg:MasterSpread><MasterSpread Self="ub5" AppliedMaster="ub5" BindingLocation="0">` +
			`<Page Self="um1" GeometricBounds="0 0 200 100" ItemTransform="1 0 0 1 0 -100"/>` + rectItem("um2", "") + `</MasterSpread></idPkg:MasterSpread>`,
		"Spreads/Spread_ucc.xml": strings.Replace(spreadPart(
			strings.Replace(rectItem("u10", `ParentStory="u100" NextTextFrame="u11" PreviousTextFrame="u11"`), "Rectangle", "TextFrame", 2)+
				strings.Replace(rectItem("u11", `ParentStory="u100" NextTextFrame="u10" PreviousTextFrame="u10"`), "Rectangle", "TextFrame", 2)),
			`AppliedMaster="n"`, `AppliedMaster="ub5"`, 1),
		"Stories/Story_u100.xml": `<idPkg:Story><Story Self="u100"><ParagraphStyleRange AppliedParagraphStyle="ParagraphStyle/A">` +
			`<CharacterStyleRange><Content>` + strings.Repeat("ring ", 200) + `</Content></CharacterStyleRange></ParagraphStyleRange></Story></idPkg:Story>`,
	})
	if err != nil || res.Pages != 1 {
		t.Fatalf("rings: %v", err)
	}
	// the master drew its rectangle once, the ring of frames took the text
	// and overset the rest
	overset := false
	for _, w := range res.Warnings {
		overset = overset || strings.Contains(w, "overset")
	}
	if !overset {
		t.Errorf("warnings %q", res.Warnings)
	}
	p := res.Doc.Views[0].Pages[0]
	if len(p.Layers) != 2 {
		t.Errorf("layers %+v", p.Layers)
	}
}

func TestHugeValues(t *testing.T) {
	// coordinates and sizes at the edge of what the format holds
	res, err := convertParts(t, map[string]string{
		"designmap.xml": designmap(`<idPkg:Spread src="Spreads/Spread_ucc.xml"/>`),
		"Spreads/Spread_ucc.xml": `<idPkg:Spread ` + pkgNS + `><Spread Self="ucc"><Page Self="udb" AppliedMaster="n" GeometricBounds="0 0 1e9 1e9"/>` +
			`<Page Self="udc" AppliedMaster="n" GeometricBounds="0 0 100 100" ItemTransform="1e300 0 0 1e300 0 0"/>` +
			rectItem("u10", `StrokeWeight="1e30" StrokeColor="Color/Black" CornerRadius="1e308" TopLeftCornerOption="RoundedCorner" TopRightCornerOption="RoundedCorner" BottomLeftCornerOption="RoundedCorner" BottomRightCornerOption="RoundedCorner"`) +
			`</Spread></idPkg:Spread>`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Pages != 1 {
		t.Errorf("%d pages", res.Pages)
	}
}
