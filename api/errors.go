package api

import (
	"encoding/json"
	"net/http"
)

type httpError struct {
	message any
	code    int
}

func (e httpError) Error() string {
	if data, err := json.Marshal(e.message); err == nil {
		return string(data)
	}
	return ""
}

func (a *API) handleError(w http.ResponseWriter, err error) {
	e, ok := err.(httpError)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(e.code)
	w.Write([]byte(e.Error()))
}
