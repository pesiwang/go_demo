package main

import (
	"fmt"
	"time"
)

func printPanic(err interface{}) {
	p := fmt.Sprintf("panic becuase: %v\n", err)
	fmt.Println("format panic msg:", p)
}

func set_data(x int) {
	defer func() {
		// recover() 可以将捕获到的panic信息打印
		if err := recover(); err != nil {
			printPanic(err)
			panicErrMsg := fmt.Sprintf("panic because: %v\n", err)
			stack := "this is set_data() panic stack"
			msg := panicErrMsg + stack
			fmt.Println(msg)
		}
	}()

	// 故意制造数组越界，触发 panic
	var arr [10]int
	arr[x] = 88
}

func main() {

	defer func() {
		// 无法捕获其他 goroutine 触发的 panic
		if err := recover(); err != nil {
			panicErrMsg := fmt.Sprintf("panic because: %v\n", err)
			stack := "this is main goroutine panic stack"
			msg := panicErrMsg + stack
			fmt.Println(msg)
		}
	}()

	// panic("this is panic msg")

	go func() {
		set_data(27)
	}()

	time.Sleep(time.Second * 3)

	data := map[string]string{}

	fmt.Println(data)
	fmt.Println("main goroutine quit")
}
