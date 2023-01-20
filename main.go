package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func main() {
	data 	   := getJSON()
	table    := makeTabel()
	flatJSON := make(map[string]interface{})
	keys     := make([]string, 0)

	flatJSON, keys = flattenJSON(data, "  ", flatJSON, keys)
	for _, key := range keys {
		// fmt.Println(key, ":", flatJSON[key])
		// when null value, skip
		if _, ok := flatJSON[key]; ok {
			key = fmt.Sprintf("%v\t%v", key, flatJSON[key])
			fmt.Fprintln(table, key)
		}
	}
	// couse table.Flush isn't stable
	// os.Stdout.Sync()
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
			"state": "NY",
			"postcode": "23410"
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
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatal("Error in JSON, please check the struture.")
		return data
	}

	return data
}

func flattenJSON(data interface{}, parentKey string, flatJSON map[string]interface{}, keys []string) (map[string]interface{}, []string) {
	switch data.(type) {
	// if the data is a map of strings to interfaces
	case map[string]interface{}:
		// iterate through the map
		for key, value := range data.(map[string]interface{}) {
			// append the current key to the keys slice
			keys = append(keys, parentKey+key)
			// recursively call the flattenJSON function with 
			// the value, updated parentKey, flatJSON and keys
			flatJSON, keys = flattenJSON(value, parentKey+key+"_", flatJSON, keys)
		}
	// if the data is a slice of interfaces
	case []interface{}:
		// iterate through the slice
		for i, value := range data.([]interface{}) {
			// append the current index to the keys slice
			keys = append(keys, parentKey+strconv.Itoa(i))
			// recursively call the flattenJSON function with 
			// the value, updated parentKey, flatJSON and keys
			flatJSON, keys = flattenJSON(value, parentKey+strconv.Itoa(i)+"_", flatJSON, keys)
		}
	// if the data is neither a map nor 
	// a slice, it is a leaf node
	default:
		// add the key-value pair to the flatJSON map
		flatJSON[strings.TrimSuffix(parentKey, "_")] = data
	}
	// return the flatJSON and keys
	return flatJSON, keys
}


func makeTabel() *tabwriter.Writer {
	// see https://pkg.go.dev/text/tabwriter#NewWriter
	table := tabwriter.NewWriter(os.Stdout, 10, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Fprintln(table, "  Key\tValue")
	fmt.Fprintln(table, "  ---\t-----")

	return table
}
