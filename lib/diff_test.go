package lib_test

import (
	"fmt"

	. "github.com/metablink/godiff/lib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type rowTestCase struct {
	from []string
	to   []string
}

func populateTestCases(testCaseHeaders [][]string, testCaseRows [][]string) (records []map[string]string, errors []error) {
	testCaseCount := len(testCaseHeaders)

	records = make([]map[string]string, testCaseCount)
	errors = make([]error, testCaseCount)

	for idx := 0; idx < testCaseCount; idx++ {
		header := testCaseHeaders[idx]
		row := testCaseRows[idx]

		records[idx], errors[idx] = BindHeader(header, row)
	}

	return records, errors
}

var _ = Describe("Diff", func() {

	Describe("Test Bind Header", func() {

		var (
			records []map[string]string
			errors  []error
		)

		Context("Blank Pair", func() {
			var (
				record map[string]string
				err    error
			)

			BeforeEach(func() {
				record, err = BindHeader([]string{}, []string{})
			})

			It("should return no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return a blank map", func() {
				Expect(record).To(Equal(make(map[string]string)))
			})
		})

		Context("Mismatched Column Count", func() {
			testCaseHeaders := [][]string{
				[]string{
					"fieldOne",
					"fieldTwo",
					"fieldThree",
				},
				[]string{
					"fieldOne",
					"fieldTwo",
				},
				[]string{
					"fieldOne",
					"fieldTwo",
					"fieldThree",
				},
			}

			testCaseRows := [][]string{
				[]string{
					"valOne",
				},
				[]string{
					"valOne",
					"valTwo",
					"valThree",
				},
				[]string{},
			}

			BeforeEach(func() {
				records, errors = populateTestCases(testCaseHeaders, testCaseRows)
			})

			It("should cause an error", func() {
				for _, err := range errors {
					Expect(err).To(HaveOccurred())
					Expect(fmt.Sprint(err)).To(MatchRegexp(".*invalid row.*column count.*"))
				}
			})

			It("should return a nil record", func() {
				for _, record := range records {
					Expect(record).To(BeNil())
				}
			})
		})

		Context("Valid Inputs", func() {

			testCaseHeaders := [][]string{
				[]string{
					"fieldOne",
					"fieldTwo",
					"fieldThree",
				},
				[]string{
					"fieldOne",
					"fieldTwo",
				},
				[]string{
					"fieldOne",
				},
			}

			testCaseRows := [][]string{
				[]string{
					"valOne",
					"valTwo",
					"valThree",
				},
				[]string{
					"valOne",
					"valTwo",
				},
				[]string{
					"valOne",
				},
			}

			BeforeEach(func() {
				records, errors = populateTestCases(testCaseHeaders, testCaseRows)
			})

			It("should return no error", func() {
				for _, err := range errors {
					Expect(err).NotTo(HaveOccurred())
				}
			})

			It("should match field count", func() {
				for rowIdx, record := range records {
					row := testCaseRows[rowIdx]
					Expect(record).To(HaveLen(len(row)))
				}
			})

			It("should contain expected field/value pair", func() {
				for rowIdx, record := range records {
					header := testCaseHeaders[rowIdx]
					row := testCaseRows[rowIdx]
					for fieldIdx, fieldKey := range header {
						fieldVal := row[fieldIdx]
						Expect(record).To(HaveKeyWithValue(fieldKey, fieldVal))
					}
				}
			})
		})
	})

	Describe("Test Diff File", func() {
		// TODO
	})

	Describe("Test Diff RowProvider", func() {
		// TODO
	})

	Describe("Test Diff Row", func() {
		// TODO
	})
})
