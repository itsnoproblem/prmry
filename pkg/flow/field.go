package flow

const (
	FieldSourceInput SourceType = "input message"
	FieldSourceFlow  SourceType = "interaction result from another flow"
)

type Field struct {
	Source SourceType
	Value  string
}

type SourceType string

func (s SourceType) String() string {
	return string(s)
}
