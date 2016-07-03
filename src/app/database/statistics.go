package database

import (
	"time"
)

type EventRecord struct {
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
}

func (d *Database) TotalsForDay(day time.Time) []EventRecord {
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

	results := make([]EventRecord, 0)
	for rows.Next() {
		var record EventRecord
		err := rows.Scan(&record.Title, &record.Duration)
		if err != nil {
			panic(err)
		}
		results = append(results, record)
	}
	return results
}

func (d *Database) EventsAllChannel() <-chan EventRecord {
	rows, err := d.connection.Query(`
		SELECT title, sum(duration)
		FROM event_log
		GROUP BY title
		ORDER BY sum(duration) DESC`)

	if err != nil {
		panic(err)
	}

	out := make(chan EventRecord)

	go func() {
		for rows.Next() {
			rec := EventRecord{}
			err := rows.Scan(
				&rec.Title,
				&rec.Duration,
			)

			if err != nil {
				panic(err)
			}

			out <- rec
		}
		close(out)
	}()

	return out
}

func (d *Database) EventsAll() []EventRecord {
	out := make([]EventRecord, 0)
	for record := range d.EventsAllChannel() {
		out = append(out, record)
	}
	return out
}
