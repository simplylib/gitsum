package main

import (
	"log"

	"github.com/ctii/gitsum/cmd"
)

func main() {
	if err := cmd.Main(); err != nil {
		log.Println(err)
	}
}
