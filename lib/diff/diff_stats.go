package diff

import (
	"fmt"
	"strings"

	"github.com/metablink/godiff/lib"
)

// DiffStats stores diff information
type DiffStats struct {

	// Tracks added/removed column names
	AddedColumns   []string
	RemovedColumns []string

	// Tracks matched columns as keys
	// and difference counts values
	MatchedColumns map[string]int

	// The name of the key column.
	// Should be unique. We note duplicates.
	KeyColumn string

	// Columns that should be excluded from analysis.
	IgnoreColumns map[string]bool

	DuplicateToKeys   map[string]int
	DuplicateFromKeys map[string]int

	AddedRowCount   int
	RemovedRowCount int
	UpdatedRowCount int

	lastFromKey string
	lastToKey   string
}

// NewDiffStats returns a new configured DiffStats object in its initial state
func NewDiffStats(fromHeader []string, toHeader []string, keyColumn string, ignoreColumns map[string]bool) *DiffStats {

	fromMap := lib.SliceToSet(fromHeader)
	toMap := lib.SliceToSet(toHeader)

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
		KeyColumn:         keyColumn,
		IgnoreColumns:     ignoreColumns,
		AddedColumns:      addedColumns,
		RemovedColumns:    removedColumns,
		MatchedColumns:    matchedColumns,
		DuplicateFromKeys: map[string]int{},
		DuplicateToKeys:   map[string]int{},
	}

	return &stats
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
		lib.SetToSlice(s.IgnoreColumns), outputFunction)

	printIfNonempty("Added Columns:\t[%v]\n", s.AddedColumns, outputFunction)
	printIfNonempty("Removed Columns:\t[%v]\n", s.RemovedColumns, outputFunction)

	outputFunction("\nColumn Difference Counts:\n")
	for key, count := range s.MatchedColumns {
		outputFunction(fmt.Sprintf("\t%v:\t%v\n", key, count))
	}

	outputFunction(fmt.Sprintf("\nAdded Rows:\t%v\n", s.AddedRowCount))
	outputFunction(fmt.Sprintf("Removed Rows:\t%v\n", s.RemovedRowCount))
	outputFunction(fmt.Sprintf("Updated Rows:\t%v\n", s.UpdatedRowCount))
}

func printIfNonempty(format string, fields []string, outputFunction func(string, ...interface{}) (int, error)) {
	if len(fields) > 0 {
		outputFunction(fmt.Sprintf(format, strings.Join(fields, ", ")))
	}
}
