package main

import (
	"app/xdotool"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

type Database struct {
	connection           *sql.DB
	stmt_log_insert      *sql.Stmt
	stmt_daily_insert    *sql.Stmt
	stmt_daily_summarize *sql.Stmt
	C                    chan bool
	mutex                sync.Mutex
}

func (d *Database) WriteLock() *Database {
	d.mutex.Lock()
	return d
}

func (d *Database) WriteUnlock() *Database {
	d.mutex.Unlock()
	return d
}

var db *Database = nil

func GetDB() *Database {
	if db == nil {
		db = &Database{}
		db.Open()
	}
	return db
}

func (d *Database) Open() bool {
	fmt.Printf("opening database\n")
	db, err := sql.Open("sqlite3", "./timesnoop.dat")
	checkErr(err)

	_, err = db.Exec("PRAGMA foreign_keys = on")
	d.connection = db
	d.CreateSchema()
	d.PrepareStatements()

	return false
}

func (d *Database) RecordEvent(event xdotool.FocusEvent) {
	d.WriteLock()
	_, err := d.stmt_log_insert.Exec(event.Start, event.Duration, event.Title)
	d.WriteUnlock()

	checkErr(err)
}

func DatabaseInitialize(events <-chan xdotool.FocusEvent) *Database {
	database := GetDB()

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
