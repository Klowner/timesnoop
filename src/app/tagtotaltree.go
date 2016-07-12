package main

import (
	"encoding/json"
	"fmt"
)

type tagTotalTreeNode struct {
	Children []*tagTotalTreeNode
	Value    *TagTotal
}

func (t *tagTotalTreeNode) MarshalJSON() ([]byte, error) {
	if len(t.Children) > 0 {
		return json.Marshal(struct {
			*TagTotal
			Children []*tagTotalTreeNode `json:"children"`
		}{t.Value, t.Children})
	} else {
		return json.Marshal(t.Value)
	}
}

func (t *tagTotalTreeNode) ShiftDownDuration() {
	for i, _ := range t.Children {
		t.Children[i].ShiftDownDuration()
	}

	// If this node has a duration as well as children, then
	// the duration needs to be moved to a new child node.
	if len(t.Children) > 0 {
		t.Children = append(t.Children, &tagTotalTreeNode{
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

func (t *tagTotalTreeNode) TotalUpDurations() float64 {
	for i, _ := range t.Children {
		t.Value.Duration += t.Children[i].TotalUpDurations()
	}

	return t.Value.Duration
}

func BuildTagTotalsTree(tag_totals []TagTotal) []tagTotalTreeNode {
	tags := GetDB().GetTags(false)

	nodes := make([]tagTotalTreeNode, len(tag_totals))

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
	}

	j := 0
	out := make([]tagTotalTreeNode, root_node_count)
	for _, tag := range tags {
		if tag.Depth == 0 {
			out[j] = nodes[tagid_to_index[tag.Id]]
			j++
		}
	}

	return out
}

func ShiftDownDurations(tree_nodes []tagTotalTreeNode) []tagTotalTreeNode {
	for i, _ := range tree_nodes {
		tree_nodes[i].ShiftDownDuration()
	}
	return tree_nodes
}

func TotalUpDurations(tree_nodes []tagTotalTreeNode) []tagTotalTreeNode {
	sum := 0.0
	for i, _ := range tree_nodes {
		sum += tree_nodes[i].TotalUpDurations()
	}
	return tree_nodes
}
