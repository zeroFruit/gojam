package jam

import (
	"math/rand"
	"strconv"
	"time"

	leveldbwrapper "github.com/DE-labtory/leveldb-wrapper"
)

type PayloadCollection struct {
	Payloads []Payload `json:"payloads"`
}

type Data struct {
	Key     int `json:"key"`
	Content int `json:"content"`
}

type Payload struct {
	Timestamp time.Time `json:"timestamp"`
	Data      `json:"data"`
}

type Bucket struct {
	Name string
}

type S3Service struct {
	db *leveldbwrapper.DBHandle
}

func NewS3Service(db *leveldbwrapper.DBHandle) *S3Service {
	return &S3Service{
		db: db,
	}
}

func (s S3Service) Upload(pl Payload) error {
	if err := s.db.Put(ByteArray(pl.Key), ByteArray(pl.Content), true); err != nil {
		return err
	}

	return nil
}

func networkLatency() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(5)
}

func ByteArray(k int) []byte {
	return []byte(strconv.Itoa(k))
}
