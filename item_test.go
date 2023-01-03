package main

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

func (sample Sample) GetKey() int64 {
	return int64(sample.Int)
}

type InvalidSample struct {
	Int     int
	Strings []string
}

func (invalidSample InvalidSample) GetKey() int64 {
	return int64(invalidSample.Int)
}

func TestItem(t *testing.T) {
	t.Run("Test serializeItem and deserializeItem", func(t *testing.T) {
		str := "hello, world"
		sample := new(Sample)
		sample.String = str
		deserializedSample := deserializeItem[Sample](serializeItem(sample))
		if deserializedSample.String != str {
			t.Errorf("string field should be %s\n", str)
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
		sample := new(Sample)
		sample.String = "Not long"
		if err := isValidStringLength(sample); err != nil {
			t.Errorf("Error should not be raised")
		}

		invalidSample := new(Sample)
		invalidSample.String = "Very loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong"
		if err := isValidStringLength(invalidSample); err == nil {
			t.Errorf("Error should be raised")
		}
	})
}
