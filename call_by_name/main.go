package main

import (
	"fmt"
	"reflect"
)

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area(scale float64, unit string) {
	area := r.Width * r.Height * scale
	fmt.Printf("The area of the rectangle is %.2f %s^2.\n", area, unit)
}

func main() {
	r := Rectangle{Width: 3.0, Height: 4.0}

	methodName := "Area"
	method := reflect.ValueOf(r).MethodByName(methodName)
	arguments := []reflect.Value{reflect.ValueOf(2.0), reflect.ValueOf("m")}
	method.Call(arguments)
}
