package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func getUserRatings(c *gin.Context) {

	userID := cast.ToInt(c.Query("id"))

	for _, user := range users {
		if user.ID == userID {
			c.JSON(http.StatusOK, user.AvgRating)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func getDriverRatings(c *gin.Context) {

	riderId := cast.ToInt(c.Query("id"))

	for _, rider := range riders {
		if rider.ID == riderId {
			c.JSON(http.StatusOK, rider.AvgRating)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
}

func submitRating(c *gin.Context) {

	var newRating Rating
	if err := c.ShouldBindJSON(&newRating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newRating.ID = len(ratings) + 1
	ratings = append(ratings, newRating)

	/*
	 Update user and rider ratings
	 as received in the newRating
	*/

	isUserPresent := false

	var userRating []Rating
	for _, rating := range ratings {
		if rating.UserID == newRating.UserID {
			userRating = append(userRating, rating)
			isUserPresent = true
		}
	}

	if !isUserPresent {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	newUser := getUserDetailsBasedOnId(newRating.UserID)
	newUser.AvgRating = calculateAverageRating(userRating)
	updateUserDetailsBasedOnId(newUser)

	isRiderPresent := false

	var riderRating []Rating
	for _, rating := range ratings {
		if rating.RiderID == newRating.RiderID {
			riderRating = append(riderRating, rating)
			isRiderPresent = true
		}
	}

	if !isRiderPresent {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rider not found"})
		return
	}

	newRider := getRiderDetailsBasedOnId(newRating.RiderID)
	newRider.AvgRating = calculateAverageRating(riderRating)
	updateRiderDetailsBasedOnId(newRider)

	c.JSON(http.StatusOK, newRating)
}

func calculateAverageRating(ratings []Rating) float64 {
	sum := 0.0
	for _, rating := range ratings {
		sum += rating.Rating
	}
	return sum / float64(len(ratings))
}
