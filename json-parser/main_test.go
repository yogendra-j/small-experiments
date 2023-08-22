package main

import (
	"bufio"
	"bytes"
	"testing"
)

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
		{"{,  }"},
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
		{`{,"key": "value" } `, false},
		{`{"key ": "va
		lue" } `, false},
		{`{"key
		": "value" } `, false},
		{`{"key"
		
		: 
			 "value"
			 , "k": "v" }  `, true},
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

func TestJsonParser_WithMultipleStringKeyValue(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"key": "value", "key2" : "val 2  "}`, true},
		{`{"key": "value","key2" : "val 2  " }`, true},
		{`{
			"abc": "def" , 
		"ghi" : "jkl  " 
		} `, true},
		{`{"key": "value", }  `, false},
		{`{
			"key": "value",
		"key2" : "val 2  " ,
		} `, false},
		{`{
			"key": "value",
		key2 : "val 2  "
		} `, false},
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

func TestJsonParser_WithSingleBoolKeyValue(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"key": true}`, true},
		{`{"key": true }  `, true},
		{`{"key ": true
		} `, true},
		{`{,"key": true } `, false},
		{`{"key
		": true } `, false},
		{`{"key"  : true" `, false},
		{`{"key": true `, false},
		{`{key": true
		"}`, false},
		{`{"key: true"}`, false},
		{`{"key" true"}`, false},
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

func TestJsonParser_WithSingleNullValue(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"key": null}`, true},
		{`{"key": null }  `, true},
		{`{"key ": null
		} `, true},
		{`{,"key": null } `, false},
		{`{"key
		": null } `, false},
		{`{"key"  : null" `, false},
		{`{"key": null `, false},
		{`{key": null
		"}`, false},
		{`{"key: null"}`, false},
		{`{"key" null"}`, false},
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

func TestJsonParser_WithSingleIntegerValue(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"key": 101}`, true},
		{`{"key": 101 }  `, true},
		{`{"key ": 101
		} `, true},
		{`{"key": -101}`, true},
		{`{"key": -101 }  `, true},
		{`{"key ": -101
		} `, true},
		{`{"key": --101 }  `, false},
		{`{"key": -10-1 }  `, false},
		{`{,"key": 101 } `, false},
		{`{"key
		": 101 } `, false},
		{`{"key"  : 101" `, false},
		{`{"key": 101 `, false},
		{`{key": -101
		"}`, false},
		{`{"key: 101"}`, false},
		{`{"key" 101"}`, false},
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

func TestJsonParser_WithSingleFractionValue(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"key": 101.254}`, true},
		{`{"key": 101.0 }  `, true},
		{`{"key ": 101.55555
		} `, true},
		{`{"key": 101.0. }  `, false},
		{`{"key": 101.0.2 }  `, false},
		{`{"key": 101. }  `, false},
		{`{,"key": .5 } `, false},
		{`{,"key": . } `, false},
		{`{"key
		": 101.2 } `, false},
		{`{"key"  : 101.2.2" `, false},
		{`{"key": 101.2 `, false},
		{`{key": 101.2
		"}`, false},
		{`{"key: 101.2"}`, false},
		{`{"key" 101.2"}`, false},
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

func TestJsonParser_WithNumBoolStringNull(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{
			"key1": true,
			"key2": false,
			"key3": null,
			"key4": "value",
			"key5": 101
		  }`, true},
		{`{
			"key1": True,
			"key2": false,
			"key3": null,
			"key4": "value",
			"key5": 101
		  }`, false},
		{`{
			"key1": true,
			"key2": false,
			"key3": nu ll,
			"key4": "value",
			"key5": 101
		  }`, false},
		{`{
			"key1": true,
			"key2": falSe,
			"key3": null,
			"key4": "value",
			"key5": 101
		  } `, false},
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
