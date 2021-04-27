package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					if filepath.Ext(event.Name) == ".metadata" {
						log.Println(event)
						b, err := ioutil.ReadFile(event.Name)
						if err != nil {
							log.Println(err)
						}
						var m metadata
						err = json.Unmarshal(b, &m)
						if err != nil {
							log.Println(err)
						}
						base := strings.TrimSuffix(event.Name, filepath.Ext(event.Name))
						b, err = ioutil.ReadFile(base + ".content")
						if err != nil {
							log.Println(err)
						}
						var c content
						err = json.Unmarshal(b, &c)
						if err != nil {
							log.Println(err)
						}
						if m.Lastopenedpage <= len(c.Pages) {
							page := filepath.Join(base, c.Pages[m.Lastopenedpage]) + ".rm"
							log.Println(page)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/home/root/.local/share/remarkable/xochitl")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

type current struct {
	base       string
	content    *content
	metadata   *metadata
	currenPage string
}
