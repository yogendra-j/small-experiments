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
	if *filePath == "" {
		return bufio.NewReader(os.Stdin), nil
	}
	fileStream, err := os.Open(*filePath)
	if err != nil {
		return nil, errors.New("FILE READ ERROR")
	}
	return bufio.NewReaderSize(fileStream, 4096), nil
}

func getFormattedOutput(args []string) string {
	wordsFlag, linesFlag, bytesFlag, charsFlag, filepath, err := parseArgs(args)
	if err != nil {
		return "INVALID ARGUMENT"
	}
	reader, err := openFile(&filepath)
	if err != nil {
		return err.Error()
	}
	wordsCount, linesCount, bytesCount, charsCount, err := getCounts(reader)

	if err != nil {
		return err.Error()
	}

	output := ""

	if linesFlag {
		output += fmt.Sprintf("%d ", linesCount)
	}
	if wordsFlag {
		output += fmt.Sprintf("%d ", wordsCount)
	}
	if bytesFlag {
		output += fmt.Sprintf("%d ", bytesCount)
	}
	if charsFlag {
		output += fmt.Sprintf("%d ", charsCount)
	}
	if !(linesFlag || wordsFlag || bytesFlag || charsFlag) {
		output += fmt.Sprintf("%d %d %d ", linesCount, wordsCount, bytesCount)
	}
	if filepath != "" {
		output += filepath
	}
	return output
}

func getCounts(reader *bufio.Reader) (int, int, int, int, error) {
	wordsCount, linesCount, bytesCount, charsCount := 0, 0, 0, 0

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return 0, 0, 0, 0, errors.New("FILE READ ERROR")
			}
			break
		}
		linesCount++
		bytesCount += len(line)
		charsCount += len(bytes.Runes(line))
		wordsCount += len(bytes.Fields(line))
	}
	return wordsCount, linesCount, bytesCount, charsCount, nil
}

func parseArgs(args []string) (bool, bool, bool, bool, string, error) {
	wordsFlag, linesFlag, bytesFlag, charsFlag := false, false, false, false
	filepath := ""

	for _, arg := range args {
		if isFlag(arg) {
			for _, flag := range arg[1:] {
				switch flag {
				case 'w':
					wordsFlag = true
				case 'l':
					linesFlag = true
				case 'c':
					bytesFlag = true
				case 'm':
					charsFlag = true
				}
			}
		} else {
			if filepath == "" {
				filepath = arg
			} else if arg != "" {
				return false, false, false, false, "", errors.New("INVALID ARGUMENT")
			}
		}
	}
	return wordsFlag, linesFlag, bytesFlag, charsFlag, filepath, nil
}

func isFlag(arg string) bool {
	return arg[0] == '-'
}
