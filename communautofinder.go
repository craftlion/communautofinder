package communautofinder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const fetchDelayInMin = 1 // delay between two API call

const dateFormat = "2006-01-02T15:04:05" // time format accepted by communauto API

// Different type of vehicule possible to search
const (
	searchingFlex = iota
	searchingStation
)

// As soon as at least one car is found return the number of cars found
func SearchStationCar(cityId int, currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time) int {
	responseChannel := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	nbCarFound := searchCar(searchingStation, cityId, currentCoordinate, marginInKm, startDate, endDate, responseChannel, ctx)
	cancel()

	return nbCarFound
}

// As soon as at least one car is found return the number of cars found
func SearchFlexCar(cityId int, currentCoordinate Coordinate, marginInKm float64) int {
	responseChannel := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	nbCarFound := searchCar(searchingFlex, cityId, currentCoordinate, marginInKm, time.Time{}, time.Time{}, responseChannel, ctx)
	cancel()

	return nbCarFound
}

// This function is designed to be called as a goroutine. As soon as at least one car is found return the number of cars found. Or can be cancelled by the context
func SearchStationCarForGoRoutine(cityId int, currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time, responseChannel chan int, ctx context.Context) int {

	defer func() {
		if r := recover(); r != nil {
			responseChannel <- -1
			log.Printf("Pannic append : %s", r)
		}
	}()

	return searchCar(searchingStation, cityId, currentCoordinate, marginInKm, startDate, endDate, responseChannel, ctx)
}

// This function is designed to be called as a goroutine. As soon as at least one car is found return the number of cars found. Or can be cancelled by the context
func SearchFlexCarForGoRoutine(cityId int, currentCoordinate Coordinate, marginInKm float64, responseChannel chan int, ctx context.Context) int {

	defer func() {
		if r := recover(); r != nil {
			responseChannel <- -1
			log.Printf("Pannic append : %s", r)
		}
	}()

	return searchCar(searchingFlex, cityId, currentCoordinate, marginInKm, time.Time{}, time.Time{}, responseChannel, ctx)
}

// Loop until a result is found. Return the number of cars found or can be cancelled by the context
func searchCar(searchingType int, cityId int, currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time, responseChannel chan int, ctx context.Context) int {
	minCoordinate, maxCoordinate := currentCoordinate.ExpandCoordinate(marginInKm)

	var urlCalled string

	if searchingType == searchingFlex {
		urlCalled = fmt.Sprintf("https://restapifrontoffice.reservauto.net/api/v2/Vehicle/FreeFloatingAvailability?CityId=%d&MaxLatitude=%f&MinLatitude=%f&MaxLongitude=%f&MinLongitude=%f", cityId, maxCoordinate.latitude, minCoordinate.latitude, maxCoordinate.longitude, minCoordinate.longitude)
	} else if searchingType == searchingStation {
		startDateFormat := startDate.Format(dateFormat)
		endDataFormat := endDate.Format(dateFormat)

		urlCalled = fmt.Sprintf("https://restapifrontoffice.reservauto.net/api/v2/StationAvailability?CityId=%d&MaxLatitude=%f&MinLatitude=%f&MaxLongitude=%f&MinLongitude=%f&StartDate=%s&EndDate=%s", cityId, maxCoordinate.latitude, minCoordinate.latitude, maxCoordinate.longitude, minCoordinate.longitude, url.QueryEscape(startDateFormat), url.QueryEscape(endDataFormat))
	}

	msSecondeToSleep := 0

	for {

		select {
		case <-ctx.Done():
			responseChannel <- -1
			return -1
		default:

			if msSecondeToSleep > 0 {
				time.Sleep(time.Millisecond)
				msSecondeToSleep--
			} else {
				var nbCarFound int

				if searchingType == searchingFlex {
					var flexAvailable flexCarResponse

					apiCall(urlCalled, &flexAvailable)

					nbCarFound = flexAvailable.TotalNbVehicles
				} else if searchingType == searchingStation {
					var stationsAvailable stationsResponse

					apiCall(urlCalled, &stationsAvailable)

					nbCarFound = len(stationsAvailable.Stations)
				}

				if nbCarFound > 0 {
					responseChannel <- nbCarFound
					return nbCarFound
				}

				msSecondeToSleep = fetchDelayInMin * 60 * 1000 // Wait only 1ms each time to don't block the for loop and be able to catch the cancel signal
			}
		}
	}
}

// Make an api call at url passed and return the result in response object
func apiCall(url string, response interface{}) {
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
