package communautofinder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

const cityId = 59 // see available cities -> https://restapifrontoffice.reservauto.net/ReservautoFrontOffice/index.html?urls.primaryName=Branch%20version%202%20(6.93.1)#/

const fetchDelayInMin = 1 // delay between two API call

var sugar *zap.SugaredLogger

var resultsChannelStation = make(chan int)
var resultsChannelFlex = make(chan int)

func communautoAPICall(url string, response interface{}) {
	resp, err := http.Get(url)

	if err != nil {
		sugar.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		errDecode := json.NewDecoder(resp.Body).Decode(response)

		if errDecode != nil {
			sugar.Fatal(errDecode)
		}
	} else {
		sugar.Fatalf("Error %d in API call", resp.StatusCode)
	}
}

// Keep looping until a station car is not found for the specified dates and position
// return the number of car found and fill the channel with
func SearchStationCar(currentCoordinate Coordinate, marginInKm float64, startDate time.Time, endDate time.Time) int {

	defer func() {
		if r := recover(); r != nil {
			resultsChannelStation <- 0
			sugar.Errorf("Pannic append :", r)
		}
	}()

	minCoordinate, maxCoordinate := currentCoordinate.ExpandCoordinate(marginInKm)

	startDateFormat := startDate.Format("2006-01-02T15:04:05")
	endDataFormat := endDate.Format("2006-01-02T15:04:05")

	urlCalled := fmt.Sprintf("https://restapifrontoffice.reservauto.net/api/v2/StationAvailability?CityId=%d&MaxLatitude=%f&MinLatitude=%f&MaxLongitude=%f&MinLongitude=%f&StartDate=%s&EndDate=%s", cityId, maxCoordinate.latitude, minCoordinate.latitude, maxCoordinate.longitude, minCoordinate.longitude, url.QueryEscape(startDateFormat), url.QueryEscape(endDataFormat))

	for {

		var stationsAvailable stationsResponse

		communautoAPICall(urlCalled, &stationsAvailable)

		nbCarFound := len(stationsAvailable.Stations)

		if nbCarFound > 0 {
			resultsChannelStation <- nbCarFound
			return nbCarFound
		}

		time.Sleep(fetchDelayInMin * time.Minute)
	}
}

// Keep looping until a flex car is not found for the specified position
// return the number of car found and fill the channel with
func SearchFlexCar(currentCoordinate Coordinate, marginInKm float64) int {

	defer func() {
		if r := recover(); r != nil {
			resultsChannelFlex <- 0
			sugar.Errorf("Pannic append :", r)
		}
	}()

	minCoordinate, maxCoordinate := currentCoordinate.ExpandCoordinate(marginInKm)

	urlCalled := fmt.Sprintf("https://restapifrontoffice.reservauto.net/api/v2/Vehicle/FreeFloatingAvailability?CityId=%d&MaxLatitude=%f&MinLatitude=%f&MaxLongitude=%f&MinLongitude=%f", cityId, maxCoordinate.latitude, minCoordinate.latitude, maxCoordinate.longitude, minCoordinate.longitude)

	for {
		var stationsAvailable flexCarResponse

		communautoAPICall(urlCalled, &stationsAvailable)

		nbCarFound := stationsAvailable.TotalNbVehicles

		if nbCarFound > 0 {
			resultsChannelFlex <- nbCarFound
			return nbCarFound
		}

		time.Sleep(fetchDelayInMin * time.Minute)
	}
}
