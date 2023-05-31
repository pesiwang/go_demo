package main

import (
	"fmt"
	"test_sample/mymath"
)

func main() {
	fmt.Println("2 + 3 = ", mymath.MySum(2, 3))
	fmt.Println("2 - 3 = ", mymath.MySubstract(2, 3))
	fmt.Println("20 + 3 = ", mymath.MySum(20, 3))
	fmt.Println("32 - 3 = ", mymath.MySubstract(32, 3))
}
