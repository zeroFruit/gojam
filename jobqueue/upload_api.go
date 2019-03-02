package jobqueue

import (
	"fmt"

	"github.com/zeroFruit/jam"
)

type UploadApi struct {
	dispatcher *UploadWorkDispatcher
}

func NewUploadApi(dispatcher *UploadWorkDispatcher) *UploadApi {
	return &UploadApi{
		dispatcher: dispatcher,
	}
}

func (a *UploadApi) UploadCollection(col jam.PayloadCollection) ([]int, error) {
	result := a.dispatcher.Schedule(col)
	if len(result) != len(col.Payloads) {
		return result, fmt.Errorf("error occurred while uploading data - %v", result)
	}

	return result, nil
}
