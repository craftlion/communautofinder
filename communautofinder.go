package communautofinder

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const cityId = 59 // see available cities -> https://restapifrontoffice.reservauto.net/ReservautoFrontOffice/index.html?urls.primaryName=Branch%20version%202%20(6.93.1)#/

const fetchDelayInMin = 1 // delay between two API call

func communautoAPICall(url string, response interface{}) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		errDecode := json.NewDecoder(resp.Body).Decode(response)

		if errDecode != nil {
			log.Fatal(errDecode)
		}
	} else {
		log.Fatalf("Error %d in API call", resp.StatusCode)
	}
}

// Keep looping until a station car is not found for the specified dates and position
// return the number of car found
func SearchStationCar(currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time) int {

	minCoordinate, maxCoordinate := currentCoordinate.ExpandCoordinate(marginInKm)

	startDateFormat := startDate.Format("2006-01-02T15:04:05")
	endDataFormat := endDate.Format("2006-01-02T15:04:05")

	urlCalled := fmt.Sprintf("https://restapifrontoffice.reservauto.net/api/v2/StationAvailability?CityId=%d&MaxLatitude=%f&MinLatitude=%f&MaxLongitude=%f&MinLongitude=%f&StartDate=%s&EndDate=%s", cityId, maxCoordinate.latitude, minCoordinate.latitude, maxCoordinate.longitude, minCoordinate.longitude, url.QueryEscape(startDateFormat), url.QueryEscape(endDataFormat))

	for {

		var stationsAvailable stationsResponse

		communautoAPICall(urlCalled, &stationsAvailable)

		nbCarFound := len(stationsAvailable.Stations)

		if nbCarFound > 0 {
			return nbCarFound
		}

		time.Sleep(fetchDelayInMin * time.Minute)
	}
}

// Keep looping until a flex car is not found for the specified position
// return the number of car found
func SearchFlexCar(currentCoordinate Coordinate, marginInKm float64) int {

	minCoordinate, maxCoordinate := currentCoordinate.ExpandCoordinate(marginInKm)

	urlCalled := fmt.Sprintf("https://restapifrontoffice.reservauto.net/api/v2/Vehicle/FreeFloatingAvailability?CityId=%d&MaxLatitude=%f&MinLatitude=%f&MaxLongitude=%f&MinLongitude=%f", cityId, maxCoordinate.latitude, minCoordinate.latitude, maxCoordinate.longitude, minCoordinate.longitude)

	for {
		var stationsAvailable flexCarResponse

		communautoAPICall(urlCalled, &stationsAvailable)

		nbCarFound := stationsAvailable.TotalNbVehicles

		if nbCarFound > 0 {
			return nbCarFound
		}

		time.Sleep(fetchDelayInMin * time.Minute)
	}
}

// Keep looping until a station car is not found for the specified dates and position
// fill the channel with number car found
func RoutineSearchStationCar(currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time, responseChannel chan int) {

	defer func() {
		if r := recover(); r != nil {
			responseChannel <- 0
			log.Panicf("Pannic append :", r)
		}
	}()

	responseChannel <- SearchStationCar(currentCoordinate, marginInKm, startDate, endDate)
}

// Keep looping until a flex car is not found for the specified position
// fill the channel with number car found
func RoutineSearchFlexCar(currentCoordinate Coordinate, marginInKm float64, responseChannel chan int) {

	defer func() {
		if r := recover(); r != nil {
			responseChannel <- 0
			log.Panicf("Pannic append :", r)
		}
	}()

	responseChannel <- SearchFlexCar(currentCoordinate, marginInKm)
}
