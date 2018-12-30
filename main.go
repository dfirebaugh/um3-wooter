package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Response -- struct for PrintJob Json response
type Response struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	TotalTime int    `json:"time_total"`
	State     string
}

func createFile() {
	// detect if file exists
	var _, err = os.Stat("lastJob.txt")

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create("lastJob.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	fmt.Println("==> done creating file", "lastJob.txt")
}

func readFile() (b []byte) {
	createFile()
	// read the whole file at once
	b, err := ioutil.ReadFile("lastJob.txt")
	if err != nil {
		panic(err)
	}
	return
}

func writeFile(b []byte) {
	// write the whole body at once
	err := ioutil.WriteFile("lastJob.txt", b, 0644)
	if err != nil {
		panic(err)
	}
}

func getJob(URI string) Response {

	req, err := http.NewRequest("GET", URI, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	response1 := Response{}
	jsonErr := json.Unmarshal(body, &response1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return response1
}

func postJob(msg string, hook string) {
	fmt.Println(msg)

	newMsg := fmt.Sprint(`{"text":'`, msg, `'}`)
	var jsonStr = []byte(newMsg)
	req, err := http.NewRequest("POST", hook, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	slackHook := os.Getenv("SLACK_HOOK")
	slackHookAll := os.Getenv("SLACK_HOOK_ALL")
	um3URI := os.Getenv("UM3_URI")

	lastJob := readFile()
	currJob := getJob(um3URI)
	hours := secondsToHours(currJob.TotalTime)

	if string(lastJob) != currJob.UUID {
		if currJob.State == "pre_print" {
			return
		}
		writeFile([]byte(currJob.UUID))
		newMsg := fmt.Sprint("woot! new printjob: `", currJob.Name, "` -- time: ", secondsToHuman(currJob.TotalTime), "-- https://ultimaker.hackrva.org")

		// if the job will take longer than 3 hours post
		// to 3d printing channel
		if hours > 2 {
			postJob(newMsg, slackHook)
		}

		postJob(newMsg, slackHookAll)
	}
}
