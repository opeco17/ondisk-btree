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

func New[T Item](path string) (*BTree[T], error) {
	if path == "" {
		return nil, errors.New("Parameter 'path' should not be empty")
	}
	if err := isValidItemFields[T](); err != nil {
		return nil, err
	}

	btree := new(BTree[T])
	btree.path = path
	btree.degree = TREE_DEGREE

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

		rootNode := new(Node[T])
		if err = btree.writeNodeToDisk(rootNode, OFFSET_SIZE_BYTE); err != nil {
			return nil, err
		}
	}

	btree.nodeSize = calNodeSize[T](btree.maxElements())

	return btree, nil
}

func (btree *BTree[T]) Get(key KeyType) (*T, error) {
	if !btree.isOpen {
		return nil, errors.New("Tree is already closed")
	}
	return nil, nil
}

func (btree *BTree[T]) Put(item *T) error {
	if !btree.isOpen {
		return errors.New("Tree is already closed")
	}
	element := newElement(item)
	element.item = item

	isFound, travarsedNodes, travarsedIndices, travarsedOffsets, err := btree.traverse(element.getKey())
	if err != nil {
		return err
	}
	if isFound {
		if err = btree.update(element, travarsedNodes, travarsedIndices, travarsedOffsets); err != nil {
			return err
		}
		return nil
	} else {
		if err = btree.insert(element, travarsedNodes, travarsedIndices, travarsedOffsets); err != nil {
			return err
		}
		return nil
	}
}

func (btree *BTree[T]) Delete(key KeyType) error {
	if !btree.isOpen {
		return errors.New("Tree is already closed")
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

func (btree *BTree[T]) update(element *Element[T], travarsedNodes []Node[T], travarsedIndices []int, travarsedOffsets []OffsetType) error {
	numberOfTravarse := len(travarsedNodes)
	node := travarsedNodes[numberOfTravarse-1]
	index := travarsedIndices[numberOfTravarse-1]
	offset := travarsedOffsets[numberOfTravarse-1]

	node.elements[index] = *element
	if err := btree.writeNodeToDisk(&node, offset); err != nil {
		return err
	}
	return nil
}

func (btree *BTree[T]) insert(element *Element[T], travarsedNodes []Node[T], travarsedIndices []int, travarsedOffsets []OffsetType) error {
	return nil
}

func (btree *BTree[T]) traverse(key KeyType) (bool, []Node[T], []int, []OffsetType, error) {
	travarsedNodes := []Node[T]{}
	travarsedIndices := []int{}
	travarsedOffsets := []OffsetType{}

	offset := btree.getRootOffset()
	node, err := btree.readNodeFromDisk(offset)
	if err != nil {
		return false, nil, nil, nil, err
	}

	isFound, index := node.traverse(key)
	for {
		travarsedNodes = append(travarsedNodes, *node)
		travarsedIndices = append(travarsedIndices, index)
		travarsedOffsets = append(travarsedOffsets, offset)
		if isFound {
			return true, travarsedNodes, travarsedIndices, travarsedOffsets, nil
		}
		if node.isLeaf() {
			return false, travarsedNodes, travarsedIndices, travarsedOffsets, nil
		}
		offset := node.childOffsets[index]
		node, err = btree.readNodeFromDisk(offset)
		isFound, index = node.traverse(key)
		if err != nil {
			return false, nil, nil, nil, err
		}
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

	node := new(Node[T])
	node.deserialize(buff, btree.maxElements())
	return node, nil
}

func (btree *BTree[T]) writeNodeToDisk(node *Node[T], offset OffsetType) error {
	buff := node.serialize(btree.maxElements())
	btree.fp.Seek(offset, 0)
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
	binary.BigEndian.PutUint64(buff, OFFSET_SIZE_BYTE)
	_, err := btree.fp.Write(buff)
	defer btree.fp.Sync()
	if err != nil {
		return err
	}
	return nil
}
