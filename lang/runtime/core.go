package runtime

type DataType struct {
	Name string
}

type DataImpl interface {
	Repr() string
}

type Instance struct {
	Type *DataType
	Impl DataImpl
}

type Function interface {
	Call(r *Runtime, args []Instance) Instance
}

// ----------------------------------------------------------------------------

func CreateNumber(value float64) *Instance {
	return &Instance{
		Type: NumberType,
		Impl: NumberImpl{Value: value},
	}
}
