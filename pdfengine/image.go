package pdfengine

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/errs"
)

// imageHolder hold image data
type imageHolder interface {
	ID() string
	io.Reader
}

// imageCache is metadata for caching images.
type imageCache struct {
	Path  string //ID or Path
	Index int
	Rect  *canvas.Rect
}

type imgInfo struct {
	w, h int
	//src              string
	formatName       string
	colspace         string
	bitsPerComponent string
	filter           string
	decodeParms      string
	trns             []byte
	smask            []byte
	smarkObjID       int
	pal              []byte
	deviceRGBObjID   int
	data             []byte
}

// DrawImageByHolder : draw image by imageHolder
func (gp *PdfEngine) DrawImageByHolder(img imageHolder, x float64, y float64, rect *canvas.Rect) error {
	gp.UnitsToPointsVar(&x, &y)

	rect = rect.UnitsToPoints(gp.Config.Unit)

	ImageOptions := ImageOptions{
		X:    x,
		Y:    y,
		Rect: rect,
	}

	return gp.imageByHolder(img, ImageOptions)
}

func (gp *PdfEngine) ImageByHolderWithOptions(img imageHolder, opts ImageOptions) error {
	gp.UnitsToPointsVar(&opts.X, &opts.Y)

	opts.Rect = opts.Rect.UnitsToPoints(gp.Config.Unit)

	imageTransparency, err := gp.getCachedTransparency(opts.transparency)
	if err != nil {
		return err
	}

	if imageTransparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, imageTransparency.extGStateIndex)
	}

	if opts.Mask != nil {
		maskTransparency, err := gp.getCachedTransparency(opts.Mask.ImageOptions.transparency)
		if err != nil {
			return err
		}

		if maskTransparency != nil {
			opts.Mask.ImageOptions.extGStateIndexes = append(opts.Mask.ImageOptions.extGStateIndexes, maskTransparency.extGStateIndex)
		}

		gp.UnitsToPointsVar(&opts.Mask.ImageOptions.X, &opts.Mask.ImageOptions.Y)
		opts.Mask.ImageOptions.Rect = opts.Mask.ImageOptions.Rect.UnitsToPoints(gp.Config.Unit)

		extGStateIndex, err := gp.maskHolder(opts.Mask.Holder, *opts.Mask)
		if err != nil {
			return err
		}

		opts.extGStateIndexes = append(opts.extGStateIndexes, extGStateIndex)
	}

	return gp.imageByHolder(img, opts)
}

func (gp *PdfEngine) maskHolder(img imageHolder, opts maskOptions) (int, error) {
	var cacheImage *imageCache
	var cacheContentImage *cacheContentImage

	for _, imgcache := range gp.curr.ImgCaches {
		if img.ID() == imgcache.Path {
			cacheImage = &imgcache
			break
		}
	}

	if cacheImage == nil {
		maskImgobj := &imageObj{IsMask: true}
		maskImgobj.Init(func() *PdfEngine {
			return gp
		})
		maskImgobj.setProtection(gp.protection())

		err := maskImgobj.SetImage(img)
		if err != nil {
			return 0, err
		}

		if opts.Rect == nil {
			if opts.Rect, err = maskImgobj.getRect(); err != nil {
				return 0, err
			}
		}

		if err := maskImgobj.parse(); err != nil {
			return 0, err
		}

		if gp.indexOfProcSet != -1 {
			index := gp.addObj(maskImgobj)
			cacheContentImage = gp.getContent().GetCacheContentImage(index, opts.ImageOptions)
			procset := gp.pdfObjs[gp.indexOfProcSet].(*procSetObj)
			procset.RelateXobjs = append(procset.RelateXobjs, relateXobject{IndexOfObj: index})

			imgcache := imageCache{
				Index: index,
				Path:  img.ID(),
				Rect:  opts.Rect,
			}
			gp.curr.ImgCaches[index] = imgcache
			gp.curr.CountOfImg++
		}
	} else {
		if opts.Rect == nil {
			opts.Rect = gp.curr.ImgCaches[cacheImage.Index].Rect
		}

		cacheContentImage = gp.getContent().GetCacheContentImage(cacheImage.Index, opts.ImageOptions)
	}

	if cacheContentImage != nil {
		extGStateInd, err := gp.createTransparencyXObjectGroup(cacheContentImage, opts)
		if err != nil {
			return 0, err
		}

		return extGStateInd, nil
	}

	return 0, errs.UndefinedCacheContentImage
}

