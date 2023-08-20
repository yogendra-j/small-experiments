package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	filepath := os.Args[1]

	reader, err := getFileReader(&filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if jsonParser(reader) {
		fmt.Println("Valid JSON")
		os.Exit(0)
	} else {
		fmt.Println("Invalid JSON")
		os.Exit(1)
	}
}

func getFileReader(filepath *string) (*bufio.Reader, error) {
	file, err := os.Open(*filepath)
	if err != nil {
		return nil, err
	}
	return bufio.NewReaderSize(file, 1024), nil
}

func jsonParser(reader *bufio.Reader) bool {
	return seekToCharSkippingWhitespace('{', reader) && seekToCharSkippingWhitespace('}', reader)
}

func seekToCharSkippingWhitespace(char byte, reader *bufio.Reader) bool {
	for {
		bytes, err := reader.ReadBytes(char)
		trimmedBytes := strings.TrimSpace(string(bytes))
		if len(trimmedBytes) != 1 {
			return false
		}
		if err != nil && err != bufio.ErrBufferFull {
			return false
		} else if err == bufio.ErrBufferFull {
			continue
		}
		return char == trimmedBytes[0]
	}
}
