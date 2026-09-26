package xlsx

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Conditional formats (ECMA-376 Part 1 §18.3.1.10) are evaluated here on the
// cached values: rules in priority order, the first rule that sets a
// property of the differential format wins it, stopIfTrue ends the
// evaluation for the cell. Color scales become fills; data bars and icons
// are drawn over the fills. Rules whose formulas need a formula engine
// (expression rules, references into other sheets) are skipped with a
// warning; constants and references to cells of the same sheet work.

type cfResults struct {
	cells map[[2]int]*cfResult
	rows  map[int][]int
}

type cfResult struct {
	font      *dxfFont
	fill      *xfill
	border    *xborder
	numFmt    *numFormat
	bar       *dataBar
	icon      *iconRef
	hideValue bool
	stopped   bool
}

type dataBar struct {
	lo, hi   float64 // fractions of the cell width the bar spans
	color    rgb
	border   *rgb
	gradient bool
}

type iconRef struct {
	set string
	idx int
}

func (r *cfResults) at(row, col int) *cfResult {
	if r == nil {
		return nil
	}
	return r.cells[[2]int{row, col}]
}

func (r *cfResults) colsIn(row int) []int {
	if r == nil {
		return nil
	}
	return r.rows[row]
}

func (res *cfResult) hasDxf() bool {
	return res.font != nil || res.fill != nil || res.border != nil || res.numFmt != nil || res.hideValue
}

// apply lays a result's differential format over a cell format.
func (res *cfResult) apply(f *cellFmt) {
	if d := res.font; d != nil {
		if d.bold != nil {
			f.font.bold = *d.bold
		}
		if d.italic != nil {
			f.font.italic = *d.italic
		}
		if d.strike != nil {
			f.font.strike = *d.strike
		}
		if d.underline != nil {
			f.font.underline = *d.underline
		}
		if d.color.kind != colorNone {
			f.font.color = d.color
		}
	}
	if res.fill != nil {
		f.fill = *res.fill
	}
	if b := res.border; b != nil {
		for _, p := range []struct{ dst, src *borderSide }{{&f.border.left, &b.left}, {&f.border.right, &b.right}, {&f.border.top, &b.top}, {&f.border.bottom, &b.bottom}} {
			if p.src.style != "" {
				*p.dst = *p.src
			}
		}
	}
	if res.numFmt != nil {
		f.numFmt = res.numFmt
	}
	if res.hideValue {
		f.hideText = true
	}
}

// merge takes the properties a higher-priority rule has not set.
func (res *cfResult) merge(d *dxf) {
	if d == nil {
		return
	}
	if d.font != nil {
		if res.font == nil {
			cp := *d.font
			res.font = &cp
		} else {
			f := res.font
			if f.bold == nil {
				f.bold = d.font.bold
			}
			if f.italic == nil {
				f.italic = d.font.italic
			}
			if f.strike == nil {
				f.strike = d.font.strike
			}
			if f.underline == nil {
				f.underline = d.font.underline
			}
			if f.color.kind == colorNone {
				f.color = d.font.color
			}
		}
	}
	if d.fill != nil && res.fill == nil && !d.fill.empty() {
		res.fill = d.fill
	}
	if d.border != nil && res.border == nil {
		res.border = d.border
	}
	if d.numFmt != nil && res.numFmt == nil {
		res.numFmt = d.numFmt
	}
}

type cfRule struct {
	n        *ooxml.Node
	ranges   []cellRange
	priority int
	x14      *ooxml.Node // the x14 extension of a data bar
}

// cellVal is a cell value as rules see it.
type cellVal struct {
	r, c int
	kind cellKind
	num  float64
	text string
}

func (v cellVal) isNum() bool { return v.kind == cellNum }

func (v cellVal) str() string {
	switch v.kind {
	case cellNum:
		return formatGeneral(v.num, 15)
	case cellBool:
		if v.num != 0 {
			return "TRUE"
		}
		return "FALSE"
	}
	return v.text
}

