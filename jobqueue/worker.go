package jobqueue

import (
	"log"

	"github.com/zeroFruit/jam"
)

// Job is value object which represents worker's job
type Job struct {
	Request interface{}
}

// JobChannel is channel which workers get its job
type JobChannel chan Job

// WorkerQueue have available worker who done its own task
type WorkerQueue chan *UploadWorker

// RespChannel is channel which sends response of worker's task result
type RespChannel chan interface{}

type UploadWorker struct {
	id        int
	s3Service *jam.S3Service
}

func NewUploadWorker(id int, s3Service *jam.S3Service) *UploadWorker {
	return &UploadWorker{
		id:        id,
		s3Service: s3Service,
	}
}

func (w *UploadWorker) Start(pl jam.Payload, respChan RespChannel, workerQueue WorkerQueue) {
	go func(pl jam.Payload) {
		err := w.s3Service.Upload(pl)
		if err != nil {
			respChan <- err
			workerQueue <- w
			return
		}

		log.Printf("successfully upload payload key [%d]\n", pl.Key)
		respChan <- pl
		workerQueue <- w
	}(pl)
}

type UploadWorkDispatcher struct {
	workerNum   int
	workerQueue WorkerQueue
	jobChan     JobChannel
	respChan    RespChannel
}

func NewUploadWorkDispatcher(s3Service *jam.S3Service, workerNum int) *UploadWorkDispatcher {
	workerQueue := make(chan *UploadWorker, workerNum)
	for i := 0; i < workerNum; i++ {
		worker := NewUploadWorker(i, s3Service)
		workerQueue <- worker
	}

	return &UploadWorkDispatcher{
		workerNum:   workerNum,
		workerQueue: workerQueue,
		jobChan:     make(chan Job, workerNum),
		respChan:    make(chan interface{}, workerNum),
	}
}

func (d *UploadWorkDispatcher) Schedule(col jam.PayloadCollection) []int {
	go func(col jam.PayloadCollection) {
		for _, pl := range col.Payloads {
			worker := <-d.workerQueue
			worker.Start(pl, d.respChan, d.workerQueue)
		}
	}(col)

	i := 0
	result := make([]int, 0)
	for {
		select {
		case resp := <-d.respChan:
			pl, ok := resp.(jam.Payload)
			if ok {
				result = append(result, pl.Key)
			}

			if i++; i < len(col.Payloads) {
				continue
			}

			return result
		}
	}
}
