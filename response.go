package jam

import (
	"encoding/json"
	"net/http"
)

type HelloResponse struct {
	Msg string
}

type UploadResponse struct {
	Keys []int
}

func RespondWithError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

func RespondWithBody(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
