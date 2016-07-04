package main

func (d *Database) CreateSchema() {
	statements := []string{

		"PRAGMA foreign_keys = ON",

		// This table contains every title transition and will likely get
		// very large unless it is routinely condensed and cleaned.
		`CREATE TABLE IF NOT EXISTS event_log(
			ts timestamp,
			duration integer,
			title text)`,

		"CREATE INDEX IF NOT EXISTS event_log_title_idx ON event_log(title)",
		"CREATE INDEX IF NOT EXISTS event_log_at_idx ON event_log(ts)",

		// This table contains titles grouped by per-day resolution, with
		// unique titles and summed duration values.
		`CREATE TABLE IF NOT EXISTS event_daily(
			ts timestamp,
			duration integer,
			title text)`,

		"CREATE INDEX IF NOT EXISTS event_daily_title_idx ON event_log(title)",
		"CREATE INDEX IF NOT EXISTS event_daily_ts_idx ON event_log(ts)",
		"CREATE UNIQUE INDEX IF NOT EXISTS event_daily_tstitle_idx ON event_log(ts, title)",

		// Tags facilitate grouping matched patterns into a hierarchy for reporting
		`CREATE TABLE IF NOT EXISTS tags(
			id INTEGER PRIMARY KEY,
			parent_id integer,
			name text,
			color text)`,

		"CREATE INDEX IF NOT EXISTS tag_name_idx ON tags(name)",
		"CREATE INDEX IF NOT EXISTS tag_parent_id_idx ON tags(parent_id)",

		// Match expressions link patterns with specific tags, a tag can
		// be defined by any number of match expressions.
		`CREATE TABLE IF NOT EXISTS match_expressions(
			id INTEGER PRIMARY KEY,
			description text,
			expression text
		)`,

		"CREATE INDEX IF NOT EXISTS match_expression_idx ON match_expressions(description)",

		// Match expression <-> tag many-to-many table
		`CREATE TABLE IF NOT EXISTS me2tags(
			tag_id INTEGER,
			me_id INTEGER,
			FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE,
			FOREIGN KEY(me_id) REFERENCES match_expressions(id) ON DELETE CASCADE
		)`,

		"CREATE INDEX IF NOT EXISTS me2tags_tag_id_idx ON me2tags(tag_id)",
		"CREATE INDEX IF NOT EXISTS me2tags_me_id_idx ON me2tags(me_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS me2tags_unique_idx ON me2tags(tag_id, me_id)",
	}

	for _, statement := range statements {
		_, err := d.connection.Exec(statement)
		if err != nil {
			panic(err)
		}
	}
}

func (d *Database) PrepareStatements() {
	var err error

	d.stmt_log_insert, err = d.connection.Prepare(`
		INSERT INTO event_log (ts, duration, title) VALUES (?, ?, ?)`)
	checkErr(err)

	d.stmt_daily_insert, err = d.connection.Prepare(`
		INSERT OR REPLACE INTO event_log (at, duration, title) VALUES (
			COALESCE((SELECT at FROM event_log WHERE title = ? AND at > date() ORDER BY at DESC LIMIT 1), ?),
			COALESCE((SELECT sum(duration) FROM event_log WHERE title = ?), 0) + ?,
			?)`)

	d.stmt_daily_summarize, err = d.connection.Prepare(`
		`)
	checkErr(err)
}
