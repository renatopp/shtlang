package repl

func isMatchingPair(character1, character2 rune) bool {
	return (character1 == '(' && character2 == ')') ||
		(character1 == '{' && character2 == '}') ||
		(character1 == '[' && character2 == ']')
}

func areBracketsBalanced(expression string) bool {
	stack := []rune{}

	inSingleQuotes := false
	inDoubleQuotes := false
	inSingleLineComment := false
	inHashComment := false

	for _, character := range expression {
		if inSingleLineComment {
			if character == '\n' {
				inSingleLineComment = false
			}
			continue
		}

		if inHashComment {
			if character == '\n' {
				inHashComment = false
			}
			continue
		}

		if inSingleQuotes {
			if character == '\'' {
				inSingleQuotes = false
			}
			continue
		}

		if inDoubleQuotes {
			if character == '"' {
				inDoubleQuotes = false
			}
			continue
		}

		switch character {
		// case '/':
		// 	if strings.HasPrefix(expression, "//") {
		// 		inSingleLineComment = true
		// 		expression = expression[1:] // shift the string
		// 	}
		case '#':
			inHashComment = true
		case '\'':
			inSingleQuotes = true
		case '"':
			inDoubleQuotes = true
		case '[', '{', '(':
			stack = append(stack, character)
		case ']', '}', ')':
			if len(stack) == 0 {
				return false
			}
			if !isMatchingPair(stack[len(stack)-1], character) {
				return false
			}
			stack = stack[:len(stack)-1]
		}

		if len(expression) > 0 {
			expression = expression[1:] // shift the string
		}
	}

	return len(stack) == 0
}

// func main() {
// 	expressions := []string{
// 		`if (x[i] > 5) { y = "hello (world)"; z = '#'; }`,
// 		`# this is comment with brackets ({[]})`,
// 		`this is an open bracket [`,
// 		`"this string has open bracket but not closed [`,
// 	}

// 	for _, expression := range expressions {
// 		if areBracketsBalanced(expression) {
// 			fmt.Println("Balanced")
// 		} else {
// 			fmt.Println("Not Balanced")
// 		}
// 	}
// }
