package report_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/reporters"
	. "github.com/onsi/gomega"
)

func TestReport(t *testing.T) {
	RegisterFailHandler(Fail)
	ReportAfterSuite("custom report", func(r Report) {
		// 注意如果文件夹不存在它不会自动创建
		reporters.GenerateJUnitReport(r, "../report.xml")
	})
	RunSpecs(t, "Report Suite")
}
