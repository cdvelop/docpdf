package docpdf

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"

	"math"
	"os"
	"strconv"
)

type pdfWriter struct {
	f       *os.File
	w       *bufio.Writer
	r       *pdfReader
	k       float64
	tpls    []*pdfTemplate
	m       int
	n       int
	offsets map[int]int
	offset  int
	result  map[int]string
	// Keep track of which objects have already been written
	obj_stack       map[int]*pdfValue
	don_obj_stack   map[int]*pdfValue
	written_objs    map[*pdfObjectId][]byte
	written_obj_pos map[*pdfObjectId]map[int]string
	current_obj     *pdfObject
	current_obj_id  int
	tpl_id_offset   int
	use_hash        bool

	log func(...any) // log function for debugging
}

type pdfObjectId struct {
	id   int
	hash string
}

type pdfObject struct {
	id     *pdfObjectId
	buffer *bytes.Buffer
}

func (this *pdfWriter) SetTplIdOffset(n int) {
	this.tpl_id_offset = n
}

func (this *pdfWriter) Init() {
	this.k = 1
	this.obj_stack = make(map[int]*pdfValue, 0)
	this.don_obj_stack = make(map[int]*pdfValue, 0)
	this.tpls = make([]*pdfTemplate, 0)
	this.written_objs = make(map[*pdfObjectId][]byte, 0)
	this.written_obj_pos = make(map[*pdfObjectId]map[int]string, 0)
	this.current_obj = new(pdfObject)
}

func (this *pdfWriter) SetUseHash(b bool) {
	this.use_hash = b
}

func (this *pdfWriter) SetNextObjectID(id int) {
	this.n = id - 1
}

func newPdfWriter(filename string, log func(a ...any)) (*pdfWriter, error) {
	writer := &pdfWriter{
		log: log,
	}
	writer.Init()

	if filename != "" {
		var err error
		f, err := os.Create(filename)
		if err != nil {
			return nil, newErr(err, "Unable to create filename: ", filename)
		}
		writer.f = f
		writer.w = bufio.NewWriter(f)
	}
	return writer, nil
}

// Done with parsing.  Now, create templates.
type pdfTemplate struct {
	Id        int
	Reader    *pdfReader
	Resources *pdfValue
	Buffer    string
	box       map[string]float64
	Boxes     map[string]map[string]float64
	X         float64
	Y         float64
	W         float64
	H         float64
	Rotation  int
	N         int
}

func (this *pdfWriter) GetImportedObjects() map[*pdfObjectId][]byte {
	return this.written_objs
}

// For each object (uniquely identified by a sha1 hash), return the positions
// of each hash within the object, to be replaced with pdf object ids (integers)
func (this *pdfWriter) GetImportedObjHashPos() map[*pdfObjectId]map[int]string {
	return this.written_obj_pos
}

func (this *pdfWriter) ClearImportedObjects() {
	this.written_objs = make(map[*pdfObjectId][]byte, 0)
}

