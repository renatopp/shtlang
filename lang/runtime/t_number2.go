package runtime

// import (
// 	"fmt"
// 	"math"
// 	"sht/lang/ast"
// 	"sht/lang/runtime/meta"
// )

// var Number = _setupNumber()

// type NumberInfo struct {
// 	Instance *Instance
// 	Type     *DataType

// 	ZERO *Instance
// 	ONE  *Instance
// }

// type NumberImpl struct {
// 	Value float64
// }

// func _setupNumber() *NumberInfo {
// 	dataType := &DataType{
// 		Name:        "Number",
// 		Properties:  map[string]ast.Node{},
// 		StaticFns:   map[string]Function{},
// 		InstanceFns: map[string]Function{},
// 		Meta:        map[meta.MetaName]Function{},
// 	}

// 	n := &NumberInfo{
// 		Instance: Type.Create(dataType, true),
// 		Type:     dataType,
// 	}

// 	n.ZERO = n.Create(0, true)
// 	n.ONE = n.Create(1, true)

// 	dataType.Meta[meta.SetProperty] = n.invalid(string(meta.SetProperty))
// 	dataType.Meta[meta.GetProperty] = n.invalid(string(meta.GetProperty))
// 	dataType.Meta[meta.SetItem] = n.invalid(string(meta.SetItem))
// 	dataType.Meta[meta.GetItem] = n.invalid(string(meta.GetItem))
// 	dataType.Meta[meta.Call] = n.invalid(string(meta.Call))

// 	dataType.Meta[meta.Boolean] = CreateNativeFunction(n.MetaBoolean)
// 	dataType.Meta[meta.String] = CreateNativeFunction(n.MetaString)
// 	dataType.Meta[meta.Repr] = CreateNativeFunction(n.MetaRepr)

// 	dataType.Meta[meta.Add] = CreateNativeFunction(n.MetaAdd)
// 	dataType.Meta[meta.Sub] = CreateNativeFunction(n.MetaSub)
// 	dataType.Meta[meta.Mul] = CreateNativeFunction(n.MetaMul)
// 	dataType.Meta[meta.Div] = CreateNativeFunction(n.MetaDiv)
// 	dataType.Meta[meta.IntDiv] = CreateNativeFunction(n.MetaIntDiv)
// 	dataType.Meta[meta.Mod] = CreateNativeFunction(n.MetaMod)
// 	dataType.Meta[meta.Pow] = CreateNativeFunction(n.MetaPow)
// 	dataType.Meta[meta.Eq] = CreateNativeFunction(n.MetaEq)
// 	dataType.Meta[meta.Neq] = CreateNativeFunction(n.MetaNeq)
// 	dataType.Meta[meta.Gt] = CreateNativeFunction(n.MetaGt)
// 	dataType.Meta[meta.Lt] = CreateNativeFunction(n.MetaLt)
// 	dataType.Meta[meta.Gte] = CreateNativeFunction(n.MetaGte)
// 	dataType.Meta[meta.Lte] = CreateNativeFunction(n.MetaLte)
// 	dataType.Meta[meta.Pos] = CreateNativeFunction(n.MetaPos)
// 	dataType.Meta[meta.Neg] = CreateNativeFunction(n.MetaNeg)
// 	dataType.Meta[meta.Not] = CreateNativeFunction(n.MetaNot)
// 	dataType.Meta[meta.PostInc] = CreateNativeFunction(n.MetaPostInc)
// 	dataType.Meta[meta.PostDec] = CreateNativeFunction(n.MetaPostDec)

// 	return n
// }

// // ----------------------------------------------------------------------------
// // Number Implementation
// // ----------------------------------------------------------------------------
// func (n NumberImpl) Repr() string {
// 	if math.Mod(n.Value, 1.0) == 0 {
// 		return fmt.Sprintf("%.0f", n.Value)
// 	} else {
// 		return fmt.Sprintf("%f", n.Value)
// 	}
// }

