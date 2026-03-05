package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Envelope struct {
	Name    string
	Payload any
}

func NewEnvelope(name string, payload any) *Envelope {
	return &Envelope{
		Name:    name,
		Payload: payload,
	}
}

func ReadParamInt(r *http.Request, key string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	value, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || value < 1 {
		return 0, errors.New("unable to parse param or it is invalid")
	}

	return value, nil
}

func WriteJSON(w http.ResponseWriter, status int, value *Envelope, headers http.Header) error {
	data := map[string]any{
		value.Name: value.Payload,
	}

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
