package main

import (
	"fmt"
	"strconv"
)

type InterfaceAnimal interface {
	eat(food string)
}

type Animal struct {
	name    string
	subject string
}

func (a *Animal) eat(food string) {
	fmt.Println(a.name + ", " + food + ", " + a.subject)
	a.name = a.name + "_1"
}

type Cat struct {
	// Animal, *Animal 都可以实现继承，继承的函数能否修改成员变量并生效跟Animal, *Animal 无关，只跟对应的函数的接收是不是指针有关
	// 但是 Cat 作为参数传入到函数时，会根据函数的 receiver 要求传变量或指针
	// 如果是 *Animal 且 具体的函数要求指针，func(Cat) or func(&Cat) work
	// 如果是 Animal 且 函数要求指针，only func(&Cat) work
	*Animal
	age int
}

func (c *Cat) sleep() {
	fmt.Println(c.name + " " + strconv.Itoa(c.age))
}

// test_eat_animal(Cat) 报错，Cat 和 Animal 是两种不同类型
func test_eat_animal(a Animal) {
	a.eat("eat_in_test")
}

// test_eat_interface(Cat) test_eat_interface(Animal) work !
func test_eat_interface(a InterfaceAnimal) {
	a.eat("eat_in_test")
}

func main() {
	animal := Animal{name: "animal_name", subject: "animal_subject"}
	animal.eat("animal.eat")
	fmt.Println(animal)

	cat := Cat{Animal: &Animal{name: "cat_name", subject: "cat_subject"}, age: 999}
	cat.eat("cat.eat")
	cat.sleep()

	test_eat_animal(animal)
	// test_eat_animal(cat) // 报错, 原型 test_eat_animal(a Animal)

	test_eat_interface(&animal)
	test_eat_interface(&cat)

}
