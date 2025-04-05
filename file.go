package docpdf

import (
	"bytes"
	"io"
	"strconv"
)

type (
	countingWriter struct {
		offset int64
		writer writer
	}
)

// WritePdfFile : write pdf file if err occurred it will be logged
func (gp *pdfEngine) WritePdfFile() {
	data, err := gp.GetBytesPdfReturnErr()
	if err != nil {
		gp.log(err)
		return
	}

	err = gp.fileWrite.FileWrite(data)
	if err != nil {
		gp.log(err)
	}
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

func (gp *pdfEngine) compilePdf(w writer) (n int64, err error) {
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

func newCountingWriter(w writer) *countingWriter {
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
