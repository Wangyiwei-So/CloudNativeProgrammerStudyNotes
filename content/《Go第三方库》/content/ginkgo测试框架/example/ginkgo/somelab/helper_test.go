package somelab_test

import (
	"ginkgotest/somelab"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {
	var result float32
	BeforeEach(func() {
		result = 3
	})
	Describe("Add函数", func() {
		It("getname函数测试", func() {
			Expect(result).To(Equal(somelab.Add(1, 2)))
		})
	})
})
