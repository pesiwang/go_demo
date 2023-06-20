package main

import (
	"fmt"
	"runtime"
	"time"
)

func test(p *int) {
	fmt.Printf("arg p:%v\n", p)
}

func main() {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			if i == 5 {
				fmt.Println("exit")
				runtime.Goexit()
				// return
			}
		}
	}()

	i := 2
	p := &i
	pp := &p

	fmt.Printf("i:%v, p:%v, pp:%v\n", i, p, pp)
	test(&i)

	for {
		time.Sleep(time.Second)
		fmt.Println("main goroutine")
	}
}
