package main

import "fmt"

func easyToFindFunctionName() {
	fmt.Println("executed")
}

func EasyToFindFunctionName() {
	easyToFindFunctionName()
}

func main() {
	EasyToFindFunctionName()
}
