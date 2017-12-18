package diff_test

import (
	"github.com/metablink/godiff/lib/diff"
	. "github.com/metablink/godiff/lib/providers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var fromData = [][]string{
	[]string{"field1", "field2", "field3"},
	[]string{"five", "two", "three"},
	[]string{"two", "three", "one"},
	[]string{"three", "five", "two"},
	[]string{"three", "six", "two"},
}

var toData = [][]string{
	[]string{"field2", "field1", "field3"},
	[]string{"two", "one", "three"},
	[]string{"three", "two", "one"},
	[]string{"two", "three", "one"},
	[]string{"one", "four", "three"},
}

var _ = Describe("Diff Stats", func() {
	var (
		fromRp *MapRowProvider
		toRp   *MapRowProvider
	)

	BeforeEach(func() {
		fromRp = &MapRowProvider{RowSrc: &StringRowProvider{Rows: fromData}}
		toRp = &MapRowProvider{RowSrc: &StringRowProvider{Rows: toData}}
	})

	Context("Test NewDiffStats", func() {

		It("should return no error", func() {
			_, err := diff.Diff(fromRp, toRp, "", map[string]bool{})
			Expect(err).NotTo(HaveOccurred())
			// Expect(ds.)
		})

		It("should return a blank map", func() {
			Expect(map[string]string{}).To(Equal(make(map[string]string)))
		})
	})

})
