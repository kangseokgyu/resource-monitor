package main

import (
	"fmt"

	"github.com/kangseokgyu/resource-monitor/memory"
)

func main() {
	out, err := memory.Get("vm_stat")
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	swap, err1 := memory.Get("sysctl", "vm.swapusage")
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(swap)
}
