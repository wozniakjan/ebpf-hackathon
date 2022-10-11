package main

import (
	"fmt"
	"time"
)

func easyToFindFunctionName(arg uint32) {
	fmt.Println(arg)
}

func EasyToFindFunctionName(arg uint32) {
	easyToFindFunctionName(arg)
}

func main() {
	t1 := time.NewTicker(time.Second * 3)
	t2 := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t1.C:
			EasyToFindFunctionName(1)
		case <-t2.C:
			EasyToFindFunctionName(2)
		}
	}
}
