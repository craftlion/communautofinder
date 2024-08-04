package communautofinder

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestUseExemple(t *testing.T) {

	const cityId = 59 // see available cities -> https://restapifrontoffice.reservauto.net/ReservautoFrontOffice/index.html?urls.primaryName=Branch%20version%202%20(6.93.1)#/

	var currentCoordinate Coordinate = New(45.538638, -73.570039)
	startDate := time.Now().AddDate(0, 0, 20)
	endDate := time.Now().AddDate(0, 0, 21)

	// Search flex car
	nbCarFoundFlex := SearchFlexCar(cityId, currentCoordinate, 10)
	fmt.Printf("Flex cars found : %d \n", nbCarFoundFlex)

	// Search station car
	nbCarFoundStation := SearchStationCar(cityId, currentCoordinate, 10, startDate, endDate, []VehiculeType{3, 5})
	fmt.Printf("Station cars found : %d \n", nbCarFoundStation)

	/////////////////////////////////

	var resultsChannelStation = make(chan int, 1)
	var resultsChannelFlex = make(chan int, 1)

	defer close(resultsChannelStation)
	defer close(resultsChannelFlex)

	ctx, cancel := context.WithCancel(context.Background())

	// Search flex car with go routine
	go SearchFlexCarForGoRoutine(cityId, currentCoordinate, 10, resultsChannelFlex, ctx, cancel)
	nbCarFoundFlex = <-resultsChannelFlex
	fmt.Printf("Flex cars found : %d \n", nbCarFoundFlex)

	// Search station car with go routine
	go SearchStationCarForGoRoutine(cityId, currentCoordinate, 10, startDate, endDate, resultsChannelStation, ctx, cancel)
	nbCarFoundStation = <-resultsChannelStation
	fmt.Printf("Station cars found : %d \n", nbCarFoundStation)

	cancel()
}
