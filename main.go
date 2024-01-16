package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

type Human Person

func main() {
	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Alice", Age: 30}

	// person1 和 person2 的值相同，类型相同
	fmt.Println(reflect.DeepEqual(person1, person2)) // 输出：true

	// 将 person2 转换为具有相同内容的不同类型
	human := Human(person2)
	fmt.Println(reflect.DeepEqual(person1, human)) // 输出：false

}
