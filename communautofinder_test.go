package communautofinder

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestUseExemple(t *testing.T) {

	var currentCoordinate Coordinate = New(45.538638, -73.570039)
	startDate := time.Now().AddDate(0, 0, 28)
	endDate := time.Now().AddDate(0, 0, 29)

	// Search flex car
	nbCarFoundFlex := SearchFlexCar(currentCoordinate, 10)
	fmt.Printf("Flex cars found : %d \n", nbCarFoundFlex)

	// Search station car
	nbCarFoundStation := SearchStationCar(currentCoordinate, 10, startDate, endDate)
	fmt.Printf("Station cars found : %d \n", nbCarFoundStation)

	/////////////////////////////////

	var resultsChannelStation = make(chan int)
	var resultsChannelFlex = make(chan int)

	defer close(resultsChannelStation)
	defer close(resultsChannelFlex)

	ctx, cancel := context.WithCancel(context.Background())

	// Search flex car with go routine
	go SearchFlexCarForGoRoutine(currentCoordinate, 10, resultsChannelFlex, ctx)
	nbCarFoundFlex = <-resultsChannelFlex
	fmt.Printf("Flex cars found : %d \n", nbCarFoundFlex)

	// Search station car with go routine
	go SearchStationCarForGoRoutine(currentCoordinate, 10, startDate, endDate, resultsChannelStation, ctx)
	nbCarFoundStation = <-resultsChannelStation
	fmt.Printf("Station cars found : %d \n", nbCarFoundStation)

	cancel()
}
