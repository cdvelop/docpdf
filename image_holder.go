package docpdf

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// imageHolder hold image data
type imageHolder interface {
	ID() string
	io.Reader
}

// imageHolderByBytes create imageHolder by []byte
func imageHolderByBytes(b []byte) (imageHolder, error) {
	return newImageBuff(b)
}

// imageHolderByPath create imageHolder by image path
func imageHolderByPath(path string) (imageHolder, error) {
	return newImageBuffByPath(path)
}

// imageHolderByReader create imageHolder by io.Reader
func imageHolderByReader(r io.Reader) (imageHolder, error) {
	return newImageBuffByReader(r)
}

// imageBuff image holder (impl imageHolder)
type imageBuff struct {
	id string
	bytes.Buffer
}

func newImageBuff(b []byte) (*imageBuff, error) {
	h := md5.New()
	_, err := h.Write(b)
	if err != nil {
		return nil, err
	}
	var i imageBuff
	i.id = fmt.Sprintf("%x", h.Sum(nil))
	i.Write(b)
	return &i, nil
}

func newImageBuffByPath(path string) (*imageBuff, error) {
	var i imageBuff
	i.id = path
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	i.Write(b)
	return &i, nil
}

func newImageBuffByReader(r io.Reader) (*imageBuff, error) {

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	_, err = h.Write(b)
	if err != nil {
		return nil, err
	}
	var i imageBuff
	i.id = fmt.Sprintf("%x", h.Sum(nil))
	i.Write(b)
	return &i, nil
}

func (i *imageBuff) ID() string {
	return i.id
}
