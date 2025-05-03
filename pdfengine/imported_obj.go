package pdfengine

import (
	"io"
)

// importedObj : imported object
type importedObj struct { //impl iObj
	Data string
}

func (c *importedObj) Init(funcGetRoot func() *PdfEngine) {

}

func (c *importedObj) GetType() string {
	return "Imported"
}

func (c *importedObj) Write(w Writer, objID int) error {
	if c != nil {
		io.WriteString(w, c.Data)
	}
	return nil
}
