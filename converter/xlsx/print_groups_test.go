package xlsx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// The last page may contain only totals or notes: table titles are needed
// on intermediate detail pages, but must not follow the table's end.
func TestFillScopedPrintTitlesAndGroups(t *testing.T) {
	sheet := rowXML(1, "Invoice") + rowXML(2, "{{customer}}") +
		rowXML(3, "Item heading") + rowXML(4, "{{#items}}Line {{_index}}", "{{name}}{{/items}}") +
		rowXML(5, "") + rowXML(6, "Subtotal") + rowXML(7, "Tax") + rowXML(8, "Grand total") +
		rowXML(9, "") + rowXML(10, "Bank") + rowXML(11, "Account") +
		rowXML(12, "Notes") + rowXML(13, "Note begins\nNote ends") + rowXML(14, "") + rowXML(15, "Thanks") +
		`</sheetData><sheetFormatPr defaultRowHeight="20"/>` +
		`<mergeCells><mergeCell ref="A13:B14"/></mergeCells>` +
		`<pageMargins left="0" right="0" top="0" bottom="0" header="0" footer="0"/>` +
		`<pageSetup paperWidth="200pt" paperHeight="200pt" scale="100"/>`
	names := `<definedNames>` +
		`<definedName name="_xlnm.Print_Area" localSheetId="0">S1!$A$1:$B$15</definedName>` +
		`<definedName name="_xlnm.Print_Titles" localSheetId="0">S1!$3:$3</definedName>` +
		`<definedName name="_bdf.Print_TitleScope" localSheetId="0">S1!$A$3:$B$4</definedName>` +
		`<definedName name="_bdf.Print_KeepTogether" localSheetId="0">S1!$A$6:$B$8,S1!$A$10:$B$11,S1!$A$12:$B$15</definedName>` +
		`</definedNames>`
	book := book([]string{sheet}, names)
	for _, n := range []int{0, 1, 5, 8, 9, 10, 14, 15, 16, 25, 35} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			items := make([]any, n)
			for i := range items {
				items[i] = map[string]any{"name": "x"}
			}
			data := map[string]any{"customer": "Customer department", "items": items}
			res, err := Fill(bytes.NewReader(book), int64(len(book)), &TemplateData{Data: data}, testOptions())
			if err != nil {
				t.Fatal(err)
			}
			if len(res.Warnings) != 0 {
				t.Fatalf("warnings: %v", res.Warnings)
			}
			pages := pageText(t, res.Doc)["S1"]
			all := strings.Join(pages, "\n") + "\n"
			for _, once := range []string{"Invoice", "Customer department", "Subtotal", "Tax", "Grand total", "Bank", "Account", "Notes", "Note begins", "Note ends", "Thanks"} {
				if got := strings.Count(all, once); got != 1 {
					t.Errorf("%q appears %d times:\n%s", once, got, all)
				}
			}
			for i := 1; i <= n; i++ {
				if got := strings.Count(all, fmt.Sprintf("Line %d\n", i)); got != 1 {
					t.Errorf("item %d appears %d times", i, got)
				}
			}
			for i, page := range pages {
				if i > 0 && strings.Contains(page, "Item heading") != strings.Contains(page, "Line ") {
					t.Errorf("page %d repeats the table heading without details, or lacks it:\n%s", i+1, page)
				}
				for _, group := range [][]string{{"Subtotal", "Tax", "Grand total"}, {"Bank", "Account"}, {"Notes", "Note begins", "Note ends", "Thanks"}} {
					seen := 0
					for _, text := range group {
						if strings.Contains(page, text) {
							seen++
						}
					}
					if seen != 0 && seen != len(group) {
						t.Errorf("page %d splits %v:\n%s", i+1, group, page)
					}
				}
			}
			// Remapping must track the end of the expanded detail section.
			encoded, _ := json.Marshal(data)
			c, _, err := filled(t, book, string(encoded))
			if err != nil {
				t.Fatal(err)
			}
			if got := c.printNames[0].titleScope.r1; got != 2+n {
				t.Errorf("title scope ends at %d, want %d", got, 2+n)
			}
		})
	}
}

