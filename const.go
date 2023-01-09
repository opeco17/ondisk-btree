package btree

import "reflect"

type OffsetType = int64
type KeyType = int64
type LengthInNodeType = int64

const DEFAULT_DATA_PATH = "btree.bin"
const OFFSET_SIZE_BYTE = 8
const LENGTH_IN_NODE_BYTE = 8
const DEFAULT_DEGREE = 3
const DEFAULT_STRING_MAX_LENGTH = 256

var AVAILABLE_TYPES = []reflect.Kind{
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Float32,
	reflect.Float64,
	reflect.Bool,
	reflect.String,
}
