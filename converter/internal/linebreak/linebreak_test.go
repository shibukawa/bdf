package linebreak

import "testing"

func TestAllowed(t *testing.T) {
	for _, c := range []struct {
		a, b rune
		want bool
	}{
		{' ', 'a', true}, {'a', ' ', false}, {'a', 'b', false}, {'-', 'b', true},
		{'日', '本', true}, {'本', '。', false}, {'「', '本', false}, {'本', 'A', true}, {'ー', 'ト', true}, {'テ', 'ー', false},
		{'a', '\u00a0', false}, {'\u3000', '本', true},
	} {
		if got := Allowed(c.a, c.b); got != c.want {
			t.Errorf("Allowed(%q, %q) = %v", c.a, c.b, got)
		}
	}
	for _, c := range []struct {
		a, b rune
		want bool
	}{{'日', '本', true}, {'本', 'A', true}, {'a', 'b', false}, {' ', '本', false}} {
		if got := Joins(c.a, c.b); got != c.want {
			t.Errorf("Joins(%q, %q) = %v", c.a, c.b, got)
		}
	}
}
