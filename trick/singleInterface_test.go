package trick

import (
	"fmt"
	"testing"
)

type DoSomething interface {
	DoSomething() string
}

type Concrete struct{}

func (a *Concrete) Print(msg string) string {
	return msg
}

func (a *Concrete) Hello() {
	fmt.Println("Hello")
}

type AdapterDoSomething func(msg string) string

func (s AdapterDoSomething) DoSomething() string {
	return s("123")
}

func DoXX(do DoSomething) {
	fmt.Println(do.DoSomething())
}

func Test_SingleInterface(t *testing.T) {
	a := &Concrete{}
	DoXX(AdapterDoSomething(a.Print))
}
