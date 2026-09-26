package psd

// descriptor is a Photoshop action descriptor, the structure that newer
// additional layer information (artboards, fill layers) is stored in. The
// values read are float64 ('doub', 'UntF'), int64 ('long', 'comp'), bool,
// string ('TEXT', and the enumeration value of 'enum'), *descriptor
// ('Objc', 'GlbO') and []any ('VlLs'). Reading stops at a type it does not
// know, keeping the items before it.
type descriptor struct {
	class string
	items map[string]any
}

// rgb returns the colour of an RGBC descriptor ('Rd  ', 'Grn ', 'Bl  ',
// 0-255).
func (d *descriptor) rgb() ([3]float64, bool) {
	r, ok1 := d.items["Rd  "].(float64)
	g, ok2 := d.items["Grn "].(float64)
	b, ok3 := d.items["Bl  "].(float64)
	return [3]float64{r, g, b}, ok1 && ok2 && ok3
}

// id reads a key or class ID: a length and a string, or a length of 0 and a
// four-character code.
func (c *cursor) id() string {
	n := int64(c.u32())
	if n == 0 {
		n = 4
	}
	if n > c.left() || n > 1024 {
		c.err = errFormat
		return ""
	}
	return string(c.read(int(n)))
}

func readDescriptor(c *cursor) *descriptor {
	d := &descriptor{items: map[string]any{}}
	if !d.read(c, 0) && len(d.items) == 0 {
		return nil
	}
	return d
}

func (d *descriptor) read(c *cursor, depth int) bool {
	c.unicode() // name
	d.class = c.id()
	n := int(c.u32())
	for i := 0; i < n && c.err == nil; i++ {
		key := c.id()
		v, ok := readValue(c, depth)
		if !ok || c.err != nil {
			return false
		}
		d.items[key] = v
	}
	return c.err == nil
}

func readValue(c *cursor, depth int) (any, bool) {
	if depth > 16 {
		return nil, false
	}
	switch string(c.read(4)) {
	case "Objc", "GlbO":
		sub := &descriptor{items: map[string]any{}}
		return sub, sub.read(c, depth+1)
	case "VlLs":
		n := int(c.u32())
		if int64(n) > c.left() {
			return nil, false
		}
		list := make([]any, 0, n)
		for i := 0; i < n; i++ {
			v, ok := readValue(c, depth+1)
			if !ok {
				return nil, false
			}
			list = append(list, v)
		}
		return list, true
	case "doub":
		return c.f64(), true
	case "UntF":
		c.read(4) // unit
		return c.f64(), true
	case "long":
		return int64(c.i32()), true
	case "comp":
		return int64(c.u64()), true
	case "bool":
		return c.u8() != 0, true
	case "TEXT":
		return c.unicode(), true
	case "enum":
		c.id() // type
		return c.id(), true
	case "type", "GlbC":
		c.unicode()
		return c.id(), true
	case "alis", "tdta":
		c.skip(int64(c.u32()))
		return nil, true
	case "UnFl":
		c.read(4)
		n := int64(c.u32())
		c.skip(8 * n)
		return nil, true
	}
	return nil, false
}
