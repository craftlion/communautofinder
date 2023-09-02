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

const cityId = 59 // see available cities -> https://restapifrontoffice.reservauto.net/ReservautoFrontOffice/index.html?urls.primaryName=Branch%20version%202%20(6.93.1)#/

const fetchDelayInMin = 1 // delay between two API call

const dateFormat = "2006-01-02T15:04:05"

// Different type of vehicule possible to search
const (
	searchingFlex = iota
	searchingStation
)

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

func searchCar(searchingType int, currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time, responseChannel chan int, ctx context.Context) int {
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

					communautoAPICall(urlCalled, &flexAvailable)

					nbCarFound = flexAvailable.TotalNbVehicles
				} else if searchingType == searchingStation {
					var stationsAvailable stationsResponse

					communautoAPICall(urlCalled, &stationsAvailable)

					nbCarFound = len(stationsAvailable.Stations)
				}

				if nbCarFound > 0 {
					responseChannel <- nbCarFound
					return nbCarFound
				}

				msSecondeToSleep = fetchDelayInMin * 60 * 1000
			}
		}
	}
}

func SearchStationCar(currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time) int {
	responseChannel := make(chan int)
	ctx, _ := context.WithCancel(context.Background())
	return searchCar(searchingStation, currentCoordinate, marginInKm, startDate, endDate, responseChannel, ctx)
}

func SearchFlexCar(currentCoordinate Coordinate, marginInKm float64) int {
	responseChannel := make(chan int)
	ctx, _ := context.WithCancel(context.Background())
	return searchCar(searchingFlex, currentCoordinate, marginInKm, time.Time{}, time.Time{}, responseChannel, ctx)
}

// This function is designed to be called as a goroutine. It returns the result in a channel and can be canceled using a cancel context.
func SearchStationCarForGoRoutine(currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time, responseChannel chan int, ctx context.Context) int {

	defer func() {
		if r := recover(); r != nil {
			responseChannel <- -1
			log.Printf("Pannic append : %s", r)
		}
	}()

	return searchCar(searchingStation, currentCoordinate, marginInKm, startDate, endDate, responseChannel, ctx)
}

// This function is designed to be called as a goroutine. It returns the result in a channel and can be canceled using a cancel context.
func SearchFlexCarForGoRoutine(currentCoordinate Coordinate, marginInKm float64, responseChannel chan int, ctx context.Context) int {

	defer func() {
		if r := recover(); r != nil {
			responseChannel <- -1
			log.Printf("Pannic append : %s", r)
		}
	}()

	return searchCar(searchingFlex, currentCoordinate, marginInKm, time.Time{}, time.Time{}, responseChannel, ctx)
}
