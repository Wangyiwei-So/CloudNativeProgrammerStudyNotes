package report_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ginkgotest/report"
)

var _ = Describe("Add", func() {
	It("add case", func() {
		Expect(4).To(Equal(report.Add(1, 3)))
	})
})
