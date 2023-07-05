package runtime

type ExecutionRecord interface{}

type BlockRecord struct {
	Current int
	Scope   *Scope
}

type IfRecord struct {
	Condition bool
	Scope     *Scope
}

type ForRecord struct {
	Scope *Scope
}

type PipeLoopRecord struct {
	Scope    *Scope
	Iterator *Instance
}
