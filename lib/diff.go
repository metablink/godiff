package lib

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
)

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
func DiffRow(from []string, to []string) {
	dmp := diffmatchpatch.New()
	maxLen := int(math.Max(float64(len(from)), float64(len(to))))

	for i := 0; i < maxLen; i++ {
		if i >= len(from) {
			// TODO
			break
		}
		if i >= len(to) {
			// TODO
			break
		}

		diffs := dmp.DiffMain(from[i], to[i], false)
		fmt.Println(dmp.DiffPrettyText(diffs))
	}
}

// func DiffRow()
