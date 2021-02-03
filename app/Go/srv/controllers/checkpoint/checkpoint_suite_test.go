package checkpoint_controller_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCheckpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Checkpoint Suite")
}
