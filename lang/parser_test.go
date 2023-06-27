package lang

import (
	"sht/lang/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLetAssignment(t *testing.T) {
	input := `let x = 1`
	tree, _ := Parse([]byte(input))
	node := tree.Children()[0].(*ast.Assignment)
	assert.NotEqual(t, node, nil)
	assert.Equal(t, node.Definition, true)
	assert.Equal(t, node.Constant, false)
	assert.NotEqual(t, node.Identifier.(*ast.Identifier), nil)
	assert.Equal(t, node.Identifier.(*ast.Identifier).Value, "x")
	assert.NotEqual(t, node.Expression.(*ast.Number), nil)
	assert.Equal(t, node.Expression.(*ast.Number).Value, float64(1))
}

func TestConstAssignment(t *testing.T) {
	input := `const x = 1`
	tree, _ := Parse([]byte(input))
	node := tree.Children()[0].(*ast.Assignment)
	assert.NotEqual(t, node, nil)
	assert.Equal(t, node.Definition, true)
	assert.Equal(t, node.Constant, true)
	assert.NotEqual(t, node.Identifier.(*ast.Identifier), nil)
	assert.Equal(t, node.Identifier.(*ast.Identifier).Value, "x")
	assert.NotEqual(t, node.Expression.(*ast.Number), nil)
	assert.Equal(t, node.Expression.(*ast.Number).Value, float64(1))
}

func TestAssignmentExpression(t *testing.T) {
	input := `x += 1`
	tree, _ := Parse([]byte(input))
	node := tree.Children()[0].(*ast.Assignment)
	assert.NotEqual(t, node, nil)
	assert.Equal(t, node.Definition, false)
	assert.Equal(t, node.Constant, false)
	assert.Equal(t, node.Literal, "+=")
	assert.NotEqual(t, node.Identifier.(*ast.Identifier), nil)
	assert.Equal(t, node.Identifier.(*ast.Identifier).Value, "x")
	assert.NotEqual(t, node.Expression.(*ast.Number), nil)
	assert.Equal(t, node.Expression.(*ast.Number).Value, float64(1))
}
