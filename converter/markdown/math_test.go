package markdown

import (
	"strings"
	"testing"
)

func TestMathHTML(t *testing.T) {
	for _, tc := range []struct{ md, want string }{
		{"a $x^2$ b", `a <math><semantics><annotation encoding="application/x-tex">x^2</annotation></semantics></math> b`},
		{"a $`x < y`$ b", `<annotation encoding="application/x-tex">x &lt; y</annotation>`},
		{"costs $5 and $10", "costs $5 and $10"},
		{"$5 と $10 は $x$", `$5 と $10 は <math><semantics><annotation encoding="application/x-tex">x</annotation>`},
		{"$ x$ and $x $", "$ x$ and $x $"},
		{`escaped \$x$`, "$x$"},
		{"in $$E=mc^2$$ text", `in <math display="block">`},
		{"$$\n\\frac{a}{b}\n$$\n", `<math display="block"><semantics><annotation encoding="application/x-tex">\frac{a}{b}</annotation>`},
		{"$$ x $$\n", `<math display="block"><semantics><annotation encoding="application/x-tex">x</annotation>`},
		{"```math\n\\sqrt{2}\n```\n", `<math display="block"><semantics><annotation encoding="application/x-tex">\sqrt{2}</annotation>`},
		{"```go\nx := 1\n```\n", `<code class="language-go">`},
		{"`$x$`", "<code>$x$</code>"},
	} {
		out, err := ToHTML([]byte(tc.md))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), tc.want) {
			t.Errorf("%q: %q not in\n%s", tc.md, tc.want, out)
		}
	}
}

func TestMath(t *testing.T) {
	_, c := convertSrc(t, []byte("# Formulas\n\nEuler's identity $e^{i\\pi}+1=0$ holds.\n\n$$\n\\int_0^1 x\\,dx = \\frac{1}{2}\n$$\n"))
	for _, want := range []string{"Euler's identity e^(iπ)+1=0 holds.", "∫_0^1 x dx=1/2"} {
		if !strings.Contains(c.text, want) {
			t.Errorf("%q not in %q", want, c.text)
		}
	}
}
