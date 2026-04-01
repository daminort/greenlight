package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func ReadParamInt(r *http.Request, key string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	value, err := strconv.ParseInt(params.ByName(key), 10, 64)
	if err != nil || value < 1 {
		return 0, errors.New("unable to parse param or it is invalid")
	}

	return value, nil
}

func ReadParamString(r *http.Request, key string) string {
	params := httprouter.ParamsFromContext(r.Context())
	value := params.ByName(key)

	return value
}
