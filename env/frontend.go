//go:build wasm
// +build wasm

package env

import (
	"syscall/js"

	"github.com/cdvelop/docpdf/errs"
	"github.com/cdvelop/docpdf/utils"
)

// SetupDefaultLogger configures the default logger for frontend environments
func SetupDefaultLogger() func(a ...any) {
	return func(a ...any) {
		// Use console.log in browser environment
		args := make([]any, len(a))
		for i, arg := range a {
			args[i] = js.ValueOf(utils.AnyToString(arg))
		}
		js.Global().Get("console").Call("log", args...)
	}
}

func SetupDefaultFileWriter() func(filename string, data []byte) error {
	return func(filename string, data []byte) error {
		return errs.New("file writing not implemented in frontend")
	}

}
