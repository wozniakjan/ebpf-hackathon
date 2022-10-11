package main

import "fmt"

func easyToFindFunctionName(arg uint32) {
	fmt.Println(arg)
}

func EasyToFindFunctionName(arg uint32) {
	easyToFindFunctionName(arg)
}

func main() {
	EasyToFindFunctionName(1)
	EasyToFindFunctionName(2)
}
