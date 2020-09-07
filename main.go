package main

import (
	"flag"
	"log"

	"github.com/traPtitech/traPortfolio/infrastructure"
)

func main() {
	migrate := flag.Bool("migrate", false, "migration mode or not")
	flag.Parse()
	if *migrate {
		_, err := infrastructure.NewSQLHandler()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("finished")
	} else {
		infrastructure.Init()
	}
}
