package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
)

func main() {
	fmt.Println(getFormattedOutput(os.Args[1:]))
}

func openFile(filePath *string) (*bufio.Reader, error) {
	fileStream, err := os.Open(*filePath)
	if err != nil {
		return nil, errors.New("FILE NOT FOUND")
	}
	return bufio.NewReaderSize(fileStream, 4096), nil
}

func getFormattedOutput(args []string) string {
	filepath := args[0]
	reader, err := openFile(&filepath)
	if err != nil {
		return fmt.Sprintf("0 0 0 %s", filepath)
	}
	wordsCount, linesCount, bytesCount, _ := getCounts(reader)

	return fmt.Sprintf("%d %d %d %s", linesCount, wordsCount, bytesCount, filepath)
}

func getCounts(reader *bufio.Reader) (int, int, int, int) {
	wordsCount, linesCount, bytesCount, charsCount := 0, 0, 0, 0

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return 0, 0, 0, 0
			}
			break
		}
		linesCount++
		bytesCount += len(line)
		charsCount += len(bytes.Runes(line))
		wordsCount += len(bytes.Fields(line))
	}
	return wordsCount, linesCount, bytesCount, charsCount
}
