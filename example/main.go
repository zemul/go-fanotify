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

	err = notifier.AddWatch([]string{"/tmp"}, gofanotify.FAN_CLOSE_WRITE)
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

	if event.Mask&gofanotify.FAN_CLOSE_WRITE != 0 {
		s = append(s, "CLOSE_WRITE")
	}
	if event.Mask&gofanotify.FAN_OPEN != 0 {
		s = append(s, "OPEN")
	}
	if event.Mask&gofanotify.FAN_MODIFY != 0 {
		s = append(s, "MODIFY")

	}
	return fmt.Sprintf("%v", s)
}
