package jam

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// MaxBodyLength is 3MB
const MaxBodyLength int64 = 3 * 1024

type UploadApi interface {
	UploadCollection(col PayloadCollection) ([]int, error)
}

type ApiHandler struct {
	uploadApi UploadApi
}

func NewApiHandler(uploadApi UploadApi) *ApiHandler {
	return &ApiHandler{
		uploadApi: uploadApi,
	}
}

func (h *ApiHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("index route accessed")
	RespondWithBody(w, http.StatusOK, HelloResponse{Msg: "Hello from naive server"})
}

func (h *ApiHandler) HandleS3Upload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	content := &PayloadCollection{}
	if err := decodePayloadCollection(r, content); err != nil {
		RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	uploadedKeys, err := h.uploadApi.UploadCollection(*content)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	RespondWithBody(w, http.StatusOK, UploadResponse{Keys: uploadedKeys})
}

func decodePayloadCollection(r *http.Request, col *PayloadCollection) error {
	return json.NewDecoder(io.LimitReader(r.Body, MaxBodyLength)).Decode(col)
}
