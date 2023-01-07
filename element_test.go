package main

import "testing"

func TestElement(t *testing.T) {
	t.Run("Test serialize and deserialize", func(t *testing.T) {
		str := "hello, world"

		element := new(Element[Sample])
		element.item = new(Sample)
		element.item.String = str
		element.isClosed = false

		deserializedElement := new(Element[Sample])
		deserializedElement.deserialize(element.serialize())

		if deserializedElement.item.String != str {
			t.Errorf("deserializedElement.item.String should be %s", str)
		}
		if deserializedElement.isClosed != false {
			t.Errorf("deserializedElement.isClosed should be false")
		}
	})
	t.Run("Test serialize and deserialize", func(t *testing.T) {
		str := "hello, world"

		element := new(Element[Sample])
		element.item = new(Sample)
		element.item.String = str
		element.isClosed = true

		deserializedElement := new(Element[Sample])
		deserializedElement.deserialize(element.serialize())

		if deserializedElement.item.String != str {
			t.Errorf("deserializedElement.item.String should be %s", str)
		}
		if deserializedElement.isClosed != true {
			t.Errorf("deserializedElement.isClosed should be true")
		}
	})
}
