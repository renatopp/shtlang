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

data Notes {
	a = (1, 2, 3, 4)

	fn create() {
		return Notes { a: (0, 1, 2) }
	}

	fn reverse(this) {
		i := len(this.a) - 1
		for i >= 0 {
			yield this.a[i]
			i -= 1
		}
	}

	on len(this) {
		return len(this.a)
	}
	
	on iter(this) {
		pipe this.a as item {
			yield item
		}
	}

	on to(iter) {
		return Notes { a: iter | to Tuple }
	}
}

notes := List { 1, 2, 3, 4 } | to Notes

notes.a

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
	res := runtime.Run(tree)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("<<<", res)
}
