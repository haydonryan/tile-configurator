package ingest

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/haydonryan/tile-configurator/dictionary"
	"github.com/haydonryan/tile-configurator/tileproperties"

	"github.com/Navops/yaml"
	"github.com/xchapter7x/lo"
)

/// Ingest required for go-flags
type Ingest struct {
	InputFile string `short:"i" long:"ingest" description:"Filename to be ingested" required:"true"`
	Simple    bool   `short:"s" long:"simple" description:"Simplify Keys"`
	Annotate  bool   `short:"a" long:"annotate" description:"Annotate output with help"`
}

/// go-flags callhack entry point
func (c *Ingest) Execute([]string) error {

	tile := tileproperties.NewTileProperties()

	base, _ := tile.ReadJSON(c.InputFile)

	//base, _ := ReadJSON(c.InputFile)
	result, _ := ProcessIngest(base)

	OutputYaml(result, c.Simple, c.Annotate)

	return nil
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

func ProcessIngest(m interface{}) (map[string]interface{}, error) {
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
	if comment[len(comment)-1] != '\n' {
		comment = comment + string('\n')

	}

	// Strip nil and *** from output
	v, correct := value.(string)
	if correct {

		if v == "***" {
			value = ""
		}
	}
	if value == nil {
		value = ""
	}

	if comment[0] == '-' { // if first char is '-' put comment after yaml
		switch value.(type) {
		case []interface{}:
			fmt.Printf("\n# %v\n%v:\n", comment, key)
			b, _ := yaml.Marshal(value)
			fmt.Println(string(b))
			//fmt.Printf("MAP:%v: %v  #%v\n", key, value, comment)
		default:
			//fmt.Printf("%v: %v  \t\t\t#%v\n", key, value, comment)

			tmp := fmt.Sprintf("%v: %v", key, value)
			fmt.Printf("%-70v#%v", tmp, comment)
		}

	} else if comment[0] == '^' { // if first char is '-' put comment after yaml
		switch value.(type) {
		case []interface{}:

			fmt.Printf("\n# %v\n%v:\n", comment, key)
			b, _ := yaml.Marshal(value)
			fmt.Println(string(b))
			//fmt.Printf("MAP:%v: %v  #%v\n", key, value, comment)
		default:
			split := strings.Split(comment, "%v")
			// check that it was actually split

			fmt.Printf("# %v\n", strings.Trim(split[0], "^"))
			tmp := fmt.Sprintf("%v: %v", key, value)
			fmt.Printf("%-70v#%v", tmp, split[1])

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
	help := dictionary.NewDictionary()
	err := help.LoadDictionary("help.yml")
	if err != nil {
		fmt.Println(err)
	}

	// Open the dictonary file
	dict := dictionary.NewDictionary()
	err = dict.LoadDictionary("dictionary.yml")
	if err != nil {
		fmt.Println(err)
	}

	// iterate over the keys
	for _, key := range sortedKeys {
		var comment = ""
		if annotate {
			comment = help.Lookup[key]
		}

		if simpleKey, correct := dict.Lookup[key]; correct && simple {
			// swap opsmgrKey for simplified key
			PrintYamlLine(simpleKey, m[key], comment)
		} else {

			PrintYamlLine(key, m[key], comment)
		}

	}
}
