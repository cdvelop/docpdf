package docpdf

import (
	"bytes"
	"image"

	"github.com/cdvelop/docpdf/alignment"
	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/env"
	"github.com/cdvelop/docpdf/errs"
)

// docImage represents an image to be added to the document
type docImage struct {
	doc           *Document
	pathOrContent any //eg: "path/to/image.png" or []byte{...}
	width         float64
	height        float64
	keepRatio     bool
	alignment     alignment.Alignment
	x, y          float64
	hasPos        bool
	inline        bool // New property to track inline status
	valign        int  // Vertical alignment for inline images
}

// AddImage creates a new image element in the document
// supporting both absolute and relative paths or image data in []byte format
// eg: doc.AddImage("path/to/image.png") or doc.AddImage([]byte{...})
func (doc *Document) AddImage(imagePathOrContent any) *docImage {
	return &docImage{
		doc:           doc,
		pathOrContent: imagePathOrContent,
		keepRatio:     true,
		alignment:     alignment.Left,
	}
}

// Width sets the image width and maintains aspect ratio if height is not set
// eg: img.Width(50) will set the width to 50 and calculate height based on aspect ratio
func (img *docImage) Width(w float64) *docImage {
	img.width = w
	return img
}

// Height sets the image height and maintains aspect ratio if width is not set
// eg: img.Height(50) will set the height to 100 and calculate width based on aspect ratio
func (img *docImage) Height(h float64) *docImage {
	img.height = h
	return img
}

// Size sets both width and height explicitly (no aspect ratio preservation)
// eg: img.Size(50, 100) will set the width to 50 and height to 100
func (img *docImage) Size(w, h float64) *docImage {
	img.width = w
	img.height = h
	img.keepRatio = false
	return img
}

// FixedPosition places the image at specific coordinates
func (img *docImage) FixedPosition(x, y float64) *docImage {
	img.x = x
	img.y = y
	img.hasPos = true
	return img
}

// AlignLeft aligns the image to the left margin
func (img *docImage) AlignLeft() *docImage {
	img.alignment = alignment.Left
	return img
}

// AlignCenter centers the image horizontally
func (img *docImage) AlignCenter() *docImage {
	img.alignment = alignment.Center
	return img
}

// AlignRight aligns the image to the right margin
func (img *docImage) AlignRight() *docImage {
	img.alignment = alignment.Right
	return img
}

// Inline makes the image display inline with text rather than breaking to a new line
// The text will continue from the right side of the image
func (img *docImage) Inline() *docImage {
	img.inline = true
	return img
}

// VerticalAlignTop aligns the image with the top of the text line when inline
func (img *docImage) VerticalAlignTop() *docImage {
	img.valign = 0
	return img
}

// VerticalAlignMiddle aligns the image with the middle of the text line when inline
func (img *docImage) VerticalAlignMiddle() *docImage {
	img.valign = 1
	return img
}

// VerticalAlignBottom aligns the image with the bottom of the text line when inline
func (img *docImage) VerticalAlignBottom() *docImage {
	img.valign = 2
	return img
}

// Draw renders the image on the document to include page break handling
func (img *docImage) Draw() error {
	imageContent, err := env.FileExists(img.pathOrContent)
	if err != nil {
		return errs.New("Image Draw error", err)
	}

	// Get image dimensions to calculate aspect ratio if needed
	imgWidth, imgHeight, err := img.getImageDimensions(imageContent)
	if err != nil {
		return err
	}

	// Calculate final dimensions
	finalWidth, finalHeight := img.calculateDimensions(imgWidth, imgHeight)

	// Check if the image has a fixed alignment.Alignment
	if !img.hasPos {
		// Skip page break check if we're in header/footer drawing mode
		if !img.doc.inHeaderFooterDraw {
			// Check if the image fits on current page
			newY := img.doc.EnsureElementFits(finalHeight)

			// Only update Y alignment.Alignment if this is not an inline element
			if !img.inline {
				img.doc.SetY(newY)
			}
		}
	}

	// Determine alignment.Alignment (after possible page break)
	x, y := img.calculatePosition(finalWidth)

	// Adjust vertical alignment.Alignment for inline images based on alignment
	if img.inline {
		lineHeight := img.doc.GetLineHeight()

		switch img.valign {
		case 0: // alignment.Top alignment
			// No adjustment needed
		case 1: // alignment.Middle alignment
			y = y + (lineHeight-finalHeight)/2
		case 2: // alignment.Bottom alignment
			y = y + lineHeight - finalHeight
		default:
			// Default to middle alignment
			y = y + (lineHeight-finalHeight)/2
		}
	}

	// Create rectangle for the image
	rect := &canvas.Rect{
		W: finalWidth,
		H: finalHeight,
	}
	// Draw the image using the underlying PdfEngine instance
	err = img.doc.DrawImageInPdf(imageContent, x, y, rect)
	if err != nil {
		return err
	}

	// Handle alignment.Alignment updates based on inline setting
	if img.inline {
		// For inline images, advance X alignment.Alignment but keep Y unchanged
		img.doc.SetX(x + finalWidth)

		// Store that we have an inline element active
		img.doc.inlineMode = true
	} else {
		// For block images, advance Y alignment.Alignment to avoid text overlapping with the image
		if !img.hasPos {
			img.doc.newLineBreakBasedOnDefaultFont(y + finalHeight)
		}

		// Reset X alignment.Alignment to left margin since this is a block element
		img.doc.SetX(img.doc.Margins().Left)

		// Reset inline mode
		img.doc.inlineMode = false
	}

	return nil
}

// getImageDimensions returns the natural width and height of the image
func (img *docImage) getImageDimensions(imageContent []byte) (float64, float64, error) {
	reader := bytes.NewReader(imageContent)

	imgConfig, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}

	return float64(imgConfig.Width), float64(imgConfig.Height), nil
}

// calculateDimensions determines the final width and height of the image
func (img *docImage) calculateDimensions(imgWidth, imgHeight float64) (float64, float64) {
	// Default to original dimensions
	finalWidth := imgWidth
	finalHeight := imgHeight

	// Scale down if original image is too large
	contentAreaWidth := img.doc.contentAreaWidth - img.doc.Margins().Left - img.doc.Margins().Right
	if finalWidth > contentAreaWidth {
		ratio := contentAreaWidth / finalWidth
		finalWidth = contentAreaWidth
		finalHeight = finalHeight * ratio
	}

	// Apply user-specified dimensions
	if img.width > 0 && img.height > 0 {
		// Both dimensions specified
		finalWidth = img.width
		finalHeight = img.height
	} else if img.width > 0 && img.keepRatio {
		// Only width specified, calculate height to maintain aspect ratio
		ratio := img.width / finalWidth
		finalWidth = img.width
		finalHeight = finalHeight * ratio
	} else if img.height > 0 && img.keepRatio {
		// Only height specified, calculate width to maintain aspect ratio
		ratio := img.height / finalHeight
		finalHeight = img.height
		finalWidth = finalWidth * ratio
	}

	return finalWidth, finalHeight
}

// calculatePosition determines where to place the image
func (img *docImage) calculatePosition(width float64) (float64, float64) {
	if img.hasPos {
		return img.x, img.y
	}

	x := img.doc.Margins().Left
	y := img.doc.GetY()

	// Apply alignment
	switch img.alignment {
	case alignment.Center:
		x = img.doc.Margins().Left + (img.doc.contentAreaWidth-width)/2
	case alignment.Right:
		x = img.doc.Margins().Left + img.doc.contentAreaWidth - width
	}

	return x, y
}
