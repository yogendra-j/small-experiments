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

func TestMain_stdin(t *testing.T) {
	os.Args = []string{"-", "-f", "1"}
	input := "f0\tf1\tf2\tf3\tf4\n0\t1\t2\t3\t4\n5\t6\t7\t8\t9\n10\t11\t12\t13\t14\n15\t16\t17\t18\t19\n20\t21\t22\t23\t24\n"

	r, w, _ := os.Pipe()
	w.WriteString(input)

	originalStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = originalStdin }()

	w.Close()

	actual := captureOutput(main)
	expected := "f0\n0\n5\n10\n15\n20\n"
	if actual != expected {
		t.Errorf("got = `%v`, want `%v`", actual, expected)
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
