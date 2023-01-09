package btree

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Item interface {
	GetKey() KeyType
}

func serializeItem[T Item](item *T) []byte {
	buff := make([]byte, calItemSize[T]())
	buffPtr := 0
	itemVal := reflect.ValueOf(item).Elem()
	itemType := reflect.TypeOf(*item)
	for i := 0; i < itemVal.NumField(); i++ {
		field := itemVal.FieldByIndex([]int{i})
		fieldType := field.Type().Kind()
		fieldSize := field.Type().Size()
		if !field.CanSet() {
			continue
		} else if fieldType == reflect.Int8 || (fieldType == reflect.Int && fieldSize == 1) {
			buff[buffPtr] = byte(field.Int())
			buffPtr += 1
		} else if fieldType == reflect.Int16 || (fieldType == reflect.Int && fieldSize == 2) {
			binary.BigEndian.PutUint16(buff[buffPtr:buffPtr+2], uint16(field.Int()))
			buffPtr += 2
		} else if fieldType == reflect.Int32 || (fieldType == reflect.Int && fieldSize == 4) {
			binary.BigEndian.PutUint32(buff[buffPtr:buffPtr+4], uint32(field.Int()))
			buffPtr += 4
		} else if fieldType == reflect.Int64 || (fieldType == reflect.Int && fieldSize == 8) {
			binary.BigEndian.PutUint64(buff[buffPtr:buffPtr+8], uint64(field.Int()))
			buffPtr += 8
		} else if fieldType == reflect.Uint8 || (fieldType == reflect.Uint && fieldSize == 1) {
			buff[buffPtr] = byte(field.Uint())
			buffPtr += 1
		} else if fieldType == reflect.Uint16 || (fieldType == reflect.Uint && fieldSize == 2) {
			binary.BigEndian.PutUint16(buff[buffPtr:buffPtr+2], uint16(field.Uint()))
			buffPtr += 2
		} else if fieldType == reflect.Uint32 || (fieldType == reflect.Uint && fieldSize == 4) {
			binary.BigEndian.PutUint32(buff[buffPtr:buffPtr+4], uint32(field.Uint()))
			buffPtr += 4
		} else if fieldType == reflect.Uint64 || (fieldType == reflect.Uint && fieldSize == 8) {
			binary.BigEndian.PutUint64(buff[buffPtr:buffPtr+8], uint64(field.Uint()))
			buffPtr += 8
		} else if fieldType == reflect.Float32 {
			binary.BigEndian.PutUint32(buff[buffPtr:buffPtr+4], math.Float32bits(float32(field.Float())))
			buffPtr += 4
		} else if fieldType == reflect.Float64 {
			binary.BigEndian.PutUint64(buff[buffPtr:buffPtr+8], math.Float64bits(field.Float()))
			buffPtr += 8
		} else if fieldType == reflect.Bool {
			if field.Bool() {
				buff[buffPtr] = byte(1)
			} else {
				buff[buffPtr] = byte(0)
			}
			buffPtr += 1
		} else if fieldType == reflect.String {
			maxLength, _ := getMaxStringLength(itemType.Field(i).Tag.Get("maxLength"))
			for _, each := range []byte(padSpaces(field.String(), maxLength)) {
				buff[buffPtr] = each
				buffPtr += 1
			}
		}
	}
	return buff
}

