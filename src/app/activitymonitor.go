package main

import (
	"app/database"
	"app/xdotool"
)

func main() {
	xdotool.Check()

	events, eventsQuit := xdotool.StartTracking()

	db := database.Initialize(events)
	routes(&db)

	switch {
	case <-db.C:
		close(eventsQuit)
	}
}
