# Ondisk BTree implementation for Go

This package provides simple ondisk BTree implementation for Go.

Since code base is relatively small, it's suitable to learn about datastructure, storage engine and database.

The implementation is mainly inspired by [Database Internals (2019 Alex Petrov)](https://www.oreilly.com/library/view/database-internals/9781492040330/).

Note that validation is minimum to keep simplicity of source code.

## Usage

```go
import btree "github.com/opeco17/ondisk-btree"

type Book struct {
	ID     int
	Name   string
	Author string `maxLength:"64"`
}

func (book Book) GetKey() int64 {
	return int64(book.ID)
}

func main() {
	btree, _ := btree.New[Book](btree.DEFAULT_DATA_PATH, btree.DEFAULT_DEGREE)
	defer btree.Close()

	btree.Put(&Book{ID: 0, Name: "Database Internals", Author: "Alex Petrov"})
	btree.Put(&Book{ID: 1, Name: "Designing Data-Intensive Applications", Author: "Martin Kleppmann"})

	btree.Delete(1)

	book, _ := btree.Get(0)
}
```
