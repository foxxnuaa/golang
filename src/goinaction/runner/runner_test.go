package runner

import (
	"goinaction/runner"
	"log"
	"os"
	"testing"
	"time"
)

const (
	timeout = 6 * time.Second
)

func TestRunner(t *testing.T) {
	log.Println("Starting work...")

	r := runner.New(timeout)

	r.Add(createTask(), createTask(), createTask())

	if err := r.Start(); err != nil {
		switch err {
		case runner.ErrTimeout:
			log.Println("Terminating due to timeout")
			os.Exit(1)

		case runner.ErrInterrupt:
			log.Println("Terminating due to interrupt")
			os.Exit(2)
		}
	}

	log.Println("Process ended.")
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)

		log.Printf("Processor - Task #%d completed.", id)
	}
}
