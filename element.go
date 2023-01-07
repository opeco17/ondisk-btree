package main

type Element[T Item] struct {
	item     *T
	isClosed bool
}

func newElement[T Item](item *T) *Element[T] {
	element := new(Element[T])
	element.isClosed = false
	return element
}

func (element *Element[T]) getKey() KeyType {
	return (*element.item).GetKey()
}

func (element *Element[T]) serialize() []byte {
	buff := serializeItem(element.item)
	if element.isClosed {
		buff = append(buff, byte(1))
	} else {
		buff = append(buff, byte(0))
	}
	return buff
}

func (element *Element[T]) deserialize(buff []byte) {
	itemSize := calItemSize[T]()
	element.item = deserializeItem[T](buff[:itemSize])
	if int(buff[itemSize]) == 1 {
		element.isClosed = true
	} else {
		element.isClosed = false
	}
}

func calElementSize[T Item]() int {
	return calItemSize[T]() + 1
}
