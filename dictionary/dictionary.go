package dictionary

import (
	"fmt"
	"io/ioutil"

	"github.com/Navops/yaml"
)

// Dictionary stores two way conversions between simple key names and opsmanager ones
type Dictionary struct {
	Lookup        map[string]string
	ReverseLookup map[string]string
}

// NewDictionary returns a new dictionary literal
func NewDictionary() *Dictionary {
	return &Dictionary{}
}

// LoadDictionary (filename string): Loads a dictionary as a yaml key pair - should not have sub elements
func (dict *Dictionary) LoadDictionary(filename string) error {
	// Open the properties file
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Could not open File")
	}

	// Read file into map
	err = yaml.Unmarshal([]byte(fileContents), &dict.Lookup)
	if err != nil {
		return fmt.Errorf("Could not unmartshall File")
	}

	dict.ReverseLookup = dict.MapReverse(dict.Lookup)
	return nil
}

// MapReverse (map[string]string) swaps value and key in a map
func (dict *Dictionary) MapReverse(m map[string]string) map[string]string {

	result := make(map[string]string)
	for k, v := range m {
		result[v] = k
	}

	return result
}
