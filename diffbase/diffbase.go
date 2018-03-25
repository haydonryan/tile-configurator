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
	return m, nil
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
