package lib

// RowProvider is an interface to provide one row at a time
type RowProvider interface {
	Next() (map[string]string, error)
}

// headerRowProvider provides base provider behavior and shouldn't be used directly
type headerRowProvider struct {
	RowProvider
	header []string
}

func (p headerRowProvider) bindHeader(row []string) (record map[string]string, err error) {
	return BindHeader(p.header, row)
}
