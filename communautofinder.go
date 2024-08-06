package communautofinder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const fetchDelayInMin = 1 // delay between two API call

const dateFormat = "2006-01-02T15:04:05" // time format accepted by communauto API

// As soon as at least one car is found return the number of cars found

func SearchStationCar(cityId CityId, currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time, vehiculeType VehiculeType) int {
	responseChannel := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	nbCarFound := searchCar(SearchingStation, cityId, currentCoordinate, marginInKm, startDate, endDate, vehiculeType, responseChannel, ctx, cancel)
	cancel()

	return nbCarFound
}

// As soon as at least one car is found return the number of cars found
func SearchFlexCar(cityId CityId, currentCoordinate Coordinate, marginInKm float64) int {
	responseChannel := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	nbCarFound := searchCar(SearchingFlex, cityId, currentCoordinate, marginInKm, time.Time{}, time.Time{}, AllTypes, responseChannel, ctx, cancel)
	cancel()

	return nbCarFound
}

// This function is designed to be called as a goroutine. As soon as at least one car is found return the number of cars found. Or can be cancelled by the context
func SearchStationCarForGoRoutine(cityId CityId, currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time, vehiculeType VehiculeType, responseChannel chan<- int, ctx context.Context, cancelCtxFunc context.CancelFunc) int {

	defer func() {
		if r := recover(); r != nil {
			responseChannel <- -1
			log.Printf("Pannic append : %s", r)
		}
	}()

	return searchCar(SearchingStation, cityId, currentCoordinate, marginInKm, startDate, endDate, vehiculeType, responseChannel, ctx, cancelCtxFunc)
}

// This function is designed to be called as a goroutine. As soon as at least one car is found return the number of cars found. Or can be cancelled by the context
func SearchFlexCarForGoRoutine(cityId CityId, currentCoordinate Coordinate, marginInKm float64, responseChannel chan<- int, ctx context.Context, cancelCtxFunc context.CancelFunc) int {

	defer func() {
		if r := recover(); r != nil {
			responseChannel <- -1
			log.Printf("Pannic append : %s", r)
		}
	}()

	return searchCar(SearchingFlex, cityId, currentCoordinate, marginInKm, time.Time{}, time.Time{}, AllTypes, responseChannel, ctx, cancelCtxFunc)
}

// Loop until a result is found. Return the number of cars found or can be cancelled by the context
func searchCar(searchingType SearchType, cityId CityId, currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time, vehiculeType VehiculeType, responseChannel chan<- int, ctx context.Context, cancelCtxFunc context.CancelFunc) int {
	minCoordinate, maxCoordinate := currentCoordinate.ExpandCoordinate(marginInKm)

	var urlCalled string

	if searchingType == SearchingFlex {
		urlCalled = fmt.Sprintf("https://restapifrontoffice.reservauto.net/api/v2/Vehicle/FreeFloatingAvailability?CityId=%d&MaxLatitude=%f&MinLatitude=%f&MaxLongitude=%f&MinLongitude=%f", cityId, maxCoordinate.latitude, minCoordinate.latitude, maxCoordinate.longitude, minCoordinate.longitude)
	} else if searchingType == SearchingStation {
		startDateFormat := startDate.Format(dateFormat)
		endDataFormat := endDate.Format(dateFormat)

		urlCalled = fmt.Sprintf("https://restapifrontoffice.reservauto.net/api/v2/StationAvailability?CityId=%d&MaxLatitude=%f&MinLatitude=%f&MaxLongitude=%f&MinLongitude=%f&StartDate=%s&EndDate=%s", cityId, maxCoordinate.latitude, minCoordinate.latitude, maxCoordinate.longitude, minCoordinate.longitude, url.QueryEscape(startDateFormat), url.QueryEscape(endDataFormat))

		if vehiculeType != AllTypes {
			urlCalled += fmt.Sprintf("&VehicleTypes=%d", vehiculeType)
		}
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
				nbCarFound := 0

				var err error

				if searchingType == SearchingFlex {
					var flexAvailable flexCarResponse

					err = apiCall(urlCalled, &flexAvailable)

					nbCarFound = flexAvailable.TotalNbVehicles
				} else if searchingType == SearchingStation {
					var stationsAvailable stationsResponse

					err = apiCall(urlCalled, &stationsAvailable)

					for _, station := range stationsAvailable.Stations {
						if station.SatisfiesFilters && station.RecommendedVehicleId != nil {
							nbCarFound++
						}
					}
				}

				if err != nil {
					cancelCtxFunc()
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
func apiCall(url string, response interface{}) error {
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

		errString := fmt.Sprintf("Error %d in API call", resp.StatusCode)
		err = errors.New(errString)

		log.Print(err)
	}

	return err
}
