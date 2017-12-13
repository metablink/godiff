package lib

import (
	"fmt"
	"log"
	"strings"
)

// DiffStats stores diff information
type DiffStats struct {

	// Data for comparison
	from RowProvider
	to   RowProvider

	// Tracks added/removed column names
	addedColumns   []string
	removedColumns []string

	// Tracks matched columns as keys
	// and difference counts values
	matchedColumns map[string]int

	// The name of the key column.
	// Should be unique. We note duplicates.
	keyColumn string

	// Columns that should be excluded from analysis.
	ignoreColumns map[string]bool

	duplicateToKeys   map[string]int
	duplicateFromKeys map[string]int

	addedRows   int
	removedRows int
	updatedRows int
}

// NewDiffStats returns a new DiffStats object
func NewDiffStats(from RowProvider, to RowProvider, keyColumn string, ignoreColumns map[string]bool) *DiffStats {

	fromMap, err := from.Next()
	if err != nil {
		log.Fatal(err)
	}

	toMap, err := to.Next()
	if err != nil {
		log.Fatal(err)
	}

	// Handle empty files
	if fromMap == nil {
		fromMap = make(map[string]string)
	}
	if toMap == nil {
		toMap = make(map[string]string)
	}

	addedColumns, removedColumns := []string{}, []string{}
	matchedColumns := map[string]int{}

	// Find the removed and matching columns
	for fromField := range fromMap {
		_, fieldMatch := toMap[fromField]

		if !fieldMatch {
			removedColumns = append(removedColumns, fromField)
			continue
		}

		// Only count differences when the column isn't ignored
		if _, ignore := ignoreColumns[fromField]; !ignore {
			// Initially all discrepancy counts are 0
			matchedColumns[fromField] = 0
		}
	}

	// Find the added columns
	for toField := range toMap {
		_, fieldMatch := fromMap[toField]

		if !fieldMatch {
			addedColumns = append(addedColumns, toField)
		}
	}

	stats := DiffStats{
		from:           from,
		to:             to,
		keyColumn:      keyColumn,
		ignoreColumns:  ignoreColumns,
		addedColumns:   addedColumns,
		removedColumns: removedColumns,
		matchedColumns: matchedColumns,
	}

	return &stats
}

// Diff runs a cell-by-cell diff between data sources
func (s *DiffStats) Diff() error {
	// We assume files coming in are already sorted
	fromRow, toRow, err := s.getNextMatchingRow()

	if err != nil {
		return err
	}

	// If either file isn't empty, count the remainder
	if fromRow != nil {
		remaining, err := getRemainingRowCount(s.to)

		if err != nil {
			return err
		}

		// Add 1 since the first remaining line has already been read
		s.addedRows += remaining + 1
	}
	if toRow != nil {
		remaining, err := getRemainingRowCount(s.from)

		if err != nil {
			return err
		}

		s.removedRows += remaining + 1
	}

	return nil
}

func (s *DiffStats) getNextMatchingRow() (fromRow map[string]string, toRow map[string]string, err error) {

	// TODO Make this scalable for large files
	// We could end up storing a lot of keys

	var (
		fromKey string
		toKey   string
	)

	// Maps of currently tracked keys and their corresponding offset from the starting row
	seenFromKeys, seenToKeys := map[string]int{}, map[string]int{}
	offset := 0

	for {

		fromRow, err = s.from.Next()
		if err != nil {
			break
		}

		toRow, err = s.to.Next()
		if err != nil {
			break
		}

		// If there's no key column specified, we assume these rows match
		if s.keyColumn == "" {
			// We return here instead of break so that we don't do post-processing
			return
		}

		// Either file can end first.
		if fromRow != nil {
			fromKey = fromRow[s.keyColumn]
		}
		if toRow != nil {
			toKey = toRow[s.keyColumn]
		}

		// If we have key duplicates, always keep the lowest offset
		if _, keyExists := seenFromKeys[fromKey]; !keyExists {
			seenFromKeys[fromKey] = offset
		}
		if _, keyExists := seenToKeys[toKey]; !keyExists {
			seenToKeys[toKey] = offset
		}

		// If the column keys are equal, we've found the right rows
		if cellEquals(fromKey, toKey) {
			break
		}

		offset++
	}

	s.addedRows += offset - seenFromKeys[fromKey]
	s.removedRows += offset - seenToKeys[toKey]

	return
}

func (s *DiffStats) diffRow(fromRow map[string]string, toRow map[string]string) {

	rowUpdated := false
	for key, colDiffCount := range s.matchedColumns {
		fromVal, toVal := fromRow[key], toRow[key]

		if !cellEquals(fromVal, toVal) {
			// TODO actually track non-aggregate differences
			s.matchedColumns[key] = colDiffCount + 1
			rowUpdated = true
		}
	}

	if rowUpdated {
		s.updatedRows++
	}
}

// Print prints the diff statistics to stdout
func (s *DiffStats) Print(outputFunction func(string, ...interface{}) (int, error)) {

	// Print the key column
	// if s.keyColumn != "" {
	// 	outputFunction(fmt.Sprintf(
	// 		"Key Column:\t%v\n\tDuplicate Keys: %v\n\n",
	// 		s.keyColumn, s.duplicateKeys))
	// }

	printIfNonempty(
		"Ignoring Columns:\t[%v]\n",
		SetToSlice(s.ignoreColumns), outputFunction)

	printIfNonempty("Added Columns:\t[%v]\n", s.addedColumns, outputFunction)
	printIfNonempty("Removed Columns:\t[%v]\n", s.removedColumns, outputFunction)

	outputFunction("\nColumn Difference Counts:\n")
	for key, count := range s.matchedColumns {
		outputFunction(fmt.Sprintf("\t%v:\t%v\n", key, count))
	}

	outputFunction(fmt.Sprintf("\nAdded Rows:\t%v\n", s.addedRows))
	outputFunction(fmt.Sprintf("Removed Rows:\t%v\n", s.removedRows))
	outputFunction(fmt.Sprintf("Updated Rows:\t%v\n", s.updatedRows))
}

func printIfNonempty(format string, fields []string, outputFunction func(string, ...interface{}) (int, error)) {
	if len(fields) > 0 {
		outputFunction(fmt.Sprintf(format, strings.Join(fields, ", ")))
	}
}

func cellEquals(fromVal string, toVal string) bool {
	// TODO support numerics, whitespace ignores, etc
	return strings.Compare(fromVal, toVal) == 0
}

func getRemainingRowCount(p RowProvider) (remaining int, err error) {

	for {
		row, err := p.Next()

		if err != nil || row == nil {
			break
		}
		remaining++
	}

	return
}
