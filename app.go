package app

import (
	"github.com/Pdh362/Exp1/config"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

type cfg struct {
	LogFile      string // File name for logs: empty means log to STDOUT
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
	// Read config
	err := config.Read(cFile, &appConfig)
	if err != nil {
		return errors.Wrap(err, "App- Read config failed")
	}
	// Fire up Gin, for serving http
	gin.SetMode(appConfig.GinMode)
	Web = gin.New()

	return nil
}

// ------------------------------------------------------------------------------------------------
//
// Run:
//
// This function will fire up the web server, and hence block unless something goes wrong.
//
func Run() error {

	return nil
}
