package main

import (
	"encoding/json"
	"fmt"
)

type TreeNode struct {
	Children []*TreeNode
	Color    string
	Value    *TagTotal
}

func (t *TreeNode) MarshalJSON() ([]byte, error) {
	if len(t.Children) > 0 {
		return json.Marshal(struct {
			*TagTotal
			Children []*TreeNode `json:"children"`
		}{t.Value, t.Children})
	} else {
		return json.Marshal(t.Value)
	}
}

func (t *TreeNode) ShiftDownDuration() {
	for i, _ := range t.Children {
		t.Children[i].ShiftDownDuration()
	}

	// If this node has a duration as well as children, then
	// the duration needs to be moved to a new child node.
	if len(t.Children) > 0 {
		t.Children = append(t.Children, &TreeNode{
			Children: nil,
			Value: &TagTotal{
				Name:     "Others",
				TagId:    -1,
				Duration: t.Value.Duration,
			},
		})
		t.Value.Duration = 0
		fmt.Printf("Add others to %s\n", t.Value.Name)
	}
}

type TreeAssignmentRecord struct {
	Id       int
	ParentId int
	Depth    int
}

// returned TreeAssignments should be in order of the root parent id.
// That means in order to find the parent, we only need to back-track
// up the results until a record wither a lesser depth is found.
func (d *Database) GetTreeAssignments() []TreeAssignmentRecord {
	rows, err := d.connection.Query(`
		WITH RECURSIVE rec(i, pid, id) AS (
			VALUES (-1, 0, 0)
			UNION ALL
			SELECT rec.i+1, tags.id, tags.parent_id
			FROM tags JOIN rec ON tags.parent_id=rec.pid
			ORDER BY 1 DESC
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

func BuildTagTotalsTree(tag_totals []TagTotal) *[]TreeNode {
	tags := GetDB().GetTags()

	nodes := make([]TreeNode, len(tag_totals))

	// tag_totals are in no particular order, but tags are
	// sorted according to their parent.
	// Simplest way to organize the totals seems to be building
	// an index from tagid to the target index
	tagid_to_index := make(map[int]int)

	root_node_count := 0
	for i, tag := range tags {
		tagid_to_index[tag.Id] = i

		if tag.Depth > 0 {
			parent_index := i - 1
			for parent_index >= 0 && tags[parent_index].Depth >= tag.Depth {
				parent_index--
			}

			nodes[parent_index].Children = append(nodes[parent_index].Children, &nodes[i])

		} else {
			root_node_count++
		}
	}

	// Link values up to nodes
	for i, total := range tag_totals {
		nodes[tagid_to_index[total.TagId]].Value = &tag_totals[i]
		nodes[tagid_to_index[total.TagId]].Color = tags[tagid_to_index[total.TagId]].Color
	}

	j := 0
	out := make([]TreeNode, root_node_count)
	for _, tag := range tags {
		if tag.Depth == 0 {
			out[j] = nodes[tagid_to_index[tag.Id]]
			j++
		}
	}

	return &out
}

func ShiftDownDurations(tree_nodes *[]TreeNode) *[]TreeNode {
	for i, _ := range *tree_nodes {
		(*tree_nodes)[i].ShiftDownDuration()
	}
	return tree_nodes
}
