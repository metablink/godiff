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
		fromRow, toRow, err = getNextMatchingRow(s, from, to)

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
	if len(toRow) > 0 {
		remaining, err := getRemainingRowCount(to)

		if err != nil {
			return err
		}

		// Add 1 since the first remaining line has already been read
		s.AddedRowCount += remaining + 1
	}
	if len(fromRow) > 0 {
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

	// Grab the initial values
	fromRow, err := from.Next()
	if err != nil {
		return nil, nil, err
	}
	toRow, err := to.Next()
	if err != nil {
		return nil, nil, err
	}

	// If there's no key column specified, we assume these rows match
	if s.KeyColumn == "" {
		return fromRow, toRow, nil
	}

	// How many rows we've had to search through from the starting point
	fromOffset, toOffset := 0, 0

	for {

		// Either file can end first.
		if fromRow == nil || toRow == nil {
			break
		}

		fromKey := fromRow[s.KeyColumn]
		toKey := toRow[s.KeyColumn]

		// If we have key duplicates, record it
		if strings.Compare(s.lastFromKey, fromKey) == 0 {
			count, keyExists := s.DuplicateFromKeys[fromKey]

			if !keyExists {
				count = 0
			}

			s.DuplicateFromKeys[fromKey] = count + 1
		}
		if strings.Compare(s.lastToKey, toKey) == 0 {
			count, keyExists := s.DuplicateToKeys[toKey]

			if !keyExists {
				count = 0
			}

			s.DuplicateToKeys[fromKey] = count + 1
		}

		s.lastFromKey = fromKey
		s.lastToKey = toKey

		// If there's a match
		if strings.Compare(fromKey, toKey) == 0 {
			break
		}

		if strings.Compare(fromKey, toKey) < 0 {
			fromRow, err = from.Next()
			if err != nil {
				return nil, nil, err
			}

			fromOffset++
		} else {
			toRow, err = to.Next()
			if err != nil {
				return nil, nil, err
			}

			toOffset++
		}
	}

	// All 'from' rows evaluated have been removed
	s.RemovedRowCount += fromOffset
	// All 'to' row evaluated have been added
	s.AddedRowCount += toOffset

	return fromRow, toRow, nil
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
