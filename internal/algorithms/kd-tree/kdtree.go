import (
	"fmt"
	"math"
	"sort"
)

// 2-d points struct
type point2D struct{
	X, Y float64
}

// finding distance using haversines formula
func (p Point2D) Distance(q Point2D) float64 {
	const R = 6371.0 // Earth radius in kilometers

	lat1 := p.Y * math.Pi / 180
	lon1 := p.X * math.Pi / 180
	lat2 := q.Y * math.Pi / 180
	lon2 := q.X * math.Pi / 180

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// node struct
type node struct {
	point2D
	Left *node
	Right *node
}

// kd-tree with just root node
type KDtree struct {
	root *node
}

func new

