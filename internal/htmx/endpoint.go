package htmx

import (
	"context"
	"net/http"

	"github.com/itsnoproblem/prmry/internal/components"
)

type HandlerFunc func(ctx context.Context, request interface{}) (response interface{}, err error)
type DecoderFunc func(ctx context.Context, request *http.Request) (decoded interface{}, err error)
type EncoderFunc func(ctx context.Context, response interface{}) (component components.Component, err error)

type Endpoint struct {
	HandleRequest  HandlerFunc
	DecodeRequest  DecoderFunc
	EncodeResponse EncoderFunc
	RequiresAuth   bool
}

func NewEndpoint(handler HandlerFunc, decoder DecoderFunc, encoder EncoderFunc, requiresAuth bool) Endpoint {
	return Endpoint{
		HandleRequest:  handler,
		DecodeRequest:  decoder,
		EncodeResponse: encoder,
		RequiresAuth:   requiresAuth,
	}
}
