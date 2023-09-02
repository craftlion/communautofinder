package communautofinder

import (
	"context"
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

	ctx, _ := context.WithCancel(context.Background())

	go SearchStationCarForGoRoutine(currentCoordinate, 10, startDate, endDate, resultsChannelStation, ctx)
	go SearchFlexCarForGoRoutine(currentCoordinate, 10, resultsChannelFlex, ctx)

	nbCarFoundStation := <-resultsChannelStation
	nbCarFoundFlex := <-resultsChannelFlex

	fmt.Printf("Station cars found : %d \n", nbCarFoundStation)
	fmt.Printf("Flex cars found : %d \n", nbCarFoundFlex)
}
