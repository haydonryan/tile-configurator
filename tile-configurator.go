package main

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	flags "github.com/jessevdk/go-flags"
	"github.com/xchapter7x/lo"

	yaml "github.com/Navops/yaml" // not using go-yaml directly as it defaults to map[interface{}]interface{}
)

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
func processAllCollection(m interface{}, url string, user string, password string, tile string) {
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
					runCommand(url, user, password, tile, makeJSON(result))
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

func processGroup(m interface{}, url string, user string, password string, tile string) {
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
			runCommand(url, user, password, tile, makeJSON(iface))
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

func runCommand(url string, user string, password string, tile string, properties string) {
	var output []byte

	cmd := "om"
	var args []string

	if GlobalOptions.useKeyAndSecret {
		fmt.Printf("-----------(%v %v", user, password)
		args = []string{"-t", url, "-k", "-c", user, "-s", password, "configure-product", "--product-name", tile, "--product-properties", properties}
	} else {
		args = []string{"-t", url, "-k", "-u", user, "-p", password, "configure-product", "--product-name", tile, "--product-properties", properties}
	}

	lo.G.Debug(args)

	cmdrnuner := exec.Command(cmd, args...)

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

func main() {

	fmt.Println("\ntileConfigurator v0.0.4")
	fmt.Println("--------")

	var opts struct {
		Verbose  []bool `short:"v" long:"verbose" description:"Show verbose debug info"`
		Filename string `short:"i" long:"input" description:"filename to apply to Ops manager" required:"true"`
		URL      string `long:"url" description:"URL of the Ops manager" required:"true"`
		Tile     string `short:"t" long:"tile" description:"filename to apply to Ops manager" required:"true"`

		//User struct {
		Username string `short:"u" long:"user" description:"Username for logging into Ops Manager."`
		Key      string `short:"k" long:"key" description:"key for logging into Ops Manager"`
		//} `description:"Credentials" required:"1"`
		//Password struct {
		Password string `short:"p" long:"password" description:"Password for logging into Ops Manager"`
		Secret   string `short:"s" long:"secret" description:"Secret for logging into Ops Manager"`
		//} `description:"Secrets" required:"1"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		lo.G.Debug(err)
		os.Exit(1)
	}

	fmt.Println(string(opts.Key))

	if string(opts.Key) != "" {
		GlobalOptions.useKeyAndSecret = true
	}
	// Open the properties file
	yaml, err := readYaml(string(opts.Filename))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse through the map and process keys as either individual, collections and groups
	for key, value := range yaml {
		if key == "collections" {
			lo.G.Debug("Found collections")
			if GlobalOptions.useKeyAndSecret {
				processAllCollection(value, string(opts.URL), string(opts.Key), string(opts.Secret), string(opts.Tile))
			} else {
				processAllCollection(value, string(opts.URL), string(opts.Username), string(opts.Password), string(opts.Tile))
			}
		} else if key == "groups" {
			lo.G.Debug("Found groups block")
			if GlobalOptions.useKeyAndSecret {
				processGroup(value, string(opts.URL), string(opts.Key), string(opts.Secret), string(opts.Tile))
			} else {
				processGroup(value, string(opts.URL), string(opts.Username), string(opts.Password), string(opts.Tile))
			}

		} else {
			properties := fmt.Sprintf("{\"%v\": {\"value\":  %v}}\n", key, value)
			fmt.Printf("Applying setting %v to tile %v......", key, string(opts.Tile))
			//runCommand(value, string(opts.URL), string(opts.Username), string(opts.Password), string(opts.Tile))
			if GlobalOptions.useKeyAndSecret {
				runCommand(string(opts.URL), string(opts.Key), string(opts.Secret), string(opts.Tile), properties)

			} else {
				runCommand(string(opts.URL), string(opts.Username), string(opts.Password), string(opts.Tile), properties)
			}
			fmt.Printf("Done.\n")
		}
	}
}