// // ----------------------------------------------------------------------------
// // Number Info
// // ----------------------------------------------------------------------------
// func (n *NumberInfo) Create(value float64, constant bool) *Instance {
// 	return &Instance{
// 		Type:  n.Type,
// 		Impl:  NumberImpl{Value: value},
// 		Const: constant,
// 	}
// }

// func (n *NumberInfo) val(instance *Instance) float64 {
// 	return instance.Impl.(NumberImpl).Value
// }

// func (n *NumberInfo) MetaPos(r *Runtime, args ...*Instance) *Instance {
// 	return args[0]
// }

// func (n *NumberInfo) MetaNeg(r *Runtime, args ...*Instance) *Instance {
// 	return n.Create(-n.val(args[0]), false)
// }

// func (n *NumberInfo) MetaNot(r *Runtime, args ...*Instance) *Instance {
// 	this := n.val(args[0])
// 	return Boolean.Create(this == 0, false)
// }

// func (n *NumberInfo) MetaBoolean(r *Runtime, args ...*Instance) *Instance {
// 	this := n.val(args[0])
// 	return Boolean.Create(this != 0, false)
// }

// func (n *NumberInfo) MetaString(r *Runtime, args ...*Instance) *Instance {
// 	return String.Create(args[0].Impl.Repr(), false)
// }

// func (n *NumberInfo) MetaRepr(r *Runtime, args ...*Instance) *Instance {
// 	return String.Create(args[0].Impl.Repr(), false)
// }

// func (n *NumberInfo) MetaAdd(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("+", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return n.Create(this+other, false)
// }

// func (n *NumberInfo) MetaSub(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("-", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return n.Create(this-other, false)
// }

// func (n *NumberInfo) MetaMul(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("*", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return n.Create(this*other, false)
// }

// func (n *NumberInfo) MetaDiv(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("/", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return n.Create(this/other, false)
// }

// func (n *NumberInfo) MetaIntDiv(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("/", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return n.Create(math.Floor(this/other), false)
// }

// func (n *NumberInfo) MetaMod(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("/", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return n.Create(math.Mod(this, other), false)
// }

// func (n *NumberInfo) MetaPow(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("/", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return n.Create(math.Pow(this, other), false)
// }

// func (n *NumberInfo) MetaEq(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return Boolean.FALSE
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return Boolean.Create(this == other, false)
// }

// func (n *NumberInfo) MetaNeq(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return Boolean.TRUE
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return Boolean.Create(this != other, false)
// }

// func (n *NumberInfo) MetaGt(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType(">", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return Boolean.Create(this > other, false)
// }

// func (n *NumberInfo) MetaLt(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("<", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return Boolean.Create(this < other, false)
// }

// func (n *NumberInfo) MetaGte(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType(">=", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return Boolean.Create(this >= other, false)
// }

// func (n *NumberInfo) MetaLte(r *Runtime, args ...*Instance) *Instance {
// 	if args[0].Type != args[1].Type {
// 		return InvalidOperationType("<=", args[0], args[1])
// 	}

// 	this := n.val(args[0])
// 	other := n.val(args[1])

// 	return Boolean.Create(this <= other, false)
// }

// func (n *NumberInfo) MetaPostInc(r *Runtime, args ...*Instance) *Instance {
// 	impl := args[0].Impl.(NumberImpl)
// 	old := impl.Value
// 	impl.Value += 1
// 	return n.Create(old, false)
// }

// func (n *NumberInfo) MetaPostDec(r *Runtime, args ...*Instance) *Instance {
// 	impl := args[0].Impl.(NumberImpl)
// 	old := impl.Value
// 	impl.Value -= 1
// 	return n.Create(old, false)
// }

// func (n *NumberInfo) invalid(action string) Function {
// 	return CreateNativeFunction(func(r *Runtime, args ...*Instance) *Instance {
// 		return NotImplemented("meta", args[0])
// 	})
// }
