package pdfbdf

import "github.com/shibukawa/bdf/internal/affine"

// matrix is an affine transform: x' = a·x + c·y + e, y' = b·x + d·y + f,
// as [a b c d e f] (see affine.Matrix).
type matrix = affine.Matrix

var identity = affine.Identity
