package main

import (
	"github.com/Pdh362/Exp1"
	"github.com/Pdh362/Exp1/log"
)

func main() {

	// Start up log
	log.InitLog("EXP1", "Experiment")

	app.Init()
	app.Run()
}
