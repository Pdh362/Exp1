package watcher

import (
	"github.com/Pdh362/Exp1/log"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

var FolderWatcher fsnotify.Watcher

func StartWatcher(path string) error {

	FolderWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "Start Watcher: Failed to create NewWatcher")
	}

	// done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-FolderWatcher.Events:
				log.Standard.Println("event:", ev)
			case err := <-FolderWatcher.Errors:
				log.Standard.Println("error:", err)
			}
		}
	}()

	err = FolderWatcher.Add(path)
	if err != nil {
		log.Standard.Fatal(err)
	}
	// <-done

	return nil
}
