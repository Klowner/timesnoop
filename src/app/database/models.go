package database

import (
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
	Id          int     `json:"id"`
	Description string  `json:"description"`
	Expression  string  `json:"expression"`
	TagIds      []int64 `json:"tag_ids"`
}

func (d *Database) GetTags() []Tag {
	rows, err := d.connection.Query("SELECT id, parent_id, name, color FROM tags ORDER BY name")
	if err != nil {
		panic(err)
	}

	results := make([]Tag, 0)
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

func (d *Database) CreateTag(tag *Tag) *Tag {
	_, err := d.connection.Exec("INSERT INTO tags (parent_id, name, color) VALUES (?, ?, ?)",
		tag.ParentId,
		tag.Name,
		tag.Color,
	)

	if err != nil {
		panic(err)
	}

	return tag
}

func (d *Database) CreateMatchExpression(expr *MatchExpression) *MatchExpression {
	_, err := d.connection.Exec("INSERT INTO match_expressions (description, expression) VALUES (?, ?)",
		expr.Description,
		expr.Expression,
	)

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

func (d *Database) GetMatchExpressions() []MatchExpression {
	rows, err := d.connection.Query(`
		WITH c AS (SELECT id, description, expression FROM match_expressions)
		SELECT id, description, expression, (
			SELECT group_concat(tag_id) FROM me2tags WHERE me_id = c.id
		) FROM c
		ORDER BY description
		`)
	if err != nil {
		panic(err)
	}

	results := make([]MatchExpression, 0)
	for rows.Next() {
		var record MatchExpression
		var tag_ids []byte

		err := rows.Scan(
			&record.Id,
			&record.Description,
			&record.Expression,
			&tag_ids,
		)

		record.TagIds = processTagIds(tag_ids)

		if err != nil {
			panic(err)
		}

		results = append(results, record)
	}
	return results
}

func (d *Database) DeleteMatchExpressionById(id int64) {
	_, err := d.connection.Exec("DELETE FROM match_expressions WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
}

func (d *Database) M2TCreate(mId int64, tId int64) {
	_, err := d.connection.Exec("INSERT OR REPLACE INTO me2tags (me_id, tag_id) VALUES (?, ?)", mId, tId)
	if err != nil {
		panic(err)
	}
}

func (d *Database) M2TDestroy(mId int64, tId int64) {
	_, err := d.connection.Exec("DELETE FROM me2tags WHERE me_id = ? AND tag_id = ?", mId, tId)
	if err != nil {
		panic(err)
	}
}
