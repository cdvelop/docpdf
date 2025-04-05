package docpdf

import (
	"io"
)

// importedObj : imported object
type importedObj struct { //impl iObj
	Data string
}

func (c *importedObj) init(funcGetRoot func() *pdfEngine) {

}

func (c *importedObj) getType() string {
	return "Imported"
}

func (c *importedObj) write(w writer, objID int) error {
	if c != nil {
		io.WriteString(w, c.Data)
	}
	return nil
}
