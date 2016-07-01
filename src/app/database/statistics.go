package database

import (
	"fmt"
	"time"
)

type TotalStatisticsRecord struct {
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
}

func (d *Database) TotalsForDay(day time.Time) <-chan TotalStatisticsRecord {

	//WHERE date(at) == date(?)
	rows, err := d.connection.Query(`
		SELECT title, sum(duration)
			FROM event_log
			WHERE datetime(at, 'start of day') == datetime(?)
			GROUP BY title
			ORDER BY sum(duration) DESC
		`, day)

	if err != nil {
		panic(err)
	}

	statistics := make(chan TotalStatisticsRecord)

	go func() {
		for rows.Next() {
			var record TotalStatisticsRecord
			var date1 string
			var date2 string

			//err := rows.Scan(&title, &duration, &date)
			err := rows.Scan(&record.Title, &record.Duration)

			fmt.Printf(date1 + " " + date2 + "\n")
			if err != nil {
				panic(err)
			}

			statistics <- record
		}
		close(statistics)
	}()

	return statistics
}
