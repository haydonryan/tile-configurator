package tileproperties_test

import (
	"testing"

	. "github.com/haydonryan/tile-configurator/tileproperties"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Properties", func() {
	Context("Setup", func() {

		It("Can create a properties interface", func() {
			m := NewTileProperties()
			Expect(m).ShouldNot(BeNil())
		})
	})
	Context("Input", func() {

		Context("YAML", func() {
			It("will throw an error if the file fails to load", func() {
				m := NewTileProperties()
				_, err := m.ReadYAML("invalid")
				Expect(err).ShouldNot(BeNil())
			})

			It("will throw an error if the file is invalid", func() {
				m := NewTileProperties()
				_, err := m.ReadYAML("fixtures/invalid.yaml")
				Expect(err).ShouldNot(BeNil())
			})

			It("Can read a YAML file", func() {
				m := NewTileProperties()
				m.ReadYAML("fixtures/types.yml")
				Expect(m).ShouldNot(BeNil())
				Expect(m.Properties[".properties.optional_protections.enable.canary_poll_frequency"]).Should(Equal(69))
			})
		})
		Context("JSON", func() {
			It("will throw an error if the file fails to load", func() {
				m := NewTileProperties()
				_, err := m.ReadJSON("invalid")
				Expect(err).ShouldNot(BeNil())
			})

			It("will throw an error if the file is invalid", func() {
				m := NewTileProperties()
				_, err := m.ReadJSON("fixtures/invalid.json")
				Expect(err).ShouldNot(BeNil())
			})

			It("Can  read a JSON file", func() {
				m := NewTileProperties()
				m.ReadJSON("fixtures/types.json")
				Expect(m.Properties["properties"]).ShouldNot(BeNil())
			})

		})

	})
	Context("Output", func() {
		Context("JSON", func() {
			It("Converts the map to a json object", func() {
				m := NewTileProperties()

				m.Properties["test"] = "testy"
				str := m.MakeJSON()
				Expect(str).Should(Equal("{\"test\":\"testy\"}"))

			})
		})
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "My Test Suite")
}
