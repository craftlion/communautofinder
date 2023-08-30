package communautofinder

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
)

func testNouvelleStructureExportee(t *testing.T) {

	defer close(resultsChannelStation)
	defer close(resultsChannelFlex)

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar = logger.Sugar()

	var currentCoordinate Coordinate
	currentCoordinate.latitude = 45.538638
	currentCoordinate.longitude = -73.570039

	startDate := time.Now().AddDate(0, 0, 28)
	endDate := time.Now().AddDate(0, 0, 29)

	go SearchStationCar(currentCoordinate, 10, startDate, endDate)

	nbCarFound := <-resultsChannelStation

	fmt.Printf("Station cars found : %d \n", nbCarFound)

	go SearchFlexCar(currentCoordinate, 10)

	nbCarFound = <-resultsChannelFlex

	fmt.Printf("Flex cars found : %d \n", nbCarFound)
}
