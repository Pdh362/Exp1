package watcher

import (
	"github.com/Pdh362/Exp1/log"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"io/ioutil"
	"time"
)

var folderWatcher fsnotify.Watcher
var contents []string
var watchPath string
var dirtyFlag bool

func logContents() {
	log.Standard.Printf("%v", contents)
}

func StopWatcher() error {
	err := folderWatcher.Close()
	if err != nil {
		return errors.Wrap(err, "StopWatcher: Failed to close watch object")
	}
	return nil
}

func BuildDirFiles(path string) error {
	c, err := ioutil.ReadDir(path)

	contents = make([]string, len(c))
	for i, v := range c {
		contents[i] = v.Name()
	}

	logContents()

	return err
}

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

	FolderWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "StartWatcher: Failed to create NewWatcher")
	}

	// Create a ticker that checks our dirty flag, and updates if needed.
	ticker := time.NewTicker(refreshRate)
	go func() {
		// I'm not too happy about using range like this: the select statement is more appropriate,
		// but appears to blocking once one value is received.
		for t := range ticker.C {
			t = t
			// log.Standard.Println("Tick: " , t)
			// Check whether we need an update
			if dirtyFlag == true {
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
				// fileSystemEvent(ev)
				switch ev.Op {
				// Add a file
				case fsnotify.Create:
					dirtyFlag = true
				case fsnotify.Remove:
					dirtyFlag = true
				}

			case err := <-FolderWatcher.Errors:
				log.Standard.Println("error:", err)
			}
		}
	}()

	// Construct a list of files in the directory
	err = BuildDirFiles(watchPath)
	if err != nil {
		return errors.Wrap(err, "StartWatcher: Failed to scan directory for files: "+watchPath)
	}

	err = FolderWatcher.Add(watchPath)
	if err != nil {
		return errors.Wrap(err, "StartWatcher: Failed to add directory to watcher: "+watchPath)
		log.Standard.Fatal(err)
	}

	return nil
}
