package tileproperties

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Navops/yaml"
)

type TileProperties struct {
	Properties map[string]interface{}
}

// NewTileProperties () returns a new tile properties interface
func NewTileProperties() *TileProperties {

	return &TileProperties{}
}

func (t *TileProperties) ReadJSON(filename string) (map[string]interface{}, error) {
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
	t.Properties = m
	return m, nil
}

func (t *TileProperties) ReadYAML(filename string) (map[string]interface{}, error) {
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
	t.Properties = m
	return m, nil
}
