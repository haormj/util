package main

import (
	"context"
	"log"
	"time"

	"github.com/haormj/util"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	util.DailyTaskCST(ctx, 11, 20, 0, func() {
		log.Println("hello timer")
	})
	select {
	case <-time.After(time.Second * 10):
		cancel()
	}
}
