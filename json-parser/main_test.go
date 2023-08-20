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
		json := []byte(test.input)
		reader := bufio.NewReader(bytes.NewReader(json))

		result := jsonParser(reader)

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
		json := []byte(test.input)
		reader := bufio.NewReader(bytes.NewReader(json))

		result := jsonParser(reader)

		if result {
			t.Errorf("Failed for: '%v'", test.input)
		}
	}
}

func TestSeekToCharSkippingWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		char     byte
		expected bool
	}{
		{"{}", '{', true},
		{"{ }", '{', true},
		{"{ \n}", '\n', false},
		{"  	{ \n}", ' ', false},
		{"  	t{ \n}", 't', true},
	}

	for _, test := range tests {
		json := []byte(test.input)
		reader := bufio.NewReader(bytes.NewReader(json))

		result := seekToCharSkippingWhitespace(test.char, reader)

		if result != test.expected {
			t.Errorf("Failed for: '%v'", test.input)
		}
	}
}
