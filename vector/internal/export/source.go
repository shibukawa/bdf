package export

import (
	"fmt"

	"github.com/shibukawa/bdf"
)

// Source is the document a writer writes, with what it has decoded of it
// and the warnings of the writing. A writer embeds it.
type Source struct {
	Doc *bdf.Document
	// SkipRoles are the roles of the layers that are not written.
	SkipRoles []string
	// Warnings are the warnings of the writing, each once.
	Warnings []string
	// Seen holds the objects written before (see Enter).
	Seen map[bdf.Hash]struct{}

	objects map[bdf.Hash]*bdf.ObjectPart // nil: not usable
	coll    map[bdf.Hash][]*bdf.Path
	warned  map[string]bool
}

// Warnf adds a warning, once.
func (s *Source) Warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if !s.warned[msg] {
		if s.warned == nil {
			s.warned = map[string]bool{}
		}
		s.warned[msg] = true
		s.Warnings = append(s.Warnings, msg)
	}
}

// Object is an object of the document, decoded once; nil (with a warning)
// for one that is missing or cannot be read.
func (s *Source) Object(h bdf.Hash) *bdf.ObjectPart {
	if o, ok := s.objects[h]; ok {
		return o
	}
	if s.objects == nil {
		s.objects = map[bdf.Hash]*bdf.ObjectPart{}
	}
	s.objects[h] = nil
	p := s.Doc.Part(h)
	if p == nil {
		s.Warnf("object %s is missing", h)
		return nil
	}
	o, err := bdf.DecodeObject(p.Data)
	if err != nil {
		s.Warnf("object %s: %v", h, err)
		return nil
	}
	s.objects[h] = o
	return o
}

// Path is a path of an object: inline, or of a path collection of the
// document, decoded once; nil (with a warning) for one there is not.
func (s *Source) Path(obj *bdf.ObjectPart, i int) *bdf.Path {
	if i < 0 || i >= len(obj.Paths) {
		s.Warnf("bad path reference")
		return nil
	}
	e := obj.Paths[i]
	if e.Inline != nil {
		return e.Inline
	}
	ps, ok := s.coll[e.Hash]
	if !ok {
		if part := s.Doc.Part(e.Hash); part != nil {
			var err error
			if ps, err = bdf.DecodePathCollection(part.Data); err != nil {
				s.Warnf("path collection %s: %v", e.Hash, err)
			}
		} else {
			s.Warnf("path collection %s is missing", e.Hash)
		}
		if s.coll == nil {
			s.coll = map[bdf.Hash][]*bdf.Path{}
		}
		s.coll[e.Hash] = ps
	}
	if int(e.Index) < len(ps) {
		return ps[e.Index]
	}
	return nil
}

// Enter notes an object as written and reports whether it was before.
func (s *Source) Enter(h bdf.Hash) bool {
	if _, ok := s.Seen[h]; ok {
		return true
	}
	if s.Seen == nil {
		s.Seen = map[bdf.Hash]struct{}{}
	}
	s.Seen[h] = struct{}{}
	return false
}

// Skip reports whether the layers of a role are not written.
func (s *Source) Skip(role string) bool {
	for _, r := range s.SkipRoles {
		if r == role {
			return true
		}
	}
	return false
}
