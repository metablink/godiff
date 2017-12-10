package lib

import "testing"

type rowTestCase struct {
	from []string
	to   []string
}

var rowTestCases = []rowTestCase{
	rowTestCase{
		from: []string{"one", "two", "three"},
		to:   []string{"three", "two", "one"},
	},
}

func TestDiffRow(t *testing.T) {

	for _, testCase := range rowTestCases {
		// TODO
		DiffRow(testCase.from, testCase.to)
	}
}
