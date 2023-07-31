package main

import (
	"testing"
)

func TestCleanData(t *testing.T) {
	// Test with IgnoreCase set to true
	text := "Hello World"
	expected := "hello world"
	flags := Flags{IgnoreCase: true}
	cleanData(&text, flags)
	if text != expected {
		t.Errorf("Expected %q, but got %q", expected, text)
	}

	// Test with IgnorePunctuation set to true
	text = "Hello! How are you?"
	expected = "Hello How are you"
	flags = Flags{IgnorePunctuation: true}
	cleanData(&text, flags)
	if text != expected {
		t.Errorf("Expected %q, but got %q", expected, text)
	}

	// Test with both IgnoreCase and IgnorePunctuation set to true
	text = "Hello, World!"
	expected = "hello world"
	flags = Flags{IgnoreCase: true, IgnorePunctuation: true}
	cleanData(&text, flags)
	if text != expected {
		t.Errorf("Expected %q, but got %q", expected, text)
	}
}

func TestGrepImage(t *testing.T) {
	// Test a regular match
	text := "his name is john. he loves his job and his family."
	flags := Flags{Padding: 25}
	pattern := `\bhis\b`
	filename := "test.png"
	expected := "MATCH test.png:\nhis name is john. he loves h\n\nMATCH test.png:\ns name is john. he loves his job and his family\n\nMATCH test.png:\nhn. he loves his job and his family\n\n"
	result, err := grepImage(text, flags, pattern, filename)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Test an inverted match
	text = "his name is john. he loves his job and his family."
	flags = Flags{Invert: true}
	pattern = `\bhis\b`
	filename = "test.png"
	expected = "test.png without (\\bhis\\b):\n name is john. he loves  job and  family.\n\n"
	result, err = grepImage(text, flags, pattern, filename)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Test an inverted, case ignored match
	text = "His name is John. He loves his job and his family."
	flags = Flags{Invert: true}
	pattern = `(?i)\bhis\b`
	filename = "test.png"
	expected = "test.png without ((?i)\\bhis\\b):\n name is John. He loves  job and  family.\n\n"
	result, err = grepImage(text, flags, pattern, filename)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Test when no matches are found
	text = "This is a test."
	flags = Flags{}
	pattern = "notfound"
	filename = "test.png"
	_, err = grepImage(text, flags, pattern, filename)
	if err == nil {
		t.Errorf("Expected error for no matches, but got none.")
	}
}
