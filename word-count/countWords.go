package main

import (
	"bufio"
	"fmt"
)

func countWords(filePath *string) int {
	wordCount := 0

	fileStream := openFile(filePath)
	defer fileStream.Close()

	scanner := bufio.NewScanner(fileStream)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wordCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("An error occurred while reading the file:", err)
		return -1
	}

	return wordCount

}