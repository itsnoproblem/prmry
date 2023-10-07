package api

import (
	"encoding/json"
	"net/http"
)

type renderer struct{}

func NewRenderer() *renderer {
	return &renderer{}
}

func (rnd *renderer) Render(w http.ResponseWriter, r *http.Request, data json.RawMessage) error {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}

func (rnd *renderer) RenderError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"error": err.Error(),
	}
	data, _ := json.Marshal(response)
	w.Write(data)
}
