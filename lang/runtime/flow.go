package runtime

type FlowType string

const (
	FlowReturn   FlowType = "return"
	FlowRaise    FlowType = "raise"
	FlowYield    FlowType = "yield"
	FlowBreak    FlowType = "break"
	FlowContinue FlowType = "continue"
)

type FlowInterruption struct {
	Origin *Scope
	Type   FlowType
	Value  *Instance
}
