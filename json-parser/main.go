package main

import (
	"bufio"
	"fmt"
	"os"
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
	for {
		_, err := reader.ReadBytes('{')
		if err != nil {
			if err == bufio.ErrBufferFull {
				continue
			}
			return false
		}
		return true

	}
}
