package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	fmt.Println(getFormattedOutput(os.Args[1:]))
}

func openFile(filePath *string) *bufio.Reader {
	fileStream, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	return bufio.NewReaderSize(fileStream, 4096)
}

func getFormattedOutput(args []string) string {
	filepath := args[0]
	reader := openFile(&filepath)
	wordsCount, linesCount, bytesCount, _ := getCounts(reader)

	return fmt.Sprintf("%d %d %d %s", linesCount, wordsCount, bytesCount, filepath)
}

func getCounts(reader *bufio.Reader) (int, int, int, int) {
	wordsCount, linesCount, bytesCount, charsCount := 0, 0, 0, 0

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() != "EOF" {
				panic(err)
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
