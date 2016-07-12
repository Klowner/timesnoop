package main

import (
	"fmt"
	"regexp"
	"sync"
)

type CompiledMatchExpression struct {
	Expr   *regexp.Regexp
	Source *MatchExpression
	Tag    *Tag
}

var compiled_expressions_mutex = &sync.Mutex{}

var compiled_expressions *[]CompiledMatchExpression

func NewCompiledMatchExpression(me MatchExpression) (*CompiledMatchExpression, error) {
	regexp, err := regexp.Compile(me.Expression)

	if err != nil {
		return nil, err
	}

	tag := GetDB().GetTagById(me.TagId)

	expr := CompiledMatchExpression{regexp, &me, &tag}
	return &expr, nil
}

func GetMatchersWithParentId(parent_tag_id int) *[]CompiledMatchExpression {
	out := []CompiledMatchExpression{}

	for _, matcher := range *GetMatchers() {
		if matcher.Tag.ParentId == parent_tag_id {
			out = append(out, matcher)
		}
	}

	return &out
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

func GetMatchers() *[]CompiledMatchExpression {
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
	expressions := GetMatchers()

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

func GetTotalsByTag(in <-chan EventRecord, matchers *[]CompiledMatchExpression, includeUnmatched bool) []TagTotal {
	out := []TagTotal{}
	tagTotals := make(map[int]float64)
	//tagNames := GetDB().GetTagNames()

	for event := range in {
		for _, matcher := range *matchers {
			match := false

			if matcher.Expr.MatchString(event.Title) {
				match = true

				tagTotals[matcher.Source.TagId] += event.Duration
			}

			if !match && includeUnmatched {
				tagTotals[-1] += event.Duration
			}
		}
	}

	for _, tag := range GetDB().GetTags(includeUnmatched) {
		total := TagTotal{
			TagId:    tag.Id,
			Name:     tag.Name,
			Duration: 0,
			Color:    tag.Color,
		}

		if val, ok := tagTotals[tag.Id]; ok {
			total.Duration = val
		}

		out = append(out, total)
	}

	//if incl_unmatched {
	//out = append(out, TagTotal{
	//TagId:    -1,
	//Name:     "Uncategorized",
	//Duration: tagTotals[-1],
	//Color:    "#444444",
	//})
	//}
	/*
		for id, name := range GetDB().GetTagNames() {

			fmt.Println(name)
			total := TagTotal{
				TagId:    id,
				Name:     name,
				Duration: 0,
			}

			if val, ok := tagTotals[id]; ok {
				total.Duration = val
			}

			out = append(out, total)
		}
	*/

	return out
}
