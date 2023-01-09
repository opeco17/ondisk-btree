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
	String16   string `maxLength:"16"`
}

func (item Sample) GetKey() int64 {
	return int64(item.Int)
}

type InvalidSample struct {
	Int     int
	String  string `maxLength:"hello"`
	Strings []string
}

func (item InvalidSample) GetKey() int64 {
	return int64(item.Int)
}

func TestItem(t *testing.T) {
	t.Run("Test padSpaces", func(t *testing.T) {
		paddedString := padSpaces("hello", 10)
		if paddedString != "hello     " {
			t.Errorf("paddedString should be 'hello     '")
		}

		paddedString = padSpaces("hello", 5)
		if paddedString != "hello" {
			t.Errorf("paddedString should be 'hello     '")
		}
	})
	t.Run("Test getMaxLength", func(t *testing.T) {
		maxLength, err := getMaxLength("")
		if err != nil {
			t.Errorf("Error should not be raised")
		}
		if maxLength != DEFAULT_STRING_MAX_LENGTH {
			t.Errorf("maxLength should be default value")
		}

		maxLength, err = getMaxLength("100")
		if err != nil {
			t.Errorf("Error should not be raised")
		}
		if maxLength != 100 {
			t.Errorf("maxLength should be 100")
		}

		maxLength, err = getMaxLength("0")
		if err == nil {
			t.Errorf("Error should be raised")
		}

		maxLength, err = getMaxLength("-1")
		if err == nil {
			t.Errorf("Error should be raised")
		}

		maxLength, err = getMaxLength("hello")
		if err == nil {
			t.Errorf("Error should be raised")
		}
	})
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
	t.Run("Test isValidStringLabel", func(t *testing.T) {
		if err := isValidStringLabel[InvalidSample](); err == nil {
			t.Errorf("Error should be raised")
		}
	})
	t.Run("Test isValidStringLength", func(t *testing.T) {
		item := new(Sample)
		item.String = "valid string"
		if err := isValidStringLength(item); err != nil {
			t.Errorf("Error should not be raised")
		}

		item = new(Sample)
		item.String16 = "valid string"
		if err := isValidStringLength(item); err != nil {
			t.Errorf("Error should not be raised")
		}

		item = new(Sample)
		item.String16 = "invalid looooooooooooooooooooong string"
		if err := isValidStringLength(item); err == nil {
			t.Errorf("Error should be raised")
		}
	})
}
