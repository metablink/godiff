package lib

import (
	"encoding/csv"
)

// CsvRowProvider provides rows from a csv file
type CsvRowProvider struct {
	headerRowProvider
	reader *csv.Reader
}

// Next provides the next row in the provider or nil if it doesn't exist
func (p CsvRowProvider) Next() (record map[string]string, err error) {
	row, err := p.reader.Read()

	if err != nil {
		return nil, err
	}

	return p.bindHeader(row)
}
