//go:build !wasm
// +build !wasm

package env

import (
	"fmt"
	"os"
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
