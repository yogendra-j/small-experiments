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
	if r, err := seekToNextNonEmptyRune(scanner); r != '}' || err != nil {
		if r != '"' {
			return false
		}
		for scanner.Scan() {
			str := scanner.Text()
			r, _ := utf8.DecodeRuneInString(str)
			if r == '"' {
				break
			}
		}
		if scanner.Err() != nil {
			return false
		}
		if r, err := seekToNextNonEmptyRune(scanner); r != ':' || err != nil {
			return false
		}
		if r, err := seekToNextNonEmptyRune(scanner); r != '"' || err != nil {
			return false
		}
		for scanner.Scan() {
			str := scanner.Text()
			r, _ := utf8.DecodeRuneInString(str)
			if r == '"' {
				break
			}
		}
		if scanner.Err() != nil {
			return false
		}
		if r, err := seekToNextNonEmptyRune(scanner); r != '}' || err != nil {
			return false
		}
	}

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
