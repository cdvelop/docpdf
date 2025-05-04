package env

// FileWriter is a function type for writing PDF data to a file
type FileWriter func(filename string, data []byte) error
