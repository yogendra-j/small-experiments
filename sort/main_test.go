package main

import (
	"fmt"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		lines  []string
		sorted []string
	}{
		{
			lines:  []string{"a", "b", "c"},
			sorted: []string{"a", "b", "c"},
		},
		{
			lines:  []string{"c", "b", "a"},
			sorted: []string{"a", "b", "c"},
		},
		{
			lines:  []string{"Is", "this", "sorted?"},
			sorted: []string{"Is", "sorted?", "this"},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Sorting case-%d", i+1), func(t *testing.T) {
			sorted := sort(test.lines)
			for i, line := range sorted {
				if line != test.sorted[i] {
					t.Errorf("Expected: %v Got: %v", test.sorted, sorted)
					t.FailNow()
				}
			}
		})
	}
}
