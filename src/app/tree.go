package main

import (
	"encoding/json"
	"fmt"
)

type TotalTreeNode struct {
	Children []*TotalTreeNode
	Color    string
	Value    *TagTotal
}

func (t *TotalTreeNode) MarshalJSON() ([]byte, error) {
	if len(t.Children) > 0 {
		return json.Marshal(struct {
			*TagTotal
			Children []*TotalTreeNode `json:"children"`
		}{t.Value, t.Children})
	} else {
		return json.Marshal(t.Value)
	}
}

func (t *TotalTreeNode) ShiftDownDuration() {
	for i, _ := range t.Children {
		t.Children[i].ShiftDownDuration()
	}

	// If this node has a duration as well as children, then
	// the duration needs to be moved to a new child node.
	if len(t.Children) > 0 {
		t.Children = append(t.Children, &TotalTreeNode{
			Children: nil,
			Value: &TagTotal{
				Name:     "Others",
				TagId:    -1,
				Duration: t.Value.Duration,
			},
		})
		fmt.Printf("Add others to %s %d\n", t.Value.Name, t.Value.Duration)
		t.Value.Duration = 0
	}
}

func BuildTagTotalsTree(tag_totals []TagTotal) []TotalTreeNode {
	tags := GetDB().GetTags()

	nodes := make([]TotalTreeNode, len(tag_totals))

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
	out := make([]TotalTreeNode, root_node_count)
	for _, tag := range tags {
		if tag.Depth == 0 {
			out[j] = nodes[tagid_to_index[tag.Id]]
			j++
		}
	}

	return out
}

func ShiftDownDurations(tree_nodes []TotalTreeNode) []TotalTreeNode {
	for i, _ := range tree_nodes {
		tree_nodes[i].ShiftDownDuration()
	}
	return tree_nodes
}
