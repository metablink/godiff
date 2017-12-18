package diff

import (
	"strings"

	"github.com/metablink/godiff/lib/providers"
)

// Diff runs a cell-by-cell diff between data sources
func Diff(
	from *providers.MapRowProvider,
	to *providers.MapRowProvider,
	keyColumn string,
	ignoreColumns map[string]bool) (*DiffStats, error) {

	fromHeader, err := from.Header()
	if err != nil {
		return nil, err
	}

	toHeader, err := to.Header()
	if err != nil {
		return nil, err
	}

	s := NewDiffStats(fromHeader, toHeader, keyColumn, ignoreColumns)
	err = runDiff(s, from, to)

	return s, err
}

func runDiff(s *DiffStats, from *providers.MapRowProvider, to *providers.MapRowProvider) error {

	var (
		fromRow map[string]string
		toRow   map[string]string
		err     error
	)

	for {
		// We assume files coming in are already sorted
		fromRow, toRow, err := getNextMatchingRow(s, from, to)

		if err != nil {
			return err
		}

		// Keep running the loop until a provider runs out
		if fromRow == nil || toRow == nil {
			break
		}

		diffRow(s, fromRow, toRow)
	}

	// If either file isn't empty, count the remainder
	if fromRow != nil {
		remaining, err := getRemainingRowCount(to)

		if err != nil {
			return err
		}

		// Add 1 since the first remaining line has already been read
		s.AddedRowCount += remaining + 1
	}
	if toRow != nil {
		remaining, err := getRemainingRowCount(from)

		if err != nil {
			return err
		}

		s.RemovedRowCount += remaining + 1
	}

	return err
}

func getNextMatchingRow(
	s *DiffStats,
	from *providers.MapRowProvider,
	to *providers.MapRowProvider) (map[string]string, map[string]string, error) {

	// TODO Make this scalable for large files
	// We could end up storing a lot of keys

	var (
		fromKey     string
		toKey       string
		matchingKey string
	)

	// Maps of currently tracked keys and their corresponding offset from the starting row
	seenFromKeys, seenToKeys := map[string]int{}, map[string]int{}

	// Stores the seen, but currently unmatched rows
	fromQueue, toQueue := []map[string]string{}, []map[string]string{}

	// How many rows we've had to search through from the starting point
	offset := 0

	for {

		fromRow, err := from.Next()
		if err != nil {
			break
		}

		toRow, err := to.Next()
		if err != nil {
			break
		}

		// If there's no key column specified, we assume these rows match
		if s.KeyColumn == "" {
			// We return here instead of break so that we don't do post-processing
			return fromRow, toRow, err
		}

		// Either file can end first.
		if fromRow == nil || toRow == nil {

			// All 'from' rows evaluated have been removed
			s.RemovedRowCount += offset
			// All 'to' row evaluated have been added
			s.AddedRowCount += offset

			return fromRow, toRow, err
		}

		fromQueue = append(fromQueue, fromRow)
		toQueue = append(toQueue, toRow)

		fromKey = fromRow[s.KeyColumn]
		toKey = toRow[s.KeyColumn]

		// If we have key duplicates, always keep the lowest offset
		if _, keyExists := seenFromKeys[fromKey]; !keyExists {
			seenFromKeys[fromKey] = offset
		} else {
			// TODO duplicates
		}

		if _, keyExists := seenToKeys[toKey]; !keyExists {
			seenToKeys[toKey] = offset
		} else {
			// TODO duplicates
		}

		// If we have a match
		if _, keyExists := seenFromKeys[toKey]; keyExists {
			matchingKey = toKey
			break
		}
		if _, keyExists := seenToKeys[fromKey]; keyExists {
			matchingKey = fromKey
			break
		}

		offset++
	}

	fromRowOffset, toRowOffset := seenFromKeys[matchingKey], seenToKeys[matchingKey]

	s.AddedRowCount += offset - fromRowOffset
	s.RemovedRowCount += offset - toRowOffset

	return fromQueue[fromRowOffset], toQueue[toRowOffset], nil
}

func diffRow(s *DiffStats, fromRow map[string]string, toRow map[string]string) {

	rowUpdated := false
	for key, colDiffCount := range s.MatchedColumns {
		fromVal, toVal := fromRow[key], toRow[key]

		if !cellEquals(fromVal, toVal) {
			// TODO actually track non-aggregate differences
			s.MatchedColumns[key] = colDiffCount + 1
			rowUpdated = true
		}
	}

	if rowUpdated {
		s.UpdatedRowCount++
	}
}

func cellEquals(fromVal string, toVal string) bool {
	// TODO support numerics, whitespace ignores, etc
	return strings.Compare(fromVal, toVal) == 0
}

func getRemainingRowCount(p *providers.MapRowProvider) (remaining int, err error) {

	for {
		row, err := p.Next()

		if err != nil || row == nil {
			break
		}
		remaining++
	}

	return
}