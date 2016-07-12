package main

import (
	"encoding/json"
)

type tagTreeNode struct {
	Children []*tagTreeNode
	Value    *Tag
}

func (t *tagTreeNode) MarshalJSON() ([]byte, error) {
	if len(t.Children) > 0 {
		return json.Marshal(struct {
			*Tag
			Children []*tagTreeNode `json:"children"`
		}{t.Value, t.Children})
	} else {
		return json.Marshal(struct {
			*Tag
			Children []int64 `json:"children"`
		}{t.Value, make([]int64, 0)})
	}
}

func BuildTagTree() []tagTreeNode {
	tags := GetDB().GetTags(false)
	tagIdToIndex := make(map[int]int)

	nodes := make([]tagTreeNode, len(tags))
	rootNodeCount := 0

	for i, tag := range tags {
		tagIdToIndex[tag.Id] = i

		if tag.Depth > 0 {
			parentIndex := i - 1
			for parentIndex >= 0 && tags[parentIndex].Depth >= tag.Depth {
				parentIndex--
			}

			nodes[parentIndex].Children = append(nodes[parentIndex].Children, &nodes[i])
		} else {
			rootNodeCount++
		}
	}

	// Link tags into tree structure
	for i, tag := range tags {
		nodes[tagIdToIndex[tag.Id]].Value = &tags[i]
	}

	j := 0
	out := make([]tagTreeNode, rootNodeCount)
	for _, tag := range tags {
		if tag.Depth == 0 {
			out[j] = nodes[tagIdToIndex[tag.Id]]
			j++
		}
	}

	return out
}
