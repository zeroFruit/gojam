package naive

import (
	"fmt"

	"github.com/zeroFruit/jam"
)

type NaiveUploadApi struct {
	s3Service *jam.S3Service
}

func NewUploadApi(s3Service *jam.S3Service) *NaiveUploadApi {
	return &NaiveUploadApi{
		s3Service: s3Service,
	}
}

func (a *NaiveUploadApi) UploadCollection(col jam.PayloadCollection) ([]int, error) {
	uploadedKeys := make([]int, 0)
	for _, pl := range col.Payloads {
		if err := a.s3Service.Upload(pl); err != nil {
			return uploadedKeys, err
		}
		fmt.Printf("successfully upload payload key [%d]\n", pl.Key)
		uploadedKeys = append(uploadedKeys, pl.Key)
	}
	return uploadedKeys, nil
}
