package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const NullBulkString = "$-1\r\n"

func bulkEncode(str string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
}

var BulkDecodeError = errors.New("not a valid bulk string")

func bulkDecode(str string) (string, error) {
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
