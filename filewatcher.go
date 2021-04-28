package linestogo

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	XochitlCacheDir = `/home/root/.local/share/remarkable/xochitl`
)

// StartPolling triggers two go-routines:
//
// The first goroutine monitors the XochitlCacheDir.
// It feeds a channel with the directory containing the last modified notebook (based on the metadata)
//
// The seecond go-routine is  reading the notebook directory from the previous channel  and monitors it
// It feeeds the output channel with the name of the current file which is likely the current page.
//
// Caution: the remarkable is writing *a lot*, so a 50ms timer is set to prevent sending too many events on the chan
func StartPolling(ctx context.Context) (<-chan string, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(ctx)
	dirC, err := pollXochitlCacheDir(ctx)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	pageC, err := getCurrentPageName(ctx, dirC)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	return pageC, cancel, err
}

func pollXochitlCacheDir(ctx context.Context) (<-chan string, error) {
	outC := make(chan string)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Send initial page
	go func() {
		dir, err := findMostRecent(XochitlCacheDir)
		if err != nil {
			log.Println("could not find most recent dir:", err)
		}

		dir = filepath.Join(XochitlCacheDir, dir)
		outC <- dir
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					if filepath.Ext(event.Name) == ".metadata" {
						base := strings.TrimSuffix(event.Name, filepath.Ext(event.Name))
						outC <- base
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-ctx.Done():
				return
			}
		}
	}()
	return outC, watcher.Add(XochitlCacheDir)
}

func getCurrentPageName(ctx context.Context, dirC <-chan string) (<-chan string, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	outC := make(chan string)
	go func() {
		defer watcher.Close()
		currDir := ""
		var t time.Duration
		var last time.Time
		for {
			select {
			case f := <-dirC:
				if f != currDir {
					if currDir != "" {
						err := watcher.Remove(currDir)
						if err != nil {
							log.Println(err)
						}
					}
					err = watcher.Add(f)
					if err != nil {
						log.Println(err)
					}
					// send initial page
					page, err := findMostRecent(f)
					if err != nil {
						log.Println(err)
					}
					page = filepath.Join(f, page)
					outC <- strings.Trim(strings.Trim(page, filepath.Ext(page)), "-metadata") + ".rm"
					currDir = f
				}
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					t = time.Since(last)
					if strings.Contains(event.Name, "-metadata.json") && t > 50*time.Millisecond {
						last = time.Now()
						file := strings.Trim(event.Name, "-metadata.json") + ".rm"
						outC <- file
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-ctx.Done():
				return
			}
		}
	}()
	return outC, nil
}
func findMostRecent(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var modTime time.Time
	var names []string
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					names = names[:0]
				}
				names = append(names, strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name())))
			}
		}
	}
	if len(names) == 1 {
		return names[0], nil
	}
	return "", errors.New("expected only one result")
}
