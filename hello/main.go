package main

import (
	"fmt"
	"learngo/app"
)

func test(p *int) {
	fmt.Printf("arg p:%v\n", p)
}

func main() {
	// var a *app.userApp // error

	a := app.NewUserApp()

	a.Print()

	//a.pri() // error
}
