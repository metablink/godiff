package lib

// RowProvider is an interface to provide one row at a time
type RowProvider interface {
	Next() (record map[string]string, err error)
}

// headerRowProvider provides base provider behavior and shouldn't be used directly
type headerRowProvider struct {
	RowProvider
	header []string
}

func (p *headerRowProvider) bindHeader(row []string) (record map[string]string, err error) {
	if p.header == nil {
		// TODO allow headerless files
		// NOTE we will see unexpected behavior with duplicate keys
		p.header = row
	}

	return BindHeader(p.header, row)
}
