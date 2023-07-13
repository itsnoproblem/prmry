package components

type ErrorView struct {
	Error string
	Code  int
	BaseComponent
}

func NewErrorView(msg string, code int) ErrorView {
	return ErrorView{
		Error: msg,
		Code:  code,
	}
}
