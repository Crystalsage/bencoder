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

func Parse(bencodeString string) ([]interface{}, error) {
	parser := Parser { bencodeString, 0 }

	var parsedElements []interface{}

	for (parser.cursor != len(parser.bencodeString) - 1) {
		decodedElement, err := ParseBencode(&parser)

		if (err != nil) {
			return nil, err
		}
		parsedElements = append(parsedElements, decodedElement)
	}

	return parsedElements, nil
}

// <BE>    ::= <DICT> | <LIST> | <INT> | <STR>
func ParseBencode(parser *Parser) (interface{}, error) {
	char := parser.bencodeString[parser.cursor]

	if unicode.IsDigit(rune(char)) {
		decodedString, err := processString(parser)
		if (err != nil) {
			return nil, err
		}
		return decodedString, nil
	}

	if char == 'i' {
		decodedInteger, err := processInteger(parser)
		if (err != nil) {
			return nil, err;
		}
		return decodedInteger, nil;
	}

	if char == 'l' {
		decodedList, err := processList(parser)
		if (err != nil) {
			return nil, err;
		}
		return decodedList, nil
	}

	if char == 'd' {
		decodedDictionary, err := processDictionary(parser)
		if (err != nil) {
			return nil, err;
		}
		return decodedDictionary, nil
	}
	return nil, errors.New("Something went wrong!")
}

// <DICT>  ::= "d" 1 * (<STR> <BE>) "e"
func processDictionary(parser *Parser) (map[string]interface{}, error) {
	decodedDictionary := make(map[string]interface{})

	// skip over 'd'
	parser.cursor += 1

	for parser.bencodeString[parser.cursor] != 'e' {
		key, keyErr := processString(parser)
		value, valueErr := ParseBencode(parser)

		if keyErr != nil {
			return nil, keyErr
		} else if valueErr != nil {
			return nil, valueErr
		}

		decodedDictionary[key] = value
	}

	if (parser.bencodeString[parser.cursor] == 'e') {
		parser.cursor += 1
	}

	return decodedDictionary, nil
}

// <LIST>  ::= "l" 1 * <BE>         "e"
func processList(parser *Parser) (interface{}, error) {
	var parsedList []interface{}
	// skip over 'l'
	parser.cursor += 1

	for parser.bencodeString[parser.cursor] != 'e' {
		parsedElement, err := ParseBencode(parser)
		if (err != nil) {
			return nil, err
		}
		parsedList = append(parsedList, parsedElement)
	}

	if (parser.bencodeString[parser.cursor] != 'e') {
		return nil, errors.New("Unexpected end of list. Expected 'e'")
	}

	parser.cursor += 1

	return parsedList, nil
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
