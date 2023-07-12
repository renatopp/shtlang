// TODO: Co
package runtime

import (
	"math/rand"
	"time"
)

type R *Runtime
type S *Scope
type I *Instance
type F MetaFunction

type DataImpl interface{}

type MetaFunction func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance

func IsBool(instance *Instance) bool {
	return instance.Type == Boolean.Type
}

func IsString(instance *Instance) bool {
	return instance.Type == String.Type
}

func IsNumber(instance *Instance) bool {
	return instance.Type == Number.Type
}

func AsBool(instance *Instance) bool {
	if instance == nil {
		return false
	} else if IsBool(instance) {
		return instance.Impl.(*BooleanDataImpl).Value
	} else {
		return instance.OnBoolean(nil, nil).Impl.(*BooleanDataImpl).Value
	}
}

func AsNumber(instance *Instance) float64 {
	if instance.Type == Number.Type {
		return instance.Impl.(*NumberDataImpl).Value
	}
	return 0
}

func AsInteger(instance *Instance) int {
	return int(AsNumber(instance))
}

func AsString(instance *Instance) string {
	if instance == nil {
		return ""
	} else if IsString(instance) {
		return instance.Impl.(*StringDataImpl).Value
	} else {
		return instance.OnString(nil, nil).Impl.(*StringDataImpl).Value
	}
}

func AsFunction(instance *Instance) MetaFunction {
	return instance.Impl.(*FunctionDataImpl).Call
}

func Variable(i *Instance) *Instance {
	i.Constant = false
	return i
}

func Constant(i *Instance) *Instance {
	i.Constant = true
	return i
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Id() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

var DoneFn = Function.CreateNative("done", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return Iteration.DONE
})

var ThrowFn = Function.CreateNative("throw", []*FunctionParam{
	{Name: "message"},
}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.Create(s, AsString(args[0])), s)
})

var GetFirstFn = fn("getFirst", p("tuple")).as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if args[0].IsTuple() {
		return args[0].AsTuple().Values[0]
	} else {
		return args[0]
	}
})
