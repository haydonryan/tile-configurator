//------------------------------------------------------------------
//
//  Diffbase returns the additions that map has to a base map
//

package diffbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/xchapter7x/lo"
)

type Diffbase struct {
	BaseFile       string `short:"b" long:"base" description:"Filename for base (unconfigured tile)" required:"true"`
	ConfiguredFile string `short:"c" long:"configured" description:"Filename to configured tile" required:"true"`
}

func readJSON(filename string) (map[string]interface{}, error) {
	// Open the properties file
	base, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Could not open File")
	}

	// Read file into map of interfaces
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(base), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, fmt.Errorf("Could not unmartshall File")
	}

	d := RemoveMainProperties(m)

	return d.(map[string]interface{}), nil
}

func (c *Diffbase) Execute([]string) error {

	// open configured file
	base, _ := readJSON(c.BaseFile)
	configured, _ := readJSON(c.ConfiguredFile)

	lo.G.Debug(base)

	lo.G.Debug(configured)

	answer := Diff(base, configured)
	lo.G.Debug("\n\nDiff between the files is:")
	lo.G.Debug(answer)
	jsonString, _ := json.Marshal(answer)

	fmt.Println(string(jsonString))

	return nil
}

func SquishCollection(m interface{}) (map[string]interface{}, error) {
	//fmt.Printf("Collection: %v %T\n\n\n", m, m)
	source, correct := m.(map[string]interface{})
	result := make(map[string]interface{})
	if correct {

		for k, v := range source {
			result[k] = v.(map[string]interface{})["value"]
		}
	}

	return result, nil

}

// SquishCollections changes a collection from the way that Ops Manager returns
// it's api to how it expects it to be set
func SquishCollections(collection []interface{}) []interface{} {
	var result []interface{}
	for _, value := range collection {
		temp, _ := SquishCollection(value)
		result = append(result, temp)

	}

	return result
}

// FilterGUIDs filters out the built
// it's api to how it expects it to be set
func RemoveNonConfigurable(collection interface{}) interface{} {

	return nil
}

func RemoveMainProperties(properties interface{}) interface{} {

	p := properties.(map[string]interface{})
	//	fmt.Printf("%T\n\n", p)
	if p["properties"] != nil {
		return p["properties"]
		//return nil
	}
	return properties
}

//Diffbase returns the additions that map has to a base map
func Diff(base interface{}, changed interface{}) interface{} {
	//	difference := make(map[string]interface{})
	lo.G.Debug("Comparing:")
	lo.G.Debugf("Value1: %v Type %T\n", base, base)
	lo.G.Debugf("Value2: %v Type %T\n", changed, changed)

	if base == nil {
		lo.G.Debug("base is nil")
		return changed
	}

	if changed == nil {
		lo.G.Debug("changed is nil")
		return nil
	}

	if reflect.TypeOf(base) != reflect.TypeOf(changed) {
		return changed
	}

	switch changed.(type) {
	case map[string]interface{}:
		if base == nil {
			return changed
		} else {
			changedMap := changed.(map[string]interface{})
			baseMap := base.(map[string]interface{})
			difference := make(map[string]interface{})

			propertyType := changedMap["type"]
			if propertyType == "collection" {

				//fmt.Printf("Found\n\n\n\n\n")
				changedMap["value"] = SquishCollections(changedMap["value"].([]interface{}))
				return changedMap
			}

			if propertyType == "multi_select_options" {
				return changedMap["value"]
			}

			for key, value := range changed.(map[string]interface{}) {

				diff := Diff(baseMap[key], changedMap[key])

				if diff != nil {
					difference[key] = diff
				}

				_ = value

			}
			lo.G.Debugf("Map returning: %v Type %T\n\n", difference, difference)
			if len(difference) == 0 {
				return nil
			}

			return difference

		}
	case []interface{}:
		if base == nil {
			return changed
		} else {
			changedMap := changed.([]interface{})
			baseMap := base.([]interface{})
			var difference []interface{}

			// need to iterate through both sides to check if the changed side has new keys comparedd to all
			// of the base map.

			for changedKey, value := range changed.([]interface{}) {
				var diff interface{}
				for baseKey, _ := range base.([]interface{}) {
					diff = Diff(baseMap[baseKey], changedMap[changedKey])
					if diff == nil {
						break
					}
				}
				if diff != nil {
					difference = append(difference, diff)

				}

				_ = value
			}
			lo.G.Debugf("Array returning: %v Type %T\n\n", difference, difference)
			if len(difference) == 0 {
				return nil
			}
			return difference
		}

	default:
		if base != changed {
			lo.G.Debugf("Value returning: %v Type %T\n\n", changed, changed)
			return changed
		}

	}
	var m interface{}

	lo.G.Debugf("returning: %v Type %T\n\n", m, m)
	return m

}
