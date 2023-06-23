package runtime

type DataType struct {
	Name string
}

type DataImpl interface {
	Repr() string
}

type Instance struct {
	Type  *DataType
	Impl  DataImpl
	Const bool
}

type Function interface {
	Call(r *Runtime, args []Instance) Instance
}
