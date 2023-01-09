package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type BTree[T Item] struct {
	path     string
	isOpen   bool
	degree   int
	nodeSize int
	fp       *os.File
}

func New[T Item](path string, degree int) (*BTree[T], error) {
	if path == "" {
		return nil, errors.New("Parameter 'path' should not be empty")
	}
	if degree <= 1 {
		return nil, errors.New("Parameter 'degree' should greater than 1")
	}
	if err := isValidItemFields[T](); err != nil {
		return nil, err
	}

	btree := new(BTree[T])
	btree.path = path
	btree.degree = degree

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open data file at %s", path))
	}
	btree.fp = fp
	btree.isOpen = true

	if btree.getLastOffset() == 0 {
		if err = btree.writeRootOffsetToDisk(OFFSET_SIZE_BYTE); err != nil {
			return nil, err
		}

		rootNode := newNode[T](OFFSET_SIZE_BYTE)
		if err = btree.writeNodeToDisk(rootNode); err != nil {
			return nil, err
		}
	}

	btree.nodeSize = calNodeSize[T](btree.maxElements())

	return btree, nil
}

func (btree *BTree[T]) Show() error {
	if err := btree.show(btree.getRootOffset(), true); err != nil {
		return err
	}
	return nil
}

func (btree *BTree[T]) Get(key KeyType) (*T, error) {
	if !btree.isOpen {
		return nil, errors.New("Tree is already closed")
	}

	isFound, traversedNodes, traversedIndices, err := btree.traverse(key)
	if err != nil {
		return nil, err
	}
	if !isFound {
		return nil, errors.New(fmt.Sprintf("Item with key %d is not found", key))
	}
	element := traversedNodes[len(traversedNodes)-1].elements[traversedIndices[len(traversedNodes)-1]]
	if element.isClosed {
		return nil, errors.New(fmt.Sprintf("Item with key %d is not found", key))
	}
	return element.item, nil
}

func (btree *BTree[T]) Put(item *T) error {
	if !btree.isOpen {
		return errors.New("Tree is already closed")
	}
	element := newElement(item)

	isFound, traversedNodes, traversedIndices, err := btree.traverse(element.getKey())
	if err != nil {
		return err
	}
	if isFound {
		if err = btree.update(element, traversedNodes, traversedIndices); err != nil {
			return err
		}
		return nil
	} else {
		if err = btree.insert(element, traversedNodes, traversedIndices); err != nil {
			return err
		}
		return nil
	}
}

func (btree *BTree[T]) Delete(key KeyType) error {
	if !btree.isOpen {
		return errors.New("Tree is already closed")
	}
	isFound, traversedNodes, traversedIndices, err := btree.traverse(key)
	if err != nil {
		return err
	}
	if !isFound {
		return errors.New(fmt.Sprintf("Item with key %d is not found", key))
	}

	node := traversedNodes[len(traversedNodes)-1]
	index := traversedIndices[len(traversedNodes)-1]

	node.elements[index].isClosed = true
	if err = btree.writeNodeToDisk(node); err != nil {
		return err
	}
	return nil
}

func (btree *BTree[T]) Close() error {
	if !btree.isOpen {
		return errors.New("Tree is already closed")
	}
	err := btree.fp.Close()
	if err != nil {
		return err
	}
	btree.isOpen = false
	return nil
}

func (btree *BTree[T]) show(offset OffsetType, isRoot bool) error {
	node, err := btree.readNodeFromDisk(offset)
	if err != nil {
		return err
	}
	node.print(offset, isRoot)
	for _, childOffset := range node.childOffsets {
		if err = btree.show(childOffset, false); err != nil {
			return err
		}
	}
	return nil
}

func (btree *BTree[T]) update(element *Element[T], traversedNodes []*Node[T], traversedIndices []int) error {
	numberOfTraverse := len(traversedNodes)
	node := traversedNodes[numberOfTraverse-1]
	index := traversedIndices[numberOfTraverse-1]

	node.elements[index] = element
	if err := btree.writeNodeToDisk(node); err != nil {
		return err
	}
	return nil
}

