package app

import (
	"github.com/Pdh362/Exp1/config"
	"github.com/Pdh362/Exp1/log"
	"github.com/Pdh362/Exp1/master"
	"github.com/Pdh362/Exp1/watcher"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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
	log.InitLog("EXP1", "Folder Watch")

	// Read config
	err := config.Read(cFile, &appConfig)
	if err != nil {
		return errors.Wrap(err, "App- Read config failed")
	}

	// Fire up Gin, for serving http
	gin.SetMode(appConfig.GinMode)
	Web = gin.New()

	// Middleware - handles panics and restarts
	Web.Use(gin.Recovery())

	return nil
}

// ------------------------------------------------------------------------------------------------
func RunWatcher() error {

	// Expose an endpoint that exposes the results
	// Not actually required, as
	Web.GET("/ping", watcher.Ping)

	err := watcher.StartWatcher(config.WatchPath, 500*time.Millisecond)
	if err != nil {
		return errors.Wrap(err, "App- Failed to start watcher")
	}

	return Web.Run(":" + strconv.Itoa(config.WPort))
}

// ------------------------------------------------------------------------------------------------
func CloseWatcher() error {
	return watcher.StopWatcher()
}

// ------------------------------------------------------------------------------------------------
func RunMaster() error {
	Web.POST("/update", master.Update)
	Web.GET("/", master.Results)

	return Web.Run(":" + strconv.Itoa(config.MPort))
}

// ------------------------------------------------------------------------------------------------
func CloseMaster() error {

	return nil
}

// ------------------------------------------------------------------------------------------------
//
// Run:
//
// This function will fire up the web server, and hence block unless something goes wrong.
//
func Run() error {
	var err error

	switch strings.ToLower(config.Mode) {

	case "watcher":
		err = RunWatcher()
	case "master":
		err = RunMaster()
	}

	if err != nil {
		return errors.Wrap(err, "Error running app")
	}

	switch strings.ToLower(config.Mode) {

	case "watcher":
		err = CloseWatcher()
	case "master":
		err = CloseMaster()
	}

	return err
}
