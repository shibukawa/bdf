package form

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func loadJSON(t *testing.T, name string) map[string]any {
	t.Helper()
	b, err := os.ReadFile("../test/xlsx/data/" + name + ".json")
	if err != nil {
		t.Fatal(err)
	}
	data, err := LoadData(Source{FileName: name + ".json", Data: b})
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func write(t *testing.T, res *Result) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := res.Doc.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// TestPartial fills a compiled template partly, and the partly filled form
// with the rest of the data: the form that the template filled with all of
// the data at once is.
func TestPartial(t *testing.T) {
	for _, file := range []string{"xlsx/testdata/form-fixed.xlsx", "docx/testdata/form-fixed.docx", "pptx/testdata/form-fixed.pptx"} {
		tpl, err := os.ReadFile("../converter/" + file)
		if err != nil {
			t.Fatal(err)
		}
		files := os.DirFS("../test/xlsx/data")
		comp, err := Compile(bytes.NewReader(tpl), int64(len(tpl)), "form-fixed", &Options{Convert: testConvert()})
		if err != nil {
			t.Fatal(err)
		}
		var bdft bytes.Buffer
		if err := comp.Doc.WriteSingle(&bdft); err != nil {
			t.Fatal(err)
		}
		all := loadJSON(t, "form-fixed")
		direct, err := Fill(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()), &Options{Data: all, Files: files, Convert: testConvert()})
		if err != nil {
			t.Fatal(err)
		}
		if direct.Partial != nil {
			t.Errorf("%s: a complete fill is partly filled", file)
		}

		// the sender fills what is theirs
		mine := loadJSON(t, "form-fixed")
		delete(mine, "customer")
		delete(mine, "memo")
		res, err := Fill(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()), &Options{Data: mine, Files: files, Convert: testConvert(), Open: []string{"customer.name", "memo"}})
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		for _, w := range res.Warnings {
			t.Errorf("%s: warning: %s", file, w)
		}
		p := res.Partial
		if p == nil || !reflect.DeepEqual(p.Open, []string{"customer.name", "memo"}) || p.From == "" || p.Name != "form-fixed" {
			t.Fatalf("%s: the partly filled form: %+v", file, p)
		}
		if photo, _ := p.Filled["photo"].(string); !strings.HasPrefix(photo, "data:image/png;base64,") {
			t.Errorf("%s: the picture is not held in the form: %.40q", file, photo)
		}
		if _, ok := p.Filled["customer"]; ok {
			t.Errorf("%s: the form holds a name it leaves open", file)
		}
		got := text(t, res.Doc)
		if !strings.Contains(got, "1024") || strings.Contains(got, "ACME") || strings.Contains(got, "Thank you") {
			t.Errorf("%s: the pages of the partly filled form:\n%s", file, got)
		}
		var schema struct {
			Properties map[string]struct {
				Properties map[string]any `json:"properties"`
			} `json:"properties"`
		}
		if err := json.Unmarshal(p.Schema, &schema); err != nil {
			t.Fatal(err)
		}
		if len(schema.Properties) != 2 || schema.Properties["customer"].Properties["name"] == nil || len(schema.Properties["customer"].Properties) != 1 {
			t.Errorf("%s: the schema of the partly filled form:\n%s", file, p.Schema)
		}
		bdff := write(t, res)
		if s, err := SchemaOf(bytes.NewReader(bdff), int64(len(bdff)), nil); err != nil || bytes.Contains(s, []byte(`"jan"`)) || !bytes.Contains(s, []byte(`"memo"`)) {
			t.Errorf("%s: SchemaOf the partly filled form: %v\n%s", file, err, s)
		}
		if _, err := Compile(bytes.NewReader(bdff), int64(len(bdff)), "", nil); err == nil || !strings.Contains(err.Error(), "partly filled") {
			t.Errorf("%s: compiling a partly filled form: %v", file, err)
		}

		// the other party fills in the rest
		theirs := map[string]any{"customer": map[string]any{"name": "ACME"}, "memo": "Thank you.\nSecond line"}
		done, err := Fill(bytes.NewReader(bdff), int64(len(bdff)), &Options{Data: theirs, Convert: testConvert()})
		if err != nil {
			t.Fatalf("%s: completing: %v", file, err)
		}
		if done.Partial != nil {
			t.Errorf("%s: the completed form is still partly filled", file)
		}
		if a, b := write(t, direct), write(t, done); !bytes.Equal(a, b) {
			t.Errorf("%s: the completed form is not the form filled at once (%d and %d bytes)", file, len(a), len(b))
		}

		// in two steps: a partly filled form of a partly filled form
		step, err := Fill(bytes.NewReader(bdff), int64(len(bdff)), &Options{Data: map[string]any{"memo": "Thank you.\nSecond line"}, Convert: testConvert(), Open: []string{OpenAll}})
		if err != nil || step.Partial == nil || !reflect.DeepEqual(step.Partial.Open, []string{"customer.name"}) || step.Partial.From != p.From {
			t.Fatalf("%s: a second step: %v %+v", file, err, step.Partial)
		}
		step2 := write(t, step)
		last, err := Fill(bytes.NewReader(step2), int64(len(step2)), &Options{Data: map[string]any{"customer": map[string]any{"name": "ACME"}}, Convert: testConvert()})
		if err != nil {
			t.Fatal(err)
		}
		if a, b := write(t, direct), write(t, last); !bytes.Equal(a, b) {
			t.Errorf("%s: the form completed in steps is not the form filled at once", file)
		}

		for _, c := range []struct {
			data map[string]any
			open []string
			want string
		}{
			{map[string]any{"no": 1, "customer": map[string]any{"name": "x"}, "memo": "x"}, nil, `"no" is filled already`},
			{map[string]any{"other": 1, "customer": map[string]any{"name": "x"}, "memo": "x"}, nil, `does not leave "other" open`},
			{map[string]any{"memo": "x"}, nil, `has no "customer.name", which the form leaves open`},
			{map[string]any{"memo": "x"}, []string{"no"}, `"no" is filled already`},
			{map[string]any{"customer": "x", "memo": "x"}, nil, `"customer" of the data is not an object`},
			{map[string]any{"memo": "x"}, []string{"customer.name", "memo"}, `gives "memo", which it leaves open`},
		} {
			_, err := Fill(bytes.NewReader(bdff), int64(len(bdff)), &Options{Data: c.data, Open: c.open, Convert: testConvert()})
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("%s: %v %v: %v (want %q)", file, c.data, c.open, err, c.want)
			}
		}
	}
}

