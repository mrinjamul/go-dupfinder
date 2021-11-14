package app

import (
	"testing"
)

// TestGetFileName tests the GetFileName function
func TestGetFileName(t *testing.T) {
	fileName := GetFileName("/path/to/file.ext")
	if fileName != "file.ext" {
		t.Errorf("Expected fileName to be 'file.ext', got: %s", fileName)
	}
}

// TestContainsString tests the ContainsString function
func TestContainsString(t *testing.T) {
	if !ContainsString([]string{"a"}, "a") {
		t.Errorf("Expected ContainsString to return true")
	}
	if ContainsString([]string{"a"}, "b") {
		t.Errorf("Expected ContainsString to return false")
	}
}
