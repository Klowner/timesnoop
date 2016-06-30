package xdotool

import (
	"fmt"
	"time"
)

type FocusEvent struct {
	Title    string
	Duration float64
	Start    time.Time
}

func StartTracking() (<-chan FocusEvent, chan struct{}) {
	lastWindowTitle := ""
	lastChangeAt := time.Now()

	events := make(chan FocusEvent)
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				windowTitle := GetWindowTitle()

				if windowTitle != lastWindowTitle {
					fmt.Printf("name: " + windowTitle + "\n")

					events <- (FocusEvent{
						lastWindowTitle,
						time.Since(lastChangeAt).Seconds(),
						lastChangeAt,
					})

					lastChangeAt = time.Now()
					lastWindowTitle = windowTitle
				}
			case <-quit:
				close(events)
				ticker.Stop()
				return
			}
		}
	}()

	return events, quit
}
