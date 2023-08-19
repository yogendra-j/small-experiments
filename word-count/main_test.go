package main

import (
	"testing"
)

func TestFormattedOutput_NoFlags_ValidFile(t *testing.T) {
	filePath := "test.txt"
	expected := "7189 58164 342384 test.txt"
	actual := getFormattedOutput([]string{filePath})
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormattedOutput_NoFlags_InvalidFile(t *testing.T) {
	filePath := "test1.txt"
	expected := "FILE READ ERROR"
	actual := getFormattedOutput([]string{filePath})
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormattedOutput_BytesFlag_ValidFile(t *testing.T) {
	filePath := "test.txt"
	expected := "342384 test.txt"
	actual := getFormattedOutput([]string{"-c", filePath})
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormattedOutput_BytesFlag_InvalidFile(t *testing.T) {
	filePath := "test1.txt"
	expected := "FILE READ ERROR"
	actual := getFormattedOutput([]string{"-c", filePath})
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormattedOutput_CharsFlag_ValidFile(t *testing.T) {
	filePath := "test.txt"
	expected := "339486 test.txt"
	actual := getFormattedOutput([]string{"-m", filePath})
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormattedOutput_CharsAndLinesFlag_ValidFile(t *testing.T) {
	filePath := "test.txt"
	expected := "7189 339486 test.txt"
	actual := getFormattedOutput([]string{"-ml", filePath})
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
