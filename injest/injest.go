package injest

import (
	"fmt"
	"sort"
)

/// Injest required for go-flags
///-----------
type Injest struct {
}

/// go-flags callhack entry point
func (c *Injest) Execute([]string) error {
	fmt.Println("in execute")
	return nil
}

var help = map[string]string{

	".properties.syslog":              "Syslog Forwarding\n# -----------------\n# Choose either enabled or disbled",
	".properties.syslog.enabled.port": "The port for syslog",
	/*".properties.backups":             "Backups\n# -----------------\n# (enabled | disbled)",

	".properties.backups.azure.base_url": "Azure URL for XXX",
	//".properties.backups.azure.container":       "<nil>",
	//".properties.backups.azure.container_path":  "<nil>",
	/*	".properties.backups.azure.storage_account": "Storage Account Name in XX Format",

		".properties.backups.enable.access_key_id": "<nil>",
		".properties.backups.enable.bucket_name":   "<nil>",
		".properties.backups.enable.bucket_path":   "<nil>",
		".properties.backups.enable.endpoint_url":  "<nil>",
		".properties.backups.enable.region":        "<nil>",

		".properties.backups.gcs.bucket_name": "<nil>",
	".properties.backups.gcs.project_id": "Choose the GCS Project ID to back up your MySQL to:",*/
}

func ProcessInjest(m interface{}) {
	var sl map[string]interface{}
	sl, correct := m.(map[string]interface{})

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
				fmt.Println(k)
				delete(sl["properties"].(map[string]interface{}), k)

				//currently these types are not supported - delete the keys
			case "salted_credentials":
				delete(sl["properties"].(map[string]interface{}), k)
			case "secret":
				delete(sl["properties"].(map[string]interface{}), k)
			case "ip_ranges":
				delete(sl["properties"].(map[string]interface{}), k)
			case "network_address":
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
