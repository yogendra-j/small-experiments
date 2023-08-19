package main

import (
	"fmt"
	"io"
)

func countBytes(filePath *string) int {
	byteCount := 0

	fileStream := openFile(filePath)
	defer fileStream.Close()

	buffer := make([]byte, 1024)
	for {
		bytesRead, err := fileStream.Read(buffer)
		byteCount += bytesRead
		if err != nil {
			if err != io.EOF {
				fmt.Println("An error occurred while reading the file:", err)
				return -1
			}
			break
		}
	}

	return byteCount
}
