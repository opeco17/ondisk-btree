package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Node[T Item] struct {
	offset       OffsetType
	elements     []*Element[T]
	childOffsets []OffsetType
}

func newNode[T Item](offset OffsetType) *Node[T] {
	node := new(Node[T])
	node.offset = offset
	return node
}

func calNodeSize[T Item](maxElements int) int {
	return LENGTH_IN_NODE_BYTE + LENGTH_IN_NODE_BYTE + maxElementLengthByte[T](maxElements) + maxChildOffsetLengthByte[T](maxElements)
}

func maxElementLengthByte[T Item](maxElements int) int {
	elementSize := calElementSize[T]()
	return elementSize * (maxElements - 1)
}

func maxChildOffsetLengthByte[T Item](maxElements int) int {
	return OFFSET_SIZE_BYTE * maxElements
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
		element := new(Element[T])
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
		if key < element.getKey() {
			return false, i
		}
	}
	return false, len(node.elements)
}

func (node *Node[T]) isLeaf() bool {
	return len(node.childOffsets) == 0
}

func (node *Node[T]) isOverPopulated(maxElements int) bool {
	return len(node.elements) > (maxElements - 1)
}

func (node *Node[T]) insertElement(element *Element[T], index int) {
	if len(node.elements) == index {
		node.elements = append(node.elements, element)
	} else {
		node.elements = append(node.elements[:index+1], node.elements[index:]...)
		node.elements[index] = element
	}
}

func (node *Node[T]) insertChildOffset(childOffset OffsetType, index int) {
	if len(node.elements) == index {
		node.childOffsets = append(node.childOffsets, childOffset)
	} else {
		node.childOffsets = append(node.childOffsets[:index+1], node.childOffsets[index:]...)
		node.childOffsets[index] = childOffset
	}
}

func (node *Node[T]) print(offset OffsetType, isRoot bool) {
	ItemKeys := []string{}
	childOffsets := []string{}
	for _, element := range node.elements {
		ItemKeys = append(ItemKeys, strconv.Itoa(int(element.getKey())))
	}
	for _, childOffset := range node.childOffsets {
		childOffsets = append(childOffsets, strconv.Itoa(int(childOffset)))
	}

	if isRoot {
		fmt.Printf("Offset: %s (root)\n", strconv.Itoa(int(offset)))
	} else {
		fmt.Printf("Offset: %s\n", strconv.Itoa(int(offset)))
	}
	fmt.Printf("| Item Keys: %s\n", strings.Join(ItemKeys, ","))
	fmt.Printf("| Child Offsets: %s\n", strings.Join(childOffsets, ","))
	fmt.Println("+--------------------")
}
