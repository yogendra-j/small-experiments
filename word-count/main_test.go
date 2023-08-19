package main

import (
	"testing"
)

func TestCountBytes(t *testing.T) {
	filePath := "test.txt"
	expected := 342384
	actual := countBytes(&filePath)
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func TestCountLines(t *testing.T) {
	filePath := "test.txt"
	expected := 7189
	actual := countLines(&filePath)
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func TestCountWords(t *testing.T) {
	filePath := "test.txt"
	expected := 58164
	actual := countWords(&filePath)
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func TestCountChars(t *testing.T) {
	filePath := "test.txt"
	expected := 339486
	actual := countChars(&filePath)
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}
