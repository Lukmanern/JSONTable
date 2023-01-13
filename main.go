// package main

// import (
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"text/tabwriter"
// )

// func main() {
// 	m := map[string]int{
// 		"item1": 12,
// 		"item2": 23,
// 		"item3": 34,
// 		"item4": 45,
// 	}

// 	w := tabwriter.NewWriter(os.Stdout, 0, 20, 20, '\t', tabwriter.AlignRight)
// 	fmt.Fprintln(w, "Key\tValue")
// 	fmt.Fprintln(w, "---\t-----")

// 	for key, val := range m {
// 		fmt.Fprintln(w, key+"\t"+strconv.Itoa(val))
// 	}
// 	w.Flush()
// }

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

// 	w := tabwriter.NewWriter(os.Stdout, 0, 20, 20, '\t', tabwriter.AlignRight)
// 	fmt.Fprintln(w, "Key\tValue")
// 	fmt.Fprintln(w, "---\t-----")

// 	for key, val := range m {
// 		fmt.Fprintln(w, key+"\t"+strconv.Itoa(val))
// 	}
// 	w.Flush()

func main() {
	json := getJSON()
	result := make(map[string]interface{})
	flattenJSON(json, "", result)
	sortMapByKey(result)
	
	table := tabwriter.NewWriter(os.Stdout, 10, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Fprintln(table, "Key\tValue")
	fmt.Fprintln(table, "---\t-----")

	for key, value := range result {
		// fmt.Println(key, ":", value)
		key = fmt.Sprintf("%v\t%v", key, value)
		fmt.Fprintln(table, key)
	}
	table.Flush()
}

func flattenJSON(data interface{}, parentKey string, flatJSON map[string]interface{}) {
	switch data.(type) {
	case map[string]interface{}:
		for key, value := range data.(map[string]interface{}) {
			flattenJSON(value, parentKey+key+"_", flatJSON)
		}
	case []interface{}:
		for i, value := range data.([]interface{}) {
			flattenJSON(value, parentKey+strconv.Itoa(i)+"_", flatJSON)
		}
	default:
		flatJSON[strings.TrimSuffix(parentKey, "_")] = data
	}
}

func sortMapByKey(m map[string]interface{}) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedMap := make(map[string]interface{}, len(m))
	for _, k := range keys {
		sortedMap[k] = m[k]
	}
	for key, val := range sortedMap {
		m[key] = val
	}
}

func getJSON() interface{} {
	jsonData := `
	{
		"name": "John Smith",
		"age": 35,
		"address": {
			"street": "Main Street",
			"city": "New York",
			"state": "NY"
		},
		"phoneNumbers": [
			{
				"type": "home",
				"number": "212 555-1234"
			},
			{
				"type": "work",
				"number": "646 555-4567"
			}
		]
	}`

	var data interface{}
	json.Unmarshal([]byte(jsonData), &data)

	return data
}

