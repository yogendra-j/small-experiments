package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fset := flag.NewFlagSet("gocut", flag.ExitOnError)
	cols := fset.String("f", "", "comma/space separated list of field indices")
	sep := fset.String("d", "\t", "field delimiter")
	fset.Parse(os.Args[1:])
	files := fset.Args()
	intCols := parseToInts(cols)
	for _, file := range files {
		cut(file, intCols, *sep)
	}
}

func parseToInts(cols *string) *[]int {
	sep := " "
	if strings.Contains(*cols, ",") {
		sep = ","
	}
	fields := strings.Split(*cols, sep)
	intCols := make([]int, len(fields))
	for i, field := range fields {
		pc, err := strconv.Atoi(field)
		if err != nil {
			log.Fatal(err)
		}
		intCols[i] = pc
	}
	return &intCols
}

func cut(fileName string, cols *[]int, sep string) {
	scanner, file := getScanner(fileName)
	defer file.Close()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, sep)
		lineOut := ""
		for i, field := range fields {
			if contains(*cols, i+1) {
				lineOut += field + "\t"
			}
		}
		fmt.Println(strings.TrimSuffix(lineOut, "\t"))
	}
}

func contains(arr []int, x int) bool {
	for _, item := range arr {
		if item == x {
			return true
		}
	}
	return false
}

func getScanner(fileName string) (*bufio.Scanner, *os.File) {
	var file *os.File
	if fileName == "-" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	return scanner, file
}