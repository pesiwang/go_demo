package app

import "fmt"

type userApp struct {
}

func NewUserApp() *userApp {
	return &userApp{}
}

func (a userApp) Print() {
	fmt.Println("private struct, public function")
}

func (a userApp) pri() {
	fmt.Println("private struct, public function")
}
