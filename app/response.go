package app

import (
	"encoding/json"
	"net/http"
)

func unauthorized(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusUnauthorized)
}

func ok(rw http.ResponseWriter, body interface{}) {
	var (
		encodedBody []byte
		err         error
	)

	if body != nil {
		encodedBody, err = json.Marshal(body)
		if err != nil {
			internalServerError(rw, err)
			return
		}
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if encodedBody != nil {
		rw.Write(encodedBody)
	}
}

func internalServerError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
}

func badRequest(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusBadRequest)
}
