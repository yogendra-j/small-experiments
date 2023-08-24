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

func mapsEqual(a, b map[string]int) (bool, error) {
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
