package core
import (
	"testing"
)

func TestCompareExpr(t *testing.T) {
	cases := []struct{
		name string
		a Expr
		b Expr
		less bool
	}{

		{
			name: "symbol",
			a: NewSymbol("alpha"),
			b: NewSymbol("beta"),
			less: true,
		},
		{
			name: "symbol",
			a: NewSymbol("beta"),
			b: NewSymbol("alpha"),
			less: false,
		},
		{
			name: "strings",
			a: NewString("aa"),
			b: NewString("bb"),
			less: true,
		},
		{
			name: "strings",
			a: NewString("bb"),
			b: NewString("aa"),
			less: false,
		},
	}

	for _, tt := range cases {
		got := CanonicalCompare(tt.a, tt.b)
		if got != tt.less {
			t.Errorf("Case %s: %s < %s is %v but got %v", 
				tt.name, tt.a, tt.b, tt.less, got)
		}
	}
}

