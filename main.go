package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func main() {
	// couse table.Flush isn't stable
	defer os.Stdout.Sync() 

	var data interface{} = getJSON()
	table    := makeTabel()
	flatJSON := make(map[string]interface{})
	keys     := make([]string, 0)

	flatJSON, keys = flattenJSON(data, "", flatJSON, keys)
	for _, key := range keys {
		// fmt.Println(key, ":", flatJSON[key])
		if flatJSON[key] == nil {continue}
		key = fmt.Sprintf("%v\t%v", key, flatJSON[key])
		fmt.Fprintln(table, key)
	}
	table.Flush()
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

func flattenJSON(data interface{}, parentKey string, flatJSON map[string]interface{}, keys []string) (map[string]interface{}, []string) {
	switch data.(type) {
	case map[string]interface{}:
		for key, value := range data.(map[string]interface{}) {
			keys = append(keys, parentKey+key)
			flatJSON, keys = flattenJSON(value, parentKey+key+"_", flatJSON, keys)
		}
	case []interface{}:
		for i, value := range data.([]interface{}) {
			keys = append(keys, parentKey+strconv.Itoa(i))
			flatJSON, keys = flattenJSON(value, parentKey+strconv.Itoa(i)+"_", flatJSON, keys)
		}
	default:
		flatJSON[strings.TrimSuffix(parentKey, "_")] = data
	}
	return flatJSON, keys
}

func makeTabel() *tabwriter.Writer {
	table := tabwriter.NewWriter(os.Stdout, 10, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Fprintln(table, "Key\tValue")
	fmt.Fprintln(table, "---\t-----")

	return table
}
