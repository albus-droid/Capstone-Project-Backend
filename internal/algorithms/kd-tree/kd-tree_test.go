// kd_test.go
package kdtree

import (
	"testing"
)

// Helper: compare two slices of Point2D (order doesn’t matter)
func sameSet(a, b []point2D) bool {
	if len(a) != len(b) {
		return false
	}
	used := make([]bool, len(b))
	for _, pa := range a {
		found := false
		for j, pb := range b {
			if used[j] {
				continue
			}
			// “Exact” match on longitude+latitude
			if pa.X == pb.X && pa.Y == pb.Y {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestRangeSearch(t *testing.T) {
	points := []point2D{
		{X: -79.3832, Y: 43.6532}, // P0: Toronto
		{X: -80.0000, Y: 43.5000}, // P1: ~33 km WSW
		{X: -79.9000, Y: 43.7000}, // P2: ~59 km WNW
		{X: -79.0000, Y: 43.8000}, // P3: ~63 km E
		{X: -79.8000, Y: 43.6000}, // P4: ~46 km W
	}

	tree := New(points)

	tests := []struct {
		center   point2D
		radiusKm float64
		expected []point2D
	}{
		{
			// Searching around Toronto with radius 50 km should catch P0, P1, P4
			center:   point2D{X: -79.3832, Y: 43.6532},
			radiusKm: 50.0,
			expected: []point2D{
				{X: -79.3832, Y: 43.6532}, // P0 itself
				{X: -80.0000, Y: 43.5000}, // P1
				{X: -79.8000, Y: 43.6000}, // P4
			},
		},
		{
			// Searching around a point far away with radius 10 km should catch nobody
			center:   point2D{X: -100.0, Y: 50.0},
			radiusKm: 10.0,
			expected: []point2D{},
		},
		{
			// Radius zero at Toronto should return exactly that point
			center:   point2D{X: -79.3832, Y: 43.6532},
			radiusKm: 0.0,
			expected: []point2D{
				{X: -79.3832, Y: 43.6532},
			},
		},
	}

	for _, tc := range tests {
		got := tree.RangeSearch(tc.center, tc.radiusKm)
		if !sameSet(got, tc.expected) {
			t.Errorf("RangeSearch(center=%.4f,%.4f, radius=%.1f) = %v, want %v",
				tc.center.X, tc.center.Y, tc.radiusKm, got, tc.expected)
		}
	}
}
