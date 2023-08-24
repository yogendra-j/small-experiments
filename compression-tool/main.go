package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	filePath := os.Args[1]
	scanner, file := getScanner(&filePath)
	if scanner == nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println(scanner.Text())
}

func getScanner(filePath *string) (*bufio.Scanner, *os.File) {
	file, err := os.Open(*filePath)
	if err != nil {
		return nil, nil
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	return scanner, file
}

func buildFreqMap(scanner *bufio.Scanner) map[string]int {
	freqMap := make(map[string]int, 300)
	for scanner.Scan() {
		char := scanner.Text()
		freqMap[char]++
	}
	return freqMap
}
