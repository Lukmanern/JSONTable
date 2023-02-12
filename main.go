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

// parseJSON func used for 
// get and parse JSON
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

	var flatJSON interface{}
	err := json.Unmarshal([]byte(jsonData), &flatJSON)
	if err != nil {
		log.Fatal("Error : JSON. Please check the struture.")
		return flatJSON
	}

	return flatJSON
}

func flattenMap(data interface{}, args FlatteningArgs) (FlattenedJSON, []string) {
	// type of data-var
	dataType := reflect.TypeOf(data).Kind()
	if dataType == reflect.Map {
		// iterate through the map
		for key, value := range data.(map[string]interface{}) {
			// append the current key to the keys slice
			args.Keys = append(args.Keys, args.parentKey+key)
			// make new args
			newArgs := FlatteningArgs{
				flatJSON: args.flatJSON,
				Keys: args.Keys,
				parentKey: args.parentKey+key+"_",
			}
			// recursively call the flattenMap function with 
			// the value, updated args.parentKey, 
			// args.flatJSON and args.Keys
			args.flatJSON, args.Keys = flattenMap(value, newArgs)
		}
	} else if dataType == reflect.Slice {
		// iterate through the slice
		for i, value := range data.([]interface{}) {
			// append the current index to the args.Keys slice
			args.Keys = append(args.Keys, args.parentKey+strconv.Itoa(i))
			// make new args
			newArgs := FlatteningArgs{
				flatJSON: args.flatJSON,
				Keys: args.Keys,
				parentKey: args.parentKey+strconv.Itoa(i)+"_",
			}
			// recursively call the flattenMap function with 
			// the value, updated args.parentKey, 
			// args.flatJSON and args.Keys
			args.flatJSON, args.Keys = flattenMap(value, newArgs)
		}
	} else {
		// if the data is neither a map nor 
		// a slice, it is a leaf node
		// add the key-value pair 
		// to the args.flatJSON map
		args.flatJSON[strings.TrimSuffix(args.parentKey, "_")] = data
	}

	// return the args.flatJSON and args.Keys
	return args.flatJSON, args.Keys
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
		// when null value, skip
		if _, ok := flatJSON[key]; ok {
			key = fmt.Sprintf("%v\t%v", key, flatJSON[key])
			fmt.Fprintln(tabWriter, key)
		}
	}
	
	// Sync() for table.Flush (isn't stable)
	os.Stdout.Sync()
	tabWriter.Flush()
}