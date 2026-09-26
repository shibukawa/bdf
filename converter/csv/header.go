package csv

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// Whether the first record names the columns is guessed from how its
// values differ from those under them, column by column, in the first
// thousand records:
//
//   - a label over a column of numbers, dates, booleans, mail addresses or
//     URLs heads it; a value of the column's own kind over it is data
//     (a year over other numbers is a label: 2023, 2024 over sales);
//   - over text, a value that recurs further down is data, and a name of
//     another length over codes of one length heads them;
//   - the first record itself: labels are distinct, rarely empty (but for
//     the corner of a cross table) and rarely numbers.
//
// Files whose first record is all distinct text lean toward a header, as
// most files with text columns have one.

// maxHeaderSample bounds the records the header guess reads.
const maxHeaderSample = 1000

// hasHeader guesses whether the first of some records names the columns.
func hasHeader(recs [][]string) bool {
	if len(recs) < 2 || len(recs[0]) == 0 {
		return false
	}
	first, data := recs[0], recs[1:min(len(recs), maxHeaderSample)]
	score := 0.0
	allText := true
	seen := map[string]bool{}
	for j, v := range first {
		v = strings.TrimSpace(v)
		switch k := classify(v); {
		case k == kindEmpty:
			if j > 0 {
				score--
				allText = false
			}
			continue
		case k.numeric() && !yearLike(v):
			score--
			allText = false
		}
		if seen[v] {
			score -= 2
			allText = false
		}
		seen[v] = true
	}
	if allText {
		score += 0.5
	}
	for j, h := range first {
		h = strings.TrimSpace(h)
		hk := classify(h)
		if hk == kindEmpty {
			continue
		}
		counts := map[kind]int{}
		values := map[string]bool{}
		lengths := map[int]bool{}
		n, years := 0, 0
		for _, rec := range data {
			if j >= len(rec) {
				continue
			}
			v := strings.TrimSpace(rec[j])
			k := classify(v)
			if k == kindEmpty {
				continue
			}
			counts[k]++
			n++
			values[v] = true
			lengths[utf8.RuneCountInString(v)] = true
			if yearLike(v) {
				years++
			}
		}
		if n == 0 {
			continue
		}
		dom, most := kindText, 0
		for k := kindText; k <= kindURL; k++ {
			if counts[k] > most {
				dom, most = k, counts[k]
			}
		}
		if yearLike(h) && years < n {
			hk = kindText // a year over other values
		}
		switch {
		case float64(most) < 0.8*float64(n):
			// a mixed column: only a recurring value tells
			if values[h] {
				score--
			}
		case hk != dom && dom != kindText:
			score += 2
		case hk == dom && dom != kindText:
			score -= 2
		case hk != dom:
			score-- // a number or such over text
		case values[h]:
			score--
		case len(lengths) == 1 && n >= 3 && !lengths[utf8.RuneCountInString(h)]:
			score++
		}
	}
	return score > 0
}

// yearLike reports whether a value is a year from 1900 to 2100.
func yearLike(v string) bool {
	if len(v) != 4 || !allDigits(v) {
		return false
	}
	y, _ := strconv.Atoi(v)
	return y >= 1900 && y <= 2100
}
