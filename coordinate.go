package communautofinder

import "math"

type Coordinate struct {
	latitude  float64
	longitude float64
}

func new(latitude float64, longitude float64) Coordinate {
	return Coordinate{latitude: latitude, longitude: longitude}
}

func (baseCoordinate Coordinate) ExpandCoordinate(kilometers float64) (Coordinate, Coordinate) {

	minCoordinate := baseCoordinate
	maxCoordinate := baseCoordinate

	minCoordinate.addKilometersToCoordinate(-kilometers)
	maxCoordinate.addKilometersToCoordinate(kilometers)

	return minCoordinate, maxCoordinate
}

func (coordinateToModify *Coordinate) addKilometersToCoordinate(kilometers float64) {

	earthRadiusKm := 6371.0

	latRad := coordinateToModify.latitude * math.Pi / 180.0
	lonRad := coordinateToModify.longitude * math.Pi / 180.0

	newLatRad := latRad + (kilometers / earthRadiusKm)
	newLonRad := lonRad + (kilometers / (earthRadiusKm * math.Cos(latRad)))

	coordinateToModify.latitude = newLatRad * 180.0 / math.Pi
	coordinateToModify.longitude = newLonRad * 180.0 / math.Pi
}
