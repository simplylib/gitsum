package main

import (
	"log"

	"github.com/ctii/gitsum/cmd"
)

func main() {
	err := cmd.Main()
	if err != nil {
		log.Println(err)
	}
}
