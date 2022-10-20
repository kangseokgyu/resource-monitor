package main

import (
	"fmt"

	"github.com/kangseokgyu/resource-monitor/memory"
)

func main() {
	out, err := memory.Get()
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
