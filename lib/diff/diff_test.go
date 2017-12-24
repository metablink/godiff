package diff_test

import (
	"github.com/metablink/godiff/lib/diff"
	. "github.com/metablink/godiff/lib/providers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var fromData = [][]string{
	[]string{"field1", "field2", "field3", "field4"},
	[]string{"five", "two", "three", "one"},
	[]string{"two", "three", "one", "two"},
	[]string{"three", "five", "two", "three"},
	[]string{"three", "six", "two", "four"},
}

var toData = [][]string{
	[]string{"field2", "field1", "field3", "field5"},
	[]string{"two", "one", "three", "one"},
	[]string{"three", "two", "one", "two"},
	[]string{"two", "three", "one", "three"},
	[]string{"one", "four", "three", "four"},
	[]string{"one", "four", "three", "four"},
	[]string{"one", "four", "three", "four"},
}

var _ = Describe("Diff", func() {
	var (
		fromRp *MapRowProvider
		toRp   *MapRowProvider
	)

	BeforeEach(func() {
		fromRp = &MapRowProvider{RowSrc: &StringRowProvider{Rows: fromData}}
		toRp = &MapRowProvider{RowSrc: &StringRowProvider{Rows: toData}}
	})

	Context("Test basic Diff - longer 'to' source", func() {
		var (
			ds  *diff.DiffStats
			err error
		)

		BeforeEach(func() {
			ds, err = diff.Diff(fromRp, toRp, "", map[string]bool{})
		})

		It("should return no error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have added columns", func() {
			Expect(ds.AddedColumns).To(ConsistOf("field5"))
		})

		It("should have removed columns", func() {
			Expect(ds.RemovedColumns).To(ConsistOf("field4"))
		})

		It("should have blank key column", func() {
			Expect(ds.KeyColumn).To(BeEmpty())
		})

		It("should have rows updated", func() {
			Expect(ds.UpdatedRowCount).To(Equal(3))
		})

		It("should have rows added", func() {
			Expect(ds.AddedRowCount).To(Equal(2))
		})

		It("should have no rows removed", func() {
			Expect(ds.RemovedRowCount).To(BeZero())
		})
	})

	Context("Test basic Diff - longer 'from' source", func() {
		var (
			ds  *diff.DiffStats
			err error
		)

		BeforeEach(func() {
			ds, err = diff.Diff(toRp, fromRp, "", map[string]bool{})
		})

		It("should return no error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have rows updated", func() {
			Expect(ds.UpdatedRowCount).To(Equal(3))
		})

		It("should have no rows added", func() {
			Expect(ds.AddedRowCount).To(BeZero())
		})

		It("should have rows removed", func() {
			Expect(ds.RemovedRowCount).To(Equal(2))
		})
	})

	Context("Test Diff with key column", func() {
		var (
			ds  *diff.DiffStats
			err error
		)

		keyColumn := "field1"

		BeforeEach(func() {
			ds, err = diff.Diff(fromRp, toRp, keyColumn, map[string]bool{})
		})

		It("should return no error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have key column", func() {
			Expect(ds.KeyColumn).To(Equal(keyColumn))
		})

		It("should have rows updated", func() {
			Expect(ds.UpdatedRowCount).To(Equal(1))
		})

		It("should have rows added", func() {
			Expect(ds.AddedRowCount).To(Equal(4))
		})

		It("should have rows removed", func() {
			Expect(ds.RemovedRowCount).To(Equal(2))
		})
	})

	// Context("Test getRemainingRowCount", func() {
	// 	It("should return correct counts", func() {
	// 		Expect(err).NotTo(HaveOccurred())
	// 	})
	// })

})
