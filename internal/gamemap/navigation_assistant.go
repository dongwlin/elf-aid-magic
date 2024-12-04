package gamemap

// NavigationAssistant assists with navigation tasks.
type NavigationAssistant struct {
	currentLocation Point
}

// NewNavigationAssistant creates a new NavigationAssistant at the given location.
func NewNavigationAssistant() *NavigationAssistant {
	return &NavigationAssistant{}
}

func (na *NavigationAssistant) SetCurrentLocation(location Point) {
	na.currentLocation = location
}

// MoveTo updates the current location of the assistant.
func (na *NavigationAssistant) MoveTo(destination Point) {
	na.currentLocation = destination
}

// StraightLineDistanceTo calculates the distance from the current location to a target location.
func (na *NavigationAssistant) StraightLineDistanceTo(target Point) float64 {
	return GetStraightLineDistance(na.currentLocation, target)
}

// AngleTo calculates the angle from the current location to a target location.
func (na *NavigationAssistant) AngleTo(target Point) float64 {
	return GetAngle(na.currentLocation, target)
}
