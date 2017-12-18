package providers

// StringRowProvider provides rows from a two-dimensional string array
type StringRowProvider struct {
	RowProvider
	Rows       [][]string
	currentRow int
}

// Next provides the next row in the provider or nil if it doesn't exist
func (p *StringRowProvider) Next() (record []string, err error) {

	// Handle the case where we have exhausted rows
	if p.currentRow >= len(p.Rows) {
		return nil, nil
	}

	row := p.Rows[p.currentRow]
	p.currentRow++

	return row, nil
}
