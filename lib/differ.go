package lib

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffRowProvider runs a diff between the given RowProviders
func DiffRowProvider(from RowProvider, to RowProvider) error {
	providers := []RowProvider{from, to}

	for {
		compRows, err := getNextProviderRows(providers)

		if err != nil || compRows == nil {
			return err
		}

		// OK, I know that this array use looks gross.
		// Still trying to find an elegant way to handle redundant errors.
		DiffRow(compRows[0], compRows[1])
	}
}

func getNextProviderRows(providers []RowProvider) (compRows []map[string]string, err error) {

	compRows = make([]map[string]string, len(providers))

	for idx, provider := range providers {
		row, err := provider.Next()

		// Bail if we have an error, or have reached the end
		if err != nil || row == nil {
			return nil, err
		}

		compRows[idx] = row
	}

	return compRows, nil
}

// DiffRow runs a diff between the given Rows
func DiffRow(from map[string]string, to map[string]string) {
	dmp := diffmatchpatch.New()

	// Run through the removed values and actual differences
	for key, fromVal := range from {
		toVal, toExists := to[key]

		if !toExists {
			continue
		}

		diffs := dmp.DiffMain(fromVal, toVal, false)
		fmt.Printf("|\t%v\t|", dmp.DiffPrettyText(diffs))
	}
	fmt.Println()
}
