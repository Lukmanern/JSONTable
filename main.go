package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
)

type FlattenedJSON map[string]interface{}

type FlatteningArgs struct {
	FlattenedJSON FlattenedJSON
	Keys          []string
	Prefix        string
}

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

func flattenMap(data interface{}, parentKey string, flatJSON FlattenedJSON, keys []string) (FlattenedJSON, []string) {
	// type of data-var
	dataType := reflect.TypeOf(data).Kind()
	if dataType == reflect.Map {
		// iterate through the map
		for key, value := range data.(map[string]interface{}) {
			// append the current key to the keys slice
			keys = append(keys, parentKey+key)
			// recursively call the flatJSON function with 
			// the value, updated parentKey, flatJSON and keys
			flatJSON, keys = flattenMap(value, parentKey+key+"_", flatJSON, keys)
		}
	} else if dataType == reflect.Slice {
		// iterate through the slice
		for i, value := range data.([]interface{}) {
			// append the current index to the keys slice
			keys = append(keys, parentKey+strconv.Itoa(i))
			// recursively call the flatJSON function with 
			// the value, updated parentKey, flatJSON and keys
			flatJSON, keys = flattenMap(value, parentKey+strconv.Itoa(i)+"_", flatJSON, keys)
		}
	} else {
		// if the data is neither a map nor 
		// a slice, it is a leaf node
		// add the key-value pair to the flatJSON map
		flatJSON[strings.TrimSuffix(parentKey, "_")] = data
	}

	// return the flatJSON and keys
	return flatJSON, keys
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