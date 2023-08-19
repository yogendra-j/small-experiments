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
