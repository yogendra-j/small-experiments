package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	filepath := os.Args[1]

	scanner, file, err := getRuneScanner(&filepath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if jsonParser(scanner) {
		fmt.Println("Valid JSON")
		os.Exit(0)
	} else {
		fmt.Println("Invalid JSON")
		os.Exit(1)
	}
}

func getRuneScanner(filepath *string) (*bufio.Scanner, *os.File, error) {
	file, err := os.Open(*filepath)
	if err != nil {
		return nil, file, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	return scanner, file, nil
}

func jsonParser(scanner *bufio.Scanner) bool {
	if r, err := seekToNextNonEmptyRune(scanner); r != '{' || err != nil {
		return false
	}

	return parseObject(scanner)

}

func parseObject(scanner *bufio.Scanner) bool {
	if scanner.Text() != "{" {
		return false
	}
	if r, _ := seekToNextNonEmptyRune(scanner); r == '}' {
		seekToNextNonEmptyRune(scanner)
		return true
	}
	for r, _ := utf8.DecodeLastRuneInString(scanner.Text()); r == '"'; r, _ = utf8.DecodeLastRuneInString(scanner.Text()) {
		if !parseKeyValuePair(scanner) {
			return false
		}
		if !commaOrEnd(scanner) {
			return false
		}
	}
	if scanner.Text() != "}" {
		return false
	}
	seekToNextNonEmptyRune(scanner)
	return true
}

func seekToNextNonEmptyRune(scanner *bufio.Scanner) (rune, error) {
	for scanner.Scan() {
		str := scanner.Text()
		r, _ := utf8.DecodeRuneInString(str)
		if !unicode.IsSpace(r) {
			return r, nil
		}
	}
	return 0, scanner.Err()
}

func commaOrEnd(scanner *bufio.Scanner) bool {
	r, _ := utf8.DecodeLastRuneInString(scanner.Text())
	if unicode.IsSpace(r) {
		r, _ = seekToNextNonEmptyRune(scanner)
	}
	if r == ',' {
		r, _ := seekToNextNonEmptyRune(scanner)
		if r == '}' {
			return false
		}
	}
	return true
}

func parseKeyValuePair(scanner *bufio.Scanner) bool {
	if !parseString(scanner) {
		return false
	}
	if scanner.Err() != nil || !colon(scanner) {
		return false
	}

	r, err := seekToNextNonEmptyRune(scanner)
	if err != nil {
		return false
	}
	switch r {
	case '"':
		if !parseString(scanner) {
			return false
		}
	case '{':
		if !parseObject(scanner) {
			return false
		}
	case 't', 'f':
		if !parseBoolean(scanner) {
			return false
		}
	case 'n':
		if !parseNull(scanner) {
			return false
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-':
		if !parseNumber(scanner) {
			return false
		}
	default:
		return false
	}
	return true
}

func parseString(scanner *bufio.Scanner) bool {
	if scanner.Text() != `"` {
		return false
	}
	for scanner.Scan() {
		str := scanner.Text()
		r, _ := utf8.DecodeRuneInString(str)
		if r == '\n' {
			return false
		}
		if r == '"' {
			scanner.Scan()
			break
		}
	}
	return scanner.Err() == nil
}

func colon(scanner *bufio.Scanner) bool {
	r, _ := utf8.DecodeRuneInString(scanner.Text())
	if unicode.IsSpace(r) {
		r, _ = seekToNextNonEmptyRune(scanner)
	}
	return scanner.Err() == nil || r == ':'
}

func parseNumber(scanner *bufio.Scanner) bool {
	token := scanner.Text()
	for scanner.Scan() {
		str := scanner.Text()
		r, _ := utf8.DecodeRuneInString(str)
		if !unicode.IsDigit(r) && r != '.' && r != '-' {
			break
		}
		token += str
	}
	return scanner.Err() == nil && isValidNumber(token)
}

func parseNull(scanner *bufio.Scanner) bool {
	token := scanner.Text()
	for scanner.Scan() {
		str := scanner.Text()
		r, _ := utf8.DecodeRuneInString(str)
		if !unicode.IsLetter(r) {
			break
		}
		token += str
	}
	return scanner.Err() == nil && token == "null"
}

func parseBoolean(scanner *bufio.Scanner) bool {
	token := scanner.Text()
	for scanner.Scan() {
		str := scanner.Text()
		r, _ := utf8.DecodeRuneInString(str)
		if !unicode.IsLetter(r) {
			break
		}
		token += str
	}
	return scanner.Err() == nil && (token == "true" || token == "false")
}

func isValidNumber(token string) bool {
	if len(token) > 0 && token[0] == '-' {
		token = token[1:]
	}
	if len(token) == 0 {
		return false
	}
	if token[0] == '0' && len(token) > 1 {
		return false
	}
	if token[len(token)-1] == '.' {
		return false
	}
	if count := countOccurrences(token, '.'); count > 1 {
		return false
	}
	if count := countOccurrences(token, '-'); count > 0 {
		return false
	}
	return true
}

func countOccurrences(token string, char rune) int {
	count := 0
	for _, c := range token {
		if c == char {
			count++
		}
	}
	return count
}
