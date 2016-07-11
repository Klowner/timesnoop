package main

type MatchExpression struct {
	Id          int    `json:"id"`
	TagId       int    `json:"tag_id"`
	Description string `json:"description"`
	Expression  string `json:"expression"`
}

func (d *Database) DeleteMatchExpressionById(id int64) {
	d.WriteLock()
	_, err := d.connection.Exec("DELETE FROM matchers WHERE id = ?", id)
	d.WriteUnlock()

	if err != nil {
		panic(err)
	}
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
