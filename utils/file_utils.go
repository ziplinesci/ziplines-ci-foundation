package foundation

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/rs/zerolog/log"
)

// WatchForFileChanges waits for a change to the provided file path and then executes the function
func WatchForFileChanges(filePath string, functionOnChange func(fsnotify.Event)) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal().Err(err).Msg("Creating file system watcher failed")
		}
		defer watcher.Close()

		file := filepath.Clean(filePath)
		fileDir, _ := filepath.Split(file)
		realFile, _ := filepath.EvalSymlinks(filePath)

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						eventsWG.Done()
						return
					}
					currentFile, _ := filepath.EvalSymlinks(filePath)
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if (filepath.Clean(event.Name) == file && event.Op&writeOrCreateMask != 0) || (currentFile != "" && currentFile != realFile) {
						realFile = currentFile
						functionOnChange(event)
					} else if filepath.Clean(event.Name) == file && event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok {
						log.Warn().Err(err).Msg("Watcher error")
					}
					eventsWG.Done()
					return
				}
			}
		}()
		watcher.Add(fileDir)
		initWG.Done()
		eventsWG.Wait()
	}()
	initWG.Wait()
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(directory string) bool {
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// PathExists checks if a directory exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
