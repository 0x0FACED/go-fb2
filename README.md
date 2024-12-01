# gofb2

This library contains a **small possible** description of the FB2 format.

**Why small possible?**

There are a lot of tags and they can be nested in many different ways. To make parsing of any fb2 book into a structure possible, there are many aspects to consider, which will probably be improved in the future.


## Features

1. Marshaling FB2 object to byte array.
2. Marshaling FB2 object to file. 
3. Unramshaling `io.Reader` to FB2 object
4. Unmarshaling `*.fb2` file to FB2 object
5. Pure Go

## Example usage

```go
package main

import (
	"fmt"
	"log"

	gofb2 "github.com/0x0FACED/go-fb2/v2"
)

func main() {
	// Path to book (fb2 ext)
	file := "/path/to/file.fb2"

	// Create parser
	parser := gofb2.NewParser()
	isPretty := true
	isStrict := true

	// Parsing file to fb2 struct
	fb2Book, err := parser.UnmarshalFromFile(file, isStrict)
	if err != nil {
		log.Fatalf("Error reading FB2: %v", err)
	}

	log.Printf("Book is unmarshaled: %+v\n", fb2Book)

	// Marshal fb2 object to file
	err = parser.MarshalToFile(fb2Book, "output.fb2", isPretty)
	if err != nil {
		log.Fatalf("Error marshaling to file: %v", err)
	}
}
```

## Install

```sh
go get -u github.com/0x0FACED/go-fb2
```

## Issues

Is there a mistake? Create an [issue](https://github.com/0x0FACED/go-fb2/issues)!

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).