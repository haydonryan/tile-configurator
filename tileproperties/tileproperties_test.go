package tileproperties_test

import (
	"testing"

	. "github.com/haydonryan/tile-configurator/tileproperties"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Properties", func() {
	It("Can create a properties interface", func() {
		m := NewTileProperties()
		Expect(m).ShouldNot(BeNil())
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "My Test Suite")
}
