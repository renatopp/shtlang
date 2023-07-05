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

fn fib(n) {
	a, b, i := 0, 1, 2

	if n > a yield a
	if n > b yield b

	for i <= n {
		a, b = b, a + b	
		yield b
		i += 1
	}
}

fn fibx(n) {
	pipe fib(n) as x {
		yield x
	}
}



i := -1
fib(10)
| map x : {
	i += 1
	return (i, x)
}
| each x, y : print(x, y)
| reduce(List { 0, 0 }) acc, val : List { acc[0] + val[0], acc[1] + val[1] }
| map x : x[1]
| each x : print(x)
| to Boolean

pipe fibx(10) as i {
	print(i)
}

ret := aasdf?
print(ret)
print(ret!)
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
