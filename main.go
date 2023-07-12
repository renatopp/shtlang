package main

import (
	"fmt"
	"sht/lang"
	"sht/lang/runtime"
	"sht/lang/utils"
)

// var sample1 = `a++`

// var sample1 = `a!`

var sample1 = `a?`
var sample2 = `
a := List {}

a.push(1)
a.push(2)
a.push(3)
a.push(4)
a.pop()
a.pop(0)

print(a)

`

func main() {
	input := []byte(sample2)

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

	utils.PrintAst(tree)
}

func testRuntime(input []byte) {
	tree, err := lang.Parse(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	runtime := runtime.CreateRuntime()
	res, err := runtime.Run(tree)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("<<<", res)
}
