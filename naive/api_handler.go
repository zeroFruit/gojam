package naive

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/zeroFruit/jam"
)

// MaxBodyLength is 3MB
const MaxBodyLength int64 = 3 * 1024

type ApiHandler struct {
	s3Service *jam.S3Service
}

func NewApiHandler(s3Service *jam.S3Service) *ApiHandler {
	return &ApiHandler{
		s3Service: s3Service,
	}
}

func (h *ApiHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("index route accessed")
	jam.RespondWithBody(w, http.StatusOK, jam.HelloResponse{Msg: "Hello from naive server"})
}

func (h *ApiHandler) HandleS3Upload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	content := &jam.PayloadCollection{}
	if err := decodePayloadCollection(r, content); err != nil {
		jam.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	uploadedKeys, err := h.uploadCollectionToS3(*content)
	if err != nil {
		jam.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	jam.RespondWithBody(w, http.StatusOK, jam.UploadResponse{Keys: uploadedKeys})
}

func decodePayloadCollection(r *http.Request, col *jam.PayloadCollection) error {
	return json.NewDecoder(io.LimitReader(r.Body, MaxBodyLength)).Decode(col)
}

func (h *ApiHandler) uploadCollectionToS3(col jam.PayloadCollection) ([]int, error) {
	uploadedKeys := make([]int, 0)
	for _, pl := range col.Payloads {
		if err := h.s3Service.Upload(pl); err != nil {
			return uploadedKeys, err
		}
		fmt.Printf("successfully upload payload key [%d]\n", pl.Key)
		uploadedKeys = append(uploadedKeys, pl.Key)
	}
	return uploadedKeys, nil
}
