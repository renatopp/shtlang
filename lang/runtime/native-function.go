package runtime

type nativeFunction func(r *Runtime, s *Scope, args ...*Instance) *Instance

type NativeFunction struct {
	fn nativeFunction
}

func (f *NativeFunction) Call(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return f.fn(r, s, args...)
}

func CreateNativeFunction(fn nativeFunction) Function {
	return &NativeFunction{fn}
}
