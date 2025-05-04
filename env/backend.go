//go:build !wasm
// +build !wasm

package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// SetupDefaultLogger configures the default logger for backend environments
func SetupDefaultLogger() func(a ...any) {
	return func(a ...any) {
		fmt.Println(a...)
	}
}

func SetupDefaultFileWriter() func(filename string, data []byte) error {
	return func(filename string, data []byte) error {
		return os.WriteFile(filename, data, 0644)
	}
}

// FileExists checks if a file exists and returns its contents.
// Accepts string (path) or []byte (content).
// For paths, verifies existence and reads the file.
// For []byte, returns the provided content directly.
func FileExists(pathOrContent any) ([]byte, error) {
	switch v := pathOrContent.(type) {
	case string:
		// Handle path string
		// Get absolute path
		absolutePath, err := filepath.Abs(v)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path: %w", err)
		}

		// Check if file exists and is not a directory
		info, err := os.Stat(absolutePath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("file does not exist: %s", absolutePath)
			}
			return nil, fmt.Errorf("failed to stat file: %w", err)
		}

		if info.IsDir() {
			return nil, fmt.Errorf("path is a directory, not a file: %s", absolutePath)
		}

		// Read file content
		content, err := os.ReadFile(absolutePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		return content, nil

	case []byte:
		// If content is provided directly, return it as is
		return v, nil

	default:
		return nil, fmt.Errorf("unsupported type: %T, expected string or []byte", pathOrContent)
	}
}

// GetSize returns the size of a file or byte slice.
// For backend: uses os.Stat for files, len() for byte slices.
func GetSize(pathOrContent any) (int64, error) {
	switch v := pathOrContent.(type) {
	case string:
		// Assume it's a file path
		stat, err := os.Stat(v)
		if err != nil {
			return -1, err
		}
		return stat.Size(), nil
	case []byte:
		// It's already content
		return int64(len(v)), nil
	default:
		return -1, fmt.Errorf("unsupported type for GetSize: %T", pathOrContent)
	}
}
