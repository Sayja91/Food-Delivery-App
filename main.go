package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var users []User
var riders []Rider
var menus []Menu
var orders []Order
var coupons []Coupon
var ratings []Rating
var restaurants []Restaurant

var defaultSpeed int
var nearestAlgoType string

func main() {

	r := gin.Default()

	/*
	 setting up global variables/configs
	*/

	defaultSpeed = 5
	nearestAlgoType = "nextAvailable"

	/*
	 Routes to register users, riders,
	 restaurants, and menus
	*/

	r.POST("/register/user", registerUser)
	r.POST("/register/rider", registerRider)
	r.POST("/register/restaurant", registerRestaurant)
	r.POST("/register/restaurant/menu", addMenu)

	/*
	 Routes providing customer services
	*/

	r.GET("/restaurants/suggest", suggestRestaurants)
	r.GET("/restaurant/menu", getMenu)
	r.POST("/order", placeOrder)
	r.PATCH("/rider/location", updateRiderLocation)
	r.GET("/user/orders", getUserOrders)
	r.GET("/rider/orders", getRiderOrders)

	/*
	 Bonus section handling for
	 coupons and ratings
	*/

	r.POST("/user/coupon", assignCouponToUser)
	r.POST("/submit-rating", submitRating)
	r.GET("/user/ratings", getUserRatings)
	r.GET("/rider/ratings", getDriverRatings)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not run the server: %v", err)
	}
}
