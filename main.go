package main

import (
	"fmt"
	"os"
)

func main() {
	var args []string = os.Args
	if len(args) < 2 {
		fmt.Println("Missing filename!")
		os.Exit(1)
	}
	
	bencodeString, err := readFile(args[1])
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Error reading file!")
	}
	
	parsedElements, err := Parse(bencodeString)
	if (err != nil) {
		fmt.Println(err)
	}
	fmt.Println(parsedElements)
}

func readFile(filename string) (string, error) {
	bencodeString, err := os.ReadFile(filename)
	return string(bencodeString), err
}


