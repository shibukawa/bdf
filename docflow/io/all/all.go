// Package all imports every input and output of docflow, and its render
// node, so that a program that serves them all imports this package
// alone: the IPP printers and CUPS servers (docflow/io/ipp), mail
// (docflow/io/smtp), S3 and its trigger (docflow/io/s3), HTTP
// (docflow/io/http), the notify node (docflow/io/notify) and the render
// node (docflow/render). A program that serves some of them imports
// those packages instead.
package all

import (
	_ "github.com/shibukawa/bdf/docflow/io/http"
	_ "github.com/shibukawa/bdf/docflow/io/ipp"
	_ "github.com/shibukawa/bdf/docflow/io/notify"
	_ "github.com/shibukawa/bdf/docflow/io/s3"
	_ "github.com/shibukawa/bdf/docflow/io/smtp"
	_ "github.com/shibukawa/bdf/docflow/render"
)
