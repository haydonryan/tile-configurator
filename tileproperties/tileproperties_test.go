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
	Context("Input", func() {
		It("Can  read a YAML file", func() {
			m := NewTileProperties()
			m.ReadYAML("fixtures/types.yml")
			Expect(m).ShouldNot(BeNil())
			Expect(m.Properties[".properties.optional_protections.enable.canary_poll_frequency"]).Should(Equal(69))

		})
		It("Can  read a JSON file", func() {
			m := NewTileProperties()
			m.ReadJSON("fixtures/types.json")
			//Expect(m).Should(ContainElement(".properties.backup_options"))
			Expect(m.Properties["properties"]).ShouldNot(BeNil())
		})
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "My Test Suite")
}
