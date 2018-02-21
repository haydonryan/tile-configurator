package dictionary_test

import (
	"testing"

	. "github.com/haydonryan/tile-configurator/dictionary"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dictionary", func() {
	Context("Loading Dictionary", func() {
		It("can create a dictionary interface", func() {

			m := NewDictionary()

			Expect(m).ShouldNot(BeNil())
		})

		It("can load a dictionary", func() {
			d := NewDictionary()
			err := d.LoadDictionary("fixtures/test.yml")
			Expect(err).Should(BeNil())

		})
		It("will throw an error if the dictionary fails to load", func() {
			d := NewDictionary()
			err := d.LoadDictionary("invalid")
			Expect(err).ShouldNot(BeNil())

		})
		It("will throw an error if the dictionary is invalid", func() {
			d := NewDictionary()
			err := d.LoadDictionary("fixtures/invalid.yml")
			Expect(err).ShouldNot(BeNil())

		})
	})

	It("dictionary should have populated Simple", func() {
		d := NewDictionary()
		d.LoadDictionary("fixtures/test.yml")

		Expect(d.Simple[".properties.syslog.enabled.port"]).Should(Equal("syslog_port"))

	})
	It("dictionary should have populated Opsman with reverse lookup", func() {
		d := NewDictionary()
		d.LoadDictionary("fixtures/test.yml")
		Expect(d.Opsman["syslog_port"]).Should(Equal(".properties.syslog.enabled.port"))

	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "My Test Suite")
}
