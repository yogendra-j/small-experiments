package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 3 {
		filePath := os.Args[1]
		scanner, file := getScanner(&filePath)

		defer file.Close()

		filePath = os.Args[2]
		writer, file := getWriter(&filePath)

		defer file.Close()

		table := buildAndWriteHuffmanTable(scanner, writer)

		filePath = os.Args[1]
		scanner, file = getScanner(&filePath)

		defer file.Close()

		encodeAndWriteFile(scanner, table, writer)
	}
	if len(os.Args) == 4 {
		filePath := os.Args[2]
		scanner, file := getScanner(&filePath)
		defer file.Close()
		scanner.Split(bufio.ScanBytes)

		filePath = os.Args[3]
		writer, file := getWriter(&filePath)
		defer file.Close()
		decompressAndWriteFile(scanner, writer)
	}
}

func decompressAndWriteFile(scanner *bufio.Scanner, writer *bufio.Writer) {

	table := readAndBuildHuffmanTable(scanner)

	decodeAndWriteFile(scanner, table, writer)
}

func buildAndWriteHuffmanTable(scanner *bufio.Scanner, writer *bufio.Writer) (table *map[string]string) {

	freqMap := buildFreqMap(scanner)
	tree := buildHuffmanTree(&freqMap)
	table = buildHuffmanTable(tree)

	writeHuffmanTable(table, writer)

	return
}

func getScanner(filePath *string) (*bufio.Scanner, *os.File) {
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	return scanner, file
}

func getWriter(filePath *string) (*bufio.Writer, *os.File) {
	file, err := os.Create(*filePath)
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}

	writer := bufio.NewWriter(file)
	return writer, file
}

func buildFreqMap(scanner *bufio.Scanner) map[string]int {
	freqMap := make(map[string]int, 300)
	for scanner.Scan() {
		char := scanner.Text()
		freqMap[char]++
	}
	freqMap["EOF"] = 0
	return freqMap
}

func buildHuffmanTree(freqMap *map[string]int) *node {
	h := &nodeHeap{}
	for char, freq := range *freqMap {
		*h = append(*h, &node{char: char, freq: freq})
	}
	heap.Init(h)

	for h.Len() > 1 {
		n1 := heap.Pop(h).(*node)
		n2 := heap.Pop(h).(*node)
		heap.Push(h, &node{char: "", freq: n1.freq + n2.freq, left: n1, right: n2})
	}

	return heap.Pop(h).(*node)
}

type node struct {
	char  string
	freq  int
	left  *node
	right *node
}

type nodeHeap []*node

func (h nodeHeap) Len() int { return len(h) }

func (h nodeHeap) Less(i, j int) bool {
	return h[i].freq < h[j].freq
}

func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *nodeHeap) Push(x interface{}) {
	*h = append(*h, x.(*node))
}

func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return x
}

type nodeCodePair struct {
	node *node
	code string
}

func buildHuffmanTable(tree *node) *map[string]string {
	table := make(map[string]string)

	level := []*nodeCodePair{{node: tree, code: ""}}

	for len(level) > 0 {
		nextLevel := []*nodeCodePair{}
		for _, ncp := range level {
			if ncp.node.char != "" {
				table[ncp.node.char] = ncp.code
			} else {
				nextLevel = append(nextLevel, &nodeCodePair{node: ncp.node.left, code: ncp.code + "0"})
				nextLevel = append(nextLevel, &nodeCodePair{node: ncp.node.right, code: ncp.code + "1"})
			}
		}
		level = nextLevel
	}

	return &table
}

func writeHuffmanTable(table *map[string]string, writer *bufio.Writer) {
	for char, code := range *table {
		writer.WriteString(char + "$#" + code + "$#")
	}
	writer.WriteString("$#")
	writer.Flush()
}

func readAndBuildHuffmanTable(scanner *bufio.Scanner) *map[string]string {
	table := make(map[string]string)
	for {
		chr := readOneCol(scanner)
		if chr == "" {
			break
		}
		code := readOneCol(scanner)
		table[code] = chr
	}
	return &table
}

func readOneCol(scanner *bufio.Scanner) string {
	readBytes := []byte{}
	for scanner.Scan() {
		b := scanner.Bytes()
		if b[len(b)-1] == '#' && len(readBytes) > 0 && readBytes[len(readBytes)-1] == '$' {
			return string(readBytes[:len(readBytes)-1])
		}
		readBytes = append(readBytes, b...)
	}
	return ""
}

func encodeAndWriteFile(scanner *bufio.Scanner, table *map[string]string, writer *bufio.Writer) {
	var bufEncodedStr uint8 = 0
	var bitsInBuffer uint8 = 0
	for scanner.Scan() {
		char := scanner.Text()
		writeBitsToFile((*table)[char], writer, &bitsInBuffer, &bufEncodedStr)
	}
	writeBitsToFile((*table)["EOF"], writer, &bitsInBuffer, &bufEncodedStr)
	if bitsInBuffer > 0 {
		bufEncodedStr = bufEncodedStr << (8 - bitsInBuffer)
		writer.WriteByte(bufEncodedStr)
	}
	writer.Flush()
}

func writeBitsToFile(bitString string, writer *bufio.Writer, bitsInBuffer *uint8, bufEncodedStr *uint8) {
	for _, bitFlag := range bitString {
		if bitFlag == '1' {
			*bufEncodedStr = (*bufEncodedStr << 1) | 1
		} else {
			*bufEncodedStr = *bufEncodedStr << 1
		}
		*bitsInBuffer++
		if *bitsInBuffer == 8 {
			writer.WriteByte(*bufEncodedStr)
			*bufEncodedStr = 0
			*bitsInBuffer = 0
		}
	}
}

func decodeAndWriteFile(scanner *bufio.Scanner, table *map[string]string, writer *bufio.Writer) {
	currCode := ""
	defer writer.Flush()
	for scanner.Scan() {
		b := scanner.Bytes()[0]
		for i := 7; i >= 0; i-- {
			if b&(1<<uint(i)) > 0 {
				currCode += "1"
			} else {
				currCode += "0"
			}
			if char, ok := (*table)[currCode]; ok && char != "EOF" {
				writer.WriteString(char)
				currCode = ""
			} else if char == "EOF" {
				return
			}
		}
	}
}
