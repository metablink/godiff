package lib

import (
	"errors"
	"fmt"
	"strings"
)

// StringToSet converts a delim-separated string into a set
func StringToSet(str string, delim string) map[string]bool {
	split := strings.Split(str, delim)

	if len(split) == 1 {
		return map[string]bool{}
	}

	return SliceToSet(split)
}

// SliceToSet converts a slice into a set
func SliceToSet(slice []string) map[string]bool {
	if len(slice) == 1 && slice[0] == "" {
		return map[string]bool{}
	}

	converted := make(map[string]bool, len(slice))

	// Create a map to make lookups easier
	for _, field := range slice {
		converted[field] = true
	}

	return converted
}

// SetToSlice converts a set into a slice
func SetToSlice(inMap map[string]bool) (slice []string) {
	for key := range inMap {
		slice = append(slice, key)
	}

	return
}

// BindHeader returns a map of header column keys to row values
func BindHeader(header []string, row []string) (record map[string]string, err error) {
	rowLen := len(row)
	headerLen := len(header)

	if rowLen != headerLen {
		errMsg := fmt.Sprintf("invalid row: column count %d doesn't match the %d column header", rowLen, headerLen)
		return nil, errors.New(errMsg)
	}

	record = make(map[string]string)

	for fieldIdx, fieldName := range header {
		record[fieldName] = row[fieldIdx]
	}

	return record, nil
}