// Create a pdfTemplate object from a page number (e.g. 1) and a boxName (e.g. MediaBox)
func (this *pdfWriter) ImportPage(reader *pdfReader, pageno int, boxName string) (int, error) {
	var err error

	// Set default scale to 1
	this.k = 1

	// Get all page boxes
	pageBoxes, err := reader.getPageBoxes(1, this.k)
	if err != nil {
		return -1, newErr(err, "Failed to get page boxes")
	}

	// If requested box name does not exist for this page, use an alternate box
	if _, ok := pageBoxes[boxName]; !ok {
		if boxName == "/BleedBox" || boxName == "/TrimBox" || boxName == "ArtBox" {
			boxName = "/CropBox"
		} else if boxName == "/CropBox" {
			boxName = "/MediaBox"
		}
	}

	// If the requested box name or an alternate box name cannot be found, trigger an error
	// TODO: Improve error handling
	if _, ok := pageBoxes[boxName]; !ok {
		return -1, newErr("box not found: " + boxName)
	}

	pageResources, err := reader.getPageResources(pageno)
	if err != nil {
		return -1, newErr(err, "Failed to get page resources")
	}

	content, err := reader.getContent(pageno)
	if err != nil {
		return -1, newErr(err, "Failed to get content")
	}

	// Set template values
	tpl := &pdfTemplate{}
	tpl.Reader = reader
	tpl.Resources = pageResources
	tpl.Buffer = content
	tpl.box = pageBoxes[boxName]
	tpl.Boxes = pageBoxes
	tpl.X = 0
	tpl.Y = 0
	tpl.W = tpl.box["w"]
	tpl.H = tpl.box["h"]

	// Set template rotation
	rotation, err := reader.getPageRotation(pageno)
	if err != nil {
		return -1, newErr(err, "Failed to get page rotation")
	}
	angle := rotation.Int % 360

	// Normalize angle
	if angle != 0 {
		steps := angle / 90
		w := tpl.W
		h := tpl.H

		if steps%2 == 0 {
			tpl.W = w
			tpl.H = h
		} else {
			tpl.W = h
			tpl.H = w
		}

		if angle < 0 {
			angle += 360
		}

		tpl.Rotation = angle * -1
	}

	this.tpls = append(this.tpls, tpl)

	// Return last template id
	return len(this.tpls) - 1, nil
}

// Create a new object and keep track of the offset for the xref table
func (this *pdfWriter) newObj(objId int, onlyNewObj bool) {
	if objId < 0 {
		this.n++
		objId = this.n
	}

	if !onlyNewObj {
		// set current object id integer
		this.current_obj_id = objId

		// Create new pdfObject and pdfObjectId
		this.current_obj = new(pdfObject)
		this.current_obj.buffer = new(bytes.Buffer)
		this.current_obj.id = new(pdfObjectId)
		this.current_obj.id.id = objId
		this.current_obj.id.hash = this.shaOfInt(objId)

		this.written_obj_pos[this.current_obj.id] = make(map[int]string, 0)
	}
}

func (this *pdfWriter) endObj() {
	this.out("endobj")

	this.written_objs[this.current_obj.id] = this.current_obj.buffer.Bytes()
	this.current_obj_id = -1
}

func (this *pdfWriter) shaOfInt(i int) string {
	hasher := sha1.New()
	hasher.Write([]byte(strconv.Itoa(i) + "-" + this.r.sourceFile))
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha
}

func (this *pdfWriter) outObjRef(objId int) {
	sha := this.shaOfInt(objId)

	// Keep track of object hash and position - to be replaced with actual object id (integer)
	this.written_obj_pos[this.current_obj.id][this.current_obj.buffer.Len()] = sha

	if this.use_hash {
		this.current_obj.buffer.WriteString(sha)
	} else {

		this.current_obj.buffer.WriteString(strconv.Itoa(objId))
	}
	this.current_obj.buffer.WriteString(" 0 R ")
}

// Output PDF data with a newline
func (this *pdfWriter) out(s string) {
	this.current_obj.buffer.WriteString(s)
	this.current_obj.buffer.WriteString("\n")
}

// Output PDF data
func (this *pdfWriter) straightOut(s string) {
	this.current_obj.buffer.WriteString(s)
}

