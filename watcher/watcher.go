package watcher

import (
	"github.com/Pdh362/Exp1/log"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"io/ioutil"
)

var folderWatcher fsnotify.Watcher
var contents []string
var watchPath string

func Index(t string) int {
	for i, v := range contents {
		if v == t {
			return i
		}
	}
	return -1
}

func addFile(path string) error {
	// Check we don't already have this file (could happen!)
	for _, v := range contents {
		if v == path {
			logContents()
			return nil
		}
	}

	// Append to the end of the array.
	contents = append(contents, path)

	logContents()
	return nil
}

func removeFile(path string) error {
	i := Index(path)
	if i == -1 {
		return errors.New("StopWatcher: Failed to close watch object")
	}

	// This is idiomatic Go, but certainly looks inefficient to me. BUT, it avoids a re-sort.
	// See https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-array-in-golang/37335777

	// We are removing an element from the slice, in this case the file just deleted.
	contents = append(contents[:i], contents[i+1:]...)

	logContents()
	return nil
}

func EventOccurred(event fsnotify.Event) {
	log.Standard.Println("event:", event)

	// These events alter the directory contents - so, we need to update our directory listing
	switch event.Op {
	// Add a file
	case fsnotify.Create:
		addFile(event.Name)
	// Remove an entry
	case fsnotify.Remove:
		removeFile(event.Name)
	// A rename is always followed by a CREATE event, so we'll act as if the original file is removed.
	case fsnotify.Rename:
		removeFile(event.Name)
	}
}

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

func StartWatcher(p string) error {
	watchPath = p

	FolderWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "StartWatcher: Failed to create NewWatcher")
	}

	// Channel function to handle incoming file change events.
	go func() {
		for {
			select {
			case ev := <-FolderWatcher.Events:
				EventOccurred(ev)
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
