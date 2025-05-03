package pdfengine

import (
	"fmt"
	"strings"
)

type annotObj struct {
	linkOption
	GetRoot func() *PdfEngine
}

func (o annotObj) Init(f func() *PdfEngine) {
}

func (o annotObj) GetType() string {
	return "Annot"
}

func (o annotObj) Write(w Writer, objID int) error {
	if o.url != "" {
		return o.writeExternalLink(w, o.linkOption, objID)
	} else {
		return o.writeInternalLink(w, o.linkOption)
	}
}

func (o annotObj) writeExternalLink(w Writer, l linkOption, objID int) error {
	protection := o.GetRoot().protection()
	url := l.url
	if protection != nil {
		tmp, err := rc4Cip(protection.objectkey(objID), []byte(url))
		if err != nil {
			return err
		}
		url = string(tmp)
	}
	url = strings.Replace(url, "\\", "\\\\", -1)
	url = strings.Replace(url, "(", "\\(", -1)
	url = strings.Replace(url, ")", "\\)", -1)
	url = strings.Replace(url, "\r", "\\r", -1)

	_, err := fmt.Fprintf(w, "<</Type /Annot /Subtype /Link /canvas.Rect [%.2f %.2f %.2f %.2f] /Border [0 0 0] /A <</S /URI /URI (%s)>>>>\n",
		l.x, l.y, l.x+l.w, l.y-l.h, url)
	return err
}

func (o annotObj) writeInternalLink(w Writer, l linkOption) error {
	a, ok := o.GetRoot().anchors[l.anchor]
	if !ok {
		return nil
	}
	_, err := fmt.Fprintf(w, "<</Type /Annot /Subtype /Link /canvas.Rect [%.2f %.2f %.2f %.2f] /Border [0 0 0] /Dest [%d 0 R /XYZ 0 %.2f null]>>\n",
		l.x, l.y, l.x+l.w, l.y-l.h, a.page+1, a.y)
	return err
}
