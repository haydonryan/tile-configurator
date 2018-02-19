package injest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/Navops/yaml"
)

/// Injest required for go-flags
///-----------
type Injest struct {
	InputFile string `short:"i" long:"injest" description:"Filename to be injested" required:"true"`
}

/// go-flags callhack entry point
func (c *Injest) Execute([]string) error {

	base, _ := readJSON(c.InputFile)
	ProcessInjest(base)

	return nil
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

func readYaml(filename string) (map[string]interface{}, error) {
	// Open the properties file
	base, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Could not open File")
	}

	// Read file into map of interfaces
	m := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(base), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, fmt.Errorf("Could not unmartshall File")
	}
	return m, nil
}

var help = map[string]string{}

func ProcessInjest(m interface{}) {
	var sl map[string]interface{}
	sl, correct := m.(map[string]interface{})

	// Open the help file
	help, err := readYaml("help.yml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(help)
	// unwrap properties
	if correct {

		keys := make([]string, 0)
		for l, _ := range sl["properties"].(map[string]interface{}) {
			keys = append(keys, l)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := sl["properties"].(map[string]interface{})[k]
			//for k, v := range sl["properties"].(map[string]interface{}) {

			if v.(map[string]interface{})["configurable"] == false { //remove anything that isn't configurable
				//fmt.Println("removing: ", k)
				continue
			}

			switch v.(map[string]interface{})["type"] {
			case "simple_credentials":
				//fmt.Println(k)
				fmt.Println("< Simple_credentials not supported >")
				delete(sl["properties"].(map[string]interface{}), k)

				//currently these types are not supported - delete the keys
			case "salted_credentials":
				fmt.Println("< salted_credentials not supported >")
				delete(sl["properties"].(map[string]interface{}), k)
			case "secret":
				fmt.Println("< secret not supported >")
				delete(sl["properties"].(map[string]interface{}), k)
			case "ip_ranges":
				fmt.Println("< ip_ranges not supported >")
				delete(sl["properties"].(map[string]interface{}), k)
			case "network_address":
				fmt.Println("< network_address not supported >")
				delete(sl["properties"].(map[string]interface{}), k)
				//case "dropdown_select":
				//fmt.Println("< dropdown_select not nsupported >")
				delete(sl["properties"].(map[string]interface{}), k)
			case "text":
				fallthrough
			case "email": //email are just text fields but have regex checking
				fallthrough
			case "selector": //selectors have
				fallthrough
			case "integer":
				fallthrough
			case "string":
				fallthrough
			case "boolean":
				// See if there is comments that are available for this key
				if val, ok := help[k]; ok {
					fmt.Printf("\n# %v\n%v: %v\n", val, k, v.(map[string]interface{})["value"])
				} else {
					fmt.Printf("%v: %v\n", k, v.(map[string]interface{})["value"])
				}

				delete(v.(map[string]interface{}), "optional")
				delete(v.(map[string]interface{}), "type")

			case "collection":
				fmt.Printf("collection: ")
				delete(v.(map[string]interface{}), "optional")
				delete(v.(map[string]interface{}), "type")
				fmt.Printf("%v, %v\n", k, v)
				fmt.Println()

			default:
				propertyType := v.(map[string]interface{})["type"]
				fmt.Printf("Default: %v\n", propertyType)
				fmt.Printf("%v, %v\n", k, v)
				fmt.Println()
			}

		}
	}

}
