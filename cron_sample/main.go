package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func print() {
	fmt.Println(time.Now().Unix())
}

func main() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.NewWithLocation(loc)

	err := crontab.AddFunc("*/2 * * * * *", print)
	crontab.Start()
	if err != nil {
		fmt.Printf("crontab.AddFunc err %v\n", err)
	}

	select {}
}
