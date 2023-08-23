package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	filepath := ""
	if len(os.Args) < 2 {
		// file name from stdin
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		filepath = scanner.Text()
	} else {
		filepath = os.Args[1]
	}
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
	if r, err := seekToNextNonEmptyRune(scanner); (r != '{' && r != '[') || err != nil {
		return false
	}

	if !parseValue(scanner) {
		return false
	}
	return scanner.Text() == ""

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
		if !commaOrEnd(scanner, '}') {
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

func commaOrEnd(scanner *bufio.Scanner, endToken rune) bool {
	r, _ := utf8.DecodeLastRuneInString(scanner.Text())
	if unicode.IsSpace(r) {
		r, _ = seekToNextNonEmptyRune(scanner)
	}
	if r == ',' {
		r, _ := seekToNextNonEmptyRune(scanner)
		if r == endToken {
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
	_, err := seekToNextNonEmptyRune(scanner)
	if err != nil {
		return false
	}

	return parseValue(scanner)
}

func parseString(scanner *bufio.Scanner) bool {
	if scanner.Text() != `"` {
		return false
	}
	token := ""
	for scanner.Scan() {
		str := scanner.Text()
		r, _ := utf8.DecodeRuneInString(str)
		if r == '\n' {
			return false
		}
		if r == '"' && isNotEscaped(token) {
			scanner.Scan()
			break
		}
		token += str
	}
	return scanner.Err() == nil && isValidString(token)
}

func isValidString(token string) bool {
	//check for invalid escape characters (\x etc) and unescaped control characters
	for i, c := range token {
		if isUnescapedControlCharacter(c) || isInvalidEscapeCharacter(token, i) {
			return false
		}
	}
	return true
}

func isNotEscaped(token string) bool {
	if len(token) == 0 {
		return true
	}
	count := 0
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '\\' {
			count++
		} else {
			break
		}
	}
	return count%2 == 0
}

func isInvalidEscapeCharacter(token string, i int) bool {
	controlChars := []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't', 'u'}
	return token[i] == '\\' && i < len(token)-1 && !contains(controlChars, rune(token[i+1]))
}

func contains[T comparable](arr []T, r T) bool {
	for _, a := range arr {
		if a == r {
			return true
		}
	}
	return false
}

func isUnescapedControlCharacter(c rune) bool {
	return unicode.IsControl(c)
}

func parseValue(scanner *bufio.Scanner) bool {
	r, _ := utf8.DecodeLastRuneInString(scanner.Text())
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
	case '[':
		if !parseArray(scanner) {
			return false
		}
	default:
		return false
	}
	return true
}

func colon(scanner *bufio.Scanner) bool {
	r, _ := utf8.DecodeRuneInString(scanner.Text())
	if unicode.IsSpace(r) {
		r, _ = seekToNextNonEmptyRune(scanner)
	}
	return scanner.Err() == nil && r == ':'
}

func parseNumber(scanner *bufio.Scanner) bool {
	token := scanner.Text()
	for scanner.Scan() {
		str := scanner.Text()
		r, _ := utf8.DecodeRuneInString(str)
		if !unicode.IsDigit(r) && r != '.' && r != '-' && r != 'e' && r != 'E' && r != '+' {
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
	token = strings.ToLower(token)
	exp1, exp2 := "", ""
	if countOccurrences(token, 'e') >= 2 {
		return false
	}
	if countOccurrences(token, 'e') == 1 {
		exp1, exp2 = strings.Split(token, "e")[0], strings.Split(token, "e")[1]
		if !isValidPower(exp2) {
			return false
		}
	} else if countOccurrences(token, 'e') == 0 {
		exp1, exp2 = token, ""
	}
	if exp1[0] == '0' && len(exp1) > 1 && exp1[1] != '.' {
		return false
	}
	if exp1[len(exp1)-1] == '.' {
		return false
	}
	if countOccurrences(exp1, '.') > 1 {
		return false
	}
	if countOccurrences(exp1, '-') > 0 {
		return false
	}
	return true
}

func isValidPower(token string) bool {
	if len(token) == 0 {
		return false
	}
	if token[0] == '+' || token[0] == '-' {
		token = token[1:]
	}
	if len(token) == 0 {
		return false
	}
	for _, c := range token {
		if !unicode.IsDigit(c) {
			return false
		}
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

func parseArray(scanner *bufio.Scanner) bool {
	if scanner.Text() != "[" {
		return false
	}
	if r, _ := seekToNextNonEmptyRune(scanner); r == ']' {
		seekToNextNonEmptyRune(scanner)
		return true
	}
	for scanner.Text() != "]" {
		if !parseValue(scanner) {
			return false
		}
		if !commaOrEnd(scanner, ']') {
			return false
		}
	}
	if scanner.Text() != "]" {
		return false
	}
	seekToNextNonEmptyRune(scanner)
	return true
}
