package runtime

type Instance struct {
	Type DataType
	Impl DataImpl
}

func (i *Instance) Repr() string {
	return AsString(i.Type.OnRepr(nil, nil, i))
}
