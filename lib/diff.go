package lib

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
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
func DiffFile(fromFile *os.File, toFile *os.File) {

	fromReader := csv.NewReader(fromFile)
	toReader := csv.NewReader(toFile)

	fromProvider := CsvRowProvider{reader: fromReader}
	toProvider := CsvRowProvider{reader: toReader}

	DiffRowProvider(fromProvider, toProvider)
}

// DiffRowProvider runs a diff between the given RowProviders
func DiffRowProvider(from RowProvider, to RowProvider) {

	for {

		fromRow, fromErr := from.Next()
		toRow, toErr := to.Next()

		if fromErr != nil {
			log.Fatal(fromErr)
		}
		if toErr != nil {
			log.Fatal(toErr)
		}

		if fromRow == nil || toRow == nil {
			break
		}

		DiffRow(fromRow, toRow)
	}

	return
}

// DiffRow runs a diff between the given Rows
func DiffRow(from map[string]string, to map[string]string) {
	dmp := diffmatchpatch.New()

	addedKeys := make(map[string]bool)
	removedKeys := make(map[string]bool)

	// Run through the removed values and actual differences
	for key, fromVal := range from {
		toVal, toExists := to[key]

		if !toExists {
			removedKeys[key] = true
			continue
		}

		diffs := dmp.DiffMain(fromVal, toVal, false)
		fmt.Println(dmp.DiffPrettyText(diffs))
	}

	// Run through the added values
	for key := range to {
		_, fromExists := to[key]

		if !fromExists {
			addedKeys[key] = true
			continue
		}
	}
}

// func DiffRow()
