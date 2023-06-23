package runtime

var BooleanType = &DataType{Name: "Boolean"}
var Boolean = Type.Create(BooleanType, true)

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

func (b BooleanImpl) Repr() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}
