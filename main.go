package main

import (
	"log"

	"github.com/adycahyoputro/merchant/server"
)

func main() {
	// run the server
	if err := server.Run(); err != nil {
		log.Fatalf("error running the server: %s", err)
	}
}