package main

import (
	"log"
	"time"

	"github.com/haormj/util/routine"
)

type HelloJob struct{}

func (*HelloJob) Do() {
	time.Sleep(time.Second * 4)
	log.Println("done")
}

func main() {
	// the number of goroutines
	maxWorkers := 1
	// the total of cache jobs
	maxJobs := 2
	d := routine.NewDispatcher(maxWorkers, maxJobs)
	d.Run()

	for {
		helloJob := HelloJob{}
		if d.Add(&helloJob) {
			log.Println("add job successfully")
		} else {
			log.Println("add job failed")
		}
		time.Sleep(time.Second)
	}
}
