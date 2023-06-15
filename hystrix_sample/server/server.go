package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

type Handle struct{}

func (h *Handle) ServeHTTP(r http.ResponseWriter, request *http.Request) {
	h.Common(r, request)
}

func (h *Handle) Common(r http.ResponseWriter, request *http.Request) {
	msg := "success"

	_ = hystrix.Do("mycommand", func() error {
		_, err := http.Get("https://www.baidu.com")
		if err != nil {
			fmt.Printf("请求失败:%v\n", err)
			return err
		}
		return nil
	}, func(err error) error {
		fmt.Printf("handle  error:%v\n", err)
		msg = "error"
		return nil
	})
	r.Write([]byte(msg))
}

func main() {
	hystrix.ConfigureCommand("mycommand", hystrix.CommandConfig{
		Timeout:                int(3 * time.Second),
		MaxConcurrentRequests:  10,
		SleepWindow:            5000,
		RequestVolumeThreshold: 20,
		ErrorPercentThreshold:  30,
	})
	http.ListenAndServe(":8090", &Handle{})
}
