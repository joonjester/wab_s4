package main

import (
	"code_first/controller"
	"fmt"
)

func main() {
	fmt.Println("Code First application is getting started")
	controller.SetUpServer()
}