package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{[]string{"-", "-f", "1", "sample.tsv"}, "f0\n0\n5\n10\n15\n20\n"},
		{[]string{"-", "-f", "2", "sample.tsv"}, "f1\n1\n6\n11\n16\n21\n"},
		{[]string{"-", "-f", "1,2", "sample.tsv"}, "f0\tf1\n0\t1\n5\t6\n10\t11\n15\t16\n20\t21\n"},
		{[]string{"-", "-f", "1,2", "-d", "	", "sample.tsv"}, "f0\tf1\n0\t1\n5\t6\n10\t11\n15\t16\n20\t21\n"},
		{[]string{"-", "-f", "1 3", "-d", "\t", "sample.tsv"}, "f0\tf2\n0\t2\n5\t7\n10\t12\n15\t17\n20\t22\n"},
		{[]string{"-", "-f", "1 3", "-d", ",", "sample.csv"}, "f0,f2\n0,2\n5,7\n10,12\n15,17\n20,22\n"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test: %v", i), func(t *testing.T) {
			os.Args = test.args
			actual := captureOutput(main)
			if actual != test.expected {
				t.Errorf("main() = %v, want %v", actual, test.expected)
			}
		})
	}
}

func captureOutput(f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
