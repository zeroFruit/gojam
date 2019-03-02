package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zeroFruit/jam/jobqueue"

	leveldbwrapper "github.com/DE-labtory/leveldb-wrapper"
	"github.com/zeroFruit/jam"
	"github.com/zeroFruit/jam/naive"
)

type runningMode string

const (
	naiveMode    runningMode = "naive"
	jobQueueMode             = "jobqueue"
)

func main() {
	os.RemoveAll(jam.DBPath)
	log.Printf("Cleaned up db files [%s]\n", jam.DBPath)

	mode := flag.String("mode", "naive", "./jam -mode <SERVER_MODE>")
	workerNum := flag.Int("worker", 5, "./jam -mode jobqueue -worker <WORKER_NUM>")
	flag.Parse()

	dbProvider := leveldbwrapper.CreateNewDBProvider(jam.DBPath)
	s3Service := jam.NewS3Service(dbProvider.GetDBHandle(*mode))

	var api jam.UploadApi

	switch runningMode(*mode) {
	case naiveMode:
		api = naiveApi(s3Service)
		log.Printf("Now server running on [localhost:8080] - [%s] mode\n", *mode)
	case jobQueueMode:
		api = jobQueueApi(s3Service, *workerNum)
		log.Printf("Now server running on [localhost:8080] - [%s] mode, [%d] workers\n", *mode, *workerNum)
	default:
		panic(fmt.Sprintf("undefineed server mode - %s", *mode))
	}

	handler := jam.NewApiHandler(api)

	http.HandleFunc("/", handler.HandleIndex)
	http.HandleFunc("/upload", handler.HandleS3Upload)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func naiveApi(s3Service *jam.S3Service) *naive.UploadApi {
	return naive.NewUploadApi(s3Service)
}

func jobQueueApi(s3Service *jam.S3Service, workerNum int) *jobqueue.UploadApi {
	dispatcher := jobqueue.NewUploadWorkDispatcher(s3Service, workerNum)
	return jobqueue.NewUploadApi(dispatcher)
}
