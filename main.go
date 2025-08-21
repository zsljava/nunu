package main

import (
	"fmt"
	"github.com/zsljava/nunu/cmd/nunu"
)

func main() {
	err := nunu.Execute()
	if err != nil {
		fmt.Println("execute error: ", err.Error())
	}
}
