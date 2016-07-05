package main

import (
	"database/sql"
	"strconv"
	"strings"
)

type Tag struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
	Color    string `json:"color"`
}

type MatchExpression struct {
	Id          int    `json:"id"`
	TagId       int    `json:"tag_id"`
	Description string `json:"description"`
	Expression  string `json:"expression"`
}

type TagTotal struct {
	TagId    int     `json:"id"`
	Name     string  `json:"name"`
	Duration float64 `json:"duration"`
}

func (d *Database) GetTags() []Tag {
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

func (d *Database) CreateMatchExpression(expr *MatchExpression) *MatchExpression {
	d.WriteLock()
	_, err := d.connection.Exec("INSERT INTO matchers (description, expression, tag_id) VALUES (?, ?, ?)",
		expr.Description,
		expr.Expression,
		expr.TagId,
	)
	d.WriteUnlock()

	if err != nil {
		panic(err)
	}

	return expr
}

func processTagIds(ids []byte) []int64 {
	string_ids := strings.Split(string(ids), ",")
	out := make([]int64, 0)

	for _, i := range string_ids {
		j, err := strconv.Atoi(i)
		if err != nil {
			return out
		}
		out = append(out, int64(j))
	}
	return out
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

func (d *Database) GetMatchExpressions() []MatchExpression {
	rows, err := d.connection.Query(`
		SELECT id, tag_id, description, expression FROM matchers
		ORDER BY tag_id
		`)
	if err != nil {
		panic(err)
	}

	results := make([]MatchExpression, 0)
	for rows.Next() {
		var record MatchExpression

		err := rows.Scan(
			&record.Id,
			&record.TagId,
			&record.Description,
			&record.Expression,
		)

		if err != nil {
			panic(err)
		}

		results = append(results, record)
	}

	return results
}

func (d *Database) DeleteMatchExpressionById(id int64) {
	d.WriteLock()
	_, err := d.connection.Exec("DELETE FROM matchers WHERE id = ?", id)
	d.WriteUnlock()

	if err != nil {
		panic(err)
	}
}

func (d *Database) DeleteTagById(id int64) {
	d.WriteLock()
	_, err := d.connection.Exec("DELETE FROM tags WHERE id = ?", id)
	d.WriteUnlock()

	if err != nil {
		panic(err)
	}
}
