// filepath: c:\Users\Cesar\Packages\Internal\docpdf\pdfengine\pdfEngineExtra.go
package pdfengine

import "bytes"

// GetBuf returns the buffer
func (gp *PdfEngine) GetBuf() *bytes.Buffer {
	return &gp.buf
}

// GetEncryptionObjID returns the encryption object ID
func (gp *PdfEngine) GetEncryptionObjID() int {
	return gp.encryptionObjID
}

// GetIndexOfContent returns the index of content
func (gp *PdfEngine) GetIndexOfContent() int {
	return gp.indexOfContent
}

// GetIndexEncodingObjFonts returns the index encoding object fonts
func (gp *PdfEngine) GetIndexEncodingObjFonts() []int {
	return gp.indexEncodingObjFonts
}
