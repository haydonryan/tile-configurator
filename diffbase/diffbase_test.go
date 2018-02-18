package diffbase_test

import (
	"reflect"

	"github.com/haydonryan/tile-configurator/diffbase"
	//. "github.com/haydonryan/tile-configurator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diffbase", func() {

	Context("Diff()", func() {
		It("should return nil if both maps are empty", func() {
			base := make(map[string]interface{})

			ret := diffbase.Diff(base, base)
			Expect(ret).Should(BeNil())
		})

		It("should return if new map has extra key", func() {
			base := make(map[string]interface{})
			changed := make(map[string]interface{})

			changed["test"] = "test string"

			ret := diffbase.Diff(base, changed)
			Expect(ret).Should(Equal(changed))
		})

		It("should return all keys if new map has multiple keys", func() {
			base := make(map[string]interface{})
			changed := make(map[string]interface{})

			changed["test"] = "test string"
			changed["test2"] = "test string2"

			ret := diffbase.Diff(base, changed)
			Expect(ret).Should(Equal(changed))
		})
		Context("Subkeys", func() {

			It("should return the difference to include a subkey of the main key", func() {
				base := make(map[string]interface{})
				changed := make(map[string]interface{})
				sub := make(map[string]interface{})

				sub["subkey"] = "test string"
				changed["test2"] = sub

				ret := diffbase.Diff(base, changed)
				Expect(ret).Should(Equal(changed))
			})

			It("should return only extra subkey if both maps have subkeys", func() {
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

				ret := diffbase.Diff(base, changed)

				Expect(reflect.DeepEqual(ret, expected)).Should(Equal(true))

			})
		})
		Context("Arrays", func() {

			It("should return nil if both arrays are the same", func() {

				// var names []interface{}
				// names = append(names, "first")
				// names = append(names, "second")

				// setup base
				var names []interface{}
				names = append(names, "first")
				names = append(names, "second")
				base := make(map[string]interface{})
				base["array"] = names

				// setup output to be different to base
				var names2 []interface{}
				names2 = append(names2, "first")
				names2 = append(names2, "second")
				changed := make(map[string]interface{})
				changed["array"] = names2

				// setup expected output
				//expected := make(map[string]interface{})
				// sub3 := make(map[string]interface{})
				// sub3["difkey"] = "diff string"
				// expected["test"] = sub3

				ret := diffbase.Diff(base, changed)

				//Expect(reflect.DeepEqual(ret, expected)).Should(Equal(true))
				Expect(ret).Should(BeNil())

			})

			It("should return an array element if arrays have different values", func() {

				// var names []interface{}
				// names = append(names, "first")
				// names = append(names, "second")

				// setup base
				var names []interface{}
				names = append(names, "first")
				names = append(names, "second")
				base := make(map[string]interface{})
				base["array"] = names

				// setup output to be different to base
				var names2 []interface{}
				names2 = append(names2, "first")
				names2 = append(names2, "different")
				changed := make(map[string]interface{})
				changed["array"] = names2

				// setup expected output
				expected := make(map[string]interface{})
				var names3 []interface{}
				names3 = append(names3, "different")
				expected["array"] = names3

				ret := diffbase.Diff(base, changed)

				// fmt.Printf("base: %v\n", base)
				// fmt.Printf("configured: %v\n", changed)
				// fmt.Printf("expected: %v\n", expected)
				// fmt.Printf("returned: %v\n", ret)
				Expect(reflect.DeepEqual(ret, expected)).Should(Equal(true))

			})
		})
	})
})
