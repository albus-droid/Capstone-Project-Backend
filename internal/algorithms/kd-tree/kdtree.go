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
	point point2D
	left  *node
	right *node
	axis  int // 0 = split by X (longitude), 1 = split by Y (latitude)
}

// kd-tree with just root node
type KDtree struct {
	root *node
}

// New builds a balanced KDTree from the given slice of points.
func New(points []point2D) *KDTree {
	return &KDTree{root: build(points, 0)}
}

// build recursively constructs a balanced KD-tree by median split.
// axis toggles between 0 (X) and 1 (Y).
func build(points []point2D, axis int) *node {
	if len(points) == 0 {
		return nil
	}
	// Sort by the current axis
	sort.Slice(points, func(i, j int) bool {
		if axis == 0 {
			return points[i].X < points[j].X
		}
		return points[i].Y < points[j].Y
	})
	mid := len(points) / 2
	nextAxis := (axis + 1) % 2

	return &node{
		point: points[mid],
		axis:  axis,
		left:  build(points[:mid], nextAxis),
		right: build(points[mid+1:], nextAxis),
	}
}


// RangeSearch returns all points within radiusKm of center, using Haversine + pruning.
func (t *KDTree) RangeSearch(center point2D, radiusKm float64) []point2D {
	var out []point2D
	search(t.root, center, radiusKm, &out)
	return out
}

// search is the recursive helper that visits nodes and prunes subtrees when possible.
func search(n *node, center point2D, radiusKm float64, out *[]point2D) {
	if n == nil {
		return
	}

	// 1) If this node’s point is within radius, include it.
	if center.Distance(n.point) <= radiusKm {
		*out = append(*out, n.point)
	}

	// 2) Compute axis difference (delta)
	var coordCenter, coordSplit float64
	if n.axis == 0 {
		coordCenter = center.X
		coordSplit = n.point.X
	} else {
		coordCenter = center.Y
		coordSplit = n.point.Y
	}
	delta := coordCenter - coordSplit

	// 3) Always search the side where center falls, then prune the other side
	if delta <= 0 {
		// search “left” subtree first
		search(n.left, center, radiusKm, out)
		// only if |delta| ≤ approxDegrees(radiusKm) do we need to check right subtree
		if math.Abs(delta) <= approxDegrees(radiusKm) {
			search(n.right, center, radiusKm, out)
		}
	} else {
		// search “right” subtree first
		search(n.right, center, radiusKm, out)
		if math.Abs(delta) <= approxDegrees(radiusKm) {
			search(n.left, center, radiusKm, out)
		}
	}
}

// approxDegrees converts a distance in kilometers to a rough degree difference.
// (1° latitude ≈ 111 km; 1° longitude near equator ≈ 111 km, farther away it shrinks,
// but as a pruning heuristic this is safe.)
func approxDegrees(km float64) float64 {
	return km / 111.0
}

// (Optional) String helper to visualize the tree (preorder).
func (t *KDTree) String() string {
	return fmt.Sprintf("[%s]", printNode(t.root))
}
func printNode(n *node) string {
	if n == nil {
		return ""
	}
	leftStr := printNode(n.left)
	rightStr := printNode(n.right)
	return fmt.Sprintf("(%0.4f,%0.4f axis=%d) %s %s",
		n.point.X, n.point.Y, n.axis, leftStr, rightStr)
}