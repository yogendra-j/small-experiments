package main

import (
	"bufio"
	"fmt"
)

func countLines(filePath *string) int {
	lineCount := 0

	fileStream := openFile(filePath)
	defer fileStream.Close()

	scanner := bufio.NewScanner(fileStream)

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("An error occurred while reading the file:", err)
		return -1
	}

	return lineCount

}
