package parquet

import (
	"errors"
	"fmt"
	"strings"
)

// The schema of a Parquet file is a tree listed depth first: the root is a
// group of the columns, groups hold fields, and the leaves are the columns
// whose values the file stores, one column chunk each in every row group.
// A value's definition level counts the optional and repeated fields on
// its path that are there (below it, the value is null), and its
// repetition level tells in which of the repeated fields on its path a new
// element starts (0: a new row).

// node is a field of the schema.
type node struct {
	name     string
	rep      int32 // repRequired, repOptional, repRepeated
	leaf     bool
	phys     int32 // of leaves
	typeLen  int32 // of FIXED_LEN_BYTE_ARRAY leaves
	conv     int32
	logical  logicalType
	children []*node
	index    int // among its parent's children
	// defLevel and repLevel count the optional and repeated fields from
	// the top down to this one, this one included.
	defLevel, repLevel int
	col                int // of leaves: the column number
	geo                bool
}

// maxSchemaDepth bounds the nesting of fields.
const maxSchemaDepth = 64

var errSchema = errors.New("parquet: malformed schema")

// buildSchema makes the tree of the schema elements, and returns its root
// and its leaves in column order.
func buildSchema(elems []schemaElement) (*node, []*node, error) {
	if len(elems) == 0 {
		return nil, nil, errSchema
	}
	var leaves []*node
	pos := 0
	var build func(parent *node, depth int) (*node, error)
	build = func(parent *node, depth int) (*node, error) {
		if pos >= len(elems) || depth > maxSchemaDepth {
			return nil, errSchema
		}
		e := &elems[pos]
		pos++
		n := &node{name: e.name, rep: e.repetition, conv: e.converted, logical: e.logical}
		if parent == nil || n.rep < repRequired || n.rep > repRepeated {
			n.rep = repRequired
		}
		if parent != nil {
			n.defLevel, n.repLevel = parent.defLevel, parent.repLevel
			switch n.rep {
			case repOptional:
				n.defLevel++
			case repRepeated:
				n.defLevel++
				n.repLevel++
			}
		}
		if n.logical.kind == logNone && n.conv == convDecimal {
			n.logical = logicalType{kind: logDecimal, scale: e.scale, precision: e.precision}
		}
		if e.numChildren <= 0 && e.typ >= 0 {
			if parent == nil {
				return nil, errSchema
			}
			n.leaf, n.phys, n.typeLen = true, e.typ, e.typeLength
			if n.phys < typeBoolean || n.phys > typeFixed || n.phys == typeFixed && n.typeLen <= 0 {
				return nil, fmt.Errorf("parquet: column %q has an unknown type", e.name)
			}
			n.col = len(leaves)
			leaves = append(leaves, n)
			return n, nil
		}
		if int(e.numChildren) > len(elems)-pos {
			return nil, errSchema
		}
		for i := range int(e.numChildren) {
			c, err := build(n, depth+1)
			if err != nil {
				return nil, err
			}
			c.index = i
			n.children = append(n.children, c)
		}
		return n, nil
	}
	root, err := build(nil, 0)
	if err != nil {
		return nil, nil, err
	}
	if root.leaf || pos != len(elems) {
		return nil, nil, errSchema
	}
	return root, leaves, nil
}

// leavesOf returns the leaves under a field (itself when it is one).
func leavesOf(n *node) []*node {
	if n.leaf {
		return []*node{n}
	}
	var out []*node
	for _, c := range n.children {
		out = append(out, leavesOf(c)...)
	}
	return out
}

// pathTo returns the fields from top down to a leaf.
func pathTo(top, leaf *node) []*node {
	if top == leaf {
		return []*node{top}
	}
	for _, c := range top.children {
		if p := pathTo(c, leaf); p != nil {
			return append([]*node{top}, p...)
		}
	}
	return nil
}

func (n *node) isList() bool {
	if n.leaf || len(n.children) != 1 || n.children[0].rep != repRepeated {
		return false
	}
	if n.logical.kind == logList || n.conv == convList {
		return true
	}
	// a map without values is a list of its keys, as Arrow reads it
	kv := n.children[0]
	return n.mapAnnotated() && !kv.leaf && len(kv.children) == 1
}

func (n *node) isMap() bool {
	if n.leaf || !n.mapAnnotated() || len(n.children) != 1 {
		return false
	}
	kv := n.children[0]
	return kv.rep == repRepeated && !kv.leaf && len(kv.children) == 2
}

func (n *node) mapAnnotated() bool {
	return n.logical.kind == logMap || n.conv == convMap || n.conv == convMapKeyValue
}

