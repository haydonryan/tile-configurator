package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/Navops/yaml"
	"github.com/xchapter7x/lo"
)

type Deploy struct {
	Filename string `short:"i" long:"input" description:"Filename to apply to Ops manager" required:"true"`
	URL      string `long:"url" description:"URL of the Ops manager" required:"true"`
	Tile     string `short:"t" long:"tile" description:"filename to apply to Ops manager" required:"true"`

	Username string `short:"u" long:"user" description:"Username for logging into Ops Manager."`
	Key      string `short:"k" long:"key" description:"key for logging into Ops Manager"`

	Password string `short:"p" long:"password" description:"Password for logging into Ops Manager"`
	Secret   string `short:"s" long:"secret" description:"Secret for logging into Ops Manager"`
	OM       string `short:"o" long:"om" description:"Name of OM on this system" default:"om-linux"`
}

func (c *Deploy) Execute([]string) error {

	if string(c.Key) != "" {
		GlobalOptions.useKeyAndSecret = true
	}
	// Open the properties file
	yaml, err := readYaml(string(c.Filename))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Parse through the map and process keys as either individual, collections and groups
	for key, value := range yaml {
		if key == "collections" {
			lo.G.Debug("Found collections")
			if GlobalOptions.useKeyAndSecret {
				c.processAllCollection(value, string(c.URL), string(c.Key), string(c.Secret), string(c.Tile))
			} else {
				c.processAllCollection(value, string(c.URL), string(c.Username), string(c.Password), string(c.Tile))
			}
		} else if key == "groups" {
			lo.G.Debug("Found groups block")
			if GlobalOptions.useKeyAndSecret {
				c.processGroup(value, string(c.URL), string(c.Key), string(c.Secret), string(c.Tile))
			} else {
				c.processGroup(value, string(c.URL), string(c.Username), string(c.Password), string(c.Tile))
			}

		} else {
			properties := fmt.Sprintf("{\"%v\": {\"value\":  \"%v\"}}\n", key, value)
			fmt.Printf("Applying setting %v to tile %v......", key, string(c.Tile))
			//runCommand(value, string(opts.URL), string(opts.Username), string(opts.Password), string(opts.Tile))
			if GlobalOptions.useKeyAndSecret {
				c.runCommand(string(c.URL), string(c.Key), string(c.Secret), string(c.Tile), properties)

			} else {
				c.runCommand(string(c.URL), string(c.Username), string(c.Password), string(c.Tile), properties)
			}
			fmt.Printf("Done.\n")
		}
	}

	return nil
}

func makeJSON(v interface{}) string {
	var enc []byte
	enc, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	return string(enc)
}

var GlobalOptions struct {
	useKeyAndSecret bool
}

//each collection block can have multiple collections.
func (c *Deploy) processAllCollection(m interface{}, url string, user string, password string, tile string) {
	var groupSlice []interface{}
	groupSlice, correct := m.([]interface{})

	if correct {
		for _, collection := range groupSlice {
			col, correct := collection.(map[string]interface{})
			//fmt.Printf("key: %v", key) //key not neeed as it's '0'
			if correct {
				for i, _ := range col {
					result := make(map[string]interface{})
					result[i] = make(map[string]interface{})
					result[i].(map[string]interface{})["value"] = col[i]

					lo.G.Debug("Call OM tool with: ")
					lo.G.Debug("[-", makeJSON(result), "-]")

					fmt.Printf("Applying setting %v to tile %v......", i, tile)
					c.runCommand(url, user, password, tile, makeJSON(result))
					fmt.Printf("Done.\n")
				}
			} else {
				fmt.Printf("errored, %T", m)
			}
		}
	} else {
		fmt.Printf("CollectionsAll: %v %T\n\n", m, m)
		fmt.Println("error")
	}
}

func (c *Deploy) processGroup(m interface{}, url string, user string, password string, tile string) {
	//Groups are defined as an array of maps.  Name is within the map
	var groupSlice []interface{}
	groupSlice, correct := m.([]interface{})

	if correct {
		for _, group := range groupSlice { //for each array process the map in a single group
			var key string
			var iface map[string]interface{}
			iface, correct = group.(map[string]interface{})
			if correct {
				// this map contains the properties that need to be run together.
				key = iface["name"].(string)
				delete(iface, "name")
				// For each of the properties add a value sub key.
				for i, _ := range iface {
					result := make(map[string]interface{})
					result["value"] = iface[i]
					iface[i] = result

				}
			} else {
				fmt.Printf("incorrect %v %T", iface, iface)
			}

			lo.G.Debug("Call OM tool with: ")
			lo.G.Debug("[-", makeJSON(iface), "-]")

			fmt.Printf("Applying Group %v to tile %v.....", key, tile)
			c.runCommand(url, user, password, tile, makeJSON(iface))
			fmt.Printf("Done.\n")
		}

	} else {
		fmt.Printf("error: %v %T \n\n", m, m)
	}
}

func processHash(m interface{}) {
	fmt.Println("Error: Hashes not supported")
	os.Exit(1)

}

func (c *Deploy) runCommand(url string, user string, password string, tile string, properties string) {
	var output []byte

	var args []string

	if GlobalOptions.useKeyAndSecret {
		fmt.Printf("-----------(%v %v", user, password)
		args = []string{"-t", url, "-k", "-c", user, "-s", password, "configure-product", "--product-name", tile, "--product-properties", properties}
	} else {
		args = []string{"-t", url, "-k", "-u", user, "-p", password, "configure-product", "--product-name", tile, "--product-properties", properties}
	}

	lo.G.Debug(args)

	cmdrnuner := exec.Command(c.OM, args...)

	if output, err := cmdrnuner.CombinedOutput(); err != nil {
		fmt.Println("Call to OM returned error:")
		fmt.Fprintln(os.Stderr, err)
		fmt.Println(string(output))
		os.Exit(1)
	} else {
		lo.G.Debug("success")
		lo.G.Debug(string(output))
	}
	_ = output
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
