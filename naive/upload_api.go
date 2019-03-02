package naive

import (
	"log"

	"github.com/zeroFruit/jam"
)

type UploadApi struct {
	s3Service *jam.S3Service
}

func NewUploadApi(s3Service *jam.S3Service) *UploadApi {
	return &UploadApi{
		s3Service: s3Service,
	}
}

func (a *UploadApi) UploadCollection(col jam.PayloadCollection) ([]int, error) {
	uploadedKeys := make([]int, 0)
	for _, pl := range col.Payloads {
		if err := a.s3Service.Upload(pl); err != nil {
			return uploadedKeys, err
		}
		log.Printf("successfully upload payload key [%d]\n", pl.Key)
		uploadedKeys = append(uploadedKeys, pl.Key)
	}
	return uploadedKeys, nil
}
