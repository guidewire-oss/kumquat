package store_test

import (
	"kumquat/store"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFile(t *testing.T) {
	// Setup a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // nolint:errcheck
	// Test data
	testFileName := "testfile.txt"
	testContent := "Hello, world!"
	// Test successful write
	err = store.WriteToFile(testFileName, tempDir, testContent)
	if err != nil {
		t.Errorf("Failed to write to file: %v", err)
	} else {
		// Verify file content
		fullPath := filepath.Join(tempDir, testFileName)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read from file: %v", err)
		} else if string(content) != testContent {
			t.Errorf("Content mismatch: expected %v, got %v", testContent, string(content))
		}
	}

}
