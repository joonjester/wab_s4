package main

import (
	"code_first/server"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Code First application is getting started")

	err := server.InitializeAcc(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	server.Router()
}
