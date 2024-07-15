package main

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FindDeliveryTimeStrategy interface {
	calculateDeliveryTime(restaurantId, userId, riderId int) int
}

type riderCloseToRestaurant struct{}

type riderFarFromRestaurant struct{}

func (r riderCloseToRestaurant) calculateDeliveryTime(restaurantId, userId, riderId int) int {

	/*
		if rider is close to restaurant, then the time taken to deliver the food will only
		include the time taken to deliver the food from restaurant to user + average preparing time of restaurant
	*/
	restaurant := getRestaurantDetailsBasedOnId(restaurantId)
	userDetails := getUserDetailsBasedOnId(userId)

	distanceBetweenRestaurntAndUser := calculateDistance(restaurant.Location.Lat, restaurant.Location.Lon, userDetails.Location.Lat, userDetails.Location.Lon)

	timeTakenToDeliverInHours := distanceBetweenRestaurntAndUser / defaultSpeed
	timeTakenToDeliverInMinutes := timeTakenToDeliverInHours * 60
	estimatedDeliveryTime := timeTakenToDeliverInMinutes + restaurant.AveragePreparingTime
	return estimatedDeliveryTime
}

func (r riderFarFromRestaurant) calculateDeliveryTime(restaurantId, userId, riderId int) int {

	/*
			if rider is far from restaurant during non peak hours, then the time taken to deliver the food will only
			include the time taken to deliver the food from restaurant to user  + rider distance from restaurant
		Assuming prep time<time for rider to reach restaurant
	*/

	restaurant := getRestaurantDetailsBasedOnId(restaurantId)
	userDetails := getUserDetailsBasedOnId(userId)
	riderDetails := getRiderDetailsBasedOnId(riderId)

	distanceBetweenRestaurntAndUser := calculateDistance(restaurant.Location.Lat, restaurant.Location.Lon, userDetails.Location.Lat, userDetails.Location.Lon)

	distanceBetweenRestaurantAndRider := calculateDistance(restaurant.Location.Lat, restaurant.Location.Lon, riderDetails.Location.Lat, riderDetails.Location.Lon)

	// fmt.Println("Distance between restaurant and rider", distanceBetweenRestaurantAndRider, "Distance between restaurant and user", distanceBetweenRestaurntAndUser)

	timeTakenToDeliverInHours := (distanceBetweenRestaurntAndUser + distanceBetweenRestaurantAndRider/defaultSpeed)

	// fmt.Println("timeTakenToDeliverInHours", timeTakenToDeliverInHours)

	timeTakenToDeliverInMinutes := timeTakenToDeliverInHours * 60

	// fmt.Println("timeTakenToDeliverInMinutes", timeTakenToDeliverInMinutes)

	estimatedDeliveryTime := timeTakenToDeliverInMinutes

	return estimatedDeliveryTime
}

func registerRider(c *gin.Context) {
	var newRider Rider
	if err := c.BindJSON(&newRider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newRider.ID = len(riders) + 1
	newRider.IsAvailable = true
	riders = append(riders, newRider)
	c.JSON(http.StatusOK, newRider)
}

func updateRiderLocation(c *gin.Context) {
	var updatedRider Rider
	if err := c.BindJSON(&updatedRider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, rider := range riders {
		if rider.ID == updatedRider.ID {
			riders[i].Location.Lat = updatedRider.Location.Lat
			riders[i].Location.Lon = updatedRider.Location.Lon
			c.JSON(http.StatusOK, riders[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Rider not found"})
}

func (sl nextAvailableRider) findRider(restaurantID int) Rider {

	for i, rider := range riders {
		if rider.IsAvailable {
			rider.IsAvailable = false
			riders[i] = rider
			return rider
		}
	}
	return Rider{}

}

func (sl nearestAvailableRider) findRider(restaurantID int) Rider {

	restaurant := getRestaurantDetailsBasedOnId(restaurantID)
	minTimeTaken := math.MaxInt64
	r := Rider{}
	index := 0
	foundRider := false

	for i, rider := range riders {

		if rider.IsAvailable {
			distance := calculateDistance(restaurant.Location.Lat, restaurant.Location.Lon, rider.Location.Lat, rider.Location.Lon)

			timeTakenToDeliver := distance / defaultSpeed

			if timeTakenToDeliver < minTimeTaken {
				minTimeTaken = timeTakenToDeliver
				r = rider
				index = i
				foundRider = true
			}
		}
	}

	if !foundRider {
		return Rider{}
	}

	r.IsAvailable = false
	riders[index] = r
	return r

}

func getRiderDetailsBasedOnId(id int) Rider {
	for _, rider := range riders {
		if rider.ID == id {
			return rider
		}
	}
	return Rider{}

}

func updateRiderDetailsBasedOnId(newRider Rider) {
	for i, rider := range riders {
		if rider.ID == newRider.ID {
			riders[i] = newRider
			break
		}
	}

}
