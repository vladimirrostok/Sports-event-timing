package result_controller_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResult(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Result Suite")
}
