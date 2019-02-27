package jam

import (
	"fmt"

	leveldbwrapper "github.com/DE-labtory/leveldb-wrapper"
)

func ValidateBenchmarkResult(db *leveldbwrapper.DBHandle, users, collectionPerUser int) {
	errList := make([]error, 0)
	for i := 0; i < users*collectionPerUser; i++ {
		if _, err := db.Get(ByteArray(i)); err != nil {
			errList = append(errList, err)
		}
	}
	printResult(errList, users*collectionPerUser)
}

func printResult(errList []error, testCount int) {
	for i, err := range errList {
		fmt.Printf("error [%d] - %s\n", i, err.Error())
	}
	fmt.Printf("error ratio: %d / %d", len(errList), testCount)
}