func (v cellVal) blank() bool {
	return v.kind == cellBlank || v.kind == cellStr && strings.TrimSpace(v.text) == ""
}

// evalConditionalFormats evaluates the conditional formats of the sheet.
func (s *sheetCtx) evalConditionalFormats() *cfResults {
	ws := s.ws
	if len(ws.cfs) == 0 && len(ws.x14cfs) == 0 {
		return nil
	}
	x14 := map[string]*ooxml.Node{}
	for _, cf := range ws.x14cfs {
		for _, r := range cf.Children("cfRule") {
			x14[r.AttrStr("id", "")] = r
		}
	}
	var rules []*cfRule
	extended := map[string]bool{}
	for _, cf := range ws.cfs {
		ranges := parseSqref(cf.AttrStr("sqref", ""))
		for _, rn := range cf.Children("cfRule") {
			rule := &cfRule{n: rn, ranges: ranges, priority: int(rn.AttrInt("priority", 1<<30))}
			for _, ext := range rn.Path("extLst").Children("ext") {
				if id := strings.TrimSpace(ext.Child("id").Content()); id != "" {
					rule.x14 = x14[id]
					extended[id] = true
				}
			}
			rules = append(rules, rule)
		}
	}
	// Excel 2010 rules of their own (icon sets with custom icons, ...),
	// not extensions of the rules above
	for _, cf := range ws.x14cfs {
		ranges := parseSqref(cf.Child("sqref").Content())
		for _, rn := range cf.Children("cfRule") {
			if !extended[rn.AttrStr("id", "")] {
				rules = append(rules, &cfRule{n: rn, ranges: ranges, priority: int(rn.AttrInt("priority", 1<<30))})
			}
		}
	}
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].priority < rules[j].priority })
	res := &cfResults{cells: map[[2]int]*cfResult{}, rows: map[int][]int{}}
	get := func(r, c int) *cfResult {
		k := [2]int{r, c}
		if x, ok := res.cells[k]; ok {
			return x
		}
		x := &cfResult{}
		res.cells[k] = x
		res.rows[r] = append(res.rows[r], c)
		return x
	}
	for _, rule := range rules {
		vals := s.rangeValues(rule.ranges)
		s.evalRule(rule, vals, get)
	}
	for r := range res.rows {
		sort.Ints(res.rows[r])
	}
	return res
}

// rangeValues returns the cells of ranges within the view (with values or
// not: blanks are values too for some rules).
func (s *sheetCtx) rangeValues(ranges []cellRange) []cellVal {
	var out []cellVal
	seen := map[[2]int]bool{}
	for _, rg := range ranges {
		r1, c1 := min(rg.r1, s.nRows-1), min(rg.c1, s.nCols-1)
		for r := rg.r0; r <= r1; r++ {
			rw := s.ws.rowAt(r)
			for c := rg.c0; c <= c1; c++ {
				k := [2]int{r, c}
				if seen[k] {
					continue
				}
				seen[k] = true
				v := cellVal{r: r, c: c}
				if rw != nil {
					if cl := s.ws.cellAt(r, c); cl != nil {
						v.kind, v.num = cl.kind, cl.num
						if cl.text != nil {
							v.text = cl.text.plain
						}
					}
				}
				out = append(out, v)
			}
			if rw == nil && rg.r1-rg.r0 > 1000 {
				// a huge range over empty rows: only rows that exist matter
				next := s.nextRowAfter(r)
				if next < 0 {
					break
				}
				r = next - 1
			}
		}
	}
	return out
}

// nextRowAfter returns the index of the first row element after r (-1).
func (s *sheetCtx) nextRowAfter(r int) int {
	i := sort.Search(len(s.ws.rows), func(i int) bool { return s.ws.rows[i].idx > r })
	if i < len(s.ws.rows) {
		return s.ws.rows[i].idx
	}
	return -1
}

