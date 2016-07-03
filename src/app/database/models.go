package database

type Tag struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
	Color    string `json:"color"`
}

type MatchExpression struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Expression  string `json:"expression"`
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
	_, err := d.connection.Exec("INSERT INTO match_expressions (description, expression) VALUES (?, ?, ?)",
		expr.Description,
		expr.Expression,
	)

	if err != nil {
		panic(err)
	}

	return expr
}

func (d *Database) GetMatchExpressions() []MatchExpression {
	rows, err := d.connection.Query("SELECT id, tag_id, description, expression FROM match_expressions")
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
