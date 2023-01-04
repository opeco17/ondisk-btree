package main

import (
	"fmt"
	"os"
	"testing"
)

func TestBTree(t *testing.T) {
	os.Remove(DEFAULT_DATA_PATH)

	btree, err := New[Sample](DEFAULT_DATA_PATH)
	defer btree.Close()
	if err != nil {
		fmt.Println(err)
	}
	btree.Put(&Sample{})
}