func (s *sheetCtx) evalRule(rule *cfRule, vals []cellVal, get func(r, c int) *cfResult) {
	n := rule.n
	typ := n.AttrStr("type", "")
	var d *dxf
	if id := int(n.AttrInt("dxfId", -1)); id >= 0 && id < len(s.c.st.dxfs) {
		d = &s.c.st.dxfs[id]
	} else if dn := n.Child("dxf"); dn != nil {
		// Excel 2010 rules carry their format inline
		dd := s.c.st.parseDxf(dn)
		d = &dd
	}
	stop := n.AttrBool("stopIfTrue", false)
	nums := numbers(vals)
	origin := [2]int{}
	if len(rule.ranges) > 0 {
		origin = [2]int{rule.ranges[0].r0, rule.ranges[0].c0}
	}
	var match func(v cellVal) bool
	switch typ {
	case "cellIs":
		fs := ruleFormulas(n)
		op := n.AttrStr("operator", "equal")
		match = func(v cellVal) bool {
			var args []cellVal
			for _, f := range fs {
				a, ok := s.evalOperand(f.Content(), origin, v.r, v.c)
				if !ok {
					fv, ok := s.evalFormula(f.Content(), origin, v.r, v.c)
					if !ok {
						return false
					}
					a = cellVal{kind: cellBlank}
					switch fv.kind {
					case fNum:
						a = cellVal{kind: cellNum, num: fv.num}
					case fStr:
						a = cellVal{kind: cellStr, text: fv.str}
					case fBool:
						a = cellVal{kind: cellBool, num: fv.num}
					case fErr:
						return false
					}
				}
				args = append(args, a)
			}
			return compare(op, v, args)
		}
	case "expression":
		fs := ruleFormulas(n)
		if len(fs) == 0 {
			return
		}
		f := fs[0].Content()
		match = func(v cellVal) bool {
			fv, ok := s.evalFormula(f, origin, v.r, v.c)
			if !ok {
				return false
			}
			b, _ := truthy(fv)
			return b
		}
	case "containsText", "notContainsText", "beginsWith", "endsWith":
		text := strings.ToLower(n.AttrStr("text", ""))
		match = func(v cellVal) bool {
			str := strings.ToLower(v.str())
			switch typ {
			case "containsText":
				return strings.Contains(str, text)
			case "notContainsText":
				return !strings.Contains(str, text)
			case "beginsWith":
				return strings.HasPrefix(str, text)
			}
			return strings.HasSuffix(str, text)
		}
	case "containsBlanks":
		match = func(v cellVal) bool { return v.blank() }
	case "notContainsBlanks":
		match = func(v cellVal) bool { return !v.blank() }
	case "containsErrors":
		match = func(v cellVal) bool { return v.kind == cellErr }
	case "notContainsErrors":
		match = func(v cellVal) bool { return v.kind != cellErr }
	case "top10":
		rank := int(n.AttrInt("rank", 10))
		sorted := append([]float64(nil), nums...)
		sort.Float64s(sorted)
		if n.AttrBool("percent", false) {
			rank = max(1, int(float64(len(sorted))*float64(rank)/100))
		}
		if len(sorted) == 0 || rank <= 0 {
			return
		}
		rank = min(rank, len(sorted))
		bottom := n.AttrBool("bottom", false)
		thr := sorted[len(sorted)-rank]
		if bottom {
			thr = sorted[rank-1]
		}
		match = func(v cellVal) bool {
			if !v.isNum() {
				return false
			}
			if bottom {
				return v.num <= thr
			}
			return v.num >= thr
		}
	case "aboveAverage":
		if len(nums) == 0 {
			return
		}
		mean, sd := meanStd(nums)
		above := n.AttrBool("aboveAverage", true)
		eq := n.AttrBool("equalAverage", false)
		k := n.AttrFloat("stdDev", 0)
		thr := mean
		if k != 0 {
			if above {
				thr = mean + k*sd
			} else {
				thr = mean - k*sd
			}
		}
		match = func(v cellVal) bool {
			if !v.isNum() {
				return false
			}
			switch {
			case above && eq:
				return v.num >= thr
			case above:
				return v.num > thr
			case eq:
				return v.num <= thr
			}
			return v.num < thr
		}
	case "duplicateValues", "uniqueValues":
		count := map[string]int{}
		key := func(v cellVal) string {
			if v.isNum() {
				return "n" + strconv.FormatFloat(v.num, 'g', 15, 64)
			}
			return "s" + strings.ToLower(v.str())
		}
		for _, v := range vals {
			if !v.blank() {
				count[key(v)]++
			}
		}
		dup := typ == "duplicateValues"
		match = func(v cellVal) bool {
			if v.blank() {
				return false
			}
			return (count[key(v)] > 1) == dup
		}
	case "colorScale":
		s.colorScale(rule, vals, nums, get)
		return
	case "dataBar":
		s.dataBar(rule, vals, nums, get)
		return
	case "iconSet":
		s.iconSet(rule, vals, nums, get)
		return
	default:
		s.c.warnOnce("cftype:"+typ, "sheet %q: conditional formats of type %s are not evaluated", s.ws.name, typ)
		return
	}
	for _, v := range vals {
		res := get(v.r, v.c)
		if res.stopped || !match(v) {
			continue
		}
		res.merge(d)
		if stop {
			res.stopped = true
		}
	}
}

