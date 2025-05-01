package canvas

import (
	"fmt"
	"math"
)

// Point is an X,Y pair
type Point struct {
	X, Y int
}

// DistanceTo calculates the distance to another point.
func (p Point) DistanceTo(other Point) float64 {
	dx := math.Pow(float64(p.X-other.X), 2)
	dy := math.Pow(float64(p.Y-other.Y), 2)
	return math.Pow(dx+dy, 0.5)
}

// Equals returns if a point equals another point.
func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

// String returns a string representation of the point.
func (p Point) String() string {
	return fmt.Sprintf("P{%d,%d}", p.X, p.Y)
}
