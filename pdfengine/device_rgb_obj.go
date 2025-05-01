package pdfengine

import (
	"fmt"
	"io"
)

// deviceRGBObj  DeviceRGB
type deviceRGBObj struct {
	data    []byte
	getRoot func() *PdfEngine
}

func (d *deviceRGBObj) init(funcGetRoot func() *PdfEngine) {
	d.getRoot = funcGetRoot
}

func (d *deviceRGBObj) protection() *pdfProtection {
	return d.getRoot().protection()
}

func (d *deviceRGBObj) getType() string {
	return "devicergb"
}

// สร้าง ข้อมูลใน pdf
func (d *deviceRGBObj) write(w io.Writer, objID int) error {

	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "/Length %d\n", len(d.data))
	io.WriteString(w, ">>\n")
	io.WriteString(w, "stream\n")
	if d.protection() != nil {
		tmp, err := rc4Cip(d.protection().objectkey(objID), d.data)
		if err != nil {
			return err
		}
		w.Write(tmp)
		io.WriteString(w, "\n")
	} else {
		w.Write(d.data)
	}
	io.WriteString(w, "endstream\n")

	return nil
}
