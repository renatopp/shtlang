package test

import (
	"sht/lang"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessBoolean(t *testing.T) {
	TRUE := "true"
	FALSE := "false"
	cases := []struct{ input, expected string }{
		{`true`, TRUE},
		{`false`, FALSE},

		{`true and false`, FALSE},
		{`false and true`, FALSE},
		{`false and false`, FALSE},
		{`true and true`, TRUE},

		{`true nand false`, TRUE},
		{`false nand true`, TRUE},
		{`false nand false`, TRUE},
		{`true nand true`, FALSE},

		{`true or false`, TRUE},
		{`false or true`, TRUE},
		{`false or false`, FALSE},
		{`true or true`, TRUE},

		{`true nor false`, FALSE},
		{`false nor true`, FALSE},
		{`false nor false`, TRUE},
		{`true nor true`, FALSE},

		{`true xor false`, TRUE},
		{`false xor true`, TRUE},
		{`false xor false`, FALSE},
		{`true xor true`, FALSE},

		{`true nxor false`, FALSE},
		{`false nxor true`, FALSE},
		{`false nxor false`, TRUE},
		{`true nxor true`, TRUE},

		{`!true`, FALSE},
		{`!false`, TRUE},
	}

	for _, c := range cases {
		result, err := lang.Eval([]byte(c.input))

		assert.NoError(t, err)
		assert.Equal(t, c.expected, result)
	}
}

func TestErrorBoolean(t *testing.T) {
	cases := []struct{ input string }{
		{`true + true`},
		{`false - true`},
		{`false * false`},
		{`true / false`},
	}

	for _, c := range cases {
		result, err := lang.Eval([]byte(c.input))

		assert.NoError(t, err)
		assert.Equal(t, "ERR!", result[:4])
	}
}
