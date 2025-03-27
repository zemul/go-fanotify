package main

import (
	"fmt"
	gofanotify "github.com/zemul/go-fanotify"
	"log"
	"os"
	"os/signal"
)

func main() {
	notifier, err := gofanotify.New()
	if err != nil {
		log.Fatalf("init: %v", err)
	}
	defer notifier.Close()

	eventSet := gofanotify.NewEventSet(gofanotify.FileWriteComplete)
	err = notifier.AddWatch([]string{"/tmp"}, eventSet)
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	eventCh := notifier.ReadEvents()
	for {
		select {
		case event := <-eventCh:
			if event.Err != nil {
				log.Printf("error: %v", event.Err)
				continue
			}
			fmt.Printf("[%d] %s event: %v\n",
				event.PID,
				event.Path,
				eventToString(event))
		case <-sig:
			return
		}
	}
}

func eventToString(event gofanotify.Event) string {
	var s []string
	if event.Closed {
		s = append(s, "CLOSE_WRITE")
	}
	if event.Opened {
		s = append(s, "OPEN")
	}
	if event.Modified {
		s = append(s, "MODIFY")

	}
	return fmt.Sprintf("%v", s)
}
