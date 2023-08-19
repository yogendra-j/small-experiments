package main

import (
	"bufio"
	"fmt"
)

func countChars(filePath *string) int {
	charCount := 0

	fileStream := openFile(filePath)
	defer fileStream.Close()

	scanner := bufio.NewScanner(fileStream)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		charCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("An error occurred while reading the file:", err)
		return -1
	}

	return charCount
}
