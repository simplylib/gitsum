package main

import (
	"log"

	"github.com/ctII/gitsum/cmd"
)

func main() {
	if err := cmd.Main(); err != nil {
		log.Println(err)
	}
}
