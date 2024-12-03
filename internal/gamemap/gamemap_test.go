package gamemap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLocation(t *testing.T) {
	testCases := []struct {
		Name         string
		LocationName string
		ExpectExists bool
		ExpectPoint  Point
	}{
		{
			Name:         "Existing Location",
			LocationName: ManderMine,
			ExpectExists: true,
			ExpectPoint:  Point{14830, 3830},
		},
		{
			Name:         "Nonexistent Location",
			LocationName: "Nonexistent Location",
			ExpectExists: false,
			ExpectPoint:  Point{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			point, exists := GetLocation(tc.LocationName)
			require.Equal(t, tc.ExpectExists, exists, "Expected existence match for location: %s", tc.LocationName)
			require.Equal(t, tc.ExpectPoint, point, "Expected Point to match for location: %s", tc.LocationName)
		})
	}
}

func TestGetDistance(t *testing.T) {
	testCases := []struct {
		Name           string
		PointA, PointB Point
		ExpectDistance float64
	}{
		{
			Name:           "3-4-5 Triangle",
			PointA:         Point{0, 0},
			PointB:         Point{3, 4},
			ExpectDistance: 5.0,
		},
		{
			Name:           "Zero Distance",
			PointA:         Point{1, 1},
			PointB:         Point{1, 1},
			ExpectDistance: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			distance := GetDistance(tc.PointA, tc.PointB)
			require.Equal(t, tc.ExpectDistance, distance, "Expected distance to match")
		})
	}
}

func TestGetAngle(t *testing.T) {
	testCases := []struct {
		Name           string
		PointA, PointB Point
		ExpectAngle    float64
	}{
		{
			Name:        "Angle 0 Degrees",
			PointA:      Point{0, 0},
			PointB:      Point{1, 0},
			ExpectAngle: 0.0,
		},
		{
			Name:        "Angle 90 Degrees",
			PointA:      Point{0, 0},
			PointB:      Point{0, 1},
			ExpectAngle: 90.0,
		},
		{
			Name:        "Angle 180 Degrees",
			PointA:      Point{0, 0},
			PointB:      Point{-1, 0},
			ExpectAngle: 180.0,
		},
		{
			Name:        "Angle 270 Degrees",
			PointA:      Point{0, 0},
			PointB:      Point{0, -1},
			ExpectAngle: 270.0,
		},
		{
			Name:        "Angle 45 Degrees",
			PointA:      Point{0, 0},
			PointB:      Point{1, 1},
			ExpectAngle: 45.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			angle := GetAngle(tc.PointA, tc.PointB)
			require.InDelta(t, tc.ExpectAngle, angle, 0.0001, "Expected angle to be approximately %f degrees", tc.ExpectAngle)
		})
	}
}
