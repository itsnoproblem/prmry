package flow

const (
	FieldSourceInput      SourceType = "input message"
	FieldSourceInputTag   SourceType = "input tag"
	FieldSourceFlow       SourceType = "interaction result from another flow"
	FieldSourceModeration SourceType = "moderation result"
)

type Field struct {
	Source SourceType
	Value  string
}

type SourceType string

func (s SourceType) String() string {
	return string(s)
}

func SupportedFields() map[string]string {
	return map[string]string{
		FieldSourceInput.String():      "The Input Message",
		FieldSourceInputTag.String():   "An Input Tag",
		FieldSourceFlow.String():       "Output from another Flow",
		FieldSourceModeration.String(): "A Moderation result (of the input message)",
	}
}
