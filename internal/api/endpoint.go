package api

import (
	"context"
	"net/http"
)

type Endpoint interface {
	HandleRequest(ctx context.Context, request interface{}) (response interface{}, err error)
	DecodeRequest(ctx context.Context, request *http.Request) (decoded interface{}, err error)
	FormatResponse(response interface{}) (formatted interface{}, err error)
}

type endpoint struct {
	h HandlerFunc
	d DecoderFunc
	f FormatterFunc
}

func (e endpoint) HandleRequest(ctx context.Context, request interface{}) (response interface{}, err error) {
	return e.h(ctx, request)
}

func (e endpoint) DecodeRequest(ctx context.Context, request *http.Request) (decoded interface{}, err error) {
	return e.d(ctx, request)
}

func (e endpoint) FormatResponse(response interface{}) (formatted interface{}, err error) {
	return e.f(response)
}

func NewEndpoint(h HandlerFunc, d DecoderFunc, f FormatterFunc) Endpoint {
	return endpoint{
		h: h,
		d: d,
		f: f,
	}
}

type Response interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

type HandlerFunc func(ctx context.Context, request interface{}) (response interface{}, err error)
type DecoderFunc func(ctx context.Context, request *http.Request) (decoded interface{}, err error)
type FormatterFunc func(response interface{}) (formatted interface{}, err error)
