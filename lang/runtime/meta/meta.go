package meta

type MetaName string

var (
	// Setters and getters
	SetProperty MetaName = "set"
	GetProperty MetaName = "get"
	SetItem     MetaName = "setItem"
	GetItem     MetaName = "getItem"

	// Calling
	New  MetaName = "new"
	Call MetaName = "call"

	// Convertions
	Boolean MetaName = "boolean"
	String  MetaName = "string"
	Repr    MetaName = "repr"
	To      MetaName = "to" // <value> to <type> 		! STATIC

	// Other meta
	Iter MetaName = "iter" // for i in x
	Bang MetaName = "bang" // !

	// Operators
	Add     MetaName = "add"     // +
	Sub     MetaName = "sub"     // -
	Mul     MetaName = "mul"     // *
	Div     MetaName = "div"     // /
	IntDiv  MetaName = "intDiv"  // //
	Mod     MetaName = "mod"     // %
	Pow     MetaName = "pow"     // **
	Eq      MetaName = "eq"      // ==
	Neq     MetaName = "neq"     // !=
	Gt      MetaName = "gt"      // >
	Lt      MetaName = "lt"      // <
	Gte     MetaName = "gte"     // >=
	Lte     MetaName = "lte"     // <=
	Pos     MetaName = "pos"     // +
	Neg     MetaName = "neg"     // -
	Not     MetaName = "not"     // !
	PostInc MetaName = "postInc" // ++
	PostDec MetaName = "postDec" // --
	Concat  MetaName = "concat"  // ..
)

func FromUnaryOperator(op string) MetaName {
	switch op {
	case "+":
		return Pos
	case "-":
		return Neg
	case "!":
		return Not
	}

	return ""
}

func FromBinaryOperator(op string) MetaName {
	switch op {
	case "+":
		return Add
	case "-":
		return Sub
	case "*":
		return Mul
	case "/":
		return Div
	case "//":
		return IntDiv
	case "%":
		return Mod
	case "**":
		return Pow
	case "==":
		return Eq
	case "!=":
		return Neq
	case ">":
		return Gt
	case "<":
		return Lt
	case ">=":
		return Gte
	case "<=":
		return Lte
	case "++":
		return PostInc
	case "--":
		return PostDec
	case "..":
		return Concat
	}

	return ""
}
