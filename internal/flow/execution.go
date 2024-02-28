package flow

type Execution struct {
	FlowID      string
	FlowName    string
	Model       string
	Temperature float64
	Executes    bool
	Prompt      string
}
