package main

import (
	"fmt"
	"errors"
	"strconv"
	"unicode"
)

type Parser struct {
	// The string being parsed
	bencodeString string

	// Where are we in the bencoded string?
	cursor int
}

func Parse(parser *Parser) string {
	return ParseBencode(parser.bencodeString)
}

func ParseBencode(bencodeString string) []string {
	var parsedElements []string

	parser := Parser { bencodeString, 0 }

	for parser.cursor != len(bencodeString) - 1 {
		char := bencodeString[parser.cursor]

		if unicode.IsDigit(rune(char)) {
			decodedString, err := processString(&parser)
			if (err != nil) {
				fmt.Println(err);
				return parsedElements;
			}
			parsedElements = append(parsedElements, decodedString)
			fmt.Println(decodedString)
		}

		if char == 'i' {
			decodedInteger, err := processInteger(&parser)
			if (err != nil) {
				fmt.Println(err)
				return parsedElements;
			}
			parsedElements = append(parsedElements, decodedInteger)
			fmt.Println(decodedInteger)
		}

		if char == 'l' {
			decodedList, err := processList(&parser)
			if (err != nil) {
				fmt.Println(err)
				return parsedElements;
			}
			parsedElements = append(parsedElements, decodedList)
			fmt.Println(decodedList)
		}
	}

	return parsedElements
}

// <LIST>  ::= "l" 1 * <BE>         "e"
func processList(parser *Parser) (int, error) {
	var listString string

	fmt.Println("Parsing list: " + parser.bencodeString)

	// skip over 'l'
	parser.cursor += 1

	parsedList := Parse(&parser)

	return 0, errors.New(listString)
}

// <INT>   ::= "i"     <SNUM>       "e"
func processInteger(parser *Parser) (int, error) {
	var integerString string

	// skip over 'i'
	parser.cursor += 1

	for unicode.IsDigit(rune(parser.bencodeString[parser.cursor])) {
		integerString += string(parser.bencodeString[parser.cursor])
		parser.cursor += 1
	}
	
	if (parser.bencodeString[parser.cursor] != 'e') {
		return 0, errors.New("Unexpected end of an integer. Expected 'e'")
	}

	// skip over 'e'
	parser.cursor += 1

	finalInteger, err := strconv.Atoi(integerString)
	if (err != nil) {
		return 0, errors.New(fmt.Sprintf("Error converting string length to an integer: %s", err))
	}

	return finalInteger, nil
}
 
// <STR>   ::= <NUM> ":" n * <CHAR>; where n equals the <NUM>
func processString(parser *Parser) (string, error) {
	var stringLengthString string

	for pos, char := range parser.bencodeString[parser.cursor:] {
		if unicode.IsDigit(char) {
			stringLengthString += string(char)
		} else {
			parser.cursor += pos
			break
		}
	}

	stringLength, err := strconv.Atoi(stringLengthString)
	if (err != nil) {
		return "", errors.New(fmt.Sprintf("Error converting string length to an integer: %s", err))
	}

	if parser.bencodeString[parser.cursor] != ':' {
		return "", errors.New("Not a string?")
	}

	// skip over ':'
	parser.cursor += 1

	final := parser.bencodeString[parser.cursor:parser.cursor+stringLength]

	parser.cursor += stringLength

	return final, nil
}
