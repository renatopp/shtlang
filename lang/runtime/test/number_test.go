package test

import (
	"sht/lang"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessNumber(t *testing.T) {
	cases := []struct{ input, expected string }{
		{`1`, "1"},
		{`1.5`, "1.500000"},
		{`-10`, "-10"},
		{`2.3e3`, "2300"},
		{`2.3e-3`, "0.002300"},
		{`.5`, "0.500000"},

		{`1 + 1`, "2"},
		{`1.5 + 1`, "2.500000"},
		{`2 - 1`, "1"},
		{`2 - 4`, "-2"},
		{`2 * 4`, "8"},
		{`1 + 2 * 4`, "9"},
		{`(1 + 2) * 3`, "9"},
		{`9 / 3`, "3"},
		{`9 / 2`, "4.500000"},
		{`9 // 2`, "4"},
		{`9 % 2`, "1"},
		{`2**10`, "1024"},

		{`!1`, "false"},
		{`!!0`, "false"},
		{`!0`, "true"},
		{`!!23123`, "true"},
	}

	for _, c := range cases {
		result, err := lang.Eval([]byte(c.input))

		assert.NoError(t, err)
		assert.Equal(t, c.expected, result)
	}
}

func TestErrorNumber(t *testing.T) {
	cases := []struct{ input string }{
		{`2 + true`},
		{`3 // false`},
	}

	for _, c := range cases {
		result, err := lang.Eval([]byte(c.input))

		assert.NoError(t, err)
		assert.Equal(t, "ERR!", result[:4])
	}
}
