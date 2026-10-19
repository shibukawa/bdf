package ipp

// Group returns the first group of the tag, or nil.
func (m *Message) Group(tag Tag) *Group {
	for i := range m.Groups {
		if m.Groups[i].Tag == tag {
			return &m.Groups[i]
		}
	}
	return nil
}

// Attr returns the first attribute of the name in the first group of the
// tag, or nil.
func (m *Message) Attr(tag Tag, name string) *Attribute {
	if g := m.Group(tag); g != nil {
		return g.Get(name)
	}
	return nil
}

// Get returns the first attribute of the name, or nil.
func (g *Group) Get(name string) *Attribute {
	if g == nil {
		return nil
	}
	for i := range g.Attrs {
		if g.Attrs[i].Name == name {
			return &g.Attrs[i]
		}
	}
	return nil
}

// Get returns the first attribute of the name among members of a
// collection, or nil.
func Get(members []Attribute, name string) *Attribute {
	for i := range members {
		if members[i].Name == name {
			return &members[i]
		}
	}
	return nil
}

// First returns the first value of the attribute. The attribute may be
// nil: ok is then false.
func (a *Attribute) First() (v Value, ok bool) {
	if a == nil || len(a.Values) == 0 {
		return Value{}, false
	}
	return a.Values[0], true
}

// Strings returns the character-string values of the attribute (keywords,
// names, texts, URIs …), skipping the others. The attribute may be nil.
func (a *Attribute) Strings() []string {
	if a == nil {
		return nil
	}
	var out []string
	for _, v := range a.Values {
		if v.Tag.IsString() || v.Tag == TagTextLang || v.Tag == TagNameLang || v.Tag == TagOctetString {
			out = append(out, v.Str)
		}
	}
	return out
}

// Text returns the first character-string value, or "".
func (a *Attribute) Text() string {
	if s := a.Strings(); len(s) > 0 {
		return s[0]
	}
	return ""
}

// Has reports whether a character-string value of the attribute is s.
func (a *Attribute) Has(s string) bool {
	for _, v := range a.Strings() {
		if v == s {
			return true
		}
	}
	return false
}

// Ints returns the integer, enum and boolean values of the attribute.
func (a *Attribute) Ints() []int32 {
	if a == nil {
		return nil
	}
	var out []int32
	for _, v := range a.Values {
		if v.Tag == TagInteger || v.Tag == TagEnum || v.Tag == TagBoolean {
			out = append(out, v.Int)
		}
	}
	return out
}

// Int returns the first integer, enum or boolean value.
func (a *Attribute) Int() (n int32, ok bool) {
	if s := a.Ints(); len(s) > 0 {
		return s[0], true
	}
	return 0, false
}

// Bool returns the first boolean value.
func (a *Attribute) Bool() (b, ok bool) {
	n, ok := a.Int()
	return n != 0, ok
}

// Collections returns the collection values of the attribute.
func (a *Attribute) Collections() [][]Attribute {
	if a == nil {
		return nil
	}
	var out [][]Attribute
	for _, v := range a.Values {
		if v.Tag == TagBeginCollection {
			out = append(out, v.Coll)
		}
	}
	return out
}

// Resolutions returns the resolution values of the attribute.
func (a *Attribute) Resolutions() []Resolution {
	if a == nil {
		return nil
	}
	var out []Resolution
	for _, v := range a.Values {
		if v.Tag == TagResolution {
			out = append(out, v.Res)
		}
	}
	return out
}

// Ranges returns the rangeOfInteger values of the attribute.
func (a *Attribute) Ranges() []Range {
	if a == nil {
		return nil
	}
	var out []Range
	for _, v := range a.Values {
		if v.Tag == TagRange {
			out = append(out, v.Range)
		}
	}
	return out
}

// OutOfBand returns the tag of an out-of-band first value (unsupported,
// unknown, no-value …), or 0.
func (a *Attribute) OutOfBand() Tag {
	if v, ok := a.First(); ok && v.Tag.IsOutOfBand() {
		return v.Tag
	}
	return 0
}
