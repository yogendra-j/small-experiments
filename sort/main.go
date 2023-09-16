package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	filePath := os.Args[1]

	lr, file := lineScanner(filePath)
	defer file.Close()

	handleDownStreamExit()

	lines := make([]string, 0)

	for lr.Scan() {
		lines = append(lines, lr.Text())
	}

	lines = sort(lines)

	for _, line := range lines {
		fmt.Println(line)
	}
}

func handleDownStreamExit() {

	signalChannel := make(chan os.Signal, 1)

	// Notify the signal channel when a SIGPIPE is received
	signal.Notify(signalChannel, syscall.SIGPIPE)

	// Start a goroutine that will terminate the program when a SIGPIPE is received
	go func() {
		<-signalChannel
		os.Exit(0)
	}()
}

func lineScanner(filePath string) (*bufio.Scanner, *os.File) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	lr := bufio.NewScanner(file)
	lr.Split(bufio.ScanLines)
	return lr, file
}

func sort(lines []string) []string {
	return mergeSort(lines)
}

func mergeSort(lines []string) []string {
	if len(lines) <= 1 {
		return lines
	}

	mid := len(lines) / 2

	left := mergeSort(lines[:mid])
	right := mergeSort(lines[mid:])

	return merge(left, right)
}

func merge(left, right []string) []string {
	result := make([]string, 0, len(left)+len(right))

	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}
