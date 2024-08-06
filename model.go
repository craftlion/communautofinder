package communautofinder

// Different type of vehicule possible to search

type SearchType int

const (
	SearchingFlex    SearchType = 0
	SearchingStation SearchType = 1
)

type CityId int

const (
	Montreal CityId = 59
)

type VehiculeType int

const (
	AllTypes       VehiculeType = 0
	FamilyCar      VehiculeType = 1
	UtilityVehicle VehiculeType = 2
	MidSize        VehiculeType = 3
	Minivan        VehiculeType = 5
)
