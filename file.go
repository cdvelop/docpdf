package docpdf

import (
	"bytes"
	"io"
	"strconv"

	"github.com/cdvelop/docpdf/errs"
)

type (
	countingWriter struct {
		offset int64
		writer io.Writer
	}
)

// fileWriter is a function type for writing PDF data to a file
type fileWriter func(filename string, data []byte) error

// SetFileWriter sets a custom function for writing PDF files
func (gp *pdfEngine) SetFileWriter(writer fileWriter) {
	gp.fileWriter = writer
}

// WritePdf writes the PDF file to the specified path using the configured FileWriter
func (gp *pdfEngine) WritePdf(pdfPath string) error {
	// Get PDF bytes
	data := gp.GetBytesPdf()

	if len(data) == 0 {
		return errs.ErrEmptyPdf
	}

	// Default behavior
	return gp.fileWriter(pdfPath, data)
}

// WritePdf implements the io.WriterTo interface and can
// be used to stream the PDF as it is compiled to an io.Writer.
func (gp *pdfEngine) WriteTo(w io.Writer) (n int64, err error) {
	return gp.compilePdf(w)
}

// Write streams the pdf as it is compiled to an io.Writer
//
// Deprecated: use the WritePdf method instead.
func (gp *pdfEngine) Write(w io.Writer) error {
	_, err := gp.compilePdf(w)
	return err
}

func (gp *pdfEngine) Read(p []byte) (int, error) {
	if gp.buf.Len() == 0 && gp.buf.Cap() == 0 {
		if _, err := gp.compilePdf(&gp.buf); err != nil {
			return 0, err
		}
	}
	return gp.buf.Read(p)
}

// Close clears the gopdf buffer.
func (gp *pdfEngine) Close() error {
	gp.buf = bytes.Buffer{}
	return nil
}

func (gp *pdfEngine) compilePdf(w io.Writer) (n int64, err error) {
	gp.prepare()
	err = gp.Close()
	if err != nil {
		return 0, err
	}
	max := len(gp.pdfObjs)
	writer := newCountingWriter(w)
	io.WriteString(writer, "%PDF-1.7\n%����\n\n")
	linelens := make([]int64, max)
	i := 0

	for i < max {
		objID := i + 1
		linelens[i] = writer.offset
		pdfObj := gp.pdfObjs[i]
		io.WriteString(writer, strconv.Itoa(objID))
		io.WriteString(writer, " 0 obj\n")
		pdfObj.write(writer, objID)
		io.WriteString(writer, "endobj\n\n")
		i++
	}
	gp.xref(writer, writer.offset, linelens, i)
	return writer.offset, nil
}

func newCountingWriter(w io.Writer) *countingWriter {
	return &countingWriter{writer: w}
}

func (cw *countingWriter) Write(b []byte) (int, error) {
	n, err := cw.writer.Write(b)
	cw.offset += int64(n)
	return n, err
}

// GetBytesPdfReturnErr : get bytes of pdf file
func (gp *pdfEngine) GetBytesPdfReturnErr() ([]byte, error) {
	err := gp.Close()
	if err != nil {
		return nil, err
	}
	_, err = gp.compilePdf(&gp.buf)
	return gp.buf.Bytes(), err
}

// GetBytesPdf : get bytes of pdf file
func (gp *pdfEngine) GetBytesPdf() []byte {
	b, err := gp.GetBytesPdfReturnErr()
	if err != nil {
		gp.log(err)
	}
	return b
}
