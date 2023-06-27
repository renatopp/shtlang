package test

import (
	"sht/lang"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessString(t *testing.T) {
	cases := []struct{ input, expected string }{
		{`'hello' + 'world'`, "helloworld"},
		{`'hello' .. 'world'`, "helloworld"},
		{`'hello' .. 1`, "hello1"},
	}

	for _, c := range cases {
		result, err := lang.Eval([]byte(c.input))

		assert.NoError(t, err)
		assert.Equal(t, c.expected, result)
	}
}

func TestErrorString(t *testing.T) {
	cases := []struct{ input string }{
		{`'hello' - 'world'`},
		{`'hello' + 2`},
	}

	for _, c := range cases {
		result, err := lang.Eval([]byte(c.input))

		assert.NoError(t, err)
		assert.Equal(t, "ERR!", result[:4])
	}
}
