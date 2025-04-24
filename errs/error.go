package errs

import (
	"bytes"
	"strconv"
)

// errMessage representa un error simple para TinyGo
type errMessage struct {
	message string
}

func (e *errMessage) Error() string {
	return e.message
}

func New(args ...any) *errMessage {

	var out bytes.Buffer
	var space string

	// Check if we have at least one argument
	if len(args) == 0 {
		return &errMessage{}
	}

	// Process remaining arguments
	for argNumber, arg := range args {
		switch v := arg.(type) {
		case string:
			if v == "" {
				continue
			}
			out.WriteString(space + v)
		case []string:
			for _, s := range v {
				if s == "" {
					continue
				}
				out.WriteString(space + s)
				space = " "
			}
		// Other cases remain the same
		case rune:
			if v == ':' {
				out.WriteString(":")
				continue
			}
			out.WriteString(space + string(v))
		case int:
			out.WriteString(space + strconv.Itoa(v))
		case float64:
			out.WriteString(space + strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			out.WriteString(space + strconv.FormatBool(v))
		case error:
			out.WriteString(space + v.Error())
		default:
			out.WriteString(space + "error not supported arg number: " + strconv.Itoa(argNumber))
		}
		space = " "
	}

	return &errMessage{
		message: out.String(),
	}
}
