package xmp

import (
	"slices"
	"testing"
)

func TestDublinCore(t *testing.T) {
	packet := `<?xpacket begin="` + "\ufeff" + `" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
 <rdf:Description rdf:about="" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmp:CreateDate="2026-09-26T12:00:00+09:00"
   xmlns:dc="http://purl.org/dc/elements/1.1/" dc:format="image/vnd.adobe.photoshop">
  <dc:title><rdf:Alt><rdf:li xml:lang="en">Poster</rdf:li><rdf:li xml:lang="x-default">ポスター</rdf:li></rdf:Alt></dc:title>
  <dc:creator><rdf:Seq><rdf:li>Ann</rdf:li><rdf:li>Bo</rdf:li></rdf:Seq></dc:creator>
  <dc:subject><rdf:Bag><rdf:li>red</rdf:li><rdf:li>blue</rdf:li></rdf:Bag></dc:subject>
  <dc:description><rdf:Alt><rdf:li xml:lang="en">A test</rdf:li></rdf:Alt></dc:description>
  <xmp:ModifyDate>2026-09-27T08:00:00Z</xmp:ModifyDate>
  <xmp:CreatorTool>ignored</xmp:CreatorTool>
 </rdf:Description>
</rdf:RDF></x:xmpmeta><?xpacket end="w"?>`
	dc := DublinCore([]byte(packet))
	for _, c := range []struct {
		name string
		got  []string
		want []string
	}{
		{"title", dc.Title, []string{"ポスター"}},
		{"creator", dc.Creator, []string{"Ann", "Bo"}},
		{"subject", dc.Subject, []string{"red", "blue"}},
		{"description", dc.Description, []string{"A test"}},
		{"format", dc.Format, []string{"image/vnd.adobe.photoshop"}},
		{"created", dc.Created, []string{"2026-09-26T12:00:00+09:00"}},
		{"modified", dc.Modified, []string{"2026-09-27T08:00:00Z"}},
	} {
		if !slices.Equal(c.got, c.want) {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
	if !DublinCore([]byte("<not xmp")).IsZero() {
		t.Error("malformed packet gave values")
	}
}
