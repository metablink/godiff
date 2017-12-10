package lib

import "encoding/csv"

// RowProvider is an interface to provide one row at a time
type RowProvider interface {
	Next() ([]string, error)
}

// CsvRowProvider provides one csv file row at a time
type CsvRowProvider struct {
	RowProvider
	header []string
	reader *csv.Reader
}

// Next provides the next row in the provider or nil if it doesn't exist
func (p CsvRowProvider) Next() (record []string, err error) {
	return p.reader.Read()
}
