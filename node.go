package main

import (
	"encoding/binary"
)

type Node[T Item] struct {
	items        []*T
	childOffsets []OffsetType
}

func calNodeSize[T Item](maxItems int) int {
	return LENGTH_IN_NODE_BYTE + LENGTH_IN_NODE_BYTE + maxItemLengthByte[T](maxItems) + maxChildOffsetLengthByte[T](maxItems)
}

func maxItemLengthByte[T Item](maxItems int) int {
	itemSize := calItemSize[T]()
	return itemSize * (maxItems - 1)
}

func maxChildOffsetLengthByte[T Item](maxItems int) int {
	return OFFSET_SIZE_BYTE * maxItems
}

// Disk layout: {itemLength}{childOffsetLength}{item1}{item2}...{childOffset1}{childOffset2}...
func (node *Node[T]) serialize(maxItems int) []byte {
	buff := []byte{}

	lengthBuff := make([]byte, LENGTH_IN_NODE_BYTE*2)
	binary.BigEndian.PutUint64(lengthBuff[0:LENGTH_IN_NODE_BYTE], uint64(len(node.items)))
	binary.BigEndian.PutUint64(lengthBuff[LENGTH_IN_NODE_BYTE:LENGTH_IN_NODE_BYTE*2], uint64(len(node.childOffsets)))

	itemBuff := make([]byte, maxItemLengthByte[T](maxItems))
	itemCount := 0
	for _, item := range node.items {
		for _, b := range serializeItem(item) {
			itemBuff[itemCount] = b
			itemCount += 1
		}
	}

	childOffsetBuff := make([]byte, maxChildOffsetLengthByte[T](maxItems))
	for i, childOffset := range node.childOffsets {
		binary.BigEndian.PutUint64(childOffsetBuff[OFFSET_SIZE_BYTE*i:OFFSET_SIZE_BYTE*(i+1)], uint64(childOffset))
	}

	buff = append(buff, lengthBuff...)
	buff = append(buff, itemBuff...)
	buff = append(buff, childOffsetBuff...)
	return buff
}

func (node *Node[T]) deserialize(buff []byte, maxItems int) {
	itemSize := calItemSize[T]()

	lengthStartAt := 0
	itemsStartAt := LENGTH_IN_NODE_BYTE * 2
	childOffsetsStartAt := LENGTH_IN_NODE_BYTE*2 + maxItemLengthByte[T](maxItems)

	itemLength := binary.BigEndian.Uint64(buff[lengthStartAt : lengthStartAt+LENGTH_IN_NODE_BYTE])
	childOffsetLength := binary.BigEndian.Uint64(buff[lengthStartAt+LENGTH_IN_NODE_BYTE : lengthStartAt+LENGTH_IN_NODE_BYTE*2])

	for i := 0; i < int(itemLength); i++ {
		item := deserializeItem[T](buff[itemsStartAt+itemSize*i : itemsStartAt+itemSize*(i+1)])
		node.items = append(node.items, item)
	}

	for i := 0; i < int(childOffsetLength); i++ {
		childOffset := OffsetType(binary.BigEndian.Uint64(buff[childOffsetsStartAt+OFFSET_SIZE_BYTE*i : childOffsetsStartAt+OFFSET_SIZE_BYTE*(i+1)]))
		node.childOffsets = append(node.childOffsets, childOffset)
	}
}
