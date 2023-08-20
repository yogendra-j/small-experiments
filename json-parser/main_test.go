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
			t.Errorf("Failed for %v", test.input)
		}
	}
}

func TestJsonParser_EmptyInvalidJson(t *testing.T) {
	tests := []struct {
		input string
	}{
		{""},
		{"  "},
		{"\n\n\n \n ddddddddddddddddddddddddddddddddddddddddddddddddddddd fdffffffffffffffffffff dfffffffffffffffdfdfdfdfdfddddddddddddddddddddddfdfdf dfdfdfdfefujdskgjhfjghjfhgjsfhgljs jklghsjkhgjshgjhsljghshgjsfhklgjhsfjghsjfhglsjhgjlhsfjlghsjlghsjhglshgljdlghaifghoaidgfladgfghfkgsghyogendras jaiswal test log input ect \n\n"},
	}

	for _, test := range tests {
		json := []byte(test.input)
		reader := bufio.NewReader(bytes.NewReader(json))

		result := jsonParser(reader)

		if result {
			t.Errorf("Failed for %v", test.input)
		}
	}
}
