package lib

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
)

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

// DiffFile runs a diff between the given Files
func DiffFile(fromFile *os.File, toFile *os.File) error {

	fromReader := csv.NewReader(fromFile)
	toReader := csv.NewReader(toFile)

	fromProvider := &CsvRowProvider{reader: fromReader}
	toProvider := &CsvRowProvider{reader: toReader}

	return DiffRowProvider(fromProvider, toProvider)
}

// DiffRowProvider runs a diff between the given RowProviders
func DiffRowProvider(from RowProvider, to RowProvider) error {
	providers := []RowProvider{from, to}

	for {
		compRows := make([]map[string]string, len(providers))
		for idx, provider := range providers {
			row, err := provider.Next()

			// Bail if we have an error, or have reached the end
			if err != nil || row == nil {
				return err
			}

			compRows[idx] = row
		}

		DiffRow(compRows[0], compRows[1])
	}
}

// DiffRow runs a diff between the given Rows
func DiffRow(from map[string]string, to map[string]string) {
	dmp := diffmatchpatch.New()

	// Run through the removed values and actual differences
	for key, fromVal := range from {
		toVal, toExists := to[key]

		if !toExists {
			continue
		}

		diffs := dmp.DiffMain(fromVal, toVal, false)
		fmt.Printf("|\t%v\t|", dmp.DiffPrettyText(diffs))
	}
	fmt.Println()
}
