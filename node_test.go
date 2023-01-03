package main

import (
	"testing"
)

func TestNode(t *testing.T) {
	t.Run("Test serialize and deserialize", func(t *testing.T) {
		node := new(Node[Sample])
		maxItems := 10
		for i := 0; i < 3; i++ {
			sample := new(Sample)
			sample.Int = i
			node.items = append(node.items, sample)
		}
		for i := 0; i < 4; i++ {
			node.childOffsets = append(node.childOffsets, int64(i))
		}
		buff := node.serialize(maxItems)

		deserializedNode := new(Node[Sample])
		deserializedNode.deserialize(buff, maxItems)
		for i := 0; i < 3; i++ {
			if deserializedNode.items[i].Int != i {
				t.Errorf("deserializedNode.items[%d].Int should be %d", i, i)
			}
		}
		for i := 0; i < 4; i++ {
			if deserializedNode.childOffsets[i] != int64(i) {
				t.Errorf("deserializedNode.childOffsets[%d] should be %d", i, i)
			}
		}
	})
	t.Run("Test serialize and deserialize with max items and child offsets", func(t *testing.T) {
		node := new(Node[Sample])
		maxItems := 10
		for i := 0; i < maxItems-1; i++ {
			sample := new(Sample)
			sample.Int = i
			node.items = append(node.items, sample)
		}
		for i := 0; i < maxItems; i++ {
			node.childOffsets = append(node.childOffsets, int64(i))
		}
		buff := node.serialize(maxItems)

		deserializedNode := new(Node[Sample])
		deserializedNode.deserialize(buff, maxItems)
		for i := 0; i < maxItems-1; i++ {
			if deserializedNode.items[i].Int != i {
				t.Errorf("deserializedNode.items[%d].Int should be %d", i, i)
			}
		}
		for i := 0; i < maxItems; i++ {
			if deserializedNode.childOffsets[i] != int64(i) {
				t.Errorf("deserializedNode.childOffsets[%d] should be %d", i, i)
			}
		}
	})
}
