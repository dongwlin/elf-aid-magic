package gamemap

import "math"

// Point represents a coordinate on the map.
type Point struct {
	X int
	Y int
}

const (
	ManderMine                   = "ManderMine"
	ClarityDataCenterAdminBureau = "ClarityDataCenterAdminBureau"
	ShoggolithCity               = "ShoggolithCity"
	WildernessStation            = "WildernessStation"
	YunxiuBridge                 = "YunxiuBridge"
	BRCLOutpost                  = "BRCLOutpost"
	Freeport                     = "Freeport"
	AnitaWeaponResearchInstitute = "AnitaWeaponResearchInstitute"
	Onederland                   = "Onederland"
	AnitaRocketBase              = "AnitaRocketBase"
	GongluCity                   = "GongluCity"
	CapeCity                     = "CapeCity"
	ConfluenceTower              = "ConfluenceTower"
	AnitaEnergyResearchInstitute = "AnitaEnergyResearchInstitute"
)

var locationMap = map[string]Point{
	ManderMine:                   {14830, 3830},
	ClarityDataCenterAdminBureau: {8830, 665},
	ShoggolithCity:               {11665, 1835},
	WildernessStation:            {16330, 1835},
	YunxiuBridge:                 {20330, 0},
	BRCLOutpost:                  {13830, 1835},
	Freeport:                     {5000, 2995},
	AnitaWeaponResearchInstitute: {6665, 3830},
	Onederland:                   {15830, 6660},
	AnitaRocketBase:              {0, 7165},
	GongluCity:                   {9995, 10160},
	CapeCity:                     {2995, 12830},
	ConfluenceTower:              {660, 12830},
	AnitaEnergyResearchInstitute: {4495, 7495},
}

// GetLocation returns the Point for a given location name.
func GetLocation(name string) (Point, bool) {
	point, exists := locationMap[name]
	return point, exists
}

// GetDistance calculates the Euclidean distance between two points.
func GetDistance(a, b Point) float64 {
	dx := float64(b.X - a.X)
	dy := float64(b.Y - a.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// GetAngle calculates the angle in degrees from point a to point b.
func GetAngle(a, b Point) float64 {
	dx := float64(b.X - a.X)
	dy := float64(b.Y - a.Y)
	angle := math.Atan2(dy, dx) * (180 / math.Pi) // Convert radians to degrees
	if angle < 0 {
		angle += 360 // Ensure the angle is non-negative
	}
	return angle
}
