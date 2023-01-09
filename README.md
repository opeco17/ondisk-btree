# Ondisk BTree implementation for Go

This package provides simple ondisk BTree implementation for Go.

Since code base is relatively small, it's suitable to learn about datastructure, storage engine and database.

The implementation is mainly inspired by [Database Internals (2019 Alex Petrov)](https://www.oreilly.com/library/view/database-internals/9781492040330/).

## Usage

```go
import ondisk-btree

btree := New()
```