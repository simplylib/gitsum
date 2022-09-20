package main

import (
	"log"

	"github.com/simplylib/gitsum/cmd"
)

func main() {
	if err := cmd.Main(); err != nil {
		log.Println(err)
	}
}
