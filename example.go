package main

import (
	"fmt"
)

func main() {
	var op = "="
	switch op {
	case "+":
	case "-":
		fmt.Println("-")
	default:
		fmt.Println("eror")
	}
}
