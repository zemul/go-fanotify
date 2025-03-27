package main

import (
	"fmt"
	"github.com/zemul/go-fanotify/fanotify"
	"log"
	"os"
	"os/signal"
)

func main() {
	notifier, err := fanotify.New()
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	defer notifier.Close()
	// 2. 添加监控 (示例监控/tmp)
	err = notifier.AddWatch("/tmp",
		unix.FAN_OPEN|
			unix.FAN_MODIFY|
			unix.FAN_CLOSE_WRITE)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 处理信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// 4. 事件循环
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
				maskToString(event.Mask))

		case <-sig:
			return
		}
	}
}

func maskToString(mask uint64) string {
	var s []string
	if mask&unix.FAN_OPEN != 0 {
		s = append(s, "OPEN")
	}
	if mask&unix.FAN_MODIFY != 0 {
		s = append(s, "MODIFY")
	}
	if mask&unix.FAN_CLOSE_WRITE != 0 {
		s = append(s, "CLOSE_WRITE")
	}
	return fmt.Sprintf("%v", s)
}
