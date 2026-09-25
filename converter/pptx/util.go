package pptx

import (
	"strconv"
	"strings"
)

func atof(s string, def float64) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return def
	}
	return f
}

func ftoa(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// parsePct parses an ST_Percentage value into a fraction.
func parsePct(v string) float64 {
	v = strings.TrimSpace(v)
	if strings.HasSuffix(v, "%") {
		return atof(strings.TrimSuffix(v, "%"), 0) / 100
	}
	return atof(v, 0) / 100000
}
