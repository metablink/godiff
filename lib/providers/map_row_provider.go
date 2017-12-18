package providers

import "github.com/metablink/godiff/lib"

// MapRowProvider provides row values mapped to associated header keys
type MapRowProvider struct {
	RowSrc RowProvider
	header []string
}

// Header returns the header value or nil if it doesn't exist
func (p *MapRowProvider) Header() (header []string, err error) {
	if p.header == nil {
		p.header, err = p.RowSrc.Next()
	}

	return p.header, err
}

// Next returns the next row values, mapped to the header keys.
func (p *MapRowProvider) Next() (record map[string]string, err error) {

	// Get the header
	header, err := p.Header()
	if header == nil || err != nil {
		return nil, err
	}

	// Get the row from the internal RowProvider
	row, err := p.RowSrc.Next()
	if row == nil || err != nil {
		return nil, err
	}

	return lib.BindHeader(header, row)
}
