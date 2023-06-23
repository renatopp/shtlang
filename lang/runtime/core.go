package runtime

type DataType struct {
	Name string
}

type DataImpl interface {
	Repr() string
}

type Instance[T DataImpl] struct {
	Type DataType
	Impl T
}

type AnyInstance Instance[DataImpl]

type Function interface {
	Call(r *Runtime, args []AnyInstance) AnyInstance
}