func TestPrintKeepGroups(t *testing.T) {
	for _, tc := range []struct {
		name   string
		groups [][2]int
		breaks []int
		want   []printBand
	}{
		{"move whole group", [][2]int{{3, 5}}, nil, []printBand{{0, 3}, {3, 8}}},
		{"exact fit", [][2]int{{2, 4}}, nil, []printBand{{0, 5}, {5, 8}}},
		{"oversized group", [][2]int{{1, 7}}, nil, []printBand{{0, 5}, {5, 8}}},
		{"manual break wins", [][2]int{{1, 3}}, []int{2}, []printBand{{0, 2}, {2, 7}, {7, 8}}},
		{"adjacent groups stay separate", [][2]int{{2, 3}, {4, 6}}, nil, []printBand{{0, 4}, {4, 8}}},
		{"overlapping groups join", [][2]int{{4, 5}, {2, 4}}, nil, []printBand{{0, 2}, {2, 7}, {7, 8}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := &printPlan{ps: &printSetup{rowBreaks: tc.breaks}, ph: 50, pw: 100, scale: 1,
				keepRows: joinKeepRows(tc.groups)}
			got, ok := p.cut(newAxis(8, 10, nil), 0, 7, true, 0)
			if !ok || !reflect.DeepEqual(got, tc.want) {
				t.Errorf("pages = %v, %v; want %v", got, ok, tc.want)
			}
		})
	}
}

func TestScopedHeadingStaysWithFirstDetail(t *testing.T) {
	intro := strings.Replace(rowXML(1, "Customer"), `<row r="1">`, `<row r="1" ht="80" customHeight="1">`, 1)
	sheet := intro + rowXML(2, "Heading") + rowXML(3, "{{#items}}{{.}}", "{{/items}}") + rowXML(4, "Total") +
		`</sheetData><sheetFormatPr defaultRowHeight="20"/>` +
		`<pageMargins left="0" right="0" top="0" bottom="0" header="0" footer="0"/>` +
		`<pageSetup paperWidth="200pt" paperHeight="100pt" scale="100"/>`
	names := `<definedNames>` +
		`<definedName name="_xlnm.Print_Area" localSheetId="0">S1!$A$1:$B$4</definedName>` +
		`<definedName name="_xlnm.Print_Titles" localSheetId="0">S1!$2:$2</definedName>` +
		`<definedName name="_bdf.Print_TitleScope" localSheetId="0">S1!$A$2:$B$3</definedName></definedNames>`
	data := book([]string{sheet}, names)
	res, err := Fill(bytes.NewReader(data), int64(len(data)), &TemplateData{Data: map[string]any{"items": []any{"AAA", "BBB", "CCC"}}}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	pages := pageText(t, res.Doc)["S1"]
	if len(pages) != 2 || strings.Contains(pages[0], "Heading") || strings.Count(pages[1], "Heading") != 1 || !strings.Contains(pages[1], "AAA") {
		t.Errorf("the table heading is orphaned or duplicated: %q", pages)
	}
}

func TestRemovedPrintTitleScope(t *testing.T) {
	sheet := rowXML(1, "Heading") + rowXML(2, "{{#items}}{{.}}", "{{/items}}") + rowXML(3, "Notes")
	names := `<definedNames>` +
		`<definedName name="_xlnm.Print_Titles" localSheetId="0">S1!$1:$1</definedName>` +
		`<definedName name="_bdf.Print_TitleScope" localSheetId="0">S1!$A$2:$B$2</definedName></definedNames>`
	c, _, err := filled(t, book([]string{sheet}, names), `{"items":[]}`)
	if err != nil {
		t.Fatal(err)
	}
	p := &printPlan{titleRows: c.printNames[0].titleRows, titleScope: c.printNames[0].titleScope}
	if p.titleScope == nil || p.rowTitles(1) {
		t.Error("removing all detail rows made table headings repeat without a scope")
	}
}
