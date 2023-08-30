package communautofinder

import (
	"fmt"
	"testing"
	"time"
)

func TestNouvelleStructureExportee(t *testing.T) {

	var resultsChannelStation = make(chan int)
	var resultsChannelFlex = make(chan int)

	defer close(resultsChannelStation)
	defer close(resultsChannelFlex)

	var currentCoordinate Coordinate = New(45.538638, -73.570039)

	startDate := time.Now().AddDate(0, 0, 28)
	endDate := time.Now().AddDate(0, 0, 29)

	go RoutineSearchStationCar(currentCoordinate, 10, startDate, endDate, resultsChannelStation)
	go RoutineSearchFlexCar(currentCoordinate, 10, resultsChannelFlex)

	nbCarFoundStation := <-resultsChannelStation
	nbCarFoundFlex := <-resultsChannelFlex

	fmt.Printf("Station cars found : %d \n", nbCarFoundStation)
	fmt.Printf("Flex cars found : %d \n", nbCarFoundFlex)
}
