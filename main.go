package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Duration(58*time.Minute + 14*time.Second).Nanoseconds())
}