func (btree *BTree[T]) insert(element *Element[T], traversedNodes []*Node[T], traversedIndices []int) error {
	leafNode := traversedNodes[len(traversedNodes)-1]
	leafIndex := traversedIndices[len(traversedNodes)-1]
	leafNode.insertElement(element, leafIndex)

	if !leafNode.isOverPopulated(btree.maxElements()) {
		btree.writeNodeToDisk(leafNode)
		return nil
	}

	// Split non-root nodes
	for i := len(traversedNodes) - 1; i > 0; i-- {
		node := traversedNodes[i]
		parentNode := traversedNodes[i-1]
		parentIndex := traversedIndices[i-1]

		if node.isOverPopulated(btree.maxElements()) {
			newOffset := btree.getLastOffset()
			newNode := btree.split(node, parentNode, parentIndex, newOffset)
			btree.writeNodeToDisk(node)
			btree.writeNodeToDisk(newNode)
			if !parentNode.isOverPopulated(btree.maxElements()) {
				btree.writeNodeToDisk(parentNode)
			}
		} else {
			return nil
		}
	}

	// Split root node
	rootNode := traversedNodes[0]
	if rootNode.isOverPopulated(btree.maxElements()) {
		newRootOffset := btree.getLastOffset()
		newRootNode := newNode[T](newRootOffset)
		newRootNode.childOffsets = []OffsetType{rootNode.offset}
		btree.writeRootOffsetToDisk(newRootOffset)
		btree.writeNodeToDisk(newRootNode)

		newOffset := btree.getLastOffset()
		newNode := btree.split(rootNode, newRootNode, 0, newOffset)
		btree.writeNodeToDisk(rootNode)
		btree.writeNodeToDisk(newNode)
		btree.writeNodeToDisk(newRootNode)
	}
	return nil
}

func (btree *BTree[T]) split(node *Node[T], parentNode *Node[T], parentIndex int, newNodeOffset OffsetType) *Node[T] {
	middleElement := node.elements[btree.minElements()]
	newNode := newNode[T](newNodeOffset)

	newNode.elements = node.elements[btree.minElements()+1:]
	node.elements = node.elements[:btree.minElements()]
	if !node.isLeaf() {
		newNode.childOffsets = node.childOffsets[btree.minElements()+1:]
		node.childOffsets = node.childOffsets[:btree.minElements()+1]
	}
	parentNode.insertElement(middleElement, parentIndex)
	parentNode.insertChildOffset(newNodeOffset, parentIndex+1)

	return newNode
}

func (btree *BTree[T]) traverse(key KeyType) (bool, []*Node[T], []int, error) {
	traversedNodes := []*Node[T]{}
	traversedIndices := []int{}

	offset := btree.getRootOffset()
	node, err := btree.readNodeFromDisk(offset)
	if err != nil {
		return false, nil, nil, err
	}

	isFound, index := node.traverse(key)
	for {
		traversedNodes = append(traversedNodes, node)
		traversedIndices = append(traversedIndices, index)

		if isFound {
			return true, traversedNodes, traversedIndices, nil
		}
		if node.isLeaf() {
			return false, traversedNodes, traversedIndices, nil
		}
		offset = node.childOffsets[index]
		node, err = btree.readNodeFromDisk(offset)
		if err != nil {
			return false, nil, nil, err
		}
		isFound, index = node.traverse(key)
	}
}

func (btree *BTree[T]) minElements() int {
	return btree.degree - 1
}

func (btree *BTree[T]) maxElements() int {
	return btree.degree*2 - 1
}

func (btree *BTree[T]) getLastOffset() OffsetType {
	file, _ := os.Stat(btree.path)
	return file.Size()
}

func (btree *BTree[T]) getRootOffset() OffsetType {
	btree.fp.Seek(0, 0)
	buff := make([]byte, OFFSET_SIZE_BYTE)
	btree.fp.Read(buff)
	return OffsetType(binary.BigEndian.Uint64(buff))
}

func (btree *BTree[T]) readNodeFromDisk(offset OffsetType) (*Node[T], error) {
	btree.fp.Seek(offset, 0)
	nodeSize := calNodeSize[T](btree.maxElements())
	buff := make([]byte, nodeSize)

	btree.fp.Seek(offset, 0)
	btree.fp.Read(buff)

	node := newNode[T](offset)
	node.deserialize(buff, btree.maxElements())
	return node, nil
}

func (btree *BTree[T]) writeNodeToDisk(node *Node[T]) error {
	buff := node.serialize(btree.maxElements())
	btree.fp.Seek(node.offset, 0)
	_, err := btree.fp.Write(buff)
	defer btree.fp.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (btree *BTree[T]) writeRootOffsetToDisk(rootOffset OffsetType) error {
	btree.fp.Seek(0, 0)
	buff := make([]byte, OFFSET_SIZE_BYTE)
	binary.BigEndian.PutUint64(buff, uint64(rootOffset))
	_, err := btree.fp.Write(buff)
	defer btree.fp.Sync()
	if err != nil {
		return err
	}
	return nil
}
