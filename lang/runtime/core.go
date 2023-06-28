// TODO: Co
package runtime

import (
	"math/rand"
	"time"
)

const RETURN_KEY = "0_return"
const RAISE_KEY = "0_raise"
const SCOPE_NAME_KEY = "0_scope_name"
const SCOPE_DEPTH_KEY = "0_scope_depth"
const SCOPE_ID_KEY = "0_scope_id"
const SCOPE_FN_KEY = "0_fn"

type DataImpl interface{}

type InternalFunction func(r *Runtime, s *Scope, args ...*Instance) *Instance

type Function interface {
	// GetName() string
	// GetToken() tokens.Token
	Call(r *Runtime, s *Scope, args ...*Instance) *Instance
}

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
		return instance.Impl.(BooleanDataImpl).Value
	} else {
		return instance.Type.OnBoolean(nil, nil, instance).Impl.(BooleanDataImpl).Value
	}
}

func AsNumber(instance *Instance) float64 {
	if instance.Type == Number.Type {
		return instance.Impl.(NumberDataImpl).Value
	}
	return 0
}

func AsString(instance *Instance) string {
	if instance == nil {
		return ""
	} else if IsString(instance) {
		return instance.Impl.(StringDataImpl).Value
	} else {
		return instance.Type.OnString(nil, nil, instance).Impl.(StringDataImpl).Value
	}
}

func Variable(i *Instance) *Reference {
	return &Reference{
		Value:    i,
		Constant: false,
	}
}

func Constant(i *Instance) *Reference {
	return &Reference{
		Value:    i,
		Constant: true,
	}
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
