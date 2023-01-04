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

	btree.nodeSize = calNodeSize[T](btree.maxItems())

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
	btree.traverse((*item).GetKey())
	return nil
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

func (btree *BTree[T]) traverse(key KeyType) (bool, []Node[T], []int, error) {
	travarsedNodes := []Node[T]{}
	travarsedIndices := []int{}

	rootOffset := btree.getRootOffset()
	node, err := btree.readNodeFromDisk(rootOffset)
	if err != nil {
		return false, nil, nil, err
	}

	isFound, index := node.traverse(key)
	for {
		travarsedNodes = append(travarsedNodes, *node)
		travarsedIndices = append(travarsedIndices, index)
		if isFound {
			return true, travarsedNodes, travarsedIndices, nil
		}
		if node.isLeaf() {
			return false, travarsedNodes, travarsedIndices, nil
		}
		node, err = btree.readNodeFromDisk(node.childOffsets[index])
		if err != nil {
			return false, nil, nil, err
		}
	}
}

func (btree *BTree[T]) minItems() int {
	return btree.degree - 1
}

func (btree *BTree[T]) maxItems() int {
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
	nodeSize := calNodeSize[T](btree.maxItems())
	buff := make([]byte, nodeSize)

	btree.fp.Seek(offset, 0)
	btree.fp.Read(buff)

	node := new(Node[T])
	node.deserialize(buff, btree.maxItems())
	return node, nil
}

func (btree *BTree[T]) writeNodeToDisk(node *Node[T], offset OffsetType) error {
	buff := node.serialize(btree.maxItems())
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
