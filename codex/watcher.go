package codex

import (
	"github.com/radovskyb/watcher"
	"log"
	"time"
)

func watchChanges(loc string, onModified func()) {
	w := watcher.New()
	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case <-w.Event:
				onModified()
			case err := <-w.Error:
				log.Println("Error:", err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Add(loc); err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()

}
