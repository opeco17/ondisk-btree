package main

import (
	"os"
	"testing"
)

func TestBTree(t *testing.T) {
	t.Run("Insert from -50 to 50", func(t *testing.T) {
		os.Remove(DEFAULT_DATA_PATH)

		start := -50
		end := 50
		btree, err := New[Sample](DEFAULT_DATA_PATH, DEFAULT_DEGREE)
		if err != nil {
			t.Errorf("Error should not be raised")
		}
		for i := start; i <= end; i++ {
			item := new(Sample)
			item.Int = i
			if err = btree.Put(item); err != nil {
				t.Errorf("Error should not be raised")
			}
		}
		for i := start; i <= end; i++ {
			item, err := btree.Get(KeyType(i))
			if err != nil {
				t.Errorf("Error should not be raised")
			}
			if item.Int != i {
				t.Errorf("item.Int should be %d", i)
			}
		}
	})
	t.Run("Insert from 50 to -50", func(t *testing.T) {
		os.Remove(DEFAULT_DATA_PATH)

		start := 50
		end := -50
		btree, err := New[Sample](DEFAULT_DATA_PATH, DEFAULT_DEGREE)
		if err != nil {
			t.Errorf("Error should not be raised")
		}
		for i := start; i >= end; i-- {
			item := new(Sample)
			item.Int = i
			if err = btree.Put(item); err != nil {
				t.Errorf("Error should not be raised")
			}
		}
		for i := start; i >= end; i-- {
			item, err := btree.Get(KeyType(i))
			if err != nil {
				t.Errorf("Error should not be raised")
			}
			if item.Int != i {
				t.Errorf("item.Int should be %d", i)
			}
		}
	})
}
