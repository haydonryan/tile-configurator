package injest_test

import (
	"testing"

	. "github.com/haydonryan/tile-configurator/injest"
	"github.com/xchapter7x/lo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Injest", func() {
	Context("JSON", func() {
		It("should return nil if an empty map is passed to it and return an error", func() {

			input := make(map[string]interface{})
			result, err := ProcessInjest(input)

			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())

		})

		It("should remove the properties key", func() {
			filename := "fixtures/types.json"
			input, err := ReadJSON(filename)
			Expect(err).To(BeNil())

			result, err := ProcessInjest(input)
			Expect(err).To(BeNil())
			lo.G.Debug(result)
			//expect properties to not exist
			Expect(result).NotTo(ContainElement("properties"))
		})

		It("should contain a subkey and return copy the 'value field' as the key's value if it is a selector", func() { //ie .properties.backup_options
			filename := "fixtures/types.json"
			input, err := ReadJSON(filename)
			Expect(err).To(BeNil())

			result, err := ProcessInjest(input)
			Expect(err).To(BeNil())
			Expect(result[".properties.backup_options"]).To(Equal("enable"))
		})

		It("should ignore properties that are not configurable", func() { //ie .properties.backup_options
			filename := "fixtures/two.json"
			input, err := ReadJSON(filename)
			Expect(err).To(BeNil())

			result, err := ProcessInjest(input)
			Expect(err).To(BeNil())
			Expect(result[".properties.backup_value"]).To(BeNil())
		})

		It("should process integers", func() { //ie .properties.backup_options
			filename := "fixtures/types.json"
			input, err := ReadJSON(filename)
			Expect(err).To(BeNil())

			result, err := ProcessInjest(input)
			Expect(err).To(BeNil())
			Expect(result[".properties.backups.scp.port"]).To(Equal(int(22)))

		})
		It("should process strings", func() { //ie .properties.backup_options
			filename := "fixtures/types.json"
			input, err := ReadJSON(filename)
			Expect(err).To(BeNil())

			result, err := ProcessInjest(input)
			Expect(err).To(BeNil())
			Expect(result[".properties.backup_options.enable.cron_schedule"]).To(Equal(string("stringvalue")))

		})

		Context("Collections", func() {
			It("should be a map of an array", func() {
				filename := "fixtures/types.json"
				input, err := ReadJSON(filename)
				Expect(err).To(BeNil())

				result, err := ProcessInjest(input)
				lo.G.Debug(result)
				Expect(err).To(BeNil())

				_, correct := result[".properties.plan_collection"].([]interface{})
				Expect(correct).To(BeTrue())

			})

			It("deals with values that are nil", func() {
				filename := "fixtures/two.json"
				input, err := ReadJSON(filename)
				Expect(err).To(BeNil())

				result, err := ProcessInjest(input)
				lo.G.Debug(result)
				Expect(err).To(BeNil())
				Expect(result[".properties.buffer_pool_size.bytes.buffer_pool_size_bytes"]).To(Equal(0))
			})

		})

		//It("should ignore things that are not configurable")

	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "My Test Suite")
}
