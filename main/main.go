package main

import (
	"github.com/Pdh362/Exp1"
	"github.com/Pdh362/Exp1/log"
)

func main() {
	err := app.Init("config.json")
	if err != nil {
		log.Standard.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Standard.Fatal(err)
	}
}
