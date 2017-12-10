package lib

import "testing"

type rowTestCase struct {
	from []string
	to   []string
}

var rowTestCases = []rowTestCase{}

func TestDiffRow(t *testing.T) {

	for _, testCase := range rowTestCases {
		// TODO
		DiffRow(testCase.from, testCase.to)
	}
}
