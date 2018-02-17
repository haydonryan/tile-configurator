//------------------------------------------------------------------
//
//  Diffbase returns the additions that map has to a base map
//

package diffbase

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
