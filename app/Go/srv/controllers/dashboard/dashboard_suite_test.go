package dashboard_controller_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDashboard(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dashboard Suite")
}
