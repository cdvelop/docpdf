package docpdf

import (
	"fmt"
	"io"
	"strings"
)

// encryptionObj  encryption object res
type encryptionObj struct {
	uValue []byte //U entry in pdf document
	oValue []byte //O entry in pdf document
	pValue int    //P entry in pdf document
}

func (e *encryptionObj) init(func() *pdfEngine) {

}

func (e *encryptionObj) getType() string {
	return "Encryption"
}

func (e *encryptionObj) write(w writer, objID int) error {
	io.WriteString(w, "<<\n")
	io.WriteString(w, "/Filter /Standard\n")
	io.WriteString(w, "/V 1\n")
	io.WriteString(w, "/R 2\n")
	fmt.Fprintf(w, "/O (%s)\n", e.escape(e.oValue))
	fmt.Fprintf(w, "/U (%s)\n", e.escape(e.uValue))
	fmt.Fprintf(w, "/P %d\n", e.pValue)
	io.WriteString(w, ">>\n")
	return nil
}

func (e *encryptionObj) escape(b []byte) string {
	s := string(b)
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "(", "\\(", -1)
	s = strings.Replace(s, ")", "\\)", -1)
	s = strings.Replace(s, "\r", "\\r", -1)
	return s
}
