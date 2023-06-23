package runtime

type nativeFunction func(r *Runtime, args ...*Instance) *Instance

type NativeFunction struct {
	fn nativeFunction
}

func (f *NativeFunction) Call(r *Runtime, args ...*Instance) *Instance {
	return f.fn(r, args...)
}

func CreateNativeFunction(fn nativeFunction) Function {
	return &NativeFunction{fn}
}
