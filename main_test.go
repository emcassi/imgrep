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

func TestGrep(t *testing.T) {

	// Test a regular match
	text := "his name is john. he loves his job and his family."
	flags := Flags{Padding: 0}
	pattern := `\bhis\b`
	filename := "test.png"
	expected := []string{
		"his",
		"his",
		"his",
	}

	result, err := grep(text, flags, pattern, filename)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if !areListsEqual(expected, result) {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Test with invert flag set to true
	text = "This is a test."
	flags = Flags{Invert: true}
	pattern = "test"
	filename = "test.png"
	expected = []string{
		"This is a .",
	}
	result, err = grep(text, flags, pattern, filename)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if !areListsEqual(expected, result) {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Test padding functionality
	text = "This is a test. Padding test."
	flags = Flags{Padding: 6}
	pattern = "test"
	filename = "test.png"
	expected = []string{
		" is a test. Padd",
		"dding test",
	}
	result, err = grep(text, flags, pattern, filename)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if !areListsEqual(expected, result) {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Test when no matches are found
	text = "This is a test."
	flags = Flags{}
	pattern = "notfound"
	filename = "test.png"
	_, err = grep(text, flags, pattern, filename)
	if err == nil {
		t.Errorf("Expected error for no matches, but got none.")
	}
}

func areListsEqual(a, b []string) bool {
	
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
