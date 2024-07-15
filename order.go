package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type FindRiderStrategy interface {
	findRider(restaurantID int) Rider
}

type nextAvailableRider struct{}

type nearestAvailableRider struct{}

func placeOrder(c *gin.Context) {

	var newOrder Order
	if err := c.BindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newOrder.ID = len(orders) + 1

	var rider Rider

	/*
		Assign the nearest rider based on the algorithm
		configured in the system parameters

		nextAvailable -  Assign the next available rider
		nearestAvailable - Assign the nearest rider based on the nearest distance algorithm
	*/

	if nearestAlgoType == "nextAvailable" {
		rider = FindRiderStrategy(nextAvailableRider{}).findRider(newOrder.RestaurantID)
	} else {
		rider = FindRiderStrategy(nearestAvailableRider{}).findRider(newOrder.RestaurantID)
	}

	if rider.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No rider available for the order. Please try again later."})
		return
	}

	newOrder.RiderID = rider.ID
	newOrder.OrderStatus = "Placed"

	couponCode := c.Query("coupon_code")

	if couponCode != "" {

		for i, coupon := range coupons {
			if coupon.UserID == newOrder.UserID && coupon.Code == couponCode {
				newOrder.TotalAmount *= (1 - coupon.Discount/100)

				/*
					Remove the coupon from the list of available coupons
				*/

				coupons = append(coupons[:i], coupons[i+1:]...)
				newOrder.CouponID = coupon.ID
				break
			}
		}

	}
	/*
		Calculate estimated delivery time based on the distance between the restaurant
		and the user + the restaurant and the rider
		assuming delivery time is > food preparation time
	*/

	newOrder.ExpectedDeliveryTime = FindDeliveryTimeStrategy(riderFarFromRestaurant{}).calculateDeliveryTime(newOrder.RestaurantID, newOrder.UserID, newOrder.RiderID)

	orders = append(orders, newOrder)
	c.JSON(http.StatusOK, newOrder)
}

func getUserOrders(c *gin.Context) {

	id := cast.ToInt(c.Query("id"))
	isUserPresent := false
	var userOrders []Order
	for _, order := range orders {
		if order.UserID == id {
			userOrders = append(userOrders, order)
			isUserPresent = true
		}
	}

	if !isUserPresent {
		c.JSON(http.StatusNotFound, gin.H{"error": "No orders found for the given user id."})
		return
	}

	c.JSON(http.StatusOK, userOrders)
}

func getRiderOrders(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))
	isRiderPresent := false
	var riderOrders []Order
	for _, order := range orders {
		if order.RiderID == id {
			riderOrders = append(riderOrders, order)
			isRiderPresent = true
		}
	}

	if !isRiderPresent {
		c.JSON(http.StatusNotFound, gin.H{"error": "No orders found for the given rider id."})
		return
	}
	c.JSON(http.StatusOK, riderOrders)
}
