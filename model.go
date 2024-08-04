package communautofinder

// Different type of vehicule possible to search

type SearchType int

const (
	searchingFlex    SearchType = 0
	searchingStation SearchType = 1
)

type CityId int

const (
	Montreal CityId = 59
)

type VehiculeType int

const (
	FamilyCar      VehiculeType = 1
	UtilityVehicle VehiculeType = 2
	MidSize        VehiculeType = 3
	Minivan        VehiculeType = 5
)