package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	leveldbwrapper "github.com/DE-labtory/leveldb-wrapper"
	"github.com/zeroFruit/jam"
	"github.com/zeroFruit/jam/naive"
)

type runningMode string

const (
	naiveServer runningMode = "naive"
)

func main() {
	os.RemoveAll(jam.DBPath)
	log.Printf("Cleaned up db files [%s]\n", jam.DBPath)

	mode := flag.String("mode", "naive", "./jam -mode <SERVER_MODE>")

	dbProvider := leveldbwrapper.CreateNewDBProvider(jam.DBPath)

	s3Service := jam.NewS3Service(dbProvider.GetDBHandle(*mode))

	var api jam.UploadApi

	switch runningMode(*mode) {
	case naiveServer:
		api = naive.NewUploadApi(s3Service)
	default:
		panic(fmt.Sprintf("undefineed server mode - %s", *mode))
	}

	handler := jam.NewApiHandler(api)

	http.HandleFunc("/", handler.HandleIndex)
	http.HandleFunc("/upload", handler.HandleS3Upload)

	log.Printf("Now server running on [localhost:8080] - [%s] mode\n", *mode)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
