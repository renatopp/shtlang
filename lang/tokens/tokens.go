package tokens

type Type string

const (
	Invalid Type = "invalid"

	// Spacing
	Eof     = "eof"
	Newline = "newline" // "\n"
	Space   = "space"   // " \t"

	// Variable-related
	Identifier = "identifier" // [a-zA-Z_$][a-zA-Z0-9_$]*
	Number     = "number"     // 123, 123.456, 123e456, -.2
	String     = "string"     // '.*'

	// Symbols
	Semicolon = "semicolon" // ";"
	Comma     = "comma"     // ","
	Colon     = "colon"     // ":"
	Bang      = "bang"      // "!"
	Question  = "question"  // "?"
	Dot       = "dot"       // "."
	Backslash = "backslash" // "\"
	At        = "at"        // "@"
	Hash      = "hash"      // "#"
	Percent   = "percent"   // "%"
	Caret     = "caret"     // "^"
	Ampersand = "ampersand" // "&"
	Pipe      = "pipe"      // "|"
	Plus      = "plus"      // "+"
	Minus     = "minus"     // "-"
	Asterisk  = "asterisk"  // "*"
	Slash     = "slash"     // "/"
	Greater   = "greater"   // ">"
	Less      = "less"      // "<"
	Equal     = "equal"     // "="
	Tilde     = "tilde"     // "~"
	Lbrace    = "lbrace"    // "{"
	Rbrace    = "rbrace"    // "}"
	Lparen    = "lparen"    // "("
	Rparen    = "rparen"    // ")"
	Lbracket  = "lbracket"  // "["
	Rbracket  = "rbracket"  // "]"
)
