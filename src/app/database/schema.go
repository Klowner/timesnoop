package database

func (d *Database) CreateSchema() {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS event_log(
			at timestamp,
			duration integer,
			title text)`,

		"CREATE INDEX IF NOT EXISTS event_log_title_idx ON event_log(title)",
		"CREATE INDEX IF NOT EXISTS event_log_at_idx ON event_log(at)",
		"CREATE UNIQUE INDEX IF NOT EXISTS event_log_at_title_idx ON event_log(at, title)",

		`CREATE TABLE IF NOT EXISTS match_expression(
			tag_id integer,
			description text,
			expression text
		)`,

		"CREATE INDEX IF NOT EXISTS match_expression_name_idx ON match_expression(name)",
		"CREATE INDEX IF NOT EXISTS match_expression_tag_id_idx ON match_expression(tag_id)",

		`CREATE TABLE IF NOT EXISTS match_tag(
			name text PRIMARY KEY,
			parent_name text,
			color text)`,

		"CREATE INDEX IF NOT EXISTS match_tag_name ON match_tag(name)",
		"CREATE INDEX IF NOT EXISTS match_Tag_parent_name ON match_tag(parent_name)",
	}

	for _, statement := range statements {
		_, err := d.connection.Exec(statement)
		if err != nil {
			panic(err)
		}
	}
}
