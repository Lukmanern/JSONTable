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

type FlattenedJSON map[string]interface{}

func main() {
	rawJSON   := parseJSON()
	flatJSON  := make(FlattenedJSON, 0)
	jsonKeys  := make([]string, 0)
	flatJSON, jsonKeys = flattenMap(rawJSON, "  ", flatJSON, jsonKeys)

	showTable(rawJSON, jsonKeys, flatJSON)
}

func parseJSON() interface{} {
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

func flattenMap(data interface{}, parentKey string, flattenJSON FlattenedJSON, keys []string) (FlattenedJSON, []string) {
	switch data.(type) {
	// if the data is a map of strings to interfaces
	case map[string]interface{}:
		// iterate through the map
		for key, value := range data.(map[string]interface{}) {
			// append the current key to the keys slice
			keys = append(keys, parentKey+key)
			// recursively call the flattenJSON function with 
			// the value, updated parentKey, flatJSON and keys
			flattenJSON, keys = flattenMap(value, parentKey+key+"_", flattenJSON, keys)
		}
	// if the data is a slice of interfaces
	case []interface{}:
		// iterate through the slice
	
		for i, value := range data.([]interface{}) {
			// append the current index to the keys slice
			keys = append(keys, parentKey+strconv.Itoa(i))
			// recursively call the flattenJSON function with 
			// the value, updated parentKey, flatJSON and keys
			flattenJSON, keys = flattenMap(value, parentKey+strconv.Itoa(i)+"_", flattenJSON, keys)
		}
	// if the data is neither a map nor 
	// a slice, it is a leaf node
	default:
		// add the key-value pair to the flattenJSON map
		flattenJSON[strings.TrimSuffix(parentKey, "_")] = data
	}
	// return the flattenJSON and keys
	return flattenJSON, keys
}

func makeTable() *tabwriter.Writer {
	// see https://pkg.go.dev/text/tabwriter#NewWriter
	table := tabwriter.NewWriter(os.Stdout, 10, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Fprintln(table, "  Key\tValue")
	fmt.Fprintln(table, "  ---\t-----")

	return table
}

func showTable(rawJSON interface{}, jsonKeys []string, flatJSON FlattenedJSON) {
	tabWriter := makeTable()

	for _, key := range jsonKeys {
		// fmt.Println(key, ":", flatJSON[key])
		// when null value, skip
		if _, ok := flatJSON[key]; ok {
			key = fmt.Sprintf("%v\t%v", key, flatJSON[key])
			fmt.Fprintln(tabWriter, key)
		}
	}
	// couse table.Flush isn't stable
	os.Stdout.Sync()
	tabWriter.Flush()
}