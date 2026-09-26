package pdf

import (
	"bytes"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
)

func init() {
	conv.Register(&conv.Format{
		Name:        "pdf",
		Description: "PDF",
		Extensions:  []string{".pdf"},
		Params: []conv.Param{
			{Name: "kind", Usage: "view kind: fixed (default) or flow"},
			{Name: "no-share", Usage: "true: do not move the instruction prefix pages have in common into a shared object"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			return bytes.Contains(head, []byte("%PDF-"))
		},
		CheckPassword: func(r io.ReaderAt, size int64, password string) (bool, error) {
			_, protected, err := readContext(io.NewSectionReader(r, 0, size), password)
			return protected, err
		},
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			opts, err := options(o)
			if err != nil {
				return nil, err
			}
			res, err := Convert(io.NewSectionReader(r, 0, size), opts)
			if err != nil {
				return nil, err
			}
			return result(res), nil
		},
		Stream: func(r io.ReaderAt, size int64, o *conv.Options) (conv.Stream, error) {
			opts, err := options(o)
			if err != nil {
				return nil, err
			}
			s, err := NewStream(io.NewSectionReader(r, 0, size), opts)
			if err != nil {
				return nil, err
			}
			return stream{s}, nil
		},
	})
}

func options(o *conv.Options) (*Options, error) {
	noShare, err := o.BoolParam("no-share")
	if err != nil {
		return nil, fmt.Errorf("pdf: %w", err)
	}
	return &Options{Pages: o.Pages, Title: o.Title, Kind: o.Param("kind"),
		NoTextIndex: o.NoTextIndex, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType,
		Images: o.Images, NoSharePrefix: noShare, Password: o.Password, Warn: o.Warn}, nil
}

func result(res *Result) *conv.Result {
	return &conv.Result{Doc: res.Doc, Warnings: res.Warnings, Protected: res.Protected,
		Summary: fmt.Sprintf("%d page(s), %d shared prefix(es) saving %d bytes", res.Pages, res.SharedPrefixes, res.SharedBytes)}
}

// stream is a Stream as the converter package has it.
type stream struct{ *Stream }

func (s stream) Finish() (*conv.Result, error) {
	res, err := s.Stream.Finish()
	if err != nil {
		return nil, err
	}
	return result(res), nil
}
