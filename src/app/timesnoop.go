package main

import (
	"app/xdotool"
)

func main() {
	xdotool.Check()

	events, eventsQuit := xdotool.StartTracking()

	db := DatabaseInitialize(events)
	routes(db)

	switch {
	case <-db.C:
		close(eventsQuit)
	}
}
