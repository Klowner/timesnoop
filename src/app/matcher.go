package main

import (
	"fmt"
	"regexp"
	"sync"
)

type CompiledMatchExpression struct {
	Expr   *regexp.Regexp
	Source *MatchExpression
}

type TagTotal struct {
	TagId    int     `json:"id"`
	Name     string  `json:"name"`
	Duration float64 `json:"duration"`
}

var compiled_expressions_mutex = &sync.Mutex{}

var compiled_expressions *[]CompiledMatchExpression

func NewCompiledMatchExpression(me MatchExpression) (*CompiledMatchExpression, error) {
	regexp, err := regexp.Compile(me.Expression)

	if err != nil {
		return nil, err
	}

	expr := CompiledMatchExpression{regexp, &me}
	return &expr, nil
}

func ReloadExpressions() int {
	count := 0
	db := GetDB()

	compiled_exprs := make([]CompiledMatchExpression, 0)

	for _, expression := range db.GetMatchExpressions() {
		compiled_expr, err := NewCompiledMatchExpression(expression)
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
	return compiled_expressions
}

func AppendExpression(expr *MatchExpression) {
	compiled_expr, err := NewCompiledMatchExpression(*expr)
	if err != nil {
		panic(err)
	}

	*compiled_expressions = append(*compiled_expressions, *compiled_expr)
}

func EventRecordFilterUnmatched(in <-chan EventRecord) <-chan EventRecord {
	expressions := GetExpressions()

	// nothing to filter against, so all inputs pass
	if len(*expressions) == 0 {
		fmt.Printf("bypassing filter because no filters exist")
		return in
	}

	out := make(chan EventRecord)

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

func GetTotalsByTag(in <-chan EventRecord) []TagTotal {
	db := GetDB()
	out := []TagTotal{}
	matchers := GetExpressions() // Get all expression matchers
	totals := make(map[string]float64)

	for event := range in {
		for _, matcher := range *matchers {
			if matcher.Expr.MatchString(event.Title) {
				// the event matches the current expression, then the
				// event's time must be added to the tally under the
				// matcher's assigned tags

				tags := db.GetTagsForMatcher(matcher.Source.Id)
				for _, tag := range tags {
					fmt.Printf("adding %f to %s\n", event.Duration, tag.Name)
					totals[tag.Name] += event.Duration
				}
			}
		}
	}

	for tagName, duration := range totals {
		out = append(out, TagTotal{
			0,
			tagName,
			duration,
		})
	}

	return out
}
