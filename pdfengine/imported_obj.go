package pdfengine

import (
	"io"
)

// importedObj : imported object
type importedObj struct { //impl iObj
	Data string
}

func (c *importedObj) init(funcGetRoot func() *PdfEngine) {

}

func (c *importedObj) getType() string {
	return "Imported"
}

func (c *importedObj) write(w io.Writer, objID int) error {
	if c != nil {
		io.WriteString(w, c.Data)
	}
	return nil
}