// isVariant reports whether a group is a Variant: one annotated so, or a
// group of the binary fields metadata and value (with typed_value when it
// is shredded) as Spark wrote them before the annotation.
func (n *node) isVariant() bool {
	if n.leaf {
		return false
	}
	meta, value, _ := n.variantFields()
	if meta == nil || !meta.leaf || meta.phys != typeBinary {
		return false
	}
	if n.logical.kind == logVariant {
		return true
	}
	return value != nil && value.leaf && value.phys == typeBinary && len(n.children) == 2 &&
		meta.logical.kind == logNone && meta.conv == convNone && value.logical.kind == logNone && value.conv == convNone
}

// variantFields returns the metadata, value and typed_value fields of a
// Variant group, or of a field of a shredded object.
func (n *node) variantFields() (meta, value, typed *node) {
	for _, c := range n.children {
		switch c.name {
		case "metadata":
			meta = c
		case "value":
			value = c
		case "typed_value":
			typed = c
		}
	}
	return
}

// listElement returns the node of a list's elements, and whether an
// element is that node's value in the repeated group (the three-level
// list) rather than the repeated field itself (the two-level lists of
// the backward-compatibility rules of LogicalTypes.md).
func (n *node) listElement() (*node, bool) {
	rep := n.children[0]
	switch {
	case rep.leaf, len(rep.children) != 1:
		return rep, false
	case rep.children[0].rep == repRepeated:
		return rep, false
	case rep.name == "array", rep.name == n.name+"_tuple":
		return rep, false
	}
	return rep.children[0], true
}

// typeName names the type of a field: the logical type of a leaf, or the
// types of the elements of lists and maps and of the fields of groups.
func typeName(n *node) string {
	var b strings.Builder
	writeType(&b, n, false)
	return b.String()
}

func writeType(b *strings.Builder, n *node, inst bool) {
	if b.Len() > maxTypeName {
		return
	}
	switch {
	case n.rep == repRepeated && !inst:
		b.WriteString("list<")
		writeType(b, n, true)
		b.WriteString(">")
	case n.leaf:
		b.WriteString(leafTypeName(n))
	case n.isVariant():
		b.WriteString("variant")
	case n.isList():
		el, inner := n.listElement()
		b.WriteString("list<")
		writeType(b, el, !inner)
		b.WriteString(">")
	case n.isMap():
		kv := n.children[0]
		b.WriteString("map<")
		writeType(b, kv.children[0], false)
		b.WriteString(", ")
		writeType(b, kv.children[1], false)
		b.WriteString(">")
	default:
		b.WriteString("struct<")
		for i, c := range n.children {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(c.name)
			b.WriteString(": ")
			writeType(b, c, false)
		}
		b.WriteString(">")
	}
}

// maxTypeName bounds the length of the type names shown.
const maxTypeName = 200

var unitNames = map[int16]string{unitMillis: "ms", unitMicros: "us", unitNanos: "ns"}

// leafTypeName names the type of a leaf as its values are shown: by its
// logical type, or by its physical type when the logical type does not
// apply to it.
func leafTypeName(n *node) string {
	l := n.logical
	switch f := formatOf(n); f.kind {
	case fString:
		if l.kind == logEnum || l.kind == logNone && n.conv == convEnum {
			return "enum"
		}
		return "string"
	case fJSON:
		return "json"
	case fBytes:
		return "bson"
	case fUUID:
		return "uuid"
	case fFloat16:
		return "float16"
	case fDecimal:
		return fmt.Sprintf("decimal(%d, %d)", l.precision, l.scale)
	case fDate:
		return "date"
	case fTime:
		return "time[" + unitNames[f.unit] + "]"
	case fTimestamp:
		if l.kind == logTimestamp && !l.utc {
			return "timestamp[" + unitNames[f.unit] + "]"
		}
		return "timestamp[" + unitNames[f.unit] + ", UTC]"
	case fInt96:
		return "timestamp[ns]"
	case fInt, fUint:
		switch {
		case l.kind == logInteger:
			if l.signed {
				return fmt.Sprintf("int%d", l.bitWidth)
			}
			return fmt.Sprintf("uint%d", l.bitWidth)
		case n.conv >= convUint8 && n.conv <= convUint64:
			return fmt.Sprintf("uint%d", 8<<(n.conv-convUint8))
		case n.conv >= convInt8 && n.conv <= convInt64:
			return fmt.Sprintf("int%d", 8<<(n.conv-convInt8))
		case n.phys == typeInt32:
			return "int32"
		}
		return "int64"
	case fNull:
		return "null"
	case fGeometry:
		if l.kind == logGeography {
			return "geography"
		}
		return "geometry"
	case fInterval:
		return "interval"
	case fBool:
		return "boolean"
	case fFloat:
		return "float"
	case fDouble:
		return "double"
	}
	if n.phys == typeBinary {
		return "binary"
	}
	return fmt.Sprintf("fixed[%d]", n.typeLen)
}