// ruleFormulas returns the formulas of a rule (formula, or xm:f in Excel
// 2010 rules).
func ruleFormulas(n *ooxml.Node) []*ooxml.Node {
	if fs := n.Children("formula"); len(fs) > 0 {
		return fs
	}
	return n.Children("f")
}

func numbers(vals []cellVal) []float64 {
	var out []float64
	for _, v := range vals {
		if v.isNum() {
			out = append(out, v.num)
		}
	}
	return out
}

func meanStd(v []float64) (float64, float64) {
	sum := 0.0
	for _, x := range v {
		sum += x
	}
	mean := sum / float64(len(v))
	if len(v) < 2 {
		return mean, 0
	}
	ss := 0.0
	for _, x := range v {
		ss += (x - mean) * (x - mean)
	}
	return mean, math.Sqrt(ss / float64(len(v)-1))
}

// evalFormula evaluates a rule's formula for the cell at r, c; false when
// the formula cannot be evaluated (with a warning, once per formula).
func (s *sheetCtx) evalFormula(f string, origin [2]int, r, c int) (fval, bool) {
	key := s.ws.name + "\x00" + f
	node, ok := s.formulas[key]
	if !ok {
		var err error
		node, err = parseFormula(f)
		if err != nil {
			node = nil
		}
		s.formulas[key] = node
	}
	if node != nil {
		env := &fenv{s: s, dr: r - origin[0], dc: c - origin[1], row: r, col: c}
		v, err := env.eval(node)
		if err == nil {
			return env.scalar(v), true
		}
		s.formulas[key] = nil
	}
	s.c.warnOnce("cf:"+f, "sheet %q: conditional format formula %q is not evaluated", s.ws.name, f)
	return fval{}, false
}

// evalOperand evaluates a formula of a cellIs rule: a number, a string, a
// boolean, or a reference to a cell of the sheet (relative references move
// with the cell, from the rule's first cell).
func (s *sheetCtx) evalOperand(f string, origin [2]int, r, c int) (cellVal, bool) {
	f = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(f), "="))
	if num, err := strconv.ParseFloat(f, 64); err == nil {
		return cellVal{kind: cellNum, num: num}, true
	}
	if len(f) >= 2 && f[0] == '"' && f[len(f)-1] == '"' {
		return cellVal{kind: cellStr, text: strings.ReplaceAll(f[1:len(f)-1], `""`, `"`)}, true
	}
	switch strings.ToUpper(f) {
	case "TRUE":
		return cellVal{kind: cellBool, num: 1}, true
	case "FALSE":
		return cellVal{kind: cellBool}, true
	}
	neg := false
	if strings.HasPrefix(f, "-") {
		neg, f = true, f[1:]
	}
	ref := f
	if strings.ContainsAny(ref, "!:( +*/&<>=") {
		return cellVal{}, false
	}
	col, row, ok := parseRef(ref)
	if !ok {
		return cellVal{}, false
	}
	// relative parts move with the evaluated cell
	colAbs := strings.HasPrefix(ref, "$")
	rowAbs := strings.Contains(strings.TrimPrefix(ref, "$"), "$")
	if !colAbs {
		col += c - origin[1]
	}
	if !rowAbs {
		row += r - origin[0]
	}
	v := cellVal{r: row, c: col}
	if cl := s.ws.cellAt(row, col); cl != nil {
		v.kind, v.num = cl.kind, cl.num
		if cl.text != nil {
			v.text = cl.text.plain
		}
	}
	if neg {
		if !v.isNum() {
			return cellVal{}, false
		}
		v.num = -v.num
	}
	return v, true
}

