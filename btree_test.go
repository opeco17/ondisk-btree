package main

import (
	"os"
	"testing"
)

func TestBTree(t *testing.T) {
	t.Run("Test put, delete and get (from -50 to 50)", func(t *testing.T) {
		os.Remove(DEFAULT_DATA_PATH)

		begin := -50
		end := 50
		deleteBegin := -10
		deleteEnd := 10

		btree, err := New[Sample](DEFAULT_DATA_PATH, DEFAULT_DEGREE)
		if err != nil {
			t.Errorf("Error should not be raised")
		}
		// Put
		for i := begin; i <= end; i++ {
			item := new(Sample)
			item.Int = i
			if err = btree.Put(item); err != nil {
				t.Errorf("Error should not be raised")
			}
		}

		// Get
		for i := begin; i <= end; i++ {
			item, err := btree.Get(KeyType(i))
			if err != nil {
				t.Errorf("Error should not be raised")
			}
			if item.Int != i {
				t.Errorf("item.Int should be %d", i)
			}
		}

		// Delete
		for i := deleteBegin; i <= deleteEnd; i++ {
			if err = btree.Delete(KeyType(i)); err != nil {
				t.Errorf("Error should not be raised")
			}
		}

		// Get
		for i := begin; i <= end; i++ {
			item, err := btree.Get(KeyType(i))
			if deleteBegin <= i && i <= deleteEnd {
				if err == nil {
					t.Errorf("Error should be raised")
				}
				if item != nil {
					t.Errorf("item should be null")
				}
			} else {
				if err != nil {
					t.Errorf("Error should not be raised")
				}
				if item.Int != i {
					t.Errorf("item.Int should be %d", i)
				}
			}
		}
	})
	t.Run("Test put, delete and get (from 50 to -50)", func(t *testing.T) {
		os.Remove(DEFAULT_DATA_PATH)

		begin := 50
		end := -50
		deleteBegin := 10
		deleteEnd := -10

		btree, err := New[Sample](DEFAULT_DATA_PATH, DEFAULT_DEGREE)
		if err != nil {
			t.Errorf("Error should not be raised")
		}
		// Put
		for i := begin; i >= end; i-- {
			item := new(Sample)
			item.Int = i
			if err = btree.Put(item); err != nil {
				t.Errorf("Error should not be raised")
			}
		}

		// Get
		for i := begin; i >= end; i-- {
			item, err := btree.Get(KeyType(i))
			if err != nil {
				t.Errorf("Error should not be raised")
			}
			if item.Int != i {
				t.Errorf("item.Int should be %d", i)
			}
		}

		// Delete
		for i := deleteBegin; i >= deleteEnd; i-- {
			if err = btree.Delete(KeyType(i)); err != nil {
				t.Errorf("Error should not be raised")
			}
		}

		// Get
		for i := begin; i >= end; i-- {
			item, err := btree.Get(KeyType(i))
			if deleteEnd <= i && i <= deleteBegin {
				if err == nil {
					t.Errorf("Error should be raised")
				}
				if item != nil {
					t.Errorf("item should be null")
				}
			} else {
				if err != nil {
					t.Errorf("Error should not be raised")
				}
				if item.Int != i {
					t.Errorf("item.Int should be %d", i)
				}
			}
		}
	})
}
