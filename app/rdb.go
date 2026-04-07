package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
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

	kv, exp, err := decodeHashTable(dump, htOffset)
	if err != nil {
		return err
	}

	for k, v := range kv {
		db.set(k, v)
		if ns, found := exp[k]; found {
			now := time.Now().UnixMilli()
			fmt.Println(now, ns)
			if now >= ns {
				delete(db, k)
			} else {
				db.expire(k, int(ns-now))
			}
		}
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

func decodeHashTable(dump []byte, offset int) (kv map[string]string, exp map[string]int64, err error) {
	kv = make(map[string]string)
	exp = make(map[string]int64)

	i := offset
	i++
	hashTableLen, offset, _ := getLen(dump[i:])
	i += offset
	expLen, offset, _ := getLen(dump[i:])
	i += offset
	for hashTableLen > 0 {
		var ns int64
		if dump[i] == OP_CODE_EXPIRETIMEMS {
			if expLen == 0 {
				err = fmt.Errorf("decode error: wrong number of expire keys")
				return
			}
			i++
			ns = decodeTimestamp(dump[i : i+8])
			i += 8
			expLen--
		} else if dump[i] == OP_CODE_EXPIRETIME {
			if expLen == 0 {
				err = fmt.Errorf("decode error: wrong number of expire keys")
				return
			}
			i++
			ns = decodeTimestamp(dump[i : i+4])
			i += 4
			expLen--
		}

		if valType := dump[i]; valType != 0 {
			err = fmt.Errorf("only supporting strings yet")
			return
		}
		i++
		// get key
		key, offset := decodeFirstElement(dump[i:])
		i += offset
		// get value
		value, offset := decodeFirstElement(dump[i:])
		i += offset
		kv[key] = value
		if ns != 0 {
			exp[key] = ns
		}
		hashTableLen--
	}
	return
}

func decodeFirstElement(dump []byte) (string, int) {
	len, offset, valueType := getLen(dump)
	if valueType == ValueTypeStr {
		return decodeStr(dump[offset:], len), len + offset
	} else if valueType == ValueTypeInt {
		return strconv.FormatInt(decodeInt(dump[offset:], len), 10), len + offset
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
		return int(decodeInt(dump[1:], 4)), 5, ValueTypeStr
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

func decodeInt(dump []byte, len int) int64 {
	var result int64
	for i := 0; i < len; i++ {
		num := int(dump[i])
		result = (256 * result) + int64(num)
	}
	return result
}

func decodeTimestamp(dump []byte) int64 {
	bytes := make([]byte, len(dump))
	copy(bytes, dump)
	slices.Reverse(bytes)
	res := decodeInt(bytes, len(bytes))
	if len(bytes) == 8 {
		return res
	}
	return res * 1e3
}
