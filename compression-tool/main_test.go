package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestBuildFreq(t *testing.T) {
	tests := []struct {
		text string
		want map[string]int
	}{
		{"abcdef", map[string]int{"a": 1, "b": 1, "c": 1, "d": 1, "e": 1, "f": 1}},
		{"bababc", map[string]int{"a": 2, "b": 3, "c": 1}},
		{"", map[string]int{}},
		{" ", map[string]int{" ": 1}},
		{"Word 1 eur\no €\n \n ww", map[string]int{"W": 1, "o": 2, "r": 2, "d": 1, " ": 5, "1": 1, "e": 1, "u": 1, "€": 1, "w": 2, "\n": 3}},
	}

	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(test.text))
			scanner.Split(bufio.ScanRunes)
			got := buildFreqMap(scanner)
			if pass, err := mapsEqual(got, test.want); !pass {
				t.Errorf(err.Error())
			}
		})
	}
}

func mapsEqual[T1 comparable, T2 comparable](a, b map[T1]T2) (bool, error) {
	if len(a) != len(b) {
		return false, fmt.Errorf("lengths differ: got %v, want %v", len(a), len(b))
	}

	for k, v := range a {
		if b[k] != v {

			return false, fmt.Errorf("key count miss match %v: got %v, want %v", k, b[k], v)
		}
	}

	return true, nil
}

func TestBuildHuffmanTree(t *testing.T) {
	freqMap := map[string]int{"a": 1, "b": 1, "c": 3, "d": 4, "e": 8, " ": 10}

	tree := buildHuffmanTree(&freqMap)

	if tree.char != "" {
		t.Errorf("root node should not have a character")
	}
	if tree.freq != 27 {
		t.Errorf("Expected root node to have a frequency of 27, got %v", tree.freq)
	}
	if tree.right.char != "" {
		t.Errorf("Expected right node to have a character of '', got %v", tree.right.char)
	}
	if tree.right.freq != 17 {
		t.Errorf("Expected right node to have a frequency of 17, got %v", tree.right.freq)
	}
	if tree.left.char != " " {
		t.Errorf("Expected left node to have a character of ' ', got %v", tree.left.char)
	}
	if tree.left.freq != 10 {
		t.Errorf("Expected left node to have a frequency of 10, got %v", tree.left.freq)
	}
	if tree.right.right.freq != 9 {
		t.Errorf("Expected node to have a frequency of 9, got %v", tree.right.right.freq)
	}
	if tree.right.right.char != "" {
		t.Errorf("Expected node to have a character of '', got %v", tree.right.right.char)
	}
	if tree.right.left.freq != 8 {
		t.Errorf("Expected node to have a frequency of 8, got %v", tree.right.left.freq)
	}
	if tree.right.left.char != "e" {
		t.Errorf("Expected node to have a character of 'e', got %v", tree.right.left.char)
	}
}

func TestBuildHuffmanTable(t *testing.T) {

	tree := &node{char: "", freq: 27,
		left: &node{char: " ", freq: 10,
			left:  nil,
			right: nil},
		right: &node{char: "", freq: 17,
			left: &node{char: "e", freq: 8,
				left:  nil,
				right: nil},
			right: &node{char: "", freq: 9,
				left: &node{char: "d", freq: 4,
					left:  nil,
					right: nil},
				right: &node{char: "", freq: 5,
					left: &node{char: "", freq: 2,
						left: &node{char: "b", freq: 1,
							left:  nil,
							right: nil},
						right: &node{char: "a", freq: 1,

							left:  nil,
							right: nil},
					},
					right: &node{char: "c", freq: 3,
						left:  nil,
						right: nil},
				},
			},
		},
	}

	table := *buildHuffmanTable(tree)

	if table[" "] != "0" {
		t.Errorf("Expected ' ' to have a code of '0', got %v", table[" "])
	}
	if table["e"] != "10" {
		t.Errorf("Expected 'e' to have a code of '10', got %v", table["e"])
	}
	if table["d"] != "110" {
		t.Errorf("Expected 'd' to have a code of '110', got %v", table["d"])
	}
	if table["c"] != "1111" {
		t.Errorf("Expected 'c' to have a code of '111', got %v", table["c"])
	}
	if table["b"] != "11100" {
		t.Errorf("Expected 'b' to have a code of '1100', got %v", table["b"])
	}
	if table["a"] != "11101" {
		t.Errorf("Expected 'a' to have a code of '1101', got %v", table["a"])
	}
}

func TestWriteHuffmanTable(t *testing.T) {
	ogTable := map[string]string{"a": "0", "b": "10", "#": "110", "$": "1110", "e": "1111", " ": "1", "\n": "01", "\t": "001"}
	strWriter := &strings.Builder{}
	b := bufio.NewWriter(strWriter)
	writeHuffmanTable(&ogTable, b)

	s := bufio.NewScanner(strings.NewReader(strWriter.String()))
	s.Split(bufio.ScanBytes)
	builtTable := readAndBuildHuffmanTable(s)

	if pass, err := mapsEqual(*builtTable, ogTable); !pass {
		t.Errorf(err.Error())
	}
}
