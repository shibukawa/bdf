//go:build tinygo

package form

import "fmt"

// decodeYAML reads a YAML document: not in builds for TinyGo, whose
// programs give the data as JSON.
func decodeYAML(data []byte) (any, error) {
	return nil, fmt.Errorf("YAML data is not read in this build; give the data as JSON")
}
