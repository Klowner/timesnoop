package database

import (
	"time"
)

type TotalsRecord struct {
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
}

func (d *Database) TotalsForDay(day time.Time) []TotalsRecord {

	rows, err := d.connection.Query(`
		SELECT title, sum(duration)
			FROM event_log
			WHERE datetime(ts, 'start of day') == datetime(?)
			GROUP BY title
			ORDER BY sum(duration) DESC
		`, day)

	if err != nil {
		panic(err)
	}

	results := make([]TotalsRecord, 0)
	for rows.Next() {
		var record TotalsRecord
		err := rows.Scan(&record.Title, &record.Duration)
		if err != nil {
			panic(err)
		}
		results = append(results, record)
	}
	return results
}
