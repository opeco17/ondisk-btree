package main

import (
	"os"
	"testing"
)

func TestBTree(t *testing.T) {
	os.Remove(DEFAULT_DATA_PATH)

	// btree, err := New[Sample](DEFAULT_DATA_PATH)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// btree.Put(&(Sample{String: "Hello, World"}))
}
