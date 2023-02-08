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
	flatJSON 	FlattenedJSON
	Keys        []string
	parentKey   string
}

func main() {
	mainArgs := FlatteningArgs{
		flatJSON: make(FlattenedJSON, 0),
		Keys: make([]string, 0),
		parentKey: "  ",
	}

	flatJSON, jsonKeys := flattenMap(parseJSON(), mainArgs)

	showTable(jsonKeys, flatJSON)
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

func flattenMap(data interface{}, mainArgs FlatteningArgs) (FlattenedJSON, []string) {
	// type of data-var
	dataType := reflect.TypeOf(data).Kind()
	if dataType == reflect.Map {
		// iterate through the map
		for key, value := range data.(map[string]interface{}) {
			// append the current key to the keys slice
			mainArgs.Keys = append(mainArgs.Keys, mainArgs.parentKey+key)
			// make new args
			newMainArgs := FlatteningArgs{
				flatJSON: mainArgs.flatJSON,
				Keys: mainArgs.Keys,
				parentKey: mainArgs.parentKey+key+"_",
			}
			// recursively call the mainArgs.flatJSON function with 
			// the value, updated mainArgs.parentKey, mainArgs.flatJSON and mainArgs.Keys
			mainArgs.flatJSON, mainArgs.Keys = flattenMap(value, newMainArgs)
		}
	} else if dataType == reflect.Slice {
		// iterate through the slice
		for i, value := range data.([]interface{}) {
			// append the current index to the mainArgs.Keys slice
			mainArgs.Keys = append(mainArgs.Keys, mainArgs.parentKey+strconv.Itoa(i))
			// make new args
			newMainArgs := FlatteningArgs{
				flatJSON: mainArgs.flatJSON,
				Keys: mainArgs.Keys,
				parentKey: mainArgs.parentKey+strconv.Itoa(i)+"_",
			}
			// recursively call the mainArgs.flatJSON function with 
			// the value, updated mainArgs.parentKey, mainArgs.flatJSON and mainArgs.Keys
			mainArgs.flatJSON, mainArgs.Keys = flattenMap(value, newMainArgs)
		}
	} else {
		// if the data is neither a map nor 
		// a slice, it is a leaf node
		// add the key-value pair to the mainArgs.flatJSON map
		mainArgs.flatJSON[strings.TrimSuffix(mainArgs.parentKey, "_")] = data
	}

	// return the mainArgs.flatJSON and mainArgs.Keys
	return mainArgs.flatJSON, mainArgs.Keys
}

// table structure
func makeTable() *tabwriter.Writer {
	// see https://pkg.go.dev/text/tabwriter#NewWriter
	table := tabwriter.NewWriter(os.Stdout, 10, 8, 2, '\t', tabwriter.AlignRight)
	fmt.Fprintln(table, "  Key\tValue")
	fmt.Fprintln(table, "  ---\t-----")

	return table
}

func showTable(jsonKeys []string, flatJSON FlattenedJSON) {
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