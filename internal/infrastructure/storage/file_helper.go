package storage

import (
	"bufio"
	"os"
	"path/filepath"
)

// FileHelper provides methods for file I/O
type FileHelper struct {
	BaseDir string
}

func NewFileHelper(baseDir string) *FileHelper {
	return &FileHelper{BaseDir: baseDir}
}

// ReadFileAsync reads lines from a file
func (h *FileHelper) ReadFileAsync(directory, fileName string) ([]string, error) {
	filePath := filepath.Join(h.BaseDir, directory, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Fallback to relative path if not found in base dir
		filePath = filepath.Join(directory, fileName)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// WriteFileNewContents writes content to a file, overwriting it
func (h *FileHelper) WriteFileNewContents(contents []string, directory, fileName string) error {
	filePath := filepath.Join(h.BaseDir, directory, fileName)
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range contents {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}

// WriteFileAppend appends content to a file
func (h *FileHelper) WriteFileAppend(contents []string, directory, fileName string) error {
	filePath := filepath.Join(h.BaseDir, directory, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range contents {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}