// TestPartialNames checks the names that a partial fill leaves open.
func TestPartialNames(t *testing.T) {
	tpl, err := os.ReadFile("../converter/xlsx/testdata/form-rows.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	all := loadJSON(t, "form-rows")
	for _, c := range []struct {
		open []string
		drop []string
		want string
	}{
		{[]string{"items.name"}, nil, `"items.name" is not a name of the data`},
		{[]string{"nothing"}, nil, `"nothing" is not a name`},
		{[]string{"items"}, nil, `gives "items", which it leaves open`},
		{[]string{"items"}, []string{"items", "no"}, `has no "no"; give it`},
		{[]string{"items"}, []string{"items"}, ""},
	} {
		data := map[string]any{}
		for k, v := range all {
			data[k] = v
		}
		for _, k := range c.drop {
			delete(data, k)
		}
		res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &Options{Data: data, Open: c.open, Convert: testConvert()})
		if c.want == "" {
			if err != nil || res.Partial == nil || res.Partial.From != "" || res.Partial.Format != FormatXLSX {
				t.Errorf("%v: %v %+v", c.open, err, res)
			}
			continue
		}
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: %v (want %q)", c.open, err, c.want)
		}
	}
}

// TestPartialContract fills the part of one party of a contract with
// signature fields, and the other fills in theirs.
func TestPartialContract(t *testing.T) {
	tpl, err := os.ReadFile("../converter/pdfform/testdata/contract.pdf")
	if err != nil {
		t.Fatal(err)
	}
	sign := func(c string) string {
		svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 30"><path d="M5 25 L50 5 L95 25" fill="none" stroke="` + c + `" stroke-width="3"/></svg>`
		return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))
	}
	a := map[string]any{"name": "Taro Yamada", "sign": sign("#000")}
	b := map[string]any{"name": "ACME", "sign": sign("#00f")}
	direct, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &Options{Data: map[string]any{"party_a": a, "party_b": b}, Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &Options{Data: map[string]any{"party_a": a}, Open: []string{"party_b"}, Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range res.Warnings {
		t.Errorf("warning: %s", w)
	}
	if res.Partial == nil || res.Partial.Format != FormatPDF || !bytes.Contains(res.Partial.Schema, []byte(`"input": "signature"`)) || bytes.Contains(res.Partial.Schema, []byte(`party_a`)) {
		t.Fatalf("the partly filled contract: %+v\n%s", res.Partial, res.Partial.Schema)
	}
	bdff := write(t, res)
	done, err := Fill(bytes.NewReader(bdff), int64(len(bdff)), &Options{Data: map[string]any{"party_b": b}, Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	if x, y := write(t, direct), write(t, done); !bytes.Equal(x, y) {
		t.Errorf("the completed contract is not the contract filled at once (%d and %d bytes)", len(x), len(y))
	}
	if got := text(t, done.Doc); !strings.Contains(got, "Taro Yamada") || !strings.Contains(got, "ACME") {
		t.Errorf("the completed contract:\n%s", got)
	}
}
