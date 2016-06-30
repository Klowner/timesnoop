package database

import (
	"regexp"
)

type LabeledExpression struct {
	Label  string
	RegExp regexp.Regexp
}

type labeledExpressionRecord struct {
	label      string
	expression string
}

func (d *Database) AddExpression(label string, expression string) (int64, error) {

	// Attempt compilation to verify validity
	_, err := regexp.Compile(expression)
	if err != nil {
		return -1, err
	}

	res, err := d.connection.Exec(
		"INSERT OR REPLACE INTO label_regexps (label, expression) VALUES (?, ?)",
		label,
		expression)

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (d *Database) LoadExpressions() <-chan LabeledExpression {
	rows, err := d.connection.Query("SELECT label, expression FROM label_regexps ORDER BY expression")
	if err != nil {
		panic(err)
	}

	records := make(chan labeledExpressionRecord)

	go func() {
		for rows.Next() {
			var record labeledExpressionRecord
			err := rows.Scan(&record.label, &record.expression)

			if err != nil {
				panic(err)
			}

			records <- record
		}
		close(records)
	}()

	expressions := make(chan LabeledExpression)

	go func() {
		for record := range records {
			expr, err := regexp.Compile(record.expression)

			if err != nil {
				panic(err)
			}

			expressions <- LabeledExpression{
				record.label,
				*expr,
			}
		}
		close(expressions)
	}()

	return expressions
}
