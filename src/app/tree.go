package main

import (
	"fmt"
)

type TreeNode struct {
	Children *[]TreeNode
	Value    *TagTotal
}

type TreeAssignmentRecord struct {
	Id       int
	ParentId int
	Depth    int
}

// returned TreeAssignments should be in order of the root parent id.

func (d *Database) GetTreeAssignments() []TreeAssignmentRecord {
	rows, err := d.connection.Query(`
		WITH RECURSIVE rec(i, pid, id) AS (
			VALUES (-1, 0, 0)
			UNION ALL
			SELECT rec.i+1, tags.id, tags.parent_id
			FROM tags JOIN rec ON tags.parent_id=rec.pid
			ORDER BY 1
			LIMIT -1 OFFSET 1
		)
		SELECT * FROM rec`)
	if err != nil {
		panic(err)
	}

	results := make([]TreeAssignmentRecord, 0)

	defer rows.Close()
	for rows.Next() {
		var record TreeAssignmentRecord

		err := rows.Scan(
			&record.Depth,
			&record.Id,
			&record.ParentId,
		)

		if err != nil {
			panic(err)
		}

		results = append(results, record)
	}

	fmt.Printf("%s\n", results)

	return results
}

/*
func BuildTree(records []interface{}) []TreeNode {
	tree_assignments := GetDB().GetTreeAssignments()
	parentid_to_index := make(map[int]int)
	result := make([]TreeNode, 0)
	flat_nodes := make([]TreeNode, len(records))

	// build index for mapping parent_id to subscript id
	for i, assignment := range tree_assignments {
		parentid_to_index[assignment.ParentId] = i
		flat_nodes[i].Value
	}

	for _, record := range records {
		record.Id
	}

	fmt.Printf("%s\n", flat_nodes)
	return result
}
*/

func BuildTagTotalsTree(tag_totals []TagTotal) []TreeNode {

	// Build index for mapping parent_id to subscript id
	parentid_to_index := make(map[int]int)
	for i, assignment := range GetDB().GetTreeAssignments() {
		parentid_to_index[assignment.ParentId] = i
	}

	flat_nodes := make([]TreeNode, len(tag_totals))
	//flat_nodes := make([]TreeNode, len(tag_totals))

	//tags := GetDB().GetTags()

	//	tags := GetDB().GetTags()

	for i, _ := range tag_totals {
		flat_nodes[i].Value = &tag_totals[i]
	}

	//for i, _ := range tags {
	//flat_nodes[i].Value = &tags[i]
	//}

	//// stitch up children into the appropriate parents
	//for i, _ := range flat_nodes {
	//target := &flat_nodes[i]
	//_ = target
	//}

	return flat_nodes
}
