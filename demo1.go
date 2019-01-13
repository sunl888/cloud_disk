package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.OpenFile("hole", os.O_CREATE|os.O_RDWR, 0600)
	defer file.Close()

	b := []byte{3}
	n, err := file.WriteAt(b, 0)

	fmt.Println(n, err)
	// output: 1 <nil>
}
