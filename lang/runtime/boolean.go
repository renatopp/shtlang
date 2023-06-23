package runtime

var BooleanType = &DataType{Name: "Boolean"}
var Boolean = CreateType(BooleanType)

func CreateBoolean(value bool, constant bool) *Instance {
	return &Instance{
		Type:  BooleanType,
		Impl:  BooleanImpl{Value: value},
		Const: constant,
	}
}

type BooleanImpl struct {
	Value bool
}

func (n BooleanImpl) Repr() string {
	if n.Value {
		return "true"
	} else {
		return "false"
	}
}
