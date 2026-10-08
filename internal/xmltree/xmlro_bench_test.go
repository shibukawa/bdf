package xmltree

import (
	"bytes"
	"encoding/xml"
	"strconv"
	"testing"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"
)

// genBody builds a Word document body the way Word writes one: paragraphs
// of runs with properties, some text with entities, no whitespace.
func genBody(paras int) []byte {
	var b bytes.Buffer
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\r\n")
	b.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" mc:Ignorable="w14"><w:body>`)
	for i := range paras {
		b.WriteString(`<w:p w14:paraId="` + strconv.Itoa(0x100000+i) + `" w:rsidR="00A1B2C3"><w:pPr><w:pStyle w:val="Normal"/><w:spacing w:after="120" w:line="276" w:lineRule="auto"/><w:jc w:val="both"/></w:pPr>`)
		for j := range 4 {
			b.WriteString(`<w:r w:rsidRPr="00A1B2C3"><w:rPr>`)
			if j%2 == 0 {
				b.WriteString(`<w:b/><w:sz w:val="24"/>`)
			} else {
				b.WriteString(`<w:i/><w:color w:val="1F497D"/>`)
			}
			b.WriteString(`</w:rPr><w:t xml:space="preserve">Paragraph ` + strconv.Itoa(i) + ` run ` + strconv.Itoa(j) + `, some words &amp; a few more words to make the text long </w:t></w:r>`)
		}
		if i%10 == 0 {
			b.WriteString(`<w:r><w:drawing><wp:inline xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"><wp:extent cx="914400" cy="914400"/><wp:docPr id="1" name="Picture 1"/></wp:inline></w:drawing></w:r>`)
		}
		b.WriteString(`</w:p>`)
	}
	b.WriteString(`<w:sectPr><w:pgSz w:w="11906" w:h="16838"/><w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440"/></w:sectPr></w:body></w:document>`)
	return b.Bytes()
}

var benchBody = genBody(4000)

// The token loops alone, with nothing kept: how fast each reader moves.

func BenchmarkBody_StdToken(b *testing.B) {
	b.SetBytes(int64(len(benchBody)))
	b.ReportAllocs()
	for range b.N {
		d := xml.NewDecoder(bytes.NewReader(benchBody))
		d.Strict = false
		for {
			if _, err := d.Token(); err != nil {
				break
			}
		}
	}
}

func BenchmarkBody_XmlroNext(b *testing.B) {
	b.SetBytes(int64(len(benchBody)))
	b.ReportAllocs()
	for range b.N {
		r := xmlro.NewBytesReader(benchBody, xmlro.Options{})
		for {
			k, err := r.Next()
			if err != nil {
				b.Fatal(err)
			}
			if k == xmlro.EOF {
				break
			}
		}
	}
}

// The trees: as encoding/xml built them, and as they are built now.

func BenchmarkBody_StdTree(b *testing.B) {
	b.SetBytes(int64(len(benchBody)))
	b.ReportAllocs()
	for range b.N {
		if _, err := stdTree(benchBody, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBody_Tree(b *testing.B) {
	b.SetBytes(int64(len(benchBody)))
	b.ReportAllocs()
	for range b.N {
		if _, err := Parse(benchBody); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBody_TreeFromReader(b *testing.B) {
	b.SetBytes(int64(len(benchBody)))
	b.ReportAllocs()
	for range b.N {
		if _, err := ParsePickingReader(bytes.NewReader(benchBody), nil); err != nil {
			b.Fatal(err)
		}
	}
}
