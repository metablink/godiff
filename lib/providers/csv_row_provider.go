package providers

import (
	"encoding/csv"
	"io"
)

// CsvRowProvider provides rows from a csv file
type CsvRowProvider struct {
	RowProvider
	Reader *csv.Reader
}

// Next provides the next row in the provider or nil if it doesn't exist
func (p *CsvRowProvider) Next() (record []string, err error) {
	row, err := p.Reader.Read()

	// Handle the case where we have exhausted rows
	if err == io.EOF {
		return nil, nil
	}

	return row, err
}
