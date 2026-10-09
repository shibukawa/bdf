package ooxml

import (
	"io"

	"github.com/shibukawa/bdf/internal/xmltree"
)

// Node is a generic XML element (see xmltree.Node): the type lives outside
// the converters so that the formula engine reads Office Math from it.
type Node = xmltree.Node

// Segment is a piece of the content of an element.
type Segment = xmltree.Segment

// Attr is an attribute of an element, and Name its name.
type (
	Attr = xmltree.Attr
	Name = xmltree.Name
)

const (
	// NSA14 is the namespace of the Office 2010 DrawingML extensions.
	NSA14 = xmltree.NSA14
	// EMUPerPoint is the number of English Metric Units in a point.
	EMUPerPoint = xmltree.EMUPerPoint
)

// Parse reads a document into a node tree (see xmltree.Parse).
func Parse(data []byte) (*Node, error) { return xmltree.Parse(data) }

// ParseChoosing is Parse that replaces mc:AlternateContent with its first
// mc:Choice whose required namespaces all pass supported.
func ParseChoosing(data []byte, supported func(prefix string) bool) (*Node, error) {
	return xmltree.ParseChoosing(data, supported)
}

// ParsePicking is Parse that replaces mc:AlternateContent with its first
// mc:Choice that pick accepts.
func ParsePicking(data []byte, pick func(choice *Node) bool) (*Node, error) {
	return xmltree.ParsePicking(data, pick)
}

// ParsePickingReader is ParsePicking from a reader.
func ParsePickingReader(r io.Reader, pick func(choice *Node) bool) (*Node, error) {
	return xmltree.ParsePickingReader(r, pick)
}

// MathChoice picks the choices that hold Office Math in DrawingML text.
func MathChoice(c *Node) bool { return xmltree.MathChoice(c) }
