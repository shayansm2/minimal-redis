package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RespStr string
type BulkStr string

var RespParseError = fmt.Errorf("not a valid RESP array")

func respArrayBulkStringParse(str string) ([]string, error) {
	header, body, found := strings.Cut(str, "\r\n")
	if !found {
		return nil, RespParseError
	}
	count, found := strings.CutPrefix(header, "*")
	if !found {
		return nil, RespParseError
	}
	length, err := strconv.Atoi(count)
	if err != nil {
		return nil, RespParseError
	}

	bulkArray := strings.Split(body, "$")
	bulkArray = bulkArray[1:] // drop first empty string

	if len(bulkArray) != length {
		return nil, RespParseError
	}
	result := make([]string, length)
	for i, bulkStr := range bulkArray {
		result[i], err = bulkStringDecode("$" + bulkStr)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func toRespArray(array []string) string {
	return fmt.Sprintf("*%d\r\n%s", len(array), strings.Join(array, ""))
}

func toRespSimpleString(str RespStr) string {
	return fmt.Sprintf("+%s\r\n", str)
}

func toRespError(err error) string {
	return fmt.Sprintf("-%s\r\n", err)
}

func toRespInteger(num int) string {
	return fmt.Sprintf(":%d\r\n", num)
}

const NullBulkString = "$-1\r\n"
const NullArray = "*-1\r\n"

func toBulkString(str BulkStr) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
}

func toBulkArray(arr []BulkStr) []string {
	result := make([]string, len(arr))
	for i, str := range arr {
		result[i] = toBulkString(str)
	}
	return result
}

var BulkDecodeError = errors.New("not a valid bulk string")

func bulkStringDecode(str string) (string, error) {
	parts := strings.Split(str, "\r\n")
	if len(parts) != 3 {
		return "", BulkDecodeError
	}
	decoded := parts[1]
	count, found := strings.CutPrefix(parts[0], "$")
	if !found {
		return "", BulkDecodeError
	}
	length, err := strconv.Atoi(count)
	if err != nil {
		return "", BulkDecodeError
	}
	if len(decoded) != length {
		return "", BulkDecodeError
	}
	return decoded, nil
}

func encode(value any) (string, error) {
	if value == nil {
		return NullBulkString, nil
	}
	switch v := value.(type) {
	case int:
		return toRespInteger(v), nil
	case RespStr:
		return toRespSimpleString(v), nil
	case BulkStr:
		return toBulkString(v), nil
	case []BulkStr:
		return toRespArray(toBulkArray(v)), nil
	case []string:
		if v == nil {
			return NullArray, nil
		}
		return toRespArray(v), nil
	case error:
		return toRespError(v), nil
	default:
		return "", fmt.Errorf("Err not implemented type")
	}
}
