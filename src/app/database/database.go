package database

import (
	"app/xdotool"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	connection  *sql.DB
	stmt_insert *sql.Stmt
	stmt_merge  *sql.Stmt
	C           chan bool
}

func (d *Database) Open() bool {
	db, err := sql.Open("sqlite3", "./activitymonitor.db")
	checkErr(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS event_log(at timestamp, duration integer, title text)")
	checkErr(err)

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS event_log_title_idx ON event_log(title)")
	checkErr(err)

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS event_log_at_idx ON event_log(at)")
	checkErr(err)

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS event_log_at_title_idx ON event_log(at, title)")
	checkErr(err)

	stmt_insert, err := db.Prepare("INSERT INTO event_log(at, duration, title) values (?, ?, ?)")
	checkErr(err)

	stmt_merge, err := db.Prepare(`
		INSERT OR REPLACE INTO event_log (at, duration, title) VALUES (
			COALESCE((SELECT at FROM event_log WHERE title = ? AND at > date() ORDER BY at DESC LIMIT 1), ?),
			COALESCE((SELECT sum(duration) FROM event_log WHERE title = ?), 0) + ?,
			?
		)`)
	checkErr(err)

	d.connection = db
	d.stmt_insert = stmt_insert
	d.stmt_merge = stmt_merge

	return false
}

func (d *Database) RecordEvent(event xdotool.FocusEvent) {
	_, err := d.stmt_insert.Exec(event.Start, event.Duration, event.Title)
	checkErr(err)
}

func Initialize(events <-chan xdotool.FocusEvent) Database {
	database := Database{}
	database.Open()

	go func() {
		for event := range events {
			// Database will remain open as long as the events channel remains open
			database.RecordEvent(event)
		}
		defer database.connection.Close()
		close(database.C)
	}()

	return database
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
