package database

type Tag struct {
	Name       string `json:"name"`
	ParentName string `json:"parent"`
	Color      string `json:"color"`
}

func (d *Database) GetTags() []Tag {
	rows, err := d.connection.Query("SELECT name, parent_name, color FROM tags")
	if err != nil {
		panic(err)
	}

	results := make([]Tag, 0)
	for rows.Next() {
		var record Tag
		err := rows.Scan(&record.Name, &record.ParentName, &record.Color)
		if err != nil {
			panic(err)
		}
		results = append(results, record)
	}
	return results
}
