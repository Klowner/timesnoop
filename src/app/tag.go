package main

import (
	"database/sql"
)

type Tag struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
	Color    string `json:"color"`
	Depth    int    `json:"depth"`
}

// Tags returned are sorted parent first
func (d *Database) GetTags() []Tag {
	rows, err := d.connection.Query(`
		WITH RECURSIVE rec(id, parent_id, name, color, r_depth) AS (
			VALUES (0, 0, null, null, -1)
			UNION ALL
			SELECT tags.id, tags.parent_id, tags.name, tags.color, rec.r_depth+1
			FROM tags JOIN rec ON tags.parent_id=rec.id
			ORDER BY 5 DESC
			LIMIT -1 OFFSET 1
		)
		SELECT id, parent_id, name, color, r_depth FROM rec`)

	if err != nil {
		panic(err)
	}

	results := make([]Tag, 0)

	defer rows.Close()
	for rows.Next() {
		var record Tag

		err := rows.Scan(
			&record.Id,
			&record.ParentId,
			&record.Name,
			&record.Color,
			&record.Depth,
		)

		if err != nil {
			panic(err)
		}

		results = append(results, record)
	}

	return results
}

func (d *Database) GetTagById(id int) Tag {
	rows, err := d.connection.Query("SELECT id, parent_id, name, color FROM tags WHERE id=? LIMIT 1", id)
	if err != nil {
		panic(err)
	}
	record := Tag{}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(
		&record.Id,
		&record.ParentId,
		&record.Name,
		&record.Color)

	if err != nil {
		panic(err)
	}
	return record
}

func (d *Database) GetTagsByName() []Tag {
	rows, err := d.connection.Query("SELECT id, parent_id, name, color FROM tags ORDER BY name")
	if err != nil {
		panic(err)
	}

	results := make([]Tag, 0)
	defer rows.Close()
	for rows.Next() {
		var record Tag

		err := rows.Scan(
			&record.Id,
			&record.ParentId,
			&record.Name,
			&record.Color)

		if err != nil {
			panic(err)
		}
		results = append(results, record)
	}

	results = append(results, Tag{
		Id:       -1,
		ParentId: 0,
		Name:     "Uncategorized",
		Color:    "#ffffff",
	})

	return results
}

func (d *Database) GetTagNames() map[int]string {
	rows, err := d.connection.Query("SELECT id, name FROM tags ORDER BY id")
	if err != nil {
		panic(err)
	}

	records := make(map[int]string)

	defer rows.Close()
	for rows.Next() {
		var name string
		var id int
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		records[id] = name
	}

	records[-1] = "Uncategorized"

	return records
}

func (d *Database) CreateTag(tag *Tag) *Tag {
	d.WriteLock()
	_, err := d.connection.Exec("INSERT INTO tags (parent_id, name, color) VALUES (?, ?, ?)",
		tag.ParentId,
		tag.Name,
		tag.Color,
	)
	d.WriteUnlock()

	if err != nil {
		panic(err)
	}

	return tag
}
func (d *Database) DeleteTagById(id int64) {
	d.WriteLock()
	_, err := d.connection.Exec("DELETE FROM tags WHERE id = ?", id)
	d.WriteUnlock()

	if err != nil {
		panic(err)
	}
}

func queryResultToTag(rows *sql.Rows) *Tag {
	record := Tag{}

	err := rows.Scan(
		&record.Id,
		&record.ParentId,
		&record.Name,
		&record.Color,
	)

	if err != nil {
		panic(err)
	}

	return &record
}
