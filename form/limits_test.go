package form

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form/tag"
)

func TestSchemaLimits(t *testing.T) {
	for _, tc := range []struct {
		name   string
		fields []tag.Field
		want   string
	}{
		{"dotted name", []tag.Field{{Name: strings.Repeat("a.", maxSchemaDepth) + "value"}}, "nested"},
		{"section stack", []tag.Field{{Name: "value", In: make([]string, maxSchemaDepth+1)}}, "nested"},
		{"combined path", []tag.Field{{Name: "value", In: []string{strings.Repeat("a.", maxSchemaDepth-1) + "a"}}}, "nested"},
		{"field count", make([]tag.Field, maxSchemaNodes+1), "fields"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Schema(tc.fields, nil); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Schema = %v, want %q", err, tc.want)
			}
		})
	}
	// Many individually shallow fields still need a total node budget.
	fields := make([]tag.Field, maxSchemaNodes/2)
	for i := range fields {
		fields[i] = tag.Field{Name: fmt.Sprintf("object%d.value", i)}
	}
	if _, err := Schema(fields, nil); err == nil || !strings.Contains(err.Error(), "nodes") {
		t.Errorf("Schema = %v, want a node-limit error", err)
	}
}

func TestCompileSize(t *testing.T) {
	data := []byte("%PDF-1.7\n%%EOF")
	for _, tc := range []struct {
		name string
		size int64
		want string
	}{
		{"negative", -1, "negative template size"},
		{"too large", MaxTemplateSize + 1, "larger"},
		{"truncated", int64(len(data)) + 1, "unexpected EOF"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Compile(bytes.NewReader(data), tc.size, "template.pdf", nil)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Compile = %v, want %q", err, tc.want)
			}
		})
	}
}

// The manifest gives both a stored and a decoded size. A small decoded
// size must not allow Reader.Part to allocate a much larger stored part.
func TestReadTemplatePartLimits(t *testing.T) {
	src := []byte("template")
	sourceHash := bdf.HashOf(src)
	def, err := json.Marshal(Template{Version: TemplateVersion, Format: FormatPDF, Source: sourceHash.String()})
	if err != nil {
		t.Fatal(err)
	}
	desc := bdf.PartEntry{H: bdf.HashOf(def), T: bdf.PartTemplate, Enc: bdf.EncIdentity, Len: len(def), Size: len(def)}
	source := bdf.PartEntry{H: sourceHash, T: bdf.PartSource, Enc: bdf.EncIdentity, Len: len(src), Size: len(src), Off: int64(len(def))}
	for _, tc := range []struct {
		name   string
		mutate func(*bdf.PartEntry, *bdf.PartEntry)
		want   string
	}{
		{"description stored", func(d, s *bdf.PartEntry) { d.Enc, d.Len = bdf.EncDeflateRaw, maxTemplatePart+1 }, "larger"},
		{"description decoded", func(d, s *bdf.PartEntry) { d.Size = maxTemplatePart + 1 }, "larger"},
		{"description negative", func(d, s *bdf.PartEntry) { d.Size = -1 }, "negative size"},
		{"description mismatch", func(d, s *bdf.PartEntry) { d.Size-- }, "inconsistent"},
		{"source stored", func(d, s *bdf.PartEntry) { s.Enc, s.Len = bdf.EncDeflateRaw, MaxTemplateSize+1 }, "larger"},
		{"source decoded", func(d, s *bdf.PartEntry) { s.Size = MaxTemplateSize + 1 }, "larger"},
		{"source negative", func(d, s *bdf.PartEntry) { s.Len = -1 }, "negative size"},
		{"source mismatch", func(d, s *bdf.PartEntry) { s.Size-- }, "inconsistent"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, s := desc, source
			tc.mutate(&d, &s)
			data := compiledFile(t, []bdf.PartEntry{d, s}, append(bytes.Clone(def), src...))
			r := &recordingReaderAt{Reader: bytes.NewReader(data)}
			_, _, err := ReadTemplate(r, 1<<35)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("ReadTemplate = %v, want %q", err, tc.want)
			}
			if r.maxRead > 4096 {
				t.Errorf("a read of %d bytes was allocated before validating the part", r.maxRead)
			}
		})
	}
	t.Run("duplicate entry", func(t *testing.T) {
		oversized := desc
		oversized.Enc, oversized.Len = bdf.EncDeflateRaw, maxTemplatePart+1
		data := compiledFile(t, []bdf.PartEntry{desc, source, oversized}, def)
		r := &recordingReaderAt{Reader: bytes.NewReader(data)}
		_, _, err := ReadTemplate(r, 1<<35)
		if err == nil || !strings.Contains(err.Error(), "larger") || r.maxRead > 4096 {
			t.Errorf("ReadTemplate = %v; largest read %d", err, r.maxRead)
		}
	})
	t.Run("damaged description", func(t *testing.T) {
		changed := bytes.Replace(def, []byte("pdf"), []byte("bad"), 1)
		data := compiledFile(t, []bdf.PartEntry{desc, source}, append(changed, src...))
		_, _, err := ReadTemplate(bytes.NewReader(data), int64(len(data)))
		if err == nil || !strings.Contains(err.Error(), "damaged") {
			t.Errorf("ReadTemplate = %v, want a damaged-description error", err)
		}
	})
}

func compiledFile(t *testing.T, entries []bdf.PartEntry, parts []byte) []byte {
	t.Helper()
	m := bdf.Manifest{BDF: bdf.FormatVersion, Opset: bdf.OpsetVersion, Unit: "pt", Parts: entries}
	j, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	hdr := make([]byte, bdf.HeaderSize)
	copy(hdr, bdf.Magic[:])
	binary.LittleEndian.PutUint16(hdr[4:], bdf.FormatVersion)
	binary.LittleEndian.PutUint64(hdr[8:], bdf.HeaderSize)
	binary.LittleEndian.PutUint64(hdr[16:], uint64(len(j)))
	return append(append(hdr, j...), parts...)
}

type recordingReaderAt struct {
	*bytes.Reader
	maxRead int
}

func (r *recordingReaderAt) ReadAt(p []byte, off int64) (int, error) {
	r.maxRead = max(r.maxRead, len(p))
	if len(p) > 4096 {
		return 0, io.EOF
	}
	return r.Reader.ReadAt(p, off)
}
