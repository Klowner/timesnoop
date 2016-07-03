package main

import (
	"app/database"
	"regexp"
	"sync"
)

type CompiledMatchExpression struct {
	Expr   *regexp.Regexp
	Source *database.MatchExpression
}

var compiled_expressions_mutex = &sync.Mutex{}

var compiled_expressions *[]CompiledMatchExpression

func NewCompiledMatchExpression(me *database.MatchExpression) *CompiledMatchExpression {
	regexp, err := regexp.Compile(me.Expression)

	if err != nil {
		panic(err)
	}

	return &CompiledMatchExpression{
		regexp,
		me,
	}
}

func ReloadExpressions() int {
	count := 0
	db := database.GetDB()

	compiled_exprs := make([]CompiledMatchExpression, 0)

	for _, expression := range db.GetMatchExpressions() {
		compiled_expr := NewCompiledMatchExpression(&expression)
		compiled_exprs = append(compiled_exprs, *compiled_expr)
		count++
	}

	compiled_expressions = &compiled_exprs

	return count
}

func GetExpressions() *[]CompiledMatchExpression {
	if compiled_expressions == nil {
		return &[]CompiledMatchExpression{}
	}
	return compiled_expressions
}
