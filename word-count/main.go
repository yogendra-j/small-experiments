package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	needWordCountFlagPtr := flag.Bool("w", false, "Count words in the input")
	needLineCountFlagPtr := flag.Bool("l", false, "Count lines in the input")
	needByteCountFlagPtr := flag.Bool("c", false, "Count bytes in the input")
	needCharCountFlagPtr := flag.Bool("m", false, "Count characters in the input")

	flag.Parse()

	filePath := flag.Arg(0)

	output := filePath

	if *needByteCountFlagPtr {
		output = fmt.Sprintf("%d %s", countBytes(&filePath), output)
	}
	if *needWordCountFlagPtr {
		output = fmt.Sprintf("%d %s", countWords(&filePath), output)
	}
	if *needLineCountFlagPtr {
		output = fmt.Sprintf("%d %s", countLines(&filePath), output)
	}
	if *needCharCountFlagPtr {
		output = fmt.Sprintf("%d %s", countChars(&filePath), output)
	}
	if isAllFalse(*needByteCountFlagPtr, *needWordCountFlagPtr, *needLineCountFlagPtr) {
		output = fmt.Sprintf("%d %d %d %s", countWords(&filePath), countLines(&filePath), countBytes(&filePath), output)
	}

	fmt.Println(output)

}

func openFile(filePath *string) *os.File {
	fileStream, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	return fileStream
}

func isAllFalse(args ...bool) bool {
	for _, v := range args {
		if v {
			return false
		}
	}
	return true
}