// compare applies a cellIs operator; strings compare case-insensitively
// and sort after numbers.
func compare(op string, v cellVal, args []cellVal) bool {
	cmp := func(a, b cellVal) int {
		an, bn := a.isNum() || a.kind == cellBlank, b.isNum() || b.kind == cellBlank
		switch {
		case an && bn:
			return sign(a.num - b.num)
		case an:
			return -1
		case bn:
			return 1
		}
		return strings.Compare(strings.ToLower(a.str()), strings.ToLower(b.str()))
	}
	if len(args) == 0 {
		return false
	}
	c0 := cmp(v, args[0])
	switch op {
	case "lessThan":
		return c0 < 0
	case "lessThanOrEqual":
		return c0 <= 0
	case "equal":
		return c0 == 0
	case "notEqual":
		return c0 != 0
	case "greaterThanOrEqual":
		return c0 >= 0
	case "greaterThan":
		return c0 > 0
	case "between", "notBetween":
		if len(args) < 2 {
			return false
		}
		lo, hi := args[0], args[1]
		if cmp(lo, hi) > 0 {
			lo, hi = hi, lo
		}
		in := cmp(v, lo) >= 0 && cmp(v, hi) <= 0
		return in == (op == "between")
	}
	return false
}

func sign(v float64) int {
	switch {
	case v < 0:
		return -1
	case v > 0:
		return 1
	}
	return 0
}

// cfvo evaluates a threshold of a color scale, data bar or icon set.
func (s *sheetCtx) cfvo(n *ooxml.Node, nums []float64, def string) (float64, bool) {
	if len(nums) == 0 {
		return 0, false
	}
	lo, hi := nums[0], nums[0]
	for _, v := range nums {
		lo, hi = math.Min(lo, v), math.Max(hi, v)
	}
	val := n.AttrStr("val", "")
	if f := n.Child("f"); f != nil {
		val = f.Content() // Excel 2010 thresholds are formulas
	}
	num := func() (float64, bool) {
		f, err := strconv.ParseFloat(strings.TrimPrefix(strings.TrimSpace(val), "="), 64)
		return f, err == nil
	}
	switch n.AttrStr("type", def) {
	case "min", "autoMin":
		return lo, true
	case "max", "autoMax":
		return hi, true
	case "num":
		return num()
	case "formula":
		if v, ok := num(); ok {
			return v, true
		}
		if fv, ok := s.evalFormula(val, [2]int{}, 0, 0); ok && fv.kind == fNum {
			return fv.num, true
		}
		return 0, false
	case "percent":
		p, ok := num()
		return lo + (hi-lo)*p/100, ok
	case "percentile":
		p, ok := num()
		return percentile(nums, p), ok
	}
	return 0, false
}

func percentile(nums []float64, p float64) float64 {
	sorted := append([]float64(nil), nums...)
	sort.Float64s(sorted)
	if len(sorted) == 1 {
		return sorted[0]
	}
	k := p / 100 * float64(len(sorted)-1)
	i := int(math.Floor(k))
	if i >= len(sorted)-1 {
		return sorted[len(sorted)-1]
	}
	return sorted[i] + (sorted[i+1]-sorted[i])*(k-float64(i))
}

