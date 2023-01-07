package main

import (
	"encoding/binary"
)

type Node[T Item] struct {
	elements     []Element[T]
	childOffsets []OffsetType
}

func calNodeSize[T Item](maxElements int) int {
	return LENGTH_IN_NODE_BYTE + LENGTH_IN_NODE_BYTE + maxElementLengthByte[T](maxElements) + maxChildOffsetLengthByte[T](maxElements)
}

func maxElementLengthByte[T Item](maxElements int) int {
	elementSize := calElementSize[T]()
	return elementSize * (maxElements - 1)
}

func maxChildOffsetLengthByte[T Item](maxElement int) int {
	return OFFSET_SIZE_BYTE * maxElement
}

// Disk layout: {elementLength}{childOffsetLength}{element1}{element2}...{childOffset1}{childOffset2}...
func (node *Node[T]) serialize(maxElements int) []byte {
	buff := []byte{}

	lengthBuff := make([]byte, LENGTH_IN_NODE_BYTE*2)
	binary.BigEndian.PutUint64(lengthBuff[0:LENGTH_IN_NODE_BYTE], uint64(len(node.elements)))
	binary.BigEndian.PutUint64(lengthBuff[LENGTH_IN_NODE_BYTE:LENGTH_IN_NODE_BYTE*2], uint64(len(node.childOffsets)))

	elementBuff := make([]byte, maxElementLengthByte[T](maxElements))
	elementCount := 0
	for _, element := range node.elements {
		for _, b := range element.serialize() {
			elementBuff[elementCount] = b
			elementCount += 1
		}
	}

	childOffsetBuff := make([]byte, maxChildOffsetLengthByte[T](maxElements))
	for i, childOffset := range node.childOffsets {
		binary.BigEndian.PutUint64(childOffsetBuff[OFFSET_SIZE_BYTE*i:OFFSET_SIZE_BYTE*(i+1)], uint64(childOffset))
	}

	buff = append(buff, lengthBuff...)
	buff = append(buff, elementBuff...)
	buff = append(buff, childOffsetBuff...)
	return buff
}

func (node *Node[T]) deserialize(buff []byte, maxElements int) {
	elementSize := calElementSize[T]()

	lengthStartAt := 0
	elementsStartAt := LENGTH_IN_NODE_BYTE * 2
	childOffsetsStartAt := LENGTH_IN_NODE_BYTE*2 + maxElementLengthByte[T](maxElements)

	elementLength := binary.BigEndian.Uint64(buff[lengthStartAt : lengthStartAt+LENGTH_IN_NODE_BYTE])
	childOffsetLength := binary.BigEndian.Uint64(buff[lengthStartAt+LENGTH_IN_NODE_BYTE : lengthStartAt+LENGTH_IN_NODE_BYTE*2])

	for i := 0; i < int(elementLength); i++ {
		element := Element[T]{}
		element.deserialize(buff[elementsStartAt+elementSize*i : elementsStartAt+elementSize*(i+1)])
		node.elements = append(node.elements, element)
	}

	for i := 0; i < int(childOffsetLength); i++ {
		childOffset := OffsetType(binary.BigEndian.Uint64(buff[childOffsetsStartAt+OFFSET_SIZE_BYTE*i : childOffsetsStartAt+OFFSET_SIZE_BYTE*(i+1)]))
		node.childOffsets = append(node.childOffsets, childOffset)
	}
}

func (node *Node[T]) traverse(key KeyType) (bool, int) {
	for i, element := range node.elements {
		if key == element.getKey() {
			return true, i
		}
		if key > element.getKey() {
			return false, i
		}
	}
	return false, len(node.elements)
}

func (node *Node[T]) isLeaf() bool {
	return len(node.childOffsets) == 0
}
