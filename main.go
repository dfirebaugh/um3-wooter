package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Response -- struct for PrintJob Json response
type Response struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	TotalTime int64  `json:"time_total"`
	State     string
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
	finishes := howLong(currJob.TotalTime)

	if string(lastJob) != currJob.UUID {
		if currJob.State != "printing" {
			return
		}
		writeFile([]byte(currJob.UUID))
		newMsg := fmt.Sprint("woot! new printjob:\n\t`",
			currJob.Name, "`\n ```time: ",
			secondsToHuman(currJob.TotalTime),
			"\nfinishes: ",
			finishes,
			"``` https://ultimaker.hackrva.org")

		// if the job will take longer than 3 hours post
		// to 3d printing channel
		if hours > 3 {
			postJob(newMsg, slackHook)
		}

		postJob(newMsg, slackHookAll)
	}
}
