package naive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/zeroFruit/jam"
)

var localAddr = "http://127.0.0.1:8080"

func BenchmarkApiHandler(b *testing.B) {
	os.RemoveAll(jam.DBPath)

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	b.ResetTimer()

	errs := make([]error, 0)

	for n := 0; n < b.N; n++ {
		col := createPayloadCollection(n)
		if err := sendUploadRequest(client, localAddr, col); err != nil {
			errs = append(errs, err)
		}
	}
	fmt.Printf("b.N - %d\n", b.N)
	printErrors(errs)
}

func createPayloadCollection(n int) *jam.PayloadCollection {
	return &jam.PayloadCollection{
		Payloads: []jam.Payload{
			{
				Timestamp: time.Now(),
				Data:      jam.Data{Key: n, Content: n},
			},
		},
	}
}

func sendUploadRequest(client *http.Client, addr string, col *jam.PayloadCollection) error {
	b, err := json.Marshal(col)
	if err != nil {
		return err
	}

	res, err := client.Post(fmt.Sprintf("%s/upload", addr), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("response status code is not ok")
	}

	return nil
}

func printErrors(errs []error) {
	if len(errs) == 0 {
		fmt.Println("no error")
		return
	}

	for i, err := range errs {
		fmt.Printf("error [%d] - %s\n", i, err.Error())
	}
}
