package store

import (
	"os"
	"path/filepath"
)

// writeFile takes a file name, file path, and a string content, then writes the content to the specified file.
func WriteToFile(fileName, filePath, content string) error {
	// Combine the file path and file name to get the full file path
	fullPath := filepath.Join(filePath, fileName)
	// Create the directory path if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	// Convert the string content to a byte slice, as WriteFile expects []byte
	byteContent := []byte(content)

	// Write the content to the file
	err := os.WriteFile(fullPath, byteContent, 0644)
	if err != nil {
		return err
	}
	return nil
}
