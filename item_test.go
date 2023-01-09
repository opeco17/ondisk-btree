package btree

import (
	"testing"
)

type Sample struct {
	privateInt int
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	Bool       bool
	String     string
}

func (item Sample) GetKey() int64 {
	return int64(item.Int)
}

type InvalidSample struct {
	Int     int
	Strings []string
}

func (item InvalidSample) GetKey() int64 {
	return int64(item.Int)
}

func TestItem(t *testing.T) {
	t.Run("Test serializeItem and deserializeItem", func(t *testing.T) {
		str := "hello, world"
		item := new(Sample)
		item.String = str
		deserializedItem := deserializeItem[Sample](serializeItem(item))
		if deserializedItem.String != str {
			t.Errorf("string field should be %s", str)
		}
	})
	t.Run("Test isValidItemFields", func(t *testing.T) {
		if err := isValidItemFields[Sample](); err != nil {
			t.Errorf("Error should not be raised")
		}
		if err := isValidItemFields[InvalidSample](); err == nil {
			t.Errorf("Error should be raised")
		}
	})
	t.Run("Test isValidStringLength", func(t *testing.T) {
		item := new(Sample)
		item.String = "Not long"
		if err := isValidStringLength(item); err != nil {
			t.Errorf("Error should not be raised")
		}

		invalidItem := new(Sample)
		invalidItem.String = "Very loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong"
		if err := isValidStringLength(invalidItem); err == nil {
			t.Errorf("Error should be raised")
		}
	})
}
