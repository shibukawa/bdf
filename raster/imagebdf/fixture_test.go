package imagebdf

import (
	"testing"

	"github.com/shibukawa/bdf"
)

func openDoc(t testing.TB, path string) *bdf.Document {
	t.Helper()
	r, err := bdf.OpenSingleFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	d, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// Use fixture fonts so cache and parallel rendering tests do not depend on
// the machine's installed fonts.
var testFonts = []string{"../../converter/pptx/testdata/fonts", "../../converter/docx/testdata/fonts", "../../fixture/testdata/fonts"}
