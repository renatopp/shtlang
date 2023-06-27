package test

import (
	"sht/lang"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	cases := []struct{ input, expected string }{
		{`fn foo(a) { a }; foo(5)`, "5"},
		{`fn foo(a) { a + 4 }; foo(5)`, "9"},
		{`fn foo(a, b) { a + b }; foo(1, 2)`, "3"},
		{`fn foo(a, b, c=2) { (a + b)*c }; foo(1, 2)`, "6"},
		{`fn foo(a, b, c=2) { (a + b)*c }; foo(1, 2, 3)`, "9"},
		{`fn foo(a=1, b=1, c=2) { (a + b)*c }; foo(1)`, "4"},
		{`fn foo(a=1, b=1, c=2) { (a + b)*c }; foo()`, "4"},

		{`let foo = fn(a) { a }; foo(5)`, "5"},
		{`let foo = fn { 2 }; foo(5)`, "2"},

		{`fn adder(a) { fn (b, c) { (b+c)*a } }; let add = adder(5); add(1, 2)`, "15"},
		{`fn func() { return 5; 1 }; func()`, "5"},
		{`
			let a = 1
			let b = 2

			fn scoped(c) {
				let b = 5
				
				a + b + c
			}

			scoped(10)
		`, "16"},
		{`
			let b = 2
			fn scoped() {
				let b = 5
			}

			scoped(10)
			b
		`, "2"},
		{`
			let b = 2
			fn scoped() {
				b = 5
			}

			scoped(10)
			b
		`, "5"},
	}

	for _, c := range cases {
		result, err := lang.Eval([]byte(c.input))

		assert.NoError(t, err)
		assert.Equal(t, c.expected, result)
	}
}
