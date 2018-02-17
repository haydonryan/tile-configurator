package diffbase_test

import (
	"reflect"

	"github.com/haydonryan/tile-configurator/diffbase"
	//. "github.com/haydonryan/tile-configurator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diffbase", func() {

	Context("Diffbase()", func() {
		It("should not return nil", func() {
			base := make(map[string]interface{})

			ret := diffbase.Diffbase(base, base)
			Expect(ret).ShouldNot(BeNil())
		})

		It("should return if new map has extra key", func() {
			base := make(map[string]interface{})
			changed := make(map[string]interface{})

			changed["test"] = "test string"

			ret := diffbase.Diffbase(base, changed)
			Expect(ret).Should(Equal(changed))
		})

		It("should return if new map has multiple keys", func() {
			base := make(map[string]interface{})
			changed := make(map[string]interface{})

			changed["test"] = "test string"
			changed["test2"] = "test string2"

			ret := diffbase.Diffbase(base, changed)
			Expect(ret).Should(Equal(changed))
		})

		It("should return the difference to include a subkey of the main key", func() {
			base := make(map[string]interface{})
			changed := make(map[string]interface{})
			sub := make(map[string]interface{})

			sub["subkey"] = "test string"
			changed["test2"] = sub

			ret := diffbase.Diffbase(base, changed)
			Expect(ret).Should(Equal(changed))
		})

		It("should return only one subkey if the other subkey is present", func() {
			base := make(map[string]interface{})
			changed := make(map[string]interface{})
			sub := make(map[string]interface{})
			sub2 := make(map[string]interface{})
			sub3 := make(map[string]interface{})

			expected := make(map[string]interface{})

			sub["subkey"] = "test string"
			base["test"] = sub
			sub2["subkey"] = "test string"
			sub2["difkey"] = "diff string"
			changed["test"] = sub2

			sub3["difkey"] = "diff string"
			expected["test"] = sub3

			ret := diffbase.Diffbase(base, changed)

			Expect(reflect.DeepEqual(ret, expected)).Should(Equal(true))

		})
	})
})
