package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("@every 5s", func() { fmt.Println("Hello World") })
	c.Start()
	time.Sleep(2 * time.Minute)
}