// Output a pdfValue
func (this *pdfWriter) writeValue(value *pdfValue) {
	switch value.Type {
	case pdf_type_token:
		this.straightOut(value.Token + " ")
		break

	case pdf_type_numeric:
		this.straightOut(strconv.Itoa(value.Int) + " ")
		break

	case pdf_type_real:
		this.straightOut(strconv.FormatFloat(value.Real, 'f', -1, 64) + " ")
		break

	case pdf_type_array:
		this.straightOut("[")
		for i := 0; i < len(value.Array); i++ {
			this.writeValue(value.Array[i])
		}
		this.out("]")
		break

	case pdf_type_dictionary:
		this.straightOut("<<")
		for k, v := range value.Dictionary {
			this.straightOut(k + " ")
			this.writeValue(v)
		}
		this.straightOut(">>")
		break

	case pdf_type_objref:
		// An indirect object reference.  Fill the object stack if needed.
		// Check to see if object already exists on the don_obj_stack.
		if _, ok := this.don_obj_stack[value.Id]; !ok {
			this.newObj(-1, true)
			this.obj_stack[value.Id] = &pdfValue{Type: pdf_type_objref, Gen: value.Gen, Id: value.Id, NewId: this.n}
			this.don_obj_stack[value.Id] = &pdfValue{Type: pdf_type_objref, Gen: value.Gen, Id: value.Id, NewId: this.n}
		}

		// Get object ID from don_obj_stack
		objId := this.don_obj_stack[value.Id].NewId
		this.outObjRef(objId)
		//this.out(fmt.Sprintf("%d 0 R", objId))
		break

	case pdf_type_string:
		// A string
		this.straightOut("(" + value.String + ")")
		break

	case pdf_TYPE_stream:
		// A stream.  First, output the stream dictionary, then the stream data itself.
		this.writeValue(value.Value)
		this.out("stream")
		this.out(string(value.Stream.Bytes))
		this.out("endstream")
		break

	case pdf_type_hex:
		this.straightOut("<" + value.String + ">")
		break

	case pdf_type_boolean:
		if value.Bool {
			this.straightOut("true ")
		} else {
			this.straightOut("false ")
		}
		break

	case pdf_type_null:
		// The null object
		this.straightOut("null ")
		break
	}
}

// Output Form XObjects (1 for each template)
// returns a map of template names (e.g. /GOFPDITPL1) to pdfObjectId
func (this *pdfWriter) PutFormXobjects(reader *pdfReader) (map[string]*pdfObjectId, error) {
	// Set current reader
	this.r = reader

	var err error
	var result = make(map[string]*pdfObjectId, 0)

	compress := true
	filter := ""
	if compress {
		filter = "/Filter /FlateDecode "
	}

	for i := 0; i < len(this.tpls); i++ {
		tpl := this.tpls[i]
		if tpl == nil {
			return nil, newErr("Template is nil")
		}
		var p string
		if compress {
			var b bytes.Buffer
			w := zlib.NewWriter(&b)
			w.Write([]byte(tpl.Buffer))
			w.Close()

			p = b.String()
		} else {
			p = tpl.Buffer
		}

		// Create new PDF object
		this.newObj(-1, false)

		cN := this.n // remember current "n"

		tpl.N = this.n

		// Return xobject form name and object position
		pdfObjId := new(pdfObjectId)
		pdfObjId.id = cN
		pdfObjId.hash = this.shaOfInt(cN)
		result["/GOFPDITPL"+strconv.Itoa(i+this.tpl_id_offset)] = pdfObjId

		this.out("<<" + filter + "/Type /XObject")
		this.out("/Subtype /Form")
		this.out("/FormType 1")

		this.out("/BBox [" +
			strconv.FormatFloat(tpl.box["llx"]*this.k, 'f', 2, 64) + " " +
			strconv.FormatFloat(tpl.box["lly"]*this.k, 'f', 2, 64) + " " +
			strconv.FormatFloat((tpl.box["urx"]+tpl.X)*this.k, 'f', 2, 64) + " " +
			strconv.FormatFloat((tpl.box["ury"]-tpl.Y)*this.k, 'f', 2, 64) +
			"]")

		var c, s, tx, ty float64
		c = 1

		// Handle rotated pages
		if tpl.box != nil {
			tx = -tpl.box["llx"]
			ty = -tpl.box["lly"]

			if tpl.Rotation != 0 {
				angle := float64(tpl.Rotation) * math.Pi / 180.0
				c = math.Cos(float64(angle))
				s = math.Sin(float64(angle))

				switch tpl.Rotation {
				case -90:
					tx = -tpl.box["lly"]
					ty = tpl.box["urx"]
					break

				case -180:
					tx = tpl.box["urx"]
					ty = tpl.box["ury"]
					break

				case -270:
					tx = tpl.box["ury"]
					ty = -tpl.box["llx"]
				}
			}
		} else {
			tx = -tpl.box["x"] * 2
			ty = tpl.box["y"] * 2
		}

		tx *= this.k
		ty *= this.k

		if c != 1 || s != 0 || tx != 0 || ty != 0 {
			this.out("/Matrix [" +
				strconv.FormatFloat(c, 'f', 5, 64) + " " +
				strconv.FormatFloat(s, 'f', 5, 64) + " " +
				strconv.FormatFloat(-s, 'f', 5, 64) + " " +
				strconv.FormatFloat(c, 'f', 5, 64) + " " +
				strconv.FormatFloat(tx, 'f', 5, 64) + " " +
				strconv.FormatFloat(ty, 'f', 5, 64) +
				"]")
		}

		// Now write resources
		this.out("/Resources ")

		if tpl.Resources != nil {
			this.writeValue(tpl.Resources) // "n" will be changed
		} else {
			return nil, newErr("Template resources are empty")
		}

		nN := this.n // remember new "n"
		this.n = cN  // reset to current "n"

		this.out("/Length " + strconv.Itoa(len(p)) + " >>")

		this.out("stream")
		this.out(p)
		this.out("endstream")

		this.endObj()

		this.n = nN // reset to new "n"

		// Put imported objects, starting with the ones from the XObject's Resources,
		// then from dependencies of those resources).
		err = this.putImportedObjects(reader)
		if err != nil {
			return nil, newErr(err, "Failed to put imported objects")
		}
	}

	return result, nil
}

