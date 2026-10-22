package docflow

// The names the IPP adapter (docflow/io/ipp, registered by
// adapters_test.go) gives its kinds and its ingress's path, which the
// tests of this package cannot import.
const (
	KindPrinter = "printer"
	KindCups    = "cups"
	IPPPath     = "/ipp/print"
	RouteBranch = "branch"
)
