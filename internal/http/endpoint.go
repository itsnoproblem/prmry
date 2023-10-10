package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/itsnoproblem/prmry/internal/components"
)

type HandlerFunc func(ctx context.Context, request interface{}) (response interface{}, err error)
type DecoderFunc func(ctx context.Context, request *http.Request) (decoded interface{}, err error)
type HTMXEncoderFunc func(ctx context.Context, response interface{}) (components.Component, error)
type JSONEncoderFunc func(ctx context.Context, response interface{}) (json.RawMessage, error)

type HTMXEndpoint struct {
	HandleRequest  HandlerFunc
	DecodeRequest  DecoderFunc
	EncodeResponse HTMXEncoderFunc
	RequiresAuth   bool
}

type JSONEndpoint struct {
	HandleRequest  HandlerFunc
	DecodeRequest  DecoderFunc
	EncodeResponse JSONEncoderFunc
	RequiresAuth   bool
}

func NewHTMXEndpoint(handler HandlerFunc, decoder DecoderFunc, encoder HTMXEncoderFunc, requiresAuth bool) HTMXEndpoint {
	return HTMXEndpoint{
		HandleRequest:  handler,
		DecodeRequest:  decoder,
		EncodeResponse: encoder,
		RequiresAuth:   requiresAuth,
	}
}

func NewJSONEndpoint(handler HandlerFunc, decoder DecoderFunc, encoder JSONEncoderFunc, requiresAuth bool) JSONEndpoint {
	return JSONEndpoint{
		HandleRequest:  handler,
		DecodeRequest:  decoder,
		EncodeResponse: encoder,
		RequiresAuth:   requiresAuth,
	}
}
