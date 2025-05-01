package canvas

import (
	"fmt"

	"github.com/cdvelop/docpdf/mathutils"
)

// BoxCorners is a box with independent corners.
type BoxCorners struct {
	TopLeft, TopRight, BottomRight, BottomLeft Point
}

// Box return the BoxCorners as a regular box.
func (bc BoxCorners) Box() Box {
	return Box{
		Top:    mathutils.MinInt(bc.TopLeft.Y, bc.TopRight.Y),
		Left:   mathutils.MinInt(bc.TopLeft.X, bc.BottomLeft.X),
		Right:  mathutils.MaxInt(bc.TopRight.X, bc.BottomRight.X),
		Bottom: mathutils.MaxInt(bc.BottomLeft.Y, bc.BottomRight.Y),
	}
}

// Width returns the width
func (bc BoxCorners) Width() int {
	minLeft := mathutils.MinInt(bc.TopLeft.X, bc.BottomLeft.X)
	maxRight := mathutils.MaxInt(bc.TopRight.X, bc.BottomRight.X)
	return maxRight - minLeft
}

// Height returns the height
func (bc BoxCorners) Height() int {
	minTop := mathutils.MinInt(bc.TopLeft.Y, bc.TopRight.Y)
	maxBottom := mathutils.MaxInt(bc.BottomLeft.Y, bc.BottomRight.Y)
	return maxBottom - minTop
}

// Center returns the center of the box
func (bc BoxCorners) Center() (x, y int) {

	left := mathutils.MeanInt(bc.TopLeft.X, bc.BottomLeft.X)
	right := mathutils.MeanInt(bc.TopRight.X, bc.BottomRight.X)
	x = ((right - left) >> 1) + left

	top := mathutils.MeanInt(bc.TopLeft.Y, bc.TopRight.Y)
	bottom := mathutils.MeanInt(bc.BottomLeft.Y, bc.BottomRight.Y)
	y = ((bottom - top) >> 1) + top

	return
}

// Rotate rotates the box.
func (bc BoxCorners) Rotate(thetaDegrees float64) BoxCorners {
	cx, cy := bc.Center()

	thetaRadians := mathutils.DegreesToRadians(thetaDegrees)

	tlx, tly := mathutils.RotateCoordinate(cx, cy, bc.TopLeft.X, bc.TopLeft.Y, thetaRadians)
	trx, try := mathutils.RotateCoordinate(cx, cy, bc.TopRight.X, bc.TopRight.Y, thetaRadians)
	brx, bry := mathutils.RotateCoordinate(cx, cy, bc.BottomRight.X, bc.BottomRight.Y, thetaRadians)
	blx, bly := mathutils.RotateCoordinate(cx, cy, bc.BottomLeft.X, bc.BottomLeft.Y, thetaRadians)

	return BoxCorners{
		TopLeft:     Point{tlx, tly},
		TopRight:    Point{trx, try},
		BottomRight: Point{brx, bry},
		BottomLeft:  Point{blx, bly},
	}
}

// Equals returns if the box equals another box.
func (bc BoxCorners) Equals(other BoxCorners) bool {
	return bc.TopLeft.Equals(other.TopLeft) &&
		bc.TopRight.Equals(other.TopRight) &&
		bc.BottomRight.Equals(other.BottomRight) &&
		bc.BottomLeft.Equals(other.BottomLeft)
}

func (bc BoxCorners) String() string {
	return fmt.Sprintf("BoxC{%s,%s,%s,%s}", bc.TopLeft.String(), bc.TopRight.String(), bc.BottomRight.String(), bc.BottomLeft.String())
}
