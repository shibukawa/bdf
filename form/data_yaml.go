//go:build !tinygo

package form

import (
	"bytes"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

// decodeYAML reads a YAML document.
func decodeYAML(data []byte) (any, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var v any
	if err := dec.Decode(&v); err != nil && err != io.EOF {
		return nil, err
	}
	var more any
	if err := dec.Decode(&more); err != io.EOF {
		return nil, fmt.Errorf("more than one YAML document")
	}
	return v, nil
}
