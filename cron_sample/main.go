package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func print() {
	fmt.Println(time.Now())
}

func TriggerNewUserAwakeDaily() {
	fmt.Println("TriggerNewUserAwakeDaily", time.Now().Unix())
}
func main() {
	cronMain := cron.New(cron.WithSeconds())

	_, err := cronMain.AddFunc("*/2 * * * * *", print)
	if err != nil {
		fmt.Printf("crontab.AddFunc err %v\n", err)
	}

	_, err = cronMain.AddFunc("0 16 20 * * ?", TriggerNewUserAwakeDaily) // 每天晚上 20 点开始触发新用户唤醒
	if err != nil {
		fmt.Printf("cronMain.AddFunc err %v\n", err)
	}

	cronMain.Start()
	fmt.Printf("started\n")

	select {}
}
