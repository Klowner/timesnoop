package main

import (
	"app/database"
	"fmt"
	"regexp"
	"sync"
)

type CompiledMatchExpression struct {
	Expr   *regexp.Regexp
	Source *database.MatchExpression
}

var compiled_expressions_mutex = &sync.Mutex{}

var compiled_expressions *[]CompiledMatchExpression

func NewCompiledMatchExpression(me *database.MatchExpression) (*CompiledMatchExpression, error) {
	regexp, err := regexp.Compile(me.Expression)

	if err != nil {
		return nil, err
	}

	return &CompiledMatchExpression{
		regexp,
		me,
	}, nil
}

func ReloadExpressions() int {
	count := 0
	db := database.GetDB()

	compiled_exprs := make([]CompiledMatchExpression, 0)

	for _, expression := range db.GetMatchExpressions() {
		compiled_expr, err := NewCompiledMatchExpression(&expression)
		if err != nil {
			panic(err)
		}
		compiled_exprs = append(compiled_exprs, *compiled_expr)
		count++
	}

	compiled_expressions = &compiled_exprs

	return count
}

func GetExpressions() *[]CompiledMatchExpression {
	if compiled_expressions == nil {
		ReloadExpressions()
	}
	fmt.Printf("%d expressions\n", len(*compiled_expressions))
	return compiled_expressions
}

func AppendExpression(expr *database.MatchExpression) {
	compiled_expr, err := NewCompiledMatchExpression(expr)
	if err != nil {
		panic(err)
	}
	*compiled_expressions = append(*compiled_expressions, *compiled_expr)
}

func EventRecordFilterUnmatched(in <-chan database.EventRecord) <-chan database.EventRecord {
	expressions := GetExpressions()

	// nothing to filter against, so all inputs pass
	if len(*expressions) == 0 {
		fmt.Printf("bypassing filter because no filters exist")
		return in
	}

	out := make(chan database.EventRecord)

	go func() {
		for event := range in {
			match := false
			// If any expressions match, then we continue
			for _, expression := range *expressions {
				if expression.Expr.MatchString(event.Title) {
					match = true
					break
				}
			}
			if !match {
				out <- event
			}
		}
		close(out)
	}()

	return out
}
