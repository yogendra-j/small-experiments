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
