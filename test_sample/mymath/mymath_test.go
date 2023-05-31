package mymath

import "testing"

func TestMySum(t *testing.T) {
	val := MySum(2, 2)
	if val != 4 {
		t.Fatal("MySum case 1 test failed")
	}

	val2 := MySum(3, 2)
	if val2 != 5 {
		t.Fatal("MySum case 2 test failed")
	}
}

func TestMySubstract(t *testing.T) {
	val := MySubstract(2, 2)
	if val != 0 {
		t.Fatal("MySubstract case 1 test failed")
	}

	val2 := MySubstract(3, 2)
	if val2 != 1 {
		t.Fatal("MySubstract case 2 test failed")
	}
}

func BenchmarkMySum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MySum(i, i)
	}
}

func BenchmarkMyDivide(b *testing.B) {
	for i := 1; i < b.N; i++ {
		MyDivide(i, i)
	}
}
