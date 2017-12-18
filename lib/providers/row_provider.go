package providers

// RowProvider is an interface to provide one row at a time
type RowProvider interface {
	Next() (record []string, err error)
	Header() []string
}
