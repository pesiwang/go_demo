package main

import (
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	_ "net/http/pprof"
)

var datas []string

// 运行程序
// 打开链接 http://localhost:6060/debug/pprof/ 可以查看pprof 数据

// 获取60秒的 profile 文件
// curl -o profile.out http://127.0.0.1:6060/debug/pprof/profile?seconds=60

// 分析 profile 文件
// go tool pprof -http=:8080 ./profile.out

// 获取 goroutine 信息
// curl -o goroutine.txt http://127.0.0.1:6060/debug/pprof/profile?debug=1

func main() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	go func() {
		for {
			log.Printf("len: %d", Add("go-programming-tour-book"))
			time.Sleep(time.Second * 3)
		}
	}()

	for i := 0; i < 3; i++ {
		go needLock(i)
	}
	_ = http.ListenAndServe("0.0.0.0:6060", nil)
}

var mutex sync.RWMutex

func needLock(i int) {
	log.Printf("lock %v", i)
	mutex.Lock()
	if i == 0 {
		log.Printf("unlock %v", i)
		mutex.Unlock()
	}
}

func Add(str string) int {
	data := []byte(str)
	datas = append(datas, string(data))
	return len(datas)
}
