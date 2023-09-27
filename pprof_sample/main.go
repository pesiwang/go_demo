package main

import (
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"
)

var datas []string

// 运行程序
// 打开链接 http://localhost:6060/debug/pprof/ 可以查看pprof 数据

func main() {
	go func() {
		for {
			log.Printf("len: %d", Add("go-programming-tour-book"))
			time.Sleep(time.Second * 3)
		}
	}()

	_ = http.ListenAndServe("0.0.0.0:6060", nil)
}

func Add(str string) int {
	data := []byte(str)
	datas = append(datas, string(data))
	return len(datas)
}
