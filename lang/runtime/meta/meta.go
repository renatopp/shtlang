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
	Bool   MetaName = "bool"
	String MetaName = "string"
	Repr   MetaName = "repr"
	To     MetaName = "to" // <value> to <type> 		! STATIC

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
