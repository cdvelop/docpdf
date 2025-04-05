package docpdf

import (
	"fmt"
	"io"
)

// outlinesObj : outlines dictionary
type outlinesObj struct { //impl iObj
	getRoot func() *pdfEngine

	index   int
	first   int
	last    int
	count   int
	lastObj *outlineObj
}

func (o *outlinesObj) init(funcGetRoot func() *pdfEngine) {
	o.getRoot = funcGetRoot
	o.first = -1
	o.last = -1
}

func (o *outlinesObj) getType() string {
	return "Outlines"
}

func (o *outlinesObj) write(w writer, objID int) error {
	content := "<<\n"
	content += fmt.Sprintf("\t/Type /%s\n", o.getType())
	content += fmt.Sprintf("\t/Count %d\n", o.count)

	if o.first >= 0 {
		content += fmt.Sprintf("\t/First %d 0 R\n", o.first)
	}

	if o.last >= 0 {
		content += fmt.Sprintf("\t/Last %d 0 R\n", o.last)
	}

	content += ">>\n"

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

	return nil
}

func (o *outlinesObj) SetIndexObjOutlines(index int) {
	o.index = index
}

func (o *outlinesObj) AddOutline(dest int, title string) {
	oo := &outlineObj{title: title, dest: dest, parent: o.index, prev: o.last, next: -1}
	o.last = o.getRoot().addObj(oo) + 1
	if o.first <= 0 {
		o.first = o.last
	}
	if o.lastObj != nil {
		o.lastObj.next = o.last
	}
	o.lastObj = oo
	o.count++
}

// AddOutlinesWithPosition add outlines with position
func (o *outlinesObj) AddOutlinesWithPosition(dest int, title string, y float64) *outlineObj {
	oo := &outlineObj{title: title, dest: dest, parent: o.index, prev: o.last, next: -1, height: y}
	o.last = o.getRoot().addObj(oo) + 1
	if o.first <= 0 {
		o.first = o.last
	}
	if o.lastObj != nil {
		o.lastObj.next = o.last
	}
	o.lastObj = oo
	o.count++
	oo.index = o.last
	return oo
}

func (o *outlinesObj) Count() int {
	return o.count
}

// outlineObj include attribute of outline
type outlineObj struct { //impl iObj
	title  string
	index  int
	dest   int
	parent int
	prev   int
	next   int
	first  int
	last   int
	height float64
}

func (o *outlineObj) init(funcGetRoot func() *pdfEngine) {
}

func (o *outlineObj) SetFirst(first int) {
	o.first = first
}

func (o *outlineObj) SetLast(last int) {
	o.last = last
}

func (o *outlineObj) SetPrev(prev int) {
	o.prev = prev
}

func (o *outlineObj) SetNext(next int) {
	o.next = next
}

func (o *outlineObj) SetParent(parent int) {
	o.parent = parent
}

func (o *outlineObj) GetIndex() int {
	return o.index
}

func (o *outlineObj) getType() string {
	return "Outline"
}

//func (o *outlineObj) write(w writer, objID int) error {
//	io.WriteString(w, "<<\n")
//	fmt.Fprintf(w, "  /Parent %d 0 R\n", o.parent)
//	if o.prev >= 0 {
//		fmt.Fprintf(w, "  /Prev %d 0 R\n", o.prev)
//	}
//	if o.next >= 0 {
//		fmt.Fprintf(w, "  /Next %d 0 R\n", o.next)
//	}
//	fmt.Fprintf(w, "  /Dest [ %d 0 R /XYZ null null null ]\n", o.dest)
//	fmt.Fprintf(w, "  /Title <FEFF%s>\n", encodeUtf8(o.title))
//	io.WriteString(w, ">>\n")
//	return nil
//}

func (o *outlineObj) write(w writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Parent %d 0 R\n", o.parent)
	if o.prev >= 0 {
		fmt.Fprintf(w, "  /Prev %d 0 R\n", o.prev)
	}
	if o.next >= 0 {
		fmt.Fprintf(w, "  /Next %d 0 R\n", o.next)
	}
	if o.first > 0 {
		fmt.Fprintf(w, "  /First %d 0 R\n", o.first)
	}
	if o.last > 0 {
		fmt.Fprintf(w, "  /Last %d 0 R\n", o.last)
	}
	fmt.Fprintf(w, "  /Dest [ %d 0 R /XYZ 90 %f 0 ]\n", o.dest, o.height)
	fmt.Fprintf(w, "  /Title <FEFF%s>\n", encodeUtf8(o.title))
	io.WriteString(w, ">>\n")
	return nil
}

// outlineNode is a node of outline
type outlineNode struct {
	Obj      *outlineObj
	Children []*outlineNode
}

// outlineNodes are all nodes of outline
type outlineNodes []*outlineNode

// Parse parse outline nodes
func (objs outlineNodes) Parse() {
	for i, obj := range objs {
		if i == 0 {
			obj.Obj.SetPrev(-1)
		} else {
			obj.Obj.SetNext(objs[i-1].Obj.GetIndex())
		}
		if i == len(objs)-1 {
			obj.Obj.SetNext(-1)
		} else {
			obj.Obj.SetNext(objs[i+1].Obj.GetIndex())
		}
		obj.Parse()
	}

}

// Parse parse outline
func (obj outlineNode) Parse() {
	if obj.Children == nil || len(obj.Children) == 0 {
		return
	}
	for i, children := range obj.Children {
		if i == 0 {
			obj.Obj.SetFirst(children.Obj.GetIndex())
			children.Obj.SetPrev(-1)
		}
		if i == len(obj.Children)-1 {
			obj.Obj.SetLast(children.Obj.GetIndex())
			children.Obj.SetNext(-1)
		}
		if i != 0 {
			children.Obj.SetPrev(obj.Children[i-1].Obj.GetIndex())
		}
		if i != len(obj.Children)-1 {
			children.Obj.SetNext(obj.Children[i+1].Obj.GetIndex())
		}
		children.Obj.SetParent(obj.Obj.GetIndex())
		children.Parse()
	}
}
