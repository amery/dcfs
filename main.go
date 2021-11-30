package main

import (
	"log"

	"github.com/amery/dcfs/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
