package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	ValueTypeStr = iota
	ValueTypeInt
)

const HEADER_LEN = 9

const OP_CODE_AUX = 0xFA
const OP_CODE_RESIZEDB = 0xFB
const OP_CODE_EXPIRETIMEMS = 0xFC
const OP_CODE_EXPIRETIME = 0xFD
const OP_CODE_SELECTDB = 0xFE

var dir string
var dbFileName string

func init() {
	dir = getConfigs().get("dir", ".")
	dbFileName = getConfigs().get("dbfilename", "")
}

func loadRDB() error {
	if dbFileName == "" {
		return nil
	}
	file, err := os.Open(fmt.Sprintf("%s/%s", dir, dbFileName))
	if err != nil {
		fmt.Printf("no rbd found in %s/%s", dir, dbFileName)
		return nil
	}
	return load(file)
}

func load(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}
	dump := make([]byte, stat.Size())
	_, err = file.Read(dump)
	if err != nil {
		fmt.Println(err)
	}

	meOffset := metadataSectionEndOffset(dump)
	dbOffset, _ := findOffset(dump, meOffset+1, OP_CODE_SELECTDB)
	htOffset, err := findOffset(dump, dbOffset+1, OP_CODE_RESIZEDB)
	if err != nil {
		return err
	}

	kv, err := decodeHashTable(dump, htOffset)
	if err != nil {
		return err
	}

	for k, v := range kv {
		db.set(k, v, nil)
	}
	return nil
}

var OffsetNotFound = errors.New("Offset Not Found")

func metadataSectionEndOffset(dump []byte) int {
	i := HEADER_LEN
	for dump[i] == OP_CODE_AUX {
		i++
		len, offset, _ := getLen(dump[i:])
		i += len + offset
		len, offset, _ = getLen(dump[i:])
		i += len + offset
	}
	return i
}

func findOffset(dump []byte, start int, find byte) (int, error) {
	for i := start; i < len(dump); i++ {
		if dump[i] == find {
			return i, nil
		}
	}
	return start - 1, OffsetNotFound
}

func decodeHashTable(dump []byte, offset int) (map[string]string, error) {
	kv := make(map[string]string)
	i := offset
	i++
	// find number of key-values
	hashTableLen, offset, _ := getLen(dump[i:])
	i += offset
	expLen, offset, _ := getLen(dump[i:])
	i += offset
	// for i in number of key values
	for hashTableLen > 0 {
		// get expiry
		if dump[i] == OP_CODE_EXPIRETIMEMS {
			if expLen == 0 {
				return nil, fmt.Errorf("decode error: wrong number of expire keys")
			}
			i++
			i += 8
			expLen--
		} else if dump[i] == OP_CODE_EXPIRETIME {
			if expLen == 0 {
				return nil, fmt.Errorf("decode error: wrong number of expire keys")
			}
			i++
			i += 4
			expLen--
		}

		if valType := dump[i]; valType != 0 {
			return nil, fmt.Errorf("only supporting strings yet")
		}
		i++
		// get key
		key, offset := decodeFirstElement(dump[i:])
		i += offset
		// get value
		value, offset := decodeFirstElement(dump[i:])
		i += offset
		kv[key] = value
		hashTableLen--
	}
	return kv, nil
}

func decodeFirstElement(dump []byte) (string, int) {
	len, offset, valueType := getLen(dump)
	if valueType == ValueTypeStr {
		return decodeStr(dump[offset:], len), len + offset
	} else if valueType == ValueTypeInt {
		return strconv.Itoa(decodeInt(dump[offset:], len)), len + offset
	}
	return "", 0 // not handled yet
}

func getLen(dump []byte) (len, offset, valType int) {
	if dump[0] < 0b01000000 {
		return int(dump[0]), 1, ValueTypeStr
	}
	if dump[0] < 0b10000000 {
		return (int(dump[0])-0b01000000)*256 + int(dump[1]), 2, ValueTypeStr
	}
	if dump[0] < 0b11000000 {
		return decodeInt(dump[1:], 4), 5, ValueTypeStr
	}

	switch dump[0] {
	case 0b11000000:
		return 1, 1, ValueTypeInt
	case 0b11000001:
		return 2, 1, ValueTypeInt
	case 0b11000010:
		return 4, 1, ValueTypeInt
	default:
		return 0, 0, 0 // not handled yet
	}
}

func decodeStr(dump []byte, len int) string {
	return string(dump[0:len])
}

func decodeInt(dump []byte, len int) int {
	var result int
	for i := 0; i < len; i++ {
		num := int(dump[i])
		result = (256 * result) + num
	}
	return result
}
