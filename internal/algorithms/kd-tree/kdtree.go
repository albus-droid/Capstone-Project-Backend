import (
	"fmt"
	"math"
	"sort"
)

// 2-d points struct
type point2D struct{
	X float64
	Y float64
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

