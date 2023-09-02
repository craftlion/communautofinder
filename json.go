package communautofinder

// --------------------------------------------
// Types used to decode Communauto API response
// --------------------------------------------

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type stationsResponse struct {
	Stations []station `json:"stations"`
}

type station struct {
	StationId              int         `json:"stationId"`
	StationNb              string      `json:"stationNb"`
	StationName            string      `json:"stationName"`
	StationLocation        location    `json:"stationLocation"`
	CityId                 int         `json:"cityId"`
	RecommendedVehicleId   int         `json:"recommendedVehicleId"`
	HasAllRequestedOptions bool        `json:"hasAllRequestedOptions"`
	SatisfiesFilters       bool        `json:"satisfiesFilters"`
	VehiclePromotions      interface{} `json:"vehiclePromotions"` // You can replace interface{} with an appropriate type if known
	HasZone                bool        `json:"hasZone"`
}

type flexCarResponse struct {
	TotalNbVehicles int       `json:"totalNbVehicles"`
	Vehicles        []vehicle `json:"vehicles"`
}

type vehicle struct {
	VehicleId                 int      `json:"vehicleId"`
	VehicleNb                 int      `json:"vehicleNb"`
	CityId                    int      `json:"cityId"`
	VehiclePropulsionTypeId   int      `json:"vehiclePropulsionTypeId"`
	VehicleTypeId             int      `json:"vehicleTypeId"`
	VehicleBodyTypeId         int      `json:"vehicleBodyTypeId"`
	VehicleTransmissionTypeId int      `json:"vehicleTransmissionTypeId"`
	VehicleTireTypeId         int      `json:"vehicleTireTypeId"`
	VehicleAccessories        []int    `json:"vehicleAccessories"`
	VehicleLocation           location `json:"vehicleLocation"`
	SatisfiesFilters          bool     `json:"satisfiesFilters"`
}
