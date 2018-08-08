package watcher

import (
	"bytes"
	"encoding/json"
	"github.com/Pdh362/Exp1/config"
	"github.com/Pdh362/Exp1/log"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//
// Global variables
//
var folderWatcher fsnotify.Watcher
var contents []string
var watchPath string
var dirtyFlag bool
var ticker *time.Ticker

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"response": "pong",
	})
	c.Next()
}

func GetContents() gin.H {
	return gin.H{
		"path":    watchPath,
		"results": contents}
}

func PostWatcherResults() error {
	jsonData := GetContents()
	jsonVal, _ := json.Marshal(jsonData)

	// log.Standard.Printf("json data is %s", jsonVal)

	// Inform the master that we have updated results
	addr := "http://localhost:" + strconv.Itoa(config.MPort) + "/update"
	_, err := http.Post(addr, "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		return errors.Wrap(err, "PostWatcherResults- Failed to post results")
	}

	// This would tell us if the server failed to respond
	// data, _ := ioutil.ReadAll(response.Body)
	// log.Standard.Printf("Response was :%s", data)

	return nil
}

//----------------------------------------------------------------------------------------------------------------------

func logContents() {
	log.Standard.Printf("%v", contents)
}

//----------------------------------------------------------------------------------------------------------------------
//
// BuildDirFiles:
//
// Rebuild the contents of the specified directory
//
func BuildDirFiles(path string) error {
	c, err := ioutil.ReadDir(path)

	contents = make([]string, len(c))
	for i, v := range c {
		contents[i] = v.Name()
	}

	// logContents()
	PostWatcherResults()

	return err
}

//----------------------------------------------------------------------------------------------------------------------
//
// StartWatcher:
//
// Start up a watcher. This works in 2 steps :-
//
// 1 -	A file system notification uses a channel to trigger events that 'something' has changed in the directory.
//		A 'dirty' flag is set to true, to indicate that we need to rebuild the directory listing.
//
// 2 -	A ticker channel receives updates every X microseconds, and checks the dirty flag.
// 		If true, set it back to false and rebuild our list of files.
//
// The benefit of decoupling the directory update from the file system notification, is that we can throttle how
// often the update occurs. An optimal algorithm would dynamically adjust this value, based on the workload of the
// service - but for this example, it is fixed.
//
func StartWatcher(p string, refreshRate time.Duration) error {
	watchPath = p
	dirtyFlag = true

	// Create a watcher, that monitors the specified folder.
	FolderWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "StartWatcher: Failed to create NewWatcher")
	}

	// Create a ticker that checks our dirty flag, and updates if needed.
	ticker = time.NewTicker(refreshRate)
	go func() {
		// I'm not too happy about using range like this: the select statement is more appropriate,
		// but appears to blocking once one value is received. Investigate?
		for t := range ticker.C {
			t = t
			// Check whether an update has been flagged
			if dirtyFlag == true {
				// It has: reset the flag, and rebuild the folder contents
				dirtyFlag = false
				err = BuildDirFiles(watchPath)
			}
		}
	}()

	// Channel function to handle incoming file change events.
	go func() {
		for {
			select {
			case ev := <-FolderWatcher.Events:
				// We'll flag an update only for these specific file system events.
				// Others don't affect the folder contents, and can be ignored.
				switch ev.Op {
				case fsnotify.Create:
					dirtyFlag = true
				case fsnotify.Remove:
					dirtyFlag = true
				}

			// Trap any errors, and log them
			case err := <-FolderWatcher.Errors:
				log.Standard.Println("error:", err)
			}
		}
	}()

	// Add the folder to watch, which will start the watching process.
	err = FolderWatcher.Add(watchPath)
	if err != nil {
		return errors.Wrap(err, "StartWatcher: Failed to add directory to watcher: "+watchPath)
	}

	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//
// StopWatcher:
//
// Cleanup and shutdown the watcher.
//
func StopWatcher() error {

	// Close down the folder watcher
	err := folderWatcher.Close()
	if err != nil {
		return errors.Wrap(err, "StopWatcher: Failed to close watch object")
	}

	// Close the ticker
	ticker.Stop()

	return nil
}
