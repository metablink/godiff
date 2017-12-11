package lib

// StringRowProvider provides rows from a two-dimensional string array
type StringRowProvider struct {
	headerRowProvider
	rows       [][]string
	currentRow int
}

// Next provides the next row in the provider or nil if it doesn't exist
func (p *StringRowProvider) Next() (record map[string]string, err error) {

	// Handle the case where we have exhausted rows
	if p.currentRow >= len(p.rows) {
		return nil, nil
	}

	row := p.rows[p.currentRow]
	p.currentRow++

	return p.bindHeader(row)
}
