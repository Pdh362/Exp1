package app

import (
	"github.com/Pdh362/Exp1/config"
	"github.com/Pdh362/Exp1/log"
	"github.com/Pdh362/Exp1/watcher"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

type cfg struct {
	GinMode      string // Mode for Gin middleware
	GinConnLimit int    // Master connection limit
}

var appConfig cfg

var Web *gin.Engine

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
//
// Init:
//
//
func Init(cFile string) error {
	// Start up log
	log.InitLog("EXP1", "Experiment")

	// Read config
	err := config.Read(cFile, &appConfig)
	if err != nil {
		return errors.Wrap(err, "App- Read config failed")
	}

	// Fire up Gin, for serving http
	gin.SetMode(appConfig.GinMode)
	Web = gin.New()

	// Middleware
	Web.Use(gin.Recovery())

	return nil
}

// ------------------------------------------------------------------------------------------------
//
// Run:
//
// This function will fire up the web server, and hence block unless something goes wrong.
//
func Run() error {

	err := watcher.StartWatcher("./", 500*time.Millisecond)
	if err != nil {
		return errors.Wrap(err, "App- Failed to start watcher")
	}

	Web.Run(":8000")

	return watcher.StopWatcher()
}