func deserializeItem[T Item](buff []byte) *T {
	item := new(T)
	itemVal := reflect.ValueOf(item).Elem()
	itemType := reflect.TypeOf(*item)
	var buffPtr uint64 = 0
	for i := 0; i < itemVal.NumField(); i++ {
		field := itemVal.FieldByIndex([]int{i})
		fieldType := field.Type().Kind()
		fieldSize := field.Type().Size()

		if !field.CanSet() {
			continue
		} else if fieldType == reflect.Int8 || (fieldType == reflect.Int && fieldSize == 1) {
			field.SetInt(int64(buff[buffPtr]))
			buffPtr += 1
		} else if fieldType == reflect.Int16 || (fieldType == reflect.Int && fieldSize == 2) {
			field.SetInt(int64(binary.BigEndian.Uint16(buff[buffPtr : buffPtr+2])))
			buffPtr += 2
		} else if fieldType == reflect.Int32 || (fieldType == reflect.Int && fieldSize == 4) {
			field.SetInt(int64(binary.BigEndian.Uint32(buff[buffPtr : buffPtr+4])))
			buffPtr += 4
		} else if fieldType == reflect.Int64 || (fieldType == reflect.Int && fieldSize == 8) {
			field.SetInt(int64(binary.BigEndian.Uint64(buff[buffPtr : buffPtr+8])))
			buffPtr += 8
		} else if fieldType == reflect.Uint8 || (fieldType == reflect.Uint && fieldSize == 1) {
			field.SetUint(uint64(buff[buffPtr]))
			buffPtr += 1
		} else if fieldType == reflect.Uint16 || (fieldType == reflect.Uint && fieldSize == 2) {
			field.SetUint(uint64(binary.BigEndian.Uint16(buff[buffPtr : buffPtr+2])))
			buffPtr += 2
		} else if fieldType == reflect.Uint32 || (fieldType == reflect.Uint && fieldSize == 4) {
			field.SetUint(uint64(binary.BigEndian.Uint32(buff[buffPtr : buffPtr+4])))
			buffPtr += 4
		} else if fieldType == reflect.Uint64 || (fieldType == reflect.Uint && fieldSize == 8) {
			field.SetUint(uint64(binary.BigEndian.Uint64(buff[buffPtr : buffPtr+8])))
			buffPtr += 8
		} else if fieldType == reflect.Float32 {
			field.SetFloat(float64(math.Float32frombits(binary.BigEndian.Uint32(buff[buffPtr : buffPtr+4]))))
			buffPtr += 4
		} else if fieldType == reflect.Float64 {
			field.SetFloat(math.Float64frombits(binary.BigEndian.Uint64(buff[buffPtr : buffPtr+8])))
			buffPtr += 8
		} else if fieldType == reflect.Bool {
			if int(buff[buffPtr]) == 1 {
				field.SetBool(true)
			} else {
				field.SetBool(false)
			}
			buffPtr += 1
		} else if fieldType == reflect.String {
			maxLength, _ := getMaxStringLength(itemType.Field(i).Tag.Get("maxLength"))
			field.SetString(strings.TrimSpace(string(buff[buffPtr : buffPtr+uint64(maxLength)])))
			buffPtr += uint64(maxLength)
		}
	}
	return item
}

func calItemSize[T Item]() int {
	size := 0
	item := new(T)
	itemVal := reflect.ValueOf(item).Elem()
	itemType := reflect.TypeOf(*item)
	for i := 0; i < itemVal.NumField(); i++ {
		field := itemVal.FieldByIndex([]int{i})
		fieldType := field.Type().Kind()
		fieldSize := field.Type().Size()
		if !field.CanSet() {
			continue
		} else if fieldType == reflect.String {
			maxLength, _ := getMaxStringLength(itemType.Field(i).Tag.Get("maxLength"))
			size += maxLength
		} else {
			size += int(fieldSize)
		}
	}
	return size
}

func isValidItemFields[T Item]() error {
	itemVal := reflect.ValueOf(new(T)).Elem()
	for i := 0; i < itemVal.NumField(); i++ {
		field := itemVal.FieldByIndex([]int{i})
		if !field.CanSet() {
			continue
		}
		if slices.Contains(AVAILABLE_TYPES, field.Type().Kind()) {
			continue
		}
		return errors.New(fmt.Sprintf("Type %s is not allowed", field.Type().Kind()))
	}
	return nil
}

func isValidStringLabel[T Item]() error {
	item := new(T)
	itemVal := reflect.ValueOf(item).Elem()
	itemType := reflect.TypeOf(*item)
	for i := 0; i < itemVal.NumField(); i++ {
		maxLengthLabel := itemType.Field(i).Tag.Get("maxLength")
		if _, err := getMaxStringLength(maxLengthLabel); err != nil {
			return err
		}
	}
	return nil
}

func isValidStringLength[T Item](item *T) error {
	itemVal := reflect.ValueOf(item).Elem()
	itemType := reflect.TypeOf(*item)
	for i := 0; i < itemVal.NumField(); i++ {
		field := itemVal.FieldByIndex([]int{i})
		if !field.CanSet() {
			continue
		}
		maxLength, _ := getMaxStringLength(itemType.Field(i).Tag.Get("maxLength"))
		if field.Type().Kind() == reflect.String && (len(field.String()) > maxLength) {
			return errors.New(fmt.Sprintf("Length of string field should be less than %d", maxLength))
		}
	}
	return nil
}

func padSpaces(v string, maxLength int) string {
	numberOfPaddings := maxLength - len(v)
	for i := 0; i < numberOfPaddings; i++ {
		v += " "
	}
	return v
}

func getMaxStringLength(label string) (int, error) {
	if label == "" {
		return DEFAULT_STRING_MAX_LENGTH, nil
	} else {
		maxLength, err := strconv.Atoi(label)
		if err != nil {
			return 0, errors.New("maxLength should be number")
		}
		if maxLength <= 0 {
			return 0, errors.New("maxLength should be greater than 0")
		}
		return maxLength, nil
	}
}