func (s *sheetCtx) colorScale(rule *cfRule, vals []cellVal, nums []float64, get func(r, c int) *cfResult) {
	cs := rule.n.Child("colorScale")
	vos, cols := cs.Children("cfvo"), cs.Children("color")
	if len(vos) < 2 || len(cols) < len(vos) {
		return
	}
	var th []float64
	var cc []rgb
	for i, vo := range vos {
		v, ok := s.cfvo(vo, nums, "min")
		if !ok {
			return
		}
		th = append(th, v)
		cc = append(cc, s.c.st.color(parseColor(cols[i]), white))
	}
	for _, v := range vals {
		if !v.isNum() {
			continue
		}
		res := get(v.r, v.c)
		if res.stopped || res.fill != nil {
			continue
		}
		col := cc[len(cc)-1]
		switch {
		case v.num <= th[0]:
			col = cc[0]
		case v.num >= th[len(th)-1]:
		default:
			for i := 0; i+1 < len(th); i++ {
				if v.num <= th[i+1] {
					t := 0.0
					if th[i+1] > th[i] {
						t = (v.num - th[i]) / (th[i+1] - th[i])
					}
					col = mix(cc[i], cc[i+1], t)
					break
				}
			}
		}
		res.fill = &xfill{pattern: "solid", fg: colorRef{kind: colorRGB, rgb: col}}
	}
}

func (s *sheetCtx) dataBar(rule *cfRule, vals []cellVal, nums []float64, get func(r, c int) *cfResult) {
	db := rule.n.Child("dataBar")
	vos := db.Children("cfvo")
	if len(vos) < 2 {
		return
	}
	lo, ok1 := s.cfvo(vos[0], nums, "min")
	hi, ok2 := s.cfvo(vos[1], nums, "max")
	if !ok1 || !ok2 {
		return
	}
	minLen, maxLen := db.AttrFloat("minLength", 10)/100, db.AttrFloat("maxLength", 90)/100
	color := s.c.st.color(parseColor(db.Child("color")), hexRGB(0x638EC6))
	gradient := true
	var border *rgb
	negColor := hexRGB(0xFF0000)
	if x := rule.x14.Child("dataBar"); x != nil {
		minLen, maxLen = x.AttrFloat("minLength", 0)/100, x.AttrFloat("maxLength", 100)/100
		gradient = x.AttrBool("gradient", true)
		if x.AttrBool("border", false) {
			bc := s.c.st.color(parseColor(x.Child("borderColor")), color)
			border = &bc
		}
		if nc := x.Child("negativeFillColor"); nc != nil {
			negColor = s.c.st.color(parseColor(nc), negColor)
		}
		// Excel 2010 bars scale from zero when the values straddle it
		if lo > 0 && x.Child("cfvo").AttrStr("type", "") == "autoMin" {
			lo = 0
		}
	}
	show := db.AttrBool("showValue", true)
	for _, v := range vals {
		if !v.isNum() {
			continue
		}
		res := get(v.r, v.c)
		if res.stopped || res.bar != nil {
			continue
		}
		bar := &dataBar{color: color, border: border, gradient: gradient}
		if lo < 0 && hi > 0 && rule.x14 != nil {
			// an axis where zero falls: negative bars go left of it
			axis := -lo / (hi - lo)
			f := v.num / (hi - lo)
			if v.num < 0 {
				bar.lo, bar.hi = axis+f, axis
				bar.color = negColor
			} else {
				bar.lo, bar.hi = axis, axis+f
			}
		} else {
			t := 1.0
			if hi > lo {
				t = clamp01((v.num - lo) / (hi - lo))
			}
			bar.lo, bar.hi = 0, minLen+t*(maxLen-minLen)
		}
		res.bar = bar
		if !show {
			res.hideValue = true
		}
	}
}

