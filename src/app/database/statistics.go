package database

import (
	"time"
)

type TotalStatisticsRecord struct {
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
}

func (d *Database) TotalsForDay(day time.Time) <-chan TotalStatisticsRecord {

	rows, err := d.connection.Query(`
		SELECT title, sum(duration)
			FROM event_log
			WHERE date(at) == date(?)
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

			//err := rows.Scan(&title, &duration, &date)
			err := rows.Scan(&record.Title, &record.Duration)

			if err != nil {
				panic(err)
			}

			statistics <- record
		}
		close(statistics)
	}()

	return statistics
}