func (gp *PdfEngine) createTransparencyXObjectGroup(image *cacheContentImage, opts maskOptions) (int, error) {
	bbox := opts.BBox
	if bbox == nil {
		bbox = &[4]float64{
			// correct BBox values is [opts.X, gp.curr.pageSize.H - opts.Y - opts.canvas.Rect.H, opts.X + opts.canvas.Rect.W, gp.curr.pageSize.H - opts.Y]
			// but if compress pdf through ghostscript result file can't open correctly in mac viewer, because mac viewer can't parse BBox value correctly
			// all other viewers parse BBox correctly (like Adobe Acrobat Reader, Chrome, even Internet Explorer)
			// that's why we need to set [0, 0, gp.curr.pageSize.W, gp.curr.pageSize.H]
			-gp.curr.pageSize.W * 2,
			-gp.curr.pageSize.H * 2,
			gp.curr.pageSize.W * 2,
			gp.curr.pageSize.H * 2,
			// Also, Chrome pdf viewer incorrectly recognize BBox value, that's why we need to set twice as much value
			// for every mask element will be displayed
		}
	}

	groupOpts := transparencyXObjectGroupOptions{
		BBox:             *bbox,
		ExtGStateIndexes: opts.extGStateIndexes,
		XObjects:         []cacheContentImage{*image},
	}

	transparencyXObjectGroup, err := getCachedTransparencyXObjectGroup(groupOpts, gp)
	if err != nil {
		return 0, err
	}

	sMaskOptions := sMaskOptions{
		Subtype:                       sMaskLuminositySubtype,
		TransparencyXObjectGroupIndex: transparencyXObjectGroup.Index,
	}
	sMask := getCachedMask(sMaskOptions, gp)

	extGStateOpts := extGStateOptions{SMaskIndex: &sMask.Index}
	extGState, err := getCachedExtGState(extGStateOpts, gp)
	if err != nil {
		return 0, err
	}

	return extGState.Index + 1, nil
}

func (gp *PdfEngine) imageByHolder(img imageHolder, opts ImageOptions) error {
	cacheImageIndex := -1

	for _, imgcache := range gp.curr.ImgCaches {
		if img.ID() == imgcache.Path {
			cacheImageIndex = imgcache.Index
			break
		}
	}

	if cacheImageIndex == -1 { //new image

		//create img object
		imgobj := new(imageObj)
		if opts.Mask != nil {
			imgobj.SplittedMask = true
		}

		imgobj.Init(func() *PdfEngine {
			return gp
		})
		imgobj.setProtection(gp.protection())

		err := imgobj.SetImage(img)
		if err != nil {
			return err
		}

		if opts.Rect == nil {
			if opts.Rect, err = imgobj.getRect(); err != nil {
				return err
			}
		}

		err = imgobj.parse()
		if err != nil {
			return err
		}
		index := gp.addObj(imgobj)
		if gp.indexOfProcSet != -1 {
			//ยัดรูป
			procset := gp.pdfObjs[gp.indexOfProcSet].(*procSetObj)
			gp.getContent().AppendStreamImage(index, opts)
			procset.RelateXobjs = append(procset.RelateXobjs, relateXobject{IndexOfObj: index})
			//เก็บข้อมูลรูปเอาไว้
			var imgcache imageCache
			imgcache.Index = index
			imgcache.Path = img.ID()
			imgcache.Rect = opts.Rect
			gp.curr.ImgCaches[index] = imgcache
			gp.curr.CountOfImg++
		}

		if imgobj.haveSMask() {
			smaskObj, err := imgobj.createSMask()
			if err != nil {
				return err
			}
			imgobj.imginfo.smarkObjID = gp.addObj(smaskObj)
		}

		if imgobj.isColspaceIndexed() {
			dRGB, err := imgobj.createDeviceRGB()
			if err != nil {
				return err
			}
			dRGB.getRoot = func() *PdfEngine {
				return gp
			}
			imgobj.imginfo.deviceRGBObjID = gp.addObj(dRGB)
		}

	} else { //same img
		if opts.Rect == nil {
			opts.Rect = gp.curr.ImgCaches[cacheImageIndex].Rect
		}

		gp.getContent().AppendStreamImage(cacheImageIndex, opts)
	}
	return nil
}

// DrawImageInPdf : draw image in pdf by image content
func (gp *PdfEngine) DrawImageInPdf(imageContent []byte, x float64, y float64, rect *canvas.Rect) error {
	gp.UnitsToPointsVar(&x, &y)
	rect = rect.UnitsToPoints(gp.Config.Unit)
	imgh, err := ImageHolderByBytes(imageContent)
	if err != nil {
		return err
	}

	ImageOptions := ImageOptions{
		X:    x,
		Y:    y,
		Rect: rect,
	}

	return gp.imageByHolder(imgh, ImageOptions)
}

func (gp *PdfEngine) ImageFrom(img image.Image, x float64, y float64, rect *canvas.Rect) error {
	return gp.ImageFromWithOption(img, imageFromOption{
		Format: "png",
		X:      x,
		Y:      y,
		Rect:   rect,
	})
}

func (gp *PdfEngine) ImageFromWithOption(img image.Image, opts imageFromOption) error {
	if img == nil {
		return errs.New("Invalid image")
	}

	gp.UnitsToPointsVar(&opts.X, &opts.Y)
	opts.Rect = opts.Rect.UnitsToPoints(gp.Config.Unit)
	r, w := io.Pipe()
	go func() {
		bw := bufio.NewWriter(w)
		var err error
		switch opts.Format {
		case "png":
			err = png.Encode(bw, img)
		case "jpeg":
			err = jpeg.Encode(bw, img, nil)
		}

		bw.Flush()
		if err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}()

	imgh, err := imageHolderByReader(bufio.NewReader(r))
	if err != nil {
		return err
	}

	ImageOptions := ImageOptions{
		X:    opts.X,
		Y:    opts.Y,
		Rect: opts.Rect,
	}

	return gp.imageByHolder(imgh, ImageOptions)
}

// ImageHolderByBytes create imageHolder by []byte
func ImageHolderByBytes(b []byte) (imageHolder, error) {
	return newImageBuff(b)
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
