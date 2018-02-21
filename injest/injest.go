package injest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/Navops/yaml"
	"github.com/xchapter7x/lo"
)

/// Injest required for go-flags
///-----------
type Injest struct {
	InputFile string `short:"i" long:"injest" description:"Filename to be injested" required:"true"`
	Simple    bool   `short:"s" long:"simple" description:"Simplify Keys"`
	Annotate  bool   `short:"a" long:"annotate" description:"Annotate output with help"`
}

/// go-flags callhack entry point
func (c *Injest) Execute([]string) error {

	base, _ := ReadJSON(c.InputFile)
	result, _ := ProcessInjest(base)
	// 	b, _ := yaml.Marshal(result)
	// fmt.Println(string(b))
	// fmt.Printf("%v %v", c.Simple, c.Annotate)
	// fmt.Println("----------")
	//result = forceSimpleKeys(result)
	OutputYaml(result, c.Simple, c.Annotate)

	return nil
}

func ReadJSON(filename string) (map[string]interface{}, error) {
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

func CreateCollection(m interface{}) (map[string]interface{}, error) {
	//fmt.Printf("Collection: %v %T\n\n\n", m, m)
	source, correct := m.(map[string]interface{})
	result := make(map[string]interface{})
	if correct {

		for k, v := range source {
			if k == "guid" {
				continue //ignore guids - we don't want them
			}
			//fmt.Printf("Element: %v %v\n\n\n", k, v.(map[string]interface{})["value"])
			result[k] = v.(map[string]interface{})["value"]

		}
	}

	return result, nil

}
func CreateCollections(m interface{}) ([]interface{}, error) {
	lo.G.Debugf("Collection: %v %T\n\n\n", m, m)
	source, correct := m.(map[string]interface{})
	var result []interface{}
	if correct {

		for k, v := range source {
			if k == "value" {
				// this is an array of maps

				value, correct := v.([]interface{})
				if correct {
					lo.G.Debugf("CA: %v %v\n\n\n", k, v)
					for element, value := range value {
						lo.G.Debugf("E: %v %v\n\n\n", element, value)
						e, _ := CreateCollection(value)
						result = append(result, e)

						lo.G.Debugf("YAY: %v %T\n", e, e)

						_ = element
					}
				}

			}

		}
	}
	return result, nil

}

func ProcessInjest(m interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	source, correct := m.(map[string]interface{})
	if correct {

		// Remove properties key and just look at keys under it.
		sub := source["properties"]
		if sub == nil {
			return nil, errors.New("properties was not found")
		}

		subkey := sub.(map[string]interface{})
		//		fmt.Println("subkey: ", subkey, "\n")

		for k, v := range subkey {
			property := v.(map[string]interface{}) //maintain type as long as possible
			switch property["type"] {
			case "integer":
				if property["configurable"] == true { // ignore unconfigurable elements (note creds)
					if property["value"] != nil {

						result[k] = int(property["value"].(float64)) // no idea why this is read in as a float64
					} else {
						result[k] = 0
					}
				}
			case "collection":
				result[k], _ = CreateCollections(v)

			//case "selector":
			case "simple_credentials":
				//fmt.Println("simple_credentials")
				if property["configurable"] == true { // ignore unconfigurable elements (note creds)
					result[k] = property["value"].(map[string]interface{})["password"]
				}

			case "string_list":
				fmt.Println("string_list not implemented yet")
			case "rsa_cert_credentials":
				fmt.Println("rsa_cert_credentials not implemented yet")
			case "secret":
				//fmt.Println("Secret")
				if property["configurable"] == true { // ignore unconfigurable elements (note creds)
					result[k] = property["value"].(map[string]interface{})["secret"]
				}

			default: // strings
				if property["configurable"] == true { // ignore unconfigurable elements (note creds)
					result[k] = property["value"]
				}
			}

			//fmt.Println("subkey: ", property)

		}
	} else {
		return nil, errors.New("parameter passed was incorrect")
	}
	return result, nil
}

func PrintYamlLine(key string, value interface{}, comment string) {

	//	fmt.Printf("%T %v\n", value, value)

	if len(comment) == 0 {

		switch value.(type) {
		case []interface{}:
			fmt.Printf("\n# %v\n%v:\n", comment, key)
			b, _ := yaml.Marshal(value)
			fmt.Println(string(b))
		default:
			fmt.Printf("%v: %v\n", key, value)
		}
		return
	}
	if comment[0] == '-' { // if first char is '-' put comment after yaml
		switch value.(type) {
		case []interface{}:
			fmt.Printf("\n# %v\n%v:\n", comment, key)
			b, _ := yaml.Marshal(value)
			fmt.Println(string(b))
			//fmt.Printf("MAP:%v: %v  #%v\n", key, value, comment)
		default:
			fmt.Printf("%v: %v  #%v\n", key, value, comment)
		}
	} else {
		switch value.(type) {
		case []interface{}:

			fmt.Printf("\n# %v\n%v:\n", comment, key)
			b, _ := yaml.Marshal(value)
			fmt.Println(string(b))
		default:
			fmt.Printf("\n# %v\n%v: %v\n", comment, key, value)
		}

	}
}

func OutputYaml(m map[string]interface{}, simple bool, annotate bool) {

	// Rip out keys into a slice so we can sort it.
	sortedKeys := make([]string, 0)
	for l, _ := range m {
		sortedKeys = append(sortedKeys, l)
	}
	sort.Strings(sortedKeys)

	// Open the help file
	help, err := readYaml("help.yml")
	if err != nil {
		fmt.Println(err)
	}

	// Open the dictonary file
	dict, err := readYaml("dictonary.yml")
	if err != nil {
		fmt.Println(err)
	}

	// iterate over the keys
	for _, k := range sortedKeys {
		//fmt.Printf("%v %T\n", m[k], m[k])
		if comment, ok := help[k]; ok && annotate { // Add comments
			if simpleKey, correct := dict[k].(string); correct && simple {
				// swap opsmgrKey for simplified key
				PrintYamlLine(simpleKey, m[k], comment.(string))
				// if comment.(string)[0] == '-' { // if first char is '-' put comment after yaml
				// 	fmt.Printf("%v: %v  #%v\n", simpleKey, m[k], comment)
				// } else {
				// 	fmt.Printf("\n# %v\n%v: %v\n", comment, simpleKey, m[k])
				// }
			} else {
				// Could be a collection, nil
				PrintYamlLine(k, m[k], comment.(string))
				/*if comment.(string)[0] == '-' { // if first char is '-' put comment after yaml
					fmt.Printf("%v: %v  #%s\n", k, m[k], comment)
				} else {
					fmt.Printf("\n# %v\n%v: %v\n", comment, k, m[k])
				}*/
			}
		} else {
			if simpleKey, correct := dict[k].(string); correct && simple {
				// swap opsmgrKey for simplified key
				//fmt.Printf("%v: %v\n", simpleKey, m[k])
				PrintYamlLine(simpleKey, m[k], "")
			} else {
				//fmt.Printf("%v: %v\n", k, m[k])
				PrintYamlLine(k, m[k], "")
			}
		}
	}
}

func OldProcessInjest(m interface{}) {
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

			// Supported
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
