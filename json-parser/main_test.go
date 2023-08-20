package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestJsonParser_EmptyValidJson(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"{}"},
		{" {} "},
		{"{ }"},
		{"{ \n} "},
	}

	for _, test := range tests {
		scanner := bufio.NewScanner(bytes.NewReader([]byte(test.input)))
		scanner.Split(bufio.ScanRunes)
		result := jsonParser(scanner)

		if !result {
			t.Errorf("Failed for: '%v'", test.input)
		}
	}
}

func TestJsonParser_EmptyInvalidJson(t *testing.T) {
	tests := []struct {
		input string
	}{
		{""},
		{"  "},
		{"{  "},
		{"  }"},
		{"\n\n\n \n ddddddddddddddddddddddddddddddddddddddddddddddddddddd fdffffffffffffffffffff dfffffffffffffffdfdfdfdfdfddddddddddddddddddddddfdfdf dfdfdfdfefujdskgjhfjghjfhgjsfhgljs jklghsjkhgjshgjhsljghshgjsfhklgjhsfjghsjfhglsjhgjlhsfjlghsjlghsjhglshgljdlghaifghoaidgfladgfghfkgsghyogendras jaiswal test log input ect \n\n"},
	}

	for _, test := range tests {
		scanner := bufio.NewScanner(bytes.NewReader([]byte(test.input)))
		scanner.Split(bufio.ScanRunes)
		result := jsonParser(scanner)

		if result {
			t.Errorf("Failed for: '%v'", test.input)
		}
	}
}

func TestJsonParser_WithOneStringKeyValue(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"key": "value"}`, true},
		{`{"key": "value" }`, true},
		{`{"key": "value" } `, true},
		{`{"key": "value" }  `, true},
		{`{"key"
		
		: 
			 "value" }  `, true},
		{`{"key"  : "value" `, false},
		{`{"key": "value `, false},
		{`{"key": "value} `, false},
		{`{"key": "value } `, false},
		{`{key": "value
		"}`, false},
		{`{"key: "value"}`, false},
		{`{"key": value"}`, false},
		{`{"key" "value"}`, false},
	}

	for _, test := range tests {
		scanner := bufio.NewScanner(bytes.NewReader([]byte(test.input)))
		scanner.Split(bufio.ScanRunes)
		result := jsonParser(scanner)

		if result != test.expected {
			t.Errorf("Failed for: '%v'", test.input)
		}
	}

}

func TestSeekToNextNonEmptyChar(t *testing.T) {
	tests := []struct {
		input    string
		expected rune
	}{
		{" {}", '{'},
		{"{ }", '{'},
		{"{ \n} ", '{'},
		{"  	\n \t{s dfdf sdf\n} ", '{'},
		{"  	\n \ts{s dfdf sdf\n} ", 's'},
		{"", 0},
		{" ", 0},
		{"  ", 0},
		{"\n", 0},
		{"\n\n", 0},
		{"\t", 0},
		{"\t\t", 0},
	}

	for _, test := range tests {
		json := []byte(test.input)
		scanner := bufio.NewScanner(bytes.NewReader(json))
		scanner.Split(bufio.ScanRunes)

		result, _ := seekToNextNonEmptyRune(scanner)

		if result != test.expected {
			t.Errorf("Failed for: '%v'; actual: '%v'", test.input, result)
		}
	}
}
