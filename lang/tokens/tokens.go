package tokens

type Type string

const (
	_ Type = ""

	// Spacing
	Eof     = "eof"
	Newline = "newline" // "\n"

	// Variable-related
	Keyword    = "keyword"
	Identifier = "identifier" // [a-zA-Z_$][a-zA-Z0-9_$]*
	Number     = "number"     // 123, 123.456, 123e456, -.2
	String     = "string"     // '.*'

	// Operators
	Operator   = "operator"   // +, -, *, /, //, %, **, ++, --, <, <=, >, >=, ==, !=, .., ??
	Assignment = "assignment" // =, +=, -=, *=, /=, //=,

	// Separators
	Semicolon = "semicolon" // ";"
	Comma     = "comma"     // ","
	Colon     = "colon"     // ":"
	Dot       = "dot"       // "."

	// Special
	Bang     = "bang"     // "!"
	Question = "question" // "?"
	At       = "at"       // "@"
	Pipe     = "pipe"     // "|"
	Arrow    = "arrow"    // "=>"
	Spread   = "spread"   // "..."

	// Blocks
	Lbrace   = "lbrace"   // "{"
	Rbrace   = "rbrace"   // "}"
	Lparen   = "lparen"   // "("
	Rparen   = "rparen"   // ")"
	Lbracket = "lbracket" // "["
	Rbracket = "rbracket" // "]"
)

func JoinTypes(types ...Type) string {
	str := ""
	for _, t := range types {
		str += string(t) + ", "
	}
	return str[:len(str)-2]
}