func (this *pdfWriter) putImportedObjects(reader *pdfReader) error {
	var err error
	var nObj *pdfValue

	// obj_stack will have new items added to it in the inner loop, so do another loop to check for extras
	// TODO make the order of this the same every time
	for {
		atLeastOne := false

		// FIXME:  How to determine number of objects before this loop?
		for i := 0; i < 9999; i++ {
			k := i
			v := this.obj_stack[i]

			if v == nil {
				continue
			}

			atLeastOne = true

			nObj, err = reader.resolveObject(v)
			if err != nil {
				return newErr(err, "Unable to resolve object")
			}

			// New object with "NewId" field
			this.newObj(v.NewId, false)

			if nObj.Type == pdf_TYPE_stream {
				this.writeValue(nObj)
			} else {
				this.writeValue(nObj.Value)
			}

			this.endObj()

			// Remove from stack
			this.obj_stack[k] = nil
		}

		if !atLeastOne {
			break
		}
	}

	return nil
}

// Get the calculated size of a template
// If one size is given, this method calculates the other one
func (this *pdfWriter) getTemplateSize(tplid int, _w float64, _h float64) map[string]float64 {
	result := make(map[string]float64, 2)

	tpl := this.tpls[tplid]

	w := tpl.W
	h := tpl.H

	if _w == 0 && _h == 0 {
		_w = w
		_h = h
	}

	if _w == 0 {
		_w = _h * w / h
	}

	if _h == 0 {
		_h = _w * h / w
	}

	result["w"] = _w
	result["h"] = _h

	return result
}

func (this *pdfWriter) UseTemplate(tplid int, _x float64, _y float64, _w float64, _h float64) (string, float64, float64, float64, float64) {
	tpl := this.tpls[tplid]

	w := tpl.W
	h := tpl.H

	_x += tpl.X
	_y += tpl.Y

	wh := this.getTemplateSize(0, _w, _h)

	_w = wh["w"]
	_h = wh["h"]

	tData := make(map[string]float64, 9)
	tData["x"] = 0.0
	tData["y"] = 0.0
	tData["w"] = _w
	tData["h"] = _h
	tData["scaleX"] = (_w / w)
	tData["scaleY"] = (_h / h)
	tData["tx"] = _x
	tData["ty"] = (0 - _y - _h)
	tData["lty"] = (0 - _y - _h) - (0-h)*(_h/h)

	return "/GOFPDITPL" + strconv.Itoa(tplid+this.tpl_id_offset), tData["scaleX"], tData["scaleY"], tData["tx"] * this.k, tData["ty"] * this.k
}
