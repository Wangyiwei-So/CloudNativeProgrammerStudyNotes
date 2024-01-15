package somelab_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ginkgotest/somelab"
)

var _ = Describe("User测试", func() {
	var user *somelab.User
	var name string
	BeforeEach(func() {
		user = &somelab.User{}
		name = "wyw"
	})
	Describe("User名字测试", func() {
		Context("User名字set/get测试", func() {
			It("case1", func() {
				user.SetName(name)
				Expect("wyw").To(Equal(user.GetName())) //故意搞错
			})
			It("case2", func() {
				user.SetName(name)
				Expect("xqq").To(Equal(user.GetName())) //故意搞错
			})
		})
	})
})
