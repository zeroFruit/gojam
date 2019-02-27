package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
	"github.com/zeroFruit/jam"
)

var dataDirPath = "./tmp"

func cleanUp() error {
	_, err := os.Stat(dataDirPath)
	if os.IsExist(err) {
		return os.RemoveAll(dataDirPath)
	}
	return nil
}

func makeUsersTraffic(target string, users, collectionPerUser int) error {
	for u := 0; u < users; u++ {
		col := makePayloadCollection(u, collectionPerUser)
		body, err := json.Marshal(&col)
		if err != nil {
			return err
		}

		err = writeTrafficData(target, u, body)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeTrafficData(target string, user int, body []byte) error {
	err := os.MkdirAll(dataDirPath, 0755)
	if err != nil {
		return err
	}
	// Use the user id as the filename
	filename := fmt.Sprintf("%s/user%d.json", dataDirPath, user)

	// Write the JSON to the file
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		return err
	}

	// Get the absolute path to the file
	filePath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	// Print the attack target
	fmt.Println(target)

	// Print '@' followed by the absolute path to our JSON file, followed by
	// two newlines, which is the delimiter Vegeta uses
	fmt.Printf("@%s\n\n", filePath)

	return nil
}

func makePayloadCollection(user, collectionPerUser int) jam.PayloadCollection {
	pls := make([]jam.Payload, 0)

	for c := user * collectionPerUser; c < (user+1)*collectionPerUser; c++ {
		p := jam.Payload{
			Timestamp: time.Now(),
			Data: jam.Data{
				Key:     c,
				Content: c,
			},
		}
		pls = append(pls, p)
	}
	return jam.PayloadCollection{Payloads: pls}
}

func main() {
	err := cleanUp()
	if err != nil {
		panic(err)
	}

	// receive user, collection size as parameters
	var users int
	var collectionPerUser int

	// generate vegita body data

	app := cli.NewApp()
	app.Name = "jam traffic data generator"
	app.Usage = ""

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "users",
			Value:       10,
			Usage:       "number of users to simulate",
			Destination: &users,
		},
		cli.IntFlag{
			Name:        "collections",
			Value:       5,
			Usage:       "number of collection per user",
			Destination: &collectionPerUser,
		},
	}

	app.Action = func(c *cli.Context) error {
		// Combine verb and URL to a target for Vegeta
		method := c.Args().Get(0)
		url := c.Args().Get(1)

		target := fmt.Sprintf("%s %s", method, url)

		if len(target) <= 1 {
			return cli.NewExitError("usage: you must specify the target in format <METHOD> <URL>", 1)
		}

		err := makeUsersTraffic(target, users, collectionPerUser)
		if err != nil {
			return err
		}

		return nil
	}

	app.Run(os.Args)
}
