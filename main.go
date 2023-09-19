package main

import "fmt"

func main() {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7}

	copy(s[2:(len(s)-1)], s[3:])
	fmt.Println(s)
}
