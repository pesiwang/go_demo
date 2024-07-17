package main

import (
	"fmt"
	"time"
)

func main() {
	innerChan := make(chan interface{}, 100)
	middleChan := make(chan interface{})
	outChan := make(chan interface{}, 1024)
	go func() {
		for i := 0; i < 1000; i++ {
			innerChan <- i
		}
		close(innerChan)
	}()

	go func() {
		for v := range innerChan {
			time.Sleep(1 * time.Millisecond)
			select {
			case middleChan <- v:
			default:
				// fmt.Println("middleChan block")
				middleChan <- v
			}
		}
		close(middleChan)
	}()

	start := time.Now()
	for v := range middleChan {
		outChan <- v
	}

	fmt.Println("cost:", time.Since(start).Milliseconds(), "ms")
}
