package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var RespParseError = errors.New("not a valid RESP array")

// does only support array of strings (not integers)
func respArrayParse(str string) ([]string, error) {
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
		result[i], err = bulkDecode("$" + bulkStr)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func respStringEncode(str string) string {
	return fmt.Sprintf("+%s\r\n", str)
}
