package runtime

type DataImpl interface{}

type InternalFunction func(r *Runtime, args ...*Instance) *Instance

type Function interface {
	Call(r *Runtime, args ...*Instance) *Instance
}

func IsBool(instance *Instance) bool {
	// return instance.Type == Boolean.Type
	return false
}

func IsString(instance *Instance) bool {
	return instance.Type == String.Type
}

func AsBool(instance *Instance) bool {
	if instance == nil {
		return false
	} else if IsBool(instance) {
		return instance.Impl.(BooleanDataImpl).Value
	} else {
		return instance.Type.OnBoolean(nil, instance).Impl.(BooleanDataImpl).Value
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
		return instance.Type.OnString(nil, instance).Impl.(StringDataImpl).Value
	}
}
