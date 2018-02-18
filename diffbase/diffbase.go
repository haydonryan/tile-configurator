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
func Diffbase(base map[string]interface{}, changed map[string]interface{}) map[string]interface{} {
	difference := make(map[string]interface{})

	for key, value := range changed {

		switch value.(type) {
		case map[string]interface{}:
			// base[key] could be nil here, if it is then just return whole tree
			if base[key] == nil {
				difference[key] = changed[key]
			} else {

				ret := Diffbase(base[key].(map[string]interface{}), changed[key].(map[string]interface{}))
				if len(ret) != 0 {
					difference[key] = ret
				}
			}
		default:
			if base[key] != changed[key] {
				difference[key] = value
			}
		}

	}
	return difference
}
