package main

import (
	"fmt"
	"sht/lang"
	"sht/lang/ast"
	"sht/lang/runtime"
	"strings"
)

var sample1 = `1.2e2`

func main() {
	input := []byte(sample1)

	fmt.Println("")
	fmt.Println("-----------------------------------")
	fmt.Println("              oh SHT!              ")
	fmt.Println("-----------------------------------")
	fmt.Println("")
	testTokens(input)
	fmt.Println("")
	fmt.Println("-----------------------------------")
	fmt.Println("")
	testParser(input)
	fmt.Println("")
	fmt.Println("-----------------------------------")
	fmt.Println("")
	testRuntime(input)
	fmt.Println("")
	fmt.Println("-----------------------------------")
}

func testTokens(input []byte) {
	tokens, err := lang.Tokenize(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, token := range tokens {
		fmt.Println(token.Pretty())
	}
}

func testParser(input []byte) {
	tree, err := lang.Parse(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	tree.Traverse(0, func(level int, node ast.Node) {
		fmt.Println(strings.Repeat("  ", level) + node.String())
	})
}

func testRuntime(input []byte) {
	tree, err := lang.Parse(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	runtime := runtime.CreateRuntime()
	res := runtime.Eval(tree, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
