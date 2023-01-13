# JSONTable

View json in command-line table.

## input

```sh
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
	}
```

# output

```sh
      Key                     Value
      ---                     -----
      name                    John Smith
      age                     35
      address_street          Main Street
      address_city            New York
      address_state           NY
      phoneNumbers_0_number   212 555-1234
      phoneNumbers_0_type     home
      phoneNumbers_1_type     work
      phoneNumbers_1_number   646 555-4567
```
