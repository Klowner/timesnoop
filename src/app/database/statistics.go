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

func (d *Database) eventsQueryAsChannel(query string, params ...interface{}) <-chan EventRecord {
	rows, err := d.connection.Query(query, params...)

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

func EventsChannelToSlice(in <-chan EventRecord) []EventRecord {
	out := make([]EventRecord, 0)
	for record := range in {
		out = append(out, record)
	}
	return out
}

func (d *Database) EventsAllChannel() <-chan EventRecord {
	return d.eventsQueryAsChannel(`
		SELECT title, sum(duration)
		FROM event_log
		GROUP BY title
		ORDER BY sum(duration) DESC
	`)
}

func (d *Database) EventsAll() []EventRecord {
	return EventsChannelToSlice(d.EventsAllChannel())
}

func (d *Database) EventsTotalsForDayChannel(day time.Time) <-chan EventRecord {
	return d.eventsQueryAsChannel(`
		SELECT title, sum(duration),
		FROM event_log
		WHERE datetime(ts, 'start of day') == datetime(?)
		GROUP BY title
		ORDER BY sum(duration)
	`, day)
}

func (d *Database) EventsTotalsForDay(day time.Time) []EventRecord {
	return EventsChannelToSlice(d.EventsTotalsForDayChannel(day))
}
