package runtime

var TypeType = &DataType{Name: "Type"}
var Type = CreateType(TypeType)

func CreateType(dataType *DataType) *Instance {
	return &Instance{
		Type:  TypeType,
		Impl:  TypeImpl{DataType: *dataType},
		Const: true,
	}
}

type TypeImpl struct {
	DataType DataType
}

func (t TypeImpl) Repr() string {
	return t.DataType.Name
}
