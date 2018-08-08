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
type cfg struct {
	GinMode          string // Mode for Gin middleware
	DisableTimeStamp bool   // Whether to log timestamps
}

var appConfig cfg

var Web *gin.Engine

// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
//
// Init:
//
// Initialise the common elements of both a master/watcher.
//
func Init(cFile string) error {
	// Read config
	err := config.Read(cFile, &appConfig)
	if err != nil {
		return errors.Wrap(err, "App- Read config failed")
	}

	// Start up log code
	log.InitLog("EXP1", "Folder Watch", appConfig.DisableTimeStamp)

	log.Standard.Printf("[%s mode][mport=%v][wport=%v]", config.Mode, config.MPort, config.WPort)

	// Fire up Gin, for serving http
	gin.SetMode(appConfig.GinMode)
	Web = gin.New()

	// Middleware - handles panics and restarts
	Web.Use(gin.Recovery())

	return nil
}

// ------------------------------------------------------------------------------------------------
// RunWatcher:
//
// Main run loop when in 'watcher' mode.
// Will usually block until error/restart/panic occurs.
//
func RunWatcher() error {

	// Expose an endpoint that exposes the results
	// Not actually required, as
	Web.GET("/ping", watcher.Ping)

	err := watcher.StartWatcher(config.WatchPath, time.Duration(config.RefreshRate)*time.Millisecond)
	if err != nil {
		return errors.Wrap(err, "App- Failed to start watcher")
	}

	return Web.Run(":" + strconv.Itoa(config.WPort))
}

// ------------------------------------------------------------------------------------------------
// CloseWatcher:
//
// Close down when in 'watcher' mode.
//
func CloseWatcher() error {
	return watcher.StopWatcher()
}

// ------------------------------------------------------------------------------------------------
// Runmaster:
//
// Main run loop when in 'master' mode.
// Will usually block until error/restart/panic occurs.
//
func RunMaster() error {
	Web.POST("/update", master.Update)
	Web.GET("/", master.Results)

	return Web.Run(":" + strconv.Itoa(config.MPort))
}

// ------------------------------------------------------------------------------------------------
// CloseMaster:
//
// Close down when in 'master' mode.
//
func CloseMaster() error {

	return nil
}

// ------------------------------------------------------------------------------------------------
//
// Run:
//
// Run, and the close down, the code relevant to the mode we are running.
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
