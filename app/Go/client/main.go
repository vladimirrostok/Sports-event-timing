package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	checkpoint_controller "sports/backend/srv/controllers/checkpoint"
	result_controller "sports/backend/srv/controllers/result"
	sportsmen_controller "sports/backend/srv/controllers/sportsmen"
	"strconv"
	"time"
)

type id struct {
	ID string `json:"id"`
}

func main() {
	addr := "https://localhost:8000"
	currentNum := uint32(1)

	// Provide the default source to a deterministic state.
	rand.Seed(time.Now().UnixNano())

	// Disable cert verification to use self-signed certificates for internal service needs.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Create new checkpoint.
	checkpointRequestBody := checkpoint_controller.NewCheckpointRequest{
		Name: "corridor",
	}

	requestBody, err := json.Marshal(checkpointRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(
		addr+"/checkpoints",
		"application/json; charset=UTF-8",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	checkPointID := id{}
	err = json.Unmarshal(data, &checkPointID)
	if err != nil {
		log.Fatal(err)
	}

	// Generate new results and finished results.
	for _ = range time.Tick(3 * time.Second) {
		go func(){
		currentNum++

		// Add new sportsmen
		newSportsmenRequest := sportsmen_controller.NewSportsmenRequest{
			FirstName:   fmt.Sprintf("Name%s", strconv.Itoa(int(currentNum))),
			LastName:    fmt.Sprintf("Lastname%s", strconv.Itoa(int(currentNum))),
			StartNumber: uint32(rand.Intn(1000)),
		}

		requestBody, err := json.Marshal(newSportsmenRequest)
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.Post(
			addr+"/sportsmens",
			"application/json; charset=UTF-8",
			bytes.NewReader(requestBody),
		)
		if err != nil {
			log.Fatal(err)
		}

		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		sportsmenID := id{}
		err = json.Unmarshal(data, &sportsmenID)
		if err != nil {
			log.Fatal(err)
		}

		// Add new result
		newResultRequest := result_controller.NewResultRequest{
			SportsmenID:  sportsmenID.ID,
			CheckpointID: checkPointID.ID,
			Time:         makeTimestamp(),
		}

		requestBody, err = json.Marshal(newResultRequest)
		if err != nil {
			log.Fatal(err)
		}

		res, err = http.Post(
			addr+"/results",
			"application/json; charset=UTF-8",
			bytes.NewReader(requestBody),
		)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(4 * time.Second)

		// Finish new result after a little time.
		newFinishRequest := result_controller.FinishRequest{
			SportsmenID:  sportsmenID.ID,
			CheckpointID: checkPointID.ID,
			Time:         makeTimestamp(),
		}

		requestBody, err = json.Marshal(newFinishRequest)
		if err != nil {
			log.Fatal(err)
		}

		res, err = http.Post(
			addr+"/finish",
			"application/json; charset=UTF-8",
			bytes.NewBuffer(requestBody),
		)
		if err != nil {
			log.Fatal(err)
		}
	}()
	}
}

func makeTimestamp() int64 {
	return (time.Now().UnixNano() / int64(time.Millisecond))
}
