package runtime

type ExecutionState interface{}

type BlockState struct {
	Current int
	Scope   *Scope
}

type IfState struct {
	Condition bool
	Scope     *Scope
}
