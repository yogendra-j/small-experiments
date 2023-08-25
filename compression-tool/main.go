package main

import (
	"bufio"
	"container/heap"
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

	fmt.Println(buildFreqMap(scanner)["t"])
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
