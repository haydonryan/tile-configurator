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
	fmt.Println("Diffing two files")

	// open configured file
	base, _ := readJSON(c.BaseFile)
	configured, _ := readJSON(c.ConfiguredFile)

	lo.G.Debug(base)

	lo.G.Debug(configured)

	answer := Diff(base, configured)
	fmt.Println("\n\nDiff between the files is:")
	//fmt.Println(answer)
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
	// compare if types are the same
	// not done yet.

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
			//difference := make([]interface{})

			for key, value := range changed.([]interface{}) {
				diff := Diff(baseMap[key], changedMap[key])

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