func (s *sheetCtx) iconSet(rule *cfRule, vals []cellVal, nums []float64, get func(r, c int) *cfResult) {
	is := rule.n.Child("iconSet")
	set := is.AttrStr("iconSet", "3TrafficLights1")
	vos := is.Children("cfvo")
	if len(vos) == 0 {
		return
	}
	type thr struct {
		v   float64
		gte bool
	}
	var th []thr
	for i, vo := range vos {
		if i == 0 {
			continue
		}
		v, ok := s.cfvo(vo, nums, "percent")
		if !ok {
			return
		}
		th = append(th, thr{v, vo.AttrBool("gte", true)})
	}
	reverse := is.AttrBool("reverse", false)
	show := is.AttrBool("showValue", true)
	// custom icon sets pick each icon from a set
	var custom []iconRef
	if is.AttrBool("custom", false) {
		for _, ic := range is.Children("cfIcon") {
			custom = append(custom, iconRef{set: ic.AttrStr("iconSet", ""), idx: int(ic.AttrInt("iconId", 0))})
		}
	}
	for _, v := range vals {
		if !v.isNum() {
			continue
		}
		res := get(v.r, v.c)
		if res.stopped || res.icon != nil {
			continue
		}
		idx := 0
		for i, t := range th {
			if v.num > t.v || t.gte && v.num == t.v {
				idx = i + 1
			}
		}
		if reverse {
			idx = len(vos) - 1 - idx
		}
		ic := iconRef{set: set, idx: idx}
		if idx < len(custom) {
			ic = custom[idx]
		}
		if ic.set != "NoIcons" {
			res.icon = &ic
		}
		if !show {
			res.hideValue = true
		}
	}
}

// paintConditional draws the data bars and icons of conditional formats.
func (s *sheetCtx) paintConditional() {
	if s.cf == nil {
		return
	}
	keys := make([][2]int, 0, len(s.cf.cells))
	for k, res := range s.cf.cells {
		if res.bar != nil || res.icon != nil {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i][0] != keys[j][0] {
			return keys[i][0] < keys[j][0]
		}
		return keys[i][1] < keys[j][1]
	})
	for _, k := range keys {
		res := s.cf.cells[k]
		b := s.cellBox(k[0], k[1])
		if m, ok := s.mergeAt(k[0], k[1]); ok {
			if m.r0 != k[0] || m.c0 != k[1] {
				continue
			}
			b = s.rangeBox(m)
		}
		if b.w <= 0 || b.h <= 0 {
			continue
		}
		if bar := res.bar; bar != nil {
			s.drawBar(b, bar)
		}
		if ic := res.icon; ic != nil {
			s.drawIcon(b, ic)
		}
	}
}

func (s *sheetCtx) drawBar(b box, bar *dataBar) {
	inner := box{b.x + padL, b.y + 2*pxPt, b.w - padL - padR, b.h - 4*pxPt}
	x0, x1 := inner.x+inner.w*bar.lo, inner.x+inner.w*bar.hi
	if x1-x0 < 0.1 || inner.h <= 0 {
		return
	}
	for _, t := range s.tilesIn(x0, inner.y, x1, inner.y+inner.h) {
		if bar.gradient {
			p := bdf.LinearGradient(f32(x0-t.ox), 0, f32(x1-t.ox), 0,
				bdf.Stop{Offset: 0, Color: bar.color.bdf()}, bdf.Stop{Offset: 1, Color: mix(bar.color, white, 0.9).bdf()})
			t.cv.Obj.FillPaint(t.cv.Obj.AddPaint(p))
			t.st.hasFill = false
		} else {
			t.fillColor(bar.color.bdf())
		}
		t.rect(x0, inner.y, x1-x0, inner.h)
		if bar.border != nil {
			t.cv.Obj.Line(f32(pxPt), bdf.CapButt, bdf.JoinMiter, 10)
			t.cv.Obj.StrokeColor(bar.border.bdf())
			t.cv.Obj.StrokeRect(f32(x0-t.ox+pxPt/2), f32(inner.y-t.oy+pxPt/2), f32(x1-x0-pxPt), f32(inner.h-pxPt))
		}
	}
}
